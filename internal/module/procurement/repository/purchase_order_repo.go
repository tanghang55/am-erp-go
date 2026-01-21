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
	supplierMap, err := r.listSuppliersByIDs(supplierIDs)
	if err != nil {
		return nil, 0, err
	}

	for i := range orders {
		orders[i].Items = itemsMap[orders[i].ID]
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

	if order.SupplierID != nil {
		supplierMap, err := r.listSuppliersByIDs([]uint64{*order.SupplierID})
		if err != nil {
			return nil, err
		}
		if supplier, ok := supplierMap[*order.SupplierID]; ok {
			order.Supplier = &supplier
		}
	}

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

		return tx.Create(&items).Error
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

		return tx.Create(&items).Error
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
	SkuSellerSku string `gorm:"column:sku_seller_sku"`
	SkuTitle     string `gorm:"column:sku_title"`
}

func (r *purchaseOrderRepository) listItemsByOrderIDs(orderIDs []uint64) (map[uint64][]domain.PurchaseOrderItem, error) {
	result := map[uint64][]domain.PurchaseOrderItem{}
	if len(orderIDs) == 0 {
		return result, nil
	}

	var rows []purchaseOrderItemRow
	if err := r.db.Table("purchase_order_item").
		Select("purchase_order_item.*, product.seller_sku AS sku_seller_sku, product.title AS sku_title").
		Joins("LEFT JOIN product ON product.id = purchase_order_item.sku_id").
		Where("purchase_order_item.purchase_order_id IN ?", orderIDs).
		Order("purchase_order_item.id ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		item := row.PurchaseOrderItem
		if row.SkuSellerSku != "" || row.SkuTitle != "" {
			item.Sku = &domain.SkuSnapshot{
				ID:        item.SkuID,
				SellerSku: row.SkuSellerSku,
				Title:     row.SkuTitle,
			}
		}
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
