package dto

type CreatePostRequest struct {
	Title      string `json:"title" binding:"required,min=5,max=200"`
	Slug       string `json:"slug" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Thumbnail  string `json:"thumbnail" binding:"omitempty,url"`
	CategoryID uint   `json:"category_id" binding:"required"`
	Status     string `json:"status" binding:"required,oneof=draft published archived"`
}

type UpdatePostRequest struct {
	Title      string `json:"title" binding:"omitempty,min=5,max=200"`
	Slug       string `json:"slug" binding:"omitempty"`
	Content    string `json:"content" binding:"omitempty"`
	Thumbnail  string `json:"thumbnail" binding:"omitempty,url"`
	CategoryID uint   `json:"category_id" binding:"omitempty"`
	Status     string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type PostResponse struct {
    ID         uint   `json:"id"`
    Title      string `json:"title"`
    Slug       string `json:"slug"`
    Content    string `json:"content"`
    Thumbnail  string `json:"thumbnail"`
    CategoryID uint   `json:"category_id"`
    Category   string `json:"category"` // tên category
    AuthorID   uint   `json:"author_id"`
    Author     string `json:"author"`   // tên tác giả
    Status     string `json:"status"`
    Views      int    `json:"views"`
    CreatedAt  string `json:"created_at"`
    UpdatedAt  string `json:"updated_at"`
}