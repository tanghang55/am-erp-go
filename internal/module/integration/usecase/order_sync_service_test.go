package usecase

import (
	"context"
	"regexp"
	"testing"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"
	salesDomain "am-erp-go/internal/module/sales/domain"
)

func TestSyncOrdersSuccessUpdatesState(t *testing.T) {
	now := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)
	fakeRepo := &fakeSyncRepository{}
	fakeImporter := &fakeOrderImportGateway{
		skuBySellerSKU: map[string]uint64{
			"SKU-A": 101,
		},
	}
	fakeProvider := &fakeOrdersProvider{
		code: "AMAZON",
		orders: []integrationDomain.ExternalOrder{
			{
				OrderID:       "111-2222222-3333333",
				MarketplaceID: "ATVPDKIKX0DER",
				PurchaseAt:    now.Add(-2 * time.Hour),
				LastUpdatedAt: now.Add(-1 * time.Hour),
				Currency:      "USD",
				Items: []integrationDomain.ExternalOrderItem{
					{
						OrderItemID: "item-1",
						SellerSKU:   "SKU-A",
						Quantity:    2,
						Amount:      39.98,
					},
				},
			},
		},
	}

	uc := NewOrderSyncService(fakeRepo, fakeImporter, fakeProvider, OrderSyncConfig{
		ProviderCode:       "AMAZON",
		Channel:            "SP_API_SELF",
		SourceType:         "AMAZON_API",
		DefaultCurrency:    "USD",
		LookbackMinutes:    10,
		InitialLookbackDay: 7,
	}, func() time.Time { return now }, nil)

	result, err := uc.SyncOrders(context.Background(), integrationDomain.SyncTriggerManual, nil)
	if err != nil {
		t.Fatalf("SyncOrders returned error: %v", err)
	}
	if result.ImportedLines != 1 {
		t.Fatalf("expected imported lines 1, got %d", result.ImportedLines)
	}
	if result.ErrorLines != 0 {
		t.Fatalf("expected 0 error lines, got %d", result.ErrorLines)
	}
	if len(fakeImporter.upsertedLines) != 1 {
		t.Fatalf("expected one upserted line, got %d", len(fakeImporter.upsertedLines))
	}
	if fakeRepo.savedState == nil || fakeRepo.savedState.LastUpdatedAfter == nil {
		t.Fatalf("expected sync state to be saved")
	}
	if fakeRepo.savedState.LastUpdatedAfter.UTC() != now.Add(-1*time.Hour).UTC() {
		t.Fatalf("unexpected last_updated_after: %s", fakeRepo.savedState.LastUpdatedAfter.UTC())
	}
}

func TestBuildBatchNoUsesReadableNumericFormat(t *testing.T) {
	value := buildBatchNo(time.Date(2026, 3, 10, 13, 8, 0, 0, time.UTC), "amazon")
	matched, err := regexp.MatchString(`^AMAZON202603101308\d{4}$`, value)
	if err != nil {
		t.Fatalf("regexp error: %v", err)
	}
	if !matched {
		t.Fatalf("unexpected sync batch number: %s", value)
	}
}

func TestSyncOrdersUnknownSkuRecordsRowError(t *testing.T) {
	now := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)
	fakeRepo := &fakeSyncRepository{}
	fakeImporter := &fakeOrderImportGateway{
		skuBySellerSKU: map[string]uint64{},
	}
	fakeProvider := &fakeOrdersProvider{
		code: "AMAZON",
		orders: []integrationDomain.ExternalOrder{
			{
				OrderID:       "111-4444444-5555555",
				MarketplaceID: "ATVPDKIKX0DER",
				PurchaseAt:    now.Add(-1 * time.Hour),
				LastUpdatedAt: now.Add(-30 * time.Minute),
				Currency:      "USD",
				Items: []integrationDomain.ExternalOrderItem{
					{
						OrderItemID: "item-2",
						SellerSKU:   "SKU-NOT-EXIST",
						Quantity:    1,
						Amount:      19.99,
					},
				},
			},
		},
	}

	uc := NewOrderSyncService(fakeRepo, fakeImporter, fakeProvider, OrderSyncConfig{
		ProviderCode:       "AMAZON",
		Channel:            "SP_API_SELF",
		SourceType:         "AMAZON_API",
		DefaultCurrency:    "USD",
		LookbackMinutes:    10,
		InitialLookbackDay: 7,
	}, func() time.Time { return now }, nil)

	result, err := uc.SyncOrders(context.Background(), integrationDomain.SyncTriggerManual, nil)
	if err != nil {
		t.Fatalf("SyncOrders returned error: %v", err)
	}
	if result.ImportedLines != 0 {
		t.Fatalf("expected imported lines 0, got %d", result.ImportedLines)
	}
	if result.ErrorLines != 1 {
		t.Fatalf("expected one error line, got %d", result.ErrorLines)
	}
	if len(fakeImporter.rowErrors) != 1 {
		t.Fatalf("expected one stored row error, got %d", len(fakeImporter.rowErrors))
	}
	if len(fakeImporter.upsertedLines) != 0 {
		t.Fatalf("expected no upserted line, got %d", len(fakeImporter.upsertedLines))
	}
}

func TestNewOrderSyncServiceUsesConfigCurrencyWhenDefaultCurrencyMissing(t *testing.T) {
	uc := NewOrderSyncService(
		&fakeSyncRepository{},
		&fakeOrderImportGateway{},
		&fakeOrdersProvider{code: "AMAZON"},
		OrderSyncConfig{
			ProviderCode: "AMAZON",
			Channel:      "SP_API_SELF",
		},
		nil,
		stubOrderSyncCurrencyProvider{currency: "EUR"},
	)

	if uc.cfg.DefaultCurrency != "EUR" {
		t.Fatalf("expected config currency EUR, got %s", uc.cfg.DefaultCurrency)
	}
}

func TestSyncOrdersUsesSKUMappingBeforeFallbackLookup(t *testing.T) {
	now := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)
	fakeRepo := &fakeSyncRepository{}
	fakeImporter := &fakeOrderImportGateway{
		skuBySellerSKU: map[string]uint64{
			"SKU-A": 101,
		},
	}
	fakeProvider := &fakeOrdersProvider{
		code: "AMAZON",
		orders: []integrationDomain.ExternalOrder{
			{
				OrderID:       "111-2222222-3333333",
				MarketplaceID: "ATVPDKIKX0DER",
				PurchaseAt:    now.Add(-2 * time.Hour),
				LastUpdatedAt: now.Add(-1 * time.Hour),
				Currency:      "USD",
				Items: []integrationDomain.ExternalOrderItem{
					{
						OrderItemID: "item-1",
						SellerSKU:   "SKU-A",
						Quantity:    2,
						Amount:      39.98,
					},
				},
			},
		},
	}
	resolver := &fakeSKUMappingResolver{
		mappings: map[string]uint64{
			"AMAZON|US|SKU-A": 999,
		},
	}

	uc := NewOrderSyncService(fakeRepo, fakeImporter, fakeProvider, OrderSyncConfig{
		ProviderCode:       "AMAZON",
		Channel:            "SP_API_SELF",
		SourceType:         "AMAZON_API",
		DefaultCurrency:    "USD",
		LookbackMinutes:    10,
		InitialLookbackDay: 7,
	}, func() time.Time { return now }, nil)
	uc.BindSKUMappingResolver(resolver)

	_, err := uc.SyncOrders(context.Background(), integrationDomain.SyncTriggerManual, nil)
	if err != nil {
		t.Fatalf("SyncOrders returned error: %v", err)
	}
	if len(fakeImporter.upsertedProductIDs) != 1 || fakeImporter.upsertedProductIDs[0] != 999 {
		t.Fatalf("expected mapped product id 999, got %v", fakeImporter.upsertedProductIDs)
	}
	if fakeImporter.resolveCalls != 0 {
		t.Fatalf("expected fallback lookup not called when mapping hit, got %d", fakeImporter.resolveCalls)
	}
	if resolver.calls != 1 {
		t.Fatalf("expected resolver called once, got %d", resolver.calls)
	}
}

type stubOrderSyncCurrencyProvider struct {
	currency string
}

func (s stubOrderSyncCurrencyProvider) GetDefaultBaseCurrency() string {
	return s.currency
}

type fakeSyncRepository struct {
	state      *integrationDomain.OrderSyncState
	savedState *integrationDomain.OrderSyncState
	createdRun *integrationDomain.OrderSyncRun
	updatedRun *integrationDomain.OrderSyncRun
}

func (f *fakeSyncRepository) GetState(provider string, channel string) (*integrationDomain.OrderSyncState, error) {
	_ = provider
	_ = channel
	return f.state, nil
}

func (f *fakeSyncRepository) SaveState(state *integrationDomain.OrderSyncState) error {
	f.savedState = state
	f.state = state
	return nil
}

func (f *fakeSyncRepository) CreateRun(run *integrationDomain.OrderSyncRun) error {
	run.ID = 1
	f.createdRun = run
	return nil
}

func (f *fakeSyncRepository) UpdateRun(run *integrationDomain.OrderSyncRun) error {
	f.updatedRun = run
	return nil
}

func (f *fakeSyncRepository) ListRuns(provider string, channel string, params *integrationDomain.ListRunsParams) ([]integrationDomain.OrderSyncRun, int64, error) {
	_ = provider
	_ = channel
	_ = params
	return []integrationDomain.OrderSyncRun{}, 0, nil
}

type fakeOrderImportGateway struct {
	skuBySellerSKU     map[string]uint64
	upsertedLines      []salesDomain.ImportOrderLine
	upsertedProductIDs []uint64
	rowErrors          []salesDomain.ReportImportRowError
	createdImports     []*salesDomain.ReportImport
	resolveCalls       int
}

func (f *fakeOrderImportGateway) CreateImport(batch *salesDomain.ReportImport) error {
	batch.ID = uint64(len(f.createdImports) + 1)
	f.createdImports = append(f.createdImports, batch)
	return nil
}

func (f *fakeOrderImportGateway) UpdateImport(batch *salesDomain.ReportImport) error {
	return nil
}

func (f *fakeOrderImportGateway) InsertImportRowErrors(rows []salesDomain.ReportImportRowError) error {
	f.rowErrors = append(f.rowErrors, rows...)
	return nil
}

func (f *fakeOrderImportGateway) ResolveProductIDBySellerSKU(sellerSKU string, marketplace string) (uint64, error) {
	_ = marketplace
	f.resolveCalls++
	return f.skuBySellerSKU[sellerSKU], nil
}

func (f *fakeOrderImportGateway) UpsertImportedOrderLine(
	line *salesDomain.ImportOrderLine,
	productID uint64,
	batchNo string,
	operatorID *uint64,
) error {
	_ = batchNo
	_ = operatorID
	f.upsertedLines = append(f.upsertedLines, *line)
	f.upsertedProductIDs = append(f.upsertedProductIDs, productID)
	return nil
}

type fakeOrdersProvider struct {
	code   string
	orders []integrationDomain.ExternalOrder
}

type fakeSKUMappingResolver struct {
	mappings map[string]uint64
	calls    int
}

func (f *fakeSKUMappingResolver) ResolveMappedProductID(providerCode string, sellerSKU string, marketplace string) (uint64, error) {
	f.calls++
	key := providerCode + "|" + marketplace + "|" + sellerSKU
	return f.mappings[key], nil
}

func (f *fakeOrdersProvider) Code() string {
	return f.code
}

func (f *fakeOrdersProvider) ListOrders(ctx context.Context, req integrationDomain.ExternalListOrdersRequest) ([]integrationDomain.ExternalOrder, error) {
	_ = ctx
	_ = req
	return f.orders, nil
}
