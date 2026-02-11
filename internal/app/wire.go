// +build wireinject

package app

import (
	"github.com/google/wire"
	"auth_info/internal/config"
)

func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(NewApp)
	return nil, nil
}
