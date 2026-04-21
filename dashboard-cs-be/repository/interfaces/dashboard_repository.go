package interfaces

import "dashboard-cs-be/entities"

type DashboardRepository interface {
	// GetSummary returns aggregated ticket + CSAT totals for a date range.
	GetSummary(from, to string) (*entities.SummaryRow, error)

	// GetDailyTrend returns created/solved counts grouped by calendar date.
	GetDailyTrend(from, to string) ([]entities.TrendRow, error)

	// GetTicketsPerHour returns created/solved counts grouped by hour-of-day.
	GetTicketsPerHour(from, to string) ([]entities.HourlyRow, error)

	// GetPrioritySummary returns ticket counts per priority bucket.
	GetPrioritySummary(from, to string) (*entities.PriorityRow, error)

	// GetChannelSLA returns SLA %, open, and closed counts per channel.
	// Tidak lagi include top_corporate & top_kip — sudah dipisah ke GetChannelDetail.
	GetChannelSLA(from, to string) ([]entities.ChannelSLARow, error)

	// GetTopCorporate returns paginated top corporate customers for a specific channel.
	GetTopCorporate(f entities.ChannelDetailFilter) ([]entities.TopCorporateRow, error)

	// GetTopKIP returns paginated top KIP topics for a specific channel.
	GetTopKIP(f entities.ChannelDetailFilter) ([]entities.TopKIPRow, error)

	// GetRealtime returns live metrics (today only, no date filter).
	GetRealtime() (*entities.RealtimeRow, error)
}