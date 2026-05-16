package dict

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"auth_info/internal/apperr"
)

type stubDictRepo struct {
	dictType *DictType
	deleted  bool
	err      error
}

func (s *stubDictRepo) ListDictTypes(context.Context) ([]DictType, error) { return nil, s.err }
func (s *stubDictRepo) GetDictTypeByCode(context.Context, string) (*DictType, error) {
	return s.dictType, s.err
}
func (s *stubDictRepo) CreateDictType(context.Context, *DictType) error { return s.err }
func (s *stubDictRepo) UpdateDictType(context.Context, uint, string, string, int) (bool, error) {
	return false, s.err
}
func (s *stubDictRepo) DeleteDictType(context.Context, uint) (bool, error) { return s.deleted, s.err }
func (s *stubDictRepo) ListDictItems(context.Context, string) ([]DictItem, error) {
	return nil, s.err
}
func (s *stubDictRepo) CreateDictItem(context.Context, *DictItem) error { return s.err }
func (s *stubDictRepo) UpdateDictItem(context.Context, uint, string, string, string, int, int) (bool, error) {
	return false, s.err
}
func (s *stubDictRepo) DeleteDictItem(context.Context, uint) (bool, error) { return false, s.err }

func TestCreateDictType_Conflict(t *testing.T) {
	uc := NewUseCase(&stubDictRepo{dictType: &DictType{Code: "gender"}}, zap.NewNop())
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
