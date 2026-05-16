package hello

import (
	"github.com/gin-gonic/gin"

	hellohdl "auth_info/internal/handler/hello"
)

// Register 注册 hello 模块路由（鉴权后）。
func Register(protected *gin.RouterGroup, h *hellohdl.Handler) {
	protected.GET("/hello", h.Hello)
}
