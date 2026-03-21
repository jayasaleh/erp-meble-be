package handlers

import (
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	service services.StockService
}

func NewStockHandler(service services.StockService) *StockHandler {
	return &StockHandler{service: service}
}

// GetStocks godoc
// @Summary      Get current stock
// @Description  Get currently available stock by warehouse and/or product
// @Tags         stocks
// @Produce      json
// @Param        warehouse_id   query   int  false  "Warehouse ID"
// @Param        product_id     query   int  false  "Product ID"
// @Success      200  {object}  utils.Response{data=[]dto.InventoryResponse}
// @Router       /stocks [get]
func (h *StockHandler) GetStocks(c *gin.Context) {
	warehouseID, _ := strconv.Atoi(c.Query("warehouse_id"))
	productID, _ := strconv.Atoi(c.Query("product_id"))

	stocks, err := h.service.GetStocks(uint(warehouseID), uint(productID))
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch stocks", err.Error())
		return
	}

	utils.OK(c, "Stocks fetched successfully", stocks)
}

// GetStockHistory godoc
// @Summary      Get stock movement history
// @Description  Get history of stock changes
// @Tags         stocks
// @Produce      json
// @Param        warehouse_id   query   int  false  "Warehouse ID"
// @Param        product_id     query   int  false  "Product ID"
// @Param        page           query   int  false  "Page number"
// @Param        limit          query   int  false  "Items per page"
// @Success      200  {object}  utils.Response{data=[]dto.StockMovementResponse}
// @Router       /stocks/history [get]
func (h *StockHandler) GetStockHistory(c *gin.Context) {
	warehouseID, _ := strconv.Atoi(c.Query("warehouse_id"))
	productID, _ := strconv.Atoi(c.Query("product_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	history, total, err := h.service.GetStockHistory(uint(warehouseID), uint(productID), limit, page)
	if err != nil {
		utils.InternalServerError(c, "Failed to fetch history", err.Error())
		return
	}

	utils.OKWithMeta(c, "History fetched successfully", history, utils.Meta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: (int(total) + limit - 1) / limit,
	})
}

// CreateStockIn godoc
// @Summary      Input Stock (Manual)
// @Description  Manually add stock (e.g. from supplier without PO, or found items)
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateStockInRequest  true  "Request Body"
// @Success      201  {object}  utils.Response
// @Router       /stocks/in [post]
func (h *StockHandler) CreateStockIn(c *gin.Context) {
	var req dto.CreateStockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request", err.Error())
		return
	}

	userID := c.GetUint("userID") // Assumes Auth Middleware sets this
	if err := h.service.CreateStockIn(userID, req); err != nil {
		utils.InternalServerError(c, "Failed to create stock in", err.Error())
		return
	}

	utils.Created(c, "Stock in recorded successfully", nil)
}

// CreateStockOut godoc
// @Summary      Output Stock (Manual)
// @Description  Manually remove stock (e.g. usage, damage, expired)
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateStockOutRequest  true  "Request Body"
// @Success      201  {object}  utils.Response
// @Router       /stocks/out [post]
func (h *StockHandler) CreateStockOut(c *gin.Context) {
	var req dto.CreateStockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request", err.Error())
		return
	}

	userID := c.GetUint("userID")
	if err := h.service.CreateStockOut(userID, req); err != nil {
		utils.InternalServerError(c, "Failed to create stock out", err.Error())
		return
	}

	utils.Created(c, "Stock out recorded successfully", nil)
}

// CreateStockOpname godoc
// @Summary      Stock Adjustment (Opname)
// @Description  Adjust stock to match physical count
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateStockOpnameRequest  true  "Request Body"
// @Success      201  {object}  utils.Response
// @Router       /stocks/adjustment [post]
func (h *StockHandler) CreateStockOpname(c *gin.Context) {
	var req dto.CreateStockOpnameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request", err.Error())
		return
	}

	userID := c.GetUint("userID")
	if err := h.service.CreateStockOpname(userID, req); err != nil {
		utils.InternalServerError(c, "Failed to create stock adjustment", err.Error())
		return
	}

	utils.Created(c, "Stock adjustment recorded successfully", nil)
}

// CreateStockTransfer godoc
// @Summary      Transfer Stock
// @Description  Transfer stock between warehouses
// @Tags         stocks
// @Accept       json
// @Produce      json
// @Param        req  body      dto.CreateStockTransferRequest  true  "Request Body"
// @Success      201  {object}  utils.Response
// @Router       /stocks/transfer [post]
func (h *StockHandler) CreateStockTransfer(c *gin.Context) {
	var req dto.CreateStockTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request", err.Error())
		return
	}

	userID := c.GetUint("userID")
	if err := h.service.CreateStockTransfer(userID, req); err != nil {
		utils.InternalServerError(c, "Failed to transfer stock", err.Error())
		return
	}

	utils.Created(c, "Stock transfer recorded successfully", nil)
}
