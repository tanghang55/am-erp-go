package transaction

import (
	"context"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	shipmentRepo "am-erp-go/internal/module/shipment/repository"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"

	"gorm.io/gorm"
)

type PlatformReceiveTxManager struct {
	db          *gorm.DB
	auditLogger shipmentUsecase.AuditLogger
}

func NewPlatformReceiveTxManager(db *gorm.DB, auditLogger ...shipmentUsecase.AuditLogger) *PlatformReceiveTxManager {
	var logger shipmentUsecase.AuditLogger
	if len(auditLogger) > 0 {
		logger = auditLogger[0]
	}
	return &PlatformReceiveTxManager{db: db, auditLogger: logger}
}

func (m *PlatformReceiveTxManager) Run(
	ctx context.Context,
	fn func(inventoryUsecase.PlatformReceiveTransactionalDeps) (*inventoryDomain.InventoryMovement, error),
) (*inventoryDomain.InventoryMovement, error) {
	var movement *inventoryDomain.InventoryMovement
	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		deps := inventoryUsecase.PlatformReceiveTransactionalDeps{
			BalanceRepo:  inventoryRepo.NewInventoryBalanceRepository(tx),
			MovementRepo: inventoryRepo.NewInventoryMovementRepository(tx),
			LotRepo:      inventoryRepo.NewInventoryLotRepository(tx),
			Recorder: shipmentUsecase.NewShipmentPlatformReceiveRecorder(
				shipmentRepo.NewShipmentRepo(tx),
				shipmentRepo.NewShipmentItemRepo(tx),
				m.auditLogger,
			),
		}
		result, err := fn(deps)
		if err != nil {
			return err
		}
		movement = result
		return nil
	})
	if err != nil {
		return nil, err
	}
	return movement, nil
}
