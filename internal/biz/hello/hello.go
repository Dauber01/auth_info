package hello

import (
	"go.uber.org/zap"

	"auth_info/internal/logger"
)

// UseCase hello 业务逻辑
type UseCase struct{}

// NewUseCase Wire Provider
func NewUseCase() *UseCase {
	return &UseCase{}
}

// SayHello 处理 hello world 核心业务，name 为空时默认 "World"
func (uc *UseCase) SayHello(name string) string {
	if name == "" {
		name = "World"
	}
	msg := "Hello, " + name + "!"
	logger.GetLogger().Info("SayHello", zap.String("name", name), zap.String("msg", msg))
	return msg
}
