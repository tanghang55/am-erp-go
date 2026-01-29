package repository

import (
	"am-erp-go/internal/module/product/domain"

	"gorm.io/gorm"
)

type productPackagingRepository struct {
	db *gorm.DB
}

func NewProductPackagingRepository(db *gorm.DB) domain.ProductPackagingRepository {
	return &productPackagingRepository{db: db}
}

func (r *productPackagingRepository) ListByProductID(productID uint64) ([]domain.ProductPackagingItem, error) {
	var items []domain.ProductPackagingItem

	// 查询产品包材配置，并关联包材详情
	err := r.db.Table("product_packaging_items as ppi").
		Select(`ppi.id, ppi.product_id, ppi.packaging_item_id, ppi.quantity_per_unit, ppi.notes,
			pi.id as 'packaging_item.id',
			pi.item_code as 'packaging_item.item_code',
			pi.item_name as 'packaging_item.item_name',
			pi.specification as 'packaging_item.specification',
			pi.unit as 'packaging_item.unit',
			pi.unit_cost as 'packaging_item.unit_cost',
			pi.currency as 'packaging_item.currency',
			pi.quantity_on_hand as 'packaging_item.quantity_on_hand'`).
		Joins("LEFT JOIN packaging_item pi ON pi.id = ppi.packaging_item_id").
		Where("ppi.product_id = ?", productID).
		Scan(&items).Error

	if err != nil {
		return nil, err
	}

	// 手动构造PackagingItem对象
	for i := range items {
		if items[i].PackagingItemID > 0 {
			items[i].PackagingItem = &domain.PackagingItemDetail{}
		}
	}

	return items, nil
}

func (r *productPackagingRepository) ReplaceAll(productID uint64, items []domain.ProductPackagingItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除该产品的所有包材配置
		if err := tx.Where("product_id = ?", productID).Delete(&domain.ProductPackagingItem{}).Error; err != nil {
			return err
		}

		// 2. 如果没有新配置，直接返回
		if len(items) == 0 {
			return nil
		}

		// 3. 插入新配置
		for i := range items {
			items[i].ProductID = productID
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil
	})
}
