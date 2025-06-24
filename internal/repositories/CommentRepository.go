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
    result := r.db.Delete(&entities.Comment{}, id)
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return result.Error
}

func (r *CommentRepository) ListByPostID(postID uint, page, pageSize int) ([]entities.Comment, int64, error) {
    var comments []entities.Comment
    var total int64

    query := r.db.Model(&entities.Comment{}).Where("post_id = ?", postID)
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }
    offset := (page - 1) * pageSize
    err := query.Order("created_at asc").Limit(pageSize).Offset(offset).Find(&comments).Error
    return comments, total, err
}