package dict

import "gorm.io/gorm"

// DictType 字典类型，定义可配置信息的分类（如数据威胁等级、社交媒体类型）
type DictType struct {
	gorm.Model
	Code        string `gorm:"uniqueIndex;size:64;not null"  json:"code"`
	Name        string `gorm:"size:128;not null"             json:"name"`
	Description string `gorm:"size:256;default:''"           json:"description"`
	Sort        int    `gorm:"default:0"                     json:"sort"`
}

// DictItem 字典数据，存储各分类下的具体配置项
type DictItem struct {
	gorm.Model
	TypeCode    string `gorm:"index;size:64;not null"        json:"type_code"`
	ItemKey     string `gorm:"size:64;not null"              json:"item_key"`
	ItemValue   string `gorm:"size:256;not null"             json:"item_value"`
	Description string `gorm:"size:256;default:''"           json:"description"`
	Sort        int    `gorm:"default:0"                     json:"sort"`
	Status      int    `gorm:"default:1"                     json:"status"`
}
