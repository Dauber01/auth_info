package dict

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"auth_info/internal/apperr"
	"auth_info/internal/data"
)

type stubDictRepo struct {
	dictType *data.DictType
	deleted  bool
	err      error
}

func (s *stubDictRepo) ListDictTypes(context.Context) ([]data.DictType, error) { return nil, s.err }
func (s *stubDictRepo) GetDictTypeByCode(context.Context, string) (*data.DictType, error) {
	return s.dictType, s.err
}
func (s *stubDictRepo) CreateDictType(context.Context, *data.DictType) error { return s.err }
func (s *stubDictRepo) UpdateDictType(context.Context, uint, string, string, int) (bool, error) {
	return false, s.err
}
func (s *stubDictRepo) DeleteDictType(context.Context, uint) (bool, error) { return s.deleted, s.err }
func (s *stubDictRepo) ListDictItems(context.Context, string) ([]data.DictItem, error) { return nil, s.err }
func (s *stubDictRepo) CreateDictItem(context.Context, *data.DictItem) error { return s.err }
func (s *stubDictRepo) UpdateDictItem(context.Context, uint, string, string, string, int, int) (bool, error) {
	return false, s.err
}
func (s *stubDictRepo) DeleteDictItem(context.Context, uint) (bool, error) { return false, s.err }

func TestCreateDictType_Conflict(t *testing.T) {
	uc := NewUseCase(&stubDictRepo{dictType: &data.DictType{Code: "gender"}}, zap.NewNop())
	err := uc.CreateDictType(context.Background(), "gender", "Gender", "", 1)
	if !apperr.IsCode(err, apperr.CodeConflict) {
		t.Fatalf("expected CodeConflict, got %v", err)
	}
}

func TestDeleteDictType_NotFound(t *testing.T) {
	uc := NewUseCase(&stubDictRepo{deleted: false}, zap.NewNop())
	err := uc.DeleteDictType(context.Background(), 1)
	if !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected CodeNotFound, got %v", err)
	}
}
