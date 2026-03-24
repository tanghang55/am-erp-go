package domain

import "time"

type CostType string

const (
	CostTypePurchase CostType = "PURCHASE"
	CostTypeLanded   CostType = "LANDED"
	CostTypeAverage  CostType = "AVERAGE"
)

type CostingSnapshot struct {
	ID              uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID         string     `json:"trace_id" gorm:"column:trace_id;size:64;not null"`
	ProductID       uint64     `json:"product_id" gorm:"column:product_id;not null"`
	SellerSKU       string     `json:"seller_sku" gorm:"->;column:seller_sku"`
	ProductTitle    string     `json:"product_title" gorm:"->;column:product_title"`
	ProductImageURL string     `json:"product_image_url" gorm:"->;column:product_image_url"`
	CostType        CostType   `json:"cost_type" gorm:"column:cost_type;type:enum('PURCHASE','LANDED','AVERAGE');not null"`
	UnitCost        float64    `json:"unit_cost" gorm:"column:unit_cost;type:decimal(15,4);not null"`
	Currency        string     `json:"currency" gorm:"column:currency;size:10;not null"`
	EffectiveFrom   time.Time  `json:"effective_from" gorm:"column:effective_from;not null"`
	EffectiveTo     *time.Time `json:"effective_to" gorm:"column:effective_to"`
	Notes           *string    `json:"notes" gorm:"column:notes;type:text"`
	CreatedBy       uint64     `json:"created_by" gorm:"column:created_by;not null"`
	GmtCreate       time.Time  `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified     time.Time  `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (CostingSnapshot) TableName() string {
	return "costing_snapshot"
}

type CostingSnapshotListParams struct {
	Page      int
	PageSize  int
	ProductID *uint64
	CostType  CostType
	IsCurrent *bool
}
