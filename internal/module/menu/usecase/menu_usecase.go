package usecase

import (
	"strings"

	identitydomain "am-erp-go/internal/module/identity/domain"
	menudomain "am-erp-go/internal/module/menu/domain"
)

type MenuUsecase struct {
	menuRepo menudomain.MenuRepository
	userRepo identitydomain.UserRepository
}

func NewMenuUsecase(menuRepo menudomain.MenuRepository, userRepo identitydomain.UserRepository) *MenuUsecase {
	return &MenuUsecase{
		menuRepo: menuRepo,
		userRepo: userRepo,
	}
}

func (uc *MenuUsecase) GetMenuTree(userID uint64) ([]*menudomain.MenuTree, error) {
	roles, err := uc.userRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	isAdmin := false
	for _, role := range roles {
		if role.IsAdmin() {
			isAdmin = true
			break
		}
	}

	var menus []menudomain.Menu
	if isAdmin {
		menus, err = uc.menuRepo.GetAllMenus()
	} else {
		permissions, err := uc.userRepo.GetUserPermissions(userID)
		if err != nil {
			return nil, err
		}
		permissionCodes := make([]string, len(permissions))
		for i, p := range permissions {
			permissionCodes[i] = p.Code
		}
		menus, err = uc.menuRepo.GetMenusByPermissionCodes(permissionCodes)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	return buildMenuTree(menus, nil), nil
}

func (uc *MenuUsecase) ListMenus(params *menudomain.MenuListParams) ([]menudomain.MenuListItem, int64, error) {
	menus, total, err := uc.menuRepo.List(params)
	if err != nil {
		return nil, 0, err
	}

	allMenus, err := uc.menuRepo.GetAllMenusRaw()
	if err != nil {
		return nil, 0, err
	}
	idMap := map[uint64]menudomain.Menu{}
	for _, menu := range allMenus {
		idMap[menu.ID] = menu
	}

	items := make([]menudomain.MenuListItem, 0, len(menus))
	for _, menu := range menus {
		parentTitle := ""
		if menu.ParentID != nil {
			if parent, ok := idMap[*menu.ParentID]; ok {
				parentTitle = parent.Title
			}
		}
		items = append(items, menudomain.MenuListItem{
			Menu:        menu,
			ParentTitle: parentTitle,
			FullPath:    buildFullPath(menu, idMap),
		})
	}

	return items, total, nil
}

func (uc *MenuUsecase) CreateMenu(menu *menudomain.Menu) error {
	return uc.menuRepo.Create(menu)
}

func (uc *MenuUsecase) UpdateMenu(menu *menudomain.Menu) error {
	return uc.menuRepo.Update(menu)
}

func (uc *MenuUsecase) DeleteMenu(id uint64) error {
	return uc.menuRepo.Delete(id)
}

func (uc *MenuUsecase) UpdateMenuStatus(id uint64, status string) error {
	return uc.menuRepo.UpdateStatus(id, status)
}

func buildMenuTree(menus []menudomain.Menu, parentID *uint64) []*menudomain.MenuTree {
	var result []*menudomain.MenuTree

	for _, menu := range menus {
		isMatch := false
		if parentID == nil && menu.ParentID == nil {
			isMatch = true
		} else if parentID != nil && menu.ParentID != nil && *parentID == *menu.ParentID {
			isMatch = true
		}

		if isMatch {
			node := &menudomain.MenuTree{
				ID:             menu.ID,
				ParentID:       menu.ParentID,
				Title:          menu.Title,
				TitleEn:        menu.TitleEn,
				Code:           menu.Code,
				Path:           menu.Path,
				Component:      menu.Component,
				Icon:           menu.Icon,
				Sort:           menu.Sort,
				IsHidden:       menu.IsHidden,
				PermissionCode: menu.PermissionCode,
				Children:       buildMenuTree(menus, &menu.ID),
			}
			result = append(result, node)
		}
	}

	return result
}

func buildFullPath(menu menudomain.Menu, idMap map[uint64]menudomain.Menu) string {
	parts := []string{menu.Title}
	current := menu
	for current.ParentID != nil {
		parent, ok := idMap[*current.ParentID]
		if !ok {
			break
		}
		parts = append([]string{parent.Title}, parts...)
		current = parent
	}
	return strings.Join(parts, " / ")
}
