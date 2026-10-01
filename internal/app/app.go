// Package app owns application startup, shutdown and dependency assembly.
package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"auth_info/internal/config"
	"auth_info/internal/server"
)

// App coordinates servers and releases resources after server shutdown.
type App struct {
	http      *http.Server
	rpc       *server.GRPC
	log       *zap.Logger
	resources *Lifecycle
	timeout   time.Duration
	listen    func(string, string) (net.Listener, error)
	mu        sync.Mutex
	started   bool
	stopping  bool
	listeners []net.Listener
	serving   sync.WaitGroup
	stopOnce  sync.Once
	stopErr   error
	done      chan struct{}
}

// NewApp accepts already constructed transports, never business handlers or routers.
func NewApp(cfg *config.Config, log *zap.Logger, httpServer *http.Server, rpc *server.GRPC,
	resources *Lifecycle) (*App, error) {
	if httpServer == nil || rpc == nil || rpc.Server == nil || resources == nil || log == nil {
		return nil, fmt.Errorf("application dependencies are incomplete")
	}
	if cfg.Server.ShutdownTimeout <= 0 {
		return nil, fmt.Errorf("shutdown timeout must be positive")
	}
	return &App{
		http: httpServer, rpc: rpc, log: log, resources: resources, timeout: cfg.Server.ShutdownTimeout,
		listen: net.Listen, done: make(chan struct{}),
	}, nil
}

// Run acquires both listeners before starting either protocol and always cleans up on exit.
func (a *App) Run() error {
	a.mu.Lock()
	if a.stopping {
		a.mu.Unlock()
		return a.Stop()
	}
	if a.started {
		a.mu.Unlock()
		return fmt.Errorf("application already started")
	}
	a.started = true
	httpListener, err := a.listen("tcp", a.http.Addr)
	if err != nil {
		a.mu.Unlock()
		return errors.Join(fmt.Errorf("listen http: %w", err), a.Stop())
	}
	a.listeners = append(a.listeners, httpListener)
	rpcListener, err := a.listen("tcp", a.rpc.Addr)
	if err != nil {
		a.mu.Unlock()
		return errors.Join(fmt.Errorf("listen grpc: %w", err), a.Stop())
	}
	a.listeners = append(a.listeners, rpcListener)
	errCh := make(chan error, 2)
	a.serving.Add(2)
	go func() {
		defer a.serving.Done()
		if err := a.http.Serve(httpListener); !normalStop(err) {
			errCh <- fmt.Errorf("serve http: %w", err)
		}
	}()
	go func() {
		defer a.serving.Done()
		if err := a.rpc.Server.Serve(rpcListener); !normalStop(err) {
			errCh <- fmt.Errorf("serve grpc: %w", err)
		}
	}()
	a.log.Info("application started", zap.String("http", httpListener.Addr().String()),
		zap.String("grpc", rpcListener.Addr().String()))
	a.mu.Unlock()
	select {
	case err := <-errCh:
		return errors.Join(err, a.Stop())
	case <-a.done:
		return a.Stop()
	}
}

func normalStop(err error) bool {
	return err == nil || errors.Is(err, http.ErrServerClosed) || errors.Is(err, grpc.ErrServerStopped) ||
		errors.Is(err, net.ErrClosed)
}

// Stop is concurrent-safe and waits for shutdown and cleanup to finish on every call.
func (a *App) Stop() error {
	a.stopOnce.Do(func() {
		a.mu.Lock()
		a.stopping = true
		a.mu.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
		defer cancel()
		rpcDone := make(chan struct{})
		go func() { a.rpc.Server.GracefulStop(); close(rpcDone) }()
		if err := a.http.Shutdown(ctx); err != nil {
			a.stopErr = errors.Join(a.stopErr, fmt.Errorf("shutdown http: %w", err), a.http.Close())
		}
		select {
		case <-rpcDone:
		case <-ctx.Done():
			a.log.Warn("grpc shutdown deadline reached")
			// grpc can hold its internal mutex while GracefulStop waits for an
			// uncooperative handler. Stop closes transports, but may then block
			// behind that mutex; it must not extend our drain budget.
			go a.rpc.Server.Stop()
		}
		for _, listener := range a.listeners {
			if err := listener.Close(); !normalStop(err) {
				a.stopErr = errors.Join(a.stopErr, err)
			}
		}
		servingDone := make(chan struct{})
		go func() { a.serving.Wait(); close(servingDone) }()
		select {
		case <-servingDone:
		case <-ctx.Done():
			// Go cannot terminate handlers that ignore cancellation. Their
			// serving goroutines finish when those handlers eventually return.
		}
		a.stopErr = errors.Join(a.stopErr, a.resources.Close())
		close(a.done)
	})
	return a.stopErr
}
