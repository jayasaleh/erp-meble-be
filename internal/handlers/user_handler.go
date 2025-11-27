package handlers

import (
	"strconv"

	"real-erp-mebel/be/internal/dto"
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
// @Success 200 {object} utils.Response{data=dto.UserResponse}
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	id := h.getUserIDFromContext(userID)
	if id == 0 {
		utils.Unauthorized(c, "Invalid user context")
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "User retrieved successfully", user)
}

// ListUsers mendapatkan daftar semua users dengan pagination dan filter
// @Summary List users
// @Description Mendapatkan daftar semua users dengan pagination dan filter
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search by name or email"
// @Param peran query string false "Filter by role" Enums(owner, kasir, admin_gudang, finance)
// @Param aktif query bool false "Filter by active status"
// @Success 200 {object} utils.Response{data=dto.ListUsersResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req dto.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Invalid query parameters", err)
		return
	}

	// Check permission (only owner and admin can list all users)
	if !h.hasPermission(c, []string{"owner", "admin_gudang"}) {
		utils.Forbidden(c, "You don't have permission to list users")
		return
	}

	response, err := h.userService.ListUsers(req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Users retrieved successfully", response)
}

// GetUserByID mendapatkan user berdasarkan ID
// @Summary Get user by ID
// @Description Mendapatkan detail user berdasarkan ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response{data=dto.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	// Check permission (users can only view their own profile, owner/admin can view all)
	userIDValue, _ := c.Get("user_id")
	userID := h.getUserIDFromContext(userIDValue)
	if userID != uint(id) && !h.hasPermission(c, []string{"owner", "admin_gudang"}) {
		utils.Forbidden(c, "You don't have permission to view this user")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "User retrieved successfully", user)
}

// CreateUser membuat user baru
// @Summary Create user
// @Description Membuat user baru (hanya owner/admin)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Create user request"
// @Success 201 {object} utils.Response{data=dto.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Check permission (only owner and admin can create users)
	if !h.hasPermission(c, []string{"owner", "admin_gudang"}) {
		utils.Forbidden(c, "You don't have permission to create users")
		return
	}

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Created(c, "User created successfully", user)
}

// UpdateUser mengupdate user
// @Summary Update user
// @Description Mengupdate data user (users can update their own profile, owner/admin can update any)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "Update user request"
// @Success 200 {object} utils.Response{data=dto.UserResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	// Check permission
	userIDValue, _ := c.Get("user_id")
	userID := h.getUserIDFromContext(userIDValue)
	isOwnerOrAdmin := h.hasPermission(c, []string{"owner", "admin_gudang"})

	if userID != uint(id) && !isOwnerOrAdmin {
		utils.Forbidden(c, "You don't have permission to update this user")
		return
	}

	// Regular users can only update their own name, not role or active status
	if userID == uint(id) && !isOwnerOrAdmin {
		var req dto.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "Invalid request body", err)
			return
		}

		// Remove fields that regular users can't update
		req.Peran = nil
		req.Aktif = nil

		user, err := h.userService.UpdateUser(uint(id), req)
		if err != nil {
			h.handleError(c, err)
			return
		}

		utils.OK(c, "User updated successfully", user)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "User updated successfully", user)
}

// DeleteUser menghapus user (soft delete)
// @Summary Delete user
// @Description Menghapus user (hanya owner/admin, tidak bisa hapus owner)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Check permission (only owner and admin can delete users)
	if !h.hasPermission(c, []string{"owner", "admin_gudang"}) {
		utils.Forbidden(c, "You don't have permission to delete users")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID", err)
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "User deleted successfully", nil)
}

// ChangePassword mengubah password user
// @Summary Change password
// @Description Mengubah password user yang sedang login
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDValue, _ := c.Get("user_id")
	userID := h.getUserIDFromContext(userIDValue)
	if userID == 0 {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Jangan kirim error detail yang mungkin berisi password
		utils.BadRequest(c, "Invalid request body", nil)
		return
	}

	if err := h.userService.ChangePassword(userID, req); err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Password changed successfully", nil)
}

// Helper functions

// getUserIDFromContext mengkonversi user_id dari context ke uint
func (h *UserHandler) getUserIDFromContext(userID interface{}) uint {
	if userID == nil {
		return 0
	}

	switch v := userID.(type) {
	case uint:
		return v
	case float64:
		return uint(v)
	case int:
		return uint(v)
	case int64:
		return uint(v)
	default:
		h.logger.Error("Invalid user_id type in context",
			zap.Any("user_id", userID),
		)
		return 0
	}
}

// hasPermission mengecek apakah user memiliki permission berdasarkan role
func (h *UserHandler) hasPermission(c *gin.Context, allowedRoles []string) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}

	roleStr, ok := role.(string)
	if !ok {
		return false
	}

	for _, allowedRole := range allowedRoles {
		if roleStr == allowedRole {
			return true
		}
	}

	return false
}

// handleError menangani error dan mengembalikan response yang sesuai
func (h *UserHandler) handleError(c *gin.Context, err error) {
	appErr := utils.GetAppError(err)

	switch appErr.Code {
	case utils.ErrCodeUnauthorized:
		utils.Unauthorized(c, appErr.Message)
	case utils.ErrCodeNotFound:
		utils.NotFound(c, appErr.Message)
	case utils.ErrCodeValidationError, utils.ErrCodeInvalidInput:
		utils.BadRequest(c, appErr.Message, appErr.Err)
	case utils.ErrCodeForbidden:
		utils.Forbidden(c, appErr.Message)
	default:
		h.logger.Error("Internal server error",
			zap.String("code", appErr.Code),
			zap.String("message", appErr.Message),
			zap.Error(appErr.Err),
		)
		utils.InternalServerError(c, "Internal server error", nil)
	}
}
