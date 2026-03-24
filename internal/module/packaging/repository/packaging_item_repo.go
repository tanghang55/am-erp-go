package repository

import (
	"am-erp-go/internal/module/packaging/domain"

	"gorm.io/gorm"
)

type packagingItemRepository struct {
	db *gorm.DB
}

func packagingItemColumn(column string) string {
	return "packaging_item." + column
}

func packagingItemReferenceCountSQL() string {
	return `
		SELECT
			COALESCE((SELECT COUNT(1) FROM product_packaging_items WHERE packaging_item_id = ?), 0) +
			COALESCE((SELECT COUNT(1) FROM packaging_ledger WHERE packaging_item_id = ?), 0) +
			COALESCE((SELECT COUNT(1) FROM packaging_procurement_plan WHERE packaging_item_id = ?), 0) +
			COALESCE((SELECT COUNT(1) FROM packaging_purchase_order_item WHERE packaging_item_id = ?), 0) +
			COALESCE((SELECT COUNT(1) FROM package_spec_packaging_items WHERE packaging_item_id = ?), 0) AS total
	`
}

func NewPackagingItemRepository(db *gorm.DB) domain.PackagingItemRepository {
	return &packagingItemRepository{db: db}
}

func (r *packagingItemRepository) List(params *domain.PackagingItemListParams) ([]domain.PackagingItem, int64, error) {
	var items []domain.PackagingItem
	var total int64

	query := r.db.Table("packaging_item").
		Select("packaging_item.*, COALESCE(supplier.name, packaging_item.supplier_name) AS supplier_name").
		Joins("LEFT JOIN supplier ON supplier.id = packaging_item.supplier_id")

	// 关键词搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(packagingItemColumn("item_code")+" LIKE ? OR "+packagingItemColumn("item_name")+" LIKE ?", keyword, keyword)
	}

	// 类别筛选
	if params.Category != "" {
		query = query.Where(packagingItemColumn("category")+" = ?", params.Category)
	}

	// 状态筛选
	if params.Status != "" {
		query = query.Where(packagingItemColumn("status")+" = ?", params.Status)
	}

	// 低库存筛选
	if params.LowStock {
		query = query.Where("quantity_on_hand <= reorder_point AND reorder_point IS NOT NULL")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order(packagingItemColumn("gmt_create") + " DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *packagingItemRepository) GetByID(id uint64) (*domain.PackagingItem, error) {
	var item domain.PackagingItem
	if err := r.db.Table("packaging_item").
		Select("packaging_item.*, COALESCE(supplier.name, packaging_item.supplier_name) AS supplier_name").
		Joins("LEFT JOIN supplier ON supplier.id = packaging_item.supplier_id").
		Where("packaging_item.id = ?", id).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *packagingItemRepository) Create(item *domain.PackagingItem) error {
	return r.db.Create(item).Error
}

func (r *packagingItemRepository) Update(item *domain.PackagingItem) error {
	return r.db.Save(item).Error
}

func (r *packagingItemRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.PackagingItem{}, id).Error
}

func (r *packagingItemRepository) CountReferences(id uint64) (int64, error) {
	type referenceCount struct {
		Total int64 `gorm:"column:total"`
	}

	var result referenceCount
	err := r.db.Raw(packagingItemReferenceCountSQL(), id, id, id, id, id).Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Total, nil
}

func (r *packagingItemRepository) GetLowStockItems() ([]domain.PackagingItem, error) {
	var items []domain.PackagingItem
	err := r.db.Where(packagingItemColumn("quantity_on_hand") + " <= " + packagingItemColumn("reorder_point") + " AND " + packagingItemColumn("reorder_point") + " IS NOT NULL AND " + packagingItemColumn("status") + " = 'ACTIVE'").
		Find(&items).Error
	return items, err
}

func (r *packagingItemRepository) UpdateQuantity(id uint64, quantity int64) error {
	return r.db.Model(&domain.PackagingItem{}).
		Where("id = ?", id).
		UpdateColumn("quantity_on_hand", gorm.Expr("quantity_on_hand + ?", quantity)).
		Error
}
