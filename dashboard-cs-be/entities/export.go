package entities

// ─────────────────────────────────────────────────────────────────────────────
// Export – query result shapes (repository layer)
// ─────────────────────────────────────────────────────────────────────────────

// ExportFilter adalah parameter filter yang diterima dari query string.
type ExportFilter struct {
	From     string   // "2006-01-02"
	To       string   // "2006-01-02"
	Channels []string // kosong = all channels
}

// ExportChannelRow — data per channel untuk Sheet 2.
type ExportChannelRow struct {
	Channel      string  `db:"channel"`
	TotalTickets int     `db:"total_tickets"`
	Open         int     `db:"open_count"`
	Closed       int     `db:"closed_count"`
	SLAPercent   float64 `db:"sla_percentage"`
	FCRPercent   float64 `db:"fcr_percentage"`
}

// ExportCustomerRow — data per customer untuk Sheet 3.
type ExportCustomerRow struct {
	CustomerID   string  `db:"customer_id"`
	Name         string  `db:"name"`
	Phone        string  `db:"phone"`
	CustomerType string  `db:"customer_type"`
	Total        int     `db:"total_tickets"`
	Open         int     `db:"open_count"`
	Closed       int     `db:"closed_count"`
	SLAAchieved  int     `db:"sla_achieved"`
	SLABreached  int     `db:"sla_breached"`
}

// ExportTopicRow — data per topic + channel untuk Sheet 4.
type ExportTopicRow struct {
	Topic      string  `db:"topic"`
	Channel    string  `db:"channel"`
	Total      int     `db:"total_tickets"`
	Open       int     `db:"open_count"`
	Closed     int     `db:"closed_count"`
	FCRPercent float64 `db:"fcr_percentage"`
}

// ExportPriorityRow — data per priority untuk Sheet 5.
type ExportPriorityRow struct {
	Priority       string  `db:"priority"`
	Total          int     `db:"total_tickets"`
	Open           int     `db:"open_count"`
	Closed         int     `db:"closed_count"`
	SLAAchieved    int     `db:"sla_achieved"`
	SLABreached    int     `db:"sla_breached"`
	AvgResolutionM float64 `db:"avg_resolution_minutes"` // dalam menit, -1 jika tidak ada data
}

// ─────────────────────────────────────────────────────────────────────────────
// Export – aggregated payload (usecase builds, handler streams)
// ─────────────────────────────────────────────────────────────────────────────

// ExportPayload adalah keseluruhan data yang dibutuhkan untuk generate Excel.
type ExportPayload struct {
	Filter    ExportFilter
	Summary   SummaryRow          // reuse dari dashboard.go
	Channels  []ExportChannelRow
	Customers []ExportCustomerRow
	Topics    []ExportTopicRow
	Priorities []ExportPriorityRow
}