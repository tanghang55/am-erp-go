package transaction

import (
	"context"

	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	shipmentRepo "am-erp-go/internal/module/shipment/repository"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"

	"gorm.io/gorm"
)

type ShipmentStateTxManager struct {
	db                  *gorm.DB
	seedLotCostResolver inventoryUsecase.SeedLotUnitCostResolver
}

func NewShipmentStateTxManager(
	db *gorm.DB,
	seedLotCostResolver inventoryUsecase.SeedLotUnitCostResolver,
) *ShipmentStateTxManager {
	return &ShipmentStateTxManager{
		db:                  db,
		seedLotCostResolver: seedLotCostResolver,
	}
}

func (m *ShipmentStateTxManager) Run(
	ctx context.Context,
	fn func(shipmentUsecase.ShipmentStateTransactionalDeps) error,
) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		inventoryUC := inventoryUsecase.NewInventoryUsecase(
			inventoryRepo.NewInventoryBalanceRepository(tx),
			inventoryRepo.NewInventoryMovementRepository(tx),
			inventoryRepo.NewInventoryLotRepository(tx),
		)
		if m.seedLotCostResolver != nil {
			inventoryUC.BindSeedLotUnitCostResolver(m.seedLotCostResolver)
		}
		deps := shipmentUsecase.ShipmentStateTransactionalDeps{
			ShipmentRepo:     shipmentRepo.NewShipmentRepo(tx),
			ShipmentItemRepo: shipmentRepo.NewShipmentItemRepo(tx),
			InventoryService: inventoryUC,
		}
		return fn(deps)
	})
}
