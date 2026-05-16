package auth

import (
	"github.com/gin-gonic/gin"

	authhdl "auth_info/internal/handler/auth"
)

// Register 注册 auth 模块的公开路由（无需鉴权）。
func Register(api *gin.RouterGroup, h *authhdl.Handler) {
	g := api.Group("/auth")
	{
		g.POST("/register", h.Register)
		g.POST("/login", h.Login)
	}
}
