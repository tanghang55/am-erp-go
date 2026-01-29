package usecase

import (
	"errors"
	"time"

	"am-erp-go/internal/module/logistics/domain"

	"gorm.io/gorm"
)

var (
	ErrRateNotFound = errors.New("shipping rate not found")
)

type ShippingRateUsecase struct {
	rateRepo     domain.ShippingRateRepository
	providerRepo domain.LogisticsProviderRepository
}

func NewShippingRateUsecase(
	rateRepo domain.ShippingRateRepository,
	providerRepo domain.LogisticsProviderRepository,
) *ShippingRateUsecase {
	return &ShippingRateUsecase{
		rateRepo:     rateRepo,
		providerRepo: providerRepo,
	}
}

func (uc *ShippingRateUsecase) Create(params *domain.CreateShippingRateParams) (*domain.ShippingRate, error) {
	// 验证供应商是否存在
	_, err := uc.providerRepo.GetByID(params.ProviderID)
	if err != nil {
		return nil, ErrProviderNotFound
	}

	now := time.Now()
	rate := &domain.ShippingRate{
		ProviderID:             params.ProviderID,
		OriginWarehouseID:      params.OriginWarehouseID,
		DestinationWarehouseID: params.DestinationWarehouseID,
		TransportMode:          params.TransportMode,
		ServiceID:              params.ServiceID,
		PricingMethod:          params.PricingMethod,
		BaseRate:               params.BaseRate,
		Currency:               params.Currency,
		MinWeight:              params.MinWeight,
		TransitDays:            params.TransitDays,
		EffectiveDate:          params.EffectiveDate,
		ExpiryDate:             params.ExpiryDate,
		Status:                 params.Status,
		Remark:                 params.Remark,
		CreatedBy:              params.OperatorID,
		OtherFee:               params.OtherFee,
		UpdatedBy:              params.OperatorID,
		GmtCreate:              now,
		GmtModified:            now,
	}


	if rate.Currency == "" {
		rate.Currency = "CNY"
	}

	if err := uc.rateRepo.Create(rate); err != nil {
		return nil, err
	}

	return rate, nil
}

func (uc *ShippingRateUsecase) Update(id uint64, params *domain.UpdateShippingRateParams) error {
	rate, err := uc.rateRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRateNotFound
		}
		return err
	}

	// 验证供应商是否存在
	if params.ProviderID != nil {
		_, err := uc.providerRepo.GetByID(*params.ProviderID)
		if err != nil {
			return ErrProviderNotFound
		}
		rate.ProviderID = *params.ProviderID
	}

	// 更新字段
	if params.OriginWarehouseID != nil {
		rate.OriginWarehouseID = *params.OriginWarehouseID
	}
	if params.DestinationWarehouseID != nil {
		rate.DestinationWarehouseID = *params.DestinationWarehouseID
	}
	if params.TransportMode != nil {
		rate.TransportMode = *params.TransportMode
	}
	if params.ServiceID != nil {
		rate.ServiceID = params.ServiceID
	}
	if params.PricingMethod != nil {
		rate.PricingMethod = *params.PricingMethod
	}
	if params.BaseRate != nil {
		rate.BaseRate = *params.BaseRate
	}
	if params.Currency != nil {
		rate.Currency = *params.Currency
	}
	if params.TransitDays != nil {
		rate.TransitDays = params.TransitDays
	}
	if params.EffectiveDate != nil {
		rate.EffectiveDate = *params.EffectiveDate
	}
	if params.ExpiryDate != nil {
		rate.ExpiryDate = params.ExpiryDate
	}
	if params.Status != nil {
		rate.Status = *params.Status
	}
	if params.Remark != nil {
		rate.Remark = params.Remark
	}
	rate.UpdatedBy = params.OperatorID
	rate.OtherFee =  params.OtherFee

	return uc.rateRepo.Update(rate)
}

func (uc *ShippingRateUsecase) Delete(id uint64) error {
	_, err := uc.rateRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRateNotFound
		}
		return err
	}

	return uc.rateRepo.Delete(id)
}

func (uc *ShippingRateUsecase) Get(id uint64) (*domain.ShippingRate, error) {
	rate, err := uc.rateRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRateNotFound
		}
		return nil, err
	}
	return rate, nil
}

func (uc *ShippingRateUsecase) List(params *domain.ShippingRateListParams) ([]*domain.ShippingRate, int64, error) {
	return uc.rateRepo.List(params)
}

func (uc *ShippingRateUsecase) QueryLatestRate(params *domain.QueryLatestRateParams) (*domain.ShippingRate, error) {
	return uc.rateRepo.QueryLatestRate(params)
}
