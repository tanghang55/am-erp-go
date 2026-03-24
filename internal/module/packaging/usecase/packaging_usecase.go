package usecase

import (
	"am-erp-go/internal/module/packaging/domain"
	"errors"
	"fmt"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/validation"
)

var ErrPackagingSupplierRequired = errors.New("supplier is required")
var ErrPackagingItemCodeInvalid = errors.New("packaging item code only supports letters, numbers, hyphen and underscore")
var ErrPackagingItemReferenced = errors.New("packaging item is still referenced by business data")

type PackagingUsecase struct {
	itemRepo   domain.PackagingItemRepository
	ledgerRepo domain.PackagingLedgerRepository
}

func NewPackagingUsecase(
	itemRepo domain.PackagingItemRepository,
	ledgerRepo domain.PackagingLedgerRepository,
) *PackagingUsecase {
	return &PackagingUsecase{
		itemRepo:   itemRepo,
		ledgerRepo: ledgerRepo,
	}
}

// ============= Packaging Item =============

func (uc *PackagingUsecase) ListItems(params *domain.PackagingItemListParams) ([]domain.PackagingItem, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	items, total, err := uc.itemRepo.List(params)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].Deletable = true
		refCount, err := uc.itemRepo.CountReferences(items[i].ID)
		if err != nil {
			return nil, 0, err
		}
		items[i].ReferenceCount = refCount
		if refCount > 0 {
			items[i].Deletable = false
			items[i].DeleteBlockReason = "已被业务数据引用，不可删除"
		}
	}
	return items, total, nil
}

func (uc *PackagingUsecase) GetItem(id uint64) (*domain.PackagingItem, error) {
	item, err := uc.itemRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	item.Deletable = true
	refCount, err := uc.itemRepo.CountReferences(item.ID)
	if err != nil {
		return nil, err
	}
	item.ReferenceCount = refCount
	if refCount > 0 {
		item.Deletable = false
		item.DeleteBlockReason = "已被业务数据引用，不可删除"
	}
	return item, nil
}

func (uc *PackagingUsecase) CreateItem(item *domain.PackagingItem) error {
	if item.SupplierID == nil || *item.SupplierID == 0 {
		return ErrPackagingSupplierRequired
	}
	if !validation.IsValidCode(strings.TrimSpace(item.ItemCode)) {
		return ErrPackagingItemCodeInvalid
	}
	// 生成trace_id
	item.TraceID = fmt.Sprintf("PKG-%d-%d", time.Now().Unix(), item.CreatedBy)
	return uc.itemRepo.Create(item)
}

func (uc *PackagingUsecase) UpdateItem(item *domain.PackagingItem) error {
	if item.SupplierID == nil || *item.SupplierID == 0 {
		return ErrPackagingSupplierRequired
	}
	if !validation.IsValidCode(strings.TrimSpace(item.ItemCode)) {
		return ErrPackagingItemCodeInvalid
	}
	return uc.itemRepo.Update(item)
}

func (uc *PackagingUsecase) DeleteItem(id uint64) error {
	refCount, err := uc.itemRepo.CountReferences(id)
	if err != nil {
		return err
	}
	if refCount > 0 {
		return ErrPackagingItemReferenced
	}
	return uc.itemRepo.Delete(id)
}

func (uc *PackagingUsecase) GetLowStockItems() ([]domain.PackagingItem, error) {
	return uc.itemRepo.GetLowStockItems()
}

// ============= Packaging Ledger =============

func (uc *PackagingUsecase) ListLedgers(params *domain.PackagingLedgerListParams) ([]domain.PackagingLedger, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.ledgerRepo.List(params)
}

func (uc *PackagingUsecase) GetLedger(id uint64) (*domain.PackagingLedger, error) {
	return uc.ledgerRepo.GetByID(id)
}

// CreateInboundLedger 创建入库流水
func (uc *PackagingUsecase) CreateInboundLedger(ledger *domain.PackagingLedger, userID uint64) error {
	// 获取当前库存
	item, err := uc.itemRepo.GetByID(ledger.PackagingItemID)
	if err != nil {
		return fmt.Errorf("包材不存在: %w", err)
	}

	// 设置流水信息
	ledger.TraceID = fmt.Sprintf("PKG-IN-%d-%d", time.Now().Unix(), userID)
	ledger.TransactionType = "IN"
	ledger.QuantityBefore = item.QuantityOnHand
	ledger.QuantityAfter = item.QuantityOnHand + uint64(ledger.Quantity)
	// TotalCost由数据库自动计算
	ledger.CreatedBy = userID

	if ledger.OccurredAt.IsZero() {
		ledger.OccurredAt = time.Now()
	}

	// 创建流水
	if err := uc.ledgerRepo.Create(ledger); err != nil {
		return err
	}

	// 更新库存
	return uc.itemRepo.UpdateQuantity(ledger.PackagingItemID, ledger.Quantity)
}

// CreateOutboundLedger 创建出库流水
func (uc *PackagingUsecase) CreateOutboundLedger(ledger *domain.PackagingLedger, userID uint64) error {
	// 获取当前库存
	item, err := uc.itemRepo.GetByID(ledger.PackagingItemID)
	if err != nil {
		return fmt.Errorf("包材不存在: %w", err)
	}

	// 检查库存是否充足
	if item.QuantityOnHand < uint64(ledger.Quantity) {
		return fmt.Errorf("库存不足，当前库存: %d", item.QuantityOnHand)
	}

	// 设置流水信息
	ledger.TraceID = fmt.Sprintf("PKG-OUT-%d-%d", time.Now().Unix(), userID)
	ledger.TransactionType = "OUT"
	ledger.QuantityBefore = item.QuantityOnHand
	ledger.QuantityAfter = item.QuantityOnHand - uint64(ledger.Quantity)
	ledger.Quantity = -ledger.Quantity // 出库为负数
	// TotalCost由数据库自动计算
	ledger.CreatedBy = userID

	if ledger.OccurredAt.IsZero() {
		ledger.OccurredAt = time.Now()
	}

	// 创建流水
	if err := uc.ledgerRepo.Create(ledger); err != nil {
		return err
	}

	// 更新库存
	return uc.itemRepo.UpdateQuantity(ledger.PackagingItemID, ledger.Quantity)
}

// CreateAdjustmentLedger 创建调整流水
func (uc *PackagingUsecase) CreateAdjustmentLedger(ledger *domain.PackagingLedger, userID uint64) error {
	// 获取当前库存
	item, err := uc.itemRepo.GetByID(ledger.PackagingItemID)
	if err != nil {
		return fmt.Errorf("包材不存在: %w", err)
	}

	// 设置流水信息
	ledger.TraceID = fmt.Sprintf("PKG-ADJ-%d-%d", time.Now().Unix(), userID)
	ledger.TransactionType = "ADJUSTMENT"
	ledger.QuantityBefore = item.QuantityOnHand

	var newQuantity uint64
	if ledger.Quantity > 0 {
		newQuantity = item.QuantityOnHand + uint64(ledger.Quantity)
	} else {
		absQty := uint64(-ledger.Quantity)
		if item.QuantityOnHand < absQty {
			return fmt.Errorf("调整后库存不能为负数")
		}
		newQuantity = item.QuantityOnHand - absQty
	}
	ledger.QuantityAfter = newQuantity
	// TotalCost由数据库自动计算
	ledger.CreatedBy = userID

	if ledger.OccurredAt.IsZero() {
		ledger.OccurredAt = time.Now()
	}

	// 创建流水
	if err := uc.ledgerRepo.Create(ledger); err != nil {
		return err
	}

	// 更新库存
	return uc.itemRepo.UpdateQuantity(ledger.PackagingItemID, ledger.Quantity)
}

func (uc *PackagingUsecase) GetUsageSummary(dateFrom, dateTo *time.Time) ([]domain.UsageSummaryItem, error) {
	return uc.ledgerRepo.GetUsageSummary(dateFrom, dateTo)
}
