package app

import (
	"github.com/gin-gonic/gin"

	"auth_info/internal/handler"
)

func registerDictRoutes(protected *gin.RouterGroup, dictHandler *handler.DictHandler) {
	dict := protected.Group("/dict")
	{
		dict.GET("/types", dictHandler.ListDictTypes)
		dict.POST("/types", dictHandler.CreateDictType)
		dict.PUT("/types/:id", dictHandler.UpdateDictType)
		dict.DELETE("/types/:id", dictHandler.DeleteDictType)

		dict.GET("/items", dictHandler.ListDictItems)
		dict.POST("/items", dictHandler.CreateDictItem)
		dict.PUT("/items/:id", dictHandler.UpdateDictItem)
		dict.DELETE("/items/:id", dictHandler.DeleteDictItem)
	}
}
