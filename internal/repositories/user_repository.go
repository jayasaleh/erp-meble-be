package repositories

import (
	"real-erp-mebel/be/internal/database"
	"real-erp-mebel/be/internal/models"

	"gorm.io/gorm"
)

// UserRepository adalah interface untuk user repository
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll(page, pageSize int, search, peran string, aktif *bool) ([]models.User, int64, error)
	Update(user *models.User) error
	UpdatePassword(id uint, hashedPassword string) error
	Delete(id uint) error
	Count(search, peran string, aktif *bool) (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance UserRepository baru
func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.DB,
	}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(page, pageSize int, search, peran string, aktif *bool) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	// Apply filters
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nama ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	if peran != "" {
		query = query.Where("peran = ?", peran)
	}

	if aktif != nil {
		query = query.Where("aktif = ?", *aktif)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("dibuat_pada DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) Count(search, peran string, aktif *bool) (int64, error) {
	var count int64

	query := r.db.Model(&models.User{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nama ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	if peran != "" {
		query = query.Where("peran = ?", peran)
	}

	if aktif != nil {
		query = query.Where("aktif = ?", *aktif)
	}

	err := query.Count(&count).Error
	return count, err
}
