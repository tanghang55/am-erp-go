package domain

// SupplierRepository defines supplier persistence operations.
type SupplierRepository interface {
	List(params *SupplierListParams) ([]Supplier, int64, error)
	GetByID(id uint64) (*Supplier, error)
	Create(supplier *Supplier) error
	Update(supplier *Supplier) error
	Delete(id uint64) error
}

type SupplierTypeRepository interface {
	ListBySupplierID(id uint64) ([]string, error)
	ListBySupplierIDs(ids []uint64) (map[uint64][]string, error)
	ReplaceBySupplierID(id uint64, types []string) error
}

type SupplierContactRepository interface {
	ListBySupplierID(id uint64) ([]SupplierContact, error)
	Create(contact *SupplierContact) error
	Update(contact *SupplierContact) error
	Delete(id uint64, supplierID uint64) error
}

type SupplierAccountRepository interface {
	ListBySupplierID(id uint64) ([]SupplierAccount, error)
	Create(account *SupplierAccount) error
	Update(account *SupplierAccount) error
	Delete(id uint64, supplierID uint64) error
}

type SupplierTagRepository interface {
	ListBySupplierID(id uint64) ([]SupplierTag, error)
	Create(tag *SupplierTag) error
	Update(tag *SupplierTag) error
	Delete(id uint64, supplierID uint64) error
}

type QuoteRepository interface {
	ListByProductIDs(productIDs []uint64) (map[uint64][]ProductSupplierQuote, error)
	ListProductsWithQuotes(params *QuoteListParams) ([]ProductQuoteRow, int64, error)
	GetByProductSupplier(productID, supplierID uint64) (*ProductSupplierQuote, error)
	Create(quote *ProductSupplierQuote) error
	Update(quote *ProductSupplierQuote) error
	Delete(productID, supplierID uint64) error
}

// SupplierListParams defines list filters for suppliers.
type SupplierListParams struct {
	Page     int
	PageSize int
	Keyword  string
	Status   string
	Types    []string
}
