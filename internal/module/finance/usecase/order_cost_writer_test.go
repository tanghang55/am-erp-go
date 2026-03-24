package usecase

import (
	"testing"
	"time"

	"am-erp-go/internal/module/finance/domain"
	salesUsecase "am-erp-go/internal/module/sales/usecase"
)

type stubOrderCostDetailRepo struct {
	details []domain.OrderCostDetail
	rows    []domain.ReturnableOrderCostDetail
	err     error
}

func (s *stubOrderCostDetailRepo) CreateBatch(details []domain.OrderCostDetail) error {
	s.details = append(s.details, details...)
	return s.err
}

func (s *stubOrderCostDetailRepo) SumQtyByDateAndMarketplace(_ time.Time, _ *string) (uint64, error) {
	return 0, nil
}

func (s *stubOrderCostDetailRepo) ListReturnableBySalesOrderItemID(_ uint64) ([]domain.ReturnableOrderCostDetail, error) {
	return s.rows, s.err
}

func TestOrderCostWriterRecordSalesShipCost(t *testing.T) {
	repo := &stubOrderCostDetailRepo{}
	writer := NewOrderCostWriter(repo)
	now := time.Now()

	err := writer.RecordSalesShipCost(&salesUsecase.SalesShipCostRecordParams{
		SalesOrderID:     12,
		SalesOrderItemID: 33,
		ProductID:        99,
		WarehouseID:      2,
		Marketplace:      "US",
		Currency:         "USD",
		OccurredAt:       now,
		Allocations: []salesUsecase.SalesShipCostAllocation{
			{InventoryLotID: 1001, Qty: 3, UnitCost: 1.8},
			{InventoryLotID: 1002, Qty: 2, UnitCost: 2.1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.details) != 2 {
		t.Fatalf("expected 2 details, got %d", len(repo.details))
	}
	if repo.details[0].InventoryLotID != 1001 || repo.details[0].QtyOut != 3 {
		t.Fatalf("unexpected first detail: %+v", repo.details[0])
	}
	if repo.details[1].InventoryLotID != 1002 || repo.details[1].QtyOut != 2 {
		t.Fatalf("unexpected second detail: %+v", repo.details[1])
	}
	if repo.details[0].OriginalAmount != round6(3*1.8) {
		t.Fatalf("unexpected original amount for first detail: %v", repo.details[0].OriginalAmount)
	}
}

func TestOrderCostWriterRecordSalesReturnCost(t *testing.T) {
	now := time.Now()
	repo := &stubOrderCostDetailRepo{
		rows: []domain.ReturnableOrderCostDetail{
			{
				OrderCostDetail: domain.OrderCostDetail{
					ID:               11,
					SalesOrderID:     12,
					SalesOrderItemID: 33,
					ProductID:        99,
					WarehouseID:      2,
					InventoryLotID:   1001,
					QtyOut:           3,
					UnitCostOriginal: 1.8,
					OriginalCurrency: "USD",
					BaseCurrency:     "USD",
					FxRate:           1,
					FxTime:           now,
				},
				AvailableQty: 2,
			},
			{
				OrderCostDetail: domain.OrderCostDetail{
					ID:               12,
					SalesOrderID:     12,
					SalesOrderItemID: 33,
					ProductID:        99,
					WarehouseID:      2,
					InventoryLotID:   1002,
					QtyOut:           2,
					UnitCostOriginal: 2.1,
					OriginalCurrency: "USD",
					BaseCurrency:     "USD",
					FxRate:           1,
					FxTime:           now,
				},
				AvailableQty: 2,
			},
		},
	}
	writer := NewOrderCostWriter(repo)

	cogs, err := writer.RecordSalesReturnCost(&salesUsecase.SalesReturnCostRecordParams{
		SalesOrderID:     12,
		SalesOrderItemID: 33,
		ProductID:        99,
		WarehouseID:      2,
		Marketplace:      "US",
		Currency:         "USD",
		QtyReturned:      3,
		OccurredAt:       now,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cogs != round6(2*1.8+1*2.1) {
		t.Fatalf("unexpected reversed cogs: %v", cogs)
	}
	if len(repo.details) != 2 {
		t.Fatalf("expected 2 reversal details, got %d", len(repo.details))
	}
	if repo.details[0].ReversalOfID == nil || *repo.details[0].ReversalOfID != 11 {
		t.Fatalf("expected first reversal_of_id=11, got %+v", repo.details[0].ReversalOfID)
	}
	if repo.details[0].QtyOut != 2 || repo.details[1].QtyOut != 1 {
		t.Fatalf("unexpected reversal qtys: %+v", repo.details)
	}
}

func TestOrderCostWriterResolveSalesReturnUnitCost(t *testing.T) {
	now := time.Now()
	repo := &stubOrderCostDetailRepo{
		rows: []domain.ReturnableOrderCostDetail{
			{
				OrderCostDetail: domain.OrderCostDetail{
					ID:               21,
					SalesOrderItemID: 44,
					InventoryLotID:   1001,
					QtyOut:           2,
					UnitCostOriginal: 10.5,
					FxTime:           now,
				},
				AvailableQty: 2,
			},
			{
				OrderCostDetail: domain.OrderCostDetail{
					ID:               22,
					SalesOrderItemID: 44,
					InventoryLotID:   1002,
					QtyOut:           1,
					UnitCostOriginal: 14.0,
					FxTime:           now,
				},
				AvailableQty: 1,
			},
		},
	}
	writer := NewOrderCostWriter(repo)

	unitCost, err := writer.ResolveSalesReturnUnitCost(&salesUsecase.SalesReturnCostRecordParams{
		SalesOrderItemID: 44,
		QtyReturned:      2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if unitCost == nil {
		t.Fatalf("expected unit cost")
	}
	if *unitCost != 10.5 {
		t.Fatalf("expected unit cost 10.5, got %v", *unitCost)
	}
}
