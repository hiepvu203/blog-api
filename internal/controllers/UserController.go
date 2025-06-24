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
	if validationErrs := utils.BindAndValidate(ctx, &req); validationErrs != nil {
		utils.SendFail(ctx, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
		return
	}

	user, err := c.authService.Register(req.Email, req.Password, req.Username)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			status = http.StatusConflict
		}
		utils.SendFail(ctx, status, "REGISTER_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(ctx, http.StatusCreated, "201", "registered", gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	})
}

func (c *UserController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if validationErrs := utils.BindAndValidate(ctx, &req); validationErrs != nil {
		utils.SendFail(ctx, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
		return
	}

	user, token, errs := c.authService.Login(req.Email, req.Password)
	if len(errs) > 0  {
		var fieldErrors []map[string]string
		for _, msg := range errs {
			fieldErrors = append(fieldErrors, map[string]string{"field": "credentials", "message": msg})
		}
		utils.SendFail(ctx, http.StatusUnauthorized, "401", "incorrect login information", fieldErrors)
        return
	}

	utils.SendSuccess(ctx, http.StatusOK, "200", "logged", gin.H{
        "token": token,
        "user": gin.H{
            "id":    user.ID,
            "email": user.Email,
            "role":  user.Role,
        },
    })
}

func (c *UserController) GetMe(context *gin.Context){
	uid, ok := utils.GetUserIDFromContext(context)
	if !ok {
		return
	}

    user, err := c.UserService.GetUserByID(uint(uid))
    if err != nil {
        utils.SendFail(context, http.StatusNotFound, "404", utils.ErrUserNotFound, nil)
        return
    }

    utils.SendSuccess(context, http.StatusOK, "200", "user information successfully retrieved", gin.H{
        "id":       user.ID,
        "email":    user.Email,
        "username": user.Username,
        "role":     user.Role,
    })
}

func (c *UserController) ChangePassword(context *gin.Context) {
	var req dto.ChangePasswordRequest
	if validationErrs := utils.BindAndValidate(context, &req); validationErrs != nil {
		utils.SendFail(context, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
		return
	}

	uid, ok := utils.GetUserIDFromContext(context)
	if !ok { 
		return 
	}

    err := c.UserService.ChangePassword(uint(uid), req.OldPassword, req.NewPassword)
    if err != nil {
        if err.Error() == "old password is incorrect" {
            utils.SendFail(context, http.StatusBadRequest, "400", utils.ErrOldPasswordIncorrect, nil)
            return
        }
        utils.SendFail(context, http.StatusBadRequest, "400", err.Error(), nil)
        return
    }

    utils.SendSuccess(context, http.StatusOK, "200", utils.MsgPasswordChanged, nil)
}

func (c *UserController) ListUsers(ctx *gin.Context) {
    page, pageSize, ok := utils.GetPaginationParams(ctx)
	if !ok {
		return
	}

    users, total, err := c.UserService.GetAllUsers(page, pageSize)
    if err != nil {
        utils.SendFail(ctx, http.StatusInternalServerError, "500", utils.ErrCouldNotFetchUsers, nil)
        return
    }
    meta := gin.H{"page": page, "page_size": pageSize, "total": total}
    var resp []dto.UserResponse
    for _, u := range users {
        resp = append(resp, dto.UserResponse{
            ID:       uint(u.ID),
            Username: u.Username,
            Email:    u.Email,
            Role:     u.Role,
        })
    }
    utils.SendSuccess(ctx, http.StatusOK, "200", utils.UserFetchOK, gin.H{"users": resp, "meta": meta})
}

func (c *UserController) ChangeUserRole(ctx *gin.Context) {
	var req dto.ChangeUserRole
	if validationErrs := utils.BindAndValidate(ctx, &req); validationErrs != nil {
		utils.SendFail(ctx, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
		return
	}

	userID, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidUserID)
	if !ok {
		return
	}

    err := c.UserService.ChangeUserRole(userID, req.Role)
    if err != nil {
        utils.SendFail(ctx, http.StatusBadRequest, "400", err.Error(), nil)
        return
    }

    utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgUserRoleUpdated, nil)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
    userID, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidUserID)
	if !ok {
		return
	}

    err := c.UserService.DeleteUser(userID)
    if err != nil {
        utils.SendFail(ctx, http.StatusNotFound, "404", err.Error(), nil)
        return
    }

    utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgUserDeleted, nil)
}

func (c *UserController) GetUserDetail(ctx *gin.Context) {
    userID, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidUserID)
	if !ok {
		return
	}
    user, err := c.UserService.GetUserByID(userID)
    if err != nil {
        utils.SendFail(ctx, http.StatusNotFound, "404", utils.ErrUserNotFound, nil)
        return
    }
    utils.SendSuccess(ctx, http.StatusOK, "200", "user information successfully retrieved", gin.H{
        "id":       user.ID,
        "email":    user.Email,
        "username": user.Username,
        "role":     user.Role,
    })
}

func (c *UserController) DeleteMe(ctx *gin.Context) {
    uid, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return
	}
    err := c.UserService.DeleteUser(uint(uid))
    if err != nil {
        utils.SendFail(ctx, http.StatusNotFound, "404", err.Error(), nil)
        return
    }
    utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgUserDeleted, nil)
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
        utils.SendFail(ctx, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
        return
    }
    err := c.UserService.UpdateCanPost(userID, req.CanPost)
    if err != nil {
        utils.SendFail(ctx, http.StatusInternalServerError, "500", err.Error(), nil)
        return
    }
    utils.SendSuccess(ctx, http.StatusOK, "200", "successfully updated posting permission", nil)
}