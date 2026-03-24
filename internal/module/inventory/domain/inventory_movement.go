package domain

import "time"

type MovementType string
type StockPool string

const (
	MovementTypePurchaseReceipt     MovementType = "PURCHASE_RECEIPT"
	MovementTypeSalesShipment       MovementType = "SALES_SHIPMENT"
	MovementTypeSalesAllocate       MovementType = "SALES_ALLOCATE"
	MovementTypeSalesRelease        MovementType = "SALES_RELEASE"
	MovementTypeSalesShip           MovementType = "SALES_SHIP"
	MovementTypeStockTakeAdjustment MovementType = "STOCK_TAKE_ADJUSTMENT"
	MovementTypeManualAdjustment    MovementType = "MANUAL_ADJUSTMENT"
	MovementTypeDamageWriteOff      MovementType = "DAMAGE_WRITE_OFF"
	MovementTypeReturnReceipt       MovementType = "RETURN_RECEIPT"
	MovementTypeTransferOut         MovementType = "TRANSFER_OUT"
	MovementTypeTransferIn          MovementType = "TRANSFER_IN"

	// 新增库存状态流转类型
	MovementTypePurchaseShip        MovementType = "PURCHASE_SHIP"         // 供应商发货 → 采购在途
	MovementTypeWarehouseReceive    MovementType = "WAREHOUSE_RECEIVE"     // 到仓收货: 采购在途 → 待检
	MovementTypeInspectionPass      MovementType = "INSPECTION_PASS"       // 质检通过: 待检 → 原料库存
	MovementTypeInspectionFail      MovementType = "INSPECTION_FAIL"       // 质检不合格: 待检 → 损坏
	MovementTypeInspectionLoss      MovementType = "INSPECTION_LOSS"       // 采购质检损失: 待检 → 损失结案
	MovementTypeAssemblyConsume     MovementType = "ASSEMBLY_CONSUME"      // 组装耗料: 子件原料库存 → 组装消耗
	MovementTypeAssemblyComplete    MovementType = "ASSEMBLY_COMPLETE"     // 组装完成: 主件产出 → 待出库存
	MovementTypePackingSkipComplete MovementType = "PACKING_SKIP_COMPLETE" // 免打包直通: 原料库存 → 待出库存
	MovementTypeShipmentAllocate    MovementType = "SHIPMENT_ALLOCATE"     // 发货单确认: 待出 → 待出锁定
	MovementTypeShipmentRelease     MovementType = "SHIPMENT_RELEASE"      // 发货单取消: 待出锁定 → 待出
	MovementTypeLogisticsShip       MovementType = "LOGISTICS_SHIP"        // 物流发货: 待出库存 → 物流在途
	MovementTypePlatformReceive     MovementType = "PLATFORM_RECEIVE"      // 平台上架: 物流在途 → 可售库存
	MovementTypeReturnInspect       MovementType = "RETURN_INSPECT"        // 退货质检: 退货库存 → 待检/损坏

	// 发货单库存流转类型
	MovementTypeShipmentShip MovementType = "SHIPMENT_SHIP" // 发货单发货: 待出锁定 → 在途
)

const (
	StockPoolAvailable StockPool = "AVAILABLE"
	StockPoolSellable  StockPool = "SELLABLE"
)

type InventoryMovement struct {
	ID                            uint64       `json:"id" gorm:"primaryKey;column:id"`
	TraceID                       *string      `json:"trace_id" gorm:"column:trace_id;size:50;index"`
	ProductID                     uint64       `json:"product_id" gorm:"column:product_id;not null;index"`
	WarehouseID                   uint64       `json:"warehouse_id" gorm:"column:warehouse_id;not null;index"`
	MovementType                  MovementType `json:"movement_type" gorm:"column:movement_type;type:enum('PURCHASE_RECEIPT','SALES_SHIPMENT','SALES_ALLOCATE','SALES_RELEASE','SALES_SHIP','STOCK_TAKE_ADJUSTMENT','MANUAL_ADJUSTMENT','DAMAGE_WRITE_OFF','RETURN_RECEIPT','TRANSFER_OUT','TRANSFER_IN','PURCHASE_SHIP','WAREHOUSE_RECEIVE','INSPECTION_PASS','INSPECTION_FAIL','INSPECTION_LOSS','ASSEMBLY_CONSUME','ASSEMBLY_COMPLETE','PACKING_SKIP_COMPLETE','SHIPMENT_ALLOCATE','SHIPMENT_RELEASE','LOGISTICS_SHIP','PLATFORM_RECEIVE','RETURN_INSPECT','SHIPMENT_SHIP');not null;index"`
	ReferenceType                 *string      `json:"reference_type" gorm:"column:reference_type;size:50;index"`
	ReferenceID                   *uint64      `json:"reference_id" gorm:"column:reference_id;index"`
	ReferenceNumber               *string      `json:"reference_number" gorm:"column:reference_number;size:100;index"`
	StockPool                     *StockPool   `json:"stock_pool" gorm:"column:stock_pool;type:enum('AVAILABLE','SELLABLE')"`
	Quantity                      int          `json:"quantity" gorm:"column:quantity;not null"`
	BeforeAvailable               uint         `json:"before_available" gorm:"column:before_available;not null"`
	AfterAvailable                uint         `json:"after_available" gorm:"column:after_available;not null"`
	BeforeReserved                uint         `json:"before_reserved" gorm:"column:before_reserved;not null"`
	AfterReserved                 uint         `json:"after_reserved" gorm:"column:after_reserved;not null"`
	BeforePurchasingInTransit     uint         `json:"before_purchasing_in_transit" gorm:"column:before_purchasing_in_transit;not null;default:0"`
	AfterPurchasingInTransit      uint         `json:"after_purchasing_in_transit" gorm:"column:after_purchasing_in_transit;not null;default:0"`
	BeforePendingInspection       uint         `json:"before_pending_inspection" gorm:"column:before_pending_inspection;not null;default:0"`
	AfterPendingInspection        uint         `json:"after_pending_inspection" gorm:"column:after_pending_inspection;not null;default:0"`
	BeforeRawMaterial             uint         `json:"before_raw_material" gorm:"column:before_raw_material;not null;default:0"`
	AfterRawMaterial              uint         `json:"after_raw_material" gorm:"column:after_raw_material;not null;default:0"`
	BeforeSellable                uint         `json:"before_sellable" gorm:"column:before_sellable;not null;default:0"`
	AfterSellable                 uint         `json:"after_sellable" gorm:"column:after_sellable;not null;default:0"`
	BeforeSellableReserved        uint         `json:"before_sellable_reserved" gorm:"column:before_sellable_reserved;not null;default:0"`
	AfterSellableReserved         uint         `json:"after_sellable_reserved" gorm:"column:after_sellable_reserved;not null;default:0"`
	BeforePendingShipment         uint         `json:"before_pending_shipment" gorm:"column:before_pending_shipment;not null;default:0"`
	AfterPendingShipment          uint         `json:"after_pending_shipment" gorm:"column:after_pending_shipment;not null;default:0"`
	BeforePendingShipmentReserved uint         `json:"before_pending_shipment_reserved" gorm:"column:before_pending_shipment_reserved;not null;default:0"`
	AfterPendingShipmentReserved  uint         `json:"after_pending_shipment_reserved" gorm:"column:after_pending_shipment_reserved;not null;default:0"`
	BeforeLogisticsInTransit      uint         `json:"before_logistics_in_transit" gorm:"column:before_logistics_in_transit;not null;default:0"`
	AfterLogisticsInTransit       uint         `json:"after_logistics_in_transit" gorm:"column:after_logistics_in_transit;not null;default:0"`
	BeforeDamaged                 uint         `json:"before_damaged" gorm:"column:before_damaged;not null"`
	AfterDamaged                  uint         `json:"after_damaged" gorm:"column:after_damaged;not null"`
	UnitCost                      *float64     `json:"unit_cost" gorm:"column:unit_cost;type:decimal(12,4)"`
	TotalCost                     *float64     `json:"total_cost" gorm:"column:total_cost;type:decimal(12,4)"`
	Remark                        *string      `json:"remark" gorm:"column:remark;type:text"`
	OperatorID                    *uint64      `json:"operator_id" gorm:"column:operator_id;index"`
	OperatedAt                    time.Time    `json:"operated_at" gorm:"column:operated_at;not null;index"`
	GmtCreate                     time.Time    `json:"created_at" gorm:"column:gmt_create"`
	GmtModified                   time.Time    `json:"updated_at" gorm:"column:gmt_modified"`

	// Associations (not stored in DB)
	Product   *ProductSnapshot   `json:"product,omitempty" gorm:"-"`
	Warehouse *WarehouseSnapshot `json:"warehouse,omitempty" gorm:"-"`
	Operator  *OperatorSnapshot  `json:"operator,omitempty" gorm:"-"`
	// FIFO 扣减明细（仅运行时返回，不落库）
	LotAllocations []InventoryLotAllocation `json:"lot_allocations,omitempty" gorm:"-"`
}

func (InventoryMovement) TableName() string {
	return "inventory_movement"
}

type ProductSnapshot struct {
	ID        uint64 `json:"id"`
	SellerSku string `json:"seller_sku"`
	Title     string `json:"title"`
	Asin      string `json:"asin,omitempty"`
}

type WarehouseSnapshot struct {
	ID   uint64 `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type OperatorSnapshot struct {
	ID       uint64  `json:"id"`
	Username string  `json:"username"`
	RealName *string `json:"real_name"`
}

type InventoryLotAllocation struct {
	InventoryLotID uint64  `json:"inventory_lot_id"`
	Qty            uint64  `json:"qty"`
	UnitCost       float64 `json:"unit_cost"`
}

type MovementListParams struct {
	Page         int
	PageSize     int
	ProductID    *uint64
	WarehouseID  *uint64
	MovementType *MovementType
	DateFrom     *string
	DateTo       *string
}

type CreateMovementParams struct {
	TraceID         *string
	ProductID       uint64
	WarehouseID     uint64
	MovementType    MovementType
	Quantity        int
	ReferenceType   *string
	ReferenceID     *uint64
	ReferenceNumber *string
	StockPool       *StockPool
	UnitCost        *float64
	Remark          *string
	OperatorID      *uint64
	OperatedAt      *time.Time
}

type TransferParams struct {
	ProductID       uint64
	FromWarehouseID uint64
	ToWarehouseID   uint64
	Quantity        uint
	UnitCost        *float64
	Remark          *string
	OperatorID      *uint64
	ReferenceType   *string
	ReferenceNumber *string
}

// 库存状态转换参数
type StockTransitionParams struct {
	ProductID       uint64
	WarehouseID     uint64
	Quantity        uint
	UnitCost        *float64
	Remark          *string
	OperatorID      *uint64
	ReferenceType   *string
	ReferenceID     *uint64
	ReferenceNumber *string
}

// 退货质检参数
type ReturnInspectParams struct {
	ProductID       uint64
	WarehouseID     uint64
	PassQuantity    uint // 质检通过数量 → 待检库存
	FailQuantity    uint // 质检不合格数量 → 损坏库存
	Remark          *string
	OperatorID      *uint64
	ReferenceType   *string
	ReferenceNumber *string
}
