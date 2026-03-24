package usecase

import (
	"errors"
	"fmt"
	"testing"

	"am-erp-go/internal/module/product/domain"
	supplierdomain "am-erp-go/internal/module/supplier/domain"
)

type stubProductValidationRepo struct {
	created  *domain.Product
	updated  *domain.Product
	products map[uint64]domain.Product
	configs  map[uint64]domain.ProductConfigItem
}

type stubProductQuoteRepo struct {
	quotes   map[string]supplierdomain.ProductSupplierQuote
	created  *supplierdomain.ProductSupplierQuote
	updated  *supplierdomain.ProductSupplierQuote
	createErr error
	updateErr error
	getErr    error
}

type stubPrimaryImageRepo struct {
	items          map[uint64][]domain.ProductImage
	replacedByID   map[uint64][]string
	listErr        error
	replaceAllErr  error
}

func (s *stubProductQuoteRepo) key(productID, supplierID uint64) string {
	return fmt.Sprintf("%d:%d", productID, supplierID)
}

func (s *stubProductQuoteRepo) GetByProductSupplier(productID, supplierID uint64) (*supplierdomain.ProductSupplierQuote, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.quotes == nil {
		return nil, supplierdomain.ErrQuoteNotFound
	}
	quote, ok := s.quotes[s.key(productID, supplierID)]
	if !ok {
		return nil, supplierdomain.ErrQuoteNotFound
	}
	copy := quote
	return &copy, nil
}

func (s *stubProductQuoteRepo) Create(quote *supplierdomain.ProductSupplierQuote) error {
	if s.createErr != nil {
		return s.createErr
	}
	if s.quotes == nil {
		s.quotes = map[string]supplierdomain.ProductSupplierQuote{}
	}
	copy := *quote
	if copy.ID == 0 {
		copy.ID = uint64(len(s.quotes) + 1)
	}
	s.quotes[s.key(copy.ProductID, copy.SupplierID)] = copy
	s.created = &copy
	return nil
}

func (s *stubProductQuoteRepo) Update(quote *supplierdomain.ProductSupplierQuote) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	if s.quotes == nil {
		s.quotes = map[string]supplierdomain.ProductSupplierQuote{}
	}
	copy := *quote
	s.quotes[s.key(copy.ProductID, copy.SupplierID)] = copy
	s.updated = &copy
	return nil
}

func (s *stubPrimaryImageRepo) ListByProductID(productID uint64) ([]domain.ProductImage, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	items := s.items[productID]
	result := make([]domain.ProductImage, len(items))
	copy(result, items)
	return result, nil
}

func (s *stubPrimaryImageRepo) ReplaceAll(productID uint64, orderedUrls []string) error {
	if s.replaceAllErr != nil {
		return s.replaceAllErr
	}
	if s.replacedByID == nil {
		s.replacedByID = map[uint64][]string{}
	}
	copied := append([]string(nil), orderedUrls...)
	s.replacedByID[productID] = copied
	images := make([]domain.ProductImage, 0, len(orderedUrls))
	for i, url := range orderedUrls {
		isPrimary := uint8(0)
		if i == 0 {
			isPrimary = 1
		}
		images = append(images, domain.ProductImage{
			ProductID: productID,
			ImageUrl:  url,
			SortOrder: uint32(i + 1),
			IsPrimary: isPrimary,
		})
	}
	if s.items == nil {
		s.items = map[uint64][]domain.ProductImage{}
	}
	s.items[productID] = images
	return nil
}

type stubProductDefaultsProvider struct {
	baseCurrency string
}

func (s *stubProductDefaultsProvider) GetDefaultBaseCurrency() string {
	return s.baseCurrency
}

func (s *stubProductValidationRepo) List(params *domain.ProductListParams) ([]domain.Product, int64, error) {
	items := make([]domain.Product, 0)
	for _, product := range s.products {
		if params != nil {
			if params.ComboID != nil {
				if product.ComboID == nil || *product.ComboID != *params.ComboID {
					continue
				}
			}
			if len(params.Statuses) > 0 {
				matched := false
				for _, status := range params.Statuses {
					if product.Status == status {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}
		}
		items = append(items, product)
	}
	return items, int64(len(items)), nil
}

func (s *stubProductValidationRepo) GetByID(id uint64) (*domain.Product, error) {
	product, ok := s.products[id]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := product
	return &copy, nil
}

func (s *stubProductValidationRepo) ListByIDs(ids []uint64) ([]domain.Product, error) {
	items := make([]domain.Product, 0, len(ids))
	for _, id := range ids {
		product, ok := s.products[id]
		if !ok {
			continue
		}
		items = append(items, product)
	}
	return items, nil
}

func (s *stubProductValidationRepo) ListByParentID(parentID uint64) ([]domain.Product, error) {
	items := make([]domain.Product, 0)
	for _, product := range s.products {
		if product.ParentID != nil && *product.ParentID == parentID {
			items = append(items, product)
		}
	}
	return items, nil
}

func (s *stubProductValidationRepo) CountReferencesByIDs(ids []uint64) (map[uint64]int64, error) {
	result := make(map[uint64]int64, len(ids))
	for _, id := range ids {
		result[id] = 0
		product, ok := s.products[id]
		if !ok {
			continue
		}
		if product.ComboID != nil && *product.ComboID != 0 {
			result[id]++
		}
		if product.ParentID != nil && *product.ParentID != 0 {
			result[id]++
		}
	}
	return result, nil
}

func (s *stubProductValidationRepo) Create(product *domain.Product) error {
	if product.ID == 0 {
		product.ID = uint64(len(s.products) + 1)
	}
	copy := *product
	s.created = &copy
	if s.products == nil {
		s.products = map[uint64]domain.Product{}
	}
	s.products[product.ID] = copy
	return nil
}

func (s *stubProductValidationRepo) CountByConfigReference(configType domain.ProductConfigType, configID uint64) (int64, error) {
	var count int64
	for _, product := range s.products {
		switch configType {
		case domain.ProductConfigTypeBrand:
			if product.BrandID != nil && *product.BrandID == configID {
				count++
			}
		case domain.ProductConfigTypeSalesStatus:
			item, ok := s.configs[configID]
			if ok && product.Status == item.ItemCode {
				count++
			}
		case domain.ProductConfigTypeDimensionUnit:
			if product.DimensionUnitID != nil && *product.DimensionUnitID == configID {
				count++
			}
		case domain.ProductConfigTypeWeightUnit:
			if product.WeightUnitID != nil && *product.WeightUnitID == configID {
				count++
			}
		}
	}
	return count, nil
}

func (s *stubProductValidationRepo) CountByCategoryID(categoryID uint64) (int64, error) {
	var count int64
	for _, product := range s.products {
		if product.CategoryID != nil && *product.CategoryID == categoryID {
			count++
		}
	}
	return count, nil
}

func (s *stubProductValidationRepo) Update(product *domain.Product) error {
	copy := *product
	s.updated = &copy
	return nil
}

func (s *stubProductValidationRepo) Delete(id uint64) error {
	return nil
}

func (s *stubProductValidationRepo) UpdateParentIDBatch(productIDs []uint64, parentID *uint64) error {
	for _, id := range productIDs {
		product, ok := s.products[id]
		if !ok {
			continue
		}
		product.ParentID = parentID
		s.products[id] = product
	}
	return nil
}

func (s *stubProductValidationRepo) UpdateImageUrl(id uint64, imageUrl string) error {
	return nil
}

func (s *stubProductValidationRepo) UpdateUnitCost(productID uint64, unitCost float64) error {
	product, ok := s.products[productID]
	if !ok {
		return nil
	}
	product.UnitCost = &unitCost
	s.products[productID] = product
	return nil
}

func (s *stubProductValidationRepo) GetDefaultSupplierID(productID uint64) (uint64, error) {
	return 0, nil
}

func (s *stubProductValidationRepo) UpdateDefaultSupplierID(productID, supplierID uint64) error {
	return nil
}

func (s *stubProductValidationRepo) UpdateComboInfo(comboID uint64, mainProductID uint64, productIDs []uint64) error {
	return nil
}

func (s *stubProductValidationRepo) ClearComboInfo(comboID uint64) error {
	return nil
}

type stubProductParentRepo struct {
	parents map[uint64]domain.ProductParent
	deleted uint64
}

type stubProductConfigRepo struct {
	items   map[uint64]domain.ProductConfigItem
	deleted uint64
}

func (s *stubProductConfigRepo) List(params *domain.ProductConfigListParams) ([]domain.ProductConfigItem, int64, error) {
	items := make([]domain.ProductConfigItem, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	return items, int64(len(items)), nil
}

func (s *stubProductConfigRepo) GetByID(id uint64) (*domain.ProductConfigItem, error) {
	item, ok := s.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := item
	return &copy, nil
}

func (s *stubProductConfigRepo) Create(item *domain.ProductConfigItem) error {
	if item.ID == 0 {
		item.ID = uint64(len(s.items) + 1)
	}
	if s.items == nil {
		s.items = map[uint64]domain.ProductConfigItem{}
	}
	copy := *item
	s.items[item.ID] = copy
	return nil
}

func (s *stubProductConfigRepo) Update(item *domain.ProductConfigItem) error {
	if s.items == nil {
		s.items = map[uint64]domain.ProductConfigItem{}
	}
	copy := *item
	s.items[item.ID] = copy
	return nil
}

func (s *stubProductConfigRepo) Delete(id uint64) error {
	s.deleted = id
	return nil
}

type stubProductCategoryRepo struct {
	items map[uint64]domain.ProductCategory
}

func (s *stubProductCategoryRepo) ListAll() ([]domain.ProductCategory, error) {
	items := make([]domain.ProductCategory, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	return items, nil
}

func (s *stubProductCategoryRepo) GetByID(id uint64) (*domain.ProductCategory, error) {
	item, ok := s.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := item
	return &copy, nil
}

func (s *stubProductCategoryRepo) Create(item *domain.ProductCategory) error {
	if item.ID == 0 {
		item.ID = uint64(len(s.items) + 1)
	}
	if s.items == nil {
		s.items = map[uint64]domain.ProductCategory{}
	}
	copy := *item
	s.items[item.ID] = copy
	return nil
}

func (s *stubProductCategoryRepo) Update(item *domain.ProductCategory) error {
	if s.items == nil {
		s.items = map[uint64]domain.ProductCategory{}
	}
	copy := *item
	s.items[item.ID] = copy
	return nil
}

func (s *stubProductCategoryRepo) Delete(id uint64) error {
	return nil
}

func (s *stubProductCategoryRepo) CountChildren(id uint64) (int64, error) {
	var count int64
	for _, item := range s.items {
		if item.ParentID != nil && *item.ParentID == id {
			count++
		}
	}
	return count, nil
}

func (s *stubProductParentRepo) List(params *domain.ProductParentListParams) ([]domain.ProductParent, int64, error) {
	items := make([]domain.ProductParent, 0, len(s.parents))
	for _, parent := range s.parents {
		items = append(items, parent)
	}
	return items, int64(len(items)), nil
}

func (s *stubProductParentRepo) GetByID(id uint64) (*domain.ProductParent, error) {
	parent, ok := s.parents[id]
	if !ok {
		return nil, errors.New("not found")
	}
	copy := parent
	return &copy, nil
}

func (s *stubProductParentRepo) Create(parent *domain.ProductParent) error {
	return nil
}

func (s *stubProductParentRepo) Update(parent *domain.ProductParent) error {
	copy := *parent
	s.parents[parent.ID] = copy
	return nil
}

func (s *stubProductParentRepo) Delete(id uint64) error {
	s.deleted = id
	delete(s.parents, id)
	return nil
}

func TestCreateProductRejectsMissingSupplier(t *testing.T) {
	repo := &stubProductValidationRepo{}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)

	err := uc.CreateProduct(&domain.Product{
		SellerSku:   "SKU-1",
		Asin:        "ASIN-1",
		Title:       "Product 1",
		Marketplace: "US",
	})
	if !errors.Is(err, ErrSupplierRequired) {
		t.Fatalf("expected ErrSupplierRequired, got %v", err)
	}
	if repo.created != nil {
		t.Fatalf("expected create not called")
	}
}

func TestCreateProductCreatesPendingDefaultSupplierQuote(t *testing.T) {
	repo := &stubProductValidationRepo{}
	quoteRepo := &stubProductQuoteRepo{}
	imageRepo := &stubPrimaryImageRepo{}
	defaults := &stubProductDefaultsProvider{baseCurrency: "CNY"}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)
	uc.BindQuoteRepository(quoteRepo)
	uc.BindImageRepository(imageRepo)
	uc.BindDefaultsProvider(defaults)

	supplierID := uint64(9)
	product := &domain.Product{
		SellerSku:   "SKU-1",
		Asin:        "ASIN-1",
		Title:       "Product 1",
		Marketplace: "US",
		SupplierID:  &supplierID,
		UnitCost:    float64Ptr(12.8),
		Status:      domain.ProductStatusDraft,
		ImageUrl:    "/uploads/products/main-create.png",
	}

	effects, err := uc.CreateProductWithEffects(product)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if effects == nil || !effects.AutoQuoteCreated {
		t.Fatalf("expected auto quote created, got %+v", effects)
	}
	if quoteRepo.created == nil {
		t.Fatalf("expected quote created")
	}
	if quoteRepo.created.ProductID != product.ID {
		t.Fatalf("expected product_id %d, got %d", product.ID, quoteRepo.created.ProductID)
	}
	if quoteRepo.created.SupplierID != supplierID {
		t.Fatalf("expected supplier_id %d, got %d", supplierID, quoteRepo.created.SupplierID)
	}
	if quoteRepo.created.Status != supplierdomain.ProductSupplierQuoteStatusPending {
		t.Fatalf("expected pending status, got %s", quoteRepo.created.Status)
	}
	if quoteRepo.created.Currency != "CNY" {
		t.Fatalf("expected currency CNY, got %s", quoteRepo.created.Currency)
	}
	if quoteRepo.created.Price != 12.8 {
		t.Fatalf("expected price 12.8, got %v", quoteRepo.created.Price)
	}
	if repo.created == nil || repo.created.UnitCost == nil || *repo.created.UnitCost != 12.8 {
		t.Fatalf("expected product unit_cost cached as 12.8, got %+v", repo.created)
	}
	if got := imageRepo.replacedByID[product.ID]; len(got) != 1 || got[0] != "/uploads/products/main-create.png" {
		t.Fatalf("expected primary image list initialized, got %+v", got)
	}
}

func TestUpdateProductCreatesPendingQuoteWhenSupplierChanged(t *testing.T) {
	supplierID := uint64(9)
	productID := uint64(10)
	repo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			productID: {
				ID:          productID,
				SellerSku:   "SKU-1",
				Asin:        "ASIN-1",
				Title:       "Product 1",
				Marketplace: "US",
				SupplierID:  &supplierID,
				Status:      domain.ProductStatusOnSale,
			},
		},
	}
	quoteRepo := &stubProductQuoteRepo{
		quotes: map[string]supplierdomain.ProductSupplierQuote{
			"10:9": {
				ID:         1,
				ProductID:  productID,
				SupplierID: supplierID,
				Status:     supplierdomain.ProductSupplierQuoteStatusActive,
				Currency:   "USD",
				QtyMOQ:     1,
			},
		},
	}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)
	uc.BindQuoteRepository(quoteRepo)

	newSupplierID := uint64(11)
	product := &domain.Product{
		ID:          productID,
		SellerSku:   "SKU-1",
		Asin:        "ASIN-1",
		Title:       "Product 1",
		Marketplace: "US",
		SupplierID:  &newSupplierID,
		UnitCost:    float64Ptr(22.6),
		Status:      domain.ProductStatusOnSale,
	}

	effects, err := uc.UpdateProductWithEffects(product)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if effects == nil || !effects.AutoQuoteCreated {
		t.Fatalf("expected auto quote created, got %+v", effects)
	}
	if quoteRepo.created == nil {
		t.Fatalf("expected new supplier quote created")
	}
	if quoteRepo.created.SupplierID != newSupplierID {
		t.Fatalf("expected supplier_id %d, got %d", newSupplierID, quoteRepo.created.SupplierID)
	}
	if quoteRepo.created.Price != 22.6 {
		t.Fatalf("expected created quote price 22.6, got %v", quoteRepo.created.Price)
	}
	if repo.updated == nil || repo.updated.UnitCost == nil || *repo.updated.UnitCost != 22.6 {
		t.Fatalf("expected cached product cost 22.6, got %+v", repo.updated)
	}
}

func TestUpdateProductUpdatesDefaultSupplierQuotePrice(t *testing.T) {
	supplierID := uint64(9)
	productID := uint64(10)
	existingPrice := 9.5
	repo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			productID: {
				ID:          productID,
				SellerSku:   "SKU-1",
				Asin:        "ASIN-1",
				Title:       "Product 1",
				Marketplace: "US",
				SupplierID:  &supplierID,
				UnitCost:    &existingPrice,
				Status:      domain.ProductStatusOnSale,
			},
		},
	}
	quoteRepo := &stubProductQuoteRepo{
		quotes: map[string]supplierdomain.ProductSupplierQuote{
			"10:9": {
				ID:         1,
				ProductID:  productID,
				SupplierID: supplierID,
				Price:      9.5,
				Status:     supplierdomain.ProductSupplierQuoteStatusActive,
				Currency:   "USD",
				QtyMOQ:     1,
			},
		},
	}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)
	uc.BindQuoteRepository(quoteRepo)

	product := &domain.Product{
		ID:          productID,
		SellerSku:   "SKU-1",
		Asin:        "ASIN-1",
		Title:       "Product 1",
		Marketplace: "US",
		SupplierID:  &supplierID,
		UnitCost:    float64Ptr(15.3),
		Status:      domain.ProductStatusOnSale,
	}

	effects, err := uc.UpdateProductWithEffects(product)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if effects == nil || effects.AutoQuoteCreated {
		t.Fatalf("expected existing quote updated without auto create, got %+v", effects)
	}
	if quoteRepo.updated == nil {
		t.Fatalf("expected quote updated")
	}
	if quoteRepo.updated.Price != 15.3 {
		t.Fatalf("expected updated quote price 15.3, got %v", quoteRepo.updated.Price)
	}
	if repo.updated == nil || repo.updated.UnitCost == nil || *repo.updated.UnitCost != 15.3 {
		t.Fatalf("expected cached product cost 15.3, got %+v", repo.updated)
	}
}

func TestUpdateProductSyncsNewPrimaryImageToImageList(t *testing.T) {
	supplierID := uint64(9)
	productID := uint64(10)
	repo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			productID: {
				ID:          productID,
				SellerSku:   "SKU-1",
				Asin:        "ASIN-1",
				Title:       "Product 1",
				Marketplace: "US",
				SupplierID:  &supplierID,
				ImageUrl:    "/uploads/products/old-main.png",
				Status:      domain.ProductStatusOnSale,
			},
		},
	}
	quoteRepo := &stubProductQuoteRepo{
		quotes: map[string]supplierdomain.ProductSupplierQuote{
			"10:9": {
				ID:         1,
				ProductID:  productID,
				SupplierID: supplierID,
				Price:      9.5,
				Status:     supplierdomain.ProductSupplierQuoteStatusActive,
				Currency:   "USD",
				QtyMOQ:     1,
			},
		},
	}
	imageRepo := &stubPrimaryImageRepo{
		items: map[uint64][]domain.ProductImage{
			productID: {
				{ProductID: productID, ImageUrl: "/uploads/products/old-main.png", SortOrder: 1, IsPrimary: 1},
				{ProductID: productID, ImageUrl: "/uploads/products/detail-2.png", SortOrder: 2, IsPrimary: 0},
			},
		},
	}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)
	uc.BindQuoteRepository(quoteRepo)
	uc.BindImageRepository(imageRepo)

	err := uc.UpdateProduct(&domain.Product{
		ID:          productID,
		SellerSku:   "SKU-1",
		Asin:        "ASIN-1",
		Title:       "Product 1",
		Marketplace: "US",
		SupplierID:  &supplierID,
		ImageUrl:    "/uploads/products/new-main.png",
		Status:      domain.ProductStatusOnSale,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := imageRepo.replacedByID[productID]
	expected := []string{"/uploads/products/new-main.png", "/uploads/products/old-main.png", "/uploads/products/detail-2.png"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d images, got %+v", len(expected), got)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("expected image order %+v, got %+v", expected, got)
		}
	}
}

func TestUpdateProductClearPrimaryPromotesNextImage(t *testing.T) {
	supplierID := uint64(9)
	productID := uint64(10)
	repo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			productID: {
				ID:          productID,
				SellerSku:   "SKU-1",
				Asin:        "ASIN-1",
				Title:       "Product 1",
				Marketplace: "US",
				SupplierID:  &supplierID,
				ImageUrl:    "/uploads/products/old-main.png",
				Status:      domain.ProductStatusOnSale,
			},
		},
	}
	quoteRepo := &stubProductQuoteRepo{
		quotes: map[string]supplierdomain.ProductSupplierQuote{
			"10:9": {
				ID:         1,
				ProductID:  productID,
				SupplierID: supplierID,
				Price:      9.5,
				Status:     supplierdomain.ProductSupplierQuoteStatusActive,
				Currency:   "USD",
				QtyMOQ:     1,
			},
		},
	}
	imageRepo := &stubPrimaryImageRepo{
		items: map[uint64][]domain.ProductImage{
			productID: {
				{ProductID: productID, ImageUrl: "/uploads/products/old-main.png", SortOrder: 1, IsPrimary: 1},
				{ProductID: productID, ImageUrl: "/uploads/products/detail-2.png", SortOrder: 2, IsPrimary: 0},
				{ProductID: productID, ImageUrl: "/uploads/products/detail-3.png", SortOrder: 3, IsPrimary: 0},
			},
		},
	}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)
	uc.BindQuoteRepository(quoteRepo)
	uc.BindImageRepository(imageRepo)

	err := uc.UpdateProduct(&domain.Product{
		ID:          productID,
		SellerSku:   "SKU-1",
		Asin:        "ASIN-1",
		Title:       "Product 1",
		Marketplace: "US",
		SupplierID:  &supplierID,
		ImageUrl:    "",
		Status:      domain.ProductStatusOnSale,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := imageRepo.replacedByID[productID]
	expected := []string{"/uploads/products/detail-2.png", "/uploads/products/detail-3.png"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d images, got %+v", len(expected), got)
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("expected image order %+v, got %+v", expected, got)
		}
	}
}

func TestCreateProductRejectsOverlongAsin(t *testing.T) {
	repo := &stubProductValidationRepo{}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)
	supplierID := uint64(9)

	err := uc.CreateProduct(&domain.Product{
		SellerSku:   "SKU-TEST-1",
		Asin:        "ASIN-TOO-LONG-123456789",
		Title:       "Product 1",
		Marketplace: "US",
		SupplierID:  &supplierID,
	})

	if !errors.Is(err, ErrProductAsinInvalid) {
		t.Fatalf("expected ErrProductAsinInvalid, got %v", err)
	}
}

func TestUpdateProductRejectsMissingSupplier(t *testing.T) {
	repo := &stubProductValidationRepo{}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)

	err := uc.UpdateProduct(&domain.Product{
		ID:          10,
		SellerSku:   "SKU-1",
		Asin:        "ASIN-1",
		Title:       "Product 1",
		Marketplace: "US",
	})
	if !errors.Is(err, ErrSupplierRequired) {
		t.Fatalf("expected ErrSupplierRequired, got %v", err)
	}
	if repo.updated != nil {
		t.Fatalf("expected update not called")
	}
}

func TestUpdateProductRejectsOverlongAsin(t *testing.T) {
	repo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "SKU-TEST-1", Asin: "ASIN-1", Title: "Product 1", Marketplace: "US", Status: domain.ProductStatusDraft},
		},
	}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)
	supplierID := uint64(9)

	err := uc.UpdateProduct(&domain.Product{
		ID:          1,
		SellerSku:   "SKU-TEST-1",
		Asin:        "ASIN-TOO-LONG-123456789",
		Title:       "Product 1",
		Marketplace: "US",
		Status:      domain.ProductStatusDraft,
		SupplierID:  &supplierID,
	})

	if !errors.Is(err, ErrProductAsinInvalid) {
		t.Fatalf("expected ErrProductAsinInvalid, got %v", err)
	}
}

func TestCreateProductConfigRejectsInvalidCode(t *testing.T) {
	uc := NewProductUsecase(nil, nil, &stubProductConfigRepo{}, nil, nil)

	err := uc.CreateProductConfig(&domain.ProductConfigItem{
		ConfigType: domain.ProductConfigTypeBrand,
		ItemCode:   "品牌@001",
		ItemName:   "测试品牌",
	})
	if !errors.Is(err, ErrProductConfigCodeInvalid) {
		t.Fatalf("expected ErrProductConfigCodeInvalid, got %v", err)
	}
}

func TestCreateProductConfigRejectsSalesStatus(t *testing.T) {
	uc := NewProductUsecase(nil, nil, &stubProductConfigRepo{}, nil, nil)

	err := uc.CreateProductConfig(&domain.ProductConfigItem{
		ConfigType: domain.ProductConfigTypeSalesStatus,
		ItemCode:   "CUSTOM_STATUS",
		ItemName:   "自定义状态",
		Status:     "ACTIVE",
	})
	if !errors.Is(err, ErrProductConfigSystemFixed) {
		t.Fatalf("expected ErrProductConfigSystemFixed, got %v", err)
	}
}

func TestCreateProductCategoryRejectsInvalidCode(t *testing.T) {
	parentID := uint64(1)
	uc := NewProductUsecase(nil, nil, nil, &stubProductCategoryRepo{
		items: map[uint64]domain.ProductCategory{
			1: {ID: 1, CategoryCode: "ROOT", CategoryName: "一级品类", Level: 1, Status: "ACTIVE"},
		},
	}, nil)

	err := uc.CreateProductCategory(&domain.ProductCategory{
		ParentID:     &parentID,
		CategoryCode: "品类@二级",
		CategoryName: "二级品类",
		Status:       "ACTIVE",
	})
	if !errors.Is(err, ErrProductCategoryCodeInvalid) {
		t.Fatalf("expected ErrProductCategoryCodeInvalid, got %v", err)
	}
}

func TestCreateProductCategoryAllowsLevel1Category(t *testing.T) {
	repo := &stubProductCategoryRepo{items: map[uint64]domain.ProductCategory{}}
	uc := NewProductUsecase(nil, nil, nil, repo, nil)

	item := &domain.ProductCategory{
		CategoryCode: "ROOT_ELEC",
		CategoryName: "电子产品",
		Status:       "ACTIVE",
	}

	if err := uc.CreateProductCategory(item); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if item.Level != 1 {
		t.Fatalf("expected level 1, got %d", item.Level)
	}
	if item.ParentID != nil {
		t.Fatalf("expected nil parent_id, got %v", *item.ParentID)
	}
}

func TestCreateProductCategoryRejectsLevel3Parent(t *testing.T) {
	parentID := uint64(3)
	uc := NewProductUsecase(nil, nil, nil, &stubProductCategoryRepo{
		items: map[uint64]domain.ProductCategory{
			3: {ID: 3, CategoryCode: "CAT_L3_01", CategoryName: "三级品类", Level: 3, Status: "ACTIVE"},
		},
	}, nil)

	err := uc.CreateProductCategory(&domain.ProductCategory{
		ParentID:     &parentID,
		CategoryCode: "CAT_L4_01",
		CategoryName: "四级品类",
		Status:       "ACTIVE",
	})
	if !errors.Is(err, ErrProductCategoryParentInvalid) {
		t.Fatalf("expected ErrProductCategoryParentInvalid, got %v", err)
	}
}

func TestUpdateProductCategoryAllowsLevel1Category(t *testing.T) {
	uc := NewProductUsecase(nil, nil, nil, &stubProductCategoryRepo{
		items: map[uint64]domain.ProductCategory{
			1: {ID: 1, CategoryCode: "ROOT", CategoryName: "一级品类", Level: 1, Status: "ACTIVE"},
		},
	}, nil)

	err := uc.UpdateProductCategory(&domain.ProductCategory{
		ID:           1,
		ParentID:     nil,
		CategoryCode: "ROOT",
		CategoryName: "一级品类-更新",
		Status:       "ACTIVE",
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestUpdateProductCategoryRejectsReferencedCodeChange(t *testing.T) {
	rootID := uint64(1)
	categoryID := uint64(9)
	uc := NewProductUsecase(
		&stubProductValidationRepo{
			products: map[uint64]domain.Product{
				1: {ID: 1, SellerSku: "SKU-1", CategoryID: &categoryID},
			},
		},
		nil,
		nil,
		&stubProductCategoryRepo{
			items: map[uint64]domain.ProductCategory{
				rootID:     {ID: rootID, CategoryCode: "ROOT", CategoryName: "一级", Level: 1, Status: "ACTIVE"},
				categoryID: {ID: categoryID, ParentID: &rootID, CategoryCode: "CAT-OLD", CategoryName: "二级", Level: 2, Status: "ACTIVE"},
			},
		},
		nil,
	)

	err := uc.UpdateProductCategory(&domain.ProductCategory{
		ID:           categoryID,
		ParentID:     &rootID,
		CategoryCode: "CAT-NEW",
		CategoryName: "二级",
		Status:       "ACTIVE",
	})
	if !errors.Is(err, ErrProductCategoryCodeImmutable) {
		t.Fatalf("expected ErrProductCategoryCodeImmutable, got %v", err)
	}
}

func TestDeleteProductCategoryAllowsLevel1CategoryWithoutChildrenOrReferences(t *testing.T) {
	uc := NewProductUsecase(nil, nil, nil, &stubProductCategoryRepo{
		items: map[uint64]domain.ProductCategory{
			1: {ID: 1, CategoryCode: "ROOT", CategoryName: "一级品类", Level: 1, Status: "ACTIVE"},
		},
	}, nil)

	err := uc.DeleteProductCategory(1)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestUpdateProductConfigRejectsReferencedSalesStatusCodeChange(t *testing.T) {
	statusID := uint64(7)
	uc := NewProductUsecase(
		&stubProductValidationRepo{
			products: map[uint64]domain.Product{
				1: {ID: 1, SellerSku: "SKU-1", Status: "ON_SALE"},
			},
			configs: map[uint64]domain.ProductConfigItem{
				statusID: {ID: statusID, ConfigType: domain.ProductConfigTypeSalesStatus, ItemCode: "ON_SALE", ItemName: "正常销售", Status: "ACTIVE"},
			},
		},
		nil,
		&stubProductConfigRepo{
			items: map[uint64]domain.ProductConfigItem{
				statusID: {ID: statusID, ConfigType: domain.ProductConfigTypeSalesStatus, ItemCode: "ON_SALE", ItemName: "正常销售", Status: "ACTIVE"},
			},
		},
		nil,
		nil,
	)

	err := uc.UpdateProductConfig(&domain.ProductConfigItem{
		ID:         statusID,
		ConfigType: domain.ProductConfigTypeSalesStatus,
		ItemCode:   "SELLING",
		ItemName:   "销售中",
		Status:     "ACTIVE",
	})
	if !errors.Is(err, ErrProductConfigSystemFixed) {
		t.Fatalf("expected ErrProductConfigSystemFixed, got %v", err)
	}
}

func TestUpdateProductConfigRejectsSalesStatusCodeChange(t *testing.T) {
	statusID := uint64(7)
	uc := NewProductUsecase(
		nil,
		nil,
		&stubProductConfigRepo{
			items: map[uint64]domain.ProductConfigItem{
				statusID: {ID: statusID, ConfigType: domain.ProductConfigTypeSalesStatus, ItemCode: "ON_SALE", ItemName: "正常销售", Status: "ACTIVE"},
			},
		},
		nil,
		nil,
	)

	err := uc.UpdateProductConfig(&domain.ProductConfigItem{
		ID:         statusID,
		ConfigType: domain.ProductConfigTypeSalesStatus,
		ItemCode:   "REPLENISHING",
		ItemName:   "正常销售",
		Status:     "ACTIVE",
	})
	if !errors.Is(err, ErrProductConfigSystemFixed) {
		t.Fatalf("expected ErrProductConfigSystemFixed, got %v", err)
	}
}

func TestDeleteProductConfigRejectsReferencedBrand(t *testing.T) {
	brandID := uint64(9)
	uc := NewProductUsecase(
		&stubProductValidationRepo{
			products: map[uint64]domain.Product{
				1: {ID: 1, SellerSku: "SKU-1", BrandID: &brandID},
			},
		},
		nil,
		&stubProductConfigRepo{
			items: map[uint64]domain.ProductConfigItem{
				brandID: {ID: brandID, ConfigType: domain.ProductConfigTypeBrand, ItemCode: "BRAND-1", ItemName: "品牌1"},
			},
		},
		nil,
		nil,
	)

	err := uc.DeleteProductConfig(brandID)
	if !errors.Is(err, ErrProductConfigReferenced) {
		t.Fatalf("expected ErrProductConfigReferenced, got %v", err)
	}
}

func TestDeleteProductCategoryRejectsReferencedCategory(t *testing.T) {
	categoryID := uint64(9)
	parentID := uint64(1)
	uc := NewProductUsecase(
		&stubProductValidationRepo{
			products: map[uint64]domain.Product{
				1: {ID: 1, SellerSku: "SKU-1", CategoryID: &categoryID},
			},
		},
		nil,
		nil,
		&stubProductCategoryRepo{
			items: map[uint64]domain.ProductCategory{
				parentID:   {ID: parentID, CategoryCode: "ROOT", CategoryName: "一级品类", Level: 1, Status: "ACTIVE"},
				categoryID: {ID: categoryID, ParentID: &parentID, CategoryCode: "CAT-2", CategoryName: "二级品类", Level: 2, Status: "ACTIVE"},
			},
		},
		nil,
	)

	err := uc.DeleteProductCategory(categoryID)
	if !errors.Is(err, ErrProductCategoryReferenced) {
		t.Fatalf("expected ErrProductCategoryReferenced, got %v", err)
	}
}

func TestListProductsIncludesDeleteState(t *testing.T) {
	comboID := uint64(66)
	parentID := uint64(8)
	repo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "P-1", SupplierID: uint64Ptr(1)},
			2: {ID: 2, SellerSku: "P-2", SupplierID: uint64Ptr(1), ComboID: &comboID},
			3: {ID: 3, SellerSku: "P-3", SupplierID: uint64Ptr(1), ParentID: &parentID},
		},
	}

	uc := NewProductUsecase(repo, nil, nil, nil, nil)
	items, total, err := uc.ListProducts(&domain.ProductListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("ListProducts returned error: %v", err)
	}
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}

	deleteState := make(map[uint64]domain.Product)
	for _, item := range items {
		deleteState[item.ID] = item
	}
	if !deleteState[1].Deletable || deleteState[1].ReferenceCount != 0 {
		t.Fatalf("expected product 1 deletable, got %+v", deleteState[1])
	}
	if deleteState[2].Deletable || deleteState[2].ReferenceCount == 0 {
		t.Fatalf("expected combo product locked for delete, got %+v", deleteState[2])
	}
	if deleteState[3].Deletable || deleteState[3].DeleteBlockReason == "" {
		t.Fatalf("expected grouped product locked for delete, got %+v", deleteState[3])
	}
}

func TestDeleteProductRejectsReferencedProduct(t *testing.T) {
	comboID := uint64(88)
	repo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "P-1", SupplierID: uint64Ptr(1), ComboID: &comboID},
		},
	}

	uc := NewProductUsecase(repo, nil, nil, nil, nil)
	err := uc.DeleteProduct(1)
	if !errors.Is(err, ErrProductReferenced) {
		t.Fatalf("expected ErrProductReferenced, got %v", err)
	}
}

func TestListProductConfigsIncludesDeleteState(t *testing.T) {
	brandID := uint64(9)
	dimensionID := uint64(10)
	uc := NewProductUsecase(
		&stubProductValidationRepo{
			products: map[uint64]domain.Product{
				1: {ID: 1, SellerSku: "SKU-1", BrandID: &brandID},
			},
		},
		nil,
		&stubProductConfigRepo{
			items: map[uint64]domain.ProductConfigItem{
				brandID:     {ID: brandID, ConfigType: domain.ProductConfigTypeBrand, ItemCode: "BRAND-1", ItemName: "品牌1"},
				dimensionID: {ID: dimensionID, ConfigType: domain.ProductConfigTypeDimensionUnit, ItemCode: "CM", ItemName: "厘米"},
			},
		},
		nil,
		nil,
	)

	items, total, err := uc.ListProductConfigs(&domain.ProductConfigListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}

	var brandItem, dimensionItem *domain.ProductConfigItem
	for i := range items {
		switch items[i].ID {
		case brandID:
			brandItem = &items[i]
		case dimensionID:
			dimensionItem = &items[i]
		}
	}
	if brandItem == nil || dimensionItem == nil {
		t.Fatalf("expected both config items present, got %+v", items)
	}
	if brandItem.Deletable {
		t.Fatalf("expected referenced brand config to be non-deletable")
	}
	if brandItem.ReferenceCount != 1 {
		t.Fatalf("expected referenced brand count 1, got %d", brandItem.ReferenceCount)
	}
	if brandItem.DeleteBlockReason == "" {
		t.Fatalf("expected referenced brand to have delete block reason")
	}
	if !dimensionItem.Deletable {
		t.Fatalf("expected unused dimension unit to be deletable")
	}
	if dimensionItem.ReferenceCount != 0 {
		t.Fatalf("expected unused dimension unit count 0, got %d", dimensionItem.ReferenceCount)
	}
}

func TestListProductCategoriesIncludesDeleteState(t *testing.T) {
	rootID := uint64(1)
	referencedID := uint64(2)
	emptyID := uint64(3)
	uc := NewProductUsecase(
		&stubProductValidationRepo{
			products: map[uint64]domain.Product{
				1: {ID: 1, SellerSku: "SKU-1", CategoryID: &referencedID},
			},
		},
		nil,
		nil,
		&stubProductCategoryRepo{
			items: map[uint64]domain.ProductCategory{
				rootID:       {ID: rootID, CategoryCode: "ROOT", CategoryName: "一级品类", Level: 1, Status: "ACTIVE"},
				referencedID: {ID: referencedID, ParentID: &rootID, CategoryCode: "CAT-L2-A", CategoryName: "二级A", Level: 2, Status: "ACTIVE"},
				emptyID:      {ID: emptyID, ParentID: &rootID, CategoryCode: "CAT-L2-B", CategoryName: "二级B", Level: 2, Status: "ACTIVE"},
			},
		},
		nil,
	)

	items, err := uc.ListProductCategories()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one root category, got %d", len(items))
	}
	root := items[0]
	if root.Deletable {
		t.Fatalf("expected level 1 root to be non-deletable")
	}
	if root.DeleteBlockReason == "" {
		t.Fatalf("expected root category to have delete block reason")
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 child categories, got %d", len(root.Children))
	}

	var referencedItem, emptyItem *domain.ProductCategory
	for _, child := range root.Children {
		switch child.ID {
		case referencedID:
			referencedItem = child
		case emptyID:
			emptyItem = child
		}
	}
	if referencedItem == nil || emptyItem == nil {
		t.Fatalf("expected both child categories present, got %+v", root.Children)
	}
	if referencedItem.Deletable {
		t.Fatalf("expected referenced category to be non-deletable")
	}
	if referencedItem.ReferenceCount != 1 {
		t.Fatalf("expected referenced category count 1, got %d", referencedItem.ReferenceCount)
	}
	if emptyItem.Deletable != true {
		t.Fatalf("expected unused category to be deletable")
	}
	if emptyItem.ReferenceCount != 0 {
		t.Fatalf("expected unused category count 0, got %d", emptyItem.ReferenceCount)
	}
}

func TestUpdateProductRejectsActiveComboMainWithoutActiveChildren(t *testing.T) {
	comboID := uint64(11)
	supplierID := uint64(8)
	repo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "COMBO-MAIN", Marketplace: "US", SupplierID: &supplierID, ComboID: &comboID, IsComboMain: 1, Status: domain.ProductStatusDraft},
			2: {ID: 2, SellerSku: "CHILD-1", Marketplace: "US", SupplierID: &supplierID, ComboID: &comboID, IsComboMain: 0, Status: domain.ProductStatusOffShelf},
			3: {ID: 3, SellerSku: "CHILD-2", Marketplace: "US", SupplierID: &supplierID, ComboID: &comboID, IsComboMain: 0, Status: domain.ProductStatusDraft},
		},
	}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)

	err := uc.UpdateProduct(&domain.Product{
		ID:          1,
		SellerSku:   "COMBO-MAIN",
		Asin:        "ASIN-COMBO",
		Title:       "Combo Main",
		Marketplace: "US",
		SupplierID:  &supplierID,
		ComboID:     &comboID,
		IsComboMain: 1,
		Status:      domain.ProductStatusOnSale,
	})
	if !errors.Is(err, ErrComboMainRequiresActiveChildren) {
		t.Fatalf("expected ErrComboMainRequiresActiveChildren, got %v", err)
	}
	if repo.updated != nil {
		t.Fatalf("expected update not called")
	}
}

func TestUpdateProductRejectsDeactivatingLastActiveComboChild(t *testing.T) {
	comboID := uint64(12)
	supplierID := uint64(8)
	repo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "COMBO-MAIN", Marketplace: "US", SupplierID: &supplierID, ComboID: &comboID, IsComboMain: 1, Status: domain.ProductStatusOnSale},
			2: {ID: 2, SellerSku: "CHILD-1", Marketplace: "US", SupplierID: &supplierID, ComboID: &comboID, IsComboMain: 0, Status: domain.ProductStatusOnSale},
			3: {ID: 3, SellerSku: "CHILD-2", Marketplace: "US", SupplierID: &supplierID, ComboID: &comboID, IsComboMain: 0, Status: domain.ProductStatusOffShelf},
		},
	}
	uc := NewProductUsecase(repo, nil, nil, nil, nil)

	err := uc.UpdateProduct(&domain.Product{
		ID:          2,
		SellerSku:   "CHILD-1",
		Asin:        "ASIN-CHILD-1",
		Title:       "Child 1",
		Marketplace: "US",
		SupplierID:  &supplierID,
		ComboID:     &comboID,
		Status:      domain.ProductStatusOffShelf,
	})
	if !errors.Is(err, ErrComboActiveRequiresActiveChildren) {
		t.Fatalf("expected ErrComboActiveRequiresActiveChildren, got %v", err)
	}
	if repo.updated != nil {
		t.Fatalf("expected update not called")
	}
}

func TestGetProductParentIncludesChildren(t *testing.T) {
	parentID := uint64(10)
	productRepo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "SKU-1", Marketplace: "US", Status: domain.ProductStatusOnSale, ParentID: &parentID},
			2: {ID: 2, SellerSku: "SKU-2", Marketplace: "US", Status: domain.ProductStatusDraft, ParentID: &parentID},
		},
	}
	parentRepo := &stubProductParentRepo{
		parents: map[uint64]domain.ProductParent{
			parentID: {ID: parentID, ParentAsin: "PARENT", Marketplace: "US"},
		},
	}

	uc := NewProductUsecase(productRepo, parentRepo, nil, nil, nil)
	parent, err := uc.GetProductParent(parentID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(parent.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(parent.Children))
	}
	if parent.ChildCount != 2 || parent.ActiveChildCount != 1 || parent.InactiveChildCount != 1 {
		t.Fatalf("unexpected counts: %+v", parent)
	}
}

func TestDeleteProductParentRejectsWhenChildrenExist(t *testing.T) {
	parentID := uint64(10)
	productRepo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "SKU-1", Marketplace: "US", ParentID: &parentID},
		},
	}
	parentRepo := &stubProductParentRepo{
		parents: map[uint64]domain.ProductParent{
			parentID: {ID: parentID, ParentAsin: "PARENT", Marketplace: "US"},
		},
	}

	uc := NewProductUsecase(productRepo, parentRepo, nil, nil, nil)
	err := uc.DeleteProductParent(parentID)
	if !errors.Is(err, ErrProductParentHasChildren) {
		t.Fatalf("expected ErrProductParentHasChildren, got %v", err)
	}
	if parentRepo.deleted != 0 {
		t.Fatalf("expected parent not deleted")
	}
}

func TestAttachProductParentChildrenRejectsChildrenFromOtherParent(t *testing.T) {
	parentID := uint64(10)
	otherParentID := uint64(20)
	productRepo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "SKU-1", Marketplace: "US", ParentID: &otherParentID},
		},
	}
	parentRepo := &stubProductParentRepo{
		parents: map[uint64]domain.ProductParent{
			parentID: {ID: parentID, ParentAsin: "PARENT", Marketplace: "US"},
		},
	}

	uc := NewProductUsecase(productRepo, parentRepo, nil, nil, nil)
	_, err := uc.AttachProductParentChildren(parentID, []uint64{1})
	if !errors.Is(err, ErrProductParentChildAlreadyAssigned) {
		t.Fatalf("expected ErrProductParentChildAlreadyAssigned, got %v", err)
	}
}

func TestAttachProductParentChildrenRejectsMarketplaceMismatch(t *testing.T) {
	parentID := uint64(10)
	productRepo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "SKU-1", Marketplace: "CA"},
		},
	}
	parentRepo := &stubProductParentRepo{
		parents: map[uint64]domain.ProductParent{
			parentID: {ID: parentID, ParentAsin: "PARENT", Marketplace: "US"},
		},
	}

	uc := NewProductUsecase(productRepo, parentRepo, nil, nil, nil)
	_, err := uc.AttachProductParentChildren(parentID, []uint64{1})
	if !errors.Is(err, ErrProductParentChildMarketplaceMismatch) {
		t.Fatalf("expected ErrProductParentChildMarketplaceMismatch, got %v", err)
	}
}

func TestAttachAndDetachProductParentChildren(t *testing.T) {
	parentID := uint64(10)
	productRepo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "SKU-1", Marketplace: "US", Status: domain.ProductStatusOnSale},
			2: {ID: 2, SellerSku: "SKU-2", Marketplace: "US", Status: domain.ProductStatusDraft},
		},
	}
	parentRepo := &stubProductParentRepo{
		parents: map[uint64]domain.ProductParent{
			parentID: {ID: parentID, ParentAsin: "PARENT", Marketplace: "US"},
		},
	}

	uc := NewProductUsecase(productRepo, parentRepo, nil, nil, nil)
	parent, err := uc.AttachProductParentChildren(parentID, []uint64{1, 2})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(parent.Children) != 2 {
		t.Fatalf("expected 2 children after attach, got %d", len(parent.Children))
	}

	parent, err = uc.DetachProductParentChild(parentID, 2)
	if err != nil {
		t.Fatalf("expected no error when detach, got %v", err)
	}
	if len(parent.Children) != 1 {
		t.Fatalf("expected 1 child after detach, got %d", len(parent.Children))
	}
	if productRepo.products[2].ParentID != nil {
		t.Fatalf("expected child parent_id cleared")
	}
}

func TestAttachProductParentChildrenAllowsUsedChildren(t *testing.T) {
	parentID := uint64(10)
	productRepo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "SKU-1", Marketplace: "US", Status: domain.ProductStatusOnSale},
		},
	}
	parentRepo := &stubProductParentRepo{
		parents: map[uint64]domain.ProductParent{
			parentID: {ID: parentID, ParentAsin: "PARENT", Marketplace: "US"},
		},
	}

	uc := NewProductUsecase(productRepo, parentRepo, nil, nil, nil)
	parent, err := uc.AttachProductParentChildren(parentID, []uint64{1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(parent.Children) != 1 {
		t.Fatalf("expected 1 child after attach, got %d", len(parent.Children))
	}
	if productRepo.products[1].ParentID == nil || *productRepo.products[1].ParentID != parentID {
		t.Fatalf("expected child parent_id updated")
	}
}

func TestDetachProductParentChildAllowsUsedChildren(t *testing.T) {
	parentID := uint64(10)
	productRepo := &stubProductValidationRepo{
		products: map[uint64]domain.Product{
			1: {ID: 1, SellerSku: "SKU-1", Marketplace: "US", Status: domain.ProductStatusOnSale, ParentID: &parentID},
		},
	}
	parentRepo := &stubProductParentRepo{
		parents: map[uint64]domain.ProductParent{
			parentID: {ID: parentID, ParentAsin: "PARENT", Marketplace: "US"},
		},
	}

	uc := NewProductUsecase(productRepo, parentRepo, nil, nil, nil)
	parent, err := uc.DetachProductParentChild(parentID, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(parent.Children) != 0 {
		t.Fatalf("expected no child after detach, got %d", len(parent.Children))
	}
	if productRepo.products[1].ParentID != nil {
		t.Fatalf("expected child parent_id cleared")
	}
}

func float64Ptr(value float64) *float64 {
	return &value
}
