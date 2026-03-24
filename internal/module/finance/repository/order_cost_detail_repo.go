package repository

import (
	"am-erp-go/internal/module/finance/domain"
	"time"

	"gorm.io/gorm"
)

type orderCostDetailRepository struct {
	db *gorm.DB
}

func NewOrderCostDetailRepository(db *gorm.DB) domain.OrderCostDetailRepository {
	return &orderCostDetailRepository{db: db}
}

func (r *orderCostDetailRepository) CreateBatch(details []domain.OrderCostDetail) error {
	if len(details) == 0 {
		return nil
	}
	return r.db.Create(&details).Error
}

func (r *orderCostDetailRepository) SumQtyByDateAndMarketplace(bizDate time.Time, marketplace *string) (uint64, error) {
	type row struct {
		Qty uint64 `gorm:"column:qty"`
	}
	result := row{}
	query := r.db.Table("finance_order_cost_detail").
		Select("COALESCE(SUM(qty_out), 0) AS qty").
		Where("DATE(occurred_at) = ?", bizDate.Format("2006-01-02")).
		Where("status = ?", domain.OrderCostDetailStatusNormal).
		Where("reversal_of_id IS NULL")
	if marketplace != nil && *marketplace != "" {
		query = query.Where("marketplace = ?", *marketplace)
	}
	if err := query.Scan(&result).Error; err != nil {
		return 0, err
	}
	return result.Qty, nil
}

func (r *orderCostDetailRepository) ListReturnableBySalesOrderItemID(salesOrderItemID uint64) ([]domain.ReturnableOrderCostDetail, error) {
	rows := make([]domain.ReturnableOrderCostDetail, 0)
	if salesOrderItemID == 0 {
		return rows, nil
	}

	err := r.db.Table("finance_order_cost_detail AS src").
		Select(`
			src.*,
			GREATEST(
				src.qty_out - COALESCE(SUM(CASE WHEN rv.status = 'NORMAL' THEN rv.qty_out ELSE 0 END), 0),
				0
			) AS available_qty
		`).
		Joins("LEFT JOIN finance_order_cost_detail rv ON rv.reversal_of_id = src.id").
		Where("src.sales_order_item_id = ?", salesOrderItemID).
		Where("src.status = ?", domain.OrderCostDetailStatusNormal).
		Where("src.reversal_of_id IS NULL").
		Group("src.id").
		Having("available_qty > 0").
		Order("src.occurred_at ASC, src.id ASC").
		Scan(&rows).Error
	return rows, err
}
