package repository

import (
	"am-erp-go/internal/module/product/domain"

	"gorm.io/gorm"
)

type productParentRepository struct {
	db *gorm.DB
}

func NewProductParentRepository(db *gorm.DB) domain.ProductParentRepository {
	return &productParentRepository{db: db}
}

func (r *productParentRepository) List(params *domain.ProductParentListParams) ([]domain.ProductParent, int64, error) {
	var parents []domain.ProductParent
	var total int64

	query := r.db.Model(&domain.ProductParent{})

	// 关键词搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("parent_asin LIKE ? OR title LIKE ? OR brand LIKE ?", keyword, keyword, keyword)
	}

	// 站点筛选
	if params.Marketplace != "" {
		query = query.Where("marketplace = ?", params.Marketplace)
	}

	// 状态筛选
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("gmt_modified DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&parents).Error; err != nil {
		return nil, 0, err
	}

	return parents, total, nil
}

func (r *productParentRepository) GetByID(id uint64) (*domain.ProductParent, error) {
	var parent domain.ProductParent
	if err := r.db.First(&parent, id).Error; err != nil {
		return nil, err
	}
	return &parent, nil
}

func (r *productParentRepository) Create(parent *domain.ProductParent) error {
	return r.db.Create(parent).Error
}

func (r *productParentRepository) Update(parent *domain.ProductParent) error {
	return r.db.Save(parent).Error
}

func (r *productParentRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.ProductParent{}, id).Error
}
