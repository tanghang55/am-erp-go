package repository

import (
	"errors"

	integrationDomain "am-erp-go/internal/module/integration/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type refundSyncRepository struct {
	db *gorm.DB
}

func NewRefundSyncRepository(db *gorm.DB) *refundSyncRepository {
	return &refundSyncRepository{db: db}
}

func (r *refundSyncRepository) GetState(provider string, channel string) (*integrationDomain.RefundSyncState, error) {
	var state integrationDomain.RefundSyncState
	err := r.db.Where("provider = ?", provider).Where("channel = ?", channel).Take(&state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &state, nil
}

func (r *refundSyncRepository) SaveState(state *integrationDomain.RefundSyncState) error {
	if state == nil {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "provider"}, {Name: "channel"}},
		DoUpdates: clause.Assignments(map[string]any{
			"last_posted_after":     state.LastPostedAfter,
			"last_sync_started_at":  state.LastSyncStarted,
			"last_sync_finished_at": state.LastSyncFinished,
			"gmt_modified":          gorm.Expr("CURRENT_TIMESTAMP"),
		}),
	}).Create(state).Error
}

func (r *refundSyncRepository) CreateRun(run *integrationDomain.RefundSyncRun) error {
	if run == nil {
		return nil
	}
	return r.db.Create(run).Error
}

func (r *refundSyncRepository) UpdateRun(run *integrationDomain.RefundSyncRun) error {
	if run == nil {
		return nil
	}
	return r.db.Save(run).Error
}

func (r *refundSyncRepository) ListRuns(provider string, channel string, params *integrationDomain.ListRunsParams) ([]integrationDomain.RefundSyncRun, int64, error) {
	var (
		list  []integrationDomain.RefundSyncRun
		total int64
	)
	if params == nil {
		params = &integrationDomain.ListRunsParams{Page: 1, PageSize: 20}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	query := r.db.Model(&integrationDomain.RefundSyncRun{}).Where("provider = ?", provider).Where("channel = ?", channel)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []integrationDomain.RefundSyncRun{}, 0, nil
	}
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(params.PageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	if len(list) == 0 {
		return []integrationDomain.RefundSyncRun{}, total, nil
	}
	return list, total, nil
}

func (r *refundSyncRepository) UpsertEvents(events []integrationDomain.ThirdPartyRefundEvent) error {
	if len(events) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "provider"},
			{Name: "channel"},
			{Name: "refund_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"order_id":      gorm.Expr("VALUES(order_id)"),
			"order_item_id": gorm.Expr("VALUES(order_item_id)"),
			"seller_sku":    gorm.Expr("VALUES(seller_sku)"),
			"marketplace":   gorm.Expr("VALUES(marketplace)"),
			"product_id":    gorm.Expr("VALUES(product_id)"),
			"qty_refunded":  gorm.Expr("VALUES(qty_refunded)"),
			"refund_amount": gorm.Expr("VALUES(refund_amount)"),
			"currency":      gorm.Expr("VALUES(currency)"),
			"posted_at":     gorm.Expr("VALUES(posted_at)"),
			"status":        gorm.Expr("VALUES(status)"),
			"error_message": gorm.Expr("VALUES(error_message)"),
			"raw_payload":   gorm.Expr("VALUES(raw_payload)"),
			"gmt_modified":  gorm.Expr("CURRENT_TIMESTAMP"),
		}),
	}).Create(&events).Error
}
