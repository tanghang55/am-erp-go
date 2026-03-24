package transaction

import (
	"context"

	productRepo "am-erp-go/internal/module/product/repository"
	productUsecase "am-erp-go/internal/module/product/usecase"
	supplierRepo "am-erp-go/internal/module/supplier/repository"

	"gorm.io/gorm"
)

type ProductUpsertTxManager struct {
	db *gorm.DB
}

func NewProductUpsertTxManager(db *gorm.DB) *ProductUpsertTxManager {
	return &ProductUpsertTxManager{db: db}
}

func (m *ProductUpsertTxManager) Run(
	ctx context.Context,
	fn func(productUsecase.ProductUpsertTransactionalDeps) error,
) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		deps := productUsecase.ProductUpsertTransactionalDeps{
			ProductRepo: productRepo.NewProductRepository(tx),
			QuoteRepo:   supplierRepo.NewQuoteRepository(tx),
			ImageRepo:   productRepo.NewProductImageRepository(tx),
		}
		return fn(deps)
	})
}
