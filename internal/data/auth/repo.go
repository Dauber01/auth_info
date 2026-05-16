package auth

import (
	"context"
	"errors"

	"gorm.io/gorm"

	bizauth "auth_info/internal/biz/auth"
)

// UserRepo 实现 bizauth.UserRepository。
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepository Wire Provider
func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*bizauth.User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return toBiz(&user), nil
}

func (r *UserRepo) Create(ctx context.Context, user *bizauth.User) error {
	model := fromBiz(user)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	// 回写自增字段到调用方持有的 DTO
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt
	return nil
}
