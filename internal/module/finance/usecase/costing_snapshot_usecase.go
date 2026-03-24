package usecase

import (
	"am-erp-go/internal/module/finance/domain"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CostingSnapshotUsecase struct {
	repo domain.CostingSnapshotRepository
}

func NewCostingSnapshotUsecase(repo domain.CostingSnapshotRepository) *CostingSnapshotUsecase {
	return &CostingSnapshotUsecase{repo: repo}
}

type CreateCostingSnapshotInput struct {
	ProductID     uint64
	CostType      domain.CostType
	UnitCost      float64
	Currency      string
	EffectiveFrom *time.Time
	EffectiveTo   *time.Time
	Notes         *string
	CreatedBy     uint64
}

type UpdateCostingSnapshotInput struct {
	UnitCost      *float64
	Currency      *string
	EffectiveFrom *time.Time
	EffectiveTo   **time.Time
	Notes         **string
}

func (uc *CostingSnapshotUsecase) List(params *domain.CostingSnapshotListParams) ([]domain.CostingSnapshot, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.repo.List(params)
}

func (uc *CostingSnapshotUsecase) Get(id uint64) (*domain.CostingSnapshot, error) {
	return uc.repo.GetByID(id)
}

func (uc *CostingSnapshotUsecase) Create(input *CreateCostingSnapshotInput) (*domain.CostingSnapshot, error) {
	if input.ProductID == 0 {
		return nil, fmt.Errorf("product_id is required")
	}
	if input.CostType != domain.CostTypePurchase && input.CostType != domain.CostTypeLanded && input.CostType != domain.CostTypeAverage {
		return nil, fmt.Errorf("invalid cost type")
	}
	if input.UnitCost <= 0 {
		return nil, fmt.Errorf("unit_cost must be greater than 0")
	}
	if input.CreatedBy == 0 {
		return nil, fmt.Errorf("created_by is required")
	}

	effectiveFrom := time.Now()
	if input.EffectiveFrom != nil {
		effectiveFrom = *input.EffectiveFrom
	}
	if input.EffectiveTo != nil && input.EffectiveTo.Before(effectiveFrom) {
		return nil, fmt.Errorf("effective_to must be greater than effective_from")
	}

	currency := input.Currency
	if currency == "" {
		currency = getDefaultBaseCurrency()
	}

	snapshot := &domain.CostingSnapshot{
		TraceID:       fmt.Sprintf("COST-%d-%d", time.Now().UnixNano(), input.CreatedBy),
		ProductID:     input.ProductID,
		CostType:      input.CostType,
		UnitCost:      input.UnitCost,
		Currency:      currency,
		EffectiveFrom: effectiveFrom,
		EffectiveTo:   input.EffectiveTo,
		Notes:         input.Notes,
		CreatedBy:     input.CreatedBy,
	}

	if input.EffectiveTo == nil {
		if err := uc.repo.ExpireCurrent(input.ProductID, input.CostType, effectiveFrom, nil); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Create(snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (uc *CostingSnapshotUsecase) Update(id uint64, input *UpdateCostingSnapshotInput) (*domain.CostingSnapshot, error) {
	snapshot, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if input.UnitCost != nil {
		if *input.UnitCost <= 0 {
			return nil, fmt.Errorf("unit_cost must be greater than 0")
		}
		snapshot.UnitCost = *input.UnitCost
	}
	if input.Currency != nil {
		if *input.Currency == "" {
			return nil, fmt.Errorf("currency is required")
		}
		snapshot.Currency = *input.Currency
	}
	if input.EffectiveFrom != nil {
		snapshot.EffectiveFrom = *input.EffectiveFrom
	}
	if input.EffectiveTo != nil {
		snapshot.EffectiveTo = *input.EffectiveTo
	}
	if input.Notes != nil {
		snapshot.Notes = *input.Notes
	}

	if snapshot.EffectiveTo != nil && snapshot.EffectiveTo.Before(snapshot.EffectiveFrom) {
		return nil, fmt.Errorf("effective_to must be greater than effective_from")
	}

	if snapshot.EffectiveTo == nil {
		excludeID := snapshot.ID
		if err := uc.repo.ExpireCurrent(snapshot.ProductID, snapshot.CostType, snapshot.EffectiveFrom, &excludeID); err != nil {
			return nil, err
		}
	}

	if err := uc.repo.Update(snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (uc *CostingSnapshotUsecase) Delete(id uint64) error {
	return uc.repo.Delete(id)
}

func (uc *CostingSnapshotUsecase) GetCurrent(productID uint64, costType domain.CostType) (*domain.CostingSnapshot, error) {
	snapshot, err := uc.repo.GetCurrent(productID, costType, time.Now())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return snapshot, nil
}

func (uc *CostingSnapshotUsecase) GetAllCurrent(productID uint64) ([]domain.CostingSnapshot, error) {
	return uc.repo.ListCurrentBySKU(productID, time.Now())
}
