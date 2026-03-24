package domain

import (
	"context"
	"time"
)

type SyncTrigger string

const (
	SyncTriggerManual    SyncTrigger = "MANUAL"
	SyncTriggerScheduled SyncTrigger = "SCHEDULED"
)

type OrderSyncRunStatus string

const (
	OrderSyncRunStatusRunning OrderSyncRunStatus = "RUNNING"
	OrderSyncRunStatusSuccess OrderSyncRunStatus = "SUCCESS"
	OrderSyncRunStatusPartial OrderSyncRunStatus = "PARTIAL_SUCCESS"
	OrderSyncRunStatusFailed  OrderSyncRunStatus = "FAILED"
)

type OrderSyncState struct {
	ID               uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Provider         string     `json:"provider" gorm:"column:provider;size:32;not null"`
	Channel          string     `json:"channel" gorm:"column:channel;size:32;not null"`
	LastUpdatedAfter *time.Time `json:"last_updated_after" gorm:"column:last_updated_after"`
	LastSyncStarted  *time.Time `json:"last_sync_started_at" gorm:"column:last_sync_started_at"`
	LastSyncFinished *time.Time `json:"last_sync_finished_at" gorm:"column:last_sync_finished_at"`
	GmtCreate        time.Time  `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time  `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (OrderSyncState) TableName() string {
	return "third_party_order_sync_state"
}

type OrderSyncRun struct {
	ID                      uint64             `json:"id" gorm:"primaryKey;autoIncrement"`
	Provider                string             `json:"provider" gorm:"column:provider;size:32;not null"`
	Channel                 string             `json:"channel" gorm:"column:channel;size:32;not null"`
	TriggerType             SyncTrigger        `json:"trigger_type" gorm:"column:trigger_type;size:16;not null"`
	Status                  OrderSyncRunStatus `json:"status" gorm:"column:status;size:20;not null"`
	RequestLastUpdatedAfter *time.Time         `json:"request_last_updated_after" gorm:"column:request_last_updated_after"`
	BatchNo                 *string            `json:"batch_no" gorm:"column:batch_no;size:64"`
	FetchedOrders           uint32             `json:"fetched_orders" gorm:"column:fetched_orders;not null"`
	FetchedItems            uint32             `json:"fetched_items" gorm:"column:fetched_items;not null"`
	ImportedItems           uint32             `json:"imported_items" gorm:"column:imported_items;not null"`
	ErrorItems              uint32             `json:"error_items" gorm:"column:error_items;not null"`
	Message                 *string            `json:"message" gorm:"column:message;size:500"`
	StartedAt               time.Time          `json:"started_at" gorm:"column:started_at;not null"`
	FinishedAt              *time.Time         `json:"finished_at" gorm:"column:finished_at"`
	GmtCreate               time.Time          `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified             time.Time          `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (OrderSyncRun) TableName() string {
	return "third_party_order_sync_run"
}

type ListRunsParams struct {
	Page     int
	PageSize int
}

type ExternalListOrdersRequest struct {
	MarketplaceIDs   []string
	LastUpdatedAfter time.Time
}

type ExternalOrder struct {
	OrderID       string
	MarketplaceID string
	PurchaseAt    time.Time
	LastUpdatedAt time.Time
	Currency      string
	Items         []ExternalOrderItem
}

type ExternalOrderItem struct {
	OrderItemID string
	SellerSKU   string
	Quantity    uint64
	Amount      float64
}

type OrdersProvider interface {
	Code() string
	ListOrders(ctx context.Context, req ExternalListOrdersRequest) ([]ExternalOrder, error)
}

type OrderSyncRepository interface {
	GetState(provider string, channel string) (*OrderSyncState, error)
	SaveState(state *OrderSyncState) error
	CreateRun(run *OrderSyncRun) error
	UpdateRun(run *OrderSyncRun) error
	ListRuns(provider string, channel string, params *ListRunsParams) ([]OrderSyncRun, int64, error)
}
