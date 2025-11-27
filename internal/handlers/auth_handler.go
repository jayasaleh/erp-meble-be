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

// Login godoc
// @Summary      Login user
// @Description  Autentikasi user dengan email dan password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request   body      dto.LoginRequest  true  "Login request"
// @Success      200       {object}  utils.Response{data=dto.LoginResponse}
// @Failure      400       {object}  utils.Response
// @Failure      401       {object}  utils.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Login successful", response)
}

// Register godoc
// @Summary      Register new user
// @Description  Mendaftarkan user baru dengan email, password, dan nama
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request   body      dto.RegisterRequest  true  "Register request"
// @Success      201       {object}  utils.Response{data=dto.UserResponse}
// @Failure      400       {object}  utils.Response
// @Failure      409       {object}  utils.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Created(c, "User registered successfully", user)
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	appErr := utils.GetAppError(err)

	switch appErr.Code {
	case utils.ErrCodeUnauthorized:
		utils.Unauthorized(c, appErr.Message)
	case utils.ErrCodeNotFound:
		utils.NotFound(c, appErr.Message)
	case utils.ErrCodeValidationError, utils.ErrCodeInvalidInput:
		utils.BadRequest(c, appErr.Message, appErr.Err)
	case utils.ErrCodeConflict:
		utils.Conflict(c, appErr.Message)
	default:
		h.logger.Error("Internal server error",
			zap.String("code", appErr.Code),
			zap.String("message", appErr.Message),
			zap.Error(appErr.Err),
		)
		utils.InternalServerError(c, "Internal server error", nil)
	}
}
