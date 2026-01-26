package domain

import "context"

type WarehouseRepository interface {
	List(params *WarehouseListParams) ([]*Warehouse, int64, error)
	GetByID(id uint64) (*Warehouse, error)
	Create(ctx context.Context, warehouse *Warehouse) error
	Update(ctx context.Context, warehouse *Warehouse) error
	Delete(ctx context.Context, id uint64) error
	GetActiveWarehouses() ([]*Warehouse, error)
}

type InventoryBalanceRepository interface {
	GetOrCreate(ctx context.Context, skuID, warehouseID uint64) (*InventoryBalance, error)
	Update(ctx context.Context, balance *InventoryBalance) error
	List(params *BalanceListParams) ([]*InventoryBalance, int64, error)
	GetBySkuAndWarehouse(skuID, warehouseID uint64) (*InventoryBalance, error)
}

type InventoryMovementRepository interface {
	Create(ctx context.Context, movement *InventoryMovement) error
	List(params *MovementListParams) ([]*InventoryMovement, int64, error)
	GetByID(id uint64) (*InventoryMovement, error)
}
