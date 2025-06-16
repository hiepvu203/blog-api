package services

import "blog-api/internal/repositories"

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService{
	return &UserService{userRepo: userRepo}
}

