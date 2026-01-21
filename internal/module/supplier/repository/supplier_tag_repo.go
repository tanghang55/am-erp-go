package repository

import (
	"am-erp-go/internal/module/supplier/domain"

	"gorm.io/gorm"
)

type supplierTagRepository struct {
	db *gorm.DB
}

func NewSupplierTagRepository(db *gorm.DB) domain.SupplierTagRepository {
	return &supplierTagRepository{db: db}
}

func (r *supplierTagRepository) ListBySupplierID(id uint64) ([]domain.SupplierTag, error) {
	var items []domain.SupplierTag
	if err := r.db.Where("supplier_id = ?", id).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *supplierTagRepository) Create(tag *domain.SupplierTag) error {
	return r.db.Create(tag).Error
}

func (r *supplierTagRepository) Update(tag *domain.SupplierTag) error {
	return r.db.Model(&domain.SupplierTag{}).
		Where("id = ? AND supplier_id = ?", tag.ID, tag.SupplierID).
		Updates(tag).Error
}

func (r *supplierTagRepository) Delete(id uint64, supplierID uint64) error {
	return r.db.Where("id = ? AND supplier_id = ?", id, supplierID).
		Delete(&domain.SupplierTag{}).Error
}
