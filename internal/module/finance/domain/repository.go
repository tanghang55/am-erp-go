package domain

import "time"

type CashLedgerRepository interface {
	List(params *CashLedgerListParams) ([]CashLedger, int64, error)
	GetByID(id uint64) (*CashLedger, error)
	Create(entry *CashLedger) error
	Update(entry *CashLedger) error
	Delete(id uint64) error
	MarkReversed(id uint64, reversedAt time.Time) error
	GetSummary(params *CashLedgerListParams) (*CashLedgerSummary, error)
	GetSummaryByCategory(params *CashLedgerListParams) ([]CategorySummaryItem, error)
}

type CostingSnapshotRepository interface {
	List(params *CostingSnapshotListParams) ([]CostingSnapshot, int64, error)
	GetByID(id uint64) (*CostingSnapshot, error)
	Create(snapshot *CostingSnapshot) error
	Update(snapshot *CostingSnapshot) error
	Delete(id uint64) error
	ExpireCurrent(productID uint64, costType CostType, effectiveTo time.Time, excludeID *uint64) error
	GetCurrent(productID uint64, costType CostType, now time.Time) (*CostingSnapshot, error)
	ListCurrentBySKU(productID uint64, now time.Time) ([]CostingSnapshot, error)
}

type CostEventRepository interface {
	Create(event *CostEvent) error
	GetLatestPackingMaterialPerUnit(productID uint64, occurredAt time.Time) (*float64, error)
}

type OrderCostDetailRepository interface {
	CreateBatch(details []OrderCostDetail) error
	SumQtyByDateAndMarketplace(bizDate time.Time, marketplace *string) (uint64, error)
	ListReturnableBySalesOrderItemID(salesOrderItemID uint64) ([]ReturnableOrderCostDetail, error)
}

type ProfitLedgerDailyAgg struct {
	Marketplace   string
	BaseCurrency  string
	SalesIncome   float64
	COGS          float64
	OrderExpense  float64
	PublicExpense float64
	OrderCount    uint64
	ShippedQty    uint64
}

type ProfitLedgerRepository interface {
	Create(entry *ProfitLedger) error
	CreateBatch(entries []ProfitLedger) error
	AggregateDaily(bizDate time.Time, marketplace *string) ([]ProfitLedgerDailyAgg, error)
}

type DailyProfitSnapshotRepository interface {
	DeleteByDate(bizDate time.Time, marketplace *string) error
	CreateBatch(snapshots []DailyProfitSnapshot) error
	List(params *DailyProfitSnapshotListParams) ([]DailyProfitSnapshot, error)
}

type ProfitQueryRepository interface {
	ListOrderProfits(params *OrderProfitListParams) ([]OrderProfitSummary, int64, error)
	GetOrderProfitDetail(salesOrderID uint64) (*OrderProfitDetail, error)
}

type ProductCostRepository interface {
	ListLedger(params *ProductCostLedgerListParams) ([]ProductCostLedgerItem, int64, error)
	GetSummary(params *ProductCostLedgerListParams) (*ProductCostSummary, error)
}

type ExchangeRateRepository interface {
	List(params *ExchangeRateListParams) ([]ExchangeRate, int64, error)
	Create(rate *ExchangeRate) error
	GetByID(id uint64) (*ExchangeRate, error)
	UpdateStatus(id uint64, status ExchangeRateStatus, operatorID uint64) error
	FindEffectiveRate(fromCurrency, toCurrency string, occurredAt time.Time) (*ExchangeRate, error)
}
