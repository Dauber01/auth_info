package dict

import (
	"context"

	"auth_info/internal/data"
)

// DictRepository 字典数据访问接口（由 biz/dict 层定义，data 层实现）
type DictRepository interface {
	ListDictTypes(ctx context.Context) ([]data.DictType, error)
	GetDictTypeByCode(ctx context.Context, code string) (*data.DictType, error)
	CreateDictType(ctx context.Context, dictType *data.DictType) error
	UpdateDictType(ctx context.Context, id uint, name, description string, sort int) (bool, error)
	DeleteDictType(ctx context.Context, id uint) (bool, error)
	ListDictItems(ctx context.Context, typeCode string) ([]data.DictItem, error)
	CreateDictItem(ctx context.Context, item *data.DictItem) error
	UpdateDictItem(ctx context.Context, id uint, itemKey, itemValue, description string, sort, status int) (bool, error)
	DeleteDictItem(ctx context.Context, id uint) (bool, error)
}
