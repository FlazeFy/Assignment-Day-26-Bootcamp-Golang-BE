package models

import (
	"gorm.io/gorm"
)

type (
	Inventaris struct {
		gorm.Model
		Jumlah int    `json:"jumlah" gorm:"not null"`
		Lokasi string `json:"lokasi" gorm:"type:varchar(75);not null"`
		// FK - Product
		ProdukID uint   `json:"produk_id"`
		Produk   Produk `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)
