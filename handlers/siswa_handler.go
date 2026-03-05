package handlers

import (
	"learn-golang/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) GetAllSiswa(ctx *gin.Context) {
	var siswaList []models.Siswa

	if err := h.DB.Preload("User").Find(&siswaList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil data siswa.",
		})
		return
	}

	type SiswaResponse struct {
		ID					 uint   `json:"ID"`
		UserID       uint   `json:"user_id"`
		Name         string `json:"name"`
		Email        string `json:"email"`
		NIS          string `json:"nis"`
		Kelas        string `json:"kelas"`
		Jurusan      string `json:"jurusan"`
		NoTelepon    string `json:"no_telepon"`
		Alamat       string `json:"alamat"`
		TanggalLahir string `json:"tanggal_lahir"`
	}

	response := make([]SiswaResponse, len(siswaList))
	for i, s := range siswaList {
		response[i] = SiswaResponse{
			ID:           s.ID,
			UserID:       s.UserID,
			Name:         s.User.Name,
			Email:        s.User.Email,
			NIS:          s.NIS,
			Kelas:        s.Kelas,
			Jurusan:      s.Jurusan,
			NoTelepon:    s.NoTelepon,
			Alamat:       s.Alamat,
			TanggalLahir: s.TanggalLahir,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error": false,
		"data": response,
	})
}

func (h *AuthHandler) GetSiswaByID(ctx *gin.Context) {
	var siswa models.Siswa
	id := ctx.Param("id")

	if err := h.DB.Preload("User").First(&siswa, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Siswa tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error": false,
		"data": gin.H{
			"ID":           siswa.ID,
			"user_id":      siswa.UserID,
			"name":         siswa.User.Name,
			"email":        siswa.User.Email,
			"nis":          siswa.NIS,
			"kelas":        siswa.Kelas,
			"jurusan":      siswa.Jurusan,
			"no_telepon":   siswa.NoTelepon,
			"alamat":       siswa.Alamat,
			"tanggal_lahir": siswa.TanggalLahir,
		},
	})
}

func (h *AuthHandler) UpdateSiswa(ctx *gin.Context) {
	var siswa models.Siswa
	id := ctx.Param("id")

	if err := h.DB.Preload("User").First(&siswa, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Siswa tidak ditemukan",
		})
		return
	}

	var input struct {
		NIS          string `json:"nis"`
		Kelas        string `json:"kelas"`
		Jurusan      string `json:"jurusan"`
		NoTelepon    string `json:"no_telepon"`
		Alamat       string `json:"alamat"`
		TanggalLahir string `json:"tanggal_lahir"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Data tidak valid.",
		})
		return
	}

	if input.NIS != "" && input.NIS != siswa.NIS {
		var existing models.Siswa
		if err := h.DB.Where("nis = ? AND id != ?", input.NIS, siswa.ID).First(&existing).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":   true,
				"message": "NIS sudah terdaftar.",
			})
			return
		}
	}

	if input.NIS != "" {
		siswa.NIS = input.NIS
	}
	if input.Kelas != "" {
		siswa.Kelas = input.Kelas
	}
	if input.Jurusan != "" {
		siswa.Jurusan = input.Jurusan
	}
	if input.NoTelepon != "" {
		siswa.NoTelepon = input.NoTelepon
	}
	if input.Alamat != "" {
		siswa.Alamat = input.Alamat
	}
	if input.TanggalLahir != "" {
		siswa.TanggalLahir = input.TanggalLahir
	}

	if err := h.DB.Save(&siswa).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengupdate data siswa.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error": false,
		"message": "Data siswa berhasil diperbarui.",
		"data": gin.H{
			"ID":           siswa.ID,
			"user_id":      siswa.UserID,
			"name":         siswa.User.Name,
			"email":        siswa.User.Email,
			"nis":          siswa.NIS,
			"kelas":        siswa.Kelas,
			"jurusan":      siswa.Jurusan,
			"no_telepon":   siswa.NoTelepon,
			"alamat":       siswa.Alamat,
			"tanggal_lahir": siswa.TanggalLahir,
		},
	})
}

func (h *AuthHandler) DeleteSiswa(ctx *gin.Context) {
	var siswa models.Siswa
	id := ctx.Param("id")

	if err := h.DB.Preload("User").First(&siswa, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Siswa tidak ditemukan",
		})
		return
	}

	if err := h.DB.Delete(&siswa).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menghapus data siswa.",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error": false,
		"message": "Data siswa berhasil dihapus.",
	})
}