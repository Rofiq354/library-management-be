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
			"message": "Gagal mengambil data user. Silakan coba lagi.",
		})
		return
	}

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
			"message": "User tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"ID":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// POST /superadmin/users
func (h *AuthHandler) CreateUser(ctx *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		NIS      string `json:"nis"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Data yang dikirim tidak valid. Mohon periksa kembali inputan Anda.",
		})
		return
	}

	if input.Name == "" || input.Email == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Nama, email, dan password wajib diisi",
		})
		return
	}

	role := input.Role
	if role == "" {
		role = "user"
	}

	validRoles := map[string]bool{"admin": true, "superadmin": true, "user": true}
	if !validRoles[role] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Role harus salah satu dari: admin, superadmin, user",
		})
		return
	}

	if role == "user" && input.NIS == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "NIS wajib diisi untuk siswa.",
		})
		return
	}

	var existing models.User
	err := h.DB.Where("email = ?", input.Email).First(&existing).Error
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "Email sudah digunakan",
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

	if role == "user" {
		var existingSiswa models.Siswa
		err := h.DB.Where("nis = ?", input.NIS).First(&existingSiswa).Error
		if err == nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   true,
				"message": "NIS sudah terdaftar",
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

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if role == "user" {
		tx := h.DB.Begin()

		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"message": "Gagal menambahkan user, Silahkan coba lagi.",
			})
			return
		}

		siswa := models.Siswa{
			UserID: user.ID,
			NIS:    input.NIS,
		}

		if err := tx.Create(&siswa).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": true,
				"message": "Gagal menambahkan user, Silahkan coba lagi.",
			})
			return
		}

		tx.Commit()
	} else {
		if err := h.DB.Create(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Gagal menambahkan user, Silahkan coba lagi.",
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": user.Name + " (" + role + ") berhasil ditambahkan",
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
			"message": "User tidak ditemukan",
		})
		return
	}

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Data yang dikirim tidak valid. Mohon periksa kembali inputan Anda.",
		})
		return
	}

	if input.Name == "" || input.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Nama dan email wajib diisi",
		})
		return
	}

	if input.Role != "" {
		validRoles := map[string]bool{"admin": true, "superadmin": true, "user": true}
		if !validRoles[input.Role] {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Role harus salah satu dari: admin, superadmin, user",
			})
			return
		}
	}

	if input.Email != user.Email {
		var existing models.User
		err := h.DB.Where("email = ? AND id != ?", input.Email, id).First(&existing).Error
		if err == nil{
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   true,
				"message": "Email sudah digunakan",
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
	}

	// Update fields
	user.Name = input.Name
	user.Email = input.Email
	if input.Role != "" {
		user.Role = input.Role
	}

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
		user.Password = string(hashedPassword)
	}

	if err := h.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal memperbarui data user. Silakan coba lagi.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
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
			"message": "User tidak ditemukan",
		})
		return
	}

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
			"message": "Gagal menghapus user. Silakan coba lagi.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "User " + user.Name + " berhasil dihapus",
	})
}