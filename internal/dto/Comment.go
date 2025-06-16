package dto

type CreateCommentRequest struct {
	PostID  uint   `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1"`
}

