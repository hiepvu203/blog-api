package dto

import "blog-api/internal/entities"

type CreatePostRequest struct {
	Title      string `json:"title" binding:"required,min=2,max=200"`
	Slug       string `json:"slug" binding:"required,slug"`
	Content    string `json:"content" binding:"required"`
	Thumbnail  string `json:"thumbnail" binding:"required,url"`	
	CategoryID uint   `json:"category_id" binding:"required,number"`
	Status     string `json:"status" binding:"required,oneof=draft published"`
}

type UpdatePostRequest struct {
	Title      string `json:"title,omitempty" binding:"required,min=2,max=200"`
	Slug       string `json:"slug" binding:"required"`
	Content    string `json:"content,omitempty" binding:"required"`
	Thumbnail  string `json:"thumbnail,omitempty" binding:"omitempty,url"`
	CategoryID uint   `json:"category_id,omitempty" binding:"required,number"`
	Status     string `json:"status,omitempty" binding:"omitempty,oneof=draft published"`
}

type PostResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Content    string `json:"content"`
	Thumbnail  string `json:"thumbnail"`
	CategoryID uint   `json:"category_id"`
	Category   string `json:"category"`
	AuthorID   uint   `json:"author_id"`
	Author     string `json:"author"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// Convert entities.Post -> dto.PostResponse
func NewPostResponse(p *entities.Post) PostResponse {
    return PostResponse{
        ID:         p.ID,
        Title:      p.Title,
        Slug:       p.Slug,
        Content:    p.Content,
        Thumbnail:  p.Thumbnail,
        CategoryID: p.CategoryID,
        Category:   p.Category.Name,
        AuthorID:   p.AuthorID,
        Author:     p.Author.Username,
        Status:     p.Status,
        CreatedAt:  p.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt:  p.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
}