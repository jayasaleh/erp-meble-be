package models

import (
	"time"

	"gorm.io/gorm"
)

// User adalah model untuk pengguna dan autentikasi
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Name      string         `json:"name" gorm:"not null"`
	Role      string         `gorm:"type:varchar(50);default:'kasir'" json:"role"` // owner, kasir, admin_gudang, finance - string dulu, bisa jadi FK nanti
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
