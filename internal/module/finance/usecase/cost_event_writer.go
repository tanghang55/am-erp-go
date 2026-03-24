package usecase

import (
	"fmt"
	"time"

	"am-erp-go/internal/module/finance/domain"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	procurementUsecase "am-erp-go/internal/module/procurement/usecase"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"
)

type CostEventWriter struct {
	repo domain.CostEventRepository
}

func NewCostEventWriter(repo domain.CostEventRepository) *CostEventWriter {
	return &CostEventWriter{repo: repo}
}

func (w *CostEventWriter) RecordPurchaseOrderEvent(params *procurementUsecase.PurchaseOrderCostEventParams) error {
	if w.repo == nil || params == nil || params.QtyEvent == 0 || params.ProductID == 0 {
		return nil
	}

	occurredAt := params.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}

	originalCurrency := normalizeCurrency(params.Currency)
	if originalCurrency == "" {
		originalCurrency = getDefaultBaseCurrency()
	}
	baseCurrency := getDefaultBaseCurrency()
	fxSnapshot, err := resolveFXRate(baseCurrency, originalCurrency, occurredAt)
	if err != nil {
		return err
	}

	originalAmount := round6(float64(params.QtyEvent) * params.UnitCost)
	baseAmount := round6(originalAmount * fxSnapshot.Rate)
	eventType, err := toCostEventType(params.EventType)
	if err != nil {
		return err
	}

	poID := params.PurchaseOrderID
	poItemID := params.PurchaseOrderItemID
	marketplace := params.Marketplace
	remark := fmt.Sprintf("auto from purchase order %s", params.EventType)

	event := &domain.CostEvent{
		TraceID:             fmt.Sprintf("COST-EVT-%d-%d", time.Now().UnixNano(), params.ProductID),
		EventType:           eventType,
		Status:              domain.CostEventStatusNormal,
		PurchaseOrderID:     &poID,
		PurchaseOrderItemID: &poItemID,
		ProductID:           params.ProductID,
		WarehouseID:         params.WarehouseID,
		Marketplace:         &marketplace,
		QtyEvent:            params.QtyEvent,
		OriginalCurrency:    originalCurrency,
		OriginalAmount:      originalAmount,
		BaseCurrency:        baseCurrency,
		FxRate:              fxSnapshot.Rate,
		BaseAmount:          baseAmount,
		FxSource:            fxSnapshot.Source,
		FxVersion:           fxSnapshot.Version,
		FxTime:              fxSnapshot.EffectiveAt,
		OccurredAt:          occurredAt,
		OperatorID:          params.OperatorID,
		Remark:              &remark,
	}

	return w.repo.Create(event)
}

func (w *CostEventWriter) RecordShipmentCostAllocation(params *shipmentUsecase.ShipmentCostAllocationRecordParams) error {
	if w.repo == nil || params == nil || params.ShipmentID == 0 || len(params.Lines) == 0 {
		return nil
	}

	for _, line := range params.Lines {
		if line.ProductID == 0 || line.OriginalAmount <= 0 || line.BaseAmount <= 0 {
			continue
		}
		shipmentID := params.ShipmentID
		shipmentItemID := line.ShipmentItemID
		warehouseID := line.WarehouseID
		marketplace := line.Marketplace
		remark := fmt.Sprintf("auto allocated from shipment %s", params.ShipmentNumber)

		event := &domain.CostEvent{
			TraceID:          fmt.Sprintf("SHIP-COST-%d-%d-%d", time.Now().UnixNano(), params.ShipmentID, line.ProductID),
			EventType:        domain.CostEventTypeShipmentAllocated,
			Status:           domain.CostEventStatusNormal,
			ShipmentID:       &shipmentID,
			ShipmentItemID:   &shipmentItemID,
			ProductID:        line.ProductID,
			WarehouseID:      &warehouseID,
			Marketplace:      stringPtrOrNil(marketplace),
			QtyEvent:         line.Quantity,
			OriginalCurrency: normalizeCurrency(params.OriginalCurrency),
			OriginalAmount:   round6(line.OriginalAmount),
			BaseCurrency:     normalizeCurrency(params.BaseCurrency),
			FxRate:           params.FxRate,
			BaseAmount:       round6(line.BaseAmount),
			FxSource:         params.FxSource,
			FxVersion:        params.FxVersion,
			FxTime:           params.FxTime,
			OccurredAt:       params.OccurredAt,
			OperatorID:       params.OperatorID,
			Remark:           &remark,
		}
		if err := w.repo.Create(event); err != nil {
			return err
		}
	}
	return nil
}

func (w *CostEventWriter) RecordPackingMaterialCost(params *inventoryUsecase.PackingMaterialCostRecordParams) error {
	if w.repo == nil || params == nil || params.ProductID == 0 || params.Quantity == 0 || len(params.Lines) == 0 {
		return nil
	}

	occurredAt := params.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}

	baseCurrency := getDefaultBaseCurrency()
	if baseCurrency == "" {
		baseCurrency = "USD"
	}

	originalCurrency := ""
	originalAmount := 0.0
	baseAmount := 0.0
	fxRate := 1.0
	fxSource := "PACKING_NORMALIZED"
	fxVersion := "mixed_currency"
	fxTime := occurredAt
	mixedCurrency := false

	for _, line := range params.Lines {
		if line.Quantity == 0 || line.UnitCost <= 0 {
			continue
		}
		currency := normalizeCurrency(line.Currency)
		if currency == "" {
			currency = baseCurrency
		}
		lineOriginalAmount := round6(float64(line.Quantity) * line.UnitCost)
		lineFX, err := resolveFXRate(baseCurrency, currency, occurredAt)
		if err != nil {
			return err
		}
		lineBaseAmount := round6(lineOriginalAmount * lineFX.Rate)
		baseAmount = round6(baseAmount + lineBaseAmount)

		if originalCurrency == "" {
			originalCurrency = currency
			originalAmount = round6(originalAmount + lineOriginalAmount)
			fxRate = lineFX.Rate
			fxSource = lineFX.Source
			fxVersion = lineFX.Version
			fxTime = lineFX.EffectiveAt
			continue
		}
		if originalCurrency != currency {
			mixedCurrency = true
		} else {
			originalAmount = round6(originalAmount + lineOriginalAmount)
		}
	}

	if baseAmount <= 0 {
		return nil
	}
	if mixedCurrency || originalCurrency == "" {
		originalCurrency = baseCurrency
		originalAmount = baseAmount
		fxRate = 1
		fxSource = "PACKING_NORMALIZED"
		fxVersion = "mixed_currency"
		fxTime = occurredAt
	}

	warehouseID := params.WarehouseID
	movementID := params.InventoryMovementID
	remark := fmt.Sprintf("auto allocated from product packing %s", params.ReferenceNumber)
	event := &domain.CostEvent{
		TraceID:             fmt.Sprintf("PACK-COST-%d-%d-%d", time.Now().UnixNano(), params.InventoryMovementID, params.ProductID),
		EventType:           domain.CostEventTypePackingMaterial,
		Status:              domain.CostEventStatusNormal,
		InventoryMovementID: &movementID,
		ProductID:           params.ProductID,
		WarehouseID:         &warehouseID,
		QtyEvent:            params.Quantity,
		OriginalCurrency:    originalCurrency,
		OriginalAmount:      originalAmount,
		BaseCurrency:        baseCurrency,
		FxRate:              fxRate,
		BaseAmount:          baseAmount,
		FxSource:            fxSource,
		FxVersion:           fxVersion,
		FxTime:              fxTime,
		OccurredAt:          occurredAt,
		OperatorID:          params.OperatorID,
		Remark:              &remark,
	}
	return w.repo.Create(event)
}

func toCostEventType(t procurementUsecase.PurchaseOrderCostEventType) (domain.CostEventType, error) {
	switch t {
	case procurementUsecase.PurchaseOrderCostEventOrdered:
		return domain.CostEventTypePOOrdered, nil
	case procurementUsecase.PurchaseOrderCostEventShipped:
		return domain.CostEventTypePOShipped, nil
	case procurementUsecase.PurchaseOrderCostEventReceived:
		return domain.CostEventTypePOReceived, nil
	default:
		return "", fmt.Errorf("unsupported purchase order cost event type: %s", t)
	}
}

func stringPtrOrNil(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
