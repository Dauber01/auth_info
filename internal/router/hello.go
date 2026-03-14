package router

import (
	"github.com/gin-gonic/gin"

	"auth_info/internal/handler"
)

func RegisterHelloRoutes(protected *gin.RouterGroup, helloHandler *handler.HelloHandler) {
	protected.GET("/hello", helloHandler.Hello)
}
