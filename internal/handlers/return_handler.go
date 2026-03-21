package handlers

import (
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReturnHandler struct {
	service services.ReturnService
}

func NewReturnHandler(service services.ReturnService) *ReturnHandler {
	return &ReturnHandler{service: service}
}

// ===========================
// RETUR PENJUALAN
// ===========================

func (h *ReturnHandler) CreateReturPenjualan(c *gin.Context) {
	userID := utils.GetUserIDValidity(c)
	if userID == 0 {
		utils.Unauthorized(c, "Unauthorized")
		return
	}
	var req dto.CreateReturPenjualanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Request tidak valid", err.Error())
		return
	}
	result, err := h.service.CreateReturPenjualan(userID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}
	utils.Created(c, "Retur penjualan berhasil dibuat (status: pending)", result)
}

func (h *ReturnHandler) GetReturPenjualan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID tidak valid", nil)
		return
	}
	result, err := h.service.GetReturPenjualanByID(uint(id))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.OK(c, "Detail retur penjualan", result)
}

func (h *ReturnHandler) ListReturPenjualan(c *gin.Context) {
	var req dto.ListReturPenjualanRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Parameter tidak valid", err.Error())
		return
	}
	results, total, err := h.service.ListReturPenjualan(&req)
	if err != nil {
		utils.InternalServerError(c, "Gagal mengambil data", err.Error())
		return
	}
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	utils.OKWithMeta(c, "Daftar retur penjualan", results, utils.Meta{
		Page: page, Limit: limit, Total: int(total), TotalPage: totalPages,
	})
}

func (h *ReturnHandler) ApproveReturPenjualan(c *gin.Context) {
	userID := utils.GetUserIDValidity(c)
	if userID == 0 {
		utils.Unauthorized(c, "Unauthorized")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.service.ApproveReturPenjualan(uint(id), userID); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}
	utils.OK(c, "Retur penjualan diapprove — stok sudah dikembalikan ke gudang", nil)
}

// ===========================
// RETUR PEMBELIAN
// ===========================

func (h *ReturnHandler) CreateReturPembelian(c *gin.Context) {
	userID := utils.GetUserIDValidity(c)
	if userID == 0 {
		utils.Unauthorized(c, "Unauthorized")
		return
	}
	var req dto.CreateReturPembelianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Request tidak valid", err.Error())
		return
	}
	result, err := h.service.CreateReturPembelian(userID, &req)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}
	utils.Created(c, "Retur pembelian berhasil dibuat (status: pending)", result)
}

func (h *ReturnHandler) GetReturPembelian(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID tidak valid", nil)
		return
	}
	result, err := h.service.GetReturPembelianByID(uint(id))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}
	utils.OK(c, "Detail retur pembelian", result)
}

func (h *ReturnHandler) ListReturPembelian(c *gin.Context) {
	var req dto.ListReturPembelianRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Parameter tidak valid", err.Error())
		return
	}
	results, total, err := h.service.ListReturPembelian(&req)
	if err != nil {
		utils.InternalServerError(c, "Gagal mengambil data", err.Error())
		return
	}
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	utils.OKWithMeta(c, "Daftar retur pembelian", results, utils.Meta{
		Page: page, Limit: limit, Total: int(total), TotalPage: totalPages,
	})
}

func (h *ReturnHandler) ApproveReturPembelian(c *gin.Context) {
	userID := utils.GetUserIDValidity(c)
	if userID == 0 {
		utils.Unauthorized(c, "Unauthorized")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID tidak valid", nil)
		return
	}
	if err := h.service.ApproveReturPembelian(uint(id), userID); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}
	utils.OK(c, "Retur pembelian diapprove — stok sudah dikurangi dari gudang", nil)
}
