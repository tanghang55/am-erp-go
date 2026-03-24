package usecase

import (
	"testing"
	"time"

	"am-erp-go/internal/module/finance/domain"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"
)

type costingSnapshotWriteRepoStub struct {
	created []*domain.CostingSnapshot
	expired []struct {
		productID   uint64
		costType    domain.CostType
		effectiveTo time.Time
	}
}

func (s *costingSnapshotWriteRepoStub) List(params *domain.CostingSnapshotListParams) ([]domain.CostingSnapshot, int64, error) {
	return nil, 0, nil
}
func (s *costingSnapshotWriteRepoStub) GetByID(id uint64) (*domain.CostingSnapshot, error) {
	return nil, nil
}

type costEventPackingReadRepoStub struct {
	perUnit map[uint64]float64
}

func (s *costEventPackingReadRepoStub) Create(event *domain.CostEvent) error {
	return nil
}

func (s *costEventPackingReadRepoStub) GetLatestPackingMaterialPerUnit(productID uint64, occurredAt time.Time) (*float64, error) {
	if s == nil || s.perUnit == nil {
		return nil, nil
	}
	value, ok := s.perUnit[productID]
	if !ok {
		return nil, nil
	}
	return &value, nil
}
func (s *costingSnapshotWriteRepoStub) Create(snapshot *domain.CostingSnapshot) error {
	s.created = append(s.created, snapshot)
	return nil
}
func (s *costingSnapshotWriteRepoStub) Update(snapshot *domain.CostingSnapshot) error {
	return nil
}
func (s *costingSnapshotWriteRepoStub) Delete(id uint64) error {
	return nil
}
func (s *costingSnapshotWriteRepoStub) ExpireCurrent(productID uint64, costType domain.CostType, effectiveTo time.Time, excludeID *uint64) error {
	s.expired = append(s.expired, struct {
		productID   uint64
		costType    domain.CostType
		effectiveTo time.Time
	}{productID: productID, costType: costType, effectiveTo: effectiveTo})
	return nil
}
func (s *costingSnapshotWriteRepoStub) GetCurrent(productID uint64, costType domain.CostType, now time.Time) (*domain.CostingSnapshot, error) {
	return nil, nil
}
func (s *costingSnapshotWriteRepoStub) ListCurrentBySKU(productID uint64, now time.Time) ([]domain.CostingSnapshot, error) {
	return nil, nil
}

func TestShipmentLandedSnapshotWriterCreatesLandedSnapshot(t *testing.T) {
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

	repo := &costingSnapshotWriteRepoStub{}
	writer := NewShipmentLandedSnapshotWriter(repo, nil)
	now := time.Date(2026, 3, 9, 16, 0, 0, 0, time.Local)
	operatorID := uint64(9)

	err := writer.UpsertShipmentLandedSnapshots(&shipmentUsecase.ShipmentCostAllocationRecordParams{
		ShipmentID:     1,
		ShipmentNumber: "SHP-100",
		BaseCurrency:   "USD",
		OccurredAt:     now,
		OperatorID:     &operatorID,
		Lines: []shipmentUsecase.ShipmentCostAllocationLine{
			{
				ProductID:      1001,
				Quantity:       2,
				ItemUnitCost:   5,
				ItemCurrency:   "USD",
				BaseAmount:     4,
				OriginalAmount: 4,
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.expired) != 1 {
		t.Fatalf("expected 1 expire call, got %d", len(repo.expired))
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(repo.created))
	}
	if repo.created[0].CostType != domain.CostTypeLanded {
		t.Fatalf("expected landed cost type, got %s", repo.created[0].CostType)
	}
	if repo.created[0].UnitCost != 7 {
		t.Fatalf("expected landed unit cost 7, got %v", repo.created[0].UnitCost)
	}
}

func TestShipmentLandedSnapshotWriterIncludesPackingMaterialUnitCost(t *testing.T) {
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

	repo := &costingSnapshotWriteRepoStub{}
	costRepo := &costEventPackingReadRepoStub{perUnit: map[uint64]float64{1001: 1.5}}
	writer := NewShipmentLandedSnapshotWriter(repo, costRepo)
	now := time.Date(2026, 3, 23, 11, 0, 0, 0, time.Local)
	operatorID := uint64(9)

	err := writer.UpsertShipmentLandedSnapshots(&shipmentUsecase.ShipmentCostAllocationRecordParams{
		ShipmentID:     1,
		ShipmentNumber: "SHP-200",
		BaseCurrency:   "USD",
		OccurredAt:     now,
		OperatorID:     &operatorID,
		Lines: []shipmentUsecase.ShipmentCostAllocationLine{
			{
				ProductID:      1001,
				Quantity:       2,
				ItemUnitCost:   5,
				ItemCurrency:   "USD",
				BaseAmount:     4,
				OriginalAmount: 4,
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(repo.created))
	}
	if repo.created[0].UnitCost != 8.5 {
		t.Fatalf("expected landed unit cost 8.5, got %v", repo.created[0].UnitCost)
	}
}
