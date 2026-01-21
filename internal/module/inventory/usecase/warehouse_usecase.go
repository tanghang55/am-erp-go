package usecase

import "am-erp-go/internal/module/inventory/domain"

type WarehouseUsecase struct {
	repo domain.WarehouseRepository
}

func NewWarehouseUsecase(repo domain.WarehouseRepository) *WarehouseUsecase {
	return &WarehouseUsecase{repo: repo}
}

func (u *WarehouseUsecase) GetActiveWarehouses() ([]*domain.Warehouse, error) {
	return u.repo.GetActiveWarehouses()
}
