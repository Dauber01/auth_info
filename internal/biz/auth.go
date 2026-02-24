package biz

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"auth_info/internal/config"
	"auth_info/internal/logger"
	"auth_info/internal/model"
)

// Claims JWT 自定义声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthUseCase 鉴权业务逻辑
type AuthUseCase struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthUseCase Wire Provider
func NewAuthUseCase(db *gorm.DB, cfg *config.Config) *AuthUseCase {
	return &AuthUseCase{db: db, cfg: cfg}
}

// Register 注册新用户（bcrypt 加密密码）
func (uc *AuthUseCase) Register(username, password string) error {
	var existing model.User
	if err := uc.db.Where("username = ?", username).First(&existing).Error; err == nil {
		return errors.New("username already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Username: username,
		Password: string(hash),
		Role:     "user",
	}
	if err = uc.db.Create(&user).Error; err != nil {
		return err
	}

	logger.GetLogger().Info("user registered", zap.String("username", username))
	return nil
}

// Login 验证用户名密码，成功后返回 JWT Token
func (uc *AuthUseCase) Login(username, password string) (string, error) {
	var user model.User
	if err := uc.db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := uc.generateToken(&user)
	if err != nil {
		return "", err
	}

	logger.GetLogger().Info("user logged in", zap.String("username", username))
	return token, nil
}

// ParseToken 解析并验证 JWT Token
func (uc *AuthUseCase) ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(uc.cfg.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (uc *AuthUseCase) generateToken(user *model.User) (string, error) {
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
