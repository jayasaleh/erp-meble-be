package models

import (
	"time"

	"gorm.io/gorm"
)

// Produk adalah model untuk master data barang
type Produk struct {
	ID             uint           `gorm:"primaryKey;column:id" json:"id"`
	SKU            string         `gorm:"uniqueIndex;not null;column:sku" json:"sku"`
	Barcode        string         `gorm:"uniqueIndex;column:barcode" json:"barcode"` // Nullable, unique jika ada
	Nama           string         `gorm:"not null;column:nama" json:"nama"`
	Kategori       string         `gorm:"type:varchar(100);column:kategori" json:"kategori"` // String dulu, bisa jadi FK nanti
	Merek          string         `gorm:"type:varchar(100);column:merek" json:"merek"`       // String dulu, bisa jadi FK nanti
	IDPemasok      *uint          `gorm:"index;column:id_supplier" json:"id_pemasok"`
	Pemasok        *Pemasok       `gorm:"foreignKey:IDPemasok" json:"pemasok,omitempty"`
	HargaModal     float64        `gorm:"type:decimal(15,2);not null;column:harga_modal" json:"harga_modal"` // Harga modal
	HargaJual      float64        `gorm:"type:decimal(15,2);not null;column:harga_jual" json:"harga_jual"`   // Harga jual default
	StokMinimum    int            `gorm:"not null;default:0;column:stok_minimum" json:"stok_minimum"`        // Minimum stock
	IzinDiskon     bool           `gorm:"default:true;column:izin_diskon" json:"izin_diskon"`                // Apakah bisa diskon
	Aktif          bool           `gorm:"default:true;column:aktif" json:"aktif"`
	DibuatPada     time.Time      `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada time.Time      `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
	DihapusPada    gorm.DeletedAt `gorm:"index;column:dihapus_pada" json:"-"`

	// Relationship untuk multiple images
	Images []GambarProduk `gorm:"foreignKey:IDProduk;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

// TableName mengembalikan nama tabel untuk model Produk
func (Produk) TableName() string {
	return "produk"
}

// Product adalah alias untuk backward compatibility (akan dihapus nanti)
type Product = Produk

// GambarProduk adalah model untuk gambar produk (multiple images)
type GambarProduk struct {
	ID             uint      `gorm:"primaryKey;column:id" json:"id"`
	IDProduk       uint      `gorm:"index;not null;column:id_produk" json:"id_produk"`
	Produk         Produk    `gorm:"foreignKey:IDProduk" json:"produk,omitempty"`
	PathGambar     string    `gorm:"type:text;not null;column:path_gambar" json:"path_gambar"` // Path ke file gambar
	GambarUtama    bool      `gorm:"default:false;column:gambar_utama" json:"gambar_utama"`    // Apakah gambar utama
	Urutan         int       `gorm:"default:0;column:urutan" json:"urutan"`                    // Urutan untuk sorting
	DibuatPada     time.Time `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada time.Time `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
}

// TableName mengembalikan nama tabel untuk model GambarProduk
func (GambarProduk) TableName() string {
	return "gambar_produk"
}

// ProductImage adalah alias untuk backward compatibility (akan dihapus nanti)
type ProductImage = GambarProduk

// Pemasok adalah model untuk data supplier/pemasok
type Pemasok struct {
	ID             uint           `gorm:"primaryKey;column:id" json:"id"`
	Nama           string         `gorm:"not null;column:nama" json:"nama"`
	Kontak         string         `gorm:"type:varchar(100);column:kontak" json:"kontak"`
	Telepon        string         `gorm:"type:varchar(50);column:telepon" json:"telepon"`
	Email          string         `gorm:"type:varchar(100);column:email" json:"email"`
	Alamat         string         `gorm:"type:text;column:alamat" json:"alamat"`
	Aktif          bool           `gorm:"default:true;column:aktif" json:"aktif"`
	DibuatPada     time.Time      `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada time.Time      `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
	DihapusPada    gorm.DeletedAt `gorm:"index;column:dihapus_pada" json:"-"`
}

// TableName mengembalikan nama tabel untuk model Pemasok
func (Pemasok) TableName() string {
	return "pemasok"
}

// Supplier adalah alias untuk backward compatibility (akan dihapus nanti)
type Supplier = Pemasok

// Gudang adalah model untuk data gudang
type Gudang struct {
	ID             uint           `gorm:"primaryKey;column:id" json:"id"`
	Kode           string         `gorm:"uniqueIndex;not null;column:kode" json:"kode"` // Kode gudang
	Nama           string         `gorm:"not null;column:nama" json:"nama"`
	Alamat         string         `gorm:"type:text;column:alamat" json:"alamat"`
	Aktif          bool           `gorm:"default:true;column:aktif" json:"aktif"`
	DibuatPada     time.Time      `gorm:"column:dibuat_pada" json:"dibuat_pada"`
	DiperbaruiPada time.Time      `gorm:"column:diperbarui_pada" json:"diperbarui_pada"`
	DihapusPada    gorm.DeletedAt `gorm:"index;column:dihapus_pada" json:"-"`
}

// TableName mengembalikan nama tabel untuk model Gudang
func (Gudang) TableName() string {
	return "gudang"
}

// Warehouse adalah alias untuk backward compatibility (akan dihapus nanti)
type Warehouse = Gudang
