package repository

import (
	menudomain "am-erp-go/internal/module/menu/domain"

	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) menudomain.MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) GetMenusByPermissionCodes(permissionCodes []string) ([]menudomain.Menu, error) {
	var menus []menudomain.Menu
	if len(permissionCodes) == 0 {
		return menus, nil
	}
	err := r.db.Where("permission_code IN ? AND status = 'ACTIVE' AND is_hidden = 0", permissionCodes).
		Order("sort ASC").
		Find(&menus).Error
	return menus, err
}

func (r *menuRepository) GetAllMenus() ([]menudomain.Menu, error) {
	var menus []menudomain.Menu
	err := r.db.Where("status = 'ACTIVE' AND is_hidden = 0").
		Order("sort ASC").
		Find(&menus).Error
	return menus, err
}

func (r *menuRepository) GetAllMenusRaw() ([]menudomain.Menu, error) {
	var menus []menudomain.Menu
	err := r.db.Order("sort ASC").
		Find(&menus).Error
	return menus, err
}

func (r *menuRepository) List(params *menudomain.MenuListParams) ([]menudomain.Menu, int64, error) {
	var menus []menudomain.Menu
	var total int64

	page := params.Page
	pageSize := params.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := r.db.Model(&menudomain.Menu{})
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(
			"title LIKE ? OR title_en LIKE ? OR code LIKE ? OR path LIKE ? OR permission_code LIKE ?",
			keyword, keyword, keyword, keyword, keyword,
		)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.IsHidden != nil {
		query = query.Where("is_hidden = ?", *params.IsHidden)
	}
	if params.ParentID != nil {
		query = query.Where("parent_id = ?", *params.ParentID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("sort ASC, gmt_modified DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

func (r *menuRepository) GetByID(id uint64) (*menudomain.Menu, error) {
	var menu menudomain.Menu
	if err := r.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *menuRepository) Create(menu *menudomain.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) Update(menu *menudomain.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id uint64) error {
	return r.db.Delete(&menudomain.Menu{}, id).Error
}

func (r *menuRepository) UpdateStatus(id uint64, status string) error {
	return r.db.Model(&menudomain.Menu{}).Where("id = ?", id).Update("status", status).Error
}
