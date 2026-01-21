package repository

import (
	"am-erp-go/internal/module/product/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) List(params *domain.ProductListParams) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.Model(&domain.Product{})

	// 关键词搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("seller_sku LIKE ? OR asin LIKE ? OR title LIKE ?", keyword, keyword, keyword)
	}

	// 站点筛选
	if params.Marketplace != "" {
		query = query.Where("marketplace = ?", params.Marketplace)
	}

	// 状态筛选
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 供应商筛选
	if params.SupplierID != nil {
		query = query.Where("supplier_id = ?", *params.SupplierID)
	}

	// 组合筛选
	if params.ComboID != nil {
		query = query.Where("combo_id = ?", *params.ComboID)
	}

	// 是否主产品
	if params.IsComboMain != nil {
		query = query.Where("is_combo_main = ?", *params.IsComboMain)
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
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) GetByID(id uint64) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) ListByIDs(ids []uint64) ([]domain.Product, error) {
	if len(ids) == 0 {
		return []domain.Product{}, nil
	}

	var products []domain.Product
	if err := r.db.Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.Product{}, id).Error
}

func (r *productRepository) UpdateImageUrl(id uint64, imageUrl string) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", id).
		Update("image_url", imageUrl).Error
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
