package repository

import (
	"time"

	"am-erp-go/internal/module/product/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) applyProductDisplayJoins(query *gorm.DB) *gorm.DB {
	return query.
		Joins("LEFT JOIN supplier ON supplier.id = product.supplier_id").
		Joins("LEFT JOIN product_supplier_quote AS default_quote ON default_quote.product_id = product.id AND default_quote.supplier_id = product.supplier_id").
		Joins("LEFT JOIN product_config_item AS brand_cfg ON brand_cfg.id = product.brand_id").
		Joins("LEFT JOIN product_category AS category_cfg ON category_cfg.id = product.category_id").
		Joins("LEFT JOIN product_config_item AS dimension_unit_cfg ON dimension_unit_cfg.id = product.dimension_unit_id").
		Joins("LEFT JOIN product_config_item AS weight_unit_cfg ON weight_unit_cfg.id = product.weight_unit_id")
}

func productDisplaySelect() string {
	return "product.*, " +
		"COALESCE(default_quote.price, product.unit_cost) AS unit_cost, " +
		"supplier.name AS supplier_name, supplier.supplier_code AS supplier_code, " +
		"brand_cfg.item_name AS brand_name, category_cfg.category_name AS category_name, " +
		"dimension_unit_cfg.item_name AS dimension_unit_name, weight_unit_cfg.item_name AS weight_unit_name, " +
		"(SELECT audit_log.username FROM audit_log " +
		"WHERE audit_log.entity_type = 'Product' " +
		"AND audit_log.entity_id = CAST(product.id AS CHAR CHARACTER SET utf8mb4) COLLATE utf8mb4_unicode_ci " +
		"AND audit_log.action IN ('UPDATE_PRODUCT', 'CREATE_PRODUCT') " +
		"ORDER BY audit_log.gmt_modified DESC, audit_log.id DESC LIMIT 1) AS updated_by_name"
}

func productListByIDsSelect() string {
	return "product.*, " +
		"COALESCE(default_quote.price, product.unit_cost) AS unit_cost, " +
		"supplier.name AS supplier_name, supplier.supplier_code AS supplier_code, " +
		"brand_cfg.item_name AS brand_name, category_cfg.category_name AS category_name, " +
		"dimension_unit_cfg.item_name AS dimension_unit_name, weight_unit_cfg.item_name AS weight_unit_name"
}

func productListOrderBy() string {
	return "product.gmt_create DESC, product.id DESC"
}

func (r *productRepository) List(params *domain.ProductListParams) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.Model(&domain.Product{})

	invQuery := r.db.Table("inventory_balance").
		Select("product_id, SUM(available_quantity) AS inventory_available, SUM(reserved_quantity) AS inventory_reserved, SUM(purchasing_in_transit + pending_inspection + pending_shipment + logistics_in_transit) AS inventory_inbound")
	if params.WarehouseID != nil {
		invQuery = invQuery.Where("warehouse_id = ?", *params.WarehouseID)
	}
	invQuery = invQuery.Group("product_id")

	// 仓库库存筛选（只返回有待出库存的产品，用于发货单选择）
	if params.WarehouseID != nil {
		query = query.Joins("INNER JOIN inventory_balance ON inventory_balance.product_id = product.id AND inventory_balance.warehouse_id = ? AND inventory_balance.pending_shipment > 0", *params.WarehouseID)
	}

	// 关键词搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("product.seller_sku LIKE ? OR product.asin LIKE ? OR product.title LIKE ?", keyword, keyword, keyword)
	}

	// 站点筛选
	if params.Marketplace != "" {
		query = query.Where("product.marketplace = ?", params.Marketplace)
	}

	// 状态筛选
	if len(params.Statuses) > 0 {
		query = query.Where("product.status IN ?", params.Statuses)
	} else if params.Status != "" {
		query = query.Where("product.status = ?", params.Status)
	}

	// 供应商筛选
	if params.SupplierID != nil {
		query = query.Where("product.supplier_id = ?", *params.SupplierID)
	}

	if params.BrandID != nil {
		query = query.Where("product.brand_id = ?", *params.BrandID)
	}

	if params.CategoryID != nil {
		query = query.Where("product.category_id = ?", *params.CategoryID)
	}

	if params.PackingRequired != nil {
		query = query.Where("product.is_packing_required = ?", *params.PackingRequired)
	}

	if params.ParentID != nil {
		query = query.Where("product.parent_id = ?", *params.ParentID)
	}

	if params.OnlyParentless {
		query = query.Where("product.parent_id IS NULL")
	}

	// 组合筛选
	if params.ComboID != nil {
		query = query.Where("product.combo_id = ?", *params.ComboID)
	}

	// 是否主产品
	if params.IsComboMain != nil {
		query = query.Where("product.is_combo_main = ?", *params.IsComboMain)
	}

	if params.OnlyStandalone {
		query = query.Where("product.combo_id IS NULL AND product.is_combo_main = 0")
	}

	if params.OnlyWithPackaging {
		query = query.Where("EXISTS (SELECT 1 FROM product_packaging_items WHERE product_packaging_items.product_id = product.id)")
	}

	// 排除组合产品的子产品（用于发货单等场景）
	if params.ExcludeComboChild {
		query = query.Where("product.combo_id IS NOT NULL OR product.is_combo_main =1 ")
	}

	query = r.applyProductDisplayJoins(query).
		Joins("LEFT JOIN (?) AS inv ON inv.product_id = product.id", invQuery).
		Select(productDisplaySelect() + ", COALESCE(inv.inventory_available, 0) AS inventory_available, COALESCE(inv.inventory_reserved, 0) AS inventory_reserved, COALESCE(inv.inventory_inbound, 0) AS inventory_inbound")

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order(productListOrderBy()).
		Offset(offset).
		Limit(params.PageSize).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) GetByID(id uint64) (*domain.Product, error) {
	var product domain.Product

	invQuery := r.db.Table("inventory_balance").
		Select("product_id, SUM(available_quantity) AS inventory_available, SUM(reserved_quantity) AS inventory_reserved, SUM(purchasing_in_transit + pending_inspection + pending_shipment + logistics_in_transit) AS inventory_inbound").
		Group("product_id")

	if err := r.applyProductDisplayJoins(r.db.Table("product")).
		Select(productDisplaySelect()+", COALESCE(inv.inventory_available, 0) AS inventory_available, COALESCE(inv.inventory_reserved, 0) AS inventory_reserved, COALESCE(inv.inventory_inbound, 0) AS inventory_inbound").
		Joins("LEFT JOIN (?) AS inv ON inv.product_id = product.id", invQuery).
		Where("product.id = ?", id).
		First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) ListByIDs(ids []uint64) ([]domain.Product, error) {
	if len(ids) == 0 {
		return []domain.Product{}, nil
	}

	var products []domain.Product
	if err := r.applyProductDisplayJoins(r.db.Table("product")).
		Select(productListByIDsSelect()).
		Where("product.id IN ?", ids).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) ListByParentID(parentID uint64) ([]domain.Product, error) {
	var products []domain.Product

	invQuery := r.db.Table("inventory_balance").
		Select("product_id, SUM(available_quantity) AS inventory_available, SUM(reserved_quantity) AS inventory_reserved, SUM(purchasing_in_transit + pending_inspection + pending_shipment + logistics_in_transit) AS inventory_inbound").
		Group("product_id")

	if err := r.applyProductDisplayJoins(r.db.Table("product")).
		Select(productDisplaySelect()+", COALESCE(inv.inventory_available, 0) AS inventory_available, COALESCE(inv.inventory_reserved, 0) AS inventory_reserved, COALESCE(inv.inventory_inbound, 0) AS inventory_inbound").
		Joins("LEFT JOIN (?) AS inv ON inv.product_id = product.id", invQuery).
		Where("product.parent_id = ?", parentID).
		Order(productListOrderBy()).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) CountReferencesByIDs(ids []uint64) (map[uint64]int64, error) {
	result := make(map[uint64]int64, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	for _, id := range ids {
		result[id] = 0
	}

	type refSpec struct {
		table  string
		column string
	}
	specs := []refSpec{
		{table: "product_combo", column: "main_product_id"},
		{table: "product_combo", column: "product_id"},
		{table: "product_packaging_items", column: "product_id"},
		{table: "product_supplier_quote", column: "product_id"},
		{table: "procurement_replenishment_policy", column: "product_id"},
		{table: "procurement_replenishment_plan", column: "product_id"},
		{table: "purchase_order_item", column: "product_id"},
		{table: "sales_order_item", column: "product_id"},
		{table: "shipment_item", column: "product_id"},
		{table: "inventory_balance", column: "product_id"},
		{table: "inventory_lot", column: "product_id"},
		{table: "inventory_movement", column: "product_id"},
		{table: "finance_cost_event", column: "product_id"},
		{table: "finance_order_cost_detail", column: "product_id"},
		{table: "costing_snapshot", column: "product_id"},
	}

	type row struct {
		ProductID uint64 `gorm:"column:product_id"`
		RefCount  int64  `gorm:"column:ref_count"`
	}

	for _, spec := range specs {
		var rows []row
		if err := r.db.Table(spec.table).
			Select(spec.column+" AS product_id, COUNT(*) AS ref_count").
			Where(spec.column+" IN ?", ids).
			Group(spec.column).
			Scan(&rows).Error; err != nil {
			return nil, err
		}
		for _, item := range rows {
			result[item.ProductID] += item.RefCount
		}
	}

	return result, nil
}

func (r *productRepository) CountByConfigReference(configType domain.ProductConfigType, configID uint64) (int64, error) {
	var total int64

	column := ""
	switch configType {
	case domain.ProductConfigTypeBrand:
		column = "brand_id"
	case domain.ProductConfigTypeSalesStatus:
		column = "status"
	case domain.ProductConfigTypeDimensionUnit:
		column = "dimension_unit_id"
	case domain.ProductConfigTypeWeightUnit:
		column = "weight_unit_id"
	default:
		return 0, nil
	}

	query := r.db.Model(&domain.Product{})
	if configType == domain.ProductConfigTypeSalesStatus {
		var item domain.ProductConfigItem
		if err := r.db.First(&item, configID).Error; err != nil {
			return 0, err
		}
		if err := query.Where(column+" = ?", item.ItemCode).Count(&total).Error; err != nil {
			return 0, err
		}
		return total, nil
	}

	if err := query.Where(column+" = ?", configID).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *productRepository) CountByCategoryID(categoryID uint64) (int64, error) {
	var total int64
	if err := r.db.Model(&domain.Product{}).
		Where("category_id = ?", categoryID).
		Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *productRepository) Create(product *domain.Product) error {
	now := time.Now()
	payload := map[string]any{
		"seller_sku":             product.SellerSku,
		"asin":                   product.Asin,
		"title":                  product.Title,
		"fnsku":                  product.Fnsku,
		"marketplace":            product.Marketplace,
		"parent_id":              product.ParentID,
		"combo_id":               product.ComboID,
		"is_combo_main":          product.IsComboMain,
		"supplier_id":            product.SupplierID,
		"brand_id":               product.BrandID,
		"category_id":            product.CategoryID,
		"dimension_unit_id":      product.DimensionUnitID,
		"weight_unit_id":         product.WeightUnitID,
		"is_inspection_required": product.IsInspectionRequired,
		"is_packing_required":    product.IsPackingRequired,
		"unit_cost":              product.UnitCost,
		"weight":                 product.Weight,
		"length":                 product.Length,
		"width":                  product.Width,
		"height":                 product.Height,
		"dimensions":             product.Dimensions,
		"status":                 product.Status,
		"image_url":              product.ImageUrl,
		"images":                 product.Images,
		"remark":                 product.Remark,
		"gmt_create":             now,
		"gmt_modified":           now,
	}
	if err := r.db.Table("product").Create(payload).Error; err != nil {
		return err
	}
	var created domain.Product
	if err := r.db.Table("product").
		Where("seller_sku = ? AND marketplace = ?", product.SellerSku, product.Marketplace).
		Order("id DESC").
		First(&created).Error; err != nil {
		return err
	}
	*product = created
	return nil
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", product.ID).
		Updates(map[string]any{
			"seller_sku":             product.SellerSku,
			"asin":                   product.Asin,
			"title":                  product.Title,
			"fnsku":                  product.Fnsku,
			"marketplace":            product.Marketplace,
			"parent_id":              product.ParentID,
			"combo_id":               product.ComboID,
			"is_combo_main":          product.IsComboMain,
			"supplier_id":            product.SupplierID,
			"brand_id":               product.BrandID,
			"category_id":            product.CategoryID,
			"dimension_unit_id":      product.DimensionUnitID,
			"weight_unit_id":         product.WeightUnitID,
			"is_inspection_required": product.IsInspectionRequired,
			"is_packing_required":    product.IsPackingRequired,
			"unit_cost":              product.UnitCost,
			"weight":                 product.Weight,
			"length":                 product.Length,
			"width":                  product.Width,
			"height":                 product.Height,
			"dimensions":             product.Dimensions,
			"status":                 product.Status,
			"image_url":              product.ImageUrl,
			"images":                 product.Images,
			"remark":                 product.Remark,
		}).Error
}

func (r *productRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepository) UpdateParentIDBatch(productIDs []uint64, parentID *uint64) error {
	if len(productIDs) == 0 {
		return nil
	}

	return r.db.Model(&domain.Product{}).
		Where("id IN ?", productIDs).
		Update("parent_id", parentID).Error
}

func (r *productRepository) UpdateImageUrl(id uint64, imageUrl string) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", id).
		Update("image_url", imageUrl).Error
}

func (r *productRepository) UpdateUnitCost(productID uint64, unitCost float64) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", productID).
		Update("unit_cost", unitCost).Error
}

func (r *productRepository) GetDefaultSupplierID(productID uint64) (uint64, error) {
	var row struct {
		SupplierID uint64 `gorm:"column:supplier_id"`
	}
	if err := r.db.Model(&domain.Product{}).
		Select("supplier_id").
		Where("id = ?", productID).
		First(&row).Error; err != nil {
		return 0, err
	}
	return row.SupplierID, nil
}

func (r *productRepository) UpdateDefaultSupplierID(productID, supplierID uint64) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", productID).
		Update("supplier_id", supplierID).Error
}

func (r *productRepository) UpdateComboInfo(comboID uint64, mainProductID uint64, productIDs []uint64) error {
	if len(productIDs) == 0 {
		return nil
	}

	if err := r.db.Model(&domain.Product{}).
		Where("id IN ?", productIDs).
		Update("combo_id", comboID).Error; err != nil {
		return err
	}

	if err := r.db.Model(&domain.Product{}).
		Where("id = ?", mainProductID).
		Update("is_combo_main", 1).Error; err != nil {
		return err
	}

	return r.db.Model(&domain.Product{}).
		Where("id IN ?", productIDs).
		Where("id <> ?", mainProductID).
		Update("is_combo_main", 0).Error
}

func (r *productRepository) ClearComboInfo(comboID uint64) error {
	return r.db.Model(&domain.Product{}).
		Where("combo_id = ?", comboID).
		Updates(map[string]any{
			"combo_id":      gorm.Expr("NULL"),
			"is_combo_main": 0,
		}).Error
}
