package usecase

import (
	"time"

	"am-erp-go/internal/module/finance/domain"
)

type DailyProfitUsecase struct {
	profitLedgerRepo domain.ProfitLedgerRepository
	snapshotRepo     domain.DailyProfitSnapshotRepository
	orderCostRepo    domain.OrderCostDetailRepository
}

func NewDailyProfitUsecase(
	profitLedgerRepo domain.ProfitLedgerRepository,
	snapshotRepo domain.DailyProfitSnapshotRepository,
	orderCostRepo domain.OrderCostDetailRepository,
) *DailyProfitUsecase {
	return &DailyProfitUsecase{
		profitLedgerRepo: profitLedgerRepo,
		snapshotRepo:     snapshotRepo,
		orderCostRepo:    orderCostRepo,
	}
}

type RebuildDailyProfitInput struct {
	BizDate     time.Time
	Marketplace *string
	BuilderID   *uint64
}

type ProfitDashboardSummary struct {
	SalesIncomeAmount        float64 `json:"sales_income_amount"`
	COGSAmount               float64 `json:"cogs_amount"`
	GrossProfitAmount        float64 `json:"gross_profit_amount"`
	OrderExpenseAmount       float64 `json:"order_expense_amount"`
	OrderNetProfitAmount     float64 `json:"order_net_profit_amount"`
	PublicExpenseAmount      float64 `json:"public_expense_amount"`
	OperatingNetProfitAmount float64 `json:"operating_net_profit_amount"`
	OrderCount               uint64  `json:"order_count"`
	ShippedQty               uint64  `json:"shipped_qty"`
}

type ProfitDashboard struct {
	Items   []domain.DailyProfitSnapshot `json:"items"`
	Summary ProfitDashboardSummary       `json:"summary"`
}

type ProfitDashboardInput struct {
	DateFrom    time.Time
	DateTo      time.Time
	Marketplace string
}

func (uc *DailyProfitUsecase) Rebuild(input *RebuildDailyProfitInput) ([]domain.DailyProfitSnapshot, error) {
	bizDate := time.Now()
	marketplace := (*string)(nil)
	builderID := (*uint64)(nil)
	if input != nil {
		if !input.BizDate.IsZero() {
			bizDate = input.BizDate
		}
		marketplace = input.Marketplace
		builderID = input.BuilderID
	}
	bizDate = toBizDate(bizDate)

	aggs, err := uc.profitLedgerRepo.AggregateDaily(bizDate, marketplace)
	if err != nil {
		return nil, err
	}

	snapshots := make([]domain.DailyProfitSnapshot, 0, len(aggs))
	for _, agg := range aggs {
		market := agg.Marketplace
		if market == "" {
			market = "ALL"
		}
		var marketPtr *string
		if market != "ALL" {
			marketPtr = &market
		}
		shippedQty, err := uc.orderCostRepo.SumQtyByDateAndMarketplace(bizDate, marketPtr)
		if err != nil {
			return nil, err
		}

		gross := round6(agg.SalesIncome - agg.COGS)
		orderNet := round6(gross - agg.OrderExpense)
		operatingNet := round6(orderNet - agg.PublicExpense)
		snapshots = append(snapshots, domain.DailyProfitSnapshot{
			BizDate:                  bizDate,
			Marketplace:              market,
			BaseCurrency:             ensureBaseCurrency(agg.BaseCurrency),
			SalesIncomeAmount:        round6(agg.SalesIncome),
			COGSAmount:               round6(agg.COGS),
			GrossProfitAmount:        gross,
			OrderExpenseAmount:       round6(agg.OrderExpense),
			OrderNetProfitAmount:     orderNet,
			PublicExpenseAmount:      round6(agg.PublicExpense),
			OperatingNetProfitAmount: operatingNet,
			OrderCount:               agg.OrderCount,
			ShippedQty:               shippedQty,
			SnapshotStatus:           domain.DailyProfitSnapshotStatusRecalculated,
			SourceVersion:            fxVersionManual,
			BuiltAt:                  time.Now(),
			BuilderID:                builderID,
		})
	}

	if err := uc.snapshotRepo.DeleteByDate(bizDate, marketplace); err != nil {
		return nil, err
	}
	if err := uc.snapshotRepo.CreateBatch(snapshots); err != nil {
		return nil, err
	}
	return snapshots, nil
}

func (uc *DailyProfitUsecase) Dashboard(input *ProfitDashboardInput) (*ProfitDashboard, error) {
	dateFrom := time.Now().AddDate(0, 0, -6)
	dateTo := time.Now()
	marketplace := ""
	if input != nil {
		if !input.DateFrom.IsZero() {
			dateFrom = input.DateFrom
		}
		if !input.DateTo.IsZero() {
			dateTo = input.DateTo
		}
		marketplace = input.Marketplace
	}
	dateFrom = toBizDate(dateFrom)
	dateTo = toBizDate(dateTo)
	if dateTo.Before(dateFrom) {
		dateFrom, dateTo = dateTo, dateFrom
	}

	items, err := uc.snapshotRepo.List(&domain.DailyProfitSnapshotListParams{
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		Marketplace: marketplace,
	})
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].BaseCurrency = ensureBaseCurrency(items[i].BaseCurrency)
	}

	summary := ProfitDashboardSummary{}
	for _, item := range items {
		summary.SalesIncomeAmount += item.SalesIncomeAmount
		summary.COGSAmount += item.COGSAmount
		summary.GrossProfitAmount += item.GrossProfitAmount
		summary.OrderExpenseAmount += item.OrderExpenseAmount
		summary.OrderNetProfitAmount += item.OrderNetProfitAmount
		summary.PublicExpenseAmount += item.PublicExpenseAmount
		summary.OperatingNetProfitAmount += item.OperatingNetProfitAmount
		summary.OrderCount += item.OrderCount
		summary.ShippedQty += item.ShippedQty
	}

	summary.SalesIncomeAmount = round6(summary.SalesIncomeAmount)
	summary.COGSAmount = round6(summary.COGSAmount)
	summary.GrossProfitAmount = round6(summary.GrossProfitAmount)
	summary.OrderExpenseAmount = round6(summary.OrderExpenseAmount)
	summary.OrderNetProfitAmount = round6(summary.OrderNetProfitAmount)
	summary.PublicExpenseAmount = round6(summary.PublicExpenseAmount)
	summary.OperatingNetProfitAmount = round6(summary.OperatingNetProfitAmount)

	return &ProfitDashboard{
		Items:   items,
		Summary: summary,
	}, nil
}
