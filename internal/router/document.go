package router

import (
	"github.com/gin-gonic/gin"

	"auth_info/internal/handler"
)

func RegisterDocumentRoutes(protected *gin.RouterGroup, documentHandler *handler.DocumentHandler) {
	doc := protected.Group("/document")
	{
		doc.POST("/generate-pdf", documentHandler.GeneratePDF)
		doc.POST("/generate-word", documentHandler.GenerateWord)
	}
}
