package repository

import (
	"am-erp-go/internal/module/supplier/domain"

	"gorm.io/gorm"
)

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) domain.SupplierRepository {
	return &supplierRepository{db: db}
}

func (r *supplierRepository) List(params *domain.SupplierListParams) ([]domain.Supplier, int64, error) {
	var suppliers []domain.Supplier
	var total int64

	query := r.db.Model(&domain.Supplier{})

	// 关键词搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("name LIKE ?", keyword)
	}

	// 状态筛选
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 类型筛选
	if len(params.Types) > 0 {
		query = query.Joins("JOIN supplier_type st ON st.supplier_id = supplier.id").
			Where("st.type IN ?", params.Types).
			Group("supplier.id")
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
		Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

func (r *supplierRepository) GetByID(id uint64) (*domain.Supplier, error) {
	var supplier domain.Supplier
	if err := r.db.First(&supplier, id).Error; err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *supplierRepository) Create(supplier *domain.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *supplierRepository) Update(supplier *domain.Supplier) error {
	return r.db.Save(supplier).Error
}

func (r *supplierRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Supplier{}, id).Error
}
