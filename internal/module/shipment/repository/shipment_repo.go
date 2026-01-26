package repository

import (
	"am-erp-go/internal/module/shipment/domain"

	"gorm.io/gorm"
)

type ShipmentRepo struct {
	db *gorm.DB
}

func NewShipmentRepo(db *gorm.DB) *ShipmentRepo {
	return &ShipmentRepo{db: db}
}

func (r *ShipmentRepo) Create(shipment *domain.Shipment) error {
	return r.db.Create(shipment).Error
}

func (r *ShipmentRepo) Update(shipment *domain.Shipment) error {
	return r.db.Save(shipment).Error
}

func (r *ShipmentRepo) GetByID(id uint64) (*domain.Shipment, error) {
	var shipment domain.Shipment
	err := r.db.Where("id = ?", id).First(&shipment).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *ShipmentRepo) GetByShipmentNumber(shipmentNumber string) (*domain.Shipment, error) {
	var shipment domain.Shipment
	err := r.db.Where("shipment_number = ?", shipmentNumber).First(&shipment).Error
	if err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (r *ShipmentRepo) List(params *domain.ShipmentListParams) ([]*domain.Shipment, int64, error) {
	var shipments []*domain.Shipment
	var total int64

	query := r.db.Model(&domain.Shipment{})

	// Filters
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *params.WarehouseID)
	}
	if params.OrderNumber != nil && *params.OrderNumber != "" {
		query = query.Where("order_number = ?", *params.OrderNumber)
	}
	if params.TrackingNumber != nil && *params.TrackingNumber != "" {
		query = query.Where("tracking_number = ?", *params.TrackingNumber)
	}
	if params.Keyword != nil && *params.Keyword != "" {
		query = query.Where("shipment_number LIKE ? OR order_number LIKE ? OR tracking_number LIKE ?",
			"%"+*params.Keyword+"%", "%"+*params.Keyword+"%", "%"+*params.Keyword+"%")
	}
	if params.DateFrom != nil && *params.DateFrom != "" {
		query = query.Where("gmt_create >= ?", *params.DateFrom)
	}
	if params.DateTo != nil && *params.DateTo != "" {
		query = query.Where("gmt_create <= ?", *params.DateTo)
	}

	// Count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Query
	err := query.Order("gmt_create DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&shipments).Error

	if err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *ShipmentRepo) Delete(id uint64) error {
	return r.db.Where("id = ?", id).Delete(&domain.Shipment{}).Error
}
