package repository

import (
	"am-erp-go/internal/module/product/domain"

	"gorm.io/gorm"
)

type productPackagingRepository struct {
	db *gorm.DB
}

type productPackagingRow struct {
	ID              uint64  `gorm:"column:id"`
	ProductID       uint64  `gorm:"column:product_id"`
	PackagingItemID uint64  `gorm:"column:packaging_item_id"`
	QuantityPerUnit float64 `gorm:"column:quantity_per_unit"`

	PackagingItemDetail bool `gorm:"-"`
	PackagingItemDetailData
}

type PackagingItemDetailData struct {
	ID             uint64  `gorm:"column:packaging_item_id_value"`
	ItemCode       string  `gorm:"column:packaging_item_item_code"`
	ItemName       string  `gorm:"column:packaging_item_item_name"`
	Specification  *string `gorm:"column:packaging_item_specification"`
	Unit           string  `gorm:"column:packaging_item_unit"`
	UnitCost       float64 `gorm:"column:packaging_item_unit_cost"`
	Currency       string  `gorm:"column:packaging_item_currency"`
	QuantityOnHand uint64  `gorm:"column:packaging_item_quantity_on_hand"`
}

func NewProductPackagingRepository(db *gorm.DB) domain.ProductPackagingRepository {
	return &productPackagingRepository{db: db}
}

func (r *productPackagingRepository) ListByProductID(productID uint64) ([]domain.ProductPackagingItem, error) {
	var rows []productPackagingRow

	// 查询产品包材配置，并关联包材详情
	err := r.db.Table("product_packaging_items as ppi").
		Select(`ppi.id, ppi.product_id, ppi.packaging_item_id, ppi.quantity_per_unit,
			pi.id as packaging_item_id_value,
			pi.item_code as packaging_item_item_code,
			pi.item_name as packaging_item_item_name,
			pi.specification as packaging_item_specification,
			pi.unit as packaging_item_unit,
			pi.unit_cost as packaging_item_unit_cost,
			pi.currency as packaging_item_currency,
			pi.quantity_on_hand as packaging_item_quantity_on_hand`).
		Joins("LEFT JOIN packaging_item pi ON pi.id = ppi.packaging_item_id").
		Where("ppi.product_id = ?", productID).
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return mapProductPackagingRows(rows), nil
}

func mapProductPackagingRows(rows []productPackagingRow) []domain.ProductPackagingItem {
	items := make([]domain.ProductPackagingItem, 0, len(rows))
	for _, row := range rows {
		item := domain.ProductPackagingItem{
			ID:              row.ID,
			ProductID:       row.ProductID,
			PackagingItemID: row.PackagingItemID,
			QuantityPerUnit: row.QuantityPerUnit,
		}
		if row.PackagingItemDetailData.ID > 0 {
			item.PackagingItem = &domain.PackagingItemDetail{
				ID:             row.PackagingItemDetailData.ID,
				ItemCode:       row.PackagingItemDetailData.ItemCode,
				ItemName:       row.PackagingItemDetailData.ItemName,
				Specification:  valueOrEmpty(row.PackagingItemDetailData.Specification),
				Unit:           row.PackagingItemDetailData.Unit,
				UnitCost:       row.PackagingItemDetailData.UnitCost,
				Currency:       row.PackagingItemDetailData.Currency,
				QuantityOnHand: row.PackagingItemDetailData.QuantityOnHand,
			}
		}
		items = append(items, item)
	}
	return items
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
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
