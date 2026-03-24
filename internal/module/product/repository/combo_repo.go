package repository

import (
	"am-erp-go/internal/module/product/domain"
	"strings"

	"gorm.io/gorm"
)

type comboRepository struct {
	db *gorm.DB
}

func qualifiedComboIDColumn() string {
	return "product_combo.combo_id"
}

func NewProductComboRepository(db *gorm.DB) domain.ProductComboRepository {
	return &comboRepository{db: db}
}

func (r *comboRepository) ListComboIDs(params *domain.ComboListParams) ([]uint64, int64, error) {
	if params == nil {
		params = &domain.ComboListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 20
	}

	query := r.db.Model(&domain.ProductComboItem{}).
		Joins("JOIN product AS main_product ON main_product.id = product_combo.main_product_id")

	if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"main_product.seller_sku LIKE ? OR main_product.title LIKE ? OR CAST(product_combo.combo_id AS CHAR) LIKE ?",
			like,
			like,
			like,
		)
	}

	if marketplace := strings.TrimSpace(params.Marketplace); marketplace != "" {
		query = query.Where("main_product.marketplace = ?", marketplace)
	}
	if len(params.Statuses) > 0 {
		query = query.Where("main_product.status IN ?", params.Statuses)
	}

	var total int64
	if err := query.
		Distinct(qualifiedComboIDColumn()).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []struct {
		ComboID uint64 `gorm:"column:combo_id"`
	}
	listQuery := query.
		Select(qualifiedComboIDColumn()).
		Group(qualifiedComboIDColumn()).
		Order(qualifiedComboIDColumn() + " DESC")
	if params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		listQuery = listQuery.Offset(offset).Limit(params.PageSize)
	}
	if err := listQuery.Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	comboIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		comboIDs = append(comboIDs, row.ComboID)
	}

	return comboIDs, total, nil
}

func (r *comboRepository) GetItemsByComboID(comboID uint64) ([]domain.ProductComboItem, error) {
	var items []domain.ProductComboItem
	if err := r.db.Where("combo_id = ?", comboID).
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *comboRepository) CreateCombo(mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) (uint64, error) {
	var comboID uint64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		mainItem := domain.ProductComboItem{
			ComboID:       0,
			MainProductID: mainProductID,
			ProductID:     mainProductID,
			QtyRatio:      1,
		}
		if err := tx.Create(&mainItem).Error; err != nil {
			return err
		}

		comboID = mainItem.ID
		if err := tx.Model(&domain.ProductComboItem{}).
			Where("id = ?", mainItem.ID).
			Update("combo_id", comboID).Error; err != nil {
			return err
		}

		if len(productIDs) == 0 {
			return nil
		}

		items := make([]domain.ProductComboItem, 0, len(productIDs))
		for _, id := range productIDs {
			ratio := uint64(1)
			if qtyRatios != nil {
				if v, ok := qtyRatios[id]; ok && v > 0 {
					ratio = v
				}
			}
			items = append(items, domain.ProductComboItem{
				ComboID:       comboID,
				MainProductID: mainProductID,
				ProductID:     id,
				QtyRatio:      ratio,
			})
		}

		return tx.Create(&items).Error
	})

	if err != nil {
		return 0, err
	}
	return comboID, nil
}

func (r *comboRepository) ReplaceComboItems(comboID uint64, mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("combo_id = ?", comboID).
			Delete(&domain.ProductComboItem{}).Error; err != nil {
			return err
		}

		mainItem := domain.ProductComboItem{
			ComboID:       comboID,
			MainProductID: mainProductID,
			ProductID:     mainProductID,
			QtyRatio:      1,
		}
		if err := tx.Create(&mainItem).Error; err != nil {
			return err
		}

		if len(productIDs) == 0 {
			return nil
		}

		items := make([]domain.ProductComboItem, 0, len(productIDs))
		for _, id := range productIDs {
			ratio := uint64(1)
			if qtyRatios != nil {
				if v, ok := qtyRatios[id]; ok && v > 0 {
					ratio = v
				}
			}
			items = append(items, domain.ProductComboItem{
				ComboID:       comboID,
				MainProductID: mainProductID,
				ProductID:     id,
				QtyRatio:      ratio,
			})
		}

		return tx.Create(&items).Error
	})
}

func (r *comboRepository) DeleteCombo(comboID uint64) error {
	return r.db.Where("combo_id = ?", comboID).
		Delete(&domain.ProductComboItem{}).Error
}
