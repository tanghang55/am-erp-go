package repository

import (
	"errors"

	integrationDomain "am-erp-go/internal/module/integration/domain"

	"gorm.io/gorm"
)

type orderSyncRepository struct {
	db *gorm.DB
}

func NewOrderSyncRepository(db *gorm.DB) *orderSyncRepository {
	return &orderSyncRepository{db: db}
}

func (r *orderSyncRepository) GetState(provider string, channel string) (*integrationDomain.OrderSyncState, error) {
	var state integrationDomain.OrderSyncState
	if err := r.db.Where("provider = ? AND channel = ?", provider, channel).Take(&state).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &state, nil
}

func (r *orderSyncRepository) SaveState(state *integrationDomain.OrderSyncState) error {
	if state == nil {
		return nil
	}
	if state.ID == 0 {
		return r.db.Create(state).Error
	}
	return r.db.Save(state).Error
}

func (r *orderSyncRepository) CreateRun(run *integrationDomain.OrderSyncRun) error {
	if run == nil {
		return nil
	}
	return r.db.Create(run).Error
}

func (r *orderSyncRepository) UpdateRun(run *integrationDomain.OrderSyncRun) error {
	if run == nil {
		return nil
	}
	return r.db.Save(run).Error
}

func (r *orderSyncRepository) ListRuns(provider string, channel string, params *integrationDomain.ListRunsParams) ([]integrationDomain.OrderSyncRun, int64, error) {
	var list []integrationDomain.OrderSyncRun
	var total int64

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

	query := r.db.Model(&integrationDomain.OrderSyncRun{})
	if provider != "" {
		query = query.Where("provider = ?", provider)
	}
	if channel != "" {
		query = query.Where("channel = ?", channel)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []integrationDomain.OrderSyncRun{}, 0, nil
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
