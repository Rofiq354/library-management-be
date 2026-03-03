package handlers

import (
	"learn-golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /books
func (h *AuthHandler) GetAllBooks(ctx *gin.Context) {
	var books []models.Book

	h.DB.Preload("Category").Find(&books)

	ctx.JSON(http.StatusOK, gin.H{
		"data": books,
	})
}

// GET /books/:id
func (h *AuthHandler) GetBookByID(ctx *gin.Context) {
	var book models.Book
	id := ctx.Param("id")

	if err := h.DB.Preload("Category").First(&book, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Buku tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": book,
	})
}

// POST /admin/books
func (h *AuthHandler) CreateBook(ctx *gin.Context) {
	var input models.Book

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Data yang dikirim tidak valid. Mohon periksa kembali inputan Anda.",
		})
		return
	}

	// Cek apakah category_id valid
	var category models.Category
	if err := h.DB.First(&category, input.CategoryID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Kategori yang dipilih tidak ditemukan",
		})
		return
	}

	if err := h.DB.Create(&input).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menambahkan buku. Silakan coba lagi.",
		})
		return
	}

	// Load category untuk response
	h.DB.Preload("Category").First(&input, input.ID)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Buku berhasil ditambahkan",
		"data":    input,
	})
}

// PUT /admin/books/:id
func (h *AuthHandler) UpdateBook(ctx *gin.Context) {
	var book models.Book
	id := ctx.Param("id")

	if err := h.DB.First(&book, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Buku tidak ditemukan",
		})
		return
	}

	var input models.Book
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Data yang dikirim tidak valid. Mohon periksa kembali inputan Anda.",
		})
		return
	}

	// Cek apakah category_id valid (kalau diisi)
	if input.CategoryID != 0 {
		var category models.Category
		if err := h.DB.First(&category, input.CategoryID).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Kategori yang dipilih tidak ditemukan",
			})
			return
		}
	}

	if err := h.DB.Model(&book).Updates(input).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal memperbarui buku. Silakan coba lagi.",
		})
		return
	}

	// Load ulang dengan relasi category untuk response
	h.DB.Preload("Category").First(&book, book.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Buku berhasil diperbarui",
		"data":    book,
	})
}

// DELETE /admin/books/:id
func (h *AuthHandler) DeleteBook(ctx *gin.Context) {
	var book models.Book
	id := ctx.Param("id")

	if err := h.DB.First(&book, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Buku tidak ditemukan",
		})
		return
	}

	if err := h.DB.Delete(&book).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menghapus buku. Silakan coba lagi.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Buku berhasil dihapus",
	})
}