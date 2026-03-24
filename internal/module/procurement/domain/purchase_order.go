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
	ID                        uint64              `json:"id" gorm:"primaryKey;autoIncrement"`
	PoNumber                  string              `json:"po_number" gorm:"column:po_number;size:50;not null"`
	BatchNo                   string              `json:"batch_no" gorm:"column:batch_no;size:50"`
	SupplierID                *uint64             `json:"supplier_id" gorm:"column:supplier_id"`
	WarehouseID               *uint64             `json:"warehouse_id,omitempty" gorm:"column:warehouse_id"`
	Marketplace               string              `json:"marketplace" gorm:"size:10"`
	Status                    PurchaseOrderStatus `json:"status" gorm:"type:enum('DRAFT','ORDERED','SHIPPED','RECEIVED','CLOSED');default:'DRAFT'"`
	Currency                  string              `json:"currency" gorm:"size:10"`
	TotalAmount               float64             `json:"total_amount" gorm:"type:decimal(12,4);default:0"`
	OrderedAt                 *time.Time          `json:"ordered_at" gorm:"column:ordered_at"`
	OrderedBy                 *uint64             `json:"ordered_by" gorm:"column:ordered_by"`
	ShippedAt                 *time.Time          `json:"shipped_at" gorm:"column:shipped_at"`
	ShippedBy                 *uint64             `json:"shipped_by" gorm:"column:shipped_by"`
	ReceivedAt                *time.Time          `json:"received_at" gorm:"column:received_at"`
	ReceivedBy                *uint64             `json:"received_by" gorm:"column:received_by"`
	InspectedAt               *time.Time          `json:"inspected_at" gorm:"column:inspected_at"`
	InspectedBy               *uint64             `json:"inspected_by" gorm:"column:inspected_by"`
	ClosedAt                  *time.Time          `json:"closed_at" gorm:"column:closed_at"`
	CompletedBy               *uint64             `json:"completed_by" gorm:"column:completed_by"`
	IsForceCompleted          uint8               `json:"is_force_completed" gorm:"column:is_force_completed;not null;default:0"`
	ForceCompletedAt          *time.Time          `json:"force_completed_at" gorm:"column:force_completed_at"`
	ForceCompletedBy          *uint64             `json:"force_completed_by" gorm:"column:force_completed_by"`
	ForceCompleteReason       string              `json:"force_complete_reason" gorm:"column:force_complete_reason;type:text"`
	Remark                    string              `json:"remark" gorm:"type:text"`
	CreatedBy                 *uint64             `json:"created_by" gorm:"column:created_by"`
	UpdatedBy                 *uint64             `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate                 time.Time           `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified               time.Time           `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
	Supplier                  *SupplierSnapshot   `json:"supplier,omitempty" gorm:"-"`
	Items                     []PurchaseOrderItem `json:"items,omitempty" gorm:"-"`
	QtyPendingInspectionTotal uint64              `json:"qty_pending_inspection_total" gorm:"-"`
	SourcePlanIDs             []uint64            `json:"source_plan_ids,omitempty" gorm:"-"`
	OrderedByName             string              `json:"ordered_by_name,omitempty" gorm:"-"`
	ShippedByName             string              `json:"shipped_by_name,omitempty" gorm:"-"`
	ReceivedByName            string              `json:"received_by_name,omitempty" gorm:"-"`
	InspectedByName           string              `json:"inspected_by_name,omitempty" gorm:"-"`
	CompletedByName           string              `json:"completed_by_name,omitempty" gorm:"-"`
	ForceCompletedByName      string              `json:"force_completed_by_name,omitempty" gorm:"-"`
}

func (PurchaseOrder) TableName() string {
	return "purchase_order"
}

type PurchaseOrderItem struct {
	ID                   uint64           `json:"id" gorm:"primaryKey;autoIncrement"`
	PurchaseOrderID      uint64           `json:"purchase_order_id" gorm:"column:purchase_order_id;not null"`
	ProductID            uint64           `json:"product_id" gorm:"column:product_id;not null"`
	SupplierID           *uint64          `json:"supplier_id,omitempty" gorm:"-"`
	SourcePlanID         *uint64          `json:"source_plan_id,omitempty" gorm:"-"`
	QtyOrdered           uint64           `json:"qty_ordered" gorm:"column:qty_ordered;not null"`
	QtyReceived          uint64           `json:"qty_received" gorm:"column:qty_received;not null"`
	QtyInspectionPass    uint64           `json:"qty_inspection_pass" gorm:"column:qty_inspection_pass;not null;default:0"`
	QtyInspectionFail    uint64           `json:"qty_inspection_fail" gorm:"column:qty_inspection_fail;not null;default:0"`
	QtyPendingInspection uint64           `json:"qty_pending_inspection" gorm:"-"`
	UnitCost             float64          `json:"unit_cost" gorm:"column:unit_cost;type:decimal(12,4);default:0"`
	Currency             string           `json:"currency" gorm:"column:currency;size:10"`
	Subtotal             float64          `json:"subtotal" gorm:"column:subtotal;type:decimal(12,4);default:0"`
	GmtCreate            time.Time        `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified          time.Time        `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
	Product              *ProductSnapshot `json:"product,omitempty" gorm:"-"`
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

type PurchaseOrderInspectParams struct {
	PassQties  map[uint64]uint64
	FailQties  map[uint64]uint64
	OperatorID *uint64
}

type PurchaseOrderForceCompleteParams struct {
	Reason     string
	OperatorID *uint64
}

type PurchaseOrderShipParams struct {
	WarehouseID uint64
	OperatorID  *uint64
}

type SupplierSnapshot struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type ProductSnapshot struct {
	ID                   uint64 `json:"id"`
	SellerSku            string `json:"seller_sku"`
	Title                string `json:"title"`
	ImageURL             string `json:"image_url,omitempty"`
	IsInspectionRequired uint8  `json:"is_inspection_required"`
	IsPackingRequired    uint8  `json:"is_packing_required"`
}

type InventoryMovement struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID         string    `json:"trace_id" gorm:"column:trace_id;size:50"`
	ProductID       uint64    `json:"product_id" gorm:"column:product_id;not null"`
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
