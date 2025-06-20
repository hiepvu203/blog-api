package services

import (
	"blog-api/internal/dto"
	"blog-api/internal/entities"
	"blog-api/internal/repositories"
)

type CommentService struct {
	repo *repositories.CommentRepository
}

func NewCommentService(repo *repositories.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(req *dto.CreateCommentRequest, userID uint) error {
	comment := &entities.Comment{
		PostID:  req.PostID,
		UserID:  userID,
		Content: req.Content,
	}
	return s.repo.Create(comment)
}

func (s *CommentService) UpdateComment(id uint, content string) error {
    return s.repo.Update(id, content)
}

func (s *CommentService) DeleteComment(id uint) error {
    return s.repo.Delete(id)
}

func (s *CommentService) GetCommentsByPostID(postID uint, page, pageSize int) ([]entities.Comment, int64, error) {
    return s.repo.ListByPostID(postID, page, pageSize)
}