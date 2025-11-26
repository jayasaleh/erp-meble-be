package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger menginisialisasi logger
func InitLogger(env string) error {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	// Replace global logger
	zap.ReplaceGlobals(Logger)

	return nil
}

// GetLogger mengembalikan logger instance
func GetLogger() *zap.Logger {
	if Logger == nil {
		// Fallback ke development logger jika belum diinisialisasi
		logger, _ := zap.NewDevelopment()
		return logger
	}
	return Logger
}

