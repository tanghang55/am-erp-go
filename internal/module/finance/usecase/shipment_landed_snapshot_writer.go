package usecase

import (
	"fmt"
	"strings"
	"time"

	"am-erp-go/internal/module/finance/domain"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"
)

type ShipmentLandedSnapshotWriter struct {
	repo     domain.CostingSnapshotRepository
	costRepo domain.CostEventRepository
}

func NewShipmentLandedSnapshotWriter(repo domain.CostingSnapshotRepository, costRepo domain.CostEventRepository) *ShipmentLandedSnapshotWriter {
	return &ShipmentLandedSnapshotWriter{repo: repo, costRepo: costRepo}
}

func (w *ShipmentLandedSnapshotWriter) UpsertShipmentLandedSnapshots(params *shipmentUsecase.ShipmentCostAllocationRecordParams) error {
	if w.repo == nil || params == nil || len(params.Lines) == 0 {
		return nil
	}

	baseCurrency := normalizeCurrency(params.BaseCurrency)
	if baseCurrency == "" {
		baseCurrency = getDefaultBaseCurrency()
	}
	effectiveFrom := params.OccurredAt
	if effectiveFrom.IsZero() {
		effectiveFrom = time.Now()
	}

	for _, line := range params.Lines {
		if line.ProductID == 0 || line.Quantity == 0 {
			continue
		}
		sourceCurrency := normalizeCurrency(line.ItemCurrency)
		if sourceCurrency == "" {
			sourceCurrency = baseCurrency
		}
		sourceFx, err := resolveFXRate(baseCurrency, sourceCurrency, effectiveFrom)
		if err != nil {
			return err
		}
		sourceUnitBase := round6(line.ItemUnitCost * sourceFx.Rate)
		packingUnitCost := 0.0
		if w.costRepo != nil {
			value, err := w.costRepo.GetLatestPackingMaterialPerUnit(line.ProductID, effectiveFrom)
			if err != nil {
				return err
			}
			if value != nil {
				packingUnitCost = *value
			}
		}
		allocatedPerUnit := round6(line.BaseAmount / float64(line.Quantity))
		landedUnitCost := round6(sourceUnitBase + packingUnitCost + allocatedPerUnit)
		if landedUnitCost <= 0 {
			continue
		}
		notes := fmt.Sprintf("auto from shipment %s quantity allocation", strings.TrimSpace(params.ShipmentNumber))
		if err := w.repo.ExpireCurrent(line.ProductID, domain.CostTypeLanded, effectiveFrom, nil); err != nil {
			return err
		}
		snapshot := &domain.CostingSnapshot{
			TraceID:       fmt.Sprintf("LANDED-SHIP-%d-%d", time.Now().UnixNano(), line.ProductID),
			ProductID:     line.ProductID,
			CostType:      domain.CostTypeLanded,
			UnitCost:      landedUnitCost,
			Currency:      baseCurrency,
			EffectiveFrom: effectiveFrom,
			Notes:         &notes,
			CreatedBy:     valueOrZero(params.OperatorID),
		}
		if snapshot.CreatedBy == 0 {
			snapshot.CreatedBy = 1
		}
		if err := w.repo.Create(snapshot); err != nil {
			return err
		}
	}

	return nil
}

func valueOrZero(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}
