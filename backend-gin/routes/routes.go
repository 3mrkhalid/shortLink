package routes

import (
	
	"shortlink/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	fmt.Println("Routes loaded")


	auth := router.Group("/api/v1/users")
	{
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/register", handlers.RegisterHandler)
	}
}