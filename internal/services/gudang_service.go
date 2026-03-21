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

type GudangService interface {
	CreateGudang(req *dto.CreateGudangRequest) (*dto.GudangResponse, error)
	GetGudangByID(id uint) (*dto.GudangResponse, error)
	ListGudangs(req *dto.ListGudangRequest) (*dto.ListGudangResponse, error)
	UpdateGudang(id uint, req *dto.UpdateGudangRequest) (*dto.GudangResponse, error)
	DeleteGudang(id uint) error
}

type gudangService struct {
	repo repositories.GudangRepository
}

func NewGudangService(repo repositories.GudangRepository) GudangService {
	return &gudangService{repo}
}

func (s *gudangService) CreateGudang(req *dto.CreateGudangRequest) (*dto.GudangResponse, error) {
	// Validate unique code
	existing, err := s.repo.FindByCode(req.Kode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("kode gudang sudah digunakan")
	}

	gudang := &models.Gudang{
		Kode:       req.Kode,
		Nama:       req.Nama,
		Alamat:     req.Alamat,
		Keterangan: req.Keterangan,
		Aktif:      req.Aktif,
		DibuatPada: time.Now(),
		DiperbaruiPada: time.Now(),
	}

	if err := s.repo.Create(gudang); err != nil {
		return nil, err
	}

	return s.toResponse(gudang), nil
}

func (s *gudangService) GetGudangByID(id uint) (*dto.GudangResponse, error) {
	gudang, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("gudang tidak ditemukan")
		}
		return nil, err
	}

	return s.toResponse(gudang), nil
}

func (s *gudangService) ListGudangs(req *dto.ListGudangRequest) (*dto.ListGudangResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	filters := make(map[string]interface{})
	if req.Search != "" {
		filters["search"] = req.Search
	}
	if req.Aktif != nil {
		filters["aktif"] = *req.Aktif
	}

	gudangs, total, err := s.repo.List(filters, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.GudangResponse, len(gudangs))
	for i, gudang := range gudangs {
		responses[i] = *s.toResponse(&gudang)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &dto.ListGudangResponse{
		Gudangs:    responses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *gudangService) UpdateGudang(id uint, req *dto.UpdateGudangRequest) (*dto.GudangResponse, error) {
	gudang, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("gudang tidak ditemukan")
		}
		return nil, err
	}

	if req.Kode != nil {
		// Validate unique code if changed
		if *req.Kode != gudang.Kode {
			existing, err := s.repo.FindByCode(*req.Kode)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if existing != nil {
				return nil, errors.New("kode gudang sudah digunakan")
			}
			gudang.Kode = *req.Kode
		}
	}

	if req.Nama != nil {
		gudang.Nama = *req.Nama
	}
	if req.Alamat != nil {
		gudang.Alamat = *req.Alamat
	}
	if req.Keterangan != nil {
		gudang.Keterangan = *req.Keterangan
	}
	if req.Aktif != nil {
		gudang.Aktif = *req.Aktif
	}

	gudang.DiperbaruiPada = time.Now()

	if err := s.repo.Update(gudang); err != nil {
		return nil, err
	}

	return s.toResponse(gudang), nil
}

func (s *gudangService) DeleteGudang(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("gudang tidak ditemukan")
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *gudangService) toResponse(gudang *models.Gudang) *dto.GudangResponse {
	return &dto.GudangResponse{
		ID:         gudang.ID,
		Kode:       gudang.Kode,
		Nama:       gudang.Nama,
		Alamat:     gudang.Alamat,
		Keterangan: gudang.Keterangan,
		Aktif:      gudang.Aktif,
	}
}
