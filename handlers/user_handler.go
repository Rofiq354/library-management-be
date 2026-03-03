package handlers

import (
	"errors"
	"learn-golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GET /users
func (h *AuthHandler) GetAllUsers(ctx *gin.Context) {
	var users []models.User

	h.DB.Find(&users)

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

// GET /users/:id
func (h *AuthHandler) GetUserByID(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")

	if err := h.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// POST /users
func (h *AuthHandler) CreateUser(ctx *gin.Context) {
	var input models.User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid",
		})
		return
	}

	if input.Name == "" || input.Email == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Semua field wajib diisi",
		})
		return
	}

	// Cek apakah email sudah ada
	var existing models.User
	err := h.DB.Where("email = ?", input.Email).First(&existing).Error
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Email sudah digunakan",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Terjadi kesalahan pada database",
		})
		return
	}

	// HASH PASSWORD
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

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
			"error": "Gagal menyimpan user",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": user.Name + " berhasil ditambahkan",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// PUT /users/:id
func (h *AuthHandler) UpdateUser(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")

	if err := h.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	var input models.User
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.DB.Model(&user).Updates(input)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User berhasil diupdate",
		"data":    user,
	})
}

// DELETE /users/:id
func (h *AuthHandler) DeleteUser(ctx *gin.Context) {
	var user models.User
	id := ctx.Param("id")

	if err := h.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User tidak ditemukan",
		})
		return
	}

	h.DB.Delete(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User berhasil dihapus",
	})
}