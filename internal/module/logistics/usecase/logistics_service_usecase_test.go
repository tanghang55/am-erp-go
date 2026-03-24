package usecase

import (
	"errors"
	"testing"

	"am-erp-go/internal/module/logistics/domain"

	"gorm.io/gorm"
)

type stubLogisticsServiceRepo struct {
	createItem     *domain.LogisticsService
	updateItem     *domain.LogisticsService
	getByID        *domain.LogisticsService
	getByCode      *domain.LogisticsService
	getErr         error
	referenceCount int64
}

func (s *stubLogisticsServiceRepo) Create(service *domain.LogisticsService) error {
	copy := *service
	s.createItem = &copy
	return nil
}

func (s *stubLogisticsServiceRepo) Update(service *domain.LogisticsService) error {
	copy := *service
	s.updateItem = &copy
	return nil
}

func (s *stubLogisticsServiceRepo) Delete(id uint64) error { return nil }

func (s *stubLogisticsServiceRepo) GetByID(id uint64) (*domain.LogisticsService, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.getByID, nil
}

func (s *stubLogisticsServiceRepo) GetByCode(code string) (*domain.LogisticsService, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if s.getByCode == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return s.getByCode, nil
}

func (s *stubLogisticsServiceRepo) List(params *domain.LogisticsServiceListParams) ([]*domain.LogisticsService, int64, error) {
	return []*domain.LogisticsService{
		{ID: 1, ServiceCode: "SEA001", ServiceName: "引用中服务"},
		{ID: 2, ServiceCode: "AIR001", ServiceName: "可删除服务"},
	}, 2, nil
}

func (s *stubLogisticsServiceRepo) GetActiveServices() ([]*domain.LogisticsService, error) {
	return nil, nil
}

func (s *stubLogisticsServiceRepo) GetServicesByTransportMode(transportMode domain.TransportMode) ([]*domain.LogisticsService, error) {
	return nil, nil
}

func (s *stubLogisticsServiceRepo) CountReferences(id uint64) (int64, error) {
	if id == 1 {
		return s.referenceCount, nil
	}
	return 0, nil
}

func TestCreateLogisticsServiceRejectsInvalidCode(t *testing.T) {
	repo := &stubLogisticsServiceRepo{}
	uc := NewLogisticsServiceUsecase(repo)

	_, err := uc.Create(&domain.CreateLogisticsServiceParams{
		ServiceCode:   "慢船@001",
		ServiceName:   "Demo",
		TransportMode: domain.TransportModeSea,
		Status:        domain.ServiceStatusActive,
	})
	if !errors.Is(err, ErrServiceCodeInvalid) {
		t.Fatalf("expected ErrServiceCodeInvalid, got %v", err)
	}
	if repo.createItem != nil {
		t.Fatalf("expected create not called")
	}
}

func TestUpdateLogisticsServiceRejectsInvalidCode(t *testing.T) {
	repo := &stubLogisticsServiceRepo{
		getByID: &domain.LogisticsService{
			ID:          1,
			ServiceCode: "SEA001",
			ServiceName: "Demo",
		},
	}
	uc := NewLogisticsServiceUsecase(repo)
	invalid := "INVALID#001"

	err := uc.Update(1, &domain.UpdateLogisticsServiceParams{
		ServiceCode: &invalid,
	})
	if !errors.Is(err, ErrServiceCodeInvalid) {
		t.Fatalf("expected ErrServiceCodeInvalid, got %v", err)
	}
	if repo.updateItem != nil {
		t.Fatalf("expected update not called")
	}
}

func TestListLogisticsServicesIncludesDeleteState(t *testing.T) {
	repo := &stubLogisticsServiceRepo{referenceCount: 4}
	uc := NewLogisticsServiceUsecase(repo)

	items, total, err := uc.List(&domain.LogisticsServiceListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	if items[0].Deletable {
		t.Fatalf("expected referenced service to be non-deletable")
	}
	if items[0].ReferenceCount != 4 {
		t.Fatalf("expected reference count 4, got %d", items[0].ReferenceCount)
	}
	if items[1].Deletable != true {
		t.Fatalf("expected unused service to be deletable")
	}
}

func TestDeleteLogisticsServiceRejectsReferencedItem(t *testing.T) {
	repo := &stubLogisticsServiceRepo{
		getByID:        &domain.LogisticsService{ID: 1, ServiceCode: "SEA001"},
		referenceCount: 1,
	}
	uc := NewLogisticsServiceUsecase(repo)

	err := uc.Delete(1)
	if !errors.Is(err, ErrServiceReferenced) {
		t.Fatalf("expected ErrServiceReferenced, got %v", err)
	}
}
