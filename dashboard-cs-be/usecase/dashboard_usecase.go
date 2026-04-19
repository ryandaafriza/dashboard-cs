package usecase

import "dashboard-cs-be/entities"

type DashboardUsecase interface {
	GetDashboard(from, to string) (*entities.DashboardResponse, error)
	GetRealtime() (*entities.RealtimeResponse, error)
}
