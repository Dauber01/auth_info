package service

import (
	"context"

	"google.golang.org/grpc"

	proto "auth_info/api/gen/api/proto"
	bizhello "auth_info/internal/biz/hello"
)

// HelloService 实现 gRPC HelloServiceServer 接口
type HelloService struct {
	proto.UnimplementedHelloServiceServer
	uc *bizhello.UseCase
}

// NewHelloService Wire Provider
func NewHelloService(uc *bizhello.UseCase) *HelloService {
	return &HelloService{uc: uc}
}

// SayHello 实现 gRPC SayHello 方法
func (s *HelloService) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	msg := s.uc.SayHello(req.Name)
	return &proto.HelloReply{
		Message: msg,
		Code:    0,
	}, nil
}

// RegisterGRPCServices 注册所有 gRPC 服务
func RegisterGRPCServices(server *grpc.Server, svc *HelloService) {
	proto.RegisterHelloServiceServer(server, svc)
}
