package dto

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20,username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,strongpwd"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,strongpwd"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}

type UserUpdateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20,username"` 
	Email    string `json:"email" binding:"required,email"`           
	Password string `json:"password" binding:"required,min=6"`        
	Role     string `json:"role" binding:"required,oneof=admin client"` 
}

type AdminUpdateUserRequest struct {
	UserUpdateRequest
	Role     string `json:"role" binding:"required,oneof=admin client"`
}

type ForgotPasswordRequest struct {
	Email 	string 	`json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token 			string `json:"token" binding:"required"`
	NewPassword		string `json:"new_password" binding:"required,strongpwd"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,strongpwd"`
}

type ChangeUserRole struct {
	Role string `json:"role" binding:"required,oneof=admin client"`
}

type UpdateCanPostRequest struct {
    CanPost bool `json:"can_post" binding:"required"`
}