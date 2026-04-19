package entities

// ─────────────────────────────────────────────────────────────────────────────
// Domain / DB row structs  (repository layer uses these)
// ─────────────────────────────────────────────────────────────────────────────

type Customer struct {
	ID        string  `db:"id"`
	Name      string  `db:"name"`
	Phone     string  `db:"phone"`
	Email     *string `db:"email"`
}

type Ticket struct {
	ID         string  `db:"id"`
	CustomerID string  `db:"customer_id"`
	Channel    string  `db:"channel"`
	Priority   string  `db:"priority"`
	Status     string  `db:"status"`
	Topic      string  `db:"topic"`
	IsFCR      bool    `db:"is_fcr"`
}

type SLARule struct {
	ID                  int    `db:"id"`
	Priority            string `db:"priority"`
	MaxDurationMinutes  int    `db:"max_duration_minutes"`
	IsActive            bool   `db:"is_active"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Query result shapes (returned by repository SQL queries)
// ─────────────────────────────────────────────────────────────────────────────

type SummaryRow struct {
	TotalTickets   int     `db:"total_tickets"`
	Open           int     `db:"open_count"`
	Closed         int     `db:"closed_count"`
	CSATPercentage float64 `db:"csat_percentage"`
	CSATScore      float64 `db:"csat_score"`
}

type TrendRow struct {
	Date    string `db:"date"`
	Created int    `db:"created"`
	Solved  int    `db:"solved"`
}

type HourlyRow struct {
	Hour    string `db:"hour"`
	Created int    `db:"created"`
	Solved  int    `db:"solved"`
}

type PriorityRow struct {
	Roaming    int `db:"roaming"`
	ExtraQuota int `db:"extra_quota"`
	CC         int `db:"cc"`
	VIP        int `db:"vip"`
	P1         int `db:"p1"`
	Urgent     int `db:"urgent"`
}

type ChannelSLARow struct {
	Channel string  `db:"channel"`
	SLA     float64 `db:"sla_percentage"`
	Open    int     `db:"open_count"`
	Closed  int     `db:"closed_count"`
}

type TopCorporateRow struct {
	Channel       string  `db:"channel"`
	CompanyName   string  `db:"company_name"`
	Interactions  int     `db:"interactions"`
	Tickets       int     `db:"ticket_count"`
	FCRPercentage float64 `db:"fcr_percentage"`
}

type TopKIPRow struct {
	Channel       string  `db:"channel"`
	Topic         string  `db:"topic"`
	Interactions  int     `db:"interactions"`
	Tickets       int     `db:"ticket_count"`
	FCRPercentage float64 `db:"fcr_percentage"`
}

type RealtimeRow struct {
	SLATodayPercentage float64 `db:"sla_today_percentage"`
	SLATodayDelta      float64 `db:"sla_today_delta"`
	CreatedTodayTotal  int     `db:"created_today_total"`
	OpenTickets        int     `db:"open_tickets"`
	Unassigned         int     `db:"unassigned"`
}

// ─────────────────────────────────────────────────────────────────────────────
// API Response structs  (usecase layer builds these, handler serialises them)
// ─────────────────────────────────────────────────────────────────────────────

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// --- GET /api/v1/realtime ---

type RealtimeResponse struct {
	SLAToday        SLAToday `json:"sla_today"`
	CreatedToday    Created  `json:"created_today"`
	OpenTickets     int      `json:"open_tickets"`
	Unassigned      int      `json:"unassigned"`
	IncidentsActive int      `json:"incidents_active"` // future: query from incidents table
}

type SLAToday struct {
	Percentage float64 `json:"percentage"`
	Delta      float64 `json:"delta"`
}

type Created struct {
	Total int `json:"total"`
	Delta int `json:"delta"`
}

// --- GET /api/v1/dashboard ---

type DashboardResponse struct {
	Filter          FilterMeta       `json:"filter"`
	Summary         SummaryAggregate `json:"summary"`
	DailyTrend      []DailyTrendPoint `json:"daily_trend"`
	TicketsPerHour  []TicketHourPoint `json:"tickets_per_hour"`
	PrioritySummary PrioritySummary  `json:"priority_summary"`
	Channels        []ChannelDetail  `json:"channels"`
}

type FilterMeta struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type SummaryAggregate struct {
	TotalTickets   int     `json:"total_tickets"`
	Open           int     `json:"open"`
	Closed         int     `json:"closed"`
	CSATPercentage float64 `json:"csat_percentage"`
	CSATScore      float64 `json:"csat_score"`
}

type DailyTrendPoint struct {
	Date    string `json:"date"`
	Created int    `json:"created"`
	Solved  int    `json:"solved"`
}

type TicketHourPoint struct {
	Hour    string `json:"hour"`
	Created int    `json:"created"`
	Solved  int    `json:"solved"`
}

type PrioritySummary struct {
	Roaming    int `json:"roaming"`
	ExtraQuota int `json:"extra_quota"`
	CC         int `json:"cc"`
	VIP        int `json:"vip"`
	P1         int `json:"p1"`
	Urgent     int `json:"urgent"`
}

type ChannelDetail struct {
	Channel      string         `json:"channel"`
	SLA          float64        `json:"sla"`
	Open         int            `json:"open"`
	Closed       int            `json:"closed"`
	TopCorporate []TopCorporate `json:"top_corporate"`
	TopKIP       []TopKIP       `json:"top_kip"`
}

type TopCorporate struct {
	CompanyName   string  `json:"company_name"`
	Interactions  int     `json:"interactions"`
	Tickets       int     `json:"tickets"`
	FCRPercentage float64 `json:"fcr_percentage"`
}

type TopKIP struct {
	Topic         string  `json:"topic"`
	Interactions  int     `json:"interactions"`
	Tickets       int     `json:"tickets"`
	FCRPercentage float64 `json:"fcr_percentage"`
}
