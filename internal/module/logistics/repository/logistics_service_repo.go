package repository

import (
	"am-erp-go/internal/module/logistics/domain"

	"gorm.io/gorm"
)

type LogisticsServiceRepository struct {
	db *gorm.DB
}

func NewLogisticsServiceRepository(db *gorm.DB) domain.LogisticsServiceRepository {
	return &LogisticsServiceRepository{db: db}
}

func (r *LogisticsServiceRepository) Create(service *domain.LogisticsService) error {
	return r.db.Create(service).Error
}

func (r *LogisticsServiceRepository) Update(service *domain.LogisticsService) error {
	return r.db.Save(service).Error
}

func (r *LogisticsServiceRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.LogisticsService{}, id).Error
}

func (r *LogisticsServiceRepository) GetByID(id uint64) (*domain.LogisticsService, error) {
	var service domain.LogisticsService
	err := r.db.First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *LogisticsServiceRepository) GetByCode(code string) (*domain.LogisticsService, error) {
	var service domain.LogisticsService
	err := r.db.Where("service_code = ?", code).First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *LogisticsServiceRepository) List(params *domain.LogisticsServiceListParams) ([]*domain.LogisticsService, int64, error) {
	var services []*domain.LogisticsService
	var total int64

	query := r.db.Model(&domain.LogisticsService{})

	// 筛选条件
	if params.TransportMode != nil {
		query = query.Where("transport_mode = ?", *params.TransportMode)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.Keyword != nil && *params.Keyword != "" {
		keyword := "%" + *params.Keyword + "%"
		query = query.Where("service_name LIKE ? OR service_code LIKE ? OR destination_region LIKE ?", keyword, keyword, keyword)
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
		Find(&services).Error; err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (r *LogisticsServiceRepository) GetActiveServices() ([]*domain.LogisticsService, error) {
	var services []*domain.LogisticsService
	err := r.db.Where("status = ?", domain.ServiceStatusActive).
		Order("transport_mode, service_name").
		Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (r *LogisticsServiceRepository) GetServicesByTransportMode(transportMode domain.TransportMode) ([]*domain.LogisticsService, error) {
	var services []*domain.LogisticsService
	err := r.db.Where("status = ? AND transport_mode = ?", domain.ServiceStatusActive, transportMode).
		Order("service_name").
		Find(&services).Error
	if err != nil {
		return nil, err
	}
	return services, nil
}
