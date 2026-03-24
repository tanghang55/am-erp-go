package repository

import (
	"context"

	"am-erp-go/internal/module/inventory/domain"

	"gorm.io/gorm"
)

type InventoryLotRepository struct {
	db *gorm.DB
}

func NewInventoryLotRepository(db *gorm.DB) domain.InventoryLotRepository {
	return &InventoryLotRepository{db: db}
}

func (r *InventoryLotRepository) List(params *domain.InventoryLotListParams) ([]*domain.InventoryLot, int64, error) {
	var lots []*domain.InventoryLot
	var total int64

	query := r.db.Model(&domain.InventoryLot{})

	if params.ProductID != nil {
		query = query.Where("product_id = ?", *params.ProductID)
	}
	if params.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *params.WarehouseID)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.Keyword != nil && *params.Keyword != "" {
		keyword := "%" + *params.Keyword + "%"
		query = query.
			Joins("LEFT JOIN product ON product.id = inventory_lot.product_id").
			Where("inventory_lot.lot_no LIKE ? OR inventory_lot.source_number LIKE ? OR product.seller_sku LIKE ? OR product.title LIKE ?", keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Order("received_at ASC, id ASC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&lots).Error; err != nil {
		return nil, 0, err
	}

	r.loadAssociations(lots)
	return lots, total, nil
}

func (r *InventoryLotRepository) ListByProductAndWarehouse(ctx context.Context, productID, warehouseID uint64) ([]*domain.InventoryLot, error) {
	var lots []*domain.InventoryLot
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
		Order("received_at ASC, id ASC").
		Find(&lots).Error; err != nil {
		return nil, err
	}
	return lots, nil
}

func (r *InventoryLotRepository) Create(ctx context.Context, lot *domain.InventoryLot) error {
	return r.db.WithContext(ctx).Create(lot).Error
}

func (r *InventoryLotRepository) Update(ctx context.Context, lot *domain.InventoryLot) error {
	return r.db.WithContext(ctx).Save(lot).Error
}

func (r *InventoryLotRepository) loadAssociations(lots []*domain.InventoryLot) {
	if len(lots) == 0 {
		return
	}

	productIDs := make([]uint64, 0, len(lots))
	warehouseIDs := make([]uint64, 0, len(lots))
	for _, lot := range lots {
		productIDs = append(productIDs, lot.ProductID)
		warehouseIDs = append(warehouseIDs, lot.WarehouseID)
	}

	type productInfo struct {
		ID        uint64
		SellerSku string
		Title     string
		Asin      string
	}
	var products []productInfo
	r.db.Table("product").Select("id, seller_sku, title, asin").Where("id IN ?", productIDs).Find(&products)
	productMap := make(map[uint64]*domain.ProductSnapshot, len(products))
	for _, product := range products {
		productMap[product.ID] = &domain.ProductSnapshot{
			ID:        product.ID,
			SellerSku: product.SellerSku,
			Title:     product.Title,
			Asin:      product.Asin,
		}
	}

	type warehouseInfo struct {
		ID   uint64
		Code string
		Name string
	}
	var warehouses []warehouseInfo
	r.db.Table("warehouse").Select("id, code, name").Where("id IN ?", warehouseIDs).Find(&warehouses)
	warehouseMap := make(map[uint64]*domain.WarehouseSnapshot, len(warehouses))
	for _, warehouse := range warehouses {
		warehouseMap[warehouse.ID] = &domain.WarehouseSnapshot{
			ID:   warehouse.ID,
			Code: warehouse.Code,
			Name: warehouse.Name,
		}
	}

	for _, lot := range lots {
		lot.Product = productMap[lot.ProductID]
		lot.Warehouse = warehouseMap[lot.WarehouseID]
	}
}
