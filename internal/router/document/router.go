package document

import (
	"github.com/gin-gonic/gin"

	dochdl "auth_info/internal/handler/document"
)

// Register 注册文档模块路由（鉴权后）。
func Register(protected *gin.RouterGroup, h *dochdl.Handler) {
	g := protected.Group("/document")
	{
		g.POST("/generate-pdf", h.GeneratePDF)
		g.POST("/generate-word", h.GenerateWord)
	}
}
