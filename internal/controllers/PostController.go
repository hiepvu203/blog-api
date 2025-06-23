package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"errors"
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
	
	// Step 1: Get DTO validation errors
	allErrors := utils.BindAndValidate(ctx, &req)

	// Step 2: Get user ID from context
	uid, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		return
	}

	// Step 3: Get business logic validation errors from the service
	serviceErrs, err := c.service.ValidatePostCreation(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("database", err.Error()))
		return
	}
	allErrors = append(allErrors, serviceErrs...)

	// Step 4: If there are any errors from steps 1 or 3, return them all.
	if len(allErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": allErrors})
		return
	}

	// Step 5: If all validations passed, create the post.
	if err := c.service.CreatePost(&req, uint(uid)); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("post", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"message": utils.MsgPostCreated}))
}

func (c *PostController) UpdatePost(ctx *gin.Context) {
	id, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidPostID)
	if !ok {
		return
	}

	var req dto.UpdatePostRequest
	if validationErrs := utils.BindAndValidate(ctx, &req); len(validationErrs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs})
		return
	}

	if err := c.service.UpdatePost(uint(id), &req); err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(appErr.Field, appErr.Message))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("post", err.Error()))
		}
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgPostUpdated}))
}

func (c *PostController) DeletePost(ctx *gin.Context) {
	id, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidPostID)
	if !ok {
		return
	}

	if err := c.service.DeletePost(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("post", utils.ErrPostNotFound))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgPostDeleted}))
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
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("posts", utils.ErrCouldNotFetchPosts))
		return
	}
	var resp []dto.PostResponse
	for _, p := range posts {
		resp = append(resp, dto.NewPostResponse(&p))
	}

	if total == 0 {
		ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
			"success": true,
			"data":    resp,
			"total":   0,
			"page":    page,
			"page_size": pageSize,
			"message": utils.NotFoundArticles,
		}))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{
		"success": true,
		"data": resp,
		"total": total,
		"page": page,
		"page_size": pageSize,
		"message": utils.SearchSuccess,
	}))
}

func (c *PostController) GetPostDetail(ctx *gin.Context) {
	id, ok := utils.GetUintIDParam(ctx, "post_id", utils.ErrInvalidPostID)
	if !ok {
		return
	}
	post, err := c.service.GetPostByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ErrorResponse("post", utils.ErrPostNotFound))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse(dto.NewPostResponse(post)))
}