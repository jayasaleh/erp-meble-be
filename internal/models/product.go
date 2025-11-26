package models

import (
	"time"

	"gorm.io/gorm"
)

// Produk adalah model untuk master data barang
type Produk struct {
	ID            uint           `gorm:"primaryKey;column:id" json:"id"`
	SKU           string         `gorm:"uniqueIndex;not null;column:sku" json:"sku"`
	Barcode       string         `gorm:"uniqueIndex;column:barcode" json:"barcode"` // Nullable, unique jika ada
	Name          string         `gorm:"not null;column:nama" json:"name"`
	Category      string         `gorm:"type:varchar(100);column:kategori" json:"category"` // String dulu, bisa jadi FK nanti
	Brand         string         `gorm:"type:varchar(100);column:merek" json:"brand"`       // String dulu, bisa jadi FK nanti
	SupplierID    *uint          `gorm:"index;column:id_supplier" json:"supplier_id"`
	Supplier      *Pemasok       `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	BasePrice     float64        `gorm:"type:decimal(15,2);not null;column:harga_modal" json:"base_price"`   // Harga modal
	SellingPrice  float64        `gorm:"type:decimal(15,2);not null;column:harga_jual" json:"selling_price"` // Harga jual default
	MinStock      int            `gorm:"not null;default:0;column:stok_minimum" json:"min_stock"`            // Minimum stock
	AllowDiscount bool           `gorm:"default:true;column:izin_diskon" json:"allow_discount"`              // Apakah bisa diskon
	IsActive      bool           `gorm:"default:true;column:aktif" json:"is_active"`
	CreatedAt     time.Time      `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:diperbarui_pada" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;column:dihapus_pada" json:"-"`

	// Relationship untuk multiple images
	Images []GambarProduk `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

// TableName mengembalikan nama tabel untuk model Produk
func (Produk) TableName() string {
	return "produk"
}

// Product adalah alias untuk backward compatibility (akan dihapus nanti)
type Product = Produk

// GambarProduk adalah model untuk gambar produk (multiple images)
type GambarProduk struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	ProductID uint      `gorm:"index;not null;column:id_produk" json:"product_id"`
	Product   Produk    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ImagePath string    `gorm:"type:text;not null;column:path_gambar" json:"image_path"` // Path ke file gambar
	IsPrimary bool      `gorm:"default:false;column:gambar_utama" json:"is_primary"`     // Apakah gambar utama
	Order     int       `gorm:"default:0;column:urutan" json:"order"`                    // Urutan untuk sorting
	CreatedAt time.Time `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:diperbarui_pada" json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model GambarProduk
func (GambarProduk) TableName() string {
	return "gambar_produk"
}

// ProductImage adalah alias untuk backward compatibility (akan dihapus nanti)
type ProductImage = GambarProduk

// Pemasok adalah model untuk data supplier/pemasok
type Pemasok struct {
	ID            uint           `gorm:"primaryKey;column:id" json:"id"`
	Name          string         `gorm:"not null;column:nama" json:"name"`
	ContactPerson string         `gorm:"type:varchar(100);column:kontak" json:"contact_person"`
	Phone         string         `gorm:"type:varchar(50);column:telepon" json:"phone"`
	Email         string         `gorm:"type:varchar(100);column:email" json:"email"`
	Address       string         `gorm:"type:text;column:alamat" json:"address"`
	IsActive      bool           `gorm:"default:true;column:aktif" json:"is_active"`
	CreatedAt     time.Time      `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:diperbarui_pada" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;column:dihapus_pada" json:"-"`
}

// TableName mengembalikan nama tabel untuk model Pemasok
func (Pemasok) TableName() string {
	return "pemasok"
}

// Supplier adalah alias untuk backward compatibility (akan dihapus nanti)
type Supplier = Pemasok

// Gudang adalah model untuk data gudang
type Gudang struct {
	ID        uint           `gorm:"primaryKey;column:id" json:"id"`
	Code      string         `gorm:"uniqueIndex;not null;column:kode" json:"code"` // Kode gudang
	Name      string         `gorm:"not null;column:nama" json:"name"`
	Address   string         `gorm:"type:text;column:alamat" json:"address"`
	IsActive  bool           `gorm:"default:true;column:aktif" json:"is_active"`
	CreatedAt time.Time      `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:diperbarui_pada" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:dihapus_pada" json:"-"`
}

// TableName mengembalikan nama tabel untuk model Gudang
func (Gudang) TableName() string {
	return "gudang"
}

// Warehouse adalah alias untuk backward compatibility (akan dihapus nanti)
type Warehouse = Gudang
