package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/repository/interfaces"
)

type mysqlExportRepository struct {
	db *sql.DB
}

// NewMySQLExportRepository constructs the MySQL-backed export repository.
func NewMySQLExportRepository(db *sql.DB) interfaces.ExportRepository {
	return &mysqlExportRepository{db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper: build channel WHERE clause
// Mengembalikan SQL snippet + args tambahan untuk channel filter.
// Jika Channels kosong → tidak ada filter channel (all).
// ─────────────────────────────────────────────────────────────────────────────

func channelClause(f entities.ExportFilter, tableAlias string) (string, []interface{}) {
	if len(f.Channels) == 0 {
		return "", nil
	}
	col := "channel"
	if tableAlias != "" {
		col = tableAlias + ".channel"
	}
	placeholders := make([]string, len(f.Channels))
	args := make([]interface{}, len(f.Channels))
	for i, ch := range f.Channels {
		placeholders[i] = "?"
		args[i] = ch
	}
	return fmt.Sprintf(" AND %s IN (%s)", col, strings.Join(placeholders, ",")), args
}

// buildArgs menggabungkan from, to, lalu args channel.
func buildArgs(from, to string, extra []interface{}) []interface{} {
	base := []interface{}{from, to}
	return append(base, extra...)
}

// ─────────────────────────────────────────────────────────────────────────────
// GetExportSummary
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlExportRepository) GetExportSummary(f entities.ExportFilter) (*entities.SummaryRow, error) {
	chClause, chArgs := channelClause(f, "t")
	query := fmt.Sprintf(`
		SELECT
			COALESCE(COUNT(t.id), 0)                                  AS total_tickets,
			COALESCE(SUM(t.status = 'open'), 0)                       AS open_count,
			COALESCE(SUM(t.status = 'closed'), 0)                     AS closed_count,
			COALESCE(ROUND(AVG(cr.score) / 5.0 * 100, 2), 0)         AS csat_percentage,
			COALESCE(ROUND(AVG(cr.score), 2), 0)                      AS csat_score
		FROM tickets t
		LEFT JOIN csat_responses cr ON cr.ticket_id = t.id
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		%s`, chClause)

	row := r.db.QueryRow(query, buildArgs(f.From, f.To, chArgs)...)
	var s entities.SummaryRow
	if err := row.Scan(&s.TotalTickets, &s.Open, &s.Closed, &s.CSATPercentage, &s.CSATScore); err != nil {
		return nil, fmt.Errorf("GetExportSummary: %w", err)
	}
	return &s, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// GetExportChannels — Sheet 2
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlExportRepository) GetExportChannels(f entities.ExportFilter) ([]entities.ExportChannelRow, error) {
	chClause, chArgs := channelClause(f, "t")
	query := fmt.Sprintf(`
		SELECT
			t.channel,
			COUNT(t.id)                                                        AS total_tickets,
			COALESCE(SUM(t.status = 'open'), 0)                               AS open_count,
			COALESCE(SUM(t.status = 'closed'), 0)                             AS closed_count,
			COALESCE(ROUND(
				100.0 * SUM(
					t.status = 'closed'
					AND TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at) <= sr.max_duration_minutes
				) / NULLIF(COUNT(t.id), 0), 2), 0)                            AS sla_percentage,
			COALESCE(ROUND(
				100.0 * SUM(t.status = 'closed') / NULLIF(COUNT(t.id), 0)
			, 2), 0)                                                           AS fcr_percentage
		FROM tickets t
		JOIN sla_rules sr ON sr.priority = t.priority AND sr.is_active = TRUE
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		%s
		GROUP BY t.channel
		ORDER BY total_tickets DESC`, chClause)

	rows, err := r.db.Query(query, buildArgs(f.From, f.To, chArgs)...)
	if err != nil {
		return nil, fmt.Errorf("GetExportChannels: %w", err)
	}
	defer rows.Close()

	var result []entities.ExportChannelRow
	for rows.Next() {
		var row entities.ExportChannelRow
		if err := rows.Scan(&row.Channel, &row.TotalTickets, &row.Open, &row.Closed,
			&row.SLAPercent, &row.FCRPercent); err != nil {
			return nil, fmt.Errorf("GetExportChannels scan: %w", err)
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func (r *mysqlExportRepository) GetExportTopics(f entities.ExportFilter) ([]entities.ExportTopicRow, error) {
	chClause, chArgs := channelClause(f, "t")
	query := fmt.Sprintf(`
		SELECT
			t.topic,
			t.channel,
			COUNT(t.id)                                                            AS total_tickets,
			COALESCE(SUM(t.status = 'open'), 0)                                   AS open_count,
			COALESCE(SUM(t.status = 'closed'), 0)                                 AS closed_count,
			COALESCE(ROUND(
				100.0 * SUM(t.status = 'closed') / NULLIF(COUNT(t.id), 0)
			, 2), 0)                                                               AS fcr_percentage
		FROM tickets t
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		%s
		GROUP BY t.topic, t.channel
		ORDER BY total_tickets DESC`, chClause)

	rows, err := r.db.Query(query, buildArgs(f.From, f.To, chArgs)...)
	if err != nil {
		return nil, fmt.Errorf("GetExportTopics: %w", err)
	}
	defer rows.Close()

	var result []entities.ExportTopicRow
	for rows.Next() {
		var row entities.ExportTopicRow
		if err := rows.Scan(&row.Topic, &row.Channel, &row.Total, &row.Open,
			&row.Closed, &row.FCRPercent); err != nil {
			return nil, fmt.Errorf("GetExportTopics scan: %w", err)
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// GetExportCustomers — Sheet 3
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlExportRepository) GetExportCustomers(f entities.ExportFilter) ([]entities.ExportCustomerRow, error) {
	chClause, chArgs := channelClause(f, "t")

	// customer_type: gunakan kolom baru jika ada, fallback ke 'unknown'
	query := fmt.Sprintf(`
		SELECT
			c.id                                                               AS customer_id,
			c.name,
			c.phone,
			COALESCE(t.customer_type, 'unknown')                               AS customer_type,
			COUNT(t.id)                                                        AS total_tickets,
			COALESCE(SUM(t.status = 'open'), 0)                               AS open_count,
			COALESCE(SUM(t.status = 'closed'), 0)                             AS closed_count,
			COALESCE(SUM(
				t.status = 'closed'
				AND TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at) <= sr.max_duration_minutes
			), 0)                                                              AS sla_achieved,
			COALESCE(SUM(
				t.status = 'closed'
				AND TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at) > sr.max_duration_minutes
			), 0)                                                              AS sla_breached
		FROM tickets t
		JOIN customers c  ON c.id = t.customer_id
		JOIN sla_rules sr ON sr.priority = t.priority AND sr.is_active = TRUE
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		%s
		GROUP BY c.id, c.name, c.phone, t.customer_type
		ORDER BY total_tickets DESC`, chClause)

	rows, err := r.db.Query(query, buildArgs(f.From, f.To, chArgs)...)
	if err != nil {
		return nil, fmt.Errorf("GetExportCustomers: %w", err)
	}
	defer rows.Close()

	var result []entities.ExportCustomerRow
	for rows.Next() {
		var row entities.ExportCustomerRow
		if err := rows.Scan(&row.CustomerID, &row.Name, &row.Phone, &row.CustomerType,
			&row.Total, &row.Open, &row.Closed, &row.SLAAchieved, &row.SLABreached); err != nil {
			return nil, fmt.Errorf("GetExportCustomers scan: %w", err)
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// GetExportPriorities — Sheet 5
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlExportRepository) GetExportPriorities(f entities.ExportFilter) ([]entities.ExportPriorityRow, error) {
	chClause, chArgs := channelClause(f, "t")
	query := fmt.Sprintf(`
		SELECT
			t.priority,
			COUNT(t.id)                                                        AS total_tickets,
			COALESCE(SUM(t.status = 'open'), 0)                               AS open_count,
			COALESCE(SUM(t.status = 'closed'), 0)                             AS closed_count,
			COALESCE(SUM(
				t.status = 'closed'
				AND TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at) <= sr.max_duration_minutes
			), 0)                                                              AS sla_achieved,
			COALESCE(SUM(
				t.status = 'closed'
				AND TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at) > sr.max_duration_minutes
			), 0)                                                              AS sla_breached,
			COALESCE(
				AVG(CASE WHEN t.status = 'closed'
				    THEN TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at)
				    END),
			-1)                                                                AS avg_resolution_minutes
		FROM tickets t
		JOIN sla_rules sr ON sr.priority = t.priority AND sr.is_active = TRUE
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		%s
		GROUP BY t.priority
		ORDER BY FIELD(t.priority,'p1','urgent','vip','cc','roaming','extra_quota','normal')`,
		chClause)

	rows, err := r.db.Query(query, buildArgs(f.From, f.To, chArgs)...)
	if err != nil {
		return nil, fmt.Errorf("GetExportPriorities: %w", err)
	}
	defer rows.Close()

	var result []entities.ExportPriorityRow
	for rows.Next() {
		var row entities.ExportPriorityRow
		if err := rows.Scan(&row.Priority, &row.Total, &row.Open, &row.Closed,
			&row.SLAAchieved, &row.SLABreached, &row.AvgResolutionM); err != nil {
			return nil, fmt.Errorf("GetExportPriorities scan: %w", err)
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
