package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/inventory/usecase"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	usecase *usecase.InventoryUsecase
}

func NewInventoryHandler(usecase *usecase.InventoryUsecase) *InventoryHandler {
	return &InventoryHandler{usecase: usecase}
}

// ListMovements 获取库存流水列表
func (h *InventoryHandler) ListMovements(c *gin.Context) {
	params := &domain.MovementListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	if skuID := c.Query("sku_id"); skuID != "" {
		if id, err := strconv.ParseUint(skuID, 10, 64); err == nil {
			params.SkuID = &id
		}
	}

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		if id, err := strconv.ParseUint(warehouseID, 10, 64); err == nil {
			params.WarehouseID = &id
		}
	}

	if movementType := c.Query("movement_type"); movementType != "" {
		mt := domain.MovementType(movementType)
		params.MovementType = &mt
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		params.DateFrom = &dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		params.DateTo = &dateTo
	}

	movements, total, err := h.usecase.ListMovements(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"data":  movements,
			"total": total,
		},
	})
}

// GetMovement 获取流水详情
func (h *InventoryHandler) GetMovement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	movement, err := h.usecase.GetMovement(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "movement not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    movement,
	})
}

type createMovementRequest struct {
	SkuID           uint64  `json:"sku_id" binding:"required"`
	WarehouseID     uint64  `json:"warehouse_id" binding:"required"`
	Quantity        int     `json:"quantity" binding:"required"`
	ReferenceType   *string `json:"reference_type"`
	ReferenceID     *uint64 `json:"reference_id"`
	ReferenceNumber *string `json:"reference_number"`
	UnitCost        *float64 `json:"unit_cost"`
	Remark          *string `json:"remark"`
	OperatorID      *uint64 `json:"operator_id"`
	OperatedAt      *string `json:"operated_at"`
}

// CreatePurchaseReceipt 采购入库
func (h *InventoryHandler) CreatePurchaseReceipt(c *gin.Context) {
	h.createMovement(c, domain.MovementTypePurchaseReceipt)
}

// CreateSalesShipment 销售出库
func (h *InventoryHandler) CreateSalesShipment(c *gin.Context) {
	h.createMovement(c, domain.MovementTypeSalesShipment)
}

// CreateStockTakeAdjustment 盘点调整
func (h *InventoryHandler) CreateStockTakeAdjustment(c *gin.Context) {
	h.createMovement(c, domain.MovementTypeStockTakeAdjustment)
}

// CreateManualAdjustment 手工调整
func (h *InventoryHandler) CreateManualAdjustment(c *gin.Context) {
	h.createMovement(c, domain.MovementTypeManualAdjustment)
}

// CreateDamageWriteOff 损坏报损
func (h *InventoryHandler) CreateDamageWriteOff(c *gin.Context) {
	h.createMovement(c, domain.MovementTypeDamageWriteOff)
}

// CreateReturnReceipt 退货入库
func (h *InventoryHandler) CreateReturnReceipt(c *gin.Context) {
	h.createMovement(c, domain.MovementTypeReturnReceipt)
}

func (h *InventoryHandler) createMovement(c *gin.Context, movementType domain.MovementType) {
	var req createMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var operatedAt *time.Time
	if req.OperatedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.OperatedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid operated_at format"})
			return
		}
		operatedAt = &t
	}

	params := &domain.CreateMovementParams{
		SkuID:           req.SkuID,
		WarehouseID:     req.WarehouseID,
		MovementType:    movementType,
		Quantity:        req.Quantity,
		ReferenceType:   req.ReferenceType,
		ReferenceID:     req.ReferenceID,
		ReferenceNumber: req.ReferenceNumber,
		UnitCost:        req.UnitCost,
		Remark:          req.Remark,
		OperatorID:      req.OperatorID,
		OperatedAt:      operatedAt,
	}

	movement, err := h.usecase.CreateMovement(c, params)
	if err != nil {
		if errors.Is(err, usecase.ErrInsufficientStock) {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "insufficient stock"})
			return
		}
		if errors.Is(err, usecase.ErrInvalidQuantity) {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid quantity"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    movement,
	})
}

type transferRequest struct {
	SkuID           uint64   `json:"sku_id" binding:"required"`
	FromWarehouseID uint64   `json:"from_warehouse_id" binding:"required"`
	ToWarehouseID   uint64   `json:"to_warehouse_id" binding:"required"`
	Quantity        uint     `json:"quantity" binding:"required"`
	UnitCost        *float64 `json:"unit_cost"`
	Remark          *string  `json:"remark"`
	OperatorID      *uint64  `json:"operator_id"`
	ReferenceType   *string  `json:"reference_type"`
	ReferenceNumber *string  `json:"reference_number"`
}

// CreateTransfer 仓库间调拨
func (h *InventoryHandler) CreateTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	params := &domain.TransferParams{
		SkuID:           req.SkuID,
		FromWarehouseID: req.FromWarehouseID,
		ToWarehouseID:   req.ToWarehouseID,
		Quantity:        req.Quantity,
		UnitCost:        req.UnitCost,
		Remark:          req.Remark,
		OperatorID:      req.OperatorID,
		ReferenceType:   req.ReferenceType,
		ReferenceNumber: req.ReferenceNumber,
	}

	if err := h.usecase.Transfer(c, params); err != nil {
		if errors.Is(err, usecase.ErrInsufficientStock) {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "insufficient stock"})
			return
		}
		if errors.Is(err, usecase.ErrInvalidQuantity) {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid quantity"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ListBalances 获取库存余额列表
func (h *InventoryHandler) ListBalances(c *gin.Context) {
	params := &domain.BalanceListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		if id, err := strconv.ParseUint(warehouseID, 10, 64); err == nil {
			params.WarehouseID = &id
		}
	}

	if skuID := c.Query("sku_id"); skuID != "" {
		if id, err := strconv.ParseUint(skuID, 10, 64); err == nil {
			params.SkuID = &id
		}
	}

	if lowStock := c.Query("low_stock"); lowStock == "true" {
		b := true
		params.LowStock = &b
	}

	if threshold := c.Query("low_stock_threshold"); threshold != "" {
		if t, err := strconv.ParseUint(threshold, 10, 32); err == nil {
			th := uint(t)
			params.LowStockThreshold = &th
		}
	}

	if zeroStock := c.Query("zero_stock"); zeroStock == "true" {
		b := true
		params.ZeroStock = &b
	}

	if keyword := c.Query("keyword"); keyword != "" {
		params.Keyword = &keyword
	}

	balances, total, err := h.usecase.ListBalances(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"data":  balances,
			"total": total,
		},
	})
}

// GetSkuBalance 获取SKU在指定仓库的库存
func (h *InventoryHandler) GetSkuBalance(c *gin.Context) {
	skuIDStr := c.Param("sku_id")
	warehouseIDStr := c.Param("warehouse_id")

	skuID, err := strconv.ParseUint(skuIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid sku_id"})
		return
	}

	warehouseID, err := strconv.ParseUint(warehouseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid warehouse_id"})
		return
	}

	balance, err := h.usecase.GetSkuBalance(skuID, warehouseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "balance not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    balance,
	})
}

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}

// ==================== 库存状态流转API ====================

type stockTransitionRequest struct {
	SkuID           uint64   `json:"sku_id" binding:"required"`
	WarehouseID     uint64   `json:"warehouse_id" binding:"required"`
	Quantity        uint     `json:"quantity" binding:"required"`
	UnitCost        *float64 `json:"unit_cost"`
	Remark          *string  `json:"remark"`
	OperatorID      *uint64  `json:"operator_id"`
	ReferenceType   *string  `json:"reference_type"`
	ReferenceID     *uint64  `json:"reference_id"`
	ReferenceNumber *string  `json:"reference_number"`
}

// CreatePurchaseShip 供应商发货 → 采购在途
func (h *InventoryHandler) CreatePurchaseShip(c *gin.Context) {
	h.createStockTransition(c, domain.MovementTypePurchaseShip)
}

// CreateWarehouseReceive 到仓收货: 采购在途 → 待检
func (h *InventoryHandler) CreateWarehouseReceive(c *gin.Context) {
	h.createStockTransition(c, domain.MovementTypeWarehouseReceive)
}

// CreateInspectionPass 质检通过: 待检 → 原料库存
func (h *InventoryHandler) CreateInspectionPass(c *gin.Context) {
	h.createStockTransition(c, domain.MovementTypeInspectionPass)
}

// CreateInspectionFail 质检不合格: 待检 → 损坏
func (h *InventoryHandler) CreateInspectionFail(c *gin.Context) {
	h.createStockTransition(c, domain.MovementTypeInspectionFail)
}

// CreateAssemblyComplete 组装完成: 原料库存 → 待出库存
func (h *InventoryHandler) CreateAssemblyComplete(c *gin.Context) {
	h.createStockTransition(c, domain.MovementTypeAssemblyComplete)
}

// CreateLogisticsShip 物流发货: 待出库存 → 物流在途
func (h *InventoryHandler) CreateLogisticsShip(c *gin.Context) {
	h.createStockTransition(c, domain.MovementTypeLogisticsShip)
}

// CreatePlatformReceive 平台上架: 物流在途 → 可售库存
func (h *InventoryHandler) CreatePlatformReceive(c *gin.Context) {
	h.createStockTransition(c, domain.MovementTypePlatformReceive)
}

func (h *InventoryHandler) createStockTransition(c *gin.Context, movementType domain.MovementType) {
	var req stockTransitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	params := &domain.CreateMovementParams{
		SkuID:           req.SkuID,
		WarehouseID:     req.WarehouseID,
		MovementType:    movementType,
		Quantity:        int(req.Quantity),
		ReferenceType:   req.ReferenceType,
		ReferenceID:     req.ReferenceID,
		ReferenceNumber: req.ReferenceNumber,
		UnitCost:        req.UnitCost,
		Remark:          req.Remark,
		OperatorID:      req.OperatorID,
	}

	movement, err := h.usecase.CreateMovement(c, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    movement,
	})
}

type returnReceiveRequest struct {
	SkuID           uint64   `json:"sku_id" binding:"required"`
	WarehouseID     uint64   `json:"warehouse_id" binding:"required"`
	Quantity        uint     `json:"quantity" binding:"required"`
	UnitCost        *float64 `json:"unit_cost"`
	Remark          *string  `json:"remark"`
	OperatorID      *uint64  `json:"operator_id"`
	ReferenceType   *string  `json:"reference_type"`
	ReferenceID     *uint64  `json:"reference_id"`
	ReferenceNumber *string  `json:"reference_number"`
}

// CreateReturnReceive 退货入库 → 退货库存
func (h *InventoryHandler) CreateReturnReceive(c *gin.Context) {
	var req returnReceiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	params := &domain.StockTransitionParams{
		SkuID:           req.SkuID,
		WarehouseID:     req.WarehouseID,
		Quantity:        req.Quantity,
		UnitCost:        req.UnitCost,
		Remark:          req.Remark,
		OperatorID:      req.OperatorID,
		ReferenceType:   req.ReferenceType,
		ReferenceID:     req.ReferenceID,
		ReferenceNumber: req.ReferenceNumber,
	}

	movement, err := h.usecase.RecordReturnReceive(c, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    movement,
	})
}

type returnInspectRequest struct {
	SkuID           uint64  `json:"sku_id" binding:"required"`
	WarehouseID     uint64  `json:"warehouse_id" binding:"required"`
	PassQuantity    uint    `json:"pass_quantity"`
	FailQuantity    uint    `json:"fail_quantity"`
	Remark          *string `json:"remark"`
	OperatorID      *uint64 `json:"operator_id"`
	ReferenceType   *string `json:"reference_type"`
	ReferenceNumber *string `json:"reference_number"`
}

// CreateReturnInspect 退货质检: 退货库存 → 待检/损坏
func (h *InventoryHandler) CreateReturnInspect(c *gin.Context) {
	var req returnInspectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if req.PassQuantity == 0 && req.FailQuantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "pass_quantity or fail_quantity is required"})
		return
	}

	params := &domain.ReturnInspectParams{
		SkuID:           req.SkuID,
		WarehouseID:     req.WarehouseID,
		PassQuantity:    req.PassQuantity,
		FailQuantity:    req.FailQuantity,
		Remark:          req.Remark,
		OperatorID:      req.OperatorID,
		ReferenceType:   req.ReferenceType,
		ReferenceNumber: req.ReferenceNumber,
	}

	if err := h.usecase.ReturnInspect(c, params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
