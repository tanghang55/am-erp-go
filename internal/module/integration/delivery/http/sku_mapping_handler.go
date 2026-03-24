package http

import (
	"fmt"
	"strconv"
	"strings"

	"am-erp-go/internal/infrastructure/response"
	integrationDomain "am-erp-go/internal/module/integration/domain"
	integrationUsecase "am-erp-go/internal/module/integration/usecase"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type SKUMappingHandler struct {
	usecase     *integrationUsecase.SKUMappingUsecase
	auditLogger AuthorizationAuditLogger
}

func NewSKUMappingHandler(usecase *integrationUsecase.SKUMappingUsecase) *SKUMappingHandler {
	return &SKUMappingHandler{usecase: usecase}
}

func (h *SKUMappingHandler) BindAuditLogger(logger AuthorizationAuditLogger) {
	h.auditLogger = logger
}

func (h *SKUMappingHandler) List(c *gin.Context) {
	params := &integrationDomain.SKUMappingListParams{
		Page:         parseIntOrDefault(c.Query("page"), 1),
		PageSize:     parseIntOrDefault(c.Query("page_size"), 20),
		ProviderCode: c.Query("provider_code"),
		Marketplace:  c.Query("marketplace"),
		Status:       c.Query("status"),
		Keyword:      c.Query("keyword"),
	}
	if productIDRaw := strings.TrimSpace(c.Query("product_id")); productIDRaw != "" {
		if productID, err := strconv.ParseUint(productIDRaw, 10, 64); err == nil {
			params.ProductID = &productID
		}
	}
	list, total, err := h.usecase.List(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, params.Page, params.PageSize)
}

type createSKUMappingRequest struct {
	ProviderCode string `json:"provider_code" binding:"required"`
	Marketplace  string `json:"marketplace" binding:"required"`
	SellerSKU    string `json:"seller_sku" binding:"required"`
	ProductID    uint64 `json:"product_id" binding:"required"`
	Status       string `json:"status"`
	Remark       string `json:"remark"`
}

func (h *SKUMappingHandler) Create(c *gin.Context) {
	var req createSKUMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	operatorID := parseOperatorIDFromContext(c)
	item, err := h.usecase.Create(&integrationUsecase.CreateSKUMappingInput{
		ProviderCode: req.ProviderCode,
		Marketplace:  req.Marketplace,
		SellerSKU:    req.SellerSKU,
		ProductID:    req.ProductID,
		Status:       integrationDomain.SKUMappingStatus(strings.ToUpper(strings.TrimSpace(req.Status))),
		Remark:       req.Remark,
		OperatorID:   operatorID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	h.recordAudit(c, "CREATE_SKU_MAPPING", "IntegrationSKUMapping", fmt.Sprintf("%d", item.ID), nil, item)
	response.Success(c, item)
}

type updateSKUMappingRequest struct {
	ProductID uint64 `json:"product_id"`
	Status    string `json:"status"`
	Remark    string `json:"remark"`
}

func (h *SKUMappingHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req updateSKUMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	before, _ := h.usecase.GetByID(id)
	operatorID := parseOperatorIDFromContext(c)
	item, err := h.usecase.Update(id, &integrationUsecase.UpdateSKUMappingInput{
		ProductID:  req.ProductID,
		Status:     integrationDomain.SKUMappingStatus(strings.ToUpper(strings.TrimSpace(req.Status))),
		Remark:     req.Remark,
		OperatorID: operatorID,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	h.recordAuditIfChanged(c, "UPDATE_SKU_MAPPING", "IntegrationSKUMapping", fmt.Sprintf("%d", id), before, item)
	response.Success(c, item)
}

func (h *SKUMappingHandler) recordAudit(c *gin.Context, action, entityType, entityID string, before, after any) {
	if h == nil || h.auditLogger == nil || c == nil {
		return
	}
	_ = h.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Integration",
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
}

func (h *SKUMappingHandler) recordAuditIfChanged(c *gin.Context, action, entityType, entityID string, before, after any) {
	beforeDiff, afterDiff, changed := buildIntegrationAuditDiff(before, after)
	if !changed {
		return
	}
	h.recordAudit(c, action, entityType, entityID, beforeDiff, afterDiff)
}
