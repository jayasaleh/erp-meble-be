package dto

import "time"

// CreatePemasokRequest adalah DTO untuk membuat pemasok baru
type CreatePemasokRequest struct {
	Nama    string `json:"nama" binding:"required"`
	Kontak  string `json:"kontak"`
	Telepon string `json:"telepon"`
	Email   string `json:"email" binding:"omitempty,email"`
	Alamat  string `json:"alamat"`
	Aktif   bool   `json:"aktif"` // default true usually, but allowed to set
}

// UpdatePemasokRequest adalah DTO untuk mengupdate pemasok
type UpdatePemasokRequest struct {
	Nama    *string `json:"nama"`
	Kontak  *string `json:"kontak"`
	Telepon *string `json:"telepon"`
	Email   *string `json:"email" binding:"omitempty,email"`
	Alamat  *string `json:"alamat"`
	Aktif   *bool   `json:"aktif"`
}

// PemasokResponse adalah DTO untuk response data pemasok
type PemasokResponse struct {
	ID             uint      `json:"id"`
	Nama           string    `json:"nama"`
	Kontak         string    `json:"kontak"`
	Telepon        string    `json:"telepon"`
	Email          string    `json:"email"`
	Alamat         string    `json:"alamat"`
	Aktif          bool      `json:"aktif"`
	DibuatPada     time.Time `json:"dibuat_pada"`
	DiperbaruiPada time.Time `json:"diperbarui_pada"`
}

// ListPemasokRequest adalah DTO untuk filter list pemasok
type ListPemasokRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search"` // Nama, Kontak, or Email
	Aktif  *bool  `form:"aktif"`
}

// ListPemasokResponse adalah DTO untuk response list pemasok
type ListPemasokResponse struct {
	Pemasok    []PemasokResponse `json:"suppliers"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}
