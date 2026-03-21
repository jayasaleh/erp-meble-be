package dto

import "time"

// ===========================
// RETUR PENJUALAN (Customer → Toko)
// ===========================

// CreateReturPenjualanRequest adalah DTO untuk membuat retur dari customer
type CreateReturPenjualanRequest struct {
	IDPenjualan        uint                          `json:"id_penjualan" binding:"required"`
	Alasan             string                        `json:"alasan" binding:"required"` // rusak, tidak_sesuai, cacat, dll
	MetodePengembalian string                        `json:"metode_pengembalian" binding:"required,oneof=cash transfer tukar_barang"`
	Keterangan         string                        `json:"keterangan"`
	Items              []CreateReturPenjualanItemReq `json:"items" binding:"required,min=1,dive"`
}

// CreateReturPenjualanItemReq adalah DTO untuk item retur penjualan
type CreateReturPenjualanItemReq struct {
	IDItemPenjualan uint `json:"id_item_penjualan" binding:"required"` // Link ke item penjualan asal
	IDProduk        uint `json:"id_produk" binding:"required"`
	Jumlah          int  `json:"jumlah" binding:"required,min=1"`
}

// ReturPenjualanResponse adalah DTO untuk response retur penjualan
type ReturPenjualanResponse struct {
	ID                 uint                         `json:"id"`
	NomorRetur         string                       `json:"nomor_retur"`
	IDPenjualan        uint                         `json:"id_penjualan"`
	NomorTransaksiAsal string                       `json:"nomor_transaksi_asal"`
	NamaPelanggan      string                       `json:"nama_pelanggan"`
	KontakPelanggan    string                       `json:"kontak_pelanggan"`
	Alasan             string                       `json:"alasan"`
	Subtotal           float64                      `json:"subtotal"`
	Total              float64                      `json:"total"`
	MetodePengembalian string                       `json:"metode_pengembalian"`
	JumlahPengembalian float64                      `json:"jumlah_pengembalian"`
	Status             string                       `json:"status"`
	Keterangan         string                       `json:"keterangan"`
	NamaPetugas        string                       `json:"nama_petugas"`
	DibuatPada         time.Time                    `json:"dibuat_pada"`
	Items              []ReturPenjualanItemResponse `json:"items,omitempty"`
}

// ReturPenjualanItemResponse adalah DTO untuk item retur penjualan
type ReturPenjualanItemResponse struct {
	ID          uint    `json:"id"`
	IDProduk    uint    `json:"id_produk"`
	SKUProduk   string  `json:"sku_produk"`
	NamaProduk  string  `json:"nama_produk"`
	Jumlah      int     `json:"jumlah"`
	HargaSatuan float64 `json:"harga_satuan"`
	Subtotal    float64 `json:"subtotal"`
}

// ListReturPenjualanRequest adalah DTO untuk filter list retur penjualan
type ListReturPenjualanRequest struct {
	Page          int        `form:"page" binding:"omitempty,min=1"`
	Limit         int        `form:"limit" binding:"omitempty,min=1,max=100"`
	TanggalDari   *time.Time `form:"tanggal_dari" time_format:"2006-01-02"`
	TanggalSampai *time.Time `form:"tanggal_sampai" time_format:"2006-01-02"`
	Status        string     `form:"status" binding:"omitempty,oneof=pending approved completed rejected"`
}

// ===========================
// RETUR PEMBELIAN (Toko → Vendor/Supplier)
// ===========================

// CreateReturPembelianRequest adalah DTO untuk membuat retur ke supplier
type CreateReturPembelianRequest struct {
	IDPemasok          uint                          `json:"id_pemasok" binding:"required"`
	IDGudang           uint                          `json:"id_gudang" binding:"required"`
	Alasan             string                        `json:"alasan" binding:"required"` // rusak, tidak_sesuai, cacat, dll
	MetodePengembalian string                        `json:"metode_pengembalian" binding:"required,oneof=potong_hutang refund tukar_barang"`
	Keterangan         string                        `json:"keterangan"`
	Items              []CreateReturPembelianItemReq `json:"items" binding:"required,min=1,dive"`
}

// CreateReturPembelianItemReq adalah DTO untuk item retur ke supplier
type CreateReturPembelianItemReq struct {
	IDProduk    uint    `json:"id_produk" binding:"required"`
	Jumlah      int     `json:"jumlah" binding:"required,min=1"`
	HargaSatuan float64 `json:"harga_satuan" binding:"required,gt=0"` // Harga saat pembelian
}

// ReturPembelianResponse adalah DTO untuk response retur pembelian
type ReturPembelianResponse struct {
	ID                 uint                         `json:"id"`
	NomorRetur         string                       `json:"nomor_retur"`
	IDPemasok          uint                         `json:"id_pemasok"`
	NamaPemasok        string                       `json:"nama_pemasok"`
	Alasan             string                       `json:"alasan"`
	Subtotal           float64                      `json:"subtotal"`
	Total              float64                      `json:"total"`
	MetodePengembalian string                       `json:"metode_pengembalian"`
	Status             string                       `json:"status"`
	Keterangan         string                       `json:"keterangan"`
	NamaPembuat        string                       `json:"nama_pembuat"`
	DibuatPada         time.Time                    `json:"dibuat_pada"`
	Items              []ReturPembelianItemResponse `json:"items,omitempty"`
}

// ReturPembelianItemResponse adalah DTO untuk item retur pembelian
type ReturPembelianItemResponse struct {
	ID          uint    `json:"id"`
	IDProduk    uint    `json:"id_produk"`
	SKUProduk   string  `json:"sku_produk"`
	NamaProduk  string  `json:"nama_produk"`
	Jumlah      int     `json:"jumlah"`
	HargaSatuan float64 `json:"harga_satuan"`
	Subtotal    float64 `json:"subtotal"`
}

// ListReturPembelianRequest adalah DTO untuk filter list retur pembelian
type ListReturPembelianRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
	IDPemasok *uint  `form:"id_pemasok"`
	Status    string `form:"status" binding:"omitempty,oneof=pending approved completed rejected"`
}
