package repository

import "dashboard-cs-be/entities"

// DashboardRepository is the single data-access contract.
// The MySQL implementation lives in mysql_dashboard_repository.go.
// Swapping to a different DB engine = new file + one line change in main.go.
type DashboardRepository interface {
	// GetSummary returns aggregated ticket + CSAT totals for a date range.
	GetSummary(from, to string) (*entities.SummaryRow, error)

	// GetDailyTrend returns created/solved counts grouped by calendar date.
	GetDailyTrend(from, to string) ([]entities.TrendRow, error)

	// GetTicketsPerHour returns created/solved counts grouped by hour-of-day.
	GetTicketsPerHour(from, to string) ([]entities.HourlyRow, error)

	// GetPrioritySummary returns open ticket counts per priority bucket.
	GetPrioritySummary(from, to string) (*entities.PriorityRow, error)

	// GetChannelSLA returns SLA %, open, and closed counts per channel.
	GetChannelSLA(from, to string) ([]entities.ChannelSLARow, error)

	// GetTopCorporate returns top companies by interaction count, per channel.
	GetTopCorporate(from, to string) ([]entities.TopCorporateRow, error)

	// GetTopKIP returns top topics by interaction count, per channel.
	GetTopKIP(from, to string) ([]entities.TopKIPRow, error)

	// GetRealtime returns live metrics (today only, no date filter).
	GetRealtime() (*entities.RealtimeRow, error)
}
