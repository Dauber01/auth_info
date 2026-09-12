// Command server hosts real API boundaries with an isolated in-memory user store.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"auth_info/internal/apperr"
	bizauth "auth_info/internal/biz/auth"
	bizhello "auth_info/internal/biz/hello"
	"auth_info/internal/config"
	authhandler "auth_info/internal/handler/auth"
	hellohandler "auth_info/internal/handler/hello"
	"auth_info/internal/middleware"
	authrouter "auth_info/internal/router/auth"
	hellorouter "auth_info/internal/router/hello"
)

type memoryUsers struct {
	mu    sync.RWMutex
	users map[string]bizauth.User
}

// GetByUsername returns a copy so requests do not share mutable user records.
func (r *memoryUsers) GetByUsername(ctx context.Context, username string) (*bizauth.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[username]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

// Create stores a fixture user and preserves username uniqueness.
func (r *memoryUsers) Create(ctx context.Context, user *bizauth.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.Username]; exists {
		return apperr.New(apperr.CodeConflict, "username already exists")
	}
	copy := *user
	copy.ID = uint(len(r.users) + 1)
	r.users[user.Username] = copy
	return nil
}

func run() error {
	gin.SetMode(gin.TestMode)
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("generate fixture key: %w", err)
	}
	logger := zap.NewNop()
	cfg := &config.Config{JWT: config.JWTConfig{Secret: hex.EncodeToString(key), Expire: 1}}
	users := &memoryUsers{users: make(map[string]bizauth.User)}
	auth := bizauth.NewUseCase(users, cfg, logger)
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf")
	if err != nil {
		return fmt.Errorf("load policy model: %w", err)
	}
	if _, err := enforcer.AddPolicy("user", "/api/v1/hello", "GET"); err != nil {
		return fmt.Errorf("add fixture policy: %w", err)
	}
	engine := gin.New()
	engine.Use(middleware.ErrorHandler(logger))
	api := engine.Group("/api/v1")
	authrouter.Register(api, authhandler.NewHandler(auth))
	protected := api.Group("")
	protected.Use(middleware.JWTAuth(auth), middleware.CasbinAuth(enforcer))
	hellorouter.Register(protected, hellohandler.NewHandler(bizhello.NewUseCase(logger)))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen fixture: %w", err)
	}
	defer func() {
		// Serve also closes this listener; repeated cleanup is expected.
		_ = listener.Close()
	}()
	server := &http.Server{Handler: engine, ReadHeaderTimeout: 5 * time.Second}
	errorsCh := make(chan error, 1)
	go func() { errorsCh <- server.Serve(listener) }()
	defer func() {
		// Ensure cleanup on setup failures as well as after graceful shutdown.
		_ = server.Close()
	}()
	if err := json.NewEncoder(os.Stdout).Encode(map[string]string{
		"base_url": "http://" + listener.Addr().String(),
	}); err != nil {
		return fmt.Errorf("report fixture address: %w", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case err := <-errorsCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve fixture: %w", err)
		}
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdown); err != nil {
			return fmt.Errorf("stop fixture: %w", err)
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		// This CLI boundary reports setup failures to the Python runner's stderr.
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
