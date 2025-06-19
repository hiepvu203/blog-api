package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
    service *services.CategoryService
}

func NewCategoryController(service *services.CategoryService) *CategoryController {
    return &CategoryController{service: service}
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
    var req dto.CreateCategoryRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }
    if err := c.service.CreateCategory(&req); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }
    ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"message": utils.MsgCategoryCreated}))
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidCategoryID))
        return
    }

    var req dto.UpdateCategoryRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    if err := c.service.UpdateCategory(uint(id), &req); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, utils.ErrorResponse(utils.ErrCategoryNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgCategoryUpdated}))
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(utils.ErrInvalidCategoryID))
        return
    }

    err = c.service.DeleteCategory(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse(utils.ErrCategoryNotFound))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgCategoryDeleted}))
}

func (c *CategoryController) ListCategories(ctx *gin.Context) {
    categories, err := c.service.GetAllCategories()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCouldNotFetchCategories))
        return
    }
    var resp []dto.CategoryResponse
	for _, cat := range categories {
		resp = append(resp, dto.CategoryResponse{
			ID: cat.ID,
			Name: cat.Name,
			Slug: cat.Slug,
		})
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse(resp))
}

func (c *CategoryController) AdminListCategories(ctx *gin.Context) {
    categories, err := c.service.GetAllCategories()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(utils.ErrCouldNotFetchCategories))
        return
    }
    var resp []dto.AdminCategoryResponse
    for _, cat := range categories {
        var posts []dto.AdminCategoryPost
        for _, post := range cat.Posts {
            posts = append(posts, dto.AdminCategoryPost{
                ID:    post.ID,
                Title: post.Title,
            })
        }
        resp = append(resp, dto.AdminCategoryResponse{
            ID:        cat.ID,
            Name:      cat.Name,
            Slug:      cat.Slug,
            PostCount: len(posts),
            Posts:     posts,
        })
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(resp))
}