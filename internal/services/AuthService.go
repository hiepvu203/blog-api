package services

import (
	"blog-api/internal/entities"
	"blog-api/internal/repositories"
	"blog-api/pkg/helper"
	"blog-api/pkg/utils"

	// "blog-api/pkg/helper"
	"errors"
	// "fmt"
	"regexp"
	// "time"
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

	if existing, err := s.userRepo.FindByUsername(username); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    } else if existing != nil {
        return nil, errors.New("username already exists")
    }

	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Email:    email,
		Password: hashedPassword,
		Username: username,
		Role:     "client", 
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*entities.User, string, []string) {
    var errs []string
    user, err := s.userRepo.FindEmail(email)
    if err != nil || user == nil {
        errs = append(errs, "email not found")
        if password == "" || len(password) < 6 {
            errs = append(errs, "password must be at least 6 characters long")
        }else if !isValidEmail(email) { 
			errs = append(errs, "Email format is invalid")
		}
        return nil, "", errs
    }
    if !helper.CheckPasswordHash(password, user.Password) {
        errs = append(errs, "Password is incorrect")
    }
    if len(errs) > 0 {
        return nil, "", errs
    }

    token, err := utils.GenerateToken(uint(user.ID), user.Role)
	if err != nil {
		errs = append(errs, "Failed to generate token")
		return nil, "", errs
	}
    return user, token, nil
}

func isValidEmail(email string) bool {
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

// func (s *AuthService) ForgotPassword(email string) error {
// 	user, err := s.userRepo.FindEmail(email)
//     if err != nil || user == nil {
//         return nil
//     }

//     token, err := utils.GenerateResetToken(uint(user.ID), 15*time.Minute)
//     if err != nil {
//         return err
//     }

//     resetLink := fmt.Sprintf("https://yourdomain.com/reset-password?token=%s", token)
//     return helper.SendResetEmail(email, resetLink)
// }

// func (s *AuthService) ResetPassword(token, newPassword string) error {
// 	userID, err := utils.ValidateResetToken(token)
// 	if err != nil {
// 		return errors.New("invalid or expired token")
// 	}

// 	user, err := s.userRepo.FindByID(userID)
// 	if err != nil || user == nil {
// 		return errors.New("user not found")
// 	}

// 	hashed, err := utils.HashPassword(newPassword)
// 	if err != nil {
// 		return err
// 	}
// 	user.Password = hashed
// 	return s.userRepo.Update(user)
// }