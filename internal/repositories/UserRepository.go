package repositories

import (
	"blog-api/internal/entities"
	"errors"

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
	if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
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
    err := r.db.Preload("Posts").Preload("Comments").First(&user, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound){
        return nil, errors.New("user not found")
    }
    return &user, err
}

func (r *UserRepository) ListAll(page, pageSize int) ([]entities.User, int64, error) {
	var users []entities.User
    var total int64

    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }
    offset := (page - 1) * pageSize

    query := r.db.Model(&entities.User{})
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    err := query.Limit(pageSize).Offset(offset).Order("created_at desc").Find(&users).Error
    return users, total, err
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

func (r *UserRepository) UpdateCanPost(userID uint, canPost bool) error {
    return r.db.Model(&entities.User{}).Where("id = ?", userID).Update("can_post", canPost).Error
}