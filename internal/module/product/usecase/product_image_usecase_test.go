package usecase

import (
	"reflect"
	"testing"

	"am-erp-go/internal/module/product/domain"
)

type stubProductImageRepo struct {
	list      []domain.ProductImage
	productID uint64
	replaced  []string
	err       error
}

func (s *stubProductImageRepo) ListByProductID(productID uint64) ([]domain.ProductImage, error) {
	s.productID = productID
	return s.list, s.err
}

func (s *stubProductImageRepo) ReplaceAll(productID uint64, orderedUrls []string) error {
	s.productID = productID
	s.replaced = orderedUrls
	return s.err
}

type stubProductRepo struct {
	updatedID  uint64
	updatedUrl string
	err        error
}

func (s *stubProductRepo) List(_ *domain.ProductListParams) ([]domain.Product, int64, error) {
	return nil, 0, nil
}
func (s *stubProductRepo) GetByID(_ uint64) (*domain.Product, error) { return nil, nil }
func (s *stubProductRepo) Create(_ *domain.Product) error            { return nil }
func (s *stubProductRepo) Update(_ *domain.Product) error            { return nil }
func (s *stubProductRepo) Delete(_ uint64) error                     { return nil }
func (s *stubProductRepo) UpdateImageUrl(id uint64, url string) error {
	s.updatedID = id
	s.updatedUrl = url
	return s.err
}
func (s *stubProductRepo) GetDefaultSupplierID(_ uint64) (uint64, error) {
	return 0, nil
}
func (s *stubProductRepo) UpdateDefaultSupplierID(_, _ uint64) error { return nil }
func (s *stubProductRepo) ListByIDs(_ []uint64) ([]domain.Product, error) {
	return nil, nil
}
func (s *stubProductRepo) UpdateComboInfo(_ uint64, _ uint64, _ []uint64) error {
	return nil
}
func (s *stubProductRepo) ClearComboInfo(_ uint64) error { return nil }

func TestSaveProductImagesUpdatesPrimary(t *testing.T) {
	imageRepo := &stubProductImageRepo{
		list: []domain.ProductImage{
			{ImageUrl: "a", SortOrder: 1},
			{ImageUrl: "b", SortOrder: 2},
		},
	}
	productRepo := &stubProductRepo{}
	uc := NewProductImageUsecase(imageRepo, productRepo)

	_, err := uc.SaveProductImages(9, []string{"b", "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(imageRepo.replaced, []string{"b", "c"}) {
		t.Fatalf("expected replace list, got %#v", imageRepo.replaced)
	}
	if productRepo.updatedID != 9 || productRepo.updatedUrl != "b" {
		t.Fatalf("expected primary image sync")
	}
}

func TestSaveProductImagesRejectsOverLimit(t *testing.T) {
	imageRepo := &stubProductImageRepo{}
	productRepo := &stubProductRepo{}
	uc := NewProductImageUsecase(imageRepo, productRepo)

	urls := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}
	if _, err := uc.SaveProductImages(1, urls); err == nil {
		t.Fatalf("expected over-limit error")
	}
}
