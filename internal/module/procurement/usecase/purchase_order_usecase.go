package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/numbering"
	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/procurement/domain"
	productdomain "am-erp-go/internal/module/product/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

var (
	ErrPurchaseOrderInvalid               = errors.New("invalid order")
	ErrPurchaseOrderMissingItems          = errors.New("missing items")
	ErrPurchaseOrderMissingSupplier       = errors.New("purchase order item missing supplier")
	ErrPurchaseOrderMissingProduct        = errors.New("purchase order item missing product")
	ErrPurchaseOrderInvalidQty            = errors.New("purchase order item qty_ordered must be greater than 0")
	ErrPurchaseOrderInvalidUnitCost       = errors.New("purchase order item unit_cost must be greater than 0")
	ErrPurchaseOrderSplitRequired         = errors.New("purchase order must be split by supplier")
	ErrPurchaseOrderComboProviderNeeded   = errors.New("combo provider required")
	ErrPurchaseOrderComboNoComponents     = errors.New("combo product cannot be purchased directly")
	ErrPurchaseOrderInvalidCompleteStatus = errors.New("order status invalid for complete")
	ErrPurchaseOrderPendingInspection     = errors.New("pending inspection remains")
	ErrPurchaseOrderIncompleteReceipt     = errors.New("receipt quantity incomplete")
	ErrPurchaseOrderInspectionFailed      = errors.New("inspection loss remains")
	ErrPurchaseOrderForceCompleteReason   = errors.New("force complete reason is required")
)

type PurchaseOrderRepository interface {
	List(params *domain.PurchaseOrderListParams) ([]domain.PurchaseOrder, int64, error)
	GetByID(id uint64) (*domain.PurchaseOrder, error)
	Create(order *domain.PurchaseOrder) error
	Update(order *domain.PurchaseOrder) error
	UpdateProgress(order *domain.PurchaseOrder) error
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

type ReplenishmentPlanCleaner interface {
	DeletePlansByPurchaseOrderID(purchaseOrderID uint64) error
}

type PurchaseOrderCostEventType string

const (
	PurchaseOrderCostEventOrdered  PurchaseOrderCostEventType = "PO_ORDERED"
	PurchaseOrderCostEventShipped  PurchaseOrderCostEventType = "PO_SHIPPED"
	PurchaseOrderCostEventReceived PurchaseOrderCostEventType = "PO_RECEIVED"
)

type PurchaseOrderCostEventParams struct {
	EventType           PurchaseOrderCostEventType
	PurchaseOrderID     uint64
	PurchaseOrderItemID uint64
	ProductID           uint64
	WarehouseID         *uint64
	Marketplace         string
	Currency            string
	UnitCost            float64
	QtyEvent            uint64
	OccurredAt          time.Time
	OperatorID          *uint64
}

type PurchaseOrderCostEventRecorder interface {
	RecordPurchaseOrderEvent(params *PurchaseOrderCostEventParams) error
}

type PurchaseOrderShipTransactionalDeps struct {
	Repo              PurchaseOrderRepository
	InventoryService  InventoryService
	CostEventRecorder PurchaseOrderCostEventRecorder
}

type PurchaseOrderShipTransactionManager interface {
	Run(ctx context.Context, fn func(PurchaseOrderShipTransactionalDeps) error) error
}

type PurchaseOrderReceiveTransactionalDeps struct {
	Repo              PurchaseOrderRepository
	InventoryService  InventoryService
	CostEventRecorder PurchaseOrderCostEventRecorder
}

type PurchaseOrderReceiveTransactionManager interface {
	Run(ctx context.Context, fn func(PurchaseOrderReceiveTransactionalDeps) error) error
}

type PurchaseOrderInspectTransactionalDeps struct {
	Repo             PurchaseOrderRepository
	InventoryService InventoryService
}

type PurchaseOrderInspectTransactionManager interface {
	Run(ctx context.Context, fn func(PurchaseOrderInspectTransactionalDeps) error) error
}

type PurchaseOrderSubmitTransactionalDeps struct {
	Repo              PurchaseOrderRepository
	PlanCleaner       ReplenishmentPlanCleaner
	CostEventRecorder PurchaseOrderCostEventRecorder
}

type PurchaseOrderSubmitTransactionManager interface {
	Run(ctx context.Context, fn func(PurchaseOrderSubmitTransactionalDeps) error) error
}

type PurchaseOrderDefaultsProvider interface {
	GetDefaultBaseCurrency() string
}

type PurchaseOrderUsecase struct {
	repo              PurchaseOrderRepository
	productLookup     ProductLookup
	comboProvider     ComboProvider
	inventoryService  InventoryService
	auditLogger       AuditLogger
	planCleaner       ReplenishmentPlanCleaner
	costEventRecorder PurchaseOrderCostEventRecorder
	defaultsProvider  PurchaseOrderDefaultsProvider
	submitTxManager   PurchaseOrderSubmitTransactionManager
	shipTxManager     PurchaseOrderShipTransactionManager
	receiveTxManager  PurchaseOrderReceiveTransactionManager
	inspectTxManager  PurchaseOrderInspectTransactionManager
}

type preparedPurchaseOrderItem struct {
	item       domain.PurchaseOrderItem
	supplierID uint64
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

func (uc *PurchaseOrderUsecase) BindPlanCleaner(cleaner ReplenishmentPlanCleaner) {
	uc.planCleaner = cleaner
}

func (uc *PurchaseOrderUsecase) BindCostEventRecorder(recorder PurchaseOrderCostEventRecorder) {
	uc.costEventRecorder = recorder
}

func (uc *PurchaseOrderUsecase) BindDefaultsProvider(provider PurchaseOrderDefaultsProvider) {
	uc.defaultsProvider = provider
}

func (uc *PurchaseOrderUsecase) BindShipTransactionManager(manager PurchaseOrderShipTransactionManager) {
	uc.shipTxManager = manager
}

func (uc *PurchaseOrderUsecase) BindSubmitTransactionManager(manager PurchaseOrderSubmitTransactionManager) {
	uc.submitTxManager = manager
}

func (uc *PurchaseOrderUsecase) BindReceiveTransactionManager(manager PurchaseOrderReceiveTransactionManager) {
	uc.receiveTxManager = manager
}

func (uc *PurchaseOrderUsecase) BindInspectTransactionManager(manager PurchaseOrderInspectTransactionManager) {
	uc.inspectTxManager = manager
}

func (uc *PurchaseOrderUsecase) Create(c *gin.Context, order *domain.PurchaseOrder) (*domain.PurchaseOrder, error) {
	if order == nil {
		return nil, ErrPurchaseOrderInvalid
	}
	if len(order.Items) == 0 {
		return nil, ErrPurchaseOrderMissingItems
	}

	order.Currency = uc.defaultCurrency(order.Currency)
	expandedItems, err := uc.expandComboItems(order.Currency, order.Items)
	if err != nil {
		return nil, err
	}
	if len(expandedItems) == 0 {
		return nil, ErrPurchaseOrderMissingItems
	}

	splitOrders, err := uc.splitOrderBySupplier(order, expandedItems)
	if err != nil {
		return nil, err
	}
	if len(splitOrders) != 1 {
		return nil, ErrPurchaseOrderSplitRequired
	}

	order = splitOrders[0]
	order.Status = domain.PurchaseOrderStatusDraft
	order.PoNumber = ensurePoNumber(order.PoNumber)
	order.BatchNo = ""
	order.TotalAmount = sumOrderTotal(order.Items)

	if err := uc.repo.Create(order); err != nil {
		return nil, err
	}
	created, err := uc.repo.GetByID(order.ID)
	if err != nil {
		return nil, err
	}
	uc.recordPurchaseOrderAuditIfChanged(c, "CREATE", order.ID, nil, created)
	return created, nil
}

func (uc *PurchaseOrderUsecase) CreateBatch(c *gin.Context, orders []*domain.PurchaseOrder) ([]domain.PurchaseOrder, error) {
	if len(orders) == 0 {
		return nil, ErrPurchaseOrderMissingItems
	}

	normalized, err := uc.prepareBatchOrders(orders)
	if err != nil {
		return nil, err
	}
	if len(normalized) == 0 {
		return nil, ErrPurchaseOrderMissingItems
	}

	batchNo := numbering.Generate("PO", time.Now())
	for index := range normalized {
		normalized[index].BatchNo = batchNo
		if len(normalized) == 1 {
			normalized[index].PoNumber = batchNo
			continue
		}
		normalized[index].PoNumber = fmt.Sprintf("%s-%d", batchNo, index+1)
	}

	result := make([]domain.PurchaseOrder, 0, len(normalized))
	for _, order := range normalized {
		order.Status = domain.PurchaseOrderStatusDraft
		order.TotalAmount = sumOrderTotal(order.Items)
		if err := uc.repo.Create(order); err != nil {
			return nil, err
		}
		created, err := uc.repo.GetByID(order.ID)
		if err != nil {
			return nil, err
		}
		created.SourcePlanIDs = cloneUint64Slice(order.SourcePlanIDs)
		uc.recordPurchaseOrderAuditIfChanged(c, "CREATE", order.ID, nil, created)
		result = append(result, *created)
	}
	return result, nil
}

func (uc *PurchaseOrderUsecase) prepareBatchOrders(orders []*domain.PurchaseOrder) ([]*domain.PurchaseOrder, error) {
	type batchKey struct {
		supplierID  uint64
		warehouseID uint64
		marketplace string
		currency    string
		remark      string
	}

	grouped := map[batchKey]*domain.PurchaseOrder{}
	keys := make([]batchKey, 0, len(orders))

	for _, order := range orders {
		if order == nil {
			continue
		}
		order.Currency = uc.defaultCurrency(order.Currency)
		expandedItems, err := uc.expandComboItems(order.Currency, order.Items)
		if err != nil {
			return nil, err
		}
		splitOrders, err := uc.splitOrderBySupplier(order, expandedItems)
		if err != nil {
			return nil, err
		}
		for _, splitOrder := range splitOrders {
			if splitOrder == nil || len(splitOrder.Items) == 0 || splitOrder.SupplierID == nil {
				continue
			}
			key := batchKey{
				supplierID:  *splitOrder.SupplierID,
				marketplace: splitOrder.Marketplace,
				currency:    splitOrder.Currency,
				remark:      splitOrder.Remark,
			}
			if splitOrder.WarehouseID != nil {
				key.warehouseID = *splitOrder.WarehouseID
			}
			group := grouped[key]
			if group == nil {
				supplierID := *splitOrder.SupplierID
				group = &domain.PurchaseOrder{
					SupplierID:    &supplierID,
					WarehouseID:   splitOrder.WarehouseID,
					Marketplace:   splitOrder.Marketplace,
					Currency:      splitOrder.Currency,
					Remark:        splitOrder.Remark,
					CreatedBy:     splitOrder.CreatedBy,
					UpdatedBy:     splitOrder.UpdatedBy,
					Items:         []domain.PurchaseOrderItem{},
					SourcePlanIDs: []uint64{},
				}
				grouped[key] = group
				keys = append(keys, key)
			}
			group.Items = append(group.Items, clonePurchaseOrderItems(splitOrder.Items)...)
			group.SourcePlanIDs = appendUniqueUint64s(group.SourcePlanIDs, sourcePlanIDsFromItems(splitOrder.Items)...)
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].supplierID == keys[j].supplierID {
			if keys[i].warehouseID == keys[j].warehouseID {
				if keys[i].marketplace == keys[j].marketplace {
					if keys[i].currency == keys[j].currency {
						return keys[i].remark < keys[j].remark
					}
					return keys[i].currency < keys[j].currency
				}
				return keys[i].marketplace < keys[j].marketplace
			}
			return keys[i].warehouseID < keys[j].warehouseID
		}
		return keys[i].supplierID < keys[j].supplierID
	})

	result := make([]*domain.PurchaseOrder, 0, len(keys))
	for _, key := range keys {
		order := grouped[key]
		if order == nil || len(order.Items) == 0 {
			continue
		}
		result = append(result, order)
	}
	return result, nil
}

func (uc *PurchaseOrderUsecase) splitOrderBySupplier(order *domain.PurchaseOrder, expandedItems []domain.PurchaseOrderItem) ([]*domain.PurchaseOrder, error) {
	preparedItems, err := uc.resolvePreparedItems(order, expandedItems)
	if err != nil {
		return nil, err
	}
	if len(preparedItems) == 0 {
		return []*domain.PurchaseOrder{}, nil
	}

	grouped := map[uint64]*domain.PurchaseOrder{}
	keys := make([]uint64, 0, len(preparedItems))
	for _, prepared := range preparedItems {
		group := grouped[prepared.supplierID]
		if group == nil {
			supplierID := prepared.supplierID
			group = &domain.PurchaseOrder{
				SupplierID:  &supplierID,
				WarehouseID: order.WarehouseID,
				Marketplace: order.Marketplace,
				Currency:    order.Currency,
				Remark:      order.Remark,
				CreatedBy:   order.CreatedBy,
				UpdatedBy:   order.UpdatedBy,
				Items:       []domain.PurchaseOrderItem{},
			}
			grouped[prepared.supplierID] = group
			keys = append(keys, prepared.supplierID)
		}
		group.Items = append(group.Items, prepared.item)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	result := make([]*domain.PurchaseOrder, 0, len(keys))
	for _, supplierID := range keys {
		result = append(result, grouped[supplierID])
	}
	return result, nil
}

func (uc *PurchaseOrderUsecase) resolvePreparedItems(order *domain.PurchaseOrder, expandedItems []domain.PurchaseOrderItem) ([]preparedPurchaseOrderItem, error) {
	if len(expandedItems) == 0 {
		return []preparedPurchaseOrderItem{}, nil
	}

	productIDs := uniqueProductIDs(expandedItems)
	productMap := map[uint64]productdomain.Product{}
	if uc.productLookup != nil && len(productIDs) > 0 {
		products, err := uc.productLookup.ListByIDs(productIDs)
		if err != nil {
			return nil, err
		}
		for _, product := range products {
			productMap[product.ID] = product
		}
	}

	result := make([]preparedPurchaseOrderItem, 0, len(expandedItems))
	for _, item := range expandedItems {
		if err := validatePurchaseOrderItem(item); err != nil {
			return nil, err
		}
		supplierID, err := resolvePurchaseOrderItemSupplier(order, item, productMap[item.ProductID])
		if err != nil {
			return nil, err
		}
		item.Currency = order.Currency
		item.Subtotal = roundAmount(item.UnitCost * float64(item.QtyOrdered))
		result = append(result, preparedPurchaseOrderItem{
			item:       item,
			supplierID: supplierID,
		})
	}
	return result, nil
}

func validatePurchaseOrderItem(item domain.PurchaseOrderItem) error {
	if item.ProductID == 0 {
		return ErrPurchaseOrderMissingProduct
	}
	if item.QtyOrdered == 0 {
		return ErrPurchaseOrderInvalidQty
	}
	if item.UnitCost <= 0 {
		return ErrPurchaseOrderInvalidUnitCost
	}
	return nil
}

func resolvePurchaseOrderItemSupplier(order *domain.PurchaseOrder, item domain.PurchaseOrderItem, product productdomain.Product) (uint64, error) {
	if item.SupplierID != nil && *item.SupplierID != 0 {
		return *item.SupplierID, nil
	}
	if product.SupplierID != nil && *product.SupplierID != 0 {
		return *product.SupplierID, nil
	}
	if order != nil && order.SupplierID != nil && *order.SupplierID != 0 {
		return *order.SupplierID, nil
	}
	return 0, ErrPurchaseOrderMissingSupplier
}

func (uc *PurchaseOrderUsecase) Receive(c *gin.Context, orderID uint64, params domain.PurchaseOrderReceiveParams) error {
	before, _ := uc.repo.GetByID(orderID)
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if uc.receiveTxManager != nil {
		if err := uc.receiveTxManager.Run(ctx, func(deps PurchaseOrderReceiveTransactionalDeps) error {
			return uc.receiveWithDeps(c, orderID, params, deps)
		}); err != nil {
			return err
		}
		if after, err := uc.repo.GetByID(orderID); err == nil {
			uc.recordPurchaseOrderAuditIfChanged(c, "RECEIVE", orderID, before, after)
		}
		return nil
	}
	if err := uc.receiveWithDeps(c, orderID, params, PurchaseOrderReceiveTransactionalDeps{
		Repo:              uc.repo,
		InventoryService:  uc.inventoryService,
		CostEventRecorder: uc.costEventRecorder,
	}); err != nil {
		return err
	}
	if after, err := uc.repo.GetByID(orderID); err == nil {
		uc.recordPurchaseOrderAuditIfChanged(c, "RECEIVE", orderID, before, after)
	}
	return nil
}

func (uc *PurchaseOrderUsecase) Inspect(c *gin.Context, orderID uint64, params domain.PurchaseOrderInspectParams) error {
	before, _ := uc.repo.GetByID(orderID)
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if uc.inspectTxManager != nil {
		if err := uc.inspectTxManager.Run(ctx, func(deps PurchaseOrderInspectTransactionalDeps) error {
			return uc.inspectWithDeps(c, orderID, params, deps)
		}); err != nil {
			return err
		}
		if after, err := uc.repo.GetByID(orderID); err == nil {
			uc.recordPurchaseOrderAuditIfChanged(c, "INSPECT", orderID, before, after)
		}
		return nil
	}
	if err := uc.inspectWithDeps(c, orderID, params, PurchaseOrderInspectTransactionalDeps{
		Repo:             uc.repo,
		InventoryService: uc.inventoryService,
	}); err != nil {
		return err
	}
	if after, err := uc.repo.GetByID(orderID); err == nil {
		uc.recordPurchaseOrderAuditIfChanged(c, "INSPECT", orderID, before, after)
	}
	return nil
}

func (uc *PurchaseOrderUsecase) receiveWithDeps(
	c *gin.Context,
	orderID uint64,
	params domain.PurchaseOrderReceiveParams,
	deps PurchaseOrderReceiveTransactionalDeps,
) error {
	if orderID == 0 {
		return errors.New("invalid order id")
	}

	execUC := *uc
	execUC.repo = deps.Repo
	execUC.inventoryService = deps.InventoryService
	execUC.costEventRecorder = deps.CostEventRecorder

	order, err := execUC.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != domain.PurchaseOrderStatusShipped {
		return errors.New("order status invalid for receive")
	}

	effectiveWarehouseID := params.WarehouseID
	if effectiveWarehouseID == 0 && order.WarehouseID != nil {
		effectiveWarehouseID = *order.WarehouseID
	}
	if effectiveWarehouseID == 0 {
		return errors.New("missing warehouse")
	}

	now := time.Now()
	operatorID := resolveOperatorID(c, params.OperatorID)
	allReceived := true
	ctx := context.Background()
	hasReceipt := false

	for i := range order.Items {
		item := &order.Items[i]
		delta := params.ReceivedQties[item.ID]
		if delta > 0 {
			hasReceipt = true
			if item.QtyReceived+delta > item.QtyOrdered {
				return errors.New("received qty exceeds ordered")
			}
			item.QtyReceived += delta

			// 仓库收货：从"采购在途"转到"待检"
			// 注意：发货时(MarkShipped)已经增加了"采购在途"库存
			if execUC.inventoryService != nil {
				unitCost := item.UnitCost
				remark := fmt.Sprintf("仓库收货入待检: %s", order.PoNumber)

				warehouseReceiveParams := &inventoryDomain.CreateMovementParams{
					ProductID:       item.ProductID,
					WarehouseID:     effectiveWarehouseID,
					MovementType:    inventoryDomain.MovementTypeWarehouseReceive,
					Quantity:        int(delta),
					ReferenceType:   stringPtr("PURCHASE_ORDER"),
					ReferenceID:     &orderID,
					ReferenceNumber: &order.PoNumber,
					UnitCost:        &unitCost,
					Remark:          &remark,
					OperatorID:      operatorID,
					OperatedAt:      &now,
				}
				if _, err := execUC.inventoryService.CreateMovement(ctx, warehouseReceiveParams); err != nil {
					return fmt.Errorf("failed to create warehouse receive movement: %w", err)
				}

				if item.Product != nil && item.Product.IsInspectionRequired == 0 {
					passRemark := fmt.Sprintf("免检产品自动质检通过: %s", order.PoNumber)
					if _, err := execUC.inventoryService.CreateMovement(ctx, &inventoryDomain.CreateMovementParams{
						ProductID:       item.ProductID,
						WarehouseID:     effectiveWarehouseID,
						MovementType:    inventoryDomain.MovementTypeInspectionPass,
						Quantity:        int(delta),
						ReferenceType:   stringPtr("PURCHASE_ORDER"),
						ReferenceID:     &orderID,
						ReferenceNumber: &order.PoNumber,
						UnitCost:        &unitCost,
						Remark:          &passRemark,
						OperatorID:      operatorID,
						OperatedAt:      &now,
					}); err != nil {
						return fmt.Errorf("failed to create auto inspection pass movement: %w", err)
					}
					item.QtyInspectionPass += delta
					if err := execUC.autoMovePendingShipmentIfPackingNotRequired(ctx, orderID, order.PoNumber, item, effectiveWarehouseID, delta, operatorID, now); err != nil {
						return err
					}
				}
			}

			if err := execUC.recordCostEvent(PurchaseOrderCostEventReceived, order, item, delta, &effectiveWarehouseID, now, operatorID); err != nil {
				return err
			}
		}

		if item.QtyReceived < item.QtyOrdered {
			allReceived = false
		}
	}

	if hasReceipt {
		order.ReceivedAt = &now
		order.ReceivedBy = operatorID
	}
	if allReceived {
		order.Status = domain.PurchaseOrderStatusReceived
	}

	if err := execUC.repo.UpdateProgress(order); err != nil {
		return err
	}
	return nil
}

func (uc *PurchaseOrderUsecase) inspectWithDeps(
	c *gin.Context,
	orderID uint64,
	params domain.PurchaseOrderInspectParams,
	deps PurchaseOrderInspectTransactionalDeps,
) error {
	if orderID == 0 {
		return errors.New("invalid order id")
	}

	execUC := *uc
	execUC.repo = deps.Repo
	execUC.inventoryService = deps.InventoryService

	order, err := execUC.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != domain.PurchaseOrderStatusShipped && order.Status != domain.PurchaseOrderStatusReceived {
		return errors.New("order status invalid for inspection")
	}
	if order.WarehouseID == nil || *order.WarehouseID == 0 {
		return errors.New("missing warehouse")
	}
	if execUC.inventoryService == nil {
		return errors.New("inventory service unavailable")
	}

	now := time.Now()
	ctx := context.Background()
	totalInspected := uint64(0)
	operatorID := resolveOperatorID(c, params.OperatorID)

	for i := range order.Items {
		item := &order.Items[i]
		passQty := params.PassQties[item.ID]
		failQty := params.FailQties[item.ID]
		if passQty == 0 && failQty == 0 {
			continue
		}
		pendingQty := deriveItemPendingInspection(*item)
		if passQty+failQty > pendingQty {
			return errors.New("inspection qty exceeds pending inspection")
		}

		if passQty > 0 {
			remark := fmt.Sprintf("采购质检通过: %s", order.PoNumber)
			if _, err := execUC.inventoryService.CreateMovement(ctx, &inventoryDomain.CreateMovementParams{
				ProductID:       item.ProductID,
				WarehouseID:     *order.WarehouseID,
				MovementType:    inventoryDomain.MovementTypeInspectionPass,
				Quantity:        int(passQty),
				ReferenceType:   stringPtr("PURCHASE_ORDER"),
				ReferenceID:     &orderID,
				ReferenceNumber: &order.PoNumber,
				UnitCost:        float64Pointer(item.UnitCost),
				Remark:          &remark,
				OperatorID:      operatorID,
				OperatedAt:      &now,
			}); err != nil {
				return fmt.Errorf("failed to create inspection pass movement: %w", err)
			}
			item.QtyInspectionPass += passQty
			if err := execUC.autoMovePendingShipmentIfPackingNotRequired(ctx, orderID, order.PoNumber, item, *order.WarehouseID, passQty, operatorID, now); err != nil {
				return err
			}
		}
		if failQty > 0 {
			remark := fmt.Sprintf("采购质检不合格损失结案: %s", order.PoNumber)
			if _, err := execUC.inventoryService.CreateMovement(ctx, &inventoryDomain.CreateMovementParams{
				ProductID:       item.ProductID,
				WarehouseID:     *order.WarehouseID,
				MovementType:    inventoryDomain.MovementTypeInspectionLoss,
				Quantity:        int(failQty),
				ReferenceType:   stringPtr("PURCHASE_ORDER"),
				ReferenceID:     &orderID,
				ReferenceNumber: &order.PoNumber,
				UnitCost:        float64Pointer(item.UnitCost),
				Remark:          &remark,
				OperatorID:      operatorID,
				OperatedAt:      &now,
			}); err != nil {
				return fmt.Errorf("failed to create inspection loss movement: %w", err)
			}
			item.QtyInspectionFail += failQty
		}

		totalInspected += passQty + failQty
	}

	if totalInspected == 0 {
		return errors.New("missing inspection qty")
	}
	order.InspectedAt = &now
	order.InspectedBy = operatorID
	return execUC.repo.UpdateProgress(order)
}

func stringPtr(s string) *string {
	return &s
}

func float64Pointer(v float64) *float64 {
	return &v
}

func (uc *PurchaseOrderUsecase) autoMovePendingShipmentIfPackingNotRequired(
	ctx context.Context,
	orderID uint64,
	poNumber string,
	item *domain.PurchaseOrderItem,
	warehouseID uint64,
	qty uint64,
	operatorID *uint64,
	operatedAt time.Time,
) error {
	if qty == 0 || item == nil || item.Product == nil || item.Product.IsPackingRequired != 0 || uc.inventoryService == nil {
		return nil
	}
	remark := fmt.Sprintf("免打包产品自动转待出: %s", poNumber)
	if _, err := uc.inventoryService.CreateMovement(ctx, &inventoryDomain.CreateMovementParams{
		ProductID:       item.ProductID,
		WarehouseID:     warehouseID,
		MovementType:    inventoryDomain.MovementTypePackingSkipComplete,
		Quantity:        int(qty),
		ReferenceType:   stringPtr("PURCHASE_ORDER"),
		ReferenceID:     &orderID,
		ReferenceNumber: &poNumber,
		UnitCost:        float64Pointer(item.UnitCost),
		Remark:          &remark,
		OperatorID:      operatorID,
		OperatedAt:      &operatedAt,
	}); err != nil {
		return fmt.Errorf("failed to create auto packing skip movement: %w", err)
	}
	return nil
}

func resolveOperatorID(c *gin.Context, fallback *uint64) *uint64 {
	if fallback != nil && *fallback != 0 {
		value := *fallback
		return &value
	}
	if c == nil {
		return nil
	}
	if raw, ok := c.Get(auth.UserIDKey); ok {
		switch typed := raw.(type) {
		case uint64:
			if typed != 0 {
				value := typed
				return &value
			}
		case uint:
			if typed != 0 {
				value := uint64(typed)
				return &value
			}
		case int:
			if typed > 0 {
				value := uint64(typed)
				return &value
			}
		}
	}
	if raw, ok := c.Get("user_id"); ok {
		if typed, ok := raw.(uint64); ok && typed != 0 {
			value := typed
			return &value
		}
	}
	return nil
}

func deriveItemPendingInspection(item domain.PurchaseOrderItem) uint64 {
	inspected := item.QtyInspectionPass + item.QtyInspectionFail
	if item.QtyReceived <= inspected {
		return 0
	}
	return item.QtyReceived - inspected
}

func getPurchaseOrderReceivedTotal(order *domain.PurchaseOrder) uint64 {
	if order == nil {
		return 0
	}
	total := uint64(0)
	for _, item := range order.Items {
		total += item.QtyReceived
	}
	return total
}

func getPurchaseOrderOrderedTotal(order *domain.PurchaseOrder) uint64 {
	if order == nil {
		return 0
	}
	total := uint64(0)
	for _, item := range order.Items {
		total += item.QtyOrdered
	}
	return total
}

func getPurchaseOrderInspectionFailTotal(order *domain.PurchaseOrder) uint64 {
	if order == nil {
		return 0
	}
	total := uint64(0)
	for _, item := range order.Items {
		total += item.QtyInspectionFail
	}
	return total
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
	beforeAudit := buildPurchaseOrderAuditSnapshot(existing)
	if existing.Status != domain.PurchaseOrderStatusDraft {
		return nil, errors.New("only draft can be updated")
	}

	currency := uc.defaultCurrency(updates.Currency)
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
	updated, err := uc.repo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	uc.recordAuditIfChanged(c, "UPDATE", "PurchaseOrder", fmt.Sprintf("%d", orderID), beforeAudit, buildPurchaseOrderAuditSnapshot(updated))
	return updated, nil
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
	uc.recordPurchaseOrderAuditIfChanged(c, "DELETE", orderID, order, nil)
	return nil
}

func (uc *PurchaseOrderUsecase) Submit(c *gin.Context, orderID uint64) error {
	before, _ := uc.repo.GetByID(orderID)
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if uc.submitTxManager != nil {
		if err := uc.submitTxManager.Run(ctx, func(deps PurchaseOrderSubmitTransactionalDeps) error {
			return uc.submitWithDeps(c, orderID, deps)
		}); err != nil {
			return err
		}
		if after, err := uc.repo.GetByID(orderID); err == nil {
			uc.recordPurchaseOrderAuditIfChanged(c, "SUBMIT", orderID, before, after)
		}
		return nil
	}
	if err := uc.submitWithDeps(c, orderID, PurchaseOrderSubmitTransactionalDeps{
		Repo:              uc.repo,
		PlanCleaner:       uc.planCleaner,
		CostEventRecorder: uc.costEventRecorder,
	}); err != nil {
		return err
	}
	if after, err := uc.repo.GetByID(orderID); err == nil {
		uc.recordPurchaseOrderAuditIfChanged(c, "SUBMIT", orderID, before, after)
	}
	return nil
}

func (uc *PurchaseOrderUsecase) submitWithDeps(
	c *gin.Context,
	orderID uint64,
	deps PurchaseOrderSubmitTransactionalDeps,
) error {
	execUC := *uc
	execUC.repo = deps.Repo
	execUC.planCleaner = deps.PlanCleaner
	execUC.costEventRecorder = deps.CostEventRecorder

	order, err := execUC.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != domain.PurchaseOrderStatusDraft {
		return errors.New("order status invalid for submit")
	}
	now := time.Now()
	orderItems := clonePurchaseOrderItems(order.Items)
	order.Status = domain.PurchaseOrderStatusOrdered
	order.OrderedAt = &now
	order.OrderedBy = resolveOperatorID(c, order.UpdatedBy)
	order.Items = nil
	if err := execUC.repo.UpdateProgress(order); err != nil {
		return err
	}
	for i := range orderItems {
		if err := execUC.recordCostEvent(PurchaseOrderCostEventOrdered, order, &orderItems[i], orderItems[i].QtyOrdered, nil, now, order.OrderedBy); err != nil {
			return err
		}
	}
	return nil
}

func (uc *PurchaseOrderUsecase) MarkShipped(c *gin.Context, orderID uint64, params domain.PurchaseOrderShipParams) error {
	before, _ := uc.repo.GetByID(orderID)
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if uc.shipTxManager != nil {
		if err := uc.shipTxManager.Run(ctx, func(deps PurchaseOrderShipTransactionalDeps) error {
			return uc.markShippedWithDeps(c, orderID, params, deps)
		}); err != nil {
			return err
		}
		if after, err := uc.repo.GetByID(orderID); err == nil {
			uc.recordPurchaseOrderAuditIfChanged(c, "SHIP", orderID, before, after)
		}
		return nil
	}
	if err := uc.markShippedWithDeps(c, orderID, params, PurchaseOrderShipTransactionalDeps{
		Repo:              uc.repo,
		InventoryService:  uc.inventoryService,
		CostEventRecorder: uc.costEventRecorder,
	}); err != nil {
		return err
	}
	if after, err := uc.repo.GetByID(orderID); err == nil {
		uc.recordPurchaseOrderAuditIfChanged(c, "SHIP", orderID, before, after)
	}
	return nil
}

func (uc *PurchaseOrderUsecase) markShippedWithDeps(
	c *gin.Context,
	orderID uint64,
	params domain.PurchaseOrderShipParams,
	deps PurchaseOrderShipTransactionalDeps,
) error {
	execUC := *uc
	execUC.repo = deps.Repo
	execUC.inventoryService = deps.InventoryService
	execUC.costEventRecorder = deps.CostEventRecorder

	order, err := execUC.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != domain.PurchaseOrderStatusOrdered {
		return errors.New("order status invalid for ship")
	}
	if params.WarehouseID == 0 {
		return errors.New("missing warehouse")
	}

	now := time.Now()
	operatorID := resolveOperatorID(c, params.OperatorID)
	ctx := context.Background()

	// 供应商发货 → 货物进入"采购在途"状态
	if execUC.inventoryService != nil && params.WarehouseID != 0 {
		for _, item := range order.Items {
			if item.QtyOrdered == 0 {
				continue
			}
			unitCost := item.UnitCost
			purchaseShipParams := &inventoryDomain.CreateMovementParams{
				ProductID:       item.ProductID,
				WarehouseID:     params.WarehouseID,
				MovementType:    inventoryDomain.MovementTypePurchaseShip,
				Quantity:        int(item.QtyOrdered),
				ReferenceType:   stringPtr("PURCHASE_ORDER"),
				ReferenceID:     &orderID,
				ReferenceNumber: &order.PoNumber,
				UnitCost:        &unitCost,
				Remark:          stringPtr(fmt.Sprintf("采购发货入在途: %s", order.PoNumber)),
				OperatorID:      operatorID,
				OperatedAt:      &now,
			}
			if _, err := execUC.inventoryService.CreateMovement(ctx, purchaseShipParams); err != nil {
				return fmt.Errorf("failed to create purchase ship movement: %w", err)
			}
		}
	}
	for i := range order.Items {
		item := &order.Items[i]
		if item.QtyOrdered == 0 {
			continue
		}
		if err := execUC.recordCostEvent(PurchaseOrderCostEventShipped, order, item, item.QtyOrdered, &params.WarehouseID, now, operatorID); err != nil {
			return err
		}
	}

	order.Status = domain.PurchaseOrderStatusShipped
	order.WarehouseID = &params.WarehouseID
	order.ShippedAt = &now
	order.ShippedBy = operatorID
	order.Items = nil
	if err := execUC.repo.UpdateProgress(order); err != nil {
		return err
	}
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
	if order.Status != domain.PurchaseOrderStatusReceived {
		return ErrPurchaseOrderInvalidCompleteStatus
	}
	if getPurchaseOrderPendingInspectionTotal(order) > 0 {
		return ErrPurchaseOrderPendingInspection
	}
	if getPurchaseOrderReceivedTotal(order) < getPurchaseOrderOrderedTotal(order) {
		return ErrPurchaseOrderIncompleteReceipt
	}
	if getPurchaseOrderInspectionFailTotal(order) > 0 {
		return ErrPurchaseOrderInspectionFailed
	}
	beforeAudit := buildPurchaseOrderAuditSnapshot(order)
	now := time.Now()
	order.Status = domain.PurchaseOrderStatusClosed
	order.ClosedAt = &now
	order.CompletedBy = resolveOperatorID(c, order.UpdatedBy)
	order.IsForceCompleted = 0
	order.ForceCompletedAt = nil
	order.ForceCompletedBy = nil
	order.ForceCompleteReason = ""
	order.Items = nil
	if err := uc.repo.UpdateProgress(order); err != nil {
		return err
	}
	if after, err := uc.repo.GetByID(orderID); err == nil {
		uc.recordAuditIfChanged(c, "CLOSE", "PurchaseOrder", fmt.Sprintf("%d", orderID), beforeAudit, buildPurchaseOrderAuditSnapshot(after))
	}
	return nil
}

func (uc *PurchaseOrderUsecase) ForceComplete(c *gin.Context, orderID uint64, params domain.PurchaseOrderForceCompleteParams) error {
	order, err := uc.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status == domain.PurchaseOrderStatusClosed {
		return nil
	}
	if order.Status != domain.PurchaseOrderStatusShipped && order.Status != domain.PurchaseOrderStatusReceived {
		return ErrPurchaseOrderInvalidCompleteStatus
	}
	if strings.TrimSpace(params.Reason) == "" {
		return ErrPurchaseOrderForceCompleteReason
	}
	if getPurchaseOrderPendingInspectionTotal(order) > 0 {
		return ErrPurchaseOrderPendingInspection
	}
	receivedTotal := getPurchaseOrderReceivedTotal(order)
	if receivedTotal == 0 {
		return ErrPurchaseOrderInvalidCompleteStatus
	}
	receivedIncomplete := receivedTotal < getPurchaseOrderOrderedTotal(order)
	hasInspectionLoss := getPurchaseOrderInspectionFailTotal(order) > 0
	if !receivedIncomplete && !hasInspectionLoss {
		return ErrPurchaseOrderInvalidCompleteStatus
	}

	beforeAudit := buildPurchaseOrderAuditSnapshot(order)
	now := time.Now()
	operatorID := resolveOperatorID(c, params.OperatorID)
	order.Status = domain.PurchaseOrderStatusClosed
	order.ClosedAt = &now
	order.CompletedBy = operatorID
	order.IsForceCompleted = 1
	order.ForceCompletedAt = &now
	order.ForceCompletedBy = operatorID
	order.ForceCompleteReason = strings.TrimSpace(params.Reason)
	order.Items = nil
	if err := uc.repo.UpdateProgress(order); err != nil {
		return err
	}
	if after, err := uc.repo.GetByID(orderID); err == nil {
		uc.recordAuditIfChanged(c, "FORCE_COMPLETE", "PurchaseOrder", fmt.Sprintf("%d", orderID), beforeAudit, buildPurchaseOrderAuditSnapshot(after))
	}
	return nil
}

func (uc *PurchaseOrderUsecase) recordCostEvent(
	eventType PurchaseOrderCostEventType,
	order *domain.PurchaseOrder,
	item *domain.PurchaseOrderItem,
	qty uint64,
	warehouseID *uint64,
	occurredAt time.Time,
	operatorID *uint64,
) error {
	if uc.costEventRecorder == nil || order == nil || item == nil || qty == 0 {
		return nil
	}
	return uc.costEventRecorder.RecordPurchaseOrderEvent(&PurchaseOrderCostEventParams{
		EventType:           eventType,
		PurchaseOrderID:     order.ID,
		PurchaseOrderItemID: item.ID,
		ProductID:           item.ProductID,
		WarehouseID:         warehouseID,
		Marketplace:         order.Marketplace,
		Currency:            uc.defaultCurrency(item.Currency),
		UnitCost:            item.UnitCost,
		QtyEvent:            qty,
		OccurredAt:          occurredAt,
		OperatorID:          operatorID,
	})
}

func clonePurchaseOrderItems(items []domain.PurchaseOrderItem) []domain.PurchaseOrderItem {
	if len(items) == 0 {
		return []domain.PurchaseOrderItem{}
	}
	result := make([]domain.PurchaseOrderItem, len(items))
	copy(result, items)
	return result
}

func cloneUint64Slice(values []uint64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}
	result := make([]uint64, len(values))
	copy(result, values)
	return result
}

func (uc *PurchaseOrderUsecase) expandComboItems(currency string, items []domain.PurchaseOrderItem) ([]domain.PurchaseOrderItem, error) {
	if len(items) == 0 {
		return []domain.PurchaseOrderItem{}, nil
	}

	if uc.productLookup == nil {
		return normalizeItems(currency, items), nil
	}

	productIDs := uniqueProductIDs(items)
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
		product, ok := productMap[item.ProductID]
		if !ok || product.ComboID == nil {
			continue
		}
		if product.IsComboMain != 1 {
			comboExpanded[*product.ComboID] = true
		}
	}

	result := make(map[uint64]*domain.PurchaseOrderItem, len(items))
	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if ok && product.ComboID != nil && product.IsComboMain == 1 {
			if comboExpanded[*product.ComboID] {
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
				if comboItem.ProductID == comboItem.MainProductID {
					continue
				}
				qty := item.QtyOrdered * comboItem.QtyRatio
				if qty == 0 {
					continue
				}
				unitCost := getUnitCost(productMap[comboItem.ProductID])
				mergeItem(result, comboItem.ProductID, qty, unitCost, currency, nil, item.SourcePlanID)
			}
			continue
		}

		mergeItem(result, item.ProductID, item.QtyOrdered, item.UnitCost, currency, item.SupplierID, item.SourcePlanID)
	}

	return flattenItems(result), nil
}

func mergeItem(target map[uint64]*domain.PurchaseOrderItem, productID, qty uint64, unitCost float64, currency string, supplierID *uint64, sourcePlanID *uint64) {
	if qty == 0 {
		return
	}
	if existing, ok := target[productID]; ok {
		existing.QtyOrdered += qty
		if existing.UnitCost == 0 && unitCost > 0 {
			existing.UnitCost = unitCost
		}
		if existing.SupplierID == nil && supplierID != nil {
			existing.SupplierID = supplierID
		}
		if existing.SourcePlanID == nil && sourcePlanID != nil {
			existing.SourcePlanID = sourcePlanID
		}
		existing.Subtotal = roundAmount(existing.UnitCost * float64(existing.QtyOrdered))
		return
	}

	item := &domain.PurchaseOrderItem{
		ProductID:    productID,
		SupplierID:   supplierID,
		SourcePlanID: sourcePlanID,
		QtyOrdered:   qty,
		QtyReceived:  0,
		UnitCost:     unitCost,
		Currency:     currency,
	}
	item.Subtotal = roundAmount(unitCost * float64(qty))
	target[productID] = item
}

func flattenItems(items map[uint64]*domain.PurchaseOrderItem) []domain.PurchaseOrderItem {
	result := make([]domain.PurchaseOrderItem, 0, len(items))
	for _, item := range items {
		result = append(result, *item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ProductID < result[j].ProductID
	})
	return result
}

func sourcePlanIDsFromItems(items []domain.PurchaseOrderItem) []uint64 {
	result := make([]uint64, 0, len(items))
	seen := map[uint64]struct{}{}
	for _, item := range items {
		if item.SourcePlanID == nil || *item.SourcePlanID == 0 {
			continue
		}
		if _, exists := seen[*item.SourcePlanID]; exists {
			continue
		}
		seen[*item.SourcePlanID] = struct{}{}
		result = append(result, *item.SourcePlanID)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func appendUniqueUint64s(base []uint64, values ...uint64) []uint64 {
	if len(values) == 0 {
		return base
	}
	seen := make(map[uint64]struct{}, len(base)+len(values))
	for _, value := range base {
		seen[value] = struct{}{}
	}
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		base = append(base, value)
	}
	sort.Slice(base, func(i, j int) bool { return base[i] < base[j] })
	return base
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

func (uc *PurchaseOrderUsecase) defaultCurrency(currency string) string {
	if currency == "" {
		if uc != nil && uc.defaultsProvider != nil {
			if defaultCurrency := strings.TrimSpace(uc.defaultsProvider.GetDefaultBaseCurrency()); defaultCurrency != "" {
				return defaultCurrency
			}
		}
		return "USD"
	}
	return currency
}

func ensurePoNumber(poNumber string) string {
	if poNumber != "" {
		return poNumber
	}
	return numbering.Generate("PO", time.Now())
}

func uniqueProductIDs(items []domain.PurchaseOrderItem) []uint64 {
	seen := make(map[uint64]struct{}, len(items))
	result := make([]uint64, 0, len(items))
	for _, item := range items {
		if item.ProductID == 0 {
			continue
		}
		if _, ok := seen[item.ProductID]; ok {
			continue
		}
		seen[item.ProductID] = struct{}{}
		result = append(result, item.ProductID)
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

func (uc *PurchaseOrderUsecase) recordAuditIfChanged(c *gin.Context, action, entityType, entityID string, before, after any) {
	beforeDiff, afterDiff, changed := buildProcurementAuditDiff(before, after)
	if !changed {
		return
	}
	uc.recordAudit(c, action, entityType, entityID, beforeDiff, afterDiff)
}

func (uc *PurchaseOrderUsecase) recordPurchaseOrderAuditIfChanged(c *gin.Context, action string, orderID uint64, before, after *domain.PurchaseOrder) {
	uc.recordAuditIfChanged(
		c,
		action,
		"PurchaseOrder",
		fmt.Sprintf("%d", orderID),
		buildPurchaseOrderAuditSnapshot(before),
		buildPurchaseOrderAuditSnapshot(after),
	)
}

func buildPurchaseOrderAuditSnapshot(order *domain.PurchaseOrder) any {
	if order == nil {
		return nil
	}

	items := make([]map[string]any, 0, len(order.Items))
	for _, item := range order.Items {
		productCode := fmt.Sprintf("产品#%d", item.ProductID)
		productTitle := ""
		if item.Product != nil {
			if strings.TrimSpace(item.Product.SellerSku) != "" {
				productCode = item.Product.SellerSku
			}
			productTitle = strings.TrimSpace(item.Product.Title)
		}
		items = append(items, map[string]any{
			"seller_sku":             productCode,
			"product_title":          productTitle,
			"qty_ordered":            item.QtyOrdered,
			"qty_received":           item.QtyReceived,
			"qty_inspection_pass":    item.QtyInspectionPass,
			"qty_inspection_fail":    item.QtyInspectionFail,
			"qty_pending_inspection": deriveItemPendingInspection(item),
			"unit_cost":              roundAmount(item.UnitCost),
			"subtotal":               roundAmount(item.Subtotal),
		})
	}

	supplierName := ""
	if order.Supplier != nil && strings.TrimSpace(order.Supplier.Name) != "" {
		supplierName = order.Supplier.Name
	} else if order.SupplierID != nil && *order.SupplierID != 0 {
		supplierName = fmt.Sprintf("供应商#%d", *order.SupplierID)
	}

	completedByName := strings.TrimSpace(order.CompletedByName)
	if completedByName == "" {
		completedByName = strings.TrimSpace(order.ForceCompletedByName)
	}

	return map[string]any{
		"po_number":                    order.PoNumber,
		"supplier_name":                supplierName,
		"marketplace":                  order.Marketplace,
		"status":                       string(order.Status),
		"currency":                     order.Currency,
		"total_amount":                 roundAmount(order.TotalAmount),
		"created_at":                   order.GmtCreate,
		"ordered_at":                   order.OrderedAt,
		"ordered_by_name":              formatProcurementAuditOperatorName(order.OrderedByName, order.OrderedBy),
		"shipped_at":                   order.ShippedAt,
		"shipped_by_name":              formatProcurementAuditOperatorName(order.ShippedByName, order.ShippedBy),
		"received_at":                  order.ReceivedAt,
		"received_by_name":             formatProcurementAuditOperatorName(order.ReceivedByName, order.ReceivedBy),
		"inspected_at":                 order.InspectedAt,
		"inspected_by_name":            formatProcurementAuditOperatorName(order.InspectedByName, order.InspectedBy),
		"closed_at":                    order.ClosedAt,
		"completed_by_name":            formatProcurementAuditOperatorName(completedByName, order.CompletedBy),
		"force_completed_at":           order.ForceCompletedAt,
		"force_completed_by_name":      formatProcurementAuditOperatorName(order.ForceCompletedByName, order.ForceCompletedBy),
		"force_complete_reason":        order.ForceCompleteReason,
		"qty_pending_inspection_total": getPurchaseOrderPendingInspectionTotal(order),
		"remark":                       order.Remark,
		"items":                        items,
	}
}

func formatProcurementAuditOperatorName(name string, operatorID *uint64) string {
	if strings.TrimSpace(name) != "" {
		return strings.TrimSpace(name)
	}
	if operatorID != nil && *operatorID != 0 {
		return fmt.Sprintf("用户#%d", *operatorID)
	}
	return ""
}

func getPurchaseOrderPendingInspectionTotal(order *domain.PurchaseOrder) uint64 {
	if order == nil {
		return 0
	}
	if len(order.Items) == 0 {
		return order.QtyPendingInspectionTotal
	}
	total := uint64(0)
	for _, item := range order.Items {
		total += deriveItemPendingInspection(item)
	}
	return total
}

func buildPurchaseOrderInspectionAuditItem(item domain.PurchaseOrderItem, passQty, failQty uint64) map[string]any {
	productCode := fmt.Sprintf("产品#%d", item.ProductID)
	productTitle := ""
	if item.Product != nil {
		if strings.TrimSpace(item.Product.SellerSku) != "" {
			productCode = item.Product.SellerSku
		}
		productTitle = strings.TrimSpace(item.Product.Title)
	}

	result := map[string]any{
		"seller_sku":    productCode,
		"product_title": productTitle,
	}
	if passQty > 0 {
		result["qty_inspection_pass"] = passQty
	}
	if failQty > 0 {
		result["qty_inspection_fail"] = failQty
	}
	return result
}

func buildPurchaseOrderInspectionAuditBaseline(item map[string]any) map[string]any {
	result := map[string]any{}
	if sellerSKU, ok := item["seller_sku"]; ok {
		result["seller_sku"] = sellerSKU
	}
	if productTitle, ok := item["product_title"]; ok {
		result["product_title"] = productTitle
	}
	if _, ok := item["qty_inspection_pass"]; ok {
		result["qty_inspection_pass"] = 0
	}
	if _, ok := item["qty_inspection_fail"]; ok {
		result["qty_inspection_fail"] = 0
	}
	return result
}

func buildProcurementAuditDiff(before, after any) (any, any, bool) {
	normalizedBefore := normalizeProcurementAuditValue(before)
	normalizedAfter := normalizeProcurementAuditValue(after)
	return diffProcurementAuditValues(normalizedBefore, normalizedAfter)
}

func normalizeProcurementAuditValue(value any) any {
	if value == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var normalized any
	if err := json.Unmarshal(raw, &normalized); err != nil {
		return value
	}
	return scrubProcurementAuditValue(normalized)
}

func scrubProcurementAuditValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		cleaned := make(map[string]any, len(typed))
		for key, item := range typed {
			switch key {
			case "gmt_create", "gmt_modified", "created_at", "updated_at", "created_by", "updated_by", "id", "purchase_order_id", "supplier_id", "product_id":
				continue
			}
			cleaned[key] = scrubProcurementAuditValue(item)
		}
		return cleaned
	case []any:
		items := make([]any, 0, len(typed))
		for _, item := range typed {
			items = append(items, scrubProcurementAuditValue(item))
		}
		return items
	default:
		return value
	}
}

func diffProcurementAuditValues(before, after any) (any, any, bool) {
	if reflect.DeepEqual(before, after) {
		return nil, nil, false
	}

	beforeSlice, beforeIsSlice := before.([]any)
	afterSlice, afterIsSlice := after.([]any)
	if beforeIsSlice && afterIsSlice {
		return diffProcurementAuditArrays(beforeSlice, afterSlice)
	}

	beforeMap, beforeIsMap := before.(map[string]any)
	afterMap, afterIsMap := after.(map[string]any)
	if beforeIsMap && afterIsMap {
		beforeDiff := map[string]any{}
		afterDiff := map[string]any{}
		keys := make(map[string]struct{}, len(beforeMap)+len(afterMap))
		for key := range beforeMap {
			keys[key] = struct{}{}
		}
		for key := range afterMap {
			keys[key] = struct{}{}
		}
		for key := range keys {
			childBefore, childAfter, childChanged := diffProcurementAuditValues(beforeMap[key], afterMap[key])
			if !childChanged {
				continue
			}
			beforeDiff[key] = childBefore
			afterDiff[key] = childAfter
		}
		if len(beforeDiff) == 0 && len(afterDiff) == 0 {
			return nil, nil, false
		}
		return beforeDiff, afterDiff, true
	}

	return before, after, true
}

func diffProcurementAuditArrays(before, after []any) (any, any, bool) {
	if !isProcurementAuditItemArray(before, after) {
		return before, after, true
	}

	beforeItems := mapProcurementAuditItems(before)
	afterItems := mapProcurementAuditItems(after)
	keys := make([]string, 0, len(beforeItems)+len(afterItems))
	keySet := make(map[string]struct{}, len(beforeItems)+len(afterItems))
	for key := range beforeItems {
		keySet[key] = struct{}{}
	}
	for key := range afterItems {
		keySet[key] = struct{}{}
	}
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	beforeDiff := make([]any, 0, len(keys))
	afterDiff := make([]any, 0, len(keys))
	for _, key := range keys {
		beforeItem, beforeExists := beforeItems[key]
		afterItem, afterExists := afterItems[key]
		switch {
		case beforeExists && afterExists:
			childBefore, childAfter, changed := diffProcurementAuditItemMaps(beforeItem, afterItem)
			if !changed {
				continue
			}
			beforeDiff = append(beforeDiff, childBefore)
			afterDiff = append(afterDiff, childAfter)
		case beforeExists:
			beforeDiff = append(beforeDiff, ensureProcurementAuditItemIdentity(filterProcurementAuditItemFields(beforeItem), beforeItem, nil))
			afterDiff = append(afterDiff, ensureProcurementAuditItemIdentity(map[string]any{}, beforeItem, nil))
		case afterExists:
			beforeDiff = append(beforeDiff, ensureProcurementAuditItemIdentity(map[string]any{}, nil, afterItem))
			afterDiff = append(afterDiff, ensureProcurementAuditItemIdentity(filterProcurementAuditItemFields(afterItem), nil, afterItem))
		}
	}

	if len(beforeDiff) == 0 && len(afterDiff) == 0 {
		return nil, nil, false
	}
	return beforeDiff, afterDiff, true
}

func isProcurementAuditItemArray(before, after []any) bool {
	if len(before) == 0 && len(after) == 0 {
		return false
	}
	for _, item := range append(append([]any{}, before...), after...) {
		itemMap, ok := item.(map[string]any)
		if !ok {
			return false
		}
		if _, hasSellerSKU := itemMap["seller_sku"]; hasSellerSKU {
			continue
		}
		if _, hasProductTitle := itemMap["product_title"]; hasProductTitle {
			continue
		}
		return false
	}
	return true
}

func mapProcurementAuditItems(items []any) map[string]map[string]any {
	result := make(map[string]map[string]any, len(items))
	for index, item := range items {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		result[procurementAuditItemKey(itemMap, index)] = itemMap
	}
	return result
}

func procurementAuditItemKey(item map[string]any, index int) string {
	if sellerSKU, ok := item["seller_sku"]; ok {
		if value := strings.TrimSpace(fmt.Sprint(sellerSKU)); value != "" {
			return value
		}
	}
	if productTitle, ok := item["product_title"]; ok {
		if value := strings.TrimSpace(fmt.Sprint(productTitle)); value != "" {
			return value
		}
	}
	return fmt.Sprintf("item-%d", index)
}

func diffProcurementAuditItemMaps(before, after map[string]any) (map[string]any, map[string]any, bool) {
	childBefore, childAfter, changed := diffProcurementAuditValues(
		filterProcurementAuditItemFields(before),
		filterProcurementAuditItemFields(after),
	)
	if !changed {
		return nil, nil, false
	}

	beforeDiff, ok := childBefore.(map[string]any)
	if !ok {
		beforeDiff = map[string]any{}
	}
	afterDiff, ok := childAfter.(map[string]any)
	if !ok {
		afterDiff = map[string]any{}
	}
	return ensureProcurementAuditItemIdentity(beforeDiff, before, after), ensureProcurementAuditItemIdentity(afterDiff, before, after), true
}

func filterProcurementAuditItemFields(item map[string]any) map[string]any {
	if item == nil {
		return map[string]any{}
	}
	filtered := make(map[string]any, len(item))
	for key, value := range item {
		switch key {
		case "subtotal":
			continue
		default:
			filtered[key] = value
		}
	}
	return filtered
}

func ensureProcurementAuditItemIdentity(diff map[string]any, before, after map[string]any) map[string]any {
	if diff == nil {
		diff = map[string]any{}
	}
	if _, exists := diff["seller_sku"]; !exists {
		if value := procurementAuditItemIdentityValue("seller_sku", before, after); value != nil {
			diff["seller_sku"] = value
		}
	}
	if _, exists := diff["product_title"]; !exists {
		if value := procurementAuditItemIdentityValue("product_title", before, after); value != nil {
			diff["product_title"] = value
		}
	}
	return diff
}

func procurementAuditItemIdentityValue(key string, before, after map[string]any) any {
	if after != nil {
		if value, ok := after[key]; ok && strings.TrimSpace(fmt.Sprint(value)) != "" {
			return value
		}
	}
	if before != nil {
		if value, ok := before[key]; ok && strings.TrimSpace(fmt.Sprint(value)) != "" {
			return value
		}
	}
	return nil
}
