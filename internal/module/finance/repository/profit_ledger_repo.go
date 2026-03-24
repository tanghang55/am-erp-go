package repository

import (
	"time"

	"am-erp-go/internal/module/finance/domain"

	"gorm.io/gorm"
)

type profitLedgerRepository struct {
	db *gorm.DB
}

func NewProfitLedgerRepository(db *gorm.DB) domain.ProfitLedgerRepository {
	return &profitLedgerRepository{db: db}
}

func (r *profitLedgerRepository) Create(entry *domain.ProfitLedger) error {
	return r.db.Create(entry).Error
}

func (r *profitLedgerRepository) CreateBatch(entries []domain.ProfitLedger) error {
	if len(entries) == 0 {
		return nil
	}
	return r.db.Create(&entries).Error
}

func (r *profitLedgerRepository) AggregateDaily(bizDate time.Time, marketplace *string) ([]domain.ProfitLedgerDailyAgg, error) {
	type row struct {
		Marketplace   string  `gorm:"column:marketplace"`
		BaseCurrency  string  `gorm:"column:base_currency"`
		SalesIncome   float64 `gorm:"column:sales_income"`
		COGS          float64 `gorm:"column:cogs"`
		OrderExpense  float64 `gorm:"column:order_expense"`
		PublicExpense float64 `gorm:"column:public_expense"`
		OrderCount    uint64  `gorm:"column:order_count"`
	}

	rows := make([]row, 0)
	query := r.db.Table("finance_profit_ledger").
		Select(`
			COALESCE(marketplace, 'ALL') AS marketplace,
			base_currency,
			SUM(CASE WHEN ledger_type = 'INCOME' AND status = 'NORMAL' THEN CASE WHEN reversal_of_id IS NULL THEN base_amount ELSE -base_amount END ELSE 0 END) AS sales_income,
			SUM(CASE WHEN ledger_type = 'COGS' AND status = 'NORMAL' THEN CASE WHEN reversal_of_id IS NULL THEN base_amount ELSE -base_amount END ELSE 0 END) AS cogs,
			SUM(CASE WHEN ledger_type = 'ORDER_EXPENSE' AND status = 'NORMAL' THEN CASE WHEN reversal_of_id IS NULL THEN base_amount ELSE -base_amount END ELSE 0 END) AS order_expense,
			SUM(CASE WHEN ledger_type = 'PUBLIC_EXPENSE' AND status = 'NORMAL' THEN CASE WHEN reversal_of_id IS NULL THEN base_amount ELSE -base_amount END ELSE 0 END) AS public_expense,
			COUNT(DISTINCT CASE WHEN ledger_type = 'INCOME' AND status = 'NORMAL' AND reversal_of_id IS NULL THEN sales_order_id END) AS order_count
		`).
		Where("biz_date = ?", bizDate.Format("2006-01-02")).
		Group("COALESCE(marketplace, 'ALL'), base_currency")

	if marketplace != nil && *marketplace != "" {
		query = query.Where("marketplace = ?", *marketplace)
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]domain.ProfitLedgerDailyAgg, 0, len(rows))
	for _, item := range rows {
		result = append(result, domain.ProfitLedgerDailyAgg{
			Marketplace:   item.Marketplace,
			BaseCurrency:  item.BaseCurrency,
			SalesIncome:   item.SalesIncome,
			COGS:          item.COGS,
			OrderExpense:  item.OrderExpense,
			PublicExpense: item.PublicExpense,
			OrderCount:    item.OrderCount,
			ShippedQty:    0,
		})
	}
	return result, nil
}
