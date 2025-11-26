package models

import (
	"time"
)

// ReturPenjualan adalah model untuk retur penjualan (customer return barang)
type ReturPenjualan struct {
	ID              uint       `gorm:"primaryKey;column:id" json:"id"`
	ReturnNumber    string     `gorm:"uniqueIndex;not null;column:nomor_retur" json:"return_number"`
	SalesID         uint       `gorm:"index;not null;column:id_penjualan" json:"sales_id"`
	Sales           Penjualan  `gorm:"foreignKey:SalesID" json:"sales,omitempty"`
	CustomerName    string     `gorm:"type:varchar(100);column:nama_pelanggan" json:"customer_name"`
	CustomerContact string     `gorm:"type:varchar(50);column:kontak_pelanggan" json:"customer_contact"`
	Reason          string     `gorm:"type:varchar(100);not null;column:alasan" json:"reason"` // rusak, tidak sesuai, cacat, dll
	Subtotal        float64    `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`
	TotalAmount     float64    `gorm:"type:decimal(15,2);not null;column:total" json:"total_amount"`
	RefundMethod    string     `gorm:"type:varchar(20);not null;column:metode_pengembalian" json:"refund_method"` // cash, transfer, tukar_barang
	RefundAmount    float64    `gorm:"type:decimal(15,2);not null;column:jumlah_pengembalian" json:"refund_amount"`
	Status          string     `gorm:"type:varchar(20);default:'pending';column:status" json:"status"` // pending, approved, rejected, completed
	Notes           string     `gorm:"type:text;column:keterangan" json:"notes"`
	ProcessedBy     uint       `gorm:"index;not null;column:diproses_oleh" json:"processed_by"`
	ProcessedByUser Pengguna   `gorm:"foreignKey:ProcessedBy" json:"processed_by_user,omitempty"`
	ApprovedBy      *uint      `gorm:"index;column:disetujui_oleh" json:"approved_by"`
	ApprovedByUser  *Pengguna  `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	ApprovedAt      *time.Time `gorm:"column:disetujui_pada" json:"approved_at"`
	ProcessedAt     time.Time  `gorm:"column:diproses_pada" json:"processed_at"`
	CreatedAt       time.Time  `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:diperbarui_pada" json:"updated_at"`

	// Relationship
	Items []ItemReturPenjualan `gorm:"foreignKey:SalesReturnID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model ReturPenjualan
func (ReturPenjualan) TableName() string {
	return "retur_penjualan"
}

// SalesReturn adalah alias untuk backward compatibility (akan dihapus nanti)
type SalesReturn = ReturPenjualan

// ItemReturPenjualan adalah model untuk detail item retur penjualan
type ItemReturPenjualan struct {
	ID            uint           `gorm:"primaryKey;column:id" json:"id"`
	SalesReturnID uint           `gorm:"index;not null;column:id_retur_penjualan" json:"sales_return_id"`
	SalesReturn   ReturPenjualan `gorm:"foreignKey:SalesReturnID" json:"sales_return,omitempty"`
	SalesItemID   uint           `gorm:"index;not null;column:id_item_penjualan" json:"sales_item_id"`
	SalesItem     ItemPenjualan  `gorm:"foreignKey:SalesItemID" json:"sales_item,omitempty"`
	ProductID     uint           `gorm:"index;not null;column:id_produk" json:"product_id"`
	Product       Produk         `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity      int            `gorm:"not null;column:jumlah" json:"quantity"`
	UnitPrice     float64        `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"unit_price"`
	Subtotal      float64        `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`
	WarehouseID   uint           `gorm:"index;not null;column:id_gudang" json:"warehouse_id"`
	Warehouse     Gudang         `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	CreatedAt     time.Time      `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:diperbarui_pada" json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model ItemReturPenjualan
func (ItemReturPenjualan) TableName() string {
	return "item_retur_penjualan"
}

// SalesReturnItem adalah alias untuk backward compatibility (akan dihapus nanti)
type SalesReturnItem = ItemReturPenjualan

// ReturPembelian adalah model untuk retur pembelian (return ke supplier)
type ReturPembelian struct {
	ID             uint              `gorm:"primaryKey;column:id" json:"id"`
	ReturnNumber   string            `gorm:"uniqueIndex;not null;column:nomor_retur" json:"return_number"`
	POID           *uint             `gorm:"index;column:id_po" json:"po_id"`
	PurchaseOrder  *PesananPembelian `gorm:"foreignKey:POID" json:"purchase_order,omitempty"`
	StockInID      *uint             `gorm:"index;column:id_stock_in" json:"stock_in_id"`
	StockIn        *BarangMasuk      `gorm:"foreignKey:StockInID" json:"stock_in,omitempty"`
	SupplierID     uint              `gorm:"index;not null;column:id_supplier" json:"supplier_id"`
	Supplier       Pemasok           `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Reason         string            `gorm:"type:varchar(100);not null;column:alasan" json:"reason"` // rusak, tidak sesuai, cacat, dll
	Subtotal       float64           `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`
	TotalAmount    float64           `gorm:"type:decimal(15,2);not null;column:total" json:"total_amount"`
	RefundMethod   string            `gorm:"type:varchar(20);not null;column:metode_pengembalian" json:"refund_method"` // potong_hutang, refund, tukar_barang
	Status         string            `gorm:"type:varchar(20);default:'pending';column:status" json:"status"`            // pending, approved, rejected, completed
	Notes          string            `gorm:"type:text;column:keterangan" json:"notes"`
	CreatedBy      uint              `gorm:"index;not null;column:dibuat_oleh" json:"created_by"`
	CreatedByUser  Pengguna          `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	ApprovedBy     *uint             `gorm:"index;column:disetujui_oleh" json:"approved_by"`
	ApprovedByUser *Pengguna         `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	ApprovedAt     *time.Time        `gorm:"column:disetujui_pada" json:"approved_at"`
	CreatedAt      time.Time         `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt      time.Time         `gorm:"column:diperbarui_pada" json:"updated_at"`

	// Relationship
	Items []ItemReturPembelian `gorm:"foreignKey:PurchaseReturnID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model ReturPembelian
func (ReturPembelian) TableName() string {
	return "retur_pembelian"
}

// PurchaseReturn adalah alias untuk backward compatibility (akan dihapus nanti)
type PurchaseReturn = ReturPembelian

// ItemReturPembelian adalah model untuk detail item retur pembelian
type ItemReturPembelian struct {
	ID               uint             `gorm:"primaryKey;column:id" json:"id"`
	PurchaseReturnID uint             `gorm:"index;not null;column:id_retur_pembelian" json:"purchase_return_id"`
	PurchaseReturn   ReturPembelian   `gorm:"foreignKey:PurchaseReturnID" json:"purchase_return,omitempty"`
	StockInItemID    *uint            `gorm:"index;column:id_stock_in_item" json:"stock_in_item_id"`
	StockInItem      *ItemBarangMasuk `gorm:"foreignKey:StockInItemID" json:"stock_in_item,omitempty"`
	ProductID        uint             `gorm:"index;not null;column:id_produk" json:"product_id"`
	Product          Produk           `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity         int              `gorm:"not null;column:jumlah" json:"quantity"`
	UnitPrice        float64          `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"unit_price"`
	Subtotal         float64          `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`
	WarehouseID      uint             `gorm:"index;not null;column:id_gudang" json:"warehouse_id"`
	Warehouse        Gudang           `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	CreatedAt        time.Time        `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"column:diperbarui_pada" json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model ItemReturPembelian
func (ItemReturPembelian) TableName() string {
	return "item_retur_pembelian"
}

// PurchaseReturnItem adalah alias untuk backward compatibility (akan dihapus nanti)
type PurchaseReturnItem = ItemReturPembelian
