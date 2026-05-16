package dict

import "gorm.io/gorm"

// DictType 字典类型持久化模型。仅在 data 层使用。
type DictType struct {
	gorm.Model
	Code        string `gorm:"uniqueIndex;size:64;not null"  json:"code"`
	Name        string `gorm:"size:128;not null"             json:"name"`
	Description string `gorm:"size:256;default:''"           json:"description"`
	Sort        int    `gorm:"default:0"                     json:"sort"`
}

// DictItem 字典数据持久化模型。仅在 data 层使用。
type DictItem struct {
	gorm.Model
	TypeCode    string `gorm:"index;size:64;not null"        json:"type_code"`
	ItemKey     string `gorm:"size:64;not null"              json:"item_key"`
	ItemValue   string `gorm:"size:256;not null"             json:"item_value"`
	Description string `gorm:"size:256;default:''"           json:"description"`
	Sort        int    `gorm:"default:0"                     json:"sort"`
	Status      int    `gorm:"default:1"                     json:"status"`
}
