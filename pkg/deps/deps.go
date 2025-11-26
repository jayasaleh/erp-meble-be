package deps

// File ini digunakan untuk memastikan semua dependensi penting terdaftar di go.mod
// Dependensi akan muncul di go.mod setelah file ini di-import atau digunakan

import (
	// Logging
	_ "go.uber.org/zap"

	// Rate Limiting
	_ "github.com/ulule/limiter/v3"
	_ "github.com/ulule/limiter/v3/drivers/store/memory"

	// Swagger/API Documentation
	_ "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/files"
	// Note: github.com/swaggo/swag/cmd/swag is a CLI tool, not importable

	// Testing
	_ "github.com/stretchr/testify"
)

