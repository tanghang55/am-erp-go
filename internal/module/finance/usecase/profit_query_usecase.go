package usecase

import (
	"am-erp-go/internal/module/finance/domain"
)

type ProfitQueryUsecase struct {
	repo domain.ProfitQueryRepository
}

func NewProfitQueryUsecase(repo domain.ProfitQueryRepository) *ProfitQueryUsecase {
	return &ProfitQueryUsecase{repo: repo}
}

func (uc *ProfitQueryUsecase) ListOrderProfits(params *domain.OrderProfitListParams) ([]domain.OrderProfitSummary, int64, error) {
	if params == nil {
		params = &domain.OrderProfitListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	items, total, err := uc.repo.ListOrderProfits(params)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].BaseCurrency = ensureBaseCurrency(items[i].BaseCurrency)
	}
	return items, total, nil
}

func (uc *ProfitQueryUsecase) GetOrderProfitDetail(salesOrderID uint64) (*domain.OrderProfitDetail, error) {
	detail, err := uc.repo.GetOrderProfitDetail(salesOrderID)
	if err != nil {
		return nil, err
	}
	if detail != nil {
		detail.Summary.BaseCurrency = ensureBaseCurrency(detail.Summary.BaseCurrency)
		for i := range detail.Expenses {
			detail.Expenses[i].BaseCurrency = ensureBaseCurrency(detail.Expenses[i].BaseCurrency)
		}
	}
	return detail, nil
}
