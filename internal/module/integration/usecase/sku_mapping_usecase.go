package usecase

import (
	"fmt"
	"strings"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

type CreateSKUMappingInput struct {
	ProviderCode string
	Marketplace  string
	SellerSKU    string
	ProductID    uint64
	Status       integrationDomain.SKUMappingStatus
	Remark       string
	OperatorID   *uint64
}

type UpdateSKUMappingInput struct {
	ProductID  uint64
	Status     integrationDomain.SKUMappingStatus
	Remark     string
	OperatorID *uint64
}

type SKUMappingUsecase struct {
	repo  integrationDomain.SKUMappingRepository
	nowFn func() time.Time
}

func NewSKUMappingUsecase(repo integrationDomain.SKUMappingRepository, nowFn func() time.Time) *SKUMappingUsecase {
	if nowFn == nil {
		nowFn = time.Now
	}
	return &SKUMappingUsecase{
		repo:  repo,
		nowFn: nowFn,
	}
}

func (u *SKUMappingUsecase) List(params *integrationDomain.SKUMappingListParams) ([]integrationDomain.IntegrationSKUMapping, int64, error) {
	if u == nil || u.repo == nil {
		return nil, 0, fmt.Errorf("sku mapping usecase not configured")
	}
	if params == nil {
		params = &integrationDomain.SKUMappingListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	params.ProviderCode = normalizeProviderCode(params.ProviderCode)
	params.Marketplace = normalizeMarketplaceCode(params.Marketplace)
	params.Status = strings.ToUpper(strings.TrimSpace(params.Status))
	params.Keyword = strings.TrimSpace(params.Keyword)
	return u.repo.List(params)
}

func (u *SKUMappingUsecase) GetByID(id uint64) (*integrationDomain.IntegrationSKUMapping, error) {
	if u == nil || u.repo == nil {
		return nil, fmt.Errorf("sku mapping usecase not configured")
	}
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	return u.repo.GetByID(id)
}

func (u *SKUMappingUsecase) Create(input *CreateSKUMappingInput) (*integrationDomain.IntegrationSKUMapping, error) {
	if u == nil || u.repo == nil {
		return nil, fmt.Errorf("sku mapping usecase not configured")
	}
	if input == nil {
		return nil, fmt.Errorf("input is required")
	}
	providerCode := normalizeProviderCode(input.ProviderCode)
	marketplace := normalizeMarketplaceCode(input.Marketplace)
	sellerSKU := strings.TrimSpace(input.SellerSKU)
	if providerCode == "" {
		return nil, fmt.Errorf("provider_code is required")
	}
	if marketplace == "" {
		return nil, fmt.Errorf("marketplace is required")
	}
	if sellerSKU == "" {
		return nil, fmt.Errorf("seller_sku is required")
	}
	if input.ProductID == 0 {
		return nil, fmt.Errorf("product_id is required")
	}

	status := input.Status
	if status == "" {
		status = integrationDomain.SKUMappingStatusActive
	}
	if !isValidSKUMappingStatus(status) {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	exists, err := u.repo.GetByUnique(providerCode, marketplace, sellerSKU)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return nil, fmt.Errorf("sku mapping already exists")
	}

	item := &integrationDomain.IntegrationSKUMapping{
		ProviderCode: providerCode,
		Marketplace:  marketplace,
		SellerSKU:    sellerSKU,
		ProductID:    input.ProductID,
		Status:       status,
		Remark:       strPtrOrNil(input.Remark),
		CreatedBy:    input.OperatorID,
		UpdatedBy:    input.OperatorID,
	}
	if err := u.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (u *SKUMappingUsecase) Update(id uint64, input *UpdateSKUMappingInput) (*integrationDomain.IntegrationSKUMapping, error) {
	if u == nil || u.repo == nil {
		return nil, fmt.Errorf("sku mapping usecase not configured")
	}
	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	if input == nil {
		return nil, fmt.Errorf("input is required")
	}
	item, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("sku mapping not found")
	}
	if input.ProductID > 0 {
		item.ProductID = input.ProductID
	}
	if input.Status != "" {
		if !isValidSKUMappingStatus(input.Status) {
			return nil, fmt.Errorf("invalid status: %s", input.Status)
		}
		item.Status = input.Status
	}
	item.Remark = strPtrOrNil(input.Remark)
	item.UpdatedBy = input.OperatorID
	if err := u.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (u *SKUMappingUsecase) ResolveMappedProductID(providerCode string, sellerSKU string, marketplace string) (uint64, error) {
	if u == nil || u.repo == nil {
		return 0, nil
	}
	normProvider := normalizeProviderCode(providerCode)
	normMarketplace := normalizeMarketplaceCode(marketplace)
	normSellerSKU := strings.TrimSpace(sellerSKU)
	if normProvider == "" || normMarketplace == "" || normSellerSKU == "" {
		return 0, nil
	}
	return u.repo.ResolveActiveProductID(normProvider, normMarketplace, normSellerSKU)
}

func normalizeMarketplaceCode(v string) string {
	return strings.ToUpper(strings.TrimSpace(v))
}

func isValidSKUMappingStatus(status integrationDomain.SKUMappingStatus) bool {
	return status == integrationDomain.SKUMappingStatusActive || status == integrationDomain.SKUMappingStatusDisabled
}
