package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/logistics/domain"
	"am-erp-go/internal/module/logistics/usecase"

	"github.com/gin-gonic/gin"
)

func TestCreateProviderReturnsBadRequestWhenCodeInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewLogisticsProviderHandler(usecase.NewLogisticsProviderUsecase(&stubProviderHandlerRepo{}))
	router := gin.New()
	router.POST("/api/v1/logistics/providers", handler.CreateProvider)

	payload := []byte(`{"provider_code":"物流@001","provider_name":"Demo","provider_type":"COURIER","status":"ACTIVE"}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/logistics/providers", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDeleteProviderReturnsBadRequestWhenReferenced(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubProviderHandlerRepo{
		getByID:        &domain.LogisticsProvider{ID: 1, ProviderCode: "LP001", ProviderName: "Demo"},
		referenceCount: 2,
	}
	handler := NewLogisticsProviderHandler(usecase.NewLogisticsProviderUsecase(repo))
	router := gin.New()
	router.DELETE("/api/v1/logistics/providers/:id", handler.DeleteProvider)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/logistics/providers/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

type stubProviderHandlerRepo struct {
	getByID        *domain.LogisticsProvider
	referenceCount int64
}

func (s *stubProviderHandlerRepo) Create(provider *domain.LogisticsProvider) error { return nil }
func (s *stubProviderHandlerRepo) Update(provider *domain.LogisticsProvider) error { return nil }
func (s *stubProviderHandlerRepo) Delete(id uint64) error                          { return nil }
func (s *stubProviderHandlerRepo) GetByID(id uint64) (*domain.LogisticsProvider, error) {
	return s.getByID, nil
}
func (s *stubProviderHandlerRepo) GetByCode(code string) (*domain.LogisticsProvider, error) {
	return nil, nil
}
func (s *stubProviderHandlerRepo) List(params *domain.LogisticsProviderListParams) ([]*domain.LogisticsProvider, int64, error) {
	return nil, 0, nil
}
func (s *stubProviderHandlerRepo) CountReferences(id uint64) (int64, error) {
	return s.referenceCount, nil
}
