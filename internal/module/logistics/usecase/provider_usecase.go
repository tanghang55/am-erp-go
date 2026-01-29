package usecase

import (
	"errors"
	"time"

	"am-erp-go/internal/module/logistics/domain"

	"gorm.io/gorm"
)

var (
	ErrProviderNotFound    = errors.New("logistics provider not found")
	ErrProviderCodeExists  = errors.New("provider code already exists")
	ErrProviderHasRates    = errors.New("cannot delete provider with active shipping rates")
)

type LogisticsProviderUsecase struct {
	providerRepo domain.LogisticsProviderRepository
}

func NewLogisticsProviderUsecase(providerRepo domain.LogisticsProviderRepository) *LogisticsProviderUsecase {
	return &LogisticsProviderUsecase{
		providerRepo: providerRepo,
	}
}

func (uc *LogisticsProviderUsecase) Create(params *domain.CreateProviderParams) (*domain.LogisticsProvider, error) {
	// 检查代码是否已存在
	existing, err := uc.providerRepo.GetByCode(params.ProviderCode)
	if err == nil && existing != nil {
		return nil, ErrProviderCodeExists
	}

	now := time.Now()
	provider := &domain.LogisticsProvider{
		ProviderCode:  params.ProviderCode,
		ProviderName:  params.ProviderName,
		ProviderType:  params.ProviderType,
		ServiceTypes:  params.ServiceTypes,
		ContactPerson: params.ContactPerson,
		ContactPhone:  params.ContactPhone,
		ContactEmail:  params.ContactEmail,
		Website:       params.Website,
		Country:       params.Country,
		City:          params.City,
		Address:       params.Address,
		AccountNumber: params.AccountNumber,
		CreditDays:    0,
		Status:        params.Status,
		Remark:        params.Remark,
		CreatedBy:     params.OperatorID,
		UpdatedBy:     params.OperatorID,
		GmtCreate:     now,
		GmtModified:   now,
	}

	if params.CreditDays != nil {
		provider.CreditDays = *params.CreditDays
	}

	if err := uc.providerRepo.Create(provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (uc *LogisticsProviderUsecase) Update(id uint64, params *domain.UpdateProviderParams) error {
	provider, err := uc.providerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProviderNotFound
		}
		return err
	}

	// 检查代码是否被其他记录占用
	if params.ProviderCode != nil && *params.ProviderCode != provider.ProviderCode {
		existing, err := uc.providerRepo.GetByCode(*params.ProviderCode)
		if err == nil && existing != nil && existing.ID != id {
			return ErrProviderCodeExists
		}
	}

	// 更新字段
	if params.ProviderCode != nil {
		provider.ProviderCode = *params.ProviderCode
	}
	if params.ProviderName != nil {
		provider.ProviderName = *params.ProviderName
	}
	if params.ProviderType != nil {
		provider.ProviderType = *params.ProviderType
	}
	if params.ServiceTypes != nil {
		provider.ServiceTypes = params.ServiceTypes
	}
	if params.ContactPerson != nil {
		provider.ContactPerson = params.ContactPerson
	}
	if params.ContactPhone != nil {
		provider.ContactPhone = params.ContactPhone
	}
	if params.ContactEmail != nil {
		provider.ContactEmail = params.ContactEmail
	}
	if params.Website != nil {
		provider.Website = params.Website
	}
	if params.Country != nil {
		provider.Country = params.Country
	}
	if params.City != nil {
		provider.City = params.City
	}
	if params.Address != nil {
		provider.Address = params.Address
	}
	if params.AccountNumber != nil {
		provider.AccountNumber = params.AccountNumber
	}
	if params.CreditDays != nil {
		provider.CreditDays = *params.CreditDays
	}
	if params.Status != nil {
		provider.Status = *params.Status
	}
	if params.Remark != nil {
		provider.Remark = params.Remark
	}
	provider.UpdatedBy = params.OperatorID

	return uc.providerRepo.Update(provider)
}

func (uc *LogisticsProviderUsecase) Delete(id uint64) error {
	_, err := uc.providerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProviderNotFound
		}
		return err
	}

	// TODO: 检查是否有关联的运费报价
	// 如果有未过期的报价，不允许删除

	return uc.providerRepo.Delete(id)
}

func (uc *LogisticsProviderUsecase) Get(id uint64) (*domain.LogisticsProvider, error) {
	provider, err := uc.providerRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProviderNotFound
		}
		return nil, err
	}
	return provider, nil
}

func (uc *LogisticsProviderUsecase) List(params *domain.LogisticsProviderListParams) ([]*domain.LogisticsProvider, int64, error) {
	return uc.providerRepo.List(params)
}
