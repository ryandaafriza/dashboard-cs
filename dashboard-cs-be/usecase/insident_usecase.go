package usecase

import "dashboard-cs-be/entities"

// IncidentUsecase mendefinisikan logika bisnis untuk incident management.
type IncidentUsecase interface {
	// CreateIncident memvalidasi request dan menyimpan incident baru.
	CreateIncident(req *entities.CreateIncidentRequest) (*entities.IncidentResponse, error)

	// GetActive mengembalikan semua incident yang sedang aktif.
	// Tidak ada filter tanggal — selalu kondisi saat ini.
	GetActive() (*entities.ActiveIncidentsSummary, error)

	// GetHistory mengembalikan riwayat incident dalam rentang tanggal.
	GetHistory(from, to string) ([]entities.IncidentResponse, error)

	// ResolveIncident menandai incident sebagai selesai.
	ResolveIncident(id string) (*entities.IncidentResponse, error)
}