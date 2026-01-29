package http

import (
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/shipment/domain"
	"am-erp-go/internal/module/shipment/usecase"

	"github.com/gin-gonic/gin"
)

type PackageSpecHandler struct {
	uc *usecase.PackageSpecUseCase
}

func NewPackageSpecHandler(uc *usecase.PackageSpecUseCase) *PackageSpecHandler {
	return &PackageSpecHandler{uc: uc}
}

func (h *PackageSpecHandler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/package-specs")
	{
		group.GET("", h.List)
		group.GET("/:id", h.GetByID)
		group.POST("", h.Create)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
		group.GET("/:id/packaging-items", h.GetPackageSpecPackagingItems)
		group.PUT("/:id/packaging-items", h.SavePackageSpecPackagingItems)
	}
}

// List 获取装箱规格列表
func (h *PackageSpecHandler) List(c *gin.Context) {
	params := &domain.PackageSpecListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
	}

	if keyword := c.Query("keyword"); keyword != "" {
		params.Keyword = &keyword
	}

	if status := c.Query("status"); status != "" {
		params.Status = &status
	}

	specs, total, err := h.uc.List(params)
	if err != nil {
		response.ServerError(c, "获取装箱规格列表失败")
		return
	}

	response.SuccessPage(c, specs, total, params.Page)
}

// GetByID 获取装箱规格详情
func (h *PackageSpecHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	spec, err := h.uc.GetByID(id)
	if err != nil {
		response.NotFound(c, "装箱规格不存在")
		return
	}

	response.Success(c, spec)
}

// Create 创建装箱规格
func (h *PackageSpecHandler) Create(c *gin.Context) {
	var req struct {
		Name           string  `json:"name" binding:"required"`
		Length         float64 `json:"length" binding:"required"`
		Width          float64 `json:"width" binding:"required"`
		Height         float64 `json:"height" binding:"required"`
		Weight         float64 `json:"weight" binding:"required"`
		QuantityPerBox uint    `json:"quantity_per_box"`
		Remark         *string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := getUserIDFromContext(c)

	params := &domain.CreatePackageSpecParams{
		Name:           req.Name,
		Length:         req.Length,
		Width:          req.Width,
		Height:         req.Height,
		Weight:         req.Weight,
		QuantityPerBox: req.QuantityPerBox,
		Remark:         req.Remark,
		CreatedBy:      userID,
	}

	spec, err := h.uc.Create(params)
	if err != nil {
		response.ServerError(c, "创建装箱规格失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", spec)
}

// Update 更新装箱规格
func (h *PackageSpecHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Name           *string  `json:"name"`
		Length         *float64 `json:"length"`
		Width          *float64 `json:"width"`
		Height         *float64 `json:"height"`
		Weight         *float64 `json:"weight"`
		QuantityPerBox *uint    `json:"quantity_per_box"`
		Remark         *string  `json:"remark"`
		Status         *string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	userID := getUserIDFromContext(c)

	params := &domain.UpdatePackageSpecParams{
		Name:           req.Name,
		Length:         req.Length,
		Width:          req.Width,
		Height:         req.Height,
		Weight:         req.Weight,
		QuantityPerBox: req.QuantityPerBox,
		Remark:         req.Remark,
		Status:         req.Status,
		UpdatedBy:      userID,
	}

	spec, err := h.uc.Update(id, params)
	if err != nil {
		response.ServerError(c, "更新装箱规格失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", spec)
}

// Delete 删除装箱规格
func (h *PackageSpecHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.uc.Delete(id); err != nil {
		response.ServerError(c, "删除装箱规格失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetPackageSpecPackagingItems 获取装箱规格的包材配置列表
func (h *PackageSpecHandler) GetPackageSpecPackagingItems(c *gin.Context) {
	packageSpecID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid package spec id")
		return
	}

	items, err := h.uc.GetPackageSpecPackagingItems(packageSpecID)
	if err != nil {
		response.ServerError(c, "获取包材配置失败")
		return
	}

	response.Success(c, items)
}

// SavePackageSpecPackagingItems 保存装箱规格的包材配置
func (h *PackageSpecHandler) SavePackageSpecPackagingItems(c *gin.Context) {
	packageSpecID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid package spec id")
		return
	}

	var req struct {
		PackagingItems []domain.PackageSpecPackagingItem `json:"packaging_items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.uc.SavePackageSpecPackagingItems(packageSpecID, req.PackagingItems); err != nil {
		response.ServerError(c, "保存包材配置失败")
		return
	}

	// 返回保存后的数据
	items, err := h.uc.GetPackageSpecPackagingItems(packageSpecID)
	if err != nil {
		response.ServerError(c, "获取包材配置失败")
		return
	}

	response.SuccessWithMessage(c, "保存成功", items)
}
