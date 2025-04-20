package routes

import (
	"assignment26/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, db *gorm.DB) {
	produkController := controllers.NewProdukController(db)
	inventarisController := controllers.NewInventarisController(db)
	pesananController := controllers.NewPesananController(db)

	api := r.Group("/api")
	{
		produk := api.Group("/produk")
		{
			produk.GET("/", produkController.GetAllProduk)
			produk.GET("/:id", produkController.GetProdukById)
			produk.GET("/produk_image/:id", produkController.GetDownloadProdukImage)
			produk.GET("/kategori/:kategori", produkController.GetProdukByKategori)
			produk.POST("/", produkController.CreateProduk)
			produk.PUT("/:id", produkController.UpdateProdukById)
			produk.DELETE("/delete/:id", produkController.SoftDeleteProdukById)
			produk.DELETE("/destroy/:id", produkController.HardDeleteProdukById)
		}
		inventaris := api.Group("/inventaris")
		{
			inventaris.GET("/:id", inventarisController.GetProdukInventarisByProdukId)
			inventaris.PUT("/:id", inventarisController.UpdateInventarisById)
			inventaris.DELETE("/destroy/:id", inventarisController.HardDeleteInventarisById)
		}
		pesanan := api.Group("/pesanan")
		{
			pesanan.GET("/:id", pesananController.GetDetailPesananById)
			pesanan.POST("/:produk_id", pesananController.CreatePesanan)
			pesanan.DELETE("/destroy/:id", pesananController.HardDeletePesananById)
		}
	}
}
