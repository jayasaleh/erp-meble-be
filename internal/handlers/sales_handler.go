package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SalesHandler struct {
	service services.SalesService
}

func NewSalesHandler(service services.SalesService) *SalesHandler {
	return &SalesHandler{service: service}
}

// CreateSale godoc
// @Summary      Buat transaksi penjualan baru (POS)
// @Description  Membuat transaksi penjualan dengan FIFO auto-deduct stok.
//
//	Mendukung dua mode request:
//	1. application/json  — body JSON langsung (tanpa upload file)
//	2. multipart/form-data — JSON di field "data", file opsional di "bukti_bayar"
//
// @Tags         sales
// @Accept       json,multipart/form-data
// @Produce      json
// @Param        body        body      dto.CreateSalesRequest  false  "JSON body (Content-Type: application/json)"
// @Param        data        formData  string                  false  "JSON body (Content-Type: multipart/form-data)"
// @Param        bukti_bayar formData  file                    false  "Foto bukti transfer (opsional, hanya multipart)"
// @Success      201  {object}  utils.Response{data=dto.SalesDetailResponse}
// @Router       /sales [post]
func (h *SalesHandler) CreateSale(c *gin.Context) {
	userID := utils.GetUserIDValidity(c)
	if userID == 0 {
		utils.Unauthorized(c, "Unauthorized")
		return
	}

	var req dto.CreateSalesRequest
	var buktiBayarPath *string

	contentType := c.GetHeader("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		// ── Mode 1: JSON body biasa ──────────────────────────────────────────
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequest(c, "Format JSON tidak valid: "+err.Error(), nil)
			return
		}
	} else {
		// ── Mode 2: multipart/form-data (untuk upload bukti bayar) ───────────
		dataStr := c.PostForm("data")
		if dataStr == "" {
			utils.BadRequest(c, "Field 'data' (JSON) wajib diisi untuk request multipart", nil)
			return
		}
		if err := json.Unmarshal([]byte(dataStr), &req); err != nil {
			utils.BadRequest(c, "Format JSON tidak valid", err.Error())
			return
		}

		// Proses upload bukti bayar (opsional)
		file, fileHeader, fileErr := c.Request.FormFile("bukti_bayar")
		if fileErr == nil && file != nil {
			defer file.Close()

			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".pdf": true}
			if !allowedExts[ext] {
				utils.BadRequest(c, "Format file tidak didukung. Gunakan JPG, PNG, WEBP, atau PDF", nil)
				return
			}

			uploadDir := "uploads/bukti_bayar"
			if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
				utils.InternalServerError(c, "Gagal membuat direktori upload", err.Error())
				return
			}

			fileName := fmt.Sprintf("bukti_%d%s", time.Now().UnixNano(), ext)
			filePath := filepath.Join(uploadDir, fileName)

			if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
				utils.InternalServerError(c, "Gagal menyimpan file", err.Error())
				return
			}
			buktiBayarPath = &filePath
		}
	}

	// ── Validasi umum ────────────────────────────────────────────────────────
	if req.IDGudang == 0 {
		utils.BadRequest(c, "id_gudang wajib diisi", nil)
		return
	}
	if req.MetodePembayaran != "cash" && req.MetodePembayaran != "transfer" {
		utils.BadRequest(c, "metode_pembayaran harus 'cash' atau 'transfer'", nil)
		return
	}
	if req.JumlahPembayaran <= 0 {
		utils.BadRequest(c, "jumlah_pembayaran harus lebih dari 0", nil)
		return
	}
	if len(req.Items) == 0 {
		utils.BadRequest(c, "Items tidak boleh kosong", nil)
		return
	}

	// ── Buat transaksi ───────────────────────────────────────────────────────
	result, err := h.service.CreateSale(userID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	// Update bukti bayar jika ada (mode multipart)
	if buktiBayarPath != nil {
		if err := h.service.UpdateBuktiBayar(result.ID, *buktiBayarPath); err == nil {
			result.BuktiBayar = buktiBayarPath
		}
	}

	utils.Created(c, "Transaksi penjualan berhasil dibuat", result)
}

// GetSale godoc
// @Summary      Detail transaksi penjualan
// @Description  Ambil detail lengkap satu transaksi termasuk breakdown batch FIFO
// @Tags         sales
// @Produce      json
// @Param        id   path      int  true  "ID Penjualan"
// @Success      200  {object}  utils.Response{data=dto.SalesDetailResponse}
// @Router       /sales/{id} [get]
func (h *SalesHandler) GetSale(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID tidak valid", nil)
		return
	}

	result, err := h.service.GetSaleByID(uint(id))
	if err != nil {
		if err.Error() == "transaksi penjualan tidak ditemukan" {
			utils.NotFound(c, "Transaksi tidak ditemukan")
			return
		}
		utils.InternalServerError(c, "Gagal mengambil data transaksi", err.Error())
		return
	}

	utils.OK(c, "Detail transaksi penjualan", result)
}

// ListSales godoc
// @Summary      List transaksi penjualan
// @Description  Daftar transaksi penjualan dengan filter tanggal, gudang, kasir, metode bayar
// @Tags         sales
// @Produce      json
// @Param        page               query  int     false  "Halaman"
// @Param        limit              query  int     false  "Jumlah per halaman"
// @Param        tanggal_dari       query  string  false  "Filter dari tanggal (YYYY-MM-DD)"
// @Param        tanggal_sampai     query  string  false  "Filter sampai tanggal (YYYY-MM-DD)"
// @Param        id_kasir           query  uint    false  "Filter ID kasir"
// @Param        id_gudang          query  uint    false  "Filter ID gudang"
// @Param        metode_pembayaran  query  string  false  "cash atau transfer"
// @Success      200  {object}  utils.Response{data=dto.ListSalesResponse}
// @Router       /sales [get]
func (h *SalesHandler) ListSales(c *gin.Context) {
	var req dto.ListSalesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Parameter tidak valid", err.Error())
		return
	}

	result, err := h.service.ListSales(&req)
	if err != nil {
		utils.InternalServerError(c, "Gagal mengambil daftar transaksi", err.Error())
		return
	}

	utils.OKWithMeta(c, "Daftar transaksi penjualan", result.Sales, utils.Meta{
		Page:      result.Page,
		Limit:     result.Limit,
		Total:     int(result.Total),
		TotalPage: result.TotalPages,
	})
}

// GetInvoice godoc
// @Summary      Data invoice untuk cetak
// @Description  Ambil data invoice yang dioptimalkan untuk keperluan cetak atau ekspor PDF
// @Tags         sales
// @Produce      json
// @Param        id   path      int  true  "ID Penjualan"
// @Success      200  {object}  utils.Response{data=dto.InvoiceResponse}
// @Router       /sales/{id}/invoice [get]
func (h *SalesHandler) GetInvoice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID tidak valid", nil)
		return
	}

	result, err := h.service.GetInvoice(uint(id))
	if err != nil {
		if err.Error() == "transaksi penjualan tidak ditemukan" {
			utils.NotFound(c, "Transaksi tidak ditemukan")
			return
		}
		utils.InternalServerError(c, "Gagal mengambil data invoice", err.Error())
		return
	}

	utils.OK(c, "Data invoice", result)
}

// UploadBuktiBayar godoc
// @Summary      Upload/ganti foto bukti transfer
// @Description  Upload atau ganti foto bukti transfer untuk transaksi yang sudah ada
// @Tags         sales
// @Accept       multipart/form-data
// @Produce      json
// @Param        id           path      int   true  "ID Penjualan"
// @Param        bukti_bayar  formData  file  true  "Foto bukti transfer"
// @Success      200  {object}  utils.Response
// @Router       /sales/{id}/bukti-bayar [post]
func (h *SalesHandler) UploadBuktiBayar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID tidak valid", nil)
		return
	}

	file, fileHeader, err := c.Request.FormFile("bukti_bayar")
	if err != nil {
		utils.BadRequest(c, "File 'bukti_bayar' wajib diupload", nil)
		return
	}
	defer file.Close()

	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".pdf": true}
	if !allowedExts[ext] {
		utils.BadRequest(c, "Format file tidak didukung. Gunakan JPG, PNG, WEBP, atau PDF", nil)
		return
	}

	// Buat direktori jika belum ada
	uploadDir := "uploads/bukti_bayar"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.InternalServerError(c, "Gagal membuat direktori upload", err.Error())
		return
	}

	fileName := fmt.Sprintf("bukti_%d_%d%s", id, time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		utils.InternalServerError(c, "Gagal menyimpan file", err.Error())
		return
	}

	if err := h.service.UpdateBuktiBayar(uint(id), filePath); err != nil {
		if err.Error() == "transaksi penjualan tidak ditemukan" {
			utils.NotFound(c, "Transaksi tidak ditemukan")
			return
		}
		utils.InternalServerError(c, "Gagal mengupdate bukti bayar", err.Error())
		return
	}

	utils.OK(c, "Bukti bayar berhasil diupload", gin.H{"path": filePath})
}
