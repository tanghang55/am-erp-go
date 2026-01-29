package repository

import (
	"time"

	"am-erp-go/internal/module/logistics/domain"

	"gorm.io/gorm"
)

type ShippingRateRepository struct {
	db *gorm.DB
}

func NewShippingRateRepository(db *gorm.DB) domain.ShippingRateRepository {
	return &ShippingRateRepository{db: db}
}

func (r *ShippingRateRepository) Create(rate *domain.ShippingRate) error {
	return r.db.Create(rate).Error
}

func (r *ShippingRateRepository) Update(rate *domain.ShippingRate) error {
	return r.db.Save(rate).Error
}

func (r *ShippingRateRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.ShippingRate{}, id).Error
}

func (r *ShippingRateRepository) GetByID(id uint64) (*domain.ShippingRate, error) {
	var rate domain.ShippingRate
	err := r.db.First(&rate, id).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *ShippingRateRepository) List(params *domain.ShippingRateListParams) ([]*domain.ShippingRate, int64, error) {
	var rates []*domain.ShippingRate
	var total int64

	query := r.db.Model(&domain.ShippingRate{})

	// 筛选条件
	if params.ProviderID != nil {
		query = query.Where("provider_id = ?", *params.ProviderID)
	}
	if params.OriginWarehouseID != nil {
		query = query.Where("origin_warehouse_id = ?", *params.OriginWarehouseID)
	}
	if params.DestinationWarehouseID != nil {
		query = query.Where("destination_warehouse_id = ?", *params.DestinationWarehouseID)
	}
	if params.TransportMode != nil {
		query = query.Where("transport_mode = ?", *params.TransportMode)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.Keyword != nil && *params.Keyword != "" {
		keyword := "%" + *params.Keyword + "%"
		query = query.Where("service_name LIKE ?", keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Preload("Provider").
		Preload("OriginWarehouse").
		Preload("DestinationWarehouse").
		Preload("Service").
		Order("effective_date DESC, gmt_create DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&rates).Error; err != nil {
		return nil, 0, err
	}

	return rates, total, nil
}

// QueryLatestRate 查询最新的有效报价
func (r *ShippingRateRepository) QueryLatestRate(params *domain.QueryLatestRateParams) (*domain.ShippingRate, error) {
	var rate domain.ShippingRate

	// 默认查询日期为今天
	queryDate := time.Now().Format("2006-01-02")
	if params.QueryDate != nil {
		queryDate = *params.QueryDate
	}

	query := r.db.Model(&domain.ShippingRate{}).
		Where("status = ?", domain.RateStatusActive).
		Where("origin_warehouse_id = ?", params.OriginWarehouseID).
		Where("destination_warehouse_id = ?", params.DestinationWarehouseID).
		Where("transport_mode = ?", params.TransportMode).
		Where("effective_date <= ?", queryDate).
		Where("(expiry_date IS NULL OR expiry_date >= ?)", queryDate)

	if params.ProviderID != nil {
		query = query.Where("provider_id = ?", *params.ProviderID)
	}

	// 如果指定了重量，匹配重量区间
	if params.Weight != nil {
		query = query.Where("(min_weight IS NULL OR min_weight <= ?)", *params.Weight).
			Where("(max_weight IS NULL OR max_weight >= ?)", *params.Weight)
	}

	// 如果指定了体积，匹配体积区间
	if params.Volume != nil {
		query = query.Where("(min_volume IS NULL OR min_volume <= ?)", *params.Volume).
			Where("(max_volume IS NULL OR max_volume >= ?)", *params.Volume)
	}

	// 按生效日期降序，取最新的
	err := query.
		Preload("Provider").
		Preload("OriginWarehouse").
		Preload("DestinationWarehouse").
		Preload("Service").
		Order("effective_date DESC").
		First(&rate).Error

	if err != nil {
		return nil, err
	}

	return &rate, nil
}
