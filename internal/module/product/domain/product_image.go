package domain

import "time"

// ProductImage 产品图片实体
type ProductImage struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID   uint64    `json:"product_id" gorm:"column:product_id;not null"`
	ImageUrl    string    `json:"image_url" gorm:"column:image_url;size:500;not null"`
	SortOrder   uint32    `json:"sort_order" gorm:"column:sort_order"`
	IsPrimary   uint8     `json:"is_primary" gorm:"column:is_primary"`
	Status      string    `json:"status" gorm:"size:20"`
	Remark      string    `json:"remark" gorm:"type:text"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ProductImage) TableName() string {
	return "product_image"
}
