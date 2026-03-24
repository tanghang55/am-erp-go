package http

import (
	"errors"
	"strconv"
	"time"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/sales/domain"
	"am-erp-go/internal/module/sales/usecase"

	"github.com/gin-gonic/gin"
)

type SalesOrderHandler struct {
	usecase *usecase.SalesOrderUsecase
}

func NewSalesOrderHandler(uc *usecase.SalesOrderUsecase) *SalesOrderHandler {
	return &SalesOrderHandler{usecase: uc}
}

func (h *SalesOrderHandler) ListSalesOrders(c *gin.Context) {
	params := &domain.SalesOrderListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:  c.Query("keyword"),
	}

	if status := c.Query("status"); status != "" {
		params.Status = domain.SalesOrderStatus(status)
	}
	if marketplace := c.Query("marketplace"); marketplace != "" {
		params.Marketplace = marketplace
	}

	list, total, err := h.usecase.List(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

func (h *SalesOrderHandler) GetSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	order, err := h.usecase.Get(id)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			response.NotFound(c, "sales order not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, order)
}

type salesOrderItemRequest struct {
	ID           uint64  `json:"id"`
	LineNo       uint32  `json:"line_no"`
	ProductID    uint64  `json:"product_id"`
	QtyOrdered   uint64  `json:"qty_ordered"`
	QtyAllocated uint64  `json:"qty_allocated"`
	QtyShipped   uint64  `json:"qty_shipped"`
	QtyReturned  uint64  `json:"qty_returned"`
	UnitPrice    float64 `json:"unit_price"`
	Subtotal     float64 `json:"subtotal"`
	Remark       *string `json:"remark"`
}

type salesOrderUpsertRequest struct {
	OrderNo         string                  `json:"order_no"`
	SourceType      string                  `json:"source_type"`
	ExternalOrderNo *string                 `json:"external_order_no"`
	SalesChannel    *string                 `json:"sales_channel"`
	Marketplace     *string                 `json:"marketplace"`
	OrderDate       string                  `json:"order_date"`
	Currency        string                  `json:"currency"`
	OrderAmount     float64                 `json:"order_amount"`
	Remark          *string                 `json:"remark"`
	Items           []salesOrderItemRequest `json:"items"`
	OperatorID      *uint64                 `json:"operator_id"`
}

func (h *SalesOrderHandler) CreateSalesOrder(c *gin.Context) {
	var req salesOrderUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	orderDate, err := time.Parse(time.RFC3339, req.OrderDate)
	if err != nil {
		response.BadRequest(c, "invalid order_date format")
		return
	}

	items := make([]domain.SalesOrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.SalesOrderItem{
			LineNo:       item.LineNo,
			ProductID:    item.ProductID,
			QtyOrdered:   item.QtyOrdered,
			QtyAllocated: item.QtyAllocated,
			QtyShipped:   item.QtyShipped,
			QtyReturned:  item.QtyReturned,
			UnitPrice:    item.UnitPrice,
			Subtotal:     item.Subtotal,
			Remark:       item.Remark,
		})
	}

	order := &domain.SalesOrder{
		OrderNo:         req.OrderNo,
		SourceType:      req.SourceType,
		ExternalOrderNo: req.ExternalOrderNo,
		SalesChannel:    req.SalesChannel,
		Marketplace:     req.Marketplace,
		OrderDate:       orderDate,
		Currency:        req.Currency,
		OrderAmount:     req.OrderAmount,
		Remark:          req.Remark,
		Items:           items,
		CreatedBy:       req.OperatorID,
		UpdatedBy:       req.OperatorID,
	}

	created, err := h.usecase.Create(c, order)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidQuantity) {
			response.BadRequest(c, "invalid order items")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, created)
}

func (h *SalesOrderHandler) UpdateSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req salesOrderUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	orderDate := time.Now()
	if req.OrderDate != "" {
		parsed, parseErr := time.Parse(time.RFC3339, req.OrderDate)
		if parseErr != nil {
			response.BadRequest(c, "invalid order_date format")
			return
		}
		orderDate = parsed
	}

	items := make([]domain.SalesOrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.SalesOrderItem{
			ID:           item.ID,
			LineNo:       item.LineNo,
			ProductID:    item.ProductID,
			QtyOrdered:   item.QtyOrdered,
			QtyAllocated: item.QtyAllocated,
			QtyShipped:   item.QtyShipped,
			QtyReturned:  item.QtyReturned,
			UnitPrice:    item.UnitPrice,
			Subtotal:     item.Subtotal,
			Remark:       item.Remark,
		})
	}

	updates := &domain.SalesOrder{
		ExternalOrderNo: req.ExternalOrderNo,
		SalesChannel:    req.SalesChannel,
		Marketplace:     req.Marketplace,
		OrderDate:       orderDate,
		Currency:        req.Currency,
		OrderAmount:     req.OrderAmount,
		Remark:          req.Remark,
		Items:           items,
		UpdatedBy:       req.OperatorID,
	}

	updated, updateErr := h.usecase.Update(c, id, updates)
	if updateErr != nil {
		if errors.Is(updateErr, domain.ErrOrderNotFound) {
			response.NotFound(c, "sales order not found")
			return
		}
		if errors.Is(updateErr, domain.ErrInvalidTransition) {
			response.BadRequest(c, "order status does not allow update")
			return
		}
		response.InternalError(c, updateErr.Error())
		return
	}

	response.Success(c, updated)
}

func (h *SalesOrderHandler) ConfirmSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	operatorID := parseOptionalUint64(c)
	if err := h.usecase.Confirm(c, id, operatorID); err != nil {
		h.handleTransitionError(c, err)
		return
	}

	response.Success(c, nil)
}

type allocateRequest struct {
	WarehouseID uint64  `json:"warehouse_id" binding:"required"`
	OperatorID  *uint64 `json:"operator_id"`
	Lines       []struct {
		ItemID       uint64 `json:"item_id"`
		QtyAllocated uint64 `json:"qty_allocated"`
	} `json:"lines" binding:"required"`
}

func (h *SalesOrderHandler) AllocateSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req allocateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	lines := make([]domain.AllocateLine, 0, len(req.Lines))
	for _, line := range req.Lines {
		lines = append(lines, domain.AllocateLine{
			ItemID:       line.ItemID,
			QtyAllocated: line.QtyAllocated,
		})
	}

	err = h.usecase.Allocate(c, id, domain.AllocateParams{
		WarehouseID: req.WarehouseID,
		Lines:       lines,
	}, req.OperatorID)
	if err != nil {
		h.handleTransitionError(c, err)
		return
	}

	response.Success(c, nil)
}

type shipRequest struct {
	WarehouseID uint64  `json:"warehouse_id" binding:"required"`
	OperatorID  *uint64 `json:"operator_id"`
	Lines       []struct {
		ItemID     uint64 `json:"item_id"`
		QtyShipped uint64 `json:"qty_shipped"`
	} `json:"lines" binding:"required"`
}

func (h *SalesOrderHandler) ShipSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req shipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	lines := make([]domain.ShipLine, 0, len(req.Lines))
	for _, line := range req.Lines {
		lines = append(lines, domain.ShipLine{
			ItemID:     line.ItemID,
			QtyShipped: line.QtyShipped,
		})
	}

	err = h.usecase.Ship(c, id, domain.ShipParams{
		WarehouseID: req.WarehouseID,
		Lines:       lines,
	}, req.OperatorID)
	if err != nil {
		h.handleTransitionError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *SalesOrderHandler) DeliverSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	operatorID := parseOptionalUint64(c)
	if err := h.usecase.Deliver(c, id, operatorID); err != nil {
		h.handleTransitionError(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *SalesOrderHandler) CancelSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	operatorID := parseOptionalUint64(c)
	if err := h.usecase.Cancel(c, id, operatorID); err != nil {
		h.handleTransitionError(c, err)
		return
	}
	response.Success(c, nil)
}

type returnRequest struct {
	WarehouseID uint64  `json:"warehouse_id" binding:"required"`
	OperatorID  *uint64 `json:"operator_id"`
	Lines       []struct {
		ItemID      uint64 `json:"item_id"`
		QtyReturned uint64 `json:"qty_returned"`
	} `json:"lines" binding:"required"`
}

func (h *SalesOrderHandler) ReturnSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req returnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	lines := make([]domain.ReturnLine, 0, len(req.Lines))
	for _, line := range req.Lines {
		lines = append(lines, domain.ReturnLine{
			ItemID:      line.ItemID,
			QtyReturned: line.QtyReturned,
		})
	}

	err = h.usecase.Return(c, id, domain.ReturnParams{
		WarehouseID: req.WarehouseID,
		Lines:       lines,
	}, req.OperatorID)
	if err != nil {
		h.handleTransitionError(c, err)
		return
	}

	response.Success(c, nil)
}

func (h *SalesOrderHandler) handleTransitionError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrOrderNotFound) {
		response.NotFound(c, "sales order not found")
		return
	}
	if errors.Is(err, domain.ErrInvalidTransition) || errors.Is(err, domain.ErrInvalidQuantity) || errors.Is(err, domain.ErrItemNotFound) {
		response.BadRequest(c, err.Error())
		return
	}
	response.InternalError(c, err.Error())
}

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}

func parseOptionalUint64(c *gin.Context) *uint64 {
	var payload struct {
		OperatorID *uint64 `json:"operator_id"`
	}
	if err := c.ShouldBindJSON(&payload); err == nil {
		return payload.OperatorID
	}
	return nil
}
