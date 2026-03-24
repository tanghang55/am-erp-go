package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/numbering"
	"am-erp-go/internal/module/inventory/domain"
	packagingDomain "am-erp-go/internal/module/packaging/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrInsufficientStock                        = errors.New("insufficient stock")
	ErrInvalidQuantity                          = errors.New("invalid quantity")
	ErrPlatformReceiveRequiresShipmentReference = errors.New("platform receive requires shipment reference")
	ErrAssemblyTransactionManagerNotConfigured  = errors.New("assembly transaction manager not configured")
	ErrPackingRequirementResolverNotConfigured  = errors.New("packing requirement resolver not configured")
	ErrPackingMaterialsNotConfigured            = errors.New("packing materials not configured")
	ErrPackingMaterialQuantityInvalid           = errors.New("packing material quantity must resolve to integer")
	ErrPackingOperatorRequired                  = errors.New("packing operator is required")
)

type InventoryUsecase struct {
	balanceRepo                     domain.InventoryBalanceRepository
	movementRepo                    domain.InventoryMovementRepository
	lotRepo                         domain.InventoryLotRepository
	platformReceiveUnitCostResolver PlatformReceiveUnitCostResolver
	platformReceiveRecorder         PlatformReceiveRecorder
	platformReceiveTxManager        PlatformReceiveTransactionManager
	seedLotUnitCostResolver         SeedLotUnitCostResolver
	shipmentLotUnitCostResolver     ShipmentLotUnitCostResolver
	packingRequirementResolver      PackingRequirementResolver
	packingCostRecorder             PackingCostRecorder
	assemblyTxManager               AssemblyTransactionManager
	auditLogger                     InventoryAuditLogger
}

type PlatformReceiveUnitCostResolver interface {
	ResolvePlatformReceiveUnitCost(ctx context.Context, params *domain.CreateMovementParams) (*float64, error)
}

type PlatformReceiveRecorder interface {
	ValidatePlatformReceive(ctx context.Context, params *domain.CreateMovementParams) error
	RecordPlatformReceive(ctx context.Context, params *domain.CreateMovementParams) error
}

type SeedLotUnitCostResolver func(ctx context.Context, productID, warehouseID uint64) (*float64, error)
type ShipmentLotUnitCostResolver func(ctx context.Context, productID uint64, referenceType *string, referenceID *uint64, referenceNumber *string, occurredAt time.Time) (*float64, error)

type PackingRequirement struct {
	PackagingItemID uint64
	QuantityPerUnit float64
	ItemCode        string
	ItemName        string
	Unit            string
}

type PackingMaterialCostLine struct {
	PackagingItemID uint64
	Quantity        uint64
	UnitCost        float64
	Currency        string
	ItemCode        string
	ItemName        string
}

type PackingMaterialCostRecordParams struct {
	InventoryMovementID uint64
	ProductID           uint64
	WarehouseID         uint64
	Quantity            uint64
	ReferenceNumber     string
	OccurredAt          time.Time
	OperatorID          *uint64
	Lines               []PackingMaterialCostLine
}

type PackingRequirementResolver interface {
	ResolvePackingRequirements(ctx context.Context, productID uint64) ([]PackingRequirement, error)
}

type PackingCostRecorder interface {
	RecordPackingMaterialCost(params *PackingMaterialCostRecordParams) error
}

type AssemblyTransactionalDeps struct {
	BalanceRepo         domain.InventoryBalanceRepository
	MovementRepo        domain.InventoryMovementRepository
	LotRepo             domain.InventoryLotRepository
	PackagingItemRepo   packagingDomain.PackagingItemRepository
	PackagingLedgerRepo packagingDomain.PackagingLedgerRepository
}

type AssemblyTransactionManager interface {
	Run(ctx context.Context, fn func(AssemblyTransactionalDeps) (*domain.InventoryMovement, error)) (*domain.InventoryMovement, error)
}

type InventoryAuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

type PlatformReceiveTransactionalDeps struct {
	BalanceRepo  domain.InventoryBalanceRepository
	MovementRepo domain.InventoryMovementRepository
	LotRepo      domain.InventoryLotRepository
	Recorder     PlatformReceiveRecorder
}

type PlatformReceiveTransactionManager interface {
	Run(ctx context.Context, fn func(PlatformReceiveTransactionalDeps) (*domain.InventoryMovement, error)) (*domain.InventoryMovement, error)
}

func NewInventoryUsecase(
	balanceRepo domain.InventoryBalanceRepository,
	movementRepo domain.InventoryMovementRepository,
	lotRepo domain.InventoryLotRepository,
) *InventoryUsecase {
	return &InventoryUsecase{
		balanceRepo:  balanceRepo,
		movementRepo: movementRepo,
		lotRepo:      lotRepo,
	}
}

func (u *InventoryUsecase) BindPlatformReceiveUnitCostResolver(resolver PlatformReceiveUnitCostResolver) {
	u.platformReceiveUnitCostResolver = resolver
}

func (u *InventoryUsecase) BindPlatformReceiveRecorder(recorder PlatformReceiveRecorder) {
	u.platformReceiveRecorder = recorder
}

func (u *InventoryUsecase) BindSeedLotUnitCostResolver(resolver SeedLotUnitCostResolver) {
	u.seedLotUnitCostResolver = resolver
}

func (u *InventoryUsecase) BindShipmentLotUnitCostResolver(resolver ShipmentLotUnitCostResolver) {
	u.shipmentLotUnitCostResolver = resolver
}

func (u *InventoryUsecase) BindPlatformReceiveTransactionManager(manager PlatformReceiveTransactionManager) {
	u.platformReceiveTxManager = manager
}

func (u *InventoryUsecase) BindPackingRequirementResolver(resolver PackingRequirementResolver) {
	u.packingRequirementResolver = resolver
}

func (u *InventoryUsecase) BindPackingCostRecorder(recorder PackingCostRecorder) {
	u.packingCostRecorder = recorder
}

func (u *InventoryUsecase) BindAssemblyTransactionManager(manager AssemblyTransactionManager) {
	u.assemblyTxManager = manager
}

func (u *InventoryUsecase) BindAuditLogger(logger InventoryAuditLogger) {
	u.auditLogger = logger
}

// ListMovements 查询库存流水列表
func (u *InventoryUsecase) ListMovements(params *domain.MovementListParams) ([]*domain.InventoryMovement, int64, error) {
	return u.movementRepo.List(params)
}

// GetMovement 获取单条流水详情
func (u *InventoryUsecase) GetMovement(id uint64) (*domain.InventoryMovement, error) {
	return u.movementRepo.GetByID(id)
}

// ListLots 查询库存批次列表（FIFO按最老批次优先）
func (u *InventoryUsecase) ListLots(params *domain.InventoryLotListParams) ([]*domain.InventoryLot, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	lots, total, err := u.lotRepo.List(params)
	if err != nil {
		return nil, 0, err
	}
	if err := u.backfillShipmentLotUnitCosts(context.Background(), lots); err != nil {
		return nil, 0, err
	}
	return lots, total, nil
}

// CreateMovement 创建库存流水（通用方法）
func (u *InventoryUsecase) CreateMovement(ctx context.Context, params *domain.CreateMovementParams) (*domain.InventoryMovement, error) {
	if params.Quantity == 0 {
		return nil, ErrInvalidQuantity
	}
	if err := validateMovementParams(params); err != nil {
		return nil, err
	}
	if params.MovementType == domain.MovementTypeAssemblyComplete {
		if u.packingRequirementResolver == nil {
			return nil, ErrPackingRequirementResolverNotConfigured
		}
		if u.assemblyTxManager == nil {
			return nil, ErrAssemblyTransactionManagerNotConfigured
		}
		movement, err := u.assemblyTxManager.Run(ctx, func(deps AssemblyTransactionalDeps) (*domain.InventoryMovement, error) {
			return u.createPackingMovementWithDeps(
				ctx,
				params,
				deps.BalanceRepo,
				deps.MovementRepo,
				deps.LotRepo,
				deps.PackagingItemRepo,
				deps.PackagingLedgerRepo,
			)
		})
		if err != nil {
			return nil, err
		}
		u.recordPackingAudit(ctx, movement, params)
		return movement, nil
	}
	if params.MovementType == domain.MovementTypePlatformReceive && params.UnitCost == nil && u.platformReceiveUnitCostResolver != nil {
		unitCost, err := u.platformReceiveUnitCostResolver.ResolvePlatformReceiveUnitCost(ctx, params)
		if err != nil {
			return nil, err
		}
		if unitCost != nil {
			params.UnitCost = unitCost
		}
	}
	if params.MovementType == domain.MovementTypePlatformReceive && u.platformReceiveTxManager != nil {
		return u.platformReceiveTxManager.Run(ctx, func(deps PlatformReceiveTransactionalDeps) (*domain.InventoryMovement, error) {
			recorder := deps.Recorder
			if recorder == nil {
				recorder = u.platformReceiveRecorder
			}
			return u.createMovementWithDeps(ctx, params, deps.BalanceRepo, deps.MovementRepo, deps.LotRepo, recorder)
		})
	}
	return u.createMovementWithDeps(ctx, params, u.balanceRepo, u.movementRepo, u.lotRepo, u.platformReceiveRecorder)
}

func (u *InventoryUsecase) createPackingMovementWithDeps(
	ctx context.Context,
	params *domain.CreateMovementParams,
	balanceRepo domain.InventoryBalanceRepository,
	movementRepo domain.InventoryMovementRepository,
	lotRepo domain.InventoryLotRepository,
	packagingItemRepo packagingDomain.PackagingItemRepository,
	packagingLedgerRepo packagingDomain.PackagingLedgerRepository,
) (*domain.InventoryMovement, error) {
	requirements, err := u.packingRequirementResolver.ResolvePackingRequirements(ctx, params.ProductID)
	if err != nil {
		return nil, err
	}
	if len(requirements) == 0 {
		return nil, ErrPackingMaterialsNotConfigured
	}
	if packagingItemRepo == nil || packagingLedgerRepo == nil {
		return nil, ErrAssemblyTransactionManagerNotConfigured
	}
	operatorID := uint64(0)
	if params.OperatorID != nil {
		operatorID = *params.OperatorID
	}
	if operatorID == 0 {
		return nil, ErrPackingOperatorRequired
	}

	traceID := uuid.New().String()
	referenceType := "PACKING"
	if params.ReferenceType != nil && strings.TrimSpace(*params.ReferenceType) != "" {
		referenceType = strings.TrimSpace(*params.ReferenceType)
	}
	referenceNumber := buildPackingReferenceNumber()
	if params.ReferenceNumber != nil && strings.TrimSpace(*params.ReferenceNumber) != "" {
		referenceNumber = strings.TrimSpace(*params.ReferenceNumber)
	}
	operatedAt := time.Now()
	if params.OperatedAt != nil {
		operatedAt = *params.OperatedAt
	}

	execUC := *u
	execUC.balanceRepo = balanceRepo
	execUC.movementRepo = movementRepo
	execUC.lotRepo = lotRepo

	consumeParams := &domain.CreateMovementParams{
		TraceID:         &traceID,
		ProductID:       params.ProductID,
		WarehouseID:     params.WarehouseID,
		MovementType:    domain.MovementTypeAssemblyConsume,
		Quantity:        params.Quantity,
		ReferenceType:   &referenceType,
		ReferenceID:     params.ReferenceID,
		ReferenceNumber: &referenceNumber,
		Remark:          params.Remark,
		OperatorID:      params.OperatorID,
		OperatedAt:      &operatedAt,
	}
	if _, err := execUC.createMovementWithDeps(ctx, consumeParams, balanceRepo, movementRepo, lotRepo, nil); err != nil {
		return nil, err
	}

	mainParams := &domain.CreateMovementParams{
		TraceID:         &traceID,
		ProductID:       params.ProductID,
		WarehouseID:     params.WarehouseID,
		MovementType:    domain.MovementTypeAssemblyComplete,
		Quantity:        params.Quantity,
		ReferenceType:   &referenceType,
		ReferenceID:     params.ReferenceID,
		ReferenceNumber: &referenceNumber,
		UnitCost:        params.UnitCost,
		Remark:          params.Remark,
		OperatorID:      params.OperatorID,
		OperatedAt:      &operatedAt,
	}

	movement, err := execUC.createMovementWithDeps(ctx, mainParams, balanceRepo, movementRepo, lotRepo, nil)
	if err != nil {
		return nil, err
	}
	costLines, err := u.createPackingPackagingLedgers(ctx, movement, requirements, params.Quantity, operatorID, operatedAt, packagingItemRepo, packagingLedgerRepo)
	if err != nil {
		return nil, err
	}
	if err := u.recordPackingMaterialCost(movement, uint64(params.Quantity), costLines, params.OperatorID, operatedAt); err != nil {
		return nil, err
	}
	return movement, nil
}

func (u *InventoryUsecase) recordPackingAudit(ctx context.Context, movement *domain.InventoryMovement, params *domain.CreateMovementParams) {
	if u.auditLogger == nil || movement == nil || movement.ID == 0 {
		return
	}
	ginCtx, ok := ctx.(*gin.Context)
	if !ok || ginCtx == nil {
		return
	}
	after := map[string]any{
		"movement_id":      movement.ID,
		"trace_id":         movement.TraceID,
		"product_id":       movement.ProductID,
		"warehouse_id":     movement.WarehouseID,
		"quantity":         movement.Quantity,
		"reference_type":   movement.ReferenceType,
		"reference_number": movement.ReferenceNumber,
	}
	if params != nil && params.Remark != nil {
		after["remark"] = *params.Remark
	}
	_ = u.auditLogger.RecordFromContext(ginCtx, systemUsecase.AuditLogPayload{
		Module:     "Inventory",
		Action:     "PACK_COMPLETE",
		EntityType: "InventoryMovement",
		EntityID:   fmt.Sprintf("%d", movement.ID),
		After:      after,
	})
}

func (u *InventoryUsecase) recordPackingMaterialCost(
	movement *domain.InventoryMovement,
	productQty uint64,
	lines []PackingMaterialCostLine,
	operatorID *uint64,
	occurredAt time.Time,
) error {
	if u.packingCostRecorder == nil || movement == nil || movement.ID == 0 || productQty == 0 || len(lines) == 0 {
		return nil
	}
	return u.packingCostRecorder.RecordPackingMaterialCost(&PackingMaterialCostRecordParams{
		InventoryMovementID: movement.ID,
		ProductID:           movement.ProductID,
		WarehouseID:         movement.WarehouseID,
		Quantity:            productQty,
		ReferenceNumber:     pointerStringValue(movement.ReferenceNumber),
		OccurredAt:          occurredAt,
		OperatorID:          operatorID,
		Lines:               lines,
	})
}

func (u *InventoryUsecase) createMovementWithDeps(
	ctx context.Context,
	params *domain.CreateMovementParams,
	balanceRepo domain.InventoryBalanceRepository,
	movementRepo domain.InventoryMovementRepository,
	lotRepo domain.InventoryLotRepository,
	recorder PlatformReceiveRecorder,
) (*domain.InventoryMovement, error) {
	execUC := *u
	execUC.balanceRepo = balanceRepo
	execUC.movementRepo = movementRepo
	execUC.lotRepo = lotRepo
	if params.MovementType == domain.MovementTypePlatformReceive && recorder != nil {
		if err := recorder.ValidatePlatformReceive(ctx, params); err != nil {
			return nil, err
		}
	}

	// Get or create balance record
	balance, err := balanceRepo.GetOrCreate(ctx, params.ProductID, params.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	balanceBefore := *balance

	// Calculate new quantities based on movement type
	newBalance, err := u.calculateNewBalance(balance, params)
	if err != nil {
		return nil, err
	}

	lotAllocations := make([]domain.InventoryLotAllocation, 0)
	if lotRepo != nil {
		if err := execUC.ensureLotSeed(ctx, &balanceBefore); err != nil {
			return nil, fmt.Errorf("failed to seed inventory lots: %w", err)
		}
		allocations, err := execUC.applyLotMovement(ctx, params, &balanceBefore)
		if err != nil {
			return nil, err
		}
		lotAllocations = allocations
	}

	// Create movement record
	operatedAt := time.Now()
	if params.OperatedAt != nil {
		operatedAt = *params.OperatedAt
	}

	traceID := uuid.New().String()
	if params.TraceID != nil && strings.TrimSpace(*params.TraceID) != "" {
		traceID = strings.TrimSpace(*params.TraceID)
	}
	movement := &domain.InventoryMovement{
		TraceID:                       &traceID,
		ProductID:                     params.ProductID,
		WarehouseID:                   params.WarehouseID,
		MovementType:                  params.MovementType,
		ReferenceType:                 params.ReferenceType,
		ReferenceID:                   params.ReferenceID,
		ReferenceNumber:               params.ReferenceNumber,
		StockPool:                     params.StockPool,
		Quantity:                      params.Quantity,
		BeforeAvailable:               balance.AvailableQuantity,
		AfterAvailable:                newBalance.AvailableQuantity,
		BeforeReserved:                balance.ReservedQuantity,
		AfterReserved:                 newBalance.ReservedQuantity,
		BeforePurchasingInTransit:     balance.PurchasingInTransit,
		AfterPurchasingInTransit:      newBalance.PurchasingInTransit,
		BeforePendingInspection:       balance.PendingInspection,
		AfterPendingInspection:        newBalance.PendingInspection,
		BeforeRawMaterial:             balance.RawMaterial,
		AfterRawMaterial:              newBalance.RawMaterial,
		BeforeSellable:                balance.Sellable,
		AfterSellable:                 newBalance.Sellable,
		BeforeSellableReserved:        balance.SellableReserved,
		AfterSellableReserved:         newBalance.SellableReserved,
		BeforePendingShipment:         balance.PendingShipment,
		AfterPendingShipment:          newBalance.PendingShipment,
		BeforePendingShipmentReserved: balance.PendingShipmentReserved,
		AfterPendingShipmentReserved:  newBalance.PendingShipmentReserved,
		BeforeLogisticsInTransit:      balance.LogisticsInTransit,
		AfterLogisticsInTransit:       newBalance.LogisticsInTransit,
		BeforeDamaged:                 balance.DamagedQuantity,
		AfterDamaged:                  newBalance.DamagedQuantity,
		UnitCost:                      params.UnitCost,
		Remark:                        params.Remark,
		OperatorID:                    params.OperatorID,
		OperatedAt:                    operatedAt,
		GmtCreate:                     time.Now(),
		GmtModified:                   time.Now(),
		LotAllocations:                lotAllocations,
	}

	if params.UnitCost != nil {
		totalCost := *params.UnitCost * float64(absInt(params.Quantity))
		movement.TotalCost = &totalCost
	}

	// Save movement
	if err := movementRepo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("failed to create movement: %w", err)
	}

	// Update balance
	balance.AvailableQuantity = newBalance.AvailableQuantity
	balance.ReservedQuantity = newBalance.ReservedQuantity
	balance.SellableReserved = newBalance.SellableReserved
	balance.DamagedQuantity = newBalance.DamagedQuantity
	balance.PurchasingInTransit = newBalance.PurchasingInTransit
	balance.PendingInspection = newBalance.PendingInspection
	balance.RawMaterial = newBalance.RawMaterial
	balance.PendingShipment = newBalance.PendingShipment
	balance.PendingShipmentReserved = newBalance.PendingShipmentReserved
	balance.LogisticsInTransit = newBalance.LogisticsInTransit
	balance.Sellable = newBalance.Sellable
	balance.Returned = newBalance.Returned
	balance.TotalQuantity = calculateBalanceTotal(balance)
	balance.LastMovementAt = &operatedAt
	balance.GmtModified = time.Now()

	if err := balanceRepo.Update(ctx, balance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}
	if params.MovementType == domain.MovementTypePlatformReceive && recorder != nil {
		if err := recorder.RecordPlatformReceive(ctx, params); err != nil {
			return nil, err
		}
	}

	return movement, nil
}

// Transfer 仓库间调拨
func (u *InventoryUsecase) Transfer(ctx context.Context, params *domain.TransferParams) error {
	if params.Quantity == 0 {
		return ErrInvalidQuantity
	}

	traceID := uuid.New().String()

	// Transfer out
	outParams := &domain.CreateMovementParams{
		ProductID:       params.ProductID,
		WarehouseID:     params.FromWarehouseID,
		MovementType:    domain.MovementTypeTransferOut,
		Quantity:        -int(params.Quantity),
		ReferenceType:   params.ReferenceType,
		ReferenceNumber: params.ReferenceNumber,
		UnitCost:        params.UnitCost,
		Remark:          params.Remark,
		OperatorID:      params.OperatorID,
	}

	outMovement, err := u.CreateMovement(ctx, outParams)
	if err != nil {
		return fmt.Errorf("failed to transfer out: %w", err)
	}

	// Transfer in
	inParams := &domain.CreateMovementParams{
		ProductID:       params.ProductID,
		WarehouseID:     params.ToWarehouseID,
		MovementType:    domain.MovementTypeTransferIn,
		Quantity:        int(params.Quantity),
		ReferenceType:   params.ReferenceType,
		ReferenceNumber: params.ReferenceNumber,
		UnitCost:        params.UnitCost,
		Remark:          params.Remark,
		OperatorID:      params.OperatorID,
	}

	if _, err := u.CreateMovement(ctx, inParams); err != nil {
		return fmt.Errorf("failed to transfer in: %w", err)
	}

	// Use same trace ID for both movements
	outMovement.TraceID = &traceID

	return nil
}

// ListBalances 查询库存余额列表
func (u *InventoryUsecase) ListBalances(params *domain.BalanceListParams) ([]*domain.InventoryBalance, int64, error) {
	return u.balanceRepo.List(params)
}

// GetProductBalance 获取产品在指定仓库的库存
func (u *InventoryUsecase) GetProductBalance(productID, warehouseID uint64) (*domain.InventoryBalance, error) {
	return u.balanceRepo.GetByProductAndWarehouse(productID, warehouseID)
}

// calculateNewBalance 根据流水类型计算新的库存数量
func (u *InventoryUsecase) calculateNewBalance(balance *domain.InventoryBalance, params *domain.CreateMovementParams) (*domain.InventoryBalance, error) {
	movementType := params.MovementType
	quantity := params.Quantity
	newBalance := &domain.InventoryBalance{
		AvailableQuantity:       balance.AvailableQuantity,
		ReservedQuantity:        balance.ReservedQuantity,
		SellableReserved:        balance.SellableReserved,
		DamagedQuantity:         balance.DamagedQuantity,
		PurchasingInTransit:     balance.PurchasingInTransit,
		PendingInspection:       balance.PendingInspection,
		RawMaterial:             balance.RawMaterial,
		PendingShipment:         balance.PendingShipment,
		PendingShipmentReserved: balance.PendingShipmentReserved,
		LogisticsInTransit:      balance.LogisticsInTransit,
		Sellable:                balance.Sellable,
		Returned:                balance.Returned,
	}

	switch movementType {
	case domain.MovementTypePurchaseReceipt,
		domain.MovementTypeReturnReceipt,
		domain.MovementTypeTransferIn:
		// Increase available quantity
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		if movementType == domain.MovementTypeReturnReceipt && stockPool(params) == domain.StockPoolSellable {
			newBalance.Sellable = balance.Sellable + uint(quantity)
		} else {
			newBalance.AvailableQuantity = balance.AvailableQuantity + uint(quantity)
		}

	case domain.MovementTypeSalesShipment,
		domain.MovementTypeTransferOut:
		// Decrease available quantity
		if quantity > 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(-quantity)
		if balance.AvailableQuantity < decrease {
			return nil, ErrInsufficientStock
		}
		newBalance.AvailableQuantity = balance.AvailableQuantity - decrease

	case domain.MovementTypeSalesAllocate:
		// 销售锁定: available/sellable -> reserved
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if stockPool(params) == domain.StockPoolSellable {
			if balance.Sellable < decrease {
				return nil, ErrInsufficientStock
			}
			newBalance.Sellable = balance.Sellable - decrease
			newBalance.SellableReserved = balance.SellableReserved + decrease
		} else {
			if balance.AvailableQuantity < decrease {
				return nil, ErrInsufficientStock
			}
			newBalance.AvailableQuantity = balance.AvailableQuantity - decrease
			newBalance.ReservedQuantity = balance.ReservedQuantity + decrease
		}

	case domain.MovementTypeSalesRelease:
		// 销售解锁: reserved -> available/sellable
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if stockPool(params) == domain.StockPoolSellable {
			if balance.SellableReserved < decrease {
				return nil, ErrInsufficientStock
			}
			newBalance.SellableReserved = balance.SellableReserved - decrease
			newBalance.Sellable = balance.Sellable + decrease
		} else {
			if balance.ReservedQuantity < decrease {
				return nil, ErrInsufficientStock
			}
			newBalance.ReservedQuantity = balance.ReservedQuantity - decrease
			newBalance.AvailableQuantity = balance.AvailableQuantity + decrease
		}

	case domain.MovementTypeSalesShip:
		// 销售发货: reserved/sellable_reserved -> 出库扣减
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if stockPool(params) == domain.StockPoolSellable {
			if balance.SellableReserved < decrease {
				return nil, ErrInsufficientStock
			}
			newBalance.SellableReserved = balance.SellableReserved - decrease
		} else {
			if balance.ReservedQuantity < decrease {
				return nil, ErrInsufficientStock
			}
			newBalance.ReservedQuantity = balance.ReservedQuantity - decrease
		}

	case domain.MovementTypeDamageWriteOff:
		// Decrease available, increase damaged
		if quantity > 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(-quantity)
		if balance.AvailableQuantity < decrease {
			return nil, ErrInsufficientStock
		}
		newBalance.AvailableQuantity = balance.AvailableQuantity - decrease
		newBalance.DamagedQuantity = balance.DamagedQuantity + decrease

	case domain.MovementTypeStockTakeAdjustment,
		domain.MovementTypeManualAdjustment:
		// Can be positive or negative
		if quantity > 0 {
			newBalance.AvailableQuantity = balance.AvailableQuantity + uint(quantity)
		} else if quantity < 0 {
			decrease := uint(-quantity)
			if balance.AvailableQuantity < decrease {
				return nil, ErrInsufficientStock
			}
			newBalance.AvailableQuantity = balance.AvailableQuantity - decrease
		}

	// 新增库存状态流转
	case domain.MovementTypePurchaseShip:
		// 供应商发货 → 采购在途
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		newBalance.PurchasingInTransit = balance.PurchasingInTransit + uint(quantity)

	case domain.MovementTypeWarehouseReceive:
		// 到仓收货: 采购在途 → 待检
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.PurchasingInTransit < decrease {
			return nil, fmt.Errorf("采购在途库存不足")
		}
		newBalance.PurchasingInTransit = balance.PurchasingInTransit - decrease
		newBalance.PendingInspection = balance.PendingInspection + decrease

	case domain.MovementTypeInspectionPass:
		// 质检通过: 待检 → 原料库存
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.PendingInspection < decrease {
			return nil, fmt.Errorf("待检库存不足")
		}
		newBalance.PendingInspection = balance.PendingInspection - decrease
		newBalance.RawMaterial = balance.RawMaterial + decrease

	case domain.MovementTypeInspectionFail:
		// 质检不合格: 待检 → 损坏
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.PendingInspection < decrease {
			return nil, fmt.Errorf("待检库存不足")
		}
		newBalance.PendingInspection = balance.PendingInspection - decrease
		newBalance.DamagedQuantity = balance.DamagedQuantity + decrease

	case domain.MovementTypeInspectionLoss:
		// 采购质检损失: 待检 → 损失结案（不进入任何库存池）
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.PendingInspection < decrease {
			return nil, fmt.Errorf("待检库存不足")
		}
		newBalance.PendingInspection = balance.PendingInspection - decrease

	case domain.MovementTypeAssemblyConsume:
		// 组装耗料: 子件原料库存 → 组装消耗
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.RawMaterial < decrease {
			return nil, fmt.Errorf("原料库存不足")
		}
		newBalance.RawMaterial = balance.RawMaterial - decrease

	case domain.MovementTypeAssemblyComplete:
		// 组装完成: 主件产出 → 待出库存
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		increase := uint(quantity)
		newBalance.PendingShipment = balance.PendingShipment + increase

	case domain.MovementTypePackingSkipComplete:
		// 免打包直通: 原料库存 → 待出库存
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.RawMaterial < decrease {
			return nil, fmt.Errorf("原料库存不足")
		}
		newBalance.RawMaterial = balance.RawMaterial - decrease
		newBalance.PendingShipment = balance.PendingShipment + decrease

	case domain.MovementTypeShipmentAllocate:
		// 发货确认锁定: 待出库存 -> 待出锁定
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.PendingShipment < decrease {
			return nil, fmt.Errorf("待出库存不足")
		}
		newBalance.PendingShipment = balance.PendingShipment - decrease
		newBalance.PendingShipmentReserved = balance.PendingShipmentReserved + decrease

	case domain.MovementTypeShipmentRelease:
		// 发货取消释放: 待出锁定 -> 待出库存
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.PendingShipmentReserved < decrease {
			return nil, fmt.Errorf("待出锁定库存不足")
		}
		newBalance.PendingShipmentReserved = balance.PendingShipmentReserved - decrease
		newBalance.PendingShipment = balance.PendingShipment + decrease

	case domain.MovementTypeLogisticsShip:
		// 物流发货: 待出库存 → 物流在途
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.PendingShipment < decrease {
			return nil, fmt.Errorf("待出库存不足")
		}
		newBalance.PendingShipment = balance.PendingShipment - decrease
		newBalance.LogisticsInTransit = balance.LogisticsInTransit + decrease

	case domain.MovementTypePlatformReceive:
		// 平台上架: 物流在途 → 可售库存
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.LogisticsInTransit < decrease {
			return nil, fmt.Errorf("物流在途库存不足")
		}
		newBalance.LogisticsInTransit = balance.LogisticsInTransit - decrease
		newBalance.Sellable = balance.Sellable + decrease

	// 发货单库存流转
	case domain.MovementTypeShipmentShip:
		// 发货单发货: 待出锁定 → 在途
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.PendingShipmentReserved < decrease {
			return nil, fmt.Errorf("待出锁定库存不足")
		}
		newBalance.PendingShipmentReserved = balance.PendingShipmentReserved - decrease
		newBalance.LogisticsInTransit = balance.LogisticsInTransit + decrease

	default:
		return nil, fmt.Errorf("unsupported movement type: %s", movementType)
	}

	return newBalance, nil
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func validateMovementParams(params *domain.CreateMovementParams) error {
	if params == nil {
		return nil
	}
	if params.MovementType == domain.MovementTypePlatformReceive {
		if params.ReferenceType == nil || !strings.EqualFold(strings.TrimSpace(*params.ReferenceType), "SHIPMENT") {
			return ErrPlatformReceiveRequiresShipmentReference
		}
		if (params.ReferenceID == nil || *params.ReferenceID == 0) && (params.ReferenceNumber == nil || strings.TrimSpace(*params.ReferenceNumber) == "") {
			return ErrPlatformReceiveRequiresShipmentReference
		}
	}
	return nil
}

func (u *InventoryUsecase) ensureLotSeed(ctx context.Context, balance *domain.InventoryBalance) error {
	if u.lotRepo == nil || balance == nil {
		return nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, balance.ProductID, balance.WarehouseID)
	if err != nil {
		return err
	}
	if err := u.backfillSeedLotUnitCosts(ctx, lots, balance.ProductID, balance.WarehouseID); err != nil {
		return err
	}

	var totalAvailable uint
	var totalReserved uint
	var totalPurchasingInTransit uint
	var totalPendingInspection uint
	var totalRawMaterial uint
	var totalSellable uint
	var totalSellableReserved uint
	var totalPendingShipment uint
	var totalPendingShipmentReserved uint
	for _, lot := range lots {
		totalPurchasingInTransit += lot.QtyPurchasingInTransit
		totalPendingInspection += lot.QtyPendingInspection
		totalRawMaterial += lot.QtyRawMaterial
		totalAvailable += lot.QtyAvailable
		totalReserved += lot.QtyReserved
		totalSellable += lot.QtySellable
		totalSellableReserved += lot.QtySellableReserved
		totalPendingShipment += lot.QtyPendingShipment
		totalPendingShipmentReserved += lot.QtyPendingShipmentReserved
	}

	if balance.AvailableQuantity > totalAvailable {
		diff := balance.AvailableQuantity - totalAvailable
		lot := u.newSeedLot(balance, "INIT_AVAILABLE", diff, 0)
		u.assignSeedLotUnitCost(ctx, lot)
		if err := u.lotRepo.Create(ctx, lot); err != nil {
			return err
		}
	}

	if balance.ReservedQuantity > totalReserved {
		diff := balance.ReservedQuantity - totalReserved
		lot := u.newSeedLot(balance, "INIT_RESERVED", 0, diff)
		u.assignSeedLotUnitCost(ctx, lot)
		if err := u.lotRepo.Create(ctx, lot); err != nil {
			return err
		}
	}
	if balance.PurchasingInTransit > totalPurchasingInTransit {
		diff := balance.PurchasingInTransit - totalPurchasingInTransit
		lot := u.newSeedLot(balance, "INIT_PURCHASING_IN_TRANSIT", 0, 0)
		lot.QtyPurchasingInTransit = diff
		lot.QtyIn = diff
		lot.Status = lotStatus(lot)
		u.assignSeedLotUnitCost(ctx, lot)
		if err := u.lotRepo.Create(ctx, lot); err != nil {
			return err
		}
	}
	if balance.PendingInspection > totalPendingInspection {
		diff := balance.PendingInspection - totalPendingInspection
		lot := u.newSeedLot(balance, "INIT_PENDING_INSPECTION", 0, 0)
		lot.QtyPendingInspection = diff
		lot.QtyIn = diff
		lot.Status = lotStatus(lot)
		u.assignSeedLotUnitCost(ctx, lot)
		if err := u.lotRepo.Create(ctx, lot); err != nil {
			return err
		}
	}
	if balance.RawMaterial > totalRawMaterial {
		diff := balance.RawMaterial - totalRawMaterial
		lot := u.newSeedLot(balance, "INIT_RAW_MATERIAL", 0, 0)
		lot.QtyRawMaterial = diff
		lot.QtyIn = diff
		lot.Status = lotStatus(lot)
		u.assignSeedLotUnitCost(ctx, lot)
		if err := u.lotRepo.Create(ctx, lot); err != nil {
			return err
		}
	}
	if balance.Sellable > totalSellable {
		diff := balance.Sellable - totalSellable
		lot := u.newSeedLot(balance, "INIT_SELLABLE", 0, 0)
		lot.QtySellable = diff
		lot.QtyIn = diff
		lot.Status = lotStatus(lot)
		u.assignSeedLotUnitCost(ctx, lot)
		if err := u.lotRepo.Create(ctx, lot); err != nil {
			return err
		}
	}
	if balance.SellableReserved > totalSellableReserved {
		diff := balance.SellableReserved - totalSellableReserved
		lot := u.newSeedLot(balance, "INIT_SELLABLE_RESERVED", 0, 0)
		lot.QtySellableReserved = diff
		lot.QtyIn = diff
		lot.Status = lotStatus(lot)
		u.assignSeedLotUnitCost(ctx, lot)
		if err := u.lotRepo.Create(ctx, lot); err != nil {
			return err
		}
	}
	if balance.PendingShipment > totalPendingShipment {
		diff := balance.PendingShipment - totalPendingShipment
		lot := u.newSeedLot(balance, "INIT_PENDING_SHIPMENT", 0, 0)
		lot.QtyPendingShipment = diff
		lot.QtyIn = diff
		lot.Status = lotStatus(lot)
		u.assignSeedLotUnitCost(ctx, lot)
		if err := u.lotRepo.Create(ctx, lot); err != nil {
			return err
		}
	}
	if balance.PendingShipmentReserved > totalPendingShipmentReserved {
		diff := balance.PendingShipmentReserved - totalPendingShipmentReserved
		lot := u.newSeedLot(balance, "INIT_PENDING_SHIPMENT_RESERVED", 0, 0)
		lot.QtyPendingShipmentReserved = diff
		lot.QtyIn = diff
		lot.Status = lotStatus(lot)
		u.assignSeedLotUnitCost(ctx, lot)
		if err := u.lotRepo.Create(ctx, lot); err != nil {
			return err
		}
	}

	return nil
}

func (u *InventoryUsecase) newSeedLot(balance *domain.InventoryBalance, sourceType string, qtyAvailable uint, qtyReserved uint) *domain.InventoryLot {
	receivedAt := time.Unix(0, 0)
	lotNo := numbering.Generate("LOT", time.Now())
	status := domain.InventoryLotStatusOpen
	if qtyAvailable == 0 && qtyReserved == 0 {
		status = domain.InventoryLotStatusClosed
	}

	remark := "seeded from existing balance"
	sourceTypeVal := sourceType
	return &domain.InventoryLot{
		ProductID:    balance.ProductID,
		WarehouseID:  balance.WarehouseID,
		LotNo:        lotNo,
		SourceType:   &sourceTypeVal,
		ReceivedAt:   receivedAt,
		QtyIn:        qtyAvailable + qtyReserved,
		QtyAvailable: qtyAvailable,
		QtyReserved:  qtyReserved,
		QtyConsumed:  0,
		Status:       status,
		Remark:       &remark,
	}
}

func (u *InventoryUsecase) backfillSeedLotUnitCosts(ctx context.Context, lots []*domain.InventoryLot, productID, warehouseID uint64) error {
	if u == nil || u.lotRepo == nil || u.seedLotUnitCostResolver == nil {
		return nil
	}
	resolved, err := u.seedLotUnitCostResolver(ctx, productID, warehouseID)
	if err != nil {
		return err
	}
	if resolved == nil || *resolved <= 0 {
		return nil
	}
	for _, lot := range lots {
		if !needsSeedLotUnitCost(lot) {
			continue
		}
		value := *resolved
		lot.UnitCost = &value
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
	}
	return nil
}

func (u *InventoryUsecase) backfillShipmentLotUnitCosts(ctx context.Context, lots []*domain.InventoryLot) error {
	if u == nil || u.lotRepo == nil || u.shipmentLotUnitCostResolver == nil {
		return nil
	}
	for _, lot := range lots {
		if !needsShipmentLotUnitCostBackfill(lot) {
			continue
		}
		resolved, err := u.shipmentLotUnitCostResolver(ctx, lot.ProductID, lot.SourceType, lot.SourceID, lot.SourceNumber, lot.ReceivedAt)
		if err != nil {
			return err
		}
		if resolved == nil || *resolved <= 0 {
			continue
		}
		if lot.UnitCost != nil && math.Abs(*lot.UnitCost-*resolved) < 0.000001 {
			continue
		}
		value := *resolved
		lot.UnitCost = &value
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
	}
	return nil
}

func (u *InventoryUsecase) assignSeedLotUnitCost(ctx context.Context, lot *domain.InventoryLot) {
	if u == nil || lot == nil || u.seedLotUnitCostResolver == nil {
		return
	}
	resolved, err := u.seedLotUnitCostResolver(ctx, lot.ProductID, lot.WarehouseID)
	if err != nil || resolved == nil || *resolved <= 0 {
		return
	}
	value := *resolved
	lot.UnitCost = &value
}

func needsSeedLotUnitCost(lot *domain.InventoryLot) bool {
	if lot == nil {
		return false
	}
	if lot.SourceType == nil || !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(*lot.SourceType)), "INIT_") {
		return false
	}
	if lot.UnitCost != nil && *lot.UnitCost > 0 {
		return false
	}
	return lot.QtyAvailable > 0 ||
		lot.QtyReserved > 0 ||
		lot.QtyPurchasingInTransit > 0 ||
		lot.QtyPendingInspection > 0 ||
		lot.QtyRawMaterial > 0 ||
		lot.QtyPendingShipment > 0 ||
		lot.QtyPendingShipmentReserved > 0 ||
		lot.QtySellable > 0 ||
		lot.QtySellableReserved > 0
}

func needsShipmentLotUnitCostBackfill(lot *domain.InventoryLot) bool {
	if lot == nil {
		return false
	}
	if lot.SourceType == nil || !strings.EqualFold(strings.TrimSpace(*lot.SourceType), "SHIPMENT") {
		return false
	}
	return lot.QtySellable > 0 || lot.QtySellableReserved > 0
}

func (u *InventoryUsecase) applyLotMovement(
	ctx context.Context,
	params *domain.CreateMovementParams,
	balance *domain.InventoryBalance,
) ([]domain.InventoryLotAllocation, error) {
	if u.lotRepo == nil || params == nil || balance == nil {
		return nil, nil
	}

	switch params.MovementType {
	case domain.MovementTypePurchaseReceipt,
		domain.MovementTypeReturnReceipt,
		domain.MovementTypeTransferIn:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		if params.MovementType == domain.MovementTypeReturnReceipt && stockPool(params) == domain.StockPoolSellable {
			return nil, u.createInboundSellableLot(ctx, params, uint(params.Quantity))
		}
		return nil, u.createInboundLot(ctx, params, uint(params.Quantity))

	case domain.MovementTypePurchaseShip:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.createInboundPurchasingInTransitLot(ctx, params, uint(params.Quantity))

	case domain.MovementTypeWarehouseReceive:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.movePurchasingInTransitToPendingInspection(ctx, params, uint(params.Quantity))

	case domain.MovementTypeInspectionPass:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.movePendingInspectionToRawMaterial(ctx, params, uint(params.Quantity))

	case domain.MovementTypeInspectionFail:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.consumePendingInspectionLots(ctx, params, uint(params.Quantity))

	case domain.MovementTypeInspectionLoss:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.consumePendingInspectionLots(ctx, params, uint(params.Quantity))

	case domain.MovementTypeAssemblyConsume:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.consumeRawMaterialLots(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))

	case domain.MovementTypePlatformReceive:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.createInboundSellableLot(ctx, params, uint(params.Quantity))

	case domain.MovementTypeAssemblyComplete:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.createInboundPendingShipmentLot(ctx, params, uint(params.Quantity))

	case domain.MovementTypePackingSkipComplete:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.moveRawMaterialToPendingShipment(ctx, params, uint(params.Quantity))

	case domain.MovementTypeStockTakeAdjustment,
		domain.MovementTypeManualAdjustment:
		if params.Quantity > 0 {
			return nil, u.createInboundLot(ctx, params, uint(params.Quantity))
		}
		if params.Quantity < 0 {
			return nil, u.consumeAvailableLots(ctx, params.ProductID, params.WarehouseID, uint(-params.Quantity))
		}
		return nil, nil

	case domain.MovementTypeSalesShipment,
		domain.MovementTypeTransferOut,
		domain.MovementTypeDamageWriteOff:
		if params.Quantity > 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.consumeAvailableLots(ctx, params.ProductID, params.WarehouseID, uint(-params.Quantity))

	case domain.MovementTypeSalesAllocate:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		if stockPool(params) == domain.StockPoolSellable {
			return nil, u.reserveSellableLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))
		}
		return nil, u.reserveLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))

	case domain.MovementTypeSalesRelease:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		if stockPool(params) == domain.StockPoolSellable {
			return nil, u.releaseSellableReservedLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))
		}
		return nil, u.releaseReservedLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))

	case domain.MovementTypeSalesShip:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		if stockPool(params) == domain.StockPoolSellable {
			return u.consumeSellableReservedLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))
		}
		return u.consumeReservedLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))

	case domain.MovementTypeShipmentAllocate:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.reservePendingShipmentLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))

	case domain.MovementTypeShipmentRelease:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.releasePendingShipmentReservedLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))

	case domain.MovementTypeShipmentShip:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return u.consumePendingShipmentReservedLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))

	case domain.MovementTypeLogisticsShip:
		if params.Quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		return nil, u.consumePendingShipmentLotsFIFO(ctx, params.ProductID, params.WarehouseID, uint(params.Quantity))
	}

	return nil, nil
}

func (u *InventoryUsecase) createInboundLot(ctx context.Context, params *domain.CreateMovementParams, qty uint) error {
	if qty == 0 {
		return nil
	}

	operatedAt := time.Now()
	if params.OperatedAt != nil {
		operatedAt = *params.OperatedAt
	}
	lotNo := u.buildLotNo(params, operatedAt)
	status := domain.InventoryLotStatusOpen
	sourceType := "INVENTORY_MOVEMENT"
	if params.ReferenceType != nil && *params.ReferenceType != "" {
		sourceType = *params.ReferenceType
	}

	lot := &domain.InventoryLot{
		ProductID:                  params.ProductID,
		WarehouseID:                params.WarehouseID,
		LotNo:                      lotNo,
		SourceType:                 &sourceType,
		SourceID:                   params.ReferenceID,
		SourceNumber:               params.ReferenceNumber,
		ReceivedAt:                 operatedAt,
		UnitCost:                   params.UnitCost,
		QtyIn:                      qty,
		QtyAvailable:               qty,
		QtyReserved:                0,
		QtyPendingShipment:         0,
		QtyPendingShipmentReserved: 0,
		QtyConsumed:                0,
		Status:                     status,
	}
	return u.lotRepo.Create(ctx, lot)
}

func (u *InventoryUsecase) createInboundPurchasingInTransitLot(ctx context.Context, params *domain.CreateMovementParams, qty uint) error {
	if qty == 0 {
		return nil
	}

	operatedAt := time.Now()
	if params.OperatedAt != nil {
		operatedAt = *params.OperatedAt
	}
	lotNo := u.buildLotNo(params, operatedAt)
	status := domain.InventoryLotStatusOpen
	sourceType := "INVENTORY_MOVEMENT"
	if params.ReferenceType != nil && *params.ReferenceType != "" {
		sourceType = *params.ReferenceType
	}

	lot := &domain.InventoryLot{
		ProductID:              params.ProductID,
		WarehouseID:            params.WarehouseID,
		LotNo:                  lotNo,
		SourceType:             &sourceType,
		SourceID:               params.ReferenceID,
		SourceNumber:           params.ReferenceNumber,
		ReceivedAt:             operatedAt,
		UnitCost:               params.UnitCost,
		QtyIn:                  qty,
		QtyPurchasingInTransit: qty,
		Status:                 status,
	}
	return u.lotRepo.Create(ctx, lot)
}

func (u *InventoryUsecase) createInboundSellableLot(ctx context.Context, params *domain.CreateMovementParams, qty uint) error {
	if qty == 0 {
		return nil
	}

	operatedAt := time.Now()
	if params.OperatedAt != nil {
		operatedAt = *params.OperatedAt
	}
	lotNo := u.buildLotNo(params, operatedAt)
	status := domain.InventoryLotStatusOpen
	sourceType := "INVENTORY_MOVEMENT"
	if params.ReferenceType != nil && *params.ReferenceType != "" {
		sourceType = *params.ReferenceType
	}

	lot := &domain.InventoryLot{
		ProductID:    params.ProductID,
		WarehouseID:  params.WarehouseID,
		LotNo:        lotNo,
		SourceType:   &sourceType,
		SourceID:     params.ReferenceID,
		SourceNumber: params.ReferenceNumber,
		ReceivedAt:   operatedAt,
		UnitCost:     params.UnitCost,
		QtyIn:        qty,
		Status:       status,
		QtySellable:  qty,
	}
	return u.lotRepo.Create(ctx, lot)
}

func (u *InventoryUsecase) createInboundPendingShipmentLot(ctx context.Context, params *domain.CreateMovementParams, qty uint) error {
	if qty == 0 {
		return nil
	}

	operatedAt := time.Now()
	if params.OperatedAt != nil {
		operatedAt = *params.OperatedAt
	}
	lotNo := u.buildLotNo(params, operatedAt)
	status := domain.InventoryLotStatusOpen
	sourceType := "INVENTORY_MOVEMENT"
	if params.ReferenceType != nil && *params.ReferenceType != "" {
		sourceType = *params.ReferenceType
	}

	lot := &domain.InventoryLot{
		ProductID:          params.ProductID,
		WarehouseID:        params.WarehouseID,
		LotNo:              lotNo,
		SourceType:         &sourceType,
		SourceID:           params.ReferenceID,
		SourceNumber:       params.ReferenceNumber,
		ReceivedAt:         operatedAt,
		UnitCost:           params.UnitCost,
		QtyIn:              qty,
		QtyPendingShipment: qty,
		Status:             status,
	}
	return u.lotRepo.Create(ctx, lot)
}

func (u *InventoryUsecase) movePurchasingInTransitToPendingInspection(ctx context.Context, params *domain.CreateMovementParams, qty uint) error {
	if qty == 0 {
		return nil
	}
	if params == nil {
		return ErrInvalidQuantity
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, params.ProductID, params.WarehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range reorderLotsForReference(lots, params.ReferenceID, params.ReferenceNumber) {
		if remaining == 0 {
			break
		}
		if lot.QtyPurchasingInTransit == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyPurchasingInTransit)
		lot.QtyPurchasingInTransit -= take
		lot.QtyPendingInspection += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func reorderLotsForReference(lots []*domain.InventoryLot, referenceID *uint64, referenceNumber *string) []*domain.InventoryLot {
	if len(lots) <= 1 {
		return lots
	}
	preferred := make([]*domain.InventoryLot, 0, len(lots))
	others := make([]*domain.InventoryLot, 0, len(lots))
	for _, lot := range lots {
		if lotMatchesReference(lot, referenceID, referenceNumber) {
			preferred = append(preferred, lot)
			continue
		}
		others = append(others, lot)
	}
	if len(preferred) == 0 {
		return lots
	}
	return append(preferred, others...)
}

func lotMatchesReference(lot *domain.InventoryLot, referenceID *uint64, referenceNumber *string) bool {
	if lot == nil {
		return false
	}
	if referenceID != nil && *referenceID != 0 && lot.SourceID != nil && *lot.SourceID == *referenceID {
		return true
	}
	if referenceNumber != nil && strings.TrimSpace(*referenceNumber) != "" && lot.SourceNumber != nil && strings.TrimSpace(*lot.SourceNumber) == strings.TrimSpace(*referenceNumber) {
		return true
	}
	return false
}

func (u *InventoryUsecase) movePendingInspectionToRawMaterial(ctx context.Context, params *domain.CreateMovementParams, qty uint) error {
	if qty == 0 {
		return nil
	}
	if params == nil {
		return ErrInvalidQuantity
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, params.ProductID, params.WarehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range reorderLotsForReference(lots, params.ReferenceID, params.ReferenceNumber) {
		if remaining == 0 {
			break
		}
		if lot.QtyPendingInspection == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyPendingInspection)
		lot.QtyPendingInspection -= take
		lot.QtyRawMaterial += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) consumePendingInspectionLots(ctx context.Context, params *domain.CreateMovementParams, qty uint) error {
	if qty == 0 {
		return nil
	}
	if params == nil {
		return ErrInvalidQuantity
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, params.ProductID, params.WarehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range reorderLotsForReference(lots, params.ReferenceID, params.ReferenceNumber) {
		if remaining == 0 {
			break
		}
		if lot.QtyPendingInspection == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyPendingInspection)
		lot.QtyPendingInspection -= take
		lot.QtyConsumed += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) consumeRawMaterialLots(ctx context.Context, productID, warehouseID uint64, qty uint) error {
	if qty == 0 {
		return nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtyRawMaterial == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyRawMaterial)
		lot.QtyRawMaterial -= take
		lot.QtyConsumed += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) moveRawMaterialToPendingShipment(ctx context.Context, params *domain.CreateMovementParams, qty uint) error {
	if qty == 0 {
		return nil
	}
	if params == nil {
		return ErrInvalidQuantity
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, params.ProductID, params.WarehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range reorderLotsForReference(lots, params.ReferenceID, params.ReferenceNumber) {
		if remaining == 0 {
			break
		}
		if lot.QtyRawMaterial == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyRawMaterial)
		lot.QtyRawMaterial -= take
		lot.QtyPendingShipment += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) reserveLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) error {
	if qty == 0 {
		return nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtyAvailable == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyAvailable)
		lot.QtyAvailable -= take
		lot.QtyReserved += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) releaseReservedLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) error {
	if qty == 0 {
		return nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtyReserved == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyReserved)
		lot.QtyReserved -= take
		lot.QtyAvailable += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) reservePendingShipmentLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) error {
	if qty == 0 {
		return nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtyPendingShipment == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyPendingShipment)
		lot.QtyPendingShipment -= take
		lot.QtyPendingShipmentReserved += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) releasePendingShipmentReservedLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) error {
	if qty == 0 {
		return nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtyPendingShipmentReserved == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyPendingShipmentReserved)
		lot.QtyPendingShipmentReserved -= take
		lot.QtyPendingShipment += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) consumePendingShipmentLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) error {
	if qty == 0 {
		return nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtyPendingShipment == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyPendingShipment)
		lot.QtyPendingShipment -= take
		lot.QtyConsumed += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) reserveSellableLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) error {
	if qty == 0 {
		return nil
	}
	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}
	if err := u.backfillShipmentLotUnitCosts(ctx, lots); err != nil {
		return err
	}
	remaining := qty
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtySellable == 0 {
			continue
		}
		take := minUint(remaining, lot.QtySellable)
		lot.QtySellable -= take
		lot.QtySellableReserved += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}
	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) releaseSellableReservedLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) error {
	if qty == 0 {
		return nil
	}
	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}
	if err := u.backfillShipmentLotUnitCosts(ctx, lots); err != nil {
		return err
	}
	remaining := qty
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtySellableReserved == 0 {
			continue
		}
		take := minUint(remaining, lot.QtySellableReserved)
		lot.QtySellableReserved -= take
		lot.QtySellable += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}
	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) consumeReservedLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) ([]domain.InventoryLotAllocation, error) {
	if qty == 0 {
		return nil, nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return nil, err
	}

	remaining := qty
	allocations := make([]domain.InventoryLotAllocation, 0)
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtyReserved == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyReserved)
		lot.QtyReserved -= take
		lot.QtyConsumed += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return nil, err
		}
		unitCost := 0.0
		if lot.UnitCost != nil {
			unitCost = *lot.UnitCost
		}
		allocations = append(allocations, domain.InventoryLotAllocation{
			InventoryLotID: lot.ID,
			Qty:            uint64(take),
			UnitCost:       unitCost,
		})
		remaining -= take
	}

	if remaining > 0 {
		return nil, ErrInsufficientStock
	}
	return allocations, nil
}

func (u *InventoryUsecase) consumeSellableReservedLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) ([]domain.InventoryLotAllocation, error) {
	if qty == 0 {
		return nil, nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return nil, err
	}
	if err := u.backfillShipmentLotUnitCosts(ctx, lots); err != nil {
		return nil, err
	}

	remaining := qty
	allocations := make([]domain.InventoryLotAllocation, 0)
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtySellableReserved == 0 {
			continue
		}
		take := minUint(remaining, lot.QtySellableReserved)
		lot.QtySellableReserved -= take
		lot.QtyConsumed += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return nil, err
		}
		unitCost := 0.0
		if lot.UnitCost != nil {
			unitCost = *lot.UnitCost
		}
		allocations = append(allocations, domain.InventoryLotAllocation{
			InventoryLotID: lot.ID,
			Qty:            uint64(take),
			UnitCost:       unitCost,
		})
		remaining -= take
	}
	if remaining > 0 {
		return nil, ErrInsufficientStock
	}
	return allocations, nil
}

func (u *InventoryUsecase) consumePendingShipmentReservedLotsFIFO(ctx context.Context, productID, warehouseID uint64, qty uint) ([]domain.InventoryLotAllocation, error) {
	if qty == 0 {
		return nil, nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return nil, err
	}

	remaining := qty
	allocations := make([]domain.InventoryLotAllocation, 0)
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtyPendingShipmentReserved == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyPendingShipmentReserved)
		lot.QtyPendingShipmentReserved -= take
		lot.QtyConsumed += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return nil, err
		}
		unitCost := 0.0
		if lot.UnitCost != nil {
			unitCost = *lot.UnitCost
		}
		allocations = append(allocations, domain.InventoryLotAllocation{
			InventoryLotID: lot.ID,
			Qty:            uint64(take),
			UnitCost:       unitCost,
		})
		remaining -= take
	}

	if remaining > 0 {
		return nil, ErrInsufficientStock
	}
	return allocations, nil
}

func (u *InventoryUsecase) consumeAvailableLots(ctx context.Context, productID, warehouseID uint64, qty uint) error {
	if qty == 0 {
		return nil
	}

	lots, err := u.lotRepo.ListByProductAndWarehouse(ctx, productID, warehouseID)
	if err != nil {
		return err
	}

	remaining := qty
	for _, lot := range lots {
		if remaining == 0 {
			break
		}
		if lot.QtyAvailable == 0 {
			continue
		}
		take := minUint(remaining, lot.QtyAvailable)
		lot.QtyAvailable -= take
		lot.QtyConsumed += take
		lot.Status = lotStatus(lot)
		if err := u.lotRepo.Update(ctx, lot); err != nil {
			return err
		}
		remaining -= take
	}

	if remaining > 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (u *InventoryUsecase) buildLotNo(params *domain.CreateMovementParams, operatedAt time.Time) string {
	return numbering.Generate("LOT", operatedAt)
}

func buildPackingReferenceNumber() string {
	return numbering.Generate("PACK", time.Now())
}

func (u *InventoryUsecase) createPackingPackagingLedgers(
	ctx context.Context,
	movement *domain.InventoryMovement,
	requirements []PackingRequirement,
	productQty int,
	operatorID uint64,
	operatedAt time.Time,
	packagingItemRepo packagingDomain.PackagingItemRepository,
	packagingLedgerRepo packagingDomain.PackagingLedgerRepository,
) ([]PackingMaterialCostLine, error) {
	if movement == nil || movement.ID == 0 || productQty <= 0 {
		return nil, nil
	}

	type aggregatedRequirement struct {
		PackagingItemID uint64
		RequiredQty     uint64
		ItemCode        string
		ItemName        string
		Unit            string
	}

	aggregated := map[uint64]*aggregatedRequirement{}
	for _, requirement := range requirements {
		if requirement.PackagingItemID == 0 || requirement.QuantityPerUnit <= 0 {
			return nil, ErrPackingMaterialsNotConfigured
		}
		requiredQty, err := resolvePackingRequirementQuantity(productQty, requirement.QuantityPerUnit)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", err, requirement.ItemCode)
		}
		if requiredQty == 0 {
			continue
		}
		if existing, ok := aggregated[requirement.PackagingItemID]; ok {
			existing.RequiredQty += requiredQty
			continue
		}
		aggregated[requirement.PackagingItemID] = &aggregatedRequirement{
			PackagingItemID: requirement.PackagingItemID,
			RequiredQty:     requiredQty,
			ItemCode:        requirement.ItemCode,
			ItemName:        requirement.ItemName,
			Unit:            requirement.Unit,
		}
	}

	referenceType := "PRODUCT_PACKING"
	note := buildPackingLedgerNote(movement.ReferenceNumber)
	costLines := make([]PackingMaterialCostLine, 0, len(aggregated))

	for _, requirement := range aggregated {
		item, err := packagingItemRepo.GetByID(requirement.PackagingItemID)
		if err != nil {
			return nil, fmt.Errorf("packaging item not found: %w", err)
		}
		if strings.TrimSpace(item.Status) != "" && !strings.EqualFold(strings.TrimSpace(item.Status), "ACTIVE") {
			return nil, fmt.Errorf("包装材料已停用: %s", firstNonEmpty(item.ItemName, item.ItemCode))
		}
		if item.QuantityOnHand < requirement.RequiredQty {
			return nil, fmt.Errorf("包材库存不足: %s", firstNonEmpty(item.ItemName, item.ItemCode))
		}

		referenceID := movement.ID
		ledger := &packagingDomain.PackagingLedger{
			TraceID:         firstNonEmpty(pointerStringValue(movement.TraceID), uuid.New().String()),
			PackagingItemID: requirement.PackagingItemID,
			TransactionType: "OUT",
			Quantity:        -int64(requirement.RequiredQty),
			UnitCost:        item.UnitCost,
			QuantityBefore:  item.QuantityOnHand,
			QuantityAfter:   item.QuantityOnHand - requirement.RequiredQty,
			ReferenceType:   &referenceType,
			ReferenceID:     &referenceID,
			OccurredAt:      operatedAt,
			Notes:           note,
			CreatedBy:       operatorID,
		}
		if err := packagingLedgerRepo.Create(ledger); err != nil {
			return nil, err
		}
		if err := packagingItemRepo.UpdateQuantity(requirement.PackagingItemID, -int64(requirement.RequiredQty)); err != nil {
			return nil, err
		}
		costLines = append(costLines, PackingMaterialCostLine{
			PackagingItemID: requirement.PackagingItemID,
			Quantity:        requirement.RequiredQty,
			UnitCost:        item.UnitCost,
			Currency:        item.Currency,
			ItemCode:        item.ItemCode,
			ItemName:        item.ItemName,
		})
	}

	return costLines, nil
}

func resolvePackingRequirementQuantity(productQty int, quantityPerUnit float64) (uint64, error) {
	required := float64(productQty) * quantityPerUnit
	if required < 0 {
		return 0, ErrPackingMaterialQuantityInvalid
	}
	rounded := math.Round(required)
	if math.Abs(required-rounded) > 1e-9 {
		return 0, ErrPackingMaterialQuantityInvalid
	}
	return uint64(rounded), nil
}

func buildPackingLedgerNote(referenceNumber *string) *string {
	if referenceNumber == nil || strings.TrimSpace(*referenceNumber) == "" {
		return nil
	}
	note := fmt.Sprintf("打包单号:%s", strings.TrimSpace(*referenceNumber))
	return &note
}

func pointerStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func lotStatus(lot *domain.InventoryLot) domain.InventoryLotStatus {
	if lot == nil {
		return domain.InventoryLotStatusClosed
	}
	if lot.QtyAvailable == 0 &&
		lot.QtyReserved == 0 &&
		lot.QtyPurchasingInTransit == 0 &&
		lot.QtyPendingInspection == 0 &&
		lot.QtyRawMaterial == 0 &&
		lot.QtyPendingShipment == 0 &&
		lot.QtyPendingShipmentReserved == 0 &&
		lot.QtySellable == 0 &&
		lot.QtySellableReserved == 0 {
		return domain.InventoryLotStatusClosed
	}
	return domain.InventoryLotStatusOpen
}

func stockPool(params *domain.CreateMovementParams) domain.StockPool {
	if params != nil && params.StockPool != nil {
		return *params.StockPool
	}
	return domain.StockPoolAvailable
}

func calculateBalanceTotal(balance *domain.InventoryBalance) uint {
	if balance == nil {
		return 0
	}
	return balance.AvailableQuantity +
		balance.ReservedQuantity +
		balance.SellableReserved +
		balance.DamagedQuantity +
		balance.PurchasingInTransit +
		balance.PendingInspection +
		balance.RawMaterial +
		balance.PendingShipment +
		balance.PendingShipmentReserved +
		balance.LogisticsInTransit +
		balance.Sellable +
		balance.Returned
}

func minUint(a uint, b uint) uint {
	if a < b {
		return a
	}
	return b
}

// RecordReturnReceive 退货入库
func (u *InventoryUsecase) RecordReturnReceive(ctx context.Context, params *domain.StockTransitionParams) (*domain.InventoryMovement, error) {
	if params.Quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	// Get or create balance record
	balance, err := u.balanceRepo.GetOrCreate(ctx, params.ProductID, params.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// 退货入库: 增加退货库存
	operatedAt := time.Now()
	traceID := uuid.New().String()
	movement := &domain.InventoryMovement{
		TraceID:                &traceID,
		ProductID:              params.ProductID,
		WarehouseID:            params.WarehouseID,
		MovementType:           domain.MovementTypeReturnReceipt,
		ReferenceType:          params.ReferenceType,
		ReferenceID:            params.ReferenceID,
		ReferenceNumber:        params.ReferenceNumber,
		Quantity:               int(params.Quantity),
		BeforeAvailable:        balance.AvailableQuantity,
		AfterAvailable:         balance.AvailableQuantity,
		BeforeReserved:         balance.ReservedQuantity,
		AfterReserved:          balance.ReservedQuantity,
		BeforeSellable:         balance.Sellable,
		AfterSellable:          balance.Sellable,
		BeforeSellableReserved: balance.SellableReserved,
		AfterSellableReserved:  balance.SellableReserved,
		BeforeDamaged:          balance.DamagedQuantity,
		AfterDamaged:           balance.DamagedQuantity,
		UnitCost:               params.UnitCost,
		Remark:                 params.Remark,
		OperatorID:             params.OperatorID,
		OperatedAt:             operatedAt,
		GmtCreate:              time.Now(),
		GmtModified:            time.Now(),
	}

	if params.UnitCost != nil {
		totalCost := *params.UnitCost * float64(params.Quantity)
		movement.TotalCost = &totalCost
	}

	// Save movement
	if err := u.movementRepo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("failed to create movement: %w", err)
	}

	// Update balance - 增加退货库存
	balance.Returned = balance.Returned + params.Quantity
	balance.TotalQuantity = calculateBalanceTotal(balance)
	balance.LastMovementAt = &operatedAt
	balance.GmtModified = time.Now()

	if err := u.balanceRepo.Update(ctx, balance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	return movement, nil
}

// ReturnInspect 退货质检
func (u *InventoryUsecase) ReturnInspect(ctx context.Context, params *domain.ReturnInspectParams) error {
	totalQuantity := params.PassQuantity + params.FailQuantity
	if totalQuantity == 0 {
		return ErrInvalidQuantity
	}

	// Get balance record
	balance, err := u.balanceRepo.GetOrCreate(ctx, params.ProductID, params.WarehouseID)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	if balance.Returned < totalQuantity {
		return fmt.Errorf("退货库存不足")
	}

	traceID := uuid.New().String()
	operatedAt := time.Now()

	// 质检通过部分: 退货库存 → 待检库存
	if params.PassQuantity > 0 {
		passMovement := &domain.InventoryMovement{
			TraceID:                &traceID,
			ProductID:              params.ProductID,
			WarehouseID:            params.WarehouseID,
			MovementType:           domain.MovementTypeReturnInspect,
			ReferenceType:          params.ReferenceType,
			ReferenceNumber:        params.ReferenceNumber,
			Quantity:               int(params.PassQuantity),
			BeforeAvailable:        balance.AvailableQuantity,
			AfterAvailable:         balance.AvailableQuantity,
			BeforeReserved:         balance.ReservedQuantity,
			AfterReserved:          balance.ReservedQuantity,
			BeforeSellable:         balance.Sellable,
			AfterSellable:          balance.Sellable,
			BeforeSellableReserved: balance.SellableReserved,
			AfterSellableReserved:  balance.SellableReserved,
			BeforeDamaged:          balance.DamagedQuantity,
			AfterDamaged:           balance.DamagedQuantity,
			Remark:                 params.Remark,
			OperatorID:             params.OperatorID,
			OperatedAt:             operatedAt,
			GmtCreate:              time.Now(),
			GmtModified:            time.Now(),
		}

		if err := u.movementRepo.Create(ctx, passMovement); err != nil {
			return fmt.Errorf("failed to create pass movement: %w", err)
		}

		balance.Returned = balance.Returned - params.PassQuantity
		balance.PendingInspection = balance.PendingInspection + params.PassQuantity
	}

	// 质检不合格部分: 退货库存 → 损坏库存
	if params.FailQuantity > 0 {
		failRemark := "质检不合格"
		if params.Remark != nil {
			failRemark = *params.Remark + " - 质检不合格"
		}
		failMovement := &domain.InventoryMovement{
			TraceID:                &traceID,
			ProductID:              params.ProductID,
			WarehouseID:            params.WarehouseID,
			MovementType:           domain.MovementTypeReturnInspect,
			ReferenceType:          params.ReferenceType,
			ReferenceNumber:        params.ReferenceNumber,
			Quantity:               -int(params.FailQuantity),
			BeforeAvailable:        balance.AvailableQuantity,
			AfterAvailable:         balance.AvailableQuantity,
			BeforeReserved:         balance.ReservedQuantity,
			AfterReserved:          balance.ReservedQuantity,
			BeforeSellable:         balance.Sellable,
			AfterSellable:          balance.Sellable,
			BeforeSellableReserved: balance.SellableReserved,
			AfterSellableReserved:  balance.SellableReserved,
			BeforeDamaged:          balance.DamagedQuantity,
			AfterDamaged:           balance.DamagedQuantity + params.FailQuantity,
			Remark:                 &failRemark,
			OperatorID:             params.OperatorID,
			OperatedAt:             operatedAt,
			GmtCreate:              time.Now(),
			GmtModified:            time.Now(),
		}

		if err := u.movementRepo.Create(ctx, failMovement); err != nil {
			return fmt.Errorf("failed to create fail movement: %w", err)
		}

		balance.Returned = balance.Returned - params.FailQuantity
		balance.DamagedQuantity = balance.DamagedQuantity + params.FailQuantity
	}

	// Update balance
	balance.TotalQuantity = calculateBalanceTotal(balance)
	balance.LastMovementAt = &operatedAt
	balance.GmtModified = time.Now()

	if err := u.balanceRepo.Update(ctx, balance); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}
