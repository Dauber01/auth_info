package hello

import (
	"context"

	"go.uber.org/zap"
)

// UseCase hello 业务逻辑
type UseCase struct {
	logger *zap.Logger
}

// NewUseCase Wire Provider
func NewUseCase(logger *zap.Logger) *UseCase {
	return &UseCase{logger: logger}
}

// SayHello 处理 hello world 核心业务，name 为空时默认 "World"
func (uc *UseCase) SayHello(ctx context.Context, name string) string {
	_ = ctx
	if name == "" {
		name = "World"
	}
	msg := "Hello, " + name + "!"
	uc.logger.Info("SayHello", zap.String("name", name), zap.String("msg", msg))
	return msg
}
