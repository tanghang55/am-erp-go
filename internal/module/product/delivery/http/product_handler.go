package http

import (
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/product/domain"
	"am-erp-go/internal/module/product/usecase"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
	imageUsecase   ProductImageUsecase
}

type ProductImageUsecase interface {
	ListProductImages(productID uint64) ([]domain.ProductImage, error)
	SaveProductImages(productID uint64, urls []string) ([]domain.ProductImage, error)
}

func NewProductHandler(productUsecase *usecase.ProductUsecase, imageUsecase ProductImageUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		imageUsecase:   imageUsecase,
	}
}

// ==================== Product SKU ====================

// ListProducts 获取产品列表
func (h *ProductHandler) ListProducts(c *gin.Context) {
	params := &domain.ProductListParams{
		Page:              parseIntOrDefault(c.Query("page"), 1),
		PageSize:          parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:           c.Query("keyword"),
		Marketplace:       c.Query("marketplace"),
		Status:            c.Query("status"),
		ExcludeComboChild: c.Query("exclude_combo_child") == "true",
	}

	if supplierID := c.Query("supplier_id"); supplierID != "" {
		if id, err := strconv.ParseUint(supplierID, 10, 64); err == nil {
			params.SupplierID = &id
		}
	}

	// 仓库ID筛选（只返回该仓库有库存的产品）
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		if id, err := strconv.ParseUint(warehouseID, 10, 64); err == nil {
			params.WarehouseID = &id
		}
	}

	products, total, err := h.productUsecase.ListProducts(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, products, total, params.Page, params.PageSize)
}

// GetProduct 获取产品详情
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	product, err := h.productUsecase.GetProduct(id)
	if err != nil {
		response.NotFound(c, "product not found")
		return
	}

	response.Success(c, product)
}

// CreateProduct 创建产品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.productUsecase.CreateProduct(&product); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, product)
}

// UpdateProduct 更新产品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	product.ID = id
	if err := h.productUsecase.UpdateProduct(&product); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, product)
}

// DeleteProduct 删除产品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.productUsecase.DeleteProduct(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ==================== Product Images ====================

// ListProductImages 获取产品图片列表
func (h *ProductHandler) ListProductImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	items, err := h.imageUsecase.ListProductImages(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	urls := make([]string, 0, len(items))
	for _, item := range items {
		urls = append(urls, item.ImageUrl)
	}

	response.Success(c, urls)
}

type saveImagesRequest struct {
	ImageUrls []string `json:"image_urls"`
}

// SaveProductImages 保存产品图片排序
func (h *ProductHandler) SaveProductImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req saveImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items, err := h.imageUsecase.SaveProductImages(id, req.ImageUrls)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	urls := make([]string, 0, len(items))
	for _, item := range items {
		urls = append(urls, item.ImageUrl)
	}

	response.Success(c, urls)
}

// ==================== ProductParent ====================

// ListProductParents 获取产品父体列表
func (h *ProductHandler) ListProductParents(c *gin.Context) {
	params := &domain.ProductParentListParams{
		Page:        parseIntOrDefault(c.Query("page"), 1),
		PageSize:    parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:     c.Query("keyword"),
		Marketplace: c.Query("marketplace"),
		Status:      c.Query("status"),
	}

	parents, total, err := h.productUsecase.ListProductParents(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, parents, total, params.Page, params.PageSize)
}

// GetProductParent 获取产品父体详情
func (h *ProductHandler) GetProductParent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	parent, err := h.productUsecase.GetProductParent(id)
	if err != nil {
		response.NotFound(c, "product parent not found")
		return
	}

	response.Success(c, parent)
}

// CreateProductParent 创建产品父体
func (h *ProductHandler) CreateProductParent(c *gin.Context) {
	var parent domain.ProductParent
	if err := c.ShouldBindJSON(&parent); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.productUsecase.CreateProductParent(&parent); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, parent)
}

// UpdateProductParent 更新产品父体
func (h *ProductHandler) UpdateProductParent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var parent domain.ProductParent
	if err := c.ShouldBindJSON(&parent); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	parent.ID = id
	if err := h.productUsecase.UpdateProductParent(&parent); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, parent)
}

// DeleteProductParent 删除产品父体
func (h *ProductHandler) DeleteProductParent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.productUsecase.DeleteProductParent(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ==================== Product Packaging ====================

// GetProductPackagingItems 获取产品的包材配置列表
func (h *ProductHandler) GetProductPackagingItems(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	items, err := h.productUsecase.GetProductPackagingItems(productID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, items)
}

// SaveProductPackagingItems 保存产品的包材配置
func (h *ProductHandler) SaveProductPackagingItems(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	var req struct {
		PackagingItems []domain.ProductPackagingItem `json:"packaging_items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.productUsecase.SaveProductPackagingItems(productID, req.PackagingItems); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 返回保存后的数据
	items, err := h.productUsecase.GetProductPackagingItems(productID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, items)
}

// ==================== Helper ====================

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}
