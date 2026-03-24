package usecase

import (
	"testing"

	"am-erp-go/internal/module/finance/domain"
)

type productCostRepoStub struct {
	listItems []domain.ProductCostLedgerItem
	listTotal int64
	summary   *domain.ProductCostSummary
}

func (s *productCostRepoStub) ListLedger(params *domain.ProductCostLedgerListParams) ([]domain.ProductCostLedgerItem, int64, error) {
	return s.listItems, s.listTotal, nil
}

func (s *productCostRepoStub) GetSummary(params *domain.ProductCostLedgerListParams) (*domain.ProductCostSummary, error) {
	if s.summary != nil {
		return s.summary, nil
	}
	return &domain.ProductCostSummary{BaseCurrency: "USD"}, nil
}

func TestProductCostUsecaseListDefaultPaging(t *testing.T) {
	repo := &productCostRepoStub{
		listItems: []domain.ProductCostLedgerItem{{ProductID: 1}},
		listTotal: 1,
	}
	uc := NewProductCostUsecase(repo)

	items, total, err := uc.ListLedger(&domain.ProductCostLedgerListParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || total != 1 {
		t.Fatalf("unexpected result: items=%d total=%d", len(items), total)
	}
}
