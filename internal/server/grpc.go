package server

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"auth_info/internal/config"
	"auth_info/internal/pkg/apperr"
	"auth_info/internal/pkg/logger"
	"auth_info/internal/pkg/trace"
	"auth_info/internal/service"
	hellosvc "auth_info/internal/service/hello"
	"auth_info/internal/validation"
)

// GRPC contains the configured server and its independently configurable address.
type GRPC struct {
	Server *grpc.Server
	Addr   string
}

// NewGRPCServer constructs the transport while reusing the existing Hello service.
func NewGRPCServer(cfg *config.Config, log *zap.Logger, hello *hellosvc.Service) *GRPC {
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		unaryBoundary(log, cfg.Server.GRPCTimeout), validation.UnaryServerInterceptor(),
	))
	service.RegisterGRPCServices(s, hello)
	return &GRPC{Server: s, Addr: fmt.Sprintf(":%d", cfg.Server.GRPCPort)}
}

func unaryBoundary(log *zap.Logger, timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (reply any, err error) {
		ids := metadata.ValueFromIncomingContext(ctx, "x-trace-id")
		id := ""
		if len(ids) > 0 {
			id = ids[0]
		}
		ctx = trace.WithID(ctx, trace.Normalize(id))
		// A missing server transport in unit tests is harmless; real RPCs carry this header.
		_ = grpc.SetHeader(ctx, metadata.Pairs("x-trace-id", trace.ID(ctx)))
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		requestLog := logger.WithContext(log, ctx)
		defer func() {
			if recover() != nil {
				requestLog.Error("grpc request panicked", zap.ByteString("stack", debug.Stack()))
				reply, err = nil, status.Error(codes.Internal, "internal server error")
				return
			}
			if err == nil {
				err = ctx.Err()
			}
			if err != nil {
				requestLog.Error("grpc request failed", zap.String("method", info.FullMethod), zap.Error(err))
				if _, ok := status.FromError(err); !ok {
					err = status.Error(apperr.GRPCStatusCode(err), apperr.Message(err))
				}
				reply = nil
			}
		}()
		return next(ctx, req)
	}
}
