package http

import (
	"net/http"

	systemdomain "am-erp-go/internal/module/system/domain"

	"github.com/gin-gonic/gin"
)

type AuditLogUsecase interface {
	List(params systemdomain.AuditLogListParams) ([]*systemdomain.AuditLog, int64, error)
}

type AuditLogHandler struct {
	usecase AuditLogUsecase
}

func NewAuditLogHandler(usecase AuditLogUsecase) *AuditLogHandler {
	return &AuditLogHandler{usecase: usecase}
}

func (h *AuditLogHandler) List(c *gin.Context) {
	params := systemdomain.AuditLogListParams{
		Page:       parseIntOrDefault(c.Query("page"), 1),
		PageSize:   parseIntOrDefault(c.Query("page_size"), 20),
		Module:     c.Query("module"),
		Action:     c.Query("action"),
		Username:   c.Query("username"),
		EntityType: c.Query("entity_type"),
		EntityID:   c.Query("entity_id"),
		Keyword:    c.Query("keyword"),
		DateFrom:   c.Query("date_from"),
		DateTo:     c.Query("date_to"),
	}

	items, total, err := h.usecase.List(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"data":  items,
			"total": total,
		},
	})
}
