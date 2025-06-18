package dto

type AdminCategoryPost struct {
    ID    uint   `json:"id"`
    Title string `json:"title"`
}

type AdminCategoryResponse struct {
    ID        uint               `json:"id"`
    Name      string             `json:"name"`
    Slug      string             `json:"slug"`
    PostCount int                `json:"post_count"`
    Posts     []AdminCategoryPost `json:"posts"`
}