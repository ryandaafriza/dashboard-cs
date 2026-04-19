package usecase

import (
	"io"

	"dashboard-cs-be/entities"
)

// ExportUsecase menangani logika bisnis export Excel.
type ExportUsecase interface {
	// ExportExcel mengambil seluruh data sesuai filter, lalu menulis file .xlsx
	// ke writer yang diberikan. Nama file dikembalikan untuk di-set sebagai
	// Content-Disposition header oleh handler.
	ExportExcel(filter entities.ExportFilter, w io.Writer) (filename string, err error)
}