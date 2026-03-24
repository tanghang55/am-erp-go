package repository

import (
	"errors"
	"time"

	"am-erp-go/internal/module/procurement/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type replenishmentRepository struct {
	db *gorm.DB
}

func NewReplenishmentRepository(db *gorm.DB) domain.ReplenishmentRepository {
	return &replenishmentRepository{db: db}
}

func (r *replenishmentRepository) GetConfig() (*domain.ReplenishmentConfig, error) {
	var cfg domain.ReplenishmentConfig
	err := r.db.Order("id ASC").First(&cfg).Error
	if err == nil {
		return &cfg, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	cfg = domain.ReplenishmentConfig{
		IsEnabled:            1,
		IntervalMinutes:      1440,
		DemandWindowDays:     30,
		DefaultLeadTimeDays:  15,
		DefaultSafetyDays:    7,
		DefaultMOQ:           1,
		DefaultOrderMultiple: 1,
	}
	if err := r.db.Create(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *replenishmentRepository) SaveConfig(config *domain.ReplenishmentConfig) error {
	if config == nil {
		return nil
	}
	if config.ID == 0 {
		return r.db.Create(config).Error
	}
	return r.db.Save(config).Error
}

func (r *replenishmentRepository) ListStrategies(params *domain.ReplenishmentStrategyListParams) ([]domain.ReplenishmentStrategy, int64, error) {
	var rows []domain.ReplenishmentStrategy
	var total int64

	query := r.db.Table("procurement_replenishment_strategy strategy").
		Joins("LEFT JOIN product ON product.id = strategy.product_id").
		Joins("LEFT JOIN warehouse ON warehouse.id = strategy.warehouse_id").
		Joins("LEFT JOIN supplier ON supplier.id = strategy.supplier_id")
	if params != nil && params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where(`
			strategy.name LIKE ?
			OR product.seller_sku LIKE ?
			OR product.title LIKE ?
			OR warehouse.name LIKE ?
			OR supplier.name LIKE ?
		`, kw, kw, kw, kw, kw)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := 1
	pageSize := 20
	if params != nil {
		if params.Page > 0 {
			page = params.Page
		}
		if params.PageSize > 0 {
			pageSize = params.PageSize
		}
	}
	offset := (page - 1) * pageSize

	if err := query.Select(`
			strategy.*,
			product.seller_sku AS seller_sku,
			product.title AS product_title,
			warehouse.code AS warehouse_code,
			warehouse.name AS warehouse_name,
			supplier.supplier_code AS supplier_code,
			supplier.name AS supplier_name
		`).
		Order("strategy.priority DESC, strategy.id DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return []domain.ReplenishmentStrategy{}, total, nil
	}
	return rows, total, nil
}

func (r *replenishmentRepository) ListActiveStrategies() ([]domain.ReplenishmentStrategy, error) {
	var rows []domain.ReplenishmentStrategy
	if err := r.db.Where("is_enabled = 1").Order("priority DESC, id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []domain.ReplenishmentStrategy{}, nil
	}
	return rows, nil
}

func (r *replenishmentRepository) UpsertStrategy(strategy *domain.ReplenishmentStrategy) (*domain.ReplenishmentStrategy, error) {
	if strategy == nil {
		return nil, nil
	}
	if strategy.ID == 0 {
		if err := r.db.Create(strategy).Error; err != nil {
			return nil, err
		}
		return strategy, nil
	}

	updates := map[string]any{
		"name":                    strategy.Name,
		"priority":                strategy.Priority,
		"is_enabled":              strategy.IsEnabled,
		"product_id":              strategy.ProductID,
		"warehouse_id":            strategy.WarehouseID,
		"supplier_id":             strategy.SupplierID,
		"marketplace":             strategy.Marketplace,
		"condition_json":          strategy.ConditionJSON,
		"demand_window_days":      strategy.DemandWindowDays,
		"procurement_cycle_days":  strategy.ProcurementCycleDays,
		"pack_days":               strategy.PackDays,
		"logistics_days":          strategy.LogisticsDays,
		"safety_days":             strategy.SafetyDays,
		"zero_sales_purchase_qty": strategy.ZeroSalesPurchaseQty,
		"moq":                     strategy.MOQ,
		"order_multiple":          strategy.OrderMultiple,
		"remark":                  strategy.Remark,
		"updated_by":              strategy.UpdatedBy,
	}
	if err := r.db.Model(&domain.ReplenishmentStrategy{}).
		Where("id = ?", strategy.ID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	var row domain.ReplenishmentStrategy
	if err := r.db.First(&row, strategy.ID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *replenishmentRepository) ListPlans(params *domain.ReplenishmentPlanListParams) ([]domain.ReplenishmentPlan, int64, error) {
	var rows []domain.ReplenishmentPlan
	var total int64

	query := r.db.Table("procurement_replenishment_plan plan").
		Joins("LEFT JOIN product ON product.id = plan.product_id").
		Joins("LEFT JOIN warehouse ON warehouse.id = plan.warehouse_id").
		Joins("LEFT JOIN supplier ON supplier.id = plan.supplier_id").
		Joins("LEFT JOIN procurement_replenishment_strategy strategy ON strategy.id = plan.strategy_id").
		Joins("LEFT JOIN purchase_order po ON po.id = plan.purchase_order_id").
		Joins(`LEFT JOIN (
			SELECT
				link.plan_id,
				GROUP_CONCAT(po.po_number ORDER BY po.po_number SEPARATOR ', ') AS purchase_order_numbers
			FROM procurement_replenishment_plan_purchase_order link
			INNER JOIN purchase_order po ON po.id = link.purchase_order_id
			GROUP BY link.plan_id
		) plan_po ON plan_po.plan_id = plan.id`)
	planDate := time.Now()
	if params != nil && params.Date != nil {
		planDate = *params.Date
	}
	query = query.Where("plan.plan_date = ?", planDate.Format("2006-01-02"))
	if params != nil && params.Status != "" {
		query = query.Where("plan.status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := 1
	pageSize := 50
	if params != nil {
		if params.Page > 0 {
			page = params.Page
		}
		if params.PageSize > 0 {
			pageSize = params.PageSize
		}
	}
	offset := (page - 1) * pageSize

	if err := query.Select(`
			plan.*,
			product.seller_sku AS seller_sku,
			product.title AS product_title,
			product.image_url AS product_image_url,
			warehouse.code AS warehouse_code,
			warehouse.name AS warehouse_name,
			supplier.supplier_code AS supplier_code,
			supplier.name AS supplier_name,
			strategy.name AS strategy_name,
			po.po_number AS purchase_order_number,
			plan_po.purchase_order_numbers AS purchase_order_numbers
		`).
		Order("plan.suggested_qty DESC, plan.id ASC").
		Offset(offset).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return []domain.ReplenishmentPlan{}, total, nil
	}
	return rows, total, nil
}

func (r *replenishmentRepository) CreatePlans(plans []domain.ReplenishmentPlan) (int, error) {
	if len(plans) == 0 {
		return 0, nil
	}
	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plan_date"}, {Name: "product_id"}, {Name: "warehouse_id"}},
		DoNothing: true,
	}).Create(&plans)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

func (r *replenishmentRepository) DeletePendingPlanByID(planID uint64) (bool, error) {
	if planID == 0 {
		return false, nil
	}
	result := r.db.Where("id = ? AND status = ?", planID, domain.ReplenishmentPlanPending).
		Delete(&domain.ReplenishmentPlan{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *replenishmentRepository) DeletePlansBefore(date time.Time) error {
	return r.db.Where("plan_date < ?", date.Format("2006-01-02")).
		Delete(&domain.ReplenishmentPlan{}).Error
}

func (r *replenishmentRepository) DeletePlansByPurchaseOrderID(purchaseOrderID uint64) error {
	_ = purchaseOrderID
	return nil
}

func (r *replenishmentRepository) ListConvertiblePlans(params *domain.ReplenishmentPlanConvertParams) ([]domain.ReplenishmentPlan, error) {
	if params == nil {
		params = &domain.ReplenishmentPlanConvertParams{}
	}
	planDate := time.Now()
	if params.PlanDate != nil {
		planDate = *params.PlanDate
	}

	query := r.db.Where(
		"plan_date = ? AND status = ? AND suggested_qty > 0",
		planDate.Format("2006-01-02"),
		domain.ReplenishmentPlanPending,
	)
	if len(params.PlanIDs) > 0 {
		query = query.Where("id IN ?", params.PlanIDs)
	}

	var rows []domain.ReplenishmentPlan
	if err := query.Order("supplier_id ASC, warehouse_id ASC, product_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []domain.ReplenishmentPlan{}, nil
	}
	return rows, nil
}

func (r *replenishmentRepository) MarkPlansConverted(planIDs []uint64, purchaseOrderID uint64) error {
	if len(planIDs) == 0 || purchaseOrderID == 0 {
		return nil
	}
	now := time.Now()
	return r.db.Model(&domain.ReplenishmentPlan{}).
		Where("id IN ?", planIDs).
		Updates(map[string]any{
			"status":            domain.ReplenishmentPlanConverted,
			"purchase_order_id": purchaseOrderID,
			"converted_at":      now,
			"gmt_modified":      now,
		}).Error
}

func (r *replenishmentRepository) LinkPlansToPurchaseOrders(links []domain.ReplenishmentPlanPurchaseOrderLink) error {
	if len(links) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plan_id"}, {Name: "purchase_order_id"}},
		DoNothing: true,
	}).Create(&links).Error
}

func (r *replenishmentRepository) ListPolicies(params *domain.ReplenishmentPolicyListParams) ([]domain.ReplenishmentPolicy, int64, error) {
	var rows []domain.ReplenishmentPolicy
	var total int64

	query := r.db.Model(&domain.ReplenishmentPolicy{})
	if params != nil {
		if params.ProductID != nil {
			query = query.Where("product_id = ?", *params.ProductID)
		}
		if params.Keyword != "" {
			kw := "%" + params.Keyword + "%"
			query = query.Joins("LEFT JOIN product ON product.id = procurement_replenishment_policy.product_id").
				Where("product.seller_sku LIKE ? OR product.title LIKE ?", kw, kw)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := 1
	pageSize := 20
	if params != nil {
		if params.Page > 0 {
			page = params.Page
		}
		if params.PageSize > 0 {
			pageSize = params.PageSize
		}
	}
	offset := (page - 1) * pageSize

	if err := query.Order("gmt_modified DESC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return []domain.ReplenishmentPolicy{}, total, nil
	}
	return rows, total, nil
}

func (r *replenishmentRepository) UpsertPolicy(policy *domain.ReplenishmentPolicy) (*domain.ReplenishmentPolicy, error) {
	if policy == nil {
		return nil, nil
	}
	var existing domain.ReplenishmentPolicy
	err := r.db.Where("product_id = ?", policy.ProductID).First(&existing).Error
	if err == nil {
		policy.ID = existing.ID
		if err := r.db.Save(policy).Error; err != nil {
			return nil, err
		}
		return policy, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err := r.db.Create(policy).Error; err != nil {
		return nil, err
	}
	return policy, nil
}

func (r *replenishmentRepository) ListRuns(params *domain.ReplenishmentRunListParams) ([]domain.ReplenishmentRun, int64, error) {
	var rows []domain.ReplenishmentRun
	var total int64

	query := r.db.Model(&domain.ReplenishmentRun{})
	if params != nil {
		if params.Status != "" {
			query = query.Where("status = ?", params.Status)
		}
		if params.Triggered != "" {
			query = query.Where("trigger_type = ?", params.Triggered)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := 1
	pageSize := 20
	if params != nil {
		if params.Page > 0 {
			page = params.Page
		}
		if params.PageSize > 0 {
			pageSize = params.PageSize
		}
	}
	offset := (page - 1) * pageSize

	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return []domain.ReplenishmentRun{}, total, nil
	}
	return rows, total, nil
}

func (r *replenishmentRepository) GetRunByID(runID uint64) (*domain.ReplenishmentRun, error) {
	var run domain.ReplenishmentRun
	if err := r.db.First(&run, runID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &run, nil
}

func (r *replenishmentRepository) ListRunItems(runID uint64) ([]domain.ReplenishmentItem, error) {
	var items []domain.ReplenishmentItem
	if err := r.db.Where("run_id = ?", runID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []domain.ReplenishmentItem{}, nil
	}
	return items, nil
}

func (r *replenishmentRepository) CreateRunWithItems(run *domain.ReplenishmentRun, items []domain.ReplenishmentItem) error {
	if run == nil {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(run).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		for i := range items {
			items[i].RunID = run.ID
		}
		return tx.Create(&items).Error
	})
}

func (r *replenishmentRepository) CreateRunItems(runID uint64, items []domain.ReplenishmentItem) error {
	if runID == 0 || len(items) == 0 {
		return nil
	}
	for i := range items {
		items[i].RunID = runID
	}
	return r.db.Create(&items).Error
}

func (r *replenishmentRepository) UpdateRun(run *domain.ReplenishmentRun) error {
	if run == nil {
		return nil
	}
	return r.db.Save(run).Error
}

func (r *replenishmentRepository) ListPoliciesByProductIDs(productIDs []uint64) (map[uint64]domain.ReplenishmentPolicy, error) {
	result := map[uint64]domain.ReplenishmentPolicy{}
	if len(productIDs) == 0 {
		return result, nil
	}

	var rows []domain.ReplenishmentPolicy
	if err := r.db.Where("product_id IN ? AND is_enabled = 1", productIDs).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.ProductID] = row
	}
	return result, nil
}

func (r *replenishmentRepository) LoadDemandByWindowDays(windowDays uint32) ([]domain.ReplenishmentDemandRow, error) {
	since := time.Now().AddDate(0, 0, -int(windowDays))
	type row struct {
		ProductID   uint64 `gorm:"column:product_id"`
		WarehouseID uint64 `gorm:"column:warehouse_id"`
		ShippedQty  uint64 `gorm:"column:shipped_qty"`
	}
	var rows []row

	err := r.db.Table("inventory_movement").
		Select("product_id, warehouse_id, COALESCE(SUM(CASE WHEN quantity > 0 THEN quantity ELSE 0 END), 0) AS shipped_qty").
		Where("movement_type = ? AND operated_at >= ?", "SALES_SHIP", since).
		Group("product_id, warehouse_id").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.ReplenishmentDemandRow, 0, len(rows))
	for _, item := range rows {
		result = append(result, domain.ReplenishmentDemandRow{
			ProductID:   item.ProductID,
			WarehouseID: item.WarehouseID,
			ShippedQty:  item.ShippedQty,
		})
	}
	return result, nil
}

func (r *replenishmentRepository) LoadBalanceRows() ([]domain.ReplenishmentBalanceRow, error) {
	type row struct {
		ProductID           uint64 `gorm:"column:product_id"`
		WarehouseID         uint64 `gorm:"column:warehouse_id"`
		AvailableQuantity   uint64 `gorm:"column:available_quantity"`
		ReservedQuantity    uint64 `gorm:"column:reserved_quantity"`
		PurchasingInTransit uint64 `gorm:"column:purchasing_in_transit"`
		PendingInspection   uint64 `gorm:"column:pending_inspection"`
	}
	var rows []row
	err := r.db.Table("inventory_balance").
		Select("product_id, warehouse_id, available_quantity, reserved_quantity, purchasing_in_transit, pending_inspection").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.ReplenishmentBalanceRow, 0, len(rows))
	for _, item := range rows {
		result = append(result, domain.ReplenishmentBalanceRow{
			ProductID:           item.ProductID,
			WarehouseID:         item.WarehouseID,
			AvailableQuantity:   item.AvailableQuantity,
			ReservedQuantity:    item.ReservedQuantity,
			PurchasingInTransit: item.PurchasingInTransit,
			PendingInspection:   item.PendingInspection,
		})
	}
	return result, nil
}

func (r *replenishmentRepository) LoadProductProfiles(productIDs []uint64) (map[uint64]domain.ReplenishmentProductProfile, error) {
	result := map[uint64]domain.ReplenishmentProductProfile{}
	if len(productIDs) == 0 {
		return result, nil
	}

	type row struct {
		ProductID     uint64   `gorm:"column:product_id"`
		SupplierID    *uint64  `gorm:"column:supplier_id"`
		Marketplace   *string  `gorm:"column:marketplace"`
		ProductStatus string   `gorm:"column:product_status"`
		QuoteMOQ      *uint64  `gorm:"column:quote_moq"`
		QuoteLeadDay  *uint64  `gorm:"column:quote_lead_day"`
		QuotePrice    *float64 `gorm:"column:quote_price"`
	}
	var rows []row
	err := r.db.Table("product p").
		Select("p.id AS product_id, p.supplier_id, p.marketplace, p.status AS product_status, psq.qty_moq AS quote_moq, psq.lead_time_days AS quote_lead_day, psq.price AS quote_price").
		Joins("LEFT JOIN product_supplier_quote psq ON psq.product_id = p.id AND psq.supplier_id = p.supplier_id").
		Where("p.id IN ?", productIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, item := range rows {
		result[item.ProductID] = domain.ReplenishmentProductProfile{
			ProductID:     item.ProductID,
			SupplierID:    item.SupplierID,
			Marketplace:   item.Marketplace,
			ProductStatus: item.ProductStatus,
			QuoteMOQ:      item.QuoteMOQ,
			QuoteLeadDays: item.QuoteLeadDay,
			QuotePrice:    item.QuotePrice,
		}
	}
	return result, nil
}

func (r *replenishmentRepository) LoadPackagingRequirementsByProduct(productIDs []uint64) (map[uint64][]domain.ReplenishmentPackagingRequirement, error) {
	result := map[uint64][]domain.ReplenishmentPackagingRequirement{}
	if len(productIDs) == 0 {
		return result, nil
	}

	type row struct {
		ProductID       uint64  `gorm:"column:product_id"`
		PackagingItemID uint64  `gorm:"column:packaging_item_id"`
		QuantityPerUnit float64 `gorm:"column:quantity_per_unit"`
	}
	var rows []row
	err := r.db.Table("product_packaging_items").
		Select("product_id AS product_id, packaging_item_id, quantity_per_unit").
		Where("product_id IN ?", productIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, item := range rows {
		result[item.ProductID] = append(result[item.ProductID], domain.ReplenishmentPackagingRequirement{
			ProductID:       item.ProductID,
			PackagingItemID: item.PackagingItemID,
			QuantityPerUnit: item.QuantityPerUnit,
		})
	}
	return result, nil
}

func (r *replenishmentRepository) LoadPackagingItems(itemIDs []uint64) (map[uint64]domain.ReplenishmentPackagingItem, error) {
	result := map[uint64]domain.ReplenishmentPackagingItem{}
	if len(itemIDs) == 0 {
		return result, nil
	}

	type row struct {
		PackagingItemID uint64 `gorm:"column:id"`
		ItemCode        string `gorm:"column:item_code"`
		ItemName        string `gorm:"column:item_name"`
		QuantityOnHand  uint64 `gorm:"column:quantity_on_hand"`
	}

	var rows []row
	err := r.db.Table("packaging_item").
		Select("id, item_code, item_name, quantity_on_hand").
		Where("id IN ?", itemIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, item := range rows {
		result[item.PackagingItemID] = domain.ReplenishmentPackagingItem{
			PackagingItemID: item.PackagingItemID,
			ItemCode:        item.ItemCode,
			ItemName:        item.ItemName,
			QuantityOnHand:  item.QuantityOnHand,
		}
	}
	return result, nil
}

func (r *replenishmentRepository) ListConvertibleItems(runID uint64, itemIDs []uint64) ([]domain.ReplenishmentItem, error) {
	var rows []domain.ReplenishmentItem
	query := r.db.Where("run_id = ? AND status = ? AND suggested_qty > 0", runID, domain.ReplenishmentItemPending)
	if len(itemIDs) > 0 {
		query = query.Where("id IN ?", itemIDs)
	}
	if err := query.Order("supplier_id ASC, warehouse_id ASC, product_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []domain.ReplenishmentItem{}, nil
	}
	return rows, nil
}

func (r *replenishmentRepository) MarkItemsConverted(itemIDs []uint64, purchaseOrderID uint64) error {
	if len(itemIDs) == 0 || purchaseOrderID == 0 {
		return nil
	}
	now := time.Now()
	return r.db.Model(&domain.ReplenishmentItem{}).
		Where("id IN ?", itemIDs).
		Updates(map[string]any{
			"status":            domain.ReplenishmentItemConverted,
			"purchase_order_id": purchaseOrderID,
			"converted_at":      now,
			"gmt_modified":      now,
		}).Error
}
