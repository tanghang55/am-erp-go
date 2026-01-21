package repository

import (
	"am-erp-go/internal/module/supplier/domain"

	"gorm.io/gorm"
)

type supplierTypeRepository struct {
	db *gorm.DB
}

func NewSupplierTypeRepository(db *gorm.DB) domain.SupplierTypeRepository {
	return &supplierTypeRepository{db: db}
}

func (r *supplierTypeRepository) ListBySupplierID(id uint64) ([]string, error) {
	var items []domain.SupplierType
	if err := r.db.Where("supplier_id = ?", id).Find(&items).Error; err != nil {
		return nil, err
	}
	types := make([]string, 0, len(items))
	for _, item := range items {
		types = append(types, item.Type)
	}
	return types, nil
}

func (r *supplierTypeRepository) ListBySupplierIDs(ids []uint64) (map[uint64][]string, error) {
	if len(ids) == 0 {
		return map[uint64][]string{}, nil
	}
	var items []domain.SupplierType
	if err := r.db.Where("supplier_id IN ?", ids).Find(&items).Error; err != nil {
		return nil, err
	}
	types := make(map[uint64][]string)
	for _, item := range items {
		types[item.SupplierID] = append(types[item.SupplierID], item.Type)
	}
	return types, nil
}

func (r *supplierTypeRepository) ReplaceBySupplierID(id uint64, types []string) error {
	if err := r.db.Where("supplier_id = ?", id).Delete(&domain.SupplierType{}).Error; err != nil {
		return err
	}
	if len(types) == 0 {
		return nil
	}
	rows := make([]domain.SupplierType, 0, len(types))
	for _, item := range types {
		rows = append(rows, domain.SupplierType{SupplierID: id, Type: item})
	}
	return r.db.Create(&rows).Error
}
