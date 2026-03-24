package repository

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"

	"am-erp-go/internal/module/sales/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type salesOrderRepository struct {
	db *gorm.DB
}

func NewSalesOrderRepository(db *gorm.DB) *salesOrderRepository {
	return &salesOrderRepository{db: db}
}

func (r *salesOrderRepository) List(params *domain.SalesOrderListParams) ([]domain.SalesOrder, int64, error) {
	var orders []domain.SalesOrder
	var total int64

	query := r.db.Model(&domain.SalesOrder{})

	if params != nil {
		if params.Status != "" {
			query = query.Where("order_status = ?", params.Status)
		}
		if params.Marketplace != "" {
			query = query.Where("marketplace = ?", params.Marketplace)
		}
		if params.Keyword != "" {
			keyword := "%" + params.Keyword + "%"
			query = query.Where("order_no LIKE ? OR external_order_no LIKE ? OR remark LIKE ?", keyword, keyword, keyword)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := 1
	pageSize := 20
	if params != nil {
		if params.Page > 0 {
			page = params.Page
		}
		if params.PageSize > 0 {
			pageSize = params.PageSize
		}
	}

	offset := (page - 1) * pageSize
	if err := query.Order("gmt_modified DESC").Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	if len(orders) == 0 {
		return []domain.SalesOrder{}, total, nil
	}

	orderIDs := make([]uint64, 0, len(orders))
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	itemsMap, err := r.listItemsByOrderIDs(orderIDs)
	if err != nil {
		return nil, 0, err
	}
	for i := range orders {
		orders[i].Items = itemsMap[orders[i].ID]
	}

	return orders, total, nil
}

func (r *salesOrderRepository) GetByID(id uint64) (*domain.SalesOrder, error) {
	var order domain.SalesOrder
	if err := r.db.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	itemsMap, err := r.listItemsByOrderIDs([]uint64{id})
	if err != nil {
		return nil, err
	}
	order.Items = itemsMap[id]
	return &order, nil
}

func (r *salesOrderRepository) Create(order *domain.SalesOrder) error {
	if order == nil {
		return nil
	}
	r.prepareOrder(order)

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if len(order.Items) == 0 {
			return nil
		}

		items := make([]domain.SalesOrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			item.SalesOrderID = order.ID
			items = append(items, item)
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		order.Items = items
		return nil
	})
}

func (r *salesOrderRepository) Update(order *domain.SalesOrder) error {
	if order == nil {
		return nil
	}
	r.prepareOrder(order)

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		if order.Items == nil {
			return nil
		}

		if err := tx.Where("sales_order_id = ?", order.ID).Delete(&domain.SalesOrderItem{}).Error; err != nil {
			return err
		}
		if len(order.Items) == 0 {
			return nil
		}

		items := make([]domain.SalesOrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			item.SalesOrderID = order.ID
			items = append(items, item)
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		order.Items = items
		return nil
	})
}

func (r *salesOrderRepository) listItemsByOrderIDs(orderIDs []uint64) (map[uint64][]domain.SalesOrderItem, error) {
	result := map[uint64][]domain.SalesOrderItem{}
	if len(orderIDs) == 0 {
		return result, nil
	}

	var items []domain.SalesOrderItem
	if err := r.db.Table("sales_order_item item").
		Select(`
			item.*,
			product.seller_sku AS seller_sku,
			product.title AS product_title,
			product.image_url AS product_image_url
		`).
		Joins("LEFT JOIN product ON product.id = item.product_id").
		Where("item.sales_order_id IN ?", orderIDs).
		Order("item.id ASC").
		Scan(&items).Error; err != nil {
		return nil, err
	}
	for _, item := range items {
		result[item.SalesOrderID] = append(result[item.SalesOrderID], item)
	}
	return result, nil
}

func (r *salesOrderRepository) GetImportByFileHash(fileHash string) (*domain.ReportImport, error) {
	var batch domain.ReportImport
	if err := r.db.Where("file_hash = ?", fileHash).First(&batch).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &batch, nil
}

func (r *salesOrderRepository) CreateImport(batch *domain.ReportImport) error {
	if batch == nil {
		return nil
	}
	return r.db.Create(batch).Error
}

func (r *salesOrderRepository) UpdateImport(batch *domain.ReportImport) error {
	if batch == nil {
		return nil
	}
	return r.db.Save(batch).Error
}

func (r *salesOrderRepository) ListImports(page int, pageSize int) ([]domain.ReportImport, int64, error) {
	var list []domain.ReportImport
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := r.db.Model(&domain.ReportImport{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	if len(list) == 0 {
		return []domain.ReportImport{}, total, nil
	}
	return list, total, nil
}

func (r *salesOrderRepository) GetImportByID(id uint64) (*domain.ReportImport, error) {
	var batch domain.ReportImport
	if err := r.db.First(&batch, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &batch, nil
}

func (r *salesOrderRepository) ListImportErrors(importID uint64) ([]domain.ReportImportRowError, error) {
	var rows []domain.ReportImportRowError
	if err := r.db.Where("report_import_id = ?", importID).Order("row_no ASC, id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []domain.ReportImportRowError{}, nil
	}
	return rows, nil
}

func (r *salesOrderRepository) InsertImportRowErrors(errorsRows []domain.ReportImportRowError) error {
	if len(errorsRows) == 0 {
		return nil
	}
	return r.db.Create(&errorsRows).Error
}

func (r *salesOrderRepository) ResolveProductIDBySellerSKU(sellerSKU string, marketplace string) (uint64, error) {
	var row struct {
		ID uint64 `gorm:"column:id"`
	}
	query := r.db.Table("product").Select("id").Where("seller_sku = ?", sellerSKU)
	if strings.TrimSpace(marketplace) != "" {
		query = query.Where("marketplace = ?", strings.TrimSpace(marketplace))
	}
	if err := query.Take(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return row.ID, nil
}

func (r *salesOrderRepository) UpsertImportedOrderLine(line *domain.ImportOrderLine, productID uint64, batchNo string, operatorID *uint64) error {
	if line == nil {
		return nil
	}
	if line.SourceType == "" {
		line.SourceType = "MANUAL_IMPORT"
	}
	if line.ExternalOrderNo == "" {
		return domain.ErrImportInvalidFile
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		order, err := r.findOrCreateImportOrder(tx, line, batchNo, operatorID)
		if err != nil {
			return err
		}

		external := line.ExternalOrderNo
		item := domain.SalesOrderItem{
			SalesOrderID:    order.ID,
			LineNo:          line.LineNo,
			SourceType:      line.SourceType,
			ExternalOrderNo: &external,
			ProductID:       productID,
			QtyOrdered:      line.Qty,
			QtyAllocated:    0,
			QtyShipped:      0,
			QtyReturned:     0,
			UnitPrice:       line.UnitPrice,
			Subtotal:        line.UnitPrice * float64(line.Qty),
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "source_type"},
				{Name: "external_order_no"},
				{Name: "line_no"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"sales_order_id": order.ID,
				"product_id":     productID,
				"qty_ordered":    line.Qty,
				"unit_price":     item.UnitPrice,
				"subtotal":       item.Subtotal,
				"gmt_modified":   gorm.Expr("CURRENT_TIMESTAMP"),
			}),
		}).Create(&item).Error; err != nil {
			return err
		}

		var amount float64
		if err := tx.Model(&domain.SalesOrderItem{}).
			Where("sales_order_id = ?", order.ID).
			Select("COALESCE(SUM(subtotal),0)").
			Scan(&amount).Error; err != nil {
			return err
		}

		updates := map[string]any{
			"order_amount":    amount,
			"order_date":      line.OrderDate,
			"import_batch_no": batchNo,
			"gmt_modified":    gorm.Expr("CURRENT_TIMESTAMP"),
		}
		if line.Marketplace != "" {
			updates["marketplace"] = line.Marketplace
		}
		if line.Currency != "" {
			updates["currency"] = line.Currency
		}
		if line.SalesChannel != nil {
			updates["sales_channel"] = *line.SalesChannel
		}
		if operatorID != nil {
			updates["updated_by"] = *operatorID
		}

		return tx.Model(&domain.SalesOrder{}).Where("id = ?", order.ID).Updates(updates).Error
	})
}

func (r *salesOrderRepository) findOrCreateImportOrder(
	tx *gorm.DB,
	line *domain.ImportOrderLine,
	batchNo string,
	operatorID *uint64,
) (*domain.SalesOrder, error) {
	var order domain.SalesOrder
	err := tx.Where("source_type = ? AND external_order_no = ?", line.SourceType, line.ExternalOrderNo).Take(&order).Error
	if err == nil {
		return &order, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	currency := line.Currency
	if currency == "" {
		return nil, domain.ErrImportInvalidFile
	}
	orderNo := buildSystemOrderNo(line.SourceType, line.ExternalOrderNo)
	external := line.ExternalOrderNo
	order = domain.SalesOrder{
		OrderNo:         orderNo,
		SourceType:      line.SourceType,
		StockPool:       defaultStockPoolBySourceType(line.SourceType),
		ExternalOrderNo: &external,
		SalesChannel:    line.SalesChannel,
		Marketplace:     strPtrOrNil(line.Marketplace),
		OrderStatus:     domain.SalesOrderStatusDraft,
		OrderDate:       line.OrderDate,
		Currency:        currency,
		OrderAmount:     line.UnitPrice * float64(line.Qty),
		ImportBatchNo:   strPtrOrNil(batchNo),
		CreatedBy:       operatorID,
		UpdatedBy:       operatorID,
	}
	if err := tx.Create(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func buildSystemOrderNo(sourceType string, externalOrderNo string) string {
	orderNo := fmt.Sprintf("SOI-%s-%s", sourceType, externalOrderNo)
	orderNo = strings.ToUpper(strings.ReplaceAll(orderNo, " ", "-"))
	if len(orderNo) <= 64 {
		return orderNo
	}
	hash := sha1.Sum([]byte(sourceType + ":" + externalOrderNo))
	return fmt.Sprintf("SOI-%x", hash[:12])
}

func strPtrOrNil(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

func (r *salesOrderRepository) prepareOrder(order *domain.SalesOrder) {
	if order.SourceType == "" {
		order.SourceType = "MANUAL_CREATE"
	}
	if order.StockPool == "" {
		order.StockPool = defaultStockPoolBySourceType(order.SourceType)
	}
	if order.OrderStatus == "" {
		order.OrderStatus = domain.SalesOrderStatusDraft
	}

	for i := range order.Items {
		if order.Items[i].SourceType == "" {
			order.Items[i].SourceType = order.SourceType
		}
		if order.Items[i].ExternalOrderNo == nil && order.ExternalOrderNo != nil {
			order.Items[i].ExternalOrderNo = order.ExternalOrderNo
		}
	}
}

func defaultStockPoolBySourceType(sourceType string) domain.StockPool {
	switch strings.ToUpper(strings.TrimSpace(sourceType)) {
	case "AMAZON_API", "AMAZON_IMPORT", "THIRD_PARTY_API":
		return domain.StockPoolSellable
	default:
		return domain.StockPoolAvailable
	}
}
