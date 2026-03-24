package http

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/finance/domain"
	"am-erp-go/internal/module/finance/usecase"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FinanceHandler struct {
	cashLedgerUC      *usecase.CashLedgerUsecase
	costingSnapshotUC *usecase.CostingSnapshotUsecase
	dailyProfitUC     *usecase.DailyProfitUsecase
	profitQueryUC     *usecase.ProfitQueryUsecase
	productCostUC     *usecase.ProductCostUsecase
	exchangeRateUC    *usecase.ExchangeRateUsecase
	auditLogger       FinanceAuditLogger
}

type FinanceAuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

func NewFinanceHandler(
	cashLedgerUC *usecase.CashLedgerUsecase,
	costingSnapshotUC *usecase.CostingSnapshotUsecase,
	dailyProfitUC *usecase.DailyProfitUsecase,
	profitQueryUC *usecase.ProfitQueryUsecase,
	productCostUC *usecase.ProductCostUsecase,
	exchangeRateUC *usecase.ExchangeRateUsecase,
) *FinanceHandler {
	return &FinanceHandler{
		cashLedgerUC:      cashLedgerUC,
		costingSnapshotUC: costingSnapshotUC,
		dailyProfitUC:     dailyProfitUC,
		profitQueryUC:     profitQueryUC,
		productCostUC:     productCostUC,
		exchangeRateUC:    exchangeRateUC,
	}
}

func (h *FinanceHandler) BindAuditLogger(logger FinanceAuditLogger) {
	h.auditLogger = logger
}

// =========================
// Cash Ledger
// =========================

type createCashLedgerRequest struct {
	LedgerType    domain.LedgerType `json:"ledger_type" binding:"required"`
	Category      string            `json:"category" binding:"required"`
	Amount        float64           `json:"amount" binding:"required"`
	Currency      string            `json:"currency"`
	Marketplace   *string           `json:"marketplace"`
	OccurredNode  *string           `json:"occurred_node"`
	ReferenceType *string           `json:"reference_type"`
	ReferenceID   *uint64           `json:"reference_id"`
	Description   *string           `json:"description"`
	OccurredAt    *string           `json:"occurred_at"`
	OperatorID    *uint64           `json:"operator_id"`
}

type updateCashLedgerRequest struct {
	LedgerType    *domain.LedgerType `json:"ledger_type"`
	Category      *string            `json:"category"`
	Amount        *float64           `json:"amount"`
	Currency      *string            `json:"currency"`
	ReferenceType **string           `json:"reference_type"`
	ReferenceID   **uint64           `json:"reference_id"`
	Description   **string           `json:"description"`
	OccurredAt    *string            `json:"occurred_at"`
	OperatorID    *uint64            `json:"operator_id"`
}

type reverseCashLedgerRequest struct {
	OperatorID *uint64 `json:"operator_id"`
	Reason     string  `json:"reason"`
}

type rebuildDailyProfitRequest struct {
	BizDate     string  `json:"biz_date" binding:"required"`
	Marketplace *string `json:"marketplace"`
	BuilderID   *uint64 `json:"builder_id"`
}

type createExchangeRateRequest struct {
	FromCurrency string  `json:"from_currency" binding:"required"`
	ToCurrency   string  `json:"to_currency" binding:"required"`
	Rate         float64 `json:"rate" binding:"required"`
	EffectiveAt  *string `json:"effective_at"`
	Remark       *string `json:"remark"`
	OperatorID   *uint64 `json:"operator_id"`
}

type updateExchangeRateStatusRequest struct {
	Status     domain.ExchangeRateStatus `json:"status" binding:"required"`
	OperatorID *uint64                   `json:"operator_id"`
}

func (h *FinanceHandler) ListCashLedger(c *gin.Context) {
	params, err := buildCashLedgerListParams(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	list, total, err := h.cashLedgerUC.List(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

func (h *FinanceHandler) GetCashLedger(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	entry, err := h.cashLedgerUC.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "cash ledger not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, entry)
}

func (h *FinanceHandler) CreateCashLedger(c *gin.Context) {
	var req createCashLedgerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := resolveOperatorID(c, req.OperatorID)
	if operatorID == nil {
		response.Unauthorized(c, "operator_id is required")
		return
	}

	var occurredAt *time.Time
	if req.OccurredAt != nil && *req.OccurredAt != "" {
		parsed, err := parseDateTime(*req.OccurredAt)
		if err != nil {
			response.BadRequest(c, "invalid occurred_at")
			return
		}
		occurredAt = &parsed
	}

	created, err := h.cashLedgerUC.Create(&usecase.CreateCashLedgerInput{
		LedgerType:    req.LedgerType,
		Category:      req.Category,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Marketplace:   req.Marketplace,
		OccurredNode:  req.OccurredNode,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		Description:   req.Description,
		OccurredAt:    occurredAt,
		CreatedBy:     *operatorID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	h.recordCashLedgerAudit(c, "CREATE_CASH_LEDGER", created.ID, nil, buildCashLedgerAuditPayload(created))
	response.Success(c, created)
}

func (h *FinanceHandler) UpdateCashLedger(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req updateCashLedgerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var occurredAt *time.Time
	if req.OccurredAt != nil && *req.OccurredAt != "" {
		parsed, err := parseDateTime(*req.OccurredAt)
		if err != nil {
			response.BadRequest(c, "invalid occurred_at")
			return
		}
		occurredAt = &parsed
	}

	updated, err := h.cashLedgerUC.Update(id, &usecase.UpdateCashLedgerInput{
		LedgerType:    req.LedgerType,
		Category:      req.Category,
		Amount:        req.Amount,
		Currency:      req.Currency,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		Description:   req.Description,
		OccurredAt:    occurredAt,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrCashLedgerImmutable) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "cash ledger not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, updated)
}

func (h *FinanceHandler) DeleteCashLedger(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.cashLedgerUC.Delete(id); err != nil {
		if errors.Is(err, usecase.ErrCashLedgerImmutable) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FinanceHandler) ReverseCashLedger(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req reverseCashLedgerRequest
	_ = c.ShouldBindJSON(&req)

	operatorID := resolveOperatorID(c, req.OperatorID)
	if operatorID == nil {
		response.Unauthorized(c, "operator_id is required")
		return
	}

	beforeEntry, err := h.cashLedgerUC.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "cash ledger not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	entry, err := h.cashLedgerUC.Reverse(id, *operatorID, req.Reason)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "cash ledger not found")
			return
		}
		if errors.Is(err, usecase.ErrCashLedgerAlreadyReversed) {
			response.BadRequest(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	afterEntry, err := h.cashLedgerUC.Get(id)
	if err == nil {
		afterPayload := buildCashLedgerAuditPayload(afterEntry)
		afterPayload["reversal_reason"] = strings.TrimSpace(req.Reason)
		afterPayload["reversal_entry_id"] = entry.ID
		h.recordCashLedgerAudit(c, "REVERSE_CASH_LEDGER", id, buildCashLedgerAuditPayload(beforeEntry), afterPayload)
	}
	response.Success(c, entry)
}

func (h *FinanceHandler) GetCashLedgerSummary(c *gin.Context) {
	params, err := buildCashLedgerListParams(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	summary, err := h.cashLedgerUC.GetSummary(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

func (h *FinanceHandler) GetCashLedgerSummaryByCategory(c *gin.Context) {
	params, err := buildCashLedgerListParams(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	summary, err := h.cashLedgerUC.GetSummaryByCategory(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

func (h *FinanceHandler) RebuildDailyProfit(c *gin.Context) {
	if h.dailyProfitUC == nil {
		response.InternalError(c, "daily profit usecase not initialized")
		return
	}

	var req rebuildDailyProfitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	bizDatePtr, err := parseDateStart(req.BizDate)
	if err != nil || bizDatePtr == nil {
		response.BadRequest(c, "invalid biz_date")
		return
	}
	builderID := resolveOperatorID(c, req.BuilderID)

	items, err := h.dailyProfitUC.Rebuild(&usecase.RebuildDailyProfitInput{
		BizDate:     *bizDatePtr,
		Marketplace: req.Marketplace,
		BuilderID:   builderID,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *FinanceHandler) ListExchangeRates(c *gin.Context) {
	params := &domain.ExchangeRateListParams{
		Page:         parseIntOrDefault(c.Query("page"), 1),
		PageSize:     parseIntOrDefault(c.Query("page_size"), 20),
		FromCurrency: c.Query("from_currency"),
		ToCurrency:   c.Query("to_currency"),
	}
	if status := strings.TrimSpace(c.Query("status")); status != "" {
		params.Status = domain.ExchangeRateStatus(strings.ToUpper(status))
	}
	list, total, err := h.exchangeRateUC.List(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

func (h *FinanceHandler) CreateExchangeRate(c *gin.Context) {
	var req createExchangeRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	operatorID := resolveOperatorID(c, req.OperatorID)
	if operatorID == nil {
		response.Unauthorized(c, "operator_id is required")
		return
	}
	var effectiveAt *time.Time
	if req.EffectiveAt != nil && strings.TrimSpace(*req.EffectiveAt) != "" {
		parsed, err := parseDateTime(*req.EffectiveAt)
		if err != nil {
			response.BadRequest(c, "invalid effective_at")
			return
		}
		effectiveAt = &parsed
	}
	item, err := h.exchangeRateUC.Create(c, &usecase.CreateExchangeRateInput{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         req.Rate,
		EffectiveAt:  effectiveAt,
		Remark:       req.Remark,
		CreatedBy:    *operatorID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *FinanceHandler) UpdateExchangeRateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req updateExchangeRateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	operatorID := resolveOperatorID(c, req.OperatorID)
	if operatorID == nil {
		response.Unauthorized(c, "operator_id is required")
		return
	}
	if err := h.exchangeRateUC.UpdateStatus(c, id, &usecase.UpdateExchangeRateStatusInput{
		Status:     req.Status,
		OperatorID: *operatorID,
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "exchange rate not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FinanceHandler) GetProfitDashboard(c *gin.Context) {
	if h.dailyProfitUC == nil {
		response.InternalError(c, "daily profit usecase not initialized")
		return
	}

	dateFrom, err := parseDateStart(c.Query("date_from"))
	if err != nil {
		response.BadRequest(c, "invalid date_from")
		return
	}
	dateTo, err := parseDateEnd(c.Query("date_to"))
	if err != nil {
		response.BadRequest(c, "invalid date_to")
		return
	}
	input := &usecase.ProfitDashboardInput{
		Marketplace: c.Query("marketplace"),
	}
	if dateFrom != nil {
		input.DateFrom = *dateFrom
	}
	if dateTo != nil {
		input.DateTo = *dateTo
	}

	dashboard, err := h.dailyProfitUC.Dashboard(input)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, dashboard)
}

func (h *FinanceHandler) ListOrderProfits(c *gin.Context) {
	if h.profitQueryUC == nil {
		response.InternalError(c, "profit query usecase not initialized")
		return
	}

	dateFrom, err := parseDateStart(c.Query("date_from"))
	if err != nil {
		response.BadRequest(c, "invalid date_from")
		return
	}
	dateTo, err := parseDateEnd(c.Query("date_to"))
	if err != nil {
		response.BadRequest(c, "invalid date_to")
		return
	}

	params := &domain.OrderProfitListParams{
		Page:        parseIntOrDefault(c.Query("page"), 1),
		PageSize:    parseIntOrDefault(c.Query("page_size"), 20),
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		Marketplace: strings.TrimSpace(c.Query("marketplace")),
		Keyword:     strings.TrimSpace(c.Query("keyword")),
	}
	items, total, err := h.profitQueryUC.ListOrderProfits(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, params.Page, params.PageSize)
}

func (h *FinanceHandler) GetOrderProfitDetail(c *gin.Context) {
	if h.profitQueryUC == nil {
		response.InternalError(c, "profit query usecase not initialized")
		return
	}
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	detail, err := h.profitQueryUC.GetOrderProfitDetail(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "order profit not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, detail)
}

func (h *FinanceHandler) ListProductCostLedger(c *gin.Context) {
	if h.productCostUC == nil {
		response.InternalError(c, "product cost usecase not initialized")
		return
	}

	params := &domain.ProductCostLedgerListParams{
		Page:        parseIntOrDefault(c.Query("page"), 1),
		PageSize:    parseIntOrDefault(c.Query("page_size"), 20),
		Marketplace: strings.TrimSpace(c.Query("marketplace")),
	}

	if productID := strings.TrimSpace(c.Query("product_id")); productID != "" {
		parsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			response.BadRequest(c, "invalid product_id")
			return
		}
		params.ProductID = &parsed
	}
	if warehouseID := strings.TrimSpace(c.Query("warehouse_id")); warehouseID != "" {
		parsed, err := strconv.ParseUint(warehouseID, 10, 64)
		if err != nil {
			response.BadRequest(c, "invalid warehouse_id")
			return
		}
		params.WarehouseID = &parsed
	}
	if dateFrom, err := parseDateStart(c.Query("date_from")); err != nil {
		response.BadRequest(c, "invalid date_from")
		return
	} else {
		params.DateFrom = dateFrom
	}
	if dateTo, err := parseDateEnd(c.Query("date_to")); err != nil {
		response.BadRequest(c, "invalid date_to")
		return
	} else {
		params.DateTo = dateTo
	}

	items, total, err := h.productCostUC.ListLedger(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, params.Page, params.PageSize)
}

func (h *FinanceHandler) GetProductCostSummary(c *gin.Context) {
	if h.productCostUC == nil {
		response.InternalError(c, "product cost usecase not initialized")
		return
	}

	params := &domain.ProductCostLedgerListParams{
		Marketplace: strings.TrimSpace(c.Query("marketplace")),
	}
	if productID := strings.TrimSpace(c.Query("product_id")); productID != "" {
		parsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			response.BadRequest(c, "invalid product_id")
			return
		}
		params.ProductID = &parsed
	}
	if warehouseID := strings.TrimSpace(c.Query("warehouse_id")); warehouseID != "" {
		parsed, err := strconv.ParseUint(warehouseID, 10, 64)
		if err != nil {
			response.BadRequest(c, "invalid warehouse_id")
			return
		}
		params.WarehouseID = &parsed
	}
	if dateFrom, err := parseDateStart(c.Query("date_from")); err != nil {
		response.BadRequest(c, "invalid date_from")
		return
	} else {
		params.DateFrom = dateFrom
	}
	if dateTo, err := parseDateEnd(c.Query("date_to")); err != nil {
		response.BadRequest(c, "invalid date_to")
		return
	} else {
		params.DateTo = dateTo
	}

	summary, err := h.productCostUC.GetSummary(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

// =========================
// Costing Snapshot
// =========================

type createCostingSnapshotRequest struct {
	ProductID     uint64          `json:"product_id" binding:"required"`
	CostType      domain.CostType `json:"cost_type" binding:"required"`
	UnitCost      float64         `json:"unit_cost" binding:"required"`
	Currency      string          `json:"currency"`
	EffectiveFrom *string         `json:"effective_from"`
	EffectiveTo   *string         `json:"effective_to"`
	Notes         *string         `json:"notes"`
	OperatorID    *uint64         `json:"operator_id"`
}

type updateCostingSnapshotRequest struct {
	UnitCost      *float64 `json:"unit_cost"`
	Currency      *string  `json:"currency"`
	EffectiveFrom *string  `json:"effective_from"`
	EffectiveTo   **string `json:"effective_to"`
	Notes         **string `json:"notes"`
	OperatorID    *uint64  `json:"operator_id"`
}

func (h *FinanceHandler) ListCostingSnapshots(c *gin.Context) {
	params := &domain.CostingSnapshotListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	if productID := c.Query("product_id"); productID != "" {
		parsed, err := strconv.ParseUint(productID, 10, 64)
		if err != nil {
			response.BadRequest(c, "invalid product_id")
			return
		}
		params.ProductID = &parsed
	}
	if costType := c.Query("cost_type"); costType != "" {
		params.CostType = domain.CostType(strings.ToUpper(costType))
	}
	if isCurrent := c.Query("is_current"); isCurrent != "" {
		parsed, err := strconv.ParseBool(isCurrent)
		if err != nil {
			response.BadRequest(c, "invalid is_current")
			return
		}
		params.IsCurrent = &parsed
	}

	list, total, err := h.costingSnapshotUC.List(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

func (h *FinanceHandler) GetCostingSnapshot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	snapshot, err := h.costingSnapshotUC.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "costing snapshot not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, snapshot)
}

func (h *FinanceHandler) CreateCostingSnapshot(c *gin.Context) {
	var req createCostingSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := resolveOperatorID(c, req.OperatorID)
	if operatorID == nil {
		response.Unauthorized(c, "operator_id is required")
		return
	}

	var effectiveFrom *time.Time
	if req.EffectiveFrom != nil && *req.EffectiveFrom != "" {
		parsed, err := parseDateTime(*req.EffectiveFrom)
		if err != nil {
			response.BadRequest(c, "invalid effective_from")
			return
		}
		effectiveFrom = &parsed
	}

	var effectiveTo *time.Time
	if req.EffectiveTo != nil && *req.EffectiveTo != "" {
		parsed, err := parseDateTime(*req.EffectiveTo)
		if err != nil {
			response.BadRequest(c, "invalid effective_to")
			return
		}
		effectiveTo = &parsed
	}

	created, err := h.costingSnapshotUC.Create(&usecase.CreateCostingSnapshotInput{
		ProductID:     req.ProductID,
		CostType:      req.CostType,
		UnitCost:      req.UnitCost,
		Currency:      req.Currency,
		EffectiveFrom: effectiveFrom,
		EffectiveTo:   effectiveTo,
		Notes:         req.Notes,
		CreatedBy:     *operatorID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, created)
}

func (h *FinanceHandler) UpdateCostingSnapshot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req updateCostingSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var effectiveFrom *time.Time
	if req.EffectiveFrom != nil && *req.EffectiveFrom != "" {
		parsed, err := parseDateTime(*req.EffectiveFrom)
		if err != nil {
			response.BadRequest(c, "invalid effective_from")
			return
		}
		effectiveFrom = &parsed
	}

	var effectiveTo **time.Time
	if req.EffectiveTo != nil {
		if *req.EffectiveTo == nil {
			tmp := (*time.Time)(nil)
			effectiveTo = &tmp
		} else if **req.EffectiveTo == "" {
			tmp := (*time.Time)(nil)
			effectiveTo = &tmp
		} else {
			parsed, err := parseDateTime(**req.EffectiveTo)
			if err != nil {
				response.BadRequest(c, "invalid effective_to")
				return
			}
			tmp := &parsed
			effectiveTo = &tmp
		}
	}

	updated, err := h.costingSnapshotUC.Update(id, &usecase.UpdateCostingSnapshotInput{
		UnitCost:      req.UnitCost,
		Currency:      req.Currency,
		EffectiveFrom: effectiveFrom,
		EffectiveTo:   effectiveTo,
		Notes:         req.Notes,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "costing snapshot not found")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, updated)
}

func (h *FinanceHandler) DeleteCostingSnapshot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if _, err := h.costingSnapshotUC.Get(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "costing snapshot not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	if err := h.costingSnapshotUC.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *FinanceHandler) GetCurrentCost(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product_id")
		return
	}

	costType := domain.CostType(strings.ToUpper(c.Query("cost_type")))
	if costType == "" {
		response.BadRequest(c, "cost_type is required")
		return
	}

	snapshot, err := h.costingSnapshotUC.GetCurrent(productID, costType)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, snapshot)
}

func (h *FinanceHandler) GetAllCurrentCosts(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product_id")
		return
	}

	list, err := h.costingSnapshotUC.GetAllCurrent(productID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, list)
}

// =========================
// Helpers
// =========================

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}

func parseDateStart(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if len(raw) == 10 {
		t, err := time.ParseInLocation("2006-01-02", raw, time.Local)
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
	t, err := parseDateTime(raw)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseDateEnd(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if len(raw) == 10 {
		t, err := time.ParseInLocation("2006-01-02", raw, time.Local)
		if err != nil {
			return nil, err
		}
		end := t.Add(24*time.Hour - time.Second)
		return &end, nil
	}
	t, err := parseDateTime(raw)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseDateTime(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		time.RFC3339Nano,
	}
	var lastErr error
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, raw, time.Local)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func resolveOperatorID(c *gin.Context, operatorID *uint64) *uint64 {
	if operatorID != nil && *operatorID > 0 {
		return operatorID
	}
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(uint64); ok && id > 0 {
			return &id
		}
	}
	return nil
}

func buildCashLedgerListParams(c *gin.Context) (*domain.CashLedgerListParams, error) {
	params := &domain.CashLedgerListParams{
		Page:         parseIntOrDefault(c.Query("page"), 1),
		PageSize:     parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:      c.Query("keyword"),
		Category:     c.Query("category"),
		Marketplace:  strings.TrimSpace(c.Query("marketplace")),
		OccurredNode: strings.TrimSpace(c.Query("occurred_node")),
	}

	if ledgerType := c.Query("ledger_type"); ledgerType != "" {
		params.LedgerType = domain.LedgerType(strings.ToUpper(ledgerType))
	}
	if referenceType := c.Query("reference_type"); referenceType != "" {
		params.ReferenceType = referenceType
	}
	if referenceID := c.Query("reference_id"); referenceID != "" {
		id, err := strconv.ParseUint(referenceID, 10, 64)
		if err != nil {
			return nil, errors.New("invalid reference_id")
		}
		params.ReferenceID = &id
	}
	if dateFrom, err := parseDateStart(c.Query("date_from")); err != nil {
		return nil, errors.New("invalid date_from")
	} else {
		params.DateFrom = dateFrom
	}
	if dateTo, err := parseDateEnd(c.Query("date_to")); err != nil {
		return nil, errors.New("invalid date_to")
	} else {
		params.DateTo = dateTo
	}

	return params, nil
}

func (h *FinanceHandler) recordCashLedgerAudit(c *gin.Context, action string, ledgerID uint64, before, after any) {
	if h == nil || h.auditLogger == nil || c == nil || ledgerID == 0 {
		return
	}
	_ = h.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Finance",
		Action:     action,
		EntityType: "CashLedger",
		EntityID:   strconv.FormatUint(ledgerID, 10),
		Before:     before,
		After:      after,
	})
}

func buildCashLedgerAuditPayload(entry *domain.CashLedger) map[string]any {
	if entry == nil {
		return nil
	}
	payload := map[string]any{
		"ledger_type":       entry.LedgerType,
		"status":            entry.Status,
		"category":          entry.Category,
		"original_currency": entry.OriginalCurrency,
		"original_amount":   entry.OriginalAmount,
		"base_currency":     entry.BaseCurrency,
		"base_amount":       entry.BaseAmount,
		"fx_rate":           entry.FxRate,
		"occurred_at":       entry.OccurredAt,
	}
	if entry.Marketplace != nil {
		payload["marketplace"] = *entry.Marketplace
	}
	if entry.OccurredNode != nil {
		payload["occurred_node"] = *entry.OccurredNode
	}
	if entry.ReferenceType != nil {
		payload["reference_type"] = *entry.ReferenceType
	}
	if entry.ReferenceID != nil {
		payload["reference_id"] = *entry.ReferenceID
	}
	if entry.Description != nil {
		payload["description"] = *entry.Description
	}
	if entry.ReversalOfID != nil {
		payload["reversal_of_id"] = *entry.ReversalOfID
	}
	return payload
}
