package http

import (
	"errors"
	"io"
	"strconv"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/sales/domain"

	"github.com/gin-gonic/gin"
)

func (h *SalesOrderHandler) ImportSalesOrders(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	stream, err := file.Open()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	defer stream.Close()

	content, err := io.ReadAll(stream)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	operatorID := parseOperatorIDFromFormOrContext(c)
	batch, err := h.usecase.ImportCSV(c.Request.Context(), file.Filename, content, operatorID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrImportDuplicateFile):
			response.BadRequest(c, "duplicate import file")
		case errors.Is(err, domain.ErrImportInvalidFile):
			response.BadRequest(c, "invalid import file")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, batch)
}

func (h *SalesOrderHandler) ListImportBatches(c *gin.Context) {
	page := parseIntOrDefault(c.Query("page"), 1)
	pageSize := parseIntOrDefault(c.Query("page_size"), 20)

	list, total, err := h.usecase.ListImports(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, list, total, page, pageSize)
}

func (h *SalesOrderHandler) GetImportBatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	item, err := h.usecase.GetImport(id)
	if err != nil {
		if errors.Is(err, domain.ErrImportNotFound) {
			response.NotFound(c, "import batch not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, item)
}

func (h *SalesOrderHandler) ListImportBatchErrors(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if _, err := h.usecase.GetImport(id); err != nil {
		if errors.Is(err, domain.ErrImportNotFound) {
			response.NotFound(c, "import batch not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	rows, err := h.usecase.ListImportErrors(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, rows)
}

func parseOperatorIDFromFormOrContext(c *gin.Context) *uint64 {
	raw := c.PostForm("operator_id")
	if raw != "" {
		if parsed, err := strconv.ParseUint(raw, 10, 64); err == nil {
			return &parsed
		}
	}

	if userID, ok := c.Get(auth.UserIDKey); ok {
		if val, castOK := userID.(uint64); castOK {
			return &val
		}
	}
	return nil
}
