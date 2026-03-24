package repository

import (
	"fmt"
	"strings"

	"am-erp-go/internal/module/finance/domain"

	"gorm.io/gorm"
)

type profitQueryRepository struct {
	db *gorm.DB
}

func profitSummaryBaseCurrencyExpr() string {
	return "NULLIF(MAX(src.base_currency), '') AS base_currency"
}

func NewProfitQueryRepository(db *gorm.DB) domain.ProfitQueryRepository {
	return &profitQueryRepository{db: db}
}

func (r *profitQueryRepository) ListOrderProfits(params *domain.OrderProfitListParams) ([]domain.OrderProfitSummary, int64, error) {
	whereSQL, args := buildProfitLedgerFilters(params)
	orderExpr := "CASE WHEN pl.sales_order_id IS NOT NULL THEN pl.sales_order_id WHEN pl.reference_type = 'SALES_ORDER' THEN pl.reference_id ELSE NULL END"

	countSQL := fmt.Sprintf(
		"SELECT COUNT(1) AS total FROM (SELECT %s AS order_id FROM finance_profit_ledger pl WHERE %s GROUP BY order_id) t",
		orderExpr, whereSQL,
	)
	var total int64
	if err := r.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	listSQL := fmt.Sprintf(`
SELECT
	CAST(src.order_id AS UNSIGNED) AS sales_order_id,
	COALESCE(MAX(so.order_no), '') AS order_no,
	COALESCE(MAX(src.marketplace), 'ALL') AS marketplace,
	%s,
	SUM(CASE WHEN src.ledger_type = 'INCOME' THEN src.signed_base_amount ELSE 0 END) AS sales_income_amount,
	SUM(CASE WHEN src.ledger_type = 'COGS' THEN src.signed_base_amount ELSE 0 END) AS cogs_amount,
	SUM(CASE WHEN src.ledger_type = 'INCOME' THEN src.signed_base_amount ELSE 0 END) - SUM(CASE WHEN src.ledger_type = 'COGS' THEN src.signed_base_amount ELSE 0 END) AS gross_profit_amount,
	SUM(CASE WHEN src.ledger_type = 'ORDER_EXPENSE' THEN src.signed_base_amount ELSE 0 END) AS order_expense_amount,
	(SUM(CASE WHEN src.ledger_type = 'INCOME' THEN src.signed_base_amount ELSE 0 END) - SUM(CASE WHEN src.ledger_type = 'COGS' THEN src.signed_base_amount ELSE 0 END) - SUM(CASE WHEN src.ledger_type = 'ORDER_EXPENSE' THEN src.signed_base_amount ELSE 0 END)) AS order_net_profit_amount,
	MAX(src.occurred_at) AS occurred_at
FROM (
	SELECT
		%s AS order_id,
		pl.ledger_type,
		CASE WHEN pl.reversal_of_id IS NULL THEN pl.base_amount ELSE -pl.base_amount END AS signed_base_amount,
		pl.marketplace,
		pl.base_currency,
		pl.occurred_at
	FROM finance_profit_ledger pl
	WHERE %s
) src
LEFT JOIN sales_order so ON so.id = src.order_id
GROUP BY src.order_id
ORDER BY MAX(src.occurred_at) DESC
LIMIT ? OFFSET ?`, profitSummaryBaseCurrencyExpr(), orderExpr, whereSQL)

	listArgs := append([]interface{}{}, args...)
	listArgs = append(listArgs, params.PageSize, offset)
	items := make([]domain.OrderProfitSummary, 0)
	if err := r.db.Raw(listSQL, listArgs...).Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *profitQueryRepository) GetOrderProfitDetail(salesOrderID uint64) (*domain.OrderProfitDetail, error) {
	if salesOrderID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	orderExpr := "CASE WHEN pl.sales_order_id IS NOT NULL THEN pl.sales_order_id WHEN pl.reference_type = 'SALES_ORDER' THEN pl.reference_id ELSE NULL END"

	summarySQL := fmt.Sprintf(`
SELECT
	CAST(src.order_id AS UNSIGNED) AS sales_order_id,
	COALESCE(MAX(so.order_no), '') AS order_no,
	COALESCE(MAX(src.marketplace), 'ALL') AS marketplace,
	%s,
	SUM(CASE WHEN src.ledger_type = 'INCOME' THEN src.signed_base_amount ELSE 0 END) AS sales_income_amount,
	SUM(CASE WHEN src.ledger_type = 'COGS' THEN src.signed_base_amount ELSE 0 END) AS cogs_amount,
	SUM(CASE WHEN src.ledger_type = 'INCOME' THEN src.signed_base_amount ELSE 0 END) - SUM(CASE WHEN src.ledger_type = 'COGS' THEN src.signed_base_amount ELSE 0 END) AS gross_profit_amount,
	SUM(CASE WHEN src.ledger_type = 'ORDER_EXPENSE' THEN src.signed_base_amount ELSE 0 END) AS order_expense_amount,
	(SUM(CASE WHEN src.ledger_type = 'INCOME' THEN src.signed_base_amount ELSE 0 END) - SUM(CASE WHEN src.ledger_type = 'COGS' THEN src.signed_base_amount ELSE 0 END) - SUM(CASE WHEN src.ledger_type = 'ORDER_EXPENSE' THEN src.signed_base_amount ELSE 0 END)) AS order_net_profit_amount,
	MAX(src.occurred_at) AS occurred_at
FROM (
	SELECT
		%s AS order_id,
		pl.ledger_type,
		CASE WHEN pl.reversal_of_id IS NULL THEN pl.base_amount ELSE -pl.base_amount END AS signed_base_amount,
		pl.marketplace,
		pl.base_currency,
		pl.occurred_at
	FROM finance_profit_ledger pl
	WHERE pl.status = 'NORMAL'
		AND (pl.sales_order_id IS NOT NULL OR (pl.reference_type = 'SALES_ORDER' AND pl.reference_id IS NOT NULL))
) src
LEFT JOIN sales_order so ON so.id = src.order_id
WHERE src.order_id = ?
GROUP BY src.order_id`, profitSummaryBaseCurrencyExpr(), orderExpr)

	summary := domain.OrderProfitSummary{}
	if err := r.db.Raw(summarySQL, salesOrderID).Scan(&summary).Error; err != nil {
		return nil, err
	}
	if summary.SalesOrderID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	lines := make([]domain.OrderProfitLine, 0)
	lineSQL := `
SELECT
	pl.sales_order_item_id,
	COALESCE(MAX(soi.product_id), 0) AS product_id,
	COALESCE(MAX(p.seller_sku), '') AS seller_sku,
	COALESCE(MAX(p.title), '') AS product_title,
	COALESCE(MAX(p.image_url), '') AS product_image_url,
	COALESCE(MAX(soi.qty_shipped), 0) AS qty_shipped,
	COALESCE(MAX(soi.unit_price), 0) AS unit_price,
	SUM(CASE WHEN pl.ledger_type = 'INCOME' THEN CASE WHEN pl.reversal_of_id IS NULL THEN pl.base_amount ELSE -pl.base_amount END ELSE 0 END) AS income_amount,
	SUM(CASE WHEN pl.ledger_type = 'COGS' THEN CASE WHEN pl.reversal_of_id IS NULL THEN pl.base_amount ELSE -pl.base_amount END ELSE 0 END) AS cogs_amount,
	SUM(CASE WHEN pl.ledger_type = 'INCOME' THEN CASE WHEN pl.reversal_of_id IS NULL THEN pl.base_amount ELSE -pl.base_amount END ELSE 0 END) - SUM(CASE WHEN pl.ledger_type = 'COGS' THEN CASE WHEN pl.reversal_of_id IS NULL THEN pl.base_amount ELSE -pl.base_amount END ELSE 0 END) AS gross_profit_amount
FROM finance_profit_ledger pl
LEFT JOIN sales_order_item soi ON soi.id = pl.sales_order_item_id
LEFT JOIN product p ON p.id = soi.product_id
WHERE pl.status = 'NORMAL'
	AND pl.sales_order_id = ?
	AND pl.sales_order_item_id IS NOT NULL
	AND pl.ledger_type IN ('INCOME','COGS')
GROUP BY pl.sales_order_item_id
ORDER BY pl.sales_order_item_id ASC`
	if err := r.db.Raw(lineSQL, salesOrderID).Scan(&lines).Error; err != nil {
		return nil, err
	}

	expenses := make([]domain.OrderProfitExpense, 0)
	expenseSQL := `
SELECT
	id,
	category,
	base_currency,
	base_amount,
	occurred_at,
	remark
FROM finance_profit_ledger
WHERE status = 'NORMAL'
	AND ledger_type = 'ORDER_EXPENSE'
	AND (sales_order_id = ? OR (reference_type = 'SALES_ORDER' AND reference_id = ?))
ORDER BY occurred_at DESC, id DESC`
	if err := r.db.Raw(expenseSQL, salesOrderID, salesOrderID).Scan(&expenses).Error; err != nil {
		return nil, err
	}

	return &domain.OrderProfitDetail{
		Summary:  summary,
		Lines:    lines,
		Expenses: expenses,
	}, nil
}

func buildProfitLedgerFilters(params *domain.OrderProfitListParams) (string, []interface{}) {
	parts := []string{
		"pl.status = 'NORMAL'",
		"(pl.sales_order_id IS NOT NULL OR (pl.reference_type = 'SALES_ORDER' AND pl.reference_id IS NOT NULL))",
	}
	args := make([]interface{}, 0)
	if params != nil {
		if params.DateFrom != nil {
			parts = append(parts, "pl.biz_date >= ?")
			args = append(args, params.DateFrom.Format("2006-01-02"))
		}
		if params.DateTo != nil {
			parts = append(parts, "pl.biz_date <= ?")
			args = append(args, params.DateTo.Format("2006-01-02"))
		}
		if strings.TrimSpace(params.Marketplace) != "" {
			parts = append(parts, "pl.marketplace = ?")
			args = append(args, strings.TrimSpace(params.Marketplace))
		}
		if strings.TrimSpace(params.Keyword) != "" {
			parts = append(parts, `EXISTS (
				SELECT 1
				FROM sales_order so
				WHERE so.id = CASE
					WHEN pl.sales_order_id IS NOT NULL THEN pl.sales_order_id
					WHEN pl.reference_type = 'SALES_ORDER' THEN pl.reference_id
					ELSE NULL
				END
				AND so.order_no LIKE ?
			)`)
			args = append(args, "%"+strings.TrimSpace(params.Keyword)+"%")
		}
	}
	return strings.Join(parts, " AND "), args
}
