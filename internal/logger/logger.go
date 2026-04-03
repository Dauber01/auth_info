package logger

import (
	"auth_info/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// NewLogger Wire Provider，从 Config 构造 *zap.Logger，同时赋值全局变量供 cmd 工具继续使用
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if err := InitLogger(cfg.Log.Level); err != nil {
		return nil, err
	}
	return Logger, nil
}

func InitLogger(level string) error {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Logger = logger
	return nil
}

func GetLogger() *zap.Logger {
	if Logger == nil {
		Logger, _ = zap.NewProduction()
	}
	return Logger
}

func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}
