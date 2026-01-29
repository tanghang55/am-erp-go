package domain

import "time"

// PackageSpec 装箱规格
type PackageSpec struct {
	ID             uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string    `json:"name" gorm:"column:name;size:100;not null"`              // 名称
	Length         float64   `json:"length" gorm:"column:length;type:decimal(10,2)"`         // 长(cm)
	Width          float64   `json:"width" gorm:"column:width;type:decimal(10,2)"`           // 宽(cm)
	Height         float64   `json:"height" gorm:"column:height;type:decimal(10,2)"`         // 高(cm)
	Weight         float64   `json:"weight" gorm:"column:weight;type:decimal(10,2)"`         // 重量(kg)
	QuantityPerBox uint      `json:"quantity_per_box" gorm:"column:quantity_per_box;default:1"` // 每箱数量
	Remark         *string   `json:"remark" gorm:"column:remark;size:500"`                   // 备注
	Status         string    `json:"status" gorm:"column:status;size:20;default:'ACTIVE'"`   // 状态
	CreatedBy   *uint64   `json:"created_by" gorm:"column:created_by"`
	UpdatedBy   *uint64   `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate   time.Time `json:"created_at" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"updated_at" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (PackageSpec) TableName() string {
	return "package_spec"
}

// Volume 计算体积 (立方米)
func (p *PackageSpec) Volume() float64 {
	return (p.Length * p.Width * p.Height) / 1000000 // cm³ -> m³
}

// PackageSpecListParams 列表查询参数
type PackageSpecListParams struct {
	Page     int
	PageSize int
	Keyword  *string
	Status   *string
}

// CreatePackageSpecParams 创建参数
type CreatePackageSpecParams struct {
	Name           string
	Length         float64
	Width          float64
	Height         float64
	Weight         float64
	QuantityPerBox uint
	Remark         *string
	CreatedBy      *uint64
}

// UpdatePackageSpecParams 更新参数
type UpdatePackageSpecParams struct {
	Name           *string
	Length         *float64
	Width          *float64
	Height         *float64
	Weight         *float64
	QuantityPerBox *uint
	Remark         *string
	Status         *string
	UpdatedBy      *uint64
}

// PackageSpecRepository 仓储接口
type PackageSpecRepository interface {
	Create(spec *PackageSpec) error
	Update(spec *PackageSpec) error
	GetByID(id uint64) (*PackageSpec, error)
	List(params *PackageSpecListParams) ([]*PackageSpec, int64, error)
	Delete(id uint64) error
	ListByIDs(ids []uint64) ([]*PackageSpec, error)
}

// PackageSpecPackagingItem 装箱规格包材关联实体
type PackageSpecPackagingItem struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	PackageSpecID   uint64    `json:"package_spec_id" gorm:"column:package_spec_id;not null"`
	PackagingItemID uint64    `json:"packaging_item_id" gorm:"column:packaging_item_id;not null"`
	QuantityPerBox  float64   `json:"quantity_per_box" gorm:"column:quantity_per_box;type:decimal(10,3);not null"`
	Notes           string    `json:"notes" gorm:"size:500"`
	CreatedBy       *uint64   `json:"created_by" gorm:"column:created_by"`
	UpdatedBy       *uint64   `json:"updated_by" gorm:"column:updated_by"`
	GmtCreate       time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified     time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`

	// 关联包材详情（查询时关联）
	PackagingItem   *PackagingItemDetail `json:"packaging_item,omitempty" gorm:"-"`
}

func (PackageSpecPackagingItem) TableName() string {
	return "package_spec_packaging_items"
}

// PackagingItemDetail 包材详情（用于关联查询）
type PackagingItemDetail struct {
	ID             uint64  `json:"id"`
	ItemCode       string  `json:"item_code"`
	ItemName       string  `json:"item_name"`
	Specification  string  `json:"specification"`
	Unit           string  `json:"unit"`
	UnitCost       float64 `json:"unit_cost"`
	Currency       string  `json:"currency"`
	QuantityOnHand uint64  `json:"quantity_on_hand"`
}

// PackageSpecPackagingRepository 装箱规格包材关联仓储接口
type PackageSpecPackagingRepository interface {
	// 获取装箱规格的包材配置列表
	ListByPackageSpecID(packageSpecID uint64) ([]PackageSpecPackagingItem, error)
	// 替换装箱规格的包材配置（先删除后插入）
	ReplaceAll(packageSpecID uint64, items []PackageSpecPackagingItem) error
}
