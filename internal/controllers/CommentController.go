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

// CreateComment godoc
// @Summary Tạo bình luận mới
// @Description Tạo bình luận cho một bài viết (yêu cầu đăng nhập)
// @Tags comments
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param   post_id  path  int  true  "ID bài viết"
// @Param   comment  body  dto.CreateCommentRequest  true  "Nội dung bình luận"
// @Success 201 {object} utils.APIResponse "Tạo bình luận thành công"
// @Failure 400 {object} utils.APIResponse "Lỗi xác thực hoặc dữ liệu không hợp lệ"
// @Failure 500 {object} utils.APIResponse "Lỗi server"
// @Router /posts/{post_id}/comments [post]
func (c *CommentController) CreateComment(ctx *gin.Context) {
	var req dto.CreateCommentRequest
	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
		utils.SendFail(ctx, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
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
		utils.SendFail(ctx, http.StatusInternalServerError, "500", err.Error(), nil)
		return
	}
	utils.SendSuccess(ctx, http.StatusCreated, "201", utils.MsgCategoryCreated, nil)
}

// UpdateComment godoc
// @Summary Cập nhật bình luận
// @Description Cập nhật nội dung bình luận (yêu cầu đăng nhập, chỉ chủ sở hữu mới được sửa)
// @Tags comments
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param   comment_id  path  int  true  "ID bình luận"
// @Param   comment     body  dto.UpdateCommentRequest  true  "Nội dung cập nhật"
// @Success 200 {object} utils.APIResponse "Cập nhật thành công"
// @Failure 400 {object} utils.APIResponse "Lỗi xác thực"
// @Failure 404 {object} utils.APIResponse "Không tìm thấy bình luận"
// @Router /comments/{comment_id} [put]
func (c *CommentController) UpdateComment(ctx *gin.Context) {
	var req dto.UpdateCommentRequest
	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
		utils.SendFail(ctx, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
		return
	}

	commentID, ok := utils.GetUintIDParam(ctx, "comment_id", utils.ErrInvalidCommentID)
	if !ok {
		return
	}

	if err := c.service.UpdateComment(commentID, req.Content); err != nil {
		utils.SendFail(ctx, http.StatusNotFound, "404", utils.ErrCategoryNotFound, nil)
		return
	}
	utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgCategoryUpdated, nil)
}

// DeleteComment godoc
// @Summary Xóa bình luận
// @Description Xóa bình luận (yêu cầu đăng nhập, chỉ chủ sở hữu bình luận hoặc chủ bài viết mới được xóa)
// @Tags comments
// @Security BearerAuth
// @Produce  json
// @Param   comment_id  path  int  true  "ID bình luận"
// @Success 200 {object} utils.APIResponse "Xóa thành công"
// @Failure 404 {object} utils.APIResponse "Không tìm thấy bình luận"
// @Router /comments/{comment_id} [delete]
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	commentID, ok := utils.GetUintIDParam(ctx, "comment_id", utils.ErrInvalidCommentID)
	if !ok {
		return
	}
	if err := c.service.DeleteComment(commentID); err != nil {
		utils.SendFail(ctx, http.StatusNotFound, "404", utils.ErrCategoryNotFound, nil)
		return
	}
	utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgCategoryDeleted, nil)
}

// GetCommentsByPost godoc
// @Summary Lấy danh sách bình luận của bài viết
// @Description Lấy danh sách bình luận theo bài viết, có phân trang
// @Tags comments
// @Produce  json
// @Param   post_id   path  int  true  "ID bài viết"
// @Param   page      query int  false "Trang hiện tại"
// @Param   page_size query int  false "Số lượng mỗi trang"
// @Success 200 {object} utils.APIResponse "Danh sách bình luận"
// @Failure 500 {object} utils.APIResponse "Lỗi server"
// @Router /posts/{post_id}/comments [get]
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
		utils.SendFail(ctx, http.StatusInternalServerError, "500", err.Error(), nil)
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

	utils.SendSuccess(ctx, http.StatusOK, "COMMENTS_FETCHED", "Lấy danh sách bình luận thành công", gin.H{
    "comments": resp,
    "meta": gin.H{
        "total":    total,
        "page":     page,
        "page_size": pageSize,
    },
})
}