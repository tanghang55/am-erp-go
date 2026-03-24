package repository

import (
	"fmt"
	"strings"

	"am-erp-go/internal/module/finance/domain"

	"gorm.io/gorm"
)

type productCostRepository struct {
	db *gorm.DB
}

func productCostSummaryBaseCurrencyExpr() string {
	return "NULLIF(MAX(base_currency), '') AS base_currency"
}

func NewProductCostRepository(db *gorm.DB) domain.ProductCostRepository {
	return &productCostRepository{db: db}
}

func (r *productCostRepository) ListLedger(params *domain.ProductCostLedgerListParams) ([]domain.ProductCostLedgerItem, int64, error) {
	whereSQL, args := buildProductCostWhere(params)
	unionSQL := productCostUnionSQL()

	countSQL := fmt.Sprintf("SELECT COUNT(1) AS total FROM (%s) t WHERE %s", unionSQL, whereSQL)
	var total int64
	if err := r.db.Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	listSQL := fmt.Sprintf(`
SELECT
	direction,
	source_type,
	occurred_at,
	product_id,
	seller_sku,
	product_title,
	warehouse_id,
	warehouse_code,
	warehouse_name,
	marketplace,
	quantity,
	unit_cost_original,
	original_currency,
	original_amount,
	base_currency,
	base_amount,
	reference_type,
	reference_id,
	reference_number
FROM (%s) t
WHERE %s
ORDER BY occurred_at DESC, product_id ASC
LIMIT ? OFFSET ?`, unionSQL, whereSQL)

	listArgs := append([]interface{}{}, args...)
	listArgs = append(listArgs, params.PageSize, offset)
	items := make([]domain.ProductCostLedgerItem, 0)
	if err := r.db.Raw(listSQL, listArgs...).Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *productCostRepository) GetSummary(params *domain.ProductCostLedgerListParams) (*domain.ProductCostSummary, error) {
	whereSQL, args := buildProductCostWhere(params)
	unionSQL := productCostUnionSQL()
	sql := productCostSummarySQL(whereSQL, unionSQL)

	summary := &domain.ProductCostSummary{}
	if err := r.db.Raw(sql, args...).Scan(summary).Error; err != nil {
		return nil, err
	}
	return summary, nil
}

func productCostSummarySQL(whereSQL string, unionSQL string) string {
	return fmt.Sprintf(`
SELECT
	%s,
	COALESCE(SUM(CASE WHEN direction = 'INBOUND' THEN quantity ELSE 0 END), 0) AS inbound_qty,
	COALESCE(SUM(CASE WHEN direction = 'INBOUND' THEN base_amount ELSE 0 END), 0) AS inbound_amount,
	COALESCE(SUM(CASE WHEN direction = 'OUTBOUND' THEN quantity ELSE 0 END), 0) AS outbound_qty,
	COALESCE(SUM(CASE WHEN direction = 'OUTBOUND' THEN base_amount ELSE 0 END), 0) AS outbound_amount,
	COALESCE(SUM(CASE WHEN direction = 'NEUTRAL' THEN base_amount ELSE 0 END), 0) AS neutral_amount,
	COALESCE(SUM(CASE WHEN direction = 'INBOUND' THEN quantity ELSE 0 END), 0) - COALESCE(SUM(CASE WHEN direction = 'OUTBOUND' THEN quantity ELSE 0 END), 0) AS net_qty,
	COALESCE(SUM(CASE WHEN direction = 'INBOUND' THEN base_amount ELSE 0 END), 0) - COALESCE(SUM(CASE WHEN direction = 'OUTBOUND' THEN base_amount ELSE 0 END), 0) - COALESCE(SUM(CASE WHEN direction = 'NEUTRAL' THEN base_amount ELSE 0 END), 0) AS net_amount,
	CASE
		WHEN COALESCE(SUM(CASE WHEN direction = 'INBOUND' THEN quantity ELSE 0 END), 0) = 0 THEN 0
		ELSE COALESCE(SUM(CASE WHEN direction = 'INBOUND' THEN base_amount ELSE 0 END), 0) / COALESCE(SUM(CASE WHEN direction = 'INBOUND' THEN quantity ELSE 0 END), 0)
	END AS avg_inbound_unit_cost,
	CASE
		WHEN COALESCE(SUM(CASE WHEN direction = 'OUTBOUND' THEN quantity ELSE 0 END), 0) = 0 THEN 0
		ELSE COALESCE(SUM(CASE WHEN direction = 'OUTBOUND' THEN base_amount ELSE 0 END), 0) / COALESCE(SUM(CASE WHEN direction = 'OUTBOUND' THEN quantity ELSE 0 END), 0)
	END AS avg_outbound_unit_cost
FROM (%s) t
WHERE %s`, productCostSummaryBaseCurrencyExpr(), unionSQL, whereSQL)
}

func productCostUnionSQL() string {
	return `
SELECT
	CASE
		WHEN ce.event_type IN ('SHIPMENT_ALLOCATED', 'PACKING_MATERIAL') THEN 'NEUTRAL'
		ELSE 'INBOUND'
	END AS direction,
	ce.event_type AS source_type,
	ce.occurred_at AS occurred_at,
	ce.product_id AS product_id,
	p.seller_sku AS seller_sku,
	p.title AS product_title,
	ce.warehouse_id AS warehouse_id,
	w.code AS warehouse_code,
	w.name AS warehouse_name,
	ce.marketplace AS marketplace,
	CASE WHEN ce.event_type IN ('SHIPMENT_ALLOCATED', 'PACKING_MATERIAL') THEN 0 ELSE ce.qty_event END AS quantity,
	CASE
		WHEN ce.event_type IN ('SHIPMENT_ALLOCATED', 'PACKING_MATERIAL') THEN 0
		WHEN ce.qty_event > 0 THEN ROUND(ce.original_amount / ce.qty_event, 6)
		ELSE 0
	END AS unit_cost_original,
	ce.original_currency AS original_currency,
	ce.original_amount AS original_amount,
	ce.base_currency AS base_currency,
	ce.base_amount AS base_amount,
	CASE
		WHEN ce.event_type = 'SHIPMENT_ALLOCATED' THEN 'SHIPMENT'
		WHEN ce.event_type = 'PACKING_MATERIAL' THEN COALESCE(im.reference_type, 'PRODUCT_PACKING')
		ELSE 'PURCHASE_ORDER'
	END AS reference_type,
	CASE
		WHEN ce.event_type = 'SHIPMENT_ALLOCATED' THEN ce.shipment_id
		WHEN ce.event_type = 'PACKING_MATERIAL' THEN ce.inventory_movement_id
		ELSE ce.purchase_order_id
	END AS reference_id,
	CASE
		WHEN ce.event_type = 'SHIPMENT_ALLOCATED' THEN sh.shipment_number
		WHEN ce.event_type = 'PACKING_MATERIAL' THEN im.reference_number
		ELSE po.po_number
	END AS reference_number
FROM finance_cost_event ce
LEFT JOIN product p ON p.id = ce.product_id
LEFT JOIN warehouse w ON w.id = ce.warehouse_id
LEFT JOIN purchase_order po ON po.id = ce.purchase_order_id
LEFT JOIN shipment sh ON sh.id = ce.shipment_id
LEFT JOIN inventory_movement im ON im.id = ce.inventory_movement_id
WHERE ce.status = 'NORMAL' AND ce.event_type IN ('PO_RECEIVED', 'PO_ADJUST', 'SHIPMENT_ALLOCATED', 'PACKING_MATERIAL')
UNION ALL
SELECT
	CASE WHEN od.reversal_of_id IS NULL THEN 'OUTBOUND' ELSE 'INBOUND' END AS direction,
	CASE WHEN od.reversal_of_id IS NULL THEN 'SALES_SHIP' ELSE 'SALES_RETURN' END AS source_type,
	od.occurred_at AS occurred_at,
	od.product_id AS product_id,
	p.seller_sku AS seller_sku,
	p.title AS product_title,
	od.warehouse_id AS warehouse_id,
	w.code AS warehouse_code,
	w.name AS warehouse_name,
	od.marketplace AS marketplace,
	od.qty_out AS quantity,
	od.unit_cost_original AS unit_cost_original,
	od.original_currency AS original_currency,
	od.original_amount AS original_amount,
	od.base_currency AS base_currency,
	od.base_amount AS base_amount,
	'SALES_ORDER' AS reference_type,
	od.sales_order_id AS reference_id,
	so.order_no AS reference_number
FROM finance_order_cost_detail od
LEFT JOIN product p ON p.id = od.product_id
LEFT JOIN warehouse w ON w.id = od.warehouse_id
LEFT JOIN sales_order so ON so.id = od.sales_order_id
WHERE od.status = 'NORMAL'`
}

func buildProductCostWhere(params *domain.ProductCostLedgerListParams) (string, []interface{}) {
	parts := []string{"1=1"}
	args := make([]interface{}, 0)
	if params == nil {
		return strings.Join(parts, " AND "), args
	}
	if params.ProductID != nil && *params.ProductID > 0 {
		parts = append(parts, "product_id = ?")
		args = append(args, *params.ProductID)
	}
	if params.WarehouseID != nil && *params.WarehouseID > 0 {
		parts = append(parts, "warehouse_id = ?")
		args = append(args, *params.WarehouseID)
	}
	if strings.TrimSpace(params.Marketplace) != "" {
		parts = append(parts, "marketplace = ?")
		args = append(args, strings.TrimSpace(params.Marketplace))
	}
	if params.DateFrom != nil {
		parts = append(parts, "occurred_at >= ?")
		args = append(args, *params.DateFrom)
	}
	if params.DateTo != nil {
		parts = append(parts, "occurred_at <= ?")
		args = append(args, *params.DateTo)
	}
	return strings.Join(parts, " AND "), args
}
