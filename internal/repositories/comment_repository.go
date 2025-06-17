package repositories

import (
	"blog-api/internal/entities"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *entities.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) Update(id uint, content string) error {
    return r.db.Model(&entities.Comment{}).Where("id = ?", id).Update("content", content).Error
}

func (r *CommentRepository) Delete(id uint) error {
    return r.db.Delete(&entities.Comment{}, id).Error
}

func (r *CommentRepository) ListByPostID(postID uint) ([]entities.Comment, error) {
    var comments []entities.Comment
    err := r.db.Where("post_id = ?", postID).Order("created_at asc").Find(&comments).Error
    return comments, err
}