package http

import (
	"strconv"
	"time"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/packaging/domain"
	"am-erp-go/internal/module/packaging/usecase"

	"github.com/gin-gonic/gin"
)

type PackagingProcurementHandler struct {
	uc *usecase.PackagingProcurementUsecase
}

func NewPackagingProcurementHandler(uc *usecase.PackagingProcurementUsecase) *PackagingProcurementHandler {
	return &PackagingProcurementHandler{uc: uc}
}

func (h *PackagingProcurementHandler) ListPlans(c *gin.Context) {
	params := &domain.PackagingProcurementPlanListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 50),
	}
	if dateStr := c.Query("date"); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			params.Date = &parsed
		}
	}
	if status := c.Query("status"); status != "" {
		params.Status = domain.PackagingProcurementPlanStatus(status)
	}

	list, total, err := h.uc.ListPlans(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

func (h *PackagingProcurementHandler) ListRuns(c *gin.Context) {
	params := &domain.PackagingProcurementRunListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}
	if status := c.Query("status"); status != "" {
		params.Status = domain.PackagingProcurementRunStatus(status)
	}
	if triggerType := c.Query("trigger_type"); triggerType != "" {
		params.TriggerType = domain.PackagingProcurementTriggerType(triggerType)
	}
	list, total, err := h.uc.ListRuns(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

type generatePackagingPlanRequest struct {
	Date *string `json:"date"`
}

func (h *PackagingProcurementHandler) GeneratePlans(c *gin.Context) {
	var req generatePackagingPlanRequest
	_ = c.ShouldBindJSON(&req)

	var date *time.Time
	if req.Date != nil && *req.Date != "" {
		if parsed, err := time.Parse("2006-01-02", *req.Date); err == nil {
			date = &parsed
		}
	}

	plans, generatedCount, err := h.uc.GenerateDailyPlans(c, date)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"generated_count": generatedCount,
		"current_count":   len(plans),
		"plans":           plans,
	})
}

type convertPackagingPlanRequest struct {
	PlanIDs    []uint64 `json:"plan_ids"`
	Date       *string  `json:"date"`
	OperatorID *uint64  `json:"operator_id"`
}

func (h *PackagingProcurementHandler) ConvertPlans(c *gin.Context) {
	var req convertPackagingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.OperatorID == nil {
		req.OperatorID = getUserIDFromContext(c)
	}

	var date *time.Time
	if req.Date != nil && *req.Date != "" {
		if parsed, err := time.Parse("2006-01-02", *req.Date); err == nil {
			date = &parsed
		}
	}

	order, err := h.uc.ConvertPlans(c, &domain.PackagingPlanConvertParams{
		PlanIDs:    req.PlanIDs,
		Date:       date,
		OperatorID: req.OperatorID,
	})
	if err != nil {
		if err == usecase.ErrPackagingNoPlans {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *PackagingProcurementHandler) ListPurchaseOrders(c *gin.Context) {
	params := &domain.PackagingPurchaseOrderListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}
	if status := c.Query("status"); status != "" {
		params.Status = domain.PackagingPurchaseOrderStatus(status)
	}
	list, total, err := h.uc.ListPurchaseOrders(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

func (h *PackagingProcurementHandler) GetPurchaseOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}
	order, err := h.uc.GetPurchaseOrder(orderID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if order == nil {
		response.NotFound(c, "order not found")
		return
	}
	response.Success(c, order)
}

type submitPackagingOrderRequest struct {
	OperatorID *uint64 `json:"operator_id"`
}

func (h *PackagingProcurementHandler) SubmitPurchaseOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	var req submitPackagingOrderRequest
	_ = c.ShouldBindJSON(&req)
	if req.OperatorID == nil {
		req.OperatorID = getUserIDFromContext(c)
	}

	order, err := h.uc.SubmitPurchaseOrder(c, orderID, req.OperatorID)
	if err != nil {
		if err == usecase.ErrPackagingPurchaseNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, order)
}

type receivePackagingOrderRequest struct {
	ReceivedQties map[string]uint64 `json:"received_qties"`
	OperatorID    *uint64           `json:"operator_id"`
}

func (h *PackagingProcurementHandler) ReceivePurchaseOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid order id")
		return
	}

	var req receivePackagingOrderRequest
	_ = c.ShouldBindJSON(&req)
	if req.OperatorID == nil {
		req.OperatorID = getUserIDFromContext(c)
	}

	receivedQties := map[uint64]uint64{}
	for itemIDText, qty := range req.ReceivedQties {
		itemID, parseErr := strconv.ParseUint(itemIDText, 10, 64)
		if parseErr != nil {
			continue
		}
		receivedQties[itemID] = qty
	}

	order, err := h.uc.ReceivePurchaseOrder(c, orderID, &domain.PackagingPurchaseOrderReceiveParams{
		ReceivedQties: receivedQties,
		OperatorID:    req.OperatorID,
	})
	if err != nil {
		if err == usecase.ErrPackagingPurchaseNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, order)
}
