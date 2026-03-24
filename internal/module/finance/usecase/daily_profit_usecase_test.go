package usecase

import (
	"testing"
	"time"

	"am-erp-go/internal/module/finance/domain"
)

type stubDailyProfitLedgerRepo struct {
	aggs []domain.ProfitLedgerDailyAgg
}

func (s *stubDailyProfitLedgerRepo) Create(_ *domain.ProfitLedger) error {
	return nil
}

func (s *stubDailyProfitLedgerRepo) CreateBatch(_ []domain.ProfitLedger) error {
	return nil
}

func (s *stubDailyProfitLedgerRepo) AggregateDaily(_ time.Time, _ *string) ([]domain.ProfitLedgerDailyAgg, error) {
	return s.aggs, nil
}

type stubDailySnapshotRepo struct {
	deletedDate *time.Time
	created     []domain.DailyProfitSnapshot
	listItems   []domain.DailyProfitSnapshot
}

func (s *stubDailySnapshotRepo) DeleteByDate(bizDate time.Time, _ *string) error {
	s.deletedDate = &bizDate
	return nil
}

func (s *stubDailySnapshotRepo) CreateBatch(items []domain.DailyProfitSnapshot) error {
	s.created = append(s.created, items...)
	return nil
}

func (s *stubDailySnapshotRepo) List(_ *domain.DailyProfitSnapshotListParams) ([]domain.DailyProfitSnapshot, error) {
	return s.listItems, nil
}

type stubDailyOrderCostRepo struct {
	qty uint64
}

func (s *stubDailyOrderCostRepo) CreateBatch(_ []domain.OrderCostDetail) error {
	return nil
}

func (s *stubDailyOrderCostRepo) SumQtyByDateAndMarketplace(_ time.Time, _ *string) (uint64, error) {
	return s.qty, nil
}

func (s *stubDailyOrderCostRepo) ListReturnableBySalesOrderItemID(_ uint64) ([]domain.ReturnableOrderCostDetail, error) {
	return nil, nil
}

func TestDailyProfitUsecase_Rebuild(t *testing.T) {
	profitRepo := &stubDailyProfitLedgerRepo{
		aggs: []domain.ProfitLedgerDailyAgg{
			{
				Marketplace:   "US",
				BaseCurrency:  "USD",
				SalesIncome:   100,
				COGS:          40,
				OrderExpense:  10,
				PublicExpense: 5,
				OrderCount:    2,
			},
		},
	}
	snapshotRepo := &stubDailySnapshotRepo{}
	orderCostRepo := &stubDailyOrderCostRepo{qty: 7}
	uc := NewDailyProfitUsecase(profitRepo, snapshotRepo, orderCostRepo)

	items, err := uc.Rebuild(&RebuildDailyProfitInput{
		BizDate: time.Date(2026, 3, 3, 10, 0, 0, 0, time.Local),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(items))
	}
	if items[0].GrossProfitAmount != 60 {
		t.Fatalf("expected gross 60, got %v", items[0].GrossProfitAmount)
	}
	if items[0].OperatingNetProfitAmount != 45 {
		t.Fatalf("expected operating net 45, got %v", items[0].OperatingNetProfitAmount)
	}
	if items[0].ShippedQty != 7 {
		t.Fatalf("expected shipped qty 7, got %d", items[0].ShippedQty)
	}
	if len(snapshotRepo.created) != 1 {
		t.Fatalf("expected 1 persisted snapshot, got %d", len(snapshotRepo.created))
	}
}
