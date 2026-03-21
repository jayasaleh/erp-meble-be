package services

import (
	"errors"
	"fmt"
	"math"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"real-erp-mebel/be/internal/repositories"
	"time"

	"gorm.io/gorm"
)

type SalesService interface {
	CreateSale(userID uint, req *dto.CreateSalesRequest) (*dto.SalesDetailResponse, error)
	GetSaleByID(id uint) (*dto.SalesDetailResponse, error)
	GetInvoice(id uint) (*dto.InvoiceResponse, error)
	ListSales(req *dto.ListSalesRequest) (*dto.ListSalesResponse, error)
	UpdateBuktiBayar(id uint, filePath string) error
}

type salesService struct {
	repo      repositories.SalesRepository
	stockRepo repositories.StockRepository
	batchRepo repositories.StockBatchRepository
}

func NewSalesService(
	repo repositories.SalesRepository,
	stockRepo repositories.StockRepository,
	batchRepo repositories.StockBatchRepository,
) SalesService {
	return &salesService{
		repo:      repo,
		stockRepo: stockRepo,
		batchRepo: batchRepo,
	}
}

// CreateSale adalah fungsi utama yang memproses transaksi penjualan POS.
// Semua operasi dijalankan dalam 1 DB transaction untuk menjamin atomicity.
//
// Alur FIFO:
//  1. Validasi stok tersedia per item
//  2. Untuk setiap item: ambil batches FIFO (terlama dulu), deduct, catat breakdown
//  3. Hitung COGS per item dari batch yang terpakai
//  4. Log pergerakan stok (tipe_referensi = "sales")
//  5. Buat barang_keluar header
//  6. Buat penjualan + item_penjualan + item_penjualan_batch
//  7. Commit
func (s *salesService) CreateSale(userID uint, req *dto.CreateSalesRequest) (*dto.SalesDetailResponse, error) {
	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// Nomor transaksi: TRX/YYYYMMDD/HHMMSS/userID
	nomorTransaksi := fmt.Sprintf("TRX/%s/%d", now.Format("20060102150405"), userID)

	// 1. Buat header BarangKeluar terlebih dahulu (untuk referensi di pergerakan stok)
	headerKeluar := models.BarangKeluar{
		NomorTransaksi: fmt.Sprintf("OUT/SALES/%s/%d", now.Format("20060102150405"), userID),
		Alasan:         "penjualan",
		TipeReferensi:  "sales",
		DibuatOleh:     userID,
		DibuatPada:     now,
		DiperbaruiPada: now,
	}
	if err := s.stockRepo.CreateStockOut(tx, &headerKeluar); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal membuat barang keluar: %w", err)
	}

	// 2. Validasi stok semua item sebelum memproses
	for _, itemReq := range req.Items {
		batches, err := s.batchRepo.GetAvailableBatches(tx, itemReq.IDProduk, req.IDGudang)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal mengambil batch produk %d: %w", itemReq.IDProduk, err)
		}
		totalAvailable := 0
		for _, b := range batches {
			totalAvailable += b.JumlahSaatIni
		}
		if totalAvailable < itemReq.Jumlah {
			tx.Rollback()
			return nil, fmt.Errorf("stok tidak cukup untuk produk ID %d (dibutuhkan: %d, tersedia: %d)",
				itemReq.IDProduk, itemReq.Jumlah, totalAvailable)
		}
	}

	// 3. Build header Penjualan
	sale := models.Penjualan{
		NomorTransaksi:   nomorTransaksi,
		IDGudang:         req.IDGudang,
		NamaPelanggan:    req.NamaPelanggan,
		KontakPelanggan:  req.KontakPelanggan,
		MetodePembayaran: req.MetodePembayaran,
		JumlahPembayaran: req.JumlahPembayaran,
		Status:           "completed",
		CatatanInternal:  req.CatatanInternal,
		IDKasir:          userID,
		DibuatPada:       now,
		DiperbaruiPada:   now,
	}

	// 4. Proses setiap item: FIFO deduct, hitung COGS, buat item
	var grandSubtotal, grandDiskon, grandTotal, grandHargaModal float64
	var saleItems []models.ItemPenjualan

	for _, itemReq := range req.Items {
		// Hitung diskon item
		var jumlahDiskon float64
		if itemReq.PersenDiskon != nil && *itemReq.PersenDiskon > 0 {
			jumlahDiskon = math.Round(itemReq.HargaSatuan*float64(itemReq.Jumlah)*(*itemReq.PersenDiskon)/100*100) / 100
		}
		subtotalItem := math.Round(itemReq.HargaSatuan*float64(itemReq.Jumlah)*100)/100 - jumlahDiskon

		// FIFO: Ambil dan deduct batches
		batches, _ := s.batchRepo.GetAvailableBatches(tx, itemReq.IDProduk, req.IDGudang)
		remainingQty := itemReq.Jumlah
		var totalCOGSItem float64
		var batchUsageRecords []models.ItemPenjualanBatch

		for i := range batches {
			if remainingQty == 0 {
				break
			}
			batch := &batches[i]
			deductQty := 0
			if batch.JumlahSaatIni >= remainingQty {
				deductQty = remainingQty
				remainingQty = 0
			} else {
				deductQty = batch.JumlahSaatIni
				remainingQty -= batch.JumlahSaatIni
			}

			totalCOGSBatch := math.Round(batch.HargaModal*float64(deductQty)*100) / 100
			totalCOGSItem += totalCOGSBatch

			// Catat breakdown batch untuk item ini
			batchUsageRecords = append(batchUsageRecords, models.ItemPenjualanBatch{
				IDBatch:    batch.ID,
				Jumlah:     deductQty,
				HargaModal: batch.HargaModal,
				TotalModal: totalCOGSBatch,
				DibuatPada: now,
			})

			// Deduct batch
			batch.JumlahSaatIni -= deductQty
			if err := s.batchRepo.Update(tx, batch); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("gagal update batch #%d: %w", batch.ID, err)
			}

			// Log pergerakan stok per batch
			movement := models.PergerakanStok{
				IDProduk:       itemReq.IDProduk,
				IDGudang:       req.IDGudang,
				IDBatch:        &batch.ID,
				TipePergerakan: "out",
				TipeReferensi:  "sales",
				IDReferensi:    &headerKeluar.ID,
				Jumlah:         -deductQty,
				IDPengguna:     userID,
				Keterangan:     fmt.Sprintf("Penjualan %s (Batch #%d, HPP: %.2f)", nomorTransaksi, batch.ID, batch.HargaModal),
				DibuatPada:     now,
			}
			if err := s.stockRepo.CreateStockMovement(tx, &movement); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("gagal log pergerakan stok: %w", err)
			}
		}

		// Update total stok inventori
		if err := s.stockRepo.UpdateStockBalance(tx, itemReq.IDProduk, req.IDGudang, -itemReq.Jumlah); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal update stok inventori produk %d: %w", itemReq.IDProduk, err)
		}

		// COGS per unit (rata-rata tertimbang)
		hargaModalPerUnit := 0.0
		if itemReq.Jumlah > 0 {
			hargaModalPerUnit = math.Round(totalCOGSItem/float64(itemReq.Jumlah)*100) / 100
		}

		// Build ItemPenjualan
		item := models.ItemPenjualan{
			IDProduk:     itemReq.IDProduk,
			IDGudang:     req.IDGudang,
			Jumlah:       itemReq.Jumlah,
			HargaSatuan:  itemReq.HargaSatuan,
			HargaModal:   hargaModalPerUnit,
			PersenDiskon: itemReq.PersenDiskon,
			JumlahDiskon: jumlahDiskon,
			Subtotal:     subtotalItem,
			TotalModal:   totalCOGSItem,
			DibuatPada:   now,
			BatchUsage:   batchUsageRecords,
		}
		saleItems = append(saleItems, item)

		grandSubtotal += itemReq.HargaSatuan * float64(itemReq.Jumlah)
		grandDiskon += jumlahDiskon
		grandTotal += subtotalItem
		grandHargaModal += totalCOGSItem
	}

	// 5. Hitung kembalian
	kembalian := 0.0
	if req.MetodePembayaran == "cash" && req.JumlahPembayaran >= grandTotal {
		kembalian = math.Round((req.JumlahPembayaran-grandTotal)*100) / 100
	}

	sale.Subtotal = math.Round(grandSubtotal*100) / 100
	sale.JumlahDiskon = math.Round(grandDiskon*100) / 100
	sale.Total = math.Round(grandTotal*100) / 100
	sale.TotalHargaModal = math.Round(grandHargaModal*100) / 100
	sale.JumlahKembalian = kembalian
	sale.Items = saleItems

	// Update IDReferensi di barang_keluar setelah sale dibuat (chicken-egg, skip untuk simplicity)
	// Sudah ada nomor_transaksi di keterangan movement

	// 6. Simpan penjualan (GORM akan cascade insert Items dan BatchUsage)
	if err := s.repo.Create(tx, &sale); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal menyimpan transaksi penjualan: %w", err)
	}

	// 7. Update IDReferensi pada barang_keluar agar link ke penjualan
	if err := tx.Model(&headerKeluar).Update("id_referensi", sale.ID).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal link barang_keluar ke penjualan: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("gagal commit transaksi: %w", err)
	}

	// Ambil data lengkap untuk response
	return s.GetSaleByID(sale.ID)
}

func (s *salesService) GetSaleByID(id uint) (*dto.SalesDetailResponse, error) {
	sale, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaksi penjualan tidak ditemukan")
		}
		return nil, err
	}
	return mapSaleToDetailResponse(sale), nil
}

func (s *salesService) GetInvoice(id uint) (*dto.InvoiceResponse, error) {
	sale, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaksi penjualan tidak ditemukan")
		}
		return nil, err
	}
	return mapSaleToInvoice(sale), nil
}

func (s *salesService) ListSales(req *dto.ListSalesRequest) (*dto.ListSalesResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}

	sales, total, err := s.repo.FindAll(req)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	var responses []dto.SalesResponse
	for _, sale := range sales {
		laba := sale.Total - sale.TotalHargaModal
		responses = append(responses, dto.SalesResponse{
			ID:               sale.ID,
			NomorTransaksi:   sale.NomorTransaksi,
			NamaGudang:       sale.Gudang.Nama,
			NamaPelanggan:    sale.NamaPelanggan,
			KontakPelanggan:  sale.KontakPelanggan,
			Total:            sale.Total,
			TotalHargaModal:  sale.TotalHargaModal,
			Laba:             math.Round(laba*100) / 100,
			MetodePembayaran: sale.MetodePembayaran,
			Status:           sale.Status,
			NamaKasir:        sale.Kasir.Nama,
			DibuatPada:       sale.DibuatPada,
		})
	}

	return &dto.ListSalesResponse{
		Sales:      responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *salesService) UpdateBuktiBayar(id uint, filePath string) error {
	// Pastikan penjualan ada
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaksi penjualan tidak ditemukan")
		}
		return err
	}
	return s.repo.UpdateBuktiBayar(id, filePath)
}

// ===========================
// MAPPING HELPERS
// ===========================

func mapSaleToDetailResponse(sale *models.Penjualan) *dto.SalesDetailResponse {
	var items []dto.SalesItemResponse
	for _, item := range sale.Items {
		var batchUsages []dto.SalesBatchUsageResponse
		for _, bu := range item.BatchUsage {
			batchUsages = append(batchUsages, dto.SalesBatchUsageResponse{
				IDBatch:    bu.IDBatch,
				Jumlah:     bu.Jumlah,
				HargaModal: bu.HargaModal,
				TotalModal: bu.TotalModal,
			})
		}
		laba := item.Subtotal - item.TotalModal
		items = append(items, dto.SalesItemResponse{
			ID:           item.ID,
			IDProduk:     item.IDProduk,
			SKUProduk:    item.Produk.SKU,
			NamaProduk:   item.Produk.Nama,
			Jumlah:       item.Jumlah,
			HargaSatuan:  item.HargaSatuan,
			HargaModal:   item.HargaModal,
			PersenDiskon: item.PersenDiskon,
			JumlahDiskon: item.JumlahDiskon,
			Subtotal:     item.Subtotal,
			TotalModal:   item.TotalModal,
			Laba:         math.Round(laba*100) / 100,
			BatchUsage:   batchUsages,
		})
	}

	laba := sale.Total - sale.TotalHargaModal
	return &dto.SalesDetailResponse{
		ID:               sale.ID,
		NomorTransaksi:   sale.NomorTransaksi,
		IDGudang:         sale.IDGudang,
		NamaGudang:       sale.Gudang.Nama,
		NamaPelanggan:    sale.NamaPelanggan,
		KontakPelanggan:  sale.KontakPelanggan,
		Subtotal:         sale.Subtotal,
		JumlahDiskon:     sale.JumlahDiskon,
		Total:            sale.Total,
		TotalHargaModal:  sale.TotalHargaModal,
		Laba:             math.Round(laba*100) / 100,
		MetodePembayaran: sale.MetodePembayaran,
		JumlahPembayaran: sale.JumlahPembayaran,
		JumlahKembalian:  sale.JumlahKembalian,
		BuktiBayar:       sale.BuktiBayar,
		Status:           sale.Status,
		CatatanInternal:  sale.CatatanInternal,
		IDKasir:          sale.IDKasir,
		NamaKasir:        sale.Kasir.Nama,
		DibuatPada:       sale.DibuatPada,
		Items:            items,
	}
}

func mapSaleToInvoice(sale *models.Penjualan) *dto.InvoiceResponse {
	var items []dto.InvoiceItemResponse
	for i, item := range sale.Items {
		items = append(items, dto.InvoiceItemResponse{
			NoProduk:    i + 1,
			SKU:         item.Produk.SKU,
			NamaProduk:  item.Produk.Nama,
			Jumlah:      item.Jumlah,
			HargaSatuan: item.HargaSatuan,
			Diskon:      item.JumlahDiskon,
			Subtotal:    item.Subtotal,
		})
	}
	return &dto.InvoiceResponse{
		NomorInvoice:     sale.NomorTransaksi,
		TanggalInvoice:   sale.DibuatPada,
		NamaGudang:       sale.Gudang.Nama,
		NamaPelanggan:    sale.NamaPelanggan,
		KontakPelanggan:  sale.KontakPelanggan,
		NamaKasir:        sale.Kasir.Nama,
		Items:            items,
		Subtotal:         sale.Subtotal,
		TotalDiskon:      sale.JumlahDiskon,
		Total:            sale.Total,
		MetodePembayaran: sale.MetodePembayaran,
		JumlahPembayaran: sale.JumlahPembayaran,
		JumlahKembalian:  sale.JumlahKembalian,
		Status:           sale.Status,
	}
}
