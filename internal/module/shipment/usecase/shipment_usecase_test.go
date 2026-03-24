package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	productDomain "am-erp-go/internal/module/product/domain"
	"am-erp-go/internal/module/shipment/domain"
)

type stubShipmentRepo struct {
	created *domain.Shipment
	updated *domain.Shipment
	getByID func(id uint64) (*domain.Shipment, error)
}

func (s *stubShipmentRepo) Create(shipment *domain.Shipment) error {
	s.created = shipment
	shipment.ID = 1
	return nil
}
func (s *stubShipmentRepo) Update(shipment *domain.Shipment) error {
	s.updated = shipment
	return nil
}
func (s *stubShipmentRepo) GetByID(id uint64) (*domain.Shipment, error) {
	if s.getByID != nil {
		return s.getByID(id)
	}
	return nil, nil
}
func (s *stubShipmentRepo) GetByShipmentNumber(shipmentNumber string) (*domain.Shipment, error) {
	return nil, nil
}
func (s *stubShipmentRepo) List(params *domain.ShipmentListParams) ([]*domain.Shipment, int64, error) {
	return nil, 0, nil
}
func (s *stubShipmentRepo) Delete(id uint64) error { return nil }

type stubShipmentItemRepo struct {
	items           []domain.ShipmentItem
	getByShipmentID func(shipmentID uint64) ([]domain.ShipmentItem, error)
	updated         []domain.ShipmentItem
	deletedShipmentID uint64
}

func (s *stubShipmentItemRepo) Create(item *domain.ShipmentItem) error { return nil }
func (s *stubShipmentItemRepo) CreateBatch(items []domain.ShipmentItem) error {
	s.items = items
	return nil
}
func (s *stubShipmentItemRepo) UpdateBatch(items []domain.ShipmentItem) error {
	s.updated = items
	return nil
}
func (s *stubShipmentItemRepo) GetByShipmentID(shipmentID uint64) ([]domain.ShipmentItem, error) {
	if s.getByShipmentID != nil {
		return s.getByShipmentID(shipmentID)
	}
	return nil, nil
}
func (s *stubShipmentItemRepo) DeleteByShipmentID(shipmentID uint64) error {
	s.deletedShipmentID = shipmentID
	return nil
}

func uint64Ptr(value uint64) *uint64 {
	return &value
}

type stubShipmentInventoryService struct{}

func (s *stubShipmentInventoryService) CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error) {
	return nil, nil
}
func (s *stubShipmentInventoryService) GetProductBalance(productID, warehouseID uint64) (*inventoryDomain.InventoryBalance, error) {
	return &inventoryDomain.InventoryBalance{PendingShipment: 100, PendingShipmentReserved: 100}, nil
}

type stubShipmentProductRepo struct {
	listByIDs func(ids []uint64) ([]productDomain.Product, error)
}

func (s *stubShipmentProductRepo) ListByIDs(ids []uint64) ([]productDomain.Product, error) {
	if s.listByIDs != nil {
		return s.listByIDs(ids)
	}
	products := make([]productDomain.Product, 0, len(ids))
	for _, id := range ids {
		products = append(products, productDomain.Product{
			ID:        id,
			SellerSku: "PRODUCT",
			Status:    productDomain.ProductStatusOnSale,
		})
	}
	return products, nil
}

type stubShipmentWarehouseRepo struct{}

func (s *stubShipmentWarehouseRepo) GetByID(id uint64) (*inventoryDomain.Warehouse, error) {
	return &inventoryDomain.Warehouse{ID: id}, nil
}

type stubShipmentBaseCurrencyProvider struct{}

func (stubShipmentBaseCurrencyProvider) GetDefaultBaseCurrency() string {
	return "EUR"
}

type stubShipmentCostRecorder struct {
	params *ShipmentCostAllocationRecordParams
	err    error
}

func (s *stubShipmentCostRecorder) RecordShipmentCostAllocation(params *ShipmentCostAllocationRecordParams) error {
	s.params = params
	return s.err
}

type stubShipmentLandedRecorder struct {
	params *ShipmentCostAllocationRecordParams
	err    error
}

func (s *stubShipmentLandedRecorder) UpsertShipmentLandedSnapshots(params *ShipmentCostAllocationRecordParams) error {
	s.params = params
	return s.err
}

type trackingShipmentInventoryService struct {
	movements []*inventoryDomain.CreateMovementParams
	balance   *inventoryDomain.InventoryBalance
}

func (s *trackingShipmentInventoryService) CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error) {
	s.movements = append(s.movements, params)
	return &inventoryDomain.InventoryMovement{}, nil
}

func (s *trackingShipmentInventoryService) GetProductBalance(productID, warehouseID uint64) (*inventoryDomain.InventoryBalance, error) {
	if s.balance != nil {
		return s.balance, nil
	}
	return &inventoryDomain.InventoryBalance{PendingShipment: 100, PendingShipmentReserved: 100}, nil
}

type shipmentMarkShippedTxManagerStub struct {
	deps      ShipmentMarkShippedTransactionalDeps
	called    bool
	committed bool
}

func (s *shipmentMarkShippedTxManagerStub) Run(ctx context.Context, fn func(ShipmentMarkShippedTransactionalDeps) error) error {
	s.called = true
	shipmentRepo := &stubShipmentRepo{getByID: s.deps.ShipmentRepo.(*stubShipmentRepo).getByID}
	itemRepo := &stubShipmentItemRepo{getByShipmentID: s.deps.ShipmentItemRepo.(*stubShipmentItemRepo).getByShipmentID}
	inventory := &trackingShipmentInventoryService{}
	deps := ShipmentMarkShippedTransactionalDeps{
		ShipmentRepo:     shipmentRepo,
		ShipmentItemRepo: itemRepo,
		InventoryService: inventory,
		CostRecorder:     s.deps.CostRecorder,
		LandedRecorder:   s.deps.LandedRecorder,
	}
	if err := fn(deps); err != nil {
		return err
	}
	s.committed = true
	if base, ok := s.deps.ShipmentRepo.(*stubShipmentRepo); ok {
		base.updated = shipmentRepo.updated
	}
	if base, ok := s.deps.ShipmentItemRepo.(*stubShipmentItemRepo); ok {
		base.updated = itemRepo.updated
	}
	if base, ok := s.deps.InventoryService.(*trackingShipmentInventoryService); ok {
		base.movements = inventory.movements
	}
	return nil
}

func TestCreateShipmentFallsBackToConfiguredCurrency(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{}
	itemRepo := &stubShipmentItemRepo{}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		&stubShipmentInventoryService{},
		&stubShipmentProductRepo{},
		&stubShipmentWarehouseRepo{},
	)
	uc.BindDefaultsProvider(stubShipmentBaseCurrencyProvider{})

	shipment, err := uc.Create(nil, &domain.CreateShipmentParams{
		WarehouseID: 1,
		Items: []domain.CreateShipmentItemParams{
			{ProductID: 1001, QuantityPlanned: 5},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(shipment.Items) != 1 || shipment.Items[0].Currency != "EUR" {
		t.Fatalf("expected shipment item currency EUR, got %+v", shipment.Items)
	}
	if len(itemRepo.items) != 1 || itemRepo.items[0].Currency != "EUR" {
		t.Fatalf("expected persisted item currency EUR, got %+v", itemRepo.items)
	}
}

func TestCreateShipmentDefaultsItemUnitCostFromProduct(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{}
	itemRepo := &stubShipmentItemRepo{}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		&stubShipmentInventoryService{},
		&stubShipmentProductRepo{
			listByIDs: func(ids []uint64) ([]productDomain.Product, error) {
				return []productDomain.Product{
					{ID: 1001, SellerSku: "SKU-1001", Status: productDomain.ProductStatusOnSale, UnitCost: float64Ptr(9.9)},
				}, nil
			},
		},
		&stubShipmentWarehouseRepo{},
	)
	uc.BindDefaultsProvider(stubShipmentBaseCurrencyProvider{})

	shipment, err := uc.Create(nil, &domain.CreateShipmentParams{
		WarehouseID: 1,
		Items: []domain.CreateShipmentItemParams{
			{ProductID: 1001, QuantityPlanned: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(shipment.Items) != 1 || shipment.Items[0].UnitCost != 9.9 {
		t.Fatalf("expected shipment item unit cost 9.9, got %+v", shipment.Items)
	}
	if len(itemRepo.items) != 1 || itemRepo.items[0].UnitCost != 9.9 {
		t.Fatalf("expected persisted shipment item unit cost 9.9, got %+v", itemRepo.items)
	}
}

func TestCreateShipmentRejectsInactiveProducts(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{}
	itemRepo := &stubShipmentItemRepo{}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		&stubShipmentInventoryService{},
		&stubShipmentProductRepo{
			listByIDs: func(ids []uint64) ([]productDomain.Product, error) {
				return []productDomain.Product{
					{ID: 1001, SellerSku: "COMBO-MAIN", Status: productDomain.ProductStatusOffShelf},
				}, nil
			},
		},
		&stubShipmentWarehouseRepo{},
	)

	_, err := uc.Create(nil, &domain.CreateShipmentParams{
		WarehouseID: 1,
		Items: []domain.CreateShipmentItemParams{
			{ProductID: 1001, QuantityPlanned: 1},
		},
	})
	if err == nil {
		t.Fatal("expected inactive product validation error")
	}
	if err != nil && !strings.Contains(err.Error(), "on sale or replenishing") {
		t.Fatalf("expected inactive product error, got %v", err)
	}
}

func TestCreateShipmentPersistsDestinationAndLogisticsFields(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{}
	itemRepo := &stubShipmentItemRepo{}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		&stubShipmentInventoryService{},
		&stubShipmentProductRepo{},
		&stubShipmentWarehouseRepo{},
	)

	destinationType := domain.DestinationTypeOwnWarehouse
	destinationWarehouseID := uint64(9)
	logisticsProviderID := uint64(3)
	shippingRateID := uint64(7)
	transportMode := "AIR"
	destinationName := "深圳中转仓"
	destinationCode := "SZX-HUB"
	destinationContact := "张三"
	destinationPhone := "13800000000"
	destinationAddress := "深圳市南山区"
	expectedDeliveryDate := "2026-03-30"
	boxCount := uint(4)
	totalWeight := 12.5
	totalVolume := 0.88
	internalNotes := "fragile"

	_, err := uc.Create(nil, &domain.CreateShipmentParams{
		WarehouseID:            1,
		DestinationWarehouseID: &destinationWarehouseID,
		DestinationType:        &destinationType,
		DestinationName:        &destinationName,
		DestinationCode:        &destinationCode,
		DestinationContact:     &destinationContact,
		DestinationPhone:       &destinationPhone,
		DestinationAddress:     &destinationAddress,
		LogisticsProviderID:    &logisticsProviderID,
		ShippingRateID:         &shippingRateID,
		TransportMode:          &transportMode,
		ExpectedDeliveryDate:   &expectedDeliveryDate,
		BoxCount:               &boxCount,
		TotalWeight:            &totalWeight,
		TotalVolume:            &totalVolume,
		InternalNotes:          &internalNotes,
		Items: []domain.CreateShipmentItemParams{
			{ProductID: 1001, QuantityPlanned: 2},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if shipmentRepo.created == nil {
		t.Fatal("expected shipment to be created")
	}
	if shipmentRepo.created.DestinationWarehouseID == nil || *shipmentRepo.created.DestinationWarehouseID != destinationWarehouseID {
		t.Fatalf("expected destination warehouse id %d, got %+v", destinationWarehouseID, shipmentRepo.created.DestinationWarehouseID)
	}
	if shipmentRepo.created.LogisticsProviderID == nil || *shipmentRepo.created.LogisticsProviderID != logisticsProviderID {
		t.Fatalf("expected logistics provider id %d, got %+v", logisticsProviderID, shipmentRepo.created.LogisticsProviderID)
	}
	if shipmentRepo.created.ShippingRateID == nil || *shipmentRepo.created.ShippingRateID != shippingRateID {
		t.Fatalf("expected shipping rate id %d, got %+v", shippingRateID, shipmentRepo.created.ShippingRateID)
	}
	if shipmentRepo.created.TransportMode == nil || *shipmentRepo.created.TransportMode != transportMode {
		t.Fatalf("expected transport mode %s, got %+v", transportMode, shipmentRepo.created.TransportMode)
	}
	if shipmentRepo.created.ExpectedDeliveryDate == nil || *shipmentRepo.created.ExpectedDeliveryDate != expectedDeliveryDate {
		t.Fatalf("expected expected delivery date %s, got %+v", expectedDeliveryDate, shipmentRepo.created.ExpectedDeliveryDate)
	}
	if shipmentRepo.created.BoxCount != boxCount || shipmentRepo.created.TotalWeight != totalWeight || shipmentRepo.created.TotalVolume != totalVolume {
		t.Fatalf("expected package summary to persist, got %+v", shipmentRepo.created)
	}
	if shipmentRepo.created.InternalNotes == nil || *shipmentRepo.created.InternalNotes != internalNotes {
		t.Fatalf("expected internal notes %s, got %+v", internalNotes, shipmentRepo.created.InternalNotes)
	}
}

func TestUpdateShipmentDraftPersistsLogisticsAndReplacesItems(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:             id,
				ShipmentNumber: "SHP-DRAFT-EDIT-001",
				WarehouseID:    1,
				Status:         domain.ShipmentStatusDraft,
			}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 11, ShipmentID: shipmentID, ProductID: 1001, QuantityPlanned: 1},
			}, nil
		},
	}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		&stubShipmentInventoryService{},
		&stubShipmentProductRepo{},
		&stubShipmentWarehouseRepo{},
	)

	logisticsProviderID := uint64(5)
	shippingRateID := uint64(7)
	transportMode := "SEA"
	destinationName := "美国 FBA"
	remark := "updated"

	updated, err := uc.Update(nil, 1, &domain.UpdateShipmentParams{
		WarehouseID:         uint64Ptr(2),
		DestinationName:     &destinationName,
		LogisticsProviderID: &logisticsProviderID,
		ShippingRateID:      &shippingRateID,
		TransportMode:       &transportMode,
		Remark:              &remark,
		Items: []domain.CreateShipmentItemParams{
			{ProductID: 2002, QuantityPlanned: 3},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated == nil || updated.WarehouseID != 2 {
		t.Fatalf("expected updated warehouse 2, got %+v", updated)
	}
	if shipmentRepo.updated == nil || shipmentRepo.updated.LogisticsProviderID == nil || *shipmentRepo.updated.LogisticsProviderID != logisticsProviderID {
		t.Fatalf("expected logistics provider id %d, got %+v", logisticsProviderID, shipmentRepo.updated)
	}
	if itemRepo.deletedShipmentID != 1 {
		t.Fatalf("expected shipment items to be replaced, delete shipment id=%d", itemRepo.deletedShipmentID)
	}
	if len(itemRepo.items) != 1 || itemRepo.items[0].ProductID != 2002 || itemRepo.items[0].QuantityPlanned != 3 {
		t.Fatalf("expected recreated shipment items, got %+v", itemRepo.items)
	}
}

func TestUpdateShipmentDraftDefaultsItemUnitCostFromProduct(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:             id,
				ShipmentNumber: "SHP-DRAFT-EDIT-UNITCOST",
				WarehouseID:    1,
				Status:         domain.ShipmentStatusDraft,
			}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 11, ShipmentID: shipmentID, ProductID: 1001, QuantityPlanned: 1, UnitCost: 5.5},
			}, nil
		},
	}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		&stubShipmentInventoryService{},
		&stubShipmentProductRepo{
			listByIDs: func(ids []uint64) ([]productDomain.Product, error) {
				return []productDomain.Product{
					{ID: 2002, SellerSku: "SKU-2002", Status: productDomain.ProductStatusOnSale, UnitCost: float64Ptr(12.3)},
				}, nil
			},
		},
		&stubShipmentWarehouseRepo{},
	)
	uc.BindDefaultsProvider(stubShipmentBaseCurrencyProvider{})

	updated, err := uc.Update(nil, 1, &domain.UpdateShipmentParams{
		Items: []domain.CreateShipmentItemParams{
			{ProductID: 2002, QuantityPlanned: 2},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated == nil {
		t.Fatal("expected updated shipment")
	}
	if len(itemRepo.items) != 1 || itemRepo.items[0].UnitCost != 12.3 {
		t.Fatalf("expected recreated shipment item unit cost 12.3, got %+v", itemRepo.items)
	}
}

func TestUpdateShipmentConfirmedOnlyAllowsLogisticsFields(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:                  id,
				ShipmentNumber:      "SHP-CONFIRMED-EDIT-001",
				WarehouseID:         1,
				Status:              domain.ShipmentStatusConfirmed,
				LogisticsProviderID: uint64Ptr(2),
			}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 11, ShipmentID: shipmentID, ProductID: 1001, QuantityPlanned: 1},
			}, nil
		},
	}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		&stubShipmentInventoryService{},
		&stubShipmentProductRepo{},
		&stubShipmentWarehouseRepo{},
	)

	newProviderID := uint64(9)
	newWarehouseID := uint64(8)
	destinationName := "新目的地"

	updated, err := uc.Update(nil, 1, &domain.UpdateShipmentParams{
		WarehouseID:         &newWarehouseID,
		DestinationName:     &destinationName,
		LogisticsProviderID: &newProviderID,
		Items: []domain.CreateShipmentItemParams{
			{ProductID: 2002, QuantityPlanned: 3},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated == nil || updated.WarehouseID != 1 {
		t.Fatalf("expected confirmed shipment warehouse to remain 1, got %+v", updated)
	}
	if shipmentRepo.updated == nil || shipmentRepo.updated.LogisticsProviderID == nil || *shipmentRepo.updated.LogisticsProviderID != newProviderID {
		t.Fatalf("expected logistics provider id %d to update, got %+v", newProviderID, shipmentRepo.updated)
	}
	if itemRepo.deletedShipmentID != 0 || len(itemRepo.items) != 0 {
		t.Fatalf("expected confirmed shipment items to remain untouched, deleted=%d recreated=%+v", itemRepo.deletedShipmentID, itemRepo.items)
	}
}

func TestMarkShippedRecordsShippingCostAllocation(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:             id,
				ShipmentNumber: "SHP-001",
				WarehouseID:    9,
				Status:         domain.ShipmentStatusConfirmed,
			}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 11, ShipmentID: shipmentID, ProductID: 101, QuantityPlanned: 2},
				{ID: 12, ShipmentID: shipmentID, ProductID: 202, QuantityPlanned: 1},
			}, nil
		},
	}
	recorder := &stubShipmentCostRecorder{}
	landedRecorder := &stubShipmentLandedRecorder{}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		&stubShipmentInventoryService{},
		&stubShipmentProductRepo{},
		&stubShipmentWarehouseRepo{},
	)
	uc.BindDefaultsProvider(stubShipmentBaseCurrencyProvider{})
	uc.BindFXResolver(func(baseCurrency, originalCurrency string, occurredAt time.Time) (*ShipmentFXSnapshot, error) {
		return &ShipmentFXSnapshot{
			Rate:        1.2,
			Source:      "MANUAL",
			Version:     "v1",
			EffectiveAt: occurredAt,
		}, nil
	})
	uc.BindCostAllocationRecorder(recorder)
	uc.BindLandedSnapshotRecorder(landedRecorder)

	currency := "USD"
	shippingCost := 30.0
	if err := uc.MarkShipped(nil, 1, &domain.MarkShippedParams{
		ShippingCost: &shippingCost,
		Currency:     &currency,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if shipmentRepo.updated == nil {
		t.Fatalf("expected shipment update")
	}
	if shipmentRepo.updated.BaseCurrency != "EUR" {
		t.Fatalf("expected base currency EUR, got %s", shipmentRepo.updated.BaseCurrency)
	}
	if shipmentRepo.updated.ShippingCostBaseAmount != 36 {
		t.Fatalf("expected base amount 36, got %v", shipmentRepo.updated.ShippingCostBaseAmount)
	}
	if recorder.params == nil || len(recorder.params.Lines) != 2 {
		t.Fatalf("expected 2 allocation lines, got %+v", recorder.params)
	}
	if landedRecorder.params == nil || len(landedRecorder.params.Lines) != 2 {
		t.Fatalf("expected landed recorder to receive 2 lines, got %+v", landedRecorder.params)
	}
	if recorder.params.Lines[0].OriginalAmount != 20 || recorder.params.Lines[1].OriginalAmount != 10 {
		t.Fatalf("unexpected original allocations: %+v", recorder.params.Lines)
	}
	if recorder.params.Lines[0].BaseAmount != 24 || recorder.params.Lines[1].BaseAmount != 12 {
		t.Fatalf("unexpected base allocations: %+v", recorder.params.Lines)
	}
	if recorder.params.Lines[0].ItemCurrency != "" || recorder.params.Lines[0].ItemUnitCost != 0 {
		t.Fatalf("unexpected item cost payload: %+v", recorder.params.Lines[0])
	}
	if len(itemRepo.updated) != 2 {
		t.Fatalf("expected shipment items to be updated, got %+v", itemRepo.updated)
	}
	if itemRepo.updated[0].QuantityShipped != 2 || itemRepo.updated[1].QuantityShipped != 1 {
		t.Fatalf("expected quantity_shipped to persist planned qty, got %+v", itemRepo.updated)
	}
}

func TestMarkShippedRollsBackWhenCostAllocationFails(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:             id,
				ShipmentNumber: "SHP-TX-001",
				WarehouseID:    9,
				Status:         domain.ShipmentStatusConfirmed,
			}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 11, ShipmentID: shipmentID, ProductID: 101, QuantityPlanned: 2},
			}, nil
		},
	}
	inventory := &trackingShipmentInventoryService{}
	costRecorder := &stubShipmentCostRecorder{err: context.DeadlineExceeded}
	landedRecorder := &stubShipmentLandedRecorder{}
	txManager := &shipmentMarkShippedTxManagerStub{
		deps: ShipmentMarkShippedTransactionalDeps{
			ShipmentRepo:     shipmentRepo,
			ShipmentItemRepo: itemRepo,
			InventoryService: inventory,
			CostRecorder:     costRecorder,
			LandedRecorder:   landedRecorder,
		},
	}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		inventory,
		&stubShipmentProductRepo{},
		&stubShipmentWarehouseRepo{},
	)
	uc.BindDefaultsProvider(stubShipmentBaseCurrencyProvider{})
	uc.BindFXResolver(func(baseCurrency, originalCurrency string, occurredAt time.Time) (*ShipmentFXSnapshot, error) {
		return &ShipmentFXSnapshot{Rate: 1.2, Source: "MANUAL", Version: "v1", EffectiveAt: occurredAt}, nil
	})
	uc.BindCostAllocationRecorder(costRecorder)
	uc.BindLandedSnapshotRecorder(landedRecorder)
	uc.BindMarkShippedTransactionManager(txManager)

	currency := "USD"
	shippingCost := 30.0
	err := uc.MarkShipped(nil, 1, &domain.MarkShippedParams{
		ShippingCost: &shippingCost,
		Currency:     &currency,
	})
	if err == nil {
		t.Fatalf("expected cost allocation error")
	}
	if !txManager.called {
		t.Fatalf("expected transaction manager to be called")
	}
	if txManager.committed {
		t.Fatalf("expected transaction not to commit")
	}
	if len(inventory.movements) != 0 {
		t.Fatalf("expected no inventory movement committed, got %+v", inventory.movements)
	}
	if shipmentRepo.updated != nil {
		t.Fatalf("expected shipment not persisted, got %+v", shipmentRepo.updated)
	}
	if len(itemRepo.updated) != 0 {
		t.Fatalf("expected shipment items not persisted, got %+v", itemRepo.updated)
	}
}

func TestConfirmCreatesShipmentAllocateMovement(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:             id,
				ShipmentNumber: "SHP-CONFIRM-001",
				WarehouseID:    9,
				Status:         domain.ShipmentStatusDraft,
			}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 11, ShipmentID: shipmentID, ProductID: 101, QuantityPlanned: 2},
			}, nil
		},
	}
	inventory := &trackingShipmentInventoryService{
		balance: &inventoryDomain.InventoryBalance{PendingShipment: 5},
	}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		inventory,
		&stubShipmentProductRepo{},
		&stubShipmentWarehouseRepo{},
	)

	if err := uc.Confirm(nil, 1, &domain.ConfirmShipmentParams{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shipmentRepo.updated == nil || shipmentRepo.updated.Status != domain.ShipmentStatusConfirmed || !shipmentRepo.updated.InventoryLocked {
		t.Fatalf("expected confirmed shipment update, got %+v", shipmentRepo.updated)
	}
	if len(inventory.movements) != 1 {
		t.Fatalf("expected one inventory movement, got %+v", inventory.movements)
	}
	if inventory.movements[0].MovementType != inventoryDomain.MovementTypeShipmentAllocate || inventory.movements[0].Quantity != 2 {
		t.Fatalf("unexpected confirm movement: %+v", inventory.movements[0])
	}
}

func TestCancelConfirmedShipmentCreatesShipmentReleaseMovement(t *testing.T) {
	shipmentRepo := &stubShipmentRepo{
		getByID: func(id uint64) (*domain.Shipment, error) {
			return &domain.Shipment{
				ID:              id,
				ShipmentNumber:  "SHP-CANCEL-001",
				WarehouseID:     9,
				Status:          domain.ShipmentStatusConfirmed,
				InventoryLocked: true,
			}, nil
		},
	}
	itemRepo := &stubShipmentItemRepo{
		getByShipmentID: func(shipmentID uint64) ([]domain.ShipmentItem, error) {
			return []domain.ShipmentItem{
				{ID: 11, ShipmentID: shipmentID, ProductID: 101, QuantityPlanned: 2},
			}, nil
		},
	}
	inventory := &trackingShipmentInventoryService{
		balance: &inventoryDomain.InventoryBalance{PendingShipmentReserved: 5},
	}
	uc := NewShipmentUsecase(
		shipmentRepo,
		itemRepo,
		inventory,
		&stubShipmentProductRepo{},
		&stubShipmentWarehouseRepo{},
	)

	if err := uc.Cancel(nil, 1, &domain.CancelShipmentParams{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shipmentRepo.updated == nil || shipmentRepo.updated.Status != domain.ShipmentStatusCancelled || shipmentRepo.updated.InventoryLocked {
		t.Fatalf("expected cancelled unlocked shipment update, got %+v", shipmentRepo.updated)
	}
	if len(inventory.movements) != 1 {
		t.Fatalf("expected one inventory movement, got %+v", inventory.movements)
	}
	if inventory.movements[0].MovementType != inventoryDomain.MovementTypeShipmentRelease || inventory.movements[0].Quantity != 2 {
		t.Fatalf("unexpected cancel movement: %+v", inventory.movements[0])
	}
}
