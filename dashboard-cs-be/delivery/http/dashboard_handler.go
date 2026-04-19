package http

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/usecase"
)

type DashboardHandler struct {
	uc usecase.DashboardUsecase
}

func NewDashboardHandler(uc usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{uc: uc}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}

func ok(data interface{}) entities.Response {
	return entities.Response{Success: true, Message: "OK", Data: data}
}

func fail(msg string) entities.Response {
	return entities.Response{Success: false, Message: msg, Data: nil}
}

// GET /api/v1/dashboard?from=YYYY-MM-DD&to=YYYY-MM-DD
func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
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

	data, err := h.uc.GetDashboard(from, to)
	if err != nil {
		log.Printf("GetDashboard: %v", err)
		writeJSON(w, http.StatusBadRequest, fail(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, ok(data))
}

// GET /api/v1/realtime
func (h *DashboardHandler) GetRealtime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, fail("method not allowed"))
		return
	}

	data, err := h.uc.GetRealtime()
	if err != nil {
		log.Printf("GetRealtime: %v", err)
		writeJSON(w, http.StatusInternalServerError, fail("failed to fetch realtime data"))
		return
	}
	writeJSON(w, http.StatusOK, ok(data))
}

// GET /health
func (h *DashboardHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
