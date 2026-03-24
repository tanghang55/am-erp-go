package usecase

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/procurement/domain"
	productdomain "am-erp-go/internal/module/product/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type stubPurchaseOrderRepo struct {
	created       *domain.PurchaseOrder
	createdOrders []*domain.PurchaseOrder
	updated       *domain.PurchaseOrder
	getID         uint64
	getItem       *domain.PurchaseOrder
	err           error
	createErrors  []error
	createCalls   int
	createNumbers []string
}

func (s *stubPurchaseOrderRepo) List(params *domain.PurchaseOrderListParams) ([]domain.PurchaseOrder, int64, error) {
	return nil, 0, s.err
}

func (s *stubPurchaseOrderRepo) GetByID(id uint64) (*domain.PurchaseOrder, error) {
	s.getID = id
	if s.getItem == nil {
		for _, created := range s.createdOrders {
			if created != nil && created.ID == id {
				return clonePurchaseOrder(created), s.err
			}
		}
		if s.updated != nil && s.updated.ID == id {
			return clonePurchaseOrder(s.updated), s.err
		}
		if s.created != nil && s.created.ID == id {
			return clonePurchaseOrder(s.created), s.err
		}
	}
	for _, created := range s.createdOrders {
		if created != nil && created.ID == id {
			return clonePurchaseOrder(created), s.err
		}
	}
	return clonePurchaseOrder(s.getItem), s.err
}

func (s *stubPurchaseOrderRepo) Create(order *domain.PurchaseOrder) error {
	s.createCalls++
	if order != nil {
		s.createNumbers = append(s.createNumbers, order.PoNumber)
	}
	if s.createCalls <= len(s.createErrors) {
		return s.createErrors[s.createCalls-1]
	}
	s.created = order
	if order != nil && order.ID == 0 {
		order.ID = uint64(s.createCalls)
	}
	s.createdOrders = append(s.createdOrders, clonePurchaseOrder(order))
	s.getItem = clonePurchaseOrder(order)
	return s.err
}

func (s *stubPurchaseOrderRepo) Update(order *domain.PurchaseOrder) error {
	if s.getItem != nil && order != nil {
		order.Supplier = s.getItem.Supplier
		if order.Items == nil {
			order.Items = clonePurchaseOrder(s.getItem).Items
		}
		for i := range order.Items {
			for j := range s.getItem.Items {
				if s.getItem.Items[j].ProductID != order.Items[i].ProductID {
					continue
				}
				order.Items[i].Product = s.getItem.Items[j].Product
				break
			}
		}
	}
	s.updated = order
	s.getItem = clonePurchaseOrder(order)
	return s.err
}

func (s *stubPurchaseOrderRepo) UpdateProgress(order *domain.PurchaseOrder) error {
	if s.getItem != nil && order != nil {
		current := clonePurchaseOrder(s.getItem)
		current.PoNumber = order.PoNumber
		current.BatchNo = order.BatchNo
		current.SupplierID = order.SupplierID
		current.WarehouseID = order.WarehouseID
		current.Marketplace = order.Marketplace
		current.Status = order.Status
		current.Currency = order.Currency
		current.TotalAmount = order.TotalAmount
		current.OrderedAt = order.OrderedAt
		current.OrderedBy = order.OrderedBy
		current.ShippedAt = order.ShippedAt
		current.ShippedBy = order.ShippedBy
		current.ReceivedAt = order.ReceivedAt
		current.ReceivedBy = order.ReceivedBy
		current.InspectedAt = order.InspectedAt
		current.InspectedBy = order.InspectedBy
		current.ClosedAt = order.ClosedAt
		current.CompletedBy = order.CompletedBy
		current.IsForceCompleted = order.IsForceCompleted
		current.ForceCompletedAt = order.ForceCompletedAt
		current.ForceCompletedBy = order.ForceCompletedBy
		current.ForceCompleteReason = order.ForceCompleteReason
		current.Remark = order.Remark
		current.CreatedBy = order.CreatedBy
		current.UpdatedBy = order.UpdatedBy
		if order.Items != nil {
			current.Items = clonePurchaseOrder(order).Items
		}
		s.updated = current
		s.getItem = clonePurchaseOrder(current)
		return s.err
	}
	s.updated = order
	s.getItem = clonePurchaseOrder(order)
	return s.err
}

func (s *stubPurchaseOrderRepo) Delete(id uint64) error {
	return s.err
}

type stubProductLookup struct {
	items map[uint64]productdomain.Product
}

func (s *stubProductLookup) ListByIDs(ids []uint64) ([]productdomain.Product, error) {
	result := make([]productdomain.Product, 0, len(ids))
	for _, id := range ids {
		if item, ok := s.items[id]; ok {
			result = append(result, item)
		}
	}
	return result, nil
}

type stubComboProvider struct {
	items map[uint64][]productdomain.ProductComboItem
}

func (s *stubComboProvider) GetItemsByComboID(comboID uint64) ([]productdomain.ProductComboItem, error) {
	return s.items[comboID], nil
}

type stubInventoryService struct {
	created []*inventoryDomain.CreateMovementParams
	err     error
}

func (s *stubInventoryService) CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error) {
	if params != nil {
		s.created = append(s.created, params)
	}
	if s.err != nil {
		return nil, s.err
	}
	return &inventoryDomain.InventoryMovement{
		ID:        1,
		ProductID: params.ProductID,
		Quantity:  params.Quantity,
	}, nil
}

type stubPurchaseOrderCostEventRecorder struct {
	events []*PurchaseOrderCostEventParams
	err    error
}

type stubReplenishmentPlanCleaner struct {
	deletedPurchaseOrderIDs []uint64
	err                     error
}

type stubPurchaseOrderAuditLogger struct {
	payloads []systemUsecase.AuditLogPayload
}

func (s *stubReplenishmentPlanCleaner) DeletePlansByPurchaseOrderID(purchaseOrderID uint64) error {
	if s.err != nil {
		return s.err
	}
	s.deletedPurchaseOrderIDs = append(s.deletedPurchaseOrderIDs, purchaseOrderID)
	return nil
}

func (s *stubPurchaseOrderAuditLogger) RecordFromContext(_ *gin.Context, payload systemUsecase.AuditLogPayload) error {
	s.payloads = append(s.payloads, payload)
	return nil
}

type stubBaseCurrencyProvider struct{}

func (stubBaseCurrencyProvider) GetDefaultBaseCurrency() string {
	return "EUR"
}

func (s *stubPurchaseOrderCostEventRecorder) RecordPurchaseOrderEvent(params *PurchaseOrderCostEventParams) error {
	if params != nil {
		s.events = append(s.events, params)
	}
	return s.err
}

type purchaseOrderShipTxManagerStub struct {
	deps      PurchaseOrderShipTransactionalDeps
	called    bool
	committed bool
}

func (s *purchaseOrderShipTxManagerStub) Run(ctx context.Context, fn func(PurchaseOrderShipTransactionalDeps) error) error {
	s.called = true
	baseRepo, _ := s.deps.Repo.(*stubPurchaseOrderRepo)
	repo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(baseRepo.getItem)}
	inventory := &stubInventoryService{}
	deps := PurchaseOrderShipTransactionalDeps{
		Repo:              repo,
		InventoryService:  inventory,
		CostEventRecorder: s.deps.CostEventRecorder,
	}
	if err := fn(deps); err != nil {
		return err
	}
	s.committed = true
	if baseRepo != nil {
		baseRepo.updated = repo.updated
		baseRepo.getItem = repo.getItem
	}
	if baseInventory, ok := s.deps.InventoryService.(*stubInventoryService); ok {
		baseInventory.created = inventory.created
	}
	return nil
}

type purchaseOrderReceiveTxManagerStub struct {
	deps      PurchaseOrderReceiveTransactionalDeps
	called    bool
	committed bool
}

type purchaseOrderInspectTxManagerStub struct {
	deps      PurchaseOrderInspectTransactionalDeps
	called    bool
	committed bool
}

type purchaseOrderSubmitTxManagerStub struct {
	deps      PurchaseOrderSubmitTransactionalDeps
	called    bool
	committed bool
}

func (s *purchaseOrderSubmitTxManagerStub) Run(ctx context.Context, fn func(PurchaseOrderSubmitTransactionalDeps) error) error {
	s.called = true
	baseRepo, _ := s.deps.Repo.(*stubPurchaseOrderRepo)
	repo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(baseRepo.getItem)}
	deps := PurchaseOrderSubmitTransactionalDeps{
		Repo:              repo,
		PlanCleaner:       s.deps.PlanCleaner,
		CostEventRecorder: s.deps.CostEventRecorder,
	}
	if err := fn(deps); err != nil {
		return err
	}
	s.committed = true
	if baseRepo != nil {
		baseRepo.updated = repo.updated
		baseRepo.getItem = repo.getItem
	}
	return nil
}

func (s *purchaseOrderReceiveTxManagerStub) Run(ctx context.Context, fn func(PurchaseOrderReceiveTransactionalDeps) error) error {
	s.called = true
	baseRepo, _ := s.deps.Repo.(*stubPurchaseOrderRepo)
	repo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(baseRepo.getItem)}
	inventory := &stubInventoryService{}
	deps := PurchaseOrderReceiveTransactionalDeps{
		Repo:              repo,
		InventoryService:  inventory,
		CostEventRecorder: s.deps.CostEventRecorder,
	}
	if err := fn(deps); err != nil {
		return err
	}
	s.committed = true
	if baseRepo != nil {
		baseRepo.updated = repo.updated
		baseRepo.getItem = repo.getItem
	}
	if baseInventory, ok := s.deps.InventoryService.(*stubInventoryService); ok {
		baseInventory.created = inventory.created
	}
	return nil
}

func (s *purchaseOrderInspectTxManagerStub) Run(ctx context.Context, fn func(PurchaseOrderInspectTransactionalDeps) error) error {
	s.called = true
	baseRepo, _ := s.deps.Repo.(*stubPurchaseOrderRepo)
	repo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(baseRepo.getItem)}
	inventory := &stubInventoryService{}
	deps := PurchaseOrderInspectTransactionalDeps{
		Repo:             repo,
		InventoryService: inventory,
	}
	if err := fn(deps); err != nil {
		return err
	}
	s.committed = true
	if baseRepo != nil {
		baseRepo.updated = repo.updated
		baseRepo.getItem = repo.getItem
	}
	if baseInventory, ok := s.deps.InventoryService.(*stubInventoryService); ok {
		baseInventory.created = inventory.created
	}
	return nil
}

func floatPtr(v float64) *float64 {
	return &v
}

func clonePurchaseOrder(order *domain.PurchaseOrder) *domain.PurchaseOrder {
	if order == nil {
		return nil
	}
	cp := *order
	if len(order.Items) > 0 {
		cp.Items = make([]domain.PurchaseOrderItem, len(order.Items))
		copy(cp.Items, order.Items)
	}
	if order.SupplierID != nil {
		v := *order.SupplierID
		cp.SupplierID = &v
	}
	if order.WarehouseID != nil {
		v := *order.WarehouseID
		cp.WarehouseID = &v
	}
	if order.OrderedBy != nil {
		v := *order.OrderedBy
		cp.OrderedBy = &v
	}
	if order.ShippedBy != nil {
		v := *order.ShippedBy
		cp.ShippedBy = &v
	}
	if order.ReceivedBy != nil {
		v := *order.ReceivedBy
		cp.ReceivedBy = &v
	}
	if order.InspectedBy != nil {
		v := *order.InspectedBy
		cp.InspectedBy = &v
	}
	if order.CompletedBy != nil {
		v := *order.CompletedBy
		cp.CompletedBy = &v
	}
	if order.ForceCompletedBy != nil {
		v := *order.ForceCompletedBy
		cp.ForceCompletedBy = &v
	}
	if order.OrderedAt != nil {
		v := *order.OrderedAt
		cp.OrderedAt = &v
	}
	if order.ShippedAt != nil {
		v := *order.ShippedAt
		cp.ShippedAt = &v
	}
	if order.ReceivedAt != nil {
		v := *order.ReceivedAt
		cp.ReceivedAt = &v
	}
	if order.InspectedAt != nil {
		v := *order.InspectedAt
		cp.InspectedAt = &v
	}
	if order.ClosedAt != nil {
		v := *order.ClosedAt
		cp.ClosedAt = &v
	}
	if order.ForceCompletedAt != nil {
		v := *order.ForceCompletedAt
		cp.ForceCompletedAt = &v
	}
	return &cp
}

func TestCreatePurchaseOrderExpandsComboItems(t *testing.T) {
	comboID := uint64(10)
	mainID := uint64(100)

	productRepo := &stubProductLookup{
		items: map[uint64]productdomain.Product{
			mainID: {ID: mainID, ComboID: &comboID, IsComboMain: 1},
			200:    {ID: 200, UnitCost: floatPtr(2.5)},
			201:    {ID: 201, UnitCost: floatPtr(1.0)},
		},
	}

	comboRepo := &stubComboProvider{
		items: map[uint64][]productdomain.ProductComboItem{
			comboID: {
				{ComboID: comboID, MainProductID: mainID, ProductID: mainID, QtyRatio: 1},
				{ComboID: comboID, MainProductID: mainID, ProductID: 200, QtyRatio: 2},
				{ComboID: comboID, MainProductID: mainID, ProductID: 201, QtyRatio: 3},
			},
		},
	}

	poRepo := &stubPurchaseOrderRepo{}

	uc := NewPurchaseOrderUsecase(poRepo, productRepo, comboRepo, nil, nil)

	order := &domain.PurchaseOrder{
		SupplierID:  uint64Ptr(9),
		Currency:    "USD",
		Marketplace: "US",
		Items: []domain.PurchaseOrderItem{
			{ProductID: mainID, QtyOrdered: 5, UnitCost: 10, Currency: "USD"},
			{ProductID: 200, QtyOrdered: 1, UnitCost: 2.5, Currency: "USD"},
		},
	}

	created, err := uc.Create(nil, order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created == nil || created.PoNumber == "" {
		t.Fatalf("expected po number generated")
	}

	if poRepo.created == nil {
		t.Fatalf("expected order to be created")
	}

	if len(poRepo.created.Items) != 2 {
		t.Fatalf("expected 2 child items after combo expansion, got %d", len(poRepo.created.Items))
	}

	var item200 *domain.PurchaseOrderItem
	var item201 *domain.PurchaseOrderItem
	for i := range poRepo.created.Items {
		item := &poRepo.created.Items[i]
		if item.ProductID == 200 {
			item200 = item
		}
		if item.ProductID == 201 {
			item201 = item
		}
	}

	if item200 == nil || item200.QtyOrdered != 11 {
		t.Fatalf("expected sku 200 qty 11, got %+v", item200)
	}
	if item200.UnitCost != 2.5 {
		t.Fatalf("expected sku 200 unit cost 2.5")
	}
	if item200.Subtotal != 27.5 {
		t.Fatalf("expected sku 200 subtotal 27.5, got %.2f", item200.Subtotal)
	}

	if item201 == nil || item201.QtyOrdered != 15 {
		t.Fatalf("expected sku 201 qty 15, got %+v", item201)
	}
	if item201.UnitCost != 1.0 {
		t.Fatalf("expected sku 201 unit cost 1.0")
	}
	if item201.Subtotal != 15.0 {
		t.Fatalf("expected sku 201 subtotal 15.0, got %.2f", item201.Subtotal)
	}
	if poRepo.created.TotalAmount != 42.5 {
		t.Fatalf("expected total amount 42.5, got %.2f", poRepo.created.TotalAmount)
	}
}

func TestCreatePurchaseOrderRejectsMissingProductID(t *testing.T) {
	poRepo := &stubPurchaseOrderRepo{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, nil)

	order := &domain.PurchaseOrder{
		SupplierID:  uint64Ptr(9),
		Currency:    "USD",
		Marketplace: "US",
		Items: []domain.PurchaseOrderItem{
			{ProductID: 0, QtyOrdered: 1, UnitCost: 10, Currency: "USD"},
		},
	}

	_, err := uc.Create(nil, order)
	if !errors.Is(err, ErrPurchaseOrderMissingProduct) {
		t.Fatalf("expected missing product error, got %v", err)
	}
}

func TestCreatePurchaseOrderRejectsZeroQtyOrdered(t *testing.T) {
	poRepo := &stubPurchaseOrderRepo{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, nil)

	order := &domain.PurchaseOrder{
		SupplierID:  uint64Ptr(9),
		Currency:    "USD",
		Marketplace: "US",
		Items: []domain.PurchaseOrderItem{
			{ProductID: 100, QtyOrdered: 0, UnitCost: 10, Currency: "USD"},
		},
	}

	_, err := uc.Create(nil, order)
	if !errors.Is(err, ErrPurchaseOrderInvalidQty) {
		t.Fatalf("expected invalid qty error, got %v", err)
	}
}

func TestCreatePurchaseOrderRejectsZeroUnitCost(t *testing.T) {
	poRepo := &stubPurchaseOrderRepo{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, nil)

	order := &domain.PurchaseOrder{
		SupplierID:  uint64Ptr(9),
		Currency:    "USD",
		Marketplace: "US",
		Items: []domain.PurchaseOrderItem{
			{ProductID: 100, QtyOrdered: 1, UnitCost: 0, Currency: "USD"},
		},
	}

	_, err := uc.Create(nil, order)
	if !errors.Is(err, ErrPurchaseOrderInvalidUnitCost) {
		t.Fatalf("expected invalid unit cost error, got %v", err)
	}
}

func TestCreatePurchaseOrderRejectsMissingSupplier(t *testing.T) {
	poRepo := &stubPurchaseOrderRepo{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, nil)

	order := &domain.PurchaseOrder{
		Currency:    "USD",
		Marketplace: "US",
		Items: []domain.PurchaseOrderItem{
			{ProductID: 100, QtyOrdered: 1, UnitCost: 10, Currency: "USD"},
		},
	}

	_, err := uc.Create(nil, order)
	if !errors.Is(err, ErrPurchaseOrderMissingSupplier) {
		t.Fatalf("expected missing supplier error, got %v", err)
	}
}

func TestReceivePurchaseOrderUpdatesQtyAndCreatesMovements(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       5,
		PoNumber: "PO-TEST-001",
		Status:   domain.PurchaseOrderStatusShipped,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{ID: 10, ProductID: 200, QtyOrdered: 10, QtyReceived: 2, UnitCost: 2.5},
			{ID: 11, ProductID: 201, QtyOrdered: 5, QtyReceived: 5, UnitCost: 1.0},
		},
	}

	poRepo := &stubPurchaseOrderRepo{getItem: order}
	inventoryService := &stubInventoryService{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventoryService, nil)

	err := uc.Receive(nil, 5, domain.PurchaseOrderReceiveParams{
		WarehouseID:   7,
		ReceivedQties: map[uint64]uint64{10: 3},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if poRepo.updated == nil {
		t.Fatalf("expected updated order")
	}

	if poRepo.updated.Items[0].QtyReceived != 5 {
		t.Fatalf("expected qty_received updated to 5, got %d", poRepo.updated.Items[0].QtyReceived)
	}

	if poRepo.updated.Status != domain.PurchaseOrderStatusShipped {
		t.Fatalf("expected status to remain SHIPPED")
	}

	if len(inventoryService.created) != 1 {
		t.Fatalf("expected 1 movement created, got %d", len(inventoryService.created))
	}

	movement := inventoryService.created[0]
	if movement.ProductID != 200 || movement.Quantity != 3 || movement.WarehouseID != 7 {
		t.Fatalf("unexpected movement: %+v", movement)
	}
}

func TestCreatePurchaseOrderRejectsComboWithoutComponents(t *testing.T) {
	comboID := uint64(20)
	mainID := uint64(300)

	productRepo := &stubProductLookup{
		items: map[uint64]productdomain.Product{
			mainID: {ID: mainID, ComboID: &comboID, IsComboMain: 1},
		},
	}

	comboRepo := &stubComboProvider{
		items: map[uint64][]productdomain.ProductComboItem{
			comboID: {
				{ComboID: comboID, MainProductID: mainID, ProductID: mainID, QtyRatio: 1},
			},
		},
	}

	poRepo := &stubPurchaseOrderRepo{}
	uc := NewPurchaseOrderUsecase(poRepo, productRepo, comboRepo, nil, nil)

	order := &domain.PurchaseOrder{
		SupplierID:  uint64Ptr(9),
		Currency:    "USD",
		Marketplace: "US",
		Items: []domain.PurchaseOrderItem{
			{ProductID: mainID, QtyOrdered: 1, UnitCost: 10, Currency: "USD"},
		},
	}

	_, err := uc.Create(nil, order)
	if err == nil {
		t.Fatalf("expected error for combo without components")
	}
}

func TestCreatePurchaseOrderKeepsComboItemsWhenComponentsProvided(t *testing.T) {
	comboID := uint64(30)
	mainID := uint64(400)
	childID := uint64(401)

	productRepo := &stubProductLookup{
		items: map[uint64]productdomain.Product{
			mainID:  {ID: mainID, ComboID: &comboID, IsComboMain: 1},
			childID: {ID: childID, ComboID: &comboID, IsComboMain: 0},
		},
	}

	comboRepo := &stubComboProvider{
		items: map[uint64][]productdomain.ProductComboItem{
			comboID: {
				{ComboID: comboID, MainProductID: mainID, ProductID: mainID, QtyRatio: 1},
				{ComboID: comboID, MainProductID: mainID, ProductID: childID, QtyRatio: 2},
			},
		},
	}

	poRepo := &stubPurchaseOrderRepo{}
	uc := NewPurchaseOrderUsecase(poRepo, productRepo, comboRepo, nil, nil)

	order := &domain.PurchaseOrder{
		SupplierID:  uint64Ptr(9),
		Currency:    "USD",
		Marketplace: "US",
		Items: []domain.PurchaseOrderItem{
			{ProductID: mainID, QtyOrdered: 1, UnitCost: 10},
			{ProductID: childID, QtyOrdered: 2, UnitCost: 2},
		},
	}

	created, err := uc.Create(nil, order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created == nil || len(created.Items) != 1 {
		t.Fatalf("expected only child items to remain, got %d", len(created.Items))
	}

	var childItem *domain.PurchaseOrderItem
	for i := range created.Items {
		item := &created.Items[i]
		if item.ProductID == childID {
			childItem = item
		}
	}

	if childItem == nil || childItem.QtyOrdered != 2 {
		t.Fatalf("expected child sku qty 2, got %+v", childItem)
	}
}

func TestCreatePurchaseOrderBatchSplitsComboChildrenBySupplierAndAssignsBatchNumbers(t *testing.T) {
	comboID := uint64(40)
	mainID := uint64(500)
	childAID := uint64(501)
	childBID := uint64(502)
	supplierAID := uint64(8)
	supplierBID := uint64(9)

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
	uc := NewPurchaseOrderUsecase(poRepo, productRepo, comboRepo, nil, nil)

	orders, err := uc.CreateBatch(nil, []*domain.PurchaseOrder{
		{
			Currency: "USD",
			Items: []domain.PurchaseOrderItem{
				{ProductID: mainID, QtyOrdered: 1, UnitCost: 0, Currency: "USD"},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 2 {
		t.Fatalf("expected 2 split purchase orders, got %d", len(orders))
	}
	if len(poRepo.createdOrders) != 2 {
		t.Fatalf("expected repo create called twice, got %d", len(poRepo.createdOrders))
	}

	if orders[0].BatchNo == "" || orders[1].BatchNo == "" {
		t.Fatalf("expected batch no assigned to all orders, got %+v", orders)
	}
	if orders[0].BatchNo != orders[1].BatchNo {
		t.Fatalf("expected shared batch no, got %s and %s", orders[0].BatchNo, orders[1].BatchNo)
	}
	if orders[0].PoNumber != orders[0].BatchNo+"-1" {
		t.Fatalf("expected first po number suffixed with -1, got %s", orders[0].PoNumber)
	}
	if orders[1].PoNumber != orders[1].BatchNo+"-2" {
		t.Fatalf("expected second po number suffixed with -2, got %s", orders[1].PoNumber)
	}
	if orders[0].SupplierID == nil || *orders[0].SupplierID != supplierAID {
		t.Fatalf("expected first split order supplier %d, got %+v", supplierAID, orders[0].SupplierID)
	}
	if orders[1].SupplierID == nil || *orders[1].SupplierID != supplierBID {
		t.Fatalf("expected second split order supplier %d, got %+v", supplierBID, orders[1].SupplierID)
	}
	if len(orders[0].Items) != 1 || orders[0].Items[0].ProductID != childAID {
		t.Fatalf("expected first split order to contain child product %d, got %+v", childAID, orders[0].Items)
	}
	if len(orders[1].Items) != 1 || orders[1].Items[0].ProductID != childBID {
		t.Fatalf("expected second split order to contain child product %d, got %+v", childBID, orders[1].Items)
	}
	if orders[1].Items[0].QtyOrdered != 2 {
		t.Fatalf("expected second split order qty 2, got %+v", orders[1].Items[0])
	}
}

func TestSubmitRecordsPurchaseOrderCostEvent(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       5,
		PoNumber: "PO-TEST-005",
		Status:   domain.PurchaseOrderStatusDraft,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{ID: 10, ProductID: 200, QtyOrdered: 8, UnitCost: 2.5, Currency: "USD"},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: order}
	recorder := &stubPurchaseOrderCostEventRecorder{}
	cleaner := &stubReplenishmentPlanCleaner{}

	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, nil)
	uc.BindCostEventRecorder(recorder)
	uc.BindPlanCleaner(cleaner)

	if err := uc.Submit(nil, order.ID); err != nil {
		t.Fatalf("unexpected submit error: %v", err)
	}
	if len(recorder.events) != 1 {
		t.Fatalf("expected 1 cost event, got %d", len(recorder.events))
	}
	got := recorder.events[0]
	if got.EventType != PurchaseOrderCostEventOrdered {
		t.Fatalf("expected event type %s, got %s", PurchaseOrderCostEventOrdered, got.EventType)
	}
	if got.QtyEvent != 8 {
		t.Fatalf("expected qty_event 8, got %d", got.QtyEvent)
	}
	if got.UnitCost != 2.5 {
		t.Fatalf("expected unit_cost 2.5, got %v", got.UnitCost)
	}
	if len(cleaner.deletedPurchaseOrderIDs) != 0 {
		t.Fatalf("expected converted plans to be preserved, got cleanup %+v", cleaner.deletedPurchaseOrderIDs)
	}
}

func TestSubmitIgnoresLegacyPlanCleanerHook(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       51,
		PoNumber: "PO-TX-SUBMIT-51",
		Status:   domain.PurchaseOrderStatusDraft,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{ID: 5101, ProductID: 200, QtyOrdered: 8, UnitCost: 2.5, Currency: "USD"},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: order}
	recorder := &stubPurchaseOrderCostEventRecorder{}
	cleaner := &stubReplenishmentPlanCleaner{err: context.Canceled}
	txManager := &purchaseOrderSubmitTxManagerStub{
		deps: PurchaseOrderSubmitTransactionalDeps{
			Repo:              poRepo,
			PlanCleaner:       cleaner,
			CostEventRecorder: recorder,
		},
	}

	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, nil)
	uc.BindCostEventRecorder(recorder)
	uc.BindPlanCleaner(cleaner)
	uc.BindSubmitTransactionManager(txManager)

	err := uc.Submit(nil, order.ID)
	if err != nil {
		t.Fatalf("unexpected submit error: %v", err)
	}
	if !txManager.called {
		t.Fatalf("expected submit transaction manager to be called")
	}
	if !txManager.committed {
		t.Fatalf("expected submit transaction to commit")
	}
	if poRepo.updated == nil {
		t.Fatalf("expected persisted update")
	}
	if poRepo.getItem.Status != domain.PurchaseOrderStatusOrdered {
		t.Fatalf("expected base order status updated, got %s", poRepo.getItem.Status)
	}
	if len(recorder.events) != 1 {
		t.Fatalf("expected committed cost events, got %+v", recorder.events)
	}
	if len(cleaner.deletedPurchaseOrderIDs) != 0 {
		t.Fatalf("expected legacy plan cleaner hook unused, got %+v", cleaner.deletedPurchaseOrderIDs)
	}
}

func TestMarkShippedAndReceiveRecordCostEvents(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       6,
		PoNumber: "PO-TEST-006",
		Status:   domain.PurchaseOrderStatusOrdered,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{ID: 11, ProductID: 201, QtyOrdered: 10, QtyReceived: 2, UnitCost: 1.3, Currency: "USD"},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: order}
	recorder := &stubPurchaseOrderCostEventRecorder{}
	inventorySvc := &stubInventoryService{}

	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventorySvc, nil)
	uc.BindCostEventRecorder(recorder)

	if err := uc.MarkShipped(nil, order.ID, domain.PurchaseOrderShipParams{WarehouseID: 1}); err != nil {
		t.Fatalf("unexpected mark shipped error: %v", err)
	}
	if len(recorder.events) != 1 {
		t.Fatalf("expected 1 cost event after shipped, got %d", len(recorder.events))
	}
	if recorder.events[0].EventType != PurchaseOrderCostEventShipped {
		t.Fatalf("expected shipped event, got %s", recorder.events[0].EventType)
	}

	// simulate persisted state change for receive path
	order.Status = domain.PurchaseOrderStatusShipped
	order.Items = []domain.PurchaseOrderItem{
		{ID: 11, ProductID: 201, QtyOrdered: 10, QtyReceived: 2, UnitCost: 1.3, Currency: "USD"},
	}
	poRepo.getItem = clonePurchaseOrder(order)

	if err := uc.Receive(nil, order.ID, domain.PurchaseOrderReceiveParams{
		WarehouseID:   1,
		ReceivedQties: map[uint64]uint64{11: 3},
	}); err != nil {
		t.Fatalf("unexpected receive error: %v", err)
	}
	if len(recorder.events) != 2 {
		t.Fatalf("expected 2 cost events total, got %d", len(recorder.events))
	}
	if recorder.events[1].EventType != PurchaseOrderCostEventReceived {
		t.Fatalf("expected received event, got %s", recorder.events[1].EventType)
	}
	if recorder.events[1].QtyEvent != 3 {
		t.Fatalf("expected received qty_event 3, got %d", recorder.events[1].QtyEvent)
	}
}

func TestMarkShippedPersistsWarehouseAndRepeatedReceiveUsesRemainingQty(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       106,
		PoNumber: "PO-TEST-0106",
		Status:   domain.PurchaseOrderStatusOrdered,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{ID: 1061, ProductID: 201, QtyOrdered: 10, QtyReceived: 0, UnitCost: 1.3, Currency: "USD"},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	recorder := &stubPurchaseOrderCostEventRecorder{}
	inventorySvc := &stubInventoryService{}

	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventorySvc, nil)
	uc.BindCostEventRecorder(recorder)

	if err := uc.MarkShipped(nil, order.ID, domain.PurchaseOrderShipParams{WarehouseID: 7}); err != nil {
		t.Fatalf("unexpected mark shipped error: %v", err)
	}

	if poRepo.updated == nil || poRepo.updated.WarehouseID == nil || *poRepo.updated.WarehouseID != 7 {
		t.Fatalf("expected shipped order warehouse persisted as 7, got %+v", poRepo.updated)
	}

	if err := uc.Receive(nil, order.ID, domain.PurchaseOrderReceiveParams{
		ReceivedQties: map[uint64]uint64{1061: 4},
	}); err != nil {
		t.Fatalf("unexpected first receive error: %v", err)
	}

	if poRepo.getItem.Items[0].QtyReceived != 4 {
		t.Fatalf("expected first receive qty 4, got %d", poRepo.getItem.Items[0].QtyReceived)
	}
	if poRepo.getItem.Status != domain.PurchaseOrderStatusShipped {
		t.Fatalf("expected partially received order remain SHIPPED, got %s", poRepo.getItem.Status)
	}

	if err := uc.Receive(nil, order.ID, domain.PurchaseOrderReceiveParams{
		ReceivedQties: map[uint64]uint64{1061: 6},
	}); err != nil {
		t.Fatalf("unexpected second receive error: %v", err)
	}

	if poRepo.getItem.Items[0].QtyReceived != 10 {
		t.Fatalf("expected second receive qty reach 10, got %d", poRepo.getItem.Items[0].QtyReceived)
	}
	if poRepo.getItem.Status != domain.PurchaseOrderStatusReceived {
		t.Fatalf("expected fully received order status RECEIVED, got %s", poRepo.getItem.Status)
	}
	if poRepo.getItem.ReceivedAt == nil {
		t.Fatalf("expected received_at set after full receive")
	}

	if len(inventorySvc.created) != 3 {
		t.Fatalf("expected 3 movements total, got %d", len(inventorySvc.created))
	}
	if inventorySvc.created[1].WarehouseID != 7 || inventorySvc.created[2].WarehouseID != 7 {
		t.Fatalf("expected repeated receive movements to reuse warehouse 7, got %+v %+v", inventorySvc.created[1], inventorySvc.created[2])
	}
}

func TestMarkShippedRollsBackWhenCostEventFails(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       61,
		PoNumber: "PO-TX-SHIP-61",
		Status:   domain.PurchaseOrderStatusOrdered,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{ID: 611, ProductID: 201, QtyOrdered: 4, UnitCost: 1.3, Currency: "USD"},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: order}
	recorder := &stubPurchaseOrderCostEventRecorder{err: context.Canceled}
	inventorySvc := &stubInventoryService{}
	txManager := &purchaseOrderShipTxManagerStub{
		deps: PurchaseOrderShipTransactionalDeps{
			Repo:              poRepo,
			InventoryService:  inventorySvc,
			CostEventRecorder: recorder,
		},
	}

	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventorySvc, nil)
	uc.BindCostEventRecorder(recorder)
	uc.BindShipTransactionManager(txManager)

	err := uc.MarkShipped(nil, order.ID, domain.PurchaseOrderShipParams{WarehouseID: 1})
	if err == nil {
		t.Fatalf("expected mark shipped error")
	}
	if !txManager.called {
		t.Fatalf("expected ship transaction manager to be called")
	}
	if txManager.committed {
		t.Fatalf("expected ship transaction not to commit")
	}
	if len(inventorySvc.created) != 0 {
		t.Fatalf("expected no committed ship movements, got %+v", inventorySvc.created)
	}
	if poRepo.updated != nil {
		t.Fatalf("expected no persisted update, got %+v", poRepo.updated)
	}
	if poRepo.getItem.Status != domain.PurchaseOrderStatusOrdered {
		t.Fatalf("expected base order status unchanged, got %s", poRepo.getItem.Status)
	}
}

func TestReceiveRollsBackWhenCostEventFails(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       62,
		PoNumber: "PO-TX-RECV-62",
		Status:   domain.PurchaseOrderStatusShipped,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{ID: 621, ProductID: 202, QtyOrdered: 5, QtyReceived: 1, UnitCost: 2.1, Currency: "USD"},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: order}
	recorder := &stubPurchaseOrderCostEventRecorder{err: context.DeadlineExceeded}
	inventorySvc := &stubInventoryService{}
	txManager := &purchaseOrderReceiveTxManagerStub{
		deps: PurchaseOrderReceiveTransactionalDeps{
			Repo:              poRepo,
			InventoryService:  inventorySvc,
			CostEventRecorder: recorder,
		},
	}

	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventorySvc, nil)
	uc.BindCostEventRecorder(recorder)
	uc.BindReceiveTransactionManager(txManager)

	err := uc.Receive(nil, order.ID, domain.PurchaseOrderReceiveParams{
		WarehouseID:   1,
		ReceivedQties: map[uint64]uint64{621: 2},
	})
	if err == nil {
		t.Fatalf("expected receive error")
	}
	if !txManager.called {
		t.Fatalf("expected receive transaction manager to be called")
	}
	if txManager.committed {
		t.Fatalf("expected receive transaction not to commit")
	}
	if len(inventorySvc.created) != 0 {
		t.Fatalf("expected no committed receive movements, got %+v", inventorySvc.created)
	}
	if poRepo.updated != nil {
		t.Fatalf("expected no persisted update, got %+v", poRepo.updated)
	}
	if poRepo.getItem.Items[0].QtyReceived != 1 {
		t.Fatalf("expected base order qty_received unchanged, got %d", poRepo.getItem.Items[0].QtyReceived)
	}
}

func TestInspectPurchaseOrderCreatesOnlyReferencedInspectionMovements(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:          71,
		PoNumber:    "PO-INSP-71",
		Status:      domain.PurchaseOrderStatusReceived,
		WarehouseID: uint64Ptr(3),
		Items: []domain.PurchaseOrderItem{
			{
				ID:                711,
				ProductID:         101,
				QtyOrdered:        5,
				QtyReceived:       5,
				QtyInspectionPass: 2,
				Product:           &domain.ProductSnapshot{ID: 101, SellerSku: "SKU-101", Title: "Test Product 101", IsInspectionRequired: 1, IsPackingRequired: 1},
			},
			{
				ID:                712,
				ProductID:         102,
				QtyOrdered:        4,
				QtyReceived:       4,
				QtyInspectionPass: 2,
				Product:           &domain.ProductSnapshot{ID: 102, SellerSku: "SKU-102", Title: "Test Product 102", IsInspectionRequired: 1, IsPackingRequired: 1},
			},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	inventorySvc := &stubInventoryService{}
	auditLogger := &stubPurchaseOrderAuditLogger{}
	txManager := &purchaseOrderInspectTxManagerStub{
		deps: PurchaseOrderInspectTransactionalDeps{
			Repo:             poRepo,
			InventoryService: inventorySvc,
		},
	}

	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventorySvc, auditLogger)
	uc.BindInspectTransactionManager(txManager)

	err := uc.Inspect(newPurchaseOrderAuditContext(), order.ID, domain.PurchaseOrderInspectParams{
		PassQties:  map[uint64]uint64{711: 2},
		FailQties:  map[uint64]uint64{712: 1},
		OperatorID: uint64Ptr(9),
	})
	if err != nil {
		t.Fatalf("unexpected inspect error: %v", err)
	}
	if !txManager.called || !txManager.committed {
		t.Fatalf("expected inspect transaction manager commit")
	}
	if len(inventorySvc.created) != 2 {
		t.Fatalf("expected 2 inspection movements, got %+v", inventorySvc.created)
	}

	passMovement := inventorySvc.created[0]
	if passMovement.MovementType != inventoryDomain.MovementTypeInspectionPass {
		t.Fatalf("expected inspection pass movement, got %+v", passMovement)
	}
	if passMovement.ProductID != 101 || passMovement.WarehouseID != 3 || passMovement.Quantity != 2 {
		t.Fatalf("unexpected pass movement: %+v", passMovement)
	}
	if passMovement.ReferenceID == nil || *passMovement.ReferenceID != order.ID {
		t.Fatalf("expected pass movement reference id %d, got %+v", order.ID, passMovement.ReferenceID)
	}
	if passMovement.ReferenceNumber == nil || *passMovement.ReferenceNumber != order.PoNumber {
		t.Fatalf("expected pass movement reference number %s, got %+v", order.PoNumber, passMovement.ReferenceNumber)
	}

	failMovement := inventorySvc.created[1]
	if failMovement.MovementType != inventoryDomain.MovementTypeInspectionLoss {
		t.Fatalf("expected inspection loss movement, got %+v", failMovement)
	}
	if failMovement.ProductID != 102 || failMovement.WarehouseID != 3 || failMovement.Quantity != 1 {
		t.Fatalf("unexpected fail movement: %+v", failMovement)
	}
}

func TestReceivePurchaseOrderAutoPassesWhenInspectionNotRequired(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       52,
		PoNumber: "PO-TEST-AUTO-QC",
		Status:   domain.PurchaseOrderStatusShipped,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{
				ID:          5201,
				ProductID:   301,
				QtyOrdered:  2,
				QtyReceived: 0,
				UnitCost:    3.2,
				Currency:    "USD",
				Product: &domain.ProductSnapshot{
					ID:                   301,
					SellerSku:            "SKU-301",
					Title:                "免检产品",
					IsInspectionRequired: 0,
					IsPackingRequired:    1,
				},
			},
		},
	}

	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	inventoryService := &stubInventoryService{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventoryService, nil)

	operatorID := uint64(9)
	err := uc.Receive(newPurchaseOrderAuditContext(), order.ID, domain.PurchaseOrderReceiveParams{
		WarehouseID:   7,
		ReceivedQties: map[uint64]uint64{5201: 2},
		OperatorID:    &operatorID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if poRepo.getItem.ReceivedBy == nil || *poRepo.getItem.ReceivedBy != operatorID {
		t.Fatalf("expected received_by %d, got %+v", operatorID, poRepo.getItem.ReceivedBy)
	}
	if poRepo.getItem.Items[0].QtyInspectionPass != 2 {
		t.Fatalf("expected auto inspection pass qty 2, got %d", poRepo.getItem.Items[0].QtyInspectionPass)
	}
	if len(inventoryService.created) != 2 {
		t.Fatalf("expected receive+inspection pass movements, got %d", len(inventoryService.created))
	}
	if inventoryService.created[0].MovementType != inventoryDomain.MovementTypeWarehouseReceive {
		t.Fatalf("expected warehouse receive first, got %s", inventoryService.created[0].MovementType)
	}
	if inventoryService.created[1].MovementType != inventoryDomain.MovementTypeInspectionPass {
		t.Fatalf("expected inspection pass second, got %s", inventoryService.created[1].MovementType)
	}
}

func TestReceivePurchaseOrderAutoMovesPendingShipmentWhenPackingNotRequired(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       521,
		PoNumber: "PO-TEST-NO-PACK",
		Status:   domain.PurchaseOrderStatusShipped,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{
				ID:          5211,
				ProductID:   302,
				QtyOrdered:  3,
				QtyReceived: 0,
				UnitCost:    5.4,
				Currency:    "USD",
				Product: &domain.ProductSnapshot{
					ID:                   302,
					SellerSku:            "SKU-302",
					Title:                "免打包产品",
					IsInspectionRequired: 0,
					IsPackingRequired:    0,
				},
			},
		},
	}

	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	inventoryService := &stubInventoryService{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventoryService, nil)

	operatorID := uint64(11)
	err := uc.Receive(newPurchaseOrderAuditContext(), order.ID, domain.PurchaseOrderReceiveParams{
		WarehouseID:   8,
		ReceivedQties: map[uint64]uint64{5211: 3},
		OperatorID:    &operatorID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(inventoryService.created) != 3 {
		t.Fatalf("expected receive+inspection+packing skip movements, got %d", len(inventoryService.created))
	}
	if inventoryService.created[2].MovementType != inventoryDomain.MovementTypePackingSkipComplete {
		t.Fatalf("expected packing skip movement third, got %s", inventoryService.created[2].MovementType)
	}
	if inventoryService.created[2].Quantity != 3 {
		t.Fatalf("expected packing skip quantity 3, got %+v", inventoryService.created[2])
	}
}

func TestInspectPurchaseOrderRecordsPassFailAndLossWithoutDamagedMovement(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:          53,
		PoNumber:    "PO-TEST-LOSS-QC",
		Status:      domain.PurchaseOrderStatusReceived,
		Currency:    "USD",
		WarehouseID: uint64Ptr(4),
		Items: []domain.PurchaseOrderItem{
			{
				ID:                5301,
				ProductID:         401,
				QtyOrdered:        5,
				QtyReceived:       5,
				QtyInspectionPass: 1,
				QtyInspectionFail: 0,
				UnitCost:          4.5,
				Currency:          "USD",
				Product: &domain.ProductSnapshot{
					ID:                   401,
					SellerSku:            "SKU-401",
					Title:                "质检产品",
					IsInspectionRequired: 1,
					IsPackingRequired:    1,
				},
			},
		},
	}

	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	inventoryService := &stubInventoryService{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventoryService, nil)

	operatorID := uint64(6)
	err := uc.Inspect(newPurchaseOrderAuditContext(), order.ID, domain.PurchaseOrderInspectParams{
		PassQties:  map[uint64]uint64{5301: 2},
		FailQties:  map[uint64]uint64{5301: 1},
		OperatorID: &operatorID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if poRepo.getItem.InspectedBy == nil || *poRepo.getItem.InspectedBy != operatorID {
		t.Fatalf("expected inspected_by %d, got %+v", operatorID, poRepo.getItem.InspectedBy)
	}
	if poRepo.getItem.Items[0].QtyInspectionPass != 3 {
		t.Fatalf("expected inspection pass qty 3, got %d", poRepo.getItem.Items[0].QtyInspectionPass)
	}
	if poRepo.getItem.Items[0].QtyInspectionFail != 1 {
		t.Fatalf("expected inspection fail qty 1, got %d", poRepo.getItem.Items[0].QtyInspectionFail)
	}
	if len(inventoryService.created) != 2 {
		t.Fatalf("expected 2 inspection movements, got %d", len(inventoryService.created))
	}
	if inventoryService.created[0].MovementType != inventoryDomain.MovementTypeInspectionPass {
		t.Fatalf("expected pass movement, got %s", inventoryService.created[0].MovementType)
	}
	if inventoryService.created[1].MovementType != inventoryDomain.MovementTypeInspectionLoss {
		t.Fatalf("expected loss movement, got %s", inventoryService.created[1].MovementType)
	}
}

func TestInspectPurchaseOrderAutoMovesPendingShipmentWhenPackingNotRequired(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:          531,
		PoNumber:    "PO-TEST-INSPECT-NO-PACK",
		Status:      domain.PurchaseOrderStatusReceived,
		Currency:    "USD",
		WarehouseID: uint64Ptr(4),
		Items: []domain.PurchaseOrderItem{
			{
				ID:                5311,
				ProductID:         402,
				QtyOrdered:        4,
				QtyReceived:       4,
				QtyInspectionPass: 0,
				QtyInspectionFail: 0,
				UnitCost:          6.2,
				Currency:          "USD",
				Product: &domain.ProductSnapshot{
					ID:                   402,
					SellerSku:            "SKU-402",
					Title:                "免打包质检产品",
					IsInspectionRequired: 1,
					IsPackingRequired:    0,
				},
			},
		},
	}

	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	inventoryService := &stubInventoryService{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, inventoryService, nil)

	operatorID := uint64(12)
	err := uc.Inspect(newPurchaseOrderAuditContext(), order.ID, domain.PurchaseOrderInspectParams{
		PassQties:  map[uint64]uint64{5311: 3},
		FailQties:  map[uint64]uint64{5311: 1},
		OperatorID: &operatorID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(inventoryService.created) != 3 {
		t.Fatalf("expected pass+skip+loss movements, got %d", len(inventoryService.created))
	}
	if inventoryService.created[0].MovementType != inventoryDomain.MovementTypeInspectionPass {
		t.Fatalf("expected inspection pass first, got %s", inventoryService.created[0].MovementType)
	}
	if inventoryService.created[1].MovementType != inventoryDomain.MovementTypePackingSkipComplete {
		t.Fatalf("expected packing skip second, got %s", inventoryService.created[1].MovementType)
	}
	if inventoryService.created[2].MovementType != inventoryDomain.MovementTypeInspectionLoss {
		t.Fatalf("expected inspection loss third, got %s", inventoryService.created[2].MovementType)
	}
}

func TestCreatePurchaseOrderFallsBackToConfiguredCurrency(t *testing.T) {
	poRepo := &stubPurchaseOrderRepo{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, nil)
	uc.BindDefaultsProvider(stubBaseCurrencyProvider{})

	order := &domain.PurchaseOrder{
		SupplierID: uint64Ptr(9),
		Items: []domain.PurchaseOrderItem{
			{ProductID: 1001, QtyOrdered: 2, UnitCost: 3.5},
		},
	}

	created, err := uc.Create(nil, order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Currency != "EUR" {
		t.Fatalf("expected order currency EUR, got %s", created.Currency)
	}
	if len(created.Items) != 1 || created.Items[0].Currency != "EUR" {
		t.Fatalf("expected item currency EUR, got %+v", created.Items)
	}
}

func uint64Ptr(v uint64) *uint64 {
	return &v
}

func newPurchaseOrderAuditContext() *gin.Context {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("user_id", uint64(1))
	ctx.Set("userID", uint64(1))
	ctx.Set("username", "tester")
	return ctx
}

func TestUpdatePurchaseOrderSkipsAuditWhenNothingChanged(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:          81,
		PoNumber:    "PO-AUDIT-81",
		SupplierID:  uint64Ptr(9),
		Marketplace: "US",
		Status:      domain.PurchaseOrderStatusDraft,
		Currency:    "USD",
		TotalAmount: 24,
		Remark:      "same",
		Supplier:    &domain.SupplierSnapshot{ID: 9, Name: "供应商A"},
		Items: []domain.PurchaseOrderItem{
			{
				ID:         811,
				ProductID:  101,
				QtyOrdered: 2,
				UnitCost:   12,
				Currency:   "USD",
				Subtotal:   24,
				Product:    &domain.ProductSnapshot{ID: 101, SellerSku: "SKU-101", Title: "Test Product 101", IsPackingRequired: 1},
			},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	auditLogger := &stubPurchaseOrderAuditLogger{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, auditLogger)

	updated, err := uc.Update(newPurchaseOrderAuditContext(), order.ID, &domain.PurchaseOrder{
		SupplierID:  uint64Ptr(9),
		Marketplace: "US",
		Currency:    "USD",
		Remark:      "same",
		Items: []domain.PurchaseOrderItem{
			{
				ProductID:  101,
				QtyOrdered: 2,
				UnitCost:   12,
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated == nil {
		t.Fatalf("expected updated order")
	}
	if len(auditLogger.payloads) != 0 {
		t.Fatalf("expected no audit payload for unchanged update, got %d", len(auditLogger.payloads))
	}
}

func TestUpdatePurchaseOrderWritesChangedFieldsOnly(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:          82,
		PoNumber:    "PO-AUDIT-82",
		SupplierID:  uint64Ptr(9),
		Marketplace: "US",
		Status:      domain.PurchaseOrderStatusDraft,
		Currency:    "USD",
		TotalAmount: 24,
		Remark:      "before",
		Supplier:    &domain.SupplierSnapshot{ID: 9, Name: "供应商A"},
		Items: []domain.PurchaseOrderItem{
			{
				ID:         821,
				ProductID:  101,
				QtyOrdered: 2,
				UnitCost:   12,
				Currency:   "USD",
				Subtotal:   24,
				Product:    &domain.ProductSnapshot{ID: 101, SellerSku: "SKU-101", Title: "Test Product 101", IsPackingRequired: 1},
			},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	auditLogger := &stubPurchaseOrderAuditLogger{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, auditLogger)

	_, err := uc.Update(newPurchaseOrderAuditContext(), order.ID, &domain.PurchaseOrder{
		SupplierID:  uint64Ptr(9),
		Marketplace: "US",
		Currency:    "USD",
		Remark:      "after",
		Items: []domain.PurchaseOrderItem{
			{
				ProductID:  101,
				QtyOrdered: 3,
				UnitCost:   12,
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(auditLogger.payloads) != 1 {
		t.Fatalf("expected 1 audit payload, got %d", len(auditLogger.payloads))
	}

	before, ok := auditLogger.payloads[0].Before.(map[string]any)
	if !ok {
		t.Fatalf("expected before diff map, got %#v", auditLogger.payloads[0].Before)
	}
	after, ok := auditLogger.payloads[0].After.(map[string]any)
	if !ok {
		t.Fatalf("expected after diff map, got %#v", auditLogger.payloads[0].After)
	}

	if _, exists := after["currency"]; exists {
		t.Fatalf("did not expect unchanged currency in audit diff: %+v", after)
	}
	if after["remark"] != "after" {
		t.Fatalf("expected remark diff, got %+v", after)
	}
	if before["remark"] != "before" {
		t.Fatalf("expected remark before diff, got %+v", before)
	}
	itemsAfter, ok := after["items"].([]any)
	if !ok || len(itemsAfter) != 1 {
		t.Fatalf("expected item diff, got %+v", after["items"])
	}
	firstItem, ok := itemsAfter[0].(map[string]any)
	if !ok {
		t.Fatalf("expected first item diff map, got %+v", itemsAfter[0])
	}
	if firstItem["qty_ordered"] != float64(3) {
		t.Fatalf("expected changed qty_ordered in audit diff, got %+v", firstItem)
	}
}

func TestUpdatePurchaseOrderAuditIncludesOnlyChangedItems(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:          92,
		PoNumber:    "PO-AUDIT-92",
		SupplierID:  uint64Ptr(9),
		Marketplace: "US",
		Status:      domain.PurchaseOrderStatusDraft,
		Currency:    "USD",
		TotalAmount: 36,
		Remark:      "before",
		Supplier:    &domain.SupplierSnapshot{ID: 9, Name: "供应商A"},
		Items: []domain.PurchaseOrderItem{
			{
				ID:         921,
				ProductID:  101,
				QtyOrdered: 2,
				UnitCost:   12,
				Currency:   "USD",
				Subtotal:   24,
				Product:    &domain.ProductSnapshot{ID: 101, SellerSku: "SKU-101", Title: "Test Product 101", IsPackingRequired: 1},
			},
			{
				ID:         922,
				ProductID:  102,
				QtyOrdered: 1,
				UnitCost:   12,
				Currency:   "USD",
				Subtotal:   12,
				Product:    &domain.ProductSnapshot{ID: 102, SellerSku: "SKU-102", Title: "Test Product 102", IsPackingRequired: 1},
			},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	auditLogger := &stubPurchaseOrderAuditLogger{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, auditLogger)

	_, err := uc.Update(newPurchaseOrderAuditContext(), order.ID, &domain.PurchaseOrder{
		SupplierID:  uint64Ptr(9),
		Marketplace: "US",
		Currency:    "USD",
		Remark:      "before",
		Items: []domain.PurchaseOrderItem{
			{
				ProductID:  101,
				QtyOrdered: 3,
				UnitCost:   12,
			},
			{
				ProductID:  102,
				QtyOrdered: 1,
				UnitCost:   12,
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(auditLogger.payloads) != 1 {
		t.Fatalf("expected 1 audit payload, got %d", len(auditLogger.payloads))
	}

	before, ok := auditLogger.payloads[0].Before.(map[string]any)
	if !ok {
		t.Fatalf("expected before diff map, got %#v", auditLogger.payloads[0].Before)
	}
	after, ok := auditLogger.payloads[0].After.(map[string]any)
	if !ok {
		t.Fatalf("expected after diff map, got %#v", auditLogger.payloads[0].After)
	}

	itemsBefore, ok := before["items"].([]any)
	if !ok || len(itemsBefore) != 1 {
		t.Fatalf("expected exactly one changed item in before diff, got %+v", before["items"])
	}
	itemsAfter, ok := after["items"].([]any)
	if !ok || len(itemsAfter) != 1 {
		t.Fatalf("expected exactly one changed item in after diff, got %+v", after["items"])
	}

	beforeItem, ok := itemsBefore[0].(map[string]any)
	if !ok {
		t.Fatalf("expected before changed item map, got %+v", itemsBefore[0])
	}
	afterItem, ok := itemsAfter[0].(map[string]any)
	if !ok {
		t.Fatalf("expected after changed item map, got %+v", itemsAfter[0])
	}

	if beforeItem["seller_sku"] != "SKU-101" || afterItem["seller_sku"] != "SKU-101" {
		t.Fatalf("expected changed item identity to be preserved, got before=%+v after=%+v", beforeItem, afterItem)
	}
	if _, exists := afterItem["subtotal"]; exists {
		t.Fatalf("did not expect unchanged derived subtotal in changed item diff: %+v", afterItem)
	}
	if afterItem["qty_ordered"] != float64(3) {
		t.Fatalf("expected changed qty_ordered in changed item diff, got %+v", afterItem)
	}
}

func TestSubmitPurchaseOrderWritesTransitionDiff(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:          83,
		PoNumber:    "PO-AUDIT-83",
		SupplierID:  uint64Ptr(9),
		Marketplace: "US",
		Status:      domain.PurchaseOrderStatusDraft,
		Currency:    "USD",
		TotalAmount: 24,
		Remark:      "submit",
		Supplier:    &domain.SupplierSnapshot{ID: 9, Name: "供应商A"},
		Items: []domain.PurchaseOrderItem{
			{
				ID:         831,
				ProductID:  101,
				QtyOrdered: 2,
				UnitCost:   12,
				Currency:   "USD",
				Subtotal:   24,
				Product:    &domain.ProductSnapshot{ID: 101, SellerSku: "SKU-101", Title: "Test Product 101", IsPackingRequired: 1},
			},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	auditLogger := &stubPurchaseOrderAuditLogger{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, auditLogger)

	if err := uc.Submit(newPurchaseOrderAuditContext(), order.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(auditLogger.payloads) != 1 {
		t.Fatalf("expected 1 audit payload, got %d", len(auditLogger.payloads))
	}

	before, ok := auditLogger.payloads[0].Before.(map[string]any)
	if !ok {
		t.Fatalf("expected before diff map, got %#v", auditLogger.payloads[0].Before)
	}
	after, ok := auditLogger.payloads[0].After.(map[string]any)
	if !ok {
		t.Fatalf("expected after diff map, got %#v", auditLogger.payloads[0].After)
	}

	if before["status"] != string(domain.PurchaseOrderStatusDraft) {
		t.Fatalf("expected before status DRAFT, got %+v", before)
	}
	if after["status"] != string(domain.PurchaseOrderStatusOrdered) {
		t.Fatalf("expected after status ORDERED, got %+v", after)
	}
	if _, exists := after["po_number"]; exists {
		t.Fatalf("did not expect unchanged po_number in submit diff: %+v", after)
	}
	if _, exists := after["ordered_at"]; !exists {
		t.Fatalf("expected ordered_at in submit diff, got %+v", after)
	}
}

func TestClosePurchaseOrderSetsClosedAtAndWritesDiff(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:          91,
		PoNumber:    "PO-AUDIT-91",
		SupplierID:  uint64Ptr(9),
		Marketplace: "US",
		Status:      domain.PurchaseOrderStatusReceived,
		Currency:    "USD",
		TotalAmount: 24,
		Remark:      "close",
		Supplier:    &domain.SupplierSnapshot{ID: 9, Name: "供应商A"},
		Items: []domain.PurchaseOrderItem{
			{
				ID:                911,
				ProductID:         101,
				QtyOrdered:        2,
				QtyReceived:       2,
				QtyInspectionPass: 2,
				UnitCost:          12,
				Currency:          "USD",
				Subtotal:          24,
				Product:           &domain.ProductSnapshot{ID: 101, SellerSku: "SKU-101", Title: "Test Product 101", IsPackingRequired: 1},
			},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	auditLogger := &stubPurchaseOrderAuditLogger{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, auditLogger)

	if err := uc.Close(newPurchaseOrderAuditContext(), order.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if poRepo.updated == nil || poRepo.updated.ClosedAt == nil {
		t.Fatalf("expected closed_at to be set on update")
	}
	if len(auditLogger.payloads) != 1 {
		t.Fatalf("expected 1 audit payload, got %d", len(auditLogger.payloads))
	}

	before, ok := auditLogger.payloads[0].Before.(map[string]any)
	if !ok {
		t.Fatalf("expected before diff map, got %#v", auditLogger.payloads[0].Before)
	}
	after, ok := auditLogger.payloads[0].After.(map[string]any)
	if !ok {
		t.Fatalf("expected after diff map, got %#v", auditLogger.payloads[0].After)
	}
	if before["status"] != string(domain.PurchaseOrderStatusReceived) {
		t.Fatalf("expected before status RECEIVED, got %+v", before)
	}
	if after["status"] != string(domain.PurchaseOrderStatusClosed) {
		t.Fatalf("expected after status CLOSED, got %+v", after)
	}
	if _, exists := after["closed_at"]; !exists {
		t.Fatalf("expected closed_at in close diff, got %+v", after)
	}
}

func TestClosePurchaseOrderRejectsPendingInspection(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:          92,
		PoNumber:    "PO-AUDIT-92",
		SupplierID:  uint64Ptr(9),
		Marketplace: "US",
		Status:      domain.PurchaseOrderStatusReceived,
		Currency:    "USD",
		TotalAmount: 24,
		Supplier:    &domain.SupplierSnapshot{ID: 9, Name: "供应商A"},
		Items: []domain.PurchaseOrderItem{
			{
				ID:                921,
				ProductID:         101,
				QtyOrdered:        2,
				QtyReceived:       2,
				QtyInspectionPass: 1,
				UnitCost:          12,
				Currency:          "USD",
				Subtotal:          24,
				Product:           &domain.ProductSnapshot{ID: 101, SellerSku: "SKU-101", Title: "Test Product 101", IsPackingRequired: 1},
			},
		},
	}
	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	auditLogger := &stubPurchaseOrderAuditLogger{}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, auditLogger)

	err := uc.Close(newPurchaseOrderAuditContext(), order.ID)
	if err == nil {
		t.Fatalf("expected close to reject pending inspection")
	}
	if poRepo.updated != nil {
		t.Fatalf("expected no update when pending inspection remains")
	}
	if len(auditLogger.payloads) != 0 {
		t.Fatalf("expected no audit log when close is rejected")
	}
}

func TestClosePurchaseOrderRejectsShortReceipt(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       115,
		PoNumber: "PO-CLOSE-REJECT-SHORT",
		Status:   domain.PurchaseOrderStatusReceived,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{
				ID:                1151,
				ProductID:         501,
				QtyOrdered:        10,
				QtyReceived:       9,
				QtyInspectionPass: 9,
				UnitCost:          12,
				Currency:          "USD",
				Product:           &domain.ProductSnapshot{ID: 501, SellerSku: "SKU-501", Title: "采购异常产品", IsPackingRequired: 1},
			},
		},
	}

	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, nil)

	err := uc.Close(newPurchaseOrderAuditContext(), order.ID)
	if !errors.Is(err, ErrPurchaseOrderIncompleteReceipt) {
		t.Fatalf("expected incomplete receipt error, got %v", err)
	}
	if poRepo.updated != nil {
		t.Fatalf("expected no update when receipt is incomplete")
	}
}

func TestForceCompletePurchaseOrderAllowsShortReceiptWithReason(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       116,
		PoNumber: "PO-FORCE-CLOSE-116",
		Status:   domain.PurchaseOrderStatusShipped,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{
				ID:                1161,
				ProductID:         601,
				QtyOrdered:        10,
				QtyReceived:       9,
				QtyInspectionPass: 9,
				UnitCost:          8.8,
				Currency:          "USD",
				Product:           &domain.ProductSnapshot{ID: 601, SellerSku: "SKU-601", Title: "异常完成产品", IsPackingRequired: 1},
			},
		},
	}

	poRepo := &stubPurchaseOrderRepo{getItem: clonePurchaseOrder(order)}
	uc := NewPurchaseOrderUsecase(poRepo, nil, nil, nil, nil)

	operatorID := uint64(11)
	err := uc.ForceComplete(newPurchaseOrderAuditContext(), order.ID, domain.PurchaseOrderForceCompleteParams{
		Reason:     "供应商少发 1 个，按损失结案",
		OperatorID: &operatorID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if poRepo.getItem.Status != domain.PurchaseOrderStatusClosed {
		t.Fatalf("expected closed status, got %s", poRepo.getItem.Status)
	}
	if poRepo.getItem.IsForceCompleted != 1 {
		t.Fatalf("expected force complete flag set")
	}
	if poRepo.getItem.ForceCompletedBy == nil || *poRepo.getItem.ForceCompletedBy != operatorID {
		t.Fatalf("expected force completed by %d, got %+v", operatorID, poRepo.getItem.ForceCompletedBy)
	}
	if poRepo.getItem.ForceCompleteReason != "供应商少发 1 个，按损失结案" {
		t.Fatalf("expected force complete reason preserved, got %q", poRepo.getItem.ForceCompleteReason)
	}
}
