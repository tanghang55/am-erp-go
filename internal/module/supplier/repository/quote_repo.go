package repository

import (
	"am-erp-go/internal/module/supplier/domain"

	"gorm.io/gorm"
)

type quoteRepository struct {
	db *gorm.DB
}

func NewQuoteRepository(db *gorm.DB) domain.QuoteRepository {
	return &quoteRepository{db: db}
}

func (r *quoteRepository) ListByProductIDs(productIDs []uint64) (map[uint64][]domain.ProductSupplierQuote, error) {
	result := map[uint64][]domain.ProductSupplierQuote{}
	if len(productIDs) == 0 {
		return result, nil
	}

	var quotes []domain.ProductSupplierQuote
	if err := r.db.Where("product_id IN ?", productIDs).
		Order("gmt_modified DESC").
		Find(&quotes).Error; err != nil {
		return nil, err
	}

	for _, quote := range quotes {
		result[quote.ProductID] = append(result[quote.ProductID], quote)
	}
	return result, nil
}

func (r *quoteRepository) ListProductsWithQuotes(params *domain.QuoteListParams) ([]domain.ProductQuoteRow, int64, error) {
	if params == nil {
		params = &domain.QuoteListParams{}
	}

	listQuery := r.buildProductQuoteQuery(params)

	var total int64
	countQuery := r.buildProductQuoteQuery(params)
	if err := countQuery.Distinct("product.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := params.Page
	pageSize := params.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var rows []productQuoteRow
	if err := listQuery.Order("product.gmt_modified DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	productIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		productIDs = append(productIDs, row.ProductID)
	}

	quoteMap, err := r.listQuoteDetailsByProductIDs(productIDs)
	if err != nil {
		return nil, 0, err
	}

	result := make([]domain.ProductQuoteRow, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.ProductQuoteRow{
			ProductID:         row.ProductID,
			SellerSku:         row.SellerSku,
			Asin:              row.Asin,
			Marketplace:       row.Marketplace,
			Title:             row.Title,
			ImageUrl:          row.ImageUrl,
			DefaultSupplierID: row.DefaultSupplierID,
			Quotes:            quoteMap[row.ProductID],
		})
	}

	return result, total, nil
}

func (r *quoteRepository) buildProductQuoteQuery(params *domain.QuoteListParams) *gorm.DB {
	query := r.db.Table("product").
		Select("product.id, product.seller_sku, product.asin, product.marketplace, product.title, product.image_url, product.supplier_id,gmt_modified")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("product.seller_sku LIKE ? OR product.asin LIKE ? OR product.title LIKE ?", keyword, keyword, keyword)
	}
	if params.Marketplace != "" {
		query = query.Where("product.marketplace = ?", params.Marketplace)
	}
	if params.SupplierID != nil {
		query = query.Joins("JOIN product_supplier_quote psq ON psq.product_id = product.id AND psq.supplier_id = ?", *params.SupplierID).
			Group("product.id")
	}
	return query
}

func (r *quoteRepository) GetByProductSupplier(productID, supplierID uint64) (*domain.ProductSupplierQuote, error) {
	var quote domain.ProductSupplierQuote
	if err := r.db.Where("product_id = ? AND supplier_id = ?", productID, supplierID).
		First(&quote).Error; err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *quoteRepository) Create(quote *domain.ProductSupplierQuote) error {
	return r.db.Create(quote).Error
}

func (r *quoteRepository) Update(quote *domain.ProductSupplierQuote) error {
	return r.db.Save(quote).Error
}

func (r *quoteRepository) Delete(productID, supplierID uint64) error {
	return r.db.Where("product_id = ? AND supplier_id = ?", productID, supplierID).
		Delete(&domain.ProductSupplierQuote{}).Error
}

type productQuoteRow struct {
	ProductID         uint64 `gorm:"column:id"`
	SellerSku         string `gorm:"column:seller_sku"`
	Asin              string `gorm:"column:asin"`
	Marketplace       string `gorm:"column:marketplace"`
	Title             string `gorm:"column:title"`
	ImageUrl          string `gorm:"column:image_url"`
	DefaultSupplierID uint64 `gorm:"column:supplier_id"`
}

func (r *quoteRepository) listQuoteDetailsByProductIDs(productIDs []uint64) (map[uint64][]domain.ProductSupplierQuoteDetail, error) {
	result := map[uint64][]domain.ProductSupplierQuoteDetail{}
	if len(productIDs) == 0 {
		return result, nil
	}

	var quotes []domain.ProductSupplierQuoteDetail
	if err := r.db.Table("product_supplier_quote").
		Select("product_supplier_quote.*, supplier.name AS supplier_name, supplier.supplier_code AS supplier_code").
		Joins("JOIN supplier ON supplier.id = product_supplier_quote.supplier_id").
		Where("product_supplier_quote.product_id IN ?", productIDs).
		Order("product_supplier_quote.gmt_modified DESC").
		Find(&quotes).Error; err != nil {
		return nil, err
	}

	for _, quote := range quotes {
		result[quote.ProductID] = append(result[quote.ProductID], quote)
	}
	return result, nil
}
