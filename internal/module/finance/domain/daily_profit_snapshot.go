package domain

import "time"

type DailyProfitSnapshotStatus string

const (
	DailyProfitSnapshotStatusNormal       DailyProfitSnapshotStatus = "NORMAL"
	DailyProfitSnapshotStatusRecalculated DailyProfitSnapshotStatus = "RECALCULATED"
)

type DailyProfitSnapshot struct {
	ID                       uint64                    `json:"id" gorm:"primaryKey;autoIncrement"`
	BizDate                  time.Time                 `json:"biz_date" gorm:"column:biz_date;type:date;not null"`
	Marketplace              string                    `json:"marketplace" gorm:"column:marketplace;size:10;not null"`
	BaseCurrency             string                    `json:"base_currency" gorm:"column:base_currency;size:10;not null"`
	SalesIncomeAmount        float64                   `json:"sales_income_amount" gorm:"column:sales_income_amount;type:decimal(18,6);not null"`
	COGSAmount               float64                   `json:"cogs_amount" gorm:"column:cogs_amount;type:decimal(18,6);not null"`
	GrossProfitAmount        float64                   `json:"gross_profit_amount" gorm:"column:gross_profit_amount;type:decimal(18,6);not null"`
	OrderExpenseAmount       float64                   `json:"order_expense_amount" gorm:"column:order_expense_amount;type:decimal(18,6);not null"`
	OrderNetProfitAmount     float64                   `json:"order_net_profit_amount" gorm:"column:order_net_profit_amount;type:decimal(18,6);not null"`
	PublicExpenseAmount      float64                   `json:"public_expense_amount" gorm:"column:public_expense_amount;type:decimal(18,6);not null"`
	OperatingNetProfitAmount float64                   `json:"operating_net_profit_amount" gorm:"column:operating_net_profit_amount;type:decimal(18,6);not null"`
	OrderCount               uint64                    `json:"order_count" gorm:"column:order_count;not null"`
	ShippedQty               uint64                    `json:"shipped_qty" gorm:"column:shipped_qty;not null"`
	SnapshotStatus           DailyProfitSnapshotStatus `json:"snapshot_status" gorm:"column:snapshot_status;type:enum('NORMAL','RECALCULATED');not null"`
	SourceVersion            string                    `json:"source_version" gorm:"column:source_version;size:32;not null"`
	BuiltAt                  time.Time                 `json:"built_at" gorm:"column:built_at;not null"`
	BuilderID                *uint64                   `json:"builder_id" gorm:"column:builder_id"`
	GmtCreate                time.Time                 `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified              time.Time                 `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (DailyProfitSnapshot) TableName() string {
	return "finance_daily_profit_snapshot"
}

type DailyProfitSnapshotListParams struct {
	DateFrom    time.Time
	DateTo      time.Time
	Marketplace string
}
