package domain

import "time"

type ExchangeRateStatus string

const (
	ExchangeRateStatusActive   ExchangeRateStatus = "ACTIVE"
	ExchangeRateStatusInactive ExchangeRateStatus = "INACTIVE"
)

type ExchangeRateSource string

const (
	ExchangeRateSourceManual ExchangeRateSource = "MANUAL"
)

type ExchangeRate struct {
	ID            uint64             `json:"id" gorm:"primaryKey;autoIncrement"`
	FromCurrency  string             `json:"from_currency" gorm:"column:from_currency;size:10;not null"`
	ToCurrency    string             `json:"to_currency" gorm:"column:to_currency;size:10;not null"`
	Rate          float64            `json:"rate" gorm:"column:rate;type:decimal(18,8);not null"`
	SourceType    ExchangeRateSource `json:"source_type" gorm:"column:source_type;type:enum('MANUAL');not null"`
	SourceVersion string             `json:"source_version" gorm:"column:source_version;size:32;not null"`
	EffectiveAt   time.Time          `json:"effective_at" gorm:"column:effective_at;not null"`
	Status        ExchangeRateStatus `json:"status" gorm:"column:status;type:enum('ACTIVE','INACTIVE');not null"`
	Remark        *string            `json:"remark" gorm:"column:remark;type:text"`
	CreatedBy     uint64             `json:"created_by" gorm:"column:created_by;not null"`
	UpdatedBy     uint64             `json:"updated_by" gorm:"column:updated_by;not null"`
	GmtCreate     time.Time          `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified   time.Time          `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ExchangeRate) TableName() string {
	return "finance_exchange_rate"
}

type ExchangeRateListParams struct {
	Page         int
	PageSize     int
	FromCurrency string
	ToCurrency   string
	Status       ExchangeRateStatus
}
