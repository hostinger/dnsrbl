package logger

import (
	"os"

	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type logger struct {
	logger *zap.Logger
}

func NewLoggerFromEnv() Logger {
	switch os.Getenv("LOG_CHANNEL") {
	case "ZAP_PRODUCTION":
		l, _ := zap.NewProduction() // nolint
		return &logger{
			logger: l,
		}
	case "ZAP_DEVELOPMENT":
		l, _ := zap.NewDevelopment() // nolint
		return &logger{
			logger: l,
		}
	default:
		l, _ := zap.NewDevelopment() // nolint
		return &logger{
			logger: l,
		}
	}
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}
