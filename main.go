package main

import (
	"assignment26/config"
	"assignment26/models"
	"assignment26/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load Env
	err := godotenv.Load()

	if err != nil {
		panic("Error loading ENV")
	}

	db := config.ConnectDatabase()
	db.AutoMigrate(&models.Produk{})
	db.AutoMigrate(&models.Inventaris{})
	db.AutoMigrate(&models.Pesanan{})

	router := gin.Default()
	db.AutoMigrate(&models.Produk{})
	db.AutoMigrate(&models.Inventaris{})
	db.AutoMigrate(&models.Pesanan{})
	routes.SetUpRoutes(router, db)
	router.Run(":8080")
}
