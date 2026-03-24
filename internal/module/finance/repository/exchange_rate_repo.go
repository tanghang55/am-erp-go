package repository

import (
	"strings"
	"time"

	"am-erp-go/internal/module/finance/domain"

	"gorm.io/gorm"
)

type exchangeRateRepository struct {
	db *gorm.DB
}

func NewExchangeRateRepository(db *gorm.DB) domain.ExchangeRateRepository {
	return &exchangeRateRepository{db: db}
}

func (r *exchangeRateRepository) List(params *domain.ExchangeRateListParams) ([]domain.ExchangeRate, int64, error) {
	query := r.db.Model(&domain.ExchangeRate{})
	if params != nil {
		if from := strings.TrimSpace(params.FromCurrency); from != "" {
			query = query.Where("from_currency = ?", strings.ToUpper(from))
		}
		if to := strings.TrimSpace(params.ToCurrency); to != "" {
			query = query.Where("to_currency = ?", strings.ToUpper(to))
		}
		if params.Status != "" {
			query = query.Where("status = ?", params.Status)
		}
	}

	var total int64
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
	items := make([]domain.ExchangeRate, 0)
	if err := query.Order("effective_at DESC, id DESC").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *exchangeRateRepository) Create(rate *domain.ExchangeRate) error {
	return r.db.Create(rate).Error
}

func (r *exchangeRateRepository) GetByID(id uint64) (*domain.ExchangeRate, error) {
	var item domain.ExchangeRate
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *exchangeRateRepository) UpdateStatus(id uint64, status domain.ExchangeRateStatus, operatorID uint64) error {
	tx := r.db.Model(&domain.ExchangeRate{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"updated_by":   operatorID,
			"gmt_modified": time.Now(),
		})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *exchangeRateRepository) FindEffectiveRate(fromCurrency, toCurrency string, occurredAt time.Time) (*domain.ExchangeRate, error) {
	var item domain.ExchangeRate
	err := r.db.
		Where("from_currency = ? AND to_currency = ? AND status = ? AND effective_at <= ?",
			strings.ToUpper(strings.TrimSpace(fromCurrency)),
			strings.ToUpper(strings.TrimSpace(toCurrency)),
			domain.ExchangeRateStatusActive,
			occurredAt,
		).
		Order("effective_at DESC, id DESC").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
