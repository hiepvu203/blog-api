package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/entities"
	"blog-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	authService *services.AuthService
}

func NewUserController(authService *services.AuthService) *UserController {
	return &UserController{authService: authService}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req dto.UserRegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user := entities.User{
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
		Role: "client",
	}

	if err := c.authService.Register(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.UserResponse{
		ID: uint(user.ID),
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
	})
}