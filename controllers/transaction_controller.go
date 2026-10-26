package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"expense-tracker/config"
	"expense-tracker/models"

	"github.com/gin-gonic/gin"
)

type TransactionInput struct {
	CategoryID  uint   `json:"category_id" binding:"required"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Type        string `json:"type" binding:"required,oneof=income expense"`
	Date        string `json:"date" binding:"required"` // Format: YYYY-MM-DD
	Description string `json:"description"`
}

type CategorySummary struct {
	CategoryID uint   `json:"category_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Total      int64  `json:"total"`
	Percentage float64`json:"percentage"`
}

type SummaryResponse struct {
	TotalIncome       int64             `json:"total_income"`
	TotalExpense      int64             `json:"total_expense"`
	NetBalance        int64             `json:"net_balance"`
	ExpenseCategories []CategorySummary `json:"expense_categories"`
	IncomeCategories  []CategorySummary `json:"income_categories"`
}

// GetTransactions mengambil daftar transaksi pengguna dengan opsi filter
func GetTransactions(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	query := config.DB.Preload("Category").Where("user_id = ?", userID)

	// Filter berdasarkan Tipe (income/expense)
	typeFilter := strings.TrimSpace(c.Query("type"))
	if typeFilter == "income" || typeFilter == "expense" {
		query = query.Where("type = ?", typeFilter)
	}

	// Filter berdasarkan Kategori
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if catID, err := strconv.ParseUint(categoryIDStr, 10, 64); err == nil {
			query = query.Where("category_id = ?", catID)
		}
	}

	// Filter berdasarkan Bulan dan Tahun
	monthStr := c.Query("month")
	yearStr := c.Query("year")
	if monthStr != "" && yearStr != "" {
		month, errM := strconv.Atoi(monthStr)
		year, errY := strconv.Atoi(yearStr)
		if errM == nil && errY == nil && month >= 1 && month <= 12 {
			startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 1, 0)
			query = query.Where("date >= ? AND date < ?", startDate, endDate)
		}
	}

	// Filter rentang tanggal kustom (YYYY-MM-DD)
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			query = query.Where("date >= ?", t)
		}
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Tambah 1 hari agar mencakup hari akhir secara penuh
			query = query.Where("date <= ?", t)
		}
	}

	var transactions []models.Transaction
	if err := query.Order("date DESC, id DESC").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

// GetTransactionByID mengambil detail satu transaksi milik pengguna
func GetTransactionByID(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")

	var transaction models.Transaction
	if err := config.DB.Preload("Category").Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transaction,
	})
}

// CreateTransaction menambahkan catatan transaksi baru
func CreateTransaction(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Input tidak valid. Pastikan nominal > 0, tipe (income/expense), category_id, dan format tanggal (YYYY-MM-DD) sudah benar.",
		})
		return
	}

	// Validasi tanggal
	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal harus YYYY-MM-DD (contoh: 2026-10-26)"})
		return
	}

	// Validasi keberadaan dan akses kategori
	var category models.Category
	err = config.DB.Where("id = ? AND (user_id IS NULL OR user_id = ?)", input.CategoryID, userID).First(&category).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori tidak valid atau tidak dapat diakses"})
		return
	}

	// Pastikan tipe kategori cocok dengan tipe transaksi
	if category.Type != input.Type {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Kategori '%s' bertipe '%s', tidak cocok dengan tipe transaksi '%s'", category.Name, category.Type, input.Type),
		})
		return
	}

	transaction := models.Transaction{
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Type:        input.Type,
		Date:        parsedDate,
		Description: strings.TrimSpace(input.Description),
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi"})
		return
	}

	// Preload kategori untuk respons
	config.DB.Preload("Category").First(&transaction, transaction.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaksi berhasil dicatat",
		"data":    transaction,
	})
}

// UpdateTransaction memperbarui data transaksi milik pengguna
func UpdateTransaction(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")

	var transaction models.Transaction
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid. Periksa nominal, tipe, dan tanggal transaksi."})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal harus YYYY-MM-DD"})
		return
	}

	// Validasi kategori
	var category models.Category
	err = config.DB.Where("id = ? AND (user_id IS NULL OR user_id = ?)", input.CategoryID, userID).First(&category).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kategori tidak valid atau tidak dapat diakses"})
		return
	}

	if category.Type != input.Type {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Kategori '%s' bertipe '%s', tidak cocok dengan tipe transaksi '%s'", category.Name, category.Type, input.Type),
		})
		return
	}

	transaction.CategoryID = input.CategoryID
	transaction.Amount = input.Amount
	transaction.Type = input.Type
	transaction.Date = parsedDate
	transaction.Description = strings.TrimSpace(input.Description)

	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui transaksi"})
		return
	}

	config.DB.Preload("Category").First(&transaction, transaction.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil diperbarui",
		"data":    transaction,
	})
}

// DeleteTransaction menghapus transaksi milik pengguna
func DeleteTransaction(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")

	var transaction models.Transaction
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil dihapus",
	})
}

// GetTransactionSummary menghitung ringkasan total pemasukan, pengeluaran, saldo bersih, dan breakdown kategori
func GetTransactionSummary(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	baseQuery := config.DB.Model(&models.Transaction{}).Where("user_id = ?", userID)

	// Filter bulan & tahun opsional
	monthStr := c.Query("month")
	yearStr := c.Query("year")
	if monthStr != "" && yearStr != "" {
		month, errM := strconv.Atoi(monthStr)
		year, errY := strconv.Atoi(yearStr)
		if errM == nil && errY == nil && month >= 1 && month <= 12 {
			startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 1, 0)
			baseQuery = baseQuery.Where("date >= ? AND date < ?", startDate, endDate)
		}
	}

	// 1. Hitung Total Income
	var totalIncome int64
	var incomeResult struct {
		Sum int64
	}
	baseQuery.Where("type = ?", "income").Select("COALESCE(SUM(amount), 0) as sum").Scan(&incomeResult)
	totalIncome = incomeResult.Sum

	// 2. Hitung Total Expense
	// Re-construct query untuk expense
	expenseQuery := config.DB.Model(&models.Transaction{}).Where("user_id = ?", userID)
	if monthStr != "" && yearStr != "" {
		month, _ := strconv.Atoi(monthStr)
		year, _ := strconv.Atoi(yearStr)
		if month >= 1 && month <= 12 {
			startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 1, 0)
			expenseQuery = expenseQuery.Where("date >= ? AND date < ?", startDate, endDate)
		}
	}
	var totalExpense int64
	var expenseResult struct {
		Sum int64
	}
	expenseQuery.Where("type = ?", "expense").Select("COALESCE(SUM(amount), 0) as sum").Scan(&expenseResult)
	totalExpense = expenseResult.Sum

	// 3. Breakdown per Kategori Pengeluaran
	type BreakdownRaw struct {
		CategoryID uint   `json:"category_id"`
		Name       string `json:"name"`
		Type       string `json:"type"`
		Total      int64  `json:"total"`
	}

	var rawExpenses []BreakdownRaw
	expenseBreakdownQuery := config.DB.Table("transactions").
		Select("transactions.category_id, categories.name, transactions.type, SUM(transactions.amount) as total").
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND transactions.type = ?", userID, "expense")

	if monthStr != "" && yearStr != "" {
		month, _ := strconv.Atoi(monthStr)
		year, _ := strconv.Atoi(yearStr)
		if month >= 1 && month <= 12 {
			startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 1, 0)
			expenseBreakdownQuery = expenseBreakdownQuery.Where("transactions.date >= ? AND transactions.date < ?", startDate, endDate)
		}
	}
	expenseBreakdownQuery.Group("transactions.category_id, categories.name, transactions.type").Order("total DESC").Scan(&rawExpenses)

	var expenseCategories []CategorySummary
	for _, raw := range rawExpenses {
		pct := 0.0
		if totalExpense > 0 {
			pct = float64(raw.Total) / float64(totalExpense) * 100
		}
		expenseCategories = append(expenseCategories, CategorySummary{
			CategoryID: raw.CategoryID,
			Name:       raw.Name,
			Type:       raw.Type,
			Total:      raw.Total,
			Percentage: float64(int(pct*10)) / 10, // 1 desimal
		})
	}

	// 4. Breakdown per Kategori Pemasukan
	var rawIncomes []BreakdownRaw
	incomeBreakdownQuery := config.DB.Table("transactions").
		Select("transactions.category_id, categories.name, transactions.type, SUM(transactions.amount) as total").
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND transactions.type = ?", userID, "income")

	if monthStr != "" && yearStr != "" {
		month, _ := strconv.Atoi(monthStr)
		year, _ := strconv.Atoi(yearStr)
		if month >= 1 && month <= 12 {
			startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 1, 0)
			incomeBreakdownQuery = incomeBreakdownQuery.Where("transactions.date >= ? AND transactions.date < ?", startDate, endDate)
		}
	}
	incomeBreakdownQuery.Group("transactions.category_id, categories.name, transactions.type").Order("total DESC").Scan(&rawIncomes)

	var incomeCategories []CategorySummary
	for _, raw := range rawIncomes {
		pct := 0.0
		if totalIncome > 0 {
			pct = float64(raw.Total) / float64(totalIncome) * 100
		}
		incomeCategories = append(incomeCategories, CategorySummary{
			CategoryID: raw.CategoryID,
			Name:       raw.Name,
			Type:       raw.Type,
			Total:      raw.Total,
			Percentage: float64(int(pct*10)) / 10,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": SummaryResponse{
			TotalIncome:       totalIncome,
			TotalExpense:      totalExpense,
			NetBalance:        totalIncome - totalExpense,
			ExpenseCategories: expenseCategories,
			IncomeCategories:  incomeCategories,
		},
	})
}

