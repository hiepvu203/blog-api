package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	service *services.CommentService
}

func NewCommentController(service *services.CommentService) *CommentController {
	return &CommentController{service: service}
}

func (c *CommentController) CreateComment(ctx *gin.Context) {
    var req struct {
        Content string `json:"content" binding:"required,min=1"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }
    userID, ok := ctx.Get("userID")
    if !ok {
        ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
        return
    }
    uid, ok := userID.(float64)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Invalid userID type"))
        return
    }
    postIDParam := ctx.Param("post_id")
    var postID uint
    _, err := fmt.Sscanf(postIDParam, "%d", &postID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid post_id"))
        return
    }
    commentReq := dto.CreateCommentRequest{
        PostID:  postID,
        Content: req.Content,
    }
    if err := c.service.CreateComment(&commentReq, uint(uid)); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }
    ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"message": "Comment created successfully"}))
}

func (c *CommentController) UpdateComment(ctx *gin.Context) {
    var req struct {
        Content string `json:"content" binding:"required,min=1"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }
    commentIDParam := ctx.Param("comment_id")
    var commentID uint
    _, err := fmt.Sscanf(commentIDParam, "%d", &commentID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid comment_id"))
        return
    }
    if err := c.service.UpdateComment(commentID, req.Content); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "Comment updated successfully"}))
}

func (c *CommentController) DeleteComment(ctx *gin.Context) {
    commentIDParam := ctx.Param("comment_id")
    var commentID uint
    _, err := fmt.Sscanf(commentIDParam, "%d", &commentID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid comment_id"))
        return
    }
    if err := c.service.DeleteComment(commentID); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "Comment deleted successfully"}))
}

func (c *CommentController) GetCommentsByPost(ctx *gin.Context) {
    postIDParam := ctx.Param("post_id")
    var postID uint
    _, err := fmt.Sscanf(postIDParam, "%d", &postID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid post_id"))
        return
    }
    comments, err := c.service.GetCommentsByPostID(postID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }

    // Convert to DTO
    var resp []dto.CommentResponse
    for _, cmt := range comments {
        resp = append(resp, dto.CommentResponse{
            ID:        cmt.ID,
            PostID:    cmt.PostID,
            UserID:    cmt.UserID,
            Content:   cmt.Content,
            CreatedAt: cmt.CreatedAt,
            UpdatedAt: cmt.CreatedAt, // Nếu có trường UpdatedAt thì dùng, còn không thì giữ nguyên
        })
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(resp))
}