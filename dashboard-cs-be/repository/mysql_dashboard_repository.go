package repository

import (
	"database/sql"
	"fmt"

	"dashboard-cs-be/entities"
)

type mysqlDashboardRepository struct {
	db *sql.DB
}

// NewMySQLDashboardRepository constructs the MySQL-backed repository.
func NewMySQLDashboardRepository(db *sql.DB) DashboardRepository {
	return &mysqlDashboardRepository{db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// GetSummary
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlDashboardRepository) GetSummary(from, to string) (*entities.SummaryRow, error) {
	query := `
		SELECT
			COALESCE(COUNT(t.id), 0)                                     AS total_tickets,
			COALESCE(SUM(t.status = 'open'), 0)                          AS open_count,
			COALESCE(SUM(t.status = 'closed'), 0)                        AS closed_count,
			COALESCE(ROUND(AVG(cr.score) / 5.0 * 100, 2), 0)            AS csat_percentage,
			COALESCE(ROUND(AVG(cr.score), 2), 0)                         AS csat_score
		FROM tickets t
		LEFT JOIN csat_responses cr ON cr.ticket_id = t.id
		WHERE DATE(t.created_at) BETWEEN ? AND ?
	`
	row := r.db.QueryRow(query, from, to)

	var s entities.SummaryRow
	if err := row.Scan(
		&s.TotalTickets,
		&s.Open,
		&s.Closed,
		&s.CSATPercentage,
		&s.CSATScore,
	); err != nil {
		return nil, fmt.Errorf("GetSummary: %w", err)
	}
	return &s, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// GetDailyTrend
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlDashboardRepository) GetDailyTrend(from, to string) ([]entities.TrendRow, error) {
	query := `
		SELECT
			DATE(t.created_at)               AS date,
			COUNT(t.id)                      AS created,
			SUM(t.status = 'closed')         AS solved
		FROM tickets t
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		GROUP BY DATE(t.created_at)
		ORDER BY DATE(t.created_at)
	`
	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetDailyTrend: %w", err)
	}
	defer rows.Close()

	var result []entities.TrendRow
	for rows.Next() {
		var tr entities.TrendRow
		if err := rows.Scan(&tr.Date, &tr.Created, &tr.Solved); err != nil {
			return nil, fmt.Errorf("GetDailyTrend scan: %w", err)
		}
		result = append(result, tr)
	}
	return result, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// GetTicketsPerHour
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlDashboardRepository) GetTicketsPerHour(from, to string) ([]entities.HourlyRow, error) {
	query := `
		SELECT
			DATE_FORMAT(t.created_at, '%H:00')  AS hour,
			COUNT(t.id)                          AS created,
			SUM(t.status = 'closed')             AS solved
		FROM tickets t
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		GROUP BY DATE_FORMAT(t.created_at, '%H:00')
		ORDER BY hour
	`
	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetTicketsPerHour: %w", err)
	}
	defer rows.Close()

	var result []entities.HourlyRow
	for rows.Next() {
		var hr entities.HourlyRow
		if err := rows.Scan(&hr.Hour, &hr.Created, &hr.Solved); err != nil {
			return nil, fmt.Errorf("GetTicketsPerHour scan: %w", err)
		}
		result = append(result, hr)
	}
	return result, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// GetPrioritySummary
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlDashboardRepository) GetPrioritySummary(from, to string) (*entities.PriorityRow, error) {
	query := `
		SELECT
			COALESCE(SUM(priority = 'roaming'), 0)     AS roaming,
			COALESCE(SUM(priority = 'extra_quota'), 0) AS extra_quota,
			COALESCE(SUM(priority = 'cc'), 0)          AS cc,
			COALESCE(SUM(priority = 'vip'), 0)         AS vip,
			COALESCE(SUM(priority = 'p1'), 0)          AS p1,
			COALESCE(SUM(priority = 'urgent'), 0)      AS urgent
		FROM tickets
		WHERE DATE(created_at) BETWEEN ? AND ?
	`
	row := r.db.QueryRow(query, from, to)

	var p entities.PriorityRow
	if err := row.Scan(
		&p.Roaming, &p.ExtraQuota, &p.CC, &p.VIP, &p.P1, &p.Urgent,
	); err != nil {
		return nil, fmt.Errorf("GetPrioritySummary: %w", err)
	}
	return &p, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// GetChannelSLA
// SLA % = (tickets resolved within sla_rules.max_duration_minutes) / total * 100
// Only closed tickets are eligible for SLA measurement.
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlDashboardRepository) GetChannelSLA(from, to string) ([]entities.ChannelSLARow, error) {
	query := `
		SELECT
			t.channel,
			ROUND(
				100.0 * SUM(
					t.status = 'closed'
					AND TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at) <= sr.max_duration_minutes
				) / NULLIF(COUNT(t.id), 0),
			2) AS sla_percentage,
			SUM(t.status = 'open')   AS open_count,
			SUM(t.status = 'closed') AS closed_count
		FROM tickets t
		JOIN sla_rules sr
			ON sr.priority = t.priority AND sr.is_active = TRUE
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		GROUP BY t.channel
		ORDER BY t.channel
	`
	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetChannelSLA: %w", err)
	}
	defer rows.Close()

	var result []entities.ChannelSLARow
	for rows.Next() {
		var ch entities.ChannelSLARow
		if err := rows.Scan(&ch.Channel, &ch.SLA, &ch.Open, &ch.Closed); err != nil {
			return nil, fmt.Errorf("GetChannelSLA scan: %w", err)
		}
		result = append(result, ch)
	}
	return result, rows.Err()
}

// ─────────────────────────────────────────────────────────────────────────────
// GetTopCorporate
// interactions = total tickets from that customer on that channel
// fcr_percentage = % of those tickets where is_fcr = true
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlDashboardRepository) GetTopCorporate(from, to string) ([]entities.TopCorporateRow, error) {
	query := `
		SELECT
			t.channel,
			c.name                                                          AS company_name,
			COUNT(t.id)                                                     AS interactions,
			SUM(t.status = 'open')                                          AS ticket_count,
			ROUND(100.0 * SUM(t.status = 'closed') / NULLIF(COUNT(t.id), 0), 2) AS fcr_percentage
		FROM tickets t
		JOIN customers c ON c.id = t.customer_id
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		GROUP BY t.channel, c.id, c.name
		ORDER BY t.channel, interactions DESC
	`
	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetTopCorporate: %w", err)
	}
	defer rows.Close()

	var result []entities.TopCorporateRow
	for rows.Next() {
		var row entities.TopCorporateRow
		if err := rows.Scan(
			&row.Channel, &row.CompanyName,
			&row.Interactions, &row.Tickets, &row.FCRPercentage,
		); err != nil {
			return nil, fmt.Errorf("GetTopCorporate scan: %w", err)
		}
		result = append(result, row)
	}
	return result, rows.Err()
}

func (r *mysqlDashboardRepository) GetTopKIP(from, to string) ([]entities.TopKIPRow, error) {
	query := `
		SELECT
			t.channel,
			t.topic,
			COUNT(t.id)                                                          AS interactions,
			SUM(t.status = 'open')                                               AS ticket_count,
			ROUND(100.0 * SUM(t.status = 'closed') / NULLIF(COUNT(t.id), 0), 2) AS fcr_percentage
		FROM tickets t
		WHERE DATE(t.created_at) BETWEEN ? AND ?
		GROUP BY t.channel, t.topic
		ORDER BY t.channel, interactions DESC
	`
	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetTopKIP: %w", err)
	}
	defer rows.Close()

	var result []entities.TopKIPRow
	for rows.Next() {
		var row entities.TopKIPRow
		if err := rows.Scan(
			&row.Channel, &row.Topic,
			&row.Interactions, &row.Tickets, &row.FCRPercentage,
		); err != nil {
			return nil, fmt.Errorf("GetTopKIP scan: %w", err)
		}
		result = append(result, row)
	}
	return result, rows.Err()
}


// ─────────────────────────────────────────────────────────────────────────────
// GetRealtime  (today only — no date param)
// delta = today's created count vs yesterday's created count
// ─────────────────────────────────────────────────────────────────────────────

func (r *mysqlDashboardRepository) GetRealtime() (*entities.RealtimeRow, error) {
	query := `
		SELECT
			-- SLA Today: hanya tiket yang dibuat HARI INI
			COALESCE(ROUND(
				100.0 * SUM(
					DATE(t.created_at) = CURDATE()
					AND t.status = 'closed'
					AND TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at) <= sr.max_duration_minutes
				) / NULLIF(SUM(DATE(t.created_at) = CURDATE()), 0),
			2), 0) AS sla_today_percentage,

			-- SLA Delta: SLA hari ini vs kemarin
			COALESCE(ROUND(
				100.0 * SUM(
					DATE(t.created_at) = CURDATE()
					AND t.status = 'closed'
					AND TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at) <= sr.max_duration_minutes
				) / NULLIF(SUM(DATE(t.created_at) = CURDATE()), 0)
				-
				100.0 * SUM(
					DATE(t.created_at) = DATE_SUB(CURDATE(), INTERVAL 1 DAY)
					AND t.status = 'closed'
					AND TIMESTAMPDIFF(MINUTE, t.created_at, t.closed_at) <= sr.max_duration_minutes
				) / NULLIF(SUM(DATE(t.created_at) = DATE_SUB(CURDATE(), INTERVAL 1 DAY)), 0),
			2), 0) AS sla_today_delta,

			-- Created Today: hanya tiket dibuat hari ini
			COALESCE(SUM(DATE(t.created_at) = CURDATE()), 0) AS created_today_total,

			-- Open Tickets: SEMUA tiket open, tanpa filter tanggal
			COALESCE(SUM(t.status = 'open'), 0) AS open_tickets,

			-- Unassigned: SEMUA tiket open tanpa agent, tanpa filter tanggal
			COALESCE(SUM(t.status = 'open' AND t.agent_id IS NULL), 0) AS unassigned

		FROM tickets t
		JOIN sla_rules sr ON sr.priority = t.priority AND sr.is_active = TRUE
	`
	// Tidak ada WHERE — semua data dimuat, filter dilakukan di level SELECT
	row := r.db.QueryRow(query)

	var rt entities.RealtimeRow
	if err := row.Scan(
		&rt.SLATodayPercentage,
		&rt.SLATodayDelta,
		&rt.CreatedTodayTotal,
		&rt.OpenTickets,
		&rt.Unassigned,
	); err != nil {
		return nil, fmt.Errorf("GetRealtime: %w", err)
	}
	return &rt, nil
}
