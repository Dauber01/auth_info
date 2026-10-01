package app

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	bizauth "auth_info/internal/biz/auth"
	bizdoc "auth_info/internal/biz/document"
	"auth_info/internal/config"
	"auth_info/internal/data"
	datadoc "auth_info/internal/data/document"
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

// ProvideDocumentResources registers local template handles and idle HTTP connections.
func ProvideDocumentResources(cfg *config.Config, resources *Lifecycle) (bizdoc.Resources, error) {
	r, err := datadoc.NewResources(cfg.Document)
	if err != nil {
		return nil, err
	}
	resources.Add("document", r.Close)
	return r, nil
}

// ProvideAuthOptions isolates business token options from deployment configuration.
func ProvideAuthOptions(cfg *config.Config) bizauth.Options {
	return bizauth.Options{Secret: cfg.JWT.Secret, Expire: time.Duration(cfg.JWT.Expire) * time.Hour}
}
