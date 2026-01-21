package domain

type WarehouseRepository interface {
	GetActiveWarehouses() ([]*Warehouse, error)
}
