package utils

import (
	"errors"
	"fmt"
)

// AppError adalah custom error type untuk aplikasi
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Error codes
const (
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeDatabaseError    = "DATABASE_ERROR"
	ErrCodeValidationError  = "VALIDATION_ERROR"
	ErrCodeTokenExpired     = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid     = "TOKEN_INVALID"
)

// NewAppError membuat AppError baru
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors
var (
	ErrInvalidCredentials = NewAppError(ErrCodeUnauthorized, "Invalid credentials", nil)
	ErrUserNotFound       = NewAppError(ErrCodeNotFound, "User not found", nil)
	ErrUserExists         = NewAppError(ErrCodeConflict, "User already exists", nil)
	ErrTokenExpired       = NewAppError(ErrCodeTokenExpired, "Token expired", nil)
	ErrTokenInvalid       = NewAppError(ErrCodeTokenInvalid, "Invalid token", nil)
)

// IsAppError mengecek apakah error adalah AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError mengembalikan AppError dari error
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return NewAppError(ErrCodeInternalError, err.Error(), err)
}

