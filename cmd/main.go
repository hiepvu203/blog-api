// @title Blog API
// @version 1.0
// @description API cho blog project
// @host localhost:9090
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	_ "blog-api/docs"
	"blog-api/internal/config"
	"blog-api/internal/routes"
	"log"
	"os"
	"regexp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func init() {
	if os.Getenv("RENDER") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found (this is OK in production)")
		}
	}
}

func SlugValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	matched, _ := regexp.MatchString(`^[a-z0-9]+(?:-[a-z0-9]+)*$`, value)
	return matched
}

func UsernameValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, value)
	return matched
}

func StrongPasswordValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if len(value) < 6 {
		return false
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(value)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(value)
	hasDigit := regexp.MustCompile(`\d`).MatchString(value)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(value)
	return hasLower && hasUpper && hasDigit && hasSpecial
}

func main() {
	config.LoadEnv()
	config.ConnectDB()
	config.InitDB()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validations := map[string]validator.Func{
			"slug":     SlugValidator,
			"username": UsernameValidator,
			"strong":   StrongPasswordValidator,
		}

		for tag, fn := range validations {
			if err := v.RegisterValidation(tag, fn); err != nil {
				return
			}
		}
	}

	routes.SetupUserRoutes(r, config.DB)
	routes.SetupCategoryRoutes(r, config.DB)
	routes.SetupPostRoutes(r, config.DB)
	routes.SetupCommentRoutes(r, config.DB)

	err := r.Run(":" + os.Getenv("PORT"))
	if err != nil {
		return
	}
}
