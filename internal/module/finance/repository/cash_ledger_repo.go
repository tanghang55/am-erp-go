package repository

import (
	"am-erp-go/internal/module/finance/domain"
	"time"

	"gorm.io/gorm"
)

type cashLedgerRepository struct {
	db *gorm.DB
}

func NewCashLedgerRepository(db *gorm.DB) domain.CashLedgerRepository {
	return &cashLedgerRepository{db: db}
}

func (r *cashLedgerRepository) List(params *domain.CashLedgerListParams) ([]domain.CashLedger, int64, error) {
	var list []domain.CashLedger
	var total int64

	query := r.buildListQuery(params)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Order("occurred_at DESC, id DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	operatorNameMap, err := r.listOperatorNames(list)
	if err != nil {
		return nil, 0, err
	}
	for i := range list {
		list[i].CreatedByName = operatorNameMap[list[i].CreatedBy]
	}

	return list, total, nil
}

func (r *cashLedgerRepository) GetByID(id uint64) (*domain.CashLedger, error) {
	var entry domain.CashLedger
	if err := r.db.First(&entry, id).Error; err != nil {
		return nil, err
	}
	operatorNameMap, err := r.listOperatorNames([]domain.CashLedger{entry})
	if err != nil {
		return nil, err
	}
	entry.CreatedByName = operatorNameMap[entry.CreatedBy]
	return &entry, nil
}

func (r *cashLedgerRepository) Create(entry *domain.CashLedger) error {
	return r.db.Create(entry).Error
}

func (r *cashLedgerRepository) Update(entry *domain.CashLedger) error {
	return r.db.Save(entry).Error
}

func (r *cashLedgerRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.CashLedger{}, id).Error
}

func (r *cashLedgerRepository) MarkReversed(id uint64, reversedAt time.Time) error {
	result := r.db.Model(&domain.CashLedger{}).
		Where("id = ? AND status = ?", id, domain.CashLedgerStatusNormal).
		Updates(map[string]any{
			"status":       domain.CashLedgerStatusReversed,
			"gmt_modified": reversedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *cashLedgerRepository) GetSummary(params *domain.CashLedgerListParams) (*domain.CashLedgerSummary, error) {
	type summaryRow struct {
		TotalIncome  float64 `gorm:"column:total_income"`
		IncomeCount  int64   `gorm:"column:income_count"`
		TotalExpense float64 `gorm:"column:total_expense"`
		ExpenseCount int64   `gorm:"column:expense_count"`
	}

	var row summaryRow
	query := r.buildListQuery(params).Select(`
		COALESCE(SUM(CASE WHEN ledger_type = 'INCOME' THEN (CASE WHEN base_amount > 0 THEN base_amount ELSE amount END) ELSE 0 END), 0) AS total_income,
		SUM(CASE WHEN ledger_type = 'INCOME' THEN 1 ELSE 0 END) AS income_count,
		COALESCE(SUM(CASE WHEN ledger_type = 'EXPENSE' THEN (CASE WHEN base_amount > 0 THEN base_amount ELSE amount END) ELSE 0 END), 0) AS total_expense,
		SUM(CASE WHEN ledger_type = 'EXPENSE' THEN 1 ELSE 0 END) AS expense_count
	`)

	if err := query.Scan(&row).Error; err != nil {
		return nil, err
	}

	return &domain.CashLedgerSummary{
		TotalIncome:  row.TotalIncome,
		IncomeCount:  row.IncomeCount,
		TotalExpense: row.TotalExpense,
		ExpenseCount: row.ExpenseCount,
		NetProfit:    row.TotalIncome - row.TotalExpense,
	}, nil
}

func (r *cashLedgerRepository) GetSummaryByCategory(params *domain.CashLedgerListParams) ([]domain.CategorySummaryItem, error) {
	var rows []domain.CategorySummaryItem

	query := r.buildListQuery(params).Select(`
		ledger_type,
		category,
		COALESCE(SUM(CASE WHEN base_amount > 0 THEN base_amount ELSE amount END), 0) AS total_amount,
		COUNT(1) AS count
	`)

	if err := query.
		Group("ledger_type, category").
		Order("ledger_type ASC, total_amount DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

type cashLedgerOperatorNameRow struct {
	ID   uint64 `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (r *cashLedgerRepository) listOperatorNames(entries []domain.CashLedger) (map[uint64]string, error) {
	result := map[uint64]string{}
	if len(entries) == 0 {
		return result, nil
	}

	ids := make([]uint64, 0, len(entries))
	seen := map[uint64]struct{}{}
	for _, entry := range entries {
		if entry.CreatedBy == 0 {
			continue
		}
		if _, ok := seen[entry.CreatedBy]; ok {
			continue
		}
		seen[entry.CreatedBy] = struct{}{}
		ids = append(ids, entry.CreatedBy)
	}

	if len(ids) == 0 {
		return result, nil
	}

	var rows []cashLedgerOperatorNameRow
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

func (r *cashLedgerRepository) buildListQuery(params *domain.CashLedgerListParams) *gorm.DB {
	query := r.db.Model(&domain.CashLedger{})
	if params == nil {
		return query
	}
	if params.LedgerType != "" {
		query = query.Where("ledger_type = ?", params.LedgerType)
	}
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}
	if params.Marketplace != "" {
		query = query.Where("marketplace = ?", params.Marketplace)
	}
	if params.OccurredNode != "" {
		query = query.Where("occurred_node = ?", params.OccurredNode)
	}
	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where(
			"description LIKE ? OR trace_id LIKE ? OR category LIKE ? OR COALESCE(marketplace, '') LIKE ? OR COALESCE(occurred_node, '') LIKE ?",
			like, like, like, like, like,
		)
	}
	if params.DateFrom != nil {
		query = query.Where("occurred_at >= ?", *params.DateFrom)
	}
	if params.DateTo != nil {
		query = query.Where("occurred_at <= ?", *params.DateTo)
	}
	if params.ReferenceType != "" {
		query = query.Where("reference_type = ?", params.ReferenceType)
	}
	if params.ReferenceID != nil {
		query = query.Where("reference_id = ?", *params.ReferenceID)
	}
	return query
}
