package domain

import "time"

type ReplenishmentTriggerType string

const (
	ReplenishmentTriggerScheduled ReplenishmentTriggerType = "SCHEDULED"
	ReplenishmentTriggerManual    ReplenishmentTriggerType = "MANUAL"
)

type ReplenishmentRunStatus string

const (
	ReplenishmentRunRunning ReplenishmentRunStatus = "RUNNING"
	ReplenishmentRunSuccess ReplenishmentRunStatus = "SUCCESS"
	ReplenishmentRunFailed  ReplenishmentRunStatus = "FAILED"
)

type ReplenishmentItemStatus string

const (
	ReplenishmentItemPending   ReplenishmentItemStatus = "PENDING"
	ReplenishmentItemConverted ReplenishmentItemStatus = "CONVERTED"
	ReplenishmentItemSkipped   ReplenishmentItemStatus = "SKIPPED"
)

type ReplenishmentConfig struct {
	ID                   uint64     `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	IsEnabled            uint8      `json:"is_enabled" gorm:"column:is_enabled;not null;default:1"`
	IntervalMinutes      uint32     `json:"interval_minutes" gorm:"column:interval_minutes;not null;default:1440"`
	DemandWindowDays     uint32     `json:"demand_window_days" gorm:"column:demand_window_days;not null;default:30"`
	DefaultLeadTimeDays  uint32     `json:"default_lead_time_days" gorm:"column:default_lead_time_days;not null;default:15"`
	DefaultSafetyDays    uint32     `json:"default_safety_days" gorm:"column:default_safety_days;not null;default:7"`
	DefaultMOQ           uint32     `json:"default_moq" gorm:"column:default_moq;not null;default:1"`
	DefaultOrderMultiple uint32     `json:"default_order_multiple" gorm:"column:default_order_multiple;not null;default:1"`
	LastGeneratedDate    *time.Time `json:"last_generated_date" gorm:"column:last_generated_date;type:date"`
	LastCleanupDate      *time.Time `json:"last_cleanup_date" gorm:"column:last_cleanup_date;type:date"`
	GmtCreate            time.Time  `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified          time.Time  `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ReplenishmentConfig) TableName() string {
	return "procurement_replenishment_config"
}

type ReplenishmentPolicy struct {
	ID               uint64    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	ProductID        uint64    `json:"product_id" gorm:"column:product_id;not null"`
	IsEnabled        uint8     `json:"is_enabled" gorm:"column:is_enabled;not null;default:1"`
	DemandWindowDays *uint32   `json:"demand_window_days" gorm:"column:demand_window_days"`
	LeadTimeDays     *uint32   `json:"lead_time_days" gorm:"column:lead_time_days"`
	SafetyDays       *uint32   `json:"safety_days" gorm:"column:safety_days"`
	MOQ              *uint32   `json:"moq" gorm:"column:moq"`
	OrderMultiple    *uint32   `json:"order_multiple" gorm:"column:order_multiple"`
	Remark           *string   `json:"remark" gorm:"column:remark;size:255"`
	CreatedBy        *uint64   `json:"created_by" gorm:"column:created_by"`
	UpdatedBy        *uint64   `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate        time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ReplenishmentPolicy) TableName() string {
	return "procurement_replenishment_policy"
}

type ReplenishmentRun struct {
	ID            uint64                   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	RunNo         string                   `json:"run_no" gorm:"column:run_no;size:64;not null"`
	TriggerType   ReplenishmentTriggerType `json:"trigger_type" gorm:"column:trigger_type;type:enum('SCHEDULED','MANUAL');not null"`
	Status        ReplenishmentRunStatus   `json:"status" gorm:"column:status;type:enum('RUNNING','SUCCESS','FAILED');not null"`
	WindowDays    uint32                   `json:"window_days" gorm:"column:window_days;not null"`
	StartedAt     *time.Time               `json:"started_at" gorm:"column:started_at"`
	FinishedAt    *time.Time               `json:"finished_at" gorm:"column:finished_at"`
	InputSummary  *string                  `json:"input_summary" gorm:"column:input_summary;type:json"`
	OutputSummary *string                  `json:"output_summary" gorm:"column:output_summary;type:json"`
	ErrorMessage  *string                  `json:"error_message" gorm:"column:error_message;type:text"`
	CreatedBy     *uint64                  `json:"created_by" gorm:"column:created_by"`
	GmtCreate     time.Time                `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified   time.Time                `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
	Items         []ReplenishmentItem      `json:"items,omitempty" gorm:"-"`
}

func (ReplenishmentRun) TableName() string {
	return "procurement_replenishment_run"
}

type ReplenishmentItem struct {
	ID              uint64                  `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	RunID           uint64                  `json:"run_id" gorm:"column:run_id;not null"`
	ProductID       uint64                  `json:"product_id" gorm:"column:product_id;not null"`
	WarehouseID     uint64                  `json:"warehouse_id" gorm:"column:warehouse_id;not null"`
	SupplierID      *uint64                 `json:"supplier_id" gorm:"column:supplier_id"`
	DailyDemand     float64                 `json:"daily_demand" gorm:"column:daily_demand;type:decimal(18,4);not null"`
	CoverageDays    uint32                  `json:"coverage_days" gorm:"column:coverage_days;not null"`
	NetSupply       int64                   `json:"net_supply" gorm:"column:net_supply;not null"`
	TargetStock     uint64                  `json:"target_stock" gorm:"column:target_stock;not null"`
	ShortageQty     uint64                  `json:"shortage_qty" gorm:"column:shortage_qty;not null"`
	SuggestedQty    uint64                  `json:"suggested_qty" gorm:"column:suggested_qty;not null"`
	MOQ             uint32                  `json:"moq" gorm:"column:moq;not null"`
	OrderMultiple   uint32                  `json:"order_multiple" gorm:"column:order_multiple;not null"`
	UnitCost        *float64                `json:"unit_cost" gorm:"column:unit_cost;type:decimal(12,4)"`
	RuleSource      string                  `json:"rule_source" gorm:"column:rule_source;size:50;not null"`
	Status          ReplenishmentItemStatus `json:"status" gorm:"column:status;type:enum('PENDING','CONVERTED','SKIPPED');not null"`
	PurchaseOrderID *uint64                 `json:"purchase_order_id" gorm:"column:purchase_order_id"`
	ConvertedAt     *time.Time              `json:"converted_at" gorm:"column:converted_at"`
	Remark          *string                 `json:"remark" gorm:"column:remark;size:255"`
	GmtCreate       time.Time               `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified     time.Time               `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ReplenishmentItem) TableName() string {
	return "procurement_replenishment_item"
}

type ReplenishmentRunListParams struct {
	Page      int
	PageSize  int
	Status    ReplenishmentRunStatus
	Triggered ReplenishmentTriggerType
}

type ReplenishmentPolicyListParams struct {
	Page      int
	PageSize  int
	Keyword   string
	ProductID *uint64
}

type ReplenishmentGenerateParams struct {
	TriggerType ReplenishmentTriggerType
	OperatorID  *uint64
}

type ReplenishmentConvertParams struct {
	RunID      uint64
	ItemIDs    []uint64
	OperatorID *uint64
}
