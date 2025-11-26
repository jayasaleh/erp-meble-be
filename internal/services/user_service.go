package services

import (
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/utils"

	"go.uber.org/zap"
)

// UserService adalah interface untuk user service
type UserService interface {
	GetUserByID(id uint) (*dto.UserResponse, error)
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

	return &dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

