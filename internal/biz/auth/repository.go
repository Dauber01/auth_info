package auth

import (
	"context"

	"auth_info/internal/data"
)

// UserRepository 用户数据访问接口（由 biz/auth 层定义，data 层实现）
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*data.User, error)
	Create(ctx context.Context, user *data.User) error
}
