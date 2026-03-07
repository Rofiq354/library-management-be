package routes

import (
	"learn-golang/config"
	"learn-golang/handlers"
	"learn-golang/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(
	db *gorm.DB,
	authhandler *handlers.AuthHandler,
	cfg *config.Config,
) *gin.Engine {

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(gin.Logger(), gin.Recovery())

	// Initialize Reading History Handler
	readingHistoryHandler := handlers.NewReadingHistoryHandler(db)

	api := router.Group("/api")
	{
		// PUBLIC - HANYA LOGIN
		api.POST("/login", authhandler.Login)

		books := api.Group("/books")
		{
			books.GET("", authhandler.GetAllBooks)
			books.GET("/:id", authhandler.GetBookByID)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", authhandler.GetAllCategories)
			categories.GET("/:id", authhandler.GetCategoryByID)
		}

		// PROTECTED
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			protected.POST("/logout", authhandler.Logout)

			// READING HISTORY - Untuk semua user yang authenticated
			reading := protected.Group("/reading")
			{
				reading.POST("/start", readingHistoryHandler.StartReading)
				reading.POST("/:session_id/progress", readingHistoryHandler.UpdateProgress)
				reading.GET("/history", readingHistoryHandler.GetHistory)
				reading.GET("/:session_id", readingHistoryHandler.GetSessionDetail)
				reading.DELETE("/:session_id", readingHistoryHandler.DeleteHistory)
			}

			admin := protected.Group("/admin")
			{
				admin.POST("/books", authhandler.CreateBook)
				admin.PUT("/books/:id", authhandler.UpdateBook)
				admin.DELETE("/books/:id", authhandler.DeleteBook)

				admin.POST("/categories", authhandler.CreateCategory)
				admin.PUT("/categories/:id", authhandler.UpdateCategory)
				admin.DELETE("/categories/:id", authhandler.DeleteCategory)

				admin.GET("/siswa", authhandler.GetAllSiswa)
				admin.GET("/siswa/:id", authhandler.GetSiswaByID)
				admin.PUT("/siswa/:id", authhandler.UpdateSiswa)
				admin.DELETE("/siswa/:id", authhandler.DeleteSiswa)
			}

			superadmin := protected.Group("/superadmin")
			superadmin.Use(middleware.SuperAdminMiddleware())
			{
				superadmin.GET("/users", authhandler.GetAllUsers)
				superadmin.GET("/users/:id", authhandler.GetUserByID)
				superadmin.POST("/users", authhandler.CreateUser)
				superadmin.PUT("/users/:id", authhandler.UpdateUser)
				superadmin.DELETE("/users/:id", authhandler.DeleteUser)
			}
		}
	}

	return router
}