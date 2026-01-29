package http

import (
	"errors"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"data":  shipments,
			"total": total,
		},
		"success": true,
	})
}

// GetShipment 获取发货单详情
func (h *ShipmentHandler) GetShipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	shipment, err := h.usecase.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "shipment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    shipment,
		"success": true,
	})
}

type shipmentItemRequest struct {
	SkuID           uint64   `json:"sku_id"`
	QuantityPlanned uint     `json:"quantity_planned"`
	PackageSpecID   *uint64  `json:"package_spec_id"`
	BoxQuantity     *uint    `json:"box_quantity"`
	UnitCost        *float64 `json:"unit_cost"`
	Currency        *string  `json:"currency"`
	Remark          *string  `json:"remark"`
}

type createShipmentRequest struct {
	OrderNumber *string               `json:"order_number"`
	WarehouseID uint64                `json:"warehouse_id"`
	Items       []shipmentItemRequest `json:"items"`
	Remark      *string               `json:"remark"`
	OperatorID  *uint64               `json:"operator_id"`
}

// CreateShipment 创建发货单（打包完成）
func (h *ShipmentHandler) CreateShipment(c *gin.Context) {
	var req createShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	items := make([]domain.CreateShipmentItemParams, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.CreateShipmentItemParams{
			SkuID:           item.SkuID,
			QuantityPlanned: item.QuantityPlanned,
			PackageSpecID:   item.PackageSpecID,
			BoxQuantity:     item.BoxQuantity,
			UnitCost:        item.UnitCost,
			Currency:        item.Currency,
			Remark:          item.Remark,
		})
	}

	params := &domain.CreateShipmentParams{
		OrderNumber: req.OrderNumber,
		WarehouseID: req.WarehouseID,
		Items:       items,
		Remark:      req.Remark,
		OperatorID:  req.OperatorID,
	}

	shipment, err := h.usecase.Create(c, params)
	if err != nil {
		if errors.Is(err, usecase.ErrEmptyItems) {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    shipment,
		"success": true,
	})
}

// ConfirmShipment 确认发货单 (DRAFT → CONFIRMED)
func (h *ShipmentHandler) ConfirmShipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"success": true,
	})
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req markShippedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"success": true,
	})
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req markDeliveredRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	params := &domain.MarkDeliveredParams{
		ActualDeliveryDate: req.ActualDeliveryDate,
		Remark:             req.Remark,
		OperatorID:         req.OperatorID,
	}

	if err := h.usecase.MarkDelivered(c, id, params); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"success": true,
	})
}

type cancelShipmentRequest struct {
	Remark     *string `json:"remark"`
	OperatorID *uint64 `json:"operator_id"`
}

// CancelShipment 取消发货单（带回滚）
func (h *ShipmentHandler) CancelShipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"success": true,
	})
}

// DeleteShipment 删除发货单
func (h *ShipmentHandler) DeleteShipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	if err := h.usecase.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"success": true,
	})
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
