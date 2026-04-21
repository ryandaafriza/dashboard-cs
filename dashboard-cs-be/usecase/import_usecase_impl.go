package usecase

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"dashboard-cs-be/entities"
	"dashboard-cs-be/repository/interfaces"
)

// ─────────────────────────────────────────────────────────────────────────────
// Constants & allowed values
// ─────────────────────────────────────────────────────────────────────────────

var allowedChannels = map[string]bool{
	"email": true, "whatsapp": true, "social_media": true,
	"live_chat": true, "call_center": true,
}

var allowedPriorities = map[string]bool{
	"p1": true, "urgent": true, "vip": true,
	"cc": true, "roaming": true, "extra_quota": true, "normal": true,
}

var allowedStatuses = map[string]bool{
	"open": true, "closed": true,
}

var allowedCustomerTypes = map[string]bool{
	"consumer": true, "corporate": true,
}

// Datetime formats yang diterima dari Excel
var dateTimeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006/01/02 15:04:05",
	"2006-01-02",
	"01/02/2006 15:04:05",
	"02-01-2006 15:04:05",
}

// Header Excel yang diharapkan (lowercase, trimmed)
// Urutan kolom fleksibel – kita pakai header mapping.
var expectedHeaders = []string{
	"ticket_id", "created_at", "resolved_at", "channel", "priority",
	"status", "customer_name", "customer_phone", "customer_email",
	"customer_type", "topic", "agent_id",
}

// ─────────────────────────────────────────────────────────────────────────────
// Struct
// ─────────────────────────────────────────────────────────────────────────────

type importUsecase struct {
	repo interfaces.ImportRepository
}

// NewImportUsecase constructs the import usecase.
func NewImportUsecase(repo interfaces.ImportRepository) ImportUsecase {
	return &importUsecase{repo: repo}
}

// ─────────────────────────────────────────────────────────────────────────────
// ImportExcel – entry point
// ─────────────────────────────────────────────────────────────────────────────

func (uc *importUsecase) ImportExcel(file multipart.File, filename string) (*entities.ImportResult, error) {
	result := &entities.ImportResult{
		Filename: filename,
		Errors:   []entities.ImportRowError{},
	}

	// Buka file Excel dari reader
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("tidak dapat membuka file Excel: %w", err)
	}
	defer f.Close()

	// Ambil sheet pertama
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("file Excel tidak memiliki sheet")
	}
	sheetName := sheets[0]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca baris Excel: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("file Excel kosong atau hanya memiliki header")
	}

	// Bangun header map: nama_kolom → index
	headerMap, err := buildHeaderMap(rows[0])
	if err != nil {
		return nil, err
	}

	// Proses tiap baris data (mulai baris ke-2, index 1)
	for i, row := range rows[1:] {
		rowNum := i + 2 // nomor baris di Excel (1-based, header di baris 1)
		result.TotalRows++

		parsed, validationErr := parseRow(row, headerMap, rowNum)
		if validationErr != nil {
			result.Errors = append(result.Errors, *validationErr)
			result.ErrorCount++
			continue
		}

		action, dbErr := uc.repo.UpsertTicket(parsed)
		if dbErr != nil {
			result.Errors = append(result.Errors, entities.ImportRowError{
				Row:      rowNum,
				TicketID: parsed.TicketID,
				Reason:   dbErr.Error(),
			})
			result.ErrorCount++
			continue
		}

		switch action {
		case "inserted":
			result.Inserted++
		case "updated":
			result.Updated++
		case "skipped":
			result.Skipped++
		}
	}

	// Simpan audit log (non-fatal jika gagal)
	_ = uc.repo.SaveImportLog(result)

	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// buildHeaderMap – mapping nama kolom → index
// ─────────────────────────────────────────────────────────────────────────────

func buildHeaderMap(headerRow []string) (map[string]int, error) {
	m := make(map[string]int, len(headerRow))
	for i, h := range headerRow {
		normalized := strings.ToLower(strings.TrimSpace(h))
		m[normalized] = i
	}

	// Validasi kolom wajib
	required := []string{"ticket_id", "created_at", "channel", "priority", "status",
		"customer_name", "customer_phone", "topic"}
	for _, col := range required {
		if _, ok := m[col]; !ok {
			return nil, fmt.Errorf("kolom wajib '%s' tidak ditemukan di header Excel", col)
		}
	}
	return m, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// parseRow – parse & validasi satu baris
// ─────────────────────────────────────────────────────────────────────────────

func parseRow(row []string, hm map[string]int, rowNum int) (*entities.ImportRow, *entities.ImportRowError) {
	get := func(col string) string {
		idx, ok := hm[col]
		if !ok || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	rowErr := func(reason string, ticketID ...string) *entities.ImportRowError {
		tid := get("ticket_id")
		if len(ticketID) > 0 {
			tid = ticketID[0]
		}
		return &entities.ImportRowError{Row: rowNum, TicketID: tid, Reason: reason}
	}

	// ── ticket_id ──────────────────────────────────────────────────────────
	ticketID := get("ticket_id")
	if ticketID == "" {
		return nil, rowErr("ticket_id tidak boleh kosong")
	}

	// ── created_at ─────────────────────────────────────────────────────────
	createdAtRaw := get("created_at")
	createdAt, err := parseDateTime(createdAtRaw)
	if err != nil {
		return nil, rowErr(fmt.Sprintf("created_at tidak valid: %v", err), ticketID)
	}

	// ── resolved_at (opsional) ─────────────────────────────────────────────
	resolvedAt := ""
	if raw := get("resolved_at"); raw != "" {
		t, err := parseDateTime(raw)
		if err != nil {
			return nil, rowErr(fmt.Sprintf("resolved_at tidak valid: %v", err), ticketID)
		}
		resolvedAt = t
	}

	// ── channel ────────────────────────────────────────────────────────────
	channel := strings.ToLower(get("channel"))
	if !allowedChannels[channel] {
		return nil, rowErr(fmt.Sprintf("channel '%s' tidak valid", channel), ticketID)
	}

	// ── priority ───────────────────────────────────────────────────────────
	priority := strings.ToLower(get("priority"))
	if priority == "" {
		priority = "normal"
	}
	if !allowedPriorities[priority] {
		return nil, rowErr(fmt.Sprintf("priority '%s' tidak valid", priority), ticketID)
	}

	// ── status ─────────────────────────────────────────────────────────────
	status := strings.ToLower(get("status"))
	if !allowedStatuses[status] {
		return nil, rowErr(fmt.Sprintf("status '%s' tidak valid", status), ticketID)
	}

	// Konsistensi: jika status closed, resolved_at harus ada
	if status == "closed" && resolvedAt == "" {
		return nil, rowErr("status 'closed' membutuhkan resolved_at", ticketID)
	}

	// ── customer_name ──────────────────────────────────────────────────────
	customerName := get("customer_name")
	if customerName == "" {
		return nil, rowErr("customer_name tidak boleh kosong", ticketID)
	}

	// ── customer_phone ─────────────────────────────────────────────────────
	customerPhone := get("customer_phone")
	if customerPhone == "" {
		return nil, rowErr("customer_phone tidak boleh kosong", ticketID)
	}

	// ── customer_type ──────────────────────────────────────────────────────
	customerType := strings.ToLower(get("customer_type"))
	if customerType == "" {
		customerType = "consumer"
	}
	if !allowedCustomerTypes[customerType] {
		customerType = "consumer" // default graceful
	}

	// ── topic ──────────────────────────────────────────────────────────────
	topic := get("topic")
	if topic == "" {
		return nil, rowErr("topic tidak boleh kosong", ticketID)
	}

	return &entities.ImportRow{
		TicketID:      ticketID,
		CreatedAt:     createdAt,
		ResolvedAt:    resolvedAt,
		Channel:       channel,
		Priority:      priority,
		Status:        status,
		CustomerName:  customerName,
		CustomerPhone: customerPhone,
		CustomerEmail: get("customer_email"),
		CustomerType:  customerType,
		Topic:         topic,
		AgentID:       get("agent_id"),
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func parseDateTime(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("nilai kosong")
	}

	// Coba semua format yang dikenal
	for _, layout := range dateTimeFormats {
		if t, err := time.ParseInLocation(layout, raw, time.Local); err == nil {
			return t.Format("2006-01-02 15:04:05"), nil
		}
	}

	// Excelize kadang mengonversi tanggal jadi serial number float
	if f, err := strconv.ParseFloat(raw, 64); err == nil {
		t, err := excelize.ExcelDateToTime(f, false)
		if err == nil {
			return t.Format("2006-01-02 15:04:05"), nil
		}
	}

	return "", fmt.Errorf("format tidak dikenali: '%s'", raw)
}

func parseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes" || s == "ya"
}