package usecase

import (
	"errors"
	"strconv"
	"strings"

	"am-erp-go/internal/module/supplier/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	ErrQuoteExists            = errors.New("quote already exists")
	ErrDefaultSupplierQuote   = errors.New("default supplier quote cannot be deleted")
	ErrProductRepoUnavailable = errors.New("product repository is required")
	ErrQuotePriceInvalid      = errors.New("quote price must be greater than zero")
)

type ProductSupplierRepository interface {
	GetDefaultSupplierID(productID uint64) (uint64, error)
	UpdateDefaultSupplierID(productID, supplierID uint64) error
}

type AuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

type QuoteUsecase struct {
	quoteRepo   domain.QuoteRepository
	productRepo ProductSupplierRepository
	auditLogger AuditLogger
}

func NewQuoteUsecase(
	quoteRepo domain.QuoteRepository,
	productRepo ProductSupplierRepository,
	auditLogger AuditLogger,
) *QuoteUsecase {
	return &QuoteUsecase{
		quoteRepo:   quoteRepo,
		productRepo: productRepo,
		auditLogger: auditLogger,
	}
}

func (uc *QuoteUsecase) ListProductQuotes(params *domain.QuoteListParams) ([]domain.ProductQuoteRow, int64, error) {
	if params == nil {
		params = &domain.QuoteListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return uc.quoteRepo.ListProductsWithQuotes(params)
}

func (uc *QuoteUsecase) CreateQuote(c *gin.Context, quote *domain.ProductSupplierQuote) (*domain.ProductSupplierQuote, error) {
	if quote == nil {
		return nil, errors.New("invalid quote")
	}
	if quote.Price <= 0 {
		return nil, ErrQuotePriceInvalid
	}

	existing, err := uc.quoteRepo.GetByProductSupplier(quote.ProductID, quote.SupplierID)
	if err == nil && existing != nil {
		return nil, ErrQuoteExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	quote.Status = defaultQuoteStatus(quote.Status)
	if err := uc.quoteRepo.Create(quote); err != nil {
		return nil, err
	}

	uc.recordAudit(c, "CREATE", "ProductSupplierQuote", strconv.FormatUint(quote.ID, 10), nil, quote)
	return quote, nil
}

func (uc *QuoteUsecase) UpdateQuote(c *gin.Context, quote *domain.ProductSupplierQuote) (*domain.ProductSupplierQuote, error) {
	if quote == nil {
		return nil, errors.New("invalid quote")
	}
	if quote.Price <= 0 {
		return nil, ErrQuotePriceInvalid
	}

	existing, err := uc.quoteRepo.GetByProductSupplier(quote.ProductID, quote.SupplierID)
	if err != nil {
		return nil, err
	}
	before := *existing

	existing.Price = quote.Price
	existing.Currency = quote.Currency
	existing.QtyMOQ = quote.QtyMOQ
	existing.LeadTimeDays = quote.LeadTimeDays
	existing.Remark = quote.Remark
	if strings.TrimSpace(quote.Status) != "" {
		existing.Status = quote.Status
	}

	if err := uc.quoteRepo.Update(existing); err != nil {
		return nil, err
	}

	uc.recordAudit(c, "UPDATE", "ProductSupplierQuote", strconv.FormatUint(existing.ID, 10), before, existing)
	return existing, nil
}

func (uc *QuoteUsecase) DeleteQuote(c *gin.Context, productID, supplierID uint64) error {
	existing, err := uc.quoteRepo.GetByProductSupplier(productID, supplierID)
	if err != nil {
		return err
	}

	if uc.productRepo == nil {
		return ErrProductRepoUnavailable
	}
	defaultSupplierID, err := uc.productRepo.GetDefaultSupplierID(productID)
	if err != nil {
		return err
	}
	if defaultSupplierID == supplierID {
		return ErrDefaultSupplierQuote
	}

	if err := uc.quoteRepo.Delete(productID, supplierID); err != nil {
		return err
	}

	uc.recordAudit(c, "DELETE", "ProductSupplierQuote", strconv.FormatUint(existing.ID, 10), existing, nil)
	return nil
}

func (uc *QuoteUsecase) SetDefaultSupplier(c *gin.Context, productID, supplierID uint64) error {
	if uc.productRepo == nil {
		return ErrProductRepoUnavailable
	}

	if _, err := uc.quoteRepo.GetByProductSupplier(productID, supplierID); err != nil {
		return err
	}

	beforeID, err := uc.productRepo.GetDefaultSupplierID(productID)
	if err != nil {
		return err
	}
	if beforeID == supplierID {
		return nil
	}

	if err := uc.productRepo.UpdateDefaultSupplierID(productID, supplierID); err != nil {
		return err
	}

	before := map[string]any{"default_supplier_id": beforeID}
	after := map[string]any{"default_supplier_id": supplierID}
	uc.recordAudit(c, "UPDATE", "Product", strconv.FormatUint(productID, 10), before, after)
	return nil
}

func (uc *QuoteUsecase) recordAudit(c *gin.Context, action, entityType, entityID string, before, after any) {
	if uc.auditLogger == nil || c == nil {
		return
	}
	_ = uc.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Product",
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Before:     before,
		After:      after,
	})
}

func defaultQuoteStatus(status string) string {
	if strings.TrimSpace(status) == "" {
		return "ACTIVE"
	}
	return status
}
