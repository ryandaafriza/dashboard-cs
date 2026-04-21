package interfaces

import "dashboard-cs-be/entities"

// ImportRepository menangani operasi upsert data hasil parsing Excel ke DB.
type ImportRepository interface {
	// UpsertTicket menyimpan atau memperbarui satu tiket beserta customer-nya.
	// Logika:
	//   - Customer di-upsert berdasarkan phone (identifier unik).
	//   - Jika ticket_id belum ada → INSERT baru.
	//   - Jika ticket_id sudah ada:
	//       • Jika status di DB sudah 'closed' → SKIP (tidak boleh dibuka ulang).
	//       • Jika status di DB masih 'open'   → UPDATE semua kolom.
	//   - Agent di-upsert (insert ignore) jika agent_id tidak kosong.
	// Return: "inserted" | "updated" | "skipped"
	UpsertTicket(row *entities.ImportRow) (string, error)

	// SaveImportLog menyimpan ringkasan hasil import ke tabel import_logs.
	SaveImportLog(result *entities.ImportResult) error
}