package http

import (
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/logistics/domain"
	"am-erp-go/internal/module/logistics/usecase"

	"github.com/gin-gonic/gin"
)

type LogisticsServiceHandler struct {
	usecase *usecase.LogisticsServiceUsecase
}

func NewLogisticsServiceHandler(usecase *usecase.LogisticsServiceUsecase) *LogisticsServiceHandler {
	return &LogisticsServiceHandler{usecase: usecase}
}

// ListLogisticsServices 获取物流服务列表
func (h *LogisticsServiceHandler) ListLogisticsServices(c *gin.Context) {
	params := &domain.LogisticsServiceListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	if transportMode := c.Query("transport_mode"); transportMode != "" {
		tm := domain.TransportMode(transportMode)
		params.TransportMode = &tm
	}
	if status := c.Query("status"); status != "" {
		s := domain.ServiceStatus(status)
		params.Status = &s
	}
	if keyword := c.Query("keyword"); keyword != "" {
		params.Keyword = &keyword
	}

	services, total, err := h.usecase.List(params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, services, total, params.Page)
}

// GetLogisticsService 获取物流服务详情
func (h *LogisticsServiceHandler) GetLogisticsService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	service, err := h.usecase.Get(id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, service)
}

// CreateLogisticsService 创建物流服务
func (h *LogisticsServiceHandler) CreateLogisticsService(c *gin.Context) {
	var params domain.CreateLogisticsServiceParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	service, err := h.usecase.Create(&params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, service)
}

// UpdateLogisticsService 更新物流服务
func (h *LogisticsServiceHandler) UpdateLogisticsService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var params domain.UpdateLogisticsServiceParams
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

// DeleteLogisticsService 删除物流服务
func (h *LogisticsServiceHandler) DeleteLogisticsService(c *gin.Context) {
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

// GetActiveServices 获取所有启用的服务
func (h *LogisticsServiceHandler) GetActiveServices(c *gin.Context) {
	services, err := h.usecase.GetActiveServices()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, services)
}

// GetServicesByTransportMode 根据运输方式获取服务
func (h *LogisticsServiceHandler) GetServicesByTransportMode(c *gin.Context) {
	transportMode := c.Query("transport_mode")
	if transportMode == "" {
		response.BadRequest(c, "transport_mode is required")
		return
	}

	services, err := h.usecase.GetServicesByTransportMode(domain.TransportMode(transportMode))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, services)
}

// RegisterServiceRoutes 注册物流服务路由
func (h *LogisticsServiceHandler) RegisterServiceRoutes(group *gin.RouterGroup) {
	services := group.Group("/logistics-services")
	{
		services.GET("", h.ListLogisticsServices)
		services.GET("/active", h.GetActiveServices)
		services.GET("/by-transport-mode", h.GetServicesByTransportMode)
		services.GET("/:id", h.GetLogisticsService)
		services.POST("", h.CreateLogisticsService)
		services.PUT("/:id", h.UpdateLogisticsService)
		services.DELETE("/:id", h.DeleteLogisticsService)
	}
}
