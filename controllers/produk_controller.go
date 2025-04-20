package controllers

import (
	"assignment26/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProdukController struct {
	DB *gorm.DB
}

func NewProdukController(db *gorm.DB) *ProdukController {
	return &ProdukController{DB: db}
}

// Queries
func (c *ProdukController) GetAllProduk(ctx *gin.Context) {
	// Models
	var produk []models.Produk

	// Query
	c.DB.Find(&produk)

	// Response
	status := http.StatusNotFound
	var data interface{} = nil

	if len(produk) > 0 {
		status = http.StatusOK
		data = produk
	}

	ctx.JSON(status, gin.H{
		"data":    data,
		"message": "produk fetched",
	})
}

func (c *ProdukController) GetProdukById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Model
	var produk models.Produk

	// Query
	result := c.DB.First(&produk, "id = ? AND deleted_at IS NULL", id)

	// Response
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"data":    nil,
			"message": "produk not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    produk,
		"message": "produk fetched",
	})
}

func (c *ProdukController) GetProdukByKategori(ctx *gin.Context) {
	// Params
	kategori := ctx.Param("kategori")

	// Model
	var produk []models.Produk

	// Query
	result := c.DB.Where("kategori = ? AND deleted_at IS NULL", kategori).Find(&produk)

	// Response
	status := http.StatusNotFound
	var data interface{} = nil

	if result.Error == nil && len(produk) > 0 {
		status = http.StatusOK
		data = produk
	}

	ctx.JSON(status, gin.H{
		"data":    data,
		"message": "produk by kategori fetched",
	})
}

// Download
func (c *ProdukController) GetDownloadProdukImage(ctx *gin.Context) {
	id := ctx.Param("id")

	// Query : Get Produk
	var produk models.Produk
	if err := c.DB.First(&produk, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "produk not found",
		})
		return
	}

	// Check if image exists
	if produk.ProdukImage == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "produk image not available",
		})
		return
	}

	// Check if file exists
	imageURL := *produk.ProdukImage
	relativePath := strings.TrimPrefix(imageURL, fmt.Sprintf("http://%s/", ctx.Request.Host))
	if _, err := os.Stat(relativePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "file not found on server",
			"path":    relativePath,
		})
		return
	}

	// Response
	ctx.FileAttachment(relativePath, filepath.Base(relativePath))
}

// Commands
func (c *ProdukController) CreateProduk(ctx *gin.Context) {
	// Rules
	allowedTypes := []string{"image/jpg", "image/png", "image/jpeg"}
	maxSize := 5

	// Request Body
	name := ctx.Request.FormValue("name")
	deskripsi := ctx.Request.FormValue("deskripsi")
	hargaStr := ctx.Request.FormValue("harga")
	kategori := ctx.Request.FormValue("kategori")
	jumlahStr := ctx.Request.FormValue("jumlah")
	lokasi := ctx.Request.FormValue("lokasi")

	// Parse Int from form
	harga, err := strconv.Atoi(hargaStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "harga must be number"})
		return
	}
	jumlah, err := strconv.Atoi(jumlahStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "jumlah must be number"})
		return
	}

	// Produk Image Upload
	var produkImage *string
	file, err := ctx.FormFile("produk_image")
	if err == nil {
		// Validate File Size
		if file.Size > int64(maxSize)*1024*1024 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("file size must be less than %dMB", maxSize)})
			return
		}

		// Open File & Validate File Type
		fileHeader, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to open file"})
			return
		}
		defer fileHeader.Close()
		buffer := make([]byte, 512)
		_, err = fileHeader.Read(buffer)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read file"})
			return
		}
		fileType := http.DetectContentType(buffer)
		valid := false
		for _, t := range allowedTypes {
			if fileType == t {
				valid = true
				break
			}
		}
		if !valid {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("file type must be %s", strings.Join(allowedTypes, ", "))})
			return
		}

		// Save File
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filePath := "storages/" + filename
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save image"})
			return
		}
		url := fmt.Sprintf("http://%s/%s", ctx.Request.Host, filePath)
		produkImage = &url
	}

	// Query : Create Produk
	produk := models.Produk{
		Name:        name,
		Deskripsi:   deskripsi,
		Harga:       harga,
		Kategori:    kategori,
		ProdukImage: produkImage,
	}
	if err := c.DB.Create(&produk).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create produk"})
		return
	}

	// Query : Create Inventaris
	inventaris := models.Inventaris{
		ProdukID: produk.ID,
		Jumlah:   jumlah,
		Lokasi:   lokasi,
	}
	if err := c.DB.Create(&inventaris).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create inventaris"})
		return
	}

	// Response
	imageURL := interface{}(nil)
	if produkImage != nil {
		imageURL = fmt.Sprintf("http://%s/api/produk/produk_image/%d", ctx.Request.Host, produk.ID)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"produk":     produk,
			"inventaris": inventaris,
			"image":      imageURL,
		},
		"message": "produk created",
	})
}

func (c *ProdukController) UpdateProdukById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Model
	var produk models.Produk
	var req models.Produk

	// Validation
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	// Query
	result := c.DB.First(&produk, "id = ? AND deleted_at IS NULL", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "produk not found",
		})
		return
	}

	produk.Name = req.Name
	produk.Deskripsi = req.Deskripsi
	produk.Harga = req.Harga
	produk.Kategori = req.Kategori

	c.DB.Save(&produk)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"data":    produk,
		"message": "produk updated",
	})
}

func (c *ProdukController) SoftDeleteProdukById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Model
	var produk models.Produk

	// Query
	result := c.DB.First(&produk, "id = ? AND deleted_at IS NULL", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "produk not found",
		})
		return
	}
	c.DB.Delete(&produk)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "produk deleted",
	})
}

func (c *ProdukController) HardDeleteProdukById(ctx *gin.Context) {
	// Params
	id := ctx.Param("id")

	// Model
	var produk models.Produk

	// Query
	result := c.DB.Unscoped().First(&produk, "id = ? AND deleted_at IS NOT NULL", id)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "produk not found",
		})
		return
	}
	c.DB.Unscoped().Delete(&produk)

	// Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "produk permanentally deleted",
	})
}
