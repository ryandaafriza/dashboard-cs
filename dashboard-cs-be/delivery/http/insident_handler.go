package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/usecase"
)

// IncidentHandler menangani semua endpoint /api/v1/incidents/*.
type IncidentHandler struct {
	uc usecase.IncidentUsecase
}

// NewIncidentHandler constructs the incident handler.
func NewIncidentHandler(uc usecase.IncidentUsecase) *IncidentHandler {
	return &IncidentHandler{uc: uc}
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /api/v1/incidents/active
// Tidak ada filter tanggal — selalu kondisi saat ini.
// Response:
//   {
//     "success": true,
//     "data": {
//       "count": 2,
//       "incidents": [...]
//     }
//   }
// ─────────────────────────────────────────────────────────────────────────────

func (h *IncidentHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, fail("method not allowed"))
		return
	}

	summary, err := h.uc.GetActive()
	if err != nil {
		log.Printf("GetActive incidents: %v", err)
		writeJSON(w, http.StatusInternalServerError, fail("gagal mengambil data incident aktif"))
		return
	}
	writeJSON(w, http.StatusOK, ok(summary))
}

// ─────────────────────────────────────────────────────────────────────────────
// GET /api/v1/incidents/history?from=YYYY-MM-DD&to=YYYY-MM-DD
// ─────────────────────────────────────────────────────────────────────────────

func (h *IncidentHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, fail("method not allowed"))
		return
	}

	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")
	today := time.Now().Format("2006-01-02")
	if from == "" {
		from = today
	}
	if to == "" {
		to = today
	}

	history, err := h.uc.GetHistory(from, to)
	if err != nil {
		log.Printf("GetHistory incidents: %v", err)
		writeJSON(w, http.StatusBadRequest, fail(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, ok(history))
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /api/v1/incidents
// Body JSON:
//   {
//     "title":       "WhatsApp Gateway Down",
//     "description": "Gateway tidak bisa menerima pesan masuk",
//     "severity":    "high",
//     "started_at":  "2026-04-19 09:00:00",   ← opsional, default: now
//     "created_by":  "Admin Ops"
//   }
// ─────────────────────────────────────────────────────────────────────────────

func (h *IncidentHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, fail("method not allowed"))
		return
	}

	var req entities.CreateIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, fail("request body tidak valid: "+err.Error()))
		return
	}
	defer r.Body.Close()

	result, err := h.uc.CreateIncident(&req)
	if err != nil {
		log.Printf("Create incident: %v", err)
		writeJSON(w, http.StatusBadRequest, fail(err.Error()))
		return
	}
	writeJSON(w, http.StatusCreated, ok(result))
}

// ─────────────────────────────────────────────────────────────────────────────
// PATCH /api/v1/incidents/{id}/resolve
// Tidak perlu body — cukup ID dari URL path.
// ─────────────────────────────────────────────────────────────────────────────

func (h *IncidentHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeJSON(w, http.StatusMethodNotAllowed, fail("method not allowed"))
		return
	}

	// Ekstrak ID dari path: /api/v1/incidents/{id}/resolve
	// Setelah Trim("/"):  "api/v1/incidents/{id}/resolve"
	// Split result:       ["api","v1","incidents","{id}","resolve"]
	//                       [0]  [1]     [2]        [3]     [4]
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 5 || parts[3] == "" || parts[3] == "resolve" {
		writeJSON(w, http.StatusBadRequest, fail("incident ID tidak ditemukan di URL"))
		return
	}
	id := parts[3]

	result, err := h.uc.ResolveIncident(id)
	if err != nil {
		log.Printf("Resolve incident %s: %v", id, err)
		// Bedakan not found vs already resolved vs error lain
		msg := err.Error()
		status := http.StatusBadRequest
		if strings.Contains(msg, "tidak ditemukan") {
			status = http.StatusNotFound
		}
		writeJSON(w, status, fail(msg))
		return
	}
	writeJSON(w, http.StatusOK, ok(result))
}