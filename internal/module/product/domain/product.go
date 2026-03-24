package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Product struct {
	ID                   uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SellerSku            string    `json:"seller_sku" gorm:"column:seller_sku;size:100;not null"`
	Asin                 string    `json:"asin" gorm:"size:20;not null"`
	Title                string    `json:"title" gorm:"size:500;not null"`
	Fnsku                string    `json:"fnsku" gorm:"size:20"`
	Marketplace          string    `json:"marketplace" gorm:"type:enum('US','CA','AU','UK','DE','JP');not null"`
	ParentID             *uint64   `json:"parent_id" gorm:"column:parent_id"`
	ComboID              *uint64   `json:"combo_id" gorm:"column:combo_id"`
	IsComboMain          uint8     `json:"is_combo_main" gorm:"column:is_combo_main;default:0"`
	SupplierID           *uint64   `json:"supplier_id" gorm:"column:supplier_id"`
	SupplierName         string    `json:"supplier_name" gorm:"column:supplier_name;->"`
	SupplierCode         string    `json:"supplier_code" gorm:"column:supplier_code;->"`
	BrandID              *uint64   `json:"brand_id" gorm:"column:brand_id"`
	BrandName            string    `json:"brand_name" gorm:"column:brand_name;->"`
	CategoryID           *uint64   `json:"category_id" gorm:"column:category_id"`
	CategoryName         string    `json:"category_name" gorm:"column:category_name;->"`
	DimensionUnitID      *uint64   `json:"dimension_unit_id" gorm:"column:dimension_unit_id"`
	DimensionUnitName    string    `json:"dimension_unit_name" gorm:"column:dimension_unit_name;->"`
	WeightUnitID         *uint64   `json:"weight_unit_id" gorm:"column:weight_unit_id"`
	WeightUnitName       string    `json:"weight_unit_name" gorm:"column:weight_unit_name;->"`
	IsInspectionRequired uint8     `json:"is_inspection_required" gorm:"column:is_inspection_required;not null;default:1"`
	IsPackingRequired    uint8     `json:"is_packing_required" gorm:"column:is_packing_required;not null;default:1"`
	InventoryAvailable   uint      `json:"inventory_available" gorm:"column:inventory_available;->"`
	InventoryReserved    uint      `json:"inventory_reserved" gorm:"column:inventory_reserved;->"`
	InventoryInbound     uint      `json:"inventory_inbound" gorm:"column:inventory_inbound;->"`
	UnitCost             *float64  `json:"unit_cost" gorm:"column:unit_cost;type:decimal(15,4)"`
	Weight               *float64  `json:"weight" gorm:"type:decimal(10,2)"`
	Length               *float64  `json:"length" gorm:"column:length;type:decimal(10,2)"`
	Width                *float64  `json:"width" gorm:"column:width;type:decimal(10,2)"`
	Height               *float64  `json:"height" gorm:"column:height;type:decimal(10,2)"`
	Dimensions           string    `json:"dimensions" gorm:"size:100"`
	Status               string    `json:"status" gorm:"type:enum('DRAFT','ON_SALE','REPLENISHING','OFF_SHELF');default:'DRAFT'"`
	ImageUrl             string    `json:"image_url" gorm:"column:image_url;size:500"`
	Images               JSONArray `json:"images" gorm:"type:json"`
	Remark               string    `json:"remark" gorm:"type:text"`
	ReferenceCount       int64     `json:"reference_count" gorm:"-"`
	Deletable            bool      `json:"deletable" gorm:"-"`
	DeleteBlockReason    string    `json:"delete_block_reason,omitempty" gorm:"-"`
	UpdatedByName        string    `json:"updated_by_name" gorm:"column:updated_by_name;->"`
	GmtCreate            time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified          time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

const (
	ProductStatusDraft        = "DRAFT"
	ProductStatusOnSale       = "ON_SALE"
	ProductStatusReplenishing = "REPLENISHING"
	ProductStatusOffShelf     = "OFF_SHELF"
)

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

	ChildCount         int64     `json:"child_count" gorm:"column:child_count;->"`
	ActiveChildCount   int64     `json:"active_child_count" gorm:"column:active_child_count;->"`
	InactiveChildCount int64     `json:"inactive_child_count" gorm:"column:inactive_child_count;->"`
	Children           []Product `json:"children,omitempty" gorm:"-"`
}

func (ProductParent) TableName() string {
	return "product_parent"
}

type ProductConfigType string

const (
	ProductConfigTypeBrand         ProductConfigType = "BRAND"
	ProductConfigTypeSalesStatus   ProductConfigType = "SALES_STATUS"
	ProductConfigTypeDimensionUnit ProductConfigType = "DIMENSION_UNIT"
	ProductConfigTypeWeightUnit    ProductConfigType = "WEIGHT_UNIT"
)

type ProductConfigItem struct {
	ID                uint64            `json:"id" gorm:"primaryKey;autoIncrement"`
	ConfigType        ProductConfigType `json:"config_type" gorm:"column:config_type;size:30;not null;index"`
	ItemCode          string            `json:"item_code" gorm:"column:item_code;size:50;not null"`
	ItemName          string            `json:"item_name" gorm:"column:item_name;size:100;not null"`
	Status            string            `json:"status" gorm:"column:status;type:enum('ACTIVE','INACTIVE');default:'ACTIVE'"`
	Sort              int               `json:"sort" gorm:"column:sort;not null;default:0"`
	Remark            string            `json:"remark" gorm:"column:remark;size:500"`
	ReferenceCount    int64             `json:"reference_count" gorm:"-"`
	Deletable         bool              `json:"deletable" gorm:"-"`
	DeleteBlockReason string            `json:"delete_block_reason,omitempty" gorm:"-"`
	GmtCreate         time.Time         `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified       time.Time         `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ProductConfigItem) TableName() string {
	return "product_config_item"
}

type ProductCategory struct {
	ID                uint64             `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID          *uint64            `json:"parent_id" gorm:"column:parent_id"`
	CategoryCode      string             `json:"category_code" gorm:"column:category_code;size:50;not null"`
	CategoryName      string             `json:"category_name" gorm:"column:category_name;size:100;not null"`
	Level             uint8              `json:"level" gorm:"column:level;not null"`
	Status            string             `json:"status" gorm:"column:status;type:enum('ACTIVE','INACTIVE');default:'ACTIVE'"`
	Sort              int                `json:"sort" gorm:"column:sort;not null;default:0"`
	Remark            string             `json:"remark" gorm:"column:remark;size:500"`
	ReferenceCount    int64              `json:"reference_count" gorm:"-"`
	Deletable         bool               `json:"deletable" gorm:"-"`
	DeleteBlockReason string             `json:"delete_block_reason,omitempty" gorm:"-"`
	GmtCreate         time.Time          `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified       time.Time          `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
	Children          []*ProductCategory `json:"children,omitempty" gorm:"-"`
}

func (ProductCategory) TableName() string {
	return "product_category"
}

// ProductPackagingItem 产品包材关联实体
type ProductPackagingItem struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID       uint64    `json:"product_id" gorm:"column:product_id;not null"`
	PackagingItemID uint64    `json:"packaging_item_id" gorm:"column:packaging_item_id;not null"`
	QuantityPerUnit float64   `json:"quantity_per_unit" gorm:"column:quantity_per_unit;type:decimal(10,3);not null"`
	CreatedBy       *uint64   `json:"created_by" gorm:"column:created_by"`
	UpdatedBy       *uint64   `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate       time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified     time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`

	// 关联包材详情（查询时关联）
	PackagingItem *PackagingItemDetail `json:"packaging_item,omitempty" gorm:"-"`
}

func (ProductPackagingItem) TableName() string {
	return "product_packaging_items"
}

// PackagingItemDetail 包材详情（用于关联查询）
type PackagingItemDetail struct {
	ID             uint64  `json:"id"`
	ItemCode       string  `json:"item_code"`
	ItemName       string  `json:"item_name"`
	Specification  string  `json:"specification"`
	Unit           string  `json:"unit"`
	UnitCost       float64 `json:"unit_cost"`
	Currency       string  `json:"currency"`
	QuantityOnHand uint64  `json:"quantity_on_hand"`
}
