package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	authService *services.AuthService
	UserService *services.UserService
}

func NewUserController(authService *services.AuthService, userService *services.UserService) *UserController {
	return &UserController{
		authService: authService,
		UserService: userService,
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Username string `json:"username" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	user, err := c.authService.Register(req.Email, req.Password, req.Username)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			status = http.StatusConflict
		}
		ctx.JSON(status, utils.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	}))
}

func (c *UserController) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	user, token, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	}))
}

func (c *UserController) GetMe(context *gin.Context){
	userID, ok := context.Get("userID")
    if !ok {
        context.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
        return
    }

    uid, ok := userID.(float64)
    if !ok {
        context.JSON(http.StatusInternalServerError, utils.ErrorResponse("Invalid userID type"))
        return
    }

    user, err := c.UserService.GetUserByID(uint(uid))
    if err != nil {
        context.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
        return
    }

    context.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
        "id":       user.ID,
        "email":    user.Email,
        "username": user.Username,
        "role":     user.Role,
    }))
}

func (c *UserController) ChangePassword(context *gin.Context) {
	userID, ok := context.Get("userID")
    if !ok {
        context.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
        return
    }

    var req dto.ChangePasswordRequest
    if err := context.ShouldBindJSON(&req); err != nil {
        context.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    uid, ok := userID.(float64)
    if !ok {
        context.JSON(http.StatusInternalServerError, utils.ErrorResponse("Invalid userID type"))
        return
    }

    err := c.UserService.ChangePassword(uint(uid), req.OldPassword, req.NewPassword)
    if err != nil {
        context.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    context.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "Password changed successfully"}))
}

func (c *UserController) ListUsers(ctx *gin.Context) {
    users, err := c.UserService.GetAllUsers()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Could not fetch users"))
        return
    }

    var resp []dto.UserResponse
    for _, u := range users {
        resp = append(resp, dto.UserResponse{
            ID:       uint(u.ID),
            Username: u.Username,
            Email:    u.Email,
            Role:     u.Role,
        })
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(resp))
}

func (c *UserController) ChangeUserRole(ctx *gin.Context) {
    var req struct {
        Role string `json:"role" binding:"required,oneof=admin client"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    idParam := ctx.Param("id")
    var userID uint
    _, err := fmt.Sscanf(idParam, "%d", &userID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user id"))
        return
    }

    err = c.UserService.ChangeUserRole(userID, req.Role)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "User role updated"}))
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
    idParam := ctx.Param("id")
    var userID uint
    _, err := fmt.Sscanf(idParam, "%d", &userID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user id"))
        return
    }

    err = c.UserService.DeleteUser(userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "User deleted successfully"}))
}