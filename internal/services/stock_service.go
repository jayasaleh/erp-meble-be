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

type StockService interface {
	GetStocks(warehouseID, productID uint) ([]dto.InventoryResponse, error)
	GetStockHistory(warehouseID, productID uint, limit, page int) ([]dto.StockMovementResponse, int64, error)

	CreateStockIn(userID uint, req dto.CreateStockInRequest) error
	CreateStockOut(userID uint, req dto.CreateStockOutRequest) error
	CreateStockOpname(userID uint, req dto.CreateStockOpnameRequest) error
	CreateStockTransfer(userID uint, req dto.CreateStockTransferRequest) error
}

type stockService struct {
	repo      repositories.StockRepository
	batchRepo repositories.StockBatchRepository
}

func NewStockService(repo repositories.StockRepository, batchRepo repositories.StockBatchRepository) StockService {
	return &stockService{
		repo:      repo,
		batchRepo: batchRepo,
	}
}

func (s *stockService) GetStocks(warehouseID, productID uint) ([]dto.InventoryResponse, error) {
	// Logic to fetch and map to DTO
	var stocks []models.StokInventori
	var err error

	if productID != 0 && warehouseID != 0 {
		stock, err := s.repo.GetStockByProductAndWarehouse(productID, warehouseID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if stock != nil {
			stocks = append(stocks, *stock)
		}
	} else {
		stocks, err = s.repo.GetStockByWarehouse(warehouseID)
		if err != nil {
			return nil, err
		}
	}

	var responses []dto.InventoryResponse
	for _, item := range stocks {
		responses = append(responses, dto.InventoryResponse{
			ID:        item.ID,
			ProductID: item.IDProduk,

			// Note: Preload should populate these
			ProductSKU:   item.Produk.SKU,
			ProductName:  item.Produk.Nama,
			WarehouseID:  item.IDGudang,
			Warehouse:    item.Gudang.Nama,
			CurrentStock: item.Jumlah,
			LastUpdate:   item.DiperbaruiPada,
		})
	}
	return responses, nil
}

func (s *stockService) GetStockHistory(warehouseID, productID uint, limit, page int) ([]dto.StockMovementResponse, int64, error) {
	offset := (page - 1) * limit
	movements, total, err := s.repo.GetStockHistory(warehouseID, productID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.StockMovementResponse
	for _, m := range movements {

		var opName string
		// Basic check if user loaded
		opName = m.Pengguna.Nama // Assuming User model has Nama

		responses = append(responses, dto.StockMovementResponse{
			ID:            m.ID,
			Date:          m.DibuatPada,
			Type:          m.TipePergerakan,
			ReferenceType: m.TipeReferensi,
			ReferenceID:   m.IDReferensi,
			Quantity:      m.Jumlah,
			BalanceAfter:  m.SaldoSetelah,
			WarehouseName: m.Gudang.Nama,
			OperatorName:  opName,
			Notes:         m.Keterangan,
		})
	}
	return responses, total, nil
}

func (s *stockService) CreateStockIn(userID uint, req dto.CreateStockInRequest) error {
	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
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

func (s *stockService) CreateStockOut(userID uint, req dto.CreateStockOutRequest) error {
	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
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

func (s *stockService) CreateStockOpname(userID uint, req dto.CreateStockOpnameRequest) error {
	// Logic:
	// User sends Actual Stock.
	// System finds System Stock.
	// Diff = Actual - System.
	// If Diff > 0 -> Create Adjustment IN.
	// If Diff < 0 -> Create Adjustment OUT.

	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()
	if !req.Date.IsZero() {
		now = req.Date
	}

	// Note: We should probably create a "StockOpname" header table if we want to track the event itself.
	// Existing schema doesn't seem to have specific Opname header, or maybe I missed it.
	// checking `stock.go`... no explicit Opname header model.
	// We will use `BarangMasuk` / `BarangKeluar` (or simpler, just Movement log with type 'adjustment').
	// The implementation plan says "Log history".
	// Let's stick to logging movement with TipePergerakan "adjustment".

	for _, item := range req.Items {
		currentStock, err := s.repo.GetStockByProductAndWarehouse(item.ProductID, req.WarehouseID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return err
		}

		systemQty := 0
		if currentStock != nil {
			systemQty = currentStock.Jumlah
		}

		diff := item.ActualStock - systemQty

		if diff == 0 {
			continue // No change
		}

		// Update Balance
		if err := s.repo.UpdateStockBalance(tx, item.ProductID, req.WarehouseID, diff); err != nil {
			tx.Rollback()
			return err
		}

		movement := models.PergerakanStok{
			IDProduk:       item.ProductID,
			IDGudang:       req.WarehouseID,
			TipePergerakan: "adjustment",
			TipeReferensi:  "opname",
			// IDReferensi:    ?? No header for opname unless generic
			Jumlah:     diff,
			IDPengguna: userID,
			Keterangan: fmt.Sprintf("Opname: System %d -> Actual %d. %s", systemQty, item.ActualStock, req.Notes),
			DibuatPada: now,
		}
		if err := s.repo.CreateStockMovement(tx, &movement); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *stockService) CreateStockTransfer(userID uint, req dto.CreateStockTransferRequest) error {
	tx := s.repo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
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
