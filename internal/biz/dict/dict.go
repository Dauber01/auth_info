package dict

import (
	"errors"

	"go.uber.org/zap"

	"auth_info/internal/data"
	"auth_info/internal/logger"
)

// UseCase 字典配置业务逻辑
type UseCase struct {
	repo data.DictRepository
}

// NewUseCase Wire Provider
func NewUseCase(repo data.DictRepository) *UseCase {
	return &UseCase{repo: repo}
}

// ListDictTypes 获取所有字典类型，按 sort 正序排列
func (uc *UseCase) ListDictTypes() ([]data.DictType, error) {
	return uc.repo.ListDictTypes()
}

// CreateDictType 创建字典类型
func (uc *UseCase) CreateDictType(code, name, description string, sort int) error {
	existing, err := uc.repo.GetDictTypeByCode(code)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("dict type code already exists")
	}

	dictType := data.DictType{
		Code:        code,
		Name:        name,
		Description: description,
		Sort:        sort,
	}
	if err := uc.repo.CreateDictType(&dictType); err != nil {
		return err
	}

	logger.GetLogger().Info("dict type created", zap.String("code", code))
	return nil
}

// UpdateDictType 更新字典类型（code 不可修改）
func (uc *UseCase) UpdateDictType(id uint, name, description string, sort int) error {
	updated, err := uc.repo.UpdateDictType(id, name, description, sort)
	if err != nil {
		return err
	}
	if !updated {
		return errors.New("dict type not found")
	}

	logger.GetLogger().Info("dict type updated", zap.Uint("id", id))
	return nil
}

// DeleteDictType 软删除字典类型
func (uc *UseCase) DeleteDictType(id uint) error {
	deleted, err := uc.repo.DeleteDictType(id)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("dict type not found")
	}

	logger.GetLogger().Info("dict type deleted", zap.Uint("id", id))
	return nil
}

// ListDictItems 根据类型编码获取字典数据，按 sort 正序排列
func (uc *UseCase) ListDictItems(typeCode string) ([]data.DictItem, error) {
	return uc.repo.ListDictItems(typeCode)
}

// CreateDictItem 创建字典数据
func (uc *UseCase) CreateDictItem(typeCode, itemKey, itemValue, description string, sort int) error {
	item := data.DictItem{
		TypeCode:    typeCode,
		ItemKey:     itemKey,
		ItemValue:   itemValue,
		Description: description,
		Sort:        sort,
		Status:      1,
	}
	if err := uc.repo.CreateDictItem(&item); err != nil {
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
	updated, err := uc.repo.UpdateDictItem(id, itemKey, itemValue, description, sort, status)
	if err != nil {
		return err
	}
	if !updated {
		return errors.New("dict item not found")
	}

	logger.GetLogger().Info("dict item updated", zap.Uint("id", id))
	return nil
}

// DeleteDictItem 软删除字典数据
func (uc *UseCase) DeleteDictItem(id uint) error {
	deleted, err := uc.repo.DeleteDictItem(id)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("dict item not found")
	}

	logger.GetLogger().Info("dict item deleted", zap.Uint("id", id))
	return nil
}
