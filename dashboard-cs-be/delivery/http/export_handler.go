package http

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/usecase"
)

// ExportHandler menangani request export laporan ke Excel.
type ExportHandler struct {
	uc usecase.ExportUsecase
}

// NewExportHandler constructs the export handler.
func NewExportHandler(uc usecase.ExportUsecase) *ExportHandler {
	return &ExportHandler{uc: uc}
}

// GET /api/v1/export?from=YYYY-MM-DD&to=YYYY-MM-DD&channel=all
// GET /api/v1/export?from=YYYY-MM-DD&to=YYYY-MM-DD&channel=whatsapp
// GET /api/v1/export?from=YYYY-MM-DD&to=YYYY-MM-DD&channel=email,live_chat
//
// Response: file .xlsx (binary download)
func (h *ExportHandler) ExportExcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, fail("method not allowed"))
		return
	}

	q := r.URL.Query()

	// ── Validasi tanggal ──────────────────────────────────────────────────
	from := q.Get("from")
	to := q.Get("to")
	if from == "" || to == "" {
		writeJSON(w, http.StatusBadRequest,
			fail("parameter 'from' dan 'to' wajib diisi (format: YYYY-MM-DD)"))
		return
	}

	// ── Parse & validasi channel ──────────────────────────────────────────
	channels, err := parseChannelParam(q.Get("channel"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, fail(err.Error()))
		return
	}

	filter := entities.ExportFilter{
		From:     from,
		To:       to,
		Channels: channels,
	}

	// ── Generate Excel ke buffer (supaya header bisa di-set sebelum stream) ──
	var buf bytes.Buffer
	filename, err := h.uc.ExportExcel(filter, &buf)
	if err != nil {
		log.Printf("ExportExcel: %v", err)
		writeJSON(w, http.StatusInternalServerError,
			fail("gagal membuat laporan: "+err.Error()))
		return
	}

	// ── Set response headers & kirim file ────────────────────────────────
	w.Header().Set("Content-Type",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf("ExportExcel write response: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper
// ─────────────────────────────────────────────────────────────────────────────

var validChannels = map[string]bool{
	"email": true, "whatsapp": true, "social_media": true,
	"live_chat": true, "call_center": true,
}

// parseChannelParam mem-parse query param channel.
//   "all" atau kosong  → nil (semua channel)
//   "email,live_chat"  → ["email", "live_chat"]
func parseChannelParam(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "all") {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		ch := strings.TrimSpace(strings.ToLower(p))
		if ch == "" {
			continue
		}
		if !validChannels[ch] {
			return nil, fmt.Errorf(
				"channel '%s' tidak valid. Pilihan: email, whatsapp, social_media, live_chat, call_center, all",
				ch)
		}
		result = append(result, ch)
	}
	return result, nil
}