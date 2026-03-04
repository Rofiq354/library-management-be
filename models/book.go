package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model

	Title           string   `json:"title"`
	Author          string   `json:"author"`
	Description     string   `json:"description"`
	CoverURL        string   `json:"cover_url"`
	Stock           int      `json:"stock"`
	PublishedYear   int      `json:"published_year"`
	CategoryID      uint     `json:"category_id"`
	PublicID        string   `json:"public_id"` // ← Cloudinary public ID

	// Relasi ke Category
	Category        Category `json:"category"`
}

