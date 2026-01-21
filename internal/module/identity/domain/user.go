package domain

import "time"

type User struct {
	ID           uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string     `json:"username" gorm:"size:50;not null;uniqueIndex"`
	PasswordHash string     `json:"-" gorm:"column:password_hash;size:255;not null"`
	RealName     string     `json:"real_name" gorm:"column:real_name;size:100"`
	Email        string     `json:"email" gorm:"size:100"`
	Phone        string     `json:"phone" gorm:"size:20"`
	Status       string     `json:"status" gorm:"type:enum('ACTIVE','DISABLED');default:'ACTIVE'"`
	LastLoginAt  *time.Time `json:"last_login_at" gorm:"column:last_login_at"`
	GmtCreate    time.Time  `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified  time.Time  `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`

	Roles []Role `json:"roles" gorm:"many2many:user_role;"`
}

func (User) TableName() string {
	return "user"
}

// IsActive 检查用户是否启用
func (u *User) IsActive() bool {
	return u.Status == "ACTIVE"
}

type UserRole struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	UserID      uint64    `gorm:"column:user_id;not null"`
	RoleID      uint64    `gorm:"column:role_id;not null"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (UserRole) TableName() string {
	return "user_role"
}
