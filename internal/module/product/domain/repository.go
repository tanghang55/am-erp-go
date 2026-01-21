package domain

// ProductRepository 产品仓储接口
type ProductRepository interface {
	List(params *ProductListParams) ([]Product, int64, error)
	GetByID(id uint64) (*Product, error)
	ListByIDs(ids []uint64) ([]Product, error)
	Create(product *Product) error
	Update(product *Product) error
	Delete(id uint64) error
	UpdateImageUrl(id uint64, imageUrl string) error
	GetDefaultSupplierID(productID uint64) (uint64, error)
	UpdateDefaultSupplierID(productID, supplierID uint64) error
	UpdateComboInfo(comboID uint64, mainProductID uint64, productIDs []uint64) error
	ClearComboInfo(comboID uint64) error
}

// ProductListParams 产品列表查询参数
type ProductListParams struct {
	Page        int
	PageSize    int
	Keyword     string
	Marketplace string
	Status      string
	SupplierID  *uint64
	ComboID     *uint64
	IsComboMain *uint8
}

// ProductParentRepository 产品父体仓储接口
type ProductParentRepository interface {
	List(params *ProductParentListParams) ([]ProductParent, int64, error)
	GetByID(id uint64) (*ProductParent, error)
	Create(parent *ProductParent) error
	Update(parent *ProductParent) error
	Delete(id uint64) error
}

// ProductParentListParams 产品父体列表查询参数
type ProductParentListParams struct {
	Page        int
	PageSize    int
	Keyword     string
	Marketplace string
	Status      string
}

// ProductImageRepository 产品图片仓储接口
type ProductImageRepository interface {
	ListByProductID(productID uint64) ([]ProductImage, error)
	ReplaceAll(productID uint64, orderedUrls []string) error
}

// ProductComboRepository 产品组合仓储接口
type ProductComboRepository interface {
	ListComboIDs(params *ComboListParams) ([]uint64, int64, error)
	GetItemsByComboID(comboID uint64) ([]ProductComboItem, error)
	GetComboIDByMainProductID(mainProductID uint64) (uint64, error)
	CreateCombo(mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) (uint64, error)
	ReplaceComboItems(comboID uint64, mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) error
	DeleteCombo(comboID uint64) error
}
