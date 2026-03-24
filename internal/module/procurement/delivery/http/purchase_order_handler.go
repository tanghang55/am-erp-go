package http

import (
	"errors"
	"strconv"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/procurement/domain"
	"am-erp-go/internal/module/procurement/usecase"

	"github.com/gin-gonic/gin"
)

type PurchaseOrderHandler struct {
	usecase *usecase.PurchaseOrderUsecase
}

func NewPurchaseOrderHandler(usecase *usecase.PurchaseOrderUsecase) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{usecase: usecase}
}

// ListPurchaseOrders 获取采购单列表
func (h *PurchaseOrderHandler) ListPurchaseOrders(c *gin.Context) {
	params := &domain.PurchaseOrderListParams{
		Page:        parseIntOrDefault(c.Query("page"), 1),
		PageSize:    parseIntOrDefault(c.Query("page_size"), 10),
		Keyword:     c.Query("keyword"),
		Marketplace: c.Query("marketplace"),
	}

	if status := c.Query("status"); status != "" {
		params.Status = domain.PurchaseOrderStatus(status)
	}
	if supplierID := c.Query("supplier_id"); supplierID != "" {
		if id, err := strconv.ParseUint(supplierID, 10, 64); err == nil {
			params.SupplierID = &id
		}
	}

	orders, total, err := h.usecase.List(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, orders, total, params.Page, params.PageSize)
}

// GetPurchaseOrder 获取采购单详情
func (h *PurchaseOrderHandler) GetPurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	order, err := h.usecase.Get(id)
	if err != nil {
		response.NotFound(c, "purchase order not found")
		return
	}

	response.Success(c, order)
}

type purchaseOrderItemRequest struct {
	ProductID  uint64  `json:"product_id"`
	QtyOrdered uint64  `json:"qty_ordered"`
	UnitCost   float64 `json:"unit_cost"`
}

type purchaseOrderUpsertRequest struct {
	SupplierID  *uint64                    `json:"supplier_id"`
	Marketplace string                     `json:"marketplace"`
	Currency    string                     `json:"currency"`
	Remark      string                     `json:"remark"`
	Items       []purchaseOrderItemRequest `json:"items"`
	OperatorID  *uint64                    `json:"operator_id"`
}

type purchaseOrderBatchCreateRequest struct {
	Orders []purchaseOrderUpsertRequest `json:"orders"`
}

func mapPurchaseOrderRequestToDomain(req purchaseOrderUpsertRequest) *domain.PurchaseOrder {
	items := make([]domain.PurchaseOrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.PurchaseOrderItem{
			ProductID:  item.ProductID,
			QtyOrdered: item.QtyOrdered,
			UnitCost:   item.UnitCost,
			Currency:   req.Currency,
		})
	}
	return &domain.PurchaseOrder{
		SupplierID:  req.SupplierID,
		Marketplace: req.Marketplace,
		Currency:    req.Currency,
		Remark:      req.Remark,
		Items:       items,
		CreatedBy:   req.OperatorID,
		UpdatedBy:   req.OperatorID,
	}
}

func resolveOperatorID(c *gin.Context, explicit *uint64) *uint64 {
	if explicit != nil && *explicit != 0 {
		return explicit
	}
	if raw, ok := c.Get(auth.UserIDKey); ok {
		if value, ok := raw.(uint64); ok && value != 0 {
			return &value
		}
	}
	return nil
}

// CreatePurchaseOrder 创建采购单
func (h *PurchaseOrderHandler) CreatePurchaseOrder(c *gin.Context) {
	var req purchaseOrderUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	order := mapPurchaseOrderRequestToDomain(req)
	order.CreatedBy = resolveOperatorID(c, req.OperatorID)
	order.UpdatedBy = order.CreatedBy

	created, err := h.usecase.Create(c, order)
	if err != nil {
		if errors.Is(err, usecase.ErrPurchaseOrderInvalid) ||
			errors.Is(err, usecase.ErrPurchaseOrderMissingItems) ||
			errors.Is(err, usecase.ErrPurchaseOrderMissingProduct) ||
			errors.Is(err, usecase.ErrPurchaseOrderInvalidQty) ||
			errors.Is(err, usecase.ErrPurchaseOrderInvalidUnitCost) ||
			errors.Is(err, usecase.ErrPurchaseOrderMissingSupplier) ||
			errors.Is(err, usecase.ErrPurchaseOrderSplitRequired) ||
			errors.Is(err, usecase.ErrPurchaseOrderComboProviderNeeded) ||
			errors.Is(err, usecase.ErrPurchaseOrderComboNoComponents) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, created)
}

// CreatePurchaseOrderBatch 批量创建采购单，允许后端按供应商拆子单并生成统一批次号。
func (h *PurchaseOrderHandler) CreatePurchaseOrderBatch(c *gin.Context) {
	var req purchaseOrderBatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	orders := make([]*domain.PurchaseOrder, 0, len(req.Orders))
	for _, orderReq := range req.Orders {
		order := mapPurchaseOrderRequestToDomain(orderReq)
		order.CreatedBy = resolveOperatorID(c, orderReq.OperatorID)
		order.UpdatedBy = order.CreatedBy
		orders = append(orders, order)
	}
	created, err := h.usecase.CreateBatch(c, orders)
	if err != nil {
		if errors.Is(err, usecase.ErrPurchaseOrderInvalid) ||
			errors.Is(err, usecase.ErrPurchaseOrderMissingItems) ||
			errors.Is(err, usecase.ErrPurchaseOrderMissingProduct) ||
			errors.Is(err, usecase.ErrPurchaseOrderInvalidQty) ||
			errors.Is(err, usecase.ErrPurchaseOrderInvalidUnitCost) ||
			errors.Is(err, usecase.ErrPurchaseOrderMissingSupplier) ||
			errors.Is(err, usecase.ErrPurchaseOrderComboProviderNeeded) ||
			errors.Is(err, usecase.ErrPurchaseOrderComboNoComponents) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, created)
}

// UpdatePurchaseOrder 更新采购单
func (h *PurchaseOrderHandler) UpdatePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req purchaseOrderUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]domain.PurchaseOrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.PurchaseOrderItem{
			ProductID:  item.ProductID,
			QtyOrdered: item.QtyOrdered,
			UnitCost:   item.UnitCost,
			Currency:   req.Currency,
		})
	}

	order := &domain.PurchaseOrder{
		SupplierID:  req.SupplierID,
		Marketplace: req.Marketplace,
		Currency:    req.Currency,
		Remark:      req.Remark,
		Items:       items,
		UpdatedBy:   resolveOperatorID(c, req.OperatorID),
	}

	updated, err := h.usecase.Update(c, id, order)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, updated)
}

// DeletePurchaseOrder 删除采购单
func (h *PurchaseOrderHandler) DeletePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.usecase.Delete(c, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// SubmitPurchaseOrder 提交采购单
func (h *PurchaseOrderHandler) SubmitPurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.usecase.Submit(c, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

type purchaseOrderShipRequest struct {
	WarehouseID uint64  `json:"warehouse_id"`
	OperatorID  *uint64 `json:"operator_id"`
}

// MarkPurchaseOrderShipped 标记发货
func (h *PurchaseOrderHandler) MarkPurchaseOrderShipped(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req purchaseOrderShipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	params := domain.PurchaseOrderShipParams{
		WarehouseID: req.WarehouseID,
		OperatorID:  resolveOperatorID(c, req.OperatorID),
	}

	if err := h.usecase.MarkShipped(c, id, params); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

type purchaseOrderReceiveRequest struct {
	WarehouseID   uint64            `json:"warehouse_id"`
	ReceivedQties map[uint64]uint64 `json:"received_qties"`
	OperatorID    *uint64           `json:"operator_id"`
}

type purchaseOrderInspectRequest struct {
	PassQties  map[uint64]uint64 `json:"pass_qties"`
	FailQties  map[uint64]uint64 `json:"fail_qties"`
	OperatorID *uint64           `json:"operator_id"`
}

// ReceivePurchaseOrder 到货验收
func (h *PurchaseOrderHandler) ReceivePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req purchaseOrderReceiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	params := domain.PurchaseOrderReceiveParams{
		WarehouseID:   req.WarehouseID,
		ReceivedQties: req.ReceivedQties,
		OperatorID:    resolveOperatorID(c, req.OperatorID),
	}

	if err := h.usecase.Receive(c, id, params); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// InspectPurchaseOrder 采购质检
func (h *PurchaseOrderHandler) InspectPurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req purchaseOrderInspectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	params := domain.PurchaseOrderInspectParams{
		PassQties:  req.PassQties,
		FailQties:  req.FailQties,
		OperatorID: resolveOperatorID(c, req.OperatorID),
	}

	if err := h.usecase.Inspect(c, id, params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ClosePurchaseOrder 完成采购单
func (h *PurchaseOrderHandler) ClosePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.usecase.Close(c, id); err != nil {
		if errors.Is(err, usecase.ErrPurchaseOrderInvalidCompleteStatus) ||
			errors.Is(err, usecase.ErrPurchaseOrderPendingInspection) ||
			errors.Is(err, usecase.ErrPurchaseOrderIncompleteReceipt) ||
			errors.Is(err, usecase.ErrPurchaseOrderInspectionFailed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

type purchaseOrderForceCompleteRequest struct {
	Reason     string  `json:"reason"`
	OperatorID *uint64 `json:"operator_id"`
}

func (h *PurchaseOrderHandler) ForceCompletePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req purchaseOrderForceCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.usecase.ForceComplete(c, id, domain.PurchaseOrderForceCompleteParams{
		Reason:     req.Reason,
		OperatorID: resolveOperatorID(c, req.OperatorID),
	}); err != nil {
		if errors.Is(err, usecase.ErrPurchaseOrderInvalidCompleteStatus) ||
			errors.Is(err, usecase.ErrPurchaseOrderPendingInspection) ||
			errors.Is(err, usecase.ErrPurchaseOrderForceCompleteReason) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}
