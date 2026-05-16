package auth

import "context"

// UserRepository 用户数据访问接口（biz 定义，data 实现）。
// 输入输出均使用 biz 层 DTO，避免 biz 依赖持久化模型。
type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
}
