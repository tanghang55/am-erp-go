package usecase

import (
	"context"
	"testing"

	productDomain "am-erp-go/internal/module/product/domain"
)

type stubPackingRequirementRepo struct {
	items []productDomain.ProductPackagingItem
	err   error
}

func (s *stubPackingRequirementRepo) ListByProductID(productID uint64) ([]productDomain.ProductPackagingItem, error) {
	return s.items, s.err
}

func (s *stubPackingRequirementRepo) ReplaceAll(productID uint64, items []productDomain.ProductPackagingItem) error {
	return nil
}

func TestResolvePackingRequirementsReturnsConfiguredPackagingItems(t *testing.T) {
	resolver := NewProductPackingRequirementResolver(&stubPackingRequirementRepo{
		items: []productDomain.ProductPackagingItem{
			{
				PackagingItemID: 11,
				QuantityPerUnit: 2,
				PackagingItem: &productDomain.PackagingItemDetail{
					ItemCode: "BOX-1",
					ItemName: "纸箱",
					Unit:     "PCS",
				},
			},
		},
	})

	requirements, err := resolver.ResolvePackingRequirements(context.Background(), 18)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(requirements) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(requirements))
	}
	if requirements[0].PackagingItemID != 11 || requirements[0].QuantityPerUnit != 2 {
		t.Fatalf("unexpected requirement: %+v", requirements[0])
	}
	if requirements[0].ItemCode != "BOX-1" || requirements[0].ItemName != "纸箱" {
		t.Fatalf("expected packaging detail propagated, got %+v", requirements[0])
	}
}

func TestResolvePackingRequirementsRejectsEmptyConfig(t *testing.T) {
	resolver := NewProductPackingRequirementResolver(&stubPackingRequirementRepo{})

	_, err := resolver.ResolvePackingRequirements(context.Background(), 18)
	if err == nil {
		t.Fatal("expected missing packing materials error")
	}
	if err != ErrPackingMaterialsNotConfigured {
		t.Fatalf("unexpected error: %v", err)
	}
}
