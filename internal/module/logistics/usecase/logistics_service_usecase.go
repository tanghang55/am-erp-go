package usecase

import (
	"errors"
	"time"

	"am-erp-go/internal/module/logistics/domain"

	"gorm.io/gorm"
)

var (
	ErrServiceNotFound      = errors.New("logistics service not found")
	ErrServiceCodeDuplicate = errors.New("service code already exists")
)

type LogisticsServiceUsecase struct {
	serviceRepo domain.LogisticsServiceRepository
}

func NewLogisticsServiceUsecase(
	serviceRepo domain.LogisticsServiceRepository,
) *LogisticsServiceUsecase {
	return &LogisticsServiceUsecase{
		serviceRepo: serviceRepo,
	}
}

func (uc *LogisticsServiceUsecase) Create(params *domain.CreateLogisticsServiceParams) (*domain.LogisticsService, error) {
	// 检查服务代码是否已存在
	existingService, err := uc.serviceRepo.GetByCode(params.ServiceCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingService != nil {
		return nil, ErrServiceCodeDuplicate
	}

	now := time.Now()
	service := &domain.LogisticsService{
		ServiceCode:       params.ServiceCode,
		ServiceName:       params.ServiceName,
		TransportMode:     params.TransportMode,
		DestinationRegion: params.DestinationRegion,
		Description:       params.Description,
		Status:            params.Status,
		CreatedBy:         params.OperatorID,
		UpdatedBy:         params.OperatorID,
		GmtCreate:         now,
		GmtModified:       now,
	}

	if service.Status == "" {
		service.Status = domain.ServiceStatusActive
	}

	if err := uc.serviceRepo.Create(service); err != nil {
		return nil, err
	}

	return service, nil
}

func (uc *LogisticsServiceUsecase) Update(id uint64, params *domain.UpdateLogisticsServiceParams) error {
	service, err := uc.serviceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrServiceNotFound
		}
		return err
	}

	// 如果要修改服务代码，检查新代码是否已存在
	if params.ServiceCode != nil && *params.ServiceCode != service.ServiceCode {
		existingService, err := uc.serviceRepo.GetByCode(*params.ServiceCode)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if existingService != nil {
			return ErrServiceCodeDuplicate
		}
		service.ServiceCode = *params.ServiceCode
	}

	// 更新字段
	if params.ServiceName != nil {
		service.ServiceName = *params.ServiceName
	}
	if params.TransportMode != nil {
		service.TransportMode = *params.TransportMode
	}
	if params.DestinationRegion != nil {
		service.DestinationRegion = params.DestinationRegion
	}
	if params.Description != nil {
		service.Description = params.Description
	}
	if params.Status != nil {
		service.Status = *params.Status
	}
	service.UpdatedBy = params.OperatorID

	return uc.serviceRepo.Update(service)
}

func (uc *LogisticsServiceUsecase) Delete(id uint64) error {
	_, err := uc.serviceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrServiceNotFound
		}
		return err
	}

	return uc.serviceRepo.Delete(id)
}

func (uc *LogisticsServiceUsecase) Get(id uint64) (*domain.LogisticsService, error) {
	service, err := uc.serviceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}
	return service, nil
}

func (uc *LogisticsServiceUsecase) List(params *domain.LogisticsServiceListParams) ([]*domain.LogisticsService, int64, error) {
	return uc.serviceRepo.List(params)
}

func (uc *LogisticsServiceUsecase) GetActiveServices() ([]*domain.LogisticsService, error) {
	return uc.serviceRepo.GetActiveServices()
}

func (uc *LogisticsServiceUsecase) GetServicesByTransportMode(transportMode domain.TransportMode) ([]*domain.LogisticsService, error) {
	return uc.serviceRepo.GetServicesByTransportMode(transportMode)
}
