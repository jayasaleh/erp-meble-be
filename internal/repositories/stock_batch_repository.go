package repositories

import (
	"real-erp-mebel/be/internal/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockBatchRepository interface {
	// Create batch baru
	Create(tx *gorm.DB, batch *models.StokBatch) error
	
	// Ambil batch untuk FIFO (urutkan berdasarkan tanggal masuk ASC)
	GetAvailableBatches(tx *gorm.DB, productID, warehouseID uint) ([]models.StokBatch, error)
	
	// Update batch (untuk decrement qty)
	Update(tx *gorm.DB, batch *models.StokBatch) error
	
	// Get batch by ID
	FindByID(batchID uint) (*models.StokBatch, error)
}

type stockBatchRepository struct {
	db *gorm.DB
}

func NewStockBatchRepository(db *gorm.DB) StockBatchRepository {
	return &stockBatchRepository{db: db}
}

func (r *stockBatchRepository) Create(tx *gorm.DB, batch *models.StokBatch) error {
	return tx.Create(batch).Error
}

// GetAvailableBatches mengambil batch yang masih punya stok, ORDER BY tanggal_masuk ASC (FIFO)
// CRITICAL: Menggunakan FOR UPDATE untuk row-level locking
func (r *stockBatchRepository) GetAvailableBatches(tx *gorm.DB, productID, warehouseID uint) ([]models.StokBatch, error) {
	var batches []models.StokBatch
	
	err := tx.Where("id_produk = ? AND id_gudang = ? AND jumlah_saat_ini > 0 AND aktif = ?", 
		productID, warehouseID, true).
		Order("tanggal_masuk ASC"). // FIFO: Yang masuk pertama keluar pertama
		Clauses(clause.Locking{Strength: "UPDATE"}). // FOR UPDATE - Row Lock
		Find(&batches).Error
	
	return batches, err
}

func (r *stockBatchRepository) Update(tx *gorm.DB, batch *models.StokBatch) error {
	batch.DiperbaruiPada = time.Now()
	
	// Jika qty habis, set aktif=false
	if batch.JumlahSaatIni <= 0 {
		batch.Aktif = false
	}
	
	return tx.Save(batch).Error
}

func (r *stockBatchRepository) FindByID(batchID uint) (*models.StokBatch, error) {
	var batch models.StokBatch
	err := r.db.Preload("Produk").Preload("Gudang").First(&batch, batchID).Error
	if err != nil {
		return nil, err
	}
	return &batch, nil
}
