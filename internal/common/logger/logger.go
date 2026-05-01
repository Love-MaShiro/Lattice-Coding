package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"lattice-coding/internal/common/config"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(cfg *config.Config) *Logger {
	var config zap.Config

	if cfg.App.Env == "development" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.Level = zap.NewAtomicLevelAt(getLevel(cfg.Logging.Level))
	config.Encoding = cfg.Logging.Format
	config.OutputPaths = []string{cfg.Logging.Output}

	logger, err := config.Build()
	if err != nil {
		zap.L().Fatal("failed to initialize logger", zap.Error(err))
		os.Exit(1)
	}

	return &Logger{logger}
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l *Logger) Sync() {
	_ = l.Logger.Sync()
}
