package usecase

import (
	"testing"
	"time"

	"am-erp-go/internal/module/finance/domain"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"
)

type costEventRepoStub struct {
	created []*domain.CostEvent
}

func (s *costEventRepoStub) Create(event *domain.CostEvent) error {
	s.created = append(s.created, event)
	return nil
}

func (s *costEventRepoStub) GetLatestPackingMaterialPerUnit(productID uint64, occurredAt time.Time) (*float64, error) {
	return nil, nil
}

func TestCostEventWriterRecordShipmentCostAllocation(t *testing.T) {
	repo := &costEventRepoStub{}
	writer := NewCostEventWriter(repo)
	now := time.Date(2026, 3, 9, 10, 0, 0, 0, time.Local)
	operatorID := uint64(7)

	err := writer.RecordShipmentCostAllocation(&shipmentUsecase.ShipmentCostAllocationRecordParams{
		ShipmentID:       88,
		ShipmentNumber:   "SHP-88",
		OriginalCurrency: "USD",
		BaseCurrency:     "EUR",
		FxRate:           1.2,
		FxSource:         "MANUAL",
		FxVersion:        "v1",
		FxTime:           now,
		OccurredAt:       now,
		OperatorID:       &operatorID,
		Lines: []shipmentUsecase.ShipmentCostAllocationLine{
			{
				ShipmentItemID: 11,
				ProductID:      101,
				WarehouseID:    3,
				Quantity:       2,
				OriginalAmount: 20,
				BaseAmount:     24,
			},
			{
				ShipmentItemID: 12,
				ProductID:      202,
				WarehouseID:    3,
				Quantity:       1,
				OriginalAmount: 10,
				BaseAmount:     12,
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.created) != 2 {
		t.Fatalf("expected 2 events, got %d", len(repo.created))
	}
	if repo.created[0].EventType != domain.CostEventTypeShipmentAllocated {
		t.Fatalf("expected shipment allocated event, got %s", repo.created[0].EventType)
	}
	if repo.created[0].ShipmentID == nil || *repo.created[0].ShipmentID != 88 {
		t.Fatalf("unexpected shipment id: %+v", repo.created[0].ShipmentID)
	}
	if repo.created[0].OriginalAmount != 20 || repo.created[0].BaseAmount != 24 {
		t.Fatalf("unexpected amounts: %+v", repo.created[0])
	}
}

func TestCostEventWriterRecordPackingMaterialCost(t *testing.T) {
	SetDefaultBaseCurrencyResolver(func() string { return "USD" })
	SetExchangeRateScaleResolver(func() uint32 { return 4 })
	SetFXRateResolver(func(baseCurrency, originalCurrency string, occurredAt time.Time) (*FXRateSnapshot, error) {
		return &FXRateSnapshot{
			Rate:        1,
			Source:      "IDENTITY",
			Version:     "same_currency",
			EffectiveAt: occurredAt,
		}, nil
	})

	repo := &costEventRepoStub{}
	writer := NewCostEventWriter(repo)
	now := time.Date(2026, 3, 23, 10, 0, 0, 0, time.Local)
	operatorID := uint64(8)
	movementID := uint64(66)

	err := writer.RecordPackingMaterialCost(&inventoryUsecase.PackingMaterialCostRecordParams{
		InventoryMovementID: movementID,
		ProductID:           101,
		WarehouseID:         3,
		Quantity:            2,
		ReferenceNumber:     "PACK-66",
		OccurredAt:          now,
		OperatorID:          &operatorID,
		Lines: []inventoryUsecase.PackingMaterialCostLine{
			{
				PackagingItemID: 501,
				Quantity:        4,
				UnitCost:        1.5,
				Currency:        "USD",
				ItemCode:        "BOX-1",
				ItemName:        "纸箱",
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected 1 event, got %d", len(repo.created))
	}
	if repo.created[0].EventType != domain.CostEventTypePackingMaterial {
		t.Fatalf("expected packing material event, got %s", repo.created[0].EventType)
	}
	if repo.created[0].InventoryMovementID == nil || *repo.created[0].InventoryMovementID != movementID {
		t.Fatalf("unexpected inventory movement id: %+v", repo.created[0].InventoryMovementID)
	}
	if repo.created[0].QtyEvent != 2 {
		t.Fatalf("expected qty_event 2, got %+v", repo.created[0])
	}
	if repo.created[0].OriginalAmount != 6 || repo.created[0].BaseAmount != 6 {
		t.Fatalf("unexpected amounts: %+v", repo.created[0])
	}
}
