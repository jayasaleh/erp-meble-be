package dto

import "time"

// ===========================
// REQUEST DTOs
// ===========================

// CreateSalesRequest adalah DTO untuk membuat transaksi penjualan baru (Mode 1: POS)
// Untuk pembayaran transfer dengan bukti foto, gunakan multipart/form-data.
// Field JSON dikirim sebagai field "data" (string JSON), file sebagai "bukti_bayar".
type CreateSalesRequest struct {
	IDGudang         uint               `json:"id_gudang" binding:"required"`
	NamaPelanggan    string             `json:"nama_pelanggan"`   // Opsional
	KontakPelanggan  string             `json:"kontak_pelanggan"` // Opsional
	MetodePembayaran string             `json:"metode_pembayaran" binding:"required,oneof=cash transfer"`
	JumlahPembayaran float64            `json:"jumlah_pembayaran" binding:"required,gt=0"`
	CatatanInternal  string             `json:"catatan_internal"`
	Items            []SalesItemRequest `json:"items" binding:"required,min=1,dive"`
}

// SalesItemRequest adalah DTO untuk setiap item dalam transaksi penjualan
type SalesItemRequest struct {
	IDProduk     uint     `json:"id_produk" binding:"required"`
	Jumlah       int      `json:"jumlah" binding:"required,min=1"`
	HargaSatuan  float64  `json:"harga_satuan" binding:"required,gt=0"`
	PersenDiskon *float64 `json:"persen_diskon" binding:"omitempty,min=0,max=100"` // 0–100 (%)
}

// ListSalesRequest adalah DTO untuk filter list transaksi penjualan
type ListSalesRequest struct {
	Page             int        `form:"page" binding:"omitempty,min=1"`
	Limit            int        `form:"limit" binding:"omitempty,min=1,max=100"`
	TanggalDari      *time.Time `form:"tanggal_dari" time_format:"2006-01-02"`
	TanggalSampai    *time.Time `form:"tanggal_sampai" time_format:"2006-01-02"`
	IDKasir          *uint      `form:"id_kasir"`
	IDGudang         *uint      `form:"id_gudang"`
	MetodePembayaran string     `form:"metode_pembayaran" binding:"omitempty,oneof=cash transfer"`
}

// ===========================
// RESPONSE DTOs
// ===========================

// SalesResponse adalah DTO untuk list view transaksi penjualan
type SalesResponse struct {
	ID               uint      `json:"id"`
	NomorTransaksi   string    `json:"nomor_transaksi"`
	NamaGudang       string    `json:"nama_gudang"`
	NamaPelanggan    string    `json:"nama_pelanggan"`
	KontakPelanggan  string    `json:"kontak_pelanggan"`
	Total            float64   `json:"total"`
	TotalHargaModal  float64   `json:"total_harga_modal"` // COGS total
	Laba             float64   `json:"laba"`              // Total - TotalHargaModal
	MetodePembayaran string    `json:"metode_pembayaran"`
	Status           string    `json:"status"`
	NamaKasir        string    `json:"nama_kasir"`
	DibuatPada       time.Time `json:"dibuat_pada"`
}

// SalesItemResponse adalah DTO untuk detail item penjualan
type SalesItemResponse struct {
	ID           uint                      `json:"id"`
	IDProduk     uint                      `json:"id_produk"`
	SKUProduk    string                    `json:"sku_produk"`
	NamaProduk   string                    `json:"nama_produk"`
	Jumlah       int                       `json:"jumlah"`
	HargaSatuan  float64                   `json:"harga_satuan"`
	HargaModal   float64                   `json:"harga_modal"` // COGS per unit (rata-rata tertimbang)
	PersenDiskon *float64                  `json:"persen_diskon,omitempty"`
	JumlahDiskon float64                   `json:"jumlah_diskon"`
	Subtotal     float64                   `json:"subtotal"`
	TotalModal   float64                   `json:"total_modal"` // COGS total item ini
	Laba         float64                   `json:"laba"`        // Subtotal - TotalModal
	BatchUsage   []SalesBatchUsageResponse `json:"batch_usage,omitempty"`
}

// SalesBatchUsageResponse adalah DTO untuk detail batch FIFO yang digunakan per item
type SalesBatchUsageResponse struct {
	IDBatch    uint    `json:"id_batch"`
	Jumlah     int     `json:"jumlah"`
	HargaModal float64 `json:"harga_modal"`
	TotalModal float64 `json:"total_modal"`
}

// SalesDetailResponse adalah DTO untuk detail lengkap satu transaksi penjualan
type SalesDetailResponse struct {
	ID               uint                `json:"id"`
	NomorTransaksi   string              `json:"nomor_transaksi"`
	IDGudang         uint                `json:"id_gudang"`
	NamaGudang       string              `json:"nama_gudang"`
	NamaPelanggan    string              `json:"nama_pelanggan"`
	KontakPelanggan  string              `json:"kontak_pelanggan"`
	Subtotal         float64             `json:"subtotal"`
	JumlahDiskon     float64             `json:"jumlah_diskon"`
	Total            float64             `json:"total"`
	TotalHargaModal  float64             `json:"total_harga_modal"`
	Laba             float64             `json:"laba"`
	MetodePembayaran string              `json:"metode_pembayaran"`
	JumlahPembayaran float64             `json:"jumlah_pembayaran"`
	JumlahKembalian  float64             `json:"jumlah_kembalian"`
	BuktiBayar       *string             `json:"bukti_bayar,omitempty"`
	Status           string              `json:"status"`
	CatatanInternal  string              `json:"catatan_internal,omitempty"`
	IDKasir          uint                `json:"id_kasir"`
	NamaKasir        string              `json:"nama_kasir"`
	DibuatPada       time.Time           `json:"dibuat_pada"`
	Items            []SalesItemResponse `json:"items"`
}

// InvoiceResponse adalah DTO yang dioptimalkan untuk keperluan cetak / ekspor invoice
type InvoiceResponse struct {
	// Header Invoice
	NomorInvoice   string    `json:"nomor_invoice"`
	TanggalInvoice time.Time `json:"tanggal_invoice"`

	// Info Toko / Gudang
	NamaGudang string `json:"nama_gudang"`

	// Info Pelanggan
	NamaPelanggan   string `json:"nama_pelanggan"`
	KontakPelanggan string `json:"kontak_pelanggan"`

	// Info Kasir
	NamaKasir string `json:"nama_kasir"`

	// Line Items
	Items []InvoiceItemResponse `json:"items"`

	// Totals
	Subtotal         float64 `json:"subtotal"`
	TotalDiskon      float64 `json:"total_diskon"`
	Total            float64 `json:"total"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	JumlahPembayaran float64 `json:"jumlah_pembayaran"`
	JumlahKembalian  float64 `json:"jumlah_kembalian"`

	// Status
	Status string `json:"status"`
}

// InvoiceItemResponse adalah DTO untuk baris item di invoice
type InvoiceItemResponse struct {
	NoProduk    int     `json:"no"`
	SKU         string  `json:"sku"`
	NamaProduk  string  `json:"nama_produk"`
	Jumlah      int     `json:"jumlah"`
	HargaSatuan float64 `json:"harga_satuan"`
	Diskon      float64 `json:"diskon"`
	Subtotal    float64 `json:"subtotal"`
}

// ListSalesResponse adalah DTO untuk response list penjualan dengan pagination
type ListSalesResponse struct {
	Sales      []SalesResponse `json:"sales"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}
