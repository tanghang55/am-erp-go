package domain

import (
	"time"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
)

type PricingMethod string
type RateStatus string

const (
	PricingMethodPerKg      PricingMethod = "PER_KG"      // 按公斤
	PricingMethodPerCbm     PricingMethod = "PER_CBM"     // 按立方米
	PricingMethodPerPackage PricingMethod = "PER_PACKAGE" // 按件
	PricingMethodFixed      PricingMethod = "FIXED"       // 固定价格
)

const (
	RateStatusActive   RateStatus = "ACTIVE"
	RateStatusInactive RateStatus = "INACTIVE"
	RateStatusExpired  RateStatus = "EXPIRED"
)

type ShippingRate struct {
	ID                     uint64        `json:"id" gorm:"primaryKey;column:id"`
	ProviderID             uint64        `json:"provider_id" gorm:"column:provider_id;not null;index"`
	OriginWarehouseID      uint64        `json:"origin_warehouse_id" gorm:"column:origin_warehouse_id;not null;index"`
	DestinationWarehouseID uint64        `json:"destination_warehouse_id" gorm:"column:destination_warehouse_id;not null;index"`
	TransportMode          TransportMode `json:"transport_mode" gorm:"column:transport_mode;type:enum('EXPRESS','AIR','SEA','RAIL','TRUCK');not null;index"`
	ServiceID              *uint64       `json:"service_id" gorm:"column:service_id;index"`
	PricingMethod          PricingMethod `json:"pricing_method" gorm:"column:pricing_method;type:enum('PER_KG','PER_CBM','PER_PACKAGE','FIXED');not null"`
	BaseRate               float64       `json:"base_rate" gorm:"column:base_rate;type:decimal(10,2);not null"`
	Currency               string        `json:"currency" gorm:"column:currency;size:10;default:CNY"`
	MinWeight              *float64      `json:"min_weight" gorm:"column:min_weight;type:decimal(10,2)"`
	TransitDays            *int          `json:"transit_days" gorm:"column:transit_days"`
	EffectiveDate          string        `json:"effective_date" gorm:"column:effective_date;type:date;not null"`
	ExpiryDate             *string       `json:"expiry_date" gorm:"column:expiry_date;type:date"`
	Status                 RateStatus    `json:"status" gorm:"column:status;type:enum('ACTIVE','INACTIVE','EXPIRED');default:ACTIVE"`
	Remark                 *string       `json:"remark" gorm:"column:remark;type:text"`
	CreatedBy              *uint64       `json:"created_by" gorm:"column:created_by"`
	UpdatedBy              *uint64       `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate              time.Time     `json:"created_at" gorm:"column:gmt_create"`
	GmtModified            time.Time     `json:"updated_at" gorm:"column:gmt_modified"`
	OtherFee               float64      `json:"other_fee" gorm:"column:other_fee;type:decimal(10,2);not null"`

	// Relations
	Provider             *LogisticsProvider          `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
	OriginWarehouse      *inventoryDomain.Warehouse  `json:"origin_warehouse,omitempty" gorm:"foreignKey:OriginWarehouseID"`
	DestinationWarehouse *inventoryDomain.Warehouse  `json:"destination_warehouse,omitempty" gorm:"foreignKey:DestinationWarehouseID"`
	Service              *LogisticsService           `json:"service,omitempty" gorm:"foreignKey:ServiceID"`
}

func (ShippingRate) TableName() string {
	return "shipping_rate"
}

type ShippingRateListParams struct {
	Page                   int
	PageSize               int
	ProviderID             *uint64
	OriginWarehouseID      *uint64
	DestinationWarehouseID *uint64
	TransportMode          *TransportMode
	Status                 *RateStatus
	Keyword                *string
}

type CreateShippingRateParams struct {
	ProviderID             uint64        `json:"provider_id" binding:"required"`
	OriginWarehouseID      uint64        `json:"origin_warehouse_id" binding:"required"`
	DestinationWarehouseID uint64        `json:"destination_warehouse_id" binding:"required"`
	TransportMode          TransportMode `json:"transport_mode" binding:"required"`
	ServiceID              *uint64       `json:"service_id"`
	PricingMethod          PricingMethod `json:"pricing_method" binding:"required"`
	BaseRate               float64       `json:"base_rate" binding:"required"`
	Currency               string        `json:"currency"`
	MinWeight              *float64      `json:"min_weight"`
	TransitDays            *int          `json:"transit_days"`
	EffectiveDate          string        `json:"effective_date" binding:"required"`
	ExpiryDate             *string       `json:"expiry_date"`
	Status                 RateStatus    `json:"status"`
	Remark                 *string       `json:"remark"`
	OperatorID             *uint64       `json:"operator_id"`
	OtherFee               float64       `json:"other_fee"`

}

type UpdateShippingRateParams struct {
	ProviderID             *uint64        `json:"provider_id"`
	OriginWarehouseID      *uint64        `json:"origin_warehouse_id"`
	DestinationWarehouseID *uint64        `json:"destination_warehouse_id"`
	TransportMode          *TransportMode `json:"transport_mode"`
	ServiceID              *uint64        `json:"service_id"`
	PricingMethod          *PricingMethod `json:"pricing_method"`
	BaseRate               *float64       `json:"base_rate"`
	Currency               *string        `json:"currency"`
	MinWeight              *float64       `json:"min_weight"`
	TransitDays            *int           `json:"transit_days"`
	EffectiveDate          *string        `json:"effective_date"`
	ExpiryDate             *string        `json:"expiry_date"`
	Status                 *RateStatus    `json:"status"`
	Remark                 *string        `json:"remark"`
	OperatorID             *uint64        `json:"operator_id"`
	OtherFee               float64       `json:"other_fee"`

}

// QueryLatestRateParams 查询最新报价的参数
type QueryLatestRateParams struct {
	ProviderID             *uint64
	OriginWarehouseID      uint64
	DestinationWarehouseID uint64
	TransportMode          TransportMode
	Weight                 *float64 // 用于匹配重量区间
	Volume                 *float64 // 用于匹配体积区间
	QueryDate              *string  // 查询日期，默认为当前日期
}
