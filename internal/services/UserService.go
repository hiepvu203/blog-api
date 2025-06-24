package services

import (
	"blog-api/internal/entities"
	"blog-api/internal/repositories"
	"blog-api/pkg/helper"
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

func (s *UserService) GetAllUsers(page, pageSize int) ([]entities.User, int64, error){
	return s.userRepo.ListAll(page, pageSize)
}

func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return err
    }
    if !helper.CheckPasswordHash(oldPassword, user.Password) {
        return errors.New("old password is incorrect")
    }
    hashed, err := helper.HashPassword(newPassword)
    if err != nil {
        return err
    }
    user.Password = hashed
    return s.userRepo.Update(user)
}

func (s *UserService) ChangeUserRole(userID uint, newRole string) error {
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return err
    }
    if newRole != "admin" && newRole != "client" {
        return errors.New("invalid role")
    }
    user.Role = newRole
    return s.userRepo.Update(user)
}

func (s *UserService) UpdateCanPost(userID uint, canPost bool) error {
    return s.userRepo.UpdateCanPost(userID, canPost)
}