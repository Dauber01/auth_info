package service

import (
	"google.golang.org/grpc"

	apipb "auth_info/api/gen/api/proto"
	hellosvc "auth_info/internal/service/hello"
)

// RegisterGRPCServices 聚合各模块的 gRPC 服务注册。
func RegisterGRPCServices(server *grpc.Server, hello *hellosvc.Service) {
	apipb.RegisterHelloServiceServer(server, hello)
}
