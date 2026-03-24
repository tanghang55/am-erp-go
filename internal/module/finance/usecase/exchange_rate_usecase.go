package usecase

import (
	"fmt"
	"strings"
	"time"

	"am-erp-go/internal/module/finance/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

const (
	fxSourceIdentity  = "IDENTITY"
	fxSourceDerived   = "DERIVED"
	fxVersionIdentity = "same_currency"
	fxVersionManual   = "v1"
	fxVersionInverse  = "inverse"
	fxVersionCrossUSD = "cross_usd"
	fxPivotCurrency   = "USD"
)

type FXRateSnapshot struct {
	Rate        float64
	Source      string
	Version     string
	EffectiveAt time.Time
}

type ExchangeRateUsecase struct {
	repo        domain.ExchangeRateRepository
	auditLogger ExchangeRateAuditLogger
}

type ExchangeRateAuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

func NewExchangeRateUsecase(repo domain.ExchangeRateRepository) *ExchangeRateUsecase {
	return &ExchangeRateUsecase{repo: repo}
}

func (uc *ExchangeRateUsecase) BindAuditLogger(logger ExchangeRateAuditLogger) {
	uc.auditLogger = logger
}

type CreateExchangeRateInput struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	EffectiveAt  *time.Time
	Remark       *string
	CreatedBy    uint64
}

type UpdateExchangeRateStatusInput struct {
	Status     domain.ExchangeRateStatus
	OperatorID uint64
}

func (uc *ExchangeRateUsecase) List(params *domain.ExchangeRateListParams) ([]domain.ExchangeRate, int64, error) {
	if params == nil {
		params = &domain.ExchangeRateListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.repo.List(params)
}

func (uc *ExchangeRateUsecase) Create(c *gin.Context, input *CreateExchangeRateInput) (*domain.ExchangeRate, error) {
	if input == nil {
		return nil, fmt.Errorf("input is required")
	}
	from := normalizeCurrency(input.FromCurrency)
	to := normalizeCurrency(input.ToCurrency)
	if from == "" || to == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if from == to {
		return nil, fmt.Errorf("same currency does not require exchange rate")
	}
	if input.Rate <= 0 {
		return nil, fmt.Errorf("rate must be greater than 0")
	}
	if input.CreatedBy == 0 {
		return nil, fmt.Errorf("created_by is required")
	}
	effectiveAt := time.Now()
	if input.EffectiveAt != nil && !input.EffectiveAt.IsZero() {
		effectiveAt = *input.EffectiveAt
	}
	rate := &domain.ExchangeRate{
		FromCurrency:  from,
		ToCurrency:    to,
		Rate:          roundExchangeRate(input.Rate),
		SourceType:    domain.ExchangeRateSourceManual,
		SourceVersion: fxVersionManual,
		EffectiveAt:   effectiveAt,
		Status:        domain.ExchangeRateStatusActive,
		Remark:        input.Remark,
		CreatedBy:     input.CreatedBy,
		UpdatedBy:     input.CreatedBy,
	}
	if err := uc.repo.Create(rate); err != nil {
		return nil, err
	}
	uc.recordAudit(c, "CREATE_EXCHANGE_RATE", fmt.Sprintf("%d", rate.ID), nil, rate)
	return rate, nil
}

func (uc *ExchangeRateUsecase) UpdateStatus(c *gin.Context, id uint64, input *UpdateExchangeRateStatusInput) error {
	if id == 0 {
		return fmt.Errorf("id is required")
	}
	if input == nil {
		return fmt.Errorf("input is required")
	}
	if input.Status != domain.ExchangeRateStatusActive && input.Status != domain.ExchangeRateStatusInactive {
		return fmt.Errorf("invalid status")
	}
	if input.OperatorID == 0 {
		return fmt.Errorf("operator_id is required")
	}
	before, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}
	if err := uc.repo.UpdateStatus(id, input.Status, input.OperatorID); err != nil {
		return err
	}
	after, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}
	uc.recordAudit(c, "UPDATE_EXCHANGE_RATE_STATUS", fmt.Sprintf("%d", id), before, after)
	return nil
}

func (uc *ExchangeRateUsecase) Resolve(baseCurrency, originalCurrency string, occurredAt time.Time) (*FXRateSnapshot, error) {
	base := normalizeCurrency(baseCurrency)
	original := normalizeCurrency(originalCurrency)
	if base == "" || original == "" {
		return nil, fmt.Errorf("currency is required")
	}
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	if base == original {
		return &FXRateSnapshot{
			Rate:        1,
			Source:      fxSourceIdentity,
			Version:     fxVersionIdentity,
			EffectiveAt: occurredAt,
		}, nil
	}
	if uc.repo == nil {
		return nil, fmt.Errorf("exchange rate repository is required")
	}
	if snapshot, err := uc.resolveDirect(original, base, occurredAt); err == nil {
		return snapshot, nil
	}
	if snapshot, err := uc.resolveInverse(original, base, occurredAt); err == nil {
		return snapshot, nil
	}
	if snapshot, err := uc.resolveViaUSDPivot(original, base, occurredAt); err == nil {
		return snapshot, nil
	}
	return nil, fmt.Errorf("exchange rate not found for %s -> %s at %s", original, base, occurredAt.Format(time.RFC3339))
}

func (uc *ExchangeRateUsecase) resolveDirect(fromCurrency, toCurrency string, occurredAt time.Time) (*FXRateSnapshot, error) {
	row, err := uc.repo.FindEffectiveRate(fromCurrency, toCurrency, occurredAt)
	if err != nil {
		return nil, err
	}
	return snapshotFromRow(row), nil
}

func (uc *ExchangeRateUsecase) resolveInverse(fromCurrency, toCurrency string, occurredAt time.Time) (*FXRateSnapshot, error) {
	row, err := uc.repo.FindEffectiveRate(toCurrency, fromCurrency, occurredAt)
	if err != nil {
		return nil, err
	}
	if row == nil || row.Rate <= 0 {
		return nil, fmt.Errorf("inverse rate not found")
	}
	return &FXRateSnapshot{
		Rate:        roundExchangeRate(1 / row.Rate),
		Source:      fxSourceDerived,
		Version:     fxVersionInverse,
		EffectiveAt: row.EffectiveAt,
	}, nil
}

func (uc *ExchangeRateUsecase) resolveViaUSDPivot(fromCurrency, toCurrency string, occurredAt time.Time) (*FXRateSnapshot, error) {
	if fromCurrency == fxPivotCurrency || toCurrency == fxPivotCurrency {
		return nil, fmt.Errorf("usd pivot not applicable")
	}
	fromPivot, err := uc.resolveAgainstPivot(fromCurrency, occurredAt)
	if err != nil {
		return nil, err
	}
	toPivot, err := uc.resolveAgainstPivot(toCurrency, occurredAt)
	if err != nil {
		return nil, err
	}
	return &FXRateSnapshot{
		Rate:        roundExchangeRate(fromPivot.Rate / toPivot.Rate),
		Source:      fxSourceDerived,
		Version:     fxVersionCrossUSD,
		EffectiveAt: maxTime(fromPivot.EffectiveAt, toPivot.EffectiveAt),
	}, nil
}

func (uc *ExchangeRateUsecase) recordAudit(c *gin.Context, action, entityID string, before, after any) {
	if uc.auditLogger == nil || c == nil {
		return
	}
	_ = uc.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Finance",
		Action:     action,
		EntityType: "ExchangeRate",
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
}

func (uc *ExchangeRateUsecase) resolveAgainstPivot(currency string, occurredAt time.Time) (*FXRateSnapshot, error) {
	if currency == fxPivotCurrency {
		return &FXRateSnapshot{
			Rate:        1,
			Source:      fxSourceIdentity,
			Version:     fxVersionIdentity,
			EffectiveAt: occurredAt,
		}, nil
	}
	if row, err := uc.repo.FindEffectiveRate(currency, fxPivotCurrency, occurredAt); err == nil && row != nil && row.Rate > 0 {
		return snapshotFromRow(row), nil
	}
	row, err := uc.repo.FindEffectiveRate(fxPivotCurrency, currency, occurredAt)
	if err != nil {
		return nil, err
	}
	if row == nil || row.Rate <= 0 {
		return nil, fmt.Errorf("pivot rate not found")
	}
	return &FXRateSnapshot{
		Rate:        roundExchangeRate(1 / row.Rate),
		Source:      fxSourceDerived,
		Version:     fxVersionInverse,
		EffectiveAt: row.EffectiveAt,
	}, nil
}

func snapshotFromRow(row *domain.ExchangeRate) *FXRateSnapshot {
	source := strings.TrimSpace(string(row.SourceType))
	if source == "" {
		source = string(domain.ExchangeRateSourceManual)
	}
	version := strings.TrimSpace(row.SourceVersion)
	if version == "" {
		version = fxVersionManual
	}
	return &FXRateSnapshot{
		Rate:        roundExchangeRate(row.Rate),
		Source:      source,
		Version:     version,
		EffectiveAt: row.EffectiveAt,
	}
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
