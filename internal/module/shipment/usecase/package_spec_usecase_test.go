package usecase

import (
	"errors"
	"testing"

	"am-erp-go/internal/module/shipment/domain"
)

type stubPackageSpecRepo struct {
	listItems      []*domain.PackageSpec
	listTotal      int64
	getItem        *domain.PackageSpec
	deleteID       uint64
	referenceCount int64
	err            error
}

func (s *stubPackageSpecRepo) Create(_ *domain.PackageSpec) error { return s.err }
func (s *stubPackageSpecRepo) Update(_ *domain.PackageSpec) error { return s.err }
func (s *stubPackageSpecRepo) GetByID(_ uint64) (*domain.PackageSpec, error) {
	return s.getItem, s.err
}
func (s *stubPackageSpecRepo) List(_ *domain.PackageSpecListParams) ([]*domain.PackageSpec, int64, error) {
	return s.listItems, s.listTotal, s.err
}
func (s *stubPackageSpecRepo) Delete(id uint64) error {
	s.deleteID = id
	return s.err
}
func (s *stubPackageSpecRepo) ListByIDs(_ []uint64) ([]*domain.PackageSpec, error) {
	return s.listItems, s.err
}
func (s *stubPackageSpecRepo) CountReferences(_ uint64) (int64, error) {
	return s.referenceCount, s.err
}

type stubPackageSpecPackagingRepo struct{}

func (s *stubPackageSpecPackagingRepo) ListByPackageSpecID(_ uint64) ([]domain.PackageSpecPackagingItem, error) {
	return []domain.PackageSpecPackagingItem{}, nil
}

func (s *stubPackageSpecPackagingRepo) ReplaceAll(_ uint64, _ []domain.PackageSpecPackagingItem) error {
	return nil
}

func TestPackageSpecListIncludesDeleteState(t *testing.T) {
	repo := &stubPackageSpecRepo{
		listItems:      []*domain.PackageSpec{{ID: 1, Name: "标准箱"}},
		listTotal:      1,
		referenceCount: 2,
	}
	uc := NewPackageSpecUseCase(repo, &stubPackageSpecPackagingRepo{})

	items, _, err := uc.List(&domain.PackageSpecListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one item")
	}
	if items[0].Deletable {
		t.Fatalf("expected deletable false")
	}
	if items[0].ReferenceCount != 2 {
		t.Fatalf("expected reference_count 2, got %d", items[0].ReferenceCount)
	}
	if items[0].DeleteBlockReason == "" {
		t.Fatalf("expected delete block reason")
	}
}

func TestPackageSpecDeleteRejectsReferencedSpec(t *testing.T) {
	repo := &stubPackageSpecRepo{referenceCount: 1}
	uc := NewPackageSpecUseCase(repo, &stubPackageSpecPackagingRepo{})

	err := uc.Delete(3)
	if !errors.Is(err, ErrPackageSpecReferenced) {
		t.Fatalf("expected ErrPackageSpecReferenced, got %v", err)
	}
	if repo.deleteID != 0 {
		t.Fatalf("expected delete not called")
	}
}
