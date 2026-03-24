package usecase

import (
	"errors"
	"testing"
	"time"

	"am-erp-go/internal/module/finance/domain"
)

type cashLedgerRepoStub struct {
	getByIDFn              func(id uint64) (*domain.CashLedger, error)
	createFn               func(entry *domain.CashLedger) error
	updateFn               func(entry *domain.CashLedger) error
	deleteFn               func(id uint64) error
	markReversedFn         func(id uint64, reversedAt time.Time) error
	listFn                 func(params *domain.CashLedgerListParams) ([]domain.CashLedger, int64, error)
	getSummaryFn           func(params *domain.CashLedgerListParams) (*domain.CashLedgerSummary, error)
	getSummaryByCategoryFn func(params *domain.CashLedgerListParams) ([]domain.CategorySummaryItem, error)
}

type cashLedgerProfitRecorderStub struct {
	entries []*domain.CashLedger
	err     error
}

func (s *cashLedgerProfitRecorderStub) RecordCashLedger(entry *domain.CashLedger) error {
	if entry != nil {
		s.entries = append(s.entries, entry)
	}
	return s.err
}

func (s *cashLedgerRepoStub) List(params *domain.CashLedgerListParams) ([]domain.CashLedger, int64, error) {
	if s.listFn != nil {
		return s.listFn(params)
	}
	return []domain.CashLedger{}, 0, nil
}

func (s *cashLedgerRepoStub) GetByID(id uint64) (*domain.CashLedger, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(id)
	}
	return nil, nil
}

func (s *cashLedgerRepoStub) Create(entry *domain.CashLedger) error {
	if s.createFn != nil {
		return s.createFn(entry)
	}
	return nil
}

func (s *cashLedgerRepoStub) Update(entry *domain.CashLedger) error {
	if s.updateFn != nil {
		return s.updateFn(entry)
	}
	return nil
}

func (s *cashLedgerRepoStub) Delete(id uint64) error {
	if s.deleteFn != nil {
		return s.deleteFn(id)
	}
	return nil
}

func (s *cashLedgerRepoStub) MarkReversed(id uint64, reversedAt time.Time) error {
	if s.markReversedFn != nil {
		return s.markReversedFn(id, reversedAt)
	}
	return nil
}

func (s *cashLedgerRepoStub) GetSummary(params *domain.CashLedgerListParams) (*domain.CashLedgerSummary, error) {
	if s.getSummaryFn != nil {
		return s.getSummaryFn(params)
	}
	return &domain.CashLedgerSummary{}, nil
}

func (s *cashLedgerRepoStub) GetSummaryByCategory(params *domain.CashLedgerListParams) ([]domain.CategorySummaryItem, error) {
	if s.getSummaryByCategoryFn != nil {
		return s.getSummaryByCategoryFn(params)
	}
	return []domain.CategorySummaryItem{}, nil
}

func TestCashLedgerUsecase_UpdateImmutable(t *testing.T) {
	uc := NewCashLedgerUsecase(&cashLedgerRepoStub{})
	_, err := uc.Update(1, &UpdateCashLedgerInput{
		Amount: float64Ptr(10),
	})
	if !errors.Is(err, ErrCashLedgerImmutable) {
		t.Fatalf("expected ErrCashLedgerImmutable, got %v", err)
	}
}

func TestCashLedgerUsecase_DeleteImmutable(t *testing.T) {
	uc := NewCashLedgerUsecase(&cashLedgerRepoStub{})
	err := uc.Delete(1)
	if !errors.Is(err, ErrCashLedgerImmutable) {
		t.Fatalf("expected ErrCashLedgerImmutable, got %v", err)
	}
}

func TestCashLedgerUsecase_ReverseSuccess(t *testing.T) {
	now := time.Now()
	var created *domain.CashLedger
	markedID := uint64(0)

	repo := &cashLedgerRepoStub{
		getByIDFn: func(id uint64) (*domain.CashLedger, error) {
			if id != 10 {
				t.Fatalf("unexpected get id: %d", id)
			}
			refType := "SALES_ORDER"
			refID := uint64(99)
			desc := "origin"
			return &domain.CashLedger{
				ID:            10,
				TraceID:       "TRACE-10",
				LedgerType:    domain.LedgerTypeIncome,
				Status:        domain.CashLedgerStatusNormal,
				Category:      "SALES_REVENUE",
				Amount:        120,
				Currency:      "USD",
				ReferenceType: &refType,
				ReferenceID:   &refID,
				Description:   &desc,
				OccurredAt:    now,
				CreatedBy:     8,
			}, nil
		},
		createFn: func(entry *domain.CashLedger) error {
			created = entry
			return nil
		},
		markReversedFn: func(id uint64, _ time.Time) error {
			markedID = id
			return nil
		},
	}

	uc := NewCashLedgerUsecase(repo)
	reversed, err := uc.Reverse(10, 8, "manual reverse")
	if err != nil {
		t.Fatalf("reverse returned err: %v", err)
	}
	if reversed == nil {
		t.Fatalf("expected reversed entry")
	}
	if created == nil {
		t.Fatalf("expected create to be called")
	}
	if markedID != 10 {
		t.Fatalf("expected mark reversed id=10, got %d", markedID)
	}
	if reversed.LedgerType != domain.LedgerTypeExpense {
		t.Fatalf("expected reversed ledger_type expense, got %s", reversed.LedgerType)
	}
	if reversed.Amount != 120 {
		t.Fatalf("expected reversed amount 120, got %v", reversed.Amount)
	}
	if reversed.ReversalOfID == nil || *reversed.ReversalOfID != 10 {
		t.Fatalf("expected reversal_of_id=10, got %+v", reversed.ReversalOfID)
	}
}

func TestCashLedgerUsecase_CreateWithFxSnapshot(t *testing.T) {
	SetFXRateResolver(func(baseCurrency, originalCurrency string, occurredAt time.Time) (*FXRateSnapshot, error) {
		return &FXRateSnapshot{
			Rate:        0.14,
			Source:      string(domain.ExchangeRateSourceManual),
			Version:     fxVersionManual,
			EffectiveAt: occurredAt,
		}, nil
	})
	t.Cleanup(func() {
		SetFXRateResolver(nil)
	})

	occurredAt := time.Date(2026, 3, 3, 9, 30, 0, 0, time.Local)
	var created *domain.CashLedger
	repo := &cashLedgerRepoStub{
		createFn: func(entry *domain.CashLedger) error {
			created = entry
			return nil
		},
	}
	uc := NewCashLedgerUsecase(repo)

	entry, err := uc.Create(&CreateCashLedgerInput{
		LedgerType: domain.LedgerTypeExpense,
		Category:   "SHIPPING_FEE",
		Amount:     100,
		Currency:   "CNY",
		OccurredAt: &occurredAt,
		CreatedBy:  8,
	})
	if err != nil {
		t.Fatalf("create returned err: %v", err)
	}
	if entry == nil || created == nil {
		t.Fatalf("expected created entry")
	}
	if entry.OriginalCurrency != "CNY" {
		t.Fatalf("expected original currency CNY, got %s", entry.OriginalCurrency)
	}
	if entry.BaseCurrency != "USD" {
		t.Fatalf("expected base currency USD, got %s", entry.BaseCurrency)
	}
	if entry.FxRate <= 0 {
		t.Fatalf("expected fx rate > 0, got %v", entry.FxRate)
	}
	if entry.BaseAmount <= 0 {
		t.Fatalf("expected base amount > 0, got %v", entry.BaseAmount)
	}
	if !entry.FxTime.Equal(occurredAt) {
		t.Fatalf("expected fx_time equals occurred_at, got %s vs %s", entry.FxTime, occurredAt)
	}
}

func TestCashLedgerUsecase_CreateRejectsUnknownCurrency(t *testing.T) {
	uc := NewCashLedgerUsecase(&cashLedgerRepoStub{})
	_, err := uc.Create(&CreateCashLedgerInput{
		LedgerType: domain.LedgerTypeIncome,
		Category:   "SALES_REVENUE",
		Amount:     10,
		Currency:   "XYZ",
		CreatedBy:  8,
	})
	if err == nil {
		t.Fatalf("expected err for unknown currency")
	}
}

func TestCashLedgerUsecase_CreateExpenseRecordsProfitLedger(t *testing.T) {
	repo := &cashLedgerRepoStub{
		createFn: func(entry *domain.CashLedger) error { return nil },
	}
	recorder := &cashLedgerProfitRecorderStub{}
	uc := NewCashLedgerUsecase(repo)
	uc.BindProfitLedgerRecorder(recorder)

	_, err := uc.Create(&CreateCashLedgerInput{
		LedgerType: domain.LedgerTypeExpense,
		Category:   "AD_FEE",
		Amount:     45,
		Currency:   "USD",
		CreatedBy:  7,
	})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if len(recorder.entries) != 1 {
		t.Fatalf("expected 1 profit record call, got %d", len(recorder.entries))
	}
	if recorder.entries[0].Category != "AD_FEE" {
		t.Fatalf("unexpected recorder category: %s", recorder.entries[0].Category)
	}
}

func float64Ptr(v float64) *float64 {
	return &v
}
