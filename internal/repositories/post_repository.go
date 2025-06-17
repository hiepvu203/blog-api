package repositories

import (
	"blog-api/internal/entities"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *entities.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) Update(id uint, updated *entities.Post) error {
    return r.db.Model(&entities.Post{}).Where("id = ?", id).Updates(updated).Error
}

func (r *PostRepository) Delete(id uint) error {
    result := r.db.Delete(&entities.Post{}, id)
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return result.Error
}

func (r *PostRepository) ListAll() ([]entities.Post, error) {
    var posts []entities.Post
    err := r.db.Preload("Author").Preload("Category").Find(&posts).Error
    return posts, err
}

func (r *PostRepository) FindByID(id uint) (*entities.Post, error) {
    var post entities.Post
    err := r.db.Preload("Author").Preload("Category").First(&post, id).Error
    if err != nil {
        return nil, err
    }
    return &post, nil
}