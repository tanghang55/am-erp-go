package usecase

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"am-erp-go/internal/module/product/domain"
)

type stubComboRepo struct {
	createdMainID   uint64
	createdProducts []uint64
	createdRatios   map[uint64]uint64
	createdComboID  uint64

	replaceComboID  uint64
	replaceMainID   uint64
	replaceProducts []uint64
	replaceRatios   map[uint64]uint64

	deletedComboID uint64

	itemsByCombo map[uint64][]domain.ProductComboItem
	listComboIDs []uint64
}

func (s *stubComboRepo) ListComboIDs(params *domain.ComboListParams) ([]uint64, int64, error) {
	if len(s.listComboIDs) == 0 {
		return nil, 0, nil
	}
	return append([]uint64{}, s.listComboIDs...), int64(len(s.listComboIDs)), nil
}

func (s *stubComboRepo) GetItemsByComboID(comboID uint64) ([]domain.ProductComboItem, error) {
	return append([]domain.ProductComboItem{}, s.itemsByCombo[comboID]...), nil
}

func (s *stubComboRepo) CreateCombo(mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) (uint64, error) {
	s.createdMainID = mainProductID
	s.createdProducts = append([]uint64{}, productIDs...)
	s.createdRatios = qtyRatios
	if s.createdComboID == 0 {
		s.createdComboID = 99
	}
	items := []domain.ProductComboItem{
		{ComboID: s.createdComboID, MainProductID: mainProductID, ProductID: mainProductID, QtyRatio: 1},
	}
	for _, productID := range productIDs {
		items = append(items, domain.ProductComboItem{
			ComboID:       s.createdComboID,
			MainProductID: mainProductID,
			ProductID:     productID,
			QtyRatio:      qtyRatios[productID],
		})
	}
	if s.itemsByCombo == nil {
		s.itemsByCombo = make(map[uint64][]domain.ProductComboItem)
	}
	s.itemsByCombo[s.createdComboID] = items
	return s.createdComboID, nil
}

func (s *stubComboRepo) ReplaceComboItems(comboID uint64, mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) error {
	s.replaceComboID = comboID
	s.replaceMainID = mainProductID
	s.replaceProducts = append([]uint64{}, productIDs...)
	s.replaceRatios = qtyRatios
	items := []domain.ProductComboItem{
		{ComboID: comboID, MainProductID: mainProductID, ProductID: mainProductID, QtyRatio: 1},
	}
	for _, productID := range productIDs {
		items = append(items, domain.ProductComboItem{
			ComboID:       comboID,
			MainProductID: mainProductID,
			ProductID:     productID,
			QtyRatio:      qtyRatios[productID],
		})
	}
	if s.itemsByCombo == nil {
		s.itemsByCombo = make(map[uint64][]domain.ProductComboItem)
	}
	s.itemsByCombo[comboID] = items
	return nil
}

func (s *stubComboRepo) DeleteCombo(comboID uint64) error {
	s.deletedComboID = comboID
	delete(s.itemsByCombo, comboID)
	return nil
}

type stubComboProductRepo struct {
	updatedComboID uint64
	updatedMainID  uint64
	updatedIDs     []uint64

	clearedComboID uint64
}

func (s *stubComboProductRepo) UpdateComboInfo(comboID uint64, mainProductID uint64, productIDs []uint64) error {
	s.updatedComboID = comboID
	s.updatedMainID = mainProductID
	s.updatedIDs = append([]uint64{}, productIDs...)
	return nil
}

type stubComboProductReader struct {
	products []domain.Product
}

func (s *stubComboProductReader) ListByIDs(ids []uint64) ([]domain.Product, error) {
	results := make([]domain.Product, 0, len(ids))
	for _, id := range ids {
		for _, product := range s.products {
			if product.ID == id {
				results = append(results, product)
				break
			}
		}
	}
	return results, nil
}

func (s *stubComboProductRepo) ClearComboInfo(comboID uint64) error {
	s.clearedComboID = comboID
	return nil
}

type stubComboUsageChecker struct {
	locked bool
	err    error
	since  time.Time
	ids    []uint64
}

func (s *stubComboUsageChecker) HasBusinessUsageSince(productIDs []uint64, since time.Time) (bool, error) {
	s.ids = append([]uint64{}, productIDs...)
	s.since = since
	return s.locked, s.err
}

func TestCreateComboUpdatesProductComboInfo(t *testing.T) {
	comboRepo := &stubComboRepo{createdComboID: 88}
	productRepo := &stubComboProductRepo{}
	productReader := &stubComboProductReader{
		products: []domain.Product{
			{ID: 10, SellerSku: "MAIN-10", Title: "Main 10"},
			{ID: 11, SellerSku: "CHILD-11", Title: "Child 11"},
			{ID: 12, SellerSku: "CHILD-12", Title: "Child 12"},
		},
	}

	uc := NewProductComboUsecase(comboRepo, productRepo, productReader, nil)

	params := domain.ComboUpsertParams{
		MainProductID: 10,
		Children: []domain.ComboChildInput{
			{ProductID: 10, QtyRatio: 9},
			{ProductID: 11, QtyRatio: 2},
			{ProductID: 12, QtyRatio: 3},
			{ProductID: 11, QtyRatio: 4},
		},
	}

	combo, err := uc.CreateCombo(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if combo.ComboID != 88 {
		t.Fatalf("expected combo id 88, got %d", combo.ComboID)
	}
	if combo.MainProduct.ID != 10 {
		t.Fatalf("expected hydrated main product 10, got %d", combo.MainProduct.ID)
	}
	if len(combo.Products) != 2 || combo.Products[0].QtyRatio != 4 || combo.Products[1].QtyRatio != 3 {
		t.Fatalf("expected hydrated child ratios [4,3], got %+v", combo.Products)
	}

	if comboRepo.createdMainID != 10 {
		t.Fatalf("expected main product id 10")
	}

	if !reflect.DeepEqual(comboRepo.createdProducts, []uint64{11, 12}) {
		t.Fatalf("expected child products [11 12], got %v", comboRepo.createdProducts)
	}

	if !reflect.DeepEqual(comboRepo.createdRatios, map[uint64]uint64{11: 4, 12: 3}) {
		t.Fatalf("expected qty ratios map[11:4 12:3], got %v", comboRepo.createdRatios)
	}

	if productRepo.updatedComboID != 88 || productRepo.updatedMainID != 10 {
		t.Fatalf("expected update combo info called")
	}

	if !reflect.DeepEqual(productRepo.updatedIDs, []uint64{10, 11, 12}) {
		t.Fatalf("expected updated ids [10 11 12], got %v", productRepo.updatedIDs)
	}
}

func TestUpdateComboReplacesItemsAndClearsOld(t *testing.T) {
	comboRepo := &stubComboRepo{
		itemsByCombo: map[uint64][]domain.ProductComboItem{
			77: {
				{ComboID: 77, MainProductID: 20, ProductID: 20, QtyRatio: 1},
				{ComboID: 77, MainProductID: 20, ProductID: 21, QtyRatio: 1},
			},
		},
	}
	productRepo := &stubComboProductRepo{}
	productReader := &stubComboProductReader{
		products: []domain.Product{
			{ID: 20, SellerSku: "MAIN-20", Title: "Main 20"},
			{ID: 21, SellerSku: "CHILD-21", Title: "Child 21"},
			{ID: 22, SellerSku: "CHILD-22", Title: "Child 22"},
		},
	}

	uc := NewProductComboUsecase(comboRepo, productRepo, productReader, nil)

	params := domain.ComboUpsertParams{
		MainProductID: 20,
		Children: []domain.ComboChildInput{
			{ProductID: 21, QtyRatio: 2},
			{ProductID: 22, QtyRatio: 5},
		},
	}

	combo, err := uc.UpdateCombo(77, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if combo.ComboID != 77 {
		t.Fatalf("expected combo id 77, got %d", combo.ComboID)
	}
	if combo.MainProduct.ID != 20 {
		t.Fatalf("expected hydrated main product 20, got %d", combo.MainProduct.ID)
	}
	if len(combo.Products) != 2 || combo.Products[0].QtyRatio != 2 || combo.Products[1].QtyRatio != 5 {
		t.Fatalf("expected hydrated child ratios [2,5], got %+v", combo.Products)
	}

	if comboRepo.replaceComboID != 77 || comboRepo.replaceMainID != 20 {
		t.Fatalf("expected replace combo called with combo id 77")
	}

	if !reflect.DeepEqual(comboRepo.replaceProducts, []uint64{21, 22}) {
		t.Fatalf("expected replace products [21 22], got %v", comboRepo.replaceProducts)
	}

	if !reflect.DeepEqual(comboRepo.replaceRatios, map[uint64]uint64{21: 2, 22: 5}) {
		t.Fatalf("expected replace ratios map[21:2 22:5], got %v", comboRepo.replaceRatios)
	}

	if productRepo.clearedComboID != 77 {
		t.Fatalf("expected clear combo info for 77")
	}

	if !reflect.DeepEqual(productRepo.updatedIDs, []uint64{20, 21, 22}) {
		t.Fatalf("expected updated ids [20 21 22], got %v", productRepo.updatedIDs)
	}
}

func TestListCombosFiltersByLockedState(t *testing.T) {
	now := time.Now()
	comboRepo := &stubComboRepo{
		listComboIDs: []uint64{77, 88},
		itemsByCombo: map[uint64][]domain.ProductComboItem{
			77: {
				{ComboID: 77, MainProductID: 20, ProductID: 20, QtyRatio: 1, GmtCreate: now},
				{ComboID: 77, MainProductID: 20, ProductID: 21, QtyRatio: 2, GmtCreate: now},
			},
			88: {
				{ComboID: 88, MainProductID: 30, ProductID: 30, QtyRatio: 1, GmtCreate: now.Add(-time.Hour)},
				{ComboID: 88, MainProductID: 30, ProductID: 31, QtyRatio: 3, GmtCreate: now.Add(-time.Hour)},
			},
		},
	}
	productReader := &stubComboProductReader{
		products: []domain.Product{
			{ID: 20, SellerSku: "MAIN-20", Title: "Main 20", Marketplace: "US"},
			{ID: 21, SellerSku: "CHILD-21", Title: "Child 21", Marketplace: "US"},
			{ID: 30, SellerSku: "MAIN-30", Title: "Main 30", Marketplace: "US"},
			{ID: 31, SellerSku: "CHILD-31", Title: "Child 31", Marketplace: "US"},
		},
	}
	usageChecker := &stubComboUsageChecker{locked: true}
	uc := NewProductComboUsecase(comboRepo, nil, productReader, usageChecker)

	combos, total, err := uc.ListCombos(&domain.ComboListParams{Page: 1, PageSize: 20, Locked: "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	if len(combos) != 2 {
		t.Fatalf("expected 2 locked combos, got %d", len(combos))
	}
}

func TestUpdateComboUsesComboIDInsteadOfMainProductID(t *testing.T) {
	comboRepo := &stubComboRepo{
		itemsByCombo: map[uint64][]domain.ProductComboItem{
			77: {
				{ComboID: 77, MainProductID: 20, ProductID: 20, QtyRatio: 1},
				{ComboID: 77, MainProductID: 20, ProductID: 21, QtyRatio: 2},
			},
		},
	}
	productRepo := &stubComboProductRepo{}
	productReader := &stubComboProductReader{
		products: []domain.Product{
			{ID: 30, SellerSku: "MAIN-30", Title: "Main 30"},
			{ID: 31, SellerSku: "CHILD-31", Title: "Child 31"},
		},
	}

	uc := NewProductComboUsecase(comboRepo, productRepo, productReader, nil)

	combo, err := uc.UpdateCombo(77, domain.ComboUpsertParams{
		MainProductID: 30,
		Children: []domain.ComboChildInput{
			{ProductID: 31, QtyRatio: 6},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if comboRepo.replaceComboID != 77 {
		t.Fatalf("expected combo id 77, got %d", comboRepo.replaceComboID)
	}
	if comboRepo.replaceMainID != 30 {
		t.Fatalf("expected new main product id 30, got %d", comboRepo.replaceMainID)
	}
	if combo.MainProduct.ID != 30 {
		t.Fatalf("expected hydrated main product 30, got %d", combo.MainProduct.ID)
	}
	if len(combo.Products) != 1 || combo.Products[0].QtyRatio != 6 {
		t.Fatalf("expected hydrated child ratio 6, got %+v", combo.Products)
	}
}

func TestDeleteComboClearsProductInfo(t *testing.T) {
	comboRepo := &stubComboRepo{
		itemsByCombo: map[uint64][]domain.ProductComboItem{
			66: {
				{ComboID: 66, MainProductID: 10, ProductID: 10, QtyRatio: 1},
				{ComboID: 66, MainProductID: 10, ProductID: 11, QtyRatio: 2},
			},
		},
	}
	productRepo := &stubComboProductRepo{}

	uc := NewProductComboUsecase(comboRepo, productRepo, nil, nil)

	if err := uc.DeleteCombo(66); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if comboRepo.deletedComboID != 66 {
		t.Fatalf("expected delete combo 66")
	}

	if productRepo.clearedComboID != 66 {
		t.Fatalf("expected clear combo info 66")
	}
}

func TestGetComboIncludesQtyRatio(t *testing.T) {
	comboRepo := &stubComboRepo{
		itemsByCombo: map[uint64][]domain.ProductComboItem{
			77: {
				{ComboID: 77, MainProductID: 100, ProductID: 100, QtyRatio: 1},
				{ComboID: 77, MainProductID: 100, ProductID: 101, QtyRatio: 2},
				{ComboID: 77, MainProductID: 100, ProductID: 102, QtyRatio: 5},
			},
		},
	}
	productReader := &stubComboProductReader{
		products: []domain.Product{
			{ID: 100, SellerSku: "MAIN-100", Title: "Main"},
			{ID: 101, SellerSku: "CHILD-101", Title: "Child 101"},
			{ID: 102, SellerSku: "CHILD-102", Title: "Child 102"},
		},
	}

	uc := NewProductComboUsecase(comboRepo, nil, productReader, nil)

	combo, err := uc.GetCombo(77)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if combo.MainProduct.ID != 100 {
		t.Fatalf("expected main product 100, got %d", combo.MainProduct.ID)
	}
	if len(combo.Products) != 2 {
		t.Fatalf("expected 2 child products, got %d", len(combo.Products))
	}
	if combo.Products[0].ID != 101 || combo.Products[0].QtyRatio != 2 {
		t.Fatalf("expected child 101 ratio 2, got id=%d ratio=%d", combo.Products[0].ID, combo.Products[0].QtyRatio)
	}
	if combo.Products[1].ID != 102 || combo.Products[1].QtyRatio != 5 {
		t.Fatalf("expected child 102 ratio 5, got id=%d ratio=%d", combo.Products[1].ID, combo.Products[1].QtyRatio)
	}
}

func TestCreateComboRejectsNonStandaloneProducts(t *testing.T) {
	comboRepo := &stubComboRepo{createdComboID: 88}
	productRepo := &stubComboProductRepo{}
	existingComboID := uint64(50)
	productReader := &stubComboProductReader{
		products: []domain.Product{
			{ID: 10, SellerSku: "MAIN-10", Title: "Main 10"},
			{ID: 2, SellerSku: "CHILD-2", Title: "Child 2", ComboID: &existingComboID},
		},
	}

	uc := NewProductComboUsecase(comboRepo, productRepo, productReader, nil)

	_, err := uc.CreateCombo(domain.ComboUpsertParams{
		MainProductID: 10,
		Children: []domain.ComboChildInput{
			{ProductID: 2, QtyRatio: 3},
		},
	})
	if !errors.Is(err, ErrComboStandaloneOnly) {
		t.Fatalf("expected ErrComboStandaloneOnly, got %v", err)
	}
}

func TestCreateComboRejectsMissingMainProduct(t *testing.T) {
	uc := NewProductComboUsecase(&stubComboRepo{}, &stubComboProductRepo{}, &stubComboProductReader{}, nil)

	_, err := uc.CreateCombo(domain.ComboUpsertParams{})
	if !errors.Is(err, ErrComboMainProductRequired) {
		t.Fatalf("expected ErrComboMainProductRequired, got %v", err)
	}
}

func TestUpdateComboRejectsWhenLocked(t *testing.T) {
	now := time.Now()
	comboRepo := &stubComboRepo{
		itemsByCombo: map[uint64][]domain.ProductComboItem{
			77: {
				{ComboID: 77, MainProductID: 20, ProductID: 20, QtyRatio: 1, GmtCreate: now},
				{ComboID: 77, MainProductID: 20, ProductID: 21, QtyRatio: 2, GmtCreate: now},
			},
		},
	}
	productRepo := &stubComboProductRepo{}
	productReader := &stubComboProductReader{
		products: []domain.Product{
			{ID: 20, SellerSku: "MAIN-20", Title: "Main 20", ComboID: uint64Ptr(77)},
			{ID: 21, SellerSku: "CHILD-21", Title: "Child 21", ComboID: uint64Ptr(77)},
			{ID: 22, SellerSku: "CHILD-22", Title: "Child 22"},
		},
	}
	usageChecker := &stubComboUsageChecker{locked: true}

	uc := NewProductComboUsecase(comboRepo, productRepo, productReader, usageChecker)

	_, err := uc.UpdateCombo(77, domain.ComboUpsertParams{
		MainProductID: 20,
		Children: []domain.ComboChildInput{
			{ProductID: 21, QtyRatio: 2},
			{ProductID: 22, QtyRatio: 5},
		},
	})
	if !errors.Is(err, ErrComboLocked) {
		t.Fatalf("expected ErrComboLocked, got %v", err)
	}
}

func TestDeleteComboRejectsWhenLocked(t *testing.T) {
	now := time.Now()
	comboRepo := &stubComboRepo{
		itemsByCombo: map[uint64][]domain.ProductComboItem{
			77: {
				{ComboID: 77, MainProductID: 20, ProductID: 20, QtyRatio: 1, GmtCreate: now},
				{ComboID: 77, MainProductID: 20, ProductID: 21, QtyRatio: 2, GmtCreate: now},
			},
		},
	}
	productRepo := &stubComboProductRepo{}
	usageChecker := &stubComboUsageChecker{locked: true}

	uc := NewProductComboUsecase(comboRepo, productRepo, nil, usageChecker)

	err := uc.DeleteCombo(77)
	if !errors.Is(err, ErrComboLocked) {
		t.Fatalf("expected ErrComboLocked, got %v", err)
	}
}

func TestGetComboMarksLockedWhenBusinessUsageExists(t *testing.T) {
	now := time.Now()
	comboRepo := &stubComboRepo{
		itemsByCombo: map[uint64][]domain.ProductComboItem{
			77: {
				{ComboID: 77, MainProductID: 100, ProductID: 100, QtyRatio: 1, GmtCreate: now},
				{ComboID: 77, MainProductID: 100, ProductID: 101, QtyRatio: 2, GmtCreate: now},
			},
		},
	}
	productReader := &stubComboProductReader{
		products: []domain.Product{
			{ID: 100, SellerSku: "MAIN-100", Title: "Main"},
			{ID: 101, SellerSku: "CHILD-101", Title: "Child 101"},
		},
	}
	usageChecker := &stubComboUsageChecker{locked: true}

	uc := NewProductComboUsecase(comboRepo, nil, productReader, usageChecker)

	combo, err := uc.GetCombo(77)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !combo.Locked {
		t.Fatalf("expected combo locked")
	}
	if combo.LockReason == "" {
		t.Fatalf("expected lock reason")
	}
}

func uint64Ptr(v uint64) *uint64 {
	return &v
}
