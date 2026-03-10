package app

import (
	"github.com/gin-gonic/gin"

	"auth_info/internal/handler"
)

func registerHelloRoutes(protected *gin.RouterGroup, helloHandler *handler.HelloHandler) {
	protected.GET("/hello", helloHandler.Hello)
}
