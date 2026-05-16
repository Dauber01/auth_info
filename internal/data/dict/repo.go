package dict

import (
	"context"
	"errors"

	"gorm.io/gorm"

	bizdict "auth_info/internal/biz/dict"
)

// DictRepo 实现 bizdict.DictRepository。
type DictRepo struct {
	db *gorm.DB
}

// NewDictRepository Wire Provider
func NewDictRepository(db *gorm.DB) *DictRepo {
	return &DictRepo{db: db}
}

func (r *DictRepo) ListDictTypes(ctx context.Context) ([]bizdict.DictType, error) {
	var types []DictType
	if err := r.db.WithContext(ctx).Order("sort asc, id asc").Find(&types).Error; err != nil {
		return nil, err
	}
	return dictTypesToBiz(types), nil
}

func (r *DictRepo) GetDictTypeByCode(ctx context.Context, code string) (*bizdict.DictType, error) {
	var dictType DictType
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&dictType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dictTypeToBiz(&dictType), nil
}

func (r *DictRepo) CreateDictType(ctx context.Context, dictType *bizdict.DictType) error {
	model := dictTypeFromBiz(dictType)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	dictType.ID = model.ID
	dictType.CreatedAt = model.CreatedAt
	dictType.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *DictRepo) UpdateDictType(ctx context.Context, id uint, name, description string, sort int) (bool, error) {
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

func (r *DictRepo) DeleteDictType(ctx context.Context, id uint) (bool, error) {
	result := r.db.WithContext(ctx).Delete(&DictType{}, id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *DictRepo) ListDictItems(ctx context.Context, typeCode string) ([]bizdict.DictItem, error) {
	var items []DictItem
	if err := r.db.WithContext(ctx).Where("type_code = ?", typeCode).Order("sort asc, id asc").Find(&items).Error; err != nil {
		return nil, err
	}
	return dictItemsToBiz(items), nil
}

func (r *DictRepo) CreateDictItem(ctx context.Context, item *bizdict.DictItem) error {
	model := dictItemFromBiz(item)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	item.ID = model.ID
	item.CreatedAt = model.CreatedAt
	item.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *DictRepo) UpdateDictItem(ctx context.Context, id uint, itemKey, itemValue, description string, sort, status int) (bool, error) {
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

func (r *DictRepo) DeleteDictItem(ctx context.Context, id uint) (bool, error) {
	result := r.db.WithContext(ctx).Delete(&DictItem{}, id)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
