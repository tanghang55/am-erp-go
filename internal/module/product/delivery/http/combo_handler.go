package http

import (
	"errors"
	"fmt"
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/product/domain"
	"am-erp-go/internal/module/product/usecase"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type ComboHandler struct {
	comboUsecase *usecase.ProductComboUsecase
	auditLogger  ProductAuditLogger
}

func NewComboHandler(comboUsecase *usecase.ProductComboUsecase) *ComboHandler {
	return &ComboHandler{comboUsecase: comboUsecase}
}

func (h *ComboHandler) BindAuditLogger(logger ProductAuditLogger) {
	h.auditLogger = logger
}

// ListCombos 获取组合列表
func (h *ComboHandler) ListCombos(c *gin.Context) {
	params := &domain.ComboListParams{
		Page:        parseIntOrDefault(c.Query("page"), 1),
		PageSize:    parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:     c.Query("keyword"),
		Marketplace: c.Query("marketplace"),
		Locked:      c.Query("locked"),
	}
	if statuses := parseCSVQuery(c.Query("statuses")); len(statuses) > 0 {
		params.Statuses = statuses
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
	MainProductID uint64 `json:"main_product_id"`
	Children      []struct {
		ProductID uint64 `json:"product_id"`
		QtyRatio  uint64 `json:"qty_ratio"`
	} `json:"children"`
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
		Children:      mapComboChildren(req.Children),
	})
	if err != nil {
		if errors.Is(err, usecase.ErrComboMainProductRequired) ||
			errors.Is(err, usecase.ErrComboProductNotFound) ||
			errors.Is(err, usecase.ErrComboStandaloneOnly) ||
			errors.Is(err, usecase.ErrComboLocked) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	h.recordAudit(c, "CREATE_PRODUCT_COMBO", fmt.Sprintf("%d", combo.ComboID), nil, combo)
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
	before, _ := h.comboUsecase.GetCombo(id)

	combo, err := h.comboUsecase.UpdateCombo(id, domain.ComboUpsertParams{
		MainProductID: req.MainProductID,
		Children:      mapComboChildren(req.Children),
	})
	if err != nil {
		if errors.Is(err, usecase.ErrComboNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		if errors.Is(err, usecase.ErrComboMainProductRequired) ||
			errors.Is(err, usecase.ErrComboProductNotFound) ||
			errors.Is(err, usecase.ErrComboStandaloneOnly) ||
			errors.Is(err, usecase.ErrComboLocked) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	h.recordAuditIfChanged(c, "UPDATE_PRODUCT_COMBO", fmt.Sprintf("%d", id), before, combo)
	response.Success(c, combo)
}

// DeleteCombo 删除组合
func (h *ComboHandler) DeleteCombo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	before, _ := h.comboUsecase.GetCombo(id)
	if err := h.comboUsecase.DeleteCombo(id); err != nil {
		if errors.Is(err, usecase.ErrComboNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		if errors.Is(err, usecase.ErrComboLocked) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	h.recordAudit(c, "DELETE_PRODUCT_COMBO", fmt.Sprintf("%d", id), before, nil)
	response.Success(c, nil)
}

func (h *ComboHandler) recordAudit(c *gin.Context, action, entityID string, before, after any) {
	if h.auditLogger == nil || c == nil {
		return
	}
	_ = h.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Product",
		Action:     action,
		EntityType: "ProductCombo",
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
}

func (h *ComboHandler) recordAuditIfChanged(c *gin.Context, action, entityID string, before, after any) {
	beforeDiff, afterDiff, changed := buildAuditDiff(before, after)
	if !changed {
		return
	}
	h.recordAudit(c, action, entityID, beforeDiff, afterDiff)
}

func mapComboChildren(children []struct {
	ProductID uint64 `json:"product_id"`
	QtyRatio  uint64 `json:"qty_ratio"`
}) []domain.ComboChildInput {
	result := make([]domain.ComboChildInput, 0, len(children))
	for _, child := range children {
		result = append(result, domain.ComboChildInput{
			ProductID: child.ProductID,
			QtyRatio:  child.QtyRatio,
		})
	}
	return result
}
