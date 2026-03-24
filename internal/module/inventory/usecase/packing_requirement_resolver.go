package usecase

import (
	"context"

	productDomain "am-erp-go/internal/module/product/domain"
)

type ProductPackingRequirementResolver struct {
	productPackagingRepo productDomain.ProductPackagingRepository
}

func NewProductPackingRequirementResolver(productPackagingRepo productDomain.ProductPackagingRepository) *ProductPackingRequirementResolver {
	return &ProductPackingRequirementResolver{
		productPackagingRepo: productPackagingRepo,
	}
}

func (r *ProductPackingRequirementResolver) ResolvePackingRequirements(ctx context.Context, productID uint64) ([]PackingRequirement, error) {
	_ = ctx
	if r == nil || r.productPackagingRepo == nil {
		return nil, ErrPackingRequirementResolverNotConfigured
	}

	items, err := r.productPackagingRepo.ListByProductID(productID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrPackingMaterialsNotConfigured
	}

	requirements := make([]PackingRequirement, 0, len(items))
	for _, item := range items {
		requirement := PackingRequirement{
			PackagingItemID: item.PackagingItemID,
			QuantityPerUnit: item.QuantityPerUnit,
		}
		if item.PackagingItem != nil {
			requirement.ItemCode = item.PackagingItem.ItemCode
			requirement.ItemName = item.PackagingItem.ItemName
			requirement.Unit = item.PackagingItem.Unit
		}
		requirements = append(requirements, requirement)
	}

	return requirements, nil
}
