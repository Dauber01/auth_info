package auth

import "gorm.io/gorm"

// User 用户持久化模型。仅在 data 层使用。
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password string `gorm:"size:256;not null"            json:"-"`
	Role     string `gorm:"size:32;default:'user'"       json:"role"`
}
