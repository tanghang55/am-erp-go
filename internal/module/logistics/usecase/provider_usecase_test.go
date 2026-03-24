package usecase

import (
	"errors"
	"testing"

	"am-erp-go/internal/module/logistics/domain"

	"gorm.io/gorm"
)

type stubProviderRepo struct {
	createItem     *domain.LogisticsProvider
	updateItem     *domain.LogisticsProvider
	getByID        *domain.LogisticsProvider
	getByCode      *domain.LogisticsProvider
	getByErr       error
	referenceCount int64
}

func (s *stubProviderRepo) Create(provider *domain.LogisticsProvider) error {
	copy := *provider
	s.createItem = &copy
	return nil
}

func (s *stubProviderRepo) Update(provider *domain.LogisticsProvider) error {
	copy := *provider
	s.updateItem = &copy
	return nil
}

func (s *stubProviderRepo) Delete(id uint64) error { return nil }

func (s *stubProviderRepo) GetByID(id uint64) (*domain.LogisticsProvider, error) {
	if s.getByErr != nil {
		return nil, s.getByErr
	}
	return s.getByID, nil
}

func (s *stubProviderRepo) GetByCode(code string) (*domain.LogisticsProvider, error) {
	if s.getByErr != nil {
		return nil, s.getByErr
	}
	if s.getByCode == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return s.getByCode, nil
}

func (s *stubProviderRepo) List(params *domain.LogisticsProviderListParams) ([]*domain.LogisticsProvider, int64, error) {
	return []*domain.LogisticsProvider{
		{ID: 1, ProviderCode: "LP001", ProviderName: "引用中物流商"},
		{ID: 2, ProviderCode: "LP002", ProviderName: "可删除物流商"},
	}, 2, nil
}

func (s *stubProviderRepo) CountReferences(id uint64) (int64, error) {
	if id == 1 {
		return s.referenceCount, nil
	}
	return 0, nil
}

func TestCreateProviderRejectsInvalidCode(t *testing.T) {
	repo := &stubProviderRepo{}
	uc := NewLogisticsProviderUsecase(repo)

	_, err := uc.Create(&domain.CreateProviderParams{
		ProviderCode: "物流@001",
		ProviderName: "Demo",
		ProviderType: domain.ProviderTypeCourier,
		Status:       domain.ProviderStatusActive,
	})
	if !errors.Is(err, ErrProviderCodeInvalid) {
		t.Fatalf("expected ErrProviderCodeInvalid, got %v", err)
	}
	if repo.createItem != nil {
		t.Fatalf("expected create not called")
	}
}

func TestUpdateProviderRejectsInvalidCode(t *testing.T) {
	existing := &domain.LogisticsProvider{
		ID:           1,
		ProviderCode: "LP001",
		ProviderName: "Demo",
	}
	repo := &stubProviderRepo{getByID: existing}
	uc := NewLogisticsProviderUsecase(repo)
	invalid := "INVALID#001"

	err := uc.Update(1, &domain.UpdateProviderParams{
		ProviderCode: &invalid,
	})
	if !errors.Is(err, ErrProviderCodeInvalid) {
		t.Fatalf("expected ErrProviderCodeInvalid, got %v", err)
	}
	if repo.updateItem != nil {
		t.Fatalf("expected update not called")
	}
}

func TestListProvidersIncludesDeleteState(t *testing.T) {
	repo := &stubProviderRepo{referenceCount: 3}
	uc := NewLogisticsProviderUsecase(repo)

	items, total, err := uc.List(&domain.LogisticsProviderListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	if items[0].Deletable {
		t.Fatalf("expected referenced provider to be non-deletable")
	}
	if items[0].ReferenceCount != 3 {
		t.Fatalf("expected reference count 3, got %d", items[0].ReferenceCount)
	}
	if items[1].Deletable != true {
		t.Fatalf("expected unused provider to be deletable")
	}
}

func TestDeleteProviderRejectsReferencedItem(t *testing.T) {
	repo := &stubProviderRepo{
		getByID:        &domain.LogisticsProvider{ID: 1, ProviderCode: "LP001"},
		referenceCount: 2,
	}
	uc := NewLogisticsProviderUsecase(repo)

	err := uc.Delete(1)
	if !errors.Is(err, ErrProviderReferenced) {
		t.Fatalf("expected ErrProviderReferenced, got %v", err)
	}
}
