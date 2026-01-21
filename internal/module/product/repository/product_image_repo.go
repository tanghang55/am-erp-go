package repository

import (
	"am-erp-go/internal/module/product/domain"

	"gorm.io/gorm"
)

type productImageRepository struct {
	db *gorm.DB
}

func NewProductImageRepository(db *gorm.DB) domain.ProductImageRepository {
	return &productImageRepository{db: db}
}

func (r *productImageRepository) ListByProductID(productID uint64) ([]domain.ProductImage, error) {
	var items []domain.ProductImage
	if err := r.db.Where("product_id = ?", productID).
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *productImageRepository) ReplaceAll(productID uint64, orderedUrls []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("product_id = ?", productID).Delete(&domain.ProductImage{}).Error; err != nil {
			return err
		}
		images := make([]domain.ProductImage, 0, len(orderedUrls))
		for i, url := range orderedUrls {
			isPrimary := uint8(0)
			if i == 0 {
				isPrimary = 1
			}
			images = append(images, domain.ProductImage{
				ProductID: productID,
				ImageUrl:  url,
				SortOrder: uint32(i + 1),
				IsPrimary: isPrimary,
			})
		}
		if len(images) == 0 {
			return nil
		}
		return tx.Create(&images).Error
	})
}
