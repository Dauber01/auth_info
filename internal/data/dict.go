package data

import (
	"errors"

	"gorm.io/gorm"
)

type DictRepository interface {
	ListDictTypes() ([]DictType, error)
	GetDictTypeByCode(code string) (*DictType, error)
	CreateDictType(dictType *DictType) error
	UpdateDictType(id uint, name, description string, sort int) (bool, error)
	DeleteDictType(id uint) (bool, error)
	ListDictItems(typeCode string) ([]DictItem, error)
	CreateDictItem(item *DictItem) error
	UpdateDictItem(id uint, itemKey, itemValue, description string, sort, status int) (bool, error)
	DeleteDictItem(id uint) (bool, error)
}

type dictRepository struct {
	db *gorm.DB
}

func NewDictRepository(db *gorm.DB) DictRepository {
	return &dictRepository{db: db}
}

func (r *dictRepository) ListDictTypes() ([]DictType, error) {
	var types []DictType
	if err := r.db.Order("sort asc, id asc").Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

func (r *dictRepository) GetDictTypeByCode(code string) (*DictType, error) {
	var dictType DictType
	if err := r.db.Where("code = ?", code).First(&dictType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dictType, nil
}

func (r *dictRepository) CreateDictType(dictType *DictType) error {
	return r.db.Create(dictType).Error
}

func (r *dictRepository) UpdateDictType(id uint, name, description string, sort int) (bool, error) {
	result := r.db.Model(&DictType{}).Where("id = ?", id).Updates(map[string]any{
		"name":        name,
		"description": description,
		"sort":        sort,
	})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *dictRepository) DeleteDictType(id uint) (bool, error) {
	result := r.db.Delete(&DictType{}, id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *dictRepository) ListDictItems(typeCode string) ([]DictItem, error) {
	var items []DictItem
	if err := r.db.Where("type_code = ?", typeCode).Order("sort asc, id asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *dictRepository) CreateDictItem(item *DictItem) error {
	return r.db.Create(item).Error
}

func (r *dictRepository) UpdateDictItem(id uint, itemKey, itemValue, description string, sort, status int) (bool, error) {
	result := r.db.Model(&DictItem{}).Where("id = ?", id).Updates(map[string]any{
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

func (r *dictRepository) DeleteDictItem(id uint) (bool, error) {
	result := r.db.Delete(&DictItem{}, id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
