package domain

import "time"

type ReplenishmentStrategy struct {
	ID                   uint64    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Name                 string    `json:"name" gorm:"column:name;size:100;not null"`
	Priority             uint32    `json:"priority" gorm:"column:priority;not null;default:100"`
	IsEnabled            uint8     `json:"is_enabled" gorm:"column:is_enabled;not null;default:1"`
	ProductID            *uint64   `json:"product_id" gorm:"column:product_id"`
	WarehouseID          *uint64   `json:"warehouse_id" gorm:"column:warehouse_id"`
	SupplierID           *uint64   `json:"supplier_id" gorm:"column:supplier_id"`
	Marketplace          *string   `json:"marketplace" gorm:"column:marketplace;size:20"`
	ConditionJSON        *string   `json:"condition_json" gorm:"column:condition_json;type:json"`
	DemandWindowDays     uint32    `json:"demand_window_days" gorm:"column:demand_window_days;not null;default:30"`
	ProcurementCycleDays uint32    `json:"procurement_cycle_days" gorm:"column:procurement_cycle_days;not null;default:15"`
	PackDays             uint32    `json:"pack_days" gorm:"column:pack_days;not null;default:3"`
	LogisticsDays        uint32    `json:"logistics_days" gorm:"column:logistics_days;not null;default:7"`
	SafetyDays           uint32    `json:"safety_days" gorm:"column:safety_days;not null;default:7"`
	ZeroSalesPurchaseQty uint32    `json:"zero_sales_purchase_qty" gorm:"column:zero_sales_purchase_qty;not null;default:0"`
	MOQ                  uint32    `json:"moq" gorm:"column:moq;not null;default:1"`
	OrderMultiple        uint32    `json:"order_multiple" gorm:"column:order_multiple;not null;default:1"`
	Remark               *string   `json:"remark" gorm:"column:remark;size:255"`
	CreatedBy            *uint64   `json:"created_by" gorm:"column:created_by"`
	UpdatedBy            *uint64   `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate            time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified          time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`

	SellerSKU     string  `json:"seller_sku,omitempty" gorm:"column:seller_sku"`
	ProductTitle  string  `json:"product_title,omitempty" gorm:"column:product_title"`
	WarehouseCode string  `json:"warehouse_code,omitempty" gorm:"column:warehouse_code"`
	WarehouseName string  `json:"warehouse_name,omitempty" gorm:"column:warehouse_name"`
	SupplierCode  *string `json:"supplier_code,omitempty" gorm:"column:supplier_code"`
	SupplierName  *string `json:"supplier_name,omitempty" gorm:"column:supplier_name"`
}

func (ReplenishmentStrategy) TableName() string {
	return "procurement_replenishment_strategy"
}

type ReplenishmentPlanStatus string

const (
	ReplenishmentPlanPending   ReplenishmentPlanStatus = "PENDING"
	ReplenishmentPlanConverted ReplenishmentPlanStatus = "CONVERTED"
	ReplenishmentPlanCancelled ReplenishmentPlanStatus = "CANCELLED"
)

type ReplenishmentPlan struct {
	ID               uint64                  `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	PlanDate         time.Time               `json:"plan_date" gorm:"column:plan_date;type:date;not null"`
	ProductID        uint64                  `json:"product_id" gorm:"column:product_id;not null"`
	WarehouseID      uint64                  `json:"warehouse_id" gorm:"column:warehouse_id;not null"`
	SupplierID       *uint64                 `json:"supplier_id" gorm:"column:supplier_id"`
	StrategyID       *uint64                 `json:"strategy_id" gorm:"column:strategy_id"`
	DailyDemand      float64                 `json:"daily_demand" gorm:"column:daily_demand;type:decimal(18,4);not null"`
	DemandWindowDays uint32                  `json:"demand_window_days" gorm:"column:demand_window_days;not null"`
	CoverageDays     uint32                  `json:"coverage_days" gorm:"column:coverage_days;not null"`
	NetSupply        int64                   `json:"net_supply" gorm:"column:net_supply;not null"`
	TargetStock      uint64                  `json:"target_stock" gorm:"column:target_stock;not null"`
	ShortageQty      uint64                  `json:"shortage_qty" gorm:"column:shortage_qty;not null"`
	SuggestedQty     uint64                  `json:"suggested_qty" gorm:"column:suggested_qty;not null"`
	MOQ              uint32                  `json:"moq" gorm:"column:moq;not null"`
	OrderMultiple    uint32                  `json:"order_multiple" gorm:"column:order_multiple;not null"`
	UnitCost         *float64                `json:"unit_cost" gorm:"column:unit_cost;type:decimal(12,4)"`
	Status           ReplenishmentPlanStatus `json:"status" gorm:"column:status;type:enum('PENDING','CONVERTED','CANCELLED');not null"`
	PurchaseOrderID  *uint64                 `json:"purchase_order_id" gorm:"column:purchase_order_id"`
	ConvertedAt      *time.Time              `json:"converted_at" gorm:"column:converted_at"`
	Remark           *string                 `json:"remark" gorm:"column:remark;size:255"`
	GmtCreate        time.Time               `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time               `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`

	SellerSKU            string  `json:"seller_sku,omitempty" gorm:"column:seller_sku"`
	ProductTitle         string  `json:"product_title,omitempty" gorm:"column:product_title"`
	ProductImageURL      *string `json:"product_image_url,omitempty" gorm:"column:product_image_url"`
	WarehouseCode        string  `json:"warehouse_code,omitempty" gorm:"column:warehouse_code"`
	WarehouseName        string  `json:"warehouse_name,omitempty" gorm:"column:warehouse_name"`
	SupplierCode         *string `json:"supplier_code,omitempty" gorm:"column:supplier_code"`
	SupplierName         *string `json:"supplier_name,omitempty" gorm:"column:supplier_name"`
	StrategyName         *string `json:"strategy_name,omitempty" gorm:"column:strategy_name"`
	PurchaseOrderNumber  *string `json:"purchase_order_number,omitempty" gorm:"column:purchase_order_number"`
	PurchaseOrderNumbers string  `json:"purchase_order_numbers,omitempty" gorm:"-"`
	PackagingShortageQty uint64  `json:"packaging_shortage_qty,omitempty" gorm:"-"`
	PackagingAlert       string  `json:"packaging_alert,omitempty" gorm:"-"`
}

func (ReplenishmentPlan) TableName() string {
	return "procurement_replenishment_plan"
}

type ReplenishmentStrategyListParams struct {
	Page     int
	PageSize int
	Keyword  string
}

type ReplenishmentPlanListParams struct {
	Page     int
	PageSize int
	Date     *time.Time
	Status   ReplenishmentPlanStatus
}

type ReplenishmentPlanConvertParams struct {
	PlanIDs    []uint64
	PlanDate   *time.Time
	OperatorID *uint64
}

type ReplenishmentPlanPurchaseOrderLink struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	PlanID          uint64    `json:"plan_id" gorm:"column:plan_id;not null"`
	PurchaseOrderID uint64    `json:"purchase_order_id" gorm:"column:purchase_order_id;not null"`
	GmtCreate       time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
}

func (ReplenishmentPlanPurchaseOrderLink) TableName() string {
	return "procurement_replenishment_plan_purchase_order"
}
