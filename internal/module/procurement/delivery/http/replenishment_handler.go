package http

import (
	"errors"
	"strconv"
	"time"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/procurement/domain"
	"am-erp-go/internal/module/procurement/usecase"

	"github.com/gin-gonic/gin"
)

type ReplenishmentHandler struct {
	usecase *usecase.ReplenishmentUsecase
}

func NewReplenishmentHandler(usecase *usecase.ReplenishmentUsecase) *ReplenishmentHandler {
	return &ReplenishmentHandler{usecase: usecase}
}

type updateConfigRequest struct {
	IsEnabled       uint8  `json:"is_enabled"`
	IntervalMinutes uint32 `json:"interval_minutes"`
}

func (h *ReplenishmentHandler) GetConfig(c *gin.Context) {
	cfg, err := h.usecase.GetConfig()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *ReplenishmentHandler) UpdateConfig(c *gin.Context) {
	var req updateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	cfg := &domain.ReplenishmentConfig{
		IsEnabled:       req.IsEnabled,
		IntervalMinutes: req.IntervalMinutes,
	}
	updated, err := h.usecase.UpdateConfig(c, cfg)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, updated)
}

type upsertStrategyRequest struct {
	ID                   uint64  `json:"id"`
	Name                 string  `json:"name"`
	Priority             uint32  `json:"priority"`
	IsEnabled            uint8   `json:"is_enabled"`
	ProductID            *uint64 `json:"product_id"`
	WarehouseID          *uint64 `json:"warehouse_id"`
	SupplierID           *uint64 `json:"supplier_id"`
	Marketplace          *string `json:"marketplace"`
	ConditionJSON        *string `json:"condition_json"`
	DemandWindowDays     uint32  `json:"demand_window_days"`
	ProcurementCycleDays uint32  `json:"procurement_cycle_days"`
	PackDays             uint32  `json:"pack_days"`
	LogisticsDays        uint32  `json:"logistics_days"`
	SafetyDays           uint32  `json:"safety_days"`
	ZeroSalesPurchaseQty uint32  `json:"zero_sales_purchase_qty"`
	MOQ                  uint32  `json:"moq"`
	OrderMultiple        uint32  `json:"order_multiple"`
	Remark               *string `json:"remark"`
	OperatorID           *uint64 `json:"operator_id"`
}

func (h *ReplenishmentHandler) ListStrategies(c *gin.Context) {
	params := &domain.ReplenishmentStrategyListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:  c.Query("keyword"),
	}
	list, total, err := h.usecase.ListStrategies(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

func (h *ReplenishmentHandler) UpsertStrategy(c *gin.Context) {
	var req upsertStrategyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	strategy := &domain.ReplenishmentStrategy{
		ID:                   req.ID,
		Name:                 req.Name,
		Priority:             req.Priority,
		IsEnabled:            req.IsEnabled,
		ProductID:            req.ProductID,
		WarehouseID:          req.WarehouseID,
		SupplierID:           req.SupplierID,
		Marketplace:          req.Marketplace,
		ConditionJSON:        req.ConditionJSON,
		DemandWindowDays:     req.DemandWindowDays,
		ProcurementCycleDays: req.ProcurementCycleDays,
		PackDays:             req.PackDays,
		LogisticsDays:        req.LogisticsDays,
		SafetyDays:           req.SafetyDays,
		ZeroSalesPurchaseQty: req.ZeroSalesPurchaseQty,
		MOQ:                  req.MOQ,
		OrderMultiple:        req.OrderMultiple,
		Remark:               req.Remark,
		CreatedBy:            req.OperatorID,
		UpdatedBy:            req.OperatorID,
	}
	updated, err := h.usecase.UpsertStrategy(c, strategy)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, updated)
}

func (h *ReplenishmentHandler) ListPlans(c *gin.Context) {
	params := &domain.ReplenishmentPlanListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 50),
	}
	if dateStr := c.Query("date"); dateStr != "" {
		d, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			params.Date = &d
		}
	}
	if status := c.Query("status"); status != "" {
		params.Status = domain.ReplenishmentPlanStatus(status)
	}
	list, total, err := h.usecase.ListPlans(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

func (h *ReplenishmentHandler) DeletePlan(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid plan id")
		return
	}
	if err := h.usecase.DeletePlanByID(c, planID); err != nil {
		if errors.Is(err, usecase.ErrReplenishmentPlanNotDeletable) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *ReplenishmentHandler) ListRuns(c *gin.Context) {
	params := &domain.ReplenishmentRunListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}
	if status := c.Query("status"); status != "" {
		params.Status = domain.ReplenishmentRunStatus(status)
	}
	if triggerType := c.Query("trigger_type"); triggerType != "" {
		params.Triggered = domain.ReplenishmentTriggerType(triggerType)
	}

	list, total, err := h.usecase.ListRuns(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

type generateRequest struct {
	OperatorID *uint64 `json:"operator_id"`
}

func (h *ReplenishmentHandler) GeneratePlans(c *gin.Context) {
	var req generateRequest
	_ = c.ShouldBindJSON(&req)

	plans, generatedCount, err := h.usecase.GenerateDailyPlans(c, &domain.ReplenishmentGenerateParams{
		TriggerType: domain.ReplenishmentTriggerManual,
		OperatorID:  req.OperatorID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"generated":       generatedCount > 0,
		"generated_count": generatedCount,
		"current_count":   len(plans),
		"plans":           plans,
	})
}

type convertPlansRequest struct {
	PlanIDs    []uint64 `json:"plan_ids"`
	Date       *string  `json:"date"`
	OperatorID *uint64  `json:"operator_id"`
}

func (h *ReplenishmentHandler) ConvertPlans(c *gin.Context) {
	var req convertPlansRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var planDate *time.Time
	if req.Date != nil && *req.Date != "" {
		if parsed, err := time.Parse("2006-01-02", *req.Date); err == nil {
			planDate = &parsed
		}
	}
	orders, err := h.usecase.ConvertPlansToPurchaseOrders(c, &domain.ReplenishmentPlanConvertParams{
		PlanIDs:    req.PlanIDs,
		PlanDate:   planDate,
		OperatorID: req.OperatorID,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrReplenishmentNoPlans) || errors.Is(err, usecase.ErrReplenishmentMissingSupplier) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"created_count": len(orders),
		"orders":        orders,
	})
}

func (h *ReplenishmentHandler) CleanupPlans(c *gin.Context) {
	if err := h.usecase.EnsureDailyCleanup(); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, nil)
}
