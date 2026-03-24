package transaction

import (
	"context"

	financeRepo "am-erp-go/internal/module/finance/repository"
	financeUsecase "am-erp-go/internal/module/finance/usecase"
	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	salesRepo "am-erp-go/internal/module/sales/repository"
	salesUsecase "am-erp-go/internal/module/sales/usecase"

	"gorm.io/gorm"
)

type SalesShipTxManager struct {
	db                  *gorm.DB
	seedLotCostResolver inventoryUsecase.SeedLotUnitCostResolver
}

func NewSalesShipTxManager(
	db *gorm.DB,
	seedLotCostResolver inventoryUsecase.SeedLotUnitCostResolver,
) *SalesShipTxManager {
	return &SalesShipTxManager{
		db:                  db,
		seedLotCostResolver: seedLotCostResolver,
	}
}

func (m *SalesShipTxManager) Run(
	ctx context.Context,
	fn func(salesUsecase.SalesShipTransactionalDeps) error,
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
		deps := salesUsecase.SalesShipTransactionalDeps{
			Repo:             salesRepo.NewSalesOrderRepository(tx),
			InventoryService: inventoryUC,
			CostWriter:       financeUsecase.NewOrderCostWriter(financeRepo.NewOrderCostDetailRepository(tx)),
			ProfitWriter:     financeUsecase.NewProfitLedgerWriter(financeRepo.NewProfitLedgerRepository(tx)),
		}
		return fn(deps)
	})
}
