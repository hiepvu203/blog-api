package main

import (
	"blog-api/internal/config"
	"blog-api/internal/routes"
	"regexp"

	// "blog-api/pkg/utils"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

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
    matched, _ := regexp.MatchString(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[^a-zA-Z0-9]).{8,}$`, value)
    return matched
}

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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("slug", SlugValidator)
		v.RegisterValidation("username", UsernameValidator)
    	v.RegisterValidation("strongpwd", StrongPasswordValidator)
    }

	// routes.SetupRoutes(r)
	routes.SetupUserRoutes(r, config.DB)
	routes.SetupCategoryRoutes(r, config.DB)
	routes.SetupPostRoutes(r, config.DB)
	routes.SetupCommentRoutes(r, config.DB)

	// r.NoRoute(func(ctx *gin.Context) {
    //     ctx.JSON(404, utils.ErrorResponse("Endpoint not found"))
    // })

	r.Run(":" + os.Getenv("PORT"))
}