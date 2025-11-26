package models

import (
	"time"

	"gorm.io/gorm"
)

// Product adalah model untuk master data barang
type Product struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	SKU           string         `gorm:"uniqueIndex;not null" json:"sku"`
	Barcode       string         `gorm:"uniqueIndex" json:"barcode"` // Nullable, unique jika ada
	Name          string         `gorm:"not null" json:"name"`
	Category      string         `gorm:"type:varchar(100)" json:"category"` // String dulu, bisa jadi FK nanti
	Brand         string         `gorm:"type:varchar(100)" json:"brand"`   // String dulu, bisa jadi FK nanti
	SupplierID    *uint          `gorm:"index" json:"supplier_id"`
	Supplier      *Supplier      `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	BasePrice     float64        `gorm:"type:decimal(15,2);not null" json:"base_price"`     // Harga modal
	SellingPrice  float64        `gorm:"type:decimal(15,2);not null" json:"selling_price"` // Harga jual default
	MinStock      int            `gorm:"not null;default:0" json:"min_stock"`               // Minimum stock
	AllowDiscount bool           `gorm:"default:true" json:"allow_discount"`                // Apakah bisa diskon
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship untuk multiple images
	Images []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

// ProductImage adalah model untuk gambar produk (multiple images)
type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ImagePath string    `gorm:"type:text;not null" json:"image_path"` // Path ke file gambar
	IsPrimary bool      `gorm:"default:false" json:"is_primary"`     // Apakah gambar utama
	Order     int       `gorm:"default:0" json:"order"`              // Urutan untuk sorting
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Supplier adalah model untuk data supplier/pemasok
type Supplier struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	ContactPerson string        `gorm:"type:varchar(100)" json:"contact_person"`
	Phone        string         `gorm:"type:varchar(50)" json:"phone"`
	Email        string         `gorm:"type:varchar(100)" json:"email"`
	Address      string         `gorm:"type:text" json:"address"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Warehouse adalah model untuk data gudang
type Warehouse struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"uniqueIndex;not null" json:"code"` // Kode gudang
	Name      string         `gorm:"not null" json:"name"`
	Address   string         `gorm:"type:text" json:"address"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}


