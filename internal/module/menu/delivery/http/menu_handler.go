package http

import (
	"strconv"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/response"
	menudomain "am-erp-go/internal/module/menu/domain"

	"github.com/gin-gonic/gin"
)

type MenuUsecase interface {
	GetMenuTree(userID uint64) ([]*menudomain.MenuTree, error)
	ListMenus(params *menudomain.MenuListParams) ([]menudomain.MenuListItem, int64, error)
	CreateMenu(menu *menudomain.Menu) error
	UpdateMenu(menu *menudomain.Menu) error
	UpdateMenuStatus(id uint64, status string) error
	DeleteMenu(id uint64) error
}

type MenuHandler struct {
	menuUsecase MenuUsecase
}

func NewMenuHandler(menuUsecase MenuUsecase) *MenuHandler {
	return &MenuHandler{menuUsecase: menuUsecase}
}

func (h *MenuHandler) GetMenuTree(c *gin.Context) {
	userID, exists := c.Get(auth.UserIDKey)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	menus, err := h.menuUsecase.GetMenuTree(userID.(uint64))
	if err != nil {
		response.InternalError(c, "Failed to get menu tree")
		return
	}

	response.Success(c, menus)
}

func (h *MenuHandler) ListMenus(c *gin.Context) {
	params := &menudomain.MenuListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:  c.Query("keyword"),
		Status:   c.Query("status"),
	}

	if isHidden := c.Query("is_hidden"); isHidden != "" {
		if v, err := strconv.ParseUint(isHidden, 10, 8); err == nil {
			value := uint8(v)
			params.IsHidden = &value
		}
	}

	if parentID := c.Query("parent_id"); parentID != "" {
		if id, err := strconv.ParseUint(parentID, 10, 64); err == nil {
			params.ParentID = &id
		}
	}

	items, total, err := h.menuUsecase.ListMenus(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, items, total, params.Page, params.PageSize)
}

func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var menu menudomain.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.menuUsecase.CreateMenu(&menu); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, menu)
}

func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var menu menudomain.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	menu.ID = id
	if err := h.menuUsecase.UpdateMenu(&menu); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, menu)
}

func (h *MenuHandler) UpdateMenuStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		response.BadRequest(c, "invalid status")
		return
	}

	if err := h.menuUsecase.UpdateMenuStatus(id, req.Status); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.menuUsecase.DeleteMenu(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}
