package server

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/casbin/casbin/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	apipb "auth_info/api/gen/api/proto"
	bizauth "auth_info/internal/biz/auth"
	bizdict "auth_info/internal/biz/dict"
	bizhello "auth_info/internal/biz/hello"
	"auth_info/internal/config"
	authhdl "auth_info/internal/handler/auth"
	dicthdl "auth_info/internal/handler/dict"
	hellohdl "auth_info/internal/handler/hello"
	"auth_info/internal/pkg/trace"
	hellosvc "auth_info/internal/service/hello"
)

func TestHTTPRoutingAndMCPDeadlineIsolation(t *testing.T) {
	cfg, err := config.LoadConfig("../../config")
	if err != nil {
		t.Fatal(err)
	}
	cfg.Server.Mode = "test"
	log := zap.NewNop()
	auth := bizauth.NewUseCase(nil, bizauth.Options{Secret: "fixture", Expire: time.Hour}, log)
	enforcer, err := casbin.NewEnforcer("../../config/rbac_model.conf")
	if err != nil {
		t.Fatal(err)
	}
	deps := HTTPDeps{
		AuthUC: auth, Enforcer: enforcer, AuthHandler: authhdl.NewHandler(auth),
		HelloHandler: hellohdl.NewHandler(bizhello.NewUseCase(log)),
		DictHandler:  dicthdl.NewHandler(bizdict.NewUseCase(nil, log)),
		HelloMCPHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := r.Context().Deadline(); ok {
				t.Error("MCP inherited API deadline")
			}
			if trace.ID(r.Context()) == "" {
				t.Error("MCP missing trace")
			}
			w.WriteHeader(202)
			w.(http.Flusher).Flush()
		}),
	}
	httpServer, err := NewHTTPServer(cfg, log, deps)
	if err != nil {
		t.Fatal(err)
	}
	if httpServer.WriteTimeout != 0 || httpServer.ReadHeaderTimeout <= 0 {
		t.Fatal("streaming or header policy lost")
	}
	for path, want := range map[string]int{
		"/mcp": 202, "/api/v1/hello": 401, "/api/v1/dict/types": 401,
	} {
		w := &deadlineRecorder{ResponseRecorder: httptest.NewRecorder()}
		before := time.Now()
		httpServer.Handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != want {
			t.Fatalf("%s: %d %s", path, w.Code, w.Body.String())
		}
		if path == "/mcp" {
			if !w.deadline.IsZero() {
				t.Fatal("MCP inherited API write deadline")
			}
			continue
		}
		budget := cfg.Server.WriteTimeout
		if w.deadline.Before(before.Add(budget)) || w.deadline.After(time.Now().Add(budget)) {
			t.Fatalf("%s: unexpected write deadline %v", path, w.deadline)
		}
	}
	for _, path := range []string{"/api/v1/document/generate-pdf", "/api/v1/document/generate-word"} {
		w := httptest.NewRecorder()
		httpServer.Handler.ServeHTTP(w, httptest.NewRequest(http.MethodPost, path, strings.NewReader("{}")))
		if w.Code != http.StatusNotFound || w.Header().Get("Content-Disposition") != "" ||
			w.Header().Get(trace.Header) == "" {
			t.Fatalf("removed route %s: %d %s", path, w.Code, w.Body.String())
		}
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	httpServer.Handler.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("public auth route lost: %d", w.Code)
	}
}

type deadlineRecorder struct {
	*httptest.ResponseRecorder
	deadline time.Time
}

func (w *deadlineRecorder) SetWriteDeadline(deadline time.Time) error {
	w.deadline = deadline
	return nil
}

func TestGRPCHelloValidationAndTrace(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{GRPCTimeout: time.Second}}
	rpc := NewGRPCServer(cfg, zap.NewNop(), hellosvc.NewService(bizhello.NewUseCase(zap.NewNop())))
	rpc.Server.RegisterService(&grpc.ServiceDesc{
		ServiceName: "fixture.Validation", HandlerType: (*interface{})(nil),
		Methods: []grpc.MethodDesc{{MethodName: "Validate", Handler: func(_ any, ctx context.Context,
			decode func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
			req := new(apipb.RegisterRequest)
			if err := decode(req); err != nil {
				return nil, err
			}
			return interceptor(ctx, req, &grpc.UnaryServerInfo{FullMethod: "/fixture.Validation/Validate"},
				func(context.Context, any) (any, error) { return &apipb.OperationReply{Code: 200}, nil })
		}}},
	}, struct{}{})
	listener := bufconn.Listen(1024 * 1024)
	go func() { _ = rpc.Server.Serve(listener) }()
	t.Cleanup(func() { rpc.Server.Stop(); _ = listener.Close() })
	conn, err := grpc.NewClient("passthrough:///fixture", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return listener.DialContext(ctx) }))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	client := apipb.NewHelloServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-trace-id", "rpc-42"))
	var headers metadata.MD
	reply, err := client.SayHello(ctx, &apipb.HelloRequest{Name: "Architecture"}, grpc.Header(&headers))
	if err != nil {
		t.Fatal(err)
	}
	if reply.GetData().GetMessage() != "Hello, Architecture!" || reply.GetCode() != 0 ||
		strings.Join(headers.Get("x-trace-id"), "") != "rpc-42" {
		t.Fatal(reply, headers)
	}
	err = conn.Invoke(ctx, "/fixture.Validation/Validate",
		&apipb.RegisterRequest{Username: "a", Password: "short"}, new(apipb.OperationReply))
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("validation lost: %v", err)
	}
}

func TestGRPCBoundaryTimeoutAndPanic(t *testing.T) {
	interceptor := unaryBoundary(zap.NewNop(), time.Millisecond)
	info := &grpc.UnaryServerInfo{FullMethod: "/fixture"}
	_, err := interceptor(context.Background(), nil, info, func(ctx context.Context, _ any) (any, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	})
	if status.Code(err) != codes.DeadlineExceeded {
		t.Fatal(err)
	}
	_, err = interceptor(context.Background(), nil, info, func(context.Context, any) (any, error) { panic("fixture") })
	if status.Code(err) != codes.Internal {
		t.Fatal(err)
	}
}
