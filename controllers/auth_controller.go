package controllers

import (
	"net/http"
	"strings"

	"expense-tracker/config"
	"expense-tracker/middleware"
	"expense-tracker/models"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"` // Email atau Username
	Password   string `json:"password" binding:"required"`
}

// Register mendaftarkan pengguna baru dengan enkripsi password bcrypt
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data input tidak valid. Username min 3 karakter, email valid, dan password min 6 karakter.",
		})
		return
	}

	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	// Periksa apakah username atau email sudah digunakan
	var existingUser models.User
	if err := config.DB.Where("email = ? OR username = ?", input.Email, input.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email atau Username sudah terdaftar",
		})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
	}

	// Enkripsi password menggunakan bcrypt
	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengenkripsi kata sandi",
		})
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan akun pengguna",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil, silakan login",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}

// Login memvalidasi kredensial dan menerbitkan JWT via HttpOnly Cookie
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email/Username dan Password wajib diisi",
		})
		return
	}

	identifier := strings.TrimSpace(input.Identifier)

	var user models.User
	err := config.DB.Where("email = ? OR username = ?", strings.ToLower(identifier), identifier).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email/Username atau Password tidak sesuai",
		})
		return
	}

	// Validasi password bcrypt
	if !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email/Username atau Password tidak sesuai",
		})
		return
	}

	// Buat token JWT
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat sesi login",
		})
		return
	}

	// Simpan ke HttpOnly Cookie (berlaku 24 jam / 86400 detik)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Logout menghapus cookie auth_token
func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	// Hapus cookie dengan MaxAge = -1
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout berhasil, sesi telah dihapus",
	})
}

// GetProfile mengembalikan data profil pengguna yang sedang login
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna belum terautentikasi"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pengguna tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}

