package usecase

import "am-erp-go/internal/module/finance/domain"

type ProductCostUsecase struct {
	repo domain.ProductCostRepository
}

func NewProductCostUsecase(repo domain.ProductCostRepository) *ProductCostUsecase {
	return &ProductCostUsecase{repo: repo}
}

func (uc *ProductCostUsecase) ListLedger(params *domain.ProductCostLedgerListParams) ([]domain.ProductCostLedgerItem, int64, error) {
	if params == nil {
		params = &domain.ProductCostLedgerListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	items, total, err := uc.repo.ListLedger(params)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].BaseCurrency = ensureBaseCurrency(items[i].BaseCurrency)
	}
	return items, total, nil
}

func (uc *ProductCostUsecase) GetSummary(params *domain.ProductCostLedgerListParams) (*domain.ProductCostSummary, error) {
	if params == nil {
		params = &domain.ProductCostLedgerListParams{}
	}
	summary, err := uc.repo.GetSummary(params)
	if err != nil {
		return nil, err
	}
	if summary != nil {
		summary.BaseCurrency = ensureBaseCurrency(summary.BaseCurrency)
	}
	return summary, nil
}
