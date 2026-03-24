package http

import (
	"errors"
	"strconv"

	"am-erp-go/internal/infrastructure/response"
	systemdomain "am-erp-go/internal/module/system/domain"

	"github.com/gin-gonic/gin"
)

type ConfigCenterUsecase interface {
	ListModules() ([]systemdomain.ConfigCenterModuleSummary, error)
	GetModule(moduleCode string) (*systemdomain.ConfigCenterModule, error)
	UpdateModule(c *gin.Context, moduleCode string, values map[string]string, operatorID *uint64) (*systemdomain.ConfigCenterModule, error)
}

type ConfigCenterHandler struct {
	usecase ConfigCenterUsecase
}

func NewConfigCenterHandler(usecase ConfigCenterUsecase) *ConfigCenterHandler {
	return &ConfigCenterHandler{usecase: usecase}
}

func (h *ConfigCenterHandler) ListModules(c *gin.Context) {
	items, err := h.usecase.ListModules()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *ConfigCenterHandler) GetModule(c *gin.Context) {
	module, err := h.usecase.GetModule(c.Param("module"))
	if err != nil {
		if errors.Is(err, systemdomain.ErrConfigCenterInvalid) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, module)
}

func (h *ConfigCenterHandler) UpdateModule(c *gin.Context) {
	var input systemdomain.ConfigCenterUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	var operatorID *uint64
	if raw := c.Query("operator_id"); raw != "" {
		if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
			operatorID = &parsed
		}
	}
	module, err := h.usecase.UpdateModule(c, c.Param("module"), input.Values, operatorID)
	if err != nil {
		if errors.Is(err, systemdomain.ErrConfigCenterInvalid) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, module)
}
