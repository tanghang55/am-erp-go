package repository

import (
	"time"

	"gorm.io/gorm"
)

type comboUsageRepository struct {
	db *gorm.DB
}

func NewComboUsageRepository(db *gorm.DB) *comboUsageRepository {
	return &comboUsageRepository{db: db}
}

func (r *comboUsageRepository) HasBusinessUsageSince(productIDs []uint64, since time.Time) (bool, error) {
	if len(productIDs) == 0 {
		return false, nil
	}

	checks := []struct {
		table      string
		column     string
		timeColumn string
	}{
		{table: "purchase_order_item", column: "product_id", timeColumn: "gmt_create"},
		{table: "sales_order_item", column: "product_id", timeColumn: "gmt_create"},
		{table: "shipment_item", column: "product_id", timeColumn: "gmt_create"},
		{table: "inventory_movement", column: "product_id", timeColumn: "operated_at"},
	}

	for _, check := range checks {
		var count int64
		if err := r.db.Table(check.table).
			Where(check.column+" IN ?", productIDs).
			Where(check.timeColumn+" >= ?", since).
			Limit(1).
			Count(&count).Error; err != nil {
			return false, err
		}
		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}
