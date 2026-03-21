package services

import (
	"errors"
	"math"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"real-erp-mebel/be/internal/repositories"
	"time"

	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(req *dto.CreateProductRequest, userID uint) (*dto.ProductResponse, error)
	GetProductByID(id uint) (*dto.ProductResponse, error)
	ListProducts(req *dto.ProductListRequest) (*dto.ProductListResponse, error)
	UpdateProduct(id uint, req *dto.UpdateProductRequest, userID uint) (*dto.ProductResponse, error)
	DeleteProduct(id uint) error
}

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

// CreateProduct membuat produk baru
func (s *productService) CreateProduct(req *dto.CreateProductRequest, userID uint) (*dto.ProductResponse, error) {
	// Validasi SKU unique
	existingSKU, err := s.productRepo.FindBySKU(req.SKU)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingSKU != nil {
		return nil, errors.New("SKU sudah digunakan")
	}

	// Validasi Barcode unique (jika ada)
	if req.Barcode != nil && *req.Barcode != "" {
		existingBarcode, err := s.productRepo.FindByBarcode(*req.Barcode)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingBarcode != nil {
			return nil, errors.New("barcode sudah digunakan")
		}
	}

	// Validasi harga jual >= harga modal
	if req.HargaJual < req.HargaModal {
		return nil, errors.New("harga jual tidak boleh lebih kecil dari harga modal")
	}

	product := &models.Produk{
		SKU:            req.SKU,
		Barcode:        req.Barcode,
		Nama:           req.Nama,
		Kategori:       req.Kategori,
		Merek:          req.Merek,
		IDPemasok:      req.IDPemasok,
		HargaModal:     req.HargaModal,
		HargaJual:      req.HargaJual,
		StokMinimum:    req.StokMinimum,
		IzinDiskon:     req.IzinDiskon,
		Aktif:          req.Aktif,
		DibuatOleh:     userID,
		DiupdateOleh:   userID,
		DibuatPada:     time.Now(),
		DiperbaruiPada: time.Now(),
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return s.GetProductByID(product.ID)
}

// GetProductByID mengambil produk berdasarkan ID
func (s *productService) GetProductByID(id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err == nil && product != nil {
		// Manual Preload hack if repo doesn't support it yet, or assume repo does it.
		// Actually, let's verify repo later. For now, we rely on Repo to preload or we do it here if we had DB access (we don't).
		// Wait, ProductService struct doesn't have *gorm.DB. Repo should handle preloading.
		// Let's modify Repo next step.
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	// Get total stock
	totalStock, err := s.productRepo.GetStockByProductID(id)
	if err != nil {
		totalStock = 0
	}

	return s.toProductResponse(product, totalStock), nil
}

// ListProducts mengambil daftar produk dengan filter
func (s *productService) ListProducts(req *dto.ProductListRequest) (*dto.ProductListResponse, error) {
	// Default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.Search != "" {
		filters["search"] = req.Search
	}
	if req.Kategori != "" {
		filters["kategori"] = req.Kategori
	}
	if req.Merek != "" {
		filters["merek"] = req.Merek
	}
	if req.IDPemasok != nil {
		filters["id_pemasok"] = *req.IDPemasok
	}
	if req.Aktif != nil {
		filters["aktif"] = *req.Aktif
	}

	products, total, err := s.productRepo.List(filters, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	// Convert to response
	productResponses := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		// Get stock for each product
		stock, _ := s.productRepo.GetStockByProductID(product.ID)

		// Filter stok rendah jika diminta
		if req.StokRendah && stock >= product.StokMinimum {
			continue
		}

		productResponses = append(productResponses, *s.toProductResponse(&product, stock))
	}

	// Recalculate total if filter stok_rendah is applied
	if req.StokRendah {
		total = int64(len(productResponses))
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &dto.ProductListResponse{
		Products:   productResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateProduct mengupdate produk
func (s *productService) UpdateProduct(id uint, req *dto.UpdateProductRequest, userID uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	// Update fields yang ada
	if req.SKU != nil {
		// Validasi SKU unique
		existingSKU, err := s.productRepo.FindBySKU(*req.SKU)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingSKU != nil && existingSKU.ID != id {
			return nil, errors.New("SKU sudah digunakan")
		}
		product.SKU = *req.SKU
	}

	if req.Barcode != nil {
		if *req.Barcode != "" {
			// Validasi Barcode unique
			existingBarcode, err := s.productRepo.FindByBarcode(*req.Barcode)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if existingBarcode != nil && existingBarcode.ID != id {
				return nil, errors.New("barcode sudah digunakan")
			}
		}
		product.Barcode = req.Barcode
	}

	if req.Nama != nil {
		product.Nama = *req.Nama
	}

	if req.Kategori != nil {
		product.Kategori = *req.Kategori
	}

	if req.Merek != nil {
		product.Merek = *req.Merek
	}

	if req.IDPemasok != nil {
		product.IDPemasok = req.IDPemasok
	}

	if req.HargaModal != nil {
		product.HargaModal = *req.HargaModal
	}

	if req.HargaJual != nil {
		product.HargaJual = *req.HargaJual
	}

	// Validasi harga jual >= harga modal
	if product.HargaJual < product.HargaModal {
		return nil, errors.New("harga jual tidak boleh lebih kecil dari harga modal")
	}

	if req.StokMinimum != nil {
		product.StokMinimum = *req.StokMinimum
	}

	if req.IzinDiskon != nil {
		product.IzinDiskon = *req.IzinDiskon
	}

	if req.Aktif != nil {
		product.Aktif = *req.Aktif
	}

	product.DiupdateOleh = userID

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return s.GetProductByID(id)
}

// DeleteProduct menghapus produk (soft delete)
func (s *productService) DeleteProduct(id uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("produk tidak ditemukan")
		}
		return err
	}

	// Check if product has stock
	stock, err := s.productRepo.GetStockByProductID(id)
	if err != nil {
		return err
	}

	if stock > 0 {
		return errors.New("tidak dapat menghapus produk yang masih memiliki stok")
	}

	return s.productRepo.Delete(product.ID)
}

// Helper: Convert model to response DTO
func (s *productService) toProductResponse(product *models.Produk, totalStock int) *dto.ProductResponse {
	response := &dto.ProductResponse{
		ID:             product.ID,
		SKU:            product.SKU,
		Barcode:        product.Barcode,
		Nama:           product.Nama,
		Kategori:       product.Kategori,
		Merek:          product.Merek,
		IDPemasok:      product.IDPemasok,
		HargaModal:     product.HargaModal,
		HargaJual:      product.HargaJual,
		StokMinimum:    product.StokMinimum,
		IzinDiskon:     product.IzinDiskon,
		Aktif:          product.Aktif,
		JumlahStok:     totalStock,
		DibuatOleh:     product.DibuatOleh,
		DiupdateOleh:   product.DiupdateOleh,
		DibuatPada:     product.DibuatPada,
		DiperbaruiPada: product.DiperbaruiPada,
	}

	if product.Pembuat != nil {
		namaPembuat := product.Pembuat.Nama
		response.NamaPembuat = &namaPembuat
	}

	if product.Pengupdate != nil {
		namaPengupdate := product.Pengupdate.Nama
		response.NamaPengupdate = &namaPengupdate
	}

	// Add supplier name if loaded
	if product.Pemasok != nil {
		namaPemasok := product.Pemasok.Nama
		response.NamaPemasok = &namaPemasok
	}

	// Add images if loaded
	if len(product.Images) > 0 {
		images := make([]dto.ProductImageResponse, 0, len(product.Images))
		for _, img := range product.Images {
			images = append(images, dto.ProductImageResponse{
				ID:          img.ID,
				PathGambar:  img.PathGambar,
				GambarUtama: img.GambarUtama,
				Urutan:      img.Urutan,
			})
		}
		response.Images = images
	}

	return response
}
