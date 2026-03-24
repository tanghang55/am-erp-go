package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/numbering"
	"am-erp-go/internal/module/packaging/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

var (
	ErrPackagingNoPlans          = errors.New("no packaging plans to convert")
	ErrPackagingPurchaseNotFound = errors.New("packaging purchase order not found")
)

type PackagingProcurementUsecase struct {
	repo             domain.PackagingProcurementRepository
	itemRepo         domain.PackagingItemRepository
	ledgerRepo       domain.PackagingLedgerRepository
	auditLogger      AuditLogger
	defaultsProvider PackagingDefaultsProvider
	convertTxManager PackagingPlanConvertTransactionManager
	receiveTxManager PackagingPurchaseReceiveTransactionManager
}

type PackagingDefaultsProvider interface {
	GetDefaultBaseCurrency() string
}

type AuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

type PackagingPlanConvertTransactionalDeps struct {
	Repo domain.PackagingProcurementRepository
}

type PackagingPlanConvertTransactionManager interface {
	Run(ctx context.Context, fn func(PackagingPlanConvertTransactionalDeps) error) error
}

type PackagingPurchaseReceiveTransactionalDeps struct {
	Repo       domain.PackagingProcurementRepository
	ItemRepo   domain.PackagingItemRepository
	LedgerRepo domain.PackagingLedgerRepository
}

type PackagingPurchaseReceiveTransactionManager interface {
	Run(ctx context.Context, fn func(PackagingPurchaseReceiveTransactionalDeps) error) error
}

func NewPackagingProcurementUsecase(
	repo domain.PackagingProcurementRepository,
	itemRepo domain.PackagingItemRepository,
	ledgerRepo domain.PackagingLedgerRepository,
) *PackagingProcurementUsecase {
	return &PackagingProcurementUsecase{
		repo:       repo,
		itemRepo:   itemRepo,
		ledgerRepo: ledgerRepo,
	}
}

func (uc *PackagingProcurementUsecase) BindAuditLogger(logger AuditLogger) {
	uc.auditLogger = logger
}

func (uc *PackagingProcurementUsecase) BindDefaultsProvider(provider PackagingDefaultsProvider) {
	uc.defaultsProvider = provider
}

func (uc *PackagingProcurementUsecase) BindConvertTransactionManager(manager PackagingPlanConvertTransactionManager) {
	uc.convertTxManager = manager
}

func (uc *PackagingProcurementUsecase) BindReceiveTransactionManager(manager PackagingPurchaseReceiveTransactionManager) {
	uc.receiveTxManager = manager
}

func (uc *PackagingProcurementUsecase) ListPlans(params *domain.PackagingProcurementPlanListParams) ([]domain.PackagingProcurementPlan, int64, error) {
	if params == nil {
		params = &domain.PackagingProcurementPlanListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 50
	}
	if params.Date == nil {
		today := beginOfDay(time.Now())
		params.Date = &today
	}
	return uc.repo.ListPlans(params)
}

func (uc *PackagingProcurementUsecase) ListRuns(params *domain.PackagingProcurementRunListParams) ([]domain.PackagingProcurementRun, int64, error) {
	if params == nil {
		params = &domain.PackagingProcurementRunListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.repo.ListRuns(params)
}

func (uc *PackagingProcurementUsecase) GenerateDailyPlans(c *gin.Context, date *time.Time) (plans []domain.PackagingProcurementPlan, generatedCount int, err error) {
	planDate := beginOfDay(time.Now())
	if date != nil {
		planDate = beginOfDay(*date)
	}

	startedAt := time.Now()
	run := &domain.PackagingProcurementRun{
		RunNo:       ensurePackagingRunNo("PKRUN", startedAt),
		TriggerType: domain.PackagingProcurementTriggerManual,
		Status:      domain.PackagingProcurementRunRunning,
		StartedAt:   &startedAt,
	}
	run.InputSummary = marshalJSONPtr(map[string]any{
		"plan_date": planDate.Format("2006-01-02"),
	})
	if err = uc.repo.CreateRun(run); err != nil {
		return nil, 0, err
	}
	defer func() {
		finishedAt := time.Now()
		run.FinishedAt = &finishedAt
		run.OutputSummary = marshalJSONPtr(map[string]any{
			"generated_count": generatedCount,
			"current_count":   len(plans),
		})
		if err != nil {
			run.Status = domain.PackagingProcurementRunFailed
			msg := err.Error()
			run.ErrorMessage = &msg
		} else {
			run.Status = domain.PackagingProcurementRunSuccess
			run.ErrorMessage = nil
		}
		_ = uc.repo.UpdateRun(run)
	}()

	if err := uc.repo.CleanupPlansBefore(planDate); err != nil {
		return nil, 0, err
	}

	productDemands, err := uc.repo.LoadOrderedProductDemands(planDate)
	if err != nil {
		return nil, 0, err
	}
	if len(productDemands) == 0 {
		plans, generatedCount, err := uc.repo.SyncDailyPlans(planDate, []domain.PackagingPlanInput{})
		if err == nil {
			uc.recordAudit(c, "GENERATE_PACKAGING_PLANS", "PackagingProcurementPlan", planDate.Format("2006-01-02"), nil, map[string]any{"generated_count": generatedCount, "current_count": len(plans)})
		}
		return plans, generatedCount, err
	}

	productIDs := make([]uint64, 0, len(productDemands))
	demandByProduct := make(map[uint64]uint64, len(productDemands))
	for _, item := range productDemands {
		if item.ProductID == 0 || item.Qty == 0 {
			continue
		}
		demandByProduct[item.ProductID] = item.Qty
		productIDs = append(productIDs, item.ProductID)
	}
	if len(productIDs) == 0 {
		plans, generatedCount, err := uc.repo.SyncDailyPlans(planDate, []domain.PackagingPlanInput{})
		if err == nil {
			uc.recordAudit(c, "GENERATE_PACKAGING_PLANS", "PackagingProcurementPlan", planDate.Format("2006-01-02"), nil, map[string]any{"generated_count": generatedCount, "current_count": len(plans)})
		}
		return plans, generatedCount, err
	}

	mappings, err := uc.repo.LoadProductPackagingMappings(productIDs)
	if err != nil {
		return nil, 0, err
	}
	if len(mappings) == 0 {
		plans, generatedCount, err := uc.repo.SyncDailyPlans(planDate, []domain.PackagingPlanInput{})
		if err == nil {
			uc.recordAudit(c, "GENERATE_PACKAGING_PLANS", "PackagingProcurementPlan", planDate.Format("2006-01-02"), nil, map[string]any{"generated_count": generatedCount, "current_count": len(plans)})
		}
		return plans, generatedCount, err
	}

	requiredByItem := map[uint64]float64{}
	sourceByItem := map[uint64]map[uint64]uint64{}
	for _, mapping := range mappings {
		productDemand := demandByProduct[mapping.ProductID]
		if productDemand == 0 || mapping.PackagingItemID == 0 || mapping.QuantityPerUnit <= 0 {
			continue
		}
		requiredByItem[mapping.PackagingItemID] += float64(productDemand) * mapping.QuantityPerUnit
		if _, ok := sourceByItem[mapping.PackagingItemID]; !ok {
			sourceByItem[mapping.PackagingItemID] = map[uint64]uint64{}
		}
		sourceByItem[mapping.PackagingItemID][mapping.ProductID] += productDemand
	}

	if len(requiredByItem) == 0 {
		plans, generatedCount, err := uc.repo.SyncDailyPlans(planDate, []domain.PackagingPlanInput{})
		if err == nil {
			uc.recordAudit(c, "GENERATE_PACKAGING_PLANS", "PackagingProcurementPlan", planDate.Format("2006-01-02"), nil, map[string]any{"generated_count": generatedCount, "current_count": len(plans)})
		}
		return plans, generatedCount, err
	}

	itemIDs := make([]uint64, 0, len(requiredByItem))
	for itemID := range requiredByItem {
		itemIDs = append(itemIDs, itemID)
	}
	snapshots, err := uc.repo.LoadPackagingItemSnapshots(itemIDs)
	if err != nil {
		return nil, 0, err
	}

	sort.Slice(itemIDs, func(i, j int) bool {
		return itemIDs[i] < itemIDs[j]
	})

	inputs := make([]domain.PackagingPlanInput, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		snapshot, ok := snapshots[itemID]
		if !ok {
			continue
		}

		requiredQty := uint64(math.Ceil(requiredByItem[itemID]))
		if requiredQty == 0 {
			continue
		}

		shortageQty := uint64(0)
		if requiredQty > snapshot.QuantityOnHand {
			shortageQty = requiredQty - snapshot.QuantityOnHand
		}
		if shortageQty == 0 {
			continue
		}

		suggestedQty := shortageQty
		if snapshot.ReorderQuantity != nil && *snapshot.ReorderQuantity > suggestedQty {
			suggestedQty = *snapshot.ReorderQuantity
		}

		sources := sourceByItem[itemID]
		sourceJSON := marshalJSONPtr(sources)

		var remark *string
		if strings.ToUpper(strings.TrimSpace(snapshot.Status)) != "ACTIVE" {
			msg := "PACKAGING_ITEM_INACTIVE"
			remark = &msg
		}

		inputs = append(inputs, domain.PackagingPlanInput{
			PackagingItemID: itemID,
			RequiredQty:     requiredQty,
			OnHandQty:       snapshot.QuantityOnHand,
			ShortageQty:     shortageQty,
			SuggestedQty:    suggestedQty,
			SourceJSON:      sourceJSON,
			Remark:          remark,
		})
	}

	plans, generatedCount, err = uc.repo.SyncDailyPlans(planDate, inputs)
	if err == nil {
		uc.recordAudit(c, "GENERATE_PACKAGING_PLANS", "PackagingProcurementPlan", planDate.Format("2006-01-02"), nil, map[string]any{"generated_count": generatedCount, "current_count": len(plans)})
	}
	return plans, generatedCount, err
}

func (uc *PackagingProcurementUsecase) ConvertPlans(c *gin.Context, params *domain.PackagingPlanConvertParams) (*domain.PackagingPurchaseOrder, error) {
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if uc.convertTxManager != nil {
		var result *domain.PackagingPurchaseOrder
		if err := uc.convertTxManager.Run(ctx, func(deps PackagingPlanConvertTransactionalDeps) error {
			order, err := uc.convertPlansWithDeps(c, params, deps)
			if err != nil {
				return err
			}
			result = order
			return nil
		}); err != nil {
			return nil, err
		}
		return result, nil
	}
	return uc.convertPlansWithDeps(c, params, PackagingPlanConvertTransactionalDeps{Repo: uc.repo})
}

func (uc *PackagingProcurementUsecase) convertPlansWithDeps(
	c *gin.Context,
	params *domain.PackagingPlanConvertParams,
	deps PackagingPlanConvertTransactionalDeps,
) (*domain.PackagingPurchaseOrder, error) {
	if params == nil {
		params = &domain.PackagingPlanConvertParams{}
	}
	if params.Date == nil {
		today := beginOfDay(time.Now())
		params.Date = &today
	}

	execUC := *uc
	execUC.repo = deps.Repo

	plans, err := execUC.repo.ListConvertiblePlans(params)
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, ErrPackagingNoPlans
	}

	itemIDs := make([]uint64, 0, len(plans))
	for _, plan := range plans {
		itemIDs = append(itemIDs, plan.PackagingItemID)
	}
	snapshots, err := execUC.repo.LoadPackagingItemSnapshots(itemIDs)
	if err != nil {
		return nil, err
	}

	items := make([]domain.PackagingPurchaseOrderItem, 0, len(plans))
	planIDs := make([]uint64, 0, len(plans))
	total := 0.0
	currency := uc.defaultCurrency("")

	for _, plan := range plans {
		snapshot, ok := snapshots[plan.PackagingItemID]
		if !ok || plan.SuggestedQty == 0 {
			continue
		}
		itemCurrency := uc.defaultCurrency(snapshot.Currency)
		subtotal := roundAmount(snapshot.UnitCost * float64(plan.SuggestedQty))
		total += subtotal
		items = append(items, domain.PackagingPurchaseOrderItem{
			PackagingItemID: plan.PackagingItemID,
			QtyOrdered:      plan.SuggestedQty,
			QtyReceived:     0,
			UnitCost:        snapshot.UnitCost,
			Currency:        itemCurrency,
			Subtotal:        subtotal,
		})
		currency = itemCurrency
		planIDs = append(planIDs, plan.ID)
	}

	if len(items) == 0 {
		return nil, ErrPackagingNoPlans
	}

	order := &domain.PackagingPurchaseOrder{
		PoNumber:    ensurePackagingPoNumber("PKPO"),
		Status:      domain.PackagingPurchaseOrderDraft,
		Currency:    currency,
		TotalAmount: roundAmount(total),
		Remark:      fmt.Sprintf("AUTO_FROM_PACKAGING_PLAN:%s", params.Date.Format("2006-01-02")),
		CreatedBy:   params.OperatorID,
		UpdatedBy:   params.OperatorID,
		Items:       items,
	}

	if err := execUC.repo.CreatePurchaseOrder(order); err != nil {
		return nil, err
	}
	if err := execUC.repo.MarkPlansConverted(planIDs, order.ID); err != nil {
		return nil, err
	}
	uc.recordAudit(c, "CONVERT_PACKAGING_PLANS", "PackagingProcurementPlan", params.Date.Format("2006-01-02"), nil, map[string]any{
		"plan_ids":          planIDs,
		"purchase_order_id": order.ID,
	})
	return order, nil
}

func (uc *PackagingProcurementUsecase) ListPurchaseOrders(params *domain.PackagingPurchaseOrderListParams) ([]domain.PackagingPurchaseOrder, int64, error) {
	if params == nil {
		params = &domain.PackagingPurchaseOrderListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.repo.ListPurchaseOrders(params)
}

func (uc *PackagingProcurementUsecase) GetPurchaseOrder(id uint64) (*domain.PackagingPurchaseOrder, error) {
	return uc.repo.GetPurchaseOrder(id)
}

func (uc *PackagingProcurementUsecase) SubmitPurchaseOrder(c *gin.Context, id uint64, operatorID *uint64) (*domain.PackagingPurchaseOrder, error) {
	order, err := uc.repo.GetPurchaseOrder(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrPackagingPurchaseNotFound
	}
	if order.Status != domain.PackagingPurchaseOrderDraft {
		return nil, errors.New("only draft order can be submitted")
	}

	now := time.Now()
	order.Status = domain.PackagingPurchaseOrderOrdered
	order.OrderedAt = &now
	order.UpdatedBy = operatorID
	order.Items = nil
	if err := uc.repo.UpdatePurchaseOrder(order); err != nil {
		return nil, err
	}
	updated, err := uc.repo.GetPurchaseOrder(id)
	if err == nil {
		uc.recordAudit(c, "SUBMIT", "PackagingPurchaseOrder", fmt.Sprintf("%d", id), nil, updated)
	}
	return updated, err
}

func (uc *PackagingProcurementUsecase) ReceivePurchaseOrder(c *gin.Context, id uint64, params *domain.PackagingPurchaseOrderReceiveParams) (*domain.PackagingPurchaseOrder, error) {
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if uc.receiveTxManager != nil {
		var result *domain.PackagingPurchaseOrder
		if err := uc.receiveTxManager.Run(ctx, func(deps PackagingPurchaseReceiveTransactionalDeps) error {
			order, err := uc.receivePurchaseOrderWithDeps(c, id, params, deps)
			if err != nil {
				return err
			}
			result = order
			return nil
		}); err != nil {
			return nil, err
		}
		return result, nil
	}
	return uc.receivePurchaseOrderWithDeps(c, id, params, PackagingPurchaseReceiveTransactionalDeps{
		Repo:       uc.repo,
		ItemRepo:   uc.itemRepo,
		LedgerRepo: uc.ledgerRepo,
	})
}

func (uc *PackagingProcurementUsecase) receivePurchaseOrderWithDeps(
	c *gin.Context,
	id uint64,
	params *domain.PackagingPurchaseOrderReceiveParams,
	deps PackagingPurchaseReceiveTransactionalDeps,
) (*domain.PackagingPurchaseOrder, error) {
	if params == nil {
		params = &domain.PackagingPurchaseOrderReceiveParams{}
	}
	execUC := *uc
	execUC.repo = deps.Repo
	execUC.itemRepo = deps.ItemRepo
	execUC.ledgerRepo = deps.LedgerRepo

	order, err := execUC.repo.GetPurchaseOrder(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrPackagingPurchaseNotFound
	}
	if order.Status != domain.PackagingPurchaseOrderOrdered {
		return nil, errors.New("only ordered order can be received")
	}

	now := time.Now()
	allReceived := true
	for idx := range order.Items {
		item := &order.Items[idx]
		if item.QtyOrdered <= item.QtyReceived {
			continue
		}
		remaining := item.QtyOrdered - item.QtyReceived

		delta := remaining
		if len(params.ReceivedQties) > 0 {
			if value, ok := params.ReceivedQties[item.ID]; ok {
				delta = value
			} else {
				delta = 0
			}
		}
		if delta == 0 {
			allReceived = false
			continue
		}
		if delta > remaining {
			return nil, errors.New("received qty exceeds ordered")
		}

		packagingItem, err := execUC.itemRepo.GetByID(item.PackagingItemID)
		if err != nil {
			return nil, err
		}
		traceID := fmt.Sprintf("PKG-PO-%d-%d", now.Unix(), item.PackagingItemID)
		refType := "PACKAGING_PURCHASE_ORDER"
		notes := fmt.Sprintf("包材采购单入库: %s", order.PoNumber)
		createdBy := uint64(0)
		if params.OperatorID != nil {
			createdBy = *params.OperatorID
		}
		ledger := &domain.PackagingLedger{
			TraceID:         traceID,
			PackagingItemID: item.PackagingItemID,
			TransactionType: "IN",
			Quantity:        int64(delta),
			UnitCost:        item.UnitCost,
			QuantityBefore:  packagingItem.QuantityOnHand,
			QuantityAfter:   packagingItem.QuantityOnHand + delta,
			ReferenceType:   &refType,
			ReferenceID:     &order.ID,
			OccurredAt:      now,
			Notes:           &notes,
			CreatedBy:       createdBy,
		}
		if err := execUC.ledgerRepo.Create(ledger); err != nil {
			return nil, err
		}
		if err := execUC.itemRepo.UpdateQuantity(item.PackagingItemID, int64(delta)); err != nil {
			return nil, err
		}

		item.QtyReceived += delta
		if item.QtyOrdered > item.QtyReceived {
			allReceived = false
		}
	}

	if allReceived {
		order.Status = domain.PackagingPurchaseOrderReceived
		order.ReceivedAt = &now
	}
	order.UpdatedBy = params.OperatorID
	if err := execUC.repo.UpdatePurchaseOrder(order); err != nil {
		return nil, err
	}
	updated, err := execUC.repo.GetPurchaseOrder(id)
	if err == nil {
		uc.recordAudit(c, "RECEIVE", "PackagingPurchaseOrder", fmt.Sprintf("%d", id), nil, updated)
	}
	return updated, err
}

func beginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func roundAmount(v float64) float64 {
	return math.Round(v*10000) / 10000
}

func marshalJSONPtr(v any) *string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	value := string(bytes)
	return &value
}

func ensurePackagingPoNumber(prefix string) string {
	return numbering.Generate(prefix, time.Now())
}

func ensurePackagingRunNo(prefix string, now time.Time) string {
	return numbering.Generate(prefix, now)
}

func (uc *PackagingProcurementUsecase) recordAudit(c *gin.Context, action, entityType, entityID string, before, after any) {
	if uc.auditLogger == nil || c == nil {
		return
	}
	_ = uc.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Packaging",
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
}

func (uc *PackagingProcurementUsecase) defaultCurrency(currency string) string {
	if strings.TrimSpace(currency) != "" {
		return strings.TrimSpace(currency)
	}
	if uc != nil && uc.defaultsProvider != nil {
		if defaultCurrency := strings.TrimSpace(uc.defaultsProvider.GetDefaultBaseCurrency()); defaultCurrency != "" {
			return defaultCurrency
		}
	}
	return "USD"
}
