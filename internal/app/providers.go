package app

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	bizauth "auth_info/internal/biz/auth"
	"auth_info/internal/config"
	"auth_info/internal/data"
	"auth_info/internal/pkg/logger"
)

// ProvideLogger registers logging first so it remains available until all other resources close.
func ProvideLogger(cfg *config.Config, resources *Lifecycle) (*zap.Logger, error) {
	log, close, err := logger.New(cfg.Log)
	if err != nil {
		return nil, err
	}
	resources.Add("logger", close)
	return log, nil
}

// ProvideDB transfers the opened SQL pool into the application lifecycle.
func ProvideDB(cfg *config.Config, log *zap.Logger, resources *Lifecycle) (*gorm.DB, error) {
	db, err := data.NewDB(cfg, log)
	if err != nil {
		return nil, err
	}
	resources.Add("mysql", func() error { return data.CloseDB(db) })
	return db, nil
}

// ProvideAuthOptions isolates business token options from deployment configuration.
func ProvideAuthOptions(cfg *config.Config) bizauth.Options {
	return bizauth.Options{Secret: cfg.JWT.Secret, Expire: time.Duration(cfg.JWT.Expire) * time.Hour}
}
