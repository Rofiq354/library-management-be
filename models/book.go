package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model

	Title           string   `json:"title"`
	Author          string   `json:"author"`
	Description     string   `json:"description"`
	CoverURL        string   `json:"cover_url"`
	PDFUrl          string   `json:"pdf_url"`       
	CoverPublicID   string   `json:"cover_public_id"`
	PDFPublicID     string   `json:"pdf_public_id"`
	Stock           int      `json:"stock"`
	PublishedYear   int      `json:"published_year"`
	CategoryID      uint     `json:"category_id"`

	// Relasi ke Category
	Category        Category `json:"category"`
}

