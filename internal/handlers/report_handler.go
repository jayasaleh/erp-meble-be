package handlers

import (
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetSalesReportByPeriod godoc
// @Summary      Laporan penjualan per periode
// @Description  Total revenue, COGS, laba, margin, breakdown per hari dan per metode bayar
// @Tags         reports
// @Produce      json
// @Param        tanggal_dari    query  string  true  "Dari tanggal (YYYY-MM-DD)"
// @Param        tanggal_sampai  query  string  true  "Sampai tanggal (YYYY-MM-DD)"
// @Param        id_gudang       query  uint    false "Filter gudang"
// @Router       /reports/sales [get]
func (h *ReportHandler) GetSalesReportByPeriod(c *gin.Context) {
	var req dto.SalesReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Parameter tidak valid. Pastikan tanggal_dari dan tanggal_sampai diisi (format: YYYY-MM-DD)", err.Error())
		return
	}
	result, err := h.service.GetSalesReportByPeriod(&req)
	if err != nil {
		utils.InternalServerError(c, "Gagal menghasilkan laporan", err.Error())
		return
	}
	utils.OK(c, "Laporan penjualan per periode", result)
}

// GetSalesReportByProduct godoc
// @Summary      Laporan penjualan per produk
// @Description  Jumlah terjual, revenue, COGS, laba, margin per produk — diurutkan paling laku
// @Tags         reports
// @Produce      json
// @Param        tanggal_dari    query  string  true  "Dari tanggal (YYYY-MM-DD)"
// @Param        tanggal_sampai  query  string  true  "Sampai tanggal (YYYY-MM-DD)"
// @Param        id_gudang       query  uint    false "Filter gudang"
// @Router       /reports/sales/by-product [get]
func (h *ReportHandler) GetSalesReportByProduct(c *gin.Context) {
	var req dto.SalesReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Parameter tidak valid. Pastikan tanggal_dari dan tanggal_sampai diisi (format: YYYY-MM-DD)", err.Error())
		return
	}
	result, err := h.service.GetSalesReportByProduct(&req)
	if err != nil {
		utils.InternalServerError(c, "Gagal menghasilkan laporan", err.Error())
		return
	}
	utils.OK(c, "Laporan penjualan per produk", result)
}

// GetSalesReportByCustomer godoc
// @Summary      Laporan penjualan per pelanggan
// @Description  Total transaksi, total belanja, dan total laba per nama pelanggan
// @Tags         reports
// @Produce      json
// @Param        tanggal_dari    query  string  true  "Dari tanggal (YYYY-MM-DD)"
// @Param        tanggal_sampai  query  string  true  "Sampai tanggal (YYYY-MM-DD)"
// @Router       /reports/sales/by-customer [get]
func (h *ReportHandler) GetSalesReportByCustomer(c *gin.Context) {
	var req dto.SalesReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Parameter tidak valid. Pastikan tanggal_dari dan tanggal_sampai diisi (format: YYYY-MM-DD)", err.Error())
		return
	}
	result, err := h.service.GetSalesReportByCustomer(&req)
	if err != nil {
		utils.InternalServerError(c, "Gagal menghasilkan laporan", err.Error())
		return
	}
	utils.OK(c, "Laporan penjualan per pelanggan", result)
}

// GetReturnReport godoc
// @Summary      Laporan rekapitulasi retur
// @Description  Mendapatkan rekap kerugian dan stok karantina akibat retur penjualan (dan retur pembelian jika ada)
// @Tags         reports
// @Produce      json
// @Param        tanggal_dari    query  string  false  "Dari tanggal (YYYY-MM-DD)"
// @Param        tanggal_sampai  query  string  false  "Sampai tanggal (YYYY-MM-DD)"
// @Router       /reports/returns [get]
func (h *ReportHandler) GetReturnReport(c *gin.Context) {
	var req dto.ReturnReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Parameter tidak valid. Pastikan format tanggal YYYY-MM-DD jika diisi", err.Error())
		return
	}
	result, err := h.service.GetReturnReport(&req)
	if err != nil {
		utils.InternalServerError(c, "Gagal menghasilkan laporan", err.Error())
		return
	}
	utils.OK(c, "Laporan rekapitulasi retur", result)
}
