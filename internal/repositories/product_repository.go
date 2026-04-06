package repositories

import (
	"real-erp-mebel/be/internal/models"
	"time"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Produk) error
	FindByID(id uint) (*models.Produk, error)
	FindBySKU(sku string) (*models.Produk, error)
	FindByBarcode(barcode string) (*models.Produk, error)
	List(filters map[string]interface{}, page, limit int) ([]models.Produk, int64, error)
	Update(product *models.Produk) error
	Delete(id uint) error
	GetStockByProductID(productID uint) (int, error)
	GetImageCount(productID uint) (int64, error)
	SaveProductImages(productID uint, images []models.GambarProduk) error
	DeleteProductImage(productID uint, imageID uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create membuat produk baru
func (r *productRepository) Create(product *models.Produk) error {
	return r.db.Create(product).Error
}

// FindByID mencari produk berdasarkan ID
func (r *productRepository) FindByID(id uint) (*models.Produk, error) {
	var product models.Produk
	err := r.db.Preload("Pemasok").Preload("Images").Preload("Pembuat").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindBySKU mencari produk berdasarkan SKU
func (r *productRepository) FindBySKU(sku string) (*models.Produk, error) {
	var product models.Produk
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindByBarcode mencari produk berdasarkan barcode
func (r *productRepository) FindByBarcode(barcode string) (*models.Produk, error) {
	var product models.Produk
	err := r.db.Where("barcode = ?", barcode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// List mengambil daftar produk dengan filter dan pagination
func (r *productRepository) List(filters map[string]interface{}, page, limit int) ([]models.Produk, int64, error) {
	var products []models.Produk
	var total int64

	query := r.db.Model(&models.Produk{})

	// Apply filters
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("nama ILIKE ? OR sku ILIKE ? OR barcode ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if kategori, ok := filters["kategori"].(string); ok && kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}

	if merek, ok := filters["merek"].(string); ok && merek != "" {
		query = query.Where("merek = ?", merek)
	}

	if idPemasok, ok := filters["id_pemasok"].(uint); ok && idPemasok > 0 {
		query = query.Where("id_supplier = ?", idPemasok)
	}

	if aktif, ok := filters["aktif"].(bool); ok {
		query = query.Where("aktif = ?", aktif)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.
		Preload("Pemasok").
		Preload("Images").
		Preload("Pembuat").
		Order("dibuat_pada DESC").
		Offset(offset).
		Limit(limit).
		Find(&products).Error

	return products, total, err
}

// Update mengupdate produk
func (r *productRepository) Update(product *models.Produk) error {
	// Gunakan Updates untuk update spesifik field dan menghindari zero values pada field lain yang tidak diubah
	// Kita set DiperbaruiPada secara manual
	product.DiperbaruiPada = time.Now()
	
	return r.db.Model(product).Updates(map[string]interface{}{
		"sku":             product.SKU,
		"barcode":         product.Barcode,
		"nama":            product.Nama,
		"kategori":        product.Kategori,
		"merek":           product.Merek,
		"id_supplier":     product.IDPemasok,
		"harga_modal":     product.HargaModal,
		"harga_jual":      product.HargaJual,
		"stok_minimum":    product.StokMinimum,
		"izin_diskon":     product.IzinDiskon,
		"aktif":           product.Aktif,
		"diupdate_oleh":   product.DiupdateOleh,
		"diperbarui_pada": product.DiperbaruiPada,
	}).Error
}

// Delete menghapus produk (soft delete)
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Produk{}, id).Error
}

// GetStockByProductID mengambil total stok produk dari semua gudang
func (r *productRepository) GetStockByProductID(productID uint) (int, error) {
	var totalStock int64
	err := r.db.Model(&models.StokInventori{}).
		Where("id_produk = ?", productID).
		Select("COALESCE(SUM(jumlah), 0)").
		Scan(&totalStock).Error

	return int(totalStock), err
}

// GetImageCount menghitung jumlah gambar produk saat ini
func (r *productRepository) GetImageCount(productID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.GambarProduk{}).Where("id_produk = ?", productID).Count(&count).Error
	return count, err
}

// SaveProductImages menyimpan gambar produk tanpa menghapus gambar lama secara replace
func (r *productRepository) SaveProductImages(productID uint, images []models.GambarProduk) error {
	// Insert gambar baru
	if len(images) > 0 {
		return r.db.Create(&images).Error
	}
	return nil
}

// DeleteProductImage menghapus spesifik gambar produk
func (r *productRepository) DeleteProductImage(productID uint, imageID uint) error {
	return r.db.Where("id_produk = ? AND id = ?", productID, imageID).Delete(&models.GambarProduk{}).Error
}
