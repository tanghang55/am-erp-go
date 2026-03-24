package domain

import "time"

type OrderCostDetailStatus string

const (
	OrderCostDetailStatusNormal   OrderCostDetailStatus = "NORMAL"
	OrderCostDetailStatusReversed OrderCostDetailStatus = "REVERSED"
)

type OrderCostDetail struct {
	ID               uint64                `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID          string                `json:"trace_id" gorm:"column:trace_id;size:64;not null"`
	Status           OrderCostDetailStatus `json:"status" gorm:"column:status;type:enum('NORMAL','REVERSED');not null;default:'NORMAL'"`
	ReversalOfID     *uint64               `json:"reversal_of_id" gorm:"column:reversal_of_id"`
	SalesOrderID     uint64                `json:"sales_order_id" gorm:"column:sales_order_id;not null"`
	SalesOrderItemID uint64                `json:"sales_order_item_id" gorm:"column:sales_order_item_id;not null"`
	ProductID        uint64                `json:"product_id" gorm:"column:product_id;not null"`
	WarehouseID      uint64                `json:"warehouse_id" gorm:"column:warehouse_id;not null"`
	Marketplace      *string               `json:"marketplace" gorm:"column:marketplace;size:10"`
	InventoryLotID   uint64                `json:"inventory_lot_id" gorm:"column:inventory_lot_id;not null"`
	QtyOut           uint64                `json:"qty_out" gorm:"column:qty_out;not null"`
	UnitCostOriginal float64               `json:"unit_cost_original" gorm:"column:unit_cost_original;type:decimal(18,6);not null"`
	OriginalCurrency string                `json:"original_currency" gorm:"column:original_currency;size:10;not null"`
	OriginalAmount   float64               `json:"original_amount" gorm:"column:original_amount;type:decimal(18,6);not null"`
	BaseCurrency     string                `json:"base_currency" gorm:"column:base_currency;size:10;not null"`
	FxRate           float64               `json:"fx_rate" gorm:"column:fx_rate;type:decimal(18,8);not null"`
	BaseAmount       float64               `json:"base_amount" gorm:"column:base_amount;type:decimal(18,6);not null"`
	FxSource         string                `json:"fx_source" gorm:"column:fx_source;size:50;not null"`
	FxVersion        string                `json:"fx_version" gorm:"column:fx_version;size:32;not null"`
	FxTime           time.Time             `json:"fx_time" gorm:"column:fx_time;not null"`
	OccurredAt       time.Time             `json:"occurred_at" gorm:"column:occurred_at;not null"`
	OperatorID       *uint64               `json:"operator_id" gorm:"column:operator_id"`
	Remark           *string               `json:"remark" gorm:"column:remark;size:500"`
	GmtCreate        time.Time             `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time             `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

type ReturnableOrderCostDetail struct {
	OrderCostDetail
	AvailableQty uint64 `json:"available_qty"`
}

func (OrderCostDetail) TableName() string {
	return "finance_order_cost_detail"
}
