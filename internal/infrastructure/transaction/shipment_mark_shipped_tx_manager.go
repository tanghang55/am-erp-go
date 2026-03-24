package transaction

import (
	"context"

	financeRepo "am-erp-go/internal/module/finance/repository"
	financeUsecase "am-erp-go/internal/module/finance/usecase"
	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	shipmentRepo "am-erp-go/internal/module/shipment/repository"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"

	"gorm.io/gorm"
)

type ShipmentMarkShippedTxManager struct {
	db                  *gorm.DB
	seedLotCostResolver inventoryUsecase.SeedLotUnitCostResolver
}

func NewShipmentMarkShippedTxManager(
	db *gorm.DB,
	seedLotCostResolver inventoryUsecase.SeedLotUnitCostResolver,
) *ShipmentMarkShippedTxManager {
	return &ShipmentMarkShippedTxManager{
		db:                  db,
		seedLotCostResolver: seedLotCostResolver,
	}
}

func (m *ShipmentMarkShippedTxManager) Run(
	ctx context.Context,
	fn func(shipmentUsecase.ShipmentMarkShippedTransactionalDeps) error,
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
		deps := shipmentUsecase.ShipmentMarkShippedTransactionalDeps{
			ShipmentRepo:     shipmentRepo.NewShipmentRepo(tx),
			ShipmentItemRepo: shipmentRepo.NewShipmentItemRepo(tx),
			InventoryService: inventoryUC,
			CostRecorder:     financeUsecase.NewCostEventWriter(financeRepo.NewCostEventRepository(tx)),
			LandedRecorder:   financeUsecase.NewShipmentLandedSnapshotWriter(financeRepo.NewCostingSnapshotRepository(tx), financeRepo.NewCostEventRepository(tx)),
		}
		return fn(deps)
	})
}
