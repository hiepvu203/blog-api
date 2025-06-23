package dto

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Slug string `json:"slug" binding:"required,min=3,max=50,slug"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Slug string `json:"slug" binding:"required,min=3,max=50,slug"`
}

type CategoryResponse struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}