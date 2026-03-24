package http

import (
	"am-erp-go/internal/infrastructure/response"
	integrationDomain "am-erp-go/internal/module/integration/domain"
	integrationUsecase "am-erp-go/internal/module/integration/usecase"

	"github.com/gin-gonic/gin"
)

type RefundSyncHandler struct {
	registry *integrationUsecase.RefundSyncRegistry
}

func NewRefundSyncHandler(registry *integrationUsecase.RefundSyncRegistry) *RefundSyncHandler {
	return &RefundSyncHandler{registry: registry}
}

func (h *RefundSyncHandler) SyncRefunds(c *gin.Context) {
	provider := c.Param("provider")
	operatorID := parseOperatorID(c)
	result, err := h.registry.SyncRefunds(c.Request.Context(), provider, integrationDomain.SyncTriggerManual, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

func (h *RefundSyncHandler) GetState(c *gin.Context) {
	provider := c.Param("provider")
	state, err := h.registry.GetState(provider)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, state)
}

func (h *RefundSyncHandler) ListRuns(c *gin.Context) {
	provider := c.Param("provider")
	page := parseIntOrDefault(c.Query("page"), 1)
	pageSize := parseIntOrDefault(c.Query("page_size"), 20)
	list, total, err := h.registry.ListRuns(provider, page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, page, pageSize)
}
