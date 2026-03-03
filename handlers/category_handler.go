package handlers

import (
	"learn-golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ── Custom response structs ────────────────────────
type BookInCategory struct {
	ID            uint   `json:"ID"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Description   string `json:"description"`
	CoverURL      string `json:"cover_url"`
	Stock         int    `json:"stock"`
	PublishedYear int    `json:"published_year"`
}

type CategoryDetailResponse struct {
	ID    uint             `json:"ID"`
	Name  string           `json:"name"`
	Books []BookInCategory `json:"books"`
}

type CategoryListResponse struct {
	ID          uint   `json:"ID"`
	Name        string `json:"name"`
	TotalBooks  int64  `json:"total_books"`
}

// GET /categories
func (h *AuthHandler) GetAllCategories(ctx *gin.Context) {
	var categories []models.Category

	h.DB.Find(&categories)

	response := make([]CategoryListResponse, len(categories))
	for i, c := range categories {
		var count int64
		h.DB.Model(&models.Book{}).Where("category_id = ?", c.ID).Count(&count)
		response[i] = CategoryListResponse{
			ID: c.ID,
			Name: c.Name,
			TotalBooks: count,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GET /categories/:id
func (h *AuthHandler) GetCategoryByID(ctx *gin.Context) {
	var category models.Category
	id := ctx.Param("id")

	if err := h.DB.Preload("Books").First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	// Map ke custom response — buang category_id & nested category
	books := make([]BookInCategory, len(category.Books))
	for i, b := range category.Books {
		books[i] = BookInCategory{
			ID:            b.ID,
			Title:         b.Title,
			Author:        b.Author,
			Description:   b.Description,
			CoverURL:      b.CoverURL,
			Stock:         b.Stock,
			PublishedYear: b.PublishedYear,
		}
	}

	response := CategoryDetailResponse{
		ID:    category.ID,
		Name:  category.Name,
		Books: books,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// POST /admin/categories
func (h *AuthHandler) CreateCategory(ctx *gin.Context) {
	var input models.Category

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Cek duplikat nama category
	var existing models.Category
	if err := h.DB.Where("name = ?", input.Name).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Category with this name already exists",
		})
		return
	}

	h.DB.Create(&input)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data":    input,
	})
}

// PUT /admin/categories/:id
func (h *AuthHandler) UpdateCategory(ctx *gin.Context) {
	var category models.Category
	id := ctx.Param("id")

	if err := h.DB.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	var input models.Category
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Cek duplikat nama (kecuali dirinya sendiri)
	var existing models.Category
	if err := h.DB.Where("name = ? AND id != ?", input.Name, id).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Category with this name already exists",
		})
		return
	}

	h.DB.Model(&category).Updates(input)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data":    category,
	})
}

// DELETE /admin/categories/:id
func (h *AuthHandler) DeleteCategory(ctx *gin.Context) {
	var category models.Category
	id := ctx.Param("id")

	if err := h.DB.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	// Cek apakah category masih dipakai oleh buku
	var bookCount int64
	h.DB.Model(&models.Book{}).Where("category_id = ?", id).Count(&bookCount)
	if bookCount > 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Cannot delete category that still has books assigned to it",
		})
		return
	}

	h.DB.Delete(&category)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}