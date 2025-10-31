package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

// InitLogger initializes the Zap logger
func InitLogger(environment string) {
	var err error
	var zapConfig zap.Config

	if environment == "production" {
		// Production: JSON format, info level
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		// Development: Console format, debug level
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	Logger, err = zapConfig.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	Sugar = Logger.Sugar()
}

// CloseLogger flushes any buffered log entries
func CloseLogger() {
	if Logger != nil {
		Logger.Sync()
	}
}
