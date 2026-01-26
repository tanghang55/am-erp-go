package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	mathrand "math/rand"
	"sort"
	"time"

	"am-erp-go/internal/module/procurement/domain"
	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	productdomain "am-erp-go/internal/module/product/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

var (
	ErrPurchaseOrderInvalid             = errors.New("invalid order")
	ErrPurchaseOrderMissingItems        = errors.New("missing items")
	ErrPurchaseOrderComboProviderNeeded = errors.New("combo provider required")
	ErrPurchaseOrderComboNoComponents   = errors.New("combo product cannot be purchased directly")
)

type PurchaseOrderRepository interface {
	List(params *domain.PurchaseOrderListParams) ([]domain.PurchaseOrder, int64, error)
	GetByID(id uint64) (*domain.PurchaseOrder, error)
	Create(order *domain.PurchaseOrder) error
	Update(order *domain.PurchaseOrder) error
	Delete(id uint64) error
}

type ProductLookup interface {
	ListByIDs(ids []uint64) ([]productdomain.Product, error)
}

type ComboProvider interface {
	GetItemsByComboID(comboID uint64) ([]productdomain.ProductComboItem, error)
}

type InventoryService interface {
	CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error)
}

type AuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

type PurchaseOrderUsecase struct {
	repo             PurchaseOrderRepository
	productLookup    ProductLookup
	comboProvider    ComboProvider
	inventoryService InventoryService
	auditLogger      AuditLogger
}

func NewPurchaseOrderUsecase(
	repo PurchaseOrderRepository,
	productLookup ProductLookup,
	comboProvider ComboProvider,
	inventoryService InventoryService,
	auditLogger AuditLogger,
) *PurchaseOrderUsecase {
	return &PurchaseOrderUsecase{
		repo:             repo,
		productLookup:    productLookup,
		comboProvider:    comboProvider,
		inventoryService: inventoryService,
		auditLogger:      auditLogger,
	}
}

func (uc *PurchaseOrderUsecase) Create(c *gin.Context, order *domain.PurchaseOrder) (*domain.PurchaseOrder, error) {
	if order == nil {
		return nil, ErrPurchaseOrderInvalid
	}
	if len(order.Items) == 0 {
		return nil, ErrPurchaseOrderMissingItems
	}

	order.Currency = defaultCurrency(order.Currency)
	expandedItems, err := uc.expandComboItems(order.Currency, order.Items)
	if err != nil {
		return nil, err
	}
	if len(expandedItems) == 0 {
		return nil, ErrPurchaseOrderMissingItems
	}

	order.Items = expandedItems
	order.Status = domain.PurchaseOrderStatusDraft
	order.PoNumber = ensurePoNumber(order.PoNumber)
	order.TotalAmount = sumOrderTotal(order.Items)

	if err := uc.repo.Create(order); err != nil {
		return nil, err
	}

	uc.recordAudit(c, "CREATE", "PurchaseOrder", fmt.Sprintf("%d", order.ID), nil, order)
	return order, nil
}

func (uc *PurchaseOrderUsecase) Receive(c *gin.Context, orderID uint64, params domain.PurchaseOrderReceiveParams) error {
	if orderID == 0 {
		return errors.New("invalid order id")
	}
	if params.WarehouseID == 0 {
		return errors.New("missing warehouse")
	}

	order, err := uc.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != domain.PurchaseOrderStatusShipped {
		return errors.New("order status invalid for receive")
	}

	now := time.Now()
	allReceived := true
	ctx := context.Background()

	for i := range order.Items {
		item := &order.Items[i]
		delta := params.ReceivedQties[item.ID]
		if delta > 0 {
			if item.QtyReceived+delta > item.QtyOrdered {
				return errors.New("received qty exceeds ordered")
			}
			item.QtyReceived += delta

			// 仓库收货：从"采购在途"转到"待检"
			// 注意：发货时(MarkShipped)已经增加了"采购在途"库存
			if uc.inventoryService != nil {
				unitCost := item.UnitCost
				remark := fmt.Sprintf("仓库收货入待检: %s", order.PoNumber)

				warehouseReceiveParams := &inventoryDomain.CreateMovementParams{
					SkuID:           item.SkuID,
					WarehouseID:     params.WarehouseID,
					MovementType:    inventoryDomain.MovementTypeWarehouseReceive,
					Quantity:        int(delta),
					ReferenceType:   stringPtr("PURCHASE_ORDER"),
					ReferenceID:     &orderID,
					ReferenceNumber: &order.PoNumber,
					UnitCost:        &unitCost,
					Remark:          &remark,
					OperatorID:      params.OperatorID,
					OperatedAt:      &now,
				}
				if _, err := uc.inventoryService.CreateMovement(ctx, warehouseReceiveParams); err != nil {
					return fmt.Errorf("failed to create warehouse receive movement: %w", err)
				}
			}
		}

		if item.QtyReceived < item.QtyOrdered {
			allReceived = false
		}
	}

	if allReceived {
		order.Status = domain.PurchaseOrderStatusReceived
		order.ReceivedAt = &now
	}

	if err := uc.repo.Update(order); err != nil {
		return err
	}

	uc.recordAudit(c, "RECEIVE", "PurchaseOrder", fmt.Sprintf("%d", orderID), nil, order)
	return nil
}

func stringPtr(s string) *string {
	return &s
}

func (uc *PurchaseOrderUsecase) List(params *domain.PurchaseOrderListParams) ([]domain.PurchaseOrder, int64, error) {
	if params == nil {
		params = &domain.PurchaseOrderListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.repo.List(params)
}

func (uc *PurchaseOrderUsecase) Get(orderID uint64) (*domain.PurchaseOrder, error) {
	return uc.repo.GetByID(orderID)
}

func (uc *PurchaseOrderUsecase) Update(c *gin.Context, orderID uint64, updates *domain.PurchaseOrder) (*domain.PurchaseOrder, error) {
	if orderID == 0 {
		return nil, errors.New("invalid order id")
	}
	if updates == nil {
		return nil, ErrPurchaseOrderInvalid
	}

	existing, err := uc.repo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if existing.Status != domain.PurchaseOrderStatusDraft {
		return nil, errors.New("only draft can be updated")
	}

	currency := defaultCurrency(updates.Currency)
	if currency == "" {
		currency = existing.Currency
	}
	updates.Currency = currency

	expandedItems, err := uc.expandComboItems(currency, updates.Items)
	if err != nil {
		return nil, err
	}
	if len(expandedItems) == 0 {
		return nil, ErrPurchaseOrderMissingItems
	}

	existing.SupplierID = updates.SupplierID
	existing.Marketplace = updates.Marketplace
	existing.Currency = updates.Currency
	existing.Remark = updates.Remark
	existing.Items = expandedItems
	existing.TotalAmount = sumOrderTotal(existing.Items)

	if err := uc.repo.Update(existing); err != nil {
		return nil, err
	}

	uc.recordAudit(c, "UPDATE", "PurchaseOrder", fmt.Sprintf("%d", orderID), nil, existing)
	return existing, nil
}

func (uc *PurchaseOrderUsecase) Delete(c *gin.Context, orderID uint64) error {
	order, err := uc.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != domain.PurchaseOrderStatusDraft {
		return errors.New("only draft can be deleted")
	}
	if err := uc.repo.Delete(orderID); err != nil {
		return err
	}
	uc.recordAudit(c, "DELETE", "PurchaseOrder", fmt.Sprintf("%d", orderID), order, nil)
	return nil
}

func (uc *PurchaseOrderUsecase) Submit(c *gin.Context, orderID uint64) error {
	order, err := uc.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != domain.PurchaseOrderStatusDraft {
		return errors.New("order status invalid for submit")
	}
	now := time.Now()
	order.Status = domain.PurchaseOrderStatusOrdered
	order.OrderedAt = &now
	order.Items = nil
	if err := uc.repo.Update(order); err != nil {
		return err
	}
	uc.recordAudit(c, "SUBMIT", "PurchaseOrder", fmt.Sprintf("%d", orderID), nil, order)
	return nil
}

func (uc *PurchaseOrderUsecase) MarkShipped(c *gin.Context, orderID uint64, params domain.PurchaseOrderShipParams) error {
	order, err := uc.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != domain.PurchaseOrderStatusOrdered {
		return errors.New("order status invalid for ship")
	}

	now := time.Now()
	ctx := context.Background()

	// 供应商发货 → 货物进入"采购在途"状态
	if uc.inventoryService != nil && params.WarehouseID != 0 {
		for _, item := range order.Items {
			if item.QtyOrdered == 0 {
				continue
			}
			unitCost := item.UnitCost
			purchaseShipParams := &inventoryDomain.CreateMovementParams{
				SkuID:           item.SkuID,
				WarehouseID:     params.WarehouseID,
				MovementType:    inventoryDomain.MovementTypePurchaseShip,
				Quantity:        int(item.QtyOrdered),
				ReferenceType:   stringPtr("PURCHASE_ORDER"),
				ReferenceID:     &orderID,
				ReferenceNumber: &order.PoNumber,
				UnitCost:        &unitCost,
				Remark:          stringPtr(fmt.Sprintf("采购发货入在途: %s", order.PoNumber)),
				OperatorID:      params.OperatorID,
				OperatedAt:      &now,
			}
			if _, err := uc.inventoryService.CreateMovement(ctx, purchaseShipParams); err != nil {
				return fmt.Errorf("failed to create purchase ship movement: %w", err)
			}
		}
	}

	order.Status = domain.PurchaseOrderStatusShipped
	order.ShippedAt = &now
	order.Items = nil
	if err := uc.repo.Update(order); err != nil {
		return err
	}
	uc.recordAudit(c, "SHIP", "PurchaseOrder", fmt.Sprintf("%d", orderID), nil, order)
	return nil
}

func (uc *PurchaseOrderUsecase) Close(c *gin.Context, orderID uint64) error {
	order, err := uc.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status == domain.PurchaseOrderStatusClosed {
		return nil
	}
	order.Status = domain.PurchaseOrderStatusClosed
	order.Items = nil
	if err := uc.repo.Update(order); err != nil {
		return err
	}
	uc.recordAudit(c, "CLOSE", "PurchaseOrder", fmt.Sprintf("%d", orderID), nil, order)
	return nil
}

func (uc *PurchaseOrderUsecase) expandComboItems(currency string, items []domain.PurchaseOrderItem) ([]domain.PurchaseOrderItem, error) {
	if len(items) == 0 {
		return []domain.PurchaseOrderItem{}, nil
	}

	if uc.productLookup == nil {
		return normalizeItems(currency, items), nil
	}

	productIDs := uniqueSkuIDs(items)
	products, err := uc.productLookup.ListByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uint64]productdomain.Product, len(products))
	for _, product := range products {
		productMap[product.ID] = product
	}

	comboExpanded := make(map[uint64]bool)
	for _, item := range items {
		product, ok := productMap[item.SkuID]
		if !ok || product.ComboID == nil {
			continue
		}
		if product.IsComboMain != 1 {
			comboExpanded[*product.ComboID] = true
		}
	}

	result := make(map[uint64]*domain.PurchaseOrderItem, len(items))
	for _, item := range items {
		product, ok := productMap[item.SkuID]
		if ok && product.ComboID != nil && product.IsComboMain == 1 {
			if comboExpanded[*product.ComboID] {
				mergeItem(result, item.SkuID, item.QtyOrdered, item.UnitCost, currency)
				continue
			}
			if uc.comboProvider == nil {
				return nil, ErrPurchaseOrderComboProviderNeeded
			}
			comboItems, err := uc.comboProvider.GetItemsByComboID(*product.ComboID)
			if err != nil {
				return nil, err
			}
			hasComponent := false
			missingIDs := make([]uint64, 0, len(comboItems))
			for _, comboItem := range comboItems {
				if comboItem.ProductID == 0 {
					continue
				}
				if comboItem.ProductID != comboItem.MainProductID {
					hasComponent = true
				}
				if _, exists := productMap[comboItem.ProductID]; !exists {
					missingIDs = append(missingIDs, comboItem.ProductID)
				}
			}
			if !hasComponent {
				return nil, ErrPurchaseOrderComboNoComponents
			}
			if len(missingIDs) > 0 {
				extraProducts, err := uc.productLookup.ListByIDs(missingIDs)
				if err != nil {
					return nil, err
				}
				for _, extra := range extraProducts {
					productMap[extra.ID] = extra
				}
			}
			for _, comboItem := range comboItems {
				qty := item.QtyOrdered * comboItem.QtyRatio
				if qty == 0 {
					continue
				}
				unitCost := getUnitCost(productMap[comboItem.ProductID])
				mergeItem(result, comboItem.ProductID, qty, unitCost, currency)
			}
			continue
		}

		mergeItem(result, item.SkuID, item.QtyOrdered, item.UnitCost, currency)
	}

	return flattenItems(result), nil
}

func mergeItem(target map[uint64]*domain.PurchaseOrderItem, skuID, qty uint64, unitCost float64, currency string) {
	if qty == 0 {
		return
	}
	if existing, ok := target[skuID]; ok {
		existing.QtyOrdered += qty
		if existing.UnitCost == 0 && unitCost > 0 {
			existing.UnitCost = unitCost
		}
		existing.Subtotal = roundAmount(existing.UnitCost * float64(existing.QtyOrdered))
		return
	}

	item := &domain.PurchaseOrderItem{
		SkuID:       skuID,
		QtyOrdered:  qty,
		QtyReceived: 0,
		UnitCost:    unitCost,
		Currency:    currency,
	}
	item.Subtotal = roundAmount(unitCost * float64(qty))
	target[skuID] = item
}

func flattenItems(items map[uint64]*domain.PurchaseOrderItem) []domain.PurchaseOrderItem {
	result := make([]domain.PurchaseOrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, *item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SkuID < result[j].SkuID
	})
	return result
}

func normalizeItems(currency string, items []domain.PurchaseOrderItem) []domain.PurchaseOrderItem {
	normalized := make([]domain.PurchaseOrderItem, 0, len(items))
	for _, item := range items {
		item.Currency = currency
		item.Subtotal = roundAmount(item.UnitCost * float64(item.QtyOrdered))
		normalized = append(normalized, item)
	}
	return normalized
}

func sumOrderTotal(items []domain.PurchaseOrderItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Subtotal
	}
	return roundAmount(total)
}

func defaultCurrency(currency string) string {
	if currency == "" {
		return "USD"
	}
	return currency
}

func ensurePoNumber(poNumber string) string {
	if poNumber != "" {
		return poNumber
	}
	now := time.Now()
	// PO + 年月日时分秒 + 毫秒(3位) + 随机数(2位)
	// 例如：PO2026012314140012345
	ms := now.Nanosecond() / 1000000 // 毫秒 0-999
	rnd := mathrand.Intn(100)        // 随机数 0-99
	return fmt.Sprintf("PO%s%03d%02d", now.Format("20060102150405"), ms, rnd)
}

func uniqueSkuIDs(items []domain.PurchaseOrderItem) []uint64 {
	seen := make(map[uint64]struct{}, len(items))
	result := make([]uint64, 0, len(items))
	for _, item := range items {
		if item.SkuID == 0 {
			continue
		}
		if _, ok := seen[item.SkuID]; ok {
			continue
		}
		seen[item.SkuID] = struct{}{}
		result = append(result, item.SkuID)
	}
	return result
}

func getUnitCost(product productdomain.Product) float64 {
	if product.UnitCost == nil {
		return 0
	}
	return *product.UnitCost
}

func roundAmount(v float64) float64 {
	return math.Round(v*10000) / 10000
}

func (uc *PurchaseOrderUsecase) recordAudit(c *gin.Context, action, entityType, entityID string, before, after any) {
	if uc.auditLogger == nil || c == nil {
		return
	}
	_ = uc.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Procurement",
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
}
