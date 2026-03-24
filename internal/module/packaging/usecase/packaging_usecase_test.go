package usecase

import (
	"errors"
	"testing"
	"time"

	"am-erp-go/internal/module/packaging/domain"
)

type stubPackagingItemValidationRepo struct {
	created  *domain.PackagingItem
	updated  *domain.PackagingItem
	items    []domain.PackagingItem
	total    int64
	countMap map[uint64]int64
}

func (s *stubPackagingItemValidationRepo) List(params *domain.PackagingItemListParams) ([]domain.PackagingItem, int64, error) {
	if s.items != nil {
		return s.items, s.total, nil
	}
	return nil, 0, nil
}

func (s *stubPackagingItemValidationRepo) GetByID(id uint64) (*domain.PackagingItem, error) {
	for _, item := range s.items {
		if item.ID == id {
			copy := item
			return &copy, nil
		}
	}
	return nil, errors.New("not implemented")
}

func (s *stubPackagingItemValidationRepo) Create(item *domain.PackagingItem) error {
	copy := *item
	s.created = &copy
	return nil
}

func (s *stubPackagingItemValidationRepo) Update(item *domain.PackagingItem) error {
	copy := *item
	s.updated = &copy
	return nil
}

func (s *stubPackagingItemValidationRepo) Delete(id uint64) error { return nil }
func (s *stubPackagingItemValidationRepo) CountReferences(id uint64) (int64, error) {
	if s.countMap == nil {
		return 0, nil
	}
	return s.countMap[id], nil
}
func (s *stubPackagingItemValidationRepo) GetLowStockItems() ([]domain.PackagingItem, error) {
	return nil, nil
}
func (s *stubPackagingItemValidationRepo) UpdateQuantity(id uint64, quantity int64) error { return nil }

type stubPackagingLedgerValidationRepo struct{}

func (s *stubPackagingLedgerValidationRepo) List(params *domain.PackagingLedgerListParams) ([]domain.PackagingLedger, int64, error) {
	return nil, 0, nil
}
func (s *stubPackagingLedgerValidationRepo) GetByID(id uint64) (*domain.PackagingLedger, error) {
	return nil, errors.New("not implemented")
}
func (s *stubPackagingLedgerValidationRepo) Create(ledger *domain.PackagingLedger) error { return nil }
func (s *stubPackagingLedgerValidationRepo) GetUsageSummary(dateFrom, dateTo *time.Time) ([]domain.UsageSummaryItem, error) {
	return nil, nil
}

func TestCreateItemRejectsMissingSupplier(t *testing.T) {
	repo := &stubPackagingItemValidationRepo{}
	uc := NewPackagingUsecase(repo, &stubPackagingLedgerValidationRepo{})

	err := uc.CreateItem(&domain.PackagingItem{
		ItemCode: "PKG-1",
		ItemName: "Package 1",
		Category: "BOX",
	})
	if !errors.Is(err, ErrPackagingSupplierRequired) {
		t.Fatalf("expected ErrPackagingSupplierRequired, got %v", err)
	}
	if repo.created != nil {
		t.Fatalf("expected create not called")
	}
}

func TestUpdateItemRejectsMissingSupplier(t *testing.T) {
	repo := &stubPackagingItemValidationRepo{}
	uc := NewPackagingUsecase(repo, &stubPackagingLedgerValidationRepo{})

	err := uc.UpdateItem(&domain.PackagingItem{
		ID:       1,
		ItemCode: "PKG-1",
		ItemName: "Package 1",
		Category: "BOX",
	})
	if !errors.Is(err, ErrPackagingSupplierRequired) {
		t.Fatalf("expected ErrPackagingSupplierRequired, got %v", err)
	}
	if repo.updated != nil {
		t.Fatalf("expected update not called")
	}
}

func TestCreateItemRejectsInvalidCode(t *testing.T) {
	supplierID := uint64(1)
	repo := &stubPackagingItemValidationRepo{}
	uc := NewPackagingUsecase(repo, &stubPackagingLedgerValidationRepo{})

	err := uc.CreateItem(&domain.PackagingItem{
		ItemCode:   "包材@001",
		ItemName:   "Package 1",
		Category:   "BOX",
		SupplierID: &supplierID,
	})
	if !errors.Is(err, ErrPackagingItemCodeInvalid) {
		t.Fatalf("expected ErrPackagingItemCodeInvalid, got %v", err)
	}
	if repo.created != nil {
		t.Fatalf("expected create not called")
	}
}

func TestUpdateItemRejectsInvalidCode(t *testing.T) {
	supplierID := uint64(1)
	repo := &stubPackagingItemValidationRepo{}
	uc := NewPackagingUsecase(repo, &stubPackagingLedgerValidationRepo{})

	err := uc.UpdateItem(&domain.PackagingItem{
		ID:         1,
		ItemCode:   "PKG#001",
		ItemName:   "Package 1",
		Category:   "BOX",
		SupplierID: &supplierID,
	})
	if !errors.Is(err, ErrPackagingItemCodeInvalid) {
		t.Fatalf("expected ErrPackagingItemCodeInvalid, got %v", err)
	}
	if repo.updated != nil {
		t.Fatalf("expected update not called")
	}
}

func TestListItemsIncludesDeleteState(t *testing.T) {
	repo := &stubPackagingItemValidationRepo{
		items: []domain.PackagingItem{
			{ID: 1, ItemCode: "BOX-001", ItemName: "Box", Category: "BOX"},
			{ID: 2, ItemCode: "TAPE-001", ItemName: "Tape", Category: "TAPE"},
		},
		total: 2,
		countMap: map[uint64]int64{
			1: 3,
			2: 0,
		},
	}
	uc := NewPackagingUsecase(repo, &stubPackagingLedgerValidationRepo{})

	items, total, err := uc.ListItems(&domain.PackagingItemListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	if items[0].Deletable {
		t.Fatalf("expected referenced packaging item to be non-deletable")
	}
	if items[0].ReferenceCount != 3 {
		t.Fatalf("expected reference_count 3, got %d", items[0].ReferenceCount)
	}
	if !items[1].Deletable {
		t.Fatalf("expected unused packaging item to be deletable")
	}
}

func TestDeleteItemRejectsReferencedItem(t *testing.T) {
	repo := &stubPackagingItemValidationRepo{
		countMap: map[uint64]int64{
			5: 2,
		},
	}
	uc := NewPackagingUsecase(repo, &stubPackagingLedgerValidationRepo{})

	err := uc.DeleteItem(5)
	if err == nil {
		t.Fatalf("expected referenced packaging item delete to fail")
	}
}
