package handlers

import (
	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService services.UserService
	logger      *zap.Logger
}

// NewUserHandler membuat instance UserHandler baru
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
		logger:      utils.GetLogger(),
	}
}

// GetCurrentUser mendapatkan user yang sedang login
// @Summary Get current user
// @Description Mendapatkan informasi user yang sedang login
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	id, ok := userID.(uint)
	if !ok {
		// Try to convert from float64 (JSON number)
		if idFloat, ok := userID.(float64); ok {
			id = uint(idFloat)
		} else {
			h.logger.Error("Invalid user_id type in context",
				zap.Any("user_id", userID),
			)
			utils.Unauthorized(c, "Invalid user context")
			return
		}
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "User retrieved successfully", user)
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	appErr := utils.GetAppError(err)

	switch appErr.Code {
	case utils.ErrCodeUnauthorized:
		utils.Unauthorized(c, appErr.Message)
	case utils.ErrCodeNotFound:
		utils.NotFound(c, appErr.Message)
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

