package usecase

import (
	"errors"
	"testing"

	"am-erp-go/internal/module/logistics/domain"
)

type stubShippingRateRepo struct {
	created  *domain.ShippingRate
	rates    []*domain.ShippingRate
	countMap map[uint64]int64
}

func (s *stubShippingRateRepo) Create(rate *domain.ShippingRate) error {
	s.created = rate
	return nil
}

func (s *stubShippingRateRepo) Update(rate *domain.ShippingRate) error {
	return nil
}

func (s *stubShippingRateRepo) Delete(id uint64) error {
	return nil
}

func (s *stubShippingRateRepo) GetByID(id uint64) (*domain.ShippingRate, error) {
	return nil, nil
}

func (s *stubShippingRateRepo) List(params *domain.ShippingRateListParams) ([]*domain.ShippingRate, int64, error) {
	if s.rates != nil {
		return s.rates, int64(len(s.rates)), nil
	}
	return nil, 0, nil
}

func (s *stubShippingRateRepo) QueryLatestRate(params *domain.QueryLatestRateParams) (*domain.ShippingRate, error) {
	return nil, nil
}

func (s *stubShippingRateRepo) CountReferences(id uint64) (int64, error) {
	if s.countMap == nil {
		return 0, nil
	}
	return s.countMap[id], nil
}

type stubLogisticsProviderRepo struct{}

func (s *stubLogisticsProviderRepo) Create(provider *domain.LogisticsProvider) error { return nil }
func (s *stubLogisticsProviderRepo) Update(provider *domain.LogisticsProvider) error { return nil }
func (s *stubLogisticsProviderRepo) Delete(id uint64) error                          { return nil }
func (s *stubLogisticsProviderRepo) GetByID(id uint64) (*domain.LogisticsProvider, error) {
	return &domain.LogisticsProvider{ID: id}, nil
}
func (s *stubLogisticsProviderRepo) List(params *domain.LogisticsProviderListParams) ([]*domain.LogisticsProvider, int64, error) {
	return nil, 0, nil
}
func (s *stubLogisticsProviderRepo) GetByCode(code string) (*domain.LogisticsProvider, error) {
	return nil, nil
}
func (s *stubLogisticsProviderRepo) CountReferences(id uint64) (int64, error) {
	return 0, nil
}

type stubShippingRateBaseCurrencyProvider struct{}

func (stubShippingRateBaseCurrencyProvider) GetDefaultBaseCurrency() string {
	return "EUR"
}

func TestCreateShippingRateFallsBackToConfiguredCurrency(t *testing.T) {
	rateRepo := &stubShippingRateRepo{}
	uc := NewShippingRateUsecase(rateRepo, &stubLogisticsProviderRepo{})
	uc.BindDefaultsProvider(stubShippingRateBaseCurrencyProvider{})

	rate, err := uc.Create(&domain.CreateShippingRateParams{
		ProviderID:             1,
		OriginWarehouseID:      1,
		DestinationWarehouseID: 2,
		TransportMode:          domain.TransportModeSea,
		PricingMethod:          domain.PricingMethodPerKg,
		BaseRate:               12.5,
		EffectiveDate:          "2026-03-09",
		Status:                 domain.RateStatusActive,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rate.Currency != "EUR" {
		t.Fatalf("expected currency EUR, got %s", rate.Currency)
	}
	if rateRepo.created == nil || rateRepo.created.Currency != "EUR" {
		t.Fatalf("expected persisted currency EUR, got %+v", rateRepo.created)
	}
}

func TestListShippingRatesIncludesDeleteState(t *testing.T) {
	repo := &stubShippingRateRepo{
		rates: []*domain.ShippingRate{
			{ID: 1},
			{ID: 2},
		},
		countMap: map[uint64]int64{
			1: 2,
			2: 0,
		},
	}
	uc := NewShippingRateUsecase(repo, &stubLogisticsProviderRepo{})

	rates, total, err := uc.List(&domain.ShippingRateListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	if rates[0].Deletable {
		t.Fatalf("expected referenced rate to be non-deletable")
	}
	if rates[0].ReferenceCount != 2 {
		t.Fatalf("expected reference_count 2, got %d", rates[0].ReferenceCount)
	}
	if !rates[1].Deletable {
		t.Fatalf("expected unused rate to be deletable")
	}
}

func TestDeleteShippingRateRejectsReferencedRate(t *testing.T) {
	repo := &stubShippingRateRepo{
		countMap: map[uint64]int64{
			5: 1,
		},
	}
	uc := NewShippingRateUsecase(repo, &stubLogisticsProviderRepo{})

	err := uc.Delete(5)
	if !errors.Is(err, ErrRateReferenced) {
		t.Fatalf("expected ErrRateReferenced, got %v", err)
	}
}
