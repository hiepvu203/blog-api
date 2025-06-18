package dto

import "time"

type CreateCommentRequest struct {
	PostID  uint   `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CommentResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListCommentResponse struct {
    Comments []CommentResponse `json:"comments"`
    Total    int               `json:"total"`
}