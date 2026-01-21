package domain

// PurchaseOrderRepository 采购单仓储接口
type PurchaseOrderRepository interface {
	List(params *PurchaseOrderListParams) ([]PurchaseOrder, int64, error)
	GetByID(id uint64) (*PurchaseOrder, error)
	Create(order *PurchaseOrder) error
	Update(order *PurchaseOrder) error
	Delete(id uint64) error
}
