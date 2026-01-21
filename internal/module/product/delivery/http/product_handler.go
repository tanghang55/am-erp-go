package http

import (
	"net/http"
	"strconv"

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
		Page:        parseIntOrDefault(c.Query("page"), 1),
		PageSize:    parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:     c.Query("keyword"),
		Marketplace: c.Query("marketplace"),
		Status:      c.Query("status"),
	}

	if supplierID := c.Query("supplier_id"); supplierID != "" {
		if id, err := strconv.ParseUint(supplierID, 10, 64); err == nil {
			params.SupplierID = &id
		}
	}

	products, total, err := h.productUsecase.ListProducts(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"data":  products,
			"total": total,
		},
	})
}

// GetProduct 获取产品详情
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	product, err := h.productUsecase.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    product,
	})
}

// CreateProduct 创建产品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := h.productUsecase.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    product,
	})
}

// UpdateProduct 更新产品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	product.ID = id
	if err := h.productUsecase.UpdateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    product,
	})
}

// DeleteProduct 删除产品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	if err := h.productUsecase.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ==================== Product Images ====================

// ListProductImages 获取产品图片列表
func (h *ProductHandler) ListProductImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	items, err := h.imageUsecase.ListProductImages(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	urls := make([]string, 0, len(items))
	for _, item := range items {
		urls = append(urls, item.ImageUrl)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    urls,
	})
}

type saveImagesRequest struct {
	ImageUrls []string `json:"image_urls"`
}

// SaveProductImages 保存产品图片排序
func (h *ProductHandler) SaveProductImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req saveImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	items, err := h.imageUsecase.SaveProductImages(id, req.ImageUrls)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	urls := make([]string, 0, len(items))
	for _, item := range items {
		urls = append(urls, item.ImageUrl)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    urls,
	})
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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"data":  parents,
			"total": total,
		},
	})
}

// GetProductParent 获取产品父体详情
func (h *ProductHandler) GetProductParent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	parent, err := h.productUsecase.GetProductParent(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "product parent not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    parent,
	})
}

// CreateProductParent 创建产品父体
func (h *ProductHandler) CreateProductParent(c *gin.Context) {
	var parent domain.ProductParent
	if err := c.ShouldBindJSON(&parent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := h.productUsecase.CreateProductParent(&parent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    parent,
	})
}

// UpdateProductParent 更新产品父体
func (h *ProductHandler) UpdateProductParent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var parent domain.ProductParent
	if err := c.ShouldBindJSON(&parent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	parent.ID = id
	if err := h.productUsecase.UpdateProductParent(&parent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    parent,
	})
}

// DeleteProductParent 删除产品父体
func (h *ProductHandler) DeleteProductParent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	if err := h.productUsecase.DeleteProductParent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ==================== Helper ====================

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}
