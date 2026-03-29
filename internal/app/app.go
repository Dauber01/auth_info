package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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
	httpServer *http.Server
	config     *config.Config
	helloSvc   *service.HelloService
	stopCh     chan struct{}
	stopOnce   sync.Once
}

const grpcGracefulStopTimeout = 5 * time.Second

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
		stopCh:     make(chan struct{}),
	}, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 2)

	grpcAddr := fmt.Sprintf(":%d", a.config.Server.Port+1000)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen gRPC: %w", err)
	}

	httpAddr := fmt.Sprintf(":%d", a.config.Server.Port)
	a.httpServer = &http.Server{
		Addr:    httpAddr,
		Handler: a.engine,
	}

	go func() {
		logger.GetLogger().Info("gRPC server starting", zap.String("addr", grpcAddr))
		if serveErr := a.grpcServer.Serve(grpcListener); serveErr != nil && !errors.Is(serveErr, grpc.ErrServerStopped) {
			errCh <- fmt.Errorf("gRPC server error: %w", serveErr)
		}
	}()

	go func() {
		logger.GetLogger().Info("HTTP server starting", zap.String("addr", httpAddr))
		if serveErr := a.httpServer.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http server error: %w", serveErr)
		}
	}()

	select {
	case runErr := <-errCh:
		_ = a.Stop()
		return runErr
	case <-a.stopCh:
		return nil
	}
}

func (a *App) Stop() error {
	var stopErr error

	a.stopOnce.Do(func() {
		close(a.stopCh)

		if a.httpServer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := a.httpServer.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				stopErr = err
			}
		}

		if a.grpcServer != nil {
			grpcStopped := make(chan struct{})
			go func() {
				a.grpcServer.GracefulStop()
				close(grpcStopped)
			}()

			select {
			case <-grpcStopped:
			case <-time.After(grpcGracefulStopTimeout):
				logger.GetLogger().Warn("gRPC graceful stop timeout reached, forcing stop")
				a.grpcServer.Stop()
			}
		}
	})

	if err := logger.Sync(); stopErr == nil {
		stopErr = err
	}
	return stopErr
}
