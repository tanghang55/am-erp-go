package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/numbering"
	"am-erp-go/internal/module/procurement/domain"
	systemdomain "am-erp-go/internal/module/system/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

var (
	ErrReplenishmentNoPlans          = errors.New("no replenishment plans to convert")
	ErrReplenishmentMissingPoUC      = errors.New("purchase order usecase not configured")
	ErrReplenishmentMissingSupplier  = errors.New("some replenishment plans have no supplier")
	ErrReplenishmentPlanNotDeletable = errors.New("only pending replenishment plan can be deleted")
)

type ReplenishmentUsecase struct {
	repo             domain.ReplenishmentRepository
	poUsecase        *PurchaseOrderUsecase
	defaultsProvider ProcurementDefaultsProvider
	auditLogger      AuditLogger
}

type ProcurementDefaultsProvider interface {
	GetProcurementDefaults() systemdomain.ConfigCenterProcurement
	GetDefaultBaseCurrency() string
}

func NewReplenishmentUsecase(repo domain.ReplenishmentRepository, poUsecase *PurchaseOrderUsecase) *ReplenishmentUsecase {
	return &ReplenishmentUsecase{
		repo:      repo,
		poUsecase: poUsecase,
	}
}

func (uc *ReplenishmentUsecase) BindAuditLogger(logger AuditLogger) {
	uc.auditLogger = logger
}

func (uc *ReplenishmentUsecase) BindDefaultsProvider(provider ProcurementDefaultsProvider) {
	uc.defaultsProvider = provider
}

func (uc *ReplenishmentUsecase) GetConfig() (*domain.ReplenishmentConfig, error) {
	cfg, err := uc.repo.GetConfig()
	if err != nil {
		return nil, err
	}
	if uc.defaultsProvider != nil {
		defaults := uc.defaultsProvider.GetProcurementDefaults()
		cfg.DemandWindowDays = defaults.DemandWindowDays
		cfg.DefaultLeadTimeDays = defaults.DefaultLeadTimeDays
		cfg.DefaultSafetyDays = defaults.DefaultSafetyDays
		cfg.DefaultMOQ = defaults.DefaultMOQ
		cfg.DefaultOrderMultiple = defaults.DefaultOrderMultiple
	}
	return cfg, nil
}

func (uc *ReplenishmentUsecase) UpdateConfig(c *gin.Context, input *domain.ReplenishmentConfig) (*domain.ReplenishmentConfig, error) {
	cfg, err := uc.repo.GetConfig()
	if err != nil {
		return nil, err
	}
	if input == nil {
		return cfg, nil
	}

	cfg.IsEnabled = input.IsEnabled
	cfg.IntervalMinutes = clampU32(input.IntervalMinutes, 1, 10080, 1440)

	if err := uc.repo.SaveConfig(cfg); err != nil {
		return nil, err
	}
	if uc.defaultsProvider != nil {
		defaults := uc.defaultsProvider.GetProcurementDefaults()
		cfg.DemandWindowDays = defaults.DemandWindowDays
		cfg.DefaultLeadTimeDays = defaults.DefaultLeadTimeDays
		cfg.DefaultSafetyDays = defaults.DefaultSafetyDays
		cfg.DefaultMOQ = defaults.DefaultMOQ
		cfg.DefaultOrderMultiple = defaults.DefaultOrderMultiple
	}
	uc.recordAudit(c, "UPDATE_CONFIG", "ReplenishmentConfig", "default", nil, cfg)
	return cfg, nil
}

func (uc *ReplenishmentUsecase) ListStrategies(params *domain.ReplenishmentStrategyListParams) ([]domain.ReplenishmentStrategy, int64, error) {
	if params == nil {
		params = &domain.ReplenishmentStrategyListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.repo.ListStrategies(params)
}

func (uc *ReplenishmentUsecase) UpsertStrategy(c *gin.Context, strategy *domain.ReplenishmentStrategy) (*domain.ReplenishmentStrategy, error) {
	if strategy == nil {
		return nil, errors.New("invalid strategy")
	}
	if strategy.Name == "" {
		strategy.Name = fmt.Sprintf("Strategy-%d", time.Now().Unix())
	}
	strategy.Priority = clampU32(strategy.Priority, 1, 1000000, 100)
	strategy.DemandWindowDays = clampU32(strategy.DemandWindowDays, 1, 365, 30)
	strategy.ProcurementCycleDays = clampU32(strategy.ProcurementCycleDays, 0, 365, 15)
	strategy.PackDays = clampU32(strategy.PackDays, 0, 365, 3)
	strategy.LogisticsDays = clampU32(strategy.LogisticsDays, 0, 365, 7)
	strategy.SafetyDays = clampU32(strategy.SafetyDays, 0, 365, 7)
	strategy.ZeroSalesPurchaseQty = clampU32(strategy.ZeroSalesPurchaseQty, 0, 1000000, 0)
	strategy.MOQ = clampU32(strategy.MOQ, 1, 1000000, 1)
	strategy.OrderMultiple = clampU32(strategy.OrderMultiple, 1, 1000000, 1)
	if strategy.IsEnabled != 0 {
		strategy.IsEnabled = 1
	}

	updated, err := uc.repo.UpsertStrategy(strategy)
	if err != nil {
		return nil, err
	}
	entityID := fmt.Sprintf("%d", updated.ID)
	uc.recordAudit(c, "UPSERT_STRATEGY", "ReplenishmentStrategy", entityID, nil, updated)
	return updated, nil
}

func (uc *ReplenishmentUsecase) ListPlans(params *domain.ReplenishmentPlanListParams) ([]domain.ReplenishmentPlan, int64, error) {
	if params == nil {
		params = &domain.ReplenishmentPlanListParams{}
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
	plans, total, err := uc.repo.ListPlans(params)
	if err != nil {
		return nil, 0, err
	}
	if err := uc.enrichPlansWithPackagingAlerts(plans); err != nil {
		return nil, 0, err
	}
	return plans, total, nil
}

func (uc *ReplenishmentUsecase) ListRuns(params *domain.ReplenishmentRunListParams) ([]domain.ReplenishmentRun, int64, error) {
	if params == nil {
		params = &domain.ReplenishmentRunListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.repo.ListRuns(params)
}

func (uc *ReplenishmentUsecase) EnsureDailyCleanup() error {
	cfg, err := uc.repo.GetConfig()
	if err != nil {
		return err
	}
	today := beginOfDay(time.Now())
	if err := uc.repo.DeletePlansBefore(today); err != nil {
		return err
	}
	cfg.LastCleanupDate = &today
	return uc.repo.SaveConfig(cfg)
}

func (uc *ReplenishmentUsecase) GenerateDailyPlans(c *gin.Context, params *domain.ReplenishmentGenerateParams) (plans []domain.ReplenishmentPlan, generatedCount int, err error) {
	if params == nil {
		params = &domain.ReplenishmentGenerateParams{}
	}
	if params.TriggerType == "" {
		params.TriggerType = domain.ReplenishmentTriggerManual
	}

	startedAt := time.Now()
	run := &domain.ReplenishmentRun{
		RunNo:       buildRunNo(startedAt),
		TriggerType: params.TriggerType,
		Status:      domain.ReplenishmentRunRunning,
		WindowDays:  30,
		StartedAt:   &startedAt,
		CreatedBy:   params.OperatorID,
	}
	run.InputSummary = marshalJSONStringPtr(map[string]any{
		"trigger_type": params.TriggerType,
		"operator_id":  params.OperatorID,
		"started_at":   startedAt.Format(time.RFC3339),
	})
	if createErr := uc.repo.CreateRunWithItems(run, nil); createErr != nil {
		return nil, 0, createErr
	}
	defer func() {
		finishedAt := time.Now()
		run.FinishedAt = &finishedAt
		run.OutputSummary = marshalJSONStringPtr(map[string]any{
			"generated":       generatedCount > 0,
			"generated_count": generatedCount,
			"current_count":   len(plans),
		})
		if err != nil {
			run.Status = domain.ReplenishmentRunFailed
			msg := err.Error()
			run.ErrorMessage = &msg
		} else {
			run.Status = domain.ReplenishmentRunSuccess
			run.ErrorMessage = nil
		}
		if updateErr := uc.repo.UpdateRun(run); updateErr != nil && err == nil {
			err = updateErr
		}
	}()

	if err = uc.EnsureDailyCleanup(); err != nil {
		return nil, 0, err
	}

	cfg, err := uc.repo.GetConfig()
	if err != nil {
		return nil, 0, err
	}
	run.WindowDays = clampU32(cfg.DemandWindowDays, 1, 365, 30)

	today := beginOfDay(time.Now())
	plans, err = uc.buildPlansForDate(today, cfg)
	if err != nil {
		return nil, 0, err
	}
	generatedCount, err = uc.repo.CreatePlans(plans)
	if err != nil {
		return nil, 0, err
	}

	cfg.LastGeneratedDate = &today
	if err := uc.repo.SaveConfig(cfg); err != nil {
		return nil, 0, err
	}

	list, _, err := uc.repo.ListPlans(&domain.ReplenishmentPlanListParams{
		Page:     1,
		PageSize: 2000,
		Date:     &today,
	})
	if err != nil {
		return nil, 0, err
	}
	uc.recordAudit(c, "GENERATE_PLANS", "ReplenishmentPlan", today.Format("2006-01-02"), nil, map[string]any{
		"generated":       generatedCount > 0,
		"generated_count": generatedCount,
		"current_count":   len(list),
		"trigger_type":    params.TriggerType,
	})
	return list, generatedCount, nil
}

func (uc *ReplenishmentUsecase) buildPlansForDate(planDate time.Time, cfg *domain.ReplenishmentConfig) ([]domain.ReplenishmentPlan, error) {
	balanceRows, err := uc.repo.LoadBalanceRows()
	if err != nil {
		return nil, err
	}
	if len(balanceRows) == 0 {
		return []domain.ReplenishmentPlan{}, nil
	}

	productSet := map[uint64]struct{}{}
	keys := make([]planKey, 0, len(balanceRows))
	balanceMap := map[planKey]domain.ReplenishmentBalanceRow{}
	for _, row := range balanceRows {
		key := planKey{ProductID: row.ProductID, WarehouseID: row.WarehouseID}
		keys = append(keys, key)
		balanceMap[key] = row
		productSet[row.ProductID] = struct{}{}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].ProductID == keys[j].ProductID {
			return keys[i].WarehouseID < keys[j].WarehouseID
		}
		return keys[i].ProductID < keys[j].ProductID
	})

	productIDs := make([]uint64, 0, len(productSet))
	for productID := range productSet {
		productIDs = append(productIDs, productID)
	}
	profiles, err := uc.repo.LoadProductProfiles(productIDs)
	if err != nil {
		return nil, err
	}

	strategies, err := uc.repo.ListActiveStrategies()
	if err != nil {
		return nil, err
	}

	windows := map[uint32]struct{}{clampU32(cfg.DemandWindowDays, 1, 365, 30): {}}
	for _, strategy := range strategies {
		windows[clampU32(strategy.DemandWindowDays, 1, 365, 30)] = struct{}{}
	}
	demandByWindow := map[uint32]map[planKey]uint64{}
	for window := range windows {
		rows, err := uc.repo.LoadDemandByWindowDays(window)
		if err != nil {
			return nil, err
		}
		demandByWindow[window] = toPlanDemandMap(rows)
	}

	plans := make([]domain.ReplenishmentPlan, 0, len(keys))
	for _, key := range keys {
		profile := profiles[key.ProductID]
		strategy := matchStrategy(strategies, key, profile)
		if strategy == nil {
			// 无命中策略时使用全局默认策略
			defaultStrategy := domain.ReplenishmentStrategy{
				Name:                 "GLOBAL_DEFAULT",
				Priority:             0,
				DemandWindowDays:     clampU32(cfg.DemandWindowDays, 1, 365, 30),
				ProcurementCycleDays: clampU32(cfg.DefaultLeadTimeDays, 0, 365, 15),
				PackDays:             0,
				LogisticsDays:        0,
				SafetyDays:           clampU32(cfg.DefaultSafetyDays, 0, 365, 7),
				ZeroSalesPurchaseQty: 0,
				MOQ:                  clampU32(cfg.DefaultMOQ, 1, 1000000, 1),
				OrderMultiple:        clampU32(cfg.DefaultOrderMultiple, 1, 1000000, 1),
			}
			strategy = &defaultStrategy
		}
		if isExcludedProductStatus(profile.ProductStatus) {
			continue
		}

		windowDays := clampU32(strategy.DemandWindowDays, 1, 365, 30)
		shippedQty := demandByWindow[windowDays][key]
		dailyDemand := float64(shippedQty) / float64(windowDays)

		balance := balanceMap[key]
		netSupply := int64(balance.AvailableQuantity) +
			int64(balance.PurchasingInTransit) +
			int64(balance.PendingInspection) -
			int64(balance.ReservedQuantity)
		coverageDays := clampU32(strategy.ProcurementCycleDays+strategy.PackDays+strategy.LogisticsDays+strategy.SafetyDays, 0, 2000, 0)
		targetStock := uint64(math.Ceil(dailyDemand * float64(coverageDays)))

		shortage := int64(targetStock) - netSupply
		if shortage < 0 {
			shortage = 0
		}
		suggestedQty := uint64(0)
		zeroStock := balance.AvailableQuantity == 0 &&
			balance.PurchasingInTransit == 0 &&
			balance.PendingInspection == 0
		if shippedQty == 0 && strategy.ZeroSalesPurchaseQty > 0 && zeroStock {
			suggestedQty = roundToPurchaseRule(uint64(strategy.ZeroSalesPurchaseQty), uint64(strategy.MOQ), uint64(strategy.OrderMultiple))
			targetStock = uint64(strategy.ZeroSalesPurchaseQty)
			shortage = int64(strategy.ZeroSalesPurchaseQty)
		} else if shortage > 0 {
			suggestedQty = roundToPurchaseRule(uint64(shortage), uint64(strategy.MOQ), uint64(strategy.OrderMultiple))
		}

		if suggestedQty == 0 {
			continue
		}

		status := domain.ReplenishmentPlanPending
		var remark *string
		if profile.SupplierID == nil || *profile.SupplierID == 0 {
			status = domain.ReplenishmentPlanCancelled
			msg := "NO_SUPPLIER"
			remark = &msg
		}

		var strategyID *uint64
		if strategy.ID > 0 {
			strategyID = &strategy.ID
		}
		plans = append(plans, domain.ReplenishmentPlan{
			PlanDate:         planDate,
			ProductID:        key.ProductID,
			WarehouseID:      key.WarehouseID,
			SupplierID:       profile.SupplierID,
			StrategyID:       strategyID,
			DailyDemand:      round4(dailyDemand),
			DemandWindowDays: windowDays,
			CoverageDays:     coverageDays,
			NetSupply:        netSupply,
			TargetStock:      targetStock,
			ShortageQty:      uint64(shortage),
			SuggestedQty:     suggestedQty,
			MOQ:              strategy.MOQ,
			OrderMultiple:    strategy.OrderMultiple,
			UnitCost:         profile.QuotePrice,
			Status:           status,
			Remark:           remark,
		})
	}

	return plans, nil
}

func (uc *ReplenishmentUsecase) ConvertPlansToPurchaseOrders(c *gin.Context, params *domain.ReplenishmentPlanConvertParams) ([]domain.PurchaseOrder, error) {
	if uc.poUsecase == nil {
		return nil, ErrReplenishmentMissingPoUC
	}
	if params == nil {
		params = &domain.ReplenishmentPlanConvertParams{}
	}
	if params.PlanDate == nil {
		today := beginOfDay(time.Now())
		params.PlanDate = &today
	}

	plans, err := uc.repo.ListConvertiblePlans(params)
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, ErrReplenishmentNoPlans
	}

	sourceOrders := make([]*domain.PurchaseOrder, 0, len(plans))
	defaultCurrency := "USD"
	if uc.defaultsProvider != nil {
		if configuredCurrency := strings.TrimSpace(uc.defaultsProvider.GetDefaultBaseCurrency()); configuredCurrency != "" {
			defaultCurrency = configuredCurrency
		}
	}
	for _, plan := range plans {
		if plan.SupplierID == nil || *plan.SupplierID == 0 {
			return nil, ErrReplenishmentMissingSupplier
		}
		unitCost := 0.0
		if plan.UnitCost != nil {
			unitCost = *plan.UnitCost
		}
		supplierID := *plan.SupplierID
		remark := fmt.Sprintf("AUTO_FROM_DAILY_PLAN:%s;warehouse=%d;plan=%d", params.PlanDate.Format("2006-01-02"), plan.WarehouseID, plan.ID)
		planID := plan.ID
		sourceOrders = append(sourceOrders, &domain.PurchaseOrder{
			SupplierID:  &supplierID,
			WarehouseID: &plan.WarehouseID,
			Currency:    defaultCurrency,
			Remark:      remark,
			CreatedBy:   params.OperatorID,
			UpdatedBy:   params.OperatorID,
			Items: []domain.PurchaseOrderItem{
				{
					ProductID:    plan.ProductID,
					SourcePlanID: &planID,
					QtyOrdered:   plan.SuggestedQty,
					UnitCost:     unitCost,
					Currency:     defaultCurrency,
				},
			},
		})
	}
	result, err := uc.poUsecase.CreateBatch(c, sourceOrders)
	if err != nil {
		return nil, err
	}
	planFirstOrder := map[uint64]uint64{}
	links := make([]domain.ReplenishmentPlanPurchaseOrderLink, 0, len(result))
	for _, order := range result {
		for _, planID := range order.SourcePlanIDs {
			if planID == 0 || order.ID == 0 {
				continue
			}
			links = append(links, domain.ReplenishmentPlanPurchaseOrderLink{
				PlanID:          planID,
				PurchaseOrderID: order.ID,
			})
			if _, exists := planFirstOrder[planID]; !exists {
				planFirstOrder[planID] = order.ID
			}
		}
	}
	if err := uc.repo.LinkPlansToPurchaseOrders(links); err != nil {
		return nil, err
	}
	for _, plan := range plans {
		purchaseOrderID, ok := planFirstOrder[plan.ID]
		if !ok || purchaseOrderID == 0 {
			return nil, fmt.Errorf("replenishment plan %d missing linked purchase order", plan.ID)
		}
		if err := uc.repo.MarkPlansConverted([]uint64{plan.ID}, purchaseOrderID); err != nil {
			return nil, err
		}
	}
	uc.recordAudit(c, "CONVERT_PLANS", "ReplenishmentPlan", params.PlanDate.Format("2006-01-02"), nil, map[string]any{
		"created_count": len(result),
		"plan_ids":      params.PlanIDs,
	})

	return result, nil
}

func (uc *ReplenishmentUsecase) DeletePlansByPurchaseOrderID(purchaseOrderID uint64) error {
	return uc.repo.DeletePlansByPurchaseOrderID(purchaseOrderID)
}

func (uc *ReplenishmentUsecase) DeletePlanByID(c *gin.Context, planID uint64) error {
	deleted, err := uc.repo.DeletePendingPlanByID(planID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrReplenishmentPlanNotDeletable
	}
	uc.recordAudit(c, "DELETE_PLAN", "ReplenishmentPlan", fmt.Sprintf("%d", planID), nil, map[string]any{"deleted": true})
	return nil
}

func (uc *ReplenishmentUsecase) IsGeneratedToday() (bool, error) {
	cfg, err := uc.repo.GetConfig()
	if err != nil {
		return false, err
	}
	if cfg.LastGeneratedDate == nil {
		return false, nil
	}
	return sameDay(*cfg.LastGeneratedDate, beginOfDay(time.Now())), nil
}

type planKey struct {
	ProductID   uint64
	WarehouseID uint64
}

func toPlanDemandMap(rows []domain.ReplenishmentDemandRow) map[planKey]uint64 {
	result := map[planKey]uint64{}
	for _, row := range rows {
		result[planKey{ProductID: row.ProductID, WarehouseID: row.WarehouseID}] = row.ShippedQty
	}
	return result
}

func (uc *ReplenishmentUsecase) enrichPlansWithPackagingAlerts(plans []domain.ReplenishmentPlan) error {
	if len(plans) == 0 {
		return nil
	}

	productSet := map[uint64]struct{}{}
	for _, plan := range plans {
		if plan.ProductID == 0 || plan.SuggestedQty == 0 {
			continue
		}
		productSet[plan.ProductID] = struct{}{}
	}
	if len(productSet) == 0 {
		return nil
	}

	productIDs := make([]uint64, 0, len(productSet))
	for productID := range productSet {
		productIDs = append(productIDs, productID)
	}
	requirementsBySku, err := uc.repo.LoadPackagingRequirementsByProduct(productIDs)
	if err != nil {
		return err
	}
	if len(requirementsBySku) == 0 {
		return nil
	}

	itemSet := map[uint64]struct{}{}
	for _, requirements := range requirementsBySku {
		for _, requirement := range requirements {
			if requirement.PackagingItemID == 0 {
				continue
			}
			itemSet[requirement.PackagingItemID] = struct{}{}
		}
	}
	if len(itemSet) == 0 {
		return nil
	}

	itemIDs := make([]uint64, 0, len(itemSet))
	for itemID := range itemSet {
		itemIDs = append(itemIDs, itemID)
	}
	packagingItems, err := uc.repo.LoadPackagingItems(itemIDs)
	if err != nil {
		return err
	}
	if len(packagingItems) == 0 {
		return nil
	}

	for i := range plans {
		plan := &plans[i]
		requirements := requirementsBySku[plan.ProductID]
		if len(requirements) == 0 || plan.SuggestedQty == 0 {
			continue
		}

		totalShortage := uint64(0)
		alerts := make([]string, 0, len(requirements))
		for _, requirement := range requirements {
			item, ok := packagingItems[requirement.PackagingItemID]
			if !ok || requirement.QuantityPerUnit <= 0 {
				continue
			}
			requiredQty := uint64(math.Ceil(float64(plan.SuggestedQty) * requirement.QuantityPerUnit))
			if requiredQty == 0 || requiredQty <= item.QuantityOnHand {
				continue
			}
			shortageQty := requiredQty - item.QuantityOnHand
			totalShortage += shortageQty

			itemLabel := item.ItemCode
			if itemLabel == "" {
				itemLabel = fmt.Sprintf("ITEM-%d", item.PackagingItemID)
			}
			alerts = append(alerts, fmt.Sprintf("%s 缺口 %d", itemLabel, shortageQty))
		}

		if totalShortage > 0 {
			plan.PackagingShortageQty = totalShortage
			plan.PackagingAlert = strings.Join(alerts, "; ")
		}
	}
	return nil
}

func matchStrategy(strategies []domain.ReplenishmentStrategy, key planKey, profile domain.ReplenishmentProductProfile) *domain.ReplenishmentStrategy {
	for i := range strategies {
		strategy := strategies[i]
		if strategy.ProductID != nil && *strategy.ProductID != key.ProductID {
			continue
		}
		if strategy.WarehouseID != nil && *strategy.WarehouseID != key.WarehouseID {
			continue
		}
		if strategy.SupplierID != nil {
			if profile.SupplierID == nil || *strategy.SupplierID != *profile.SupplierID {
				continue
			}
		}
		if strategy.Marketplace != nil {
			if profile.Marketplace == nil || *strategy.Marketplace != *profile.Marketplace {
				continue
			}
		}
		return &strategy
	}
	return nil
}

func isExcludedProductStatus(status string) bool {
	s := strings.ToUpper(strings.TrimSpace(status))
	return s == "OFF_SHELF" || s == "DRAFT"
}

func beginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func sameDay(a, b time.Time) bool {
	aa := beginOfDay(a)
	bb := beginOfDay(b)
	return aa.Equal(bb)
}

func roundToPurchaseRule(shortage, moq, multiple uint64) uint64 {
	if shortage == 0 {
		return 0
	}
	if moq == 0 {
		moq = 1
	}
	if multiple == 0 {
		multiple = 1
	}
	result := shortage
	if result < moq {
		result = moq
	}
	if mod := result % multiple; mod != 0 {
		result += multiple - mod
	}
	return result
}

func clampU32(value, minVal, maxVal, fallback uint32) uint32 {
	if value == 0 {
		value = fallback
	}
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}

func buildRunNo(t time.Time) string {
	return numbering.Generate("RPL", t)
}

func marshalJSONStringPtr(v any) *string {
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	s := string(b)
	return &s
}

func (uc *ReplenishmentUsecase) recordAudit(c *gin.Context, action, entityType, entityID string, before, after any) {
	if uc.auditLogger == nil || c == nil {
		return
	}
	_ = uc.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Procurement",
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
}
