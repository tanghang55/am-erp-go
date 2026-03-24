package usecase

import (
	"am-erp-go/internal/module/finance/domain"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrCashLedgerImmutable       = errors.New("cash ledger is immutable, use reverse")
	ErrCashLedgerAlreadyReversed = errors.New("cash ledger already reversed")
)

type CashLedgerUsecase struct {
	repo                 domain.CashLedgerRepository
	profitLedgerRecorder CashLedgerProfitLedgerRecorder
}

func NewCashLedgerUsecase(repo domain.CashLedgerRepository) *CashLedgerUsecase {
	return &CashLedgerUsecase{repo: repo}
}

type CashLedgerProfitLedgerRecorder interface {
	RecordCashLedger(entry *domain.CashLedger) error
}

func (uc *CashLedgerUsecase) BindProfitLedgerRecorder(recorder CashLedgerProfitLedgerRecorder) {
	uc.profitLedgerRecorder = recorder
}

type CreateCashLedgerInput struct {
	LedgerType    domain.LedgerType
	Category      string
	Amount        float64
	Currency      string
	Marketplace   *string
	OccurredNode  *string
	ReferenceType *string
	ReferenceID   *uint64
	Description   *string
	OccurredAt    *time.Time
	CreatedBy     uint64
}

type UpdateCashLedgerInput struct {
	LedgerType    *domain.LedgerType
	Category      *string
	Amount        *float64
	Currency      *string
	ReferenceType **string
	ReferenceID   **uint64
	Description   **string
	OccurredAt    *time.Time
}

func (uc *CashLedgerUsecase) List(params *domain.CashLedgerListParams) ([]domain.CashLedger, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.repo.List(params)
}

func (uc *CashLedgerUsecase) Get(id uint64) (*domain.CashLedger, error) {
	return uc.repo.GetByID(id)
}

func (uc *CashLedgerUsecase) Create(input *CreateCashLedgerInput) (*domain.CashLedger, error) {
	if input.LedgerType != domain.LedgerTypeIncome && input.LedgerType != domain.LedgerTypeExpense {
		return nil, fmt.Errorf("invalid ledger type")
	}
	if input.Category == "" {
		return nil, fmt.Errorf("category is required")
	}
	if input.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}
	if input.CreatedBy == 0 {
		return nil, fmt.Errorf("created_by is required")
	}

	currency := input.Currency
	if strings.TrimSpace(currency) == "" {
		currency = getDefaultBaseCurrency()
	}
	originalCurrency := normalizeCurrency(currency)
	baseCurrency := getDefaultBaseCurrency()

	occurredAt := time.Now()
	if input.OccurredAt != nil {
		occurredAt = *input.OccurredAt
	}
	fxSnapshot, err := resolveFXRate(baseCurrency, originalCurrency, occurredAt)
	if err != nil {
		return nil, err
	}
	originalAmount := round6(input.Amount)
	baseAmount := round6(input.Amount * fxSnapshot.Rate)

	entry := &domain.CashLedger{
		TraceID:          fmt.Sprintf("CASH-%d-%d", time.Now().UnixNano(), input.CreatedBy),
		LedgerType:       input.LedgerType,
		Status:           domain.CashLedgerStatusNormal,
		Category:         input.Category,
		Amount:           input.Amount,
		Currency:         originalCurrency,
		OriginalCurrency: originalCurrency,
		OriginalAmount:   originalAmount,
		BaseCurrency:     baseCurrency,
		FxRate:           fxSnapshot.Rate,
		BaseAmount:       baseAmount,
		FxSource:         fxSnapshot.Source,
		FxVersion:        fxSnapshot.Version,
		FxTime:           fxSnapshot.EffectiveAt,
		Marketplace:      input.Marketplace,
		OccurredNode:     input.OccurredNode,
		ReferenceType:    input.ReferenceType,
		ReferenceID:      input.ReferenceID,
		Description:      input.Description,
		OccurredAt:       occurredAt,
		CreatedBy:        input.CreatedBy,
	}

	if err := uc.repo.Create(entry); err != nil {
		return nil, err
	}
	if uc.profitLedgerRecorder != nil && entry.LedgerType == domain.LedgerTypeExpense {
		if err := uc.profitLedgerRecorder.RecordCashLedger(entry); err != nil {
			return nil, err
		}
	}
	return entry, nil
}

func (uc *CashLedgerUsecase) Update(id uint64, input *UpdateCashLedgerInput) (*domain.CashLedger, error) {
	_ = id
	_ = input
	return nil, ErrCashLedgerImmutable
}

func (uc *CashLedgerUsecase) Delete(id uint64) error {
	_ = id
	return ErrCashLedgerImmutable
}

func (uc *CashLedgerUsecase) Reverse(id uint64, operatorID uint64, reason string) (*domain.CashLedger, error) {
	if operatorID == 0 {
		return nil, fmt.Errorf("operator_id is required")
	}

	entry, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if entry.Status == domain.CashLedgerStatusReversed {
		return nil, ErrCashLedgerAlreadyReversed
	}

	reversedType := domain.LedgerTypeExpense
	if entry.LedgerType == domain.LedgerTypeExpense {
		reversedType = domain.LedgerTypeIncome
	}

	now := time.Now()
	description := "reverse cash ledger"
	if reason != "" {
		description = "reverse cash ledger: " + reason
	}
	originalCurrency := normalizeCurrency(entry.OriginalCurrency)
	if originalCurrency == "" {
		originalCurrency = normalizeCurrency(entry.Currency)
	}
	baseCurrency := normalizeCurrency(entry.BaseCurrency)
	if baseCurrency == "" {
		baseCurrency = getDefaultBaseCurrency()
	}
	originalAmount := entry.OriginalAmount
	if originalAmount <= 0 {
		originalAmount = entry.Amount
	}
	fxSnapshot, err := resolveFXRate(baseCurrency, originalCurrency, now)
	if err != nil {
		return nil, err
	}
	baseAmount := round6(originalAmount * fxSnapshot.Rate)
	occurredNode := "REVERSED"

	reversed := &domain.CashLedger{
		TraceID:          fmt.Sprintf("CASH-REV-%d-%d", now.UnixNano(), operatorID),
		LedgerType:       reversedType,
		Status:           domain.CashLedgerStatusNormal,
		ReversalOfID:     &entry.ID,
		Category:         entry.Category,
		Amount:           round6(originalAmount),
		Currency:         originalCurrency,
		OriginalCurrency: originalCurrency,
		OriginalAmount:   round6(originalAmount),
		BaseCurrency:     baseCurrency,
		FxRate:           fxSnapshot.Rate,
		BaseAmount:       baseAmount,
		FxSource:         fxSnapshot.Source,
		FxVersion:        fxSnapshot.Version,
		FxTime:           fxSnapshot.EffectiveAt,
		Marketplace:      entry.Marketplace,
		OccurredNode:     &occurredNode,
		ReferenceType:    entry.ReferenceType,
		ReferenceID:      entry.ReferenceID,
		Description:      &description,
		OccurredAt:       now,
		CreatedBy:        operatorID,
	}

	if err = uc.repo.Create(reversed); err != nil {
		return nil, err
	}
	if err := uc.repo.MarkReversed(entry.ID, now); err != nil {
		return nil, err
	}

	return reversed, nil
}

func (uc *CashLedgerUsecase) GetSummary(params *domain.CashLedgerListParams) (*domain.CashLedgerSummary, error) {
	return uc.repo.GetSummary(params)
}

func (uc *CashLedgerUsecase) GetSummaryByCategory(params *domain.CashLedgerListParams) ([]domain.CategorySummaryItem, error) {
	return uc.repo.GetSummaryByCategory(params)
}
