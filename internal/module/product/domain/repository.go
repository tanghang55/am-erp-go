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
	Page              int
	PageSize          int
	Keyword           string
	Marketplace       string
	Status            string
	SupplierID        *uint64
	ComboID           *uint64
	IsComboMain       *uint8
	ExcludeComboChild bool    // 排除组合产品的子产品（只返回combo_id IS NULL的SKU，用于发货单选择）
	WarehouseID       *uint64 // 仓库ID，如果指定则只返回该仓库有库存的产品
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

// ProductPackagingRepository 产品包材关联仓储接口
type ProductPackagingRepository interface {
	// 获取产品的包材配置列表
	ListByProductID(productID uint64) ([]ProductPackagingItem, error)
	// 替换产品的包材配置（先删除后插入）
	ReplaceAll(productID uint64, items []ProductPackagingItem) error
}
