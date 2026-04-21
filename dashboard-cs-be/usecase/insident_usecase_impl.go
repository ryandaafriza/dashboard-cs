package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/repository/interfaces"
	mysqlrepo "dashboard-cs-be/repository/mysql"
)

const incidentDateLayout = "2006-01-02 15:04:05"

var validSeverities = map[string]bool{
	"low": true, "medium": true, "high": true, "critical": true,
}

type incidentUsecase struct {
	repo interfaces.IncidentRepository
	db   *sql.DB // dibutuhkan hanya untuk generate ID
}

// NewIncidentUsecase constructs the incident usecase.
// db diperlukan semata untuk NextIncidentID — alternatifnya bisa dipindah ke repo.
func NewIncidentUsecase(repo interfaces.IncidentRepository, db *sql.DB) IncidentUsecase {
	return &incidentUsecase{repo: repo, db: db}
}

// ─────────────────────────────────────────────────────────────────────────────
// CreateIncident
// ─────────────────────────────────────────────────────────────────────────────

func (uc *incidentUsecase) CreateIncident(req *entities.CreateIncidentRequest) (*entities.IncidentResponse, error) {
	// Validasi
	if req.Title == "" {
		return nil, fmt.Errorf("title tidak boleh kosong")
	}
	if !validSeverities[req.Severity] {
		return nil, fmt.Errorf("severity '%s' tidak valid, pilihan: low, medium, high, critical", req.Severity)
	}
	if req.CreatedBy == "" {
		return nil, fmt.Errorf("created_by tidak boleh kosong")
	}

	// Parse started_at — default ke waktu sekarang jika kosong
	startedAt := time.Now()
	if req.StartedAt != "" {
		t, err := time.ParseInLocation(incidentDateLayout, req.StartedAt, time.Local)
		if err != nil {
			return nil, fmt.Errorf("started_at tidak valid, gunakan format YYYY-MM-DD HH:MM:SS")
		}
		startedAt = t
	}

	// Generate ID
	id, err := mysqlrepo.NextIncidentID(uc.db)
	if err != nil {
		return nil, fmt.Errorf("generate incident ID: %w", err)
	}

	inc := &entities.Incident{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Severity:    req.Severity,
		StartedAt:   startedAt,
		CreatedBy:   req.CreatedBy,
	}

	saved, err := uc.repo.Create(inc)
	if err != nil {
		return nil, fmt.Errorf("CreateIncident: %w", err)
	}

	return toResponse(saved), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// GetActive
// ─────────────────────────────────────────────────────────────────────────────

func (uc *incidentUsecase) GetActive() (*entities.ActiveIncidentsSummary, error) {
	incidents, err := uc.repo.GetActive()
	if err != nil {
		return nil, fmt.Errorf("GetActive: %w", err)
	}

	responses := make([]entities.IncidentResponse, 0, len(incidents))
	for _, inc := range incidents {
		responses = append(responses, *toResponse(&inc))
	}

	return &entities.ActiveIncidentsSummary{
		Count:     len(responses),
		Incidents: responses,
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// GetHistory
// ─────────────────────────────────────────────────────────────────────────────

func (uc *incidentUsecase) GetHistory(from, to string) ([]entities.IncidentResponse, error) {
	if err := validateDateRange(from, to); err != nil {
		return nil, err
	}

	incidents, err := uc.repo.GetHistory(from, to)
	if err != nil {
		return nil, fmt.Errorf("GetHistory: %w", err)
	}

	responses := make([]entities.IncidentResponse, 0, len(incidents))
	for _, inc := range incidents {
		responses = append(responses, *toResponse(&inc))
	}
	return responses, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// ResolveIncident
// ─────────────────────────────────────────────────────────────────────────────

func (uc *incidentUsecase) ResolveIncident(id string) (*entities.IncidentResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("incident ID tidak boleh kosong")
	}

	resolved, err := uc.repo.Resolve(id)
	if err != nil {
		if errors.Is(err, mysqlrepo.ErrIncidentNotFound) {
			return nil, fmt.Errorf("incident '%s' tidak ditemukan", id)
		}
		if errors.Is(err, mysqlrepo.ErrIncidentAlreadyResolved) {
			return nil, fmt.Errorf("incident '%s' sudah diselesaikan sebelumnya", id)
		}
		return nil, fmt.Errorf("ResolveIncident: %w", err)
	}

	return toResponse(resolved), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Mapper
// ─────────────────────────────────────────────────────────────────────────────

func toResponse(inc *entities.Incident) *entities.IncidentResponse {
	r := &entities.IncidentResponse{
		ID:          inc.ID,
		Title:       inc.Title,
		Description: inc.Description,
		Severity:    inc.Severity,
		Status:      inc.Status,
		StartedAt:   inc.StartedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:   inc.CreatedBy,
		CreatedAt:   inc.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if inc.ResolvedAt != nil {
		s := inc.ResolvedAt.Format("2006-01-02 15:04:05")
		r.ResolvedAt = &s

		// Hitung durasi dalam menit (dibulatkan ke atas)
		dur := int(math.Ceil(inc.ResolvedAt.Sub(inc.StartedAt).Minutes()))
		r.DurationMin = &dur
	}

	return r
}
