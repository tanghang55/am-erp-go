package repository

import (
	"am-erp-go/internal/module/inventory/domain"

	"gorm.io/gorm"
)

type WarehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) domain.WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) GetActiveWarehouses() ([]*domain.Warehouse, error) {
	var warehouses []*domain.Warehouse
	err := r.db.Where("status = ?", domain.WarehouseStatusActive).
		Order("id ASC").
		Find(&warehouses).Error
	return warehouses, err
}
