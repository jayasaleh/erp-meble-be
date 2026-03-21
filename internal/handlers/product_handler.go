package handlers

import (
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary      Create new product
// @Description  Create a new product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      dto.CreateProductRequest  true  "Product data"
// @Success      201      {object}  utils.Response{data=dto.ProductResponse}
// @Failure      400      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Security     BearerAuth
// @Router       /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	userID := utils.GetUserIDValidity(c)
	product, err := h.productService.CreateProduct(&req, userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.Created(c, "Produk berhasil dibuat", product)
}

// GetProduct godoc
// @Summary      Get product by ID
// @Description  Get product details by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.Response{data=dto.ProductResponse}
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Security     BearerAuth
// @Router       /api/v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID", nil)
		return
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Produk berhasil ditemukan", product)
}

// ListProducts godoc
// @Summary      List products
// @Description  Get list of products with filters and pagination
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "Page number"  default(1)
// @Param        limit       query     int     false  "Items per page"  default(10)
// @Param        search      query     string  false  "Search by name, SKU, or barcode"
// @Param        kategori    query     string  false  "Filter by category"
// @Param        merek       query     string  false  "Filter by brand"
// @Param        id_pemasok  query     int     false  "Filter by supplier ID"
// @Param        aktif       query     bool    false  "Filter by active status"
// @Param        stok_rendah query     bool    false  "Filter low stock products"
// @Success      200         {object}  utils.Response{data=dto.ProductListResponse}
// @Failure      400         {object}  utils.Response
// @Failure      500         {object}  utils.Response
// @Security     BearerAuth
// @Router       /api/v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var req dto.ProductListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "Invalid query parameters", err)
		return
	}

	products, err := h.productService.ListProducts(&req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Daftar produk berhasil diambil", products)
}

// UpdateProduct godoc
// @Summary      Update product
// @Description  Update product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "Product ID"
// @Param        product  body      dto.UpdateProductRequest  true  "Product data to update"
// @Success      200      {object}  utils.Response{data=dto.ProductResponse}
// @Failure      400      {object}  utils.Response
// @Failure      404      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Security     BearerAuth
// @Router       /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID", nil)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request body", err)
		return
	}

	userID := utils.GetUserIDValidity(c)
	product, err := h.productService.UpdateProduct(uint(id), &req, userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Produk berhasil diupdate", product)
}

// DeleteProduct godoc
// @Summary      Delete product
// @Description  Soft delete product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Security     BearerAuth
// @Router       /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID", nil)
		return
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Produk berhasil dihapus", nil)
}

// Helper function to handle errors
func (h *ProductHandler) handleError(c *gin.Context, err error) {
	switch err.Error() {
	case "produk tidak ditemukan":
		utils.NotFound(c, err.Error())
	case "SKU sudah digunakan", "barcode sudah digunakan":
		utils.Conflict(c, err.Error())
	case "harga jual tidak boleh lebih kecil dari harga modal":
		utils.BadRequest(c, err.Error(), nil)
	case "tidak dapat menghapus produk yang masih memiliki stok":
		utils.Conflict(c, err.Error())
	default:
		utils.InternalServerError(c, "Terjadi kesalahan server", err)
	}
}
