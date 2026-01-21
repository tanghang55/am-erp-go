package repository

import (
	"am-erp-go/internal/module/product/domain"

	"gorm.io/gorm"
)

type comboRepository struct {
	db *gorm.DB
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
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	var total int64
	if err := r.db.Model(&domain.ProductComboItem{}).
		Distinct("combo_id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize

	var rows []struct {
		ComboID uint64 `gorm:"column:combo_id"`
	}
	if err := r.db.Model(&domain.ProductComboItem{}).
		Select("combo_id").
		Group("combo_id").
		Order("combo_id DESC").
		Offset(offset).
		Limit(params.PageSize).
		Scan(&rows).Error; err != nil {
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

func (r *comboRepository) GetComboIDByMainProductID(mainProductID uint64) (uint64, error) {
	var row struct {
		ComboID uint64 `gorm:"column:combo_id"`
	}
	if err := r.db.Model(&domain.ProductComboItem{}).
		Select("combo_id").
		Where("main_product_id = ?", mainProductID).
		Order("combo_id DESC").
		Limit(1).
		Scan(&row).Error; err != nil {
		return 0, err
	}
	return row.ComboID, nil
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
