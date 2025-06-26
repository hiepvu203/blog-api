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

// Register godoc
// @Summary Đăng ký người dùng mới
// @Description Đăng ký tài khoản với email, password và username
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body  dto.UserRegisterRequest  true  "Thông tin đăng ký"
// @Success 201 {object} map[string]interface{} "Đăng ký thành công"
// @Failure 400 {object} utils.APIResponse "Lỗi xác thực hoặc đăng ký"
// @Failure 409 {object} utils.APIResponse "Email đã tồn tại"
// @Router /users/register [post]
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

// Login godoc
// @Summary Đăng nhập
// @Description Đăng nhập với email và password
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body  dto.UserLoginRequest  true  "Thông tin đăng nhập"
// @Success 200 {object} map[string]interface{} "Đăng nhập thành công, trả về token và thông tin user"
// @Failure 400 {object} utils.APIResponse "Lỗi xác thực"
// @Failure 401 {object} utils.APIResponse "Sai thông tin đăng nhập"
// @Router /users/login [post]
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
            "id":    	user.ID,
			"username": user.Username,
            "email": 	user.Email,
            "role":  	user.Role,
        },
    })
}

// GetMe godoc
// @Summary Lấy thông tin người dùng hiện tại
// @Description Lấy thông tin user từ token
// @Tags users
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} map[string]interface{} "Thông tin user"
// @Failure 404 {object} utils.APIResponse "Không tìm thấy user"
// @Router /users/me [get]
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

// ChangePassword godoc
// @Summary Đổi mật khẩu
// @Description Đổi mật khẩu cho user hiện tại
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param   body  body  dto.ChangePasswordRequest  true  "Thông tin đổi mật khẩu"
// @Success 200 {object} utils.APIResponse "Đổi mật khẩu thành công"
// @Failure 400 {object} utils.APIResponse "Lỗi xác thực hoặc mật khẩu cũ sai"
// @Router /users/change-password [put]
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

// ListUsers godoc
// @Summary Lấy danh sách người dùng
// @Description Lấy danh sách user, có phân trang
// @Tags users
// @Security BearerAuth
// @Produce  json
// @Param   page     query   int  false  "Trang hiện tại"
// @Param   page_size query  int  false  "Số lượng mỗi trang"
// @Success 200 {object} map[string]interface{} "Danh sách user và meta"
// @Failure 500 {object} utils.APIResponse "Lỗi server"
// @Router /admin/users [get]
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

// ChangeUserRole godoc
// @Summary Đổi vai trò người dùng
// @Description Đổi role cho user theo id
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param   id    path   int  true  "ID người dùng"
// @Param   body  body   dto.ChangeUserRole  true  "Thông tin role mới"
// @Success 200 {object} utils.APIResponse "Cập nhật role thành công"
// @Failure 400 {object} utils.APIResponse "Lỗi xác thực hoặc không tìm thấy user"
// @Router /admin/users/{id}/role [put]
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

// DeleteUser godoc
// @Summary Xóa người dùng
// @Description Xóa user theo id
// @Tags users
// @Security BearerAuth
// @Produce  json
// @Param   id  path  int  true  "ID người dùng"
// @Success 200 {object} utils.APIResponse "Xóa thành công"
// @Failure 404 {object} utils.APIResponse "Không tìm thấy user"
// @Router /admin/users/{id} [delete]
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

// GetUserDetail godoc
// @Summary Lấy chi tiết người dùng
// @Description Lấy thông tin user theo id
// @Tags users
// @Security BearerAuth
// @Produce  json
// @Param   id  path  int  true  "ID người dùng"
// @Success 200 {object} map[string]interface{} "Thông tin user"
// @Failure 404 {object} utils.APIResponse "Không tìm thấy user"
// @Router /admin/users/{id} [get]
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

// DeleteMe godoc
// @Summary Xóa tài khoản của chính mình
// @Description Xóa user hiện tại (self-delete)
// @Tags users
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} utils.APIResponse "Xóa thành công"
// @Failure 404 {object} utils.APIResponse "Không tìm thấy user"
// @Router /users/me [delete]
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

// UpdateCanPost godoc
// @Summary Cập nhật quyền đăng bài
// @Description Cập nhật quyền đăng bài cho user theo id
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param   id    path   int  true  "ID người dùng"
// @Param   body  body   dto.UpdateCanPostRequest  true  "Trạng thái quyền đăng bài"
// @Success 200 {object} utils.APIResponse "Cập nhật thành công"
// @Failure 400 {object} utils.APIResponse "Lỗi xác thực"
// @Failure 500 {object} utils.APIResponse "Lỗi server"
// @Router /admin/users/{id}/ban-post [put]
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