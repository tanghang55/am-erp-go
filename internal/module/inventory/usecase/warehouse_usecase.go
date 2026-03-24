package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/validation"
	"am-erp-go/internal/module/inventory/domain"
)

var ErrWarehouseCodeInvalid = errors.New("warehouse code only supports letters, numbers, hyphen and underscore")
var ErrWarehouseReferenced = errors.New("warehouse is still referenced by business data")

type WarehouseUsecase struct {
	repo domain.WarehouseRepository
}

func NewWarehouseUsecase(repo domain.WarehouseRepository) *WarehouseUsecase {
	return &WarehouseUsecase{repo: repo}
}

func (u *WarehouseUsecase) ListWarehouses(params *domain.WarehouseListParams) ([]*domain.Warehouse, int64, error) {
	warehouses, total, err := u.repo.List(params)
	if err != nil {
		return nil, 0, err
	}
	for _, warehouse := range warehouses {
		warehouse.Deletable = true
		refCount, err := u.repo.CountReferences(warehouse.ID)
		if err != nil {
			return nil, 0, err
		}
		warehouse.ReferenceCount = refCount
		if refCount > 0 {
			warehouse.Deletable = false
			warehouse.DeleteBlockReason = "已被业务数据引用，不可删除"
		}
	}
	return warehouses, total, nil
}

func (u *WarehouseUsecase) GetWarehouse(id uint64) (*domain.Warehouse, error) {
	return u.repo.GetByID(id)
}

func (u *WarehouseUsecase) CreateWarehouse(ctx context.Context, params *domain.CreateWarehouseParams) (*domain.Warehouse, error) {
	if !validation.IsValidCode(strings.TrimSpace(params.Code)) {
		return nil, ErrWarehouseCodeInvalid
	}
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
		if !validation.IsValidCode(strings.TrimSpace(*params.Code)) {
			return nil, ErrWarehouseCodeInvalid
		}
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
	refCount, err := u.repo.CountReferences(id)
	if err != nil {
		return err
	}
	if refCount > 0 {
		return ErrWarehouseReferenced
	}
	return u.repo.Delete(ctx, id)
}

func (u *WarehouseUsecase) GetActiveWarehouses() ([]*domain.Warehouse, error) {
	return u.repo.GetActiveWarehouses()
}
