package domain

import "time"

type SKUMappingStatus string

const (
	SKUMappingStatusActive   SKUMappingStatus = "ACTIVE"
	SKUMappingStatusDisabled SKUMappingStatus = "DISABLED"
)

type IntegrationSKUMapping struct {
	ID               uint64           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProviderCode     string           `json:"provider_code" gorm:"column:provider_code;size:64;not null"`
	Marketplace      string           `json:"marketplace" gorm:"column:marketplace;size:10;not null"`
	SellerSKU        string           `json:"seller_sku" gorm:"column:seller_sku;size:100;not null"`
	ProductID        uint64           `json:"product_id" gorm:"column:product_id;not null"`
	Status           SKUMappingStatus `json:"status" gorm:"column:status;type:enum('ACTIVE','DISABLED');not null;default:'ACTIVE'"`
	Remark           *string          `json:"remark" gorm:"column:remark;size:255"`
	CreatedBy        *uint64          `json:"created_by" gorm:"column:created_by"`
	UpdatedBy        *uint64          `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate        time.Time        `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified      time.Time        `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
	ProductTitle     string           `json:"product_title,omitempty" gorm:"->;column:product_title"`
	ProductSellerSKU string           `json:"product_seller_sku,omitempty" gorm:"->;column:product_seller_sku"`
}

func (IntegrationSKUMapping) TableName() string {
	return "integration_sku_mapping"
}

type SKUMappingListParams struct {
	Page         int
	PageSize     int
	ProviderCode string
	Marketplace  string
	Status       string
	Keyword      string
	ProductID    *uint64
}

type SKUMappingRepository interface {
	Create(item *IntegrationSKUMapping) error
	Update(item *IntegrationSKUMapping) error
	GetByID(id uint64) (*IntegrationSKUMapping, error)
	GetByUnique(providerCode string, marketplace string, sellerSKU string) (*IntegrationSKUMapping, error)
	List(params *SKUMappingListParams) ([]IntegrationSKUMapping, int64, error)
	ResolveActiveProductID(providerCode string, marketplace string, sellerSKU string) (uint64, error)
}
