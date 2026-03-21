package handlers

import (
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PemasokHandler struct {
	service services.PemasokService
}

func NewPemasokHandler(service services.PemasokService) *PemasokHandler {
	return &PemasokHandler{service: service}
}

// CreatePemasok godoc
// @Summary      Create new supplier
// @Description  Create a new supplier
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreatePemasokRequest  true  "Request Body"
// @Success      201  {object}  utils.Response{data=dto.PemasokResponse}
// @Router       /suppliers [post]
func (h *PemasokHandler) CreatePemasok(c *gin.Context) {
	var req dto.CreatePemasokRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request", err)
		return
	}

	pemasok, err := h.service.CreatePemasok(&req)
	if err != nil {
		utils.InternalServerError(c, "Failed to create supplier", err.Error())
		return
	}

	utils.Created(c, "Supplier created successfully", pemasok)
}

// GetPemasok godoc
// @Summary      Get supplier by ID
// @Description  Get supplier details by ID
// @Tags         suppliers
// @Produce      json
// @Param        id   path      int  true  "Supplier ID"
// @Success      200  {object}  utils.Response{data=dto.PemasokResponse}
// @Router       /suppliers/{id} [get]
func (h *PemasokHandler) GetPemasok(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ID", nil)
		return
	}

	pemasok, err := h.service.GetPemasokByID(uint(id))
	if err != nil {
		if err.Error() == "pemasok tidak ditemukan" {
			utils.NotFound(c, "Supplier not found")
			return
		}
		utils.InternalServerError(c, "Failed to fetch supplier", err.Error())
		return
	}

	utils.OK(c, "Supplier fetched successfully", pemasok)
}

// ListPemasok godoc
// @Summary      List suppliers
// @Description  Get list of suppliers with filter and pagination
// @Tags         suppliers
// @Produce      json
// @Param        page    query     int     false  "Page number"
// @Param        limit   query     int     false  "Items per page"
// @Param        search  query     string  false  "Search by name/contact/email"
// @Param        aktif   query     bool    false  "Filter by active status"
// @Success      200     {object}  utils.Response{data=dto.ListPemasokResponse}
// @Router       /suppliers [get]
func (h *PemasokHandler) ListPemasok(c *gin.Context) {
	var req dto.ListPemasokRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Invalid query parameters", err)
		return
	}

	response, err := h.service.ListPemasok(&req)
	if err != nil {
		utils.InternalServerError(c, "Failed to list suppliers", err.Error())
		return
	}

	utils.OK(c, "Suppliers listed successfully", response)
}

// UpdatePemasok godoc
// @Summary      Update supplier
// @Description  Update supplier details
// @Tags         suppliers
// @Accept       json
// @Produce      json
// @Param        id   path      int                       true  "Supplier ID"
// @Param        req  body      dto.UpdatePemasokRequest  true  "Request Body"
// @Success      200  {object}  utils.Response{data=dto.PemasokResponse}
// @Router       /suppliers/{id} [put]
func (h *PemasokHandler) UpdatePemasok(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ID", nil)
		return
	}

	var req dto.UpdatePemasokRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request", err)
		return
	}

	pemasok, err := h.service.UpdatePemasok(uint(id), &req)
	if err != nil {
		if err.Error() == "pemasok tidak ditemukan" {
			utils.NotFound(c, "Supplier not found")
			return
		}
		utils.InternalServerError(c, "Failed to update supplier", err.Error())
		return
	}

	utils.OK(c, "Supplier updated successfully", pemasok)
}

// DeletePemasok godoc
// @Summary      Delete supplier
// @Description  Soft delete supplier
// @Tags         suppliers
// @Produce      json
// @Param        id   path      int  true  "Supplier ID"
// @Success      200  {object}  utils.Response
// @Router       /suppliers/{id} [delete]
func (h *PemasokHandler) DeletePemasok(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid ID", nil)
		return
	}

	if err := h.service.DeletePemasok(uint(id)); err != nil {
		if err.Error() == "pemasok tidak ditemukan" {
			utils.NotFound(c, "Supplier not found")
			return
		}
		utils.InternalServerError(c, "Failed to delete supplier", err.Error())
		return
	}

	utils.OK(c, "Supplier deleted successfully", nil)
}
