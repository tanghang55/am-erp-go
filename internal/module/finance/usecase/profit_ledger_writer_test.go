package usecase

import (
	"testing"
	"time"

	"am-erp-go/internal/module/finance/domain"
	salesUsecase "am-erp-go/internal/module/sales/usecase"
)

type profitLedgerRepoStub struct {
	created      []*domain.ProfitLedger
	createdBatch []domain.ProfitLedger
	err          error
}

func (s *profitLedgerRepoStub) Create(entry *domain.ProfitLedger) error {
	if entry != nil {
		s.created = append(s.created, entry)
	}
	return s.err
}

func (s *profitLedgerRepoStub) CreateBatch(entries []domain.ProfitLedger) error {
	s.createdBatch = append(s.createdBatch, entries...)
	return s.err
}

func (s *profitLedgerRepoStub) AggregateDaily(_ time.Time, _ *string) ([]domain.ProfitLedgerDailyAgg, error) {
	return nil, nil
}

func TestProfitLedgerWriter_RecordSalesShipProfit(t *testing.T) {
	repo := &profitLedgerRepoStub{}
	writer := NewProfitLedgerWriter(repo)
	now := time.Now()
	err := writer.RecordSalesShipProfit(&salesUsecase.SalesShipProfitRecordParams{
		SalesOrderID:     11,
		SalesOrderItemID: 21,
		OrderNo:          "SO-11",
		Marketplace:      "US",
		IncomeCurrency:   "USD",
		IncomeAmount:     100,
		COGSAmount:       35.5,
		OccurredAt:       now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.createdBatch) != 2 {
		t.Fatalf("expected 2 profit ledgers, got %d", len(repo.createdBatch))
	}
	if repo.createdBatch[0].LedgerType != domain.ProfitLedgerTypeIncome {
		t.Fatalf("expected first ledger INCOME, got %s", repo.createdBatch[0].LedgerType)
	}
	if repo.createdBatch[1].LedgerType != domain.ProfitLedgerTypeCOGS {
		t.Fatalf("expected second ledger COGS, got %s", repo.createdBatch[1].LedgerType)
	}
}

func TestProfitLedgerWriter_RecordSalesReturnProfit(t *testing.T) {
	repo := &profitLedgerRepoStub{}
	writer := NewProfitLedgerWriter(repo)
	now := time.Now()

	err := writer.RecordSalesReturnProfit(&salesUsecase.SalesReturnProfitRecordParams{
		SalesOrderID:     11,
		SalesOrderItemID: 21,
		OrderNo:          "SO-11",
		Marketplace:      "US",
		IncomeCurrency:   "USD",
		IncomeAmount:     20,
		COGSAmount:       8.5,
		OccurredAt:       now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.createdBatch) != 2 {
		t.Fatalf("expected 2 reversal profit ledgers, got %d", len(repo.createdBatch))
	}
	if repo.createdBatch[0].ReversalOfID == nil || repo.createdBatch[1].ReversalOfID == nil {
		t.Fatalf("expected reversal_of_id on return ledgers, got %+v", repo.createdBatch)
	}
	if repo.createdBatch[0].Category != "SALES_RETURN_REVENUE" {
		t.Fatalf("unexpected return income category: %s", repo.createdBatch[0].Category)
	}
	if repo.createdBatch[1].Category != "SALES_RETURN_COGS" {
		t.Fatalf("unexpected return cogs category: %s", repo.createdBatch[1].Category)
	}
}

func TestProfitLedgerWriter_RecordCashLedger(t *testing.T) {
	repo := &profitLedgerRepoStub{}
	writer := NewProfitLedgerWriter(repo)
	refType := "SALES_ORDER"
	refID := uint64(8)
	marketplace := "US"
	desc := "ad fee"
	entry := &domain.CashLedger{
		ID:               3,
		LedgerType:       domain.LedgerTypeExpense,
		Category:         "AD_FEE",
		OriginalCurrency: "USD",
		OriginalAmount:   20,
		BaseCurrency:     "USD",
		FxRate:           1,
		BaseAmount:       20,
		ReferenceType:    &refType,
		ReferenceID:      &refID,
		Marketplace:      &marketplace,
		Description:      &desc,
		OccurredAt:       time.Now(),
		CreatedBy:        9,
	}
	if err := writer.RecordCashLedger(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected 1 profit ledger, got %d", len(repo.created))
	}
	if repo.created[0].LedgerType != domain.ProfitLedgerTypeOrderExpense {
		t.Fatalf("expected ORDER_EXPENSE, got %s", repo.created[0].LedgerType)
	}
}
