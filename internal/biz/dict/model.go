package dict

import "time"

// DictType 字典类型 DTO
type DictType struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Code        string
	Name        string
	Description string
	Sort        int
}

// DictItem 字典数据 DTO
type DictItem struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TypeCode    string
	ItemKey     string
	ItemValue   string
	Description string
	Sort        int
	Status      int
}
