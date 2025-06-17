package services

import (
	"blog-api/internal/dto"
	"blog-api/internal/entities"
	"blog-api/internal/repositories"
)

type PostService struct {
	repo *repositories.PostRepository
}

func NewPostService(repo *repositories.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(req *dto.CreatePostRequest, authorID uint) error {
	post := &entities.Post{
		Title:      req.Title,
		Slug:       req.Slug,
		Content:    req.Content,
		Thumbnail:  req.Thumbnail,
		CategoryID: req.CategoryID,
		AuthorID:   authorID,
		Status:     req.Status,
	}
	return s.repo.Create(post)
}

func (s *PostService) UpdatePost(id uint, req *dto.UpdatePostRequest) error {
    updated := &entities.Post{}
    if req.Title != "" {
        updated.Title = req.Title
    }
    if req.Slug != "" {
        updated.Slug = req.Slug
    }
    if req.Content != "" {
        updated.Content = req.Content
    }
    if req.Thumbnail != "" {
        updated.Thumbnail = req.Thumbnail
    }
    if req.CategoryID != 0 {
        updated.CategoryID = req.CategoryID
    }
    if req.Status != "" {
        updated.Status = req.Status
    }
    return s.repo.Update(id, updated)
}

func (s *PostService) DeletePost(id uint) error {
    return s.repo.Delete(id)
}

func (s *PostService) GetAllPosts() ([]entities.Post, error) {
    return s.repo.ListAll()
}

func (s *PostService) GetPostByID(id uint) (*entities.Post, error) {
    return s.repo.FindByID(id)
}