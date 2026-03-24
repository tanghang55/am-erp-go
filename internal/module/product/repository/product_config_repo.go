package repository

import (
	"am-erp-go/internal/module/product/domain"

	"gorm.io/gorm"
)

type productConfigRepository struct {
	db *gorm.DB
}

func NewProductConfigRepository(db *gorm.DB) domain.ProductConfigRepository {
	return &productConfigRepository{db: db}
}

func (r *productConfigRepository) List(params *domain.ProductConfigListParams) ([]domain.ProductConfigItem, int64, error) {
	var items []domain.ProductConfigItem
	var total int64

	query := r.db.Model(&domain.ProductConfigItem{})

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("item_code LIKE ? OR item_name LIKE ?", keyword, keyword)
	}
	if params.ConfigType != "" {
		query = query.Where("config_type = ?", params.ConfigType)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Order("sort ASC, id ASC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *productConfigRepository) GetByID(id uint64) (*domain.ProductConfigItem, error) {
	var item domain.ProductConfigItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *productConfigRepository) Create(item *domain.ProductConfigItem) error {
	return r.db.Create(item).Error
}

func (r *productConfigRepository) Update(item *domain.ProductConfigItem) error {
	return r.db.Save(item).Error
}

func (r *productConfigRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.ProductConfigItem{}, id).Error
}
