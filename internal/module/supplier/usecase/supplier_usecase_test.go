package usecase

import (
	"reflect"
	"testing"

	"am-erp-go/internal/module/supplier/domain"
)

type stubSupplierRepo struct {
	listParams *domain.SupplierListParams
	listItems  []domain.Supplier
	listTotal  int64
	getID      uint64
	getItem    *domain.Supplier
	createItem *domain.Supplier
	updateItem *domain.Supplier
	deleteID   uint64
	err        error
}

func (s *stubSupplierRepo) List(params *domain.SupplierListParams) ([]domain.Supplier, int64, error) {
	s.listParams = params
	return s.listItems, s.listTotal, s.err
}

func (s *stubSupplierRepo) GetByID(id uint64) (*domain.Supplier, error) {
	s.getID = id
	return s.getItem, s.err
}

func (s *stubSupplierRepo) Create(supplier *domain.Supplier) error {
	s.createItem = supplier
	if supplier.ID == 0 {
		supplier.ID = 9
	}
	return s.err
}

func (s *stubSupplierRepo) Update(supplier *domain.Supplier) error {
	s.updateItem = supplier
	return s.err
}

func (s *stubSupplierRepo) Delete(id uint64) error {
	s.deleteID = id
	return s.err
}

type stubSupplierTypeRepo struct {
	listID    uint64
	listIDs   []uint64
	listMap   map[uint64][]string
	list      []string
	replace   []string
	replaceID uint64
	err       error
}

func (s *stubSupplierTypeRepo) ListBySupplierID(id uint64) ([]string, error) {
	s.listID = id
	return s.list, s.err
}

func (s *stubSupplierTypeRepo) ListBySupplierIDs(ids []uint64) (map[uint64][]string, error) {
	s.listIDs = ids
	return s.listMap, s.err
}

func (s *stubSupplierTypeRepo) ReplaceBySupplierID(id uint64, types []string) error {
	s.replaceID = id
	s.replace = types
	return s.err
}

type stubSupplierContactRepo struct {
	listID           uint64
	list             []domain.SupplierContact
	create           *domain.SupplierContact
	update           *domain.SupplierContact
	deleteID         uint64
	deleteSupplierID uint64
	err              error
}

func (s *stubSupplierContactRepo) ListBySupplierID(id uint64) ([]domain.SupplierContact, error) {
	s.listID = id
	return s.list, s.err
}

func (s *stubSupplierContactRepo) Create(contact *domain.SupplierContact) error {
	s.create = contact
	return s.err
}

func (s *stubSupplierContactRepo) Update(contact *domain.SupplierContact) error {
	s.update = contact
	return s.err
}

func (s *stubSupplierContactRepo) Delete(id uint64, supplierID uint64) error {
	s.deleteID = id
	s.deleteSupplierID = supplierID
	return s.err
}

type stubSupplierAccountRepo struct {
	listID           uint64
	list             []domain.SupplierAccount
	create           *domain.SupplierAccount
	update           *domain.SupplierAccount
	deleteID         uint64
	deleteSupplierID uint64
	err              error
}

func (s *stubSupplierAccountRepo) ListBySupplierID(id uint64) ([]domain.SupplierAccount, error) {
	s.listID = id
	return s.list, s.err
}

func (s *stubSupplierAccountRepo) Create(account *domain.SupplierAccount) error {
	s.create = account
	return s.err
}

func (s *stubSupplierAccountRepo) Update(account *domain.SupplierAccount) error {
	s.update = account
	return s.err
}

func (s *stubSupplierAccountRepo) Delete(id uint64, supplierID uint64) error {
	s.deleteID = id
	s.deleteSupplierID = supplierID
	return s.err
}

type stubSupplierTagRepo struct {
	listID           uint64
	list             []domain.SupplierTag
	create           *domain.SupplierTag
	update           *domain.SupplierTag
	deleteID         uint64
	deleteSupplierID uint64
	err              error
}

func (s *stubSupplierTagRepo) ListBySupplierID(id uint64) ([]domain.SupplierTag, error) {
	s.listID = id
	return s.list, s.err
}

func (s *stubSupplierTagRepo) Create(tag *domain.SupplierTag) error {
	s.create = tag
	return s.err
}

func (s *stubSupplierTagRepo) Update(tag *domain.SupplierTag) error {
	s.update = tag
	return s.err
}

func (s *stubSupplierTagRepo) Delete(id uint64, supplierID uint64) error {
	s.deleteID = id
	s.deleteSupplierID = supplierID
	return s.err
}

func TestListSuppliersMapsTypes(t *testing.T) {
	supplierRepo := &stubSupplierRepo{
		listItems: []domain.Supplier{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}},
		listTotal: 2,
	}
	typeRepo := &stubSupplierTypeRepo{
		listMap: map[uint64][]string{1: {"PRODUCT"}},
	}

	uc := NewSupplierUsecase(supplierRepo, typeRepo, nil, nil, nil)

	items, _, err := uc.ListSuppliers(&domain.SupplierListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 || !reflect.DeepEqual(items[0].Types, []string{"PRODUCT"}) {
		t.Fatalf("expected types mapped")
	}
}

func TestGetSupplierDetailAggregatesSubEntities(t *testing.T) {
	supplierRepo := &stubSupplierRepo{
		getItem: &domain.Supplier{ID: 7, Name: "Demo"},
	}
	typeRepo := &stubSupplierTypeRepo{list: []string{"LOGISTICS"}}
	contactRepo := &stubSupplierContactRepo{list: []domain.SupplierContact{{ID: 1, Name: "Kate"}}}
	accountRepo := &stubSupplierAccountRepo{list: []domain.SupplierAccount{{ID: 2, BankName: "HSBC"}}}
	tagRepo := &stubSupplierTagRepo{list: []domain.SupplierTag{{ID: 3, Tag: "VIP"}}}

	uc := NewSupplierUsecase(supplierRepo, typeRepo, contactRepo, accountRepo, tagRepo)

	detail, err := uc.GetSupplier(7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if detail == nil || detail.ID != 7 || len(detail.Types) != 1 || len(detail.Contacts) != 1 || len(detail.Accounts) != 1 || len(detail.Tags) != 1 {
		t.Fatalf("expected detail aggregation")
	}
}

func TestCreateSupplierReplacesTypes(t *testing.T) {
	supplierRepo := &stubSupplierRepo{}
	typeRepo := &stubSupplierTypeRepo{}

	uc := NewSupplierUsecase(supplierRepo, typeRepo, nil, nil, nil)
	supplier := &domain.Supplier{Name: "Demo"}

	created, err := uc.CreateSupplier(supplier, []string{"PRODUCT", "PACKAGING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created == nil || created.ID == 0 || typeRepo.replaceID != created.ID || len(typeRepo.replace) != 2 {
		t.Fatalf("expected types replaced")
	}
}
