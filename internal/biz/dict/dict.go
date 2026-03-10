package dict

import (
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"auth_info/internal/logger"
	modeldict "auth_info/internal/model/dict"
)

// UseCase 字典配置业务逻辑
type UseCase struct {
	db *gorm.DB
}

// NewUseCase Wire Provider
func NewUseCase(db *gorm.DB) *UseCase {
	return &UseCase{db: db}
}

// ─── DictType ─────────────────────────────────────────────────────────────────

// ListDictTypes 获取所有字典类型，按 sort 正序排列
func (uc *UseCase) ListDictTypes() ([]modeldict.DictType, error) {
	var types []modeldict.DictType
	if err := uc.db.Order("sort asc, id asc").Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

// CreateDictType 创建字典类型
func (uc *UseCase) CreateDictType(code, name, description string, sort int) error {
	var existing modeldict.DictType
	if err := uc.db.Where("code = ?", code).First(&existing).Error; err == nil {
		return errors.New("dict type code already exists")
	}

	dictType := modeldict.DictType{
		Code:        code,
		Name:        name,
		Description: description,
		Sort:        sort,
	}
	if err := uc.db.Create(&dictType).Error; err != nil {
		return err
	}

	logger.GetLogger().Info("dict type created", zap.String("code", code))
	return nil
}

// UpdateDictType 更新字典类型（code 不可修改）
func (uc *UseCase) UpdateDictType(id uint, name, description string, sort int) error {
	result := uc.db.Model(&modeldict.DictType{}).Where("id = ?", id).Updates(map[string]any{
		"name":        name,
		"description": description,
		"sort":        sort,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("dict type not found")
	}

	logger.GetLogger().Info("dict type updated", zap.Uint("id", id))
	return nil
}

// DeleteDictType 软删除字典类型
func (uc *UseCase) DeleteDictType(id uint) error {
	result := uc.db.Delete(&modeldict.DictType{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("dict type not found")
	}

	logger.GetLogger().Info("dict type deleted", zap.Uint("id", id))
	return nil
}

// ─── DictItem ─────────────────────────────────────────────────────────────────

// ListDictItems 根据类型编码获取字典数据，按 sort 正序排列
func (uc *UseCase) ListDictItems(typeCode string) ([]modeldict.DictItem, error) {
	var items []modeldict.DictItem
	if err := uc.db.Where("type_code = ?", typeCode).Order("sort asc, id asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateDictItem 创建字典数据
func (uc *UseCase) CreateDictItem(typeCode, itemKey, itemValue, description string, sort int) error {
	item := modeldict.DictItem{
		TypeCode:    typeCode,
		ItemKey:     itemKey,
		ItemValue:   itemValue,
		Description: description,
		Sort:        sort,
		Status:      1,
	}
	if err := uc.db.Create(&item).Error; err != nil {
		return err
	}

	logger.GetLogger().Info("dict item created",
		zap.String("type_code", typeCode),
		zap.String("item_key", itemKey),
	)
	return nil
}

// UpdateDictItem 更新字典数据
func (uc *UseCase) UpdateDictItem(id uint, itemKey, itemValue, description string, sort, status int) error {
	result := uc.db.Model(&modeldict.DictItem{}).Where("id = ?", id).Updates(map[string]any{
		"item_key":    itemKey,
		"item_value":  itemValue,
		"description": description,
		"sort":        sort,
		"status":      status,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("dict item not found")
	}

	logger.GetLogger().Info("dict item updated", zap.Uint("id", id))
	return nil
}

// DeleteDictItem 软删除字典数据
func (uc *UseCase) DeleteDictItem(id uint) error {
	result := uc.db.Delete(&modeldict.DictItem{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("dict item not found")
	}

	logger.GetLogger().Info("dict item deleted", zap.Uint("id", id))
	return nil
}
