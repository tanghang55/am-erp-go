package http

import (
	"errors"
	"strconv"
	"strings"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/supplier/domain"
	"am-erp-go/internal/module/supplier/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuoteHandler struct {
	quoteUsecase *usecase.QuoteUsecase
}

func NewQuoteHandler(quoteUsecase *usecase.QuoteUsecase) *QuoteHandler {
	return &QuoteHandler{quoteUsecase: quoteUsecase}
}

type quoteUpsertRequest struct {
	ProductID    uint64  `json:"product_id"`
	SupplierID   uint64  `json:"supplier_id"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	QtyMOQ       uint64  `json:"qty_moq"`
	LeadTimeDays uint64  `json:"lead_time_days"`
	Status       string  `json:"status"`
	Remark       string  `json:"remark"`
}

type quoteDeleteRequest struct {
	ProductID  uint64 `json:"product_id"`
	SupplierID uint64 `json:"supplier_id"`
}

type quoteDefaultRequest struct {
	ProductID  uint64 `json:"product_id"`
	SupplierID uint64 `json:"supplier_id"`
}

// ListProductQuotes 获取产品供应商报价列表
func (h *QuoteHandler) ListProductQuotes(c *gin.Context) {
	params := &domain.QuoteListParams{
		Page:        parseIntOrDefault(c.Query("page"), 1),
		PageSize:    parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:     c.Query("keyword"),
		Marketplace: c.Query("marketplace"),
	}
	if supplierID := c.Query("supplier_id"); supplierID != "" {
		if id, err := strconv.ParseUint(supplierID, 10, 64); err == nil {
			params.SupplierID = &id
		}
	}
	if productID := c.Query("product_id"); productID != "" {
		if id, err := strconv.ParseUint(productID, 10, 64); err == nil {
			params.ProductID = &id
		}
	}
	if productIDs := c.Query("product_ids"); productIDs != "" {
		params.ProductIDs = parseUint64List(productIDs)
	}

	rows, total, err := h.quoteUsecase.ListProductQuotes(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, rows, total, params.Page, params.PageSize)
}

func parseUint64List(raw string) []uint64 {
	parts := strings.Split(raw, ",")
	result := make([]uint64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseUint(part, 10, 64)
		if err != nil || id == 0 {
			continue
		}
		result = append(result, id)
	}
	return result
}

func (h *QuoteHandler) GetProductQuote(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Query("product_id"), 10, 64)
	if err != nil || productID == 0 {
		response.BadRequest(c, "invalid product id")
		return
	}
	supplierID, err := strconv.ParseUint(c.Query("supplier_id"), 10, 64)
	if err != nil || supplierID == 0 {
		response.BadRequest(c, "invalid supplier id")
		return
	}

	quote, err := h.quoteUsecase.GetProductQuote(productID, supplierID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, domain.ErrQuoteNotFound) {
			response.NotFound(c, "quote not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, quote)
}

// CreateQuote 创建报价
func (h *QuoteHandler) CreateQuote(c *gin.Context) {
	var req quoteUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ProductID == 0 || req.SupplierID == 0 {
		response.BadRequest(c, "invalid product or supplier id")
		return
	}

	quote := domain.ProductSupplierQuote{
		ProductID:    req.ProductID,
		SupplierID:   req.SupplierID,
		Price:        req.Price,
		Currency:     req.Currency,
		QtyMOQ:       req.QtyMOQ,
		LeadTimeDays: req.LeadTimeDays,
		Status:       req.Status,
		Remark:       req.Remark,
	}

	created, err := h.quoteUsecase.CreateQuote(c, &quote)
	if err != nil {
		if errors.Is(err, usecase.ErrQuoteExists) || errors.Is(err, usecase.ErrQuotePriceInvalid) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, created)
}

// UpdateQuote 更新报价
func (h *QuoteHandler) UpdateQuote(c *gin.Context) {
	var req quoteUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ProductID == 0 || req.SupplierID == 0 {
		response.BadRequest(c, "invalid product or supplier id")
		return
	}

	quote := domain.ProductSupplierQuote{
		ProductID:    req.ProductID,
		SupplierID:   req.SupplierID,
		Price:        req.Price,
		Currency:     req.Currency,
		QtyMOQ:       req.QtyMOQ,
		LeadTimeDays: req.LeadTimeDays,
		Status:       req.Status,
		Remark:       req.Remark,
	}

	updated, err := h.quoteUsecase.UpdateQuote(c, &quote)
	if err != nil {
		if errors.Is(err, usecase.ErrQuotePriceInvalid) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "quote not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, updated)
}

// DeleteQuote 删除报价
func (h *QuoteHandler) DeleteQuote(c *gin.Context) {
	var req quoteDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ProductID == 0 || req.SupplierID == 0 {
		response.BadRequest(c, "invalid product or supplier id")
		return
	}

	if err := h.quoteUsecase.DeleteQuote(c, req.ProductID, req.SupplierID); err != nil {
		if errors.Is(err, usecase.ErrDefaultSupplierQuote) {
			response.BadRequest(c, err.Error())
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "quote not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// SetDefaultSupplier 设置默认供应商
func (h *QuoteHandler) SetDefaultSupplier(c *gin.Context) {
	var req quoteDefaultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ProductID == 0 || req.SupplierID == 0 {
		response.BadRequest(c, "invalid product or supplier id")
		return
	}

	if err := h.quoteUsecase.SetDefaultSupplier(c, req.ProductID, req.SupplierID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "quote not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
