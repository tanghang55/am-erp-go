package domain

import "time"

type Role struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:50;not null;uniqueIndex"`
	DisplayName string    `json:"display_name" gorm:"column:display_name;size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`

	Permissions []Permission `json:"permissions" gorm:"many2many:role_permission;"`
}

func (Role) TableName() string {
	return "role"
}

// IsAdmin 检查是否为管理员角色
func (r *Role) IsAdmin() bool {
	return r.Name == "admin"
}

type RolePermission struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	RoleID       uint64    `gorm:"column:role_id;not null"`
	PermissionID uint64    `gorm:"column:permission_id;not null"`
	GmtCreate    time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified  time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (RolePermission) TableName() string {
	return "role_permission"
}
