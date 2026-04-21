package entities

// ─────────────────────────────────────────────────────────────────────────────
// Domain / DB row structs  (repository layer uses these)
// ─────────────────────────────────────────────────────────────────────────────

type Customer struct {
	ID    string  `db:"id"`
	Name  string  `db:"name"`
	Phone string  `db:"phone"`
	Email *string `db:"email"`
}

type Ticket struct {
	ID         string `db:"id"`
	CustomerID string `db:"customer_id"`
	Channel    string `db:"channel"`
	Priority   string `db:"priority"`
	Status     string `db:"status"`
	Topic      string `db:"topic"`
}

type SLARule struct {
	ID                 int    `db:"id"`
	Priority           string `db:"priority"`
	MaxDurationMinutes int    `db:"max_duration_minutes"`
	IsActive           bool   `db:"is_active"`
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
	Unassigned     int     `db:"unassigned_count"`
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
	CompanyName   string  `db:"company_name"`
	Interactions  int     `db:"interactions"`
	Tickets       int     `db:"ticket_count"`
	FCRPercentage float64 `db:"fcr_percentage"`
	Total         int     `db:"total_count"` // untuk pagination
}

type TopKIPRow struct {
	Topic         string  `db:"topic"`
	Interactions  int     `db:"interactions"`
	Tickets       int     `db:"ticket_count"`
	FCRPercentage float64 `db:"fcr_percentage"`
	Total         int     `db:"total_count"` // untuk pagination
}

type RealtimeRow struct {
	SLATodayPercentage float64 `db:"sla_today_percentage"`
	SLATodayDelta      float64 `db:"sla_today_delta"`
	CreatedTodayTotal  int     `db:"created_today_total"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Filter structs
// ─────────────────────────────────────────────────────────────────────────────

// ChannelDetailFilter adalah parameter untuk GET /api/v1/dashboard/channels
type ChannelDetailFilter struct {
	From    string
	To      string
	Channel string // wajib: email | whatsapp | social_media | live_chat | call_center
	Page    int    // default: 1
	Limit   int    // default: 10
}

// ─────────────────────────────────────────────────────────────────────────────
// API Response structs
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
	IncidentsActive int      `json:"incidents_active"`
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
	Filter          FilterMeta        `json:"filter"`
	Summary         SummaryAggregate  `json:"summary"`
	DailyTrend      []DailyTrendPoint `json:"daily_trend"`
	TicketsPerHour  []TicketHourPoint `json:"tickets_per_hour"`
	PrioritySummary PrioritySummary   `json:"priority_summary"`
	Channels        []ChannelSummary  `json:"channels"` // hanya SLA, open, closed — tanpa nested detail
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
	Unassigned     int     `json:"unassigned"`
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

// ChannelSummary — ringkasan per channel di /dashboard (tanpa top_corporate & top_kip)
type ChannelSummary struct {
	Channel string  `json:"channel"`
	SLA     float64 `json:"sla"`
	Open    int     `json:"open"`
	Closed  int     `json:"closed"`
}

// --- GET /api/v1/dashboard/channels ---

// ChannelDetailResponse — response lengkap per channel dengan pagination
type ChannelDetailResponse struct {
	Filter       ChannelDetailFilter `json:"filter"`
	TopCorporate PaginatedCorporate  `json:"top_corporate"`
	TopKIP       PaginatedKIP        `json:"top_kip"`
}

type PaginatedCorporate struct {
	Data       []TopCorporate `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

type PaginatedKIP struct {
	Data       []TopKIP   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
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