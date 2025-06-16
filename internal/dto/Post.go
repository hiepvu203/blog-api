package dto

type CreatePostRequest struct {
	Title      string `json:"title" binding:"required,min=5,max=200"`
	Slug       string `json:"slug" binding:"required,alphanum"`
	Content    string `json:"content" binding:"required"`
	Thumbnail  string `json:"thumbnail" binding:"omitempty,url"`
	CategoryID uint   `json:"category_id" binding:"required"`
	Status     string `json:"status" binding:"required,oneof=draft published archived"`
}

type UpdatePostRequest struct {
	Title      string `json:"title" binding:"omitempty,min=5,max=200"`
	Slug       string `json:"slug" binding:"omitempty,alphanum"`
	Content    string `json:"content" binding:"omitempty"`
	Thumbnail  string `json:"thumbnail" binding:"omitempty,url"`
	CategoryID uint   `json:"category_id" binding:"omitempty"`
	Status     string `json:"status" binding:"omitempty,oneof=draft published archived"`
}