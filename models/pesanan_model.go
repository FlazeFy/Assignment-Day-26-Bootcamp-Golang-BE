package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Pesanan struct {
		gorm.Model
		Jumlah       int       `json:"jumlah" gorm:"not null"`
		TanggalPesan time.Time `json:"tanggal_pesanan" gorm:"type:date;not null"`
		// FK - Product
		ProdukID uint   `json:"produk_id"`
		Produk   Produk `json:"produk" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
	PesananResponse struct {
		ID           uint      `json:"ID"`
		CreatedAt    time.Time `json:"CreatedAt"`
		UpdatedAt    time.Time `json:"UpdatedAt"`
		Jumlah       int       `json:"jumlah"`
		TanggalPesan time.Time `json:"tanggal_pesanan"`
		ProdukID     uint      `json:"produk_id"`
	}
)
