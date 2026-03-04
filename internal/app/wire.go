//go:build wireinject

package app

import (
	"github.com/google/wire"

	"auth_info/internal/biz"
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
		biz.NewHelloUseCase,
		biz.NewAuthUseCase,
		biz.NewDictUseCase,
		biz.NewDocumentUseCase,
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
