package domain

import "context"

type WarehouseRepository interface {
	List(params *WarehouseListParams) ([]*Warehouse, int64, error)
	GetByID(id uint64) (*Warehouse, error)
	Create(ctx context.Context, warehouse *Warehouse) error
	Update(ctx context.Context, warehouse *Warehouse) error
	Delete(ctx context.Context, id uint64) error
	GetActiveWarehouses() ([]*Warehouse, error)
	CountReferences(id uint64) (int64, error)
}

type InventoryBalanceRepository interface {
	GetOrCreate(ctx context.Context, productID, warehouseID uint64) (*InventoryBalance, error)
	Update(ctx context.Context, balance *InventoryBalance) error
	List(params *BalanceListParams) ([]*InventoryBalance, int64, error)
	GetByProductAndWarehouse(productID, warehouseID uint64) (*InventoryBalance, error)
}

type InventoryMovementRepository interface {
	Create(ctx context.Context, movement *InventoryMovement) error
	List(params *MovementListParams) ([]*InventoryMovement, int64, error)
	GetByID(id uint64) (*InventoryMovement, error)
}

type InventoryLotRepository interface {
	List(params *InventoryLotListParams) ([]*InventoryLot, int64, error)
	ListByProductAndWarehouse(ctx context.Context, productID, warehouseID uint64) ([]*InventoryLot, error)
	Create(ctx context.Context, lot *InventoryLot) error
	Update(ctx context.Context, lot *InventoryLot) error
}
