package transaction

import (
	"context"

	"am-erp-go/internal/module/packaging/repository"
	packagingUsecase "am-erp-go/internal/module/packaging/usecase"

	"gorm.io/gorm"
)

type packagingPlanConvertTxManager struct {
	db *gorm.DB
}

func NewPackagingPlanConvertTxManager(db *gorm.DB) packagingUsecase.PackagingPlanConvertTransactionManager {
	return &packagingPlanConvertTxManager{db: db}
}

func (m *packagingPlanConvertTxManager) Run(ctx context.Context, fn func(packagingUsecase.PackagingPlanConvertTransactionalDeps) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(packagingUsecase.PackagingPlanConvertTransactionalDeps{
			Repo: repository.NewPackagingProcurementRepository(tx),
		})
	})
}

type packagingPurchaseReceiveTxManager struct {
	db *gorm.DB
}

func NewPackagingPurchaseReceiveTxManager(db *gorm.DB) packagingUsecase.PackagingPurchaseReceiveTransactionManager {
	return &packagingPurchaseReceiveTxManager{db: db}
}

func (m *packagingPurchaseReceiveTxManager) Run(ctx context.Context, fn func(packagingUsecase.PackagingPurchaseReceiveTransactionalDeps) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(packagingUsecase.PackagingPurchaseReceiveTransactionalDeps{
			Repo:       repository.NewPackagingProcurementRepository(tx),
			ItemRepo:   repository.NewPackagingItemRepository(tx),
			LedgerRepo: repository.NewPackagingLedgerRepository(tx),
		})
	})
}
