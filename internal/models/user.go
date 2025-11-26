package models

import (
	"time"

	"gorm.io/gorm"
)

// Pengguna adalah model untuk pengguna dan autentikasi
type Pengguna struct {
	ID        uint           `json:"id" gorm:"primaryKey;column:id"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;column:email"`
	Password  string         `json:"-" gorm:"not null;column:password"`
	Name      string         `json:"name" gorm:"not null;column:nama"`
	Role      string         `gorm:"type:varchar(50);default:'kasir';column:peran" json:"role"` // owner, kasir, admin_gudang, finance - string dulu, bisa jadi FK nanti
	IsActive  bool           `gorm:"default:true;column:aktif" json:"is_active"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:dibuat_pada"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:diperbarui_pada"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;column:dihapus_pada"`
}

// TableName mengembalikan nama tabel untuk model Pengguna
func (Pengguna) TableName() string {
	return "pengguna"
}

// User adalah alias untuk backward compatibility (akan dihapus nanti)
type User = Pengguna
