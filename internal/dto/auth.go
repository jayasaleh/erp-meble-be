package dto

// LoginRequest adalah DTO untuk request login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RegisterRequest adalah DTO untuk request register
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
}

// LoginResponse adalah DTO untuk response login
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse dipindahkan ke dto/user.go untuk konsistensi

