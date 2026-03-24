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
	operatorNameMap, err := r.listOperatorNames([]domain.Shipment{shipment})
	if err != nil {
		return nil, err
	}
	shipment.CreatedByName = operatorNameMap[valueOfUint64Ptr(shipment.CreatedBy)]
	shipment.ConfirmedByName = operatorNameMap[valueOfUint64Ptr(shipment.ConfirmedBy)]
	shipment.ShippedByName = operatorNameMap[valueOfUint64Ptr(shipment.ShippedBy)]
	shipment.DeliveredByName = operatorNameMap[valueOfUint64Ptr(shipment.DeliveredBy)]
	shipment.ReceiptCompletedByName = operatorNameMap[valueOfUint64Ptr(shipment.ReceiptCompletedBy)]
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

	operatorNameMap, err := r.listOperatorNames(shipmentsToValues(shipments))
	if err != nil {
		return nil, 0, err
	}
	for i := range shipments {
		shipments[i].CreatedByName = operatorNameMap[valueOfUint64Ptr(shipments[i].CreatedBy)]
		shipments[i].ConfirmedByName = operatorNameMap[valueOfUint64Ptr(shipments[i].ConfirmedBy)]
		shipments[i].ShippedByName = operatorNameMap[valueOfUint64Ptr(shipments[i].ShippedBy)]
		shipments[i].DeliveredByName = operatorNameMap[valueOfUint64Ptr(shipments[i].DeliveredBy)]
		shipments[i].ReceiptCompletedByName = operatorNameMap[valueOfUint64Ptr(shipments[i].ReceiptCompletedBy)]
	}

	return shipments, total, nil
}

func (r *ShipmentRepo) Delete(id uint64) error {
	return r.db.Where("id = ?", id).Delete(&domain.Shipment{}).Error
}

type shipmentOperatorNameRow struct {
	ID   uint64 `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (r *ShipmentRepo) listOperatorNames(shipments []domain.Shipment) (map[uint64]string, error) {
	result := map[uint64]string{}
	if len(shipments) == 0 {
		return result, nil
	}

	ids := make([]uint64, 0, len(shipments)*5)
	seen := map[uint64]struct{}{}
	appendID := func(id *uint64) {
		if id == nil || *id == 0 {
			return
		}
		if _, ok := seen[*id]; ok {
			return
		}
		seen[*id] = struct{}{}
		ids = append(ids, *id)
	}

	for _, shipment := range shipments {
		appendID(shipment.CreatedBy)
		appendID(shipment.ConfirmedBy)
		appendID(shipment.ShippedBy)
		appendID(shipment.DeliveredBy)
		appendID(shipment.ReceiptCompletedBy)
	}

	if len(ids) == 0 {
		return result, nil
	}

	var rows []shipmentOperatorNameRow
	if err := r.db.Table("user").
		Select("id, COALESCE(NULLIF(real_name, ''), username) AS name").
		Where("id IN ?", ids).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.ID] = row.Name
	}
	return result, nil
}

func valueOfUint64Ptr(value *uint64) uint64 {
	if value == nil {
		return 0
	}
	return *value
}

func shipmentsToValues(shipments []*domain.Shipment) []domain.Shipment {
	if len(shipments) == 0 {
		return nil
	}
	values := make([]domain.Shipment, 0, len(shipments))
	for _, shipment := range shipments {
		if shipment == nil {
			continue
		}
		values = append(values, *shipment)
	}
	return values
}
