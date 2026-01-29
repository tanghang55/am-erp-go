package domain

import "time"

// PackagingLedger 包材流水实体
type PackagingLedger struct {
	ID               uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID          string         `json:"trace_id" gorm:"column:trace_id;size:64;not null"`
	PackagingItemID  uint64         `json:"packaging_item_id" gorm:"column:packaging_item_id;not null"`
	TransactionType  string         `json:"transaction_type" gorm:"column:transaction_type;size:20;not null"`
	Quantity         int64          `json:"quantity" gorm:"not null"`
	UnitCost         float64        `json:"unit_cost" gorm:"column:unit_cost;type:decimal(10,4);default:0"`
	TotalCost        float64        `json:"total_cost" gorm:"column:total_cost;->"` // 只读，数据库计算列
	QuantityBefore   uint64         `json:"quantity_before" gorm:"column:quantity_before;not null"`
	QuantityAfter    uint64         `json:"quantity_after" gorm:"column:quantity_after;not null"`
	ReferenceType    *string        `json:"reference_type" gorm:"column:reference_type;size:50"`
	ReferenceID      *uint64        `json:"reference_id" gorm:"column:reference_id"`
	OccurredAt       time.Time      `json:"occurred_at" gorm:"column:occurred_at;not null"`
	Notes            *string        `json:"notes" gorm:"type:text"`
	CreatedBy        uint64         `json:"created_by" gorm:"column:created_by;not null"`
	GmtCreate        time.Time      `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`

	// 关联
	PackagingItem    *PackagingItem `json:"packaging_item,omitempty" gorm:"foreignKey:PackagingItemID"`
}

func (PackagingLedger) TableName() string {
	return "packaging_ledger"
}

// PackagingLedgerRepository 流水仓储接口
type PackagingLedgerRepository interface {
	List(params *PackagingLedgerListParams) ([]PackagingLedger, int64, error)
	GetByID(id uint64) (*PackagingLedger, error)
	Create(ledger *PackagingLedger) error
	GetUsageSummary(dateFrom, dateTo *time.Time) ([]UsageSummaryItem, error)
}

// PackagingLedgerListParams 查询参数
type PackagingLedgerListParams struct {
	Page             int
	PageSize         int
	PackagingItemID  *uint64
	TransactionType  string
	DateFrom         *time.Time
	DateTo           *time.Time
	ReferenceType    string
	ReferenceID      *uint64
}

// UsageSummaryItem 使用情况统计项
type UsageSummaryItem struct {
	ID         uint64  `json:"id"`
	ItemCode   string  `json:"item_code"`
	ItemName   string  `json:"item_name"`
	Category   string  `json:"category"`
	TotalIn    uint64  `json:"total_in"`
	TotalOut   uint64  `json:"total_out"`
	TotalCost  float64 `json:"total_cost"`
}
