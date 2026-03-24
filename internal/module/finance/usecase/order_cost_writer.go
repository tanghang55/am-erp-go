package usecase

import (
	"fmt"
	"time"

	"am-erp-go/internal/module/finance/domain"
	salesUsecase "am-erp-go/internal/module/sales/usecase"
)

type OrderCostWriter struct {
	repo domain.OrderCostDetailRepository
}

func NewOrderCostWriter(repo domain.OrderCostDetailRepository) *OrderCostWriter {
	return &OrderCostWriter{repo: repo}
}

func (w *OrderCostWriter) RecordSalesShipCost(params *salesUsecase.SalesShipCostRecordParams) error {
	if w.repo == nil || params == nil || params.SalesOrderID == 0 || params.SalesOrderItemID == 0 {
		return nil
	}
	if params.ProductID == 0 || params.WarehouseID == 0 || len(params.Allocations) == 0 {
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

	var marketplace *string
	if params.Marketplace != "" {
		m := params.Marketplace
		marketplace = &m
	}
	remark := "auto from sales ship fifo allocation"

	details := make([]domain.OrderCostDetail, 0, len(params.Allocations))
	for _, allocation := range params.Allocations {
		if allocation.InventoryLotID == 0 || allocation.Qty == 0 {
			continue
		}
		originalAmount := round6(float64(allocation.Qty) * allocation.UnitCost)
		baseAmount := round6(originalAmount * fxSnapshot.Rate)
		details = append(details, domain.OrderCostDetail{
			TraceID:          fmt.Sprintf("ORDER-COST-%d-%d-%d", time.Now().UnixNano(), params.SalesOrderID, allocation.InventoryLotID),
			Status:           domain.OrderCostDetailStatusNormal,
			SalesOrderID:     params.SalesOrderID,
			SalesOrderItemID: params.SalesOrderItemID,
			ProductID:        params.ProductID,
			WarehouseID:      params.WarehouseID,
			Marketplace:      marketplace,
			InventoryLotID:   allocation.InventoryLotID,
			QtyOut:           allocation.Qty,
			UnitCostOriginal: round6(allocation.UnitCost),
			OriginalCurrency: originalCurrency,
			OriginalAmount:   originalAmount,
			BaseCurrency:     baseCurrency,
			FxRate:           fxSnapshot.Rate,
			BaseAmount:       baseAmount,
			FxSource:         fxSnapshot.Source,
			FxVersion:        fxSnapshot.Version,
			FxTime:           fxSnapshot.EffectiveAt,
			OccurredAt:       occurredAt,
			OperatorID:       params.OperatorID,
			Remark:           &remark,
		})
	}
	return w.repo.CreateBatch(details)
}

func (w *OrderCostWriter) ResolveSalesReturnUnitCost(params *salesUsecase.SalesReturnCostRecordParams) (*float64, error) {
	if w.repo == nil || params == nil || params.SalesOrderItemID == 0 || params.QtyReturned == 0 {
		return nil, nil
	}

	rows, err := w.repo.ListReturnableBySalesOrderItemID(params.SalesOrderItemID)
	if err != nil {
		return nil, err
	}

	remaining := params.QtyReturned
	totalQty := uint64(0)
	totalOriginalAmount := 0.0
	for _, row := range rows {
		if remaining == 0 {
			break
		}
		availableQty := row.AvailableQty
		if availableQty == 0 {
			continue
		}
		reverseQty := availableQty
		if reverseQty > remaining {
			reverseQty = remaining
		}
		remaining -= reverseQty
		totalQty += reverseQty
		totalOriginalAmount += float64(reverseQty) * row.UnitCostOriginal
	}

	if totalQty == 0 {
		return nil, nil
	}
	unitCost := round6(totalOriginalAmount / float64(totalQty))
	return &unitCost, nil
}

func (w *OrderCostWriter) RecordSalesReturnCost(params *salesUsecase.SalesReturnCostRecordParams) (float64, error) {
	if w.repo == nil || params == nil || params.SalesOrderID == 0 || params.SalesOrderItemID == 0 || params.QtyReturned == 0 {
		return 0, nil
	}

	occurredAt := params.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}

	rows, err := w.repo.ListReturnableBySalesOrderItemID(params.SalesOrderItemID)
	if err != nil {
		return 0, err
	}

	remaining := params.QtyReturned
	remark := "auto from sales return reversal"
	details := make([]domain.OrderCostDetail, 0)
	totalCOGS := 0.0
	for _, row := range rows {
		if remaining == 0 {
			break
		}
		availableQty := row.AvailableQty
		if availableQty == 0 {
			continue
		}
		reverseQty := availableQty
		if reverseQty > remaining {
			reverseQty = remaining
		}
		remaining -= reverseQty

		originalAmount := round6(float64(reverseQty) * row.UnitCostOriginal)
		baseAmount := round6(originalAmount * row.FxRate)
		reversalOfID := row.ID
		details = append(details, domain.OrderCostDetail{
			TraceID:          fmt.Sprintf("ORDER-COST-RETURN-%d-%d", time.Now().UnixNano(), reversalOfID),
			Status:           domain.OrderCostDetailStatusNormal,
			ReversalOfID:     &reversalOfID,
			SalesOrderID:     params.SalesOrderID,
			SalesOrderItemID: params.SalesOrderItemID,
			ProductID:        params.ProductID,
			WarehouseID:      params.WarehouseID,
			Marketplace:      row.Marketplace,
			InventoryLotID:   row.InventoryLotID,
			QtyOut:           reverseQty,
			UnitCostOriginal: row.UnitCostOriginal,
			OriginalCurrency: row.OriginalCurrency,
			OriginalAmount:   originalAmount,
			BaseCurrency:     row.BaseCurrency,
			FxRate:           row.FxRate,
			BaseAmount:       baseAmount,
			FxSource:         row.FxSource,
			FxVersion:        row.FxVersion,
			FxTime:           occurredAt,
			OccurredAt:       occurredAt,
			OperatorID:       params.OperatorID,
			Remark:           &remark,
		})
		totalCOGS += baseAmount
	}

	if len(details) == 0 {
		return 0, nil
	}
	if err := w.repo.CreateBatch(details); err != nil {
		return 0, err
	}
	return round6(totalCOGS), nil
}
