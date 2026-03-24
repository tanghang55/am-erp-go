package repository

import (
	"am-erp-go/internal/module/finance/domain"
	"errors"
	"math"
	"time"

	"gorm.io/gorm"
)

type costEventRepository struct {
	db *gorm.DB
}

func NewCostEventRepository(db *gorm.DB) domain.CostEventRepository {
	return &costEventRepository{db: db}
}

func (r *costEventRepository) Create(event *domain.CostEvent) error {
	return r.db.Create(event).Error
}

func (r *costEventRepository) GetLatestPackingMaterialPerUnit(productID uint64, occurredAt time.Time) (*float64, error) {
	var event domain.CostEvent
	err := r.db.
		Where("status = ? AND event_type = ? AND product_id = ? AND qty_event > 0 AND occurred_at <= ?",
			domain.CostEventStatusNormal,
			domain.CostEventTypePackingMaterial,
			productID,
			occurredAt,
		).
		Order("occurred_at DESC, id DESC").
		First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if event.QtyEvent == 0 {
		return nil, nil
	}
	perUnit := math.Round((event.BaseAmount/float64(event.QtyEvent))*1_000_000) / 1_000_000
	return &perUnit, nil
}
