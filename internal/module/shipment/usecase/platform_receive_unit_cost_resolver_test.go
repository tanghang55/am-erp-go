package usecase

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/shipment/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type stubPlatformReceiveBaseCurrencyProvider struct{}

type stubShipmentAuditLogger struct {
	payloads []systemUsecase.AuditLogPayload
}

type stubPackingMaterialUnitCostReader struct {
	getLatestPackingMaterialPerUnit func(productID uint64, occurredAt time.Time) (*float64, error)
}

func (s *stubShipmentAuditLogger) RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error {
	s.payloads = append(s.payloads, payload)
	return nil
}

func (s *stubPackingMaterialUnitCostReader) GetLatestPackingMaterialPerUnit(productID uint64, occurredAt time.Time) (*float64, error) {
	if s == nil || s.getLatestPackingMaterialPerUnit == nil {
		return nil, nil
	}
	return s.getLatestPackingMaterialPerUnit(productID, occurredAt)
}

func (stubPlatformReceiveBaseCurrencyProvider) GetDefaultBaseCurrency() string {
	return "USD"
}

func TestShipmentPlatformReceiveUnitCostResolverResolveByShipmentID(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			shippedAt := time.Date(2026, 3, 9, 10, 0, 0, 0, time.Local)
			return &domain.Shipment{
				ID:                     id,
				ShipmentNumber:         "SHP-0088",
				Status:                 domain.ShipmentStatusShipped,
				Currency:               "USD",
				BaseCurrency:           "USD",
				ShippingCost:           12,
				ShippingCostBaseAmount: 12,
				ShippingCostFxRate:     1,
				ShippingCostFxSource:   "MANUAL",
				ShippingCostFxVersion:  "v1",
				ShippingCostFxTime:     &shippedAt,
				ShippedAt:              &shippedAt,
			}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 1, ShipmentID: shipmentID, ProductID: 1001, QuantityPlanned: 4, QuantityShipped: 4, UnitCost: 5, Currency: "USD"},
				{ID: 2, ShipmentID: shipmentID, ProductID: 2002, QuantityPlanned: 2, QuantityShipped: 2, UnitCost: 3, Currency: "USD"},
			}, nil
		},
	}
	resolver := NewShipmentPlatformReceiveUnitCostResolver(
		shipmentRepo,
		itemRepo,
		stubPlatformReceiveBaseCurrencyProvider{},
		func(baseCurrency, originalCurrency string, occurredAt time.Time) (*ShipmentFXSnapshot, error) {
			return &ShipmentFXSnapshot{
				Rate:        1,
				Source:      "MANUAL",
				Version:     "v1",
				EffectiveAt: occurredAt,
			}, nil
		},
		nil,
	)
	referenceType := "SHIPMENT"
	referenceID := uint64(88)

	cost, err := resolver.ResolveByShipmentReference(context.Background(), 1001, &referenceType, &referenceID, nil, time.Date(2026, 3, 9, 12, 0, 0, 0, time.Local))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cost == nil {
		t.Fatalf("expected resolved unit cost")
	}
	expected := 7.0 // 5 purchase + 12 freight / 6 total qty
	if *cost != expected {
		t.Fatalf("expected landed unit cost %.4f, got %.4f", expected, *cost)
	}
}

func TestShipmentPlatformReceiveUnitCostResolverIncludesPackingMaterialCost(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			shippedAt := time.Date(2026, 3, 9, 10, 0, 0, 0, time.Local)
			return &domain.Shipment{
				ID:                     id,
				ShipmentNumber:         "SHP-0089",
				Status:                 domain.ShipmentStatusShipped,
				Currency:               "USD",
				BaseCurrency:           "USD",
				ShippingCost:           12,
				ShippingCostBaseAmount: 12,
				ShippingCostFxRate:     1,
				ShippingCostFxSource:   "MANUAL",
				ShippingCostFxVersion:  "v1",
				ShippingCostFxTime:     &shippedAt,
				ShippedAt:              &shippedAt,
			}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 1, ShipmentID: shipmentID, ProductID: 1001, QuantityPlanned: 4, QuantityShipped: 4, UnitCost: 5, Currency: "USD"},
				{ID: 2, ShipmentID: shipmentID, ProductID: 2002, QuantityPlanned: 2, QuantityShipped: 2, UnitCost: 3, Currency: "USD"},
			}, nil
		},
	}
	packingReader := &stubPackingMaterialUnitCostReader{
		getLatestPackingMaterialPerUnit: func(productID uint64, occurredAt time.Time) (*float64, error) {
			if productID != 1001 {
				return nil, nil
			}
			value := 1.5
			return &value, nil
		},
	}
	resolver := NewShipmentPlatformReceiveUnitCostResolver(
		shipmentRepo,
		itemRepo,
		stubPlatformReceiveBaseCurrencyProvider{},
		func(baseCurrency, originalCurrency string, occurredAt time.Time) (*ShipmentFXSnapshot, error) {
			return &ShipmentFXSnapshot{
				Rate:        1,
				Source:      "MANUAL",
				Version:     "v1",
				EffectiveAt: occurredAt,
			}, nil
		},
		packingReader,
	)
	referenceType := "SHIPMENT"
	referenceID := uint64(89)

	cost, err := resolver.ResolveByShipmentReference(context.Background(), 1001, &referenceType, &referenceID, nil, time.Date(2026, 3, 9, 12, 0, 0, 0, time.Local))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cost == nil {
		t.Fatalf("expected resolved unit cost")
	}
	expected := 8.5 // 5 purchase + 1.5 packing + 12 freight / 6 total qty
	if *cost != expected {
		t.Fatalf("expected landed unit cost %.4f, got %.4f", expected, *cost)
	}
}

func TestShipmentPlatformReceiveRecorderUpdatesQuantityReceived(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, ShipmentNumber: "SHP-0099", ReceiptStatus: domain.ShipmentReceiptStatusPending}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 1, ShipmentID: shipmentID, ProductID: 1001, QuantityShipped: 2, QuantityReceived: 0},
				{ID: 2, ShipmentID: shipmentID, ProductID: 1001, QuantityShipped: 1, QuantityReceived: 0},
			}, nil
		},
	}
	recorder := NewShipmentPlatformReceiveRecorder(shipmentRepo, itemRepo)
	refType := "SHIPMENT"
	refID := uint64(99)
	operatorID := uint64(77)

	if err := recorder.RecordPlatformReceive(context.Background(), &inventoryDomain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   1,
		MovementType:  inventoryDomain.MovementTypePlatformReceive,
		Quantity:      3,
		ReferenceType: &refType,
		ReferenceID:   &refID,
		OperatorID:    &operatorID,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(itemRepo.updated) != 2 {
		t.Fatalf("expected updated items, got %+v", itemRepo.updated)
	}
	if itemRepo.updated[0].QuantityReceived != 2 || itemRepo.updated[1].QuantityReceived != 1 {
		t.Fatalf("unexpected received quantities: %+v", itemRepo.updated)
	}
	if shipmentRepo.updated == nil || shipmentRepo.updated.ReceiptStatus != domain.ShipmentReceiptStatusCompleted {
		t.Fatalf("expected shipment receipt status completed, got %+v", shipmentRepo.updated)
	}
	if shipmentRepo.updated.ReceiptCompletedAt == nil {
		t.Fatalf("expected receipt completed time to be set")
	}
	if shipmentRepo.updated.ReceiptCompletedBy == nil || *shipmentRepo.updated.ReceiptCompletedBy != operatorID {
		t.Fatalf("expected receipt completed by %d, got %+v", operatorID, shipmentRepo.updated.ReceiptCompletedBy)
	}
}

func TestShipmentPlatformReceiveRecorderRejectsOverReceive(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, ShipmentNumber: "SHP-0100", ReceiptStatus: domain.ShipmentReceiptStatusPending}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 1, ShipmentID: shipmentID, ProductID: 1001, QuantityShipped: 1, QuantityReceived: 1},
			}, nil
		},
	}
	recorder := NewShipmentPlatformReceiveRecorder(shipmentRepo, itemRepo)
	refType := "SHIPMENT"
	refID := uint64(100)

	err := recorder.RecordPlatformReceive(context.Background(), &inventoryDomain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   1,
		MovementType:  inventoryDomain.MovementTypePlatformReceive,
		Quantity:      1,
		ReferenceType: &refType,
		ReferenceID:   &refID,
	})
	if err == nil {
		t.Fatalf("expected over-receive error")
	}
}

func TestShipmentPlatformReceiveRecorderMarksPartialReceiptStatus(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, ShipmentNumber: "SHP-0101", ReceiptStatus: domain.ShipmentReceiptStatusPending}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 1, ShipmentID: shipmentID, ProductID: 1001, QuantityShipped: 2, QuantityReceived: 0},
				{ID: 2, ShipmentID: shipmentID, ProductID: 1001, QuantityShipped: 2, QuantityReceived: 0},
			}, nil
		},
	}
	recorder := NewShipmentPlatformReceiveRecorder(shipmentRepo, itemRepo)
	refType := "SHIPMENT"
	refID := uint64(101)

	if err := recorder.RecordPlatformReceive(context.Background(), &inventoryDomain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   1,
		MovementType:  inventoryDomain.MovementTypePlatformReceive,
		Quantity:      1,
		ReferenceType: &refType,
		ReferenceID:   &refID,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shipmentRepo.updated == nil || shipmentRepo.updated.ReceiptStatus != domain.ShipmentReceiptStatusPartial {
		t.Fatalf("expected shipment receipt status partial, got %+v", shipmentRepo.updated)
	}
	if shipmentRepo.updated.ReceiptCompletedAt != nil {
		t.Fatalf("expected no completion time on partial receipt, got %+v", shipmentRepo.updated)
	}
	if shipmentRepo.updated.ReceiptCompletedBy != nil {
		t.Fatalf("expected no completion operator on partial receipt, got %+v", shipmentRepo.updated)
	}
}

func TestShipmentPlatformReceiveRecorderRecordsShipmentAudit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	userID := uint64(8)
	ctx.Set("userID", userID)
	ctx.Set("username", "admin")

	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{ID: id, ShipmentNumber: "SHP-0200", ReceiptStatus: domain.ShipmentReceiptStatusPending}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 1, ShipmentID: shipmentID, ProductID: 1001, QuantityShipped: 2, QuantityReceived: 0},
				{ID: 2, ShipmentID: shipmentID, ProductID: 2002, QuantityShipped: 1, QuantityReceived: 0},
			}, nil
		},
	}
	auditLogger := &stubShipmentAuditLogger{}
	recorder := NewShipmentPlatformReceiveRecorder(shipmentRepo, itemRepo, auditLogger)
	refType := "SHIPMENT"
	refID := uint64(200)
	operatorID := uint64(77)

	if err := recorder.RecordPlatformReceive(ctx, &inventoryDomain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   1,
		MovementType:  inventoryDomain.MovementTypePlatformReceive,
		Quantity:      2,
		ReferenceType: &refType,
		ReferenceID:   &refID,
		OperatorID:    &operatorID,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(auditLogger.payloads) != 1 {
		t.Fatalf("expected one audit payload, got %d", len(auditLogger.payloads))
	}
	payload := auditLogger.payloads[0]
	if payload.Module != "Shipment" || payload.Action != "RECEIVE" || payload.EntityType != "Shipment" || payload.EntityID != "200" {
		t.Fatalf("unexpected audit payload meta: %+v", payload)
	}
	before, ok := payload.Before.(map[string]any)
	if !ok {
		t.Fatalf("expected before snapshot map, got %#v", payload.Before)
	}
	after, ok := payload.After.(map[string]any)
	if !ok {
		t.Fatalf("expected after snapshot map, got %#v", payload.After)
	}
	if before["receipt_status"] != domain.ShipmentReceiptStatusPending || after["receipt_status"] != domain.ShipmentReceiptStatusPartial {
		t.Fatalf("expected receipt status change, got before=%v after=%v", before["receipt_status"], after["receipt_status"])
	}
	if before["received_quantity_total"] != uint64(0) || after["received_quantity_total"] != uint64(2) {
		t.Fatalf("expected received totals 0->2, got before=%v after=%v", before["received_quantity_total"], after["received_quantity_total"])
	}
	if before["remaining_quantity_total"] != uint64(3) || after["remaining_quantity_total"] != uint64(1) {
		t.Fatalf("expected remaining totals 3->1, got before=%v after=%v", before["remaining_quantity_total"], after["remaining_quantity_total"])
	}
}
