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

	// Ambil semua batch (termasuk habis) untuk halaman detail stok
	GetAllBatches(productID, warehouseID uint, limit, offset int) ([]models.StokBatch, int64, error)

	// Update batch (untuk decrement qty)
	Update(tx *gorm.DB, batch *models.StokBatch) error

	// Get batch by ID
	FindByID(batchID uint) (*models.StokBatch, error)

	// Ambil movement opname terakhir per batch
	GetLatestOpnameByBatchIDs(batchIDs []uint) (map[uint]models.PergerakanStok, error)

	// Ambil nama operator yang memasukkan batch tersebut
	GetCreatorByBatchIDs(batchIDs []uint) (map[uint]string, error)
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

// GetAllBatches mengambil semua batch (aktif maupun habis) untuk satu produk+gudang
// ORDER BY tanggal_masuk DESC (terbaru di atas) untuk tampilan detail
func (r *stockBatchRepository) GetAllBatches(productID, warehouseID uint, limit, offset int) ([]models.StokBatch, int64, error) {
	var batches []models.StokBatch
	var total int64

	query := r.db.Model(&models.StokBatch{}).
		Where("id_produk = ? AND id_gudang = ?", productID, warehouseID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Produk").Preload("Gudang").
		Order("tanggal_masuk DESC").
		Limit(limit).Offset(offset).
		Find(&batches).Error

	return batches, total, err
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

func (r *stockBatchRepository) GetLatestOpnameByBatchIDs(batchIDs []uint) (map[uint]models.PergerakanStok, error) {
	result := make(map[uint]models.PergerakanStok)
	if len(batchIDs) == 0 {
		return result, nil
	}

	var movements []models.PergerakanStok
	err := r.db.
		Where("id_batch IN ? AND tipe_referensi = ?", batchIDs, "opname").
		Order("dibuat_pada DESC").
		Find(&movements).Error
	if err != nil {
		return nil, err
	}

	for _, m := range movements {
		if m.IDBatch == nil {
			continue
		}
		// Karena urut DESC, record pertama per batch adalah opname terakhir.
		if _, exists := result[*m.IDBatch]; !exists {
			result[*m.IDBatch] = m
		}
	}

	return result, nil
}

func (r *stockBatchRepository) GetCreatorByBatchIDs(batchIDs []uint) (map[uint]string, error) {
	result := make(map[uint]string)
	if len(batchIDs) == 0 {
		return result, nil
	}

	var movements []models.PergerakanStok
	// Ambil pergerakan in atau transfer_in perdana untuk tahu siapa creator
	err := r.db.Preload("Pengguna").
		Where("id_batch IN ? AND (tipe_pergerakan = ? OR tipe_pergerakan = ? OR tipe_pergerakan = ?)", batchIDs, "in", "transfer_in", "adjustment").
		Order("id ASC").
		Find(&movements).Error
	if err != nil {
		return nil, err
	}

	for _, m := range movements {
		if m.IDBatch == nil {
			continue
		}
		// Isi hanya kalau belum ada (alias movement paling pertama yg nyiptain batch)
		if _, exists := result[*m.IDBatch]; !exists {
			if m.Pengguna.Nama != "" {
				result[*m.IDBatch] = m.Pengguna.Nama
			} else {
				result[*m.IDBatch] = "Sistem"
			}
		}
	}

	return result, nil
}
