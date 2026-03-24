package domain

type SalesOrderRepository interface {
	List(params *SalesOrderListParams) ([]SalesOrder, int64, error)
	GetByID(id uint64) (*SalesOrder, error)
	Create(order *SalesOrder) error
	Update(order *SalesOrder) error
}

type SalesImportRepository interface {
	GetImportByFileHash(fileHash string) (*ReportImport, error)
	CreateImport(batch *ReportImport) error
	UpdateImport(batch *ReportImport) error
	ListImports(page int, pageSize int) ([]ReportImport, int64, error)
	GetImportByID(id uint64) (*ReportImport, error)
	ListImportErrors(importID uint64) ([]ReportImportRowError, error)

	InsertImportRowErrors(errors []ReportImportRowError) error
	ResolveProductIDBySellerSKU(sellerSKU string, marketplace string) (uint64, error)
	UpsertImportedOrderLine(line *ImportOrderLine, productID uint64, batchNo string, operatorID *uint64) error
}
