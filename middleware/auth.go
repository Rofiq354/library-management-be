package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware -> cek token valid (semua role boleh)
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak ditemukan, silakan login dulu",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Format token salah, gunakan: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid atau sudah expired",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid",
			})
			return
		}

		ctx.Set("user_id", claims["user_id"])
		ctx.Set("email", claims["email"])
		ctx.Set("role", claims["role"])

		ctx.Next()
	}
}

// SuperAdminMiddleware -> cek role harus superadmin
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		roleValue, exists := ctx.Get("role")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok || role != "superadmin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Akses ditolak, hanya superadmin yang bisa melakukan ini",
			})
			return
		}

		ctx.Next()
	}
}