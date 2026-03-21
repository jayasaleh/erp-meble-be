package dto

// CreateGudangRequest adalah DTO untuk membuat gudang baru
type CreateGudangRequest struct {
	Kode      string `json:"kode" binding:"required"`
	Nama      string `json:"nama" binding:"required"`
	Alamat    string `json:"alamat"`
	Keterangan string `json:"keterangan"`
	Aktif     bool   `json:"aktif"`
}

// UpdateGudangRequest adalah DTO untuk mengupdate gudang
type UpdateGudangRequest struct {
	Kode      *string `json:"kode"`
	Nama      *string `json:"nama"`
	Alamat    *string `json:"alamat"`
	Keterangan *string `json:"keterangan"`
	Aktif     *bool   `json:"aktif"`
}

// GudangResponse adalah DTO untuk response gudang
type GudangResponse struct {
	ID         uint   `json:"id"`
	Kode       string `json:"kode"`
	Nama       string `json:"nama"`
	Alamat     string `json:"alamat"`
	Keterangan string `json:"keterangan"`
	Aktif      bool   `json:"aktif"`
}

// ListGudangRequest adalah DTO untuk filter list gudang
type ListGudangRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search"`
	Aktif  *bool  `form:"aktif"`
}

// ListGudangResponse adalah DTO untuk response list gudang
type ListGudangResponse struct {
	Gudangs    []GudangResponse `json:"gudangs"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}
