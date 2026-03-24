package domain

import "time"

type CostEventType string

const (
	CostEventTypePOOrdered         CostEventType = "PO_ORDERED"
	CostEventTypePOShipped         CostEventType = "PO_SHIPPED"
	CostEventTypePOReceived        CostEventType = "PO_RECEIVED"
	CostEventTypePOAdjust          CostEventType = "PO_ADJUST"
	CostEventTypePackingMaterial   CostEventType = "PACKING_MATERIAL"
	CostEventTypeShipmentAllocated CostEventType = "SHIPMENT_ALLOCATED"
)

type CostEventStatus string

const (
	CostEventStatusNormal   CostEventStatus = "NORMAL"
	CostEventStatusReversed CostEventStatus = "REVERSED"
)

type CostEvent struct {
	ID                  uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID             string          `json:"trace_id" gorm:"column:trace_id;size:64;not null"`
	EventType           CostEventType   `json:"event_type" gorm:"column:event_type;type:enum('PO_ORDERED','PO_SHIPPED','PO_RECEIVED','PO_ADJUST','PACKING_MATERIAL','SHIPMENT_ALLOCATED');not null"`
	Status              CostEventStatus `json:"status" gorm:"column:status;type:enum('NORMAL','REVERSED');not null;default:'NORMAL'"`
	ReversalOfID        *uint64         `json:"reversal_of_id" gorm:"column:reversal_of_id"`
	PurchaseOrderID     *uint64         `json:"purchase_order_id" gorm:"column:purchase_order_id"`
	PurchaseOrderItemID *uint64         `json:"purchase_order_item_id" gorm:"column:purchase_order_item_id"`
	ShipmentID          *uint64         `json:"shipment_id" gorm:"column:shipment_id"`
	ShipmentItemID      *uint64         `json:"shipment_item_id" gorm:"column:shipment_item_id"`
	InventoryMovementID *uint64         `json:"inventory_movement_id" gorm:"column:inventory_movement_id"`
	ProductID           uint64          `json:"product_id" gorm:"column:product_id;not null"`
	WarehouseID         *uint64         `json:"warehouse_id" gorm:"column:warehouse_id"`
	Marketplace         *string         `json:"marketplace" gorm:"column:marketplace;size:10"`
	QtyEvent            uint64          `json:"qty_event" gorm:"column:qty_event;not null"`
	OriginalCurrency    string          `json:"original_currency" gorm:"column:original_currency;size:10;not null"`
	OriginalAmount      float64         `json:"original_amount" gorm:"column:original_amount;type:decimal(18,6);not null"`
	BaseCurrency        string          `json:"base_currency" gorm:"column:base_currency;size:10;not null"`
	FxRate              float64         `json:"fx_rate" gorm:"column:fx_rate;type:decimal(18,8);not null"`
	BaseAmount          float64         `json:"base_amount" gorm:"column:base_amount;type:decimal(18,6);not null"`
	FxSource            string          `json:"fx_source" gorm:"column:fx_source;size:50;not null"`
	FxVersion           string          `json:"fx_version" gorm:"column:fx_version;size:32;not null"`
	FxTime              time.Time       `json:"fx_time" gorm:"column:fx_time;not null"`
	OccurredAt          time.Time       `json:"occurred_at" gorm:"column:occurred_at;not null"`
	OperatorID          *uint64         `json:"operator_id" gorm:"column:operator_id"`
	Remark              *string         `json:"remark" gorm:"column:remark;size:500"`
	GmtCreate           time.Time       `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified         time.Time       `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (CostEvent) TableName() string {
	return "finance_cost_event"
}
