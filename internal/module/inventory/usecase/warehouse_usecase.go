package usecase

import (
	"context"
	"fmt"
	"time"

	"am-erp-go/internal/module/inventory/domain"
)

type WarehouseUsecase struct {
	repo domain.WarehouseRepository
}

func NewWarehouseUsecase(repo domain.WarehouseRepository) *WarehouseUsecase {
	return &WarehouseUsecase{repo: repo}
}

func (u *WarehouseUsecase) ListWarehouses(params *domain.WarehouseListParams) ([]*domain.Warehouse, int64, error) {
	return u.repo.List(params)
}

func (u *WarehouseUsecase) GetWarehouse(id uint64) (*domain.Warehouse, error) {
	return u.repo.GetByID(id)
}

func (u *WarehouseUsecase) CreateWarehouse(ctx context.Context, params *domain.CreateWarehouseParams) (*domain.Warehouse, error) {
	warehouse := &domain.Warehouse{
		Code:          params.Code,
		Name:          params.Name,
		Type:          params.Type,
		Country:       params.Country,
		Address:       params.Address,
		ContactPerson: params.ContactPerson,
		ContactPhone:  params.ContactPhone,
		ContactEmail:  params.ContactEmail,
		Status:        params.Status,
		Remark:        params.Remark,
		GmtCreate:     time.Now(),
		GmtModified:   time.Now(),
	}

	if warehouse.Status == "" {
		warehouse.Status = domain.WarehouseStatusActive
	}

	if err := u.repo.Create(ctx, warehouse); err != nil {
		return nil, fmt.Errorf("failed to create warehouse: %w", err)
	}

	return warehouse, nil
}

func (u *WarehouseUsecase) UpdateWarehouse(ctx context.Context, id uint64, params *domain.UpdateWarehouseParams) (*domain.Warehouse, error) {
	warehouse, err := u.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("warehouse not found: %w", err)
	}

	if params.Code != nil {
		warehouse.Code = *params.Code
	}
	if params.Name != nil {
		warehouse.Name = *params.Name
	}
	if params.Type != nil {
		warehouse.Type = *params.Type
	}
	if params.Country != nil {
		warehouse.Country = params.Country
	}
	if params.Address != nil {
		warehouse.Address = params.Address
	}
	if params.ContactPerson != nil {
		warehouse.ContactPerson = params.ContactPerson
	}
	if params.ContactPhone != nil {
		warehouse.ContactPhone = params.ContactPhone
	}
	if params.ContactEmail != nil {
		warehouse.ContactEmail = params.ContactEmail
	}
	if params.Status != nil {
		warehouse.Status = *params.Status
	}
	if params.Remark != nil {
		warehouse.Remark = params.Remark
	}

	warehouse.GmtModified = time.Now()

	if err := u.repo.Update(ctx, warehouse); err != nil {
		return nil, fmt.Errorf("failed to update warehouse: %w", err)
	}

	return warehouse, nil
}

func (u *WarehouseUsecase) DeleteWarehouse(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, id)
}

func (u *WarehouseUsecase) GetActiveWarehouses() ([]*domain.Warehouse, error) {
	return u.repo.GetActiveWarehouses()
}
