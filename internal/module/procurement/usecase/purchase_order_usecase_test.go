package usecase

import (
	"context"
	"testing"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/procurement/domain"
	productdomain "am-erp-go/internal/module/product/domain"
)

type stubPurchaseOrderRepo struct {
	created *domain.PurchaseOrder
	updated *domain.PurchaseOrder
	getID   uint64
	getItem *domain.PurchaseOrder
	err     error
	createErrors   []error
	createCalls    int
	createNumbers  []string
}

func (s *stubPurchaseOrderRepo) List(params *domain.PurchaseOrderListParams) ([]domain.PurchaseOrder, int64, error) {
	return nil, 0, s.err
}

func (s *stubPurchaseOrderRepo) GetByID(id uint64) (*domain.PurchaseOrder, error) {
	s.getID = id
	return s.getItem, s.err
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
		order.ID = 1
	}
	return s.err
}

func (s *stubPurchaseOrderRepo) Update(order *domain.PurchaseOrder) error {
	s.updated = order
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
}

func (s *stubInventoryService) CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error) {
	if params != nil {
		s.created = append(s.created, params)
	}
	return &inventoryDomain.InventoryMovement{
		ID:       1,
		SkuID:    params.SkuID,
		Quantity: params.Quantity,
	}, nil
}

func floatPtr(v float64) *float64 {
	return &v
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
			{SkuID: mainID, QtyOrdered: 5, UnitCost: 10, Currency: "USD"},
			{SkuID: 200, QtyOrdered: 1, UnitCost: 2.5, Currency: "USD"},
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

	if len(poRepo.created.Items) != 3 {
		t.Fatalf("expected 3 items after combo expansion, got %d", len(poRepo.created.Items))
	}

	var item200 *domain.PurchaseOrderItem
	var item201 *domain.PurchaseOrderItem
	var itemMain *domain.PurchaseOrderItem
	for i := range poRepo.created.Items {
		item := &poRepo.created.Items[i]
		if item.SkuID == 200 {
			item200 = item
		}
		if item.SkuID == 201 {
			item201 = item
		}
		if item.SkuID == mainID {
			itemMain = item
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
	if itemMain == nil || itemMain.QtyOrdered != 5 {
		t.Fatalf("expected main sku qty 5, got %+v", itemMain)
	}

	if poRepo.created.TotalAmount != 42.5 {
		t.Fatalf("expected total amount 42.5, got %.2f", poRepo.created.TotalAmount)
	}
}

func TestReceivePurchaseOrderUpdatesQtyAndCreatesMovements(t *testing.T) {
	order := &domain.PurchaseOrder{
		ID:       5,
		PoNumber: "PO-TEST-001",
		Status:   domain.PurchaseOrderStatusShipped,
		Currency: "USD",
		Items: []domain.PurchaseOrderItem{
			{ID: 10, SkuID: 200, QtyOrdered: 10, QtyReceived: 2, UnitCost: 2.5},
			{ID: 11, SkuID: 201, QtyOrdered: 5, QtyReceived: 5, UnitCost: 1.0},
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
	if movement.SkuID != 200 || movement.Quantity != 3 || movement.WarehouseID != 7 {
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
			{SkuID: mainID, QtyOrdered: 1, UnitCost: 10, Currency: "USD"},
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
			{SkuID: mainID, QtyOrdered: 1, UnitCost: 10},
			{SkuID: childID, QtyOrdered: 2, UnitCost: 2},
		},
	}

	created, err := uc.Create(nil, order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created == nil || len(created.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(created.Items))
	}

	var mainItem *domain.PurchaseOrderItem
	var childItem *domain.PurchaseOrderItem
	for i := range created.Items {
		item := &created.Items[i]
		if item.SkuID == mainID {
			mainItem = item
		}
		if item.SkuID == childID {
			childItem = item
		}
	}

	if mainItem == nil || mainItem.QtyOrdered != 1 {
		t.Fatalf("expected main sku qty 1, got %+v", mainItem)
	}
	if childItem == nil || childItem.QtyOrdered != 2 {
		t.Fatalf("expected child sku qty 2, got %+v", childItem)
	}
}

func uint64Ptr(v uint64) *uint64 {
	return &v
}
