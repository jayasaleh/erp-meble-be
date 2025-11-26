package models

import (
	"time"

	"gorm.io/gorm"
)

// Pengguna adalah model untuk pengguna dan autentikasi
type Pengguna struct {
	ID             uint           `json:"id" gorm:"primaryKey;column:id"`
	Email          string         `json:"email" gorm:"uniqueIndex;not null;column:email"`
	Password       string         `json:"-" gorm:"not null;column:password"`
	Nama           string         `json:"nama" gorm:"not null;column:nama"`
	Peran          string         `gorm:"type:varchar(50);default:'kasir';column:peran" json:"peran"` // owner, kasir, admin_gudang, finance - string dulu, bisa jadi FK nanti
	Aktif          bool           `gorm:"default:true;column:aktif" json:"aktif"`
	DibuatPada     time.Time      `json:"dibuat_pada" gorm:"column:dibuat_pada"`
	DiperbaruiPada time.Time      `json:"diperbarui_pada" gorm:"column:diperbarui_pada"`
	DihapusPada    gorm.DeletedAt `json:"-" gorm:"index;column:dihapus_pada"`
}

// TableName mengembalikan nama tabel untuk model Pengguna
func (Pengguna) TableName() string {
	return "pengguna"
}

// User adalah alias untuk backward compatibility (akan dihapus nanti)
type User = Pengguna
