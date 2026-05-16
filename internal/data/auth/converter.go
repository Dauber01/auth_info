package auth

import bizauth "auth_info/internal/biz/auth"

// toBiz 将持久化模型转为 biz 层 DTO。
func toBiz(m *User) *bizauth.User {
	if m == nil {
		return nil
	}
	return &bizauth.User{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Username:  m.Username,
		Password:  m.Password,
		Role:      m.Role,
	}
}

// fromBiz 将 biz 层 DTO 转为持久化模型。
func fromBiz(d *bizauth.User) *User {
	if d == nil {
		return nil
	}
	m := &User{
		Username: d.Username,
		Password: d.Password,
		Role:     d.Role,
	}
	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	return m
}
