package app

import (
	"fmt"
	"net"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	bizauth "auth_info/internal/biz/auth"
	"auth_info/internal/config"
	"auth_info/internal/handler"
	"auth_info/internal/logger"
	"auth_info/internal/middleware"
	"auth_info/internal/router"
	"auth_info/internal/service"
	"auth_info/internal/validation"
)

type App struct {
	engine     *gin.Engine
	grpcServer *grpc.Server
	config     *config.Config
	helloSvc   *service.HelloService
}

// NewApp Wire Provider
func NewApp(
	cfg *config.Config,
	authUC *bizauth.UseCase,
	enforcer *casbin.Enforcer,
	helloHandler *handler.HelloHandler,
	authHandler *handler.AuthHandler,
	helloSvc *service.HelloService,
	dictHandler *handler.DictHandler,
	documentHandler *handler.DocumentHandler,
) (*App, error) {
	if err := logger.InitLogger(cfg.Log.Level); err != nil {
		return nil, err
	}

	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler())

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api/v1")
	{
		// Public routes (no auth)
		router.RegisterAuthRoutes(api, authHandler)

		// Protected routes (JWT + Casbin)
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(authUC))
		protected.Use(middleware.CasbinAuth(enforcer))
		{
			router.RegisterHelloRoutes(protected, helloHandler)
			router.RegisterDictRoutes(protected, dictHandler)
			router.RegisterDocumentRoutes(protected, documentHandler)
		}
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(validation.UnaryServerInterceptor()))
	service.RegisterGRPCServices(grpcServer, helloSvc)

	return &App{
		engine:     engine,
		grpcServer: grpcServer,
		config:     cfg,
		helloSvc:   helloSvc,
	}, nil
}

func (a *App) Run() error {
	go func() {
		grpcAddr := fmt.Sprintf(":%d", a.config.Server.Port+1000)
		listener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.GetLogger().Error("Failed to listen for gRPC", zap.Error(err))
			return
		}
		logger.GetLogger().Info("gRPC server starting", zap.String("addr", grpcAddr))
		if err := a.grpcServer.Serve(listener); err != nil {
			logger.GetLogger().Error("gRPC server error", zap.Error(err))
		}
	}()

	addr := fmt.Sprintf(":%d", a.config.Server.Port)
	logger.GetLogger().Info("Server starting", zap.String("addr", addr))
	return a.engine.Run(addr)
}

func (a *App) Stop() error {
	logger.GetLogger().Info("Server stopping")
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}
	return logger.Sync()
}
