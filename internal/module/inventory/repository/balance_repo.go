package repository

import (
	"context"
	"time"

	"am-erp-go/internal/module/inventory/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoryBalanceRepository struct {
	db *gorm.DB
}

func NewInventoryBalanceRepository(db *gorm.DB) domain.InventoryBalanceRepository {
	return &InventoryBalanceRepository{db: db}
}

func (r *InventoryBalanceRepository) GetOrCreate(ctx context.Context, skuID, warehouseID uint64) (*domain.InventoryBalance, error) {
	var balance domain.InventoryBalance

	err := r.db.WithContext(ctx).
		Where("sku_id = ? AND warehouse_id = ?", skuID, warehouseID).
		First(&balance).Error

	if err == gorm.ErrRecordNotFound {
		// Create new balance record
		now := time.Now()
		balance = domain.InventoryBalance{
			SkuID:             skuID,
			WarehouseID:       warehouseID,
			AvailableQuantity: 0,
			ReservedQuantity:  0,
			DamagedQuantity:   0,
			TotalQuantity:     0,
			GmtCreate:         now,
			GmtModified:       now,
		}

		if err := r.db.WithContext(ctx).Create(&balance).Error; err != nil {
			return nil, err
		}
		return &balance, nil
	}

	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (r *InventoryBalanceRepository) Update(ctx context.Context, balance *domain.InventoryBalance) error {
	return r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Save(balance).Error
}

func (r *InventoryBalanceRepository) List(params *domain.BalanceListParams) ([]*domain.InventoryBalance, int64, error) {
	var balances []*domain.InventoryBalance
	var total int64

	query := r.db.Model(&domain.InventoryBalance{})

	if params.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *params.WarehouseID)
	}

	if params.SkuID != nil {
		query = query.Where("sku_id = ?", *params.SkuID)
	}

	if params.ZeroStock != nil && *params.ZeroStock {
		query = query.Where("total_quantity = 0")
	}

	if params.LowStock != nil && *params.LowStock {
		threshold := uint(10)
		if params.LowStockThreshold != nil {
			threshold = *params.LowStockThreshold
		}
		query = query.Where("total_quantity > 0 AND total_quantity <= ?", threshold)
	}

	if params.Keyword != nil && *params.Keyword != "" {
		// Join with product table for keyword search
		query = query.Joins("LEFT JOIN product ON product.id = inventory_balance.sku_id").
			Where("product.seller_sku LIKE ? OR product.title LIKE ?",
				"%"+*params.Keyword+"%", "%"+*params.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Order("gmt_modified DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&balances).Error; err != nil {
		return nil, 0, err
	}

	// Load SKU and Warehouse info
	r.loadAssociations(balances)

	return balances, total, nil
}

func (r *InventoryBalanceRepository) GetBySkuAndWarehouse(skuID, warehouseID uint64) (*domain.InventoryBalance, error) {
	var balance domain.InventoryBalance

	err := r.db.Where("sku_id = ? AND warehouse_id = ?", skuID, warehouseID).
		First(&balance).Error

	if err != nil {
		return nil, err
	}

	r.loadAssociations([]*domain.InventoryBalance{&balance})
	return &balance, nil
}

func (r *InventoryBalanceRepository) loadAssociations(balances []*domain.InventoryBalance) {
	if len(balances) == 0 {
		return
	}

	// Load SKU info
	skuIDs := make([]uint64, 0)
	for _, b := range balances {
		skuIDs = append(skuIDs, b.SkuID)
	}

	type SkuInfo struct {
		ID        uint64
		SellerSku string
		Title     string
	}
	var skus []SkuInfo
	r.db.Table("product").
		Select("id, seller_sku, title").
		Where("id IN ?", skuIDs).
		Find(&skus)

	skuMap := make(map[uint64]*domain.SkuSnapshot)
	for _, sku := range skus {
		skuMap[sku.ID] = &domain.SkuSnapshot{
			ID:        sku.ID,
			SellerSku: sku.SellerSku,
			Title:     sku.Title,
		}
	}

	// Load Warehouse info
	warehouseIDs := make([]uint64, 0)
	for _, b := range balances {
		warehouseIDs = append(warehouseIDs, b.WarehouseID)
	}

	type WarehouseInfo struct {
		ID   uint64
		Code string
		Name string
	}
	var warehouses []WarehouseInfo
	r.db.Table("warehouse").
		Select("id, code, name").
		Where("id IN ?", warehouseIDs).
		Find(&warehouses)

	warehouseMap := make(map[uint64]*domain.WarehouseSnapshot)
	for _, wh := range warehouses {
		warehouseMap[wh.ID] = &domain.WarehouseSnapshot{
			ID:   wh.ID,
			Code: wh.Code,
			Name: wh.Name,
		}
	}

	// Assign to balances
	for _, b := range balances {
		b.Sku = skuMap[b.SkuID]
		b.Warehouse = warehouseMap[b.WarehouseID]
	}
}
