package http

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		qs := ""
		if r.URL.RawQuery != "" {
			qs = "?" + r.URL.RawQuery
		}
		log.Printf("[%s] %s%s  %s", r.Method, r.URL.Path, qs, time.Since(start))
	})
}

func NewRouter(
	dh *DashboardHandler,
	ih *ImportHandler,
	eh *ExportHandler,
	inch *IncidentHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", dh.HealthCheck)

	// Realtime — tanpa filter tanggal
	mux.HandleFunc("/api/v1/realtime", dh.GetRealtime)

	// Dashboard — dengan filter tanggal
	// PENTING: /api/v1/dashboard/channels harus didaftarkan SEBELUM /api/v1/dashboard
	// supaya tidak tertangkap oleh handler dashboard utama
	mux.HandleFunc("/api/v1/dashboard/channels", dh.GetChannelDetail)
	mux.HandleFunc("/api/v1/dashboard", dh.GetDashboard)

	// Import & Export
	mux.HandleFunc("/api/v1/import", ih.ImportExcel)
	mux.HandleFunc("/api/v1/export", eh.ExportExcel)

	// Incidents
	mux.HandleFunc("/api/v1/incidents", inch.Create)
	mux.HandleFunc("/api/v1/incidents/active", inch.GetActive)
	mux.HandleFunc("/api/v1/incidents/history", inch.GetHistory)
	mux.HandleFunc("/api/v1/incidents/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/incidents/")
		if strings.HasSuffix(path, "/resolve") && r.Method == http.MethodPatch {
			inch.Resolve(w, r)
			return
		}
		writeJSON(w, http.StatusNotFound, fail("endpoint tidak ditemukan"))
	})

	return loggingMiddleware(corsMiddleware(mux))
}