package repository

import (
	"errors"
	"time"

	"am-erp-go/internal/module/packaging/domain"

	"gorm.io/gorm"
)

type packagingProcurementRepository struct {
	db *gorm.DB
}

func NewPackagingProcurementRepository(db *gorm.DB) domain.PackagingProcurementRepository {
	return &packagingProcurementRepository{db: db}
}

func (r *packagingProcurementRepository) CleanupPlansBefore(date time.Time) error {
	return r.db.Where("plan_date < ?", date.Format("2006-01-02")).
		Delete(&domain.PackagingProcurementPlan{}).Error
}

func (r *packagingProcurementRepository) LoadOrderedProductDemands(planDate time.Time) ([]domain.PackagingProductDemand, error) {
	type row struct {
		ProductID uint64 `gorm:"column:product_id"`
		Qty       uint64 `gorm:"column:qty"`
	}
	var rows []row
	err := r.db.Table("purchase_order_item poi").
		Select("poi.product_id, COALESCE(SUM(poi.qty_ordered), 0) AS qty").
		Joins("INNER JOIN purchase_order po ON po.id = poi.purchase_order_id").
		Where("po.status IN ?", []string{"ORDERED", "SHIPPED", "RECEIVED", "CLOSED"}).
		Where("po.remark LIKE ?", "AUTO_FROM_DAILY_PLAN:"+planDate.Format("2006-01-02")+"%").
		Group("poi.product_id").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.PackagingProductDemand, 0, len(rows))
	for _, item := range rows {
		result = append(result, domain.PackagingProductDemand{
			ProductID: item.ProductID,
			Qty:       item.Qty,
		})
	}
	return result, nil
}

func (r *packagingProcurementRepository) LoadProductPackagingMappings(productIDs []uint64) ([]domain.ProductPackagingMapping, error) {
	if len(productIDs) == 0 {
		return []domain.ProductPackagingMapping{}, nil
	}
	type row struct {
		ProductID       uint64  `gorm:"column:product_id"`
		PackagingItemID uint64  `gorm:"column:packaging_item_id"`
		QuantityPerUnit float64 `gorm:"column:quantity_per_unit"`
	}
	var rows []row
	err := r.db.Table("product_packaging_items").
		Select("product_id, packaging_item_id, quantity_per_unit").
		Where("product_id IN ?", productIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]domain.ProductPackagingMapping, 0, len(rows))
	for _, item := range rows {
		result = append(result, domain.ProductPackagingMapping{
			ProductID:       item.ProductID,
			PackagingItemID: item.PackagingItemID,
			QuantityPerUnit: item.QuantityPerUnit,
		})
	}
	return result, nil
}

func (r *packagingProcurementRepository) LoadPackagingItemSnapshots(itemIDs []uint64) (map[uint64]domain.PackagingItemSnapshot, error) {
	result := map[uint64]domain.PackagingItemSnapshot{}
	if len(itemIDs) == 0 {
		return result, nil
	}

	type row struct {
		ID              uint64  `gorm:"column:id"`
		ItemCode        string  `gorm:"column:item_code"`
		ItemName        string  `gorm:"column:item_name"`
		Unit            string  `gorm:"column:unit"`
		QuantityOnHand  uint64  `gorm:"column:quantity_on_hand"`
		ReorderQuantity *uint64 `gorm:"column:reorder_quantity"`
		UnitCost        float64 `gorm:"column:unit_cost"`
		Currency        string  `gorm:"column:currency"`
		Status          string  `gorm:"column:status"`
	}

	var rows []row
	err := r.db.Table("packaging_item").
		Select("id, item_code, item_name, unit, quantity_on_hand, reorder_quantity, unit_cost, currency, status").
		Where("id IN ?", itemIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, item := range rows {
		result[item.ID] = domain.PackagingItemSnapshot{
			ID:              item.ID,
			ItemCode:        item.ItemCode,
			ItemName:        item.ItemName,
			Unit:            item.Unit,
			QuantityOnHand:  item.QuantityOnHand,
			ReorderQuantity: item.ReorderQuantity,
			UnitCost:        item.UnitCost,
			Currency:        item.Currency,
			Status:          item.Status,
		}
	}
	return result, nil
}

func (r *packagingProcurementRepository) SyncDailyPlans(planDate time.Time, inputs []domain.PackagingPlanInput) ([]domain.PackagingProcurementPlan, int, error) {
	createdCount := 0
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing []domain.PackagingProcurementPlan
		if err := tx.Where("plan_date = ?", planDate.Format("2006-01-02")).
			Find(&existing).Error; err != nil {
			return err
		}

		existingMap := make(map[uint64]domain.PackagingProcurementPlan, len(existing))
		for _, item := range existing {
			existingMap[item.PackagingItemID] = item
		}

		seen := make(map[uint64]struct{}, len(inputs))
		for _, input := range inputs {
			if input.PackagingItemID == 0 {
				continue
			}
			seen[input.PackagingItemID] = struct{}{}

			if row, ok := existingMap[input.PackagingItemID]; ok {
				if row.Status != domain.PackagingProcurementPlanPending {
					continue
				}
				if err := tx.Model(&domain.PackagingProcurementPlan{}).
					Where("id = ?", row.ID).
					Updates(map[string]any{
						"required_qty":  input.RequiredQty,
						"on_hand_qty":   input.OnHandQty,
						"shortage_qty":  input.ShortageQty,
						"suggested_qty": input.SuggestedQty,
						"source_json":   input.SourceJSON,
						"remark":        input.Remark,
						"gmt_modified":  time.Now(),
					}).Error; err != nil {
					return err
				}
				continue
			}

			row := domain.PackagingProcurementPlan{
				PlanDate:        planDate,
				PackagingItemID: input.PackagingItemID,
				RequiredQty:     input.RequiredQty,
				OnHandQty:       input.OnHandQty,
				ShortageQty:     input.ShortageQty,
				SuggestedQty:    input.SuggestedQty,
				Status:          domain.PackagingProcurementPlanPending,
				SourceJSON:      input.SourceJSON,
				Remark:          input.Remark,
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
			createdCount++
		}

		for _, item := range existing {
			if item.Status != domain.PackagingProcurementPlanPending {
				continue
			}
			if _, ok := seen[item.PackagingItemID]; ok {
				continue
			}
			if err := tx.Delete(&domain.PackagingProcurementPlan{}, item.ID).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	plans, err := r.listPlansByDate(planDate)
	if err != nil {
		return nil, 0, err
	}
	return plans, createdCount, nil
}

func (r *packagingProcurementRepository) ListPlans(params *domain.PackagingProcurementPlanListParams) ([]domain.PackagingProcurementPlan, int64, error) {
	if params == nil {
		params = &domain.PackagingProcurementPlanListParams{}
	}
	planDate := time.Now()
	if params.Date != nil {
		planDate = *params.Date
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}

	query := r.db.Table("packaging_procurement_plan p").
		Where("p.plan_date = ?", planDate.Format("2006-01-02"))
	if params.Status != "" {
		query = query.Where("p.status = ?", params.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []domain.PackagingProcurementPlan
	err := query.Select(`
			p.*,
			pi.item_code AS packaging_item_code,
			pi.item_name AS packaging_item_name,
			pi.unit AS packaging_item_unit,
			ppo.po_number AS packaging_purchase_order_number
		`).
		Joins("LEFT JOIN packaging_item pi ON pi.id = p.packaging_item_id").
		Joins("LEFT JOIN packaging_purchase_order ppo ON ppo.id = p.packaging_purchase_order_id").
		Order("p.shortage_qty DESC, p.id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	if len(rows) == 0 {
		return []domain.PackagingProcurementPlan{}, total, nil
	}
	return rows, total, nil
}

func (r *packagingProcurementRepository) ListRuns(params *domain.PackagingProcurementRunListParams) ([]domain.PackagingProcurementRun, int64, error) {
	if params == nil {
		params = &domain.PackagingProcurementRunListParams{}
	}
	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := r.db.Model(&domain.PackagingProcurementRun{})
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.TriggerType != "" {
		query = query.Where("trigger_type = ?", params.TriggerType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	rows := make([]domain.PackagingProcurementRun, 0)
	if err := query.Order("gmt_create DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *packagingProcurementRepository) listPlansByDate(planDate time.Time) ([]domain.PackagingProcurementPlan, error) {
	var rows []domain.PackagingProcurementPlan
	err := r.db.Table("packaging_procurement_plan p").
		Select(`
			p.*,
			pi.item_code AS packaging_item_code,
			pi.item_name AS packaging_item_name,
			pi.unit AS packaging_item_unit,
			ppo.po_number AS packaging_purchase_order_number
		`).
		Joins("LEFT JOIN packaging_item pi ON pi.id = p.packaging_item_id").
		Joins("LEFT JOIN packaging_purchase_order ppo ON ppo.id = p.packaging_purchase_order_id").
		Where("p.plan_date = ?", planDate.Format("2006-01-02")).
		Order("p.shortage_qty DESC, p.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []domain.PackagingProcurementPlan{}, nil
	}
	return rows, nil
}

func (r *packagingProcurementRepository) ListConvertiblePlans(params *domain.PackagingPlanConvertParams) ([]domain.PackagingProcurementPlan, error) {
	if params == nil {
		params = &domain.PackagingPlanConvertParams{}
	}
	planDate := time.Now()
	if params.Date != nil {
		planDate = *params.Date
	}

	query := r.db.Model(&domain.PackagingProcurementPlan{}).
		Where("plan_date = ? AND status = ? AND suggested_qty > 0", planDate.Format("2006-01-02"), domain.PackagingProcurementPlanPending)
	if len(params.PlanIDs) > 0 {
		query = query.Where("id IN ?", params.PlanIDs)
	}

	var rows []domain.PackagingProcurementPlan
	if err := query.Order("id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []domain.PackagingProcurementPlan{}, nil
	}
	return rows, nil
}

func (r *packagingProcurementRepository) MarkPlansConverted(planIDs []uint64, purchaseOrderID uint64) error {
	if len(planIDs) == 0 || purchaseOrderID == 0 {
		return nil
	}
	now := time.Now()
	return r.db.Model(&domain.PackagingProcurementPlan{}).
		Where("id IN ?", planIDs).
		Updates(map[string]any{
			"status":                      domain.PackagingProcurementPlanConverted,
			"packaging_purchase_order_id": purchaseOrderID,
			"converted_at":                now,
			"gmt_modified":                now,
		}).Error
}

func (r *packagingProcurementRepository) CreateRun(run *domain.PackagingProcurementRun) error {
	if run == nil {
		return nil
	}
	return r.db.Create(run).Error
}

func (r *packagingProcurementRepository) UpdateRun(run *domain.PackagingProcurementRun) error {
	if run == nil {
		return nil
	}
	return r.db.Save(run).Error
}

func (r *packagingProcurementRepository) CreatePurchaseOrder(order *domain.PackagingPurchaseOrder) error {
	if order == nil {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if len(order.Items) == 0 {
			return nil
		}
		items := make([]domain.PackagingPurchaseOrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			item.PackagingPurchaseOrderID = order.ID
			items = append(items, item)
		}
		return tx.Create(&items).Error
	})
}

func (r *packagingProcurementRepository) UpdatePurchaseOrder(order *domain.PackagingPurchaseOrder) error {
	if order == nil {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(order).Error; err != nil {
			return err
		}
		if order.Items == nil {
			return nil
		}

		if err := tx.Where("packaging_purchase_order_id = ?", order.ID).
			Delete(&domain.PackagingPurchaseOrderItem{}).Error; err != nil {
			return err
		}

		if len(order.Items) == 0 {
			return nil
		}
		items := make([]domain.PackagingPurchaseOrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			item.PackagingPurchaseOrderID = order.ID
			items = append(items, item)
		}
		return tx.Create(&items).Error
	})
}

func (r *packagingProcurementRepository) ListPurchaseOrders(params *domain.PackagingPurchaseOrderListParams) ([]domain.PackagingPurchaseOrder, int64, error) {
	if params == nil {
		params = &domain.PackagingPurchaseOrderListParams{}
	}
	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := r.db.Model(&domain.PackagingPurchaseOrder{})
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []domain.PackagingPurchaseOrder
	err := query.Order("gmt_modified DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	if len(orders) == 0 {
		return []domain.PackagingPurchaseOrder{}, total, nil
	}

	orderIDs := make([]uint64, 0, len(orders))
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}
	itemsMap, err := r.listOrderItemsByOrderIDs(orderIDs)
	if err != nil {
		return nil, 0, err
	}
	for i := range orders {
		orders[i].Items = itemsMap[orders[i].ID]
	}
	return orders, total, nil
}

func (r *packagingProcurementRepository) GetPurchaseOrder(id uint64) (*domain.PackagingPurchaseOrder, error) {
	var order domain.PackagingPurchaseOrder
	if err := r.db.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	itemsMap, err := r.listOrderItemsByOrderIDs([]uint64{id})
	if err != nil {
		return nil, err
	}
	order.Items = itemsMap[id]
	return &order, nil
}

type packagingOrderItemRow struct {
	domain.PackagingPurchaseOrderItem
	PackagingItemCode string `gorm:"column:packaging_item_code"`
	PackagingItemName string `gorm:"column:packaging_item_name"`
	PackagingItemUnit string `gorm:"column:packaging_item_unit"`
}

func (r *packagingProcurementRepository) listOrderItemsByOrderIDs(orderIDs []uint64) (map[uint64][]domain.PackagingPurchaseOrderItem, error) {
	result := map[uint64][]domain.PackagingPurchaseOrderItem{}
	if len(orderIDs) == 0 {
		return result, nil
	}

	var rows []packagingOrderItemRow
	err := r.db.Table("packaging_purchase_order_item ppoi").
		Select("ppoi.*, pi.item_code AS packaging_item_code, pi.item_name AS packaging_item_name, pi.unit AS packaging_item_unit").
		Joins("LEFT JOIN packaging_item pi ON pi.id = ppoi.packaging_item_id").
		Where("ppoi.packaging_purchase_order_id IN ?", orderIDs).
		Order("ppoi.id ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		item := row.PackagingPurchaseOrderItem
		item.PackagingItemCode = row.PackagingItemCode
		item.PackagingItemName = row.PackagingItemName
		item.PackagingItemUnit = row.PackagingItemUnit
		result[item.PackagingPurchaseOrderID] = append(result[item.PackagingPurchaseOrderID], item)
	}
	return result, nil
}
