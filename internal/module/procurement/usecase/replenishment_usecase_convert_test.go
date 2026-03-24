package usecase

import (
	"testing"
	"time"

	"am-erp-go/internal/module/procurement/domain"
	productdomain "am-erp-go/internal/module/product/domain"
)

type stubReplenishmentConvertRepo struct {
	plans          []domain.ReplenishmentPlan
	links          []domain.ReplenishmentPlanPurchaseOrderLink
	convertedPlans map[uint64]uint64
}

func (s *stubReplenishmentConvertRepo) GetConfig() (*domain.ReplenishmentConfig, error) {
	return &domain.ReplenishmentConfig{}, nil
}

func (s *stubReplenishmentConvertRepo) SaveConfig(config *domain.ReplenishmentConfig) error {
	return nil
}

func (s *stubReplenishmentConvertRepo) ListStrategies(params *domain.ReplenishmentStrategyListParams) ([]domain.ReplenishmentStrategy, int64, error) {
	return []domain.ReplenishmentStrategy{}, 0, nil
}

func (s *stubReplenishmentConvertRepo) ListActiveStrategies() ([]domain.ReplenishmentStrategy, error) {
	return []domain.ReplenishmentStrategy{}, nil
}

func (s *stubReplenishmentConvertRepo) UpsertStrategy(strategy *domain.ReplenishmentStrategy) (*domain.ReplenishmentStrategy, error) {
	return strategy, nil
}

func (s *stubReplenishmentConvertRepo) ListPlans(params *domain.ReplenishmentPlanListParams) ([]domain.ReplenishmentPlan, int64, error) {
	return s.plans, int64(len(s.plans)), nil
}

func (s *stubReplenishmentConvertRepo) CreatePlans(plans []domain.ReplenishmentPlan) (int, error) {
	return len(plans), nil
}

func (s *stubReplenishmentConvertRepo) DeletePendingPlanByID(planID uint64) (bool, error) {
	return false, nil
}

func (s *stubReplenishmentConvertRepo) DeletePlansBefore(date time.Time) error {
	return nil
}

func (s *stubReplenishmentConvertRepo) DeletePlansByPurchaseOrderID(purchaseOrderID uint64) error {
	return nil
}

func (s *stubReplenishmentConvertRepo) ListConvertiblePlans(params *domain.ReplenishmentPlanConvertParams) ([]domain.ReplenishmentPlan, error) {
	return s.plans, nil
}

func (s *stubReplenishmentConvertRepo) MarkPlansConverted(planIDs []uint64, purchaseOrderID uint64) error {
	if s.convertedPlans == nil {
		s.convertedPlans = map[uint64]uint64{}
	}
	for _, planID := range planIDs {
		s.convertedPlans[planID] = purchaseOrderID
	}
	return nil
}

func (s *stubReplenishmentConvertRepo) LinkPlansToPurchaseOrders(links []domain.ReplenishmentPlanPurchaseOrderLink) error {
	s.links = append(s.links, links...)
	return nil
}

func (s *stubReplenishmentConvertRepo) ListPolicies(params *domain.ReplenishmentPolicyListParams) ([]domain.ReplenishmentPolicy, int64, error) {
	return []domain.ReplenishmentPolicy{}, 0, nil
}

func (s *stubReplenishmentConvertRepo) UpsertPolicy(policy *domain.ReplenishmentPolicy) (*domain.ReplenishmentPolicy, error) {
	return policy, nil
}

func (s *stubReplenishmentConvertRepo) ListRuns(params *domain.ReplenishmentRunListParams) ([]domain.ReplenishmentRun, int64, error) {
	return []domain.ReplenishmentRun{}, 0, nil
}

func (s *stubReplenishmentConvertRepo) GetRunByID(runID uint64) (*domain.ReplenishmentRun, error) {
	return nil, nil
}

func (s *stubReplenishmentConvertRepo) ListRunItems(runID uint64) ([]domain.ReplenishmentItem, error) {
	return []domain.ReplenishmentItem{}, nil
}

func (s *stubReplenishmentConvertRepo) CreateRunWithItems(run *domain.ReplenishmentRun, items []domain.ReplenishmentItem) error {
	return nil
}

func (s *stubReplenishmentConvertRepo) CreateRunItems(runID uint64, items []domain.ReplenishmentItem) error {
	return nil
}

func (s *stubReplenishmentConvertRepo) UpdateRun(run *domain.ReplenishmentRun) error {
	return nil
}

func (s *stubReplenishmentConvertRepo) ListPoliciesByProductIDs(productIDs []uint64) (map[uint64]domain.ReplenishmentPolicy, error) {
	return map[uint64]domain.ReplenishmentPolicy{}, nil
}

func (s *stubReplenishmentConvertRepo) LoadDemandByWindowDays(windowDays uint32) ([]domain.ReplenishmentDemandRow, error) {
	return []domain.ReplenishmentDemandRow{}, nil
}

func (s *stubReplenishmentConvertRepo) LoadBalanceRows() ([]domain.ReplenishmentBalanceRow, error) {
	return []domain.ReplenishmentBalanceRow{}, nil
}

func (s *stubReplenishmentConvertRepo) LoadProductProfiles(productIDs []uint64) (map[uint64]domain.ReplenishmentProductProfile, error) {
	return map[uint64]domain.ReplenishmentProductProfile{}, nil
}

func (s *stubReplenishmentConvertRepo) LoadPackagingRequirementsByProduct(productIDs []uint64) (map[uint64][]domain.ReplenishmentPackagingRequirement, error) {
	return map[uint64][]domain.ReplenishmentPackagingRequirement{}, nil
}

func (s *stubReplenishmentConvertRepo) LoadPackagingItems(itemIDs []uint64) (map[uint64]domain.ReplenishmentPackagingItem, error) {
	return map[uint64]domain.ReplenishmentPackagingItem{}, nil
}

func (s *stubReplenishmentConvertRepo) ListConvertibleItems(runID uint64, itemIDs []uint64) ([]domain.ReplenishmentItem, error) {
	return []domain.ReplenishmentItem{}, nil
}

func (s *stubReplenishmentConvertRepo) MarkItemsConverted(itemIDs []uint64, purchaseOrderID uint64) error {
	return nil
}

func TestConvertPlansToPurchaseOrdersSplitsComboPlanAndLinksAllChildOrders(t *testing.T) {
	planDate := time.Date(2026, 3, 15, 0, 0, 0, 0, time.Local)
	planSupplierID := uint64(1)
	comboID := uint64(40)
	mainID := uint64(500)
	childAID := uint64(501)
	childBID := uint64(502)
	supplierAID := uint64(8)
	supplierBID := uint64(9)

	repo := &stubReplenishmentConvertRepo{
		plans: []domain.ReplenishmentPlan{
			{
				ID:           1,
				PlanDate:     planDate,
				ProductID:    mainID,
				WarehouseID:  3,
				SupplierID:   &planSupplierID,
				SuggestedQty: 2,
				Status:       domain.ReplenishmentPlanPending,
			},
		},
	}
	poRepo := &stubPurchaseOrderRepo{}
	productRepo := &stubProductLookup{
		items: map[uint64]productdomain.Product{
			mainID:   {ID: mainID, ComboID: &comboID, IsComboMain: 1},
			childAID: {ID: childAID, SupplierID: &supplierAID, UnitCost: floatPtr(10)},
			childBID: {ID: childBID, SupplierID: &supplierBID, UnitCost: floatPtr(20)},
		},
	}
	comboRepo := &stubComboProvider{
		items: map[uint64][]productdomain.ProductComboItem{
			comboID: {
				{ComboID: comboID, MainProductID: mainID, ProductID: mainID, QtyRatio: 1},
				{ComboID: comboID, MainProductID: mainID, ProductID: childAID, QtyRatio: 1},
				{ComboID: comboID, MainProductID: mainID, ProductID: childBID, QtyRatio: 2},
			},
		},
	}

	poUsecase := NewPurchaseOrderUsecase(poRepo, productRepo, comboRepo, nil, nil)
	uc := NewReplenishmentUsecase(repo, poUsecase)

	orders, err := uc.ConvertPlansToPurchaseOrders(nil, &domain.ReplenishmentPlanConvertParams{
		PlanDate: &planDate,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 2 {
		t.Fatalf("expected 2 split purchase orders, got %d", len(orders))
	}
	if orders[0].BatchNo == "" || orders[0].BatchNo != orders[1].BatchNo {
		t.Fatalf("expected shared batch number, got %+v", orders)
	}
	if len(repo.links) != 2 {
		t.Fatalf("expected 2 plan-order links, got %+v", repo.links)
	}
	for _, link := range repo.links {
		if link.PlanID != 1 {
			t.Fatalf("expected plan id 1 in all links, got %+v", repo.links)
		}
		if link.PurchaseOrderID == 0 {
			t.Fatalf("expected linked purchase order id, got %+v", repo.links)
		}
	}
	if repo.convertedPlans[1] == 0 {
		t.Fatalf("expected plan marked converted with first purchase order id, got %+v", repo.convertedPlans)
	}
}
