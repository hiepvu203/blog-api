package services

import (
	"fmt"
	"blog-api/internal/entities"
	"blog-api/internal/repositories"
	"blog-api/pkg/utils"
	"errors"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(user *entities.User) error {
	fmt.Printf("Checking email: %s\n", user.Email)
	existingUser, err := s.userRepo.FindEmail(user.Email)
	fmt.Printf("Result: %+v, Error: %v\n", existingUser, err)

	if err != nil {
        return fmt.Errorf("database error: %v", err)
    }

	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashPassword

	return s.userRepo.Create(user)
}