package data

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password string `gorm:"size:256;not null"            json:"-"`
	Role     string `gorm:"size:32;default:'user'"       json:"role"`
}

type DictType struct {
	gorm.Model
	Code        string `gorm:"uniqueIndex;size:64;not null"  json:"code"`
	Name        string `gorm:"size:128;not null"             json:"name"`
	Description string `gorm:"size:256;default:''"           json:"description"`
	Sort        int    `gorm:"default:0"                     json:"sort"`
}

type DictItem struct {
	gorm.Model
	TypeCode    string `gorm:"index;size:64;not null"        json:"type_code"`
	ItemKey     string `gorm:"size:64;not null"              json:"item_key"`
	ItemValue   string `gorm:"size:256;not null"             json:"item_value"`
	Description string `gorm:"size:256;default:''"           json:"description"`
	Sort        int    `gorm:"default:0"                     json:"sort"`
	Status      int    `gorm:"default:1"                     json:"status"`
}
