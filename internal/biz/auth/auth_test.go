package auth

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"auth_info/internal/apperr"
	"auth_info/internal/config"
	"auth_info/internal/data"
)

// stubUserRepo 是 UserRepository 的最小 stub，用于 biz 层单元测试
type stubUserRepo struct {
	user *data.User
	err  error
}

func (s *stubUserRepo) GetByUsername(_ context.Context, _ string) (*data.User, error) {
	return s.user, s.err
}

func (s *stubUserRepo) Create(_ context.Context, _ *data.User) error {
	return s.err
}

func newTestUseCase(repo UserRepository) *UseCase {
	return NewUseCase(repo, &config.Config{
		JWT: config.JWTConfig{Secret: "test-secret", Expire: 1},
	}, zap.NewNop())
}

func TestRegister_Conflict(t *testing.T) {
	uc := newTestUseCase(&stubUserRepo{user: &data.User{Username: "alice"}})
	err := uc.Register(context.Background(), "alice", "password")
	if !apperr.IsCode(err, apperr.CodeConflict) {
		t.Fatalf("expected CodeConflict, got %v", err)
	}
}

func TestRegister_Success(t *testing.T) {
	uc := newTestUseCase(&stubUserRepo{user: nil})
	err := uc.Register(context.Background(), "bob", "password123")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	uc := newTestUseCase(&stubUserRepo{user: nil})
	_, err := uc.Login(context.Background(), "nobody", "pass")
	if !apperr.IsCode(err, apperr.CodeUnauthenticated) {
		t.Fatalf("expected CodeUnauthenticated, got %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	// 用一个不是 bcrypt 哈希的密码，让 bcrypt.CompareHashAndPassword 必定失败
	uc := newTestUseCase(&stubUserRepo{user: &data.User{Username: "alice", Password: "not-a-hash"}})
	_, err := uc.Login(context.Background(), "alice", "wrong")
	if !apperr.IsCode(err, apperr.CodeUnauthenticated) {
		t.Fatalf("expected CodeUnauthenticated, got %v", err)
	}
}
