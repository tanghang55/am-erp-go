package repository

import (
	"am-erp-go/internal/module/shipment/domain"

	"gorm.io/gorm"
)

type packageSpecRepository struct {
	db *gorm.DB
}

func NewPackageSpecRepository(db *gorm.DB) domain.PackageSpecRepository {
	return &packageSpecRepository{db: db}
}

func (r *packageSpecRepository) Create(spec *domain.PackageSpec) error {
	return r.db.Create(spec).Error
}

func (r *packageSpecRepository) Update(spec *domain.PackageSpec) error {
	return r.db.Save(spec).Error
}

func (r *packageSpecRepository) GetByID(id uint64) (*domain.PackageSpec, error) {
	var spec domain.PackageSpec
	if err := r.db.First(&spec, id).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}

func (r *packageSpecRepository) List(params *domain.PackageSpecListParams) ([]*domain.PackageSpec, int64, error) {
	var specs []*domain.PackageSpec
	var total int64

	query := r.db.Model(&domain.PackageSpec{})

	if params.Keyword != nil && *params.Keyword != "" {
		keyword := "%" + *params.Keyword + "%"
		query = query.Where("name LIKE ?", keyword)
	}

	if params.Status != nil && *params.Status != "" {
		query = query.Where("status = ?", *params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("gmt_modified DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&specs).Error; err != nil {
		return nil, 0, err
	}

	return specs, total, nil
}

func (r *packageSpecRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.PackageSpec{}, id).Error
}

func (r *packageSpecRepository) ListByIDs(ids []uint64) ([]*domain.PackageSpec, error) {
	if len(ids) == 0 {
		return []*domain.PackageSpec{}, nil
	}

	var specs []*domain.PackageSpec
	if err := r.db.Where("id IN ?", ids).Find(&specs).Error; err != nil {
		return nil, err
	}
	return specs, nil
}
