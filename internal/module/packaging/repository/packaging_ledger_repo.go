package repository

import (
	"am-erp-go/internal/module/packaging/domain"
	"time"

	"gorm.io/gorm"
)

type packagingLedgerRepository struct {
	db *gorm.DB
}

func NewPackagingLedgerRepository(db *gorm.DB) domain.PackagingLedgerRepository {
	return &packagingLedgerRepository{db: db}
}

func (r *packagingLedgerRepository) List(params *domain.PackagingLedgerListParams) ([]domain.PackagingLedger, int64, error) {
	var ledgers []domain.PackagingLedger
	var total int64

	query := r.db.Model(&domain.PackagingLedger{}).Preload("PackagingItem")

	// 包材ID筛选
	if params.PackagingItemID != nil {
		query = query.Where("packaging_item_id = ?", *params.PackagingItemID)
	}

	// 类型筛选
	if params.TransactionType != "" {
		query = query.Where("transaction_type = ?", params.TransactionType)
	}

	// 日期范围筛选
	if params.DateFrom != nil {
		query = query.Where("occurred_at >= ?", *params.DateFrom)
	}
	if params.DateTo != nil {
		query = query.Where("occurred_at <= ?", *params.DateTo)
	}

	// 关联单据筛选
	if params.ReferenceType != "" {
		query = query.Where("reference_type = ?", params.ReferenceType)
	}
	if params.ReferenceID != nil {
		query = query.Where("reference_id = ?", *params.ReferenceID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("gmt_create DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&ledgers).Error; err != nil {
		return nil, 0, err
	}

	return ledgers, total, nil
}

func (r *packagingLedgerRepository) GetByID(id uint64) (*domain.PackagingLedger, error) {
	var ledger domain.PackagingLedger
	if err := r.db.Preload("PackagingItem").First(&ledger, id).Error; err != nil {
		return nil, err
	}
	return &ledger, nil
}

func (r *packagingLedgerRepository) Create(ledger *domain.PackagingLedger) error {
	return r.db.Create(ledger).Error
}

func (r *packagingLedgerRepository) GetUsageSummary(dateFrom, dateTo *time.Time) ([]domain.UsageSummaryItem, error) {
	var summary []domain.UsageSummaryItem

	query := r.db.Table("packaging_ledger pl").
		Select(`pi.id, pi.item_code, pi.item_name, pi.category,
			SUM(CASE WHEN pl.transaction_type = 'IN' THEN ABS(pl.quantity) ELSE 0 END) as total_in,
			SUM(CASE WHEN pl.transaction_type = 'OUT' THEN ABS(pl.quantity) ELSE 0 END) as total_out,
			SUM(CASE WHEN pl.transaction_type = 'OUT' THEN pl.total_cost ELSE 0 END) as total_cost`).
		Joins("INNER JOIN packaging_item pi ON pi.id = pl.packaging_item_id")

	if dateFrom != nil {
		query = query.Where("pl.occurred_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("pl.occurred_at <= ?", *dateTo)
	}

	err := query.Group("pi.id, pi.item_code, pi.item_name, pi.category").
		Order("total_cost DESC").
		Scan(&summary).Error

	return summary, err
}
