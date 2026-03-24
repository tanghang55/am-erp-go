package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"am-erp-go/internal/module/packaging/domain"
)

type stubPackagingProcurementRepo struct {
	plans         []domain.PackagingProcurementPlan
	snapshots     map[uint64]domain.PackagingItemSnapshot
	createdOrder  *domain.PackagingPurchaseOrder
	updatedOrder  *domain.PackagingPurchaseOrder
	order         *domain.PackagingPurchaseOrder
	err           error
	markedPlanIDs []uint64
	markedOrderID uint64
}

func (s *stubPackagingProcurementRepo) CleanupPlansBefore(date time.Time) error { return nil }
func (s *stubPackagingProcurementRepo) LoadOrderedProductDemands(planDate time.Time) ([]domain.PackagingProductDemand, error) {
	return nil, nil
}
func (s *stubPackagingProcurementRepo) LoadProductPackagingMappings(productIDs []uint64) ([]domain.ProductPackagingMapping, error) {
	return nil, nil
}
func (s *stubPackagingProcurementRepo) LoadPackagingItemSnapshots(itemIDs []uint64) (map[uint64]domain.PackagingItemSnapshot, error) {
	return s.snapshots, nil
}
func (s *stubPackagingProcurementRepo) SyncDailyPlans(planDate time.Time, inputs []domain.PackagingPlanInput) ([]domain.PackagingProcurementPlan, int, error) {
	return nil, 0, nil
}
func (s *stubPackagingProcurementRepo) ListPlans(params *domain.PackagingProcurementPlanListParams) ([]domain.PackagingProcurementPlan, int64, error) {
	return nil, 0, nil
}
func (s *stubPackagingProcurementRepo) ListRuns(params *domain.PackagingProcurementRunListParams) ([]domain.PackagingProcurementRun, int64, error) {
	return nil, 0, nil
}
func (s *stubPackagingProcurementRepo) ListConvertiblePlans(params *domain.PackagingPlanConvertParams) ([]domain.PackagingProcurementPlan, error) {
	return s.plans, s.err
}
func (s *stubPackagingProcurementRepo) MarkPlansConverted(planIDs []uint64, purchaseOrderID uint64) error {
	s.markedPlanIDs = append([]uint64{}, planIDs...)
	s.markedOrderID = purchaseOrderID
	return s.err
}
func (s *stubPackagingProcurementRepo) CreateRun(run *domain.PackagingProcurementRun) error {
	return nil
}
func (s *stubPackagingProcurementRepo) UpdateRun(run *domain.PackagingProcurementRun) error {
	return nil
}
func (s *stubPackagingProcurementRepo) CreatePurchaseOrder(order *domain.PackagingPurchaseOrder) error {
	s.createdOrder = order
	order.ID = 1
	s.order = clonePackagingPurchaseOrder(order)
	return s.err
}
func (s *stubPackagingProcurementRepo) UpdatePurchaseOrder(order *domain.PackagingPurchaseOrder) error {
	s.updatedOrder = order
	s.order = clonePackagingPurchaseOrder(order)
	return s.err
}
func (s *stubPackagingProcurementRepo) ListPurchaseOrders(params *domain.PackagingPurchaseOrderListParams) ([]domain.PackagingPurchaseOrder, int64, error) {
	return nil, 0, nil
}
func (s *stubPackagingProcurementRepo) GetPurchaseOrder(id uint64) (*domain.PackagingPurchaseOrder, error) {
	return clonePackagingPurchaseOrder(s.order), s.err
}

type stubPackagingItemRepo struct {
	items     map[uint64]*domain.PackagingItem
	updated   map[uint64]int64
	getErr    error
	updateErr error
}

func (s *stubPackagingItemRepo) List(params *domain.PackagingItemListParams) ([]domain.PackagingItem, int64, error) {
	return nil, 0, nil
}
func (s *stubPackagingItemRepo) GetByID(id uint64) (*domain.PackagingItem, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	if item, ok := s.items[id]; ok {
		cp := *item
		return &cp, nil
	}
	return nil, nil
}
func (s *stubPackagingItemRepo) Create(item *domain.PackagingItem) error           { return nil }
func (s *stubPackagingItemRepo) Update(item *domain.PackagingItem) error           { return nil }
func (s *stubPackagingItemRepo) Delete(id uint64) error                            { return nil }
func (s *stubPackagingItemRepo) CountReferences(id uint64) (int64, error)          { return 0, nil }
func (s *stubPackagingItemRepo) GetLowStockItems() ([]domain.PackagingItem, error) { return nil, nil }
func (s *stubPackagingItemRepo) UpdateQuantity(id uint64, quantity int64) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	if s.updated == nil {
		s.updated = map[uint64]int64{}
	}
	s.updated[id] += quantity
	if item, ok := s.items[id]; ok {
		item.QuantityOnHand += uint64(quantity)
	}
	return nil
}

type stubPackagingLedgerRepo struct {
	created []*domain.PackagingLedger
	err     error
}

func (s *stubPackagingLedgerRepo) Create(ledger *domain.PackagingLedger) error {
	if s.err != nil {
		return s.err
	}
	s.created = append(s.created, ledger)
	return nil
}
func (s *stubPackagingLedgerRepo) List(params *domain.PackagingLedgerListParams) ([]domain.PackagingLedger, int64, error) {
	return nil, 0, nil
}
func (s *stubPackagingLedgerRepo) GetByID(id uint64) (*domain.PackagingLedger, error) {
	return nil, nil
}
func (s *stubPackagingLedgerRepo) GetUsageSummary(dateFrom, dateTo *time.Time) ([]domain.UsageSummaryItem, error) {
	return nil, nil
}

type stubPackagingBaseCurrencyProvider struct{}

func (stubPackagingBaseCurrencyProvider) GetDefaultBaseCurrency() string {
	return "EUR"
}

type packagingConvertTxManagerStub struct {
	deps      PackagingPlanConvertTransactionalDeps
	called    bool
	committed bool
}

func (s *packagingConvertTxManagerStub) Run(ctx context.Context, fn func(PackagingPlanConvertTransactionalDeps) error) error {
	s.called = true
	baseRepo, _ := s.deps.Repo.(*stubPackagingProcurementRepo)
	txRepo := *baseRepo
	txRepo.createdOrder = nil
	txRepo.markedPlanIDs = nil
	txRepo.markedOrderID = 0
	if err := fn(PackagingPlanConvertTransactionalDeps{Repo: &txRepo}); err != nil {
		return err
	}
	s.committed = true
	*baseRepo = txRepo
	return nil
}

type packagingReceiveTxManagerStub struct {
	deps      PackagingPurchaseReceiveTransactionalDeps
	called    bool
	committed bool
}

func (s *packagingReceiveTxManagerStub) Run(ctx context.Context, fn func(PackagingPurchaseReceiveTransactionalDeps) error) error {
	s.called = true
	baseRepo, _ := s.deps.Repo.(*stubPackagingProcurementRepo)
	baseItemRepo, _ := s.deps.ItemRepo.(*stubPackagingItemRepo)
	baseLedgerRepo, _ := s.deps.LedgerRepo.(*stubPackagingLedgerRepo)

	txRepo := *baseRepo
	txRepo.order = clonePackagingPurchaseOrder(baseRepo.order)
	txRepo.updatedOrder = nil

	txItems := &stubPackagingItemRepo{
		items:     clonePackagingItems(baseItemRepo.items),
		updated:   map[uint64]int64{},
		getErr:    baseItemRepo.getErr,
		updateErr: baseItemRepo.updateErr,
	}
	txLedger := &stubPackagingLedgerRepo{created: nil, err: baseLedgerRepo.err}

	if err := fn(PackagingPurchaseReceiveTransactionalDeps{
		Repo:       &txRepo,
		ItemRepo:   txItems,
		LedgerRepo: txLedger,
	}); err != nil {
		return err
	}
	s.committed = true
	*baseRepo = txRepo
	baseItemRepo.items = txItems.items
	baseItemRepo.updated = txItems.updated
	baseLedgerRepo.created = txLedger.created
	return nil
}

func clonePackagingPurchaseOrder(order *domain.PackagingPurchaseOrder) *domain.PackagingPurchaseOrder {
	if order == nil {
		return nil
	}
	cp := *order
	if len(order.Items) > 0 {
		cp.Items = make([]domain.PackagingPurchaseOrderItem, len(order.Items))
		copy(cp.Items, order.Items)
	}
	if order.CreatedBy != nil {
		v := *order.CreatedBy
		cp.CreatedBy = &v
	}
	if order.UpdatedBy != nil {
		v := *order.UpdatedBy
		cp.UpdatedBy = &v
	}
	return &cp
}

func clonePackagingItems(items map[uint64]*domain.PackagingItem) map[uint64]*domain.PackagingItem {
	if items == nil {
		return nil
	}
	result := make(map[uint64]*domain.PackagingItem, len(items))
	for id, item := range items {
		cp := *item
		result[id] = &cp
	}
	return result
}

func TestConvertPackagingPlansFallsBackToConfiguredCurrency(t *testing.T) {
	repo := &stubPackagingProcurementRepo{
		plans: []domain.PackagingProcurementPlan{
			{ID: 1, PackagingItemID: 10, SuggestedQty: 12},
		},
		snapshots: map[uint64]domain.PackagingItemSnapshot{
			10: {ID: 10, UnitCost: 2.5, Currency: ""},
		},
	}
	uc := NewPackagingProcurementUsecase(repo, &stubPackagingItemRepo{}, &stubPackagingLedgerRepo{})
	uc.BindDefaultsProvider(stubPackagingBaseCurrencyProvider{})

	date := time.Date(2026, 3, 9, 0, 0, 0, 0, time.Local)
	order, err := uc.ConvertPlans(nil, &domain.PackagingPlanConvertParams{Date: &date})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.Currency != "EUR" {
		t.Fatalf("expected order currency EUR, got %s", order.Currency)
	}
	if len(order.Items) != 1 || order.Items[0].Currency != "EUR" {
		t.Fatalf("expected item currency EUR, got %+v", order.Items)
	}
	if repo.createdOrder == nil || repo.createdOrder.Currency != "EUR" {
		t.Fatalf("expected persisted order currency EUR, got %+v", repo.createdOrder)
	}
}

func TestConvertPackagingPlansRollsBackWhenMarkConvertedFails(t *testing.T) {
	repo := &stubPackagingProcurementRepo{
		plans: []domain.PackagingProcurementPlan{
			{ID: 1, PackagingItemID: 10, SuggestedQty: 12},
		},
		snapshots: map[uint64]domain.PackagingItemSnapshot{
			10: {ID: 10, UnitCost: 2.5, Currency: ""},
		},
	}
	txManager := &packagingConvertTxManagerStub{
		deps: PackagingPlanConvertTransactionalDeps{Repo: repo},
	}
	uc := NewPackagingProcurementUsecase(repo, &stubPackagingItemRepo{}, &stubPackagingLedgerRepo{})
	uc.BindDefaultsProvider(stubPackagingBaseCurrencyProvider{})
	uc.BindConvertTransactionManager(txManager)

	date := time.Date(2026, 3, 9, 0, 0, 0, 0, time.Local)
	repo.err = errors.New("mark converted failed")
	_, err := uc.ConvertPlans(nil, &domain.PackagingPlanConvertParams{Date: &date})
	if err == nil {
		t.Fatalf("expected convert error")
	}
	if !txManager.called {
		t.Fatalf("expected convert tx manager called")
	}
	if txManager.committed {
		t.Fatalf("expected convert tx not committed")
	}
	if repo.createdOrder != nil {
		t.Fatalf("expected no persisted created order, got %+v", repo.createdOrder)
	}
	if len(repo.markedPlanIDs) != 0 {
		t.Fatalf("expected no persisted marked plans")
	}
}

func TestReceivePurchaseOrderRollsBackWhenLedgerCreateFails(t *testing.T) {
	order := &domain.PackagingPurchaseOrder{
		ID:       8,
		PoNumber: "PKPO-8",
		Status:   domain.PackagingPurchaseOrderOrdered,
		Currency: "USD",
		Items: []domain.PackagingPurchaseOrderItem{
			{ID: 81, PackagingItemID: 10, QtyOrdered: 5, QtyReceived: 1, UnitCost: 2.5, Currency: "USD"},
		},
	}
	repo := &stubPackagingProcurementRepo{order: order}
	itemRepo := &stubPackagingItemRepo{
		items: map[uint64]*domain.PackagingItem{
			10: {ID: 10, QuantityOnHand: 3},
		},
	}
	ledgerRepo := &stubPackagingLedgerRepo{err: errors.New("ledger failed")}
	txManager := &packagingReceiveTxManagerStub{
		deps: PackagingPurchaseReceiveTransactionalDeps{
			Repo:       repo,
			ItemRepo:   itemRepo,
			LedgerRepo: ledgerRepo,
		},
	}
	uc := NewPackagingProcurementUsecase(repo, itemRepo, ledgerRepo)
	uc.BindReceiveTransactionManager(txManager)

	_, err := uc.ReceivePurchaseOrder(nil, order.ID, &domain.PackagingPurchaseOrderReceiveParams{
		ReceivedQties: map[uint64]uint64{81: 2},
	})
	if err == nil {
		t.Fatalf("expected receive error")
	}
	if !txManager.called {
		t.Fatalf("expected receive tx manager called")
	}
	if txManager.committed {
		t.Fatalf("expected receive tx not committed")
	}
	if repo.updatedOrder != nil {
		t.Fatalf("expected no persisted order update")
	}
	if len(ledgerRepo.created) != 0 {
		t.Fatalf("expected no persisted ledger")
	}
	if len(itemRepo.updated) != 0 {
		t.Fatalf("expected no persisted item quantity update")
	}
	if repo.order.Items[0].QtyReceived != 1 {
		t.Fatalf("expected order qty_received unchanged, got %d", repo.order.Items[0].QtyReceived)
	}
}
