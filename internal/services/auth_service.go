package services

import (
	"time"

	"real-erp-mebel/be/internal/config"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"real-erp-mebel/be/internal/repositories"
	"real-erp-mebel/be/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AuthService adalah interface untuk auth service
type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(req dto.RegisterRequest) (*dto.UserResponse, error)
	GenerateToken(user *models.User) (string, error)
}

type authService struct {
	userRepo repositories.UserRepository
	logger   *zap.Logger
}

// NewAuthService membuat instance AuthService baru
func NewAuthService() AuthService {
	return &authService{
		userRepo: repositories.NewUserRepository(),
		logger:   utils.GetLogger(),
	}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		s.logger.Warn("Login attempt with invalid email",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, utils.ErrInvalidCredentials
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Warn("Login attempt with invalid password",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, utils.ErrInvalidCredentials
	}

	// Generate token
	token, err := s.GenerateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate token",
			zap.Uint("user_id", user.ID),
			zap.Error(err),
		)
		return nil, utils.NewAppError(utils.ErrCodeInternalError, "Failed to generate token", err)
	}

	s.logger.Info("User logged in successfully",
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
	)

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}, nil
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		s.logger.Warn("Registration attempt with existing email",
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
		Name:     req.Name,
		Role:     "kasir", // Default role
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("Failed to create user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, utils.NewAppError(utils.ErrCodeDatabaseError, "Failed to create user", err)
	}

	s.logger.Info("User registered successfully",
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
	)

	return &dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

func (s *authService) GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

