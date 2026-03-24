package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"am-erp-go/internal/module/packaging/domain"
	"am-erp-go/internal/module/packaging/usecase"

	"github.com/gin-gonic/gin"
)

func TestCreateItemReturnsBadRequestWhenSupplierMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewPackagingHandler(usecase.NewPackagingUsecase(&stubPackagingItemValidationRepo{}, &stubPackagingLedgerValidationRepo{}))

	router := gin.New()
	router.POST("/api/v1/packaging/items", func(c *gin.Context) {
		c.Set("userID", uint64(1))
		handler.CreateItem(c)
	})

	body, _ := json.Marshal(map[string]any{
		"item_code": "PKG-1",
		"item_name": "Package 1",
		"category":  "BOX",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/packaging/items", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestUpdateItemReturnsBadRequestWhenSupplierMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewPackagingHandler(usecase.NewPackagingUsecase(&stubPackagingItemValidationRepo{}, &stubPackagingLedgerValidationRepo{}))

	router := gin.New()
	router.PUT("/api/v1/packaging/items/:id", handler.UpdateItem)

	body, _ := json.Marshal(map[string]any{
		"item_code": "PKG-1",
		"item_name": "Package 1",
		"category":  "BOX",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/packaging/items/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCreateItemReturnsBadRequestWhenCodeInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewPackagingHandler(usecase.NewPackagingUsecase(&stubPackagingItemValidationRepo{}, &stubPackagingLedgerValidationRepo{}))

	router := gin.New()
	router.POST("/api/v1/packaging/items", func(c *gin.Context) {
		c.Set("userID", uint64(1))
		handler.CreateItem(c)
	})

	body, _ := json.Marshal(map[string]any{
		"item_code":   "包材@001",
		"item_name":   "Package 1",
		"category":    "BOX",
		"supplier_id": 1,
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/packaging/items", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

type stubPackagingItemValidationRepo struct{}

func (s *stubPackagingItemValidationRepo) List(params *domain.PackagingItemListParams) ([]domain.PackagingItem, int64, error) {
	return nil, 0, nil
}
func (s *stubPackagingItemValidationRepo) GetByID(id uint64) (*domain.PackagingItem, error) {
	return nil, nil
}
func (s *stubPackagingItemValidationRepo) Create(item *domain.PackagingItem) error { return nil }
func (s *stubPackagingItemValidationRepo) Update(item *domain.PackagingItem) error { return nil }
func (s *stubPackagingItemValidationRepo) Delete(id uint64) error                  { return nil }
func (s *stubPackagingItemValidationRepo) CountReferences(id uint64) (int64, error) {
	if id == 5 {
		return 2, nil
	}
	return 0, nil
}
func (s *stubPackagingItemValidationRepo) GetLowStockItems() ([]domain.PackagingItem, error) {
	return nil, nil
}
func (s *stubPackagingItemValidationRepo) UpdateQuantity(id uint64, quantity int64) error {
	return nil
}

type stubPackagingLedgerValidationRepo struct{}

func (s *stubPackagingLedgerValidationRepo) List(params *domain.PackagingLedgerListParams) ([]domain.PackagingLedger, int64, error) {
	return nil, 0, nil
}
func (s *stubPackagingLedgerValidationRepo) GetByID(id uint64) (*domain.PackagingLedger, error) {
	return nil, nil
}
func (s *stubPackagingLedgerValidationRepo) Create(ledger *domain.PackagingLedger) error { return nil }
func (s *stubPackagingLedgerValidationRepo) GetUsageSummary(dateFrom, dateTo *time.Time) ([]domain.UsageSummaryItem, error) {
	return nil, nil
}

func TestDeleteItemReturnsBadRequestWhenReferenced(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewPackagingHandler(usecase.NewPackagingUsecase(&stubPackagingItemValidationRepo{}, &stubPackagingLedgerValidationRepo{}))

	router := gin.New()
	router.DELETE("/api/v1/packaging/items/:id", handler.DeleteItem)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/packaging/items/5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
