package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/repository/interfaces"
)

// ErrIncidentNotFound dikembalikan saat ID incident tidak ditemukan.
var ErrIncidentNotFound = errors.New("incident tidak ditemukan")

// ErrIncidentAlreadyResolved dikembalikan saat incident sudah resolved.
var ErrIncidentAlreadyResolved = errors.New("incident sudah diselesaikan sebelumnya")

type mysqlIncidentRepository struct {
	db *sql.DB
}

// NewMySQLIncidentRepository constructs the MySQL-backed incident repository.
func NewMySQLIncidentRepository(db *sql.DB) interfaces.IncidentRepository {
	return &mysqlIncidentRepository{db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// Create
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlIncidentRepository) Create(inc *entities.Incident) (*entities.Incident, error) {
	query := `
		INSERT INTO incidents (id, title, description, severity, status, started_at, created_by)
		VALUES (?, ?, ?, ?, 'active', ?, ?)`

	_, err := r.db.Exec(query,
		inc.ID, inc.Title, inc.Description,
		inc.Severity, inc.StartedAt.Format("2006-01-02 15:04:05"),
		inc.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("Create incident: %w", err)
	}

	return r.findByID(inc.ID)
}

// ─────────────────────────────────────────────────────────────────────────────
// GetActive
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlIncidentRepository) GetActive() ([]entities.Incident, error) {
	return r.query(`
		SELECT id, title, description, severity, status,
		       started_at, resolved_at, created_by, created_at, updated_at
		FROM incidents
		WHERE status = 'active'
		ORDER BY started_at DESC`)
}

// ─────────────────────────────────────────────────────────────────────────────
// GetHistory
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlIncidentRepository) GetHistory(from, to string) ([]entities.Incident, error) {
	return r.query(`
		SELECT id, title, description, severity, status,
		       started_at, resolved_at, created_by, created_at, updated_at
		FROM incidents
		WHERE DATE(started_at) BETWEEN ? AND ?
		ORDER BY started_at DESC`, from, to)
}

// ─────────────────────────────────────────────────────────────────────────────
// Resolve
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlIncidentRepository) Resolve(id string) (*entities.Incident, error) {
	// Cek keberadaan dan status saat ini
	var status string
	err := r.db.QueryRow(`SELECT status FROM incidents WHERE id = ?`, id).Scan(&status)
	if err == sql.ErrNoRows {
		return nil, ErrIncidentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("Resolve check: %w", err)
	}
	if status == "resolved" {
		return nil, ErrIncidentAlreadyResolved
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	_, err = r.db.Exec(`
		UPDATE incidents
		SET status = 'resolved', resolved_at = ?, updated_at = NOW()
		WHERE id = ?`, now, id)
	if err != nil {
		return nil, fmt.Errorf("Resolve update: %w", err)
	}

	return r.findByID(id)
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlIncidentRepository) findByID(id string) (*entities.Incident, error) {
	results, err := r.query(`
		SELECT id, title, description, severity, status,
		       started_at, resolved_at, created_by, created_at, updated_at
		FROM incidents WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrIncidentNotFound
	}
	return &results[0], nil
}

func (r *mysqlIncidentRepository) query(q string, args ...interface{}) ([]entities.Incident, error) {
	// Ganti placeholder ? sesuai jumlah args (sudah sesuai untuk MySQL)
	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("incident query: %w", err)
	}
	defer rows.Close()

	var result []entities.Incident
	for rows.Next() {
		var inc entities.Incident
		var desc sql.NullString
		var resolvedAt sql.NullTime

		if err := rows.Scan(
			&inc.ID, &inc.Title, &desc,
			&inc.Severity, &inc.Status,
			&inc.StartedAt, &resolvedAt,
			&inc.CreatedBy, &inc.CreatedAt, &inc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("incident scan: %w", err)
		}

		if desc.Valid {
			inc.Description = desc.String
		}
		if resolvedAt.Valid {
			t := resolvedAt.Time
			inc.ResolvedAt = &t
		}

		result = append(result, inc)
	}
	return result, rows.Err()
}

// nextIncidentID menghasilkan ID format INC-YYYYMMDD-XXX.
// Dipanggil dari usecase, bukan langsung dari repository.
func NextIncidentID(db *sql.DB) (string, error) {
	today := time.Now().Format("20060102")
	prefix := "INC-" + today + "-"

	var maxID sql.NullString
	_ = db.QueryRow(
		`SELECT MAX(id) FROM incidents WHERE id LIKE ?`,
		prefix+"%",
	).Scan(&maxID)

	seq := 1
	if maxID.Valid && len(maxID.String) > len(prefix) {
		suffix := strings.TrimPrefix(maxID.String, prefix)
		var n int
		if _, err := fmt.Sscanf(suffix, "%d", &n); err == nil {
			seq = n + 1
		}
	}
	return fmt.Sprintf("%s%03d", prefix, seq), nil
}
