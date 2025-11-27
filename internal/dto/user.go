package dto

// UserResponse adalah DTO untuk response user
// @Description User response dengan informasi lengkap
type UserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Email string `json:"email" example:"user@example.com"`
	Nama  string `json:"nama" example:"John Doe"`
	Peran string `json:"peran" example:"kasir" enums:"owner,kasir,admin_gudang,finance"`
	Aktif bool   `json:"aktif" example:"true"`
}

// CreateUserRequest adalah DTO untuk request create user
// @Description Request untuk membuat user baru
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Nama     string `json:"nama" binding:"required" example:"John Doe"`
	Peran    string `json:"peran" binding:"required,oneof=owner kasir admin_gudang finance" example:"kasir" enums:"owner,kasir,admin_gudang,finance"`
}

// UpdateUserRequest adalah DTO untuk request update user
// @Description Request untuk mengupdate user (semua field optional)
type UpdateUserRequest struct {
	Nama  *string `json:"nama,omitempty" example:"John Doe Updated"`
	Peran *string `json:"peran,omitempty" binding:"omitempty,oneof=owner kasir admin_gudang finance" example:"admin_gudang" enums:"owner,kasir,admin_gudang,finance"`
	Aktif *bool   `json:"aktif,omitempty" example:"true"`
}

// ChangePasswordRequest adalah DTO untuk request change password
// @Description Request untuk mengubah password
type ChangePasswordRequest struct {
	PasswordLama string `json:"password_lama" binding:"required" example:"oldpassword123"`
	PasswordBaru string `json:"password_baru" binding:"required,min=6" example:"newpassword123"`
}

// ListUsersRequest adalah DTO untuk request list users (query params)
// @Description Query parameters untuk list users
type ListUsersRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`
	Search   string `form:"search" example:"john"`
	Peran    string `form:"peran" binding:"omitempty,oneof=owner kasir admin_gudang finance" example:"kasir" enums:"owner,kasir,admin_gudang,finance"`
	Aktif    *bool  `form:"aktif" example:"true"`
}

// ListUsersResponse adalah DTO untuk response list users
// @Description Response untuk list users dengan pagination
type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination Pagination     `json:"pagination"`
}

// Pagination adalah DTO untuk pagination
// @Description Informasi pagination
type Pagination struct {
	Page       int   `json:"page" example:"1"`
	PageSize   int   `json:"page_size" example:"10"`
	Total      int64 `json:"total" example:"50"`
	TotalPages int   `json:"total_pages" example:"5"`
}
