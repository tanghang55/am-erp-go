package usecase

import (
	"fmt"
	"strings"
	"time"

	"am-erp-go/internal/module/finance/domain"
	salesUsecase "am-erp-go/internal/module/sales/usecase"
)

type ProfitLedgerWriter struct {
	repo domain.ProfitLedgerRepository
}

func NewProfitLedgerWriter(repo domain.ProfitLedgerRepository) *ProfitLedgerWriter {
	return &ProfitLedgerWriter{repo: repo}
}

func (w *ProfitLedgerWriter) RecordSalesShipProfit(params *salesUsecase.SalesShipProfitRecordParams) error {
	if w.repo == nil || params == nil || params.SalesOrderID == 0 || params.SalesOrderItemID == 0 {
		return nil
	}
	occurredAt := params.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	bizDate := toBizDate(occurredAt)

	var marketplace *string
	if params.Marketplace != "" {
		m := params.Marketplace
		marketplace = &m
	}

	refType := "SALES_ORDER"
	refID := params.SalesOrderID
	refNo := params.OrderNo
	salesOrderID := params.SalesOrderID
	salesOrderItemID := params.SalesOrderItemID

	incomeCurrency := normalizeCurrency(params.IncomeCurrency)
	if incomeCurrency == "" {
		incomeCurrency = getDefaultBaseCurrency()
	}
	incomeOriginal := round6(params.IncomeAmount)
	baseCurrency := getDefaultBaseCurrency()
	incomeFxSnapshot, err := resolveFXRate(baseCurrency, incomeCurrency, occurredAt)
	if err != nil {
		return err
	}
	incomeBase := round6(incomeOriginal * incomeFxSnapshot.Rate)

	entries := make([]domain.ProfitLedger, 0, 2)
	if incomeOriginal > 0 {
		entries = append(entries, domain.ProfitLedger{
			TraceID:          fmt.Sprintf("PROFIT-SALES-%d-INCOME", time.Now().UnixNano()),
			LedgerType:       domain.ProfitLedgerTypeIncome,
			Status:           domain.ProfitLedgerStatusNormal,
			BizDate:          bizDate,
			Marketplace:      marketplace,
			SalesOrderID:     &salesOrderID,
			SalesOrderItemID: &salesOrderItemID,
			ReferenceType:    &refType,
			ReferenceID:      &refID,
			ReferenceNumber:  &refNo,
			Category:         "SALES_REVENUE",
			OriginalCurrency: incomeCurrency,
			OriginalAmount:   incomeOriginal,
			BaseCurrency:     baseCurrency,
			FxRate:           incomeFxSnapshot.Rate,
			BaseAmount:       incomeBase,
			FxSource:         incomeFxSnapshot.Source,
			FxVersion:        incomeFxSnapshot.Version,
			FxTime:           incomeFxSnapshot.EffectiveAt,
			OccurredAt:       occurredAt,
			OperatorID:       params.OperatorID,
		})
	}

	cogsOriginal := round6(params.COGSAmount)
	if cogsOriginal > 0 {
		cogsFxSnapshot, err := resolveFXRate(baseCurrency, baseCurrency, occurredAt)
		if err != nil {
			return err
		}
		entries = append(entries, domain.ProfitLedger{
			TraceID:          fmt.Sprintf("PROFIT-SALES-%d-COGS", time.Now().UnixNano()),
			LedgerType:       domain.ProfitLedgerTypeCOGS,
			Status:           domain.ProfitLedgerStatusNormal,
			BizDate:          bizDate,
			Marketplace:      marketplace,
			SalesOrderID:     &salesOrderID,
			SalesOrderItemID: &salesOrderItemID,
			ReferenceType:    &refType,
			ReferenceID:      &refID,
			ReferenceNumber:  &refNo,
			Category:         "SALES_COGS",
			OriginalCurrency: baseCurrency,
			OriginalAmount:   cogsOriginal,
			BaseCurrency:     baseCurrency,
			FxRate:           cogsFxSnapshot.Rate,
			BaseAmount:       cogsOriginal,
			FxSource:         cogsFxSnapshot.Source,
			FxVersion:        cogsFxSnapshot.Version,
			FxTime:           cogsFxSnapshot.EffectiveAt,
			OccurredAt:       occurredAt,
			OperatorID:       params.OperatorID,
		})
	}

	return w.repo.CreateBatch(entries)
}

func (w *ProfitLedgerWriter) RecordSalesReturnProfit(params *salesUsecase.SalesReturnProfitRecordParams) error {
	if w.repo == nil || params == nil || params.SalesOrderID == 0 || params.SalesOrderItemID == 0 {
		return nil
	}
	occurredAt := params.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	bizDate := toBizDate(occurredAt)

	var marketplace *string
	if params.Marketplace != "" {
		m := params.Marketplace
		marketplace = &m
	}

	refType := "SALES_ORDER"
	refID := params.SalesOrderID
	refNo := params.OrderNo
	salesOrderID := params.SalesOrderID
	salesOrderItemID := params.SalesOrderItemID

	incomeCurrency := normalizeCurrency(params.IncomeCurrency)
	if incomeCurrency == "" {
		incomeCurrency = getDefaultBaseCurrency()
	}
	incomeOriginal := round6(params.IncomeAmount)
	baseCurrency := getDefaultBaseCurrency()
	incomeFxSnapshot, err := resolveFXRate(baseCurrency, incomeCurrency, occurredAt)
	if err != nil {
		return err
	}
	incomeBase := round6(incomeOriginal * incomeFxSnapshot.Rate)

	entries := make([]domain.ProfitLedger, 0, 2)
	if incomeOriginal > 0 {
		reversalOfID := uint64(1)
		entries = append(entries, domain.ProfitLedger{
			TraceID:          fmt.Sprintf("PROFIT-SALES-RETURN-%d-INCOME", time.Now().UnixNano()),
			LedgerType:       domain.ProfitLedgerTypeIncome,
			Status:           domain.ProfitLedgerStatusNormal,
			ReversalOfID:     &reversalOfID,
			BizDate:          bizDate,
			Marketplace:      marketplace,
			SalesOrderID:     &salesOrderID,
			SalesOrderItemID: &salesOrderItemID,
			ReferenceType:    &refType,
			ReferenceID:      &refID,
			ReferenceNumber:  &refNo,
			Category:         "SALES_RETURN_REVENUE",
			OriginalCurrency: incomeCurrency,
			OriginalAmount:   incomeOriginal,
			BaseCurrency:     baseCurrency,
			FxRate:           incomeFxSnapshot.Rate,
			BaseAmount:       incomeBase,
			FxSource:         incomeFxSnapshot.Source,
			FxVersion:        incomeFxSnapshot.Version,
			FxTime:           incomeFxSnapshot.EffectiveAt,
			OccurredAt:       occurredAt,
			OperatorID:       params.OperatorID,
		})
	}

	cogsOriginal := round6(params.COGSAmount)
	if cogsOriginal > 0 {
		cogsFxSnapshot, err := resolveFXRate(baseCurrency, baseCurrency, occurredAt)
		if err != nil {
			return err
		}
		reversalOfID := uint64(1)
		entries = append(entries, domain.ProfitLedger{
			TraceID:          fmt.Sprintf("PROFIT-SALES-RETURN-%d-COGS", time.Now().UnixNano()),
			LedgerType:       domain.ProfitLedgerTypeCOGS,
			Status:           domain.ProfitLedgerStatusNormal,
			ReversalOfID:     &reversalOfID,
			BizDate:          bizDate,
			Marketplace:      marketplace,
			SalesOrderID:     &salesOrderID,
			SalesOrderItemID: &salesOrderItemID,
			ReferenceType:    &refType,
			ReferenceID:      &refID,
			ReferenceNumber:  &refNo,
			Category:         "SALES_RETURN_COGS",
			OriginalCurrency: baseCurrency,
			OriginalAmount:   cogsOriginal,
			BaseCurrency:     baseCurrency,
			FxRate:           cogsFxSnapshot.Rate,
			BaseAmount:       cogsOriginal,
			FxSource:         cogsFxSnapshot.Source,
			FxVersion:        cogsFxSnapshot.Version,
			FxTime:           cogsFxSnapshot.EffectiveAt,
			OccurredAt:       occurredAt,
			OperatorID:       params.OperatorID,
		})
	}

	return w.repo.CreateBatch(entries)
}

func (w *ProfitLedgerWriter) RecordCashLedger(entry *domain.CashLedger) error {
	if w.repo == nil || entry == nil || entry.LedgerType != domain.LedgerTypeExpense {
		return nil
	}

	occurredAt := entry.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	bizDate := toBizDate(occurredAt)

	ledgerType := domain.ProfitLedgerTypePublicExp
	if entry.ReferenceType != nil && strings.EqualFold(strings.TrimSpace(*entry.ReferenceType), "SALES_ORDER") {
		ledgerType = domain.ProfitLedgerTypeOrderExpense
	}

	originalCurrency := normalizeCurrency(entry.OriginalCurrency)
	if originalCurrency == "" {
		originalCurrency = normalizeCurrency(entry.Currency)
	}
	if originalCurrency == "" {
		originalCurrency = getDefaultBaseCurrency()
	}

	baseCurrency := normalizeCurrency(entry.BaseCurrency)
	if baseCurrency == "" {
		baseCurrency = getDefaultBaseCurrency()
	}

	fxRate := entry.FxRate
	fxSource := entry.FxSource
	fxVersion := entry.FxVersion
	fxTime := entry.FxTime
	if fxRate <= 0 {
		fxSnapshot, err := resolveFXRate(baseCurrency, originalCurrency, entry.OccurredAt)
		if err != nil {
			return err
		}
		fxRate = fxSnapshot.Rate
		fxSource = fxSnapshot.Source
		fxVersion = fxSnapshot.Version
		fxTime = fxSnapshot.EffectiveAt
	}

	originalAmount := entry.OriginalAmount
	if originalAmount <= 0 {
		originalAmount = entry.Amount
	}
	originalAmount = round6(originalAmount)

	baseAmount := entry.BaseAmount
	if baseAmount <= 0 {
		baseAmount = round6(originalAmount * fxRate)
	}

	category := strings.TrimSpace(entry.Category)
	if category == "" {
		category = "UNKNOWN"
	}

	refType := entry.ReferenceType
	refID := entry.ReferenceID
	refNo := (*string)(nil)
	remark := (*string)(nil)
	if entry.Description != nil && strings.TrimSpace(*entry.Description) != "" {
		desc := strings.TrimSpace(*entry.Description)
		remark = &desc
	}
	marketplace := entry.Marketplace
	if fxTime.IsZero() {
		fxTime = occurredAt
	}

	var operatorID *uint64
	if entry.CreatedBy > 0 {
		id := entry.CreatedBy
		operatorID = &id
	}

	profitEntry := &domain.ProfitLedger{
		TraceID:          fmt.Sprintf("PROFIT-CASH-%d-%d", time.Now().UnixNano(), entry.ID),
		LedgerType:       ledgerType,
		Status:           domain.ProfitLedgerStatusNormal,
		BizDate:          bizDate,
		Marketplace:      marketplace,
		ReferenceType:    refType,
		ReferenceID:      refID,
		ReferenceNumber:  refNo,
		Category:         category,
		OriginalCurrency: originalCurrency,
		OriginalAmount:   originalAmount,
		BaseCurrency:     baseCurrency,
		FxRate:           fxRate,
		BaseAmount:       baseAmount,
		FxSource:         fxSource,
		FxVersion:        fxVersion,
		FxTime:           fxTime,
		OccurredAt:       occurredAt,
		OperatorID:       operatorID,
		Remark:           remark,
	}
	return w.repo.Create(profitEntry)
}

func toBizDate(t time.Time) time.Time {
	local := t.In(time.Local)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, local.Location())
}
