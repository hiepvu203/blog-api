package repositories

import (
	"blog-api/internal/entities"
	// "fmt"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
    return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *entities.Category) error {
    return r.db.Create(category).Error
}

func (r *CategoryRepository) Update(id uint, updated *entities.Category) error {
    return r.db.Model(&entities.Category{}).Where("id = ?", id).Updates(updated).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	result := r.db.Delete((&entities.Category{}), id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *CategoryRepository) ListAll() ([]entities.Category, error) {
    var categories []entities.Category
    err := r.db.Preload("Posts").Find(&categories).Error
    return categories, err

	// return nil, fmt.Errorf("simulate db error")  // case : test error
}

func (r *CategoryRepository) Exists(id uint) (bool, error) {
    var count int64
    err := r.db.Model(&entities.Category{}).Where("id = ?", id).Count(&count).Error
    return count > 0, err
}