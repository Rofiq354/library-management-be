package handlers

import (
	"errors"
	"learn-golang/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReadingHistoryHandler struct {
	DB *gorm.DB
}

func NewReadingHistoryHandler(db *gorm.DB) *ReadingHistoryHandler {
	return &ReadingHistoryHandler{
		DB: db,
	}
}

// Helper - dapatkan siswa berdasarkan user_id dari JWT
func (h *ReadingHistoryHandler) getSiswaFromContext(ctx *gin.Context) (*models.Siswa, error) {
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		return nil, errors.New("user_id tidak ditemukan di context")
	}

	userID, ok := userIDValue.(float64) // JWT MapClaims convert int to float64
	if !ok {
		return nil, errors.New("user_id format tidak valid")
	}

	var siswa models.Siswa
	if err := h.DB.Where("user_id = ?", uint(userID)).First(&siswa).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("siswa tidak ditemukan")
		}
		return nil, err
	}

	return &siswa, nil
}

// POST /api/reading/start
// Body: { "book_id": 1 }
// Mulai membaca buku
func (h *ReadingHistoryHandler) StartReading(ctx *gin.Context) {
	var input struct {
		BookID uint `json:"bookId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "book_id wajib diisi",
		})
		return
	}

	// Cek siswa
	siswa, err := h.getSiswaFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Cek book exists
	var book models.Book
	if err := h.DB.First(&book, input.BookID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Buku tidak ditemukan",
		})
		return
	}

	// Cek apakah sudah ada session yang belum selesai untuk buku yang sama
	var existingSession models.ReadingSession
	err = h.DB.Where("siswa_id = ? AND book_id = ? AND is_completed = ?", siswa.ID, input.BookID, false).First(&existingSession).Error
	
	if err == nil {
		// Session sudah ada, return yang ada
		ctx.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": "Session membaca sudah ada",
			"data":    existingSession,
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Terjadi kesalahan database",
		})
		return
	}

	// Buat session baru
	now := time.Now()
	session := models.ReadingSession{
		SiswaID:    siswa.ID,
		BookID:     input.BookID,
		StartedAt:  now,
		LastReadAt: now,
		Progress:   0,
		IsCompleted: false,
	}

	if err := h.DB.Create(&session).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal membuat session membaca",
		})
		return
	}

	// Preload relasi
	h.DB.Preload("Siswa.User").Preload("Book.Category").First(&session, session.ID)

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Session membaca dimulai",
		"data":    session,
	})
}

// POST /api/reading/:session_id/progress
// Body: { "current_page": 45, "total_pages": 200 }
// Update progress membaca
func (h *ReadingHistoryHandler) UpdateProgress(ctx *gin.Context) {
	sessionID, err := strconv.ParseUint(ctx.Param("session_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "session_id tidak valid",
		})
		return
	}

	var input struct {
		CurrentPage int `json:"current_page" binding:"required"`
		TotalPages  int `json:"total_pages" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "current_page dan total_pages wajib diisi",
		})
		return
	}

	if input.CurrentPage < 0 || input.TotalPages <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "current_page harus >= 0 dan total_pages harus > 0",
		})
		return
	}

	// Cek session exist
	var session models.ReadingSession
	if err := h.DB.First(&session, sessionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Session tidak ditemukan",
		})
		return
	}

	// Calculate progress
	progress := 0
	if input.TotalPages > 0 {
		progress = (input.CurrentPage * 100) / input.TotalPages
	}

	// Cek apakah sudah selesai
	isCompleted := progress >= 100
	var completedAt *time.Time
	if isCompleted {
		now := time.Now()
		completedAt = &now
	}

	// Update
	if err := h.DB.Model(&session).Updates(map[string]interface{}{
		"current_page": input.CurrentPage,
		"total_pages":  input.TotalPages,
		"progress":     progress,
		"last_read_at": time.Now(),
		"is_completed": isCompleted,
		"completed_at": completedAt,
	}).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal update progress",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Progress berhasil diperbarui",
		"data": gin.H{
			"current_page": input.CurrentPage,
			"total_pages":  input.TotalPages,
			"progress":     progress,
			"is_completed": isCompleted,
		},
	})
}

// GET /api/reading/history
// Ambil semua reading history siswa yang login
func (h *ReadingHistoryHandler) GetHistory(ctx *gin.Context) {
	siswa, err := h.getSiswaFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	var sessions []models.ReadingSession
	if err := h.DB.
		Preload("Book").
		Where("siswa_id = ?", siswa.ID).
		Order("last_read_at DESC").
		Find(&sessions).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil history",
		})
		return
	}

	// Format response
	type HistoryResponse struct {
		ID          uint      `json:"ID"`
		BookID      uint      `json:"bookId"`
		BookTitle   string    `json:"title"`
		BookAuthor  string    `json:"author"`
		BookCover   string    `json:"coverUrl"`
		StartedAt   time.Time `json:"readAt"`
		LastReadAt  time.Time `json:"lastReadAt"`
		CurrentPage int       `json:"currentPage"`
		TotalPages  int       `json:"totalPages"`
		Progress    int       `json:"progress"`
		IsCompleted bool      `json:"isCompleted"`
	}

	response := make([]HistoryResponse, len(sessions))
	for i, s := range sessions {
		response[i] = HistoryResponse{
			ID:          s.ID,
			BookID:      s.BookID,
			BookTitle:   s.Book.Title,
			BookAuthor:  s.Book.Author,
			BookCover:   s.Book.CoverURL,
			StartedAt:   s.StartedAt,
			LastReadAt:  s.LastReadAt,
			CurrentPage: s.CurrentPage,
			TotalPages:  s.TotalPages,
			Progress:    s.Progress,
			IsCompleted: s.IsCompleted,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "History berhasil diambil",
		"data":    response,
	})
}

// GET /api/reading/:session_id
// Ambil detail satu session
func (h *ReadingHistoryHandler) GetSessionDetail(ctx *gin.Context) {
	sessionID, err := strconv.ParseUint(ctx.Param("session_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "session_id tidak valid",
		})
		return
	}

	var session models.ReadingSession
	if err := h.DB.Preload("Book").First(&session, sessionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Session tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Detail session berhasil diambil",
		"data":    session,
	})
}

// DELETE /api/reading/:session_id
// Hapus session dari history
func (h *ReadingHistoryHandler) DeleteHistory(ctx *gin.Context) {
	sessionID, err := strconv.ParseUint(ctx.Param("session_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "session_id tidak valid",
		})
		return
	}

	var session models.ReadingSession
	if err := h.DB.First(&session, sessionID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Session tidak ditemukan",
		})
		return
	}

	if err := h.DB.Delete(&session).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menghapus history",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "History berhasil dihapus",
	})
}