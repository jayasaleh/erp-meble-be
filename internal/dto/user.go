package dto

// UserResponse adalah DTO untuk response user
type UserResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Nama   string `json:"nama"`
	Peran  string `json:"peran"`
	Aktif  bool   `json:"aktif"`
}

// CreateUserRequest adalah DTO untuk request create user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Nama     string `json:"nama" binding:"required" example:"John Doe"`
	Peran    string `json:"peran" binding:"required,oneof=owner kasir admin_gudang finance" example:"kasir"`
}

// UpdateUserRequest adalah DTO untuk request update user
type UpdateUserRequest struct {
	Nama    *string `json:"nama,omitempty" example:"John Doe Updated"`
	Peran   *string `json:"peran,omitempty" binding:"omitempty,oneof=owner kasir admin_gudang finance" example:"admin_gudang"`
	Aktif   *bool   `json:"aktif,omitempty" example:"true"`
}

// ChangePasswordRequest adalah DTO untuk request change password
type ChangePasswordRequest struct {
	PasswordLama string `json:"password_lama" binding:"required" example:"oldpassword123"`
	PasswordBaru string `json:"password_baru" binding:"required,min=6" example:"newpassword123"`
}

// ListUsersRequest adalah DTO untuk request list users (query params)
type ListUsersRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`
	Search   string `form:"search" example:"john"`
	Peran    string `form:"peran" binding:"omitempty,oneof=owner kasir admin_gudang finance" example:"kasir"`
	Aktif    *bool  `form:"aktif" example:"true"`
}

// ListUsersResponse adalah DTO untuk response list users
type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Pagination Pagination     `json:"pagination"`
}

// Pagination adalah DTO untuk pagination
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

