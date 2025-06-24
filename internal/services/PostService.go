package services

import (
	"blog-api/internal/dto"
	"blog-api/internal/entities"
	"blog-api/internal/repositories"

	// "blog-api/pkg/utils"
	"errors"
)

type PostService struct {
	repo *repositories.PostRepository
	categoryRepo   *repositories.CategoryRepository
	userRepo     *repositories.UserRepository
}

func NewPostService(repo *repositories.PostRepository, categoryRepo *repositories.CategoryRepository, userRepo *repositories.UserRepository) *PostService {
    return &PostService{repo: repo, categoryRepo: categoryRepo, userRepo: userRepo}
}

func (s *PostService) CategoryExists(id uint) (bool, error) {
	return s.categoryRepo.Exists(id)
}

func (s *PostService) IsSlugExists(slug string) (bool, error) {
	return s.repo.IsSlugExists(slug)
}

func (s *PostService) CreatePost(req *dto.CreatePostRequest, authorID uint) error {
    user, err := s.userRepo.FindByID(authorID)
    if err != nil {
        return errors.New("user not found")
    }

    if !user.CanPost {
        return errors.New("you have been blocked from posting")
    }

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
    if req.CategoryID != 0 {
        exists, err := s.categoryRepo.Exists(req.CategoryID)
        if err != nil {
            return err
        }
        if !exists {
            return errors.New("category does not exist")
        }
    }
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

func (s *PostService) GetPostByID(id uint) (*entities.Post, error) {
    return s.repo.FindByID(id)
}

func (s *PostService) ListPosts(title, content, category, author, status string, page, pageSize int) ([]entities.Post, int64, error) {
    return s.repo.ListPosts(title, content, category, author, status, page, pageSize)
}