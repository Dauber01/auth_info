package app

import (
	"github.com/gin-gonic/gin"

	"auth_info/internal/handler"
)

func registerDocumentRoutes(protected *gin.RouterGroup, documentHandler *handler.DocumentHandler) {
	doc := protected.Group("/document")
	{
		doc.POST("/generate-pdf", documentHandler.GeneratePDF)
		doc.POST("/generate-word", documentHandler.GenerateWord)
	}
}
