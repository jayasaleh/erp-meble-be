package repositories

import (
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReturnRepository interface {
	// Retur Penjualan
	CreateReturPenjualan(tx *gorm.DB, retur *models.ReturPenjualan) error
	FindReturPenjualanByID(id uint) (*models.ReturPenjualan, error)
	FindAllReturPenjualan(req *dto.ListReturPenjualanRequest) ([]models.ReturPenjualan, int64, error)
	UpdateStatusReturPenjualan(tx *gorm.DB, id uint, status string, approvedBy uint) error

	// Retur Pembelian
	CreateReturPembelian(tx *gorm.DB, retur *models.ReturPembelian) error
	FindReturPembelianByID(id uint) (*models.ReturPembelian, error)
	FindAllReturPembelian(req *dto.ListReturPembelianRequest) ([]models.ReturPembelian, int64, error)
	UpdateStatusReturPembelian(tx *gorm.DB, id uint, status string, approvedBy uint) error

	// Utility
	BeginTx() *gorm.DB
}

type returnRepository struct {
	db *gorm.DB
}

func NewReturnRepository(db *gorm.DB) ReturnRepository {
	return &returnRepository{db: db}
}

func (r *returnRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

// ===========================
// RETUR PENJUALAN
// ===========================

func (r *returnRepository) CreateReturPenjualan(tx *gorm.DB, retur *models.ReturPenjualan) error {
	return tx.Create(retur).Error
}

func (r *returnRepository) FindReturPenjualanByID(id uint) (*models.ReturPenjualan, error) {
	var retur models.ReturPenjualan
	err := r.db.
		Preload("Penjualan").
		Preload("DiprosesOlehPengguna").
		Preload("DisetujuiOlehPengguna").
		Preload("Items").
		Preload("Items.Produk").
		Preload("Items.Gudang").
		First(&retur, id).Error
	if err != nil {
		return nil, err
	}
	return &retur, nil
}

func (r *returnRepository) FindAllReturPenjualan(req *dto.ListReturPenjualanRequest) ([]models.ReturPenjualan, int64, error) {
	var returs []models.ReturPenjualan
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

	query := r.db.Model(&models.ReturPenjualan{})

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.TanggalDari != nil {
		startOfDay := time.Date(req.TanggalDari.Year(), req.TanggalDari.Month(), req.TanggalDari.Day(), 0, 0, 0, 0, req.TanggalDari.Location())
		query = query.Where("dibuat_pada >= ?", startOfDay)
	}
	if req.TanggalSampai != nil {
		endOfDay := time.Date(req.TanggalSampai.Year(), req.TanggalSampai.Month(), req.TanggalSampai.Day(), 23, 59, 59, 999999999, req.TanggalSampai.Location())
		query = query.Where("dibuat_pada <= ?", endOfDay)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Penjualan").
		Preload("DiprosesOlehPengguna").
		Order("dibuat_pada DESC").
		Limit(limit).Offset(offset).
		Find(&returs).Error

	return returs, total, err
}

func (r *returnRepository) UpdateStatusReturPenjualan(tx *gorm.DB, id uint, status string, approvedBy uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":          status,
		"diperbarui_pada": now,
	}
	if status == "approved" || status == "completed" {
		updates["disetujui_oleh"] = approvedBy
		updates["disetujui_pada"] = now
	}
	return tx.Model(&models.ReturPenjualan{}).Where("id = ?", id).Updates(updates).Error
}

// ===========================
// RETUR PEMBELIAN
// ===========================

func (r *returnRepository) CreateReturPembelian(tx *gorm.DB, retur *models.ReturPembelian) error {
	return tx.Create(retur).Error
}

func (r *returnRepository) FindReturPembelianByID(id uint) (*models.ReturPembelian, error) {
	var retur models.ReturPembelian
	err := r.db.
		Preload("Pemasok").
		Preload("DibuatOlehPengguna").
		Preload("DisetujuiOlehPengguna").
		Preload("Items").
		Preload("Items.Produk").
		Preload("Items.Gudang").
		First(&retur, id).Error
	if err != nil {
		return nil, err
	}
	return &retur, nil
}

func (r *returnRepository) FindAllReturPembelian(req *dto.ListReturPembelianRequest) ([]models.ReturPembelian, int64, error) {
	var returs []models.ReturPembelian
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

	query := r.db.Model(&models.ReturPembelian{})

	if req.IDPemasok != nil {
		query = query.Where("id_supplier = ?", *req.IDPemasok)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Pemasok").
		Preload("DibuatOlehPengguna").
		Order("dibuat_pada DESC").
		Limit(limit).Offset(offset).
		Find(&returs).Error

	return returs, total, err
}

func (r *returnRepository) UpdateStatusReturPembelian(tx *gorm.DB, id uint, status string, approvedBy uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":          status,
		"diperbarui_pada": now,
	}
	if status == "approved" || status == "completed" {
		updates["disetujui_oleh"] = approvedBy
		updates["disetujui_pada"] = now
	}
	return tx.Model(&models.ReturPembelian{}).Where("id = ?", id).Updates(updates).Error
}
