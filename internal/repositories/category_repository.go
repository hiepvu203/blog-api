package repositories

import (
	"blog-api/internal/entities"
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
    err := r.db.Find(&categories).Error
    return categories, err
}