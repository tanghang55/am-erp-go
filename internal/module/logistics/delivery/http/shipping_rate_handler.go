package http

import (
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/logistics/domain"
	"am-erp-go/internal/module/logistics/usecase"

	"github.com/gin-gonic/gin"
)

type ShippingRateHandler struct {
	usecase *usecase.ShippingRateUsecase
}

func NewShippingRateHandler(usecase *usecase.ShippingRateUsecase) *ShippingRateHandler {
	return &ShippingRateHandler{usecase: usecase}
}

// ListShippingRates 获取运费报价列表
func (h *ShippingRateHandler) ListShippingRates(c *gin.Context) {
	params := &domain.ShippingRateListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	if providerID := c.Query("provider_id"); providerID != "" {
		if id, err := strconv.ParseUint(providerID, 10, 64); err == nil {
			params.ProviderID = &id
		}
	}
	if originWarehouseID := c.Query("origin_warehouse_id"); originWarehouseID != "" {
		if id, err := strconv.ParseUint(originWarehouseID, 10, 64); err == nil {
			params.OriginWarehouseID = &id
		}
	}
	if destWarehouseID := c.Query("destination_warehouse_id"); destWarehouseID != "" {
		if id, err := strconv.ParseUint(destWarehouseID, 10, 64); err == nil {
			params.DestinationWarehouseID = &id
		}
	}
	if transportMode := c.Query("transport_mode"); transportMode != "" {
		tm := domain.TransportMode(transportMode)
		params.TransportMode = &tm
	}
	if status := c.Query("status"); status != "" {
		s := domain.RateStatus(status)
		params.Status = &s
	}
	if keyword := c.Query("keyword"); keyword != "" {
		params.Keyword = &keyword
	}

	rates, total, err := h.usecase.List(params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, rates, total, params.Page)
}

// GetShippingRate 获取运费报价详情
func (h *ShippingRateHandler) GetShippingRate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	rate, err := h.usecase.Get(id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, rate)
}

// CreateShippingRate 创建运费报价
func (h *ShippingRateHandler) CreateShippingRate(c *gin.Context) {
	var params domain.CreateShippingRateParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rate, err := h.usecase.Create(&params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, rate)
}

// UpdateShippingRate 更新运费报价
func (h *ShippingRateHandler) UpdateShippingRate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var params domain.UpdateShippingRateParams
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

// DeleteShippingRate 删除运费报价
func (h *ShippingRateHandler) DeleteShippingRate(c *gin.Context) {
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

// QueryLatestRate 查询最新有效报价
func (h *ShippingRateHandler) QueryLatestRate(c *gin.Context) {
	var params domain.QueryLatestRateParams

	originWarehouseID, err := strconv.ParseUint(c.Query("origin_warehouse_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "origin_warehouse_id is required")
		return
	}
	params.OriginWarehouseID = originWarehouseID

	destWarehouseID, err := strconv.ParseUint(c.Query("destination_warehouse_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "destination_warehouse_id is required")
		return
	}
	params.DestinationWarehouseID = destWarehouseID

	transportMode := c.Query("transport_mode")
	if transportMode == "" {
		response.BadRequest(c, "transport_mode is required")
		return
	}
	params.TransportMode = domain.TransportMode(transportMode)

	if providerID := c.Query("provider_id"); providerID != "" {
		if id, err := strconv.ParseUint(providerID, 10, 64); err == nil {
			params.ProviderID = &id
		}
	}

	if weight := c.Query("weight"); weight != "" {
		if w, err := strconv.ParseFloat(weight, 64); err == nil {
			params.Weight = &w
		}
	}

	if volume := c.Query("volume"); volume != "" {
		if v, err := strconv.ParseFloat(volume, 64); err == nil {
			params.Volume = &v
		}
	}

	if queryDate := c.Query("query_date"); queryDate != "" {
		params.QueryDate = &queryDate
	}

	rate, err := h.usecase.QueryLatestRate(&params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, rate)
}
