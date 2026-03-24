package transaction

import (
	"context"

	financeRepo "am-erp-go/internal/module/finance/repository"
	financeUsecase "am-erp-go/internal/module/finance/usecase"
	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	procurementDomain "am-erp-go/internal/module/procurement/domain"
	procurementRepo "am-erp-go/internal/module/procurement/repository"
	procurementUsecase "am-erp-go/internal/module/procurement/usecase"

	"gorm.io/gorm"
)

type PurchaseOrderShipTxManager struct {
	db *gorm.DB
}

func NewPurchaseOrderShipTxManager(db *gorm.DB) *PurchaseOrderShipTxManager {
	return &PurchaseOrderShipTxManager{db: db}
}

func (m *PurchaseOrderShipTxManager) Run(
	ctx context.Context,
	fn func(procurementUsecase.PurchaseOrderShipTransactionalDeps) error,
) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		deps := procurementUsecase.PurchaseOrderShipTransactionalDeps{
			Repo: procurementRepo.NewPurchaseOrderRepository(tx),
			InventoryService: inventoryUsecase.NewInventoryUsecase(
				inventoryRepo.NewInventoryBalanceRepository(tx),
				inventoryRepo.NewInventoryMovementRepository(tx),
				inventoryRepo.NewInventoryLotRepository(tx),
			),
			CostEventRecorder: financeUsecase.NewCostEventWriter(financeRepo.NewCostEventRepository(tx)),
		}
		return fn(deps)
	})
}

type PurchaseOrderReceiveTxManager struct {
	db *gorm.DB
}

func NewPurchaseOrderReceiveTxManager(db *gorm.DB) *PurchaseOrderReceiveTxManager {
	return &PurchaseOrderReceiveTxManager{db: db}
}

func (m *PurchaseOrderReceiveTxManager) Run(
	ctx context.Context,
	fn func(procurementUsecase.PurchaseOrderReceiveTransactionalDeps) error,
) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		deps := procurementUsecase.PurchaseOrderReceiveTransactionalDeps{
			Repo: procurementRepo.NewPurchaseOrderRepository(tx),
			InventoryService: inventoryUsecase.NewInventoryUsecase(
				inventoryRepo.NewInventoryBalanceRepository(tx),
				inventoryRepo.NewInventoryMovementRepository(tx),
				inventoryRepo.NewInventoryLotRepository(tx),
			),
			CostEventRecorder: financeUsecase.NewCostEventWriter(financeRepo.NewCostEventRepository(tx)),
		}
		return fn(deps)
	})
}

type PurchaseOrderInspectTxManager struct {
	db *gorm.DB
}

func NewPurchaseOrderInspectTxManager(db *gorm.DB) *PurchaseOrderInspectTxManager {
	return &PurchaseOrderInspectTxManager{db: db}
}

func (m *PurchaseOrderInspectTxManager) Run(
	ctx context.Context,
	fn func(procurementUsecase.PurchaseOrderInspectTransactionalDeps) error,
) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		deps := procurementUsecase.PurchaseOrderInspectTransactionalDeps{
			Repo: procurementRepo.NewPurchaseOrderRepository(tx),
			InventoryService: inventoryUsecase.NewInventoryUsecase(
				inventoryRepo.NewInventoryBalanceRepository(tx),
				inventoryRepo.NewInventoryMovementRepository(tx),
				inventoryRepo.NewInventoryLotRepository(tx),
			),
		}
		return fn(deps)
	})
}

type purchaseOrderPlanCleanerTx struct {
	repo procurementDomain.ReplenishmentRepository
}

func (c *purchaseOrderPlanCleanerTx) DeletePlansByPurchaseOrderID(purchaseOrderID uint64) error {
	return c.repo.DeletePlansByPurchaseOrderID(purchaseOrderID)
}

type PurchaseOrderSubmitTxManager struct {
	db *gorm.DB
}

func NewPurchaseOrderSubmitTxManager(db *gorm.DB) *PurchaseOrderSubmitTxManager {
	return &PurchaseOrderSubmitTxManager{db: db}
}

func (m *PurchaseOrderSubmitTxManager) Run(
	ctx context.Context,
	fn func(procurementUsecase.PurchaseOrderSubmitTransactionalDeps) error,
) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		replenishmentRepo := procurementRepo.NewReplenishmentRepository(tx)
		deps := procurementUsecase.PurchaseOrderSubmitTransactionalDeps{
			Repo:              procurementRepo.NewPurchaseOrderRepository(tx),
			PlanCleaner:       &purchaseOrderPlanCleanerTx{repo: replenishmentRepo},
			CostEventRecorder: financeUsecase.NewCostEventWriter(financeRepo.NewCostEventRepository(tx)),
		}
		return fn(deps)
	})
}
