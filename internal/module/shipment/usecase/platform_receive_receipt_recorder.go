package usecase

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/shipment/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type ShipmentPlatformReceiveRecorder struct {
	shipmentRepo     domain.ShipmentRepository
	shipmentItemRepo domain.ShipmentItemRepository
	auditLogger      AuditLogger
}

func NewShipmentPlatformReceiveRecorder(
	shipmentRepo domain.ShipmentRepository,
	shipmentItemRepo domain.ShipmentItemRepository,
	auditLogger ...AuditLogger,
) *ShipmentPlatformReceiveRecorder {
	var logger AuditLogger
	if len(auditLogger) > 0 {
		logger = auditLogger[0]
	}
	return &ShipmentPlatformReceiveRecorder{
		shipmentRepo:     shipmentRepo,
		shipmentItemRepo: shipmentItemRepo,
		auditLogger:      logger,
	}
}

func (r *ShipmentPlatformReceiveRecorder) ValidatePlatformReceive(ctx context.Context, params *inventoryDomain.CreateMovementParams) error {
	_, _, _, err := r.planPlatformReceive(ctx, params)
	return err
}

func (r *ShipmentPlatformReceiveRecorder) RecordPlatformReceive(ctx context.Context, params *inventoryDomain.CreateMovementParams) error {
	shipment, items, beforeSnapshot, err := r.planPlatformReceive(ctx, params)
	if err != nil {
		return err
	}
	if shipment == nil {
		return nil
	}
	if err := r.shipmentItemRepo.UpdateBatch(items); err != nil {
		return err
	}
	r.syncShipmentReceiptStatus(shipment, items, params)
	if err := r.shipmentRepo.Update(shipment); err != nil {
		return err
	}
	afterSnapshot := buildShipmentReceiptAuditSnapshot(shipment, items)
	r.recordPlatformReceiveAudit(ctx, shipment.ID, beforeSnapshot, afterSnapshot)
	return nil
}

func (r *ShipmentPlatformReceiveRecorder) planPlatformReceive(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*domain.Shipment, []domain.ShipmentItem, map[string]any, error) {
	if r == nil || params == nil || params.Quantity <= 0 || params.ProductID == 0 {
		return nil, nil, nil, nil
	}
	shipment, err := r.findShipment(params.ReferenceID, params.ReferenceNumber)
	if err != nil {
		return nil, nil, nil, err
	}
	if shipment == nil {
		return nil, nil, nil, nil
	}
	items, err := r.shipmentItemRepo.GetByShipmentID(shipment.ID)
	if err != nil {
		return nil, nil, nil, err
	}
	beforeSnapshot := buildShipmentReceiptAuditSnapshot(shipment, items)
	remaining := uint(params.Quantity)
	matched := false
	for i := range items {
		if items[i].ProductID != params.ProductID {
			continue
		}
		matched = true
		receivable := receivableQty(items[i])
		if receivable == 0 {
			continue
		}
		take := minShipmentReceiveUint(remaining, uint(receivable))
		items[i].QuantityReceived += take
		remaining -= take
		if remaining == 0 {
			break
		}
	}
	if !matched {
		return nil, nil, nil, fmt.Errorf("货件中不存在产品 %d", params.ProductID)
	}
	if remaining > 0 {
		return nil, nil, nil, fmt.Errorf("平台上架数量超过货件未接收数量: 产品 %d 超出 %d", params.ProductID, remaining)
	}
	return shipment, items, beforeSnapshot, nil
}

func (r *ShipmentPlatformReceiveRecorder) findShipment(referenceID *uint64, referenceNumber *string) (*domain.Shipment, error) {
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

func receivableQty(item domain.ShipmentItem) uint64 {
	shipped := shippedQty(item)
	if shipped <= uint64(item.QuantityReceived) {
		return 0
	}
	return shipped - uint64(item.QuantityReceived)
}

func minShipmentReceiveUint(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func (r *ShipmentPlatformReceiveRecorder) syncShipmentReceiptStatus(
	shipment *domain.Shipment,
	items []domain.ShipmentItem,
	params *inventoryDomain.CreateMovementParams,
) {
	if shipment == nil {
		return
	}
	totalShipped := uint64(0)
	totalReceived := uint64(0)
	for _, item := range items {
		totalShipped += shippedQty(item)
		totalReceived += uint64(item.QuantityReceived)
	}
	switch {
	case totalShipped > 0 && totalReceived >= totalShipped:
		shipment.ReceiptStatus = domain.ShipmentReceiptStatusCompleted
		now := time.Now()
		shipment.ReceiptCompletedAt = &now
		shipment.ReceiptCompletedBy = params.OperatorID
	case totalReceived > 0:
		shipment.ReceiptStatus = domain.ShipmentReceiptStatusPartial
		shipment.ReceiptCompletedAt = nil
		shipment.ReceiptCompletedBy = nil
	default:
		shipment.ReceiptStatus = domain.ShipmentReceiptStatusPending
		shipment.ReceiptCompletedAt = nil
		shipment.ReceiptCompletedBy = nil
	}
}

func buildShipmentReceiptAuditSnapshot(shipment *domain.Shipment, items []domain.ShipmentItem) map[string]any {
	snapshot := map[string]any{}
	if shipment == nil {
		return snapshot
	}
	totalShipped := uint64(0)
	totalReceived := uint64(0)
	for _, item := range items {
		totalShipped += shippedQty(item)
		totalReceived += uint64(item.QuantityReceived)
	}
	snapshot["receipt_status"] = shipment.ReceiptStatus
	snapshot["received_quantity_total"] = totalReceived
	snapshot["remaining_quantity_total"] = totalShipped - totalReceived
	if shipment.ReceiptCompletedAt != nil {
		snapshot["receipt_completed_at"] = shipment.ReceiptCompletedAt.Format(time.RFC3339)
	}
	return snapshot
}

func (r *ShipmentPlatformReceiveRecorder) recordPlatformReceiveAudit(
	ctx context.Context,
	shipmentID uint64,
	before map[string]any,
	after map[string]any,
) {
	if r == nil || r.auditLogger == nil || shipmentID == 0 {
		return
	}
	ginCtx, ok := ctx.(*gin.Context)
	if !ok || ginCtx == nil {
		return
	}
	beforeDiff, afterDiff := diffShipmentReceiptAudit(before, after)
	if len(beforeDiff) == 0 && len(afterDiff) == 0 {
		return
	}
	_ = r.auditLogger.RecordFromContext(ginCtx, systemUsecase.AuditLogPayload{
		Module:     "Shipment",
		Action:     "RECEIVE",
		EntityType: "Shipment",
		EntityID:   fmt.Sprintf("%d", shipmentID),
		Before:     beforeDiff,
		After:      afterDiff,
	})
}

func diffShipmentReceiptAudit(before map[string]any, after map[string]any) (map[string]any, map[string]any) {
	beforeDiff := map[string]any{}
	afterDiff := map[string]any{}
	keys := map[string]struct{}{}
	for key := range before {
		keys[key] = struct{}{}
	}
	for key := range after {
		keys[key] = struct{}{}
	}
	for key := range keys {
		beforeValue, beforeOk := before[key]
		afterValue, afterOk := after[key]
		if beforeOk && afterOk && reflect.DeepEqual(beforeValue, afterValue) {
			continue
		}
		if beforeOk {
			beforeDiff[key] = beforeValue
		}
		if afterOk {
			afterDiff[key] = afterValue
		}
	}
	return beforeDiff, afterDiff
}
