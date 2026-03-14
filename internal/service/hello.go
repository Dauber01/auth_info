package service

import (
	"context"

	"google.golang.org/grpc"

	apipb "auth_info/api/gen/api/proto"
	bizhello "auth_info/internal/biz/hello"
)

// HelloService 实现 gRPC HelloServiceServer 接口
type HelloService struct {
	apipb.UnimplementedHelloServiceServer
	uc *bizhello.UseCase
}

// NewHelloService Wire Provider
func NewHelloService(uc *bizhello.UseCase) *HelloService {
	return &HelloService{uc: uc}
}

// SayHello 实现 gRPC SayHello 方法
func (s *HelloService) SayHello(ctx context.Context, req *apipb.HelloRequest) (*apipb.HelloReply, error) {
	msg := s.uc.SayHello(ctx, req.GetName())
	return &apipb.HelloReply{
		Code:    0,
		Message: "success",
		Data: &apipb.HelloData{
			Message: msg,
		},
	}, nil
}

// RegisterGRPCServices 注册所有 gRPC 服务
func RegisterGRPCServices(server *grpc.Server, svc *HelloService) {
	apipb.RegisterHelloServiceServer(server, svc)
}
