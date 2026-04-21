package usecase

import "dashboard-cs-be/entities"

type DashboardUsecase interface {
	// GetDashboard mengembalikan summary, trend, hourly, priority, dan channel SLA
	// sesuai filter tanggal. Channel hanya berisi SLA%, open, closed — tanpa detail.
	GetDashboard(from, to string) (*entities.DashboardResponse, error)

	// GetChannelDetail mengembalikan top_corporate dan top_kip untuk satu channel
	// spesifik dengan pagination.
	GetChannelDetail(f entities.ChannelDetailFilter) (*entities.ChannelDetailResponse, error)

	// GetRealtime mengembalikan metrik live (SLA hari ini, created today, incidents).
	GetRealtime() (*entities.RealtimeResponse, error)
}