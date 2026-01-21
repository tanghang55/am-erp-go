package usecase

import "am-erp-go/internal/module/product/domain"

const maxProductImages = 10

type ProductImageUsecase struct {
	imageRepo   domain.ProductImageRepository
	productRepo domain.ProductRepository
}

func NewProductImageUsecase(imageRepo domain.ProductImageRepository, productRepo domain.ProductRepository) *ProductImageUsecase {
	return &ProductImageUsecase{
		imageRepo:   imageRepo,
		productRepo: productRepo,
	}
}

func (uc *ProductImageUsecase) ListProductImages(productID uint64) ([]domain.ProductImage, error) {
	return uc.imageRepo.ListByProductID(productID)
}

func (uc *ProductImageUsecase) SaveProductImages(productID uint64, orderedUrls []string) ([]domain.ProductImage, error) {
	if len(orderedUrls) > maxProductImages {
		return nil, domain.ErrProductImageLimit
	}
	if err := uc.imageRepo.ReplaceAll(productID, orderedUrls); err != nil {
		return nil, err
	}
	primary := ""
	if len(orderedUrls) > 0 {
		primary = orderedUrls[0]
	}
	if err := uc.productRepo.UpdateImageUrl(productID, primary); err != nil {
		return nil, err
	}
	return uc.imageRepo.ListByProductID(productID)
}
