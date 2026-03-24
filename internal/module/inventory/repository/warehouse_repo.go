package repository

import (
	"context"
	"strings"

	"am-erp-go/internal/module/inventory/domain"

	"gorm.io/gorm"
)

type WarehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) domain.WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) List(params *domain.WarehouseListParams) ([]*domain.Warehouse, int64, error) {
	var warehouses []*domain.Warehouse
	var total int64

	query := r.db.Model(&domain.Warehouse{})

	// 过滤条件
	if params.Type != nil {
		query = query.Where("type = ?", *params.Type)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.Keyword != nil && *params.Keyword != "" {
		keyword := "%" + *params.Keyword + "%"
		query = query.Where("code LIKE ? OR name LIKE ?", keyword, keyword)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(params.PageSize).Find(&warehouses).Error; err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

func (r *WarehouseRepository) GetByID(id uint64) (*domain.Warehouse, error) {
	var warehouse domain.Warehouse
	if err := r.db.Where("id = ?", id).First(&warehouse).Error; err != nil {
		return nil, err
	}
	return &warehouse, nil
}

func (r *WarehouseRepository) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	return r.db.WithContext(ctx).Create(warehouse).Error
}

func (r *WarehouseRepository) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	return r.db.WithContext(ctx).Save(warehouse).Error
}

func (r *WarehouseRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Warehouse{}, id).Error
}

func (r *WarehouseRepository) GetActiveWarehouses() ([]*domain.Warehouse, error) {
	var warehouses []*domain.Warehouse
	err := r.db.Where("status = ?", domain.WarehouseStatusActive).
		Order("id ASC").
		Find(&warehouses).Error
	return warehouses, err
}

func (r *WarehouseRepository) CountReferences(id uint64) (int64, error) {
	var total int64
	queries := []string{
		"SELECT COUNT(1) FROM inventory_balance WHERE warehouse_id = ?",
		"SELECT COUNT(1) FROM inventory_lot WHERE warehouse_id = ?",
		"SELECT COUNT(1) FROM inventory_movement WHERE warehouse_id = ?",
		"SELECT COUNT(1) FROM shipping_rate WHERE origin_warehouse_id = ? OR destination_warehouse_id = ?",
		"SELECT COUNT(1) FROM shipment WHERE warehouse_id = ? OR destination_warehouse_id = ?",
		"SELECT COUNT(1) FROM procurement_replenishment_strategy WHERE warehouse_id = ?",
		"SELECT COUNT(1) FROM procurement_replenishment_plan WHERE warehouse_id = ?",
	}

	for _, query := range queries {
		var count int64
		var err error
		if strings.Contains(query, "OR") {
			err = r.db.Raw(query, id, id).Scan(&count).Error
		} else {
			err = r.db.Raw(query, id).Scan(&count).Error
		}
		if err != nil {
			return 0, err
		}
		total += count
	}

	return total, nil
}
