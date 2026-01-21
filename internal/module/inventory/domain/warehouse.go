package domain

import "time"

type WarehouseType string
type WarehouseStatus string

const (
	WarehouseTypeFBA        WarehouseType = "FBA"
	WarehouseTypeThirdParty WarehouseType = "THIRD_PARTY"
	WarehouseTypeOwn        WarehouseType = "OWN"
)

const (
	WarehouseStatusActive   WarehouseStatus = "ACTIVE"
	WarehouseStatusInactive WarehouseStatus = "INACTIVE"
	WarehouseStatusClosed   WarehouseStatus = "CLOSED"
)

type Warehouse struct {
	ID            uint64          `json:"id" gorm:"primaryKey;column:id"`
	Code          string          `json:"code" gorm:"column:code;uniqueIndex;size:50;not null"`
	Name          string          `json:"name" gorm:"column:name;size:100;not null"`
	Type          WarehouseType   `json:"type" gorm:"column:type;type:enum('FBA','THIRD_PARTY','OWN');default:OWN"`
	Country       *string         `json:"country" gorm:"column:country;size:10"`
	Address       *string         `json:"address" gorm:"column:address;type:text"`
	ContactPerson *string         `json:"contact_person" gorm:"column:contact_person;size:100"`
	ContactPhone  *string         `json:"contact_phone" gorm:"column:contact_phone;size:50"`
	ContactEmail  *string         `json:"contact_email" gorm:"column:contact_email;size:100"`
	Status        WarehouseStatus `json:"status" gorm:"column:status;type:enum('ACTIVE','INACTIVE','CLOSED');default:ACTIVE"`
	Remark        *string         `json:"remark" gorm:"column:remark;type:text"`
	CreatedBy     *uint64         `json:"created_by" gorm:"column:created_by"`
	UpdatedBy     *uint64         `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate     time.Time       `json:"gmt_create" gorm:"column:gmt_create"`
	GmtModified   time.Time       `json:"gmt_modified" gorm:"column:gmt_modified"`
}

func (Warehouse) TableName() string {
	return "warehouse"
}
