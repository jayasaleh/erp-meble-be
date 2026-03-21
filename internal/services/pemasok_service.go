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

type PemasokService interface {
	CreatePemasok(req *dto.CreatePemasokRequest) (*dto.PemasokResponse, error)
	GetPemasokByID(id uint) (*dto.PemasokResponse, error)
	ListPemasok(req *dto.ListPemasokRequest) (*dto.ListPemasokResponse, error)
	UpdatePemasok(id uint, req *dto.UpdatePemasokRequest) (*dto.PemasokResponse, error)
	DeletePemasok(id uint) error
}

type pemasokService struct {
	repo repositories.PemasokRepository
}

func NewPemasokService(repo repositories.PemasokRepository) PemasokService {
	return &pemasokService{repo: repo}
}

func (s *pemasokService) CreatePemasok(req *dto.CreatePemasokRequest) (*dto.PemasokResponse, error) {
	pemasok := &models.Pemasok{
		Nama:           req.Nama,
		Kontak:         req.Kontak,
		Telepon:        req.Telepon,
		Email:          req.Email,
		Alamat:         req.Alamat,
		Aktif:          true,
		DibuatPada:     time.Now(),
		DiperbaruiPada: time.Now(),
	}

	if req.Aktif {
		pemasok.Aktif = req.Aktif
	}

	if err := s.repo.Create(pemasok); err != nil {
		return nil, err
	}

	return s.toResponse(pemasok), nil
}

func (s *pemasokService) GetPemasokByID(id uint) (*dto.PemasokResponse, error) {
	pemasok, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pemasok tidak ditemukan")
		}
		return nil, err
	}

	return s.toResponse(pemasok), nil
}

func (s *pemasokService) ListPemasok(req *dto.ListPemasokRequest) (*dto.ListPemasokResponse, error) {
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

	pemasokList, total, err := s.repo.List(filters, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PemasokResponse, len(pemasokList))
	for i, p := range pemasokList {
		responses[i] = *s.toResponse(&p)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &dto.ListPemasokResponse{
		Pemasok:    responses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *pemasokService) UpdatePemasok(id uint, req *dto.UpdatePemasokRequest) (*dto.PemasokResponse, error) {
	pemasok, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pemasok tidak ditemukan")
		}
		return nil, err
	}

	if req.Nama != nil {
		pemasok.Nama = *req.Nama
	}
	if req.Kontak != nil {
		pemasok.Kontak = *req.Kontak
	}
	if req.Telepon != nil {
		pemasok.Telepon = *req.Telepon
	}
	if req.Email != nil {
		pemasok.Email = *req.Email
	}
	if req.Alamat != nil {
		pemasok.Alamat = *req.Alamat
	}
	if req.Aktif != nil {
		pemasok.Aktif = *req.Aktif
	}

	pemasok.DiperbaruiPada = time.Now()

	if err := s.repo.Update(pemasok); err != nil {
		return nil, err
	}

	return s.toResponse(pemasok), nil
}

func (s *pemasokService) DeletePemasok(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("pemasok tidak ditemukan")
		}
		return err
	}

	// TODO: Check constraints (products linked to this supplier)
	// For now, allow soft delete.

	return s.repo.Delete(id)
}

func (s *pemasokService) toResponse(p *models.Pemasok) *dto.PemasokResponse {
	return &dto.PemasokResponse{
		ID:             p.ID,
		Nama:           p.Nama,
		Kontak:         p.Kontak,
		Telepon:        p.Telepon,
		Email:          p.Email,
		Alamat:         p.Alamat,
		Aktif:          p.Aktif,
		DibuatPada:     p.DibuatPada,
		DiperbaruiPada: p.DiperbaruiPada,
	}
}
