package main

import (
	"shortlink/config"
	"shortlink/middlewares"
	"shortlink/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()
	config.DBConn()
	config.CreateUserIndexes()

	router := gin.Default()
	router.Use(middlewares.CorsMiddleware())
	routes.SetupRoutes(router)
	

	router.Run(":5000")
}
