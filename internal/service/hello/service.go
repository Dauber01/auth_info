package hello

import (
	"context"

	apipb "auth_info/api/gen/api/proto"
	bizhello "auth_info/internal/biz/hello"
)

// Service 实现 gRPC HelloServiceServer 接口。
type Service struct {
	apipb.UnimplementedHelloServiceServer
	uc *bizhello.UseCase
}

// NewService Wire Provider
func NewService(uc *bizhello.UseCase) *Service {
	return &Service{uc: uc}
}

// SayHello 实现 gRPC SayHello 方法。
func (s *Service) SayHello(ctx context.Context, req *apipb.HelloRequest) (*apipb.HelloReply, error) {
	msg := s.uc.SayHello(ctx, req.GetName())
	return &apipb.HelloReply{
		Code:    0,
		Message: "success",
		Data: &apipb.HelloData{
			Message: msg,
		},
	}, nil
}
