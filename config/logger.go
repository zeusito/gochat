package config

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

// NewLogger Creates a new instance of Uber Zap
func NewLogger() *zap.SugaredLogger {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}

// CloseLogger Flushes any pending logs
func CloseLogger(logger *zap.SugaredLogger) {
	_ = logger.Sync()
}
