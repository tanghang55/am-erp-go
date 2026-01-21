package domain

import "time"

type AuditLog struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID    string    `json:"trace_id" gorm:"column:trace_id;size:64"`
	UserID     *uint64   `json:"user_id" gorm:"column:user_id"`
	Username   string    `json:"username" gorm:"column:username;size:50"`
	Module     string    `json:"module" gorm:"column:module;size:50"`
	Action     string    `json:"action" gorm:"column:action;size:100"`
	EntityType string    `json:"entity_type" gorm:"column:entity_type;size:50"`
	EntityID   string    `json:"entity_id" gorm:"column:entity_id;size:100"`
	Changes    string    `json:"changes" gorm:"column:changes;type:json"`
	IPAddress  string    `json:"ip_address" gorm:"column:ip_address;size:45"`
	UserAgent  string    `json:"user_agent" gorm:"column:user_agent;type:text"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}

type AuditLogListParams struct {
	Page       int
	PageSize   int
	Module     string
	Action     string
	Username   string
	EntityType string
	EntityID   string
	Keyword    string
	DateFrom   string
	DateTo     string
}
