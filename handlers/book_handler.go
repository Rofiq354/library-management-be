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
			"error": "Buku tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": book,
	})
}

// POST /books
func (h *AuthHandler) CreateBook(ctx *gin.Context) {
	var input models.Book

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.DB.Create(&input)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Buku berhasil ditambahkan",
		"data":    input,
	})
}

// PUT /books/:id
func (h *AuthHandler) UpdateBook(ctx *gin.Context) {
	var book models.Book
	id := ctx.Param("id")

	if err := h.DB.First(&book, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Buku tidak ditemukan",
		})
		return
	}

	var input models.Book
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.DB.Model(&book).Updates(input)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Buku berhasil diupdate",
		"data":    book,
	})
}

// DELETE /books/:id
func (h *AuthHandler) DeleteBook(ctx *gin.Context) {
	var book models.Book
	id := ctx.Param("id")

	if err := h.DB.First(&book, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Buku tidak ditemukan",
		})
		return
	}

	h.DB.Delete(&book)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Buku berhasil dihapus",
	})
}