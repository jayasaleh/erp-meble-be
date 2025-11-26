package models

import (
	"time"
)

// PesananPembelian adalah model untuk Purchase Order
type PesananPembelian struct {
	ID             uint       `gorm:"primaryKey;column:id" json:"id"`
	PONumber       string     `gorm:"uniqueIndex;not null;column:nomor_po" json:"po_number"`
	SupplierID     uint       `gorm:"index;not null;column:id_supplier" json:"supplier_id"`
	Supplier       Pemasok    `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	OrderDate      time.Time  `gorm:"not null;column:tanggal_pesan" json:"order_date"`
	DueDate        *time.Time `gorm:"column:tanggal_jatuh_tempo" json:"due_date"`
	Status         string     `gorm:"type:varchar(30);default:'draft';column:status" json:"status"` // draft, sent, approved, partially_received, completed, cancelled
	TotalAmount    float64    `gorm:"type:decimal(15,2);not null;column:total" json:"total_amount"`
	CreatedBy      uint       `gorm:"index;not null;column:dibuat_oleh" json:"created_by"`
	CreatedByUser  Pengguna   `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	ApprovedBy     *uint      `gorm:"index;column:disetujui_oleh" json:"approved_by"`
	ApprovedByUser *Pengguna  `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	ApprovedAt     *time.Time `gorm:"column:disetujui_pada" json:"approved_at"`
	CreatedAt      time.Time  `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:diperbarui_pada" json:"updated_at"`

	// Relationship
	Items []ItemPesananPembelian `gorm:"foreignKey:POID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model PesananPembelian
func (PesananPembelian) TableName() string {
	return "pesanan_pembelian"
}

// PurchaseOrder adalah alias untuk backward compatibility (akan dihapus nanti)
type PurchaseOrder = PesananPembelian

// ItemPesananPembelian adalah model untuk detail item Purchase Order
type ItemPesananPembelian struct {
	ID               uint             `gorm:"primaryKey;column:id" json:"id"`
	POID             uint             `gorm:"index;not null;column:id_po" json:"po_id"`
	PurchaseOrder    PesananPembelian `gorm:"foreignKey:POID" json:"purchase_order,omitempty"`
	ProductID        uint             `gorm:"index;not null;column:id_produk" json:"product_id"`
	Product          Produk           `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity         int              `gorm:"not null;column:jumlah" json:"quantity"`
	UnitPrice        float64          `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"unit_price"`
	Subtotal         float64          `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"` // quantity × unit_price
	ReceivedQuantity int              `gorm:"default:0;column:jumlah_diterima" json:"received_quantity"`   // Jumlah yang sudah diterima
	CreatedAt        time.Time        `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"column:diperbarui_pada" json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model ItemPesananPembelian
func (ItemPesananPembelian) TableName() string {
	return "item_pesanan_pembelian"
}

// PurchaseOrderItem adalah alias untuk backward compatibility (akan dihapus nanti)
type PurchaseOrderItem = ItemPesananPembelian
