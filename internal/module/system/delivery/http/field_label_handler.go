package http

import (
	"net/http"
	"strconv"

	systemdomain "am-erp-go/internal/module/system/domain"

	"github.com/gin-gonic/gin"
)

type FieldLabelUsecase interface {
	GetLabels(locale string) (map[string]string, error)
	List(page, pageSize int, keyword string) ([]*systemdomain.FieldLabel, int64, error)
	Create(label *systemdomain.FieldLabel) error
	Update(label *systemdomain.FieldLabel) error
	Delete(id uint64) error
}

type FieldLabelHandler struct {
	usecase FieldLabelUsecase
}

func NewFieldLabelHandler(usecase FieldLabelUsecase) *FieldLabelHandler {
	return &FieldLabelHandler{usecase: usecase}
}

func (h *FieldLabelHandler) GetLabels(c *gin.Context) {
	locale := c.Query("locale")
	if locale == "" {
		locale = "zh-CN"
	}

	labels, err := h.usecase.GetLabels(locale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"locale": locale,
			"labels": labels,
		},
	})
}

func (h *FieldLabelHandler) List(c *gin.Context) {
	page := parseIntOrDefault(c.Query("page"), 1)
	pageSize := parseIntOrDefault(c.Query("page_size"), 20)
	keyword := c.Query("keyword")

	items, total, err := h.usecase.List(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"items":     items,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (h *FieldLabelHandler) Create(c *gin.Context) {
	var label systemdomain.FieldLabel
	if err := c.ShouldBindJSON(&label); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := h.usecase.Create(&label); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": label})
}

func (h *FieldLabelHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var label systemdomain.FieldLabel
	if err := c.ShouldBindJSON(&label); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	label.ID = id
	if err := h.usecase.Update(&label); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": label})
}

func (h *FieldLabelHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}
