package biz

import (
	"go.uber.org/zap"

	"auth_info/internal/logger"
)

// HelloUseCase hello 业务逻辑
type HelloUseCase struct{}

// NewHelloUseCase Wire Provider
func NewHelloUseCase() *HelloUseCase {
	return &HelloUseCase{}
}

// SayHello 处理 hello world 核心业务，name 为空时默认 "World"
func (uc *HelloUseCase) SayHello(name string) string {
	if name == "" {
		name = "World"
	}
	msg := "Hello, " + name + "!"
	logger.GetLogger().Info("SayHello", zap.String("name", name), zap.String("msg", msg))
	return msg
}
