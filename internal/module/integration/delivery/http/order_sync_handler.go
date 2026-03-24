package http

import (
	"strconv"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/response"
	integrationDomain "am-erp-go/internal/module/integration/domain"
	integrationUsecase "am-erp-go/internal/module/integration/usecase"

	"github.com/gin-gonic/gin"
)

type OrderSyncHandler struct {
	registry *integrationUsecase.OrderSyncRegistry
}

func NewOrderSyncHandler(registry *integrationUsecase.OrderSyncRegistry) *OrderSyncHandler {
	return &OrderSyncHandler{registry: registry}
}

func (h *OrderSyncHandler) SyncOrders(c *gin.Context) {
	provider := c.Param("provider")
	operatorID := parseOperatorID(c)
	result, err := h.registry.SyncOrders(c.Request.Context(), provider, integrationDomain.SyncTriggerManual, operatorID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

func (h *OrderSyncHandler) GetState(c *gin.Context) {
	provider := c.Param("provider")
	state, err := h.registry.GetState(provider)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, state)
}

func (h *OrderSyncHandler) ListRuns(c *gin.Context) {
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

func parseOperatorID(c *gin.Context) *uint64 {
	if userID, ok := c.Get(auth.UserIDKey); ok {
		if val, castOK := userID.(uint64); castOK {
			return &val
		}
	}
	return nil
}

func parseIntOrDefault(raw string, fallback int) int {
	if v, err := strconv.Atoi(raw); err == nil && v > 0 {
		return v
	}
	return fallback
}
