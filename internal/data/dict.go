package data

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type DictRepository interface {
	ListDictTypes(ctx context.Context) ([]DictType, error)
	GetDictTypeByCode(ctx context.Context, code string) (*DictType, error)
	CreateDictType(ctx context.Context, dictType *DictType) error
	UpdateDictType(ctx context.Context, id uint, name, description string, sort int) (bool, error)
	DeleteDictType(ctx context.Context, id uint) (bool, error)
	ListDictItems(ctx context.Context, typeCode string) ([]DictItem, error)
	CreateDictItem(ctx context.Context, item *DictItem) error
	UpdateDictItem(ctx context.Context, id uint, itemKey, itemValue, description string, sort, status int) (bool, error)
	DeleteDictItem(ctx context.Context, id uint) (bool, error)
}

type dictRepository struct {
	db *gorm.DB
}

func NewDictRepository(db *gorm.DB) DictRepository {
	return &dictRepository{db: db}
}

func (r *dictRepository) ListDictTypes(ctx context.Context) ([]DictType, error) {
	var types []DictType
	if err := r.db.WithContext(ctx).Order("sort asc, id asc").Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

func (r *dictRepository) GetDictTypeByCode(ctx context.Context, code string) (*DictType, error) {
	var dictType DictType
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&dictType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dictType, nil
}

func (r *dictRepository) CreateDictType(ctx context.Context, dictType *DictType) error {
	return r.db.WithContext(ctx).Create(dictType).Error
}

func (r *dictRepository) UpdateDictType(ctx context.Context, id uint, name, description string, sort int) (bool, error) {
	result := r.db.WithContext(ctx).Model(&DictType{}).Where("id = ?", id).Updates(map[string]any{
		"name":        name,
		"description": description,
		"sort":        sort,
	})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *dictRepository) DeleteDictType(ctx context.Context, id uint) (bool, error) {
	result := r.db.WithContext(ctx).Delete(&DictType{}, id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *dictRepository) ListDictItems(ctx context.Context, typeCode string) ([]DictItem, error) {
	var items []DictItem
	if err := r.db.WithContext(ctx).Where("type_code = ?", typeCode).Order("sort asc, id asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *dictRepository) CreateDictItem(ctx context.Context, item *DictItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *dictRepository) UpdateDictItem(ctx context.Context, id uint, itemKey, itemValue, description string, sort, status int) (bool, error) {
	result := r.db.WithContext(ctx).Model(&DictItem{}).Where("id = ?", id).Updates(map[string]any{
		"item_key":    itemKey,
		"item_value":  itemValue,
		"description": description,
		"sort":        sort,
		"status":      status,
	})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *dictRepository) DeleteDictItem(ctx context.Context, id uint) (bool, error) {
	result := r.db.WithContext(ctx).Delete(&DictItem{}, id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
