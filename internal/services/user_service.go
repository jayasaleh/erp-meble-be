package services

import (
	"math"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/utils"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService adalah interface untuk user service
type UserService interface {
	GetUserByID(id uint) (*dto.UserResponse, error)
	ListUsers(req dto.ListUsersRequest) (*dto.ListUsersResponse, error)
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uint) error
	ChangePassword(userID uint, req dto.ChangePasswordRequest) error
}

type userService struct {
	userRepo repositories.UserRepository
	logger   *zap.Logger
}

// NewUserService membuat instance UserService baru
func NewUserService() UserService {
	return &userService{
		userRepo: repositories.NewUserRepository(),
		logger:   utils.GetLogger(),
	}
}

func (s *userService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Warn("User not found",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		return nil, utils.ErrUserNotFound
	}

	return s.toUserResponse(user), nil
}

func (s *userService) ListUsers(req dto.ListUsersRequest) (*dto.ListUsersResponse, error) {
	// Set default values
	page := req.Page
	if page < 1 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Get users
	users, total, err := s.userRepo.FindAll(page, pageSize, req.Search, req.Peran, req.Aktif)
	if err != nil {
		s.logger.Error("Failed to list users",
			zap.Error(err),
		)
		return nil, utils.NewAppError(utils.ErrCodeDatabaseError, "Failed to list users", err)
	}

	// Convert to DTO
	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.toUserResponse(&user)
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.ListUsersResponse{
		Users: userResponses,
		Pagination: dto.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		s.logger.Warn("User already exists",
			zap.String("email", req.Email),
		)
		return nil, utils.ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, utils.NewAppError(utils.ErrCodeInternalError, "Failed to hash password", err)
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Nama:     req.Nama,
		Peran:    req.Peran,
		Aktif:    true, // Default aktif
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("Failed to create user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, utils.NewAppError(utils.ErrCodeDatabaseError, "Failed to create user", err)
	}

	s.logger.Info("User created successfully",
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
	)

	return s.toUserResponse(user), nil
}

func (s *userService) UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Warn("User not found",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		return nil, utils.ErrUserNotFound
	}

	// Update fields if provided
	if req.Nama != nil {
		user.Nama = *req.Nama
	}

	if req.Peran != nil {
		user.Peran = *req.Peran
	}

	if req.Aktif != nil {
		user.Aktif = *req.Aktif
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		s.logger.Error("Failed to update user",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		return nil, utils.NewAppError(utils.ErrCodeDatabaseError, "Failed to update user", err)
	}

	s.logger.Info("User updated successfully",
		zap.Uint("user_id", id),
	)

	return s.toUserResponse(user), nil
}

func (s *userService) DeleteUser(id uint) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		s.logger.Warn("User not found",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		return utils.ErrUserNotFound
	}

	// Prevent deleting owner
	if user.Peran == "owner" {
		s.logger.Warn("Cannot delete owner user",
			zap.Uint("user_id", id),
		)
		return utils.NewAppError(utils.ErrCodeValidationError, "Cannot delete owner user", nil)
	}

	// Delete user (soft delete)
	if err := s.userRepo.Delete(id); err != nil {
		s.logger.Error("Failed to delete user",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
		return utils.NewAppError(utils.ErrCodeDatabaseError, "Failed to delete user", err)
	}

	s.logger.Info("User deleted successfully",
		zap.Uint("user_id", id),
	)

	return nil
}

func (s *userService) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Warn("User not found",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return utils.ErrUserNotFound
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.PasswordLama)); err != nil {
		s.logger.Warn("Invalid old password",
			zap.Uint("user_id", userID),
		)
		return utils.NewAppError(utils.ErrCodeValidationError, "Invalid old password", nil)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PasswordBaru), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return utils.NewAppError(utils.ErrCodeInternalError, "Failed to hash password", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(userID, string(hashedPassword)); err != nil {
		s.logger.Error("Failed to update password",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return utils.NewAppError(utils.ErrCodeDatabaseError, "Failed to update password", err)
	}

	s.logger.Info("Password changed successfully",
		zap.Uint("user_id", userID),
	)

	return nil
}

// toUserResponse mengkonversi model User ke DTO UserResponse
func (s *userService) toUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Nama:  user.Nama,
		Peran: user.Peran,
		Aktif: user.Aktif,
	}
}
