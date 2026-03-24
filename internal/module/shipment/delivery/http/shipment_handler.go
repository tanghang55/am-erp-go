package http

import (
	"errors"
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/shipment/domain"
	"am-erp-go/internal/module/shipment/usecase"

	"github.com/gin-gonic/gin"
)

type ShipmentHandler struct {
	usecase *usecase.ShipmentUsecase
}

func NewShipmentHandler(usecase *usecase.ShipmentUsecase) *ShipmentHandler {
	return &ShipmentHandler{usecase: usecase}
}

func (h *ShipmentHandler) respondShipmentError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if errors.Is(err, usecase.ErrShipmentNotFound) {
		response.NotFound(c, "shipment not found")
		return
	}
	if errors.Is(err, usecase.ErrEmptyItems) ||
		errors.Is(err, usecase.ErrShipmentProductNotFound) ||
		errors.Is(err, usecase.ErrShipmentInactiveProduct) ||
		errors.Is(err, usecase.ErrInsufficientInventory) ||
		errors.Is(err, usecase.ErrInvalidStatus) {
		response.BadRequest(c, err.Error())
		return
	}
	response.InternalError(c, err.Error())
}

// ListShipments 获取发货单列表
func (h *ShipmentHandler) ListShipments(c *gin.Context) {
	params := &domain.ShipmentListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	if status := c.Query("status"); status != "" {
		s := domain.ShipmentStatus(status)
		params.Status = &s
	}
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		if id, err := strconv.ParseUint(warehouseID, 10, 64); err == nil {
			params.WarehouseID = &id
		}
	}
	if orderNumber := c.Query("order_number"); orderNumber != "" {
		params.OrderNumber = &orderNumber
	}
	if trackingNumber := c.Query("tracking_number"); trackingNumber != "" {
		params.TrackingNumber = &trackingNumber
	}
	if keyword := c.Query("keyword"); keyword != "" {
		params.Keyword = &keyword
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		params.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		params.DateTo = &dateTo
	}

	shipments, total, err := h.usecase.List(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, shipments, total, params.Page, params.PageSize)
}

// GetShipment 获取发货单详情
func (h *ShipmentHandler) GetShipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	shipment, err := h.usecase.Get(id)
	if err != nil {
		response.NotFound(c, "shipment not found")
		return
	}

	response.Success(c, shipment)
}

type shipmentItemRequest struct {
	ProductID       uint64   `json:"product_id"`
	QuantityPlanned uint     `json:"quantity_planned"`
	PackageSpecID   *uint64  `json:"package_spec_id"`
	BoxQuantity     *uint    `json:"box_quantity"`
	UnitCost        *float64 `json:"unit_cost"`
	Currency        *string  `json:"currency"`
	Remark          *string  `json:"remark"`
}

type createShipmentRequest struct {
	OrderNumber  *string `json:"order_number"`
	SalesChannel *string `json:"sales_channel"`
	WarehouseID  uint64  `json:"warehouse_id"`

	DestinationWarehouseID *uint64                 `json:"destination_warehouse_id"`
	DestinationType        *domain.DestinationType `json:"destination_type"`
	DestinationName        *string                 `json:"destination_name"`
	DestinationCode        *string                 `json:"destination_code"`
	DestinationContact     *string                 `json:"destination_contact"`
	DestinationPhone       *string                 `json:"destination_phone"`
	DestinationAddress     *string                 `json:"destination_address"`

	LogisticsProviderID  *uint64 `json:"logistics_provider_id"`
	ShippingRateID       *uint64 `json:"shipping_rate_id"`
	TransportMode        *string `json:"transport_mode"`
	Carrier              *string `json:"carrier"`
	TrackingNumber       *string `json:"tracking_number"`
	ExpectedShipDate     *string `json:"expected_ship_date"`
	ExpectedDeliveryDate *string `json:"expected_delivery_date"`

	BoxCount    *uint    `json:"box_count"`
	TotalWeight *float64 `json:"total_weight"`
	TotalVolume *float64 `json:"total_volume"`

	Items         []shipmentItemRequest `json:"items"`
	Remark        *string               `json:"remark"`
	InternalNotes *string               `json:"internal_notes"`
	OperatorID    *uint64               `json:"operator_id"`
}

type updateShipmentRequest struct {
	OrderNumber  *string `json:"order_number"`
	SalesChannel *string `json:"sales_channel"`
	WarehouseID  *uint64 `json:"warehouse_id"`

	DestinationWarehouseID *uint64                 `json:"destination_warehouse_id"`
	DestinationType        *domain.DestinationType `json:"destination_type"`
	DestinationName        *string                 `json:"destination_name"`
	DestinationCode        *string                 `json:"destination_code"`
	DestinationContact     *string                 `json:"destination_contact"`
	DestinationPhone       *string                 `json:"destination_phone"`
	DestinationAddress     *string                 `json:"destination_address"`

	LogisticsProviderID  *uint64 `json:"logistics_provider_id"`
	ShippingRateID       *uint64 `json:"shipping_rate_id"`
	TransportMode        *string `json:"transport_mode"`
	Carrier              *string `json:"carrier"`
	TrackingNumber       *string `json:"tracking_number"`
	ExpectedDeliveryDate *string `json:"expected_delivery_date"`

	BoxCount    *uint    `json:"box_count"`
	TotalWeight *float64 `json:"total_weight"`
	TotalVolume *float64 `json:"total_volume"`

	Items         []shipmentItemRequest `json:"items"`
	Remark        *string               `json:"remark"`
	InternalNotes *string               `json:"internal_notes"`
	OperatorID    *uint64               `json:"operator_id"`
}

// CreateShipment 创建发货单（打包完成）
func (h *ShipmentHandler) CreateShipment(c *gin.Context) {
	var req createShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]domain.CreateShipmentItemParams, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.CreateShipmentItemParams{
			ProductID:       item.ProductID,
			QuantityPlanned: item.QuantityPlanned,
			PackageSpecID:   item.PackageSpecID,
			BoxQuantity:     item.BoxQuantity,
			UnitCost:        item.UnitCost,
			Currency:        item.Currency,
			Remark:          item.Remark,
		})
	}

	params := &domain.CreateShipmentParams{
		OrderNumber:            req.OrderNumber,
		SalesChannel:           req.SalesChannel,
		WarehouseID:            req.WarehouseID,
		DestinationWarehouseID: req.DestinationWarehouseID,
		DestinationType:        req.DestinationType,
		DestinationName:        req.DestinationName,
		DestinationCode:        req.DestinationCode,
		DestinationContact:     req.DestinationContact,
		DestinationPhone:       req.DestinationPhone,
		DestinationAddress:     req.DestinationAddress,
		LogisticsProviderID:    req.LogisticsProviderID,
		ShippingRateID:         req.ShippingRateID,
		TransportMode:          req.TransportMode,
		Carrier:                req.Carrier,
		TrackingNumber:         req.TrackingNumber,
		ExpectedDeliveryDate:   req.ExpectedDeliveryDate,
		BoxCount:               req.BoxCount,
		TotalWeight:            req.TotalWeight,
		TotalVolume:            req.TotalVolume,
		Items:                  items,
		Remark:                 req.Remark,
		InternalNotes:          req.InternalNotes,
		OperatorID:             req.OperatorID,
	}

	shipment, err := h.usecase.Create(c, params)
	if err != nil {
		h.respondShipmentError(c, err)
		return
	}

	response.Success(c, shipment)
}

// UpdateShipment 编辑发货单
func (h *ShipmentHandler) UpdateShipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req updateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]domain.CreateShipmentItemParams, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.CreateShipmentItemParams{
			ProductID:       item.ProductID,
			QuantityPlanned: item.QuantityPlanned,
			PackageSpecID:   item.PackageSpecID,
			BoxQuantity:     item.BoxQuantity,
			UnitCost:        item.UnitCost,
			Currency:        item.Currency,
			Remark:          item.Remark,
		})
	}

	params := &domain.UpdateShipmentParams{
		OrderNumber:            req.OrderNumber,
		SalesChannel:           req.SalesChannel,
		WarehouseID:            req.WarehouseID,
		DestinationWarehouseID: req.DestinationWarehouseID,
		DestinationType:        req.DestinationType,
		DestinationName:        req.DestinationName,
		DestinationCode:        req.DestinationCode,
		DestinationContact:     req.DestinationContact,
		DestinationPhone:       req.DestinationPhone,
		DestinationAddress:     req.DestinationAddress,
		LogisticsProviderID:    req.LogisticsProviderID,
		ShippingRateID:         req.ShippingRateID,
		TransportMode:          req.TransportMode,
		Carrier:                req.Carrier,
		TrackingNumber:         req.TrackingNumber,
		ExpectedDeliveryDate:   req.ExpectedDeliveryDate,
		BoxCount:               req.BoxCount,
		TotalWeight:            req.TotalWeight,
		TotalVolume:            req.TotalVolume,
		Items:                  items,
		Remark:                 req.Remark,
		InternalNotes:          req.InternalNotes,
		OperatorID:             req.OperatorID,
	}

	shipment, err := h.usecase.Update(c, id, params)
	if err != nil {
		h.respondShipmentError(c, err)
		return
	}

	response.Success(c, shipment)
}

// ConfirmShipment 确认发货单 (DRAFT → CONFIRMED)
func (h *ShipmentHandler) ConfirmShipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		OperatorID *uint64 `json:"operator_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body
		req.OperatorID = nil
	}

	params := &domain.ConfirmShipmentParams{
		OperatorID: req.OperatorID,
	}

	if err := h.usecase.Confirm(c, id, params); err != nil {
		h.respondShipmentError(c, err)
		return
	}

	response.Success(c, nil)
}

type markShippedRequest struct {
	Carrier        *string  `json:"carrier"`
	TrackingNumber *string  `json:"tracking_number"`
	ShippingCost   *float64 `json:"shipping_cost"`
	Currency       *string  `json:"currency"`
	ShipDate       *string  `json:"ship_date"`
	Remark         *string  `json:"remark"`
	OperatorID     *uint64  `json:"operator_id"`
}

// MarkShipmentShipped 标记发货 (PACKED → SHIPPED)
func (h *ShipmentHandler) MarkShipmentShipped(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req markShippedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	params := &domain.MarkShippedParams{
		Carrier:        req.Carrier,
		TrackingNumber: req.TrackingNumber,
		ShippingCost:   req.ShippingCost,
		Currency:       req.Currency,
		ShipDate:       req.ShipDate,
		Remark:         req.Remark,
		OperatorID:     req.OperatorID,
	}

	if err := h.usecase.MarkShipped(c, id, params); err != nil {
		h.respondShipmentError(c, err)
		return
	}

	response.Success(c, nil)
}

type markDeliveredRequest struct {
	ActualDeliveryDate *string `json:"actual_delivery_date"`
	Remark             *string `json:"remark"`
	OperatorID         *uint64 `json:"operator_id"`
}

// MarkShipmentDelivered 标记送达 (SHIPPED → DELIVERED)
func (h *ShipmentHandler) MarkShipmentDelivered(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req markDeliveredRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	params := &domain.MarkDeliveredParams{
		ActualDeliveryDate: req.ActualDeliveryDate,
		Remark:             req.Remark,
		OperatorID:         req.OperatorID,
	}

	if err := h.usecase.MarkDelivered(c, id, params); err != nil {
		h.respondShipmentError(c, err)
		return
	}

	response.Success(c, nil)
}

type cancelShipmentRequest struct {
	Remark     *string `json:"remark"`
	OperatorID *uint64 `json:"operator_id"`
}

// CancelShipment 取消发货单（带回滚）
func (h *ShipmentHandler) CancelShipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req cancelShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body
		req.Remark = nil
		req.OperatorID = nil
	}

	params := &domain.CancelShipmentParams{
		Remark:     req.Remark,
		OperatorID: req.OperatorID,
	}

	if err := h.usecase.Cancel(c, id, params); err != nil {
		h.respondShipmentError(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteShipment 删除发货单
func (h *ShipmentHandler) DeleteShipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.usecase.Delete(c, id); err != nil {
		h.respondShipmentError(c, err)
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

func getUserIDFromContext(c *gin.Context) *uint64 {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint64); ok {
			return &id
		}
	}
	return nil
}
