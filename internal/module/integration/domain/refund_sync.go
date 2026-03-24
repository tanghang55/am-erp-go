package domain

import (
	"context"
	"time"
)

type RefundEventStatus string

const (
	RefundEventStatusMapped   RefundEventStatus = "MAPPED"
	RefundEventStatusUnmapped RefundEventStatus = "UNMAPPED"
)

type RefundSyncState struct {
	ID               uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Provider         string     `json:"provider" gorm:"column:provider;size:32;not null"`
	Channel          string     `json:"channel" gorm:"column:channel;size:32;not null"`
	LastPostedAfter  *time.Time `json:"last_posted_after" gorm:"column:last_posted_after"`
	LastSyncStarted  *time.Time `json:"last_sync_started_at" gorm:"column:last_sync_started_at"`
	LastSyncFinished *time.Time `json:"last_sync_finished_at" gorm:"column:last_sync_finished_at"`
	GmtCreate        time.Time  `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time  `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (RefundSyncState) TableName() string {
	return "third_party_refund_sync_state"
}

type RefundSyncRun struct {
	ID                 uint64             `json:"id" gorm:"primaryKey;autoIncrement"`
	Provider           string             `json:"provider" gorm:"column:provider;size:32;not null"`
	Channel            string             `json:"channel" gorm:"column:channel;size:32;not null"`
	TriggerType        SyncTrigger        `json:"trigger_type" gorm:"column:trigger_type;size:16;not null"`
	Status             OrderSyncRunStatus `json:"status" gorm:"column:status;size:20;not null"`
	RequestPostedAfter *time.Time         `json:"request_posted_after" gorm:"column:request_posted_after"`
	FetchedRefunds     uint32             `json:"fetched_refunds" gorm:"column:fetched_refunds;not null"`
	ImportedRefunds    uint32             `json:"imported_refunds" gorm:"column:imported_refunds;not null"`
	ErrorRefunds       uint32             `json:"error_refunds" gorm:"column:error_refunds;not null"`
	Message            *string            `json:"message" gorm:"column:message;size:500"`
	StartedAt          time.Time          `json:"started_at" gorm:"column:started_at;not null"`
	FinishedAt         *time.Time         `json:"finished_at" gorm:"column:finished_at"`
	GmtCreate          time.Time          `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified        time.Time          `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (RefundSyncRun) TableName() string {
	return "third_party_refund_sync_run"
}

type ThirdPartyRefundEvent struct {
	ID           uint64            `json:"id" gorm:"primaryKey;autoIncrement"`
	Provider     string            `json:"provider" gorm:"column:provider;size:32;not null"`
	Channel      string            `json:"channel" gorm:"column:channel;size:32;not null"`
	RefundID     string            `json:"refund_id" gorm:"column:refund_id;size:128;not null"`
	OrderID      string            `json:"order_id" gorm:"column:order_id;size:64;not null"`
	OrderItemID  *string           `json:"order_item_id" gorm:"column:order_item_id;size:64"`
	SellerSKU    string            `json:"seller_sku" gorm:"column:seller_sku;size:100;not null"`
	Marketplace  string            `json:"marketplace" gorm:"column:marketplace;size:10;not null"`
	ProductID    *uint64           `json:"product_id" gorm:"column:product_id"`
	QtyRefunded  uint64            `json:"qty_refunded" gorm:"column:qty_refunded;not null"`
	RefundAmount float64           `json:"refund_amount" gorm:"column:refund_amount;type:decimal(18,4);not null"`
	Currency     string            `json:"currency" gorm:"column:currency;size:3;not null"`
	PostedAt     time.Time         `json:"posted_at" gorm:"column:posted_at;not null"`
	Status       RefundEventStatus `json:"status" gorm:"column:status;type:enum('MAPPED','UNMAPPED');not null"`
	ErrorMessage *string           `json:"error_message" gorm:"column:error_message;size:500"`
	RawPayload   *string           `json:"raw_payload" gorm:"column:raw_payload;type:text"`
	GmtCreate    time.Time         `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified  time.Time         `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ThirdPartyRefundEvent) TableName() string {
	return "third_party_refund_event"
}

type ExternalListRefundsRequest struct {
	MarketplaceIDs []string
	PostedAfter    time.Time
}

type ExternalRefund struct {
	RefundID      string
	OrderID       string
	OrderItemID   string
	SellerSKU     string
	MarketplaceID string
	Quantity      uint64
	Amount        float64
	Currency      string
	PostedAt      time.Time
}

type RefundsProvider interface {
	Code() string
	ListRefunds(ctx context.Context, req ExternalListRefundsRequest) ([]ExternalRefund, error)
}

type RefundSyncRepository interface {
	GetState(provider string, channel string) (*RefundSyncState, error)
	SaveState(state *RefundSyncState) error
	CreateRun(run *RefundSyncRun) error
	UpdateRun(run *RefundSyncRun) error
	ListRuns(provider string, channel string, params *ListRunsParams) ([]RefundSyncRun, int64, error)
	UpsertEvents(events []ThirdPartyRefundEvent) error
}
