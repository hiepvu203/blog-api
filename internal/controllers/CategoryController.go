package controllers

import (
	"blog-api/internal/dto"
	"blog-api/internal/services"
	"blog-api/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
    ctx.JSON(http.StatusCreated, utils.SuccessResponse(gin.H{"message": "Category created successfully"}))
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid category id"))
        return
    }

    var req dto.UpdateCategoryRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
        return
    }

    if err := c.service.UpdateCategory(uint(id), &req); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "Category updated successfully"}))
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid category id"))
        return
    }

    err = c.service.DeleteCategory(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, utils.ErrorResponse("Category not found"))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"message": "Category deleted successfully"}))
}

func (c *CategoryController) ListCategories(ctx *gin.Context) {
    categories, err := c.service.GetAllCategories()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse("Could not fetch categories"))
        return
    }
    ctx.JSON(http.StatusOK, utils.SuccessResponse(categories))
}