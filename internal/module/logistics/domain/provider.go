package domain

import "time"

type ProviderType string
type ProviderStatus string
type TransportMode string

const (
	ProviderTypeFreightForwarder ProviderType = "FREIGHT_FORWARDER" // 货代
	ProviderTypeCourier          ProviderType = "COURIER"           // 快递
	ProviderTypeShippingLine     ProviderType = "SHIPPING_LINE"     // 船公司
	ProviderTypeAirline          ProviderType = "AIRLINE"           // 航空公司
)

const (
	ProviderStatusActive   ProviderStatus = "ACTIVE"
	ProviderStatusInactive ProviderStatus = "INACTIVE"
)

const (
	TransportModeExpress TransportMode = "EXPRESS" // 快递
	TransportModeAir     TransportMode = "AIR"     // 空运
	TransportModeSea     TransportMode = "SEA"     // 海运
	TransportModeRail    TransportMode = "RAIL"    // 铁路
	TransportModeTruck   TransportMode = "TRUCK"   // 卡车
)

type LogisticsProvider struct {
	ID            uint64         `json:"id" gorm:"primaryKey;column:id"`
	ProviderCode  string         `json:"provider_code" gorm:"column:provider_code;uniqueIndex;size:50;not null"`
	ProviderName  string         `json:"provider_name" gorm:"column:provider_name;size:200;not null"`
	ProviderType  ProviderType   `json:"provider_type" gorm:"column:provider_type;type:enum('FREIGHT_FORWARDER','COURIER','SHIPPING_LINE','AIRLINE');not null"`
	ServiceTypes  *string        `json:"service_types" gorm:"column:service_types;size:200"`
	ContactPerson *string        `json:"contact_person" gorm:"column:contact_person;size:100"`
	ContactPhone  *string        `json:"contact_phone" gorm:"column:contact_phone;size:50"`
	ContactEmail  *string        `json:"contact_email" gorm:"column:contact_email;size:100"`
	Website       *string        `json:"website" gorm:"column:website;size:200"`
	Country       *string        `json:"country" gorm:"column:country;size:50"`
	City          *string        `json:"city" gorm:"column:city;size:100"`
	Address       *string        `json:"address" gorm:"column:address;type:text"`
	AccountNumber *string        `json:"account_number" gorm:"column:account_number;size:100"`
	CreditDays    int            `json:"credit_days" gorm:"column:credit_days;default:0"`
	Status        ProviderStatus `json:"status" gorm:"column:status;type:enum('ACTIVE','INACTIVE');default:ACTIVE"`
	Remark        *string        `json:"remark" gorm:"column:remark;type:text"`
	CreatedBy     *uint64        `json:"created_by" gorm:"column:created_by"`
	UpdatedBy     *uint64        `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate     time.Time      `json:"created_at" gorm:"column:gmt_create"`
	GmtModified   time.Time      `json:"updated_at" gorm:"column:gmt_modified"`
}

func (LogisticsProvider) TableName() string {
	return "logistics_provider"
}

type LogisticsProviderListParams struct {
	Page         int
	PageSize     int
	ProviderType *ProviderType
	Status       *ProviderStatus
	Keyword      *string
}

type CreateProviderParams struct {
	ProviderCode  string         `json:"provider_code" binding:"required"`
	ProviderName  string         `json:"provider_name" binding:"required"`
	ProviderType  ProviderType   `json:"provider_type" binding:"required"`
	ServiceTypes  *string        `json:"service_types"`
	ContactPerson *string        `json:"contact_person"`
	ContactPhone  *string        `json:"contact_phone"`
	ContactEmail  *string        `json:"contact_email"`
	Website       *string        `json:"website"`
	Country       *string        `json:"country"`
	City          *string        `json:"city"`
	Address       *string        `json:"address"`
	AccountNumber *string        `json:"account_number"`
	CreditDays    *int           `json:"credit_days"`
	Status        ProviderStatus `json:"status"`
	Remark        *string        `json:"remark"`
	OperatorID    *uint64        `json:"operator_id"`
}

type UpdateProviderParams struct {
	ProviderCode  *string         `json:"provider_code"`
	ProviderName  *string         `json:"provider_name"`
	ProviderType  *ProviderType   `json:"provider_type"`
	ServiceTypes  *string         `json:"service_types"`
	ContactPerson *string         `json:"contact_person"`
	ContactPhone  *string         `json:"contact_phone"`
	ContactEmail  *string         `json:"contact_email"`
	Website       *string         `json:"website"`
	Country       *string         `json:"country"`
	City          *string         `json:"city"`
	Address       *string         `json:"address"`
	AccountNumber *string         `json:"account_number"`
	CreditDays    *int            `json:"credit_days"`
	Status        *ProviderStatus `json:"status"`
	Remark        *string         `json:"remark"`
	OperatorID    *uint64         `json:"operator_id"`
}
