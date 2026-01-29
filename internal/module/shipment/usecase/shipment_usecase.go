package usecase

import (
	"context"
	"errors"
	"fmt"
	mathrand "math/rand"
	"time"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	productDomain "am-erp-go/internal/module/product/domain"
	"am-erp-go/internal/module/shipment/domain"

	"github.com/gin-gonic/gin"
)

var (
	ErrShipmentNotFound      = errors.New("shipment not found")
	ErrInvalidStatus         = errors.New("invalid shipment status")
	ErrShipmentAlreadyExists = errors.New("shipment already exists")
	ErrEmptyItems            = errors.New("shipment must have at least one item")
	ErrInsufficientInventory = errors.New("insufficient inventory")
)

// InventoryService interface for inventory operations
type InventoryService interface {
	CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error)
	GetSkuBalance(skuID, warehouseID uint64) (*inventoryDomain.InventoryBalance, error)
}

// ProductRepository interface for loading product info
type ProductRepository interface {
	ListByIDs(ids []uint64) ([]productDomain.Product, error)
}

// WarehouseRepository interface for loading warehouse info
type WarehouseRepository interface {
	GetByID(id uint64) (*inventoryDomain.Warehouse, error)
}

type ShipmentUsecase struct {
	shipmentRepo     domain.ShipmentRepository
	shipmentItemRepo domain.ShipmentItemRepository
	inventoryService InventoryService
	productRepo      ProductRepository
	warehouseRepo    WarehouseRepository
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

// List 获取发货单列表
func (uc *ShipmentUsecase) List(params *domain.ShipmentListParams) ([]*domain.Shipment, int64, error) {
	return uc.shipmentRepo.List(params)
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

	// Load items
	items, err := uc.shipmentItemRepo.GetByShipmentID(id)
	if err != nil {
		return nil, err
	}

	// Load SKU data for items
	if len(items) > 0 && uc.productRepo != nil {
		skuIDs := make([]uint64, 0, len(items))
		for _, item := range items {
			skuIDs = append(skuIDs, item.SkuID)
		}

		products, err := uc.productRepo.ListByIDs(skuIDs)
		if err == nil {
			// Create a map for quick lookup
			productMap := make(map[uint64]*productDomain.Product)
			for i := range products {
				productMap[products[i].ID] = &products[i]
			}

			// Attach SKU data to items
			for i := range items {
				if product, ok := productMap[items[i].SkuID]; ok {
					items[i].Sku = product
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

	// 检查待发库存是否充足
	if uc.inventoryService != nil {
		for _, item := range params.Items {
			balance, err := uc.inventoryService.GetSkuBalance(item.SkuID, params.WarehouseID)
			if err != nil {
				return nil, fmt.Errorf("获取 SKU %d 库存失败: %w", item.SkuID, err)
			}
			if balance.PendingShipment < item.QuantityPlanned {
				return nil, fmt.Errorf("待发库存不足: SKU %d 待发库存 %d, 需要 %d", item.SkuID, balance.PendingShipment, item.QuantityPlanned)
			}
		}
	}

	// Generate shipment number
	shipmentNumber := generateShipmentNumber()

	// Create shipment
	shipment := &domain.Shipment{
		ShipmentNumber: shipmentNumber,
		OrderNumber:    params.OrderNumber,
		WarehouseID:    params.WarehouseID,
		Status:         domain.ShipmentStatusDraft,
		CreatedBy:      params.OperatorID,
		UpdatedBy:      params.OperatorID,
		Remark:         params.Remark,
	}

	if err := uc.shipmentRepo.Create(shipment); err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	// Create items
	items := make([]domain.ShipmentItem, 0, len(params.Items))
	for _, itemParam := range params.Items {
		item := domain.ShipmentItem{
			ShipmentID:      shipment.ID,
			SkuID:           itemParam.SkuID,
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
		}
		if itemParam.Currency != nil {
			item.Currency = *itemParam.Currency
		} else {
			item.Currency = "USD"
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
	return shipment, nil
}

// Confirm 确认发货单 (DRAFT → CONFIRMED)
// 锁定库存，检查原料库存是否充足
func (uc *ShipmentUsecase) Confirm(c *gin.Context, id uint64, params *domain.ConfirmShipmentParams) error {
	shipment, err := uc.shipmentRepo.GetByID(id)
	if err != nil {
		return ErrShipmentNotFound
	}

	// 只有DRAFT状态才能确认
	if shipment.Status != domain.ShipmentStatusDraft {
		return fmt.Errorf("只有草稿状态的发货单才能确认，当前状态: %s", shipment.Status)
	}

	// Load items
	items, err := uc.shipmentItemRepo.GetByShipmentID(id)
	if err != nil {
		return err
	}

	// 检查待发库存是否充足（只能发待发库存）
	if uc.inventoryService != nil {
		for _, item := range items {
			balance, err := uc.inventoryService.GetSkuBalance(item.SkuID, shipment.WarehouseID)
			if err != nil {
				return fmt.Errorf("failed to get inventory balance for SKU %d: %w", item.SkuID, err)
			}
			if balance.PendingShipment < item.QuantityPlanned {
				return fmt.Errorf("待发库存不足: SKU %d 待发库存 %d, 需要 %d", item.SkuID, balance.PendingShipment, item.QuantityPlanned)
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

	if err := uc.shipmentRepo.Update(shipment); err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}

	return nil
}

// MarkShipped 标记发货 (CONFIRMED → SHIPPED)
// 执行库存流转: 待出库存 → 物流在途
func (uc *ShipmentUsecase) MarkShipped(c *gin.Context, id uint64, params *domain.MarkShippedParams) error {
	shipment, err := uc.shipmentRepo.GetByID(id)
	if err != nil {
		return ErrShipmentNotFound
	}

	// 只有CONFIRMED状态才能发货
	if shipment.Status != domain.ShipmentStatusConfirmed {
		return fmt.Errorf("只有已确认的发货单才能标记发货，当前状态: %s", shipment.Status)
	}

	ctx := c.Request.Context()

	// Load items
	items, err := uc.shipmentItemRepo.GetByShipmentID(id)
	if err != nil {
		return err
	}

	// 创建库存流转: 待出库存 → 物流在途（发货时才真正扣减库存）
	for _, item := range items {
		if uc.inventoryService != nil {
			// 使用quantity_shipped，如果没有设置则使用quantity_planned
			quantity := item.QuantityShipped
			if quantity == 0 {
				quantity = item.QuantityPlanned
			}

			// 再次检查待发库存是否充足
			balance, err := uc.inventoryService.GetSkuBalance(item.SkuID, shipment.WarehouseID)
			if err != nil {
				return fmt.Errorf("failed to get inventory balance for SKU %d: %w", item.SkuID, err)
			}
			if balance.PendingShipment < quantity {
				return fmt.Errorf("待发库存不足: SKU %d 待发库存 %d, 需要 %d", item.SkuID, balance.PendingShipment, quantity)
			}

			shipParams := &inventoryDomain.CreateMovementParams{
				SkuID:        item.SkuID,
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
			if _, err := uc.inventoryService.CreateMovement(ctx, shipParams); err != nil {
				return fmt.Errorf("failed to create ship movement: %w", err)
			}
		}
	}

	// Update shipment - 发货时才标记库存已扣减
	shipment.InventoryDeducted = true

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
	if params.Currency != nil {
		shipment.Currency = *params.Currency
	}
	if params.ShipDate != nil {
		shipment.ShipDate = params.ShipDate
	}
	shipment.UpdatedBy = params.OperatorID

	if err := uc.shipmentRepo.Update(shipment); err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}

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
		return fmt.Errorf("只有已发货的发货单才能标记送达，当前状态: %s", shipment.Status)
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

	return nil
}

// Cancel 取消发货单（带回滚）
// 根据当前状态决定回滚策略
func (uc *ShipmentUsecase) Cancel(c *gin.Context, id uint64, params *domain.CancelShipmentParams) error {
	shipment, err := uc.shipmentRepo.GetByID(id)
	if err != nil {
		return ErrShipmentNotFound
	}

	// SHIPPED和DELIVERED状态不允许取消
	if shipment.Status == domain.ShipmentStatusShipped || shipment.Status == domain.ShipmentStatusDelivered {
		return fmt.Errorf("货物已发出，无法取消，当前状态: %s", shipment.Status)
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
		// 已确认但未发货，只需解除锁定
		shipment.InventoryLocked = false
	}

	// Update shipment
	shipment.Status = domain.ShipmentStatusCancelled
	shipment.UpdatedBy = params.OperatorID
	if params.Remark != nil {
		shipment.Remark = params.Remark
	}

	if err := uc.shipmentRepo.Update(shipment); err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}

	return nil
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
		return fmt.Errorf("只有草稿或已取消的发货单才能删除，当前状态: %s", shipment.Status)
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
// 格式: SH + 年月日时分秒 + 毫秒(3位) + 随机数(2位)
func generateShipmentNumber() string {
	now := time.Now()
	ms := now.Nanosecond() / 1000000 // 毫秒 0-999
	rnd := mathrand.Intn(100)        // 随机数 0-99
	return fmt.Sprintf("SH%s%03d%02d", now.Format("20060102150405"), ms, rnd)
}
