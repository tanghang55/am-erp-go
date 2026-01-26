package domain

import "time"

type InventoryBalance struct {
	ID                uint64     `json:"id" gorm:"primaryKey;column:id"`
	SkuID             uint64     `json:"sku_id" gorm:"column:sku_id;not null;index"`
	WarehouseID       uint64     `json:"warehouse_id" gorm:"column:warehouse_id;not null;index"`
	AvailableQuantity uint       `json:"available_quantity" gorm:"column:available_quantity;not null;index"`
	ReservedQuantity  uint       `json:"reserved_quantity" gorm:"column:reserved_quantity;not null"`
	DamagedQuantity   uint       `json:"damaged_quantity" gorm:"column:damaged_quantity;not null"`

	// 库存状态字段
	PurchasingInTransit  uint `json:"purchasing_in_transit" gorm:"column:purchasing_in_transit;not null;default:0"`  // 采购在途库存
	PendingInspection    uint `json:"pending_inspection" gorm:"column:pending_inspection;not null;default:0"`        // 待检库存
	RawMaterial          uint `json:"raw_material" gorm:"column:raw_material;not null;default:0"`                    // 原料库存(上架库存)
	PendingShipment      uint `json:"pending_shipment" gorm:"column:pending_shipment;not null;default:0"`            // 待出库存
	LogisticsInTransit   uint `json:"logistics_in_transit" gorm:"column:logistics_in_transit;not null;default:0"`    // 物流在途库存
	Sellable             uint `json:"sellable" gorm:"column:sellable;not null;default:0"`                            // 可售库存
	Returned             uint `json:"returned" gorm:"column:returned;not null;default:0"`                            // 退货库存

	TotalQuantity  uint       `json:"total_quantity" gorm:"column:total_quantity;not null;index"`
	LastMovementAt *time.Time `json:"last_movement_at" gorm:"column:last_movement_at"`
	GmtCreate      time.Time  `json:"created_at" gorm:"column:gmt_create"`
	GmtModified    time.Time  `json:"updated_at" gorm:"column:gmt_modified;index"`

	// Associations (not stored in DB)
	Sku       *SkuSnapshot       `json:"sku,omitempty" gorm:"-"`
	Warehouse *WarehouseSnapshot `json:"warehouse,omitempty" gorm:"-"`
}

func (InventoryBalance) TableName() string {
	return "inventory_balance"
}

type BalanceListParams struct {
	Page                int
	PageSize            int
	WarehouseID         *uint64
	SkuID               *uint64
	LowStock            *bool
	LowStockThreshold   *uint
	ZeroStock           *bool
	Keyword             *string
}
