package usecase

import (
	"context"
	"errors"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/sales/domain"
	systemdomain "am-erp-go/internal/module/system/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type stubSalesOrderRepo struct {
	order     *domain.SalesOrder
	updated   *domain.SalesOrder
	created   *domain.SalesOrder
	err       error
	getErr    error
	createErr error
	updateErr error
}

func (s *stubSalesOrderRepo) List(params *domain.SalesOrderListParams) ([]domain.SalesOrder, int64, error) {
	if s.order == nil {
		return []domain.SalesOrder{}, 0, s.err
	}
	return []domain.SalesOrder{*s.order}, 1, s.err
}

func TestBuildImportBatchNoUsesReadableNumericFormat(t *testing.T) {
	value := buildImportBatchNo()
	matched, err := regexp.MatchString(`^IMP\d{16}$`, value)
	if err != nil {
		t.Fatalf("regexp error: %v", err)
	}
	if !matched {
		t.Fatalf("unexpected import batch number: %s", value)
	}
}

func (s *stubSalesOrderRepo) GetByID(id uint64) (*domain.SalesOrder, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.order, s.err
}

func (s *stubSalesOrderRepo) Create(order *domain.SalesOrder) error {
	s.created = order
	if s.createErr != nil {
		return s.createErr
	}
	return s.err
}

func (s *stubSalesOrderRepo) Update(order *domain.SalesOrder) error {
	s.updated = order
	if s.updateErr != nil {
		return s.updateErr
	}
	return s.err
}

type stubSalesInventoryService struct {
	movement *inventoryDomain.InventoryMovement
	err      error
	params   []*inventoryDomain.CreateMovementParams
}

func (s *stubSalesInventoryService) CreateMovement(_ context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error) {
	if params != nil {
		s.params = append(s.params, params)
	}
	if s.movement != nil {
		return s.movement, s.err
	}
	return &inventoryDomain.InventoryMovement{}, s.err
}

type trackingSalesInventoryService struct {
	params []*inventoryDomain.CreateMovementParams
}

func (s *trackingSalesInventoryService) CreateMovement(_ context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error) {
	if params != nil {
		s.params = append(s.params, params)
	}
	qty := uint64(0)
	if params != nil && params.Quantity > 0 {
		qty = uint64(params.Quantity)
	}
	return &inventoryDomain.InventoryMovement{
		OperatedAt: time.Now(),
		LotAllocations: []inventoryDomain.InventoryLotAllocation{
			{InventoryLotID: 9001, Qty: qty, UnitCost: 3.2},
		},
	}, nil
}

type stubSalesShipCostRecorder struct {
	records        []*SalesShipCostRecordParams
	returns        []*SalesReturnCostRecordParams
	resolves       []*SalesReturnCostRecordParams
	cogs           float64
	returnUnitCost *float64
	err            error
}

func (s *stubSalesShipCostRecorder) RecordSalesShipCost(params *SalesShipCostRecordParams) error {
	if params != nil {
		s.records = append(s.records, params)
	}
	return s.err
}

func (s *stubSalesShipCostRecorder) ResolveSalesReturnUnitCost(params *SalesReturnCostRecordParams) (*float64, error) {
	if params != nil {
		s.resolves = append(s.resolves, params)
	}
	return s.returnUnitCost, s.err
}

func (s *stubSalesShipCostRecorder) RecordSalesReturnCost(params *SalesReturnCostRecordParams) (float64, error) {
	if params != nil {
		s.returns = append(s.returns, params)
	}
	return s.cogs, s.err
}

type stubSalesShipProfitRecorder struct {
	records []*SalesShipProfitRecordParams
	returns []*SalesReturnProfitRecordParams
	err     error
}

func (s *stubSalesShipProfitRecorder) RecordSalesShipProfit(params *SalesShipProfitRecordParams) error {
	if params != nil {
		s.records = append(s.records, params)
	}
	return s.err
}

type stubSalesAuditLogger struct {
	payloads []systemUsecase.AuditLogPayload
}

func (s *stubSalesAuditLogger) RecordFromContext(_ *gin.Context, payload systemUsecase.AuditLogPayload) error {
	s.payloads = append(s.payloads, payload)
	return nil
}

func (s *stubSalesShipProfitRecorder) RecordSalesReturnProfit(params *SalesReturnProfitRecordParams) error {
	if params != nil {
		s.returns = append(s.returns, params)
	}
	return s.err
}

type stubSalesConfigProvider struct{}

func (stubSalesConfigProvider) GetDefaultBaseCurrency() string {
	return "EUR"
}

func (stubSalesConfigProvider) GetSalesImportDefaults() systemdomain.ConfigCenterSalesImport {
	return systemdomain.ConfigCenterSalesImport{
		DefaultChannel:     "MANUAL-IMPORT",
		DefaultMarketplace: "DE",
	}
}

type salesShipTxManagerStub struct {
	deps      SalesShipTransactionalDeps
	called    bool
	committed bool
}

func (s *salesShipTxManagerStub) Run(ctx context.Context, fn func(SalesShipTransactionalDeps) error) error {
	s.called = true
	repo := &stubSalesOrderRepo{order: cloneSalesOrder(s.deps.Repo.(*stubSalesOrderRepo).order)}
	inventory := &trackingSalesInventoryService{}
	deps := SalesShipTransactionalDeps{
		Repo:             repo,
		InventoryService: inventory,
		CostWriter:       s.deps.CostWriter,
		ProfitWriter:     s.deps.ProfitWriter,
	}
	if err := fn(deps); err != nil {
		return err
	}
	s.committed = true
	if base, ok := s.deps.Repo.(*stubSalesOrderRepo); ok {
		base.order = repo.order
		base.updated = repo.updated
	}
	if base, ok := s.deps.InventoryService.(*trackingSalesInventoryService); ok {
		base.params = inventory.params
	}
	return nil
}

func cloneSalesOrder(order *domain.SalesOrder) *domain.SalesOrder {
	if order == nil {
		return nil
	}
	cp := *order
	if len(order.Items) > 0 {
		cp.Items = make([]domain.SalesOrderItem, len(order.Items))
		copy(cp.Items, order.Items)
	}
	if order.ExternalOrderNo != nil {
		v := *order.ExternalOrderNo
		cp.ExternalOrderNo = &v
	}
	if order.SalesChannel != nil {
		v := *order.SalesChannel
		cp.SalesChannel = &v
	}
	if order.Marketplace != nil {
		v := *order.Marketplace
		cp.Marketplace = &v
	}
	if order.Remark != nil {
		v := *order.Remark
		cp.Remark = &v
	}
	if order.ImportBatchNo != nil {
		v := *order.ImportBatchNo
		cp.ImportBatchNo = &v
	}
	return &cp
}

type salesAllocateTxManagerStub struct {
	deps      SalesAllocateTransactionalDeps
	called    bool
	committed bool
}

func (s *salesAllocateTxManagerStub) Run(ctx context.Context, fn func(SalesAllocateTransactionalDeps) error) error {
	s.called = true
	baseRepo, _ := s.deps.Repo.(*stubSalesOrderRepo)
	repo := &stubSalesOrderRepo{
		order:     cloneSalesOrder(baseRepo.order),
		updateErr: baseRepo.updateErr,
	}
	inventory := &trackingSalesInventoryService{}
	deps := SalesAllocateTransactionalDeps{
		Repo:             repo,
		InventoryService: inventory,
	}
	if err := fn(deps); err != nil {
		return err
	}
	s.committed = true
	if baseRepo != nil {
		baseRepo.order = repo.order
		baseRepo.updated = repo.updated
	}
	if baseInventory, ok := s.deps.InventoryService.(*trackingSalesInventoryService); ok {
		baseInventory.params = inventory.params
	}
	return nil
}

type salesReturnTxManagerStub struct {
	deps      SalesReturnTransactionalDeps
	called    bool
	committed bool
}

func (s *salesReturnTxManagerStub) Run(ctx context.Context, fn func(SalesReturnTransactionalDeps) error) error {
	s.called = true
	baseRepo, _ := s.deps.Repo.(*stubSalesOrderRepo)
	repo := &stubSalesOrderRepo{
		order:     cloneSalesOrder(baseRepo.order),
		updateErr: baseRepo.updateErr,
	}
	inventory := &trackingSalesInventoryService{}
	deps := SalesReturnTransactionalDeps{
		Repo:             repo,
		InventoryService: inventory,
		CostWriter:       s.deps.CostWriter,
		ProfitWriter:     s.deps.ProfitWriter,
	}
	if err := fn(deps); err != nil {
		return err
	}
	s.committed = true
	if baseRepo != nil {
		baseRepo.order = repo.order
		baseRepo.updated = repo.updated
	}
	if baseInventory, ok := s.deps.InventoryService.(*trackingSalesInventoryService); ok {
		baseInventory.params = inventory.params
	}
	return nil
}

func TestTransition_DraftToConfirmed(t *testing.T) {
	repo := &stubSalesOrderRepo{order: &domain.SalesOrder{ID: 1, OrderStatus: domain.SalesOrderStatusDraft}}
	uc := NewSalesOrderUsecase(repo)
	auditLogger := &stubSalesAuditLogger{}
	uc.BindAuditLogger(auditLogger)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	if err := uc.Confirm(c, 1, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.updated == nil || repo.updated.OrderStatus != domain.SalesOrderStatusConfirmed {
		t.Fatalf("expected status CONFIRMED, got %+v", repo.updated)
	}
	if len(auditLogger.payloads) != 1 || auditLogger.payloads[0].Action != "CONFIRM" {
		t.Fatalf("expected confirm audit log, got %+v", auditLogger.payloads)
	}
}

func TestCancelRejectsShipped(t *testing.T) {
	repo := &stubSalesOrderRepo{order: &domain.SalesOrder{ID: 1, OrderStatus: domain.SalesOrderStatusShipped}}
	uc := NewSalesOrderUsecase(repo)

	err := uc.Cancel(nil, 1, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAllocateSupportsPartial(t *testing.T) {
	repo := &stubSalesOrderRepo{order: &domain.SalesOrder{
		ID:          1,
		OrderStatus: domain.SalesOrderStatusConfirmed,
		Items:       []domain.SalesOrderItem{{ID: 10, QtyOrdered: 10, QtyAllocated: 0}},
	}}
	uc := NewSalesOrderUsecase(repo)

	err := uc.Allocate(nil, 1, domain.AllocateParams{Lines: []domain.AllocateLine{{ItemID: 10, QtyAllocated: 6}}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updated == nil || repo.updated.Items[0].QtyAllocated != 6 {
		t.Fatalf("expected qty_allocated=6, got %+v", repo.updated)
	}
	if repo.updated.OrderStatus != domain.SalesOrderStatusAllocated {
		t.Fatalf("expected ALLOCATED status, got %s", repo.updated.OrderStatus)
	}
}

func TestAllocateRollsBackWhenOrderUpdateFails(t *testing.T) {
	repo := &stubSalesOrderRepo{order: &domain.SalesOrder{
		ID:          21,
		OrderStatus: domain.SalesOrderStatusConfirmed,
		OrderNo:     "SO-ALLOC-TX-21",
		StockPool:   domain.StockPoolAvailable,
		Items: []domain.SalesOrderItem{
			{ID: 31, ProductID: 1, QtyOrdered: 1, QtyAllocated: 0},
		},
	}, updateErr: errors.New("update failed")}
	inventorySvc := &trackingSalesInventoryService{}
	uc := NewSalesOrderUsecase(repo)
	uc.BindInventoryService(inventorySvc)
	uc.BindAllocateTransactionManager(&salesAllocateTxManagerStub{
		deps: SalesAllocateTransactionalDeps{
			Repo:             repo,
			InventoryService: inventorySvc,
		},
	})

	err := uc.Allocate(nil, 21, domain.AllocateParams{
		WarehouseID: 1,
		Lines:       []domain.AllocateLine{{ItemID: 31, QtyAllocated: 1}},
	}, nil)
	if err == nil || err.Error() != "update failed" {
		t.Fatalf("expected update failed error, got %v", err)
	}
	if len(inventorySvc.params) != 0 {
		t.Fatalf("expected inventory allocate to rollback, got %d params", len(inventorySvc.params))
	}
	if repo.order.Items[0].QtyAllocated != 0 {
		t.Fatalf("expected base order qty_allocated unchanged, got %d", repo.order.Items[0].QtyAllocated)
	}
	if repo.updated != nil {
		t.Fatalf("expected no committed repo update, got %+v", repo.updated)
	}
}

func TestReturnRejectsExceedShipped(t *testing.T) {
	repo := &stubSalesOrderRepo{order: &domain.SalesOrder{
		ID:          1,
		OrderStatus: domain.SalesOrderStatusDelivered,
		Items:       []domain.SalesOrderItem{{ID: 10, QtyShipped: 5, QtyReturned: 1}},
	}}
	uc := NewSalesOrderUsecase(repo)

	err := uc.Return(nil, 1, domain.ReturnParams{Lines: []domain.ReturnLine{{ItemID: 10, QtyReturned: 5}}}, nil)
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestShipRecordsFIFOOrderCostDetails(t *testing.T) {
	marketplace := "US"
	order := &domain.SalesOrder{
		ID:          8,
		OrderNo:     "SO-0008",
		OrderStatus: domain.SalesOrderStatusAllocated,
		Currency:    "USD",
		Marketplace: &marketplace,
		Items: []domain.SalesOrderItem{
			{ID: 18, ProductID: 201, QtyAllocated: 10, QtyShipped: 0, UnitPrice: 12.8},
		},
	}
	repo := &stubSalesOrderRepo{order: order}
	now := time.Now()
	inventorySvc := &stubSalesInventoryService{
		movement: &inventoryDomain.InventoryMovement{
			OperatedAt: now,
			LotAllocations: []inventoryDomain.InventoryLotAllocation{
				{InventoryLotID: 101, Qty: 6, UnitCost: 3.2},
				{InventoryLotID: 102, Qty: 4, UnitCost: 3.5},
			},
		},
	}
	costRecorder := &stubSalesShipCostRecorder{}
	profitRecorder := &stubSalesShipProfitRecorder{}

	uc := NewSalesOrderUsecase(repo)
	uc.BindInventoryService(inventorySvc)
	uc.BindShipCostDetailWriter(costRecorder)
	uc.BindShipProfitWriter(profitRecorder)

	err := uc.Ship(nil, 8, domain.ShipParams{
		WarehouseID: 1,
		Lines:       []domain.ShipLine{{ItemID: 18, QtyShipped: 10}},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected ship error: %v", err)
	}
	if len(costRecorder.records) != 1 {
		t.Fatalf("expected 1 cost record call, got %d", len(costRecorder.records))
	}
	record := costRecorder.records[0]
	if len(record.Allocations) != 2 {
		t.Fatalf("expected 2 lot allocations, got %d", len(record.Allocations))
	}
	if record.Allocations[0].InventoryLotID != 101 || record.Allocations[0].Qty != 6 {
		t.Fatalf("unexpected first allocation: %+v", record.Allocations[0])
	}
	if record.Allocations[1].InventoryLotID != 102 || record.Allocations[1].Qty != 4 {
		t.Fatalf("unexpected second allocation: %+v", record.Allocations[1])
	}
	if len(profitRecorder.records) != 1 {
		t.Fatalf("expected 1 profit record call, got %d", len(profitRecorder.records))
	}
	if profitRecorder.records[0].IncomeAmount != 128 {
		t.Fatalf("expected income amount 128, got %v", profitRecorder.records[0].IncomeAmount)
	}
	if profitRecorder.records[0].COGSAmount != round6(6*3.2+4*3.5) {
		t.Fatalf("unexpected cogs amount: %v", profitRecorder.records[0].COGSAmount)
	}
}

func TestShipRollsBackWhenCostWriterFails(t *testing.T) {
	marketplace := "US"
	order := &domain.SalesOrder{
		ID:          18,
		OrderNo:     "SO-0018",
		OrderStatus: domain.SalesOrderStatusAllocated,
		Currency:    "USD",
		Marketplace: &marketplace,
		Items: []domain.SalesOrderItem{
			{ID: 28, ProductID: 301, QtyAllocated: 2, QtyShipped: 0, UnitPrice: 12.8},
		},
	}
	repo := &stubSalesOrderRepo{order: order}
	inventorySvc := &trackingSalesInventoryService{}
	costRecorder := &stubSalesShipCostRecorder{err: context.DeadlineExceeded}
	profitRecorder := &stubSalesShipProfitRecorder{}
	txManager := &salesShipTxManagerStub{
		deps: SalesShipTransactionalDeps{
			Repo:             repo,
			InventoryService: inventorySvc,
			CostWriter:       costRecorder,
			ProfitWriter:     profitRecorder,
		},
	}

	uc := NewSalesOrderUsecase(repo)
	uc.BindInventoryService(inventorySvc)
	uc.BindShipCostDetailWriter(costRecorder)
	uc.BindShipProfitWriter(profitRecorder)
	uc.BindShipTransactionManager(txManager)

	err := uc.Ship(nil, 18, domain.ShipParams{
		WarehouseID: 1,
		Lines:       []domain.ShipLine{{ItemID: 28, QtyShipped: 2}},
	}, nil)
	if err == nil {
		t.Fatalf("expected ship error")
	}
	if !txManager.called {
		t.Fatalf("expected transaction manager to be called")
	}
	if txManager.committed {
		t.Fatalf("expected transaction not to commit")
	}
	if len(inventorySvc.params) != 0 {
		t.Fatalf("expected no inventory movement committed, got %+v", inventorySvc.params)
	}
	if repo.updated != nil {
		t.Fatalf("expected order not persisted, got %+v", repo.updated)
	}
	if repo.order.Items[0].QtyShipped != 0 {
		t.Fatalf("expected in-memory order untouched, got %+v", repo.order.Items[0])
	}
}

func TestParseImportRowsUsesBusinessConfigDefaults(t *testing.T) {
	csv := []byte("order_date,order_no,line_no,seller_sku,qty,marketplace,currency,sales_channel\n2026-03-08,SO-1001,1,SKU-001,2,,,\n")
	uc := NewSalesOrderUsecase(&stubSalesOrderRepo{})
	uc.BindConfigProvider(stubSalesConfigProvider{})

	lines, rowErrors, totalRows, err := parseImportRows(csv, uc.getImportDefaults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if totalRows != 1 || len(rowErrors) != 0 || len(lines) != 1 {
		t.Fatalf("unexpected import parse result: total=%d errors=%d lines=%d", totalRows, len(rowErrors), len(lines))
	}
	if lines[0].Marketplace != "DE" {
		t.Fatalf("expected default marketplace DE, got %s", lines[0].Marketplace)
	}
	if lines[0].Currency != "EUR" {
		t.Fatalf("expected default currency EUR, got %s", lines[0].Currency)
	}
	if lines[0].SalesChannel == nil || *lines[0].SalesChannel != "MANUAL-IMPORT" {
		t.Fatalf("expected default sales channel, got %+v", lines[0].SalesChannel)
	}
}

func TestParseImportRowsAllowsMissingMarketplaceColumnWhenDefaultConfigured(t *testing.T) {
	csv := []byte("order_date,order_no,line_no,seller_sku,qty,currency,sales_channel\n2026-03-08,SO-1002,1,SKU-001,2,,\n")
	uc := NewSalesOrderUsecase(&stubSalesOrderRepo{})
	uc.BindConfigProvider(stubSalesConfigProvider{})

	lines, rowErrors, totalRows, err := parseImportRows(csv, uc.getImportDefaults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if totalRows != 1 || len(rowErrors) != 0 || len(lines) != 1 {
		t.Fatalf("unexpected import parse result: total=%d errors=%d lines=%d", totalRows, len(rowErrors), len(lines))
	}
	if lines[0].Marketplace != "DE" {
		t.Fatalf("expected default marketplace DE, got %s", lines[0].Marketplace)
	}
}

func TestParseImportRowsMarksAmazonImportAsSellableBusinessSource(t *testing.T) {
	csv := []byte("order_date,order_no,external_order_no,line_no,seller_sku,qty,marketplace,currency,unit_price,sales_channel\n2026-03-10,AMZ-UAT-0001,AMZ-UAT-0001,1,SKU-001,1,US,USD,19.99,Amazon\n")

	lines, rowErrors, totalRows, err := parseImportRows(csv, importDefaults{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if totalRows != 1 || len(rowErrors) != 0 || len(lines) != 1 {
		t.Fatalf("unexpected import parse result: total=%d errors=%d lines=%d", totalRows, len(rowErrors), len(lines))
	}
	if lines[0].SourceType != "AMAZON_IMPORT" {
		t.Fatalf("expected AMAZON_IMPORT, got %s", lines[0].SourceType)
	}
}

func TestReturnRecordsReverseCostAndProfit(t *testing.T) {
	marketplace := "US"
	order := &domain.SalesOrder{
		ID:          9,
		OrderNo:     "SO-0009",
		OrderStatus: domain.SalesOrderStatusDelivered,
		Currency:    "USD",
		Marketplace: &marketplace,
		Items: []domain.SalesOrderItem{
			{ID: 19, ProductID: 301, QtyShipped: 5, QtyReturned: 0, UnitPrice: 10.5},
		},
	}
	repo := &stubSalesOrderRepo{order: order}
	now := time.Now()
	inventorySvc := &stubSalesInventoryService{
		movement: &inventoryDomain.InventoryMovement{OperatedAt: now},
	}
	returnUnitCost := 7.3
	costRecorder := &stubSalesShipCostRecorder{cogs: 14.6, returnUnitCost: &returnUnitCost}
	profitRecorder := &stubSalesShipProfitRecorder{}

	uc := NewSalesOrderUsecase(repo)
	uc.BindInventoryService(inventorySvc)
	uc.BindShipCostDetailWriter(costRecorder)
	uc.BindShipProfitWriter(profitRecorder)

	err := uc.Return(nil, 9, domain.ReturnParams{
		WarehouseID: 1,
		Lines:       []domain.ReturnLine{{ItemID: 19, QtyReturned: 2}},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected return error: %v", err)
	}
	if len(costRecorder.resolves) != 1 {
		t.Fatalf("expected 1 return unit-cost resolve call, got %d", len(costRecorder.resolves))
	}
	if len(inventorySvc.params) != 1 || inventorySvc.params[0].UnitCost == nil || *inventorySvc.params[0].UnitCost != returnUnitCost {
		t.Fatalf("expected return movement unit cost %v, got %+v", returnUnitCost, inventorySvc.params)
	}
	if len(costRecorder.returns) != 1 {
		t.Fatalf("expected 1 return cost record call, got %d", len(costRecorder.returns))
	}
	if costRecorder.returns[0].QtyReturned != 2 {
		t.Fatalf("expected returned qty 2, got %d", costRecorder.returns[0].QtyReturned)
	}
	if len(profitRecorder.returns) != 1 {
		t.Fatalf("expected 1 return profit record, got %d", len(profitRecorder.returns))
	}
	if profitRecorder.returns[0].IncomeAmount != 21 {
		t.Fatalf("expected return income amount 21, got %v", profitRecorder.returns[0].IncomeAmount)
	}
	if profitRecorder.returns[0].COGSAmount != 14.6 {
		t.Fatalf("expected return cogs amount 14.6, got %v", profitRecorder.returns[0].COGSAmount)
	}
	if repo.updated == nil || repo.updated.Items[0].QtyReturned != 2 {
		t.Fatalf("expected qty_returned=2, got %+v", repo.updated)
	}
}

func TestReturnRollsBackWhenFinanceRecorderFails(t *testing.T) {
	marketplace := "US"
	order := &domain.SalesOrder{
		ID:          29,
		OrderNo:     "SO-0029",
		OrderStatus: domain.SalesOrderStatusDelivered,
		Currency:    "USD",
		Marketplace: &marketplace,
		Items: []domain.SalesOrderItem{
			{ID: 39, ProductID: 401, QtyOrdered: 2, QtyAllocated: 2, QtyShipped: 2, QtyReturned: 0, UnitPrice: 15},
		},
	}
	repo := &stubSalesOrderRepo{order: order}
	inventorySvc := &trackingSalesInventoryService{}
	costRecorder := &stubSalesShipCostRecorder{err: context.Canceled}
	profitRecorder := &stubSalesShipProfitRecorder{}
	txManager := &salesReturnTxManagerStub{
		deps: SalesReturnTransactionalDeps{
			Repo:             repo,
			InventoryService: inventorySvc,
			CostWriter:       costRecorder,
			ProfitWriter:     profitRecorder,
		},
	}

	uc := NewSalesOrderUsecase(repo)
	uc.BindInventoryService(inventorySvc)
	uc.BindShipCostDetailWriter(costRecorder)
	uc.BindShipProfitWriter(profitRecorder)
	uc.BindReturnTransactionManager(txManager)

	err := uc.Return(nil, 29, domain.ReturnParams{
		WarehouseID: 1,
		Lines:       []domain.ReturnLine{{ItemID: 39, QtyReturned: 1}},
	}, nil)
	if err == nil {
		t.Fatalf("expected return error")
	}
	if !txManager.called {
		t.Fatalf("expected return transaction manager to be called")
	}
	if txManager.committed {
		t.Fatalf("expected return transaction not to commit")
	}
	if len(inventorySvc.params) != 0 {
		t.Fatalf("expected no return inventory movement committed, got %+v", inventorySvc.params)
	}
	if repo.updated != nil {
		t.Fatalf("expected order not persisted, got %+v", repo.updated)
	}
	if repo.order.Items[0].QtyReturned != 0 {
		t.Fatalf("expected in-memory return qty untouched, got %+v", repo.order.Items[0])
	}
}

func TestRecordShipProfitFallsBackToConfiguredCurrency(t *testing.T) {
	marketplace := "US"
	order := &domain.SalesOrder{
		ID:          10,
		OrderNo:     "SO-0010",
		OrderStatus: domain.SalesOrderStatusAllocated,
		Marketplace: &marketplace,
		Items: []domain.SalesOrderItem{
			{ID: 20, ProductID: 501, QtyAllocated: 1, UnitPrice: 9.9},
		},
	}
	movement := &inventoryDomain.InventoryMovement{
		OperatedAt: time.Now(),
		LotAllocations: []inventoryDomain.InventoryLotAllocation{
			{InventoryLotID: 201, Qty: 1, UnitCost: 2.2},
		},
	}
	profitRecorder := &stubSalesShipProfitRecorder{}

	uc := NewSalesOrderUsecase(&stubSalesOrderRepo{})
	uc.BindConfigProvider(stubSalesConfigProvider{})
	uc.BindShipProfitWriter(profitRecorder)

	if err := uc.recordShipProfit(order, &order.Items[0], 1, movement, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profitRecorder.records) != 1 {
		t.Fatalf("expected one profit record, got %d", len(profitRecorder.records))
	}
	if profitRecorder.records[0].IncomeCurrency != "EUR" {
		t.Fatalf("expected fallback currency EUR, got %s", profitRecorder.records[0].IncomeCurrency)
	}
}

func TestRecordReturnFinanceFallsBackToConfiguredCurrency(t *testing.T) {
	marketplace := "US"
	order := &domain.SalesOrder{
		ID:          11,
		OrderNo:     "SO-0011",
		OrderStatus: domain.SalesOrderStatusDelivered,
		Marketplace: &marketplace,
		Items: []domain.SalesOrderItem{
			{ID: 21, ProductID: 601, QtyShipped: 2, UnitPrice: 10},
		},
	}
	movement := &inventoryDomain.InventoryMovement{OperatedAt: time.Now()}
	costRecorder := &stubSalesShipCostRecorder{}
	profitRecorder := &stubSalesShipProfitRecorder{}

	uc := NewSalesOrderUsecase(&stubSalesOrderRepo{})
	uc.BindConfigProvider(stubSalesConfigProvider{})
	uc.BindShipCostDetailWriter(costRecorder)
	uc.BindShipProfitWriter(profitRecorder)

	if err := uc.recordReturnFinance(order, &order.Items[0], 1, 1, movement, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(costRecorder.returns) != 1 {
		t.Fatalf("expected one return cost record, got %d", len(costRecorder.returns))
	}
	if costRecorder.returns[0].Currency != "EUR" {
		t.Fatalf("expected fallback currency EUR, got %s", costRecorder.returns[0].Currency)
	}
	if len(profitRecorder.returns) != 1 {
		t.Fatalf("expected one return profit record, got %d", len(profitRecorder.returns))
	}
	if profitRecorder.returns[0].IncomeCurrency != "EUR" {
		t.Fatalf("expected fallback currency EUR, got %s", profitRecorder.returns[0].IncomeCurrency)
	}
}

func TestAllocateAmazonOrderUsesSellableStockPool(t *testing.T) {
	repo := &stubSalesOrderRepo{order: &domain.SalesOrder{
		ID:          12,
		SourceType:  "AMAZON_API",
		StockPool:   domain.StockPoolSellable,
		OrderStatus: domain.SalesOrderStatusConfirmed,
		Items:       []domain.SalesOrderItem{{ID: 22, ProductID: 701, QtyOrdered: 3, QtyAllocated: 0}},
	}}
	inventorySvc := &stubSalesInventoryService{}
	uc := NewSalesOrderUsecase(repo)
	uc.BindInventoryService(inventorySvc)

	err := uc.Allocate(nil, 12, domain.AllocateParams{
		WarehouseID: 9,
		Lines:       []domain.AllocateLine{{ItemID: 22, QtyAllocated: 2}},
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inventorySvc.params) != 1 || inventorySvc.params[0].StockPool == nil {
		t.Fatalf("expected one inventory call with stock pool, got %+v", inventorySvc.params)
	}
	if *inventorySvc.params[0].StockPool != inventoryDomain.StockPoolSellable {
		t.Fatalf("expected sellable stock pool, got %v", *inventorySvc.params[0].StockPool)
	}
}
