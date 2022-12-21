package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func New() *zap.Logger {
	logger := zap.NewExample()
	defer logger.Sync()

	return logger
}

