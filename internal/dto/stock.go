package dto

import "time"

// InventoryResponse adalah response untuk stok saat ini
type InventoryResponse struct {
	ID           uint      `json:"id"`
	ProductID    uint      `json:"product_id"`
	ProductSKU   string    `json:"product_sku"`
	ProductName  string    `json:"product_name"`
	WarehouseID  uint      `json:"warehouse_id"`
	Warehouse    string    `json:"warehouse_name"`
	CurrentStock int       `json:"current_stock"`
	LastUpdate   time.Time `json:"last_update"`
}

// StockMovementResponse adalah response untuk history pergerakan stok
type StockMovementResponse struct {
	ID            uint      `json:"id"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`     // in, out, adjustment, dll
	ReferenceType string    `json:"ref_type"` // sales, po, manual, dll
	ReferenceID   *uint     `json:"ref_id"`
	ProductID     uint      `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Quantity      int       `json:"quantity"`      // + / -
	SystemStock   int       `json:"system_stock"`  // BalanceBefore (System stock during opname)
	BalanceAfter  int       `json:"balance_after"` // Running balance (Actual physical stock after opname)
	CostPrice     float64   `json:"cost_price"`    // Cost price per piece to evaluate loss/surplus
	WarehouseID   uint      `json:"warehouse_id"`
	WarehouseName string    `json:"warehouse_name"`
	OperatorName  string    `json:"operator_name"`
	Notes         string    `json:"notes"`
	BatchID       *uint     `json:"batch_id"`
}

// BatchResponse adalah response untuk detail batch FIFO
type BatchResponse struct {
	ID            uint       `json:"id"`
	ProductID     uint       `json:"product_id"`
	ProductName   string     `json:"product_name"`
	ProductSKU    string     `json:"product_sku"`
	WarehouseID   uint       `json:"warehouse_id"`
	WarehouseName string     `json:"warehouse_name"`
	EntryDate     time.Time  `json:"entry_date"`
	ExpiryDate    *time.Time `json:"expiry_date"`
	InitialQty    int        `json:"initial_qty"`
	CurrentQty    int        `json:"current_qty"`
	CostPrice     float64    `json:"cost_price"`
	ReferenceType string     `json:"reference_type"`
	ReferenceID   *uint      `json:"reference_id"`
	Notes         string     `json:"notes"`
	IsActive      bool       `json:"is_active"`
	LastOpnameAt  *time.Time `json:"last_opname_at"`
	LastOpnameQty *int       `json:"last_opname_qty"`
	OperatorName  string     `json:"operator_name"`
	CreatedAt     time.Time  `json:"created_at"`
}

// CreateStockInRequest adalah request untuk barang masuk manual
type CreateStockInRequest struct {
	WarehouseID uint               `json:"warehouse_id" binding:"required"`
	Date        time.Time          `json:"date"` // Optional, default now
	Notes       string             `json:"notes"`
	Items       []StockRequestItem `json:"items" binding:"required,dive"`
}

// CreateStockOutRequest adalah request untuk barang keluar manual (usage/damaged etc, not sales)
type CreateStockOutRequest struct {
	WarehouseID uint               `json:"warehouse_id" binding:"required"`
	Date        time.Time          `json:"date"`                      // Optional, default now
	Reason      string             `json:"reason" binding:"required"` // rusak, expired, usage, dll
	Items       []StockRequestItem `json:"items" binding:"required,dive"`
}

// CreateStockOpnameRequest adalah request untuk penyesuaian stok (stock take)
// User input stok fisik, sistem hitung selisih
type CreateStockOpnameRequest struct {
	WarehouseID uint              `json:"warehouse_id" binding:"required"`
	Date        time.Time         `json:"date"`
	Notes       string            `json:"notes"`
	Items       []StockOpnameItem `json:"items" binding:"required,dive"`
}

// CreateStockTransferRequest adalah DTO untuk memindahkan stok antar gudang
type CreateStockTransferRequest struct {
	SourceWarehouseID uint               `json:"source_warehouse_id" binding:"required"`
	TargetWarehouseID uint               `json:"target_warehouse_id" binding:"required"`
	Date              *time.Time         `json:"date"`
	Notes             string             `json:"notes"`
	Items             []StockRequestItem `json:"items" binding:"required,dive"`
}

type StockRequestItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type StockOpnameItem struct {
	ProductID   uint `json:"product_id" binding:"required"`
	BatchID     uint `json:"batch_id"`
	ActualStock int  `json:"actual_stock" binding:"required,min=0"`
}
