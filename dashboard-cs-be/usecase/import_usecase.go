package usecase

import (
	"mime/multipart"

	"dashboard-cs-be/entities"
)

// ImportUsecase menangani logika bisnis import Excel.
type ImportUsecase interface {
	// ImportExcel mem-parsing file Excel yang diterima dari multipart upload,
	// memvalidasi setiap baris, lalu mengorkestrasikan upsert ke repository.
	// Parameter filename digunakan untuk logging & audit trail.
	ImportExcel(file multipart.File, filename string) (*entities.ImportResult, error)
}