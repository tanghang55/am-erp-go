package usecase

import (
	"reflect"
	"testing"

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

	comboIDByMain uint64

	deletedComboID uint64
}

func (s *stubComboRepo) ListComboIDs(params *domain.ComboListParams) ([]uint64, int64, error) {
	return nil, 0, nil
}

func (s *stubComboRepo) GetItemsByComboID(comboID uint64) ([]domain.ProductComboItem, error) {
	return nil, nil
}

func (s *stubComboRepo) GetComboIDByMainProductID(mainProductID uint64) (uint64, error) {
	return s.comboIDByMain, nil
}

func (s *stubComboRepo) CreateCombo(mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) (uint64, error) {
	s.createdMainID = mainProductID
	s.createdProducts = append([]uint64{}, productIDs...)
	s.createdRatios = qtyRatios
	if s.createdComboID == 0 {
		s.createdComboID = 99
	}
	return s.createdComboID, nil
}

func (s *stubComboRepo) ReplaceComboItems(comboID uint64, mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) error {
	s.replaceComboID = comboID
	s.replaceMainID = mainProductID
	s.replaceProducts = append([]uint64{}, productIDs...)
	s.replaceRatios = qtyRatios
	return nil
}

func (s *stubComboRepo) DeleteCombo(comboID uint64) error {
	s.deletedComboID = comboID
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

func (s *stubComboProductRepo) ClearComboInfo(comboID uint64) error {
	s.clearedComboID = comboID
	return nil
}

func TestCreateComboUpdatesProductComboInfo(t *testing.T) {
	comboRepo := &stubComboRepo{createdComboID: 88}
	productRepo := &stubComboProductRepo{}

	uc := NewProductComboUsecase(comboRepo, productRepo, nil)

	params := domain.ComboUpsertParams{
		MainProductID: 10,
		ProductIDs:    []uint64{10, 11, 12, 11},
	}

	combo, err := uc.CreateCombo(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if combo.ComboID != 88 {
		t.Fatalf("expected combo id 88, got %d", combo.ComboID)
	}

	if comboRepo.createdMainID != 10 {
		t.Fatalf("expected main product id 10")
	}

	if !reflect.DeepEqual(comboRepo.createdProducts, []uint64{11, 12}) {
		t.Fatalf("expected child products [11 12], got %v", comboRepo.createdProducts)
	}

	if productRepo.updatedComboID != 88 || productRepo.updatedMainID != 10 {
		t.Fatalf("expected update combo info called")
	}

	if !reflect.DeepEqual(productRepo.updatedIDs, []uint64{10, 11, 12}) {
		t.Fatalf("expected updated ids [10 11 12], got %v", productRepo.updatedIDs)
	}
}

func TestUpdateComboReplacesItemsAndClearsOld(t *testing.T) {
	comboRepo := &stubComboRepo{comboIDByMain: 77}
	productRepo := &stubComboProductRepo{}

	uc := NewProductComboUsecase(comboRepo, productRepo, nil)

	params := domain.ComboUpsertParams{
		MainProductID: 20,
		ProductIDs:    []uint64{21, 22},
	}

	combo, err := uc.UpdateComboByMainProductID(20, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if combo.ComboID != 77 {
		t.Fatalf("expected combo id 77, got %d", combo.ComboID)
	}

	if comboRepo.replaceComboID != 77 || comboRepo.replaceMainID != 20 {
		t.Fatalf("expected replace combo called with combo id 77")
	}

	if !reflect.DeepEqual(comboRepo.replaceProducts, []uint64{21, 22}) {
		t.Fatalf("expected replace products [21 22], got %v", comboRepo.replaceProducts)
	}

	if productRepo.clearedComboID != 77 {
		t.Fatalf("expected clear combo info for 77")
	}

	if !reflect.DeepEqual(productRepo.updatedIDs, []uint64{20, 21, 22}) {
		t.Fatalf("expected updated ids [20 21 22], got %v", productRepo.updatedIDs)
	}
}

func TestDeleteComboClearsProductInfo(t *testing.T) {
	comboRepo := &stubComboRepo{}
	productRepo := &stubComboProductRepo{}

	uc := NewProductComboUsecase(comboRepo, productRepo, nil)

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
