package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	service *services.PostService
}

func NewPostController(service *services.PostService) *PostController {
	return &PostController{service: service}
}

func (c *PostController) CreatePost(ctx *gin.Context) {
	var req dto.CreatePostRequest
	
	validationErrs := utils.BindAndValidate(ctx, &req)

	uid, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return
	}

	if validationErrs == nil {
		validationErrs = make(map[string]string)
	}

	categoryExists, err := c.service.CategoryExists(req.CategoryID)
	if err != nil {
		utils.SendFail(ctx, http.StatusInternalServerError, "DB_ERROR", "Lỗi kiểm tra danh mục", nil)
		return
	}
	if !categoryExists {
		validationErrs["category_id"] = "category does not exist"
	}

	slugExists, err := c.service.IsSlugExists(req.Slug)
	if err != nil {
		utils.SendFail(ctx, http.StatusInternalServerError, "DB_ERROR", "Lỗi kiểm tra slug", nil)
		return
	}
	if slugExists {
		validationErrs["slug"] = "slug already exists"
	}

	if len(validationErrs) > 0 {
		utils.SendFail(ctx, http.StatusBadRequest, "VALIDATION_F400AILED", "VALIDATION_FAILED", validationErrs)
		return
	}

	// Create the post
	if err := c.service.CreatePost(&req, uint(uid)); err != nil {
		utils.SendFail(ctx, http.StatusInternalServerError, "500", err.Error(), nil)
		return
	}

	utils.SendSuccess(ctx, http.StatusCreated, "201", utils.MsgPostCreated, nil)
}

func (c *PostController) UpdatePost(ctx *gin.Context) {
	id, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidPostID)
	if !ok {
		return
	}

	var req dto.UpdatePostRequest
	if validationErrs := utils.BindAndValidate(ctx, &req); validationErrs != nil {
		utils.SendFail(ctx, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
		return
	}

	if err := c.service.UpdatePost(uint(id), &req); err != nil {
		utils.SendFail(ctx, http.StatusBadRequest, "400", err.Error(), nil)
		return
	}
	utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgPostUpdated, nil)
}

func (c *PostController) DeletePost(ctx *gin.Context) {
	id, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidPostID)
	if !ok {
		return
	}

	if err := c.service.DeletePost(uint(id)); err != nil {
		utils.SendFail(ctx, http.StatusNotFound, "404", utils.ErrPostNotFound, nil)
		return
	}
	utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgPostDeleted, nil)
}

func (c *PostController) GetAllPosts(ctx *gin.Context) {
	title := ctx.Query("title")
	content := ctx.Query("content")
	category := ctx.Query("category")
	author := ctx.Query("author")

	page, pageSize, ok := utils.GetPaginationParams(ctx)
	if !ok {
		return
	}

	posts, total, err := c.service.ListPosts(title, content, category, author, "published", page, pageSize)
	if err != nil {
		utils.SendFail(ctx, http.StatusInternalServerError, "500", utils.ErrCouldNotFetchPosts, nil)
		return
	}
	var resp []dto.PostResponse
	for _, p := range posts {
		resp = append(resp, dto.NewPostResponse(&p))
	}

	meta := gin.H{"page": page, "page_size": pageSize, "total": total}
	data := gin.H{"posts": resp, "meta": meta}
	utils.SendSuccess(ctx, http.StatusOK, "200", utils.SearchSuccess, data)
}

func (c *PostController) GetPostDetail(ctx *gin.Context) {
	id, ok := utils.GetUintIDParam(ctx, "post_id", utils.ErrInvalidPostID)
	if !ok {
		return
	}
	post, err := c.service.GetPostByID(uint(id))
	if err != nil {
		utils.SendFail(ctx, http.StatusNotFound, "404", utils.ErrPostNotFound, nil)
		return
	}
	utils.SendSuccess(ctx, http.StatusOK, "200", "article details successfully retrieved", gin.H{"post": dto.NewPostResponse(post)})
}