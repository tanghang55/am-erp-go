package repository

import (
	"am-erp-go/internal/module/logistics/domain"

	"gorm.io/gorm"
)

type LogisticsProviderRepository struct {
	db *gorm.DB
}

func NewLogisticsProviderRepository(db *gorm.DB) domain.LogisticsProviderRepository {
	return &LogisticsProviderRepository{db: db}
}

func (r *LogisticsProviderRepository) Create(provider *domain.LogisticsProvider) error {
	return r.db.Create(provider).Error
}

func (r *LogisticsProviderRepository) Update(provider *domain.LogisticsProvider) error {
	return r.db.Save(provider).Error
}

func (r *LogisticsProviderRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.LogisticsProvider{}, id).Error
}

func (r *LogisticsProviderRepository) GetByID(id uint64) (*domain.LogisticsProvider, error) {
	var provider domain.LogisticsProvider
	err := r.db.First(&provider, id).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *LogisticsProviderRepository) GetByCode(code string) (*domain.LogisticsProvider, error) {
	var provider domain.LogisticsProvider
	err := r.db.Where("provider_code = ?", code).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *LogisticsProviderRepository) List(params *domain.LogisticsProviderListParams) ([]*domain.LogisticsProvider, int64, error) {
	var providers []*domain.LogisticsProvider
	var total int64

	query := r.db.Model(&domain.LogisticsProvider{})

	// 筛选条件
	if params.ProviderType != nil {
		query = query.Where("provider_type = ?", *params.ProviderType)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.Keyword != nil && *params.Keyword != "" {
		keyword := "%" + *params.Keyword + "%"
		query = query.Where("provider_code LIKE ? OR provider_name LIKE ?", keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("gmt_create DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&providers).Error; err != nil {
		return nil, 0, err
	}

	return providers, total, nil
}
