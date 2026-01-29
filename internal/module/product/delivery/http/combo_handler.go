package http

import (
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/product/domain"
	"am-erp-go/internal/module/product/usecase"

	"github.com/gin-gonic/gin"
)

type ComboHandler struct {
	comboUsecase *usecase.ProductComboUsecase
}

func NewComboHandler(comboUsecase *usecase.ProductComboUsecase) *ComboHandler {
	return &ComboHandler{comboUsecase: comboUsecase}
}

// ListCombos 获取组合列表
func (h *ComboHandler) ListCombos(c *gin.Context) {
	params := &domain.ComboListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	combos, total, err := h.comboUsecase.ListCombos(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, combos, total, params.Page, params.PageSize)
}

// GetCombo 获取组合详情
func (h *ComboHandler) GetCombo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	combo, err := h.comboUsecase.GetCombo(id)
	if err != nil {
		response.NotFound(c, "combo not found")
		return
	}

	response.Success(c, combo)
}

type comboUpsertRequest struct {
	MainProductID uint64   `json:"main_product_id"`
	ProductIDs    []uint64 `json:"product_ids"`
}

// CreateCombo 创建组合
func (h *ComboHandler) CreateCombo(c *gin.Context) {
	var req comboUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	combo, err := h.comboUsecase.CreateCombo(domain.ComboUpsertParams{
		MainProductID: req.MainProductID,
		ProductIDs:    req.ProductIDs,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, combo)
}

// UpdateCombo 更新组合
func (h *ComboHandler) UpdateCombo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req comboUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.MainProductID == 0 {
		req.MainProductID = id
	}

	combo, err := h.comboUsecase.UpdateComboByMainProductID(id, domain.ComboUpsertParams{
		MainProductID: req.MainProductID,
		ProductIDs:    req.ProductIDs,
	})
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, combo)
}

// DeleteCombo 删除组合
func (h *ComboHandler) DeleteCombo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.comboUsecase.DeleteCombo(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
