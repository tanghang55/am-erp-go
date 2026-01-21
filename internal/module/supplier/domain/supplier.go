package domain

import "time"

// Supplier represents a supplier profile.
type Supplier struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierCode string    `json:"supplier_code" gorm:"column:supplier_code;size:50;not null"`
	Name         string    `json:"name" gorm:"size:200;not null"`
	Status       string    `json:"status" gorm:"size:20;default:'ACTIVE'"`
	Remark       string    `json:"remark" gorm:"type:text"`
	GmtCreate    time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified  time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (Supplier) TableName() string {
	return "supplier"
}

// SupplierType represents a supplier type relation.
type SupplierType struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierID  uint64    `json:"supplier_id" gorm:"column:supplier_id;not null"`
	Type        string    `json:"type" gorm:"size:20;not null"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (SupplierType) TableName() string {
	return "supplier_type"
}

// SupplierContact represents a supplier contact.
type SupplierContact struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierID  uint64    `json:"supplier_id" gorm:"column:supplier_id;not null"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Phone       string    `json:"phone" gorm:"size:50"`
	Email       string    `json:"email" gorm:"size:100"`
	Position    string    `json:"position" gorm:"size:100"`
	IsPrimary   uint8     `json:"is_primary" gorm:"column:is_primary"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (SupplierContact) TableName() string {
	return "supplier_contact"
}

// SupplierAccount represents a supplier account.
type SupplierAccount struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierID   uint64    `json:"supplier_id" gorm:"column:supplier_id;not null"`
	BankName     string    `json:"bank_name" gorm:"column:bank_name;size:100;not null"`
	BankAccount  string    `json:"bank_account" gorm:"column:bank_account;size:100;not null"`
	Currency     string    `json:"currency" gorm:"size:20"`
	TaxNo        string    `json:"tax_no" gorm:"column:tax_no;size:100"`
	PaymentTerms string    `json:"payment_terms" gorm:"column:payment_terms;size:100"`
	GmtCreate    time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified  time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (SupplierAccount) TableName() string {
	return "supplier_account"
}

// SupplierTag represents a supplier tag.
type SupplierTag struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierID  uint64    `json:"supplier_id" gorm:"column:supplier_id;not null"`
	Tag         string    `json:"tag" gorm:"size:100;not null"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (SupplierTag) TableName() string {
	return "supplier_tag"
}

// SupplierDetail is the aggregated supplier detail response.
type SupplierDetail struct {
	Supplier
	Types    []string          `json:"types"`
	Contacts []SupplierContact `json:"contacts"`
	Accounts []SupplierAccount `json:"accounts"`
	Tags     []SupplierTag     `json:"tags"`
}

// SupplierListItem is the supplier list response with types.
type SupplierListItem struct {
	Supplier
	Types []string `json:"types"`
}
