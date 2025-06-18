package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"net/http"
	"strconv"

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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}
	authorID, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(utils.ErrUnauthorized))
		return
	}
	uid, ok := authorID.(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrInvalidUserIDType))
		return
	}
	if err := c.service.CreatePost(&req, uint(uid)); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"message": utils.MsgPostCreated}))
}

func (c *PostController) UpdatePost(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidPostID))
        return
    }

    var req dto.UpdatePostRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    if err := c.service.UpdatePost(uint(id), &req); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgPostUpdated}))
}

func (c *PostController) DeletePost(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidPostID))
        return
    }

    if err := c.service.DeletePost(uint(id)); err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse(utils.ErrPostNotFound))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgPostDeleted}))
}

func (c *PostController) GetAllPosts(ctx *gin.Context) {
    title := ctx.Query("title")
    content := ctx.Query("content")
    category := ctx.Query("category")
    author := ctx.Query("author")

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

    posts, total, err := c.service.ListPosts(title, content, category, author, "published", page, pageSize)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCouldNotFetchPosts))
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
    idParam := ctx.Param("post_id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidPostID))
        return
    }
    post, err := c.service.GetPostByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse(utils.ErrPostNotFound))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(dto.NewPostResponse(post)))
}