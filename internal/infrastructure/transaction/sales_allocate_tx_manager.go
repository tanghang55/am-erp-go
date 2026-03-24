package transaction

import (
	"context"

	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	salesRepo "am-erp-go/internal/module/sales/repository"
	salesUsecase "am-erp-go/internal/module/sales/usecase"

	"gorm.io/gorm"
)

type SalesAllocateTxManager struct {
	db                  *gorm.DB
	seedLotCostResolver inventoryUsecase.SeedLotUnitCostResolver
}

func NewSalesAllocateTxManager(
	db *gorm.DB,
	seedLotCostResolver inventoryUsecase.SeedLotUnitCostResolver,
) *SalesAllocateTxManager {
	return &SalesAllocateTxManager{
		db:                  db,
		seedLotCostResolver: seedLotCostResolver,
	}
}

func (m *SalesAllocateTxManager) Run(
	ctx context.Context,
	fn func(salesUsecase.SalesAllocateTransactionalDeps) error,
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
		deps := salesUsecase.SalesAllocateTransactionalDeps{
			Repo:             salesRepo.NewSalesOrderRepository(tx),
			InventoryService: inventoryUC,
		}
		return fn(deps)
	})
}
