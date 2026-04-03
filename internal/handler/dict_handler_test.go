package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	bizdict "auth_info/internal/biz/dict"
	"auth_info/internal/data"
	"auth_info/internal/middleware"
)

type dictRepoStub struct {
	types []data.DictType
	err   error
}

func (s *dictRepoStub) ListDictTypes(_ context.Context) ([]data.DictType, error) { return s.types, s.err }
func (s *dictRepoStub) GetDictTypeByCode(_ context.Context, _ string) (*data.DictType, error) {
	return nil, s.err
}
func (s *dictRepoStub) CreateDictType(_ context.Context, _ *data.DictType) error { return s.err }
func (s *dictRepoStub) UpdateDictType(_ context.Context, _ uint, _, _ string, _ int) (bool, error) {
	return false, s.err
}
func (s *dictRepoStub) DeleteDictType(_ context.Context, _ uint) (bool, error) { return false, s.err }
func (s *dictRepoStub) ListDictItems(_ context.Context, _ string) ([]data.DictItem, error) {
	return nil, s.err
}
func (s *dictRepoStub) CreateDictItem(_ context.Context, _ *data.DictItem) error { return s.err }
func (s *dictRepoStub) UpdateDictItem(_ context.Context, _ uint, _, _, _ string, _ int, _ int) (bool, error) {
	return false, s.err
}
func (s *dictRepoStub) DeleteDictItem(_ context.Context, _ uint) (bool, error) { return false, s.err }

func newDictHandlerForTest(repo bizdict.DictRepository) *DictHandler {
	uc := bizdict.NewUseCase(repo, zap.NewNop())
	return NewDictHandler(uc)
}

func TestListDictTypes_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := newDictHandlerForTest(&dictRepoStub{types: []data.DictType{{Code: "gender", Name: "Gender"}}})

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
	h := newDictHandlerForTest(&dictRepoStub{})

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
