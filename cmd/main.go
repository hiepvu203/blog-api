package main

import (
	"blog-api/internal/config"
	"blog-api/internal/routes"
	"blog-api/pkg/utils"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main(){
	config.LoadEnv()
	config.ConnectDB()
	config.InitDB()

	r := gin.Default()
	// r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:4200"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	// routes.SetupRoutes(r)
	routes.SetupUserRoutes(r, config.DB)
	routes.SetupCategoryRoutes(r, config.DB)
	routes.SetupPostRoutes(r, config.DB)
	routes.SetupCommentRoutes(r, config.DB)

	r.NoRoute(func(ctx *gin.Context) {
        ctx.JSON(404, utils.ErrorResponse("Endpoint not found"))
    })

	r.Run(":" + os.Getenv("PORT"))
}