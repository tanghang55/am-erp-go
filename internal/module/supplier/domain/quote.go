package domain

import "time"

type ProductSupplierQuote struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID    uint64    `json:"product_id" gorm:"column:product_id;not null"`
	SupplierID   uint64    `json:"supplier_id" gorm:"column:supplier_id;not null"`
	Price        float64   `json:"price" gorm:"column:price;type:decimal(15,4);not null"`
	Currency     string    `json:"currency" gorm:"column:currency;size:3;not null"`
	QtyMOQ       uint64    `json:"qty_moq" gorm:"column:qty_moq"`
	LeadTimeDays uint64    `json:"lead_time_days" gorm:"column:lead_time_days"`
	Status       string    `json:"status" gorm:"size:20"`
	Remark       string    `json:"remark" gorm:"type:text"`
	GmtCreate    time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified  time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (ProductSupplierQuote) TableName() string {
	return "product_supplier_quote"
}

type ProductSupplierQuoteDetail struct {
	ProductSupplierQuote
	SupplierName string `json:"supplier_name" gorm:"column:supplier_name"`
	SupplierCode string `json:"supplier_code" gorm:"column:supplier_code"`
}

type ProductQuoteRow struct {
	ProductID         uint64                       `json:"product_id"`
	SellerSku         string                       `json:"seller_sku"`
	Asin              string                       `json:"asin"`
	Marketplace       string                       `json:"marketplace"`
	Title             string                       `json:"title"`
	ImageUrl          string                       `json:"image_url"`
	DefaultSupplierID uint64                       `json:"default_supplier_id"`
	Quotes            []ProductSupplierQuoteDetail `json:"quotes"`
}

type QuoteListParams struct {
	Page        int
	PageSize    int
	Keyword     string
	Marketplace string
	SupplierID  *uint64
}
