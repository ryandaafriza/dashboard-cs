package http

import (
	"log"
	"net/http"
	"strings"

	"dashboard-cs-be/usecase"
)

// ImportHandler menangani upload & import file Excel.
type ImportHandler struct {
	uc usecase.ImportUsecase
}

// NewImportHandler constructs the import handler.
func NewImportHandler(uc usecase.ImportUsecase) *ImportHandler {
	return &ImportHandler{uc: uc}
}

// POST /api/v1/import
// Content-Type: multipart/form-data
// Field: file  (file .xlsx)
//
// Contoh curl:
//   curl -X POST http://localhost:8080/api/v1/import \
//        -F "file=@data_tiket.xlsx"
func (h *ImportHandler) ImportExcel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, fail("method not allowed"))
		return
	}

	// Batasi ukuran request: maks 20 MB
	const maxSize = 20 << 20 // 20 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	if err := r.ParseMultipartForm(maxSize); err != nil {
		writeJSON(w, http.StatusBadRequest, fail("file terlalu besar atau bukan multipart/form-data (maks 20 MB)"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, fail("field 'file' tidak ditemukan dalam form"))
		return
	}
	defer file.Close()

	// Validasi ekstensi file
	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") {
		writeJSON(w, http.StatusBadRequest, fail("hanya file .xlsx yang didukung"))
		return
	}

	result, err := h.uc.ImportExcel(file, filename)
	if err != nil {
		log.Printf("ImportExcel [%s]: %v", filename, err)
		writeJSON(w, http.StatusBadRequest, fail(err.Error()))
		return
	}

	writeJSON(w, http.StatusOK, ok(result))
}