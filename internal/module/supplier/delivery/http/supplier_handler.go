package http

import (
	"strconv"
	"strings"

	"am-erp-go/internal/infrastructure/response"
	"am-erp-go/internal/module/supplier/domain"
	"am-erp-go/internal/module/supplier/usecase"

	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	supplierUsecase *usecase.SupplierUsecase
}

func NewSupplierHandler(supplierUsecase *usecase.SupplierUsecase) *SupplierHandler {
	return &SupplierHandler{supplierUsecase: supplierUsecase}
}

// ==================== Supplier ====================

type supplierUpsertRequest struct {
	SupplierCode string   `json:"supplier_code"`
	Name         string   `json:"name"`
	Status       string   `json:"status"`
	Remark       string   `json:"remark"`
	Types        []string `json:"types"`
}

type supplierContactRequest struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Position  string `json:"position"`
	IsPrimary uint8  `json:"is_primary"`
}

type supplierAccountRequest struct {
	ID           uint64 `json:"id"`
	BankName     string `json:"bank_name"`
	BankAccount  string `json:"bank_account"`
	Currency     string `json:"currency"`
	TaxNo        string `json:"tax_no"`
	PaymentTerms string `json:"payment_terms"`
}

type supplierTagRequest struct {
	ID  uint64 `json:"id"`
	Tag string `json:"tag"`
}

// ListSuppliers 获取供应商列表
func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
	params := &domain.SupplierListParams{
		Page:     parseIntOrDefault(c.Query("page"), 1),
		PageSize: parseIntOrDefault(c.Query("page_size"), 20),
		Keyword:  c.Query("keyword"),
		Status:   c.Query("status"),
	}
	if rawType := strings.TrimSpace(c.Query("type")); rawType != "" {
		params.Types = splitComma(rawType)
	}

	suppliers, total, err := h.supplierUsecase.ListSuppliers(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPage(c, suppliers, total, params.Page, params.PageSize)
}

// GetSupplier 获取供应商详情
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	supplier, err := h.supplierUsecase.GetSupplier(id)
	if err != nil {
		response.NotFound(c, "supplier not found")
		return
	}

	response.Success(c, supplier)
}

// CreateSupplier 创建供应商
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req supplierUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	supplier := domain.Supplier{
		SupplierCode: req.SupplierCode,
		Name:         req.Name,
		Status:       defaultStatus(req.Status),
		Remark:       req.Remark,
	}

	created, err := h.supplierUsecase.CreateSupplier(&supplier, req.Types)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, created)
}

// UpdateSupplier 更新供应商
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	supplier := domain.Supplier{
		ID:           id,
		SupplierCode: req.SupplierCode,
		Name:         req.Name,
		Status:       defaultStatus(req.Status),
		Remark:       req.Remark,
	}

	updated, err := h.supplierUsecase.UpdateSupplier(&supplier, req.Types)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, updated)
}

// DeleteSupplier 删除供应商
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.supplierUsecase.DeleteSupplier(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateSupplierContact 创建供应商联系人
func (h *SupplierHandler) CreateSupplierContact(c *gin.Context) {
	supplierID, err := parseUintParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	contact := domain.SupplierContact{
		Name:      req.Name,
		Phone:     req.Phone,
		Email:     req.Email,
		Position:  req.Position,
		IsPrimary: req.IsPrimary,
	}

	created, err := h.supplierUsecase.CreateSupplierContact(supplierID, &contact)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, created)
}

// UpdateSupplierContact 更新供应商联系人
func (h *SupplierHandler) UpdateSupplierContact(c *gin.Context) {
	supplierID, err := parseUintParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID == 0 {
		response.BadRequest(c, "invalid contact id")
		return
	}

	contact := domain.SupplierContact{
		ID:        req.ID,
		Name:      req.Name,
		Phone:     req.Phone,
		Email:     req.Email,
		Position:  req.Position,
		IsPrimary: req.IsPrimary,
	}

	updated, err := h.supplierUsecase.UpdateSupplierContact(supplierID, &contact)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, updated)
}

// DeleteSupplierContact 删除供应商联系人
func (h *SupplierHandler) DeleteSupplierContact(c *gin.Context) {
	supplierID, err := parseUintParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID == 0 {
		response.BadRequest(c, "invalid contact id")
		return
	}

	if err := h.supplierUsecase.DeleteSupplierContact(supplierID, req.ID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateSupplierAccount 创建供应商账户
func (h *SupplierHandler) CreateSupplierAccount(c *gin.Context) {
	supplierID, err := parseUintParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	account := domain.SupplierAccount{
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		Currency:     req.Currency,
		TaxNo:        req.TaxNo,
		PaymentTerms: req.PaymentTerms,
	}

	created, err := h.supplierUsecase.CreateSupplierAccount(supplierID, &account)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, created)
}

// UpdateSupplierAccount 更新供应商账户
func (h *SupplierHandler) UpdateSupplierAccount(c *gin.Context) {
	supplierID, err := parseUintParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID == 0 {
		response.BadRequest(c, "invalid account id")
		return
	}

	account := domain.SupplierAccount{
		ID:           req.ID,
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		Currency:     req.Currency,
		TaxNo:        req.TaxNo,
		PaymentTerms: req.PaymentTerms,
	}

	updated, err := h.supplierUsecase.UpdateSupplierAccount(supplierID, &account)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, updated)
}

// DeleteSupplierAccount 删除供应商账户
func (h *SupplierHandler) DeleteSupplierAccount(c *gin.Context) {
	supplierID, err := parseUintParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID == 0 {
		response.BadRequest(c, "invalid account id")
		return
	}

	if err := h.supplierUsecase.DeleteSupplierAccount(supplierID, req.ID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateSupplierTag 创建供应商标签
func (h *SupplierHandler) CreateSupplierTag(c *gin.Context) {
	supplierID, err := parseUintParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tag := domain.SupplierTag{Tag: req.Tag}
	created, err := h.supplierUsecase.CreateSupplierTag(supplierID, &tag)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, created)
}

// UpdateSupplierTag 更新供应商标签
func (h *SupplierHandler) UpdateSupplierTag(c *gin.Context) {
	supplierID, err := parseUintParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID == 0 {
		response.BadRequest(c, "invalid tag id")
		return
	}

	tag := domain.SupplierTag{ID: req.ID, Tag: req.Tag}
	updated, err := h.supplierUsecase.UpdateSupplierTag(supplierID, &tag)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, updated)
}

// DeleteSupplierTag 删除供应商标签
func (h *SupplierHandler) DeleteSupplierTag(c *gin.Context) {
	supplierID, err := parseUintParam(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req supplierTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.ID == 0 {
		response.BadRequest(c, "invalid tag id")
		return
	}

	if err := h.supplierUsecase.DeleteSupplierTag(supplierID, req.ID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ==================== Helper ====================

func parseIntOrDefault(s string, defaultVal int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return defaultVal
}

func parseUintParam(raw string) (uint64, error) {
	return strconv.ParseUint(raw, 10, 64)
}

func splitComma(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func defaultStatus(status string) string {
	if strings.TrimSpace(status) == "" {
		return "ACTIVE"
	}
	return status
}
