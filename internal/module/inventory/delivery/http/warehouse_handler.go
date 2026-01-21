package http

import (
	"net/http"

	"am-erp-go/internal/module/inventory/usecase"

	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	usecase *usecase.WarehouseUsecase
}

func NewWarehouseHandler(usecase *usecase.WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{usecase: usecase}
}

// GetActiveWarehouses 获取所有启用的仓库
func (h *WarehouseHandler) GetActiveWarehouses(c *gin.Context) {
	warehouses, err := h.usecase.GetActiveWarehouses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    warehouses,
	})
}
