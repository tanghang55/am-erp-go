package http

import (
	"strconv"
	"time"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/packaging/domain"
	"am-erp-go/internal/module/packaging/usecase"

	"github.com/gin-gonic/gin"
)

type PackagingHandler struct {
	uc *usecase.PackagingUsecase
}

func NewPackagingHandler(uc *usecase.PackagingUsecase) *PackagingHandler {
	return &PackagingHandler{uc: uc}
}

// ============= Packaging Item Handlers =============

// ListItems 获取包材列表
func (h *PackagingHandler) ListItems(c *gin.Context) {
	params := &domain.PackagingItemListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:  c.Query("keyword"),
		Category: c.Query("category"),
		Status:   c.Query("status"),
		LowStock: c.Query("low_stock") == "true",
	}

	items, total, err := h.uc.ListItems(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Paginated(c, items, total)
}

// GetItem 获取包材详情
func (h *PackagingHandler) GetItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	item, err := h.uc.GetItem(id)
	if err != nil {
		response.NotFound(c, "包材不存在")
		return
	}

	response.Success(c, item)
}

// CreateItem 创建包材
func (h *PackagingHandler) CreateItem(c *gin.Context) {
	var item domain.PackagingItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	if userID == nil {
		response.Unauthorized(c, "未授权")
		return
	}
	item.CreatedBy = *userID

	if err := h.uc.CreateItem(&item); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "创建成功", item)
}

// UpdateItem 更新包材
func (h *PackagingHandler) UpdateItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var item domain.PackagingItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item.ID = id
	if err := h.uc.UpdateItem(&item); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "更新成功", item)
}

// DeleteItem 删除包材
func (h *PackagingHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.uc.DeleteItem(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "删除成功", nil)
}

// GetLowStockItems 获取低库存包材
func (h *PackagingHandler) GetLowStockItems(c *gin.Context) {
	items, err := h.uc.GetLowStockItems()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, items)
}

// ============= Packaging Ledger Handlers =============

// ListLedgers 获取流水列表
func (h *PackagingHandler) ListLedgers(c *gin.Context) {
	params := &domain.PackagingLedgerListParams{
		Page:            parseIntOrDefault(c.Query("page"), 1),
		PageSize:        parseIntOrDefault(c.Query("page_size"), 20),
		TransactionType: c.Query("transaction_type"),
		ReferenceType:   c.Query("reference_type"),
	}

	if itemID := c.Query("packaging_item_id"); itemID != "" {
		if id, err := strconv.ParseUint(itemID, 10, 64); err == nil {
			params.PackagingItemID = &id
		}
	}

	if refID := c.Query("reference_id"); refID != "" {
		if id, err := strconv.ParseUint(refID, 10, 64); err == nil {
			params.ReferenceID = &id
		}
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			params.DateFrom = &t
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			params.DateTo = &t
		}
	}

	ledgers, total, err := h.uc.ListLedgers(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Paginated(c, ledgers, total)
}

// GetLedger 获取流水详情
func (h *PackagingHandler) GetLedger(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	ledger, err := h.uc.GetLedger(id)
	if err != nil {
		response.NotFound(c, "流水不存在")
		return
	}

	response.Success(c, ledger)
}

// CreateInboundLedger 创建入库流水
func (h *PackagingHandler) CreateInboundLedger(c *gin.Context) {
	var ledger domain.PackagingLedger
	if err := c.ShouldBindJSON(&ledger); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	if userID == nil {
		response.Unauthorized(c, "未授权")
		return
	}

	if err := h.uc.CreateInboundLedger(&ledger, *userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "入库成功", ledger)
}

// CreateOutboundLedger 创建出库流水
func (h *PackagingHandler) CreateOutboundLedger(c *gin.Context) {
	var ledger domain.PackagingLedger
	if err := c.ShouldBindJSON(&ledger); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	if userID == nil {
		response.Unauthorized(c, "未授权")
		return
	}

	if err := h.uc.CreateOutboundLedger(&ledger, *userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "出库成功", ledger)
}

// CreateAdjustmentLedger 创建调整流水
func (h *PackagingHandler) CreateAdjustmentLedger(c *gin.Context) {
	var ledger domain.PackagingLedger
	if err := c.ShouldBindJSON(&ledger); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)
	if userID == nil {
		response.Unauthorized(c, "未授权")
		return
	}

	if err := h.uc.CreateAdjustmentLedger(&ledger, *userID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessMessage(c, "调整成功", ledger)
}

// GetUsageSummary 获取使用情况统计
func (h *PackagingHandler) GetUsageSummary(c *gin.Context) {
	var dateFrom, dateTo *time.Time

	if from := c.Query("date_from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			dateFrom = &t
		}
	}

	if to := c.Query("date_to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			dateTo = &t
		}
	}

	summary, err := h.uc.GetUsageSummary(dateFrom, dateTo)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, summary)
}

// Helper functions
func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}

func getUserIDFromContext(c *gin.Context) *uint64 {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(uint64); ok {
			return &id
		}
	}
	return nil
}
