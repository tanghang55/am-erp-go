package usecase

import (
	"context"
	"fmt"
	"strings"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

type RefundSync interface {
	SyncRefunds(ctx context.Context, trigger integrationDomain.SyncTrigger, operatorID *uint64) (*RefundSyncResult, error)
	GetState() (*integrationDomain.RefundSyncState, error)
	ListRuns(page int, pageSize int) ([]integrationDomain.RefundSyncRun, int64, error)
}

type RefundSyncRegistry struct {
	services map[string]RefundSync
}

func NewRefundSyncRegistry() *RefundSyncRegistry {
	return &RefundSyncRegistry{
		services: map[string]RefundSync{},
	}
}

func (r *RefundSyncRegistry) Register(provider string, service RefundSync) {
	key := normalizeRefundProvider(provider)
	if key == "" || service == nil {
		return
	}
	r.services[key] = service
}

func (r *RefundSyncRegistry) SyncRefunds(ctx context.Context, provider string, trigger integrationDomain.SyncTrigger, operatorID *uint64) (*RefundSyncResult, error) {
	svc, err := r.getService(provider)
	if err != nil {
		return nil, err
	}
	return svc.SyncRefunds(ctx, trigger, operatorID)
}

func (r *RefundSyncRegistry) GetState(provider string) (*integrationDomain.RefundSyncState, error) {
	svc, err := r.getService(provider)
	if err != nil {
		return nil, err
	}
	return svc.GetState()
}

func (r *RefundSyncRegistry) ListRuns(provider string, page int, pageSize int) ([]integrationDomain.RefundSyncRun, int64, error) {
	svc, err := r.getService(provider)
	if err != nil {
		return nil, 0, err
	}
	return svc.ListRuns(page, pageSize)
}

func (r *RefundSyncRegistry) getService(provider string) (RefundSync, error) {
	if r == nil {
		return nil, fmt.Errorf("refund sync registry not configured")
	}
	svc := r.services[normalizeRefundProvider(provider)]
	if svc == nil {
		return nil, fmt.Errorf("provider not supported: %s", provider)
	}
	return svc, nil
}

func normalizeRefundProvider(provider string) string {
	return strings.TrimSpace(strings.ToLower(provider))
}
