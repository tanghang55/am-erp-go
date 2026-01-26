package domain

import "time"

type ShipmentStatus string

const (
	ShipmentStatusDraft      ShipmentStatus = "DRAFT"      // 草稿
	ShipmentStatusConfirmed  ShipmentStatus = "CONFIRMED"  // 已确认(库存已锁定)
	ShipmentStatusPacked     ShipmentStatus = "PACKED"     // 已打包(库存已扣减到待出)
	ShipmentStatusShipped    ShipmentStatus = "SHIPPED"    // 已发货(库存已扣减到在途)
	ShipmentStatusDelivered  ShipmentStatus = "DELIVERED"  // 已送达
	ShipmentStatusCancelled  ShipmentStatus = "CANCELLED"  // 已取消
)

type DestinationType string

const (
	DestinationTypePlatformWarehouse DestinationType = "PLATFORM_WAREHOUSE" // 平台仓库
	DestinationTypeCustomer          DestinationType = "CUSTOMER"           // 客户
	DestinationTypeOwnWarehouse      DestinationType = "OWN_WAREHOUSE"      // 自有仓库
	DestinationTypeSupplier          DestinationType = "SUPPLIER"           // 供应商
	DestinationTypeOther             DestinationType = "OTHER"              // 其他
)

type Shipment struct {
	ID             uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	ShipmentNumber string         `json:"shipment_number" gorm:"column:shipment_number;size:50;uniqueIndex;not null"`

	// 订单信息（通用）
	OrderNumber    *string        `json:"order_number" gorm:"column:order_number;size:100;index"`
	SalesChannel   *string        `json:"sales_channel" gorm:"column:sales_channel;size:50"`

	// 发货仓库
	WarehouseID    uint64         `json:"warehouse_id" gorm:"column:warehouse_id;not null;index"`

	// 收货方信息（通用）
	DestinationType    *DestinationType `json:"destination_type" gorm:"column:destination_type;type:enum('PLATFORM_WAREHOUSE','CUSTOMER','OWN_WAREHOUSE','SUPPLIER','OTHER');default:'PLATFORM_WAREHOUSE'"`
	DestinationName    *string          `json:"destination_name" gorm:"column:destination_name;size:200"`
	DestinationContact *string          `json:"destination_contact" gorm:"column:destination_contact;size:100"`
	DestinationPhone   *string          `json:"destination_phone" gorm:"column:destination_phone;size:50"`
	DestinationAddress *string          `json:"destination_address" gorm:"column:destination_address;type:text"`
	DestinationCode    *string          `json:"destination_code" gorm:"column:destination_code;size:50"`

	// 物流信息（通用）
	Carrier        *string        `json:"carrier" gorm:"column:carrier;size:50"`
	ShippingMethod *string        `json:"shipping_method" gorm:"column:shipping_method;size:50"`
	TrackingNumber *string        `json:"tracking_number" gorm:"column:tracking_number;size:200;index"`

	// 包装信息
	BoxCount    uint    `json:"box_count" gorm:"column:box_count;default:0"`
	TotalWeight float64 `json:"total_weight" gorm:"column:total_weight;type:decimal(10,2);default:0"`
	TotalVolume float64 `json:"total_volume" gorm:"column:total_volume;type:decimal(10,3);default:0"`

	// 费用
	ShippingCost float64 `json:"shipping_cost" gorm:"column:shipping_cost;type:decimal(12,4);default:0"`
	Currency     string  `json:"currency" gorm:"column:currency;size:10;default:'USD'"`

	// 时间节点
	ShipDate             *string    `json:"ship_date" gorm:"column:ship_date;type:date"`
	ExpectedDeliveryDate *string    `json:"expected_delivery_date" gorm:"column:expected_delivery_date;type:date"`
	ActualDeliveryDate   *string    `json:"actual_delivery_date" gorm:"column:actual_delivery_date;type:date"`

	// 状态（简化）
	Status ShipmentStatus `json:"status" gorm:"column:status;type:enum('DRAFT','CONFIRMED','PACKED','SHIPPED','DELIVERED','CANCELLED');default:'DRAFT';index"`

	// 库存标记
	InventoryLocked   bool `json:"inventory_locked" gorm:"column:inventory_locked;default:0"`
	InventoryDeducted bool `json:"inventory_deducted" gorm:"column:inventory_deducted;default:0"`

	// 备注
	Remark        *string `json:"remark" gorm:"column:remark;type:text"`
	InternalNotes *string `json:"internal_notes" gorm:"column:internal_notes;type:text"`

	// 审计
	CreatedBy   *uint64   `json:"created_by" gorm:"column:created_by"`
	UpdatedBy   *uint64   `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate   time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`

	// Relations
	Items     []ShipmentItem `json:"items,omitempty" gorm:"-"`
	Warehouse interface{}    `json:"warehouse,omitempty" gorm:"-"` // 仓库信息
}

func (Shipment) TableName() string {
	return "shipment"
}

type ShipmentItem struct {
	ID               uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ShipmentID       uint64    `json:"shipment_id" gorm:"column:shipment_id;not null;index"`
	SkuID            uint64    `json:"sku_id" gorm:"column:sku_id;not null;index"`

	// 数量
	QuantityPlanned  uint      `json:"quantity_planned" gorm:"column:quantity_planned;not null"`
	QuantityShipped  uint      `json:"quantity_shipped" gorm:"column:quantity_shipped;default:0"`

	// 装箱信息
	PackageSpecID    *uint64   `json:"package_spec_id" gorm:"column:package_spec_id;index"`
	BoxQuantity      uint      `json:"box_quantity" gorm:"column:box_quantity;default:0"` // 装箱数量

	// 成本（保留字段，前端可不展示）
	UnitCost         float64   `json:"unit_cost" gorm:"column:unit_cost;type:decimal(12,4);default:0"`
	Currency         string    `json:"currency" gorm:"column:currency;size:10;default:'USD'"`

	// 备注
	Remark           *string   `json:"remark" gorm:"column:remark;size:500"`

	GmtCreate        time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`

	// Relations (not stored in DB)
	Sku              interface{}  `json:"sku,omitempty" gorm:"-"`         // 产品信息
	PackageSpec      *PackageSpec `json:"package_spec,omitempty" gorm:"-"` // 装箱规格
}

func (ShipmentItem) TableName() string {
	return "shipment_item"
}

// Request/Response types

type ShipmentListParams struct {
	Page           int
	PageSize       int
	Status         *ShipmentStatus
	WarehouseID    *uint64
	OrderNumber    *string
	TrackingNumber *string
	Keyword        *string
	DateFrom       *string
	DateTo         *string
}

type CreateShipmentParams struct {
	OrderNumber *string
	WarehouseID uint64
	Items       []CreateShipmentItemParams
	Remark      *string
	OperatorID  *uint64
}

type CreateShipmentItemParams struct {
	SkuID           uint64
	QuantityPlanned uint
	PackageSpecID   *uint64
	BoxQuantity     *uint
	UnitCost        *float64
	Currency        *string
	Remark          *string
}

type ConfirmShipmentParams struct {
	OperatorID *uint64
}

type PackShipmentParams struct {
	OperatorID *uint64
}

type MarkShippedParams struct {
	Carrier        *string
	TrackingNumber *string
	ShippingCost   *float64
	Currency       *string
	ShipDate       *string
	Remark         *string
	OperatorID     *uint64
}

type MarkDeliveredParams struct {
	ActualDeliveryDate *string
	Remark             *string
	OperatorID         *uint64
}

type CancelShipmentParams struct {
	Remark     *string
	OperatorID *uint64
}

// Repository interface
type ShipmentRepository interface {
	Create(shipment *Shipment) error
	Update(shipment *Shipment) error
	GetByID(id uint64) (*Shipment, error)
	GetByShipmentNumber(shipmentNumber string) (*Shipment, error)
	List(params *ShipmentListParams) ([]*Shipment, int64, error)
	Delete(id uint64) error
}

type ShipmentItemRepository interface {
	Create(item *ShipmentItem) error
	CreateBatch(items []ShipmentItem) error
	GetByShipmentID(shipmentID uint64) ([]ShipmentItem, error)
	DeleteByShipmentID(shipmentID uint64) error
}
