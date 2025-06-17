package services

import (
	"blog-api/internal/dto"
	"blog-api/internal/entities"
	"blog-api/internal/repositories"
	"errors"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService{
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(id uint) (*entities.User, error){
	return s.userRepo.FindByID(id);
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) GetAllUsers() ([]entities.User, error){
	return s.userRepo.ListAll()
}

func (s *UserService) UpdateUser(currentUserID uint, targetUserID uint, updates dto.UserUpdateRequest) error {
	// get the user to update
	targetUser, err := s.userRepo.FindByID(targetUserID)

	if err != nil{
		return err
	}

	// Kiểm tra quyền:
    // - Admin: được sửa mọi thứ
    // - Client: chỉ được sửa chính mình và không được thay đổi role

	currentUser, err := s.userRepo.FindByID(currentUserID)
	if err != nil{
		return err
	}

	if currentUser.Role != "admin" && currentUser.ID != targetUser.ID {
		return errors.New("permission denied: you can only update your own profile")
	}

	if updates.Username != "" {
        targetUser.Username = updates.Username
    }
    
	// Chỉ admin được đổi role
    if currentUser.Role == "admin" && updates.Role != "" {
        targetUser.Role = updates.Role
    }

    return s.userRepo.Update(targetUser)
}