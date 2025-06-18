package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"fmt"
	"net/http"
	"strconv"

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
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ParseValidationError(err)))
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
        context.JSON(http.StatusUnauthorized, utils.ErrorResponse(utils.ErrUnauthorized))
        return
    }

    uid, ok := userID.(float64)
    if !ok {
        context.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrInvalidUserIDType))
        return
    }

    user, err := c.UserService.GetUserByID(uint(uid))
    if err != nil {
        context.JSON(http.StatusNotFound, utils.ErrorResponse(utils.ErrUserNotFound))
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
        context.JSON(http.StatusUnauthorized, utils.ErrorResponse(utils.ErrUnauthorized))
        return
    }

    var req dto.ChangePasswordRequest
    if err := context.ShouldBindJSON(&req); err != nil {
        context.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ParseValidationError(err)))
        return
    }

    uid, ok := userID.(float64)
    if !ok {
        context.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrInvalidUserIDType))
        return
    }

    err := c.UserService.ChangePassword(uint(uid), req.OldPassword, req.NewPassword)
    if err != nil {
        if err.Error() == "old password is incorrect" {
            context.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrOldPasswordIncorrect))
            return
        }
        context.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    context.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgPasswordChanged}))
}

func (c *UserController) ListUsers(ctx *gin.Context) {
    page := 1
    pageSize := 10
    if p := ctx.Query("page"); p != "" {
        if v, err := strconv.Atoi(p); err == nil && v > 0 {
            page = v
        } else {
            ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidPageParam))
            return
        }
    }
    if ps := ctx.Query("page_size"); ps != "" {
        if v, err := strconv.Atoi(ps); err == nil && v > 0 {
            pageSize = v
        } else {
            ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidPageSizeParam))
            return
        }
    }

    users, total, err := c.UserService.GetAllUsers(page, pageSize)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCouldNotFetchUsers))
        return
    }
    if total == 0 {
        ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
            "data":    []dto.UserResponse{},
            "total":   0,
            "page":    page,
            "page_size": pageSize,
            "message": utils.ErrUserNotFound,
        }))
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

    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
        "data":    resp,
        "total":   total,
        "page":    page,
        "page_size": pageSize,
        "message": utils.UserFetchOK,
    }))
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
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidUserID))
        return
    }

    err = c.UserService.ChangeUserRole(userID, req.Role)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgUserRoleUpdated}))
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
    idParam := ctx.Param("id")
    var userID uint
    _, err := fmt.Sscanf(idParam, "%d", &userID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidUserID))
        return
    }

    err = c.UserService.DeleteUser(userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgUserDeleted}))
}

func (c *UserController) GetUserDetail(ctx *gin.Context) {
    idParam := ctx.Param("id")
    var userID uint
    _, err := fmt.Sscanf(idParam, "%d", &userID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidUserID))
        return
    }

    user, err := c.UserService.GetUserByID(userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse(utils.ErrUserNotFound))
        return
    }

    // Count posts and comments user's
    postCount := len(user.Posts)
    commentCount := len(user.Comments)

    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
        "id":           user.ID,
        "username":     user.Username,
        "email":        user.Email,
        "role":         user.Role,
        "created_at":   user.CreatedAt,
        "post_count":   postCount,
        "comment_count": commentCount,
    }))
}

func (c *UserController) DeleteMe(ctx *gin.Context) {
    userID, ok := ctx.Get("userID")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(utils.ErrUnauthorized))
        return
    }
    uid, ok := userID.(float64)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrInvalidUserIDType))
        return
    }
    err := c.UserService.DeleteUser(uint(uid))
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error()))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgUserDeleted}))
}