package dto

import "time"

// CreateProductRequest adalah DTO untuk membuat produk baru
type CreateProductRequest struct {
	SKU         string  `json:"sku" binding:"required"`
	Barcode     *string `json:"barcode"`
	Nama        string  `json:"nama" binding:"required"`
	Kategori    string  `json:"kategori"`
	Merek       string  `json:"merek"`
	IDPemasok   *uint   `json:"id_pemasok"`
	HargaModal  float64 `json:"harga_modal" binding:"required,min=0"`
	HargaJual   float64 `json:"harga_jual" binding:"required,min=0"`
	StokMinimum int     `json:"stok_minimum" binding:"min=0"`
	IzinDiskon  bool    `json:"izin_diskon"`
	Aktif       bool    `json:"aktif"`
}

// UpdateProductRequest adalah DTO untuk mengupdate produk
type UpdateProductRequest struct {
	SKU         *string  `json:"sku"`
	Barcode     *string  `json:"barcode"`
	Nama        *string  `json:"nama"`
	Kategori    *string  `json:"kategori"`
	Merek       *string  `json:"merek"`
	IDPemasok   *uint    `json:"id_pemasok"`
	HargaModal  *float64 `json:"harga_modal" binding:"omitempty,min=0"`
	HargaJual   *float64 `json:"harga_jual" binding:"omitempty,min=0"`
	StokMinimum *int     `json:"stok_minimum" binding:"omitempty,min=0"`
	IzinDiskon  *bool    `json:"izin_diskon"`
	Aktif       *bool    `json:"aktif"`
}

// ProductResponse adalah DTO untuk response produk
type ProductResponse struct {
	ID             uint                   `json:"id"`
	SKU            string                 `json:"sku"`
	Barcode        *string                `json:"barcode"`
	Nama           string                 `json:"nama"`
	Kategori       string                 `json:"kategori"`
	Merek          string                 `json:"merek"`
	IDPemasok      *uint                  `json:"id_pemasok"`
	NamaPemasok    *string                `json:"nama_pemasok,omitempty"`
	HargaModal     float64                `json:"harga_modal"`
	HargaJual      float64                `json:"harga_jual"`
	StokMinimum    int                    `json:"stok_minimum"`
	IzinDiskon     bool                   `json:"izin_diskon"`
	Aktif          bool                   `json:"aktif"`
	JumlahStok     int                    `json:"jumlah_stok"`
	DibuatOleh     uint                   `json:"dibuat_oleh"`
	NamaPembuat    *string                `json:"nama_pembuat,omitempty"`
	DiupdateOleh   uint                   `json:"diupdate_oleh"`
	NamaPengupdate *string                `json:"nama_pengupdate,omitempty"`
	Images         []ProductImageResponse `json:"images,omitempty"`
	DibuatPada     time.Time              `json:"dibuat_pada"`
	DiperbaruiPada time.Time              `json:"diperbarui_pada"`
}

// ProductImageResponse adalah DTO untuk response gambar produk
type ProductImageResponse struct {
	ID          uint   `json:"id"`
	PathGambar  string `json:"path_gambar"`
	GambarUtama bool   `json:"gambar_utama"`
	Urutan      int    `json:"urutan"`
}

// ProductListRequest adalah DTO untuk filter list produk
type ProductListRequest struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	Limit      int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search     string `form:"search"`
	Kategori   string `form:"kategori"`
	Merek      string `form:"merek"`
	IDPemasok  *uint  `form:"id_pemasok"`
	Aktif      *bool  `form:"aktif"`
	StokRendah bool   `form:"stok_rendah"` // Filter produk dengan stok < stok minimum
}

// ProductListResponse adalah DTO untuk response list produk
type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}
