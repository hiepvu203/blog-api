package main

import (
	"blog-api/internal/config"
	"blog-api/internal/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main(){
	config.LoadEnv()
	config.ConnectDB()
	// config.InitDB()

	r := gin.Default()
	// routes.SetupRoutes(r)
	routes.SetupUserRoutes(r, config.DB)
	routes.SetupCategoryRoutes(r, config.DB)
	routes.SetupPostRoutes(r, config.DB)

	r.Run(":" + os.Getenv("PORT"))
}