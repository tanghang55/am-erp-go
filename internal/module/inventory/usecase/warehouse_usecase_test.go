package usecase

import (
	"context"
	"errors"
	"testing"

	"am-erp-go/internal/module/inventory/domain"
)

type stubWarehouseRepo struct {
	listItems        []*domain.Warehouse
	listTotal        int64
	getByID          *domain.Warehouse
	referenceByID    map[uint64]int64
	deleteCalled     bool
	activeWarehouses []*domain.Warehouse
}

func (s *stubWarehouseRepo) List(params *domain.WarehouseListParams) ([]*domain.Warehouse, int64, error) {
	return s.listItems, s.listTotal, nil
}

func (s *stubWarehouseRepo) GetByID(id uint64) (*domain.Warehouse, error) {
	return s.getByID, nil
}

func (s *stubWarehouseRepo) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	return nil
}
func (s *stubWarehouseRepo) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	return nil
}
func (s *stubWarehouseRepo) Delete(ctx context.Context, id uint64) error {
	s.deleteCalled = true
	return nil
}
func (s *stubWarehouseRepo) GetActiveWarehouses() ([]*domain.Warehouse, error) {
	return s.activeWarehouses, nil
}
func (s *stubWarehouseRepo) CountReferences(id uint64) (int64, error) {
	if s.referenceByID == nil {
		return 0, nil
	}
	return s.referenceByID[id], nil
}

func TestListWarehousesIncludesDeleteState(t *testing.T) {
	repo := &stubWarehouseRepo{
		listItems: []*domain.Warehouse{
			{ID: 1, Code: "WH001", Name: "引用中仓库"},
			{ID: 2, Code: "WH002", Name: "可删除仓库"},
		},
		listTotal: 2,
		referenceByID: map[uint64]int64{
			1: 5,
		},
	}
	uc := NewWarehouseUsecase(repo)

	items, total, err := uc.ListWarehouses(&domain.WarehouseListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	if items[0].Deletable {
		t.Fatalf("expected referenced warehouse to be non-deletable")
	}
	if items[0].ReferenceCount != 5 {
		t.Fatalf("expected reference count 5, got %d", items[0].ReferenceCount)
	}
	if items[1].Deletable != true {
		t.Fatalf("expected unused warehouse to be deletable")
	}
}

func TestDeleteWarehouseRejectsReferencedItem(t *testing.T) {
	repo := &stubWarehouseRepo{
		referenceByID: map[uint64]int64{1: 1},
	}
	uc := NewWarehouseUsecase(repo)

	err := uc.DeleteWarehouse(context.Background(), 1)
	if !errors.Is(err, ErrWarehouseReferenced) {
		t.Fatalf("expected ErrWarehouseReferenced, got %v", err)
	}
	if repo.deleteCalled {
		t.Fatalf("expected delete not called")
	}
}
