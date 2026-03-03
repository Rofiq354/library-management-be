package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DBHost string
	DBUser string
	DBPass string
	DBName string
	JWTSecret string
	FrontendOrigin string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:   os.Getenv("PORT"),
		DBHost: os.Getenv("DB_HOST"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		FrontendOrigin: os.Getenv("FRONTEND_ORIGIN"),
	}, nil
}