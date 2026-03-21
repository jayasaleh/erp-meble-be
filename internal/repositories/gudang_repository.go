package repositories

import (
	"real-erp-mebel/be/internal/models"

	"gorm.io/gorm"
)

type GudangRepository interface {
	Create(gudang *models.Gudang) error
	FindByID(id uint) (*models.Gudang, error)
	List(filters map[string]interface{}, page int, limit int) ([]models.Gudang, int64, error)
	Update(gudang *models.Gudang) error
	Delete(id uint) error
	FindByCode(code string) (*models.Gudang, error)
}

type gudangRepository struct {
	db *gorm.DB
}

func NewGudangRepository(db *gorm.DB) GudangRepository {
	return &gudangRepository{db}
}

func (r *gudangRepository) Create(gudang *models.Gudang) error {
	return r.db.Create(gudang).Error
}

func (r *gudangRepository) FindByID(id uint) (*models.Gudang, error) {
	var gudang models.Gudang
	err := r.db.First(&gudang, id).Error
	if err != nil {
		return nil, err
	}
	return &gudang, nil
}

func (r *gudangRepository) List(filters map[string]interface{}, page int, limit int) ([]models.Gudang, int64, error) {
	var gudangs []models.Gudang
	var total int64

	query := r.db.Model(&models.Gudang{})

	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("nama ILIKE ? OR kode ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if aktif, ok := filters["aktif"].(bool); ok {
		query = query.Where("aktif = ?", aktif)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&gudangs).Error
	if err != nil {
		return nil, 0, err
	}

	return gudangs, total, nil
}

func (r *gudangRepository) Update(gudang *models.Gudang) error {
	return r.db.Updates(gudang).Error
}

func (r *gudangRepository) Delete(id uint) error {
	return r.db.Delete(&models.Gudang{}, id).Error
}

func (r *gudangRepository) FindByCode(code string) (*models.Gudang, error) {
	var gudang models.Gudang
	err := r.db.Where("kode = ?", code).First(&gudang).Error
	if err != nil {
		return nil, err
	}
	return &gudang, nil
}
