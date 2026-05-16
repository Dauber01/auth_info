package dict

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	bizdict "auth_info/internal/biz/dict"
	"auth_info/internal/middleware"
)

type dictRepoStub struct {
	types []bizdict.DictType
	err   error
}

func (s *dictRepoStub) ListDictTypes(_ context.Context) ([]bizdict.DictType, error) {
	return s.types, s.err
}
func (s *dictRepoStub) GetDictTypeByCode(_ context.Context, _ string) (*bizdict.DictType, error) {
	return nil, s.err
}
func (s *dictRepoStub) CreateDictType(_ context.Context, _ *bizdict.DictType) error { return s.err }
func (s *dictRepoStub) UpdateDictType(_ context.Context, _ uint, _, _ string, _ int) (bool, error) {
	return false, s.err
}
func (s *dictRepoStub) DeleteDictType(_ context.Context, _ uint) (bool, error) { return false, s.err }
func (s *dictRepoStub) ListDictItems(_ context.Context, _ string) ([]bizdict.DictItem, error) {
	return nil, s.err
}
func (s *dictRepoStub) CreateDictItem(_ context.Context, _ *bizdict.DictItem) error { return s.err }
func (s *dictRepoStub) UpdateDictItem(_ context.Context, _ uint, _, _, _ string, _ int, _ int) (bool, error) {
	return false, s.err
}
func (s *dictRepoStub) DeleteDictItem(_ context.Context, _ uint) (bool, error) { return false, s.err }

func newHandlerForTest(repo bizdict.DictRepository) *Handler {
	uc := bizdict.NewUseCase(repo, zap.NewNop())
	return NewHandler(uc)
}

func TestListDictTypes_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newHandlerForTest(&dictRepoStub{types: []bizdict.DictType{{Code: "gender", Name: "Gender"}}})

	engine := gin.New()
	engine.Use(middleware.ErrorHandler(zap.NewNop()))
	engine.GET("/dict/types", h.ListDictTypes)

	req := httptest.NewRequest(http.MethodGet, "/dict/types", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCreateDictType_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newHandlerForTest(&dictRepoStub{})

	engine := gin.New()
	engine.Use(middleware.ErrorHandler(zap.NewNop()))
	engine.POST("/dict/types", h.CreateDictType)

	req := httptest.NewRequest(http.MethodPost, "/dict/types", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
