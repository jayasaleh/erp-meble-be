package models

import (
	"time"
)

// StockIn adalah model untuk transaksi barang masuk
type StockIn struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	TransactionNumber string     `gorm:"uniqueIndex;not null" json:"transaction_number"`
	SupplierID        *uint      `gorm:"index" json:"supplier_id"`
	Supplier          *Supplier  `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	POID              *uint      `gorm:"index" json:"po_id"` // Link ke Purchase Order
	ReceivedBy        uint       `gorm:"not null" json:"received_by"`
	ReceivedByUser    User       `gorm:"foreignKey:ReceivedBy" json:"received_by_user,omitempty"`
	ReceivedAt        time.Time  `json:"received_at"`
	ApprovedBy        *uint      `gorm:"index" json:"approved_by"`
	ApprovedByUser    *User      `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at"`
	Status            string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, rejected
	Notes             string     `gorm:"type:text" json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// Relationship
	Items []StockInItem `gorm:"foreignKey:StockInID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// StockInItem adalah model untuk detail item barang masuk
type StockInItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	StockInID   uint      `gorm:"index;not null" json:"stock_in_id"`
	StockIn     StockIn   `gorm:"foreignKey:StockInID" json:"stock_in,omitempty"`
	ProductID   uint      `gorm:"index;not null" json:"product_id"`
	Product     Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	UnitPrice   float64   `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	POPrice     *float64  `gorm:"type:decimal(15,2)" json:"po_price"` // Harga dari PO untuk perbandingan
	WarehouseID uint      `gorm:"index;not null" json:"warehouse_id"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Location    string    `gorm:"type:text" json:"location"` // "Rak A, Slot B" - bisa jadi FK nanti
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockOut adalah model untuk transaksi barang keluar
type StockOut struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	TransactionNumber string    `gorm:"uniqueIndex;not null" json:"transaction_number"`
	Reason            string    `gorm:"type:varchar(50);not null" json:"reason"` // penjualan, mutasi, produksi, rusak, adjustment
	ReferenceID       *uint     `gorm:"index" json:"reference_id"`               // Link ke sales, transfer, dll
	ReferenceType     string    `gorm:"type:varchar(50)" json:"reference_type"`  // sales, transfer, dll
	CreatedBy         uint      `gorm:"not null" json:"created_by"`
	CreatedByUser     User      `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relationship
	Items []StockOutItem `gorm:"foreignKey:StockOutID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

// StockOutItem adalah model untuk detail item barang keluar
type StockOutItem struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	StockOutID    uint         `gorm:"index;not null" json:"stock_out_id"`
	StockOut      StockOut     `gorm:"foreignKey:StockOutID" json:"stock_out,omitempty"`
	ProductID     uint         `gorm:"index;not null" json:"product_id"`
	Product       Product      `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity      int          `gorm:"not null" json:"quantity"`
	WarehouseID   uint         `gorm:"index;not null" json:"warehouse_id"`
	Warehouse     Warehouse    `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	StockInItemID *uint        `gorm:"index" json:"stock_in_item_id"` // Untuk FIFO tracking
	StockInItem   *StockInItem `gorm:"foreignKey:StockInItemID" json:"stock_in_item,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// InventoryStock adalah model untuk stok saat ini per produk per gudang
type InventoryStock struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ProductID      uint       `gorm:"uniqueIndex:idx_product_warehouse;not null" json:"product_id"`
	Product        Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	WarehouseID    uint       `gorm:"uniqueIndex:idx_product_warehouse;not null" json:"warehouse_id"`
	Warehouse      Warehouse  `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	Quantity       int        `gorm:"not null;default:0" json:"quantity"` // Stok saat ini
	LastMovementAt *time.Time `json:"last_movement_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// StockMovement adalah model untuk kartu stok (semua pergerakan)
type StockMovement struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProductID     uint      `gorm:"index;not null" json:"product_id"`
	Product       Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	WarehouseID   uint      `gorm:"index;not null" json:"warehouse_id"`
	Warehouse     Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse,omitempty"`
	MovementType  string    `gorm:"type:varchar(20);not null" json:"movement_type"`  // in, out, transfer_in, transfer_out, adjustment
	ReferenceType string    `gorm:"type:varchar(50);not null" json:"reference_type"` // stock_in, stock_out, sales, transfer, opname
	ReferenceID   *uint     `gorm:"index" json:"reference_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`      // Positif untuk in, negatif untuk out
	BalanceAfter  int       `gorm:"not null" json:"balance_after"` // Running balance
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}
