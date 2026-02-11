package app

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"go.uber.org/zap"

	"auth_info/internal/config"
	"auth_info/internal/handler"
	"auth_info/internal/logger"
	"auth_info/internal/middleware"
	"auth_info/internal/service"
)

type App struct {
	engine      *gin.Engine
	grpcServer  *grpc.Server
	config      *config.Config
	helloSvc    *service.HelloService
}

func NewApp(cfg *config.Config) (*App, error) {
	// 初始化日志
	if err := logger.InitLogger(cfg.Log.Level); err != nil {
		return nil, err
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	engine := gin.New()

	// 使用中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler())

	// 添加 Swagger 路由
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 初始化 handler
	helloHandler := handler.NewHelloHandler()

	// 路由注册
	api := engine.Group("/api/v1")
	{
		api.GET("/hello", helloHandler.Hello)
	}

	// 初始化 gRPC 服务
	grpcServer := grpc.NewServer()
	helloSvc := service.NewHelloService()

	// 注册 gRPC 服务
	service.RegisterGRPCServices(grpcServer, helloSvc)

	app := &App{
		engine:     engine,
		grpcServer: grpcServer,
		config:     cfg,
		helloSvc:   helloSvc,
	}

	return app, nil
}

func (a *App) Run() error {
	// 在 goroutine 中启动 gRPC 服务器
	go func() {
		grpcAddr := fmt.Sprintf(":%d", a.config.Server.Port+1000) // gRPC 端口 = HTTP 端口 + 1000
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

	// 启动 HTTP 服务器
	addr := fmt.Sprintf(":%d", a.config.Server.Port)
	logger.GetLogger().Info("Server starting", zap.String("addr", addr))
	return a.engine.Run(addr)
}

func (a *App) Stop() error {
	logger.GetLogger().Info("Server stopping")

	// 优雅关闭 gRPC 服务器
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

	return logger.Sync()
}
