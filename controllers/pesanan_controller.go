package controllers

import (
	"assignment26/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PesananController struct {
	DB *gorm.DB
}

func NewPesananController(db *gorm.DB) *PesananController {
	return &PesananController{DB: db}
}

// Queries
func (c *PesananController) GetDetailPesananById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Model
	var pesanan models.Pesanan

	// Query
	result := c.DB.Preload("Produk").First(&pesanan, "id = ? AND deleted_at IS NULL", id)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": "pesanan not found",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"data":    pesanan,
		"message": "pesanan with produk fetched",
	})
}

// Commands
func (c *PesananController) CreatePesanan(ctx *gin.Context) {
	id := ctx.Param("produk_id")

	var req models.Pesanan
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	// Validate : Jumlah is valid
	if req.Jumlah <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid jumlah",
		})
		return
	}

	// Query : Check produk in inventaris
	var inventaris models.Inventaris
	if err := c.DB.Where("produk_id = ? AND deleted_at IS NULL", id).First(&inventaris).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "produk not found in inventaris",
		})
		return
	}

	// Validate : Check produk avaiablity
	if req.Jumlah > inventaris.Jumlah {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "jumlah pesanan melebihi stok",
		})
		return
	}

	// Query : Add Pesanan
	pesanan := models.Pesanan{
		Jumlah:       req.Jumlah,
		TanggalPesan: time.Now(),
		ProdukID:     inventaris.ProdukID,
	}
	if err := c.DB.Create(&pesanan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create pesanan",
		})
		return
	}

	// Query : Update Inventaris
	inventaris.Jumlah -= req.Jumlah
	if err := c.DB.Save(&inventaris).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update inventaris",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": models.PesananResponse{
			ID:           pesanan.ID,
			CreatedAt:    pesanan.CreatedAt,
			UpdatedAt:    pesanan.UpdatedAt,
			Jumlah:       pesanan.Jumlah,
			TanggalPesan: pesanan.TanggalPesan,
			ProdukID:     pesanan.ProdukID,
		},
		"message": "pesanan created",
	})
}

func (c *PesananController) HardDeletePesananById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Model
	var pesanan models.Pesanan

	// Query
	result := c.DB.Unscoped().First(&pesanan, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "pesanan not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&pesanan)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pesanan permanentally deleted",
	})
}
