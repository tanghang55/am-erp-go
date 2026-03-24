package usecase

import (
	"errors"
	"testing"
	"time"

	"am-erp-go/internal/module/finance/domain"
)

type exchangeRateRepoStub struct {
	listFn          func(params *domain.ExchangeRateListParams) ([]domain.ExchangeRate, int64, error)
	createFn        func(rate *domain.ExchangeRate) error
	getByIDFn       func(id uint64) (*domain.ExchangeRate, error)
	updateStatusFn  func(id uint64, status domain.ExchangeRateStatus, operatorID uint64) error
	findEffectiveFn func(fromCurrency, toCurrency string, occurredAt time.Time) (*domain.ExchangeRate, error)
}

func (s *exchangeRateRepoStub) List(params *domain.ExchangeRateListParams) ([]domain.ExchangeRate, int64, error) {
	if s.listFn != nil {
		return s.listFn(params)
	}
	return nil, 0, nil
}

func (s *exchangeRateRepoStub) Create(rate *domain.ExchangeRate) error {
	if s.createFn != nil {
		return s.createFn(rate)
	}
	return nil
}

func (s *exchangeRateRepoStub) GetByID(id uint64) (*domain.ExchangeRate, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(id)
	}
	return nil, nil
}

func (s *exchangeRateRepoStub) UpdateStatus(id uint64, status domain.ExchangeRateStatus, operatorID uint64) error {
	if s.updateStatusFn != nil {
		return s.updateStatusFn(id, status, operatorID)
	}
	return nil
}

func (s *exchangeRateRepoStub) FindEffectiveRate(fromCurrency, toCurrency string, occurredAt time.Time) (*domain.ExchangeRate, error) {
	if s.findEffectiveFn != nil {
		return s.findEffectiveFn(fromCurrency, toCurrency, occurredAt)
	}
	return nil, nil
}

func TestExchangeRateUsecaseResolveReturnsIdentityForSameCurrency(t *testing.T) {
	uc := NewExchangeRateUsecase(&exchangeRateRepoStub{})
	now := time.Date(2026, 3, 8, 12, 0, 0, 0, time.Local)

	snapshot, err := uc.Resolve("USD", "usd", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snapshot.Rate != 1 {
		t.Fatalf("expected identity rate 1, got %v", snapshot.Rate)
	}
	if snapshot.Source != fxSourceIdentity {
		t.Fatalf("expected identity source, got %s", snapshot.Source)
	}
	if !snapshot.EffectiveAt.Equal(now) {
		t.Fatalf("expected effective_at to equal occurred_at")
	}
}

func TestExchangeRateUsecaseResolveUsesLatestEffectiveRate(t *testing.T) {
	now := time.Date(2026, 3, 8, 12, 0, 0, 0, time.Local)
	var gotFrom, gotTo string
	var gotTime time.Time
	uc := NewExchangeRateUsecase(&exchangeRateRepoStub{
		findEffectiveFn: func(fromCurrency, toCurrency string, occurredAt time.Time) (*domain.ExchangeRate, error) {
			gotFrom, gotTo, gotTime = fromCurrency, toCurrency, occurredAt
			return &domain.ExchangeRate{
				ID:           1,
				FromCurrency: "CNY",
				ToCurrency:   "USD",
				Rate:         0.14,
				SourceType:   domain.ExchangeRateSourceManual,
				EffectiveAt:  now.Add(-time.Hour),
				Status:       domain.ExchangeRateStatusActive,
			}, nil
		},
	})

	snapshot, err := uc.Resolve("USD", "CNY", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotFrom != "CNY" || gotTo != "USD" {
		t.Fatalf("unexpected pair %s -> %s", gotFrom, gotTo)
	}
	if !gotTime.Equal(now) {
		t.Fatalf("unexpected occurred_at passed to repository")
	}
	if snapshot.Rate != 0.14 {
		t.Fatalf("expected rate 0.14, got %v", snapshot.Rate)
	}
	if snapshot.Source != string(domain.ExchangeRateSourceManual) {
		t.Fatalf("unexpected source %s", snapshot.Source)
	}
}

func TestExchangeRateUsecaseResolveReturnsErrorWhenRateMissing(t *testing.T) {
	uc := NewExchangeRateUsecase(&exchangeRateRepoStub{
		findEffectiveFn: func(fromCurrency, toCurrency string, occurredAt time.Time) (*domain.ExchangeRate, error) {
			return nil, errors.New("record not found")
		},
	})

	_, err := uc.Resolve("USD", "CNY", time.Now())
	if err == nil {
		t.Fatal("expected missing rate error")
	}
}

func TestExchangeRateUsecaseResolveFallsBackToInverseRate(t *testing.T) {
	SetExchangeRateScaleResolver(func() uint32 { return 4 })
	defer SetExchangeRateScaleResolver(nil)

	now := time.Date(2026, 3, 9, 10, 0, 0, 0, time.Local)
	callCount := 0
	uc := NewExchangeRateUsecase(&exchangeRateRepoStub{
		findEffectiveFn: func(fromCurrency, toCurrency string, occurredAt time.Time) (*domain.ExchangeRate, error) {
			callCount++
			if fromCurrency == "USD" && toCurrency == "EUR" {
				return nil, errors.New("record not found")
			}
			if fromCurrency == "EUR" && toCurrency == "USD" {
				return &domain.ExchangeRate{
					ID:            2,
					FromCurrency:  "EUR",
					ToCurrency:    "USD",
					Rate:          1.25,
					SourceType:    domain.ExchangeRateSourceManual,
					SourceVersion: "v2",
					EffectiveAt:   now.Add(-2 * time.Hour),
					Status:        domain.ExchangeRateStatusActive,
				}, nil
			}
			return nil, errors.New("record not found")
		},
	})

	snapshot, err := uc.Resolve("EUR", "USD", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Fatalf("expected 2 lookups, got %d", callCount)
	}
	if snapshot.Rate != roundRate(1/1.25, 4) {
		t.Fatalf("expected inverse rate, got %v", snapshot.Rate)
	}
	if snapshot.Source != fxSourceDerived {
		t.Fatalf("expected derived source, got %s", snapshot.Source)
	}
	if snapshot.Version != fxVersionInverse {
		t.Fatalf("expected inverse version, got %s", snapshot.Version)
	}
}

func TestExchangeRateUsecaseResolveFallsBackToUSDCrossRate(t *testing.T) {
	SetExchangeRateScaleResolver(func() uint32 { return 4 })
	defer SetExchangeRateScaleResolver(nil)

	now := time.Date(2026, 3, 9, 10, 0, 0, 0, time.Local)
	uc := NewExchangeRateUsecase(&exchangeRateRepoStub{
		findEffectiveFn: func(fromCurrency, toCurrency string, occurredAt time.Time) (*domain.ExchangeRate, error) {
			switch {
			case fromCurrency == "CNY" && toCurrency == "EUR":
				return nil, errors.New("record not found")
			case fromCurrency == "EUR" && toCurrency == "CNY":
				return nil, errors.New("record not found")
			case fromCurrency == "CNY" && toCurrency == "USD":
				return &domain.ExchangeRate{
					ID:            3,
					FromCurrency:  "CNY",
					ToCurrency:    "USD",
					Rate:          0.14,
					SourceType:    domain.ExchangeRateSourceManual,
					SourceVersion: "seed",
					EffectiveAt:   now.Add(-3 * time.Hour),
					Status:        domain.ExchangeRateStatusActive,
				}, nil
			case fromCurrency == "EUR" && toCurrency == "USD":
				return &domain.ExchangeRate{
					ID:            4,
					FromCurrency:  "EUR",
					ToCurrency:    "USD",
					Rate:          1.08,
					SourceType:    domain.ExchangeRateSourceManual,
					SourceVersion: "seed",
					EffectiveAt:   now.Add(-time.Hour),
					Status:        domain.ExchangeRateStatusActive,
				}, nil
			default:
				return nil, errors.New("record not found")
			}
		},
	})

	snapshot, err := uc.Resolve("EUR", "CNY", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := roundRate(0.14/1.08, 4)
	if snapshot.Rate != expected {
		t.Fatalf("expected cross rate %v, got %v", expected, snapshot.Rate)
	}
	if snapshot.Source != fxSourceDerived {
		t.Fatalf("expected derived source, got %s", snapshot.Source)
	}
	if snapshot.Version != fxVersionCrossUSD {
		t.Fatalf("expected cross-usd version, got %s", snapshot.Version)
	}
	if !snapshot.EffectiveAt.Equal(now.Add(-time.Hour)) {
		t.Fatalf("expected effective_at to use latest component time, got %s", snapshot.EffectiveAt)
	}
}

func TestExchangeRateUsecaseCreateRoundsRateByConfiguredScale(t *testing.T) {
	SetExchangeRateScaleResolver(func() uint32 { return 4 })
	defer SetExchangeRateScaleResolver(nil)

	var created *domain.ExchangeRate
	uc := NewExchangeRateUsecase(&exchangeRateRepoStub{
		createFn: func(rate *domain.ExchangeRate) error {
			created = rate
			rate.ID = 1
			return nil
		},
	})

	remark := "seed"
	createdAt := time.Date(2026, 3, 9, 11, 0, 0, 0, time.Local)
	rate, err := uc.Create(nil, &CreateExchangeRateInput{
		FromCurrency: "CNY",
		ToCurrency:   "USD",
		Rate:         0.137891,
		EffectiveAt:  &createdAt,
		Remark:       &remark,
		CreatedBy:    1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created == nil {
		t.Fatal("expected created rate captured")
	}
	if created.Rate != 0.1379 {
		t.Fatalf("expected persisted rate 0.1379, got %v", created.Rate)
	}
	if rate.Rate != 0.1379 {
		t.Fatalf("expected returned rate 0.1379, got %v", rate.Rate)
	}
}
