package usecase

import (
    "testing"

    "am-erp-go/internal/module/identity/domain"
    menudomain "am-erp-go/internal/module/menu/domain"
)

type stubMenuRepo struct {
	all []menudomain.Menu
	created   *menudomain.Menu
	updated   *menudomain.Menu
	deletedID *uint64
	statusID  uint64
	status    string
}

func (s *stubMenuRepo) GetMenusByPermissionCodes(_ []string) ([]menudomain.Menu, error) {
    return s.all, nil
}

func (s *stubMenuRepo) GetAllMenus() ([]menudomain.Menu, error) {
    return s.all, nil
}

func (s *stubMenuRepo) GetAllMenusRaw() ([]menudomain.Menu, error) {
    return s.all, nil
}

func (s *stubMenuRepo) List(_ *menudomain.MenuListParams) ([]menudomain.Menu, int64, error) {
	return s.all, int64(len(s.all)), nil
}

func (s *stubMenuRepo) GetByID(_ uint64) (*menudomain.Menu, error) {
    return nil, nil
}

func (s *stubMenuRepo) Create(menu *menudomain.Menu) error {
	s.created = menu
	return nil
}

func (s *stubMenuRepo) Update(menu *menudomain.Menu) error {
	s.updated = menu
	return nil
}

func (s *stubMenuRepo) Delete(id uint64) error {
	s.deletedID = &id
	return nil
}

func (s *stubMenuRepo) UpdateStatus(id uint64, status string) error {
	s.statusID = id
	s.status = status
	return nil
}

type stubUserRepo struct {
    roles       []domain.Role
    permissions []domain.Permission
}

func (s *stubUserRepo) FindByUsername(_ string) (*domain.User, error) { return nil, nil }
func (s *stubUserRepo) FindByID(_ uint64) (*domain.User, error)        { return nil, nil }
func (s *stubUserRepo) GetUserRoles(_ uint64) ([]domain.Role, error)   { return s.roles, nil }
func (s *stubUserRepo) GetUserPermissions(_ uint64) ([]domain.Permission, error) {
    return s.permissions, nil
}

func TestMenuTreeBuildsHierarchy(t *testing.T) {
    rootID := uint64(1)
    menus := []menudomain.Menu{
        {ID: rootID, Title: "System", Code: "system", Sort: 1},
        {ID: 2, Title: "Menu List", Code: "menu-list", ParentID: &rootID, Sort: 1},
    }

    uc := NewMenuUsecase(&stubMenuRepo{all: menus}, &stubUserRepo{})

    tree, err := uc.GetMenuTree(1)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(tree) != 1 || len(tree[0].Children) != 1 {
        t.Fatalf("expected one root with one child, got: %#v", tree)
    }
}

func TestListMenusAddsParentAndFullPath(t *testing.T) {
	rootID := uint64(10)
	menus := []menudomain.Menu{
		{ID: rootID, Title: "System", Code: "system"},
		{ID: 11, Title: "Menu List", Code: "menu-list", ParentID: &rootID},
	}
	repo := &stubMenuRepo{all: menus}
	uc := NewMenuUsecase(repo, &stubUserRepo{})

	items, total, err := uc.ListMenus(&menudomain.MenuListParams{Page: 1, PageSize: 20})
	if err != nil || total != int64(len(menus)) {
		t.Fatalf("unexpected result: total=%d err=%v", total, err)
	}
	if items[1].ParentTitle != "System" || items[1].FullPath != "System / Menu List" {
		t.Fatalf("expected parent title and full path, got: %#v", items[1])
	}
}

func TestCreateMenuCallsRepo(t *testing.T) {
	repo := &stubMenuRepo{}
	uc := NewMenuUsecase(repo, &stubUserRepo{})

	menu := menudomain.Menu{Title: "Menu", Code: "menu"}
	if err := uc.CreateMenu(&menu); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created != &menu {
		t.Fatalf("expected repo to receive menu")
	}
}

func TestUpdateMenuCallsRepo(t *testing.T) {
	repo := &stubMenuRepo{}
	uc := NewMenuUsecase(repo, &stubUserRepo{})

	menu := menudomain.Menu{ID: 11, Title: "Menu", Code: "menu"}
	if err := uc.UpdateMenu(&menu); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updated != &menu {
		t.Fatalf("expected repo to receive updated menu")
	}
}

func TestDeleteMenuCallsRepo(t *testing.T) {
	repo := &stubMenuRepo{}
	uc := NewMenuUsecase(repo, &stubUserRepo{})

	if err := uc.DeleteMenu(12); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.deletedID == nil || *repo.deletedID != 12 {
		t.Fatalf("expected repo to receive deleted id")
	}
}

func TestUpdateMenuStatusCallsRepo(t *testing.T) {
	repo := &stubMenuRepo{}
	uc := NewMenuUsecase(repo, &stubUserRepo{})

	if err := uc.UpdateMenuStatus(13, "DISABLED"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.statusID != 13 || repo.status != "DISABLED" {
		t.Fatalf("expected repo to receive status change")
	}
}
