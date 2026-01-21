package domain

import "time"

type Permission struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Code        string    `json:"code" gorm:"size:100;not null;uniqueIndex"`
	Module      string    `json:"module" gorm:"size:50;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"type:enum('ACTIVE','DISABLED');default:'ACTIVE'"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (Permission) TableName() string {
	return "permission"
}

// IsActive 检查权限是否启用
func (p *Permission) IsActive() bool {
	return p.Status == "ACTIVE"
}
