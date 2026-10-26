package main

import (
	"log"
	"net/http"
	"os"

	"expense-tracker/config"
	"expense-tracker/controllers"
	"expense-tracker/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Memuat konfigurasi environment
	config.LoadEnv(".env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Inisialisasi Database
	if os.Getenv("DATABASE_URL") != "" {
		_, err := config.ConnectDatabase()
		if err != nil {
			log.Fatalf("[FATAL] Gagal inisialisasi database: %v", err)
		}
	} else {
		log.Fatal("[FATAL] DATABASE_URL belum diatur. Buat file .env dari .env.example untuk menghubungkan database.")
	}

	// 3. Inisialisasi Router Gin
	r := gin.Default()

	// 4. File Statis & Halaman HTML Frontend
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/landing.html")
	r.StaticFile("/login", "./static/login.html")
	r.StaticFile("/register", "./static/register.html")
	r.StaticFile("/dashboard", "./static/index.html")

	// Endpoint health check
	r.GET("/api/health", func(c *gin.Context) {
		dbStatus := "terputus"
		if config.DB != nil {
			dbStatus = "terhubung"
		}
		c.JSON(http.StatusOK, gin.H{
			"app":      "Expense Tracker",
			"status":   "running",
			"database": dbStatus,
		})
	})

	// 5. Rute Autentikasi Publik
	authPublic := r.Group("/api/auth")
	{
		authPublic.POST("/register", controllers.Register)
		authPublic.POST("/login", controllers.Login)
		authPublic.POST("/logout", controllers.Logout)
	}

	// 6. Rute Terproteksi (Memerlukan Cookie JWT HttpOnly)
	apiProtected := r.Group("/api")
	apiProtected.Use(middleware.AuthMiddleware())
	{
		// Profil Pengguna
		apiProtected.GET("/auth/me", controllers.GetProfile)

		// Manajemen Kategori
		apiProtected.GET("/categories", controllers.GetCategories)
		apiProtected.POST("/categories", controllers.CreateCategory)
		apiProtected.PUT("/categories/:id", controllers.UpdateCategory)
		apiProtected.DELETE("/categories/:id", controllers.DeleteCategory)

		// Manajemen & Ringkasan Transaksi
		apiProtected.GET("/transactions", controllers.GetTransactions)
		apiProtected.GET("/transactions/summary", controllers.GetTransactionSummary)
		apiProtected.GET("/transactions/:id", controllers.GetTransactionByID)
		apiProtected.POST("/transactions", controllers.CreateTransaction)
		apiProtected.PUT("/transactions/:id", controllers.UpdateTransaction)
		apiProtected.DELETE("/transactions/:id", controllers.DeleteTransaction)
	}

	log.Printf("Server berjalan di port :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
