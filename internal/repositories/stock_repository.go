package repositories

import (
	"real-erp-mebel/be/internal/models"
	"time"

	"gorm.io/gorm"
)

type StockRepository interface {
	// Query
	GetStockByProductAndWarehouse(productID, warehouseID uint) (*models.StokInventori, error)
	GetStockByWarehouse(warehouseID uint) ([]models.StokInventori, error)
	GetStockHistory(warehouseID, productID uint, limit, offset int) ([]models.PergerakanStok, int64, error)

	// Transaction
	BeginTx() *gorm.DB

	// Atomic Update (used within Tx)
	UpdateStockBalance(tx *gorm.DB, productID, warehouseID uint, delta int) error
	CreateStockMovement(tx *gorm.DB, movement *models.PergerakanStok) error

	// Headers (used within Tx)
	CreateStockIn(tx *gorm.DB, header *models.BarangMasuk) error
	CreateStockOut(tx *gorm.DB, header *models.BarangKeluar) error
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *stockRepository) GetStockByProductAndWarehouse(productID, warehouseID uint) (*models.StokInventori, error) {
	var stock models.StokInventori
	err := r.db.Preload("Produk").Preload("Gudang").
		Where("id_produk = ? AND id_gudang = ?", productID, warehouseID).
		First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) GetStockByWarehouse(warehouseID uint) ([]models.StokInventori, error) {
	var stocks []models.StokInventori
	query := r.db.Preload("Produk").Preload("Gudang")

	if warehouseID != 0 {
		query = query.Where("id_gudang = ?", warehouseID)
	}

	err := query.Find(&stocks).Error
	return stocks, err
}

func (r *stockRepository) GetStockHistory(warehouseID, productID uint, limit, offset int) ([]models.PergerakanStok, int64, error) {
	var movements []models.PergerakanStok
	var total int64

	query := r.db.Model(&models.PergerakanStok{})

	if warehouseID != 0 {
		query = query.Where("id_gudang = ?", warehouseID)
	}
	if productID != 0 {
		query = query.Where("id_produk = ?", productID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Produk").Preload("Gudang").Preload("Pengguna").
		Order("dibuat_pada DESC").
		Limit(limit).Offset(offset).
		Find(&movements).Error

	return movements, total, err
}

// UpdateStockBalance handles atomic update (UPSERT logic essentially)
func (r *stockRepository) UpdateStockBalance(tx *gorm.DB, productID, warehouseID uint, delta int) error {
	// Check if record exists
	var stock models.StokInventori
	err := tx.Where("id_produk = ? AND id_gudang = ?", productID, warehouseID).First(&stock).Error

	now := time.Now()

	if err == gorm.ErrRecordNotFound {
		// Create new record
		// NOTE: Check if we are decrementing from 0 (creating negative stock) - maybe allow?
		// For now we assume controller checks logic, but physically it shouldn't happen usually.
		newStock := models.StokInventori{
			IDProduk:               productID,
			IDGudang:               warehouseID,
			Jumlah:                 delta,
			PergerakanTerakhirPada: &now,
			DiperbaruiPada:         now,
		}
		return tx.Create(&newStock).Error
	} else if err != nil {
		return err
	}

	// Update existing
	// Using gorm.Expr to ensure atomic update in DB
	return tx.Model(&stock).Updates(map[string]interface{}{
		"jumlah":                   gorm.Expr("jumlah + ?", delta),
		"pergerakan_terakhir_pada": now,
		"diperbarui_pada":          now,
	}).Error
}

func (r *stockRepository) CreateStockMovement(tx *gorm.DB, movement *models.PergerakanStok) error {
	// Needs to fetch current balance to set SaldoSetelah correctly?
	// ACTUALLY: UpdateStockBalance is atomic, so we need to know the RESULTING balance.
	// Alternative: Read lock row -> calc -> write.
	// Simplified approach for now:
	// We trust the `UpdateStockBalance` has run appropriately.
	// To get accurate "SaldoSetelah", we should re-read the stock record.

	var stock models.StokInventori
	if err := tx.Where("id_produk = ? AND id_gudang = ?", movement.IDProduk, movement.IDGudang).First(&stock).Error; err != nil {
		return err // Should exist because we called UpdateStockBalance before this
	}

	movement.SaldoSetelah = stock.Jumlah
	return tx.Create(movement).Error
}

func (r *stockRepository) CreateStockIn(tx *gorm.DB, header *models.BarangMasuk) error {
	return tx.Create(header).Error
}

func (r *stockRepository) CreateStockOut(tx *gorm.DB, header *models.BarangKeluar) error {
	return tx.Create(header).Error
}
