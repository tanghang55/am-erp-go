package usecase

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"am-erp-go/internal/module/inventory/domain"
	packagingDomain "am-erp-go/internal/module/packaging/domain"
)

type inventoryBalanceRepoStub struct {
	balance   *domain.InventoryBalance
	updated   *domain.InventoryBalance
	balances  map[string]*domain.InventoryBalance
	updatedBy map[string]*domain.InventoryBalance
}

func (s *inventoryBalanceRepoStub) GetOrCreate(ctx context.Context, productID, warehouseID uint64) (*domain.InventoryBalance, error) {
	if s.balances != nil {
		key := balanceStubKey(productID, warehouseID)
		if balance, ok := s.balances[key]; ok {
			return balance, nil
		}
		balance := &domain.InventoryBalance{ProductID: productID, WarehouseID: warehouseID}
		s.balances[key] = balance
		return balance, nil
	}
	if s.balance != nil {
		return s.balance, nil
	}
	s.balance = &domain.InventoryBalance{ProductID: productID, WarehouseID: warehouseID}
	return s.balance, nil
}
func (s *inventoryBalanceRepoStub) Update(ctx context.Context, balance *domain.InventoryBalance) error {
	if s.balances != nil {
		key := balanceStubKey(balance.ProductID, balance.WarehouseID)
		s.balances[key] = balance
		if s.updatedBy == nil {
			s.updatedBy = map[string]*domain.InventoryBalance{}
		}
		s.updatedBy[key] = balance
	}
	s.updated = balance
	return nil
}
func (s *inventoryBalanceRepoStub) List(params *domain.BalanceListParams) ([]*domain.InventoryBalance, int64, error) {
	return nil, 0, nil
}
func (s *inventoryBalanceRepoStub) GetByProductAndWarehouse(productID, warehouseID uint64) (*domain.InventoryBalance, error) {
	if s.balances != nil {
		return s.balances[balanceStubKey(productID, warehouseID)], nil
	}
	return s.balance, nil
}

func balanceStubKey(productID, warehouseID uint64) string {
	return fmt.Sprintf("%d-%d", productID, warehouseID)
}

type inventoryMovementRepoStub struct {
	created []*domain.InventoryMovement
}

func (s *inventoryMovementRepoStub) Create(ctx context.Context, movement *domain.InventoryMovement) error {
	if movement.ID == 0 {
		movement.ID = uint64(len(s.created) + 1)
	}
	s.created = append(s.created, movement)
	return nil
}
func (s *inventoryMovementRepoStub) List(params *domain.MovementListParams) ([]*domain.InventoryMovement, int64, error) {
	return nil, 0, nil
}
func (s *inventoryMovementRepoStub) GetByID(id uint64) (*domain.InventoryMovement, error) {
	return nil, nil
}

type inventoryLotRepoStub struct {
	lots               []*domain.InventoryLot
	created            []*domain.InventoryLot
	updated            []*domain.InventoryLot
	enforceUniqueLotNo bool
}

type platformReceiveUnitCostResolverStub struct {
	called bool
	cost   *float64
	err    error
}

type platformReceiveRecorderStub struct {
	validated   bool
	called      bool
	params      *domain.CreateMovementParams
	err         error
	validateErr error
}

type platformReceiveTxManagerStub struct {
	deps          PlatformReceiveTransactionalDeps
	called        bool
	committed     bool
	balanceState  *domain.InventoryBalance
	movementState []*domain.InventoryMovement
	lotState      []*domain.InventoryLot
}

type packingRequirementResolverStub struct {
	requirements []PackingRequirement
	err          error
}

type assemblyTxManagerStub struct {
	deps               AssemblyTransactionalDeps
	called             bool
	committed          bool
	balanceState       *domain.InventoryBalance
	movementState      []*domain.InventoryMovement
	lotState           []*domain.InventoryLot
	packagingItemState map[uint64]*packagingDomain.PackagingItem
	packagingLedgers   []*packagingDomain.PackagingLedger
}

type packagingItemRepoStub struct {
	items map[uint64]*packagingDomain.PackagingItem
}

type packagingLedgerRepoStub struct {
	created []*packagingDomain.PackagingLedger
}

type packingCostRecorderStub struct {
	called bool
	params *PackingMaterialCostRecordParams
	err    error
}

func (s *platformReceiveUnitCostResolverStub) ResolvePlatformReceiveUnitCost(_ context.Context, _ *domain.CreateMovementParams) (*float64, error) {
	s.called = true
	return s.cost, s.err
}

func (s *platformReceiveRecorderStub) ValidatePlatformReceive(_ context.Context, params *domain.CreateMovementParams) error {
	s.validated = true
	s.params = params
	return s.validateErr
}

func (s *platformReceiveRecorderStub) RecordPlatformReceive(_ context.Context, params *domain.CreateMovementParams) error {
	s.called = true
	s.params = params
	return s.err
}

type seedLotUnitCostResolverStub struct {
	called bool
	cost   *float64
	err    error
}

type shipmentLotUnitCostResolverStub struct {
	called bool
	cost   *float64
	err    error
}

func (s *seedLotUnitCostResolverStub) Resolve(_ context.Context, _ uint64, _ uint64) (*float64, error) {
	s.called = true
	return s.cost, s.err
}

func (s *shipmentLotUnitCostResolverStub) Resolve(_ context.Context, _ uint64, _ *string, _ *uint64, _ *string, _ time.Time) (*float64, error) {
	s.called = true
	return s.cost, s.err
}

func (s *packingRequirementResolverStub) ResolvePackingRequirements(_ context.Context, _ uint64) ([]PackingRequirement, error) {
	return s.requirements, s.err
}

func (s *packagingItemRepoStub) List(params *packagingDomain.PackagingItemListParams) ([]packagingDomain.PackagingItem, int64, error) {
	return nil, 0, nil
}

func (s *packagingItemRepoStub) GetByID(id uint64) (*packagingDomain.PackagingItem, error) {
	if s.items == nil {
		return nil, fmt.Errorf("not found")
	}
	item, ok := s.items[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return item, nil
}

func (s *packagingItemRepoStub) Create(item *packagingDomain.PackagingItem) error { return nil }
func (s *packagingItemRepoStub) Update(item *packagingDomain.PackagingItem) error { return nil }
func (s *packagingItemRepoStub) Delete(id uint64) error                           { return nil }
func (s *packagingItemRepoStub) CountReferences(id uint64) (int64, error)         { return 0, nil }
func (s *packagingItemRepoStub) GetLowStockItems() ([]packagingDomain.PackagingItem, error) {
	return nil, nil
}
func (s *packagingItemRepoStub) UpdateQuantity(id uint64, quantity int64) error {
	if s.items == nil {
		return fmt.Errorf("not found")
	}
	item, ok := s.items[id]
	if !ok {
		return fmt.Errorf("not found")
	}
	if quantity < 0 {
		decrease := uint64(-quantity)
		if item.QuantityOnHand < decrease {
			return fmt.Errorf("insufficient quantity")
		}
		item.QuantityOnHand -= decrease
		return nil
	}
	item.QuantityOnHand += uint64(quantity)
	return nil
}

func (s *packagingLedgerRepoStub) List(params *packagingDomain.PackagingLedgerListParams) ([]packagingDomain.PackagingLedger, int64, error) {
	return nil, 0, nil
}

func (s *packagingLedgerRepoStub) GetByID(id uint64) (*packagingDomain.PackagingLedger, error) {
	return nil, nil
}

func (s *packagingLedgerRepoStub) Create(ledger *packagingDomain.PackagingLedger) error {
	s.created = append(s.created, ledger)
	return nil
}

func (s *packagingLedgerRepoStub) GetUsageSummary(dateFrom, dateTo *time.Time) ([]packagingDomain.UsageSummaryItem, error) {
	return nil, nil
}

func (s *packingCostRecorderStub) RecordPackingMaterialCost(params *PackingMaterialCostRecordParams) error {
	s.called = true
	s.params = params
	return s.err
}

func (s *platformReceiveTxManagerStub) Run(ctx context.Context, fn func(PlatformReceiveTransactionalDeps) (*domain.InventoryMovement, error)) (*domain.InventoryMovement, error) {
	s.called = true
	txBalance := &inventoryBalanceRepoStub{}
	if base, ok := s.deps.BalanceRepo.(*inventoryBalanceRepoStub); ok && base.balance != nil {
		clone := *base.balance
		txBalance.balance = &clone
	}
	if base, ok := s.deps.BalanceRepo.(*inventoryBalanceRepoStub); ok && base.balances != nil {
		txBalance.balances = cloneBalances(base.balances)
	}
	txMovement := &inventoryMovementRepoStub{}
	if base, ok := s.deps.MovementRepo.(*inventoryMovementRepoStub); ok {
		txMovement.created = append(txMovement.created, base.created...)
	}
	txLot := &inventoryLotRepoStub{}
	if base, ok := s.deps.LotRepo.(*inventoryLotRepoStub); ok {
		txLot.enforceUniqueLotNo = base.enforceUniqueLotNo
		txLot.lots = cloneLots(base.lots)
	}
	txRecorder := s.deps.Recorder
	deps := PlatformReceiveTransactionalDeps{
		BalanceRepo:  txBalance,
		MovementRepo: txMovement,
		LotRepo:      txLot,
		Recorder:     txRecorder,
	}
	movement, err := fn(deps)
	if err != nil {
		return nil, err
	}
	s.committed = true
	if base, ok := s.deps.BalanceRepo.(*inventoryBalanceRepoStub); ok {
		base.balance = txBalance.balance
		base.updated = txBalance.updated
		if txBalance.balances != nil {
			base.balances = txBalance.balances
			base.updatedBy = txBalance.updatedBy
		}
		s.balanceState = txBalance.balance
	}
	if base, ok := s.deps.MovementRepo.(*inventoryMovementRepoStub); ok {
		base.created = txMovement.created
		s.movementState = txMovement.created
	}
	if base, ok := s.deps.LotRepo.(*inventoryLotRepoStub); ok {
		base.lots = txLot.lots
		base.created = txLot.created
		base.updated = txLot.updated
		s.lotState = txLot.lots
	}
	return movement, nil
}

func (s *assemblyTxManagerStub) Run(ctx context.Context, fn func(AssemblyTransactionalDeps) (*domain.InventoryMovement, error)) (*domain.InventoryMovement, error) {
	s.called = true
	txBalance := &inventoryBalanceRepoStub{}
	if base, ok := s.deps.BalanceRepo.(*inventoryBalanceRepoStub); ok && base.balance != nil {
		clone := *base.balance
		txBalance.balance = &clone
	}
	if base, ok := s.deps.BalanceRepo.(*inventoryBalanceRepoStub); ok && base.balances != nil {
		txBalance.balances = cloneBalances(base.balances)
	}
	txMovement := &inventoryMovementRepoStub{}
	if base, ok := s.deps.MovementRepo.(*inventoryMovementRepoStub); ok {
		txMovement.created = append(txMovement.created, base.created...)
	}
	txLot := &inventoryLotRepoStub{}
	if base, ok := s.deps.LotRepo.(*inventoryLotRepoStub); ok {
		txLot.enforceUniqueLotNo = base.enforceUniqueLotNo
		txLot.lots = cloneLots(base.lots)
	}
	txPackagingItems := &packagingItemRepoStub{}
	if base, ok := s.deps.PackagingItemRepo.(*packagingItemRepoStub); ok {
		txPackagingItems.items = clonePackagingItems(base.items)
	}
	txPackagingLedgers := &packagingLedgerRepoStub{}
	if base, ok := s.deps.PackagingLedgerRepo.(*packagingLedgerRepoStub); ok {
		txPackagingLedgers.created = append(txPackagingLedgers.created, base.created...)
	}
	deps := AssemblyTransactionalDeps{
		BalanceRepo:         txBalance,
		MovementRepo:        txMovement,
		LotRepo:             txLot,
		PackagingItemRepo:   txPackagingItems,
		PackagingLedgerRepo: txPackagingLedgers,
	}
	movement, err := fn(deps)
	if err != nil {
		return nil, err
	}
	s.committed = true
	if base, ok := s.deps.BalanceRepo.(*inventoryBalanceRepoStub); ok {
		base.balance = txBalance.balance
		base.updated = txBalance.updated
		if txBalance.balances != nil {
			base.balances = txBalance.balances
			base.updatedBy = txBalance.updatedBy
		}
		s.balanceState = txBalance.balance
	}
	if base, ok := s.deps.MovementRepo.(*inventoryMovementRepoStub); ok {
		base.created = txMovement.created
		s.movementState = txMovement.created
	}
	if base, ok := s.deps.LotRepo.(*inventoryLotRepoStub); ok {
		base.lots = txLot.lots
		base.created = txLot.created
		base.updated = txLot.updated
		s.lotState = txLot.lots
	}
	if base, ok := s.deps.PackagingItemRepo.(*packagingItemRepoStub); ok {
		base.items = txPackagingItems.items
		s.packagingItemState = txPackagingItems.items
	}
	if base, ok := s.deps.PackagingLedgerRepo.(*packagingLedgerRepoStub); ok {
		base.created = txPackagingLedgers.created
		s.packagingLedgers = txPackagingLedgers.created
	}
	return movement, nil
}

func cloneLots(src []*domain.InventoryLot) []*domain.InventoryLot {
	if len(src) == 0 {
		return nil
	}
	dst := make([]*domain.InventoryLot, 0, len(src))
	for _, lot := range src {
		if lot == nil {
			dst = append(dst, nil)
			continue
		}
		cp := *lot
		if lot.SourceType != nil {
			v := *lot.SourceType
			cp.SourceType = &v
		}
		if lot.SourceNumber != nil {
			v := *lot.SourceNumber
			cp.SourceNumber = &v
		}
		if lot.UnitCost != nil {
			v := *lot.UnitCost
			cp.UnitCost = &v
		}
		if lot.Remark != nil {
			v := *lot.Remark
			cp.Remark = &v
		}
		dst = append(dst, &cp)
	}
	return dst
}

func cloneBalances(src map[string]*domain.InventoryBalance) map[string]*domain.InventoryBalance {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]*domain.InventoryBalance, len(src))
	for key, balance := range src {
		if balance == nil {
			dst[key] = nil
			continue
		}
		cp := *balance
		if balance.LastMovementAt != nil {
			t := *balance.LastMovementAt
			cp.LastMovementAt = &t
		}
		dst[key] = &cp
	}
	return dst
}

func clonePackagingItems(src map[uint64]*packagingDomain.PackagingItem) map[uint64]*packagingDomain.PackagingItem {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[uint64]*packagingDomain.PackagingItem, len(src))
	for key, item := range src {
		if item == nil {
			continue
		}
		clone := *item
		dst[key] = &clone
	}
	return dst
}

func (s *inventoryLotRepoStub) List(params *domain.InventoryLotListParams) ([]*domain.InventoryLot, int64, error) {
	return s.lots, int64(len(s.lots)), nil
}
func (s *inventoryLotRepoStub) ListByProductAndWarehouse(ctx context.Context, productID, warehouseID uint64) ([]*domain.InventoryLot, error) {
	return s.lots, nil
}
func (s *inventoryLotRepoStub) Create(ctx context.Context, lot *domain.InventoryLot) error {
	if s.enforceUniqueLotNo {
		for _, existing := range s.lots {
			if existing.LotNo == lot.LotNo {
				return fmt.Errorf("duplicate lot no: %s", lot.LotNo)
			}
		}
	}
	if lot.ID == 0 {
		lot.ID = uint64(len(s.lots) + len(s.created) + 1)
	}
	s.created = append(s.created, lot)
	s.lots = append(s.lots, lot)
	return nil
}
func (s *inventoryLotRepoStub) Update(ctx context.Context, lot *domain.InventoryLot) error {
	s.updated = append(s.updated, lot)
	return nil
}

func TestCreateMovementPlatformReceiveCreatesSellableLot(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:          1001,
			WarehouseID:        9,
			LogisticsInTransit: 5,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	unitCost := 7.5
	referenceType := "SHIPMENT"
	referenceID := uint64(77)

	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   9,
		MovementType:  domain.MovementTypePlatformReceive,
		Quantity:      3,
		ReferenceType: &referenceType,
		ReferenceID:   &referenceID,
		UnitCost:      &unitCost,
		OperatedAt:    timePtr(time.Date(2026, 3, 9, 18, 0, 0, 0, time.Local)),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if movement.AfterSellable != 3 {
		t.Fatalf("expected after sellable 3, got %d", movement.AfterSellable)
	}
	if len(lotRepo.created) != 1 || lotRepo.created[0].QtySellable != 3 {
		t.Fatalf("expected one sellable lot, got %+v", lotRepo.created)
	}
}

func TestCreateMovementSalesAllocateSellableUsesSellablePool(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:   1001,
			WarehouseID: 9,
			Sellable:    5,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{ID: 1, ProductID: 1001, WarehouseID: 9, QtySellable: 5, Status: domain.InventoryLotStatusOpen},
		},
	}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	pool := domain.StockPoolSellable

	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypeSalesAllocate,
		Quantity:     2,
		StockPool:    &pool,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if movement.AfterSellable != 3 || movement.AfterSellableReserved != 2 {
		t.Fatalf("unexpected sellable movement result: %+v", movement)
	}
	if lotRepo.lots[0].QtySellable != 3 || lotRepo.lots[0].QtySellableReserved != 2 {
		t.Fatalf("unexpected lot state: %+v", lotRepo.lots[0])
	}
}

func TestCreateMovementAssemblyCompleteConsumesPackagingAndCreatesPendingShipmentLot(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balances: map[string]*domain.InventoryBalance{
			balanceStubKey(1001, 9): {ProductID: 1001, WarehouseID: 9, RawMaterial: 6, PendingShipment: 0},
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	sourceType := "INIT_RAW_MATERIAL"
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{ID: 1, ProductID: 1001, WarehouseID: 9, SourceType: &sourceType, QtyRawMaterial: 6, UnitCost: float64Ptr(4.5), Status: domain.InventoryLotStatusOpen},
		},
	}
	packagingItemRepo := &packagingItemRepoStub{
		items: map[uint64]*packagingDomain.PackagingItem{
			501: {ID: 501, ItemCode: "BOX-1", ItemName: "纸箱", QuantityOnHand: 10, UnitCost: 1.2, Status: "ACTIVE"},
			502: {ID: 502, ItemCode: "TAPE-1", ItemName: "胶带", QuantityOnHand: 3, UnitCost: 0.8, Status: "ACTIVE"},
		},
	}
	packagingLedgerRepo := &packagingLedgerRepoStub{}
	txManager := &assemblyTxManagerStub{
		deps: AssemblyTransactionalDeps{
			BalanceRepo:         balanceRepo,
			MovementRepo:        movementRepo,
			LotRepo:             lotRepo,
			PackagingItemRepo:   packagingItemRepo,
			PackagingLedgerRepo: packagingLedgerRepo,
		},
	}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	packingCostRecorder := &packingCostRecorderStub{}
	uc.BindPackingRequirementResolver(&packingRequirementResolverStub{
		requirements: []PackingRequirement{
			{PackagingItemID: 501, QuantityPerUnit: 2, ItemCode: "BOX-1", ItemName: "纸箱"},
			{PackagingItemID: 502, QuantityPerUnit: 1, ItemCode: "TAPE-1", ItemName: "胶带"},
		},
	})
	uc.BindAssemblyTransactionManager(txManager)
	uc.BindPackingCostRecorder(packingCostRecorder)
	unitCost := 6.2

	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypeAssemblyComplete,
		Quantity:     3,
		UnitCost:     &unitCost,
		OperatorID:   uint64Ptr(1),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !txManager.called || !txManager.committed {
		t.Fatalf("expected assembly transaction committed")
	}
	if movement.MovementType != domain.MovementTypeAssemblyComplete || movement.ProductID != 1001 {
		t.Fatalf("expected main assembly movement, got %+v", movement)
	}
	if len(movementRepo.created) != 2 {
		t.Fatalf("expected 2 movements, got %d", len(movementRepo.created))
	}
	if movementRepo.created[0].MovementType != domain.MovementTypeAssemblyConsume || movementRepo.created[0].ProductID != 1001 || movementRepo.created[0].Quantity != 3 {
		t.Fatalf("unexpected consume movement: %+v", movementRepo.created[0])
	}
	if movementRepo.created[1].MovementType != domain.MovementTypeAssemblyComplete || movementRepo.created[1].ProductID != 1001 || movementRepo.created[1].Quantity != 3 {
		t.Fatalf("unexpected output movement: %+v", movementRepo.created[1])
	}
	if movementRepo.created[0].TraceID == nil || movementRepo.created[1].TraceID == nil {
		t.Fatalf("expected shared trace ids")
	}
	if *movementRepo.created[0].TraceID != *movementRepo.created[1].TraceID {
		t.Fatalf("expected same trace id across assembly movements")
	}
	if balanceRepo.balances[balanceStubKey(1001, 9)].RawMaterial != 3 {
		t.Fatalf("expected product raw material 3, got %+v", balanceRepo.balances[balanceStubKey(1001, 9)])
	}
	if balanceRepo.balances[balanceStubKey(1001, 9)].PendingShipment != 3 {
		t.Fatalf("expected main pending shipment 3, got %+v", balanceRepo.balances[balanceStubKey(1001, 9)])
	}
	if lotRepo.lots[0].QtyRawMaterial != 3 || lotRepo.lots[0].QtyConsumed != 3 {
		t.Fatalf("unexpected raw material lot: %+v", lotRepo.lots[0])
	}
	if len(lotRepo.created) != 1 || lotRepo.created[0].ProductID != 1001 || lotRepo.created[0].QtyPendingShipment != 3 {
		t.Fatalf("expected one main output pending shipment lot, got %+v", lotRepo.created)
	}
	if packagingItemRepo.items[501].QuantityOnHand != 4 || packagingItemRepo.items[502].QuantityOnHand != 0 {
		t.Fatalf("unexpected packaging stock: %+v %+v", packagingItemRepo.items[501], packagingItemRepo.items[502])
	}
	if len(packagingLedgerRepo.created) != 2 {
		t.Fatalf("expected 2 packaging ledgers, got %d", len(packagingLedgerRepo.created))
	}
	if packagingLedgerRepo.created[0].Quantity >= 0 || packagingLedgerRepo.created[1].Quantity >= 0 {
		t.Fatalf("expected outbound packaging ledgers, got %+v", packagingLedgerRepo.created)
	}
	if !packingCostRecorder.called || packingCostRecorder.params == nil {
		t.Fatalf("expected packing cost recorder called")
	}
	if packingCostRecorder.params.ProductID != 1001 || packingCostRecorder.params.Quantity != 3 {
		t.Fatalf("unexpected packing cost params: %+v", packingCostRecorder.params)
	}
	if len(packingCostRecorder.params.Lines) != 2 {
		t.Fatalf("expected 2 packing cost lines, got %+v", packingCostRecorder.params)
	}
}

func TestCreateMovementAssemblyCompleteRejectsWhenPackingMaterialsMissing(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balances: map[string]*domain.InventoryBalance{
			balanceStubKey(1001, 9): {ProductID: 1001, WarehouseID: 9, RawMaterial: 6},
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	txManager := &assemblyTxManagerStub{
		deps: AssemblyTransactionalDeps{
			BalanceRepo:         balanceRepo,
			MovementRepo:        movementRepo,
			LotRepo:             lotRepo,
			PackagingItemRepo:   &packagingItemRepoStub{},
			PackagingLedgerRepo: &packagingLedgerRepoStub{},
		},
	}

	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindPackingRequirementResolver(&packingRequirementResolverStub{})
	uc.BindAssemblyTransactionManager(txManager)

	_, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypeAssemblyComplete,
		Quantity:     3,
		OperatorID:   uint64Ptr(1),
	})
	if err == nil {
		t.Fatal("expected missing packing materials error")
	}
	if !errors.Is(err, ErrPackingMaterialsNotConfigured) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateMovementAssemblyCompleteRejectsWhenQuantityExceedsRawMaterial(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balances: map[string]*domain.InventoryBalance{
			balanceStubKey(1001, 9): {ProductID: 1001, WarehouseID: 9, RawMaterial: 2},
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	sourceType := "INIT_RAW_MATERIAL"
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{ID: 1, ProductID: 1001, WarehouseID: 9, SourceType: &sourceType, QtyRawMaterial: 2, Status: domain.InventoryLotStatusOpen},
		},
	}
	txManager := &assemblyTxManagerStub{
		deps: AssemblyTransactionalDeps{
			BalanceRepo:  balanceRepo,
			MovementRepo: movementRepo,
			LotRepo:      lotRepo,
			PackagingItemRepo: &packagingItemRepoStub{
				items: map[uint64]*packagingDomain.PackagingItem{
					501: {ID: 501, ItemCode: "BOX-1", ItemName: "纸箱", QuantityOnHand: 10, Status: "ACTIVE"},
				},
			},
			PackagingLedgerRepo: &packagingLedgerRepoStub{},
		},
	}

	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindPackingRequirementResolver(&packingRequirementResolverStub{
		requirements: []PackingRequirement{
			{PackagingItemID: 501, QuantityPerUnit: 1, ItemCode: "BOX-1", ItemName: "纸箱"},
		},
	})
	uc.BindAssemblyTransactionManager(txManager)

	_, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypeAssemblyComplete,
		Quantity:     3,
		OperatorID:   uint64Ptr(1),
	})
	if err == nil {
		t.Fatal("expected insufficient raw material error")
	}
	if !strings.Contains(err.Error(), "原料库存不足") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateMovementPurchasePipelineMaintainsLotState(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:   1004,
			WarehouseID: 9,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	unitCost := 11.2
	referenceType := "PURCHASE_ORDER"
	referenceID := uint64(501)
	referenceNumber := "PO-UAT-501"

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:       1004,
		WarehouseID:     9,
		MovementType:    domain.MovementTypePurchaseShip,
		Quantity:        5,
		ReferenceType:   &referenceType,
		ReferenceID:     &referenceID,
		ReferenceNumber: &referenceNumber,
		UnitCost:        &unitCost,
	}); err != nil {
		t.Fatalf("purchase ship failed: %v", err)
	}
	if len(lotRepo.lots) != 1 {
		t.Fatalf("expected one purchasing in transit lot, got %+v", lotRepo.lots)
	}
	if lotRepo.lots[0].QtyPurchasingInTransit != 5 {
		t.Fatalf("expected purchasing in transit qty 5, got %+v", lotRepo.lots[0])
	}

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:       1004,
		WarehouseID:     9,
		MovementType:    domain.MovementTypeWarehouseReceive,
		Quantity:        5,
		ReferenceType:   &referenceType,
		ReferenceID:     &referenceID,
		ReferenceNumber: &referenceNumber,
		UnitCost:        &unitCost,
	}); err != nil {
		t.Fatalf("warehouse receive failed: %v", err)
	}
	if lotRepo.lots[0].QtyPurchasingInTransit != 0 || lotRepo.lots[0].QtyPendingInspection != 5 {
		t.Fatalf("expected lot moved to pending inspection, got %+v", lotRepo.lots[0])
	}

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:       1004,
		WarehouseID:     9,
		MovementType:    domain.MovementTypeInspectionPass,
		Quantity:        5,
		ReferenceType:   &referenceType,
		ReferenceID:     &referenceID,
		ReferenceNumber: &referenceNumber,
		UnitCost:        &unitCost,
	}); err != nil {
		t.Fatalf("inspection pass failed: %v", err)
	}
	if lotRepo.lots[0].QtyPendingInspection != 0 || lotRepo.lots[0].QtyRawMaterial != 5 {
		t.Fatalf("expected lot moved to raw material, got %+v", lotRepo.lots[0])
	}

}

func TestCreateMovementShipmentAllocateUsesPendingShipmentPool(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:               1001,
			WarehouseID:             9,
			PendingShipment:         5,
			PendingShipmentReserved: 0,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{ID: 1, ProductID: 1001, WarehouseID: 9, QtyPendingShipment: 5, Status: domain.InventoryLotStatusOpen},
		},
	}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)

	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypeShipmentAllocate,
		Quantity:     2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if movement.AfterPendingShipment != 3 || movement.AfterPendingShipmentReserved != 2 {
		t.Fatalf("unexpected shipment allocate movement result: %+v", movement)
	}
	if lotRepo.lots[0].QtyPendingShipment != 3 || lotRepo.lots[0].QtyPendingShipmentReserved != 2 {
		t.Fatalf("unexpected pending shipment lot state: %+v", lotRepo.lots[0])
	}
}

func TestWarehouseReceivePrefersMatchingPurchaseOrderLot(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:               2006,
			WarehouseID:             3,
			PurchasingInTransit:     7,
			PendingInspection:       0,
			AvailableQuantity:       0,
			ReservedQuantity:        0,
			DamagedQuantity:         0,
			RawMaterial:             0,
			PendingShipment:         0,
			PendingShipmentReserved: 0,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	oldReferenceID := uint64(701)
	oldReferenceNumber := "PO-OLD-701"
	newReferenceID := uint64(702)
	newReferenceNumber := "PO-NEW-702"
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{
				ID:                     1,
				ProductID:              2006,
				WarehouseID:            3,
				LotNo:                  "LOT-OLD-1",
				SourceID:               &oldReferenceID,
				SourceNumber:           &oldReferenceNumber,
				QtyPurchasingInTransit: 2,
				QtyPendingInspection:   0,
			},
			{
				ID:                     2,
				ProductID:              2006,
				WarehouseID:            3,
				LotNo:                  "LOT-NEW-2",
				SourceID:               &newReferenceID,
				SourceNumber:           &newReferenceNumber,
				QtyPurchasingInTransit: 5,
				QtyPendingInspection:   0,
			},
		},
	}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	unitCost := 3.2
	referenceType := "PURCHASE_ORDER"

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:       2006,
		WarehouseID:     3,
		MovementType:    domain.MovementTypeWarehouseReceive,
		Quantity:        5,
		ReferenceType:   &referenceType,
		ReferenceID:     &newReferenceID,
		ReferenceNumber: &newReferenceNumber,
		UnitCost:        &unitCost,
	}); err != nil {
		t.Fatalf("warehouse receive failed: %v", err)
	}

	if lotRepo.lots[0].QtyPurchasingInTransit != 2 || lotRepo.lots[0].QtyPendingInspection != 0 {
		t.Fatalf("expected old lot untouched, got %+v", lotRepo.lots[0])
	}
	if lotRepo.lots[1].QtyPurchasingInTransit != 0 || lotRepo.lots[1].QtyPendingInspection != 5 {
		t.Fatalf("expected matching lot moved to pending inspection, got %+v", lotRepo.lots[1])
	}
}

func TestInspectionPassPrefersMatchingPurchaseOrderLot(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:         2007,
			WarehouseID:       3,
			PendingInspection: 7,
			RawMaterial:       0,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	oldReferenceID := uint64(801)
	oldReferenceNumber := "PO-OLD-801"
	newReferenceID := uint64(802)
	newReferenceNumber := "PO-NEW-802"
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{
				ID:                   1,
				ProductID:            2007,
				WarehouseID:          3,
				LotNo:                "LOT-OLD-1",
				SourceID:             &oldReferenceID,
				SourceNumber:         &oldReferenceNumber,
				QtyPendingInspection: 2,
				QtyRawMaterial:       0,
			},
			{
				ID:                   2,
				ProductID:            2007,
				WarehouseID:          3,
				LotNo:                "LOT-NEW-2",
				SourceID:             &newReferenceID,
				SourceNumber:         &newReferenceNumber,
				QtyPendingInspection: 5,
				QtyRawMaterial:       0,
			},
		},
	}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	referenceType := "PURCHASE_ORDER"

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:       2007,
		WarehouseID:     3,
		MovementType:    domain.MovementTypeInspectionPass,
		Quantity:        5,
		ReferenceType:   &referenceType,
		ReferenceID:     &newReferenceID,
		ReferenceNumber: &newReferenceNumber,
	}); err != nil {
		t.Fatalf("inspection pass failed: %v", err)
	}

	if lotRepo.lots[0].QtyPendingInspection != 2 || lotRepo.lots[0].QtyRawMaterial != 0 {
		t.Fatalf("expected old lot untouched, got %+v", lotRepo.lots[0])
	}
	if lotRepo.lots[1].QtyPendingInspection != 0 || lotRepo.lots[1].QtyRawMaterial != 5 {
		t.Fatalf("expected matching lot moved to raw material, got %+v", lotRepo.lots[1])
	}
}

func TestCreateMovementShipmentReleaseReturnsReservedPendingShipment(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:               1001,
			WarehouseID:             9,
			PendingShipment:         1,
			PendingShipmentReserved: 4,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{ID: 1, ProductID: 1001, WarehouseID: 9, QtyPendingShipment: 1, QtyPendingShipmentReserved: 4, Status: domain.InventoryLotStatusOpen},
		},
	}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)

	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypeShipmentRelease,
		Quantity:     2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if movement.AfterPendingShipment != 3 || movement.AfterPendingShipmentReserved != 2 {
		t.Fatalf("unexpected shipment release movement result: %+v", movement)
	}
	if lotRepo.lots[0].QtyPendingShipment != 3 || lotRepo.lots[0].QtyPendingShipmentReserved != 2 {
		t.Fatalf("unexpected pending shipment lot state: %+v", lotRepo.lots[0])
	}
}

func TestCreateMovementShipmentShipConsumesReservedPendingShipment(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:               1001,
			WarehouseID:             9,
			PendingShipmentReserved: 3,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{ID: 1, ProductID: 1001, WarehouseID: 9, QtyPendingShipmentReserved: 3, UnitCost: float64Ptr(8.8), Status: domain.InventoryLotStatusOpen},
		},
	}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)

	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypeShipmentShip,
		Quantity:     2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if movement.AfterPendingShipmentReserved != 1 {
		t.Fatalf("expected after pending shipment reserved 1, got %+v", movement)
	}
	if len(movement.LotAllocations) != 1 || movement.LotAllocations[0].Qty != 2 || movement.LotAllocations[0].UnitCost != 8.8 {
		t.Fatalf("unexpected shipment ship allocations: %+v", movement.LotAllocations)
	}
	if lotRepo.lots[0].QtyPendingShipmentReserved != 1 || lotRepo.lots[0].QtyConsumed != 2 {
		t.Fatalf("unexpected pending shipment lot state: %+v", lotRepo.lots[0])
	}
}

func TestBuildLotNoStaysWithinSchemaLimitWithLongReferenceNumber(t *testing.T) {
	uc := NewInventoryUsecase(&inventoryBalanceRepoStub{}, &inventoryMovementRepoStub{}, &inventoryLotRepoStub{})
	ref := "THIS-IS-A-VERY-LONG-REFERENCE-NUMBER-FOR-ASSEMBLY-THAT-MUST-NOT-BE-IN-LOT-NO"
	lotNo := uc.buildLotNo(&domain.CreateMovementParams{
		ProductID:       1001,
		WarehouseID:     9,
		MovementType:    domain.MovementTypeAssemblyComplete,
		ReferenceNumber: &ref,
	}, time.Date(2026, 3, 10, 13, 0, 0, 0, time.Local))
	if len(lotNo) > 64 {
		t.Fatalf("expected lot_no length <= 64, got %d: %s", len(lotNo), lotNo)
	}
	matched, err := regexp.MatchString(`^LOT202603101300\d{4}$`, lotNo)
	if err != nil {
		t.Fatalf("regexp error: %v", err)
	}
	if !matched {
		t.Fatalf("expected readable lot number format, got %s", lotNo)
	}
}

func TestCreateMovementPlatformReceiveResolvesUnitCostFromShipmentWhenMissing(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:          1001,
			WarehouseID:        9,
			LogisticsInTransit: 5,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	resolvedCost := 9.8765
	resolver := &platformReceiveUnitCostResolverStub{cost: &resolvedCost}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindPlatformReceiveUnitCostResolver(resolver)
	referenceType := "SHIPMENT"
	referenceID := uint64(88)

	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   9,
		MovementType:  domain.MovementTypePlatformReceive,
		Quantity:      3,
		ReferenceType: &referenceType,
		ReferenceID:   &referenceID,
		OperatedAt:    timePtr(time.Date(2026, 3, 9, 18, 30, 0, 0, time.Local)),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resolver.called {
		t.Fatalf("expected platform receive cost resolver to be called")
	}
	if movement.UnitCost == nil || *movement.UnitCost != resolvedCost {
		t.Fatalf("expected movement unit cost %.4f, got %+v", resolvedCost, movement.UnitCost)
	}
	if len(lotRepo.created) != 1 || lotRepo.created[0].UnitCost == nil || *lotRepo.created[0].UnitCost != resolvedCost {
		t.Fatalf("expected sellable lot unit cost %.4f, got %+v", resolvedCost, lotRepo.created)
	}
}

func TestCreateMovementPlatformReceiveRecordsShipmentReceipt(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:          1001,
			WarehouseID:        9,
			LogisticsInTransit: 3,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	recorder := &platformReceiveRecorderStub{}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindPlatformReceiveRecorder(recorder)
	referenceType := "SHIPMENT"
	referenceID := uint64(88)

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   9,
		MovementType:  domain.MovementTypePlatformReceive,
		Quantity:      2,
		ReferenceType: &referenceType,
		ReferenceID:   &referenceID,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !recorder.called || recorder.params == nil {
		t.Fatalf("expected platform receive recorder to be called")
	}
	if recorder.params.Quantity != 2 || recorder.params.ProductID != 1001 {
		t.Fatalf("unexpected recorder params: %+v", recorder.params)
	}
}

func TestSellableLotResolvedCostFlowsIntoSalesShipAllocation(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:          1001,
			WarehouseID:        9,
			LogisticsInTransit: 2,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	resolvedCost := 9.8765
	resolver := &platformReceiveUnitCostResolverStub{cost: &resolvedCost}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindPlatformReceiveUnitCostResolver(resolver)
	referenceType := "SHIPMENT"
	referenceID := uint64(88)
	pool := domain.StockPoolSellable

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   9,
		MovementType:  domain.MovementTypePlatformReceive,
		Quantity:      2,
		ReferenceType: &referenceType,
		ReferenceID:   &referenceID,
	}); err != nil {
		t.Fatalf("platform receive failed: %v", err)
	}
	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypeSalesAllocate,
		Quantity:     2,
		StockPool:    &pool,
	}); err != nil {
		t.Fatalf("sales allocate failed: %v", err)
	}
	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypeSalesShip,
		Quantity:     2,
		StockPool:    &pool,
	})
	if err != nil {
		t.Fatalf("sales ship failed: %v", err)
	}
	if len(movement.LotAllocations) != 1 {
		t.Fatalf("expected one sellable lot allocation, got %+v", movement.LotAllocations)
	}
	if movement.LotAllocations[0].UnitCost != resolvedCost {
		t.Fatalf("expected sellable ship unit cost %.4f, got %.4f", resolvedCost, movement.LotAllocations[0].UnitCost)
	}
	if lotRepo.lots[0].QtyConsumed != 2 {
		t.Fatalf("expected lot consumed qty 2, got %+v", lotRepo.lots[0])
	}
}

func TestCreateMovementPlatformReceiveRequiresShipmentReference(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:          1001,
			WarehouseID:        9,
			LogisticsInTransit: 3,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)

	_, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1001,
		WarehouseID:  9,
		MovementType: domain.MovementTypePlatformReceive,
		Quantity:     2,
	})
	if err == nil {
		t.Fatalf("expected shipment reference validation error")
	}
}

func TestCreateMovementPlatformReceiveValidatesBeforeInventoryWrite(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:          1001,
			WarehouseID:        9,
			LogisticsInTransit: 3,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	recorder := &platformReceiveRecorderStub{validateErr: ErrPlatformReceiveRequiresShipmentReference}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindPlatformReceiveRecorder(recorder)
	referenceType := "SHIPMENT"
	referenceID := uint64(88)

	_, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   9,
		MovementType:  domain.MovementTypePlatformReceive,
		Quantity:      2,
		ReferenceType: &referenceType,
		ReferenceID:   &referenceID,
	})
	if err == nil {
		t.Fatalf("expected platform receive validation error")
	}
	if !recorder.validated {
		t.Fatalf("expected platform receive validator to be called")
	}
	if recorder.called {
		t.Fatalf("expected recorder not to run after validation failure")
	}
	if len(movementRepo.created) != 0 {
		t.Fatalf("expected no movement persisted, got %+v", movementRepo.created)
	}
	if len(lotRepo.created) != 0 {
		t.Fatalf("expected no lot created, got %+v", lotRepo.created)
	}
	if balanceRepo.updated != nil {
		t.Fatalf("expected balance not updated, got %+v", balanceRepo.updated)
	}
}

func TestCreateMovementPlatformReceiveAllowsRepeatedPartialReceivesForSameShipment(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:          1001,
			WarehouseID:        9,
			LogisticsInTransit: 2,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{enforceUniqueLotNo: true}
	recorder := &platformReceiveRecorderStub{}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindPlatformReceiveRecorder(recorder)
	referenceType := "SHIPMENT"
	referenceID := uint64(88)
	referenceNumber := "SH-TEST-0001"

	for i := 0; i < 2; i++ {
		if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
			ProductID:       1001,
			WarehouseID:     9,
			MovementType:    domain.MovementTypePlatformReceive,
			Quantity:        1,
			ReferenceType:   &referenceType,
			ReferenceID:     &referenceID,
			ReferenceNumber: &referenceNumber,
		}); err != nil {
			t.Fatalf("platform receive #%d failed: %v", i+1, err)
		}
	}
	if len(lotRepo.created) != 2 {
		t.Fatalf("expected two inbound sellable lots, got %+v", lotRepo.created)
	}
	if lotRepo.created[0].LotNo == lotRepo.created[1].LotNo {
		t.Fatalf("expected unique lot numbers for repeated platform receives, got %s", lotRepo.created[0].LotNo)
	}
}

func TestCreateMovementPlatformReceiveRollsBackWhenRecorderFails(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:          1001,
			WarehouseID:        9,
			LogisticsInTransit: 2,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	recorder := &platformReceiveRecorderStub{err: fmt.Errorf("receipt write failed")}
	txManager := &platformReceiveTxManagerStub{
		deps: PlatformReceiveTransactionalDeps{
			BalanceRepo:  balanceRepo,
			MovementRepo: movementRepo,
			LotRepo:      lotRepo,
			Recorder:     recorder,
		},
	}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindPlatformReceiveRecorder(recorder)
	uc.BindPlatformReceiveTransactionManager(txManager)
	referenceType := "SHIPMENT"
	referenceID := uint64(88)

	_, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:     1001,
		WarehouseID:   9,
		MovementType:  domain.MovementTypePlatformReceive,
		Quantity:      1,
		ReferenceType: &referenceType,
		ReferenceID:   &referenceID,
	})
	if err == nil {
		t.Fatalf("expected recorder error")
	}
	if !txManager.called {
		t.Fatalf("expected transaction manager to be called")
	}
	if txManager.committed {
		t.Fatalf("expected transaction not to commit")
	}
	if len(movementRepo.created) != 0 {
		t.Fatalf("expected no movement persisted, got %+v", movementRepo.created)
	}
	if len(lotRepo.created) != 0 || len(lotRepo.lots) != 0 {
		t.Fatalf("expected no lot persisted, got created=%+v lots=%+v", lotRepo.created, lotRepo.lots)
	}
	if balanceRepo.updated != nil {
		t.Fatalf("expected balance not updated, got %+v", balanceRepo.updated)
	}
	if balanceRepo.balance == nil || balanceRepo.balance.LogisticsInTransit != 2 || balanceRepo.balance.Sellable != 0 {
		t.Fatalf("expected original balance untouched, got %+v", balanceRepo.balance)
	}
}

func TestCreateMovementSalesAllocateSellableSeedsInitLotWithResolvedCost(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:   1002,
			WarehouseID: 9,
			Sellable:    2,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	lotRepo := &inventoryLotRepoStub{}
	resolvedCost := 14.0
	resolver := &seedLotUnitCostResolverStub{cost: &resolvedCost}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindSeedLotUnitCostResolver(resolver.Resolve)
	pool := domain.StockPoolSellable

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1002,
		WarehouseID:  9,
		MovementType: domain.MovementTypeSalesAllocate,
		Quantity:     1,
		StockPool:    &pool,
	}); err != nil {
		t.Fatalf("sales allocate failed: %v", err)
	}
	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1002,
		WarehouseID:  9,
		MovementType: domain.MovementTypeSalesShip,
		Quantity:     1,
		StockPool:    &pool,
	})
	if err != nil {
		t.Fatalf("sales ship failed: %v", err)
	}
	if !resolver.called {
		t.Fatalf("expected seed lot resolver to be called")
	}
	if len(lotRepo.created) != 1 || lotRepo.created[0].UnitCost == nil || *lotRepo.created[0].UnitCost != resolvedCost {
		t.Fatalf("expected seeded lot cost %.4f, got %+v", resolvedCost, lotRepo.created)
	}
	if len(movement.LotAllocations) != 1 || movement.LotAllocations[0].UnitCost != resolvedCost {
		t.Fatalf("expected shipped allocation cost %.4f, got %+v", resolvedCost, movement.LotAllocations)
	}
}

func TestCreateMovementSalesAllocateSellableBackfillsExistingInitLotCost(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:   1003,
			WarehouseID: 9,
			Sellable:    2,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	sourceType := "INIT_SELLABLE"
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{ID: 1, ProductID: 1003, WarehouseID: 9, SourceType: &sourceType, QtySellable: 2, Status: domain.InventoryLotStatusOpen},
		},
	}
	resolvedCost := 13.25
	resolver := &seedLotUnitCostResolverStub{cost: &resolvedCost}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindSeedLotUnitCostResolver(resolver.Resolve)
	pool := domain.StockPoolSellable

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1003,
		WarehouseID:  9,
		MovementType: domain.MovementTypeSalesAllocate,
		Quantity:     1,
		StockPool:    &pool,
	}); err != nil {
		t.Fatalf("sales allocate failed: %v", err)
	}
	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1003,
		WarehouseID:  9,
		MovementType: domain.MovementTypeSalesShip,
		Quantity:     1,
		StockPool:    &pool,
	})
	if err != nil {
		t.Fatalf("sales ship failed: %v", err)
	}
	if !resolver.called {
		t.Fatalf("expected seed lot resolver to be called")
	}
	if lotRepo.lots[0].UnitCost == nil || *lotRepo.lots[0].UnitCost != resolvedCost {
		t.Fatalf("expected backfilled lot cost %.4f, got %+v", resolvedCost, lotRepo.lots[0])
	}
	if len(movement.LotAllocations) != 1 || movement.LotAllocations[0].UnitCost != resolvedCost {
		t.Fatalf("expected shipped allocation cost %.4f, got %+v", resolvedCost, movement.LotAllocations)
	}
}

func TestCreateMovementSalesAllocateSellableBackfillsExistingShipmentLotCost(t *testing.T) {
	balanceRepo := &inventoryBalanceRepoStub{
		balance: &domain.InventoryBalance{
			ProductID:   1004,
			WarehouseID: 9,
			Sellable:    2,
		},
	}
	movementRepo := &inventoryMovementRepoStub{}
	sourceType := "SHIPMENT"
	sourceID := uint64(88)
	sourceNumber := "SH-88"
	originalCost := 11.3
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{
				ID:                  1,
				ProductID:           1004,
				WarehouseID:         9,
				SourceType:          &sourceType,
				SourceID:            &sourceID,
				SourceNumber:        &sourceNumber,
				ReceivedAt:          time.Date(2026, 3, 23, 14, 26, 22, 0, time.FixedZone("CST", 8*3600)),
				UnitCost:            &originalCost,
				QtySellable:         2,
				QtySellableReserved: 0,
				Status:              domain.InventoryLotStatusOpen,
			},
		},
	}
	resolvedCost := 11.7662
	resolver := &shipmentLotUnitCostResolverStub{cost: &resolvedCost}
	uc := NewInventoryUsecase(balanceRepo, movementRepo, lotRepo)
	uc.BindShipmentLotUnitCostResolver(resolver.Resolve)
	pool := domain.StockPoolSellable

	if _, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1004,
		WarehouseID:  9,
		MovementType: domain.MovementTypeSalesAllocate,
		Quantity:     1,
		StockPool:    &pool,
	}); err != nil {
		t.Fatalf("sales allocate failed: %v", err)
	}
	movement, err := uc.CreateMovement(context.Background(), &domain.CreateMovementParams{
		ProductID:    1004,
		WarehouseID:  9,
		MovementType: domain.MovementTypeSalesShip,
		Quantity:     1,
		StockPool:    &pool,
	})
	if err != nil {
		t.Fatalf("sales ship failed: %v", err)
	}
	if !resolver.called {
		t.Fatalf("expected shipment lot resolver to be called")
	}
	if lotRepo.lots[0].UnitCost == nil || *lotRepo.lots[0].UnitCost != resolvedCost {
		t.Fatalf("expected backfilled shipment lot cost %.4f, got %+v", resolvedCost, lotRepo.lots[0])
	}
	if len(movement.LotAllocations) != 1 || movement.LotAllocations[0].UnitCost != resolvedCost {
		t.Fatalf("expected shipped allocation cost %.4f, got %+v", resolvedCost, movement.LotAllocations)
	}
}

func TestListLotsBackfillsExistingShipmentLotCost(t *testing.T) {
	sourceType := "SHIPMENT"
	sourceID := uint64(101)
	sourceNumber := "SH-101"
	originalCost := 11.3
	lotRepo := &inventoryLotRepoStub{
		lots: []*domain.InventoryLot{
			{
				ID:           1,
				ProductID:    1005,
				WarehouseID:  9,
				SourceType:   &sourceType,
				SourceID:     &sourceID,
				SourceNumber: &sourceNumber,
				ReceivedAt:   time.Date(2026, 3, 23, 14, 26, 22, 0, time.FixedZone("CST", 8*3600)),
				UnitCost:     &originalCost,
				QtySellable:  1,
				Status:       domain.InventoryLotStatusOpen,
			},
		},
	}
	resolvedCost := 11.7662
	resolver := &shipmentLotUnitCostResolverStub{cost: &resolvedCost}
	uc := NewInventoryUsecase(&inventoryBalanceRepoStub{}, &inventoryMovementRepoStub{}, lotRepo)
	uc.BindShipmentLotUnitCostResolver(resolver.Resolve)

	lots, total, err := uc.ListLots(&domain.InventoryLotListParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("list lots failed: %v", err)
	}
	if total != 1 || len(lots) != 1 {
		t.Fatalf("unexpected lots result: total=%d lots=%+v", total, lots)
	}
	if !resolver.called {
		t.Fatalf("expected shipment lot resolver to be called")
	}
	if lots[0].UnitCost == nil || *lots[0].UnitCost != resolvedCost {
		t.Fatalf("expected listed lot cost %.4f, got %+v", resolvedCost, lots[0])
	}
}

func timePtr(v time.Time) *time.Time {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}

func uint64Ptr(v uint64) *uint64 {
	return &v
}

func stringPtr(v string) *string {
	return &v
}
