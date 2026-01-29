package repository

import (
	"am-erp-go/internal/module/shipment/domain"

	"gorm.io/gorm"
)

type packageSpecPackagingRepository struct {
	db *gorm.DB
}

func NewPackageSpecPackagingRepository(db *gorm.DB) domain.PackageSpecPackagingRepository {
	return &packageSpecPackagingRepository{db: db}
}

func (r *packageSpecPackagingRepository) ListByPackageSpecID(packageSpecID uint64) ([]domain.PackageSpecPackagingItem, error) {
	var items []domain.PackageSpecPackagingItem

	// 查询装箱规格包材配置，并关联包材详情
	err := r.db.Table("package_spec_packaging_items as pspi").
		Select(`pspi.id, pspi.package_spec_id, pspi.packaging_item_id, pspi.quantity_per_box, pspi.notes,
			pi.id as 'packaging_item.id',
			pi.item_code as 'packaging_item.item_code',
			pi.item_name as 'packaging_item.item_name',
			pi.specification as 'packaging_item.specification',
			pi.unit as 'packaging_item.unit',
			pi.unit_cost as 'packaging_item.unit_cost',
			pi.currency as 'packaging_item.currency',
			pi.quantity_on_hand as 'packaging_item.quantity_on_hand'`).
		Joins("LEFT JOIN packaging_item pi ON pi.id = pspi.packaging_item_id").
		Where("pspi.package_spec_id = ?", packageSpecID).
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

func (r *packageSpecPackagingRepository) ReplaceAll(packageSpecID uint64, items []domain.PackageSpecPackagingItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除该装箱规格的所有包材配置
		if err := tx.Where("package_spec_id = ?", packageSpecID).Delete(&domain.PackageSpecPackagingItem{}).Error; err != nil {
			return err
		}

		// 2. 如果没有新配置，直接返回
		if len(items) == 0 {
			return nil
		}

		// 3. 插入新配置
		for i := range items {
			items[i].PackageSpecID = packageSpecID
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		return nil
	})
}
