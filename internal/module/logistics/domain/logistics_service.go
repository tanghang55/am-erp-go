package domain

import "time"

type ServiceStatus string

const (
	ServiceStatusActive   ServiceStatus = "ACTIVE"
	ServiceStatusInactive ServiceStatus = "INACTIVE"
)

type LogisticsService struct {
	ID                uint64        `json:"id" gorm:"primaryKey;column:id"`
	ServiceCode       string        `json:"service_code" gorm:"column:service_code;uniqueIndex;size:50;not null"`
	ServiceName       string        `json:"service_name" gorm:"column:service_name;size:100;not null"`
	TransportMode     TransportMode `json:"transport_mode" gorm:"column:transport_mode;type:enum('EXPRESS','AIR','SEA','RAIL','TRUCK');not null;index"`
	DestinationRegion *string       `json:"destination_region" gorm:"column:destination_region;size:100"`
	Description       *string       `json:"description" gorm:"column:description;type:text"`
	Status            ServiceStatus `json:"status" gorm:"column:status;type:enum('ACTIVE','INACTIVE');default:ACTIVE;index"`
	CreatedBy         *uint64       `json:"created_by" gorm:"column:created_by"`
	UpdatedBy         *uint64       `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate         time.Time     `json:"created_at" gorm:"column:gmt_create"`
	GmtModified       time.Time     `json:"updated_at" gorm:"column:gmt_modified"`
}

func (LogisticsService) TableName() string {
	return "logistics_service"
}

type LogisticsServiceRepository interface {
	Create(service *LogisticsService) error
	Update(service *LogisticsService) error
	Delete(id uint64) error
	GetByID(id uint64) (*LogisticsService, error)
	GetByCode(code string) (*LogisticsService, error)
	List(params *LogisticsServiceListParams) ([]*LogisticsService, int64, error)
	GetActiveServices() ([]*LogisticsService, error)
	GetServicesByTransportMode(transportMode TransportMode) ([]*LogisticsService, error)
}

type LogisticsServiceListParams struct {
	Page          int
	PageSize      int
	TransportMode *TransportMode
	Status        *ServiceStatus
	Keyword       *string
}

type CreateLogisticsServiceParams struct {
	ServiceCode       string        `json:"service_code" binding:"required"`
	ServiceName       string        `json:"service_name" binding:"required"`
	TransportMode     TransportMode `json:"transport_mode" binding:"required"`
	DestinationRegion *string       `json:"destination_region"`
	Description       *string       `json:"description"`
	Status            ServiceStatus `json:"status"`
	OperatorID        *uint64       `json:"operator_id"`
}

type UpdateLogisticsServiceParams struct {
	ServiceCode       *string        `json:"service_code"`
	ServiceName       *string        `json:"service_name"`
	TransportMode     *TransportMode `json:"transport_mode"`
	DestinationRegion *string        `json:"destination_region"`
	Description       *string        `json:"description"`
	Status            *ServiceStatus `json:"status"`
	OperatorID        *uint64        `json:"operator_id"`
}
