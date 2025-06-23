package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
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
	var req dto.CreateCommentRequest
	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
		return
	}

	uid, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return
	}

	postID, ok := utils.GetUintIDParam(ctx, "post_id", utils.ErrInvalidPostID)
	if !ok {
		return
	}
	req.PostID = postID

	if err := c.service.CreateComment(&req, uint(uid)); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("comment", err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"message": utils.MsgCommentCreated}))
}

func (c *CommentController) UpdateComment(ctx *gin.Context) {
	var req dto.UpdateCommentRequest
	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
		return
	}

	commentID, ok := utils.GetUintIDParam(ctx, "comment_id", utils.ErrInvalidCommentID)
	if !ok {
		return
	}

	if err := c.service.UpdateComment(commentID, req.Content); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("comment", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgCommentUpdated}))
}

func (c *CommentController) DeleteComment(ctx *gin.Context) {
	commentID, ok := utils.GetUintIDParam(ctx, "comment_id", utils.ErrInvalidCommentID)
	if !ok {
		return
	}
	if err := c.service.DeleteComment(commentID); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("comment", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgCommentDeleted}))
}

func (c *CommentController) GetCommentsByPost(ctx *gin.Context) {
	postID, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidPostID)
	if !ok {
		return
	}

	page, pageSize, ok := utils.GetPaginationParams(ctx)
	if !ok {
		return
	}

	comments, total, err := c.service.GetCommentsByPostID(postID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("comment", err.Error()))
		return
	}

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