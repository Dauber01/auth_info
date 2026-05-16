package dict

import (
	"github.com/gin-gonic/gin"

	dicthdl "auth_info/internal/handler/dict"
)

// Register 注册字典模块路由（鉴权后）。
func Register(protected *gin.RouterGroup, h *dicthdl.Handler) {
	g := protected.Group("/dict")
	{
		g.GET("/types", h.ListDictTypes)
		g.POST("/types", h.CreateDictType)
		g.PUT("/types/:id", h.UpdateDictType)
		g.DELETE("/types/:id", h.DeleteDictType)

		g.GET("/items", h.ListDictItems)
		g.POST("/items", h.CreateDictItem)
		g.PUT("/items/:id", h.UpdateDictItem)
		g.DELETE("/items/:id", h.DeleteDictItem)
	}
}
