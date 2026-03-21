package handlers

import (
	"real-erp-mebel/be/internal/dto"

	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GudangHandler interface {
	CreateGudang(c *gin.Context)
	GetGudang(c *gin.Context)
	ListGudangs(c *gin.Context)
	UpdateGudang(c *gin.Context)
	DeleteGudang(c *gin.Context)
}

type gudangHandler struct {
	service services.GudangService
}

func NewGudangHandler(service services.GudangService) GudangHandler {
	return &gudangHandler{service}
}

func (h *gudangHandler) CreateGudang(c *gin.Context) {
	var req dto.CreateGudangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	gudang, err := h.service.CreateGudang(&req)
	if err != nil {
		utils.InternalServerError(c, "Failed to create gudang", err.Error())
		return
	}

	utils.Created(c, "Gudang berhasil dibuat", gudang)
}

func (h *gudangHandler) GetGudang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID gudang tidak valid", err.Error())
		return
	}

	gudang, err := h.service.GetGudangByID(uint(id))
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.OK(c, "Detail gudang", gudang)
}

func (h *gudangHandler) ListGudangs(c *gin.Context) {
	var req dto.ListGudangRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	result, err := h.service.ListGudangs(&req)
	if err != nil {
		utils.InternalServerError(c, "Failed to list gudangs", err.Error())
		return
	}

	utils.OK(c, "Daftar gudang", result)
}

func (h *gudangHandler) UpdateGudang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID gudang tidak valid", err.Error())
		return
	}

	var req dto.UpdateGudangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	gudang, err := h.service.UpdateGudang(uint(id), &req)
	if err != nil {
		utils.InternalServerError(c, "Failed to update gudang", err.Error())
		return
	}

	utils.OK(c, "Gudang berhasil diupdate", gudang)
}

func (h *gudangHandler) DeleteGudang(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID gudang tidak valid", err.Error())
		return
	}

	if err := h.service.DeleteGudang(uint(id)); err != nil {
		utils.InternalServerError(c, "Failed to delete gudang", err.Error())
		return
	}

	utils.OK(c, "Gudang berhasil dihapus", nil)
}
