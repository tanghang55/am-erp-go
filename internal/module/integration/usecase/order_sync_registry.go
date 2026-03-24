package usecase

import (
	"context"
	"fmt"
	"strings"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

type ProviderOrderSyncService interface {
	SyncOrders(ctx context.Context, trigger integrationDomain.SyncTrigger, operatorID *uint64) (*OrderSyncResult, error)
	GetState() (*integrationDomain.OrderSyncState, error)
	ListRuns(page int, pageSize int) ([]integrationDomain.OrderSyncRun, int64, error)
}

type OrderSyncRegistry struct {
	services map[string]ProviderOrderSyncService
}

func NewOrderSyncRegistry() *OrderSyncRegistry {
	return &OrderSyncRegistry{
		services: map[string]ProviderOrderSyncService{},
	}
}

func (r *OrderSyncRegistry) Register(provider string, service ProviderOrderSyncService) {
	key := normalizeProvider(provider)
	if key == "" || service == nil {
		return
	}
	r.services[key] = service
}

func (r *OrderSyncRegistry) SyncOrders(ctx context.Context, provider string, trigger integrationDomain.SyncTrigger, operatorID *uint64) (*OrderSyncResult, error) {
	svc, err := r.getService(provider)
	if err != nil {
		return nil, err
	}
	return svc.SyncOrders(ctx, trigger, operatorID)
}

func (r *OrderSyncRegistry) GetState(provider string) (*integrationDomain.OrderSyncState, error) {
	svc, err := r.getService(provider)
	if err != nil {
		return nil, err
	}
	return svc.GetState()
}

func (r *OrderSyncRegistry) ListRuns(provider string, page int, pageSize int) ([]integrationDomain.OrderSyncRun, int64, error) {
	svc, err := r.getService(provider)
	if err != nil {
		return nil, 0, err
	}
	return svc.ListRuns(page, pageSize)
}

func (r *OrderSyncRegistry) getService(provider string) (ProviderOrderSyncService, error) {
	if r == nil {
		return nil, fmt.Errorf("order sync registry not configured")
	}
	svc, ok := r.services[normalizeProvider(provider)]
	if !ok {
		return nil, fmt.Errorf("provider not supported: %s", provider)
	}
	return svc, nil
}

func normalizeProvider(provider string) string {
	return strings.TrimSpace(strings.ToLower(provider))
}
