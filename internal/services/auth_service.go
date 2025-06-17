package services

import (
	"blog-api/internal/entities"
	"blog-api/internal/repositories"
	"blog-api/pkg/utils"
	"errors"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(email, password, username string) (*entities.User, error) {
	if existing, err := s.userRepo.FindEmail(email); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if existing != nil {
		return nil, errors.New("email already exists")
	}
	
	// if _, err := s.userRepo.FindEmail(email); err == nil {
	// 	return nil, errors.New("email already exists")
	// }

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Email:    email,
		Password: hashedPassword,
		Username: username,
		Role:     "client", // default role
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*entities.User, string, error) {
	user, err := s.userRepo.FindEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(uint(user.ID), user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}