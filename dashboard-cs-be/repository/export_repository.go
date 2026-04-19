package repository

import "dashboard-cs-be/entities"

// ExportRepository menyediakan seluruh data mentah yang dibutuhkan untuk
// generate laporan Excel. Semua method menerima ExportFilter sehingga
// filter channel & tanggal diterapkan konsisten di level SQL.
type ExportRepository interface {
	// GetExportSummary mengembalikan ringkasan tiket (total, open, closed, CSAT)
	// dengan filter tanggal + channel opsional.
	GetExportSummary(f entities.ExportFilter) (*entities.SummaryRow, error)

	// GetExportChannels mengembalikan performa per channel (Sheet 2).
	GetExportChannels(f entities.ExportFilter) ([]entities.ExportChannelRow, error)

	// GetExportCustomers mengembalikan data per customer (Sheet 3).
	GetExportCustomers(f entities.ExportFilter) ([]entities.ExportCustomerRow, error)

	// GetExportTopics mengembalikan data per topik + channel (Sheet 4).
	GetExportTopics(f entities.ExportFilter) ([]entities.ExportTopicRow, error)

	// GetExportPriorities mengembalikan data per prioritas (Sheet 5).
	GetExportPriorities(f entities.ExportFilter) ([]entities.ExportPriorityRow, error)
}