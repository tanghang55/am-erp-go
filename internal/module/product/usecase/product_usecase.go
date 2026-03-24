package usecase

import (
	"context"
	"errors"
	"strings"

	"am-erp-go/internal/infrastructure/validation"
	"am-erp-go/internal/module/product/domain"
	supplierdomain "am-erp-go/internal/module/supplier/domain"
)

var ErrSupplierRequired = errors.New("supplier is required")
var ErrProductAsinInvalid = errors.New("asin length must be between 1 and 20")
var ErrProductParentHasChildren = errors.New("product parent still has children")
var ErrProductParentImmutableIdentity = errors.New("product parent asin and marketplace cannot be changed")
var ErrProductParentChildIDsRequired = errors.New("child_ids is required")
var ErrProductParentChildNotFound = errors.New("product child not found")
var ErrProductParentChildAlreadyAssigned = errors.New("product child already belongs to another parent")
var ErrProductParentChildMarketplaceMismatch = errors.New("product child marketplace does not match parent")
var ErrProductCategoryParentRequired = errors.New("parent category is required")
var ErrProductCategoryParentInvalid = errors.New("parent category level is invalid")
var ErrProductCategoryHasChildren = errors.New("product category still has children")
var ErrProductCategoryReferenced = errors.New("product category is still referenced by products")
var ErrProductConfigTypeUnsupported = errors.New("unsupported product config type")
var ErrProductConfigCodeInvalid = errors.New("config code only supports letters, numbers, hyphen and underscore")
var ErrProductConfigReferenced = errors.New("product config is still referenced by products")
var ErrProductConfigCodeImmutable = errors.New("product config code cannot be changed after being referenced")
var ErrProductConfigSystemFixed = errors.New("product config is system fixed")
var ErrProductCategoryCodeInvalid = errors.New("category code only supports letters, numbers, hyphen and underscore")
var ErrProductCategoryCodeImmutable = errors.New("product category code cannot be changed after being referenced")
var ErrComboMainRequiresActiveChildren = errors.New("combo main product requires at least one active child")
var ErrComboActiveRequiresActiveChildren = errors.New("active combo cannot have all child products inactive")
var ErrProductReferenced = errors.New("product is still referenced by business data")
var ErrProductQuoteRepoUnavailable = errors.New("product quote repository is required")

type ProductQuoteRepository interface {
	GetByProductSupplier(productID, supplierID uint64) (*supplierdomain.ProductSupplierQuote, error)
	Create(quote *supplierdomain.ProductSupplierQuote) error
	Update(quote *supplierdomain.ProductSupplierQuote) error
}

type ProductDefaultsProvider interface {
	GetDefaultBaseCurrency() string
}

type ProductUpsertTransactionalDeps struct {
	ProductRepo domain.ProductRepository
	QuoteRepo   ProductQuoteRepository
	ImageRepo   domain.ProductImageRepository
}

type ProductUpsertTransactionManager interface {
	Run(ctx context.Context, fn func(ProductUpsertTransactionalDeps) error) error
}

type ProductUpsertEffects struct {
	AutoQuoteCreated bool
}

type ProductUsecase struct {
	productRepo          domain.ProductRepository
	productParentRepo    domain.ProductParentRepository
	productConfigRepo    domain.ProductConfigRepository
	productCategoryRepo  domain.ProductCategoryRepository
	productPackagingRepo domain.ProductPackagingRepository
	quoteRepo            ProductQuoteRepository
	imageRepo            domain.ProductImageRepository
	defaultsProvider     ProductDefaultsProvider
	upsertTxManager      ProductUpsertTransactionManager
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
	productParentRepo domain.ProductParentRepository,
	productConfigRepo domain.ProductConfigRepository,
	productCategoryRepo domain.ProductCategoryRepository,
	productPackagingRepo domain.ProductPackagingRepository,
) *ProductUsecase {
	return &ProductUsecase{
		productRepo:          productRepo,
		productParentRepo:    productParentRepo,
		productConfigRepo:    productConfigRepo,
		productCategoryRepo:  productCategoryRepo,
		productPackagingRepo: productPackagingRepo,
	}
}

func (uc *ProductUsecase) BindQuoteRepository(repo ProductQuoteRepository) {
	uc.quoteRepo = repo
}

func (uc *ProductUsecase) BindImageRepository(repo domain.ProductImageRepository) {
	uc.imageRepo = repo
}

func (uc *ProductUsecase) BindDefaultsProvider(provider ProductDefaultsProvider) {
	uc.defaultsProvider = provider
}

func (uc *ProductUsecase) BindUpsertTransactionManager(manager ProductUpsertTransactionManager) {
	uc.upsertTxManager = manager
}

// Product 相关方法

func (uc *ProductUsecase) ListProducts(params *domain.ProductListParams) ([]domain.Product, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	items, total, err := uc.productRepo.List(params)
	if err != nil {
		return nil, 0, err
	}
	if err := uc.applyProductDeleteMeta(items); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (uc *ProductUsecase) GetProduct(id uint64) (*domain.Product, error) {
	product, err := uc.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}
	items := []domain.Product{*product}
	if err := uc.applyProductDeleteMeta(items); err != nil {
		return nil, err
	}
	copyItem := items[0]
	return &copyItem, nil
}

func (uc *ProductUsecase) CreateProduct(product *domain.Product) error {
	_, err := uc.CreateProductWithEffects(product)
	return err
}

func (uc *ProductUsecase) CreateProductWithEffects(product *domain.Product) (*ProductUpsertEffects, error) {
	if product.SupplierID == nil || *product.SupplierID == 0 {
		return nil, ErrSupplierRequired
	}
	if err := uc.validateProduct(product); err != nil {
		return nil, err
	}
	effects := &ProductUpsertEffects{}
	err := uc.runProductUpsert(context.Background(), func(productRepo domain.ProductRepository, quoteRepo ProductQuoteRepository, imageRepo domain.ProductImageRepository) error {
		if err := productRepo.Create(product); err != nil {
			return err
		}
		created, err := uc.ensureDefaultSupplierQuote(product, quoteRepo)
		if err != nil {
			return err
		}
		effects.AutoQuoteCreated = created
		if err := uc.syncPrimaryProductImage(product, productRepo, imageRepo); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return effects, nil
}

func (uc *ProductUsecase) UpdateProduct(product *domain.Product) error {
	_, err := uc.UpdateProductWithEffects(product)
	return err
}

func (uc *ProductUsecase) UpdateProductWithEffects(product *domain.Product) (*ProductUpsertEffects, error) {
	if product.SupplierID == nil || *product.SupplierID == 0 {
		return nil, ErrSupplierRequired
	}
	if err := uc.validateProduct(product); err != nil {
		return nil, err
	}
	if err := uc.validateComboStatusChange(product); err != nil {
		return nil, err
	}
	effects := &ProductUpsertEffects{}
	err := uc.runProductUpsert(context.Background(), func(productRepo domain.ProductRepository, quoteRepo ProductQuoteRepository, imageRepo domain.ProductImageRepository) error {
		if err := productRepo.Update(product); err != nil {
			return err
		}
		created, err := uc.ensureDefaultSupplierQuote(product, quoteRepo)
		if err != nil {
			return err
		}
		effects.AutoQuoteCreated = created
		if err := uc.syncPrimaryProductImage(product, productRepo, imageRepo); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return effects, nil
}

func (uc *ProductUsecase) runProductUpsert(ctx context.Context, fn func(domain.ProductRepository, ProductQuoteRepository, domain.ProductImageRepository) error) error {
	if uc.quoteRepo == nil {
		return ErrProductQuoteRepoUnavailable
	}
	if uc.upsertTxManager == nil {
		return fn(uc.productRepo, uc.quoteRepo, uc.imageRepo)
	}
	return uc.upsertTxManager.Run(ctx, func(deps ProductUpsertTransactionalDeps) error {
		productRepo := deps.ProductRepo
		if productRepo == nil {
			productRepo = uc.productRepo
		}
		quoteRepo := deps.QuoteRepo
		if quoteRepo == nil {
			quoteRepo = uc.quoteRepo
		}
		imageRepo := deps.ImageRepo
		if imageRepo == nil {
			imageRepo = uc.imageRepo
		}
		return fn(productRepo, quoteRepo, imageRepo)
	})
}

func (uc *ProductUsecase) syncPrimaryProductImage(product *domain.Product, productRepo domain.ProductRepository, imageRepo domain.ProductImageRepository) error {
	if product == nil || productRepo == nil || imageRepo == nil || product.ID == 0 {
		return nil
	}
	items, err := imageRepo.ListByProductID(product.ID)
	if err != nil {
		return err
	}
	current := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		url := strings.TrimSpace(item.ImageUrl)
		if url == "" {
			continue
		}
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		current = append(current, url)
	}

	next := uc.buildPrimaryImageOrder(strings.TrimSpace(product.ImageUrl), current)
	if sameStringSlice(current, next) {
		return nil
	}
	if err := imageRepo.ReplaceAll(product.ID, next); err != nil {
		return err
	}
	primary := ""
	if len(next) > 0 {
		primary = next[0]
	}
	return productRepo.UpdateImageUrl(product.ID, primary)
}

func (uc *ProductUsecase) buildPrimaryImageOrder(primary string, current []string) []string {
	primary = strings.TrimSpace(primary)
	if primary == "" {
		if len(current) <= 1 {
			return []string{}
		}
		return append([]string(nil), current[1:]...)
	}

	next := make([]string, 0, len(current)+1)
	next = append(next, primary)
	for _, item := range current {
		item = strings.TrimSpace(item)
		if item == "" || item == primary {
			continue
		}
		next = append(next, item)
	}
	return next
}

func sameStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (uc *ProductUsecase) ensureDefaultSupplierQuote(product *domain.Product, quoteRepo ProductQuoteRepository) (bool, error) {
	if product == nil || product.SupplierID == nil || *product.SupplierID == 0 {
		return false, nil
	}
	existing, err := quoteRepo.GetByProductSupplier(product.ID, *product.SupplierID)
	if err == nil && existing != nil {
		if product.UnitCost == nil {
			return false, nil
		}
		desiredPrice := *product.UnitCost
		if existing.Price != desiredPrice {
			existing.Price = desiredPrice
			if err := quoteRepo.Update(existing); err != nil {
				return false, err
			}
		}
		return false, nil
	}
	if err != nil && !errors.Is(err, supplierdomain.ErrQuoteNotFound) && !strings.Contains(strings.ToLower(err.Error()), "record not found") {
		return false, err
	}
	currency := "USD"
	if uc.defaultsProvider != nil && strings.TrimSpace(uc.defaultsProvider.GetDefaultBaseCurrency()) != "" {
		currency = strings.ToUpper(strings.TrimSpace(uc.defaultsProvider.GetDefaultBaseCurrency()))
	}
	quote := &supplierdomain.ProductSupplierQuote{
		ProductID:    product.ID,
		SupplierID:   *product.SupplierID,
		Price:        normalizeProductUnitCost(product.UnitCost),
		Currency:     currency,
		QtyMOQ:       1,
		LeadTimeDays: 0,
		Status:       supplierdomain.ProductSupplierQuoteStatusPending,
		Remark:       "产品默认供应商自动生成，待报价",
	}
	if err := quoteRepo.Create(quote); err != nil {
		return false, err
	}
	return true, nil
}

func normalizeProductUnitCost(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func (uc *ProductUsecase) DeleteProduct(id uint64) error {
	if err := uc.ensureProductDeletable(id); err != nil {
		return err
	}
	return uc.productRepo.Delete(id)
}

func (uc *ProductUsecase) ensureProductDeletable(id uint64) error {
	if uc.productRepo == nil {
		return nil
	}
	counts, err := uc.productRepo.CountReferencesByIDs([]uint64{id})
	if err != nil {
		return err
	}
	if counts[id] > 0 {
		return ErrProductReferenced
	}
	return nil
}

func (uc *ProductUsecase) applyProductDeleteMeta(items []domain.Product) error {
	if uc.productRepo == nil || len(items) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	counts, err := uc.productRepo.CountReferencesByIDs(ids)
	if err != nil {
		return err
	}
	for i := range items {
		items[i].Deletable = true
		items[i].ReferenceCount = counts[items[i].ID]
		if items[i].ReferenceCount > 0 {
			items[i].Deletable = false
			items[i].DeleteBlockReason = "已被业务数据或组合关系引用，不可删除"
		}
	}
	return nil
}

// ProductConfig 相关方法

func (uc *ProductUsecase) ListProductConfigs(params *domain.ProductConfigListParams) ([]domain.ProductConfigItem, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 100
	}
	items, total, err := uc.productConfigRepo.List(params)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		items[i].Deletable = true
		if items[i].ConfigType == domain.ProductConfigTypeSalesStatus {
			items[i].Deletable = false
			items[i].DeleteBlockReason = "系统销售状态，不可删除"
			continue
		}
		if uc.productRepo == nil {
			continue
		}
		refCount, err := uc.productRepo.CountByConfigReference(items[i].ConfigType, items[i].ID)
		if err != nil {
			return nil, 0, err
		}
		items[i].ReferenceCount = refCount
		if refCount > 0 {
			items[i].Deletable = false
			items[i].DeleteBlockReason = "已被产品引用，不可删除"
		}
	}
	return items, total, nil
}

func (uc *ProductUsecase) GetProductConfig(id uint64) (*domain.ProductConfigItem, error) {
	return uc.productConfigRepo.GetByID(id)
}

func (uc *ProductUsecase) CreateProductConfig(item *domain.ProductConfigItem) error {
	if item.ConfigType == "" || item.ConfigType == "CATEGORY" {
		return ErrProductConfigTypeUnsupported
	}
	if item.ConfigType == domain.ProductConfigTypeSalesStatus {
		return ErrProductConfigSystemFixed
	}
	if !validation.IsValidCode(strings.TrimSpace(item.ItemCode)) {
		return ErrProductConfigCodeInvalid
	}
	return uc.productConfigRepo.Create(item)
}

func (uc *ProductUsecase) UpdateProductConfig(item *domain.ProductConfigItem) error {
	if item.ConfigType == "" || item.ConfigType == "CATEGORY" {
		return ErrProductConfigTypeUnsupported
	}
	if !validation.IsValidCode(strings.TrimSpace(item.ItemCode)) {
		return ErrProductConfigCodeInvalid
	}
	current, err := uc.productConfigRepo.GetByID(item.ID)
	if err != nil {
		return err
	}
	if current.ConfigType == domain.ProductConfigTypeSalesStatus && current.ItemCode != item.ItemCode {
		return ErrProductConfigSystemFixed
	}
	if uc.productRepo != nil && current.ItemCode != item.ItemCode {
		refCount, err := uc.productRepo.CountByConfigReference(current.ConfigType, current.ID)
		if err != nil {
			return err
		}
		if refCount > 0 {
			return ErrProductConfigCodeImmutable
		}
	}
	item.ConfigType = current.ConfigType
	return uc.productConfigRepo.Update(item)
}

func (uc *ProductUsecase) DeleteProductConfig(id uint64) error {
	item, err := uc.productConfigRepo.GetByID(id)
	if err != nil {
		return err
	}
	if item.ConfigType == domain.ProductConfigTypeSalesStatus {
		return ErrProductConfigSystemFixed
	}
	if uc.productRepo != nil {
		refCount, err := uc.productRepo.CountByConfigReference(item.ConfigType, id)
		if err != nil {
			return err
		}
		if refCount > 0 {
			return ErrProductConfigReferenced
		}
	}
	return uc.productConfigRepo.Delete(id)
}

func (uc *ProductUsecase) ListProductCategories() ([]domain.ProductCategory, error) {
	items, err := uc.productCategoryRepo.ListAll()
	if err != nil {
		return nil, err
	}

	nodeMap := make(map[uint64]*domain.ProductCategory, len(items))
	roots := make([]domain.ProductCategory, 0)
	for _, item := range items {
		copyItem := item
		copyItem.Children = nil
		if err := uc.applyCategoryDeleteMeta(&copyItem, items); err != nil {
			return nil, err
		}
		nodeMap[item.ID] = &copyItem
	}
	for _, item := range items {
		node := nodeMap[item.ID]
		if item.ParentID == nil {
			roots = append(roots, *node)
			continue
		}
		parent, ok := nodeMap[*item.ParentID]
		if !ok {
			roots = append(roots, *node)
			continue
		}
		parent.Children = append(parent.Children, node)
	}

	result := make([]domain.ProductCategory, 0, len(roots))
	for i := range roots {
		if node, ok := nodeMap[roots[i].ID]; ok {
			result = append(result, *node)
		}
	}
	return result, nil
}

func (uc *ProductUsecase) applyCategoryDeleteMeta(item *domain.ProductCategory, all []domain.ProductCategory) error {
	item.Deletable = true
	switch {
	case uc.categoryHasChildren(item.ID, all):
		item.Deletable = false
		item.DeleteBlockReason = "仍有下级品类，不可删除"
		return nil
	}
	if uc.productRepo == nil {
		return nil
	}
	refCount, err := uc.productRepo.CountByCategoryID(item.ID)
	if err != nil {
		return err
	}
	item.ReferenceCount = refCount
	if refCount > 0 {
		item.Deletable = false
		item.DeleteBlockReason = "已被产品引用，不可删除"
	}
	return nil
}

func (uc *ProductUsecase) categoryHasChildren(id uint64, all []domain.ProductCategory) bool {
	for _, item := range all {
		if item.ParentID != nil && *item.ParentID == id {
			return true
		}
	}
	return false
}

func (uc *ProductUsecase) CreateProductCategory(item *domain.ProductCategory) error {
	if !validation.IsValidCode(strings.TrimSpace(item.CategoryCode)) {
		return ErrProductCategoryCodeInvalid
	}
	if item.ParentID == nil || *item.ParentID == 0 {
		item.ParentID = nil
		item.Level = 1
		return uc.productCategoryRepo.Create(item)
	}
	parent, err := uc.productCategoryRepo.GetByID(*item.ParentID)
	if err != nil {
		return err
	}
	if parent.Level >= 3 {
		return ErrProductCategoryParentInvalid
	}
	item.Level = parent.Level + 1
	return uc.productCategoryRepo.Create(item)
}

func (uc *ProductUsecase) UpdateProductCategory(item *domain.ProductCategory) error {
	current, err := uc.productCategoryRepo.GetByID(item.ID)
	if err != nil {
		return err
	}
	if !validation.IsValidCode(strings.TrimSpace(item.CategoryCode)) {
		return ErrProductCategoryCodeInvalid
	}
	if !sameUint64Ptr(current.ParentID, item.ParentID) {
		return ErrProductCategoryParentInvalid
	}
	if uc.productRepo != nil && current.CategoryCode != item.CategoryCode {
		refCount, err := uc.productRepo.CountByCategoryID(item.ID)
		if err != nil {
			return err
		}
		if refCount > 0 {
			return ErrProductCategoryCodeImmutable
		}
	}
	item.Level = current.Level
	return uc.productCategoryRepo.Update(item)
}

func (uc *ProductUsecase) DeleteProductCategory(id uint64) error {
	childCount, err := uc.productCategoryRepo.CountChildren(id)
	if err != nil {
		return err
	}
	if childCount > 0 {
		return ErrProductCategoryHasChildren
	}
	if uc.productRepo != nil {
		refCount, err := uc.productRepo.CountByCategoryID(id)
		if err != nil {
			return err
		}
		if refCount > 0 {
			return ErrProductCategoryReferenced
		}
	}
	return uc.productCategoryRepo.Delete(id)
}

func (uc *ProductUsecase) isAllowedSalesStatusCode(code string) bool {
	switch strings.ToUpper(strings.TrimSpace(code)) {
	case domain.ProductStatusDraft, domain.ProductStatusOnSale, domain.ProductStatusReplenishing, domain.ProductStatusOffShelf:
		return true
	default:
		return false
	}
}

func sameUint64Ptr(a, b *uint64) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	default:
		return *a == *b
	}
}

// ProductGroup 相关方法

func (uc *ProductUsecase) ListProductParents(params *domain.ProductParentListParams) ([]domain.ProductParent, int64, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.productParentRepo.List(params)
}

func (uc *ProductUsecase) GetProductParent(id uint64) (*domain.ProductParent, error) {
	parent, err := uc.productParentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	children, err := uc.productRepo.ListByParentID(id)
	if err != nil {
		return nil, err
	}
	parent.Children = children
	if parent.ChildCount == 0 {
		parent.ChildCount = int64(len(children))
		for _, child := range children {
			if child.Status == "ON_SALE" || child.Status == "REPLENISHING" {
				parent.ActiveChildCount++
			} else {
				parent.InactiveChildCount++
			}
		}
	}
	return parent, nil
}

func (uc *ProductUsecase) CreateProductParent(parent *domain.ProductParent) error {
	return uc.productParentRepo.Create(parent)
}

func (uc *ProductUsecase) UpdateProductParent(parent *domain.ProductParent) error {
	current, err := uc.productParentRepo.GetByID(parent.ID)
	if err != nil {
		return err
	}
	if current.ParentAsin != parent.ParentAsin || current.Marketplace != parent.Marketplace {
		return ErrProductParentImmutableIdentity
	}
	return uc.productParentRepo.Update(parent)
}

func (uc *ProductUsecase) DeleteProductParent(id uint64) error {
	children, err := uc.productRepo.ListByParentID(id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return ErrProductParentHasChildren
	}
	return uc.productParentRepo.Delete(id)
}

func (uc *ProductUsecase) AttachProductParentChildren(parentID uint64, childIDs []uint64) (*domain.ProductParent, error) {
	if len(childIDs) == 0 {
		return nil, ErrProductParentChildIDsRequired
	}

	parent, err := uc.productParentRepo.GetByID(parentID)
	if err != nil {
		return nil, err
	}

	uniqueIDs := make([]uint64, 0, len(childIDs))
	seen := make(map[uint64]struct{}, len(childIDs))
	for _, id := range childIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}
	if len(uniqueIDs) == 0 {
		return nil, ErrProductParentChildIDsRequired
	}

	children, err := uc.productRepo.ListByIDs(uniqueIDs)
	if err != nil {
		return nil, err
	}
	if len(children) != len(uniqueIDs) {
		return nil, ErrProductParentChildNotFound
	}

	assignIDs := make([]uint64, 0, len(uniqueIDs))
	for _, child := range children {
		if child.ParentID != nil && *child.ParentID != parentID {
			return nil, ErrProductParentChildAlreadyAssigned
		}
		if child.Marketplace != parent.Marketplace {
			return nil, ErrProductParentChildMarketplaceMismatch
		}
		assignIDs = append(assignIDs, child.ID)
	}

	if err := uc.productRepo.UpdateParentIDBatch(assignIDs, &parentID); err != nil {
		return nil, err
	}

	return uc.GetProductParent(parentID)
}

func (uc *ProductUsecase) DetachProductParentChild(parentID uint64, childID uint64) (*domain.ProductParent, error) {
	child, err := uc.productRepo.GetByID(childID)
	if err != nil {
		return nil, err
	}
	if child.ParentID == nil || *child.ParentID != parentID {
		return nil, ErrProductParentChildNotFound
	}
	if err := uc.productRepo.UpdateParentIDBatch([]uint64{childID}, nil); err != nil {
		return nil, err
	}
	return uc.GetProductParent(parentID)
}

// ProductPackaging 相关方法

func (uc *ProductUsecase) GetProductPackagingItems(productID uint64) ([]domain.ProductPackagingItem, error) {
	return uc.productPackagingRepo.ListByProductID(productID)
}

func (uc *ProductUsecase) SaveProductPackagingItems(productID uint64, items []domain.ProductPackagingItem) error {
	return uc.productPackagingRepo.ReplaceAll(productID, items)
}

func (uc *ProductUsecase) validateComboStatusChange(product *domain.Product) error {
	current, err := uc.productRepo.GetByID(product.ID)
	if err != nil {
		return err
	}

	comboID := current.ComboID
	if product.ComboID != nil {
		comboID = product.ComboID
	}
	if comboID == nil || *comboID == 0 {
		return nil
	}

	isComboMain := current.IsComboMain
	if product.IsComboMain != 0 {
		isComboMain = product.IsComboMain
	}

	nextStatus := current.Status
	if strings.TrimSpace(product.Status) != "" {
		nextStatus = product.Status
	}

	relatedProducts, _, err := uc.productRepo.List(&domain.ProductListParams{
		Page:     1,
		PageSize: 1000,
		ComboID:  comboID,
	})
	if err != nil {
		return err
	}
	if len(relatedProducts) == 0 {
		return nil
	}

	if isComboMain == 1 && isActiveProductStatus(nextStatus) {
		hasActiveChild := false
		for _, item := range relatedProducts {
			if item.ID == product.ID || item.IsComboMain == 1 {
				continue
			}
			if isActiveProductStatus(item.Status) {
				hasActiveChild = true
				break
			}
		}
		if !hasActiveChild {
			return ErrComboMainRequiresActiveChildren
		}
		return nil
	}

	if isComboMain == 0 && !isActiveProductStatus(nextStatus) {
		activeComboMain := false
		activeChildCount := 0
		for _, item := range relatedProducts {
			status := item.Status
			if item.ID == product.ID {
				status = nextStatus
			}
			if item.IsComboMain == 1 {
				activeComboMain = isActiveProductStatus(status)
				continue
			}
			if isActiveProductStatus(status) {
				activeChildCount++
			}
		}
		if activeComboMain && activeChildCount == 0 {
			return ErrComboActiveRequiresActiveChildren
		}
	}

	return nil
}

func (uc *ProductUsecase) validateProduct(product *domain.Product) error {
	asin := strings.TrimSpace(product.Asin)
	if asin == "" || len(asin) > 20 {
		return ErrProductAsinInvalid
	}
	product.Asin = asin
	return nil
}

func isActiveProductStatus(status string) bool {
	return status == domain.ProductStatusOnSale || status == domain.ProductStatusReplenishing
}
