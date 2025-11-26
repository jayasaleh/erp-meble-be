package models

import (
	"time"
)

// Sales adalah model untuk transaksi penjualan (POS)
type Sales struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	TransactionNumber string    `gorm:"uniqueIndex;not null" json:"transaction_number"`
	CustomerName      string    `gorm:"type:varchar(100)" json:"customer_name"`
	CustomerContact   string    `gorm:"type:varchar(50)" json:"customer_contact"`
	Subtotal          float64   `gorm:"type:decimal(15,2);not null" json:"subtotal"`         // Subtotal sebelum diskon
	DiscountAmount    float64   `gorm:"type:decimal(15,2);default:0" json:"discount_amount"` // Jumlah diskon
	TotalAmount       float64   `gorm:"type:decimal(15,2);not null" json:"total_amount"`     // Total setelah diskon
	PaymentMethod     string    `gorm:"type:varchar(20);not null" json:"payment_method"`     // cash, transfer, qris
	PaymentAmount     float64   `gorm:"type:decimal(15,2);not null" json:"payment_amount"`   // Jumlah yang dibayar
	ChangeAmount      float64   `gorm:"type:decimal(15,2);default:0" json:"change_amount"`   // Kembalian
	CashierID         uint      `gorm:"index;not null" json:"cashier_id"`
	Cashier           User      `gorm:"foreignKey:CashierID" json:"cashier,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relationship
	Items []SalesItem `gorm:"foreignKey:SalesID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// SalesItem adalah model untuk detail item penjualan
type SalesItem struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SalesID         uint      `gorm:"index;not null" json:"sales_id"`
	Sales           Sales     `gorm:"foreignKey:SalesID" json:"sales,omitempty"`
	ProductID       uint      `gorm:"index;not null" json:"product_id"`
	Product         Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity        int       `gorm:"not null" json:"quantity"`
	UnitPrice       float64   `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	DiscountPercent *float64  `gorm:"type:decimal(5,2)" json:"discount_percent"`   // Persentase diskon
	DiscountAmount  *float64  `gorm:"type:decimal(15,2)" json:"discount_amount"`   // Jumlah diskon
	Subtotal        float64   `gorm:"type:decimal(15,2);not null" json:"subtotal"` // Total setelah diskon
	CreatedAt       time.Time `json:"created_at"`
}
