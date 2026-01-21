package repository

import (
	"am-erp-go/internal/module/supplier/domain"

	"gorm.io/gorm"
)

type supplierContactRepository struct {
	db *gorm.DB
}

func NewSupplierContactRepository(db *gorm.DB) domain.SupplierContactRepository {
	return &supplierContactRepository{db: db}
}

func (r *supplierContactRepository) ListBySupplierID(id uint64) ([]domain.SupplierContact, error) {
	var items []domain.SupplierContact
	if err := r.db.Where("supplier_id = ?", id).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *supplierContactRepository) Create(contact *domain.SupplierContact) error {
	return r.db.Create(contact).Error
}

func (r *supplierContactRepository) Update(contact *domain.SupplierContact) error {
	return r.db.Model(&domain.SupplierContact{}).
		Where("id = ? AND supplier_id = ?", contact.ID, contact.SupplierID).
		Updates(contact).Error
}

func (r *supplierContactRepository) Delete(id uint64, supplierID uint64) error {
	return r.db.Where("id = ? AND supplier_id = ?", id, supplierID).
		Delete(&domain.SupplierContact{}).Error
}
