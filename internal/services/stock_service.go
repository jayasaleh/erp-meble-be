package services

import (
	"errors"
	"fmt"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/models"
	"real-erp-mebel/be/internal/repositories"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type StockService interface {
	GetStocks(warehouseID, productID uint, limit, page int) ([]dto.InventoryResponse, int64, error)
	GetStockHistory(warehouseID, productID uint, refType string, limit, page int) ([]dto.StockMovementResponse, int64, error)
	GetStockBatches(productID, warehouseID uint, limit, page int) ([]dto.BatchResponse, int64, error)

	CreateStockIn(userID uint, req dto.CreateStockInRequest) error
	CreateStockOut(userID uint, req dto.CreateStockOutRequest) error
	CreateStockOpname(userID uint, req dto.CreateStockOpnameRequest) error
	CreateStockTransfer(userID uint, req dto.CreateStockTransferRequest) error
}

type stockService struct {
	repo      repositories.StockRepository
	batchRepo repositories.StockBatchRepository
}

func parseOpnameQtyFromNote(note string) (int, bool) {
	const key = "opname_batch_qty="
	idx := strings.Index(note, key)
	if idx == -1 {
		return 0, false
	}
	start := idx + len(key)
	end := start
	for end < len(note) && note[end] >= '0' && note[end] <= '9' {
		end++
	}
	if end == start {
		return 0, false
	}
	v, err := strconv.Atoi(note[start:end])
	if err != nil {
		return 0, false
	}
	return v, true
}

func NewStockService(repo repositories.StockRepository, batchRepo repositories.StockBatchRepository) StockService {
	return &stockService{
		repo:      repo,
		batchRepo: batchRepo,
	}
}

func (s *stockService) GetStocks(warehouseID, productID uint, limit, page int) ([]dto.InventoryResponse, int64, error) {
	// Offset calc
	offset := (page - 1) * limit
	var stocks []models.StokInventori
	var total int64
	var err error

	if productID != 0 && warehouseID != 0 {
		stock, err := s.repo.GetStockByProductAndWarehouse(productID, warehouseID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, err
		}
		if stock != nil {
			stocks = append(stocks, *stock)
			total = 1
		}
	} else {
		stocks, total, err = s.repo.GetStockByWarehouse(warehouseID, limit, offset)
		if err != nil {
			return nil, 0, err
		}
	}

	var responses []dto.InventoryResponse
	for _, item := range stocks {
		responses = append(responses, dto.InventoryResponse{
			ID:           item.ID,
			ProductID:    item.IDProduk,
			ProductSKU:   item.Produk.SKU,
			ProductName:  item.Produk.Nama,
			WarehouseID:  item.IDGudang,
			Warehouse:    item.Gudang.Nama,
			CurrentStock: item.Jumlah,
			LastUpdate:   item.DiperbaruiPada,
		})
	}
	return responses, total, nil
}

func (s *stockService) GetStockHistory(warehouseID, productID uint, refType string, limit, page int) ([]dto.StockMovementResponse, int64, error) {
	offset := (page - 1) * limit
	movements, total, err := s.repo.GetStockHistory(warehouseID, productID, refType, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.StockMovementResponse
	for _, m := range movements {
		costPrice := 0.0
		if m.Batch != nil {
			costPrice = m.Batch.HargaModal
		}

		systemStock := m.SaldoSetelah - m.Jumlah // If qty -2, bal_after 50, then sys_stock was 52

		responses = append(responses, dto.StockMovementResponse{
			ID:            m.ID,
			Date:          m.DibuatPada,
			Type:          m.TipePergerakan,
			ReferenceType: m.TipeReferensi,
			ReferenceID:   m.IDReferensi,
			ProductID:     m.IDProduk,
			ProductName:   m.Produk.Nama,
			Quantity:      m.Jumlah,
			SystemStock:   systemStock,
			BalanceAfter:  m.SaldoSetelah,
			CostPrice:     costPrice,
			WarehouseID:   m.IDGudang,
			WarehouseName: m.Gudang.Nama,
			OperatorName:  m.Pengguna.Nama,
			Notes:         m.Keterangan,
			BatchID:       m.IDBatch,
		})
	}
	return responses, total, nil
}

func (s *stockService) GetStockBatches(productID, warehouseID uint, limit, page int) ([]dto.BatchResponse, int64, error) {
	offset := (page - 1) * limit
	batches, total, err := s.batchRepo.GetAllBatches(productID, warehouseID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	batchIDs := make([]uint, 0, len(batches))
	for _, b := range batches {
		batchIDs = append(batchIDs, b.ID)
	}
	lastOpnameByBatch, err := s.batchRepo.GetLatestOpnameByBatchIDs(batchIDs)
	if err != nil {
		return nil, 0, err
	}

	creatorsByBatch, err := s.batchRepo.GetCreatorByBatchIDs(batchIDs)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.BatchResponse
	for _, b := range batches {
		var lastOpnameAt *time.Time
		var lastOpnameQty *int
		if m, ok := lastOpnameByBatch[b.ID]; ok {
			t := m.DibuatPada
			q := m.Jumlah
			if parsedQty, ok := parseOpnameQtyFromNote(m.Keterangan); ok {
				q = parsedQty
			}
			lastOpnameAt = &t
			lastOpnameQty = &q
		}

		responses = append(responses, dto.BatchResponse{
			ID:            b.ID,
			ProductID:     b.IDProduk,
			ProductName:   b.Produk.Nama,
			ProductSKU:    b.Produk.SKU,
			WarehouseID:   b.IDGudang,
			WarehouseName: b.Gudang.Nama,
			EntryDate:     b.TanggalMasuk,
			ExpiryDate:    b.TanggalKadaluarsa,
			InitialQty:    b.JumlahAwal,
			CurrentQty:    b.JumlahSaatIni,
			CostPrice:     b.HargaModal,
			ReferenceType: b.TipeReferensi,
			ReferenceID:   b.IDReferensi,
			Notes:         b.Keterangan,
			IsActive:      b.Aktif,
			LastOpnameAt:  lastOpnameAt,
			LastOpnameQty: lastOpnameQty,
			OperatorName:  creatorsByBatch[b.ID],
			CreatedAt:     b.DibuatPada,
		})
	}
	return responses, total, nil
}

func (s *stockService) CreateStockIn(userID uint, req dto.CreateStockInRequest) (err error) {
	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	// 1. Create Header BarangMasuk
	now := time.Now()
	if !req.Date.IsZero() {
		now = req.Date
	}

	header := models.BarangMasuk{
		NomorTransaksi: fmt.Sprintf("IN/MANUAL/%d/%d", now.Unix(), userID), // Simple logic
		DiterimaOleh:   userID,
		DiterimaPada:   now,
		Status:         "approved", // Direct approved for manual stock in
		Keterangan:     req.Notes,
		DibuatPada:     now,
		DiperbaruiPada: now,
	}

	if err := s.repo.CreateStockIn(tx, &header); err != nil {
		tx.Rollback()
		return err
	}

	// 2. Process Items - Buat Batch untuk setiap item (FIFO)
	for _, item := range req.Items {
		// Ambil harga modal produk untuk HPP batch ini
		var p models.Produk
		hargaModal := 0.0
		if err := tx.Select("harga_modal").First(&p, item.ProductID).Error; err == nil {
			hargaModal = p.HargaModal
		}

		// Buat batch baru untuk barang masuk ini
		batch := models.StokBatch{
			IDProduk:      item.ProductID,
			IDGudang:      req.WarehouseID,
			TanggalMasuk:  now,
			JumlahAwal:    item.Quantity,
			JumlahSaatIni: item.Quantity,
			HargaModal:    hargaModal,
			IDReferensi:   &header.ID,
			TipeReferensi: "stock_in",
			Aktif:         true,
			Keterangan:    req.Notes,
			DibuatPada:    now,
			DiperbaruiPada: now,
		}
		if err := s.batchRepo.Create(tx, &batch); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create batch: %w", err)
		}

		// Update Stock Balance (Totalan)
		if err := s.repo.UpdateStockBalance(tx, item.ProductID, req.WarehouseID, item.Quantity); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update stock balance: %w", err)
		}

		// Create Movement Log dengan link ke batch
		movement := models.PergerakanStok{
			IDProduk:       item.ProductID,
			IDGudang:       req.WarehouseID,
			IDBatch:        &batch.ID, // Link ke batch yang baru dibuat
			TipePergerakan: "in",
			TipeReferensi:  "manual_in",
			IDReferensi:    &header.ID,
			Jumlah:         item.Quantity,
			IDPengguna:     userID,
			Keterangan:     fmt.Sprintf("%s (New Batch #%d)", req.Notes, batch.ID),
			DibuatPada:     now,
		}
		if err := s.repo.CreateStockMovement(tx, &movement); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create movement: %w", err)
		}
	}

	return tx.Commit().Error
}

func (s *stockService) CreateStockOut(userID uint, req dto.CreateStockOutRequest) (err error) {
	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	now := time.Now()
	if !req.Date.IsZero() {
		now = req.Date
	}

	header := models.BarangKeluar{
		NomorTransaksi: fmt.Sprintf("OUT/MANUAL/%d/%d", now.Unix(), userID),
		Alasan:         req.Reason,
		DibuatOleh:     userID,
		DibuatPada:     now,
		DiperbaruiPada: now,
	}

	if err := s.repo.CreateStockOut(tx, &header); err != nil {
		tx.Rollback()
		return err
	}


	for _, item := range req.Items {
		// === FIFO LOGIC: Ambil batch terlama sampai qty terpenuhi ===
		batches, err := s.batchRepo.GetAvailableBatches(tx, item.ProductID, req.WarehouseID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get batches: %w", err)
		}

		remainingQty := item.Quantity
		
		// Check total available
		totalAvailable := 0
		for _, b := range batches {
			totalAvailable += b.JumlahSaatIni
		}
		if totalAvailable < item.Quantity {
			tx.Rollback()
			return fmt.Errorf("insufficient stock for product %d (need %d, available %d)", 
				item.ProductID, item.Quantity, totalAvailable)
		}

		// Deduct dari batch terlama (FIFO)
		for i := range batches {
			if remainingQty == 0 {
				break
			}

			batch := &batches[i]
			deductQty := 0
			
			if batch.JumlahSaatIni >= remainingQty {
				// Batch ini cukup
				deductQty = remainingQty
				remainingQty = 0
			} else {
				// Batch habis, ambil semua
				deductQty = batch.JumlahSaatIni
				remainingQty -= batch.JumlahSaatIni
			}

			// Update batch
			batch.JumlahSaatIni -= deductQty
			if err := s.batchRepo.Update(tx, batch); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update batch: %w", err)
			}

			// Log pergerakan dengan link ke batch
			movement := models.PergerakanStok{
				IDProduk:       item.ProductID,
				IDGudang:       req.WarehouseID,
				IDBatch:        &batch.ID, // Link ke batch FIFO
				TipePergerakan: "out",
				TipeReferensi:  "manual_out",
				IDReferensi:    &header.ID,
				Jumlah:         -deductQty,
				IDPengguna:     userID,
				Keterangan:     fmt.Sprintf("%s (Batch #%d, HPP: %.2f)", req.Reason, batch.ID, batch.HargaModal),
				DibuatPada:     now,
			}
			if err := s.repo.CreateStockMovement(tx, &movement); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create movement: %w", err)
			}
		}

		// Update StokInventori (totalan)
		if err := s.repo.UpdateStockBalance(tx, item.ProductID, req.WarehouseID, -item.Quantity); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update stock balance: %w", err)
		}
	}

	return tx.Commit().Error
}

func (s *stockService) CreateStockOpname(userID uint, req dto.CreateStockOpnameRequest) (err error) {
	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	now := time.Now()
	if !req.Date.IsZero() {
		now = req.Date
	}

	for _, item := range req.Items {
		var diff int
		var systemQty int

		if item.BatchID != 0 {
			var batch models.StokBatch
			// Load batch with FOR UPDATE
			if err := tx.Model(&models.StokBatch{}).Where("id = ?", item.BatchID).First(&batch).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("batch not found: %w", err)
			}
			systemQty = batch.JumlahSaatIni
			diff = item.ActualStock - systemQty

			if diff == 0 {
				continue // No change
			}

			// Update that specific batch
			batch.JumlahSaatIni = item.ActualStock
			batch.DiperbaruiPada = now
			if batch.JumlahSaatIni <= 0 {
				batch.Aktif = false
			} else {
				batch.Aktif = true
			}
			if err := tx.Save(&batch).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update batch: %w", err)
			}

			// Update Total Inventory Balance
			if err := s.repo.UpdateStockBalance(tx, item.ProductID, req.WarehouseID, diff); err != nil {
				tx.Rollback()
				return err
			}

			// Record Movement
			movement := models.PergerakanStok{
				IDProduk:       item.ProductID,
				IDGudang:       req.WarehouseID,
				IDBatch:        &batch.ID,
				TipePergerakan: "adjustment",
				TipeReferensi:  "opname",
				Jumlah:         diff,
				IDPengguna:     userID,
				Keterangan:     fmt.Sprintf("Opname Batch #%d: System %d -> Actual %d. %s | opname_batch_qty=%d", batch.ID, systemQty, item.ActualStock, req.Notes, batch.JumlahSaatIni),
				DibuatPada:     now,
			}
			if err := s.repo.CreateStockMovement(tx, &movement); err != nil {
				tx.Rollback()
				return err
			}

		} else {
			// BatchID = 0 means new surplus physical stock found
			diff = item.ActualStock
			if diff <= 0 {
				continue // Cannot have negative diff without batch
			}

			var p models.Produk
			hargaModal := 0.0
			if err := tx.Select("harga_modal").First(&p, item.ProductID).Error; err == nil {
				hargaModal = p.HargaModal
			}

			batch := models.StokBatch{
				IDProduk:       item.ProductID,
				IDGudang:       req.WarehouseID,
				TanggalMasuk:   now,
				JumlahAwal:     diff,
				JumlahSaatIni:  diff,
				HargaModal:     hargaModal,
				TipeReferensi:  "opname_adjustment_in",
				Aktif:          true,
				Keterangan:     req.Notes,
				DibuatPada:     now,
				DiperbaruiPada: now,
			}
			if err := s.batchRepo.Create(tx, &batch); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create opname batch: %w", err)
			}

			// Update Total Inventory Balance
			if err := s.repo.UpdateStockBalance(tx, item.ProductID, req.WarehouseID, diff); err != nil {
				tx.Rollback()
				return err
			}

			movement := models.PergerakanStok{
				IDProduk:       item.ProductID,
				IDGudang:       req.WarehouseID,
				IDBatch:        &batch.ID,
				TipePergerakan: "adjustment",
				TipeReferensi:  "opname",
				Jumlah:         diff,
				IDPengguna:     userID,
				Keterangan:     fmt.Sprintf("Opname (New Surplus Batch): Actual %d. %s | opname_batch_qty=%d", item.ActualStock, req.Notes, batch.JumlahSaatIni),
				DibuatPada:     now,
			}
			if err := s.repo.CreateStockMovement(tx, &movement); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (s *stockService) CreateStockTransfer(userID uint, req dto.CreateStockTransferRequest) (err error) {
	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()

	now := time.Now()
	if req.Date != nil && !req.Date.IsZero() {
		now = *req.Date
	}

	transferID := fmt.Sprintf("TRF/%d/%d/%d", req.SourceWarehouseID, req.TargetWarehouseID, now.Unix())

	for _, item := range req.Items {
		// 1. Check Source Stock
		sourceStock, err := s.repo.GetStockByProductAndWarehouse(item.ProductID, req.SourceWarehouseID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}
		if sourceStock == nil || sourceStock.Jumlah < item.Quantity {
			tx.Rollback()
			return fmt.Errorf("insufficient stock for product %d in source warehouse", item.ProductID)
		}

		// 2. Decrement Source Stock
		if err := s.repo.UpdateStockBalance(tx, item.ProductID, req.SourceWarehouseID, -item.Quantity); err != nil {
			tx.Rollback()
			return err
		}

		// 3. Log Source Movement (Transfer Out)
		outMovement := models.PergerakanStok{
			IDProduk:       item.ProductID,
			IDGudang:       req.SourceWarehouseID,
			TipePergerakan: "transfer_out",
			TipeReferensi:  "transfer",
			Keterangan:     fmt.Sprintf("Transfer to Warehouse %d. %s", req.TargetWarehouseID, req.Notes), // Store Target ID in notes or Reference?
			IDPengguna:     userID,
			Jumlah:         -item.Quantity,
			DibuatPada:     now,
			// Kita bisa simpan Transfer ID di Keterangan atau field lain jika ada
			// Untuk sekarang string formatting di Keterangan cukup
		}
		// Hack to store Transfer Ref. actually Request doesn't map to a table.
		// We use Keterangan to store Transfer ID for tracing.
		outMovement.Keterangan += fmt.Sprintf(" (Ref: %s)", transferID)

		if err := s.repo.CreateStockMovement(tx, &outMovement); err != nil {
			tx.Rollback()
			return err
		}

		// 4. Increment Target Stock
		if err := s.repo.UpdateStockBalance(tx, item.ProductID, req.TargetWarehouseID, item.Quantity); err != nil {
			tx.Rollback()
			return err
		}

		// 5. Log Target Movement (Transfer In)
		inMovement := models.PergerakanStok{
			IDProduk:       item.ProductID,
			IDGudang:       req.TargetWarehouseID,
			TipePergerakan: "transfer_in",
			TipeReferensi:  "transfer",
			Keterangan:     fmt.Sprintf("Transfer from Warehouse %d. %s (Ref: %s)", req.SourceWarehouseID, req.Notes, transferID),
			IDPengguna:     userID,
			Jumlah:         item.Quantity,
			DibuatPada:     now,
		}

		if err := s.repo.CreateStockMovement(tx, &inMovement); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
