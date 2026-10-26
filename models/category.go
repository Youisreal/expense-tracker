package models

import (
	"time"
)

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Type      string    `gorm:"type:varchar(20);not null" json:"type"` // "income" atau "expense"
	UserID    *uint     `gorm:"index" json:"user_id"`                  // null jika merupakan kategori bawaan (Opsi A)
	User      *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsDefault memeriksa apakah kategori merupakan kategori bawaan sistem
func (c *Category) IsDefault() bool {
	return c.UserID == nil
}

