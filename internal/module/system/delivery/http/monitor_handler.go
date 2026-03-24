package http

import (
	"am-erp-go/internal/infrastructure/response"
	systemdomain "am-erp-go/internal/module/system/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MonitorUsecase interface {
	GetOverview() (*systemdomain.MonitorOverview, error)
	ListRecentJobs(status string, traceID string, limit int) ([]*systemdomain.MonitorRecentJob, error)
	ListRecentLogs(level string, traceID string, limit int) ([]*systemdomain.MonitorRecentLog, error)
}

type MonitorHandler struct {
	usecase MonitorUsecase
}

func NewMonitorHandler(usecase MonitorUsecase) *MonitorHandler {
	return &MonitorHandler{usecase: usecase}
}

func (h *MonitorHandler) GetOverview(c *gin.Context) {
	overview, err := h.usecase.GetOverview()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, overview)
}

func (h *MonitorHandler) ListRecentJobs(c *gin.Context) {
	status := c.Query("status")
	traceID := c.Query("trace_id")
	limit := 10
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			response.BadRequest(c, "limit must be a number")
			return
		}
		limit = parsed
	}

	items, err := h.usecase.ListRecentJobs(status, traceID, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *MonitorHandler) ListRecentLogs(c *gin.Context) {
	level := c.Query("level")
	traceID := c.Query("trace_id")
	limit := 10
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			response.BadRequest(c, "limit must be a number")
			return
		}
		limit = parsed
	}

	items, err := h.usecase.ListRecentLogs(level, traceID, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, items)
}
