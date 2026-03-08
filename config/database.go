package config

import (
	"fmt"
	"learn-golang/models"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	var dsn string

	// Cek DATABASE_URL dulu
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		dsn = databaseURL
	} else {
		// Fallback ke individual env vars
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASS")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Category{},
		&models.User{},
		&models.Book{},
		&models.Siswa{},
		&models.ReadingSession{},
	)
	if err != nil {
		return nil, err
	}

	seedAdmin(db)

	return db, nil
}

func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count > 0 {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("superadmin123"), bcrypt.DefaultCost)

	admin := models.User{
		Name:     "Super Admin",
		Email:    "superadmin@library.com",
		Password: string(hashedPassword),
		Role:     "superadmin",
	}

	db.Create(&admin)
	fmt.Println("✅ Superadmin berhasil dibuat!")
	fmt.Println("		Email: superadmin@library.com")
	fmt.Println("		Password: superadmin123")
}