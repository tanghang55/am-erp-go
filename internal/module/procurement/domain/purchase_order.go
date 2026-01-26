package domain

import "time"

type PurchaseOrderStatus string

const (
	PurchaseOrderStatusDraft    PurchaseOrderStatus = "DRAFT"
	PurchaseOrderStatusOrdered  PurchaseOrderStatus = "ORDERED"
	PurchaseOrderStatusShipped  PurchaseOrderStatus = "SHIPPED"
	PurchaseOrderStatusReceived PurchaseOrderStatus = "RECEIVED"
	PurchaseOrderStatusClosed   PurchaseOrderStatus = "CLOSED"
)

type PurchaseOrder struct {
	ID          uint64              `json:"id" gorm:"primaryKey;autoIncrement"`
	PoNumber    string              `json:"po_number" gorm:"column:po_number;size:50;not null"`
	SupplierID  *uint64             `json:"supplier_id" gorm:"column:supplier_id"`
	Marketplace string              `json:"marketplace" gorm:"size:10"`
	Status      PurchaseOrderStatus `json:"status" gorm:"type:enum('DRAFT','ORDERED','SHIPPED','RECEIVED','CLOSED');default:'DRAFT'"`
	Currency    string              `json:"currency" gorm:"size:10;default:'USD'"`
	TotalAmount float64             `json:"total_amount" gorm:"type:decimal(12,4);default:0"`
	OrderedAt   *time.Time          `json:"ordered_at" gorm:"column:ordered_at"`
	ShippedAt   *time.Time          `json:"shipped_at" gorm:"column:shipped_at"`
	ReceivedAt  *time.Time          `json:"received_at" gorm:"column:received_at"`
	Remark      string              `json:"remark" gorm:"type:text"`
	CreatedBy   *uint64             `json:"created_by" gorm:"column:created_by"`
	UpdatedBy   *uint64             `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate   time.Time           `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time           `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
	Supplier    *SupplierSnapshot   `json:"supplier,omitempty" gorm:"-"`
	Items       []PurchaseOrderItem `json:"items,omitempty" gorm:"-"`
}

func (PurchaseOrder) TableName() string {
	return "purchase_order"
}

type PurchaseOrderItem struct {
	ID              uint64       `json:"id" gorm:"primaryKey;autoIncrement"`
	PurchaseOrderID uint64       `json:"purchase_order_id" gorm:"column:purchase_order_id;not null"`
	SkuID           uint64       `json:"sku_id" gorm:"column:sku_id;not null"`
	QtyOrdered      uint64       `json:"qty_ordered" gorm:"column:qty_ordered;not null"`
	QtyReceived     uint64       `json:"qty_received" gorm:"column:qty_received;not null"`
	UnitCost        float64      `json:"unit_cost" gorm:"column:unit_cost;type:decimal(12,4);default:0"`
	Currency        string       `json:"currency" gorm:"column:currency;size:10;default:'USD'"`
	Subtotal        float64      `json:"subtotal" gorm:"column:subtotal;type:decimal(12,4);default:0"`
	GmtCreate       time.Time    `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified     time.Time    `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
	Sku             *SkuSnapshot `json:"sku,omitempty" gorm:"-"`
}

func (PurchaseOrderItem) TableName() string {
	return "purchase_order_item"
}

type PurchaseOrderListParams struct {
	Page        int
	PageSize    int
	Status      PurchaseOrderStatus
	SupplierID  *uint64
	Marketplace string
	Keyword     string
}

type PurchaseOrderReceiveParams struct {
	WarehouseID   uint64
	ReceivedQties map[uint64]uint64
	OperatorID    *uint64
}

type PurchaseOrderShipParams struct {
	WarehouseID uint64
	OperatorID  *uint64
}

type SupplierSnapshot struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type SkuSnapshot struct {
	ID        uint64 `json:"id"`
	SellerSku string `json:"seller_sku"`
	Title     string `json:"title"`
}

type InventoryMovement struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID         string    `json:"trace_id" gorm:"column:trace_id;size:50"`
	SkuID           uint64    `json:"sku_id" gorm:"column:sku_id;not null"`
	WarehouseID     uint64    `json:"warehouse_id" gorm:"column:warehouse_id;not null"`
	MovementType    string    `json:"movement_type" gorm:"column:movement_type;size:50"`
	ReferenceType   string    `json:"reference_type" gorm:"column:reference_type;size:50"`
	ReferenceID     *uint64   `json:"reference_id" gorm:"column:reference_id"`
	ReferenceNumber string    `json:"reference_number" gorm:"column:reference_number;size:100"`
	Quantity        int64     `json:"quantity" gorm:"column:quantity;not null"`
	BeforeAvailable uint64    `json:"before_available" gorm:"column:before_available;not null"`
	AfterAvailable  uint64    `json:"after_available" gorm:"column:after_available;not null"`
	BeforeReserved  uint64    `json:"before_reserved" gorm:"column:before_reserved;not null"`
	AfterReserved   uint64    `json:"after_reserved" gorm:"column:after_reserved;not null"`
	BeforeDamaged   uint64    `json:"before_damaged" gorm:"column:before_damaged;not null"`
	AfterDamaged    uint64    `json:"after_damaged" gorm:"column:after_damaged;not null"`
	UnitCost        float64   `json:"unit_cost" gorm:"column:unit_cost;type:decimal(12,4);default:0"`
	TotalCost       float64   `json:"total_cost" gorm:"column:total_cost;type:decimal(12,4);default:0"`
	Remark          string    `json:"remark" gorm:"column:remark;type:text"`
	OperatorID      *uint64   `json:"operator_id" gorm:"column:operator_id"`
	OperatedAt      time.Time `json:"operated_at" gorm:"column:operated_at"`
	GmtCreate       time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified     time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (InventoryMovement) TableName() string {
	return "inventory_movement"
}
