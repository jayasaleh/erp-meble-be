package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response adalah struktur response standar untuk semua API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse mengembalikan response sukses
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengembalikan response error
func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// BadRequest mengembalikan 400 Bad Request
func BadRequest(c *gin.Context, message string, err interface{}) {
	ErrorResponse(c, http.StatusBadRequest, message, err)
}

// Unauthorized mengembalikan 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

// Forbidden mengembalikan 403 Forbidden
func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, nil)
}

// NotFound mengembalikan 404 Not Found
func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

// Conflict mengembalikan 409 Conflict
func Conflict(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message, nil)
}

// InternalServerError mengembalikan 500 Internal Server Error
func InternalServerError(c *gin.Context, message string, err interface{}) {
	ErrorResponse(c, http.StatusInternalServerError, message, err)
}

// Created mengembalikan 201 Created
func Created(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

// OK mengembalikan 200 OK
func OK(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

