package main

import (
	"log"

	"learn-golang/config"
	"learn-golang/handlers"
	"learn-golang/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)

	router := routes.SetupRouter(db, authHandler, cfg)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}