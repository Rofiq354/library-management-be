package handlers

import (
	"net/http"
	"time"

	"learn-golang/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
	JWTSecret string
}

func NewAuthHandler(db *gorm.DB, secret string) *AuthHandler {
	return  &AuthHandler{
		DB: db,
		JWTSecret: secret,
	}
}

// LOGIN -------------------------------------------
func (h *AuthHandler) Login(ctx * gin.Context) {
	var input struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Input tidak valid"})
		return
	}

	if input.Email == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Email dan password wajib diisi"})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Email atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Email atau password salah"})	
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email": user.Email,
		"role": user.Role,
		"exp": time.Now().Add(24 * time.Hour).Unix(), // token expired 24 jam
	})

	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Gagal membuat token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error": true, "message": "Login berhasil!",
		"token": tokenString,
		"user": gin.H{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,
			"role": user.Role,
		},
	})
}

// REGISTER -------------------------------------------
func (h *AuthHandler) Register(ctx * gin.Context) {
	var input models.User

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Input tidak valid"})
		return
	}

	if input.Email == "" || input.Password == "" || input.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Semua field wajib diisi"})
		return
	}

	var existing models.User
	if err := h.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Email sudah terdaftar"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Gagal hash password"})
		return
	}

	user := models.User{
		Name: input.Name,
		Email: input.Email,
		Password: string(hashedPassword),
		Role: "admin",
	}

	if err := h.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Gagal menyimpan user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error": true, "message": "Akun admin berhasil dibuat!",
	})
}

// LOGOUT -------------------------------------------
func (h *AuthHandler) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"error": true, "message": "Logout berhasil",
	})
}