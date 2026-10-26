package models

import (
	"time"
)

type Transaction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	CategoryID  uint      `gorm:"not null;index" json:"category_id"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Amount      int64     `gorm:"type:bigint;not null" json:"amount"` // Bilangan bulat Rupiah
	Type        string    `gorm:"type:varchar(20);not null" json:"type"` // "income" atau "expense"
	Date        time.Time `gorm:"not null;index" json:"date"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

