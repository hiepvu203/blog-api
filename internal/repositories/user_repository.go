package repositories

import (
	"blog-api/internal/entities"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	// log debug
    fmt.Printf("=== DEBUG ===\nEmail: %s\nError: %v\nUser: %+v\n", email, err, user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil // Trả về nil, nil nếu không tìm thấy
    }
	return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*entities.User, error) {
    var user entities.User
    err := r.db.Where("username = ?", username).First(&user).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &user, err
}

func (r *UserRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound){
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *UserRepository) ListAll() ([]entities.User, error) {
	var users []entities.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&entities.User{}, id)
	if result.RowsAffected == 0{
		return errors.New("user not found")
	}

	return result.Error
}

func (r *UserRepository) Update(user  *entities.User) error {
	return r.db.Model(&entities.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
        "username":  user.Username,
        "email":     user.Email,
        "password":  user.Password, // hashed from services
		"role":      user.Role,
    }).Error
}