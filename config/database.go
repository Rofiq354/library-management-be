package config

import (
	"fmt"
	"learn-golang/models"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Category{},
		&models.User{},
		&models.Book{},
		&models.Siswa{},
	)
	if err != nil {
		return nil, err
	}

	seedAdmin(db)

	return db, nil
}

func seedAdmin(db * gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count > 0 {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("superadmin123"), bcrypt.DefaultCost)

	admin := models.User{
		Name: "Super Admin",
		Email: "superadmin@library.com",
		Password: string(hashedPassword),
		Role: "superadmin",
	}

	db.Create(&admin)
	fmt.Println("✅ Superadmin berhasil dibuat!")
	fmt.Println("		Email: superadmin@library.com")
	fmt.Println("		Password: superadmin123")
}