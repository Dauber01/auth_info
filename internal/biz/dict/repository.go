package dict

import "context"

// DictRepository 字典数据访问接口（biz 定义，data 实现）。
// 接口的输入输出仅使用 biz 层 DTO，保持 biz 与持久化模型解耦。
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
