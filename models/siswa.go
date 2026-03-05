package models

import "gorm.io/gorm"

type Siswa struct {
	gorm.Model
	UserID       uint    `json:"user_id" gorm:"uniqueIndex"`
	User         User    `json:"user" gorm:"foreignKey:UserID;onDelete:CASCADE"`
	NIS          string  `json:"nis" gorm:"uniqueIndex"`      
	Kelas        string  `json:"kelas"`                       
	Jurusan      string  `json:"jurusan"`                     
	NoTelepon    string  `json:"no_telepon"`                  
	Alamat       string  `json:"alamat" gorm:"type:text"`    
	TanggalLahir string  `json:"tanggal_lahir"`              
}