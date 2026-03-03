package handlers

import (
	"errors"
	"learn-golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GET /superadmin/users
func (h *AuthHandler) GetAllUsers(ctx *gin.Context) {
	currentUserID, _ := ctx.Get("user_id")
	
	var users []models.User

	if err := h.DB.Where("id != ?", currentUserID).Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil data admin. Silakan coba lagi.",
		})
		return
	}

	// Sembunyikan password dari response
	type UserResponse struct {
		ID    uint   `json:"ID"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	response := make([]UserResponse, len(users))
	for i, u := range users {
		response[i] = UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// GET /superadmin/users/:id
func (h *AuthHandler) GetUserByID(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")

	if err := h.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Admin tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// POST /superadmin/users
func (h *AuthHandler) CreateUser(ctx *gin.Context) {
	var input models.User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Data yang dikirim tidak valid. Mohon periksa kembali inputan Anda.",
		})
		return
	}

	// Validasi field wajib
	if input.Name == "" || input.Email == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Nama, email, dan password wajib diisi",
		})
		return
	}

	// Cek duplikat email
	var existing models.User
	err := h.DB.Where("email = ?", input.Email).First(&existing).Error
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "Email sudah digunakan oleh admin lain",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Terjadi kesalahan pada database. Silakan coba lagi.",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal memproses password. Silakan coba lagi.",
		})
		return
	}

	role := input.Role
	if role == "" {
		role = "admin"
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menambahkan admin. Silakan coba lagi.",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": user.Name + " berhasil ditambahkan",
		"data": gin.H{
			"ID":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// PUT /superadmin/users/:id
func (h *AuthHandler) UpdateUser(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")

	if err := h.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Admin tidak ditemukan",
		})
		return
	}

	var input models.User
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Data yang dikirim tidak valid. Mohon periksa kembali inputan Anda.",
		})
		return
	}

	// Cek duplikat email (kecuali dirinya sendiri)
	if input.Email != "" {
		var existing models.User
		if err := h.DB.Where("email = ? AND id != ?", input.Email, id).First(&existing).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   true,
				"message": "Email sudah digunakan oleh admin lain",
			})
			return
		}
	}

	// Hash password baru kalau diisi
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(input.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Gagal memproses password. Silakan coba lagi.",
			})
			return
		}
		input.Password = string(hashedPassword)
	}

	if err := h.DB.Model(&user).Updates(input).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal memperbarui data admin. Silakan coba lagi.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data " + user.Name + " berhasil diperbarui",
		"data": gin.H{
			"ID":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// DELETE /superadmin/users/:id
func (h *AuthHandler) DeleteUser(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")

	if err := h.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Admin tidak ditemukan",
		})
		return
	}

	// Cegah hapus diri sendiri
	currentUserID, exists := ctx.Get("user_id")
	if exists && currentUserID == user.ID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "Anda tidak dapat menghapus akun Anda sendiri",
		})
		return
	}

	if err := h.DB.Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menghapus admin. Silakan coba lagi.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Admin " + user.Name + " berhasil dihapus",
	})
}