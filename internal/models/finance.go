package models

import (
	"time"
)

// SupplierDebt adalah model untuk hutang supplier
type SupplierDebt struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	SupplierID      uint       `gorm:"index;not null" json:"supplier_id"`
	Supplier        Supplier   `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	POID            *uint      `gorm:"index" json:"po_id"`                                  // Link ke Purchase Order
	StockInID       *uint      `gorm:"index" json:"stock_in_id"`                            // Link ke Stock In
	Amount          float64    `gorm:"type:decimal(15,2);not null" json:"amount"`           // Jumlah hutang
	PaidAmount      float64    `gorm:"type:decimal(15,2);default:0" json:"paid_amount"`     // Jumlah yang sudah dibayar
	RemainingAmount float64    `gorm:"type:decimal(15,2);not null" json:"remaining_amount"` // Sisa hutang
	DueDate         *time.Time `json:"due_date"`
	Status          string     `gorm:"type:varchar(20);default:'unpaid'" json:"status"` // unpaid, partially_paid, paid
	PaymentProof    string     `gorm:"type:text" json:"payment_proof"`                  // Path ke file bukti bayar
	PaidAt          *time.Time `json:"paid_at"`
	PaidBy          *uint      `gorm:"index" json:"paid_by"`
	PaidByUser      *User      `gorm:"foreignKey:PaidBy" json:"paid_by_user,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
