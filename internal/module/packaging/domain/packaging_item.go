package domain

import "time"

// PackagingItem 包材实体
type PackagingItem struct {
	ID               uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID          string    `json:"trace_id" gorm:"column:trace_id;size:64;not null"`
	ItemCode         string    `json:"item_code" gorm:"column:item_code;size:50;not null;uniqueIndex"`
	ItemName         string    `json:"item_name" gorm:"column:item_name;size:100;not null"`
	Category         string    `json:"category" gorm:"size:50;not null"`
	Specification    *string   `json:"specification" gorm:"size:200"`
	UnitCost         float64   `json:"unit_cost" gorm:"column:unit_cost;type:decimal(10,4);default:0"`
	Currency         string    `json:"currency" gorm:"size:10;default:'CNY'"`
	Unit             string    `json:"unit" gorm:"size:20;default:'PCS'"`
	QuantityOnHand   uint64    `json:"quantity_on_hand" gorm:"column:quantity_on_hand;default:0"`
	ReorderPoint     *uint64   `json:"reorder_point" gorm:"column:reorder_point"`
	ReorderQuantity  *uint64   `json:"reorder_quantity" gorm:"column:reorder_quantity"`
	SupplierName     *string   `json:"supplier_name" gorm:"column:supplier_name;size:100"`
	SupplierContact  *string   `json:"supplier_contact" gorm:"column:supplier_contact;size:100"`
	Status           string    `json:"status" gorm:"size:20;default:'ACTIVE'"`
	Notes            *string   `json:"notes" gorm:"type:text"`
	CreatedBy        uint64    `json:"created_by" gorm:"column:created_by;not null"`
	GmtCreate        time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (PackagingItem) TableName() string {
	return "packaging_item"
}

// PackagingItemRepository 包材仓储接口
type PackagingItemRepository interface {
	List(params *PackagingItemListParams) ([]PackagingItem, int64, error)
	GetByID(id uint64) (*PackagingItem, error)
	Create(item *PackagingItem) error
	Update(item *PackagingItem) error
	Delete(id uint64) error
	GetLowStockItems() ([]PackagingItem, error)
	UpdateQuantity(id uint64, quantity int64) error
}

// PackagingItemListParams 查询参数
type PackagingItemListParams struct {
	Page      int
	PageSize  int
	Keyword   string
	Category  string
	Status    string
	LowStock  bool
}
