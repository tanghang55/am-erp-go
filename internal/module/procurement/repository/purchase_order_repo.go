package repository

import (
	"am-erp-go/internal/module/procurement/domain"

	"gorm.io/gorm"
)

type purchaseOrderRepository struct {
	db *gorm.DB
}

func NewPurchaseOrderRepository(db *gorm.DB) domain.PurchaseOrderRepository {
	return &purchaseOrderRepository{db: db}
}

func (r *purchaseOrderRepository) List(params *domain.PurchaseOrderListParams) ([]domain.PurchaseOrder, int64, error) {
	var orders []domain.PurchaseOrder
	var total int64

	query := r.db.Model(&domain.PurchaseOrder{})

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.SupplierID != nil {
		query = query.Where("supplier_id = ?", *params.SupplierID)
	}
	if params.Marketplace != "" {
		query = query.Where("marketplace = ?", params.Marketplace)
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("po_number LIKE ? OR remark LIKE ?", keyword, keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("gmt_modified DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	if len(orders) == 0 {
		return []domain.PurchaseOrder{}, total, nil
	}

	orderIDs := make([]uint64, 0, len(orders))
	supplierIDs := make([]uint64, 0, len(orders))
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
		if order.SupplierID != nil {
			supplierIDs = append(supplierIDs, *order.SupplierID)
		}
	}

	itemsMap, err := r.listItemsByOrderIDs(orderIDs)
	if err != nil {
		return nil, 0, err
	}
	operatorNameMap, err := r.listOperatorNames(orders)
	if err != nil {
		return nil, 0, err
	}
	supplierMap, err := r.listSuppliersByIDs(supplierIDs)
	if err != nil {
		return nil, 0, err
	}

	for i := range orders {
		orders[i].Items = itemsMap[orders[i].ID]
		orders[i].QtyPendingInspectionTotal = calculateOrderPendingInspection(orders[i].Items)
		orders[i].OrderedByName = operatorNameMap[valueOfUint64Ptr(orders[i].OrderedBy)]
		orders[i].ShippedByName = operatorNameMap[valueOfUint64Ptr(orders[i].ShippedBy)]
		orders[i].ReceivedByName = operatorNameMap[valueOfUint64Ptr(orders[i].ReceivedBy)]
		orders[i].InspectedByName = operatorNameMap[valueOfUint64Ptr(orders[i].InspectedBy)]
		orders[i].CompletedByName = operatorNameMap[valueOfUint64Ptr(orders[i].CompletedBy)]
		orders[i].ForceCompletedByName = operatorNameMap[valueOfUint64Ptr(orders[i].ForceCompletedBy)]
		if orders[i].SupplierID != nil {
			if supplier, ok := supplierMap[*orders[i].SupplierID]; ok {
				orders[i].Supplier = &supplier
			}
		}
	}

	return orders, total, nil
}

func (r *purchaseOrderRepository) GetByID(id uint64) (*domain.PurchaseOrder, error) {
	var order domain.PurchaseOrder
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}

	itemsMap, err := r.listItemsByOrderIDs([]uint64{id})
	if err != nil {
		return nil, err
	}
	order.Items = itemsMap[id]
	order.QtyPendingInspectionTotal = calculateOrderPendingInspection(order.Items)

	if order.SupplierID != nil {
		supplierMap, err := r.listSuppliersByIDs([]uint64{*order.SupplierID})
		if err != nil {
			return nil, err
		}
		if supplier, ok := supplierMap[*order.SupplierID]; ok {
			order.Supplier = &supplier
		}
	}
	operatorNameMap, err := r.listOperatorNames([]domain.PurchaseOrder{order})
	if err != nil {
		return nil, err
	}
	order.OrderedByName = operatorNameMap[valueOfUint64Ptr(order.OrderedBy)]
	order.ShippedByName = operatorNameMap[valueOfUint64Ptr(order.ShippedBy)]
	order.ReceivedByName = operatorNameMap[valueOfUint64Ptr(order.ReceivedBy)]
	order.InspectedByName = operatorNameMap[valueOfUint64Ptr(order.InspectedBy)]
	order.CompletedByName = operatorNameMap[valueOfUint64Ptr(order.CompletedBy)]
	order.ForceCompletedByName = operatorNameMap[valueOfUint64Ptr(order.ForceCompletedBy)]

	return &order, nil
}

func (r *purchaseOrderRepository) Create(order *domain.PurchaseOrder) error {
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

		items := make([]domain.PurchaseOrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			item.PurchaseOrderID = order.ID
			items = append(items, item)
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		order.Items = items
		return nil
	})
}

func (r *purchaseOrderRepository) Update(order *domain.PurchaseOrder) error {
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

		if err := tx.Where("purchase_order_id = ?", order.ID).
			Delete(&domain.PurchaseOrderItem{}).Error; err != nil {
			return err
		}

		if len(order.Items) == 0 {
			return nil
		}

		items := make([]domain.PurchaseOrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			item.PurchaseOrderID = order.ID
			items = append(items, item)
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		order.Items = items
		return nil
	})
}

func (r *purchaseOrderRepository) UpdateProgress(order *domain.PurchaseOrder) error {
	if order == nil {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		headerUpdates := map[string]any{
			"po_number":             order.PoNumber,
			"batch_no":              order.BatchNo,
			"supplier_id":           order.SupplierID,
			"warehouse_id":          order.WarehouseID,
			"marketplace":           order.Marketplace,
			"status":                order.Status,
			"currency":              order.Currency,
			"total_amount":          order.TotalAmount,
			"ordered_at":            order.OrderedAt,
			"ordered_by":            order.OrderedBy,
			"shipped_at":            order.ShippedAt,
			"shipped_by":            order.ShippedBy,
			"received_at":           order.ReceivedAt,
			"received_by":           order.ReceivedBy,
			"inspected_at":          order.InspectedAt,
			"inspected_by":          order.InspectedBy,
			"closed_at":             order.ClosedAt,
			"completed_by":          order.CompletedBy,
			"is_force_completed":    order.IsForceCompleted,
			"force_completed_at":    order.ForceCompletedAt,
			"force_completed_by":    order.ForceCompletedBy,
			"force_complete_reason": order.ForceCompleteReason,
			"remark":                order.Remark,
			"created_by":            order.CreatedBy,
			"updated_by":            order.UpdatedBy,
			"gmt_modified":          gorm.Expr("CURRENT_TIMESTAMP"),
		}
		if err := tx.Model(&domain.PurchaseOrder{}).Where("id = ?", order.ID).Updates(headerUpdates).Error; err != nil {
			return err
		}
		if order.Items == nil {
			return nil
		}
		for _, item := range order.Items {
			itemUpdates := map[string]any{
				"qty_received":        item.QtyReceived,
				"qty_inspection_pass": item.QtyInspectionPass,
				"qty_inspection_fail": item.QtyInspectionFail,
				"gmt_modified":        gorm.Expr("CURRENT_TIMESTAMP"),
			}
			if err := tx.Model(&domain.PurchaseOrderItem{}).Where("id = ?", item.ID).Updates(itemUpdates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *purchaseOrderRepository) Delete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("purchase_order_id = ?", id).
			Delete(&domain.PurchaseOrderItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&domain.PurchaseOrder{}, id).Error
	})
}

type purchaseOrderItemRow struct {
	domain.PurchaseOrderItem
	ProductSellerSku            string `gorm:"column:product_seller_sku"`
	ProductTitle                string `gorm:"column:product_title"`
	ProductImageURL             string `gorm:"column:product_image_url"`
	ProductIsInspectionRequired uint8  `gorm:"column:product_is_inspection_required"`
	ProductIsPackingRequired    uint8  `gorm:"column:product_is_packing_required"`
}

func (r *purchaseOrderRepository) listItemsByOrderIDs(orderIDs []uint64) (map[uint64][]domain.PurchaseOrderItem, error) {
	result := map[uint64][]domain.PurchaseOrderItem{}
	if len(orderIDs) == 0 {
		return result, nil
	}

	var rows []purchaseOrderItemRow
	if err := r.db.Table("purchase_order_item").
		Select("purchase_order_item.*, product.seller_sku AS product_seller_sku, product.title AS product_title, product.image_url AS product_image_url, product.is_inspection_required AS product_is_inspection_required, product.is_packing_required AS product_is_packing_required").
		Joins("LEFT JOIN product ON product.id = purchase_order_item.product_id").
		Where("purchase_order_item.purchase_order_id IN ?", orderIDs).
		Order("purchase_order_item.id ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		item := row.PurchaseOrderItem
		if row.ProductSellerSku != "" || row.ProductTitle != "" {
			item.Product = &domain.ProductSnapshot{
				ID:                   item.ProductID,
				SellerSku:            row.ProductSellerSku,
				Title:                row.ProductTitle,
				ImageURL:             row.ProductImageURL,
				IsInspectionRequired: row.ProductIsInspectionRequired,
				IsPackingRequired:    row.ProductIsPackingRequired,
			}
		}
		item.QtyPendingInspection = calculateItemPendingInspection(item)
		result[item.PurchaseOrderID] = append(result[item.PurchaseOrderID], item)
	}

	return result, nil
}

type supplierRow struct {
	ID   uint64 `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (r *purchaseOrderRepository) listSuppliersByIDs(ids []uint64) (map[uint64]domain.SupplierSnapshot, error) {
	result := map[uint64]domain.SupplierSnapshot{}
	if len(ids) == 0 {
		return result, nil
	}

	var rows []supplierRow
	if err := r.db.Table("supplier").
		Select("id, name").
		Where("id IN ?", ids).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.ID] = domain.SupplierSnapshot{ID: row.ID, Name: row.Name}
	}

	return result, nil
}

type operatorNameRow struct {
	ID   uint64 `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (r *purchaseOrderRepository) listOperatorNames(orders []domain.PurchaseOrder) (map[uint64]string, error) {
	result := map[uint64]string{}
	if len(orders) == 0 {
		return result, nil
	}
	ids := make([]uint64, 0, len(orders)*6)
	seen := map[uint64]struct{}{}
	appendID := func(id *uint64) {
		if id == nil || *id == 0 {
			return
		}
		if _, ok := seen[*id]; ok {
			return
		}
		seen[*id] = struct{}{}
		ids = append(ids, *id)
	}
	for _, order := range orders {
		appendID(order.OrderedBy)
		appendID(order.ShippedBy)
		appendID(order.ReceivedBy)
		appendID(order.InspectedBy)
		appendID(order.CompletedBy)
		appendID(order.ForceCompletedBy)
	}
	if len(ids) == 0 {
		return result, nil
	}
	var rows []operatorNameRow
	if err := r.db.Table("user").
		Select("id, COALESCE(NULLIF(real_name, ''), username) AS name").
		Where("id IN ?", ids).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.ID] = row.Name
	}
	return result, nil
}

func calculateItemPendingInspection(item domain.PurchaseOrderItem) uint64 {
	inspected := item.QtyInspectionPass + item.QtyInspectionFail
	if item.QtyReceived <= inspected {
		return 0
	}
	return item.QtyReceived - inspected
}

func calculateOrderPendingInspection(items []domain.PurchaseOrderItem) uint64 {
	total := uint64(0)
	for _, item := range items {
		total += calculateItemPendingInspection(item)
	}
	return total
}

func valueOfUint64Ptr(value *uint64) uint64 {
	if value == nil {
		return 0
	}
	return *value
}
