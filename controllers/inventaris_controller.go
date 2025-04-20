package controllers

import (
	"assignment26/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventarisController struct {
	DB *gorm.DB
}

func NewInventarisController(db *gorm.DB) *InventarisController {
	return &InventarisController{DB: db}
}

// Queries
func (c *InventarisController) GetProdukInventarisByProdukId(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Model
	var produk models.Produk

	// Query
	result := c.DB.Preload("Inventaris").First(&produk, "id = ? AND deleted_at IS NULL", id)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": "produk not found",
		})
		return
	}

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"data":    produk,
		"message": "produk with inventaris fetched",
	})
}

// Commands
func (c *InventarisController) UpdateInventarisById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Model
	var inventaris models.Inventaris
	var req models.Inventaris

	// Validation
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	// Validate : Jumlah is valid
	if req.Jumlah < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid jumlah",
		})
		return
	}

	// Query
	result := c.DB.First(&inventaris, "id = ? AND deleted_at IS NULL", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "inventaris not found",
		})
		return
	}

	inventaris.Jumlah = req.Jumlah

	c.DB.Save(&inventaris)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"data":    inventaris,
		"message": "inventaris updated",
	})
}

func (c *InventarisController) HardDeleteInventarisById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Model
	var inventaris models.Inventaris

	// Query
	result := c.DB.Unscoped().First(&inventaris, "id = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "inventaris not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&inventaris)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "inventaris permanentally deleted",
	})
}
