package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/supplier/domain"
	"am-erp-go/internal/module/supplier/usecase"

	"github.com/gin-gonic/gin"
)

type stubQuoteRepo struct {
	listParams *domain.QuoteListParams
	listRows   []domain.ProductQuoteRow
	listTotal  int64
	listErr    error
	getQuote   *domain.ProductSupplierQuote
	getErr     error
	deleteErr  error
}

func (s *stubQuoteRepo) ListByProductIDs(_ []uint64) (map[uint64][]domain.ProductSupplierQuote, error) {
	return map[uint64][]domain.ProductSupplierQuote{}, nil
}

func (s *stubQuoteRepo) ListProductsWithQuotes(params *domain.QuoteListParams) ([]domain.ProductQuoteRow, int64, error) {
	s.listParams = params
	return s.listRows, s.listTotal, s.listErr
}

func (s *stubQuoteRepo) GetByProductSupplier(_, _ uint64) (*domain.ProductSupplierQuote, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.getQuote, nil
}

func (s *stubQuoteRepo) Create(_ *domain.ProductSupplierQuote) error { return nil }

func (s *stubQuoteRepo) Update(_ *domain.ProductSupplierQuote) error { return nil }

func (s *stubQuoteRepo) Delete(_, _ uint64) error { return s.deleteErr }

type stubProductSupplierRepo struct {
	defaultID uint64
	err       error
}

func (s *stubProductSupplierRepo) GetDefaultSupplierID(_ uint64) (uint64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return s.defaultID, nil
}

func (s *stubProductSupplierRepo) UpdateDefaultSupplierID(_, _ uint64) error { return s.err }

func TestListProductQuotesParsesQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	quoteRepo := &stubQuoteRepo{}
	uc := usecase.NewQuoteUsecase(quoteRepo, nil, nil)
	handler := NewQuoteHandler(uc)

	router := gin.New()
	router.GET("/api/v1/suppliers/product-quotes", handler.ListProductQuotes)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/suppliers/product-quotes?page=2&page_size=50&keyword=usb&marketplace=US&supplier_id=9", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if quoteRepo.listParams == nil || quoteRepo.listParams.Page != 2 || quoteRepo.listParams.PageSize != 50 {
		t.Fatalf("expected pagination parsed")
	}
	if quoteRepo.listParams.Keyword != "usb" || quoteRepo.listParams.Marketplace != "US" {
		t.Fatalf("expected filters parsed")
	}
	if quoteRepo.listParams.SupplierID == nil || *quoteRepo.listParams.SupplierID != 9 {
		t.Fatalf("expected supplier_id parsed")
	}
}

func TestDeleteQuoteBlocksDefaultSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)

	quoteRepo := &stubQuoteRepo{
		getQuote: &domain.ProductSupplierQuote{ID: 1, ProductID: 10, SupplierID: 5},
	}
	productRepo := &stubProductSupplierRepo{defaultID: 5}
	uc := usecase.NewQuoteUsecase(quoteRepo, productRepo, nil)
	handler := NewQuoteHandler(uc)

	router := gin.New()
	router.DELETE("/api/v1/suppliers/product-quotes", handler.DeleteQuote)

	payload := []byte(`{"product_id":10,"supplier_id":5}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/suppliers/product-quotes", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateQuoteRejectsZeroPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	quoteRepo := &stubQuoteRepo{}
	uc := usecase.NewQuoteUsecase(quoteRepo, nil, nil)
	handler := NewQuoteHandler(uc)

	router := gin.New()
	router.POST("/api/v1/suppliers/product-quotes", handler.CreateQuote)

	payload := []byte(`{"product_id":10,"supplier_id":5,"price":0,"currency":"USD","qty_moq":1,"lead_time_days":7}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers/product-quotes", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateQuoteRejectsZeroPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	quoteRepo := &stubQuoteRepo{
		getQuote: &domain.ProductSupplierQuote{ID: 1, ProductID: 10, SupplierID: 5},
	}
	uc := usecase.NewQuoteUsecase(quoteRepo, nil, nil)
	handler := NewQuoteHandler(uc)

	router := gin.New()
	router.PUT("/api/v1/suppliers/product-quotes", handler.UpdateQuote)

	payload := []byte(`{"product_id":10,"supplier_id":5,"price":0,"currency":"USD","qty_moq":1,"lead_time_days":7}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/suppliers/product-quotes", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
