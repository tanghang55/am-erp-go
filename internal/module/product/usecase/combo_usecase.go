package usecase

import (
	"errors"
	"sort"
	"time"

	"am-erp-go/internal/module/product/domain"
)

var ErrComboMainProductRequired = errors.New("missing main product")
var ErrComboNotFound = errors.New("combo not found")
var ErrComboProductNotFound = errors.New("product not found")
var ErrComboStandaloneOnly = errors.New("only standalone products can be selected")
var ErrComboLocked = errors.New("combo already has business data and cannot be modified or deleted")

type ComboProductUpdater interface {
	UpdateComboInfo(comboID uint64, mainProductID uint64, productIDs []uint64) error
	ClearComboInfo(comboID uint64) error
}

type ComboUsageChecker interface {
	HasBusinessUsageSince(productIDs []uint64, since time.Time) (bool, error)
}

type ProductComboUsecase struct {
	comboRepo   domain.ProductComboRepository
	productRepo ComboProductUpdater
	productRead ProductReader
	usageCheck  ComboUsageChecker
}

type ProductReader interface {
	ListByIDs(ids []uint64) ([]domain.Product, error)
}

func NewProductComboUsecase(comboRepo domain.ProductComboRepository, productRepo ComboProductUpdater, productRead ProductReader, usageCheck ComboUsageChecker) *ProductComboUsecase {
	return &ProductComboUsecase{
		comboRepo:   comboRepo,
		productRepo: productRepo,
		productRead: productRead,
		usageCheck:  usageCheck,
	}
}

func (uc *ProductComboUsecase) CreateCombo(params domain.ComboUpsertParams) (*domain.ProductCombo, error) {
	if params.MainProductID == 0 {
		return nil, ErrComboMainProductRequired
	}

	childIDs, ratios := normalizeChildren(params.MainProductID, params.Children)
	if err := uc.validateStandaloneProducts(0, append([]uint64{params.MainProductID}, childIDs...)); err != nil {
		return nil, err
	}

	comboID, err := uc.comboRepo.CreateCombo(params.MainProductID, childIDs, ratios)
	if err != nil {
		return nil, err
	}

	if uc.productRepo != nil {
		allIDs := append([]uint64{params.MainProductID}, childIDs...)
		if err := uc.productRepo.UpdateComboInfo(comboID, params.MainProductID, allIDs); err != nil {
			return nil, err
		}
	}

	return uc.GetCombo(comboID)
}

func (uc *ProductComboUsecase) UpdateCombo(comboID uint64, params domain.ComboUpsertParams) (*domain.ProductCombo, error) {
	if comboID == 0 {
		return nil, errors.New("missing combo id")
	}

	items, err := uc.comboRepo.GetItemsByComboID(comboID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrComboNotFound
	}

	mainProductID := params.MainProductID
	if mainProductID == 0 {
		mainProductID = items[0].MainProductID
		if mainProductID == 0 {
			for _, item := range items {
				if item.MainProductID != 0 {
					mainProductID = item.MainProductID
					break
				}
			}
		}
	}
	if mainProductID == 0 {
		return nil, ErrComboMainProductRequired
	}
	if err := uc.ensureComboMutable(comboID, items); err != nil {
		return nil, err
	}

	childIDs, ratios := normalizeChildren(mainProductID, params.Children)
	if err := uc.validateStandaloneProducts(comboID, append([]uint64{mainProductID}, childIDs...)); err != nil {
		return nil, err
	}

	if err := uc.comboRepo.ReplaceComboItems(comboID, mainProductID, childIDs, ratios); err != nil {
		return nil, err
	}

	if uc.productRepo != nil {
		if err := uc.productRepo.ClearComboInfo(comboID); err != nil {
			return nil, err
		}
		allIDs := append([]uint64{mainProductID}, childIDs...)
		if err := uc.productRepo.UpdateComboInfo(comboID, mainProductID, allIDs); err != nil {
			return nil, err
		}
	}

	return uc.GetCombo(comboID)
}

func (uc *ProductComboUsecase) DeleteCombo(comboID uint64) error {
	items, err := uc.comboRepo.GetItemsByComboID(comboID)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return ErrComboNotFound
	}
	if err := uc.ensureComboMutable(comboID, items); err != nil {
		return err
	}
	if err := uc.comboRepo.DeleteCombo(comboID); err != nil {
		return err
	}
	if uc.productRepo != nil {
		return uc.productRepo.ClearComboInfo(comboID)
	}
	return nil
}

func (uc *ProductComboUsecase) ListCombos(params *domain.ComboListParams) ([]domain.ProductCombo, int64, error) {
	if params != nil && params.Locked != "" {
		return uc.listCombosWithLockFilter(params)
	}
	comboIDs, total, err := uc.comboRepo.ListComboIDs(params)
	if err != nil {
		return nil, 0, err
	}
	if len(comboIDs) == 0 {
		return []domain.ProductCombo{}, total, nil
	}

	combos := make([]domain.ProductCombo, 0, len(comboIDs))
	for _, comboID := range comboIDs {
		combo, err := uc.GetCombo(comboID)
		if err != nil {
			return nil, 0, err
		}
		combos = append(combos, *combo)
	}
	return combos, total, nil
}

func (uc *ProductComboUsecase) listCombosWithLockFilter(params *domain.ComboListParams) ([]domain.ProductCombo, int64, error) {
	allParams := *params
	allParams.Page = 1
	allParams.PageSize = -1

	comboIDs, _, err := uc.comboRepo.ListComboIDs(&allParams)
	if err != nil {
		return nil, 0, err
	}

	wantLocked := params.Locked == "true"
	filtered := make([]domain.ProductCombo, 0, len(comboIDs))
	for _, comboID := range comboIDs {
		combo, err := uc.GetCombo(comboID)
		if err != nil {
			return nil, 0, err
		}
		if combo.Locked == wantLocked {
			filtered = append(filtered, *combo)
		}
	}

	total := int64(len(filtered))
	if total == 0 {
		return []domain.ProductCombo{}, 0, nil
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	start := (page - 1) * pageSize
	if start >= len(filtered) {
		return []domain.ProductCombo{}, total, nil
	}
	end := start + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

func (uc *ProductComboUsecase) GetCombo(comboID uint64) (*domain.ProductCombo, error) {
	items, err := uc.comboRepo.GetItemsByComboID(comboID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return &domain.ProductCombo{ComboID: comboID}, nil
	}

	mainProductID := items[0].MainProductID
	if mainProductID == 0 {
		for _, item := range items {
			if item.MainProductID != 0 {
				mainProductID = item.MainProductID
				break
			}
		}
	}

	if uc.productRead == nil {
		return &domain.ProductCombo{ComboID: comboID}, nil
	}

	productIDs := make([]uint64, 0, len(items))
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := uc.productRead.ListByIDs(productIDs)
	if err != nil {
		return nil, err
	}
	productMap := make(map[uint64]domain.Product, len(products))
	for _, product := range products {
		productMap[product.ID] = product
	}

	combo := &domain.ProductCombo{ComboID: comboID}
	if main, ok := productMap[mainProductID]; ok {
		combo.MainProduct = main
	}

	children := make([]domain.ComboChildProduct, 0, len(items))
	for _, item := range items {
		if item.ProductID == mainProductID {
			continue
		}
		if product, ok := productMap[item.ProductID]; ok {
			children = append(children, domain.ComboChildProduct{
				Product:  product,
				QtyRatio: item.QtyRatio,
			})
		}
	}
	combo.Products = children
	locked, lockReason, err := uc.getComboLockState(items)
	if err != nil {
		return nil, err
	}
	combo.Locked = locked
	combo.LockReason = lockReason
	return combo, nil
}

func (uc *ProductComboUsecase) validateStandaloneProducts(currentComboID uint64, productIDs []uint64) error {
	if uc.productRead == nil || len(productIDs) == 0 {
		return nil
	}
	products, err := uc.productRead.ListByIDs(productIDs)
	if err != nil {
		return err
	}
	productMap := make(map[uint64]domain.Product, len(products))
	for _, product := range products {
		productMap[product.ID] = product
	}
	for _, productID := range productIDs {
		product, ok := productMap[productID]
		if !ok {
			return ErrComboProductNotFound
		}
		if product.ComboID != nil && *product.ComboID != 0 && *product.ComboID != currentComboID {
			return ErrComboStandaloneOnly
		}
	}
	return nil
}

func (uc *ProductComboUsecase) ensureComboMutable(comboID uint64, items []domain.ProductComboItem) error {
	locked, _, err := uc.getComboLockState(items)
	if err != nil {
		return err
	}
	if locked {
		return ErrComboLocked
	}
	return nil
}

func (uc *ProductComboUsecase) getComboLockState(items []domain.ProductComboItem) (bool, string, error) {
	if uc.usageCheck == nil || len(items) == 0 {
		return false, "", nil
	}
	productIDs := make([]uint64, 0, len(items))
	since := items[0].GmtCreate
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
		if item.GmtCreate.Before(since) {
			since = item.GmtCreate
		}
	}
	locked, err := uc.usageCheck.HasBusinessUsageSince(productIDs, since)
	if err != nil {
		return false, "", err
	}
	if !locked {
		return false, "", nil
	}
	return true, "已有业务数据，禁止修改或删除", nil
}

func normalizeChildren(mainID uint64, children []domain.ComboChildInput) ([]uint64, map[uint64]uint64) {
	if len(children) == 0 {
		return []uint64{}, map[uint64]uint64{}
	}

	type normalizedChild struct {
		ProductID uint64
		QtyRatio  uint64
	}

	normalized := make(map[uint64]normalizedChild, len(children))
	for _, child := range children {
		if child.ProductID == 0 || child.ProductID == mainID {
			continue
		}
		ratio := child.QtyRatio
		if ratio == 0 {
			ratio = 1
		}
		normalized[child.ProductID] = normalizedChild{
			ProductID: child.ProductID,
			QtyRatio:  ratio,
		}
	}

	result := make([]uint64, 0, len(normalized))
	ratios := make(map[uint64]uint64, len(normalized))
	for productID, child := range normalized {
		result = append(result, productID)
		ratios[productID] = child.QtyRatio
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result, ratios
}
