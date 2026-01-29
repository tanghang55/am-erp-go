package http

import (
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/logistics/domain"
	"am-erp-go/internal/module/logistics/usecase"

	"github.com/gin-gonic/gin"
)

type LogisticsProviderHandler struct {
	usecase *usecase.LogisticsProviderUsecase
}

func NewLogisticsProviderHandler(usecase *usecase.LogisticsProviderUsecase) *LogisticsProviderHandler {
	return &LogisticsProviderHandler{usecase: usecase}
}

// ListProviders 获取物流供应商列表
func (h *LogisticsProviderHandler) ListProviders(c *gin.Context) {
	params := &domain.LogisticsProviderListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	if providerType := c.Query("provider_type"); providerType != "" {
		pt := domain.ProviderType(providerType)
		params.ProviderType = &pt
	}
	if status := c.Query("status"); status != "" {
		s := domain.ProviderStatus(status)
		params.Status = &s
	}
	if keyword := c.Query("keyword"); keyword != "" {
		params.Keyword = &keyword
	}

	providers, total, err := h.usecase.List(params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, providers, total, params.Page)
}

// GetProvider 获取物流供应商详情
func (h *LogisticsProviderHandler) GetProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	provider, err := h.usecase.Get(id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, provider)
}

// CreateProvider 创建物流供应商
func (h *LogisticsProviderHandler) CreateProvider(c *gin.Context) {
	var params domain.CreateProviderParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	provider, err := h.usecase.Create(&params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, provider)
}

// UpdateProvider 更新物流供应商
func (h *LogisticsProviderHandler) UpdateProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var params domain.UpdateProviderParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.usecase.Update(id, &params); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteProvider 删除物流供应商
func (h *LogisticsProviderHandler) DeleteProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func parseIntOrDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}
