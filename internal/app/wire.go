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
		// data 层
		data.NewDB,
		data.NewEnforcer,
		// biz 层
		bizhello.NewUseCase,
		bizauth.NewUseCase,
		bizdict.NewUseCase,
		bizdoc.NewUseCase,
		// handler 层
		handler.NewHelloHandler,
		handler.NewAuthHandler,
		handler.NewDictHandler,
		handler.NewDocumentHandler,
		// service 层（gRPC）
		service.NewHelloService,
		// app 装配
		NewApp,
	)
	return nil, nil
}
