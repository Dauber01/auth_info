//go:build wireinject

package app

import (
	"github.com/google/wire"

	bizauth "auth_info/internal/biz/auth"
	bizdict "auth_info/internal/biz/dict"
	bizhello "auth_info/internal/biz/hello"
	"auth_info/internal/config"
	"auth_info/internal/data"
	dataauth "auth_info/internal/data/auth"
	datadict "auth_info/internal/data/dict"
	authhdl "auth_info/internal/handler/auth"
	dicthdl "auth_info/internal/handler/dict"
	hellohdl "auth_info/internal/handler/hello"
	"auth_info/internal/mcpserver"
	"auth_info/internal/server"
	hellosvc "auth_info/internal/service/hello"
)

func initializeApp(cfg *config.Config, resources *Lifecycle) (*App, error) {
	wire.Build(
		ProvideLogger,
		ProvideDB,
		data.NewEnforcer,
		dataauth.NewUserRepository,
		datadict.NewDictRepository,
		wire.Bind(new(bizauth.UserRepository), new(*dataauth.UserRepo)),
		wire.Bind(new(bizdict.DictRepository), new(*datadict.DictRepo)),
		bizhello.NewUseCase,
		ProvideAuthOptions,
		bizauth.NewUseCase,
		bizdict.NewUseCase,
		hellohdl.NewHandler,
		authhdl.NewHandler,
		dicthdl.NewHandler,
		mcpserver.NewHelloMCPHandler,
		hellosvc.NewService,
		wire.Struct(new(server.HTTPDeps), "*"),
		server.NewHTTPServer,
		server.NewGRPCServer,
		NewApp,
	)
	return nil, nil
}
