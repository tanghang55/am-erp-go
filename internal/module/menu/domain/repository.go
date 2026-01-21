package domain

type MenuRepository interface {
	GetMenusByPermissionCodes(permissionCodes []string) ([]Menu, error)
	GetAllMenus() ([]Menu, error)
	GetAllMenusRaw() ([]Menu, error)
	List(params *MenuListParams) ([]Menu, int64, error)
	GetByID(id uint64) (*Menu, error)
	Create(menu *Menu) error
	Update(menu *Menu) error
	Delete(id uint64) error
	UpdateStatus(id uint64, status string) error
}
