package entities

import "time"

// ─────────────────────────────────────────────────────────────────────────────
// Domain struct
// ─────────────────────────────────────────────────────────────────────────────

type Incident struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Severity    string     `db:"severity"`  // low | medium | high | critical
	Status      string     `db:"status"`    // active | resolved
	StartedAt   time.Time  `db:"started_at"`
	ResolvedAt  *time.Time `db:"resolved_at"` // nil jika masih active
	CreatedBy   string     `db:"created_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Request / Response structs (usecase ↔ handler)
// ─────────────────────────────────────────────────────────────────────────────

// CreateIncidentRequest adalah payload POST /api/v1/incidents.
type CreateIncidentRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`    // low | medium | high | critical
	StartedAt   string `json:"started_at"`  // "YYYY-MM-DD HH:MM:SS", kosong = now
	CreatedBy   string `json:"created_by"`  // nama/ID admin yang input
}

// IncidentResponse adalah representasi Incident untuk JSON response.
type IncidentResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"`
	Status      string  `json:"status"`
	StartedAt   string  `json:"started_at"`
	ResolvedAt  *string `json:"resolved_at"` // null jika masih active
	CreatedBy   string  `json:"created_by"`
	CreatedAt   string  `json:"created_at"`
	DurationMin *int    `json:"duration_minutes"` // null jika masih active
}

// ActiveIncidentsSummary adalah payload GET /api/v1/incidents/active.
type ActiveIncidentsSummary struct {
	Count     int                `json:"count"`
	Incidents []IncidentResponse `json:"incidents"`
}