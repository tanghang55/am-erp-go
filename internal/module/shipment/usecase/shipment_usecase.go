package usecase

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/numbering"
	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	logisticsDomain "am-erp-go/internal/module/logistics/domain"
	productDomain "am-erp-go/internal/module/product/domain"
	"am-erp-go/internal/module/shipment/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

var (
	ErrShipmentNotFound        = errors.New("shipment not found")
	ErrInvalidStatus           = errors.New("invalid shipment status")
	ErrShipmentAlreadyExists   = errors.New("shipment already exists")
	ErrEmptyItems              = errors.New("shipment must have at least one item")
	ErrInsufficientInventory   = errors.New("insufficient inventory")
	ErrShipmentProductNotFound = errors.New("shipment product not found")
	ErrShipmentInactiveProduct = errors.New("shipment only supports on sale or replenishing products")
)

func newShipmentInvalidStatusError(message string) error {
	return fmt.Errorf("%w: %s", ErrInvalidStatus, message)
}

func newShipmentInsufficientInventoryError(message string) error {
	return fmt.Errorf("%w: %s", ErrInsufficientInventory, message)
}

// InventoryService interface for inventory operations
type InventoryService interface {
	CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error)
	GetProductBalance(productID, warehouseID uint64) (*inventoryDomain.InventoryBalance, error)
}

// ProductRepository interface for loading product info
type ProductRepository interface {
	ListByIDs(ids []uint64) ([]productDomain.Product, error)
}

// WarehouseRepository interface for loading warehouse info
type WarehouseRepository interface {
	GetByID(id uint64) (*inventoryDomain.Warehouse, error)
}

type LogisticsProviderRepository interface {
	GetByID(id uint64) (*logisticsDomain.LogisticsProvider, error)
}

type ShippingRateRepository interface {
	GetByID(id uint64) (*logisticsDomain.ShippingRate, error)
}

type ShipmentDefaultsProvider interface {
	GetDefaultBaseCurrency() string
}

type ShipmentFXSnapshot struct {
	Rate        float64
	Source      string
	Version     string
	EffectiveAt time.Time
}

type ShipmentFXResolver func(baseCurrency, originalCurrency string, occurredAt time.Time) (*ShipmentFXSnapshot, error)

type ShipmentCostAllocationLine struct {
	ShipmentItemID uint64
	ProductID      uint64
	WarehouseID    uint64
	Marketplace    string
	Quantity       uint64
	ItemUnitCost   float64
	ItemCurrency   string
	OriginalAmount float64
	BaseAmount     float64
}

type ShipmentCostAllocationRecordParams struct {
	ShipmentID       uint64
	ShipmentNumber   string
	OriginalCurrency string
	BaseCurrency     string
	FxRate           float64
	FxSource         string
	FxVersion        string
	FxTime           time.Time
	OccurredAt       time.Time
	OperatorID       *uint64
	Lines            []ShipmentCostAllocationLine
}

type ShipmentCostAllocationRecorder interface {
	RecordShipmentCostAllocation(params *ShipmentCostAllocationRecordParams) error
}

type ShipmentLandedSnapshotRecorder interface {
	UpsertShipmentLandedSnapshots(params *ShipmentCostAllocationRecordParams) error
}

type AuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

type ShipmentMarkShippedTransactionalDeps struct {
	ShipmentRepo     domain.ShipmentRepository
	ShipmentItemRepo domain.ShipmentItemRepository
	InventoryService InventoryService
	CostRecorder     ShipmentCostAllocationRecorder
	LandedRecorder   ShipmentLandedSnapshotRecorder
}

type ShipmentMarkShippedTransactionManager interface {
	Run(ctx context.Context, fn func(ShipmentMarkShippedTransactionalDeps) error) error
}

type ShipmentStateTransactionalDeps struct {
	ShipmentRepo     domain.ShipmentRepository
	ShipmentItemRepo domain.ShipmentItemRepository
	InventoryService InventoryService
}

type ShipmentConfirmTransactionManager interface {
	Run(ctx context.Context, fn func(ShipmentStateTransactionalDeps) error) error
}

type ShipmentCancelTransactionManager interface {
	Run(ctx context.Context, fn func(ShipmentStateTransactionalDeps) error) error
}

type ShipmentUsecase struct {
	shipmentRepo          domain.ShipmentRepository
	shipmentItemRepo      domain.ShipmentItemRepository
	inventoryService      InventoryService
	productRepo           ProductRepository
	warehouseRepo         WarehouseRepository
	logisticsProviderRepo LogisticsProviderRepository
	shippingRateRepo      ShippingRateRepository
	defaultsProvider      ShipmentDefaultsProvider
	fxResolver            ShipmentFXResolver
	costRecorder          ShipmentCostAllocationRecorder
	landedRecorder        ShipmentLandedSnapshotRecorder
	auditLogger           AuditLogger
	markShippedTxManager  ShipmentMarkShippedTransactionManager
	confirmTxManager      ShipmentConfirmTransactionManager
	cancelTxManager       ShipmentCancelTransactionManager
}

func NewShipmentUsecase(
	shipmentRepo domain.ShipmentRepository,
	shipmentItemRepo domain.ShipmentItemRepository,
	inventoryService InventoryService,
	productRepo ProductRepository,
	warehouseRepo WarehouseRepository,
) *ShipmentUsecase {
	return &ShipmentUsecase{
		shipmentRepo:     shipmentRepo,
		shipmentItemRepo: shipmentItemRepo,
		inventoryService: inventoryService,
		productRepo:      productRepo,
		warehouseRepo:    warehouseRepo,
	}
}

func (uc *ShipmentUsecase) BindDefaultsProvider(provider ShipmentDefaultsProvider) {
	uc.defaultsProvider = provider
}

func (uc *ShipmentUsecase) BindLogisticsProviderRepository(repo LogisticsProviderRepository) {
	uc.logisticsProviderRepo = repo
}

func (uc *ShipmentUsecase) BindShippingRateRepository(repo ShippingRateRepository) {
	uc.shippingRateRepo = repo
}

func (uc *ShipmentUsecase) BindFXResolver(resolver ShipmentFXResolver) {
	uc.fxResolver = resolver
}

func (uc *ShipmentUsecase) BindCostAllocationRecorder(recorder ShipmentCostAllocationRecorder) {
	uc.costRecorder = recorder
}

func (uc *ShipmentUsecase) BindLandedSnapshotRecorder(recorder ShipmentLandedSnapshotRecorder) {
	uc.landedRecorder = recorder
}

func (uc *ShipmentUsecase) BindAuditLogger(logger AuditLogger) {
	uc.auditLogger = logger
}

func (uc *ShipmentUsecase) BindMarkShippedTransactionManager(manager ShipmentMarkShippedTransactionManager) {
	uc.markShippedTxManager = manager
}

func (uc *ShipmentUsecase) BindConfirmTransactionManager(manager ShipmentConfirmTransactionManager) {
	uc.confirmTxManager = manager
}

func (uc *ShipmentUsecase) BindCancelTransactionManager(manager ShipmentCancelTransactionManager) {
	uc.cancelTxManager = manager
}

// List 获取发货单列表
func (uc *ShipmentUsecase) List(params *domain.ShipmentListParams) ([]*domain.Shipment, int64, error) {
	shipments, total, err := uc.shipmentRepo.List(params)
	if err != nil {
		return nil, 0, err
	}
	if len(shipments) == 0 {
		return shipments, total, nil
	}

	if uc.warehouseRepo != nil {
		warehouseMap := make(map[uint64]interface{}, len(shipments))
		for _, shipment := range shipments {
			if shipment == nil || shipment.WarehouseID == 0 {
				continue
			}
			if _, exists := warehouseMap[shipment.WarehouseID]; exists {
				shipment.Warehouse = warehouseMap[shipment.WarehouseID]
				continue
			}
			warehouse, warehouseErr := uc.warehouseRepo.GetByID(shipment.WarehouseID)
			if warehouseErr != nil {
				continue
			}
			warehouseMap[shipment.WarehouseID] = warehouse
			shipment.Warehouse = warehouse
		}
	}

	if uc.logisticsProviderRepo != nil {
		providerMap := make(map[uint64]interface{}, len(shipments))
		for _, shipment := range shipments {
			if shipment == nil || shipment.LogisticsProviderID == nil || *shipment.LogisticsProviderID == 0 {
				continue
			}
			providerID := *shipment.LogisticsProviderID
			if provider, exists := providerMap[providerID]; exists {
				shipment.LogisticsProvider = provider
				continue
			}
			provider, providerErr := uc.logisticsProviderRepo.GetByID(providerID)
			if providerErr != nil {
				continue
			}
			providerMap[providerID] = provider
			shipment.LogisticsProvider = provider
		}
	}

	if uc.shippingRateRepo != nil {
		rateMap := make(map[uint64]interface{}, len(shipments))
		for _, shipment := range shipments {
			if shipment == nil || shipment.ShippingRateID == nil || *shipment.ShippingRateID == 0 {
				continue
			}
			rateID := *shipment.ShippingRateID
			if rate, exists := rateMap[rateID]; exists {
				shipment.ShippingRate = rate
				continue
			}
			rate, rateErr := uc.shippingRateRepo.GetByID(rateID)
			if rateErr != nil {
				continue
			}
			rateMap[rateID] = rate
			shipment.ShippingRate = rate
		}
	}

	if uc.shipmentItemRepo == nil {
		return shipments, total, nil
	}

	productIDSet := make(map[uint64]struct{})
	for _, shipment := range shipments {
		if shipment == nil || shipment.ID == 0 {
			continue
		}
		items, itemErr := uc.shipmentItemRepo.GetByShipmentID(shipment.ID)
		if itemErr != nil {
			return nil, 0, itemErr
		}
		shipment.Items = items
		for _, item := range items {
			if item.ProductID == 0 {
				continue
			}
			productIDSet[item.ProductID] = struct{}{}
		}
	}

	if len(productIDSet) == 0 || uc.productRepo == nil {
		return shipments, total, nil
	}

	productIDs := make([]uint64, 0, len(productIDSet))
	for productID := range productIDSet {
		productIDs = append(productIDs, productID)
	}
	products, productErr := uc.productRepo.ListByIDs(productIDs)
	if productErr != nil {
		return nil, 0, productErr
	}

	productMap := make(map[uint64]*productDomain.Product, len(products))
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	for _, shipment := range shipments {
		if shipment == nil {
			continue
		}
		for i := range shipment.Items {
			if product, ok := productMap[shipment.Items[i].ProductID]; ok {
				shipment.Items[i].Product = product
			}
		}
	}

	return shipments, total, nil
}

// Get 获取发货单详情
func (uc *ShipmentUsecase) Get(id uint64) (*domain.Shipment, error) {
	shipment, err := uc.shipmentRepo.GetByID(id)
	if err != nil {
		return nil, ErrShipmentNotFound
	}

	// Load warehouse info
	if uc.warehouseRepo != nil {
		warehouse, err := uc.warehouseRepo.GetByID(shipment.WarehouseID)
		if err == nil {
			shipment.Warehouse = warehouse
		}
	}

	if uc.logisticsProviderRepo != nil && shipment.LogisticsProviderID != nil && *shipment.LogisticsProviderID != 0 {
		provider, providerErr := uc.logisticsProviderRepo.GetByID(*shipment.LogisticsProviderID)
		if providerErr == nil {
			shipment.LogisticsProvider = provider
		}
	}

	if uc.shippingRateRepo != nil && shipment.ShippingRateID != nil && *shipment.ShippingRateID != 0 {
		rate, rateErr := uc.shippingRateRepo.GetByID(*shipment.ShippingRateID)
		if rateErr == nil {
			shipment.ShippingRate = rate
		}
	}

	// Load items
	items, err := uc.shipmentItemRepo.GetByShipmentID(id)
	if err != nil {
		return nil, err
	}

	// Load product data for items
	if len(items) > 0 && uc.productRepo != nil {
		productIDs := make([]uint64, 0, len(items))
		for _, item := range items {
			productIDs = append(productIDs, item.ProductID)
		}

		products, err := uc.productRepo.ListByIDs(productIDs)
		if err == nil {
			// Create a map for quick lookup
			productMap := make(map[uint64]*productDomain.Product)
			for i := range products {
				productMap[products[i].ID] = &products[i]
			}

			// Attach product data to items
			for i := range items {
				if product, ok := productMap[items[i].ProductID]; ok {
					items[i].Product = product
				}
			}
		}
	}

	shipment.Items = items

	return shipment, nil
}

// Create 创建发货单（草稿）
// 创建时验证待发库存是否充足
func (uc *ShipmentUsecase) Create(c *gin.Context, params *domain.CreateShipmentParams) (*domain.Shipment, error) {
	if len(params.Items) == 0 {
		return nil, ErrEmptyItems
	}
	if err := uc.validateShipmentProducts(params.Items); err != nil {
		return nil, err
	}
	productMap, err := uc.loadShipmentProductsMap(params.Items)
	if err != nil {
		return nil, err
	}

	// 检查待发库存是否充足
	if uc.inventoryService != nil {
		for _, item := range params.Items {
			balance, err := uc.inventoryService.GetProductBalance(item.ProductID, params.WarehouseID)
			if err != nil {
				return nil, fmt.Errorf("获取产品 %d 库存失败: %w", item.ProductID, err)
			}
			if balance.PendingShipment < item.QuantityPlanned {
				return nil, newShipmentInsufficientInventoryError(
					fmt.Sprintf("待发库存不足: 产品 %d 待发库存 %d, 需要 %d", item.ProductID, balance.PendingShipment, item.QuantityPlanned),
				)
			}
		}
	}

	// Generate shipment number
	shipmentNumber := generateShipmentNumber()

	// Create shipment
	shipment := &domain.Shipment{
		ShipmentNumber:         shipmentNumber,
		OrderNumber:            params.OrderNumber,
		SalesChannel:           params.SalesChannel,
		WarehouseID:            params.WarehouseID,
		DestinationWarehouseID: params.DestinationWarehouseID,
		DestinationType:        params.DestinationType,
		DestinationName:        params.DestinationName,
		DestinationContact:     params.DestinationContact,
		DestinationPhone:       params.DestinationPhone,
		DestinationAddress:     params.DestinationAddress,
		DestinationCode:        params.DestinationCode,
		LogisticsProviderID:    params.LogisticsProviderID,
		ShippingRateID:         params.ShippingRateID,
		TransportMode:          params.TransportMode,
		Carrier:                params.Carrier,
		TrackingNumber:         params.TrackingNumber,
		ExpectedDeliveryDate:   params.ExpectedDeliveryDate,
		Status:                 domain.ShipmentStatusDraft,
		ReceiptStatus:          domain.ShipmentReceiptStatusPending,
		CreatedBy:              params.OperatorID,
		UpdatedBy:              params.OperatorID,
		Remark:                 params.Remark,
		InternalNotes:          params.InternalNotes,
	}
	if params.BoxCount != nil {
		shipment.BoxCount = *params.BoxCount
	}
	if params.TotalWeight != nil {
		shipment.TotalWeight = *params.TotalWeight
	}
	if params.TotalVolume != nil {
		shipment.TotalVolume = *params.TotalVolume
	}

	if err := uc.shipmentRepo.Create(shipment); err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	// Create items
	items := make([]domain.ShipmentItem, 0, len(params.Items))
	for _, itemParam := range params.Items {
		item := domain.ShipmentItem{
			ShipmentID:      shipment.ID,
			ProductID:       itemParam.ProductID,
			QuantityPlanned: itemParam.QuantityPlanned,
		}
		if itemParam.PackageSpecID != nil {
			item.PackageSpecID = itemParam.PackageSpecID
		}
		if itemParam.BoxQuantity != nil {
			item.BoxQuantity = *itemParam.BoxQuantity
		}
		if itemParam.UnitCost != nil {
			item.UnitCost = *itemParam.UnitCost
		} else if product := productMap[itemParam.ProductID]; product != nil && product.UnitCost != nil && *product.UnitCost > 0 {
			item.UnitCost = *product.UnitCost
		}
		if itemParam.Currency != nil {
			item.Currency = uc.defaultCurrency(*itemParam.Currency)
		} else {
			item.Currency = uc.defaultCurrency("")
		}
		if itemParam.Remark != nil {
			item.Remark = itemParam.Remark
		}
		items = append(items, item)
	}

	if err := uc.shipmentItemRepo.CreateBatch(items); err != nil {
		return nil, fmt.Errorf("failed to create shipment items: %w", err)
	}

	shipment.Items = items
	uc.recordAudit(c, "CREATE", shipment.ID, nil, map[string]any{
		"shipment_number": shipment.ShipmentNumber,
		"warehouse_id":    shipment.WarehouseID,
		"status":          shipment.Status,
		"items_count":     len(items),
	})
	return shipment, nil
}

// Update 编辑发货单
// 草稿允许全量编辑；已确认仅允许调整收货/物流信息和备注，避免破坏已锁定库存。
func (uc *ShipmentUsecase) Update(c *gin.Context, id uint64, params *domain.UpdateShipmentParams) (*domain.Shipment, error) {
	shipment, err := uc.shipmentRepo.GetByID(id)
	if err != nil {
		return nil, ErrShipmentNotFound
	}
	if shipment.Status != domain.ShipmentStatusDraft && shipment.Status != domain.ShipmentStatusConfirmed {
		return nil, newShipmentInvalidStatusError(
			fmt.Sprintf("只有草稿或已确认的发货单允许编辑，当前状态: %s", shipment.Status),
		)
	}

	beforeItems, err := uc.shipmentItemRepo.GetByShipmentID(id)
	if err != nil {
		return nil, err
	}
	beforeSnapshot := uc.buildShipmentAuditSnapshot(shipment, beforeItems)
	productMap := make(map[uint64]*productDomain.Product)

	if shipment.Status == domain.ShipmentStatusDraft {
		if len(params.Items) == 0 {
			return nil, ErrEmptyItems
		}
		if err := uc.validateShipmentProducts(params.Items); err != nil {
			return nil, err
		}
		productMap, err = uc.loadShipmentProductsMap(params.Items)
		if err != nil {
			return nil, err
		}
		nextWarehouseID := shipment.WarehouseID
		if params.WarehouseID != nil && *params.WarehouseID != 0 {
			nextWarehouseID = *params.WarehouseID
		}
		if err := uc.validatePendingShipmentAvailability(params.Items, nextWarehouseID); err != nil {
			return nil, err
		}
		uc.applyShipmentEditableFields(shipment, params)
		shipment.WarehouseID = nextWarehouseID
	} else {
		uc.applyShipmentConfirmedEditableFields(shipment, params)
	}

	shipment.UpdatedBy = params.OperatorID
	if err := uc.shipmentRepo.Update(shipment); err != nil {
		return nil, fmt.Errorf("failed to update shipment: %w", err)
	}

	if shipment.Status == domain.ShipmentStatusDraft {
		if err := uc.shipmentItemRepo.DeleteByShipmentID(shipment.ID); err != nil {
			return nil, fmt.Errorf("failed to replace shipment items: %w", err)
		}
		items := make([]domain.ShipmentItem, 0, len(params.Items))
		for _, itemParam := range params.Items {
			item := domain.ShipmentItem{
				ShipmentID:      shipment.ID,
				ProductID:       itemParam.ProductID,
				QuantityPlanned: itemParam.QuantityPlanned,
			}
			if itemParam.PackageSpecID != nil {
				item.PackageSpecID = itemParam.PackageSpecID
			}
			if itemParam.BoxQuantity != nil {
				item.BoxQuantity = *itemParam.BoxQuantity
			}
			if itemParam.UnitCost != nil {
				item.UnitCost = *itemParam.UnitCost
			} else if product := productMap[itemParam.ProductID]; product != nil && product.UnitCost != nil && *product.UnitCost > 0 {
				item.UnitCost = *product.UnitCost
			}
			if itemParam.Currency != nil {
				item.Currency = uc.defaultCurrency(*itemParam.Currency)
			} else {
				item.Currency = uc.defaultCurrency("")
			}
			if itemParam.Remark != nil {
				item.Remark = itemParam.Remark
			}
			items = append(items, item)
		}
		if err := uc.shipmentItemRepo.CreateBatch(items); err != nil {
			return nil, fmt.Errorf("failed to create shipment items: %w", err)
		}
		shipment.Items = items
	} else {
		shipment.Items = beforeItems
	}

	afterSnapshot := uc.buildShipmentAuditSnapshot(shipment, shipment.Items)
	if !reflect.DeepEqual(beforeSnapshot, afterSnapshot) {
		uc.recordAudit(c, "UPDATE", shipment.ID, beforeSnapshot, afterSnapshot)
	}

	return shipment, nil
}

func (uc *ShipmentUsecase) applyShipmentEditableFields(shipment *domain.Shipment, params *domain.UpdateShipmentParams) {
	if params.OrderNumber != nil {
		shipment.OrderNumber = params.OrderNumber
	}
	if params.SalesChannel != nil {
		shipment.SalesChannel = params.SalesChannel
	}
	if params.DestinationWarehouseID != nil {
		shipment.DestinationWarehouseID = params.DestinationWarehouseID
	}
	if params.DestinationType != nil {
		shipment.DestinationType = params.DestinationType
	}
	if params.DestinationName != nil {
		shipment.DestinationName = params.DestinationName
	}
	if params.DestinationCode != nil {
		shipment.DestinationCode = params.DestinationCode
	}
	if params.DestinationContact != nil {
		shipment.DestinationContact = params.DestinationContact
	}
	if params.DestinationPhone != nil {
		shipment.DestinationPhone = params.DestinationPhone
	}
	if params.DestinationAddress != nil {
		shipment.DestinationAddress = params.DestinationAddress
	}
	if params.LogisticsProviderID != nil {
		shipment.LogisticsProviderID = params.LogisticsProviderID
	}
	if params.ShippingRateID != nil {
		shipment.ShippingRateID = params.ShippingRateID
	}
	if params.TransportMode != nil {
		shipment.TransportMode = params.TransportMode
	}
	if params.Carrier != nil {
		shipment.Carrier = params.Carrier
	}
	if params.TrackingNumber != nil {
		shipment.TrackingNumber = params.TrackingNumber
	}
	if params.ExpectedDeliveryDate != nil {
		shipment.ExpectedDeliveryDate = params.ExpectedDeliveryDate
	}
	if params.BoxCount != nil {
		shipment.BoxCount = *params.BoxCount
	}
	if params.TotalWeight != nil {
		shipment.TotalWeight = *params.TotalWeight
	}
	if params.TotalVolume != nil {
		shipment.TotalVolume = *params.TotalVolume
	}
	if params.Remark != nil {
		shipment.Remark = params.Remark
	}
	if params.InternalNotes != nil {
		shipment.InternalNotes = params.InternalNotes
	}
}

func (uc *ShipmentUsecase) applyShipmentConfirmedEditableFields(shipment *domain.Shipment, params *domain.UpdateShipmentParams) {
	if params.DestinationWarehouseID != nil {
		shipment.DestinationWarehouseID = params.DestinationWarehouseID
	}
	if params.DestinationType != nil {
		shipment.DestinationType = params.DestinationType
	}
	if params.DestinationName != nil {
		shipment.DestinationName = params.DestinationName
	}
	if params.DestinationCode != nil {
		shipment.DestinationCode = params.DestinationCode
	}
	if params.DestinationContact != nil {
		shipment.DestinationContact = params.DestinationContact
	}
	if params.DestinationPhone != nil {
		shipment.DestinationPhone = params.DestinationPhone
	}
	if params.DestinationAddress != nil {
		shipment.DestinationAddress = params.DestinationAddress
	}
	if params.LogisticsProviderID != nil {
		shipment.LogisticsProviderID = params.LogisticsProviderID
	}
	if params.ShippingRateID != nil {
		shipment.ShippingRateID = params.ShippingRateID
	}
	if params.TransportMode != nil {
		shipment.TransportMode = params.TransportMode
	}
	if params.Carrier != nil {
		shipment.Carrier = params.Carrier
	}
	if params.TrackingNumber != nil {
		shipment.TrackingNumber = params.TrackingNumber
	}
	if params.ExpectedDeliveryDate != nil {
		shipment.ExpectedDeliveryDate = params.ExpectedDeliveryDate
	}
	if params.Remark != nil {
		shipment.Remark = params.Remark
	}
	if params.InternalNotes != nil {
		shipment.InternalNotes = params.InternalNotes
	}
}

func (uc *ShipmentUsecase) validatePendingShipmentAvailability(items []domain.CreateShipmentItemParams, warehouseID uint64) error {
	if uc.inventoryService == nil {
		return nil
	}
	for _, item := range items {
		balance, err := uc.inventoryService.GetProductBalance(item.ProductID, warehouseID)
		if err != nil {
			return fmt.Errorf("获取产品 %d 库存失败: %w", item.ProductID, err)
		}
		if balance.PendingShipment < item.QuantityPlanned {
			return newShipmentInsufficientInventoryError(
				fmt.Sprintf("待发库存不足: 产品 %d 待发库存 %d, 需要 %d", item.ProductID, balance.PendingShipment, item.QuantityPlanned),
			)
		}
	}
	return nil
}

func (uc *ShipmentUsecase) buildShipmentAuditSnapshot(shipment *domain.Shipment, items []domain.ShipmentItem) map[string]any {
	if shipment == nil {
		return nil
	}
	snapshot := map[string]any{
		"shipment_number":          shipment.ShipmentNumber,
		"warehouse_id":             shipment.WarehouseID,
		"destination_warehouse_id": valueOfUint64PtrLocal(shipment.DestinationWarehouseID),
		"destination_type":         stringValueLocalDestinationType(shipment.DestinationType),
		"destination_name":         stringValueLocal(shipment.DestinationName),
		"destination_code":         stringValueLocal(shipment.DestinationCode),
		"destination_contact":      stringValueLocal(shipment.DestinationContact),
		"destination_phone":        stringValueLocal(shipment.DestinationPhone),
		"destination_address":      stringValueLocal(shipment.DestinationAddress),
		"logistics_provider_id":    valueOfUint64PtrLocal(shipment.LogisticsProviderID),
		"shipping_rate_id":         valueOfUint64PtrLocal(shipment.ShippingRateID),
		"transport_mode":           stringValueLocal(shipment.TransportMode),
		"carrier":                  stringValueLocal(shipment.Carrier),
		"tracking_number":          stringValueLocal(shipment.TrackingNumber),
		"expected_delivery_date":   stringValueLocal(shipment.ExpectedDeliveryDate),
		"box_count":                shipment.BoxCount,
		"total_weight":             shipment.TotalWeight,
		"total_volume":             shipment.TotalVolume,
		"remark":                   stringValueLocal(shipment.Remark),
		"internal_notes":           stringValueLocal(shipment.InternalNotes),
	}
	if len(items) > 0 {
		itemSnapshots := make([]map[string]any, 0, len(items))
		for _, item := range items {
			itemSnapshots = append(itemSnapshots, map[string]any{
				"product_id":        item.ProductID,
				"quantity_planned":  item.QuantityPlanned,
				"package_spec_id":   valueOfUint64PtrLocal(item.PackageSpecID),
				"box_quantity":      item.BoxQuantity,
				"remark":            stringValueLocal(item.Remark),
			})
		}
		snapshot["items"] = itemSnapshots
	}
	return snapshot
}

func (uc *ShipmentUsecase) validateShipmentProducts(items []domain.CreateShipmentItemParams) error {
	if uc == nil || uc.productRepo == nil || len(items) == 0 {
		return nil
	}

	seen := make(map[uint64]struct{}, len(items))
	productIDs := make([]uint64, 0, len(items))
	for _, item := range items {
		if item.ProductID == 0 {
			continue
		}
		if _, ok := seen[item.ProductID]; ok {
			continue
		}
		seen[item.ProductID] = struct{}{}
		productIDs = append(productIDs, item.ProductID)
	}
	if len(productIDs) == 0 {
		return ErrShipmentProductNotFound
	}

	products, err := uc.productRepo.ListByIDs(productIDs)
	if err != nil {
		return err
	}
	if len(products) != len(productIDs) {
		return ErrShipmentProductNotFound
	}
	for _, product := range products {
		if !isShipmentAllowedProductStatus(product.Status) {
			if strings.TrimSpace(product.SellerSku) != "" {
				return fmt.Errorf("%w: %s", ErrShipmentInactiveProduct, product.SellerSku)
			}
			return ErrShipmentInactiveProduct
		}
	}
	return nil
}

func (uc *ShipmentUsecase) loadShipmentProductsMap(items []domain.CreateShipmentItemParams) (map[uint64]*productDomain.Product, error) {
	result := make(map[uint64]*productDomain.Product)
	if uc == nil || uc.productRepo == nil || len(items) == 0 {
		return result, nil
	}

	seen := make(map[uint64]struct{}, len(items))
	productIDs := make([]uint64, 0, len(items))
	for _, item := range items {
		if item.ProductID == 0 {
			continue
		}
		if _, ok := seen[item.ProductID]; ok {
			continue
		}
		seen[item.ProductID] = struct{}{}
		productIDs = append(productIDs, item.ProductID)
	}
	if len(productIDs) == 0 {
		return result, nil
	}

	products, err := uc.productRepo.ListByIDs(productIDs)
	if err != nil {
		return nil, err
	}
	for i := range products {
		product := products[i]
		result[product.ID] = &product
	}
	return result, nil
}

func isShipmentAllowedProductStatus(status string) bool {
	return status == productDomain.ProductStatusOnSale || status == productDomain.ProductStatusReplenishing
}

func (uc *ShipmentUsecase) defaultCurrency(currency string) string {
	if strings.TrimSpace(currency) != "" {
		return strings.TrimSpace(currency)
	}
	if uc != nil && uc.defaultsProvider != nil {
		if defaultCurrency := strings.TrimSpace(uc.defaultsProvider.GetDefaultBaseCurrency()); defaultCurrency != "" {
			return defaultCurrency
		}
	}
	return "USD"
}

// Confirm 确认发货单 (DRAFT → CONFIRMED)
// 锁定库存，检查原料库存是否充足
func (uc *ShipmentUsecase) Confirm(c *gin.Context, id uint64, params *domain.ConfirmShipmentParams) error {
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if uc.confirmTxManager != nil {
		return uc.confirmTxManager.Run(ctx, func(deps ShipmentStateTransactionalDeps) error {
			return uc.confirmWithDeps(c, id, params, deps)
		})
	}
	return uc.confirmWithDeps(c, id, params, ShipmentStateTransactionalDeps{
		ShipmentRepo:     uc.shipmentRepo,
		ShipmentItemRepo: uc.shipmentItemRepo,
		InventoryService: uc.inventoryService,
	})
}

func (uc *ShipmentUsecase) confirmWithDeps(c *gin.Context, id uint64, params *domain.ConfirmShipmentParams, deps ShipmentStateTransactionalDeps) error {
	shipment, err := deps.ShipmentRepo.GetByID(id)
	if err != nil {
		return ErrShipmentNotFound
	}
	moveCtx := context.Background()
	if c != nil && c.Request != nil {
		moveCtx = c.Request.Context()
	}

	// 只有DRAFT状态才能确认
	if shipment.Status != domain.ShipmentStatusDraft {
		return newShipmentInvalidStatusError(
			fmt.Sprintf("只有草稿状态的发货单才能确认，当前状态: %s", shipment.Status),
		)
	}

	// Load items
	items, err := deps.ShipmentItemRepo.GetByShipmentID(id)
	if err != nil {
		return err
	}

	// 检查待发库存是否充足（只能发待发库存）
	if deps.InventoryService != nil {
		for _, item := range items {
			balance, err := deps.InventoryService.GetProductBalance(item.ProductID, shipment.WarehouseID)
			if err != nil {
				return fmt.Errorf("failed to get inventory balance for product %d: %w", item.ProductID, err)
			}
			if balance.PendingShipment < item.QuantityPlanned {
				return newShipmentInsufficientInventoryError(
					fmt.Sprintf("待发库存不足: 产品 %d 待发库存 %d, 需要 %d", item.ProductID, balance.PendingShipment, item.QuantityPlanned),
				)
			}
		}
		for _, item := range items {
			confirmParams := &inventoryDomain.CreateMovementParams{
				ProductID:    item.ProductID,
				WarehouseID:  shipment.WarehouseID,
				MovementType: inventoryDomain.MovementTypeShipmentAllocate,
				Quantity:     int(item.QuantityPlanned),
				ReferenceType: func() *string {
					s := "SHIPMENT"
					return &s
				}(),
				ReferenceID:     &shipment.ID,
				ReferenceNumber: &shipment.ShipmentNumber,
				UnitCost:        float64Ptr(item.UnitCost),
				OperatorID:      params.OperatorID,
			}
			if _, err := deps.InventoryService.CreateMovement(moveCtx, confirmParams); err != nil {
				return fmt.Errorf("failed to create shipment allocate movement: %w", err)
			}
		}
	}

	// Update shipment
	now := time.Now()
	shipment.Status = domain.ShipmentStatusConfirmed
	shipment.InventoryLocked = true
	shipment.ConfirmedAt = &now
	shipment.ConfirmedBy = params.OperatorID
	shipment.UpdatedBy = params.OperatorID

	if err := deps.ShipmentRepo.Update(shipment); err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}

	uc.recordAudit(c, "CONFIRM", shipment.ID, map[string]any{
		"status":           domain.ShipmentStatusDraft,
		"inventory_locked": false,
	}, map[string]any{
		"status":           shipment.Status,
		"inventory_locked": shipment.InventoryLocked,
	})

	return nil
}

// MarkShipped 标记发货 (CONFIRMED → SHIPPED)
// 执行库存流转: 待出库存 → 物流在途
func (uc *ShipmentUsecase) MarkShipped(c *gin.Context, id uint64, params *domain.MarkShippedParams) error {
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if uc.markShippedTxManager != nil {
		return uc.markShippedTxManager.Run(ctx, func(deps ShipmentMarkShippedTransactionalDeps) error {
			return uc.markShippedWithDeps(c, id, params, deps)
		})
	}
	return uc.markShippedWithDeps(c, id, params, ShipmentMarkShippedTransactionalDeps{
		ShipmentRepo:     uc.shipmentRepo,
		ShipmentItemRepo: uc.shipmentItemRepo,
		InventoryService: uc.inventoryService,
		CostRecorder:     uc.costRecorder,
		LandedRecorder:   uc.landedRecorder,
	})
}

func (uc *ShipmentUsecase) markShippedWithDeps(
	c *gin.Context,
	id uint64,
	params *domain.MarkShippedParams,
	deps ShipmentMarkShippedTransactionalDeps,
) error {
	execUC := *uc
	execUC.shipmentRepo = deps.ShipmentRepo
	execUC.shipmentItemRepo = deps.ShipmentItemRepo
	execUC.inventoryService = deps.InventoryService
	execUC.costRecorder = deps.CostRecorder
	execUC.landedRecorder = deps.LandedRecorder

	shipment, err := execUC.shipmentRepo.GetByID(id)
	if err != nil {
		return ErrShipmentNotFound
	}

	// 只有CONFIRMED状态才能发货
	if shipment.Status != domain.ShipmentStatusConfirmed {
		return newShipmentInvalidStatusError(
			fmt.Sprintf("只有已确认的发货单才能标记发货，当前状态: %s", shipment.Status),
		)
	}

	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}

	// Load items
	items, err := execUC.shipmentItemRepo.GetByShipmentID(id)
	if err != nil {
		return err
	}

	// 创建库存流转: 待出库存 → 物流在途（发货时才真正扣减库存）
	for i := range items {
		item := &items[i]
		if execUC.inventoryService != nil {
			// 使用quantity_shipped，如果没有设置则使用quantity_planned
			quantity := item.QuantityShipped
			if quantity == 0 {
				quantity = item.QuantityPlanned
			}

			// 再次检查待发库存是否充足
			balance, err := execUC.inventoryService.GetProductBalance(item.ProductID, shipment.WarehouseID)
			if err != nil {
				return fmt.Errorf("failed to get inventory balance for product %d: %w", item.ProductID, err)
			}
			if balance.PendingShipmentReserved < quantity {
				return newShipmentInsufficientInventoryError(
					fmt.Sprintf("待发锁定库存不足: 产品 %d 待发锁定库存 %d, 需要 %d", item.ProductID, balance.PendingShipmentReserved, quantity),
				)
			}

			shipParams := &inventoryDomain.CreateMovementParams{
				ProductID:    item.ProductID,
				WarehouseID:  shipment.WarehouseID,
				MovementType: inventoryDomain.MovementTypeShipmentShip,
				Quantity:     int(quantity),
				ReferenceType: func() *string {
					s := "SHIPMENT"
					return &s
				}(),
				ReferenceID:     &shipment.ID,
				ReferenceNumber: &shipment.ShipmentNumber,
				UnitCost:        &item.UnitCost,
				OperatorID:      params.OperatorID,
				Remark:          params.Remark,
			}
			if _, err := execUC.inventoryService.CreateMovement(ctx, shipParams); err != nil {
				return fmt.Errorf("failed to create ship movement: %w", err)
			}
			item.QuantityShipped = quantity
		}
	}
	if err := execUC.shipmentItemRepo.UpdateBatch(items); err != nil {
		return fmt.Errorf("failed to update shipment items: %w", err)
	}

	// Update shipment - 发货时才标记库存已扣减
	shipment.InventoryDeducted = true
	shipment.InventoryLocked = false

	// Update shipment
	now := time.Now()
	shipment.Status = domain.ShipmentStatusShipped
	shipment.ShippedAt = &now
	shipment.ShippedBy = params.OperatorID
	shipment.Carrier = params.Carrier
	shipment.TrackingNumber = params.TrackingNumber
	if params.ShippingCost != nil {
		shipment.ShippingCost = *params.ShippingCost
	}
	if params.Currency != nil || shipment.ShippingCost > 0 {
		shipment.Currency = uc.defaultCurrency(stringValue(params.Currency))
	}
	if shipment.ShippingCost > 0 {
		fxSnapshot, err := execUC.resolveShipmentFXSnapshot(shipment.Currency, now)
		if err != nil {
			return err
		}
		shipment.BaseCurrency = execUC.defaultCurrency("")
		shipment.ShippingCostFxRate = fxSnapshot.Rate
		shipment.ShippingCostFxSource = fxSnapshot.Source
		shipment.ShippingCostFxVersion = fxSnapshot.Version
		shipment.ShippingCostFxTime = &fxSnapshot.EffectiveAt
		shipment.ShippingCostBaseAmount = roundAmount6(shipment.ShippingCost * fxSnapshot.Rate)
	}
	if params.ShipDate != nil {
		shipment.ShipDate = params.ShipDate
	}
	shipment.UpdatedBy = params.OperatorID

	if err := execUC.shipmentRepo.Update(shipment); err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}
	if shipment.ShippingCost > 0 {
		if err := execUC.recordShipmentCostAllocation(shipment, items, now, params.OperatorID); err != nil {
			return err
		}
	}

	execUC.recordAudit(c, "SHIP", shipment.ID, map[string]any{
		"status":               domain.ShipmentStatusConfirmed,
		"inventory_locked":     true,
		"inventory_deducted":   false,
		"quantity_shipped_set": false,
	}, map[string]any{
		"status":             shipment.Status,
		"inventory_locked":   shipment.InventoryLocked,
		"inventory_deducted": shipment.InventoryDeducted,
		"shipping_cost":      shipment.ShippingCost,
		"currency":           shipment.Currency,
		"base_currency":      shipment.BaseCurrency,
	})

	return nil
}

// MarkDelivered 标记送达 (SHIPPED → DELIVERED)
// 注意: 这里不执行库存流转，因为送达到平台上架是另一个流程
func (uc *ShipmentUsecase) MarkDelivered(c *gin.Context, id uint64, params *domain.MarkDeliveredParams) error {
	shipment, err := uc.shipmentRepo.GetByID(id)
	if err != nil {
		return ErrShipmentNotFound
	}

	// 只有SHIPPED状态才能标记送达
	if shipment.Status != domain.ShipmentStatusShipped {
		return newShipmentInvalidStatusError(
			fmt.Sprintf("只有已发货的发货单才能标记送达，当前状态: %s", shipment.Status),
		)
	}

	// Update shipment
	now := time.Now()
	shipment.Status = domain.ShipmentStatusDelivered
	shipment.DeliveredAt = &now
	shipment.DeliveredBy = params.OperatorID
	if params.ActualDeliveryDate != nil {
		shipment.ActualDeliveryDate = params.ActualDeliveryDate
	}
	shipment.UpdatedBy = params.OperatorID
	if params.Remark != nil {
		shipment.Remark = params.Remark
	}

	if err := uc.shipmentRepo.Update(shipment); err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}

	uc.recordAudit(c, "DELIVER", shipment.ID, map[string]any{
		"status": domain.ShipmentStatusShipped,
	}, map[string]any{
		"status": shipment.Status,
	})

	return nil
}

// Cancel 取消发货单（带回滚）
// 根据当前状态决定回滚策略
func (uc *ShipmentUsecase) Cancel(c *gin.Context, id uint64, params *domain.CancelShipmentParams) error {
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if uc.cancelTxManager != nil {
		return uc.cancelTxManager.Run(ctx, func(deps ShipmentStateTransactionalDeps) error {
			return uc.cancelWithDeps(c, id, params, deps)
		})
	}
	return uc.cancelWithDeps(c, id, params, ShipmentStateTransactionalDeps{
		ShipmentRepo:     uc.shipmentRepo,
		ShipmentItemRepo: uc.shipmentItemRepo,
		InventoryService: uc.inventoryService,
	})
}

func (uc *ShipmentUsecase) cancelWithDeps(c *gin.Context, id uint64, params *domain.CancelShipmentParams, deps ShipmentStateTransactionalDeps) error {
	shipment, err := deps.ShipmentRepo.GetByID(id)
	if err != nil {
		return ErrShipmentNotFound
	}
	beforeStatus := shipment.Status
	beforeLocked := shipment.InventoryLocked

	// SHIPPED和DELIVERED状态不允许取消
	if shipment.Status == domain.ShipmentStatusShipped || shipment.Status == domain.ShipmentStatusDelivered {
		return newShipmentInvalidStatusError(
			fmt.Sprintf("货物已发出，无法取消，当前状态: %s", shipment.Status),
		)
	}

	// 如果已经是CANCELLED状态，直接返回
	if shipment.Status == domain.ShipmentStatusCancelled {
		return nil
	}

	// ctx := c.Request.Context()

	// 根据状态决定回滚策略
	switch shipment.Status {
	case domain.ShipmentStatusDraft:
		// 草稿直接取消，无需回滚库存

	case domain.ShipmentStatusConfirmed:
		items, err := deps.ShipmentItemRepo.GetByShipmentID(id)
		if err != nil {
			return err
		}
		for _, item := range items {
			cancelParams := &inventoryDomain.CreateMovementParams{
				ProductID:    item.ProductID,
				WarehouseID:  shipment.WarehouseID,
				MovementType: inventoryDomain.MovementTypeShipmentRelease,
				Quantity:     int(item.QuantityPlanned),
				ReferenceType: func() *string {
					s := "SHIPMENT"
					return &s
				}(),
				ReferenceID:     &shipment.ID,
				ReferenceNumber: &shipment.ShipmentNumber,
				UnitCost:        float64Ptr(item.UnitCost),
				OperatorID:      params.OperatorID,
				Remark:          params.Remark,
			}
			moveCtx := context.Background()
			if c != nil && c.Request != nil {
				moveCtx = c.Request.Context()
			}
			if _, err := deps.InventoryService.CreateMovement(moveCtx, cancelParams); err != nil {
				return fmt.Errorf("failed to create shipment release movement: %w", err)
			}
		}
		shipment.InventoryLocked = false
	}

	// Update shipment
	shipment.Status = domain.ShipmentStatusCancelled
	shipment.UpdatedBy = params.OperatorID
	if params.Remark != nil {
		shipment.Remark = params.Remark
	}

	if err := deps.ShipmentRepo.Update(shipment); err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}

	uc.recordAudit(c, "CANCEL", shipment.ID, map[string]any{
		"status":           beforeStatus,
		"inventory_locked": beforeLocked,
	}, map[string]any{
		"status":           domain.ShipmentStatusCancelled,
		"inventory_locked": false,
	})

	return nil
}

func (uc *ShipmentUsecase) resolveShipmentFXSnapshot(originalCurrency string, occurredAt time.Time) (*ShipmentFXSnapshot, error) {
	baseCurrency := uc.defaultCurrency("")
	originalCurrency = uc.defaultCurrency(originalCurrency)
	if strings.EqualFold(baseCurrency, originalCurrency) {
		return &ShipmentFXSnapshot{
			Rate:        1,
			Source:      "IDENTITY",
			Version:     "same_currency",
			EffectiveAt: occurredAt,
		}, nil
	}
	if uc.fxResolver == nil {
		return nil, fmt.Errorf("shipment fx resolver not configured")
	}
	return uc.fxResolver(baseCurrency, originalCurrency, occurredAt)
}

func (uc *ShipmentUsecase) recordShipmentCostAllocation(shipment *domain.Shipment, items []domain.ShipmentItem, occurredAt time.Time, operatorID *uint64) error {
	if uc.costRecorder == nil || shipment == nil || shipment.ID == 0 || shipment.ShippingCost <= 0 {
		return nil
	}

	totalQty := uint64(0)
	lines := make([]ShipmentCostAllocationLine, 0, len(items))
	for _, item := range items {
		qty := uint64(item.QuantityShipped)
		if qty == 0 {
			qty = uint64(item.QuantityPlanned)
		}
		if item.ProductID == 0 || qty == 0 {
			continue
		}
		lines = append(lines, ShipmentCostAllocationLine{
			ShipmentItemID: item.ID,
			ProductID:      item.ProductID,
			WarehouseID:    shipment.WarehouseID,
			Quantity:       qty,
			ItemUnitCost:   item.UnitCost,
			ItemCurrency:   item.Currency,
		})
		totalQty += qty
	}
	if totalQty == 0 || len(lines) == 0 {
		return nil
	}

	remainingOriginal := roundAmount6(shipment.ShippingCost)
	remainingBase := roundAmount6(shipment.ShippingCostBaseAmount)
	for idx := range lines {
		line := &lines[idx]
		if idx == len(lines)-1 {
			line.OriginalAmount = remainingOriginal
			line.BaseAmount = remainingBase
			break
		}
		ratio := float64(line.Quantity) / float64(totalQty)
		line.OriginalAmount = roundAmount6(shipment.ShippingCost * ratio)
		line.BaseAmount = roundAmount6(shipment.ShippingCostBaseAmount * ratio)
		remainingOriginal = roundAmount6(remainingOriginal - line.OriginalAmount)
		remainingBase = roundAmount6(remainingBase - line.BaseAmount)
	}

	params := &ShipmentCostAllocationRecordParams{
		ShipmentID:       shipment.ID,
		ShipmentNumber:   shipment.ShipmentNumber,
		OriginalCurrency: shipment.Currency,
		BaseCurrency:     shipment.BaseCurrency,
		FxRate:           shipment.ShippingCostFxRate,
		FxSource:         shipment.ShippingCostFxSource,
		FxVersion:        shipment.ShippingCostFxVersion,
		FxTime:           timeValue(shipment.ShippingCostFxTime, occurredAt),
		OccurredAt:       occurredAt,
		OperatorID:       operatorID,
		Lines:            lines,
	}
	if err := uc.costRecorder.RecordShipmentCostAllocation(params); err != nil {
		return err
	}
	if uc.landedRecorder != nil {
		if err := uc.landedRecorder.UpsertShipmentLandedSnapshots(params); err != nil {
			return err
		}
	}
	return nil
}

func roundAmount6(v float64) float64 {
	return float64(int64(v*1_000_000+0.5)) / 1_000_000
}

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func stringValueLocal(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func stringValueLocalDestinationType(v *domain.DestinationType) string {
	if v == nil {
		return ""
	}
	return string(*v)
}

func valueOfUint64PtrLocal(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func timeValue(v *time.Time, fallback time.Time) time.Time {
	if v == nil || v.IsZero() {
		return fallback
	}
	return *v
}

func (uc *ShipmentUsecase) recordAudit(c *gin.Context, action string, shipmentID uint64, before any, after any) {
	if uc.auditLogger == nil || c == nil || shipmentID == 0 {
		return
	}
	_ = uc.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Shipment",
		Action:     action,
		EntityType: "Shipment",
		EntityID:   fmt.Sprintf("%d", shipmentID),
		Before:     before,
		After:      after,
	})
}

func float64Ptr(v float64) *float64 {
	return &v
}

// Delete 删除发货单
// 只能删除DRAFT或CANCELLED状态的发货单
func (uc *ShipmentUsecase) Delete(c *gin.Context, id uint64) error {
	shipment, err := uc.shipmentRepo.GetByID(id)
	if err != nil {
		return ErrShipmentNotFound
	}

	// 只有DRAFT或CANCELLED状态才能删除
	if shipment.Status != domain.ShipmentStatusDraft && shipment.Status != domain.ShipmentStatusCancelled {
		return newShipmentInvalidStatusError(
			fmt.Sprintf("只有草稿或已取消的发货单才能删除，当前状态: %s", shipment.Status),
		)
	}

	// Delete items
	if err := uc.shipmentItemRepo.DeleteByShipmentID(id); err != nil {
		return fmt.Errorf("failed to delete shipment items: %w", err)
	}

	// Delete shipment
	if err := uc.shipmentRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete shipment: %w", err)
	}

	return nil
}

// generateShipmentNumber 生成发货单号
func generateShipmentNumber() string {
	return numbering.Generate("SH", time.Now())
}
