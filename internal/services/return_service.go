package services

import (
	"errors"
	"fmt"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"real-erp-mebel/be/internal/repositories"
	"time"

	"gorm.io/gorm"
)

type ReturnService interface {
	// Retur Penjualan (Customer → Toko)
	CreateReturPenjualan(userID uint, req *dto.CreateReturPenjualanRequest) (*dto.ReturPenjualanResponse, error)
	GetReturPenjualanByID(id uint) (*dto.ReturPenjualanResponse, error)
	ListReturPenjualan(req *dto.ListReturPenjualanRequest) ([]dto.ReturPenjualanResponse, int64, error)
	ApproveReturPenjualan(id, approvedByUserID uint) error // Stok masuk kembali

	// Retur Pembelian (Toko → Vendor)
	CreateReturPembelian(userID uint, req *dto.CreateReturPembelianRequest) (*dto.ReturPembelianResponse, error)
	GetReturPembelianByID(id uint) (*dto.ReturPembelianResponse, error)
	ListReturPembelian(req *dto.ListReturPembelianRequest) ([]dto.ReturPembelianResponse, int64, error)
	ApproveReturPembelian(id, approvedByUserID uint) error // Stok keluar via FIFO
}

type returnService struct {
	repo      repositories.ReturnRepository
	stockRepo repositories.StockRepository
	batchRepo repositories.StockBatchRepository
	salesRepo repositories.SalesRepository
}

func NewReturnService(
	repo repositories.ReturnRepository,
	stockRepo repositories.StockRepository,
	batchRepo repositories.StockBatchRepository,
	salesRepo repositories.SalesRepository,
) ReturnService {
	return &returnService{
		repo:      repo,
		stockRepo: stockRepo,
		batchRepo: batchRepo,
		salesRepo: salesRepo,
	}
}

// ===========================
// RETUR PENJUALAN
// ===========================

// CreateReturPenjualan membuat dokumen retur dari customer. Status awal: pending.
// Stok BELUM kembali — baru masuk saat ApproveReturPenjualan dipanggil.
func (s *returnService) CreateReturPenjualan(userID uint, req *dto.CreateReturPenjualanRequest) (*dto.ReturPenjualanResponse, error) {
	// Ambil data penjualan asal
	sale, err := s.salesRepo.FindByID(req.IDPenjualan)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaksi penjualan tidak ditemukan")
		}
		return nil, err
	}

	now := time.Now()
	nomorRetur := fmt.Sprintf("RETP/%s/%d", now.Format("20060102150405"), userID)

	// Hitung subtotal retur dari item yang diretur
	var subtotal float64
	var returItems []models.ItemReturPenjualan

	// Buat map item penjualan untuk lookup cepat
	itemMap := make(map[uint]models.ItemPenjualan)
	for _, item := range sale.Items {
		itemMap[item.ID] = item
	}

	for _, itemReq := range req.Items {
		origItem, ok := itemMap[itemReq.IDItemPenjualan]
		if !ok {
			return nil, fmt.Errorf("item penjualan ID %d tidak ditemukan pada transaksi ini", itemReq.IDItemPenjualan)
		}
		if itemReq.Jumlah > origItem.Jumlah {
			return nil, fmt.Errorf("jumlah retur (%d) melebihi jumlah item asli (%d) untuk produk ID %d",
				itemReq.Jumlah, origItem.Jumlah, itemReq.IDProduk)
		}

		itemSubtotal := origItem.HargaSatuan * float64(itemReq.Jumlah)
		subtotal += itemSubtotal

		returItems = append(returItems, models.ItemReturPenjualan{
			IDItemPenjualan: itemReq.IDItemPenjualan,
			IDProduk:        itemReq.IDProduk,
			Jumlah:          itemReq.Jumlah,
			HargaSatuan:     origItem.HargaSatuan,
			Subtotal:        itemSubtotal,
			IDGudang:        sale.IDGudang,
			DibuatPada:      now,
			DiperbaruiPada:  now,
		})
	}

	retur := models.ReturPenjualan{
		NomorRetur:         nomorRetur,
		IDPenjualan:        req.IDPenjualan,
		NamaPelanggan:      sale.NamaPelanggan,
		KontakPelanggan:    sale.KontakPelanggan,
		Alasan:             req.Alasan,
		Subtotal:           subtotal,
		Total:              subtotal,
		MetodePengembalian: req.MetodePengembalian,
		JumlahPengembalian: subtotal,
		Status:             "pending",
		Keterangan:         req.Keterangan,
		DiprosesOleh:       userID,
		DiprosesPada:       now,
		DibuatPada:         now,
		DiperbaruiPada:     now,
		Items:              returItems,
	}

	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.repo.CreateReturPenjualan(tx, &retur); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal membuat retur penjualan: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.GetReturPenjualanByID(retur.ID)
}

// ApproveReturPenjualan menyetujui retur dan mencatat barang ke batch KARANTINA.
// ⚠️  Barang yang diretur TIDAK dikembalikan ke stok normal (tidak menambah stok_inventori).
// Barang diretur dicatat dalam batch terpisah (Aktif=false, TipeReferensi="retur_penjualan")
// sehingga bisa direkap sebagai stok retur tersendiri (bisa dijual, dihancurkan, dsb).
func (s *returnService) ApproveReturPenjualan(id, approvedByUserID uint) error {
	retur, err := s.repo.FindReturPenjualanByID(id)
	if err != nil {
		return errors.New("retur penjualan tidak ditemukan")
	}
	if retur.Status != "pending" {
		return fmt.Errorf("retur sudah dalam status '%s', tidak dapat diapprove", retur.Status)
	}

	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// Buat header barang masuk (hanya sebagai referensi dokumen, bukan penambah stok aktif)
	nomorBM := fmt.Sprintf("IN/RETP/%s/%d", now.Format("20060102150405"), approvedByUserID)
	headerMasuk := models.BarangMasuk{
		NomorTransaksi: nomorBM,
		DiterimaOleh:   approvedByUserID,
		DiterimaPada:   now,
		Status:         "approved",
		Keterangan:     fmt.Sprintf("[RETUR] %s — barang karantina, belum masuk stok aktif", retur.NomorRetur),
		DibuatPada:     now,
		DiperbaruiPada: now,
	}
	if err := s.stockRepo.CreateStockIn(tx, &headerMasuk); err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal membuat dokumen penerimaan retur: %w", err)
	}

	// Buat batch KARANTINA per item — TIDAK menambah stok_inventori
	for _, item := range retur.Items {
		batch := models.StokBatch{
			IDProduk:       item.IDProduk,
			IDGudang:       item.IDGudang,
			TanggalMasuk:   now,
			JumlahAwal:     item.Jumlah,
			JumlahSaatIni:  item.Jumlah,
			HargaModal:     item.HargaSatuan, // Nilai modal baris retur = harga jual asal
			IDReferensi:    &headerMasuk.ID,
			TipeReferensi:  "retur_penjualan",
			Aktif:          false, // ⚠️ false = tidak masuk ke stok FIFO yang bisa dijual
			Keterangan:     fmt.Sprintf("[KARANTINA-RETUR] %s | Nomor Retur: %s", item.Produk.SKU, retur.NomorRetur),
			DibuatPada:     now,
			DiperbaruiPada: now,
		}
		if err := s.batchRepo.Create(tx, &batch); err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal membuat batch karantina retur: %w", err)
		}

		// Log pergerakan untuk audit trail (jumlah positif = barang masuk ke sistem, tapi di stok karantina)
		batchID := batch.ID
		movement := models.PergerakanStok{
			IDProduk:       item.IDProduk,
			IDGudang:       item.IDGudang,
			IDBatch:        &batchID,
			TipePergerakan: "in",
			TipeReferensi:  "retur_penjualan",
			IDReferensi:    &retur.ID,
			Jumlah:         item.Jumlah,
			IDPengguna:     approvedByUserID,
			Keterangan:     fmt.Sprintf("[KARANTINA] Retur %s diapprove — batch karantina #%d (tidak masuk stok aktif)", retur.NomorRetur, batchID),
			DibuatPada:     now,
		}
		if err := s.stockRepo.CreateStockMovement(tx, &movement); err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal log pergerakan karantina: %w", err)
		}
		// ⚠️ UpdateStockBalance TIDAK dipanggil — stok_inventori tidak berubah
	}

	// Update status retur
	if err := s.repo.UpdateStatusReturPenjualan(tx, id, "completed", approvedByUserID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *returnService) GetReturPenjualanByID(id uint) (*dto.ReturPenjualanResponse, error) {
	retur, err := s.repo.FindReturPenjualanByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("retur penjualan tidak ditemukan")
		}
		return nil, err
	}
	return mapReturPenjualanToResponse(retur), nil
}

func (s *returnService) ListReturPenjualan(req *dto.ListReturPenjualanRequest) ([]dto.ReturPenjualanResponse, int64, error) {
	returs, total, err := s.repo.FindAllReturPenjualan(req)
	if err != nil {
		return nil, 0, err
	}
	var responses []dto.ReturPenjualanResponse
	for _, r := range returs {
		responses = append(responses, *mapReturPenjualanToResponse(&r))
	}
	return responses, total, nil
}

// ===========================
// RETUR PEMBELIAN
// ===========================

// CreateReturPembelian membuat dokumen retur ke supplier. Status awal: pending.
// Stok BELUM keluar — baru keluar saat ApproveReturPembelian dipanggil.
func (s *returnService) CreateReturPembelian(userID uint, req *dto.CreateReturPembelianRequest) (*dto.ReturPembelianResponse, error) {
	now := time.Now()
	nomorRetur := fmt.Sprintf("RETB/%s/%d", now.Format("20060102150405"), userID)

	var subtotal float64
	var returItems []models.ItemReturPembelian

	for _, itemReq := range req.Items {
		s := itemReq.HargaSatuan * float64(itemReq.Jumlah)
		subtotal += s
		returItems = append(returItems, models.ItemReturPembelian{
			IDProduk:       itemReq.IDProduk,
			Jumlah:         itemReq.Jumlah,
			HargaSatuan:    itemReq.HargaSatuan,
			Subtotal:       s,
			IDGudang:       req.IDGudang,
			DibuatPada:     now,
			DiperbaruiPada: now,
		})
	}

	retur := models.ReturPembelian{
		NomorRetur:         nomorRetur,
		IDPemasok:          req.IDPemasok,
		Alasan:             req.Alasan,
		Subtotal:           subtotal,
		Total:              subtotal,
		MetodePengembalian: req.MetodePengembalian,
		Status:             "pending",
		Keterangan:         req.Keterangan,
		DibuatOleh:         userID,
		DibuatPada:         now,
		DiperbaruiPada:     now,
		Items:              returItems,
	}

	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.repo.CreateReturPembelian(tx, &retur); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal membuat retur pembelian: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return s.GetReturPembelianByID(retur.ID)
}

// ApproveReturPembelian menyetujui retur ke supplier dan mengurangi stok via FIFO.
func (s *returnService) ApproveReturPembelian(id, approvedByUserID uint) error {
	retur, err := s.repo.FindReturPembelianByID(id)
	if err != nil {
		return errors.New("retur pembelian tidak ditemukan")
	}
	if retur.Status != "pending" {
		return fmt.Errorf("retur sudah dalam status '%s', tidak dapat diapprove", retur.Status)
	}

	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// Validasi stok cukup untuk semua item dulu
	for _, item := range retur.Items {
		batches, err := s.batchRepo.GetAvailableBatches(tx, item.IDProduk, item.IDGudang)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal get batch produk %d: %w", item.IDProduk, err)
		}
		totalAvail := 0
		for _, b := range batches {
			totalAvail += b.JumlahSaatIni
		}
		if totalAvail < item.Jumlah {
			tx.Rollback()
			return fmt.Errorf("stok tidak cukup untuk produk ID %d (tersedia: %d, dibutuhkan: %d)",
				item.IDProduk, totalAvail, item.Jumlah)
		}
	}

	// Buat header barang keluar
	nomorBK := fmt.Sprintf("OUT/RETB/%s/%d", now.Format("20060102150405"), approvedByUserID)
	headerKeluar := models.BarangKeluar{
		NomorTransaksi: nomorBK,
		Alasan:         "retur_pembelian",
		TipeReferensi:  "retur_pembelian",
		IDReferensi:    &retur.ID,
		DibuatOleh:     approvedByUserID,
		DibuatPada:     now,
		DiperbaruiPada: now,
	}
	if err := s.stockRepo.CreateStockOut(tx, &headerKeluar); err != nil {
		tx.Rollback()
		return err
	}

	// Kurangi stok via FIFO untuk setiap item
	for _, item := range retur.Items {
		batches, _ := s.batchRepo.GetAvailableBatches(tx, item.IDProduk, item.IDGudang)
		remaining := item.Jumlah

		for i := range batches {
			if remaining == 0 {
				break
			}
			batch := &batches[i]
			deduct := 0
			if batch.JumlahSaatIni >= remaining {
				deduct = remaining
				remaining = 0
			} else {
				deduct = batch.JumlahSaatIni
				remaining -= batch.JumlahSaatIni
			}

			batch.JumlahSaatIni -= deduct
			if err := s.batchRepo.Update(tx, batch); err != nil {
				tx.Rollback()
				return err
			}

			movement := models.PergerakanStok{
				IDProduk:       item.IDProduk,
				IDGudang:       item.IDGudang,
				IDBatch:        &batch.ID,
				TipePergerakan: "out",
				TipeReferensi:  "retur_pembelian",
				IDReferensi:    &retur.ID,
				Jumlah:         -deduct,
				IDPengguna:     approvedByUserID,
				Keterangan:     fmt.Sprintf("Retur ke Supplier %s (Batch #%d)", retur.NomorRetur, batch.ID),
				DibuatPada:     now,
			}
			if err := s.stockRepo.CreateStockMovement(tx, &movement); err != nil {
				tx.Rollback()
				return err
			}
		}

		if err := s.stockRepo.UpdateStockBalance(tx, item.IDProduk, item.IDGudang, -item.Jumlah); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := s.repo.UpdateStatusReturPembelian(tx, id, "completed", approvedByUserID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *returnService) GetReturPembelianByID(id uint) (*dto.ReturPembelianResponse, error) {
	retur, err := s.repo.FindReturPembelianByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("retur pembelian tidak ditemukan")
		}
		return nil, err
	}
	return mapReturPembelianToResponse(retur), nil
}

func (s *returnService) ListReturPembelian(req *dto.ListReturPembelianRequest) ([]dto.ReturPembelianResponse, int64, error) {
	returs, total, err := s.repo.FindAllReturPembelian(req)
	if err != nil {
		return nil, 0, err
	}
	var responses []dto.ReturPembelianResponse
	for _, r := range returs {
		responses = append(responses, *mapReturPembelianToResponse(&r))
	}
	return responses, total, nil
}

// ===========================
// MAPPING HELPERS
// ===========================

func mapReturPenjualanToResponse(r *models.ReturPenjualan) *dto.ReturPenjualanResponse {
	var items []dto.ReturPenjualanItemResponse
	for _, item := range r.Items {
		items = append(items, dto.ReturPenjualanItemResponse{
			ID:          item.ID,
			IDProduk:    item.IDProduk,
			SKUProduk:   item.Produk.SKU,
			NamaProduk:  item.Produk.Nama,
			Jumlah:      item.Jumlah,
			HargaSatuan: item.HargaSatuan,
			Subtotal:    item.Subtotal,
		})
	}
	nomorAsal := ""
	if r.Penjualan.NomorTransaksi != "" {
		nomorAsal = r.Penjualan.NomorTransaksi
	}
	return &dto.ReturPenjualanResponse{
		ID:                 r.ID,
		NomorRetur:         r.NomorRetur,
		IDPenjualan:        r.IDPenjualan,
		NomorTransaksiAsal: nomorAsal,
		NamaPelanggan:      r.NamaPelanggan,
		KontakPelanggan:    r.KontakPelanggan,
		Alasan:             r.Alasan,
		Subtotal:           r.Subtotal,
		Total:              r.Total,
		MetodePengembalian: r.MetodePengembalian,
		JumlahPengembalian: r.JumlahPengembalian,
		Status:             r.Status,
		Keterangan:         r.Keterangan,
		NamaPetugas:        r.DiprosesOlehPengguna.Nama,
		DibuatPada:         r.DibuatPada,
		Items:              items,
	}
}

func mapReturPembelianToResponse(r *models.ReturPembelian) *dto.ReturPembelianResponse {
	var items []dto.ReturPembelianItemResponse
	for _, item := range r.Items {
		items = append(items, dto.ReturPembelianItemResponse{
			ID:          item.ID,
			IDProduk:    item.IDProduk,
			SKUProduk:   item.Produk.SKU,
			NamaProduk:  item.Produk.Nama,
			Jumlah:      item.Jumlah,
			HargaSatuan: item.HargaSatuan,
			Subtotal:    item.Subtotal,
		})
	}
	return &dto.ReturPembelianResponse{
		ID:                 r.ID,
		NomorRetur:         r.NomorRetur,
		IDPemasok:          r.IDPemasok,
		NamaPemasok:        r.Pemasok.Nama,
		Alasan:             r.Alasan,
		Subtotal:           r.Subtotal,
		Total:              r.Total,
		MetodePengembalian: r.MetodePengembalian,
		Status:             r.Status,
		Keterangan:         r.Keterangan,
		NamaPembuat:        r.DibuatOlehPengguna.Nama,
		DibuatPada:         r.DibuatPada,
		Items:              items,
	}
}
