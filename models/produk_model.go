package models

import (
	"gorm.io/gorm"
)

type (
	Produk struct {
		gorm.Model
		Name        string  `json:"name" gorm:"type:varchar(100);not null"`
		Deskripsi   string  `json:"deskripsi" gorm:"type:text;null"`
		Harga       int     `json:"harga" gorm:"not null"`
		Kategori    string  `json:"kategori" gorm:"type:varchar(50);not null"`
		ProdukImage *string `json:"produk_image" gorm:"type:varchar(500);null"`
		// Relation
		Inventaris []Inventaris `json:"inventaris" gorm:"foreignKey:ProdukID"`
	}
	CreateProduk struct {
		Name      string `json:"name"`
		Deskripsi string `json:"deskripsi"`
		Harga     int    `json:"harga"`
		Kategori  string `json:"kategori"`
		Jumlah    int    `json:"jumlah"`
		Lokasi    string `json:"lokasi"`
	}
)
