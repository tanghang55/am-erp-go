package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"am-erp-go/internal/module/inventory/domain"

	"github.com/google/uuid"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity   = errors.New("invalid quantity")
)

type InventoryUsecase struct {
	balanceRepo  domain.InventoryBalanceRepository
	movementRepo domain.InventoryMovementRepository
}

func NewInventoryUsecase(
	balanceRepo domain.InventoryBalanceRepository,
	movementRepo domain.InventoryMovementRepository,
) *InventoryUsecase {
	return &InventoryUsecase{
		balanceRepo:  balanceRepo,
		movementRepo: movementRepo,
	}
}

// ListMovements 查询库存流水列表
func (u *InventoryUsecase) ListMovements(params *domain.MovementListParams) ([]*domain.InventoryMovement, int64, error) {
	return u.movementRepo.List(params)
}

// GetMovement 获取单条流水详情
func (u *InventoryUsecase) GetMovement(id uint64) (*domain.InventoryMovement, error) {
	return u.movementRepo.GetByID(id)
}

// CreateMovement 创建库存流水（通用方法）
func (u *InventoryUsecase) CreateMovement(ctx context.Context, params *domain.CreateMovementParams) (*domain.InventoryMovement, error) {
	if params.Quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	// Get or create balance record
	balance, err := u.balanceRepo.GetOrCreate(ctx, params.SkuID, params.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// Calculate new quantities based on movement type
	newBalance, err := u.calculateNewBalance(balance, params.MovementType, params.Quantity)
	if err != nil {
		return nil, err
	}

	// Create movement record
	operatedAt := time.Now()
	if params.OperatedAt != nil {
		operatedAt = *params.OperatedAt
	}

	traceID := uuid.New().String()
	movement := &domain.InventoryMovement{
		TraceID:         &traceID,
		SkuID:           params.SkuID,
		WarehouseID:     params.WarehouseID,
		MovementType:    params.MovementType,
		ReferenceType:   params.ReferenceType,
		ReferenceID:     params.ReferenceID,
		ReferenceNumber: params.ReferenceNumber,
		Quantity:        params.Quantity,
		BeforeAvailable: balance.AvailableQuantity,
		AfterAvailable:  newBalance.AvailableQuantity,
		BeforeReserved:  balance.ReservedQuantity,
		AfterReserved:   newBalance.ReservedQuantity,
		BeforeDamaged:   balance.DamagedQuantity,
		AfterDamaged:    newBalance.DamagedQuantity,
		UnitCost:        params.UnitCost,
		Remark:          params.Remark,
		OperatorID:      params.OperatorID,
		OperatedAt:      operatedAt,
		GmtCreate:       time.Now(),
		GmtModified:     time.Now(),
	}

	if params.UnitCost != nil {
		totalCost := *params.UnitCost * float64(absInt(params.Quantity))
		movement.TotalCost = &totalCost
	}

	// Save movement
	if err := u.movementRepo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("failed to create movement: %w", err)
	}

	// Update balance
	balance.AvailableQuantity = newBalance.AvailableQuantity
	balance.ReservedQuantity = newBalance.ReservedQuantity
	balance.DamagedQuantity = newBalance.DamagedQuantity
	balance.PurchasingInTransit = newBalance.PurchasingInTransit
	balance.PendingInspection = newBalance.PendingInspection
	balance.RawMaterial = newBalance.RawMaterial
	balance.PendingShipment = newBalance.PendingShipment
	balance.LogisticsInTransit = newBalance.LogisticsInTransit
	balance.Sellable = newBalance.Sellable
	balance.Returned = newBalance.Returned
	balance.TotalQuantity = balance.AvailableQuantity + balance.ReservedQuantity + balance.DamagedQuantity +
		balance.PurchasingInTransit + balance.PendingInspection + balance.RawMaterial +
		balance.PendingShipment + balance.LogisticsInTransit + balance.Sellable + balance.Returned
	balance.LastMovementAt = &operatedAt
	balance.GmtModified = time.Now()

	if err := u.balanceRepo.Update(ctx, balance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
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
		SkuID:           params.SkuID,
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
		SkuID:           params.SkuID,
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

// GetSkuBalance 获取SKU在指定仓库的库存
func (u *InventoryUsecase) GetSkuBalance(skuID, warehouseID uint64) (*domain.InventoryBalance, error) {
	return u.balanceRepo.GetBySkuAndWarehouse(skuID, warehouseID)
}

// calculateNewBalance 根据流水类型计算新的库存数量
func (u *InventoryUsecase) calculateNewBalance(balance *domain.InventoryBalance, movementType domain.MovementType, quantity int) (*domain.InventoryBalance, error) {
	newBalance := &domain.InventoryBalance{
		AvailableQuantity:    balance.AvailableQuantity,
		ReservedQuantity:     balance.ReservedQuantity,
		DamagedQuantity:      balance.DamagedQuantity,
		PurchasingInTransit:  balance.PurchasingInTransit,
		PendingInspection:    balance.PendingInspection,
		RawMaterial:          balance.RawMaterial,
		PendingShipment:      balance.PendingShipment,
		LogisticsInTransit:   balance.LogisticsInTransit,
		Sellable:             balance.Sellable,
		Returned:             balance.Returned,
	}

	switch movementType {
	case domain.MovementTypePurchaseReceipt,
		domain.MovementTypeReturnReceipt,
		domain.MovementTypeTransferIn:
		// Increase available quantity
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		newBalance.AvailableQuantity = balance.AvailableQuantity + uint(quantity)

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

	case domain.MovementTypeAssemblyComplete:
		// 组装完成: 原料库存 → 待出库存
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.RawMaterial < decrease {
			return nil, fmt.Errorf("原料库存不足")
		}
		newBalance.RawMaterial = balance.RawMaterial - decrease
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
		// 发货单发货: 待出 → 在途
		if quantity < 0 {
			return nil, ErrInvalidQuantity
		}
		decrease := uint(quantity)
		if balance.PendingShipment < decrease {
			return nil, fmt.Errorf("待出库存不足")
		}
		newBalance.PendingShipment = balance.PendingShipment - decrease
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

// RecordReturnReceive 退货入库
func (u *InventoryUsecase) RecordReturnReceive(ctx context.Context, params *domain.StockTransitionParams) (*domain.InventoryMovement, error) {
	if params.Quantity == 0 {
		return nil, ErrInvalidQuantity
	}

	// Get or create balance record
	balance, err := u.balanceRepo.GetOrCreate(ctx, params.SkuID, params.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	// 退货入库: 增加退货库存
	operatedAt := time.Now()
	traceID := uuid.New().String()
	movement := &domain.InventoryMovement{
		TraceID:         &traceID,
		SkuID:           params.SkuID,
		WarehouseID:     params.WarehouseID,
		MovementType:    domain.MovementTypeReturnReceipt,
		ReferenceType:   params.ReferenceType,
		ReferenceID:     params.ReferenceID,
		ReferenceNumber: params.ReferenceNumber,
		Quantity:        int(params.Quantity),
		BeforeAvailable: balance.AvailableQuantity,
		AfterAvailable:  balance.AvailableQuantity,
		BeforeReserved:  balance.ReservedQuantity,
		AfterReserved:   balance.ReservedQuantity,
		BeforeDamaged:   balance.DamagedQuantity,
		AfterDamaged:    balance.DamagedQuantity,
		UnitCost:        params.UnitCost,
		Remark:          params.Remark,
		OperatorID:      params.OperatorID,
		OperatedAt:      operatedAt,
		GmtCreate:       time.Now(),
		GmtModified:     time.Now(),
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
	balance.TotalQuantity = balance.AvailableQuantity + balance.ReservedQuantity + balance.DamagedQuantity +
		balance.PurchasingInTransit + balance.PendingInspection + balance.RawMaterial +
		balance.PendingShipment + balance.LogisticsInTransit + balance.Sellable + balance.Returned
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
	balance, err := u.balanceRepo.GetOrCreate(ctx, params.SkuID, params.WarehouseID)
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
			TraceID:         &traceID,
			SkuID:           params.SkuID,
			WarehouseID:     params.WarehouseID,
			MovementType:    domain.MovementTypeReturnInspect,
			ReferenceType:   params.ReferenceType,
			ReferenceNumber: params.ReferenceNumber,
			Quantity:        int(params.PassQuantity),
			BeforeAvailable: balance.AvailableQuantity,
			AfterAvailable:  balance.AvailableQuantity,
			BeforeReserved:  balance.ReservedQuantity,
			AfterReserved:   balance.ReservedQuantity,
			BeforeDamaged:   balance.DamagedQuantity,
			AfterDamaged:    balance.DamagedQuantity,
			Remark:          params.Remark,
			OperatorID:      params.OperatorID,
			OperatedAt:      operatedAt,
			GmtCreate:       time.Now(),
			GmtModified:     time.Now(),
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
			TraceID:         &traceID,
			SkuID:           params.SkuID,
			WarehouseID:     params.WarehouseID,
			MovementType:    domain.MovementTypeReturnInspect,
			ReferenceType:   params.ReferenceType,
			ReferenceNumber: params.ReferenceNumber,
			Quantity:        -int(params.FailQuantity),
			BeforeAvailable: balance.AvailableQuantity,
			AfterAvailable:  balance.AvailableQuantity,
			BeforeReserved:  balance.ReservedQuantity,
			AfterReserved:   balance.ReservedQuantity,
			BeforeDamaged:   balance.DamagedQuantity,
			AfterDamaged:    balance.DamagedQuantity + params.FailQuantity,
			Remark:          &failRemark,
			OperatorID:      params.OperatorID,
			OperatedAt:      operatedAt,
			GmtCreate:       time.Now(),
			GmtModified:     time.Now(),
		}

		if err := u.movementRepo.Create(ctx, failMovement); err != nil {
			return fmt.Errorf("failed to create fail movement: %w", err)
		}

		balance.Returned = balance.Returned - params.FailQuantity
		balance.DamagedQuantity = balance.DamagedQuantity + params.FailQuantity
	}

	// Update balance
	balance.TotalQuantity = balance.AvailableQuantity + balance.ReservedQuantity + balance.DamagedQuantity +
		balance.PurchasingInTransit + balance.PendingInspection + balance.RawMaterial +
		balance.PendingShipment + balance.LogisticsInTransit + balance.Sellable + balance.Returned
	balance.LastMovementAt = &operatedAt
	balance.GmtModified = time.Now()

	if err := u.balanceRepo.Update(ctx, balance); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}
