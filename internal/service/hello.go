package service

import (
	"context"

	"google.golang.org/grpc"
	"go.uber.org/zap"
	"auth_info/internal/logger"
)

// HelloServiceServer 实现了 proto 生成的 HelloServiceServer 接口
// 这是一个占位符，实际的实现会在 proto 代码生成后完成
type HelloService struct {
	logger *zap.Logger
}

// NewHelloService 创建 HelloService 实例
func NewHelloService() *HelloService {
	return &HelloService{
		logger: logger.GetLogger(),
	}
}

// RegisterGRPCServices 注册所有 gRPC 服务
func RegisterGRPCServices(server *grpc.Server, svc *HelloService) {
	// 这个函数将在 proto 代码生成后被调用
	// 用于注册生成的 gRPC 服务
	// 目前是占位符
	logger.GetLogger().Info("gRPC services registered")
}

// SayHello 示例方法（占位符）
func (s *HelloService) SayHello(ctx context.Context, req interface{}) (interface{}, error) {
	s.logger.Info("SayHello called")
	return map[string]interface{}{
		"message": "Hello from gRPC",
		"code":    0,
	}, nil
}
