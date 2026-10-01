package config

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

// Validate rejects unusable settings before any resources are opened.
func (cfg Config) Validate() error {
	s := cfg.Server
	if s.Port < 1 || s.Port > 65535 || s.GRPCPort < 1 || s.GRPCPort > 65535 || s.Port == s.GRPCPort {
		return fmt.Errorf("server ports must be distinct and between 1 and 65535")
	}
	if s.Mode != "debug" && s.Mode != "release" && s.Mode != "test" {
		return fmt.Errorf("invalid server mode")
	}
	for name, duration := range map[string]time.Duration{
		"read_header_timeout": s.ReadHeaderTimeout, "read_timeout": s.ReadTimeout,
		"write_timeout": s.WriteTimeout, "idle_timeout": s.IdleTimeout,
		"request_timeout": s.RequestTimeout,
		"grpc_timeout":    s.GRPCTimeout,
	} {
		if duration < 0 {
			return fmt.Errorf("server.%s cannot be negative", name)
		}
	}
	if s.ShutdownTimeout <= 0 {
		return fmt.Errorf("server.shutdown_timeout must be positive")
	}
	if s.WriteTimeout > 0 && s.RequestTimeout > 0 && s.WriteTimeout <= s.RequestTimeout {
		return fmt.Errorf("server.write_timeout must exceed request_timeout or be disabled")
	}
	if _, err := zapcore.ParseLevel(cfg.Log.Level); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	if cfg.Log.Format != "json" && cfg.Log.Format != "console" {
		return fmt.Errorf("invalid log format")
	}
	f := cfg.Log.File
	if f.Enabled && (strings.TrimSpace(f.Path) == "" || f.MaxSizeMB <= 0 || f.MaxBackups < 1 || f.MaxAgeDays < 1) {
		return fmt.Errorf("log file requires a path and positive size, backups and age")
	}
	if cfg.MySQL.Port < 1 || cfg.MySQL.Port > 65535 {
		return fmt.Errorf("invalid mysql port")
	}
	pool := cfg.MySQL.Pool
	if pool.MaxOpenConns < 0 || pool.MaxIdleConns < 0 || pool.ConnMaxLifetime < 0 || pool.ConnMaxIdleTime < 0 {
		return fmt.Errorf("mysql pool values cannot be negative")
	}
	if pool.MaxOpenConns > 0 && pool.MaxIdleConns > pool.MaxOpenConns {
		return fmt.Errorf("mysql idle connections exceed open connections")
	}
	if strings.TrimSpace(cfg.JWT.Secret) == "" || cfg.JWT.Expire <= 0 {
		return fmt.Errorf("jwt requires secret and positive expiry")
	}
	if cfg.Casbin.Model == "" {
		return fmt.Errorf("casbin model is required")
	}
	return nil
}
