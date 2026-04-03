package dict

import (
	"context"

	"go.uber.org/zap"

	"auth_info/internal/apperr"
	"auth_info/internal/data"
)

// UseCase 字典配置业务逻辑
type UseCase struct {
	repo   DictRepository
	logger *zap.Logger
}

// NewUseCase Wire Provider
func NewUseCase(repo DictRepository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

// ListDictTypes 获取所有字典类型，按 sort 正序排列
func (uc *UseCase) ListDictTypes(ctx context.Context) ([]data.DictType, error) {
	types, err := uc.repo.ListDictTypes(ctx)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, "failed to list dict types", err)
	}
	return types, nil
}

// CreateDictType 创建字典类型
func (uc *UseCase) CreateDictType(ctx context.Context, code, name, description string, sort int) error {
	existing, err := uc.repo.GetDictTypeByCode(ctx, code)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to query dict type", err)
	}
	if existing != nil {
		return apperr.New(apperr.CodeConflict, "dict type code already exists")
	}

	dictType := data.DictType{
		Code:        code,
		Name:        name,
		Description: description,
		Sort:        sort,
	}
	if err := uc.repo.CreateDictType(ctx, &dictType); err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to create dict type", err)
	}

	uc.logger.Info("dict type created", zap.String("code", code))
	return nil
}

// UpdateDictType 更新字典类型（code 不可修改）
func (uc *UseCase) UpdateDictType(ctx context.Context, id uint, name, description string, sort int) error {
	updated, err := uc.repo.UpdateDictType(ctx, id, name, description, sort)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to update dict type", err)
	}
	if !updated {
		return apperr.New(apperr.CodeNotFound, "dict type not found")
	}

	uc.logger.Info("dict type updated", zap.Uint("id", id))
	return nil
}

// DeleteDictType 软删除字典类型
func (uc *UseCase) DeleteDictType(ctx context.Context, id uint) error {
	deleted, err := uc.repo.DeleteDictType(ctx, id)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to delete dict type", err)
	}
	if !deleted {
		return apperr.New(apperr.CodeNotFound, "dict type not found")
	}

	uc.logger.Info("dict type deleted", zap.Uint("id", id))
	return nil
}

// ListDictItems 根据类型编码获取字典数据，按 sort 正序排列
func (uc *UseCase) ListDictItems(ctx context.Context, typeCode string) ([]data.DictItem, error) {
	items, err := uc.repo.ListDictItems(ctx, typeCode)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, "failed to list dict items", err)
	}
	return items, nil
}

// CreateDictItem 创建字典数据
func (uc *UseCase) CreateDictItem(ctx context.Context, typeCode, itemKey, itemValue, description string, sort int) error {
	item := data.DictItem{
		TypeCode:    typeCode,
		ItemKey:     itemKey,
		ItemValue:   itemValue,
		Description: description,
		Sort:        sort,
		Status:      1,
	}
	if err := uc.repo.CreateDictItem(ctx, &item); err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to create dict item", err)
	}

	uc.logger.Info("dict item created",
		zap.String("type_code", typeCode),
		zap.String("item_key", itemKey),
	)
	return nil
}

// UpdateDictItem 更新字典数据
func (uc *UseCase) UpdateDictItem(ctx context.Context, id uint, itemKey, itemValue, description string, sort, status int) error {
	updated, err := uc.repo.UpdateDictItem(ctx, id, itemKey, itemValue, description, sort, status)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to update dict item", err)
	}
	if !updated {
		return apperr.New(apperr.CodeNotFound, "dict item not found")
	}

	uc.logger.Info("dict item updated", zap.Uint("id", id))
	return nil
}

// DeleteDictItem 软删除字典数据
func (uc *UseCase) DeleteDictItem(ctx context.Context, id uint) error {
	deleted, err := uc.repo.DeleteDictItem(ctx, id)
	if err != nil {
		return apperr.Wrap(apperr.CodeInternal, "failed to delete dict item", err)
	}
	if !deleted {
		return apperr.New(apperr.CodeNotFound, "dict item not found")
	}

	uc.logger.Info("dict item deleted", zap.Uint("id", id))
	return nil
}
