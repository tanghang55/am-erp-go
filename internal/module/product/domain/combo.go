package domain

import "time"

type ProductComboItem struct {
	ID            uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ComboID       uint64    `json:"combo_id" gorm:"column:combo_id;not null"`
	MainProductID uint64    `json:"main_product_id" gorm:"column:main_product_id;not null"`
	ProductID     uint64    `json:"product_id" gorm:"column:product_id;not null"`
	QtyRatio      uint64    `json:"qty_ratio" gorm:"column:qty_ratio;not null"`
	GmtCreate     time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified   time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ProductComboItem) TableName() string {
	return "product_combo"
}

type ProductCombo struct {
	ComboID     uint64    `json:"combo_id"`
	MainProduct Product   `json:"main_product"`
	Products    []Product `json:"products"`
}

type ComboListParams struct {
	Page     int
	PageSize int
}

type ComboUpsertParams struct {
	MainProductID uint64
	ProductIDs    []uint64
	QtyRatios     map[uint64]uint64
}
