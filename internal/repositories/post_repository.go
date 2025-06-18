package repositories

import (
	"blog-api/internal/entities"
	// "fmt"

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

func (r *PostRepository) FindByID(id uint) (*entities.Post, error) {
    var post entities.Post
    err := r.db.Preload("Author").Preload("Category").First(&post, id).Error
    if err != nil {
        return nil, err
    }
    return &post, nil
}

func (r *PostRepository) ListPosts(title, content, category, author string, page, pageSize int) ([]entities.Post, int64, error) {
    var posts []entities.Post
    var total int64

    query := r.db.Model(&entities.Post{}).Preload("Author").Preload("Category").Preload("Comments")
    if title != "" {
        query = query.Where("title ILIKE ?", "%"+title+"%")
    }
    if content != "" {
        query = query.Where("content ILIKE ?", "%"+content+"%")
    }
    if category != "" {
        query = query.Joins("JOIN categories ON categories.id = posts.category_id").Where("categories.slug ILIKE ?", "%"+category+"%")
    }
    if author != "" {
        query = query.Joins("JOIN users ON users.id = posts.author_id").Where("users.username ILIKE ?", "%"+author+"%")
    }

    // count result
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // page
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }
    offset := (page - 1) * pageSize

    err := query.Limit(pageSize).Offset(offset).Order("created_at desc").Find(&posts).Error
    return posts, total, err

	// return nil, 0, fmt.Errorf("simulate db error") // test case error
}