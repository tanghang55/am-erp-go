package transaction

import (
	"context"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	packagingRepo "am-erp-go/internal/module/packaging/repository"

	"gorm.io/gorm"
)

type AssemblyTxManager struct {
	db *gorm.DB
}

func NewAssemblyTxManager(db *gorm.DB) *AssemblyTxManager {
	return &AssemblyTxManager{db: db}
}

func (m *AssemblyTxManager) Run(
	ctx context.Context,
	fn func(inventoryUsecase.AssemblyTransactionalDeps) (*inventoryDomain.InventoryMovement, error),
) (*inventoryDomain.InventoryMovement, error) {
	var movement *inventoryDomain.InventoryMovement
	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		deps := inventoryUsecase.AssemblyTransactionalDeps{
			BalanceRepo:         inventoryRepo.NewInventoryBalanceRepository(tx),
			MovementRepo:        inventoryRepo.NewInventoryMovementRepository(tx),
			LotRepo:             inventoryRepo.NewInventoryLotRepository(tx),
			PackagingItemRepo:   packagingRepo.NewPackagingItemRepository(tx),
			PackagingLedgerRepo: packagingRepo.NewPackagingLedgerRepository(tx),
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
