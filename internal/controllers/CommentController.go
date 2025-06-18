package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"fmt"
	"net/http"
	"strconv"

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
        ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(utils.ErrUnauthorized))
        return
    }
    uid, ok := userID.(float64)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrInvalidUserIDType))
        return
    }
    postIDParam := ctx.Param("post_id")
    var postID uint
    _, err := fmt.Sscanf(postIDParam, "%d", &postID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidPostID))
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
    ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"message": utils.MsgCommentCreated}))
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
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidCommentID))
        return
    }
    if err := c.service.UpdateComment(commentID, req.Content); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgCommentUpdated}))
}

func (c *CommentController) DeleteComment(ctx *gin.Context) {
    commentIDParam := ctx.Param("comment_id")
    var commentID uint
    _, err := fmt.Sscanf(commentIDParam, "%d", &commentID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidCommentID))
        return
    }
    if err := c.service.DeleteComment(commentID); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgCommentDeleted}))
}

func (c *CommentController) GetCommentsByPost(ctx *gin.Context) {
    postIDParam := ctx.Param("post_id")
    var postID uint
    _, err := fmt.Sscanf(postIDParam, "%d", &postID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidPostID))
        return
    }

    // Lấy page và page_size từ query
    page := 1
    pageSize := 10
    if p := ctx.Query("page"); p != "" {
        if v, err := strconv.Atoi(p); err == nil && v > 0 {
            page = v
        }
    }
    if ps := ctx.Query("page_size"); ps != "" {
        if v, err := strconv.Atoi(ps); err == nil && v > 0 {
            pageSize = v
        }
    }

    comments, total, err := c.service.GetCommentsByPostID(postID, page, pageSize)
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
            UpdatedAt: cmt.CreatedAt,
        })
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
        "comments": resp,
        "total":    total,
        "page":     page,
        "page_size": pageSize,
    }))
}