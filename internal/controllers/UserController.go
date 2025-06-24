package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
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
	var req dto.UserRegisterRequest
	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
		return
	}

	user, err := c.authService.Register(req.Email, req.Password, req.Username)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			status = http.StatusConflict
		}
		ctx.JSON(status, utils.ErrorResponse("error", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	}))
}

func (c *UserController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
		return
	}

	user, token, errs := c.authService.Login(req.Email, req.Password)
	if len(errs) > 0  {
		var fieldErrors []utils.FieldError
        for _, msg := range errs {
            fieldErrors = append(fieldErrors, utils.FieldError{Field: "credentials", Message: msg})
        }
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "errors": fieldErrors})
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
	uid, ok := utils.GetUserIDFromContext(context)
	if !ok {
		return
	}

    user, err := c.UserService.GetUserByID(uint(uid))
    if err != nil {
        context.JSON(http.StatusNotFound, utils.ErrorResponse("user",utils.ErrUserNotFound))
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
	var req dto.ChangePasswordRequest
	if validationErrs := utils.BindAndValidate(context, &req); len(validationErrs) > 0 {
		context.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
		return
	}

	uid, ok := utils.GetUserIDFromContext(context)
	if !ok { 
		return 
	}

    err := c.UserService.ChangePassword(uint(uid), req.OldPassword, req.NewPassword)
    if err != nil {
        if err.Error() == "old password is incorrect" {
            context.JSON(http.StatusBadRequest, utils.ErrorResponse("password",utils.ErrOldPasswordIncorrect))
            return
        }
        context.JSON(http.StatusBadRequest, utils.ErrorResponse("error", err.Error()))
        return
    }

    context.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgPasswordChanged}))
}

func (c *UserController) ListUsers(ctx *gin.Context) {
    page, pageSize, ok := utils.GetPaginationParams(ctx)
	if !ok {
		return
	}

    users, total, err := c.UserService.GetAllUsers(page, pageSize)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("users",utils.ErrCouldNotFetchUsers))
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
	var req dto.ChangeUserRole
	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
		return
	}

	userID, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidUserID)
	if !ok {
		return
	}

    err := c.UserService.ChangeUserRole(userID, req.Role)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("role",err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgUserRoleUpdated}))
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
    userID, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidUserID)
	if !ok {
		return
	}

    err := c.UserService.DeleteUser(userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse("user",err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgUserDeleted}))
}

func (c *UserController) GetUserDetail(ctx *gin.Context) {
    userID, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidUserID)
	if !ok {
		return
	}

    user, err := c.UserService.GetUserByID(userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse("user",utils.ErrUserNotFound))
        return
    }

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
    uid, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return
	}
    err := c.UserService.DeleteUser(uint(uid))
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse("user",err.Error()))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgUserDeleted}))
}

// func (c *UserController) ForgotPassword(ctx *gin.Context) {
// 	var req dto.ForgotPasswordRequest
// 	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
// 		return
// 	}

// 	err := c.authService.ForgotPassword(req.Email)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("email", err.Error()))
// 		return
// 	}
// }

// func (c *UserController) ResetPassword(ctx *gin.Context) {
// 	var req dto.ResetPasswordRequest
// 	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
// 		return
// 	}

// 	err := c.authService.ResetPassword(req.Token, req.NewPassword)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("reset_password", err.Error()))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "password reset successful"}))
// }

func (c *UserController) UpdateCanPost(ctx *gin.Context) {
    userID, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidUserID)
    if !ok {
        return
    }
    var req dto.UpdateCanPostRequest
    if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
        return
    }
    err := c.UserService.UpdateCanPost(userID, req.CanPost)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("user", err.Error()))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "successfully updated posting permission"}))
}