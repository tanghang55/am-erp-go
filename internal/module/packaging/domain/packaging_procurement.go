package domain

import "time"

type PackagingProcurementPlanStatus string

const (
	PackagingProcurementPlanPending   PackagingProcurementPlanStatus = "PENDING"
	PackagingProcurementPlanConverted PackagingProcurementPlanStatus = "CONVERTED"
	PackagingProcurementPlanCancelled PackagingProcurementPlanStatus = "CANCELLED"
)

type PackagingProcurementTriggerType string

const (
	PackagingProcurementTriggerManual    PackagingProcurementTriggerType = "MANUAL"
	PackagingProcurementTriggerScheduled PackagingProcurementTriggerType = "SCHEDULED"
)

type PackagingProcurementRunStatus string

const (
	PackagingProcurementRunRunning PackagingProcurementRunStatus = "RUNNING"
	PackagingProcurementRunSuccess PackagingProcurementRunStatus = "SUCCESS"
	PackagingProcurementRunFailed  PackagingProcurementRunStatus = "FAILED"
)

type PackagingPurchaseOrderStatus string

const (
	PackagingPurchaseOrderDraft    PackagingPurchaseOrderStatus = "DRAFT"
	PackagingPurchaseOrderOrdered  PackagingPurchaseOrderStatus = "ORDERED"
	PackagingPurchaseOrderReceived PackagingPurchaseOrderStatus = "RECEIVED"
	PackagingPurchaseOrderClosed   PackagingPurchaseOrderStatus = "CLOSED"
)

type PackagingProductDemand struct {
	ProductID uint64
	Qty       uint64
}

type ProductPackagingMapping struct {
	ProductID       uint64
	PackagingItemID uint64
	QuantityPerUnit float64
}

type PackagingItemSnapshot struct {
	ID              uint64
	ItemCode        string
	ItemName        string
	Unit            string
	QuantityOnHand  uint64
	ReorderQuantity *uint64
	UnitCost        float64
	Currency        string
	Status          string
}

type PackagingPlanInput struct {
	PackagingItemID uint64
	RequiredQty     uint64
	OnHandQty       uint64
	ShortageQty     uint64
	SuggestedQty    uint64
	SourceJSON      *string
	Remark          *string
}

type PackagingProcurementPlan struct {
	ID                       uint64                         `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	PlanDate                 time.Time                      `json:"plan_date" gorm:"column:plan_date;type:date;not null"`
	PackagingItemID          uint64                         `json:"packaging_item_id" gorm:"column:packaging_item_id;not null"`
	RequiredQty              uint64                         `json:"required_qty" gorm:"column:required_qty;not null"`
	OnHandQty                uint64                         `json:"on_hand_qty" gorm:"column:on_hand_qty;not null"`
	ShortageQty              uint64                         `json:"shortage_qty" gorm:"column:shortage_qty;not null"`
	SuggestedQty             uint64                         `json:"suggested_qty" gorm:"column:suggested_qty;not null"`
	Status                   PackagingProcurementPlanStatus `json:"status" gorm:"column:status;type:enum('PENDING','CONVERTED','CANCELLED');not null"`
	PackagingPurchaseOrderID *uint64                        `json:"packaging_purchase_order_id" gorm:"column:packaging_purchase_order_id"`
	ConvertedAt              *time.Time                     `json:"converted_at" gorm:"column:converted_at"`
	SourceJSON               *string                        `json:"source_json" gorm:"column:source_json;type:json"`
	Remark                   *string                        `json:"remark" gorm:"column:remark;size:255"`
	GmtCreate                time.Time                      `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified              time.Time                      `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`

	PackagingItemCode            string `json:"packaging_item_code,omitempty" gorm:"-"`
	PackagingItemName            string `json:"packaging_item_name,omitempty" gorm:"-"`
	PackagingItemUnit            string `json:"packaging_item_unit,omitempty" gorm:"-"`
	PackagingPurchaseOrderNumber string `json:"packaging_purchase_order_number,omitempty" gorm:"column:packaging_purchase_order_number"`
}

func (PackagingProcurementPlan) TableName() string {
	return "packaging_procurement_plan"
}

type PackagingProcurementRun struct {
	ID            uint64                          `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	RunNo         string                          `json:"run_no" gorm:"column:run_no;size:64;not null"`
	TriggerType   PackagingProcurementTriggerType `json:"trigger_type" gorm:"column:trigger_type;type:enum('MANUAL','SCHEDULED');not null"`
	Status        PackagingProcurementRunStatus   `json:"status" gorm:"column:status;type:enum('RUNNING','SUCCESS','FAILED');not null"`
	StartedAt     *time.Time                      `json:"started_at" gorm:"column:started_at"`
	FinishedAt    *time.Time                      `json:"finished_at" gorm:"column:finished_at"`
	InputSummary  *string                         `json:"input_summary" gorm:"column:input_summary;type:json"`
	OutputSummary *string                         `json:"output_summary" gorm:"column:output_summary;type:json"`
	ErrorMessage  *string                         `json:"error_message" gorm:"column:error_message;type:text"`
	CreatedBy     *uint64                         `json:"created_by" gorm:"column:created_by"`
	GmtCreate     time.Time                       `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified   time.Time                       `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (PackagingProcurementRun) TableName() string {
	return "packaging_procurement_run"
}

type PackagingPurchaseOrder struct {
	ID          uint64                       `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	PoNumber    string                       `json:"po_number" gorm:"column:po_number;size:50;not null"`
	Status      PackagingPurchaseOrderStatus `json:"status" gorm:"column:status;type:enum('DRAFT','ORDERED','RECEIVED','CLOSED');not null"`
	Currency    string                       `json:"currency" gorm:"column:currency;size:10;not null"`
	TotalAmount float64                      `json:"total_amount" gorm:"column:total_amount;type:decimal(12,4);not null"`
	OrderedAt   *time.Time                   `json:"ordered_at" gorm:"column:ordered_at"`
	ReceivedAt  *time.Time                   `json:"received_at" gorm:"column:received_at"`
	Remark      string                       `json:"remark" gorm:"column:remark;type:text"`
	CreatedBy   *uint64                      `json:"created_by" gorm:"column:created_by"`
	UpdatedBy   *uint64                      `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate   time.Time                    `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time                    `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
	Items       []PackagingPurchaseOrderItem `json:"items,omitempty" gorm:"-"`
}

func (PackagingPurchaseOrder) TableName() string {
	return "packaging_purchase_order"
}

type PackagingPurchaseOrderItem struct {
	ID                       uint64    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	PackagingPurchaseOrderID uint64    `json:"packaging_purchase_order_id" gorm:"column:packaging_purchase_order_id;not null"`
	PackagingItemID          uint64    `json:"packaging_item_id" gorm:"column:packaging_item_id;not null"`
	QtyOrdered               uint64    `json:"qty_ordered" gorm:"column:qty_ordered;not null"`
	QtyReceived              uint64    `json:"qty_received" gorm:"column:qty_received;not null"`
	UnitCost                 float64   `json:"unit_cost" gorm:"column:unit_cost;type:decimal(12,4);not null"`
	Currency                 string    `json:"currency" gorm:"column:currency;size:10;not null"`
	Subtotal                 float64   `json:"subtotal" gorm:"column:subtotal;type:decimal(12,4);not null"`
	GmtCreate                time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified              time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`

	PackagingItemCode string `json:"packaging_item_code,omitempty" gorm:"-"`
	PackagingItemName string `json:"packaging_item_name,omitempty" gorm:"-"`
	PackagingItemUnit string `json:"packaging_item_unit,omitempty" gorm:"-"`
}

func (PackagingPurchaseOrderItem) TableName() string {
	return "packaging_purchase_order_item"
}

type PackagingProcurementPlanListParams struct {
	Page     int
	PageSize int
	Date     *time.Time
	Status   PackagingProcurementPlanStatus
}

type PackagingProcurementRunListParams struct {
	Page        int
	PageSize    int
	Status      PackagingProcurementRunStatus
	TriggerType PackagingProcurementTriggerType
}

type PackagingPlanConvertParams struct {
	PlanIDs    []uint64
	Date       *time.Time
	OperatorID *uint64
}

type PackagingPurchaseOrderListParams struct {
	Page     int
	PageSize int
	Status   PackagingPurchaseOrderStatus
}

type PackagingPurchaseOrderReceiveParams struct {
	ReceivedQties map[uint64]uint64
	OperatorID    *uint64
}

type PackagingProcurementRepository interface {
	CleanupPlansBefore(date time.Time) error
	LoadOrderedProductDemands(planDate time.Time) ([]PackagingProductDemand, error)
	LoadProductPackagingMappings(productIDs []uint64) ([]ProductPackagingMapping, error)
	LoadPackagingItemSnapshots(itemIDs []uint64) (map[uint64]PackagingItemSnapshot, error)
	SyncDailyPlans(planDate time.Time, inputs []PackagingPlanInput) ([]PackagingProcurementPlan, int, error)
	ListPlans(params *PackagingProcurementPlanListParams) ([]PackagingProcurementPlan, int64, error)
	ListRuns(params *PackagingProcurementRunListParams) ([]PackagingProcurementRun, int64, error)
	ListConvertiblePlans(params *PackagingPlanConvertParams) ([]PackagingProcurementPlan, error)
	MarkPlansConverted(planIDs []uint64, purchaseOrderID uint64) error
	CreateRun(run *PackagingProcurementRun) error
	UpdateRun(run *PackagingProcurementRun) error

	CreatePurchaseOrder(order *PackagingPurchaseOrder) error
	UpdatePurchaseOrder(order *PackagingPurchaseOrder) error
	ListPurchaseOrders(params *PackagingPurchaseOrderListParams) ([]PackagingPurchaseOrder, int64, error)
	GetPurchaseOrder(id uint64) (*PackagingPurchaseOrder, error)
}
