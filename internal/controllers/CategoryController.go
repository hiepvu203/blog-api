package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
	service *services.CategoryService
}

func NewCategoryController(service *services.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if validationErrs := utils.BindAndValidate(c, &req); len(validationErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
		return
	}
	if err := cc.service.CreateCategory(&req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("category", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"message": utils.MsgCategoryCreated}))
}

func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	var req dto.UpdateCategoryRequest
	if validationErrs := utils.BindAndValidate(c, &req); len(validationErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": validationErrs})
		return
	}

	id, ok := utils.GetUintIDParam(c, "id", utils.ErrInvalidCategoryID)
	if !ok {
		return
	}

	if err := cc.service.UpdateCategory(uint(id), &req); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("categoryID", utils.ErrCategoryNotFound))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("category", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgCategoryUpdated}))
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
    id, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidCategoryID)
	if !ok {
		return
	}

    err := c.service.DeleteCategory(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse("categoryID",utils.ErrCategoryNotFound))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": utils.MsgCategoryDeleted}))
}

func (c *CategoryController) ListCategories(ctx *gin.Context) {
    categories, err := c.service.GetAllCategories()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("categoryID",utils.ErrCouldNotFetchCategories))
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
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("category",utils.ErrCouldNotFetchCategories))
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