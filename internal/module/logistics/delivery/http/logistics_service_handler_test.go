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

func TestCreateLogisticsServiceReturnsBadRequestWhenCodeInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewLogisticsServiceHandler(usecase.NewLogisticsServiceUsecase(&stubLogisticsServiceHandlerRepo{}))
	router := gin.New()
	router.POST("/api/v1/logistics/logistics-services", handler.CreateLogisticsService)

	payload := []byte(`{"service_code":"慢船@001","service_name":"Demo","transport_mode":"SEA","status":"ACTIVE"}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/logistics/logistics-services", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDeleteLogisticsServiceReturnsBadRequestWhenReferenced(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubLogisticsServiceHandlerRepo{
		getByID:        &domain.LogisticsService{ID: 1, ServiceCode: "SEA001", ServiceName: "Demo"},
		referenceCount: 1,
	}
	handler := NewLogisticsServiceHandler(usecase.NewLogisticsServiceUsecase(repo))
	router := gin.New()
	router.DELETE("/api/v1/logistics/logistics-services/:id", handler.DeleteLogisticsService)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/logistics/logistics-services/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

type stubLogisticsServiceHandlerRepo struct {
	getByID        *domain.LogisticsService
	referenceCount int64
}

func (s *stubLogisticsServiceHandlerRepo) Create(service *domain.LogisticsService) error { return nil }
func (s *stubLogisticsServiceHandlerRepo) Update(service *domain.LogisticsService) error { return nil }
func (s *stubLogisticsServiceHandlerRepo) Delete(id uint64) error                        { return nil }
func (s *stubLogisticsServiceHandlerRepo) GetByID(id uint64) (*domain.LogisticsService, error) {
	return s.getByID, nil
}
func (s *stubLogisticsServiceHandlerRepo) GetByCode(code string) (*domain.LogisticsService, error) {
	return nil, nil
}
func (s *stubLogisticsServiceHandlerRepo) List(params *domain.LogisticsServiceListParams) ([]*domain.LogisticsService, int64, error) {
	return nil, 0, nil
}
func (s *stubLogisticsServiceHandlerRepo) GetActiveServices() ([]*domain.LogisticsService, error) {
	return nil, nil
}
func (s *stubLogisticsServiceHandlerRepo) GetServicesByTransportMode(transportMode domain.TransportMode) ([]*domain.LogisticsService, error) {
	return nil, nil
}
func (s *stubLogisticsServiceHandlerRepo) CountReferences(id uint64) (int64, error) {
	return s.referenceCount, nil
}
