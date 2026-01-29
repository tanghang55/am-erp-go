package http

import "github.com/gin-gonic/gin"

func (h *LogisticsProviderHandler) RegisterProviderRoutes(r *gin.RouterGroup) {
	providers := r.Group("/logistics-providers")
	{
		providers.GET("", h.ListProviders)
		providers.GET("/:id", h.GetProvider)
		providers.POST("", h.CreateProvider)
		providers.PUT("/:id", h.UpdateProvider)
		providers.DELETE("/:id", h.DeleteProvider)
	}
}

func (h *ShippingRateHandler) RegisterRateRoutes(r *gin.RouterGroup) {
	rates := r.Group("/shipping-rates")
	{
		rates.GET("", h.ListShippingRates)
		rates.GET("/:id", h.GetShippingRate)
		rates.POST("", h.CreateShippingRate)
		rates.PUT("/:id", h.UpdateShippingRate)
		rates.DELETE("/:id", h.DeleteShippingRate)
		rates.GET("/query-latest", h.QueryLatestRate) // 查询最新报价
	}
}
