package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"dashboard-cs-be/entities"
)

type mysqlImportRepository struct {
	db *sql.DB
}

// NewMySQLImportRepository constructs the MySQL-backed import repository.
func NewMySQLImportRepository(db *sql.DB) ImportRepository {
	return &mysqlImportRepository{db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// UpsertTicket
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlImportRepository) UpsertTicket(row *entities.ImportRow) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Upsert agent (INSERT IGNORE – hanya buat placeholder jika belum ada)
	if row.AgentID != "" {
		_, err = tx.Exec(
			`INSERT IGNORE INTO agents (id, name) VALUES (?, '')`,
			row.AgentID,
		)
		if err != nil {
			return "", fmt.Errorf("upsert agent: %w", err)
		}
	}

	// 2. Upsert customer berdasarkan phone number
	customerID, err := r.upsertCustomer(tx, row)
	if err != nil {
		return "", err
	}

	// 3. Cek apakah tiket sudah ada
	var existingStatus sql.NullString
	err = tx.QueryRow(
		`SELECT status FROM tickets WHERE id = ?`,
		row.TicketID,
	).Scan(&existingStatus)

	action := ""

	if err == sql.ErrNoRows {
		// ── INSERT baru ──────────────────────────────────────────────────
		if err = r.insertTicket(tx, row, customerID); err != nil {
			return "", err
		}
		action = "inserted"

	} else if err != nil {
		return "", fmt.Errorf("check ticket: %w", err)

	} else {
		// ── Tiket sudah ada ──────────────────────────────────────────────
		if existingStatus.String == "closed" {
			// Tiket closed tidak boleh dibuka ulang → skip
			_ = tx.Rollback()
			return "skipped", nil
		}
		// Status masih open → UPDATE
		if err = r.updateTicket(tx, row, customerID); err != nil {
			return "", err
		}
		action = "updated"
	}

	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("commit: %w", err)
	}
	return action, nil
}

// upsertCustomer melakukan INSERT … ON DUPLICATE KEY UPDATE berdasarkan phone.
// Mengembalikan customer_id yang dipakai.
func (r *mysqlImportRepository) upsertCustomer(tx *sql.Tx, row *entities.ImportRow) (string, error) {
	// Cek apakah sudah ada berdasarkan phone
	var existingID string
	err := tx.QueryRow(
		`SELECT id FROM customers WHERE phone = ?`,
		row.CustomerPhone,
	).Scan(&existingID)

	if err == nil {
		// Sudah ada → update nama & email jika berubah
		var emailArg interface{}
		if row.CustomerEmail != "" {
			emailArg = row.CustomerEmail
		}
		_, err = tx.Exec(
			`UPDATE customers SET name = ?, email = COALESCE(?, email), updated_at = NOW() WHERE id = ?`,
			row.CustomerName, emailArg, existingID,
		)
		if err != nil {
			return "", fmt.Errorf("update customer: %w", err)
		}
		return existingID, nil
	}

	if err != sql.ErrNoRows {
		return "", fmt.Errorf("check customer: %w", err)
	}

	// Belum ada → generate ID baru
	newID, err := r.nextCustomerID(tx)
	if err != nil {
		return "", err
	}

	var emailArg interface{}
	if row.CustomerEmail != "" {
		emailArg = row.CustomerEmail
	}

	_, err = tx.Exec(
		`INSERT INTO customers (id, name, phone, email) VALUES (?, ?, ?, ?)`,
		newID, row.CustomerName, row.CustomerPhone, emailArg,
	)
	if err != nil {
		return "", fmt.Errorf("insert customer: %w", err)
	}
	return newID, nil
}

// nextCustomerID mengambil ID numerik terbesar dari customers lalu +1.
func (r *mysqlImportRepository) nextCustomerID(tx *sql.Tx) (string, error) {
	var maxID sql.NullString
	_ = tx.QueryRow(`SELECT MAX(id) FROM customers`).Scan(&maxID)

	seq := 1
	if maxID.Valid && len(maxID.String) > 5 {
		// Format: CUST-00000001 → ambil bagian angka
		numPart := strings.TrimPrefix(maxID.String, "CUST-")
		var n int
		if _, err := fmt.Sscanf(numPart, "%d", &n); err == nil {
			seq = n + 1
		}
	}
	return fmt.Sprintf("CUST-%08d", seq), nil
}

func (r *mysqlImportRepository) insertTicket(tx *sql.Tx, row *entities.ImportRow, customerID string) error {
	var closedAt interface{}
	if row.ResolvedAt != "" {
		closedAt = row.ResolvedAt
	}

	var agentID interface{}
	if row.AgentID != "" {
		agentID = row.AgentID
	}

	_, err := tx.Exec(`
		INSERT INTO tickets
			(id, customer_id, channel, priority, status, topic, is_fcr,
			 agent_id, customer_type, created_at, closed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		row.TicketID, customerID, row.Channel, row.Priority, row.Status,
		row.Topic, row.IsFCR, agentID, row.CustomerType,
		row.CreatedAt, closedAt,
	)
	if err != nil {
		return fmt.Errorf("insert ticket %s: %w", row.TicketID, err)
	}
	return nil
}

func (r *mysqlImportRepository) updateTicket(tx *sql.Tx, row *entities.ImportRow, customerID string) error {
	var closedAt interface{}
	if row.ResolvedAt != "" {
		closedAt = row.ResolvedAt
	}

	var agentID interface{}
	if row.AgentID != "" {
		agentID = row.AgentID
	}

	_, err := tx.Exec(`
		UPDATE tickets SET
			customer_id   = ?,
			channel       = ?,
			priority      = ?,
			status        = ?,
			topic         = ?,
			is_fcr        = ?,
			agent_id      = COALESCE(?, agent_id),
			customer_type = ?,
			closed_at     = ?,
			updated_at    = NOW()
		WHERE id = ?`,
		customerID, row.Channel, row.Priority, row.Status,
		row.Topic, row.IsFCR, agentID, row.CustomerType,
		closedAt, row.TicketID,
	)
	if err != nil {
		return fmt.Errorf("update ticket %s: %w", row.TicketID, err)
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// SaveImportLog
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlImportRepository) SaveImportLog(result *entities.ImportResult) error {
	var errorsJSON interface{}
	if len(result.Errors) > 0 {
		b, err := json.Marshal(result.Errors)
		if err == nil {
			errorsJSON = string(b)
		}
	}

	_, err := r.db.Exec(`
		INSERT INTO import_logs
			(filename, total_rows, inserted, updated, skipped, error_count, errors)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		result.Filename, result.TotalRows,
		result.Inserted, result.Updated, result.Skipped,
		result.ErrorCount, errorsJSON,
	)
	if err != nil {
		return fmt.Errorf("SaveImportLog: %w", err)
	}
	return nil
}