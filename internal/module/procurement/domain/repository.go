package domain

import "time"

// PurchaseOrderRepository 采购单仓储接口
type PurchaseOrderRepository interface {
	List(params *PurchaseOrderListParams) ([]PurchaseOrder, int64, error)
	GetByID(id uint64) (*PurchaseOrder, error)
	Create(order *PurchaseOrder) error
	Update(order *PurchaseOrder) error
	UpdateProgress(order *PurchaseOrder) error
	Delete(id uint64) error
}

type ReplenishmentDemandRow struct {
	ProductID   uint64
	WarehouseID uint64
	ShippedQty  uint64
}

type ReplenishmentBalanceRow struct {
	ProductID           uint64
	WarehouseID         uint64
	AvailableQuantity   uint64
	ReservedQuantity    uint64
	PurchasingInTransit uint64
	PendingInspection   uint64
}

type ReplenishmentProductProfile struct {
	ProductID     uint64
	SupplierID    *uint64
	Marketplace   *string
	ProductStatus string
	QuoteMOQ      *uint64
	QuoteLeadDays *uint64
	QuotePrice    *float64
}

type ReplenishmentPackagingRequirement struct {
	ProductID       uint64
	PackagingItemID uint64
	QuantityPerUnit float64
}

type ReplenishmentPackagingItem struct {
	PackagingItemID uint64
	ItemCode        string
	ItemName        string
	QuantityOnHand  uint64
}

type ReplenishmentRepository interface {
	GetConfig() (*ReplenishmentConfig, error)
	SaveConfig(config *ReplenishmentConfig) error
	ListStrategies(params *ReplenishmentStrategyListParams) ([]ReplenishmentStrategy, int64, error)
	ListActiveStrategies() ([]ReplenishmentStrategy, error)
	UpsertStrategy(strategy *ReplenishmentStrategy) (*ReplenishmentStrategy, error)
	ListPlans(params *ReplenishmentPlanListParams) ([]ReplenishmentPlan, int64, error)
	CreatePlans(plans []ReplenishmentPlan) (int, error)
	DeletePendingPlanByID(planID uint64) (bool, error)
	DeletePlansBefore(date time.Time) error
	DeletePlansByPurchaseOrderID(purchaseOrderID uint64) error
	ListConvertiblePlans(params *ReplenishmentPlanConvertParams) ([]ReplenishmentPlan, error)
	MarkPlansConverted(planIDs []uint64, purchaseOrderID uint64) error
	LinkPlansToPurchaseOrders(links []ReplenishmentPlanPurchaseOrderLink) error

	// Legacy methods kept for compatibility with existing code paths.
	ListPolicies(params *ReplenishmentPolicyListParams) ([]ReplenishmentPolicy, int64, error)
	UpsertPolicy(policy *ReplenishmentPolicy) (*ReplenishmentPolicy, error)
	ListRuns(params *ReplenishmentRunListParams) ([]ReplenishmentRun, int64, error)
	GetRunByID(runID uint64) (*ReplenishmentRun, error)
	ListRunItems(runID uint64) ([]ReplenishmentItem, error)
	CreateRunWithItems(run *ReplenishmentRun, items []ReplenishmentItem) error
	CreateRunItems(runID uint64, items []ReplenishmentItem) error
	UpdateRun(run *ReplenishmentRun) error
	ListPoliciesByProductIDs(productIDs []uint64) (map[uint64]ReplenishmentPolicy, error)
	LoadDemandByWindowDays(windowDays uint32) ([]ReplenishmentDemandRow, error)
	LoadBalanceRows() ([]ReplenishmentBalanceRow, error)
	LoadProductProfiles(productIDs []uint64) (map[uint64]ReplenishmentProductProfile, error)
	LoadPackagingRequirementsByProduct(productIDs []uint64) (map[uint64][]ReplenishmentPackagingRequirement, error)
	LoadPackagingItems(itemIDs []uint64) (map[uint64]ReplenishmentPackagingItem, error)
	ListConvertibleItems(runID uint64, itemIDs []uint64) ([]ReplenishmentItem, error)
	MarkItemsConverted(itemIDs []uint64, purchaseOrderID uint64) error
}
