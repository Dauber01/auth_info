package model

import "gorm.io/gorm"

// User 用户实体
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password string `gorm:"size:256;not null"            json:"-"`
	Role     string `gorm:"size:32;default:'user'"       json:"role"`
}
