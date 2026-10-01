// Package logger builds independently owned loggers; it has no global singleton.
package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"auth_info/internal/config"
	"auth_info/internal/pkg/trace"
)

// New creates console output and an optional rotating file, with idempotent cleanup.
func New(cfg config.LogConfig) (*zap.Logger, func() error, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, nil, fmt.Errorf("log level: %w", err)
	}
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	var encoder zapcore.Encoder
	switch cfg.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encCfg)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encCfg)
	default:
		return nil, nil, fmt.Errorf("unsupported log format")
	}
	// Stdout belongs to the process and may be a pipe or terminal, not a syncable file.
	console := zapcore.Lock(zapcore.AddSync(struct{ io.Writer }{os.Stdout}))
	writers := []zapcore.WriteSyncer{console}
	var file *lumberjack.Logger
	if cfg.File.Enabled {
		if cfg.File.Path == "" || cfg.File.MaxSizeMB <= 0 || cfg.File.MaxBackups <= 0 || cfg.File.MaxAgeDays <= 0 {
			return nil, nil, fmt.Errorf("invalid file log configuration")
		}
		if err := os.MkdirAll(filepath.Dir(cfg.File.Path), 0750); err != nil {
			return nil, nil, err
		}
		// Fail at construction rather than silently losing the first log entry.
		probe, err := os.OpenFile(cfg.File.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return nil, nil, fmt.Errorf("open log file: %w", err)
		}
		if err := probe.Close(); err != nil {
			return nil, nil, err
		}
		file = &lumberjack.Logger{
			Filename: cfg.File.Path, MaxSize: cfg.File.MaxSizeMB, MaxBackups: cfg.File.MaxBackups,
			MaxAge: cfg.File.MaxAgeDays, Compress: cfg.File.Compress,
		}
		writers = append(writers, zapcore.AddSync(file))
	}
	log := zap.New(zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), level), zap.AddCaller())
	var once sync.Once
	var closeErr error
	close := func() error {
		once.Do(func() {
			if err := log.Sync(); err != nil {
				closeErr = err
			}
			if file != nil {
				closeErr = errors.Join(closeErr, file.Close())
			}
		})
		return closeErr
	}
	return log, close, nil
}

// WithContext attaches correlation to the injected logger without shared mutation.
func WithContext(log *zap.Logger, ctx context.Context) *zap.Logger {
	if log == nil {
		log = zap.NewNop()
	}
	if id := trace.ID(ctx); id != "" {
		return log.With(zap.String("trace_id", id))
	}
	return log
}
