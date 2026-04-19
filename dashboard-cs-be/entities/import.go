package entities

// ─────────────────────────────────────────────────────────────────────────────
// Import Excel – row parsed dari file .xlsx
// ─────────────────────────────────────────────────────────────────────────────

// ImportRow merepresentasikan satu baris data dari file Excel yang diupload.
// Kolom yang bertipe waktu (created_at / resolved_at) sudah diparse menjadi string
// format "YYYY-MM-DD HH:MM:SS" oleh usecase sebelum diteruskan ke repository.
type ImportRow struct {
	TicketID     string  // TKT-YYYYMMDD-XXX
	CreatedAt    string  // "2006-01-02 15:04:05"
	ResolvedAt   string  // "" jika kosong
	Channel      string  // email / whatsapp / social_media / live_chat / call_center
	Priority     string  // p1 / urgent / vip / cc / roaming / extra_quota / normal
	Status       string  // open / closed
	CustomerName string
	CustomerPhone string
	CustomerEmail string  // boleh kosong
	CustomerType string   // consumer / corporate
	Topic        string
	AgentID      string  // boleh kosong
	IsFCR        bool
}

// ─────────────────────────────────────────────────────────────────────────────
// Import Result – dikembalikan ke handler
// ─────────────────────────────────────────────────────────────────────────────

// ImportRowError menyimpan detail error satu baris.
type ImportRowError struct {
	Row      int    `json:"row"`
	TicketID string `json:"ticket_id,omitempty"`
	Reason   string `json:"reason"`
}

// ImportResult adalah ringkasan hasil proses import satu file.
type ImportResult struct {
	Filename   string           `json:"filename"`
	TotalRows  int              `json:"total_rows"`
	Inserted   int              `json:"inserted"`
	Updated    int              `json:"updated"`
	Skipped    int              `json:"skipped"`
	ErrorCount int              `json:"error_count"`
	Errors     []ImportRowError `json:"errors,omitempty"`
}