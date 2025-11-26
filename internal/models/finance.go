package models

import (
	"time"
)

// HutangPemasok adalah model untuk hutang supplier
type HutangPemasok struct {
	ID              uint       `gorm:"primaryKey;column:id" json:"id"`
	SupplierID      uint       `gorm:"index;not null;column:id_supplier" json:"supplier_id"`
	Supplier        Pemasok    `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	POID            *uint      `gorm:"index;column:id_po" json:"po_id"`                                        // Link ke Purchase Order
	StockInID       *uint      `gorm:"index;column:id_stock_in" json:"stock_in_id"`                            // Link ke Stock In
	Amount          float64    `gorm:"type:decimal(15,2);not null;column:jumlah" json:"amount"`                // Jumlah hutang
	PaidAmount      float64    `gorm:"type:decimal(15,2);default:0;column:jumlah_dibayar" json:"paid_amount"`  // Jumlah yang sudah dibayar
	RemainingAmount float64    `gorm:"type:decimal(15,2);not null;column:sisa_hutang" json:"remaining_amount"` // Sisa hutang
	DueDate         *time.Time `gorm:"column:jatuh_tempo" json:"due_date"`
	Status          string     `gorm:"type:varchar(20);default:'unpaid';column:status" json:"status"` // unpaid, partially_paid, paid
	PaymentProof    string     `gorm:"type:text;column:bukti_pembayaran" json:"payment_proof"`        // Path ke file bukti bayar
	PaidAt          *time.Time `gorm:"column:dibayar_pada" json:"paid_at"`
	PaidBy          *uint      `gorm:"index;column:dibayar_oleh" json:"paid_by"`
	PaidByUser      *Pengguna  `gorm:"foreignKey:PaidBy" json:"paid_by_user,omitempty"`
	CreatedAt       time.Time  `gorm:"column:dibuat_pada" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:diperbarui_pada" json:"updated_at"`
}

// TableName mengembalikan nama tabel untuk model HutangPemasok
func (HutangPemasok) TableName() string {
	return "hutang_pemasok"
}

// SupplierDebt adalah alias untuk backward compatibility (akan dihapus nanti)
type SupplierDebt = HutangPemasok
