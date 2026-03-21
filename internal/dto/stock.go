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
	Quantity      int       `json:"quantity"`      // + / -
	BalanceAfter  int       `json:"balance_after"` // Running balance
	WarehouseName string    `json:"warehouse_name"`
	OperatorName  string    `json:"operator_name"`
	Notes         string    `json:"notes"`
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
	ActualStock int  `json:"actual_stock" binding:"required,min=0"`
}
