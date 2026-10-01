package app

import (
	"errors"

	"auth_info/internal/config"
)

// InitializeApp transfers owned resources to App only after all providers succeed.
func InitializeApp(cfg *config.Config) (*App, error) { return initializeWith(cfg, initializeApp) }

func initializeWith(cfg *config.Config, build func(*config.Config, *Lifecycle) (*App, error)) (*App, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	resources := &Lifecycle{}
	application, err := build(cfg, resources)
	if err != nil {
		return nil, errors.Join(err, resources.Close())
	}
	return application, nil
}
