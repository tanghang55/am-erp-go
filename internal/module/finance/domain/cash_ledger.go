package domain

import "time"

type LedgerType string

const (
	LedgerTypeIncome  LedgerType = "INCOME"
	LedgerTypeExpense LedgerType = "EXPENSE"
)

type CashLedgerStatus string

const (
	CashLedgerStatusNormal   CashLedgerStatus = "NORMAL"
	CashLedgerStatusReversed CashLedgerStatus = "REVERSED"
)

type CashLedger struct {
	ID               uint64           `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID          string           `json:"trace_id" gorm:"column:trace_id;size:64;not null"`
	LedgerType       LedgerType       `json:"ledger_type" gorm:"column:ledger_type;type:enum('INCOME','EXPENSE');not null"`
	Status           CashLedgerStatus `json:"status" gorm:"column:status;type:enum('NORMAL','REVERSED');not null;default:'NORMAL'"`
	ReversalOfID     *uint64          `json:"reversal_of_id" gorm:"column:reversal_of_id"`
	Category         string           `json:"category" gorm:"column:category;size:50;not null"`
	Amount           float64          `json:"amount" gorm:"column:amount;type:decimal(15,2);not null"`
	Currency         string           `json:"currency" gorm:"column:currency;size:10;not null"`
	OriginalCurrency string           `json:"original_currency" gorm:"column:original_currency;size:10;not null"`
	OriginalAmount   float64          `json:"original_amount" gorm:"column:original_amount;type:decimal(18,6);not null;default:0"`
	BaseCurrency     string           `json:"base_currency" gorm:"column:base_currency;size:10;not null"`
	FxRate           float64          `json:"fx_rate" gorm:"column:fx_rate;type:decimal(18,8);not null;default:1"`
	BaseAmount       float64          `json:"base_amount" gorm:"column:base_amount;type:decimal(18,6);not null;default:0"`
	FxSource         string           `json:"fx_source" gorm:"column:fx_source;size:50;not null;default:'STATIC'"`
	FxVersion        string           `json:"fx_version" gorm:"column:fx_version;size:32;not null;default:'v1'"`
	FxTime           time.Time        `json:"fx_time" gorm:"column:fx_time;not null"`
	Marketplace      *string          `json:"marketplace" gorm:"column:marketplace;size:10"`
	OccurredNode     *string          `json:"occurred_node" gorm:"column:occurred_node;size:32"`
	ReferenceType    *string          `json:"reference_type" gorm:"column:reference_type;size:50"`
	ReferenceID      *uint64          `json:"reference_id" gorm:"column:reference_id"`
	Description      *string          `json:"description" gorm:"column:description;type:text"`
	OccurredAt       time.Time        `json:"occurred_at" gorm:"column:occurred_at;not null"`
	CreatedBy        uint64           `json:"created_by" gorm:"column:created_by;not null"`
	CreatedByName    string           `json:"created_by_name,omitempty" gorm:"-"`
	GmtCreate        time.Time        `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time        `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (CashLedger) TableName() string {
	return "cash_ledger"
}

type CashLedgerListParams struct {
	Page          int
	PageSize      int
	LedgerType    LedgerType
	Category      string
	Marketplace   string
	OccurredNode  string
	Keyword       string
	DateFrom      *time.Time
	DateTo        *time.Time
	ReferenceType string
	ReferenceID   *uint64
}

type CashLedgerSummary struct {
	TotalIncome  float64 `json:"total_income"`
	IncomeCount  int64   `json:"income_count"`
	TotalExpense float64 `json:"total_expense"`
	ExpenseCount int64   `json:"expense_count"`
	NetProfit    float64 `json:"net_profit"`
}

type CategorySummaryItem struct {
	LedgerType  LedgerType `json:"ledger_type"`
	Category    string     `json:"category"`
	TotalAmount float64    `json:"total_amount"`
	Count       int64      `json:"count"`
}
