package models

import (
	"time"
)

// Penjualan adalah model untuk transaksi penjualan (POS)
type Penjualan struct {
	ID                uint      `gorm:"primaryKey;column:id" json:"id"`
	TransactionNumber string    `gorm:"uniqueIndex;not null;column:nomor_transaksi" json:"transaction_number"`
	CustomerName      string    `gorm:"type:varchar(100);column:nama_pelanggan" json:"customer_name"`
	CustomerContact   string    `gorm:"type:varchar(50);column:kontak_pelanggan" json:"customer_contact"`
	Subtotal          float64   `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`                // Subtotal sebelum diskon
	DiscountAmount    float64   `gorm:"type:decimal(15,2);default:0;column:jumlah_diskon" json:"discount_amount"`   // Jumlah diskon
	TotalAmount       float64   `gorm:"type:decimal(15,2);not null;column:total" json:"total_amount"`               // Total setelah diskon
	PaymentMethod     string    `gorm:"type:varchar(20);not null;column:metode_pembayaran" json:"payment_method"`   // cash, transfer, qris
	PaymentAmount     float64   `gorm:"type:decimal(15,2);not null;column:jumlah_pembayaran" json:"payment_amount"` // Jumlah yang dibayar
	ChangeAmount      float64   `gorm:"type:decimal(15,2);default:0;column:jumlah_kembalian" json:"change_amount"`  // Kembalian
	CashierID         uint      `gorm:"index;not null;column:id_kasir" json:"cashier_id"`
	Cashier           Pengguna  `gorm:"foreignKey:CashierID" json:"cashier,omitempty"`
	CreatedAt         time.Time `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:diperbarui_pada" json:"updated_at"`

	// Relationship
	Items []ItemPenjualan `gorm:"foreignKey:SalesID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// TableName mengembalikan nama tabel untuk model Penjualan
func (Penjualan) TableName() string {
	return "penjualan"
}

// Sales adalah alias untuk backward compatibility (akan dihapus nanti)
type Sales = Penjualan

// ItemPenjualan adalah model untuk detail item penjualan
type ItemPenjualan struct {
	ID              uint      `gorm:"primaryKey;column:id" json:"id"`
	SalesID         uint      `gorm:"index;not null;column:id_penjualan" json:"sales_id"`
	Sales           Penjualan `gorm:"foreignKey:SalesID" json:"sales,omitempty"`
	ProductID       uint      `gorm:"index;not null;column:id_produk" json:"product_id"`
	Product         Produk    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity        int       `gorm:"not null;column:jumlah" json:"quantity"`
	UnitPrice       float64   `gorm:"type:decimal(15,2);not null;column:harga_satuan" json:"unit_price"`
	DiscountPercent *float64  `gorm:"type:decimal(5,2);column:persen_diskon" json:"discount_percent"` // Persentase diskon
	DiscountAmount  *float64  `gorm:"type:decimal(15,2);column:jumlah_diskon" json:"discount_amount"` // Jumlah diskon
	Subtotal        float64   `gorm:"type:decimal(15,2);not null;column:subtotal" json:"subtotal"`    // Total setelah diskon
	CreatedAt       time.Time `gorm:"column:dibuat_pada" json:"created_at"`
}

// TableName mengembalikan nama tabel untuk model ItemPenjualan
func (ItemPenjualan) TableName() string {
	return "item_penjualan"
}

// SalesItem adalah alias untuk backward compatibility (akan dihapus nanti)
type SalesItem = ItemPenjualan
