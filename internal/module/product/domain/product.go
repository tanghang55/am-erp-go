package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Product SKU实体
type Product struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SellerSku   string    `json:"seller_sku" gorm:"column:seller_sku;size:100;not null"`
	Asin        string    `json:"asin" gorm:"size:20;not null"`
	Title       string    `json:"title" gorm:"size:500;not null"`
	Fnsku       string    `json:"fnsku" gorm:"size:20"`
	Marketplace string    `json:"marketplace" gorm:"type:enum('US','CA','AU','UK','DE','JP');not null"`
	ParentID    *uint64   `json:"parent_id" gorm:"column:parent_id"`
	ComboID     *uint64   `json:"combo_id" gorm:"column:combo_id"`
	IsComboMain uint8     `json:"is_combo_main" gorm:"column:is_combo_main;default:0"`
	SupplierID  *uint64   `json:"supplier_id" gorm:"column:supplier_id"`
	UnitCost    *float64  `json:"unit_cost" gorm:"column:unit_cost;type:decimal(15,4)"`
	Weight      *float64  `json:"weight" gorm:"type:decimal(10,2)"`
	Dimensions  string    `json:"dimensions" gorm:"size:100"`
	Status      string    `json:"status" gorm:"type:enum('ACTIVE','INACTIVE','DISCONTINUED');default:'ACTIVE'"`
	ImageUrl    string    `json:"image_url" gorm:"column:image_url;size:500"`
	Images      JSONArray `json:"images" gorm:"type:json"`
	Remark      string    `json:"remark" gorm:"type:text"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (Product) TableName() string {
	return "product"
}

// JSONArray 用于处理JSON数组字段
type JSONArray []string

func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// ProductParent 产品父体实体
type ProductParent struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentAsin  string    `json:"parent_asin" gorm:"column:parent_asin;size:20;not null"`
	Title       string    `json:"title" gorm:"size:500;not null"`
	Marketplace string    `json:"marketplace" gorm:"type:enum('US','CA','AU','UK','DE','JP');not null"`
	Brand       string    `json:"brand" gorm:"size:100"`
	Category    string    `json:"category" gorm:"size:200"`
	Status      string    `json:"status" gorm:"type:enum('ACTIVE','INACTIVE','DISCONTINUED');default:'ACTIVE'"`
	ImageUrl    string    `json:"image_url" gorm:"column:image_url;size:500"`
	Remark      string    `json:"remark" gorm:"type:text"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ProductParent) TableName() string {
	return "product_parent"
}

// ProductPackagingItem 产品包材关联实体
type ProductPackagingItem struct {
	ID               uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID        uint64    `json:"product_id" gorm:"column:product_id;not null"`
	PackagingItemID  uint64    `json:"packaging_item_id" gorm:"column:packaging_item_id;not null"`
	QuantityPerUnit  float64   `json:"quantity_per_unit" gorm:"column:quantity_per_unit;type:decimal(10,3);not null"`
	Notes            string    `json:"notes" gorm:"size:500"`
	CreatedBy        *uint64   `json:"created_by" gorm:"column:created_by"`
	UpdatedBy        *uint64   `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate        time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`

	// 关联包材详情（查询时关联）
	PackagingItem    *PackagingItemDetail `json:"packaging_item,omitempty" gorm:"-"`
}

func (ProductPackagingItem) TableName() string {
	return "product_packaging_items"
}

// PackagingItemDetail 包材详情（用于关联查询）
type PackagingItemDetail struct {
	ID              uint64  `json:"id"`
	ItemCode        string  `json:"item_code"`
	ItemName        string  `json:"item_name"`
	Specification   string  `json:"specification"`
	Unit            string  `json:"unit"`
	UnitCost        float64 `json:"unit_cost"`
	Currency        string  `json:"currency"`
	QuantityOnHand  uint64  `json:"quantity_on_hand"`
}
