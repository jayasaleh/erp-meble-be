package repositories

import (
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"time"

	"gorm.io/gorm"
)

type SalesRepository interface {
	// Buat transaksi penjualan lengkap dalam TX (header + items + batch usage)
	Create(tx *gorm.DB, sale *models.Penjualan) error

	// Ambil detail penjualan by ID dengan semua relasi
	FindByID(id uint) (*models.Penjualan, error)

	// List penjualan dengan filter dan pagination
	FindAll(req *dto.ListSalesRequest) ([]models.Penjualan, int64, error)

	// Update path bukti bayar pada penjualan
	UpdateBuktiBayar(id uint, filePath string) error

	// Begin transaction
	BeginTx() *gorm.DB
}

type salesRepository struct {
	db *gorm.DB
}

func NewSalesRepository(db *gorm.DB) SalesRepository {
	return &salesRepository{db: db}
}

func (r *salesRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *salesRepository) Create(tx *gorm.DB, sale *models.Penjualan) error {
	return tx.Create(sale).Error
}

func (r *salesRepository) FindByID(id uint) (*models.Penjualan, error) {
	var sale models.Penjualan
	err := r.db.
		Preload("Gudang").
		Preload("Kasir").
		Preload("Items").
		Preload("Items.Produk").
		Preload("Items.Gudang").
		Preload("Items.BatchUsage").
		Preload("Items.BatchUsage.Batch").
		First(&sale, id).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *salesRepository) FindAll(req *dto.ListSalesRequest) ([]models.Penjualan, int64, error) {
	var sales []models.Penjualan
	var total int64

	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := r.db.Model(&models.Penjualan{})

	// Filter gudang
	if req.IDGudang != nil {
		query = query.Where("id_gudang = ?", *req.IDGudang)
	}
	// Filter kasir
	if req.IDKasir != nil {
		query = query.Where("id_kasir = ?", *req.IDKasir)
	}
	// Filter metode pembayaran
	if req.MetodePembayaran != "" {
		query = query.Where("metode_pembayaran = ?", req.MetodePembayaran)
	}
	// Filter tanggal
	if req.TanggalDari != nil {
		startOfDay := time.Date(req.TanggalDari.Year(), req.TanggalDari.Month(), req.TanggalDari.Day(), 0, 0, 0, 0, req.TanggalDari.Location())
		query = query.Where("dibuat_pada >= ?", startOfDay)
	}
	if req.TanggalSampai != nil {
		endOfDay := time.Date(req.TanggalSampai.Year(), req.TanggalSampai.Month(), req.TanggalSampai.Day(), 23, 59, 59, 999999999, req.TanggalSampai.Location())
		query = query.Where("dibuat_pada <= ?", endOfDay)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with relations
	err := query.
		Preload("Gudang").
		Preload("Kasir").
		Order("dibuat_pada DESC").
		Limit(limit).Offset(offset).
		Find(&sales).Error

	return sales, total, err
}

func (r *salesRepository) UpdateBuktiBayar(id uint, filePath string) error {
	return r.db.Model(&models.Penjualan{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"bukti_bayar":     filePath,
			"diperbarui_pada": time.Now(),
		}).Error
}
