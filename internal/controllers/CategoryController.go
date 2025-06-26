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

// CreateCategory godoc
// @Summary Tạo danh mục mới
// @Description Tạo một danh mục mới (chỉ admin)
// @Tags categories
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param   category  body  dto.CreateCategoryRequest  true  "Thông tin danh mục"
// @Success 201 {object} utils.APIResponse "Tạo danh mục thành công"
// @Failure 400 {object} utils.APIResponse "Lỗi xác thực hoặc tạo danh mục"
// @Failure 500 {object} utils.APIResponse "Lỗi server"
// @Router /admin/categories [post]
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if validationErrs := utils.BindAndValidate(c, &req); validationErrs != nil {
		utils.SendFail(c, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
		return
	}
	if err := cc.service.CreateCategory(&req); err != nil {
		utils.SendFail(c, http.StatusInternalServerError, "400", err.Error(), nil)
		return
	}
	utils.SendSuccess(c, http.StatusCreated, "201", utils.MsgCategoryCreated, nil)
}

// UpdateCategory godoc
// @Summary Cập nhật danh mục
// @Description Cập nhật thông tin danh mục (chỉ admin)
// @Tags categories
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param   id        path  int  true  "ID danh mục"
// @Param   category  body  dto.UpdateCategoryRequest  true  "Thông tin cập nhật danh mục"
// @Success 200 {object} utils.APIResponse "Cập nhật thành công"
// @Failure 400 {object} utils.APIResponse "Lỗi xác thực"
// @Failure 404 {object} utils.APIResponse "Không tìm thấy danh mục"
// @Failure 500 {object} utils.APIResponse "Lỗi server"
// @Router /admin/categories/{id} [put]
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	id, ok := utils.GetUintIDParam(c, "id", utils.ErrInvalidCategoryID)
	if !ok {
		return
	}

	var req dto.UpdateCategoryRequest
	if validationErrs := utils.BindAndValidate(c, &req); validationErrs != nil {
		utils.SendFail(c, http.StatusBadRequest, "400", "VALIDATION_FAILED", validationErrs)
		return
	}

	if err := cc.service.UpdateCategory(uint(id), &req); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendFail(c, http.StatusNotFound, "404", utils.ErrCategoryNotFound, nil)
			return
		}
		utils.SendFail(c, http.StatusInternalServerError, "500", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "200", utils.MsgCategoryUpdated, nil)
}

// DeleteCategory godoc
// @Summary Xóa danh mục
// @Description Xóa danh mục theo ID (chỉ admin)
// @Tags categories
// @Security BearerAuth
// @Produce  json
// @Param   id  path  int  true  "ID danh mục"
// @Success 200 {object} utils.APIResponse "Xóa thành công"
// @Failure 404 {object} utils.APIResponse "Không tìm thấy danh mục"
// @Router /admin/categories/{id} [delete]
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id, ok := utils.GetUintIDParam(ctx, "id", utils.ErrInvalidCategoryID)
	if !ok {
		return
	}

	err := c.service.DeleteCategory(uint(id))
	if err != nil {
		utils.SendFail(ctx, http.StatusNotFound, "404", utils.ErrCategoryNotFound, nil)
		return
	}
	utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgCategoryDeleted, nil)
}

// ListCategories godoc
// @Summary Lấy danh sách danh mục
// @Description Lấy danh sách tất cả danh mục (public)
// @Tags categories
// @Produce  json
// @Success 200 {object} utils.APIResponse "Danh sách danh mục"
// @Failure 500 {object} utils.APIResponse "Lỗi server"
// @Router /categories [get]
func (c *CategoryController) ListCategories(ctx *gin.Context) {
	categories, err := c.service.GetAllCategories()
	if err != nil {
		utils.SendFail(ctx, http.StatusInternalServerError, "500", utils.ErrCouldNotFetchCategories, nil)
		return
	}
	var resp []dto.CategoryResponse
	for _, cat := range categories {
		resp = append(resp, dto.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
			Slug: cat.Slug,
		})
	}
	utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgCategoriesFetched, gin.H{"categories": resp})
}

// AdminListCategories godoc
// @Summary Lấy danh sách danh mục (admin)
// @Description Lấy danh sách danh mục, kèm số lượng và danh sách bài viết trong từng danh mục (chỉ admin)
// @Tags categories
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} utils.APIResponse "Danh sách danh mục cho admin"
// @Failure 500 {object} utils.APIResponse "Lỗi server"
// @Router /admin/categories [get]
func (c *CategoryController) AdminListCategories(ctx *gin.Context) {
	categories, err := c.service.GetAllCategories()
	if err != nil {
		utils.SendFail(ctx, http.StatusInternalServerError, "500", utils.ErrCouldNotFetchCategories, nil)
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
	utils.SendSuccess(ctx, http.StatusOK, "200", utils.MsgAdminCategoriesFetched, gin.H{"categories": resp})
}