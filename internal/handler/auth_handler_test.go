package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	bizauth "auth_info/internal/biz/auth"
	"auth_info/internal/config"
	"auth_info/internal/data"
	"auth_info/internal/middleware"
)

type authRepoStub struct {
	user *data.User
	err  error
}

func (s *authRepoStub) GetByUsername(_ context.Context, _ string) (*data.User, error) { return s.user, s.err }
func (s *authRepoStub) Create(_ context.Context, _ *data.User) error { return s.err }

func newAuthHandlerForTest(repo bizauth.UserRepository) *AuthHandler {
	uc := bizauth.NewUseCase(repo, &config.Config{JWT: config.JWTConfig{Secret: "test", Expire: 1}}, zap.NewNop())
	return NewAuthHandler(uc)
}

func TestRegister_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newAuthHandlerForTest(&authRepoStub{})

	engine := gin.New()
	engine.Use(middleware.ErrorHandler(zap.NewNop()))
	engine.POST("/register", h.Register)

	body := map[string]any{"username": ""}
	buf, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
