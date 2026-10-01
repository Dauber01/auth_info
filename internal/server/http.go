// Package server constructs protocol servers without owning application resources.
package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	bizauth "auth_info/internal/biz/auth"
	"auth_info/internal/config"
	authhdl "auth_info/internal/handler/auth"
	dicthdl "auth_info/internal/handler/dict"
	dochdl "auth_info/internal/handler/document"
	hellohdl "auth_info/internal/handler/hello"
	"auth_info/internal/middleware"
	authrouter "auth_info/internal/router/auth"
	dictrouter "auth_info/internal/router/dict"
	docrouter "auth_info/internal/router/document"
	hellorouter "auth_info/internal/router/hello"
)

// HTTPDeps contains protocol adapters and authentication dependencies.
type HTTPDeps struct {
	AuthUC          *bizauth.UseCase
	Enforcer        *casbin.Enforcer
	HelloHandler    *hellohdl.Handler
	AuthHandler     *authhdl.Handler
	HelloMCPHandler http.Handler
	DictHandler     *dicthdl.Handler
	DocumentHandler *dochdl.Handler
}

// NewHTTPServer preserves public/protected routing and keeps MCP outside API deadlines.
func NewHTTPServer(cfg *config.Config, log *zap.Logger, deps HTTPDeps) (*http.Server, error) {
	if deps.AuthUC == nil || deps.Enforcer == nil || deps.HelloHandler == nil || deps.AuthHandler == nil ||
		deps.HelloMCPHandler == nil || deps.DictHandler == nil || deps.DocumentHandler == nil {
		return nil, fmt.Errorf("http dependencies are incomplete")
	}
	gin.SetMode(cfg.Server.Mode)
	engine := gin.New()
	engine.Use(middleware.TraceID())
	if cfg.Log.Access {
		engine.Use(middleware.AccessLog(log))
	}
	engine.Use(middleware.Recovery(log), middleware.ErrorHandler(log))
	engine.Any("/mcp", gin.WrapH(deps.HelloMCPHandler))
	api := engine.Group("/api/v1")
	authrouter.Register(api.Group("", middleware.RequestTimeout(cfg.Server.RequestTimeout)), deps.AuthHandler)
	protected := api.Group("", middleware.JWTAuth(deps.AuthUC), middleware.CasbinAuth(deps.Enforcer))
	ordinary := protected.Group("", middleware.RequestTimeout(cfg.Server.RequestTimeout))
	hellorouter.Register(ordinary, deps.HelloHandler)
	dictrouter.Register(ordinary, deps.DictHandler)
	docrouter.Register(protected.Group("", middleware.RequestTimeout(cfg.Server.DocumentTimeout)), deps.DocumentHandler)

	// A server-wide WriteTimeout would terminate MCP streams. Apply it to API responses only.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/") && cfg.Server.WriteTimeout > 0 {
			duration := cfg.Server.WriteTimeout
			if strings.HasPrefix(r.URL.Path, "/api/v1/document/") && cfg.Server.DocumentTimeout > 0 {
				if budget := cfg.Server.DocumentTimeout + 5*time.Second; budget > duration {
					duration = budget
				}
			}
			controller := http.NewResponseController(w)
			err := controller.SetWriteDeadline(time.Now().Add(duration))
			if err != nil && !errors.Is(err, http.ErrNotSupported) {
				log.Warn("set response deadline failed", zap.Error(err))
			}
		}
		engine.ServeHTTP(w, r)
	})
	return &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Server.Port), Handler: handler,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout, ReadTimeout: cfg.Server.ReadTimeout,
		IdleTimeout: cfg.Server.IdleTimeout,
	}, nil
}
