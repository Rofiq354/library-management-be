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
	ID         uint   `json:"ID"`
	Name       string `json:"name"`
	TotalBooks int64  `json:"total_books"`
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
			ID:         c.ID,
			Name:       c.Name,
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
			"error":   true,
			"message": "Kategori tidak ditemukan",
		})
		return
	}

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
			"error":   true,
			"message": "Data yang dikirim tidak valid. Mohon periksa kembali inputan Anda.",
		})
		return
	}

	var existing models.Category
	if err := h.DB.Where("name = ?", input.Name).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "Kategori dengan nama tersebut sudah ada",
		})
		return
	}

	if err := h.DB.Create(&input).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menambahkan kategori. Silakan coba lagi.",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Kategori berhasil ditambahkan",
		"data":    input,
	})
}

// PUT /admin/categories/:id
func (h *AuthHandler) UpdateCategory(ctx *gin.Context) {
	var category models.Category
	id := ctx.Param("id")

	if err := h.DB.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Kategori tidak ditemukan",
		})
		return
	}

	var input models.Category
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Data yang dikirim tidak valid. Mohon periksa kembali inputan Anda.",
		})
		return
	}

	var existing models.Category
	if err := h.DB.Where("name = ? AND id != ?", input.Name, id).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "Kategori dengan nama tersebut sudah ada",
		})
		return
	}

	if err := h.DB.Model(&category).Updates(input).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal memperbarui kategori. Silakan coba lagi.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil diperbarui",
		"data":    category,
	})
}

// DELETE /admin/categories/:id
func (h *AuthHandler) DeleteCategory(ctx *gin.Context) {
	var category models.Category
	id := ctx.Param("id")

	if err := h.DB.First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Kategori tidak ditemukan",
		})
		return
	}

	var bookCount int64
	h.DB.Model(&models.Book{}).Where("category_id = ?", id).Count(&bookCount)
	if bookCount > 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "Kategori tidak dapat dihapus karena masih memiliki buku yang terkait",
		})
		return
	}

	if err := h.DB.Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menghapus kategori. Silakan coba lagi.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil dihapus",
	})
}