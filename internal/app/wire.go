//go:build wireinject

package app

import (
	"github.com/google/wire"

	bizauth "auth_info/internal/biz/auth"
	bizdict "auth_info/internal/biz/dict"
	bizdoc "auth_info/internal/biz/document"
	bizhello "auth_info/internal/biz/hello"
	"auth_info/internal/config"
	"auth_info/internal/data"
	"auth_info/internal/handler"
	"auth_info/internal/service"
)

func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		data.NewDB,
		data.NewEnforcer,
		data.NewUserRepository,
		data.NewDictRepository,
		bizhello.NewUseCase,
		bizauth.NewUseCase,
		bizdict.NewUseCase,
		bizdoc.NewUseCase,
		handler.NewHelloHandler,
		handler.NewAuthHandler,
		handler.NewDictHandler,
		handler.NewDocumentHandler,
		service.NewHelloService,
		NewApp,
	)
	return nil, nil
}
