package app

import (
	"github.com/gin-gonic/gin"

	"auth_info/internal/handler"
)

func registerAuthRoutes(api *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
