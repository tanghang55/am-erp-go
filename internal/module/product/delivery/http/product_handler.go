package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/product/domain"
	"am-erp-go/internal/module/product/usecase"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
	imageUsecase   ProductImageUsecase
	auditLogger    ProductAuditLogger
}

type ProductImageUsecase interface {
	ListProductImages(productID uint64) ([]domain.ProductImage, error)
	SaveProductImages(productID uint64, urls []string) ([]domain.ProductImage, error)
}

type ProductAuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

func NewProductHandler(productUsecase *usecase.ProductUsecase, imageUsecase ProductImageUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		imageUsecase:   imageUsecase,
	}
}

func (h *ProductHandler) BindAuditLogger(logger ProductAuditLogger) {
	h.auditLogger = logger
}

// ==================== Product ====================

// ListProducts 获取产品列表
func (h *ProductHandler) ListProducts(c *gin.Context) {
	params := &domain.ProductListParams{
		Page:              parseIntOrDefault(c.Query("page"), 1),
		PageSize:          parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:           c.Query("keyword"),
		Marketplace:       c.Query("marketplace"),
		Status:            c.Query("status"),
		OnlyParentless:    c.Query("only_parentless") == "true",
		OnlyStandalone:    c.Query("only_standalone") == "true",
		OnlyWithPackaging: c.Query("only_with_packaging") == "true",
		ExcludeComboChild: c.Query("exclude_combo_child") == "true",
	}
	if statuses := parseCSVQuery(c.Query("statuses")); len(statuses) > 0 {
		params.Statuses = statuses
	}

	if supplierID := c.Query("supplier_id"); supplierID != "" {
		if id, err := strconv.ParseUint(supplierID, 10, 64); err == nil {
			params.SupplierID = &id
		}
	}

	if brandID := c.Query("brand_id"); brandID != "" {
		if id, err := strconv.ParseUint(brandID, 10, 64); err == nil {
			params.BrandID = &id
		}
	}

	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := strconv.ParseUint(categoryID, 10, 64); err == nil {
			params.CategoryID = &id
		}
	}

	if packingRequired := c.Query("packing_required"); packingRequired != "" {
		value := uint8(0)
		switch strings.ToLower(strings.TrimSpace(packingRequired)) {
		case "1", "true":
			value = 1
		case "0", "false":
			value = 0
		default:
			value = 0
		}
		params.PackingRequired = &value
	}

	if parentID := c.Query("parent_id"); parentID != "" {
		if id, err := strconv.ParseUint(parentID, 10, 64); err == nil {
			params.ParentID = &id
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

func parseCSVQuery(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		item := strings.ToUpper(strings.TrimSpace(part))
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
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
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	var product domain.Product
	if err := json.Unmarshal(body, &product); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	h.applyProductFlagDefaults(body, &product)

	effects, err := h.productUsecase.CreateProductWithEffects(&product)
	if err != nil {
		if errors.Is(err, usecase.ErrSupplierRequired) || errors.Is(err, usecase.ErrProductAsinInvalid) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	h.recordAudit(c, "CREATE_PRODUCT", "Product", fmt.Sprintf("%d", product.ID), nil, product)
	h.recordAutoQuoteAudit(c, product, effects)
	response.Success(c, product)
}

// UpdateProduct 更新产品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	var product domain.Product
	if err := json.Unmarshal(body, &product); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	before, _ := h.productUsecase.GetProduct(id)
	if before == nil {
		response.NotFound(c, "product not found")
		return
	}
	product.ID = id
	h.mergeProductPersistedFields(body, before, &product)
	effects, err := h.productUsecase.UpdateProductWithEffects(&product)
	if err != nil {
		if errors.Is(err, usecase.ErrSupplierRequired) ||
			errors.Is(err, usecase.ErrProductAsinInvalid) ||
			errors.Is(err, usecase.ErrComboMainRequiresActiveChildren) ||
			errors.Is(err, usecase.ErrComboActiveRequiresActiveChildren) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	after, _ := h.productUsecase.GetProduct(id)
	h.recordAuditIfChanged(c, "UPDATE_PRODUCT", "Product", fmt.Sprintf("%d", product.ID), before, after)
	h.recordAutoQuoteAudit(c, product, effects)
	if after != nil {
		response.Success(c, after)
		return
	}
	response.Success(c, product)
}

func (h *ProductHandler) applyProductFlagDefaults(body []byte, product *domain.Product) {
	if product == nil {
		return
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return
	}

	if _, ok := fields["is_inspection_required"]; !ok {
		product.IsInspectionRequired = 1
	}
	if _, ok := fields["is_packing_required"]; !ok {
		product.IsPackingRequired = 1
	}
}

func (h *ProductHandler) mergeProductPersistedFields(body []byte, before *domain.Product, product *domain.Product) {
	if before == nil || product == nil {
		return
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return
	}

	if _, ok := fields["parent_id"]; !ok {
		product.ParentID = before.ParentID
	}
	if _, ok := fields["combo_id"]; !ok {
		product.ComboID = before.ComboID
	}
	if _, ok := fields["is_combo_main"]; !ok {
		product.IsComboMain = before.IsComboMain
	}
	if _, ok := fields["is_inspection_required"]; !ok {
		product.IsInspectionRequired = before.IsInspectionRequired
	}
	if _, ok := fields["is_packing_required"]; !ok {
		product.IsPackingRequired = before.IsPackingRequired
	}
}

// DeleteProduct 删除产品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	before, _ := h.productUsecase.GetProduct(id)
	if err := h.productUsecase.DeleteProduct(id); err != nil {
		if errors.Is(err, usecase.ErrProductReferenced) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	h.recordAudit(c, "DELETE_PRODUCT", "Product", fmt.Sprintf("%d", id), before, nil)
	response.Success(c, nil)
}

// ==================== Product Config ====================

func (h *ProductHandler) ListProductConfigs(c *gin.Context) {
	params := &domain.ProductConfigListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 100),
		Keyword:  c.Query("keyword"),
		Status:   c.Query("status"),
	}
	if configType := c.Query("config_type"); configType != "" {
		params.ConfigType = domain.ProductConfigType(configType)
	}

	items, total, err := h.productUsecase.ListProductConfigs(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, items, total, params.Page, params.PageSize)
}

func (h *ProductHandler) CreateProductConfig(c *gin.Context) {
	var item domain.ProductConfigItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.productUsecase.CreateProductConfig(&item); err != nil {
		if errors.Is(err, usecase.ErrProductConfigTypeUnsupported) || errors.Is(err, usecase.ErrProductConfigCodeInvalid) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	h.recordAudit(c, "CREATE_PRODUCT_CONFIG", "ProductConfig", fmt.Sprintf("%d", item.ID), nil, item)
	response.Success(c, item)
}

func (h *ProductHandler) UpdateProductConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var item domain.ProductConfigItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	item.ID = id
	before, _ := h.productUsecase.GetProductConfig(id)

	if err := h.productUsecase.UpdateProductConfig(&item); err != nil {
		if errors.Is(err, usecase.ErrProductConfigTypeUnsupported) ||
			errors.Is(err, usecase.ErrProductConfigCodeInvalid) ||
			errors.Is(err, usecase.ErrProductConfigCodeImmutable) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	after, _ := h.productUsecase.GetProductConfig(id)
	h.recordAuditIfChanged(c, "UPDATE_PRODUCT_CONFIG", "ProductConfig", fmt.Sprintf("%d", item.ID), before, after)
	if after != nil {
		response.Success(c, after)
		return
	}
	response.Success(c, item)
}

func (h *ProductHandler) DeleteProductConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	before, _ := h.productUsecase.GetProductConfig(id)
	if err := h.productUsecase.DeleteProductConfig(id); err != nil {
		if errors.Is(err, usecase.ErrProductConfigReferenced) || errors.Is(err, usecase.ErrProductConfigSystemFixed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	h.recordAudit(c, "DELETE_PRODUCT_CONFIG", "ProductConfig", fmt.Sprintf("%d", id), before, nil)
	response.Success(c, nil)
}

func (h *ProductHandler) ListProductCategories(c *gin.Context) {
	items, err := h.productUsecase.ListProductCategories()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *ProductHandler) CreateProductCategory(c *gin.Context) {
	var item domain.ProductCategory
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.productUsecase.CreateProductCategory(&item); err != nil {
		if errors.Is(err, usecase.ErrProductCategoryParentRequired) || errors.Is(err, usecase.ErrProductCategoryParentInvalid) || errors.Is(err, usecase.ErrProductCategoryCodeInvalid) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	h.recordAudit(c, "CREATE_PRODUCT_CATEGORY", "ProductCategory", fmt.Sprintf("%d", item.ID), nil, item)
	response.Success(c, item)
}

func (h *ProductHandler) UpdateProductCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var item domain.ProductCategory
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	item.ID = id
	before := h.findCategoryByID(id)
	if err := h.productUsecase.UpdateProductCategory(&item); err != nil {
		if errors.Is(err, usecase.ErrProductCategoryParentInvalid) ||
			errors.Is(err, usecase.ErrProductCategoryCodeInvalid) ||
			errors.Is(err, usecase.ErrProductCategoryCodeImmutable) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	after := h.findCategoryByID(id)
	h.recordAuditIfChanged(c, "UPDATE_PRODUCT_CATEGORY", "ProductCategory", fmt.Sprintf("%d", item.ID), before, after)
	if after != nil {
		response.Success(c, after)
		return
	}
	response.Success(c, item)
}

func (h *ProductHandler) DeleteProductCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	before := h.findCategoryByID(id)
	if err := h.productUsecase.DeleteProductCategory(id); err != nil {
		if errors.Is(err, usecase.ErrProductCategoryHasChildren) ||
			errors.Is(err, usecase.ErrProductCategoryReferenced) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	h.recordAudit(c, "DELETE_PRODUCT_CATEGORY", "ProductCategory", fmt.Sprintf("%d", id), before, nil)
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
	beforeItems, _ := h.imageUsecase.ListProductImages(id)

	items, err := h.imageUsecase.SaveProductImages(id, req.ImageUrls)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	urls := make([]string, 0, len(items))
	for _, item := range items {
		urls = append(urls, item.ImageUrl)
	}

	before := make([]string, 0, len(beforeItems))
	for _, item := range beforeItems {
		before = append(before, item.ImageUrl)
	}
	h.recordAuditIfChanged(c, "SAVE_PRODUCT_IMAGES", "ProductImage", fmt.Sprintf("%d", id), before, urls)
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
		HasChildren: c.Query("has_children"),
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

	h.recordAudit(c, "CREATE_PRODUCT_GROUP", "ProductGroup", fmt.Sprintf("%d", parent.ID), nil, parent)
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
	before, _ := h.productUsecase.GetProductParent(id)
	if err := h.productUsecase.UpdateProductParent(&parent); err != nil {
		if errors.Is(err, usecase.ErrProductParentImmutableIdentity) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	after, _ := h.productUsecase.GetProductParent(id)
	h.recordAuditIfChanged(c, "UPDATE_PRODUCT_GROUP", "ProductGroup", fmt.Sprintf("%d", parent.ID), before, after)
	if after != nil {
		response.Success(c, after)
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

	before, _ := h.productUsecase.GetProductParent(id)
	if err := h.productUsecase.DeleteProductParent(id); err != nil {
		if errors.Is(err, usecase.ErrProductParentHasChildren) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	h.recordAudit(c, "DELETE_PRODUCT_GROUP", "ProductGroup", fmt.Sprintf("%d", id), before, nil)
	response.Success(c, nil)
}

func (h *ProductHandler) AttachProductParentChildren(c *gin.Context) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		ChildIDs []uint64 `json:"child_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	before, _ := h.productUsecase.GetProductParent(parentID)
	parent, err := h.productUsecase.AttachProductParentChildren(parentID, req.ChildIDs)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrProductParentChildIDsRequired),
			errors.Is(err, usecase.ErrProductParentChildNotFound),
			errors.Is(err, usecase.ErrProductParentChildAlreadyAssigned),
			errors.Is(err, usecase.ErrProductParentChildMarketplaceMismatch):
			response.BadRequest(c, err.Error())
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	h.recordAuditIfChanged(c, "ATTACH_PRODUCT_GROUP_CHILDREN", "ProductGroup", fmt.Sprintf("%d", parentID), before, parent)
	response.Success(c, parent)
}

func (h *ProductHandler) DetachProductParentChild(c *gin.Context) {
	parentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	childID, err := strconv.ParseUint(c.Param("childId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid child id")
		return
	}

	before, _ := h.productUsecase.GetProductParent(parentID)
	parent, err := h.productUsecase.DetachProductParentChild(parentID, childID)
	if err != nil {
		if errors.Is(err, usecase.ErrProductParentChildNotFound) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	h.recordAuditIfChanged(c, "DETACH_PRODUCT_GROUP_CHILD", "ProductGroup", fmt.Sprintf("%d", parentID), before, parent)
	response.Success(c, parent)
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
	before, _ := h.productUsecase.GetProductPackagingItems(productID)

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

	h.recordAuditIfChanged(c, "SAVE_PRODUCT_PACKAGING", "ProductPackaging", fmt.Sprintf("%d", productID), before, items)
	response.Success(c, items)
}

func (h *ProductHandler) recordAudit(c *gin.Context, action, entityType, entityID string, before, after any) {
	if h.auditLogger == nil || c == nil {
		return
	}
	_ = h.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Product",
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
}

func (h *ProductHandler) recordAutoQuoteAudit(c *gin.Context, product domain.Product, effects *usecase.ProductUpsertEffects) {
	if effects == nil || !effects.AutoQuoteCreated || product.SupplierID == nil || *product.SupplierID == 0 {
		return
	}
	h.recordAudit(c, "AUTO_CREATE_DEFAULT_SUPPLIER_QUOTE", "ProductSupplierQuote", fmt.Sprintf("%d:%d", product.ID, *product.SupplierID), nil, map[string]any{
		"product_id":  product.ID,
		"supplier_id": *product.SupplierID,
		"price":       normalizeAuditNumber(product.UnitCost),
		"status":      "PENDING",
		"remark":      "产品默认供应商自动生成，待报价",
	})
}

func normalizeAuditNumber(value *float64) any {
	if value == nil {
		return 0
	}
	return *value
}

func (h *ProductHandler) recordAuditIfChanged(c *gin.Context, action, entityType, entityID string, before, after any) {
	beforeDiff, afterDiff, changed := buildAuditDiff(before, after)
	if !changed {
		return
	}
	h.recordAudit(c, action, entityType, entityID, beforeDiff, afterDiff)
}

func (h *ProductHandler) findCategoryByID(id uint64) *domain.ProductCategory {
	items, err := h.productUsecase.ListProductCategories()
	if err != nil {
		return nil
	}
	var walk func(nodes []domain.ProductCategory) *domain.ProductCategory
	walk = func(nodes []domain.ProductCategory) *domain.ProductCategory {
		for _, item := range nodes {
			if item.ID == id {
				copyItem := item
				return &copyItem
			}
			if len(item.Children) == 0 {
				continue
			}
			children := make([]domain.ProductCategory, 0, len(item.Children))
			for _, child := range item.Children {
				if child == nil {
					continue
				}
				children = append(children, *child)
			}
			if found := walk(children); found != nil {
				return found
			}
		}
		return nil
	}
	return walk(items)
}

// ==================== Helper ====================

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}

func hasMeaningfulAuditChange(before, after any) bool {
	return !reflect.DeepEqual(normalizeAuditValue(before), normalizeAuditValue(after))
}

func buildAuditDiff(before, after any) (any, any, bool) {
	normalizedBefore := normalizeAuditValue(before)
	normalizedAfter := normalizeAuditValue(after)
	beforeDiff, afterDiff, changed := diffAuditValues(normalizedBefore, normalizedAfter)
	return beforeDiff, afterDiff, changed
}

func normalizeAuditValue(value any) any {
	if value == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var normalized any
	if err := json.Unmarshal(raw, &normalized); err != nil {
		return value
	}
	return scrubAuditValue(normalized)
}

func scrubAuditValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		cleaned := make(map[string]any, len(typed))
		for key, item := range typed {
			switch key {
			case "gmt_create", "gmt_modified", "created_at", "updated_at", "created_by", "updated_by",
				"reference_count", "deletable", "delete_block_reason", "updated_by_name",
				"inventory_available", "inventory_reserved", "inventory_inbound":
				continue
			}
			cleaned[key] = scrubAuditValue(item)
		}
		return cleaned
	case []any:
		items := make([]any, 0, len(typed))
		for _, item := range typed {
			items = append(items, scrubAuditValue(item))
		}
		return items
	default:
		return value
	}
}

func diffAuditValues(before, after any) (any, any, bool) {
	if reflect.DeepEqual(before, after) {
		return nil, nil, false
	}

	beforeMap, beforeIsMap := before.(map[string]any)
	afterMap, afterIsMap := after.(map[string]any)
	if beforeIsMap && afterIsMap {
		beforeDiff := map[string]any{}
		afterDiff := map[string]any{}
		keys := make(map[string]struct{}, len(beforeMap)+len(afterMap))
		for key := range beforeMap {
			keys[key] = struct{}{}
		}
		for key := range afterMap {
			keys[key] = struct{}{}
		}
		for key := range keys {
			childBefore, childAfter, childChanged := diffAuditValues(beforeMap[key], afterMap[key])
			if !childChanged {
				continue
			}
			beforeDiff[key] = childBefore
			afterDiff[key] = childAfter
		}
		if len(beforeDiff) == 0 && len(afterDiff) == 0 {
			return nil, nil, false
		}
		return beforeDiff, afterDiff, true
	}

	return before, after, true
}
