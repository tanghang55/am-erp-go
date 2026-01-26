package usecase

import (
	"am-erp-go/internal/module/shipment/domain"
)

type PackageSpecUseCase struct {
	repo domain.PackageSpecRepository
}

func NewPackageSpecUseCase(repo domain.PackageSpecRepository) *PackageSpecUseCase {
	return &PackageSpecUseCase{repo: repo}
}

func (uc *PackageSpecUseCase) Create(params *domain.CreatePackageSpecParams) (*domain.PackageSpec, error) {
	quantityPerBox := params.QuantityPerBox
	if quantityPerBox == 0 {
		quantityPerBox = 1
	}
	spec := &domain.PackageSpec{
		Name:           params.Name,
		Length:         params.Length,
		Width:          params.Width,
		Height:         params.Height,
		Weight:         params.Weight,
		QuantityPerBox: quantityPerBox,
		Remark:         params.Remark,
		Status:         "ACTIVE",
		CreatedBy:      params.CreatedBy,
	}

	if err := uc.repo.Create(spec); err != nil {
		return nil, err
	}

	return spec, nil
}

func (uc *PackageSpecUseCase) Update(id uint64, params *domain.UpdatePackageSpecParams) (*domain.PackageSpec, error) {
	spec, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if params.Name != nil {
		spec.Name = *params.Name
	}
	if params.Length != nil {
		spec.Length = *params.Length
	}
	if params.Width != nil {
		spec.Width = *params.Width
	}
	if params.Height != nil {
		spec.Height = *params.Height
	}
	if params.Weight != nil {
		spec.Weight = *params.Weight
	}
	if params.QuantityPerBox != nil {
		spec.QuantityPerBox = *params.QuantityPerBox
	}
	if params.Remark != nil {
		spec.Remark = params.Remark
	}
	if params.Status != nil {
		spec.Status = *params.Status
	}
	spec.UpdatedBy = params.UpdatedBy

	if err := uc.repo.Update(spec); err != nil {
		return nil, err
	}

	return spec, nil
}

func (uc *PackageSpecUseCase) GetByID(id uint64) (*domain.PackageSpec, error) {
	return uc.repo.GetByID(id)
}

func (uc *PackageSpecUseCase) List(params *domain.PackageSpecListParams) ([]*domain.PackageSpec, int64, error) {
	return uc.repo.List(params)
}

func (uc *PackageSpecUseCase) Delete(id uint64) error {
	return uc.repo.Delete(id)
}

func (uc *PackageSpecUseCase) ListByIDs(ids []uint64) ([]*domain.PackageSpec, error) {
	return uc.repo.ListByIDs(ids)
}
