package repository

import (
	"am-erp-go/internal/module/shipment/domain"

	"gorm.io/gorm"
)

type ShipmentItemRepo struct {
	db *gorm.DB
}

func NewShipmentItemRepo(db *gorm.DB) *ShipmentItemRepo {
	return &ShipmentItemRepo{db: db}
}

func (r *ShipmentItemRepo) Create(item *domain.ShipmentItem) error {
	return r.db.Create(item).Error
}

func (r *ShipmentItemRepo) CreateBatch(items []domain.ShipmentItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}

func (r *ShipmentItemRepo) GetByShipmentID(shipmentID uint64) ([]domain.ShipmentItem, error) {
	var items []domain.ShipmentItem
	err := r.db.Where("shipment_id = ?", shipmentID).Find(&items).Error
	return items, err
}

func (r *ShipmentItemRepo) DeleteByShipmentID(shipmentID uint64) error {
	return r.db.Where("shipment_id = ?", shipmentID).Delete(&domain.ShipmentItem{}).Error
}
