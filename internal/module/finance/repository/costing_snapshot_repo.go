package repository

import (
	"am-erp-go/internal/module/finance/domain"
	"time"

	"gorm.io/gorm"
)

type costingSnapshotRepository struct {
	db *gorm.DB
}

func NewCostingSnapshotRepository(db *gorm.DB) domain.CostingSnapshotRepository {
	return &costingSnapshotRepository{db: db}
}

func costingSnapshotSelectSQL() string {
	return `
			costing_snapshot.*,
			COALESCE(product.seller_sku, '') AS seller_sku,
			COALESCE(product.title, '') AS product_title,
			COALESCE(product.image_url, '') AS product_image_url
		`
}

func (r *costingSnapshotRepository) readQuery() *gorm.DB {
	return r.db.Table("costing_snapshot").
		Select(costingSnapshotSelectSQL()).
		Joins("LEFT JOIN product ON product.id = costing_snapshot.product_id")
}

func currentCostTimeWindowWhereSQL() string {
	return "costing_snapshot.effective_from <= ? AND (costing_snapshot.effective_to IS NULL OR costing_snapshot.effective_to > ?)"
}

func currentCostWhere(query *gorm.DB, productID uint64, now time.Time) *gorm.DB {
	return query.
		Where("costing_snapshot.product_id = ?", productID).
		Where(currentCostTimeWindowWhereSQL(), now, now)
}

func currentCostOrderBy() string {
	return "costing_snapshot.cost_type ASC, costing_snapshot.effective_from DESC, costing_snapshot.id DESC"
}

func (r *costingSnapshotRepository) List(params *domain.CostingSnapshotListParams) ([]domain.CostingSnapshot, int64, error) {
	var list []domain.CostingSnapshot
	var total int64

	query := r.readQuery()

	if params.ProductID != nil {
		query = query.Where("product_id = ?", *params.ProductID)
	}
	if params.CostType != "" {
		query = query.Where("cost_type = ?", params.CostType)
	}
	if params.IsCurrent != nil {
		now := time.Now()
		if *params.IsCurrent {
			query = query.Where("effective_from <= ? AND (effective_to IS NULL OR effective_to > ?)", now, now)
		} else {
			query = query.Where("effective_to IS NOT NULL AND effective_to <= ?", now)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Order("effective_from DESC, id DESC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *costingSnapshotRepository) GetByID(id uint64) (*domain.CostingSnapshot, error) {
	var snapshot domain.CostingSnapshot
	if err := r.readQuery().Where("costing_snapshot.id = ?", id).First(&snapshot).Error; err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (r *costingSnapshotRepository) Create(snapshot *domain.CostingSnapshot) error {
	return r.db.Create(snapshot).Error
}

func (r *costingSnapshotRepository) Update(snapshot *domain.CostingSnapshot) error {
	return r.db.Save(snapshot).Error
}

func (r *costingSnapshotRepository) Delete(id uint64) error {
	return r.db.Delete(&domain.CostingSnapshot{}, id).Error
}

func (r *costingSnapshotRepository) ExpireCurrent(productID uint64, costType domain.CostType, effectiveTo time.Time, excludeID *uint64) error {
	query := r.db.Model(&domain.CostingSnapshot{}).
		Where("product_id = ? AND cost_type = ?", productID, costType).
		Where("effective_to IS NULL")

	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}

	return query.Update("effective_to", effectiveTo).Error
}

func (r *costingSnapshotRepository) GetCurrent(productID uint64, costType domain.CostType, now time.Time) (*domain.CostingSnapshot, error) {
	var snapshot domain.CostingSnapshot
	err := currentCostWhere(r.readQuery(), productID, now).
		Where("costing_snapshot.cost_type = ?", costType).
		Order("costing_snapshot.effective_from DESC, costing_snapshot.id DESC").
		First(&snapshot).Error
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (r *costingSnapshotRepository) ListCurrentBySKU(productID uint64, now time.Time) ([]domain.CostingSnapshot, error) {
	var list []domain.CostingSnapshot
	if err := currentCostWhere(r.readQuery(), productID, now).
		Order(currentCostOrderBy()).
		Find(&list).Error; err != nil {
		return nil, err
	}

	// 每个成本类型只保留最新一条
	latest := make(map[domain.CostType]domain.CostingSnapshot)
	for _, item := range list {
		if _, exists := latest[item.CostType]; !exists {
			latest[item.CostType] = item
		}
	}

	result := make([]domain.CostingSnapshot, 0, len(latest))
	for _, costType := range []domain.CostType{
		domain.CostTypePurchase,
		domain.CostTypeLanded,
		domain.CostTypeAverage,
	} {
		if item, exists := latest[costType]; exists {
			result = append(result, item)
		}
	}
	return result, nil
}
