package models

import (
	"time"
)

// PurchaseOrder adalah model untuk Purchase Order
type PurchaseOrder struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	PONumber       string     `gorm:"uniqueIndex;not null" json:"po_number"`
	SupplierID     uint       `gorm:"index;not null" json:"supplier_id"`
	Supplier       Supplier   `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	OrderDate      time.Time  `gorm:"not null" json:"order_date"`
	DueDate        *time.Time `json:"due_date"`
	Status         string     `gorm:"type:varchar(30);default:'draft'" json:"status"` // draft, sent, approved, partially_received, completed, cancelled
	TotalAmount    float64    `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	CreatedBy      uint       `gorm:"index;not null" json:"created_by"`
	CreatedByUser  User       `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	ApprovedBy     *uint      `gorm:"index" json:"approved_by"`
	ApprovedByUser *User      `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relationship
	Items []PurchaseOrderItem `gorm:"foreignKey:POID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// PurchaseOrderItem adalah model untuk detail item Purchase Order
type PurchaseOrderItem struct {
	ID               uint          `gorm:"primaryKey" json:"id"`
	POID             uint          `gorm:"index;not null" json:"po_id"`
	PurchaseOrder    PurchaseOrder `gorm:"foreignKey:POID" json:"purchase_order,omitempty"`
	ProductID        uint          `gorm:"index;not null" json:"product_id"`
	Product          Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity         int           `gorm:"not null" json:"quantity"`
	UnitPrice        float64       `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	Subtotal         float64       `gorm:"type:decimal(15,2);not null" json:"subtotal"` // quantity × unit_price
	ReceivedQuantity int           `gorm:"default:0" json:"received_quantity"`          // Jumlah yang sudah diterima
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}
