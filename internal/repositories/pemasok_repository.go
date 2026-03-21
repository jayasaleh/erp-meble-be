package repositories

import (
	"real-erp-mebel/be/internal/models"

	"gorm.io/gorm"
)

type PemasokRepository interface {
	Create(pemasok *models.Pemasok) error
	FindByID(id uint) (*models.Pemasok, error)
	List(filters map[string]interface{}, page, limit int) ([]models.Pemasok, int64, error)
	Update(pemasok *models.Pemasok) error
	Delete(id uint) error
}

type pemasokRepository struct {
	db *gorm.DB
}

func NewPemasokRepository(db *gorm.DB) PemasokRepository {
	return &pemasokRepository{db: db}
}

func (r *pemasokRepository) Create(pemasok *models.Pemasok) error {
	return r.db.Create(pemasok).Error
}

func (r *pemasokRepository) FindByID(id uint) (*models.Pemasok, error) {
	var pemasok models.Pemasok
	if err := r.db.First(&pemasok, id).Error; err != nil {
		return nil, err
	}
	return &pemasok, nil
}

func (r *pemasokRepository) List(filters map[string]interface{}, page, limit int) ([]models.Pemasok, int64, error) {
	var pemasok []models.Pemasok
	var total int64

	query := r.db.Model(&models.Pemasok{})

	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("nama ILIKE ? OR kontak ILIKE ? OR email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if aktif, ok := filters["aktif"].(bool); ok {
		query = query.Where("aktif = ?", aktif)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("nama ASC").Offset(offset).Limit(limit).Find(&pemasok).Error
	return pemasok, total, err
}

func (r *pemasokRepository) Update(pemasok *models.Pemasok) error {
	return r.db.Save(pemasok).Error
}

func (r *pemasokRepository) Delete(id uint) error {
	return r.db.Delete(&models.Pemasok{}, id).Error
}
