package usecase

import (
	"context"
	"testing"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

func TestSyncRefundsSuccessUpdatesState(t *testing.T) {
	now := time.Date(2026, 3, 24, 14, 0, 0, 0, time.UTC)
	repo := &fakeRefundSyncRepository{}
	provider := &fakeRefundProvider{
		code: "AMAZON_US",
		refunds: []integrationDomain.ExternalRefund{
			{
				RefundID:      "R-1",
				OrderID:       "AMZ-ORDER-1",
				OrderItemID:   "ITEM-1",
				SellerSKU:     "SKU-REF-1",
				MarketplaceID: "ATVPDKIKX0DER",
				Quantity:      1,
				Amount:        19.9,
				Currency:      "USD",
				PostedAt:      now.Add(-30 * time.Minute),
			},
		},
	}
	resolver := &fakeRefundProductResolver{
		skuMap: map[string]uint64{
			"SKU-REF-1|US": 2001,
		},
	}
	service := NewRefundSyncService(
		repo,
		provider,
		RefundSyncConfig{
			ProviderCode:       "AMAZON_US",
			Channel:            "SP_API_SELF",
			LookbackMinutes:    10,
			InitialLookbackDay: 7,
			MarketplaceIDs:     []string{"ATVPDKIKX0DER"},
		},
		func() time.Time { return now },
		resolver,
		resolver,
	)

	result, err := service.SyncRefunds(context.Background(), integrationDomain.SyncTriggerManual, nil)
	if err != nil {
		t.Fatalf("SyncRefunds error: %v", err)
	}
	if result.FetchedRefunds != 1 || result.ImportedRefunds != 1 || result.ErrorRefunds != 0 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if repo.savedState == nil || repo.savedState.LastPostedAfter == nil {
		t.Fatalf("expected state saved")
	}
	if len(repo.upserted) != 1 {
		t.Fatalf("expected one upserted refund event, got %d", len(repo.upserted))
	}
	if repo.upserted[0].ProductID == nil || *repo.upserted[0].ProductID != 2001 {
		t.Fatalf("expected mapped product_id=2001, got %+v", repo.upserted[0].ProductID)
	}
	if repo.upserted[0].Status != integrationDomain.RefundEventStatusMapped {
		t.Fatalf("expected status MAPPED, got %s", repo.upserted[0].Status)
	}
}

func TestSyncRefundsUnmappedRecordsError(t *testing.T) {
	now := time.Date(2026, 3, 24, 15, 0, 0, 0, time.UTC)
	repo := &fakeRefundSyncRepository{}
	provider := &fakeRefundProvider{
		code: "AMAZON_US",
		refunds: []integrationDomain.ExternalRefund{
			{
				RefundID:      "R-2",
				OrderID:       "AMZ-ORDER-2",
				OrderItemID:   "ITEM-2",
				SellerSKU:     "SKU-REF-NOTFOUND",
				MarketplaceID: "ATVPDKIKX0DER",
				Quantity:      1,
				Amount:        9.9,
				Currency:      "USD",
				PostedAt:      now.Add(-20 * time.Minute),
			},
		},
	}
	service := NewRefundSyncService(
		repo,
		provider,
		RefundSyncConfig{
			ProviderCode:       "AMAZON_US",
			Channel:            "SP_API_SELF",
			LookbackMinutes:    10,
			InitialLookbackDay: 7,
			MarketplaceIDs:     []string{"ATVPDKIKX0DER"},
		},
		func() time.Time { return now },
		&fakeRefundProductResolver{skuMap: map[string]uint64{}},
		&fakeRefundProductResolver{skuMap: map[string]uint64{}},
	)

	result, err := service.SyncRefunds(context.Background(), integrationDomain.SyncTriggerManual, nil)
	if err != nil {
		t.Fatalf("SyncRefunds error: %v", err)
	}
	if result.ErrorRefunds != 1 {
		t.Fatalf("expected one error refund, got %+v", result)
	}
	if len(repo.upserted) != 1 {
		t.Fatalf("expected one upserted refund event, got %d", len(repo.upserted))
	}
	if repo.upserted[0].Status != integrationDomain.RefundEventStatusUnmapped {
		t.Fatalf("expected status UNMAPPED, got %s", repo.upserted[0].Status)
	}
	if repo.upserted[0].ErrorMessage == nil || *repo.upserted[0].ErrorMessage == "" {
		t.Fatalf("expected error message when refund unmapped")
	}
}

type fakeRefundSyncRepository struct {
	state      *integrationDomain.RefundSyncState
	savedState *integrationDomain.RefundSyncState
	createdRun *integrationDomain.RefundSyncRun
	updatedRun *integrationDomain.RefundSyncRun
	upserted   []integrationDomain.ThirdPartyRefundEvent
}

func (f *fakeRefundSyncRepository) GetState(provider string, channel string) (*integrationDomain.RefundSyncState, error) {
	_ = provider
	_ = channel
	return f.state, nil
}

func (f *fakeRefundSyncRepository) SaveState(state *integrationDomain.RefundSyncState) error {
	f.savedState = state
	f.state = state
	return nil
}

func (f *fakeRefundSyncRepository) CreateRun(run *integrationDomain.RefundSyncRun) error {
	run.ID = 1
	f.createdRun = run
	return nil
}

func (f *fakeRefundSyncRepository) UpdateRun(run *integrationDomain.RefundSyncRun) error {
	f.updatedRun = run
	return nil
}

func (f *fakeRefundSyncRepository) ListRuns(provider string, channel string, params *integrationDomain.ListRunsParams) ([]integrationDomain.RefundSyncRun, int64, error) {
	_ = provider
	_ = channel
	_ = params
	return []integrationDomain.RefundSyncRun{}, 0, nil
}

func (f *fakeRefundSyncRepository) UpsertEvents(events []integrationDomain.ThirdPartyRefundEvent) error {
	f.upserted = append(f.upserted, events...)
	return nil
}

type fakeRefundProvider struct {
	code    string
	refunds []integrationDomain.ExternalRefund
}

func (f *fakeRefundProvider) Code() string {
	return f.code
}

func (f *fakeRefundProvider) ListRefunds(ctx context.Context, req integrationDomain.ExternalListRefundsRequest) ([]integrationDomain.ExternalRefund, error) {
	_ = ctx
	_ = req
	return f.refunds, nil
}

type fakeRefundProductResolver struct {
	skuMap map[string]uint64
}

func (f *fakeRefundProductResolver) ResolveMappedProductID(providerCode string, sellerSKU string, marketplace string) (uint64, error) {
	_ = providerCode
	return f.skuMap[sellerSKU+"|"+marketplace], nil
}

func (f *fakeRefundProductResolver) ResolveProductIDBySellerSKU(sellerSKU string, marketplace string) (uint64, error) {
	return f.skuMap[sellerSKU+"|"+marketplace], nil
}
