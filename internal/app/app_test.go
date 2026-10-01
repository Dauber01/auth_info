package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apipb "auth_info/api/gen/api/proto"
	"auth_info/internal/config"
	"auth_info/internal/server"
)

func testApp(t *testing.T, resources *Lifecycle) *App {
	t.Helper()
	cfg := &config.Config{Server: config.ServerConfig{ShutdownTimeout: 100 * time.Millisecond}}
	app, err := NewApp(cfg, zap.NewNop(), &http.Server{Addr: "127.0.0.1:0"},
		&server.GRPC{Server: grpc.NewServer(), Addr: "127.0.0.1:0"}, resources)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = app.Stop() })
	return app
}

func TestStopIsReverseConcurrentAndIdempotent(t *testing.T) {
	resources := &Lifecycle{}
	var order []string
	failure := errors.New("close failed")
	resources.Add("logger", func() error { order = append(order, "logger"); return nil })
	resources.Add("database", func() error { order = append(order, "database"); return failure })
	app := testApp(t, resources)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := app.Stop(); !errors.Is(err, failure) {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if len(order) != 2 || order[0] != "database" || order[1] != "logger" {
		t.Fatal(order)
	}
	if err := app.Run(); !errors.Is(err, failure) {
		t.Fatal(err)
	}
}

func TestStartupFailureClosesFirstListenerAndResources(t *testing.T) {
	resources := &Lifecycle{}
	closed := false
	resources.Add("fixture", func() error { closed = true; return nil })
	app := testApp(t, resources)
	var first net.Listener
	calls := 0
	failure := errors.New("occupied")
	app.listen = func(network, address string) (net.Listener, error) {
		calls++
		if calls == 2 {
			return nil, failure
		}
		var err error
		first, err = net.Listen(network, address)
		return first, err
	}
	if err := app.Run(); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	if !closed {
		t.Fatal("resources leaked")
	}
	if _, err := first.Accept(); !errors.Is(err, net.ErrClosed) {
		t.Fatal(err)
	}
}

func TestRunAndStopWaitForServerCompletion(t *testing.T) {
	app := testApp(t, &Lifecycle{})
	started := make(chan struct{}, 2)
	app.listen = func(network, address string) (net.Listener, error) {
		l, err := net.Listen(network, address)
		started <- struct{}{}
		return l, err
	}
	result := make(chan error, 1)
	go func() { result <- app.Run() }()
	for i := 0; i < 2; i++ {
		select {
		case <-started:
		case <-time.After(3 * time.Second):
			t.Fatal("start timed out")
		}
	}
	if err := app.Stop(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-result:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("run leaked")
	}
}

func TestInitializeFailureReleasesAlreadyAcquiredResources(t *testing.T) {
	cfg, err := config.LoadConfig("../../config")
	if err != nil {
		t.Fatal(err)
	}
	failure := errors.New("provider failed")
	var closed []string
	_, err = initializeWith(cfg, func(_ *config.Config, l *Lifecycle) (*App, error) {
		l.Add("first", func() error { closed = append(closed, "first"); return nil })
		l.Add("second", func() error { closed = append(closed, "second"); return nil })
		return nil, failure
	})
	if !errors.Is(err, failure) || len(closed) != 2 || closed[0] != "second" {
		t.Fatalf("%v %v", closed, err)
	}
}

func TestHTTPShutdownDeadlineClosesActiveConnection(t *testing.T) {
	entered := make(chan struct{})
	exited := make(chan struct{})
	app := testApp(t, &Lifecycle{})
	app.http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-r.Context().Done()
		close(exited)
	})
	addresses := make(chan string, 2)
	app.listen = func(network, address string) (net.Listener, error) {
		l, err := net.Listen(network, address)
		if err == nil {
			addresses <- l.Addr().String()
		}
		return l, err
	}
	result := make(chan error, 1)
	go func() { result <- app.Run() }()
	nextAddress := func() string {
		select {
		case address := <-addresses:
			return address
		case err := <-result:
			t.Fatalf("server startup failed: %v", err)
		case <-time.After(3 * time.Second):
			t.Fatal("server startup timed out")
		}
		return ""
	}
	address := nextAddress()
	nextAddress()
	clientDone := make(chan struct{})
	go func() {
		defer close(clientDone)
		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get("http://" + address)
		if err == nil {
			resp.Body.Close()
		}
	}()
	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		t.Fatal("request not started")
	}
	if err := app.Stop(); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	select {
	case <-exited:
	case <-time.After(time.Second):
		t.Fatal("connection not canceled")
	}
	<-clientDone
	<-result
}

func TestGRPCShutdownDeadlineCancelsActiveRPC(t *testing.T) {
	app := testApp(t, &Lifecycle{})
	entered := make(chan struct{})
	exited := make(chan struct{})
	release := make(chan struct{})
	defer close(release)
	app.rpc.Server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "fixture.Slow", HandlerType: (*interface{})(nil),
		Methods: []grpc.MethodDesc{{MethodName: "Wait", Handler: func(_ any, ctx context.Context,
			decode func(any) error, _ grpc.UnaryServerInterceptor) (any, error) {
			if err := decode(new(apipb.HelloRequest)); err != nil {
				return nil, err
			}
			close(entered)
			<-ctx.Done()
			close(exited)
			<-release
			return nil, ctx.Err()
		}}},
	}, struct{}{})
	addresses := make(chan string, 2)
	app.listen = func(network, address string) (net.Listener, error) {
		l, err := net.Listen(network, address)
		if err == nil {
			addresses <- l.Addr().String()
		}
		return l, err
	}
	result := make(chan error, 1)
	go func() { result <- app.Run() }()
	nextAddress := func() string {
		select {
		case address := <-addresses:
			return address
		case err := <-result:
			t.Fatalf("startup: %v", err)
		case <-time.After(3 * time.Second):
			t.Fatal("startup timeout")
		}
		return ""
	}
	nextAddress()
	address := nextAddress()
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rpcDone := make(chan error, 1)
	go func() {
		rpcDone <- conn.Invoke(ctx, "/fixture.Slow/Wait", new(apipb.HelloRequest), new(apipb.HelloReply))
	}()
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal("RPC did not start")
	}
	stopped := make(chan error, 1)
	go func() { stopped <- app.Stop() }()
	select {
	case err := <-stopped:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("shutdown exceeded budget while handler ignored cancellation")
	}
	select {
	case <-exited:
	case <-ctx.Done():
		t.Fatal("RPC not canceled")
	}
	if err := <-rpcDone; err == nil {
		t.Fatal("interrupted RPC succeeded")
	}
	if err := <-result; err != nil {
		t.Fatal(err)
	}
}
