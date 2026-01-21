package usecase

import (
	"errors"
	"sort"

	"am-erp-go/internal/module/product/domain"
)

type ComboProductUpdater interface {
	UpdateComboInfo(comboID uint64, mainProductID uint64, productIDs []uint64) error
	ClearComboInfo(comboID uint64) error
}

type ProductComboUsecase struct {
	comboRepo   domain.ProductComboRepository
	productRepo ComboProductUpdater
	productRead ProductReader
}

type ProductReader interface {
	ListByIDs(ids []uint64) ([]domain.Product, error)
}

func NewProductComboUsecase(comboRepo domain.ProductComboRepository, productRepo ComboProductUpdater, productRead ProductReader) *ProductComboUsecase {
	return &ProductComboUsecase{
		comboRepo:   comboRepo,
		productRepo: productRepo,
		productRead: productRead,
	}
}

func (uc *ProductComboUsecase) CreateCombo(params domain.ComboUpsertParams) (*domain.ProductCombo, error) {
	if params.MainProductID == 0 {
		return nil, errors.New("missing main product")
	}

	childIDs := uniqueChildIDs(params.MainProductID, params.ProductIDs)
	ratios := normalizeRatios(childIDs, params.QtyRatios)

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

	return &domain.ProductCombo{ComboID: comboID}, nil
}

func (uc *ProductComboUsecase) UpdateComboByMainProductID(mainProductID uint64, params domain.ComboUpsertParams) (*domain.ProductCombo, error) {
	if mainProductID == 0 {
		return nil, errors.New("missing main product")
	}

	comboID, err := uc.comboRepo.GetComboIDByMainProductID(mainProductID)
	if err != nil {
		return nil, err
	}

	childIDs := uniqueChildIDs(mainProductID, params.ProductIDs)
	ratios := normalizeRatios(childIDs, params.QtyRatios)

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

	return &domain.ProductCombo{ComboID: comboID}, nil
}

func (uc *ProductComboUsecase) DeleteCombo(comboID uint64) error {
	if err := uc.comboRepo.DeleteCombo(comboID); err != nil {
		return err
	}
	if uc.productRepo != nil {
		return uc.productRepo.ClearComboInfo(comboID)
	}
	return nil
}

func (uc *ProductComboUsecase) ListCombos(params *domain.ComboListParams) ([]domain.ProductCombo, int64, error) {
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

	children := make([]domain.Product, 0, len(items))
	for _, item := range items {
		if item.ProductID == mainProductID {
			continue
		}
		if product, ok := productMap[item.ProductID]; ok {
			children = append(children, product)
		}
	}
	combo.Products = children
	return combo, nil
}

func uniqueChildIDs(mainID uint64, productIDs []uint64) []uint64 {
	if len(productIDs) == 0 {
		return []uint64{}
	}
	seen := make(map[uint64]struct{}, len(productIDs))
	result := make([]uint64, 0, len(productIDs))
	for _, id := range productIDs {
		if id == 0 || id == mainID {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func normalizeRatios(productIDs []uint64, ratios map[uint64]uint64) map[uint64]uint64 {
	result := make(map[uint64]uint64, len(productIDs))
	for _, id := range productIDs {
		ratio := uint64(1)
		if ratios != nil {
			if v, ok := ratios[id]; ok && v > 0 {
				ratio = v
			}
		}
		result[id] = ratio
	}
	return result
}
