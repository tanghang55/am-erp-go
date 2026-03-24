package domain

import "time"

type ProfitLedgerType string

const (
	ProfitLedgerTypeIncome       ProfitLedgerType = "INCOME"
	ProfitLedgerTypeCOGS         ProfitLedgerType = "COGS"
	ProfitLedgerTypeOrderExpense ProfitLedgerType = "ORDER_EXPENSE"
	ProfitLedgerTypePublicExp    ProfitLedgerType = "PUBLIC_EXPENSE"
)

type ProfitLedgerStatus string

const (
	ProfitLedgerStatusNormal   ProfitLedgerStatus = "NORMAL"
	ProfitLedgerStatusReversed ProfitLedgerStatus = "REVERSED"
)

type ProfitLedger struct {
	ID               uint64             `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID          string             `json:"trace_id" gorm:"column:trace_id;size:64;not null"`
	LedgerType       ProfitLedgerType   `json:"ledger_type" gorm:"column:ledger_type;type:enum('INCOME','COGS','ORDER_EXPENSE','PUBLIC_EXPENSE');not null"`
	Status           ProfitLedgerStatus `json:"status" gorm:"column:status;type:enum('NORMAL','REVERSED');not null;default:'NORMAL'"`
	ReversalOfID     *uint64            `json:"reversal_of_id" gorm:"column:reversal_of_id"`
	BizDate          time.Time          `json:"biz_date" gorm:"column:biz_date;type:date;not null"`
	Marketplace      *string            `json:"marketplace" gorm:"column:marketplace;size:10"`
	SalesOrderID     *uint64            `json:"sales_order_id" gorm:"column:sales_order_id"`
	SalesOrderItemID *uint64            `json:"sales_order_item_id" gorm:"column:sales_order_item_id"`
	ReferenceType    *string            `json:"reference_type" gorm:"column:reference_type;size:50"`
	ReferenceID      *uint64            `json:"reference_id" gorm:"column:reference_id"`
	ReferenceNumber  *string            `json:"reference_number" gorm:"column:reference_number;size:100"`
	Category         string             `json:"category" gorm:"column:category;size:50;not null"`
	OriginalCurrency string             `json:"original_currency" gorm:"column:original_currency;size:10;not null"`
	OriginalAmount   float64            `json:"original_amount" gorm:"column:original_amount;type:decimal(18,6);not null"`
	BaseCurrency     string             `json:"base_currency" gorm:"column:base_currency;size:10;not null"`
	FxRate           float64            `json:"fx_rate" gorm:"column:fx_rate;type:decimal(18,8);not null"`
	BaseAmount       float64            `json:"base_amount" gorm:"column:base_amount;type:decimal(18,6);not null"`
	FxSource         string             `json:"fx_source" gorm:"column:fx_source;size:50;not null"`
	FxVersion        string             `json:"fx_version" gorm:"column:fx_version;size:32;not null"`
	FxTime           time.Time          `json:"fx_time" gorm:"column:fx_time;not null"`
	OccurredAt       time.Time          `json:"occurred_at" gorm:"column:occurred_at;not null"`
	OperatorID       *uint64            `json:"operator_id" gorm:"column:operator_id"`
	Remark           *string            `json:"remark" gorm:"column:remark;size:500"`
	GmtCreate        time.Time          `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time          `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ProfitLedger) TableName() string {
	return "finance_profit_ledger"
}
