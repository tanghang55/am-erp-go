package usecase

import (
	"testing"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

func TestCreateSKUMappingNormalizesFields(t *testing.T) {
	now := time.Date(2026, 3, 24, 12, 0, 0, 0, time.UTC)
	repo := &fakeSKUMappingRepo{}
	uc := NewSKUMappingUsecase(repo, func() time.Time { return now })

	operatorID := uint64(7)
	item, err := uc.Create(&CreateSKUMappingInput{
		ProviderCode: " amazon_us ",
		Marketplace:  " us ",
		SellerSKU:    "  sku-001 ",
		ProductID:    1001,
		Remark:       "  首次绑定 ",
		OperatorID:   &operatorID,
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if item.ProviderCode != "AMAZON_US" {
		t.Fatalf("expected provider_code normalized, got %q", item.ProviderCode)
	}
	if item.Marketplace != "US" {
		t.Fatalf("expected marketplace normalized, got %q", item.Marketplace)
	}
	if item.SellerSKU != "sku-001" {
		t.Fatalf("expected seller_sku trimmed, got %q", item.SellerSKU)
	}
	if item.Status != integrationDomain.SKUMappingStatusActive {
		t.Fatalf("expected default ACTIVE status, got %s", item.Status)
	}
	if repo.created == nil || repo.created.CreatedBy == nil || *repo.created.CreatedBy != operatorID {
		t.Fatalf("expected created_by to be set")
	}
}

func TestUpdateSKUMappingOnlyChangesEditableFields(t *testing.T) {
	now := time.Date(2026, 3, 24, 12, 0, 0, 0, time.UTC)
	repo := &fakeSKUMappingRepo{
		byID: &integrationDomain.IntegrationSKUMapping{
			ID:           1,
			ProviderCode: "AMAZON_US",
			Marketplace:  "US",
			SellerSKU:    "sku-001",
			ProductID:    1001,
			Status:       integrationDomain.SKUMappingStatusActive,
		},
	}
	uc := NewSKUMappingUsecase(repo, func() time.Time { return now })
	operatorID := uint64(9)

	item, err := uc.Update(1, &UpdateSKUMappingInput{
		ProductID:  2002,
		Status:     integrationDomain.SKUMappingStatusDisabled,
		Remark:     "调整",
		OperatorID: &operatorID,
	})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if item.ProductID != 2002 {
		t.Fatalf("expected product_id=2002, got %d", item.ProductID)
	}
	if item.Status != integrationDomain.SKUMappingStatusDisabled {
		t.Fatalf("expected status DISABLED, got %s", item.Status)
	}
	if item.ProviderCode != "AMAZON_US" || item.SellerSKU != "sku-001" {
		t.Fatalf("expected immutable identity fields unchanged, got %#v", item)
	}
	if repo.updated == nil || repo.updated.UpdatedBy == nil || *repo.updated.UpdatedBy != operatorID {
		t.Fatalf("expected updated_by to be set")
	}
}

func TestResolveMappedProductIDOnlyReturnsActive(t *testing.T) {
	repo := &fakeSKUMappingRepo{
		resolvedID: 1002,
	}
	uc := NewSKUMappingUsecase(repo, nil)

	productID, err := uc.ResolveMappedProductID("amazon_us", "sku-001", "us")
	if err != nil {
		t.Fatalf("ResolveMappedProductID error: %v", err)
	}
	if productID != 1002 {
		t.Fatalf("expected product_id=1002, got %d", productID)
	}
	if repo.resolveProvider != "AMAZON_US" || repo.resolveMarketplace != "US" || repo.resolveSellerSKU != "sku-001" {
		t.Fatalf("unexpected normalized resolve arguments: %+v", repo)
	}
}

type fakeSKUMappingRepo struct {
	created            *integrationDomain.IntegrationSKUMapping
	updated            *integrationDomain.IntegrationSKUMapping
	byID               *integrationDomain.IntegrationSKUMapping
	resolvedID         uint64
	resolveProvider    string
	resolveMarketplace string
	resolveSellerSKU   string
}

func (f *fakeSKUMappingRepo) Create(item *integrationDomain.IntegrationSKUMapping) error {
	item.ID = 1
	f.created = item
	f.byID = item
	return nil
}

func (f *fakeSKUMappingRepo) Update(item *integrationDomain.IntegrationSKUMapping) error {
	f.updated = item
	f.byID = item
	return nil
}

func (f *fakeSKUMappingRepo) GetByID(id uint64) (*integrationDomain.IntegrationSKUMapping, error) {
	if f.byID != nil && f.byID.ID == id {
		return f.byID, nil
	}
	return nil, nil
}

func (f *fakeSKUMappingRepo) GetByUnique(providerCode string, marketplace string, sellerSKU string) (*integrationDomain.IntegrationSKUMapping, error) {
	if f.byID == nil {
		return nil, nil
	}
	if f.byID.ProviderCode == providerCode && f.byID.Marketplace == marketplace && f.byID.SellerSKU == sellerSKU {
		return f.byID, nil
	}
	return nil, nil
}

func (f *fakeSKUMappingRepo) List(params *integrationDomain.SKUMappingListParams) ([]integrationDomain.IntegrationSKUMapping, int64, error) {
	_ = params
	return []integrationDomain.IntegrationSKUMapping{}, 0, nil
}

func (f *fakeSKUMappingRepo) ResolveActiveProductID(providerCode string, marketplace string, sellerSKU string) (uint64, error) {
	f.resolveProvider = providerCode
	f.resolveMarketplace = marketplace
	f.resolveSellerSKU = sellerSKU
	return f.resolvedID, nil
}
