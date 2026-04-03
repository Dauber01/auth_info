package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"auth_info/internal/apperr"
	"auth_info/internal/config"
	"auth_info/internal/data"
)

// Claims JWT 自定义声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// UseCase 鉴权业务逻辑
type UseCase struct {
	users  UserRepository
	cfg    *config.Config
	logger *zap.Logger
}

// NewUseCase Wire Provider
func NewUseCase(users UserRepository, cfg *config.Config, logger *zap.Logger) *UseCase {
	return &UseCase{users: users, cfg: cfg, logger: logger}
}

// Register 注册新用户（bcrypt 加密密码）
func (uc *UseCase) Register(ctx context.Context, username, password string) error {
	existing, err := uc.users.GetByUsername(ctx, username)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to query user", err)
	}
	if existing != nil {
		return apperr.New(apperr.CodeConflict, "username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to hash password", err)
	}

	user := data.User{
		Username: username,
		Password: string(hash),
		Role:     "user",
	}
	if err = uc.users.Create(ctx, &user); err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to create user", err)
	}

	uc.logger.Info("user registered", zap.String("username", username))
	return nil
}

// Login 验证用户名密码，成功后返回 JWT Token
func (uc *UseCase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.users.GetByUsername(ctx, username)
	if err != nil {
		return "", apperr.Wrap(apperr.CodeInternal, "failed to query user", err)
	}
	if user == nil {
		return "", apperr.New(apperr.CodeUnauthenticated, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", apperr.New(apperr.CodeUnauthenticated, "invalid credentials")
	}

	token, err := uc.generateToken(user)
	if err != nil {
		return "", apperr.Wrap(apperr.CodeInternal, "failed to generate token", err)
	}

	uc.logger.Info("user logged in", zap.String("username", username))
	return token, nil
}

// ParseToken 解析并验证 JWT Token
func (uc *UseCase) ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(uc.cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, apperr.New(apperr.CodeUnauthenticated, "invalid token")
	}
	return claims, nil
}

func (uc *UseCase) generateToken(user *data.User) (string, error) {
	expire := time.Duration(uc.cfg.JWT.Expire) * time.Hour
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.cfg.JWT.Secret))
}
