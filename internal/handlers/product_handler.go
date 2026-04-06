package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"real-erp-mebel/be/internal/dto"
	"real-erp-mebel/be/internal/services"
	"real-erp-mebel/be/internal/utils"
	"strconv"
	"time"

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

// UploadProductImages godoc
// @Summary      Upload multiple product images
// @Description  Upload multiple images for a product replacing existing images
// @Tags         products
// @Accept       multipart/form-data
// @Produce      json
// @Param        id       path      int  true  "Product ID"
// @Param        images   formData  file true  "Multiple image files"
// @Success      200      {object}  utils.Response
// @Failure      400      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Security     BearerAuth
// @Router       /api/v1/products/{id}/images [post]
func (h *ProductHandler) UploadProductImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID", nil)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequest(c, "Gagal memproses form data", err)
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		utils.BadRequest(c, "Tidak ada file gambar yang diupload", nil)
		return
	}

	uploadDir := "uploads/products"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		utils.InternalServerError(c, "Gagal membuat direktori upload", err)
		return
	}

	var imagePaths []string
	for _, fileHeader := range files {
		fileName := fmt.Sprintf("%d_%d_%s", id, time.Now().UnixNano(), fileHeader.Filename)
		filePath := filepath.Join(uploadDir, fileName)

		if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
			utils.InternalServerError(c, "Gagal menyimpan file gambar", err)
			return
		}
		relativeURL := "/uploads/products/" + fileName
		imagePaths = append(imagePaths, relativeURL)
	}

	if err := h.productService.SaveProductImages(uint(id), imagePaths); err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Gambar produk berhasil disimpan", gin.H{"paths": imagePaths})
}

// DeleteProductImage godoc
// @Summary      Delete a single product image
// @Description  Delete a specific image from a product
// @Tags         products
// @Produce      json
// @Param        id       path      int  true  "Product ID"
// @Param        imageId  path      int  true  "Image ID"
// @Success      200      {object}  utils.Response
// @Failure      400      {object}  utils.Response
// @Failure      500      {object}  utils.Response
// @Security     BearerAuth
// @Router       /api/v1/products/{id}/images/{imageId} [delete]
func (h *ProductHandler) DeleteProductImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid product ID", nil)
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid image ID", nil)
		return
	}

	if err := h.productService.DeleteProductImage(uint(productID), uint(imageID)); err != nil {
		h.handleError(c, err)
		return
	}

	utils.OK(c, "Gambar produk berhasil dihapus", nil)
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
