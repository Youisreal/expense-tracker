package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"expense-tracker/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// LoadEnv membaca file .env sederhana jika tersedia tanpa dependensi eksternal
func LoadEnv(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		// File .env tidak ditemukan, menggunakan environment variable bawaan sistem
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			// Hanya set jika belum ada di environment OS
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("gagal membaca file .env: %v", err)
	}
}

// ConnectDatabase menginisialisasi koneksi GORM ke PostgreSQL Neon
func ConnectDatabase() (*gorm.DB, error) {
	// Memuat file .env jika ada di direktori kerja
	LoadEnv(".env")

	dbURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(dbURL) == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required but not provided")
	}

	// Inisialisasi GORM dengan log mode Warn untuk mencegah kebocoran data sensitif
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	// Konfigurasi pool koneksi
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan konfigurasi pool database: %w", err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// AutoMigrate tabel User, Category, dan Transaction
	err = db.AutoMigrate(&models.User{}, &models.Category{}, &models.Transaction{})
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan migrasi schema database: %w", err)
	}

	// Seeding kategori default sistem (Opsi A: UserID = nil)
	seedDefaultCategories(db)

	DB = db
	log.Println("Koneksi database PostgreSQL Neon dan migrasi tabel berhasil disiapkan.")
	return db, nil
}

// seedDefaultCategories menyiapkan kategori bawaan sistem dengan user_id = NULL
func seedDefaultCategories(db *gorm.DB) {
	defaultCategories := []models.Category{
		// Pengeluaran
		{Name: "Makanan & Minuman", Type: "expense", UserID: nil},
		{Name: "Tagihan & Utilitas", Type: "expense", UserID: nil},
		{Name: "Transportasi", Type: "expense", UserID: nil},
		{Name: "Hiburan & Rekreasi", Type: "expense", UserID: nil},
		{Name: "Belanja", Type: "expense", UserID: nil},
		{Name: "Kesehatan", Type: "expense", UserID: nil},
		// Pemasukan
		{Name: "Gaji & Upah", Type: "income", UserID: nil},
		{Name: "Bonus & Tunjangan", Type: "income", UserID: nil},
		{Name: "Investasi & Dividen", Type: "income", UserID: nil},
		{Name: "Lain-lain", Type: "income", UserID: nil},
	}

	for _, cat := range defaultCategories {
		var existing models.Category
		err := db.Where("name = ? AND type = ? AND user_id IS NULL", cat.Name, cat.Type).First(&existing).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			if createErr := db.Create(&cat).Error; createErr != nil {
				log.Printf("Peringatan: Gagal menambahkan kategori default %s: %v", cat.Name, createErr)
			}
		}
	}
}
