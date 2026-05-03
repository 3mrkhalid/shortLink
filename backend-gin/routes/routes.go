package routes

import (
	"fmt"
	"shortlink/handlers"
	"shortlink/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	fmt.Println("Routes loaded")

	// public routes
	auth := router.Group("/api/v1/users")
	{
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/register", handlers.RegisterHandler)
	}

	// protected routes
	protected := router.Group("/api/v1/users")
	protected.Use(middlewares.VerifyJwt)
	{
		// protected.GET("/profile", handlers.ProfileHandler)
		// protected.POST("/logout", handlers.LogoutHandler)
	}
}