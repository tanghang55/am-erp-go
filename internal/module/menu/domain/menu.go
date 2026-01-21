package domain

import "time"

type Menu struct {
	ID             uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          string    `json:"title" gorm:"size:100;not null"`
	TitleEn        string    `json:"title_en" gorm:"column:title_en;size:100"`
	Code           string    `json:"code" gorm:"size:50;not null;uniqueIndex"`
	ParentID       *uint64   `json:"parent_id" gorm:"column:parent_id"`
	Path           string    `json:"path" gorm:"size:200"`
	Icon           string    `json:"icon" gorm:"size:100"`
	Component      string    `json:"component" gorm:"size:200"`
	Sort           uint      `json:"sort" gorm:"not null"`
	IsHidden       uint8     `json:"is_hidden" gorm:"column:is_hidden;default:0"`
	PermissionCode string    `json:"permission_code" gorm:"column:permission_code;size:100"`
	Status         string    `json:"status" gorm:"size:20;default:'ACTIVE'"`
	GmtCreate      time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified    time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`

	Children []Menu `json:"children" gorm:"-"`
}

func (Menu) TableName() string {
	return "menu"
}

func (m *Menu) IsActive() bool {
	return m.Status == "ACTIVE"
}

func (m *Menu) IsVisible() bool {
	return m.IsHidden == 0
}

type MenuTree struct {
	ID             uint64      `json:"id"`
	ParentID       *uint64     `json:"parent_id"`
	Title          string      `json:"title"`
	TitleEn        string      `json:"title_en,omitempty"`
	Code           string      `json:"code"`
	Path           string      `json:"path"`
	Component      string      `json:"component"`
	Icon           string      `json:"icon"`
	Sort           uint        `json:"sort"`
	IsHidden       uint8       `json:"is_hidden"`
	PermissionCode string      `json:"permission_code,omitempty"`
	Children       []*MenuTree `json:"children,omitempty"`
}

type MenuListParams struct {
	Page     int
	PageSize int
	Keyword  string
	Status   string
	IsHidden *uint8
	ParentID *uint64
}

type MenuListItem struct {
	Menu
	ParentTitle string `json:"parent_title" gorm:"-"`
	FullPath    string `json:"full_path" gorm:"-"`
}
