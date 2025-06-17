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
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
		return
	}
	uid, ok := authorID.(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Invalid userID type"))
		return
	}
	if err := c.service.CreatePost(&req, uint(uid)); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"message": "Post created successfully"}))
}

func (c *PostController) UpdatePost(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid post id"))
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
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "Post updated successfully"}))
}

func (c *PostController) DeletePost(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid post id"))
        return
    }

    if err := c.service.DeletePost(uint(id)); err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Post not found"))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "Post deleted successfully"}))
}

func (c *PostController) GetAllPosts(ctx *gin.Context) {
    posts, err := c.service.GetAllPosts()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Could not fetch posts"))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(posts))
}

func (c *PostController) GetPostDetail(ctx *gin.Context) {
    idParam := ctx.Param("post_id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid post id"))
        return
    }
    post, err := c.service.GetPostByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Post not found"))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(post))
}