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
	dataauth "auth_info/internal/data/auth"
	datadict "auth_info/internal/data/dict"
	authhdl "auth_info/internal/handler/auth"
	dicthdl "auth_info/internal/handler/dict"
	dochdl "auth_info/internal/handler/document"
	hellohdl "auth_info/internal/handler/hello"
	"auth_info/internal/logger"
	"auth_info/internal/mcpserver"
	hellosvc "auth_info/internal/service/hello"
)

func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		logger.NewLogger,
		data.NewDB,
		data.NewEnforcer,
		dataauth.NewUserRepository,
		datadict.NewDictRepository,
		wire.Bind(new(bizauth.UserRepository), new(*dataauth.UserRepo)),
		wire.Bind(new(bizdict.DictRepository), new(*datadict.DictRepo)),
		bizhello.NewUseCase,
		bizauth.NewUseCase,
		bizdict.NewUseCase,
		bizdoc.NewUseCase,
		hellohdl.NewHandler,
		authhdl.NewHandler,
		dicthdl.NewHandler,
		dochdl.NewHandler,
		mcpserver.NewHelloMCPHandler,
		hellosvc.NewService,
		wire.Struct(new(AppDeps), "*"),
		NewApp,
	)
	return nil, nil
}
