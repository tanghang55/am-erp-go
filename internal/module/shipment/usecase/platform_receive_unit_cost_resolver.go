package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/shipment/domain"
)

type ShipmentPlatformReceiveUnitCostResolver struct {
	shipmentRepo     domain.ShipmentRepository
	shipmentItemRepo domain.ShipmentItemRepository
	defaultsProvider ShipmentDefaultsProvider
	fxResolver       ShipmentFXResolver
	packingCostRepo  ShipmentPackingUnitCostReader
}

type ShipmentPackingUnitCostReader interface {
	GetLatestPackingMaterialPerUnit(productID uint64, occurredAt time.Time) (*float64, error)
}

func NewShipmentPlatformReceiveUnitCostResolver(
	shipmentRepo domain.ShipmentRepository,
	shipmentItemRepo domain.ShipmentItemRepository,
	defaultsProvider ShipmentDefaultsProvider,
	fxResolver ShipmentFXResolver,
	packingCostRepo ShipmentPackingUnitCostReader,
) *ShipmentPlatformReceiveUnitCostResolver {
	return &ShipmentPlatformReceiveUnitCostResolver{
		shipmentRepo:     shipmentRepo,
		shipmentItemRepo: shipmentItemRepo,
		defaultsProvider: defaultsProvider,
		fxResolver:       fxResolver,
		packingCostRepo:  packingCostRepo,
	}
}

func (r *ShipmentPlatformReceiveUnitCostResolver) ResolvePlatformReceiveUnitCost(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*float64, error) {
	if r == nil || params == nil || params.ProductID == 0 {
		return nil, nil
	}
	return r.ResolveByShipmentReference(ctx, params.ProductID, params.ReferenceType, params.ReferenceID, params.ReferenceNumber, timeValue(params.OperatedAt, time.Now()))
}

func (r *ShipmentPlatformReceiveUnitCostResolver) ResolveByShipmentReference(
	ctx context.Context,
	productID uint64,
	referenceType *string,
	referenceID *uint64,
	referenceNumber *string,
	occurredAt time.Time,
) (*float64, error) {
	if !isShipmentReference(referenceType) {
		return nil, nil
	}
	shipment, err := r.findShipment(referenceID, referenceNumber)
	if err != nil {
		return nil, err
	}
	if shipment == nil {
		return nil, nil
	}
	items, err := r.shipmentItemRepo.GetByShipmentID(shipment.ID)
	if err != nil {
		return nil, err
	}

	totalQty := uint64(0)
	targetQty := uint64(0)
	totalTargetBaseAmount := 0.0
	perUnitFreightBase := 0.0
	packingUnitCost := 0.0
	packingCostEffectiveAt := occurredAt
	if shipment.ShippedAt != nil && !shipment.ShippedAt.IsZero() {
		packingCostEffectiveAt = *shipment.ShippedAt
	}
	if shipment.ShippingCostBaseAmount > 0 {
		for _, item := range items {
			qty := shippedQty(item)
			if qty == 0 {
				continue
			}
			totalQty += qty
		}
		if totalQty > 0 {
			perUnitFreightBase = roundAmount6(shipment.ShippingCostBaseAmount / float64(totalQty))
		}
	}
	if r.packingCostRepo != nil {
		value, err := r.packingCostRepo.GetLatestPackingMaterialPerUnit(productID, packingCostEffectiveAt)
		if err != nil {
			return nil, err
		}
		if value != nil {
			packingUnitCost = roundAmount6(*value)
		}
	}

	for _, item := range items {
		if item.ProductID != productID {
			continue
		}
		qty := shippedQty(item)
		if qty == 0 {
			continue
		}
		sourceUnitBase, err := r.toBaseUnitCost(item.Currency, item.UnitCost, shipment, occurredAt)
		if err != nil {
			return nil, err
		}
		totalTargetBaseAmount += roundAmount6((sourceUnitBase + packingUnitCost + perUnitFreightBase) * float64(qty))
		targetQty += qty
	}
	if targetQty == 0 {
		return nil, nil
	}
	cost := roundAmount6(totalTargetBaseAmount / float64(targetQty))
	return &cost, nil
}

func (r *ShipmentPlatformReceiveUnitCostResolver) findShipment(referenceID *uint64, referenceNumber *string) (*domain.Shipment, error) {
	if r == nil || r.shipmentRepo == nil {
		return nil, nil
	}
	if referenceID != nil && *referenceID > 0 {
		return r.shipmentRepo.GetByID(*referenceID)
	}
	if referenceNumber != nil && strings.TrimSpace(*referenceNumber) != "" {
		return r.shipmentRepo.GetByShipmentNumber(strings.TrimSpace(*referenceNumber))
	}
	return nil, nil
}

func (r *ShipmentPlatformReceiveUnitCostResolver) toBaseUnitCost(currency string, unitCost float64, shipment *domain.Shipment, occurredAt time.Time) (float64, error) {
	baseCurrency := defaultShipmentBaseCurrency(r.defaultsProvider)
	sourceCurrency := strings.TrimSpace(currency)
	if sourceCurrency == "" {
		sourceCurrency = baseCurrency
	}
	if strings.EqualFold(baseCurrency, sourceCurrency) {
		return roundAmount6(unitCost), nil
	}
	if r.fxResolver == nil {
		return 0, fmt.Errorf("shipment fx resolver not configured")
	}
	effectiveAt := occurredAt
	if shipment != nil && shipment.ShippedAt != nil && !shipment.ShippedAt.IsZero() {
		effectiveAt = *shipment.ShippedAt
	}
	snapshot, err := r.fxResolver(baseCurrency, sourceCurrency, effectiveAt)
	if err != nil {
		return 0, err
	}
	return roundAmount6(unitCost * snapshot.Rate), nil
}

func shippedQty(item domain.ShipmentItem) uint64 {
	if item.QuantityShipped > 0 {
		return uint64(item.QuantityShipped)
	}
	return uint64(item.QuantityPlanned)
}

func isShipmentReference(referenceType *string) bool {
	return referenceType != nil && strings.EqualFold(strings.TrimSpace(*referenceType), "SHIPMENT")
}

func defaultShipmentBaseCurrency(provider ShipmentDefaultsProvider) string {
	if provider != nil {
		if currency := strings.TrimSpace(provider.GetDefaultBaseCurrency()); currency != "" {
			return currency
		}
	}
	return "USD"
}
