package usecase

import (
	"am-erp-go/internal/module/finance/domain"
	"testing"
	"time"
)

type costingSnapshotRepoStub struct {
	created *domain.CostingSnapshot
}

func (s *costingSnapshotRepoStub) List(params *domain.CostingSnapshotListParams) ([]domain.CostingSnapshot, int64, error) {
	return nil, 0, nil
}

func (s *costingSnapshotRepoStub) GetByID(id uint64) (*domain.CostingSnapshot, error) {
	return nil, nil
}

func (s *costingSnapshotRepoStub) Create(snapshot *domain.CostingSnapshot) error {
	s.created = snapshot
	return nil
}

func (s *costingSnapshotRepoStub) Update(snapshot *domain.CostingSnapshot) error {
	return nil
}

func (s *costingSnapshotRepoStub) Delete(id uint64) error {
	return nil
}

func (s *costingSnapshotRepoStub) ExpireCurrent(productID uint64, costType domain.CostType, effectiveTo time.Time, excludeID *uint64) error {
	return nil
}

func (s *costingSnapshotRepoStub) GetCurrent(productID uint64, costType domain.CostType, now time.Time) (*domain.CostingSnapshot, error) {
	return nil, nil
}

func (s *costingSnapshotRepoStub) ListCurrentBySKU(productID uint64, now time.Time) ([]domain.CostingSnapshot, error) {
	return nil, nil
}

func TestCostingSnapshotCreateUsesDefaultBaseCurrency(t *testing.T) {
	SetDefaultBaseCurrencyResolver(func() string { return "USD" })
	t.Cleanup(func() {
		SetDefaultBaseCurrencyResolver(nil)
	})

	repo := &costingSnapshotRepoStub{}
	uc := NewCostingSnapshotUsecase(repo)

	snapshot, err := uc.Create(&CreateCostingSnapshotInput{
		ProductID: 1001,
		CostType:  domain.CostTypePurchase,
		UnitCost:  12.34,
		CreatedBy: 1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if snapshot.Currency != "USD" {
		t.Fatalf("expected default currency USD, got %s", snapshot.Currency)
	}
	if repo.created == nil {
		t.Fatalf("expected snapshot to be created")
	}
	if repo.created.Currency != "USD" {
		t.Fatalf("expected persisted currency USD, got %s", repo.created.Currency)
	}
}
