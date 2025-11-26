package handlers

import (
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService services.AuthService
	logger      *zap.Logger
}

// NewAuthHandler membuat instance AuthHandler baru
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
		logger:      utils.GetLogger(),
	}
}

// Login menangani request login
// @Summary Login user
// @Description Login dengan email dan password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request",
			zap.Error(err),
		)
		utils.BadRequest(c, "Invalid request", err.Error())
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Login successful", response)
}

// Register menangani request register
// @Summary Register new user
// @Description Register user baru
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} utils.Response
// @Failure 409 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid register request",
			zap.Error(err),
		)
		utils.BadRequest(c, "Invalid request", err.Error())
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Created(c, "User created successfully", user)
}

// handleError menangani error dengan mapping ke response yang sesuai
func (h *AuthHandler) handleError(c *gin.Context, err error) {
	appErr := utils.GetAppError(err)

	switch appErr.Code {
	case utils.ErrCodeUnauthorized:
		utils.Unauthorized(c, appErr.Message)
	case utils.ErrCodeNotFound:
		utils.NotFound(c, appErr.Message)
	case utils.ErrCodeConflict:
		utils.Conflict(c, appErr.Message)
	case utils.ErrCodeValidationError, utils.ErrCodeInvalidInput:
		utils.BadRequest(c, appErr.Message, appErr.Error())
	default:
		h.logger.Error("Internal server error",
			zap.String("code", appErr.Code),
			zap.String("message", appErr.Message),
			zap.Error(appErr.Err),
		)
		utils.InternalServerError(c, "Internal server error", nil)
	}
}

