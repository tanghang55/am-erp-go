package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/product/domain"

	"github.com/gin-gonic/gin"
)

type stubProductImageUsecase struct {
	listID   uint64
	saveID   uint64
	saveUrls []string
	list     []domain.ProductImage
	err      error
}

func (s *stubProductImageUsecase) ListProductImages(productID uint64) ([]domain.ProductImage, error) {
	s.listID = productID
	return s.list, s.err
}

func (s *stubProductImageUsecase) SaveProductImages(productID uint64, urls []string) ([]domain.ProductImage, error) {
	s.saveID = productID
	s.saveUrls = urls
	return s.list, s.err
}

func TestListProductImagesCallsUsecase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &stubProductImageUsecase{}
	handler := NewProductHandler(nil, stub)

	router := gin.New()
	router.GET("/api/v1/products/:id/images", handler.ListProductImages)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/9/images", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if stub.listID != 9 {
		t.Fatalf("expected product id 9, got %d", stub.listID)
	}
}

func TestSaveProductImagesParsesBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	stub := &stubProductImageUsecase{}
	handler := NewProductHandler(nil, stub)

	router := gin.New()
	router.PUT("/api/v1/products/:id/images/reorder", handler.SaveProductImages)

	payload := []byte(`{"image_urls":["a","b"]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/11/images/reorder", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if stub.saveID != 11 || len(stub.saveUrls) != 2 || stub.saveUrls[0] != "a" {
		t.Fatalf("unexpected save args: %+v", stub.saveUrls)
	}
}
