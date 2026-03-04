package handlers

import (
	"context"
	"learn-golang/models"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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
	title := ctx.PostForm("title")
	author := ctx.PostForm("author")
	description := ctx.PostForm("description")
	stock := ctx.PostForm("stock")
	publishedYear := ctx.PostForm("published_year")
	categoryIDStr := ctx.PostForm("category_id")

	// Validasi input
	if title == "" || author == "" || categoryIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Field title, author, dan category_id tidak boleh kosong",
		})
		return
	}

	stockInt, _ := strconv.Atoi(stock)
	publishedYearInt, _ := strconv.Atoi(publishedYear)
	categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)

	// Cek category
	var category models.Category
	if err := h.DB.First(&category, categoryID).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Kategori yang dipilih tidak ditemukan",
		})
		return
	}

	// Handle image upload
	coverFile, err := ctx.FormFile("cover")
	if err != nil && err != http.ErrMissingFile {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Gagal membaca file",
		})
		return
	}

	var coverURL string
	var coverPublicID string

	if coverFile != nil {
		// Upload ke Cloudinary
		uploadURL, publicID, err := h.uploadToCloudinary(coverFile, "library/books/covers")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Gagal upload cover ke Cloudinary",
			})
			return
		}
		coverURL = uploadURL
		coverPublicID = publicID
	}

	// Handle PDF upload
	pdfFile, _ := ctx.FormFile("pdf")
	var pdfURL string
	var pdfPublicID string

	if pdfFile != nil {
		uploadURL, publicID, err := h.uploadPDFToCloudinary(pdfFile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Gagal upload PDF ke Cloudinary",
			})
			return
		}
		pdfURL = uploadURL
		pdfPublicID = publicID
	}

	// Create book
	book := models.Book{
		Title:         title,
		Author:        author,
		Description:   description,
		CoverURL:      coverURL,
		PDFUrl:        pdfURL,
		Stock:         stockInt,
		PublishedYear: publishedYearInt,
		CategoryID:    uint(categoryID),
		CoverPublicID: coverPublicID,
		PDFPublicID:   pdfPublicID,
	}

	if err := h.DB.Create(&book).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menambahkan buku",
		})
		return
	}

	// Load category
	h.DB.Preload("Category").First(&book, book.ID)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Buku berhasil ditambahkan",
		"data":    book,
	})
}

// PUT /admin/books/:id
func (h *AuthHandler) UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")

	var book models.Book
	if err := h.DB.First(&book, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Buku tidak ditemukan",
		})
		return
	}

	// Update fields
	if title := ctx.PostForm("title"); title != "" {
		book.Title = title
	}
	if author := ctx.PostForm("author"); author != "" {
		book.Author = author
	}
	if description := ctx.PostForm("description"); description != "" {
		book.Description = description
	}
	if stock := ctx.PostForm("stock"); stock != "" {
		stockInt, _ := strconv.Atoi(stock)
		book.Stock = stockInt
	}

	// Handle new cover image
	coverFile, _ := ctx.FormFile("cover")
	if coverFile != nil {
		// Delete old cover
		if book.CoverPublicID != "" {
			h.Cloudinary.Upload.Destroy(context.Background(), uploader.DestroyParams{
				PublicID: book.CoverPublicID,
			})
		}

		// Upload new cover
		coverURL, publicID, err := h.uploadToCloudinary(coverFile, "library/books/covers")
		if err == nil {
			book.CoverURL = coverURL
			book.CoverPublicID = publicID
		}
	}

	// Handle new PDF
	pdfFile, _ := ctx.FormFile("pdf")
	if pdfFile != nil {
		// Delete old PDF
		if book.PDFPublicID != "" {
			h.Cloudinary.Upload.Destroy(context.Background(), uploader.DestroyParams{
				PublicID: book.PDFPublicID,
			})
		}

		// Upload new PDF
		pdfURL, publicID, err := h.uploadPDFToCloudinary(pdfFile)
		if err == nil {
			book.PDFUrl = pdfURL
			book.PDFPublicID = publicID
		}
	}

	if err := h.DB.Save(&book).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal update buku",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Buku berhasil diupdate",
		"data":    book,
	})
}

// DELETE /admin/books/:id
func (h *AuthHandler) DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")

	var book models.Book
	if err := h.DB.First(&book, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Buku tidak ditemukan",
		})
		return
	}

	// Delete cover image
	if book.CoverPublicID != "" {
		h.Cloudinary.Upload.Destroy(context.Background(), uploader.DestroyParams{
			PublicID: book.CoverPublicID,
		})
	}

	// Delete PDF
	if book.PDFPublicID != "" {
		h.Cloudinary.Upload.Destroy(context.Background(), uploader.DestroyParams{
			PublicID: book.PDFPublicID,
		})
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