package http

import (
	"errors"
	"strconv"
	"time"

	"am-erp-go/internal/infrastructure/response"
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

	if productID := c.Query("product_id"); productID != "" {
		if id, err := strconv.ParseUint(productID, 10, 64); err == nil {
			params.ProductID = &id
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
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, movements, total, params.Page, params.PageSize)
}

// GetMovement 获取流水详情
func (h *InventoryHandler) GetMovement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	movement, err := h.usecase.GetMovement(id)
	if err != nil {
		response.NotFound(c, "movement not found")
		return
	}

	response.Success(c, movement)
}

type createMovementRequest struct {
	ProductID       uint64   `json:"product_id" binding:"required"`
	WarehouseID     uint64   `json:"warehouse_id" binding:"required"`
	Quantity        int      `json:"quantity" binding:"required"`
	ReferenceType   *string  `json:"reference_type"`
	ReferenceID     *uint64  `json:"reference_id"`
	ReferenceNumber *string  `json:"reference_number"`
	UnitCost        *float64 `json:"unit_cost"`
	Remark          *string  `json:"remark"`
	OperatorID      *uint64  `json:"operator_id"`
	OperatedAt      *string  `json:"operated_at"`
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

func (h *InventoryHandler) createMovement(c *gin.Context, movementType domain.MovementType) {
	var req createMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var operatedAt *time.Time
	if req.OperatedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.OperatedAt)
		if err != nil {
			response.BadRequest(c, "invalid operated_at format")
			return
		}
		operatedAt = &t
	}

	params := &domain.CreateMovementParams{
		ProductID:       req.ProductID,
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
			response.BadRequest(c, "insufficient stock")
			return
		}
		if errors.Is(err, usecase.ErrInvalidQuantity) {
			response.BadRequest(c, "invalid quantity")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, movement)
}

type transferRequest struct {
	ProductID       uint64   `json:"product_id" binding:"required"`
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
		response.BadRequest(c, err.Error())
		return
	}

	params := &domain.TransferParams{
		ProductID:       req.ProductID,
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
			response.BadRequest(c, "insufficient stock")
			return
		}
		if errors.Is(err, usecase.ErrInvalidQuantity) {
			response.BadRequest(c, "invalid quantity")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
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

	if productID := c.Query("product_id"); productID != "" {
		if id, err := strconv.ParseUint(productID, 10, 64); err == nil {
			params.ProductID = &id
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
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, balances, total, params.Page, params.PageSize)
}

// ListLots 获取库存批次列表
func (h *InventoryHandler) ListLots(c *gin.Context) {
	params := &domain.InventoryLotListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		if id, err := strconv.ParseUint(warehouseID, 10, 64); err == nil {
			params.WarehouseID = &id
		}
	}

	if productID := c.Query("product_id"); productID != "" {
		if id, err := strconv.ParseUint(productID, 10, 64); err == nil {
			params.ProductID = &id
		}
	}

	if status := c.Query("status"); status != "" {
		lotStatus := domain.InventoryLotStatus(status)
		params.Status = &lotStatus
	}

	if keyword := c.Query("keyword"); keyword != "" {
		params.Keyword = &keyword
	}

	lots, total, err := h.usecase.ListLots(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, lots, total, params.Page, params.PageSize)
}

// GetProductBalance 获取产品在指定仓库的库存
func (h *InventoryHandler) GetProductBalance(c *gin.Context) {
	productIDStr := c.Param("product_id")
	warehouseIDStr := c.Param("warehouse_id")

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product_id")
		return
	}

	warehouseID, err := strconv.ParseUint(warehouseIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid warehouse_id")
		return
	}

	balance, err := h.usecase.GetProductBalance(productID, warehouseID)
	if err != nil {
		response.NotFound(c, "balance not found")
		return
	}

	response.Success(c, balance)
}

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}

// ==================== 库存状态流转API ====================

type stockTransitionRequest struct {
	ProductID       uint64   `json:"product_id" binding:"required"`
	WarehouseID     uint64   `json:"warehouse_id" binding:"required"`
	Quantity        uint     `json:"quantity" binding:"required"`
	UnitCost        *float64 `json:"unit_cost"`
	Remark          *string  `json:"remark"`
	OperatorID      *uint64  `json:"operator_id"`
	ReferenceType   *string  `json:"reference_type"`
	ReferenceID     *uint64  `json:"reference_id"`
	ReferenceNumber *string  `json:"reference_number"`
}

// CreateAssemblyComplete 组装完成: 原料库存 → 待出库存
func (h *InventoryHandler) CreateAssemblyComplete(c *gin.Context) {
	h.createStockTransition(c, domain.MovementTypeAssemblyComplete)
}

// CreatePlatformReceive 平台上架: 物流在途 → 可售库存
func (h *InventoryHandler) CreatePlatformReceive(c *gin.Context) {
	h.createStockTransition(c, domain.MovementTypePlatformReceive)
}

func (h *InventoryHandler) createStockTransition(c *gin.Context, movementType domain.MovementType) {
	var req stockTransitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.OperatorID == nil {
		req.OperatorID = getUserIDFromContext(c)
	}

	params := &domain.CreateMovementParams{
		ProductID:       req.ProductID,
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
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, movement)
}

func getUserIDFromContext(c *gin.Context) *uint64 {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(uint64); ok {
			return &id
		}
	}
	return nil
}
