package usecase

import (
	"am-erp-go/internal/module/product/domain"
)

type ProductUsecase struct {
	productRepo          domain.ProductRepository
	productParentRepo    domain.ProductParentRepository
	productPackagingRepo domain.ProductPackagingRepository
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
	productParentRepo domain.ProductParentRepository,
	productPackagingRepo domain.ProductPackagingRepository,
) *ProductUsecase {
	return &ProductUsecase{
		productRepo:          productRepo,
		productParentRepo:    productParentRepo,
		productPackagingRepo: productPackagingRepo,
	}
}

// Product SKU 相关方法

func (uc *ProductUsecase) ListProducts(params *domain.ProductListParams) ([]domain.Product, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.productRepo.List(params)
}

func (uc *ProductUsecase) GetProduct(id uint64) (*domain.Product, error) {
	return uc.productRepo.GetByID(id)
}

func (uc *ProductUsecase) CreateProduct(product *domain.Product) error {
	return uc.productRepo.Create(product)
}

func (uc *ProductUsecase) UpdateProduct(product *domain.Product) error {
	return uc.productRepo.Update(product)
}

func (uc *ProductUsecase) DeleteProduct(id uint64) error {
	return uc.productRepo.Delete(id)
}

// ProductParent 相关方法

func (uc *ProductUsecase) ListProductParents(params *domain.ProductParentListParams) ([]domain.ProductParent, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.productParentRepo.List(params)
}

func (uc *ProductUsecase) GetProductParent(id uint64) (*domain.ProductParent, error) {
	return uc.productParentRepo.GetByID(id)
}

func (uc *ProductUsecase) CreateProductParent(parent *domain.ProductParent) error {
	return uc.productParentRepo.Create(parent)
}

func (uc *ProductUsecase) UpdateProductParent(parent *domain.ProductParent) error {
	return uc.productParentRepo.Update(parent)
}

func (uc *ProductUsecase) DeleteProductParent(id uint64) error {
	return uc.productParentRepo.Delete(id)
}

// ProductPackaging 相关方法

func (uc *ProductUsecase) GetProductPackagingItems(productID uint64) ([]domain.ProductPackagingItem, error) {
	return uc.productPackagingRepo.ListByProductID(productID)
}

func (uc *ProductUsecase) SaveProductPackagingItems(productID uint64, items []domain.ProductPackagingItem) error {
	return uc.productPackagingRepo.ReplaceAll(productID, items)
}
