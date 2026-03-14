package validation

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"auth_info/internal/apperr"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if msg, ok := req.(proto.Message); ok {
			if err := ValidateProto(msg); err != nil {
				return nil, status.Error(apperr.GRPCStatusCode(err), apperr.Message(err))
			}
		}
		return handler(ctx, req)
	}
}
