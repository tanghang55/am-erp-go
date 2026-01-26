package repository

import (
	"context"

	"am-erp-go/internal/module/inventory/domain"

	"gorm.io/gorm"
)

type InventoryMovementRepository struct {
	db *gorm.DB
}

func NewInventoryMovementRepository(db *gorm.DB) domain.InventoryMovementRepository {
	return &InventoryMovementRepository{db: db}
}

func (r *InventoryMovementRepository) Create(ctx context.Context, movement *domain.InventoryMovement) error {
	return r.db.WithContext(ctx).Create(movement).Error
}

func (r *InventoryMovementRepository) List(params *domain.MovementListParams) ([]*domain.InventoryMovement, int64, error) {
	var movements []*domain.InventoryMovement
	var total int64

	query := r.db.Model(&domain.InventoryMovement{})

	if params.SkuID != nil {
		query = query.Where("sku_id = ?", *params.SkuID)
	}

	if params.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *params.WarehouseID)
	}

	if params.MovementType != nil {
		query = query.Where("movement_type = ?", *params.MovementType)
	}

	if params.DateFrom != nil && *params.DateFrom != "" {
		query = query.Where("operated_at >= ?", *params.DateFrom)
	}

	if params.DateTo != nil && *params.DateTo != "" {
		query = query.Where("operated_at <= ?", *params.DateTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Order("operated_at DESC, id DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&movements).Error; err != nil {
		return nil, 0, err
	}

	// Load associations
	r.loadAssociations(movements)

	return movements, total, nil
}

func (r *InventoryMovementRepository) GetByID(id uint64) (*domain.InventoryMovement, error) {
	var movement domain.InventoryMovement

	if err := r.db.First(&movement, id).Error; err != nil {
		return nil, err
	}

	r.loadAssociations([]*domain.InventoryMovement{&movement})
	return &movement, nil
}

func (r *InventoryMovementRepository) loadAssociations(movements []*domain.InventoryMovement) {
	if len(movements) == 0 {
		return
	}

	// Load SKU info
	skuIDs := make([]uint64, 0)
	for _, m := range movements {
		skuIDs = append(skuIDs, m.SkuID)
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
	for _, m := range movements {
		warehouseIDs = append(warehouseIDs, m.WarehouseID)
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

	// Load Operator info
	operatorIDs := make([]uint64, 0)
	for _, m := range movements {
		if m.OperatorID != nil {
			operatorIDs = append(operatorIDs, *m.OperatorID)
		}
	}

	type OperatorInfo struct {
		ID       uint64
		Username string
		RealName *string
	}
	var operators []OperatorInfo
	if len(operatorIDs) > 0 {
		r.db.Table("user").
			Select("id, username, real_name").
			Where("id IN ?", operatorIDs).
			Find(&operators)
	}

	operatorMap := make(map[uint64]*domain.OperatorSnapshot)
	for _, op := range operators {
		operatorMap[op.ID] = &domain.OperatorSnapshot{
			ID:       op.ID,
			Username: op.Username,
			RealName: op.RealName,
		}
	}

	// Assign to movements
	for _, m := range movements {
		m.Sku = skuMap[m.SkuID]
		m.Warehouse = warehouseMap[m.WarehouseID]
		if m.OperatorID != nil {
			m.Operator = operatorMap[*m.OperatorID]
		}
	}
}
