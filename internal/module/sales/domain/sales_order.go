package domain

import "time"

type SalesOrderStatus string
type StockPool string

const (
	SalesOrderStatusDraft     SalesOrderStatus = "DRAFT"
	SalesOrderStatusConfirmed SalesOrderStatus = "CONFIRMED"
	SalesOrderStatusAllocated SalesOrderStatus = "ALLOCATED"
	SalesOrderStatusShipped   SalesOrderStatus = "SHIPPED"
	SalesOrderStatusDelivered SalesOrderStatus = "DELIVERED"
	SalesOrderStatusCancelled SalesOrderStatus = "CANCELLED"
	SalesOrderStatusReturned  SalesOrderStatus = "RETURNED"
)

const (
	StockPoolAvailable StockPool = "AVAILABLE"
	StockPoolSellable  StockPool = "SELLABLE"
)

type SalesOrder struct {
	ID              uint64           `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo         string           `json:"order_no" gorm:"column:order_no;size:64;not null"`
	SourceType      string           `json:"source_type" gorm:"column:source_type;size:32;not null"`
	ExternalOrderNo *string          `json:"external_order_no" gorm:"column:external_order_no;size:100"`
	SalesChannel    *string          `json:"sales_channel" gorm:"column:sales_channel;size:50"`
	Marketplace     *string          `json:"marketplace" gorm:"column:marketplace;size:10"`
	StockPool       StockPool        `json:"stock_pool" gorm:"column:stock_pool;type:enum('AVAILABLE','SELLABLE');not null;default:'AVAILABLE'"`
	OrderStatus     SalesOrderStatus `json:"order_status" gorm:"column:order_status;size:20;not null"`
	OrderDate       time.Time        `json:"order_date" gorm:"column:order_date;not null"`
	ConfirmAt       *time.Time       `json:"confirm_at" gorm:"column:confirm_at"`
	AllocatedAt     *time.Time       `json:"allocated_at" gorm:"column:allocated_at"`
	ShippedAt       *time.Time       `json:"shipped_at" gorm:"column:shipped_at"`
	DeliveredAt     *time.Time       `json:"delivered_at" gorm:"column:delivered_at"`
	CancelledAt     *time.Time       `json:"cancelled_at" gorm:"column:cancelled_at"`
	Currency        string           `json:"currency" gorm:"column:currency;size:3;not null"`
	OrderAmount     float64          `json:"order_amount" gorm:"column:order_amount;type:decimal(18,4);not null"`
	Remark          *string          `json:"remark" gorm:"column:remark;type:text"`
	ImportBatchNo   *string          `json:"import_batch_no" gorm:"column:import_batch_no;size:64"`
	CreatedBy       *uint64          `json:"created_by" gorm:"column:created_by"`
	UpdatedBy       *uint64          `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate       time.Time        `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified     time.Time        `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
	Items           []SalesOrderItem `json:"items,omitempty" gorm:"-"`
}

func (SalesOrder) TableName() string {
	return "sales_order"
}

type SalesOrderItem struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SalesOrderID    uint64    `json:"sales_order_id" gorm:"column:sales_order_id;not null"`
	LineNo          uint32    `json:"line_no" gorm:"column:line_no;not null"`
	SourceType      string    `json:"source_type" gorm:"column:source_type;size:32;not null"`
	ExternalOrderNo *string   `json:"external_order_no" gorm:"column:external_order_no;size:100"`
	ProductID       uint64    `json:"product_id" gorm:"column:product_id;not null"`
	QtyOrdered      uint64    `json:"qty_ordered" gorm:"column:qty_ordered;not null"`
	QtyAllocated    uint64    `json:"qty_allocated" gorm:"column:qty_allocated;not null"`
	QtyShipped      uint64    `json:"qty_shipped" gorm:"column:qty_shipped;not null"`
	QtyReturned     uint64    `json:"qty_returned" gorm:"column:qty_returned;not null"`
	UnitPrice       float64   `json:"unit_price" gorm:"column:unit_price;type:decimal(18,4);not null"`
	Subtotal        float64   `json:"subtotal" gorm:"column:subtotal;type:decimal(18,4);not null"`
	Remark          *string   `json:"remark" gorm:"column:remark;type:text"`
	GmtCreate       time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified     time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
	SellerSKU       string    `json:"seller_sku,omitempty" gorm:"->;column:seller_sku"`
	ProductTitle    string    `json:"product_title,omitempty" gorm:"->;column:product_title"`
	ProductImageURL *string   `json:"product_image_url,omitempty" gorm:"->;column:product_image_url"`
}

func (SalesOrderItem) TableName() string {
	return "sales_order_item"
}

type ReportImportStatus string

const (
	ReportImportStatusPending        ReportImportStatus = "PENDING"
	ReportImportStatusProcessing     ReportImportStatus = "PROCESSING"
	ReportImportStatusSuccess        ReportImportStatus = "SUCCESS"
	ReportImportStatusFailed         ReportImportStatus = "FAILED"
	ReportImportStatusPartialSuccess ReportImportStatus = "PARTIAL_SUCCESS"
)

type ReportImport struct {
	ID          uint64             `json:"id" gorm:"primaryKey;autoIncrement"`
	BatchNo     string             `json:"batch_no" gorm:"column:batch_no;size:64;not null"`
	ReportType  string             `json:"report_type" gorm:"column:report_type;size:50;not null"`
	FileName    string             `json:"file_name" gorm:"column:file_name;size:255;not null"`
	FileHash    string             `json:"file_hash" gorm:"column:file_hash;size:64;not null"`
	Status      ReportImportStatus `json:"status" gorm:"column:status;size:20;not null"`
	TotalRows   uint32             `json:"total_rows" gorm:"column:total_rows;not null"`
	SuccessRows uint32             `json:"success_rows" gorm:"column:success_rows;not null"`
	ErrorRows   uint32             `json:"error_rows" gorm:"column:error_rows;not null"`
	Message     *string            `json:"message" gorm:"column:message;size:500"`
	OperatorID  *uint64            `json:"operator_id" gorm:"column:operator_id"`
	StartedAt   *time.Time         `json:"started_at" gorm:"column:started_at"`
	FinishedAt  *time.Time         `json:"finished_at" gorm:"column:finished_at"`
	GmtCreate   time.Time          `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time          `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ReportImport) TableName() string {
	return "report_import"
}

type ReportImportRowError struct {
	ID             uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ReportImportID uint64    `json:"report_import_id" gorm:"column:report_import_id;not null"`
	RowNo          uint32    `json:"row_no" gorm:"column:row_no;not null"`
	ErrorCode      *string   `json:"error_code" gorm:"column:error_code;size:50"`
	ErrorMessage   string    `json:"error_message" gorm:"column:error_message;size:500;not null"`
	RawRow         *string   `json:"raw_row" gorm:"column:raw_row;type:text"`
	GmtCreate      time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified    time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ReportImportRowError) TableName() string {
	return "report_import_row_error"
}

type ImportOrderLine struct {
	RowNo           uint32
	OrderNo         string
	SourceType      string
	ExternalOrderNo string
	LineNo          uint32
	SellerSKU       string
	Qty             uint64
	Marketplace     string
	OrderDate       time.Time
	SalesChannel    *string
	Currency        string
	UnitPrice       float64
	RawRow          string
}

type AllocateLine struct {
	ItemID       uint64
	QtyAllocated uint64
}

type AllocateParams struct {
	WarehouseID uint64
	Lines       []AllocateLine
}

type ReturnLine struct {
	ItemID      uint64
	QtyReturned uint64
}

type ReturnParams struct {
	WarehouseID uint64
	Lines       []ReturnLine
}

type ShipLine struct {
	ItemID     uint64
	QtyShipped uint64
}

type ShipParams struct {
	WarehouseID uint64
	Lines       []ShipLine
}

type SalesOrderListParams struct {
	Page        int
	PageSize    int
	Status      SalesOrderStatus
	Marketplace string
	Keyword     string
}
