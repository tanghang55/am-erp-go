package http

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    menudomain "am-erp-go/internal/module/menu/domain"

    "github.com/gin-gonic/gin"
)

type stubMenuUsecase struct {
    tree         []*menudomain.MenuTree
    listItems    []menudomain.MenuListItem
    total        int64
    listParams   *menudomain.MenuListParams
    createdMenu  *menudomain.Menu
    updatedMenu  *menudomain.Menu
    statusID     uint64
    status       string
    deletedID    uint64
    treeUserID   uint64
    err          error
}

func (s *stubMenuUsecase) GetMenuTree(userID uint64) ([]*menudomain.MenuTree, error) {
    s.treeUserID = userID
    return s.tree, s.err
}

func (s *stubMenuUsecase) ListMenus(params *menudomain.MenuListParams) ([]menudomain.MenuListItem, int64, error) {
    s.listParams = params
    return s.listItems, s.total, s.err
}

func (s *stubMenuUsecase) CreateMenu(menu *menudomain.Menu) error {
    s.createdMenu = menu
    return s.err
}

func (s *stubMenuUsecase) UpdateMenu(menu *menudomain.Menu) error {
    s.updatedMenu = menu
    return s.err
}

func (s *stubMenuUsecase) UpdateMenuStatus(id uint64, status string) error {
    s.statusID = id
    s.status = status
    return s.err
}

func (s *stubMenuUsecase) DeleteMenu(id uint64) error {
    s.deletedID = id
    return s.err
}

func TestGetMenuTreeRequiresAuth(t *testing.T) {
    gin.SetMode(gin.TestMode)

    handler := NewMenuHandler(&stubMenuUsecase{})
    router := gin.New()
    router.GET("/api/v1/menus/tree", handler.GetMenuTree)

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/api/v1/menus/tree", nil)
    router.ServeHTTP(w, req)

    if w.Code != http.StatusUnauthorized {
        t.Fatalf("expected 401, got %d", w.Code)
    }
}

func TestListMenusParsesParams(t *testing.T) {
    gin.SetMode(gin.TestMode)

    stub := &stubMenuUsecase{
        listItems: []menudomain.MenuListItem{{Menu: menudomain.Menu{ID: 1, Title: "System", Code: "system"}}},
        total:     1,
    }
    handler := NewMenuHandler(stub)

    router := gin.New()
    router.GET("/api/v1/menus", handler.ListMenus)

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/api/v1/menus?page=2&page_size=10&keyword=abc&status=ACTIVE&is_hidden=1&parent_id=5", nil)
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    if stub.listParams == nil {
        t.Fatalf("expected list params to be captured")
    }
    if stub.listParams.Page != 2 || stub.listParams.PageSize != 10 {
        t.Fatalf("unexpected pagination: %#v", stub.listParams)
    }
    if stub.listParams.Keyword != "abc" || stub.listParams.Status != "ACTIVE" {
        t.Fatalf("unexpected filters: %#v", stub.listParams)
    }
    if stub.listParams.ParentID == nil || *stub.listParams.ParentID != 5 {
        t.Fatalf("expected parent_id 5, got %#v", stub.listParams.ParentID)
    }
    if stub.listParams.IsHidden == nil || *stub.listParams.IsHidden != 1 {
        t.Fatalf("expected is_hidden 1, got %#v", stub.listParams.IsHidden)
    }
}

func TestCreateMenuCallsUsecase(t *testing.T) {
    gin.SetMode(gin.TestMode)

    stub := &stubMenuUsecase{}
    handler := NewMenuHandler(stub)
    router := gin.New()
    router.POST("/api/v1/menus", handler.CreateMenu)

    payload := []byte(`{"title":"Menu","code":"menu"}`)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/api/v1/menus", bytes.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    if stub.createdMenu == nil || stub.createdMenu.Title != "Menu" {
        t.Fatalf("expected created menu to be captured")
    }
}

func TestUpdateMenuCallsUsecase(t *testing.T) {
    gin.SetMode(gin.TestMode)

    stub := &stubMenuUsecase{}
    handler := NewMenuHandler(stub)
    router := gin.New()
    router.PUT("/api/v1/menus/:id", handler.UpdateMenu)

    payload := []byte(`{"title":"Menu"}`)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPut, "/api/v1/menus/12", bytes.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    if stub.updatedMenu == nil || stub.updatedMenu.ID != 12 {
        t.Fatalf("expected updated menu id to be 12")
    }
}

func TestUpdateMenuStatusCallsUsecase(t *testing.T) {
    gin.SetMode(gin.TestMode)

    stub := &stubMenuUsecase{}
    handler := NewMenuHandler(stub)
    router := gin.New()
    router.PATCH("/api/v1/menus/:id/status", handler.UpdateMenuStatus)

    body := map[string]string{"status": "DISABLED"}
    payload, _ := json.Marshal(body)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPatch, "/api/v1/menus/18/status", bytes.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    if stub.statusID != 18 || stub.status != "DISABLED" {
        t.Fatalf("expected status update to be captured")
    }
}

func TestDeleteMenuCallsUsecase(t *testing.T) {
    gin.SetMode(gin.TestMode)

    stub := &stubMenuUsecase{}
    handler := NewMenuHandler(stub)
    router := gin.New()
    router.DELETE("/api/v1/menus/:id", handler.DeleteMenu)

    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodDelete, "/api/v1/menus/21", nil)
    req.Header.Set("Authorization", "Bearer test")
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }
    if stub.deletedID != 21 {
        t.Fatalf("expected deleted id to be 21")
    }
}
