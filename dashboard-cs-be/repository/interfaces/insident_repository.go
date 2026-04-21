package interfaces

import "dashboard-cs-be/entities"

// IncidentRepository mendefinisikan kontrak akses data untuk incidents.
type IncidentRepository interface {
	// Create menyimpan incident baru ke DB dan mengembalikan entitas lengkapnya.
	Create(inc *entities.Incident) (*entities.Incident, error)

	// GetActive mengembalikan semua incident berstatus 'active',
	// diurutkan dari yang paling baru.
	GetActive() ([]entities.Incident, error)

	// GetHistory mengembalikan incident dalam rentang started_at,
	// termasuk yang sudah resolved. Diurutkan started_at DESC.
	GetHistory(from, to string) ([]entities.Incident, error)

	// Resolve menandai incident sebagai 'resolved' dengan waktu sekarang.
	// Mengembalikan ErrNotFound jika ID tidak ada, atau error jika
	// incident sudah resolved sebelumnya.
	Resolve(id string) (*entities.Incident, error)
}