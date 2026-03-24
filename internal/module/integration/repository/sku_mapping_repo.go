package repository

import (
	"errors"
	"strings"

	integrationDomain "am-erp-go/internal/module/integration/domain"

	"gorm.io/gorm"
)

type skuMappingRepository struct {
	db *gorm.DB
}

func NewSKUMappingRepository(db *gorm.DB) *skuMappingRepository {
	return &skuMappingRepository{db: db}
}

func (r *skuMappingRepository) Create(item *integrationDomain.IntegrationSKUMapping) error {
	if item == nil {
		return nil
	}
	return r.db.Create(item).Error
}

func (r *skuMappingRepository) Update(item *integrationDomain.IntegrationSKUMapping) error {
	if item == nil {
		return nil
	}
	return r.db.Save(item).Error
}

func (r *skuMappingRepository) GetByID(id uint64) (*integrationDomain.IntegrationSKUMapping, error) {
	var item integrationDomain.IntegrationSKUMapping
	err := r.baseQuery().Where("m.id = ?", id).Take(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *skuMappingRepository) GetByUnique(providerCode string, marketplace string, sellerSKU string) (*integrationDomain.IntegrationSKUMapping, error) {
	var item integrationDomain.IntegrationSKUMapping
	err := r.db.Model(&integrationDomain.IntegrationSKUMapping{}).
		Where("provider_code = ?", providerCode).
		Where("marketplace = ?", marketplace).
		Where("seller_sku = ?", sellerSKU).
		Take(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *skuMappingRepository) List(params *integrationDomain.SKUMappingListParams) ([]integrationDomain.IntegrationSKUMapping, int64, error) {
	var (
		list  []integrationDomain.IntegrationSKUMapping
		total int64
	)
	if params == nil {
		params = &integrationDomain.SKUMappingListParams{Page: 1, PageSize: 20}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	query := r.baseQuery()
	query = applyListFilter(query, params)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []integrationDomain.IntegrationSKUMapping{}, 0, nil
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("m.id DESC").Offset(offset).Limit(params.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	if len(list) == 0 {
		return []integrationDomain.IntegrationSKUMapping{}, total, nil
	}
	return list, total, nil
}

func (r *skuMappingRepository) ResolveActiveProductID(providerCode string, marketplace string, sellerSKU string) (uint64, error) {
	var row struct {
		ProductID uint64 `gorm:"column:product_id"`
	}
	err := r.db.Model(&integrationDomain.IntegrationSKUMapping{}).
		Select("product_id").
		Where("provider_code = ?", providerCode).
		Where("marketplace = ?", marketplace).
		Where("seller_sku = ?", sellerSKU).
		Where("status = ?", integrationDomain.SKUMappingStatusActive).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return row.ProductID, nil
}

func (r *skuMappingRepository) baseQuery() *gorm.DB {
	return r.db.Table("integration_sku_mapping m").
		Select(`
			m.id,
			m.provider_code,
			m.marketplace,
			m.seller_sku,
			m.product_id,
			m.status,
			m.remark,
			m.created_by,
			m.updated_by,
			m.gmt_create,
			m.gmt_modified,
			COALESCE(p.title, '') AS product_title,
			COALESCE(p.seller_sku, '') AS product_seller_sku
		`).
		Joins("LEFT JOIN product p ON p.id = m.product_id")
}

func applyListFilter(query *gorm.DB, params *integrationDomain.SKUMappingListParams) *gorm.DB {
	if params == nil {
		return query
	}
	if params.ProviderCode != "" {
		query = query.Where("m.provider_code = ?", strings.TrimSpace(params.ProviderCode))
	}
	if params.Marketplace != "" {
		query = query.Where("m.marketplace = ?", strings.TrimSpace(params.Marketplace))
	}
	if params.Status != "" {
		query = query.Where("m.status = ?", strings.TrimSpace(params.Status))
	}
	if params.ProductID != nil && *params.ProductID > 0 {
		query = query.Where("m.product_id = ?", *params.ProductID)
	}
	if strings.TrimSpace(params.Keyword) != "" {
		keyword := "%" + strings.TrimSpace(params.Keyword) + "%"
		query = query.Where("m.seller_sku LIKE ? OR p.seller_sku LIKE ? OR p.title LIKE ?", keyword, keyword, keyword)
	}
	return query
}
