package domain

import "time"

type InventoryLotStatus string

const (
	InventoryLotStatusOpen   InventoryLotStatus = "OPEN"
	InventoryLotStatusClosed InventoryLotStatus = "CLOSED"
)

type InventoryLot struct {
	ID                         uint64             `json:"id" gorm:"primaryKey;column:id"`
	ProductID                  uint64             `json:"product_id" gorm:"column:product_id;not null;index"`
	WarehouseID                uint64             `json:"warehouse_id" gorm:"column:warehouse_id;not null;index"`
	LotNo                      string             `json:"lot_no" gorm:"column:lot_no;size:64;not null"`
	SourceType                 *string            `json:"source_type" gorm:"column:source_type;size:50"`
	SourceID                   *uint64            `json:"source_id" gorm:"column:source_id"`
	SourceNumber               *string            `json:"source_number" gorm:"column:source_number;size:100"`
	ReceivedAt                 time.Time          `json:"received_at" gorm:"column:received_at;not null;index"`
	UnitCost                   *float64           `json:"unit_cost" gorm:"column:unit_cost;type:decimal(12,4)"`
	QtyIn                      uint               `json:"qty_in" gorm:"column:qty_in;not null;default:0"`
	QtyPurchasingInTransit     uint               `json:"qty_purchasing_in_transit" gorm:"column:qty_purchasing_in_transit;not null;default:0"`
	QtyPendingInspection       uint               `json:"qty_pending_inspection" gorm:"column:qty_pending_inspection;not null;default:0"`
	QtyRawMaterial             uint               `json:"qty_raw_material" gorm:"column:qty_raw_material;not null;default:0"`
	QtyAvailable               uint               `json:"qty_available" gorm:"column:qty_available;not null;default:0"`
	QtyReserved                uint               `json:"qty_reserved" gorm:"column:qty_reserved;not null;default:0"`
	QtyPendingShipment         uint               `json:"qty_pending_shipment" gorm:"column:qty_pending_shipment;not null;default:0"`
	QtyPendingShipmentReserved uint               `json:"qty_pending_shipment_reserved" gorm:"column:qty_pending_shipment_reserved;not null;default:0"`
	QtySellable                uint               `json:"qty_sellable" gorm:"column:qty_sellable;not null;default:0"`
	QtySellableReserved        uint               `json:"qty_sellable_reserved" gorm:"column:qty_sellable_reserved;not null;default:0"`
	QtyConsumed                uint               `json:"qty_consumed" gorm:"column:qty_consumed;not null;default:0"`
	Status                     InventoryLotStatus `json:"status" gorm:"column:status;type:enum('OPEN','CLOSED');not null;default:'OPEN'"`
	Remark                     *string            `json:"remark" gorm:"column:remark;size:500"`
	GmtCreate                  time.Time          `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified                time.Time          `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`

	Product   *ProductSnapshot   `json:"product,omitempty" gorm:"-"`
	Warehouse *WarehouseSnapshot `json:"warehouse,omitempty" gorm:"-"`
}

func (InventoryLot) TableName() string {
	return "inventory_lot"
}

type InventoryLotListParams struct {
	Page        int
	PageSize    int
	ProductID   *uint64
	WarehouseID *uint64
	Status      *InventoryLotStatus
	Keyword     *string
}
