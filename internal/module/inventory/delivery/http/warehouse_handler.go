package http

import (
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/inventory/usecase"

	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	usecase *usecase.WarehouseUsecase
}

func NewWarehouseHandler(usecase *usecase.WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{usecase: usecase}
}

// ListWarehouses 获取仓库列表
func (h *WarehouseHandler) ListWarehouses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := &domain.WarehouseListParams{
		Page:     page,
		PageSize: pageSize,
	}

	if warehouseType := c.Query("type"); warehouseType != "" {
		t := domain.WarehouseType(warehouseType)
		params.Type = &t
	}

	if status := c.Query("status"); status != "" {
		s := domain.WarehouseStatus(status)
		params.Status = &s
	}

	if keyword := c.Query("keyword"); keyword != "" {
		params.Keyword = &keyword
	}

	warehouses, total, err := h.usecase.ListWarehouses(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, warehouses, total, page, pageSize)
}

// GetWarehouse 获取仓库详情
func (h *WarehouseHandler) GetWarehouse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	warehouse, err := h.usecase.GetWarehouse(id)
	if err != nil {
		response.NotFound(c, "warehouse not found")
		return
	}

	response.Success(c, warehouse)
}

// CreateWarehouse 创建仓库
func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var req domain.CreateWarehouseParams
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	warehouse, err := h.usecase.CreateWarehouse(c, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, warehouse)
}

// UpdateWarehouse 更新仓库
func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req domain.UpdateWarehouseParams
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	warehouse, err := h.usecase.UpdateWarehouse(c, id, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, warehouse)
}

// DeleteWarehouse 删除仓库
func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.usecase.DeleteWarehouse(c, id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetActiveWarehouses 获取所有启用的仓库
func (h *WarehouseHandler) GetActiveWarehouses(c *gin.Context) {
	warehouses, err := h.usecase.GetActiveWarehouses()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, warehouses)
}
