package auth

import "time"

// User 鉴权领域 DTO，不携带 GORM 标签或框架依赖。
// data 层负责在持久化模型与本 DTO 之间转换。
type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  string
	Role      string
}
