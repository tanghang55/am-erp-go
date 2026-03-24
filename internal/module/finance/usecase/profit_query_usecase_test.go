package usecase

import (
	"testing"

	"am-erp-go/internal/module/finance/domain"
)

type profitQueryRepoStub struct {
	items  []domain.OrderProfitSummary
	detail *domain.OrderProfitDetail
}

func (s *profitQueryRepoStub) ListOrderProfits(params *domain.OrderProfitListParams) ([]domain.OrderProfitSummary, int64, error) {
	return s.items, int64(len(s.items)), nil
}

func (s *profitQueryRepoStub) GetOrderProfitDetail(salesOrderID uint64) (*domain.OrderProfitDetail, error) {
	return s.detail, nil
}

func TestProfitQueryUsecaseFillsMissingBaseCurrency(t *testing.T) {
	SetDefaultBaseCurrencyResolver(func() string { return "EUR" })
	t.Cleanup(func() {
		SetDefaultBaseCurrencyResolver(nil)
	})

	repo := &profitQueryRepoStub{
		items: []domain.OrderProfitSummary{
			{SalesOrderID: 1, BaseCurrency: ""},
		},
		detail: &domain.OrderProfitDetail{
			Summary: domain.OrderProfitSummary{SalesOrderID: 1, BaseCurrency: ""},
			Expenses: []domain.OrderProfitExpense{
				{ID: 1, BaseCurrency: ""},
			},
		},
	}
	uc := NewProfitQueryUsecase(repo)

	items, _, err := uc.ListOrderProfits(&domain.OrderProfitListParams{})
	if err != nil {
		t.Fatalf("unexpected list error: %v", err)
	}
	if items[0].BaseCurrency != "EUR" {
		t.Fatalf("expected EUR, got %s", items[0].BaseCurrency)
	}

	detail, err := uc.GetOrderProfitDetail(1)
	if err != nil {
		t.Fatalf("unexpected detail error: %v", err)
	}
	if detail.Summary.BaseCurrency != "EUR" {
		t.Fatalf("expected summary EUR, got %s", detail.Summary.BaseCurrency)
	}
	if detail.Expenses[0].BaseCurrency != "EUR" {
		t.Fatalf("expected expense EUR, got %s", detail.Expenses[0].BaseCurrency)
	}
}
