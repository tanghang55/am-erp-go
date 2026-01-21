package repository

import (
	"am-erp-go/internal/module/supplier/domain"

	"gorm.io/gorm"
)

type supplierAccountRepository struct {
	db *gorm.DB
}

func NewSupplierAccountRepository(db *gorm.DB) domain.SupplierAccountRepository {
	return &supplierAccountRepository{db: db}
}

func (r *supplierAccountRepository) ListBySupplierID(id uint64) ([]domain.SupplierAccount, error) {
	var items []domain.SupplierAccount
	if err := r.db.Where("supplier_id = ?", id).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *supplierAccountRepository) Create(account *domain.SupplierAccount) error {
	return r.db.Create(account).Error
}

func (r *supplierAccountRepository) Update(account *domain.SupplierAccount) error {
	return r.db.Model(&domain.SupplierAccount{}).
		Where("id = ? AND supplier_id = ?", account.ID, account.SupplierID).
		Updates(account).Error
}

func (r *supplierAccountRepository) Delete(id uint64, supplierID uint64) error {
	return r.db.Where("id = ? AND supplier_id = ?", id, supplierID).
		Delete(&domain.SupplierAccount{}).Error
}
