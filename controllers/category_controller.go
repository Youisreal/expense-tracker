package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"expense-tracker/config"
	"expense-tracker/models"

	"github.com/gin-gonic/gin"
)

type CategoryInput struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
	Type string `json:"type" binding:"required,oneof=income expense"`
}

type CategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default"`
	UserID    *uint  `json:"user_id"`
}

// GetCategories mengambil daftar kategori default sistem dan kategori custom milik pengguna
func GetCategories(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	typeFilter := strings.TrimSpace(c.Query("type"))

	query := config.DB.Where("user_id IS NULL OR user_id = ?", userID)
	if typeFilter == "income" || typeFilter == "expense" {
		query = query.Where("type = ?", typeFilter)
	}

	var categories []models.Category
	if err := query.Order("type ASC, name ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar kategori"})
		return
	}

	var response []CategoryResponse
	for _, cat := range categories {
		response = append(response, CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			Type:      cat.Type,
			IsDefault: cat.IsDefault(),
			UserID:    cat.UserID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// CreateCategory membuat kategori custom baru khusus untuk pengguna yang sedang login
func CreateCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid. Nama wajib diisi (2-50 karakter) dan tipe harus 'income' atau 'expense'.",
		})
		return
	}

	input.Name = strings.TrimSpace(input.Name)

	// Cek apakah kategori dengan nama dan tipe yang sama sudah ada (baik default maupun custom milik user ini)
	var existing models.Category
	err := config.DB.Where("LOWER(name) = LOWER(?) AND type = ? AND (user_id IS NULL OR user_id = ?)",
		input.Name, input.Type, userID).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Kategori dengan nama dan tipe ini sudah ada",
		})
		return
	}

	category := models.Category{
		Name:   input.Name,
		Type:   input.Type,
		UserID: &userID,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kategori baru"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori custom berhasil ditambahkan",
		"data": CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Type:      category.Type,
			IsDefault: false,
			UserID:    category.UserID,
		},
	})
}

// UpdateCategory memperbarui nama/tipe kategori custom milik pengguna (Kategori bawaan dilarang diubah)
func UpdateCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idParam := c.Param("id")

	catID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid"})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, catID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	// Larangan mengubah kategori sistem bawaan (Opsi A)
	if category.IsDefault() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Kategori bawaan sistem tidak dapat diubah"})
		return
	}

	// Validasi kepemilikan data
	if category.UserID == nil || *category.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: kategori ini bukan milik Anda"})
		return
	}

	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid. Nama wajib diisi dan tipe harus 'income' atau 'expense'."})
		return
	}

	input.Name = strings.TrimSpace(input.Name)

	category.Name = input.Name
	category.Type = input.Type

	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kategori berhasil diperbarui",
		"data": CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			Type:      category.Type,
			IsDefault: false,
			UserID:    category.UserID,
		},
	})
}

// DeleteCategory menghapus kategori custom milik pengguna (Kategori bawaan dilarang dihapus)
func DeleteCategory(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idParam := c.Param("id")

	catID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kategori tidak valid"})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, catID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}

	// Larangan menghapus kategori bawaan sistem (Opsi A)
	if category.IsDefault() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Kategori bawaan sistem tidak dapat dihapus"})
		return
	}

	// Validasi kepemilikan data
	if category.UserID == nil || *category.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak: kategori ini bukan milik Anda"})
		return
	}

	// Cek apakah ada transaksi pengguna yang sedang menggunakan kategori ini
	var transactionCount int64
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND category_id = ?", userID, category.ID).
		Count(&transactionCount)

	if transactionCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Kategori tidak dapat dihapus karena masih digunakan oleh " + strconv.FormatInt(transactionCount, 10) + " transaksi.",
		})
		return
	}

	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kategori custom berhasil dihapus",
	})
}

