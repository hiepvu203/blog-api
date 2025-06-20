package services

import (
	"blog-api/internal/dto"
	"blog-api/internal/entities"
	"blog-api/internal/repositories"
	"errors"
)

type CategoryService struct {
    repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
    return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(req *dto.CreateCategoryRequest) error {
    category := &entities.Category{
        Name: req.Name,
        Slug: req.Slug,
    }
    return s.repo.Create(category)
}

func (s *CategoryService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) error {
    
    var updated = &entities.Category{}
    if req.Name != "" {
        updated.Name = req.Name
    }
    if req.Slug != "" {
        updated.Slug = req.Slug
    }
    if updated.Name == "" && updated.Slug == "" {
        return errors.New("no fields to update")
    }
    return s.repo.Update(id, updated)
}

func (s *CategoryService) DeleteCategory(id uint) error {
    return s.repo.Delete(id)
}

func (s *CategoryService) GetAllCategories() ([]entities.Category, error) {
    return s.repo.ListAll()
}