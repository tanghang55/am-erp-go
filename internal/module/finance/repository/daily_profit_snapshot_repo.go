package repository

import (
	"am-erp-go/internal/module/finance/domain"
	"time"

	"gorm.io/gorm"
)

type dailyProfitSnapshotRepository struct {
	db *gorm.DB
}

func NewDailyProfitSnapshotRepository(db *gorm.DB) domain.DailyProfitSnapshotRepository {
	return &dailyProfitSnapshotRepository{db: db}
}

func (r *dailyProfitSnapshotRepository) DeleteByDate(bizDate time.Time, marketplace *string) error {
	query := r.db.Where("biz_date = ?", bizDate.Format("2006-01-02"))
	if marketplace != nil && *marketplace != "" {
		query = query.Where("marketplace = ?", *marketplace)
	}
	return query.Delete(&domain.DailyProfitSnapshot{}).Error
}

func (r *dailyProfitSnapshotRepository) CreateBatch(snapshots []domain.DailyProfitSnapshot) error {
	if len(snapshots) == 0 {
		return nil
	}
	return r.db.Create(&snapshots).Error
}

func (r *dailyProfitSnapshotRepository) List(params *domain.DailyProfitSnapshotListParams) ([]domain.DailyProfitSnapshot, error) {
	items := make([]domain.DailyProfitSnapshot, 0)
	query := r.db.Model(&domain.DailyProfitSnapshot{})

	if !params.DateFrom.IsZero() {
		query = query.Where("biz_date >= ?", params.DateFrom.Format("2006-01-02"))
	}
	if !params.DateTo.IsZero() {
		query = query.Where("biz_date <= ?", params.DateTo.Format("2006-01-02"))
	}
	if params.Marketplace != "" {
		query = query.Where("marketplace = ?", params.Marketplace)
	}

	if err := query.Order("biz_date ASC, marketplace ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
