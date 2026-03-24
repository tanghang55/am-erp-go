package domain

import "time"

type ProductCostDirection string

const (
	ProductCostDirectionInbound  ProductCostDirection = "INBOUND"
	ProductCostDirectionOutbound ProductCostDirection = "OUTBOUND"
	ProductCostDirectionNeutral  ProductCostDirection = "NEUTRAL"
)

type ProductCostLedgerListParams struct {
	Page        int
	PageSize    int
	ProductID   *uint64
	WarehouseID *uint64
	Marketplace string
	DateFrom    *time.Time
	DateTo      *time.Time
}

type ProductCostLedgerItem struct {
	Direction        ProductCostDirection `json:"direction"`
	SourceType       string               `json:"source_type"`
	OccurredAt       time.Time            `json:"occurred_at"`
	ProductID        uint64               `json:"product_id"`
	SellerSKU        *string              `json:"seller_sku,omitempty"`
	ProductTitle     *string              `json:"product_title,omitempty"`
	WarehouseID      *uint64              `json:"warehouse_id,omitempty"`
	WarehouseCode    *string              `json:"warehouse_code,omitempty"`
	WarehouseName    *string              `json:"warehouse_name,omitempty"`
	Marketplace      *string              `json:"marketplace,omitempty"`
	Quantity         uint64               `json:"quantity"`
	UnitCostOriginal float64              `json:"unit_cost_original"`
	OriginalCurrency string               `json:"original_currency"`
	OriginalAmount   float64              `json:"original_amount"`
	BaseCurrency     string               `json:"base_currency"`
	BaseAmount       float64              `json:"base_amount"`
	ReferenceType    *string              `json:"reference_type,omitempty"`
	ReferenceID      *uint64              `json:"reference_id,omitempty"`
	ReferenceNumber  *string              `json:"reference_number,omitempty"`
}

type ProductCostSummary struct {
	BaseCurrency        string  `json:"base_currency"`
	InboundQty          uint64  `json:"inbound_qty"`
	InboundAmount       float64 `json:"inbound_amount"`
	OutboundQty         uint64  `json:"outbound_qty"`
	OutboundAmount      float64 `json:"outbound_amount"`
	NeutralAmount       float64 `json:"neutral_amount"`
	NetQty              int64   `json:"net_qty"`
	NetAmount           float64 `json:"net_amount"`
	AvgInboundUnitCost  float64 `json:"avg_inbound_unit_cost"`
	AvgOutboundUnitCost float64 `json:"avg_outbound_unit_cost"`
}
