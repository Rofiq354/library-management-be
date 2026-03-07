package models

import (
	"time"

	"gorm.io/gorm"
)

type ReadingSession struct {
	gorm.Model

	SiswaID uint `json:"siswa_id" gorm:"index"`
	Siswa Siswa `json:"siswa" gorm:"foreignKey:SiswaID;onDelete:CASCADE"`
	BookID uint `json:"book_id" gorm:"index"`
	Book Book `json:"book" gorm:"foreignKey:BookID;onDelete:CASCADE"`

	StartedAt time.Time `json:"started_at"`
	LastReadAt time.Time `json:"last_read_at"`
	CurrentPage int `json:"current_page"`
	TotalPages int `json:"total_pages"`
	Progress int `json:"progress"`

	IsCompleted bool `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at"`

	Rating *int `json:"rating"`
	Review string `json:"review"`
}

func (ReadingSession) TableName() string {
	return "reading_sessions"
}