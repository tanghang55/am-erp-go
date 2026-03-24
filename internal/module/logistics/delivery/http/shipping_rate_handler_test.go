package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/logistics/domain"
	"am-erp-go/internal/module/logistics/usecase"

	"github.com/gin-gonic/gin"
)

type stubShippingRateHandlerRepo struct{}

func (s *stubShippingRateHandlerRepo) Create(rate *domain.ShippingRate) error { return nil }
func (s *stubShippingRateHandlerRepo) Update(rate *domain.ShippingRate) error { return nil }
func (s *stubShippingRateHandlerRepo) Delete(id uint64) error                 { return nil }
func (s *stubShippingRateHandlerRepo) GetByID(id uint64) (*domain.ShippingRate, error) {
	return &domain.ShippingRate{ID: id}, nil
}
func (s *stubShippingRateHandlerRepo) List(params *domain.ShippingRateListParams) ([]*domain.ShippingRate, int64, error) {
	return nil, 0, nil
}
func (s *stubShippingRateHandlerRepo) QueryLatestRate(params *domain.QueryLatestRateParams) (*domain.ShippingRate, error) {
	return nil, nil
}
func (s *stubShippingRateHandlerRepo) CountReferences(id uint64) (int64, error) {
	if id == 5 {
		return 2, nil
	}
	return 0, nil
}

func TestDeleteShippingRateReturnsBadRequestWhenReferenced(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewShippingRateHandler(usecase.NewShippingRateUsecase(&stubShippingRateHandlerRepo{}, &stubProviderHandlerRepo{}))

	router := gin.New()
	router.DELETE("/api/v1/logistics/shipping-rates/:id", handler.DeleteShippingRate)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/logistics/shipping-rates/5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
