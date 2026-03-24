package repository

import (
	"am-erp-go/internal/module/product/domain"

	"gorm.io/gorm"
)

type productCategoryRepository struct {
	db *gorm.DB
}

func NewProductCategoryRepository(db *gorm.DB) domain.ProductCategoryRepository {
	return &productCategoryRepository{db: db}
}

func (r *productCategoryRepository) ListAll() ([]domain.ProductCategory, error) {
	var items []domain.ProductCategory
	if err := r.db.Order("level ASC, sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *productCategoryRepository) GetByID(id uint64) (*domain.ProductCategory, error) {
	var item domain.ProductCategory
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *productCategoryRepository) Create(item *domain.ProductCategory) error {
	return r.db.Create(item).Error
}

func (r *productCategoryRepository) Update(item *domain.ProductCategory) error {
	return r.db.Save(item).Error
}

func (r *productCategoryRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.ProductCategory{}, id).Error
}

func (r *productCategoryRepository) CountChildren(id uint64) (int64, error) {
	var total int64
	if err := r.db.Model(&domain.ProductCategory{}).Where("parent_id = ?", id).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
