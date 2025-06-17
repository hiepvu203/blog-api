package dto

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}

type UserUpdateRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20"` // update only if valid
	Email    string `json:"email" binding:"omitempty,email"`           // Validate format if valid
	Password string `json:"password" binding:"omitempty,min=6"`        // hashed when save before
	Role     string `json:"role" binding:"omitempty,oneof=admin client"` // Chỉ admin được gửi field này
}

type AdminUpdateUserRequest struct {
	UserUpdateRequest
	Role             string `json:"role" binding:"omitempty,oneof=admin client"` // Only admin can post
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}