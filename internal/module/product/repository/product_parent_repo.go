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

func (r *productParentRepository) withChildStats() *gorm.DB {
	childStats := r.db.Table("product").
		Select(`
			parent_id,
			COUNT(*) AS child_count,
			SUM(CASE WHEN status IN ('ON_SALE', 'REPLENISHING') THEN 1 ELSE 0 END) AS active_child_count,
			SUM(CASE WHEN status NOT IN ('ON_SALE', 'REPLENISHING') THEN 1 ELSE 0 END) AS inactive_child_count
		`).
		Where("parent_id IS NOT NULL").
		Group("parent_id")

	return r.db.Model(&domain.ProductParent{}).
		Joins("LEFT JOIN (?) child_stats ON child_stats.parent_id = product_parent.id", childStats).
		Select(`
			product_parent.*,
			COALESCE(child_stats.child_count, 0) AS child_count,
			COALESCE(child_stats.active_child_count, 0) AS active_child_count,
			COALESCE(child_stats.inactive_child_count, 0) AS inactive_child_count
		`)
}

func (r *productParentRepository) List(params *domain.ProductParentListParams) ([]domain.ProductParent, int64, error) {
	var parents []domain.ProductParent
	var total int64

	query := r.withChildStats()

	// 关键词搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("product_parent.parent_asin LIKE ? OR product_parent.title LIKE ? OR product_parent.brand LIKE ?", keyword, keyword, keyword)
	}

	// 站点筛选
	if params.Marketplace != "" {
		query = query.Where("product_parent.marketplace = ?", params.Marketplace)
	}

	// 状态筛选
	if params.Status != "" {
		query = query.Where("product_parent.status = ?", params.Status)
	}

	if params.HasChildren == "true" {
		query = query.Where("COALESCE(child_stats.child_count, 0) > 0")
	} else if params.HasChildren == "false" {
		query = query.Where("COALESCE(child_stats.child_count, 0) = 0")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("product_parent.gmt_modified DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&parents).Error; err != nil {
		return nil, 0, err
	}

	return parents, total, nil
}

func (r *productParentRepository) GetByID(id uint64) (*domain.ProductParent, error) {
	var parent domain.ProductParent
	if err := r.withChildStats().Where("product_parent.id = ?", id).First(&parent).Error; err != nil {
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
