package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/inventory/usecase"

	"github.com/gin-gonic/gin"
)

type stubWarehouseHandlerRepo struct {
	referenceCount int64
}

func (s *stubWarehouseHandlerRepo) List(params *domain.WarehouseListParams) ([]*domain.Warehouse, int64, error) {
	return nil, 0, nil
}
func (s *stubWarehouseHandlerRepo) GetByID(id uint64) (*domain.Warehouse, error) {
	return &domain.Warehouse{ID: id, Code: "WH001", Name: "测试仓库"}, nil
}
func (s *stubWarehouseHandlerRepo) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	return nil
}
func (s *stubWarehouseHandlerRepo) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	return nil
}
func (s *stubWarehouseHandlerRepo) Delete(ctx context.Context, id uint64) error { return nil }
func (s *stubWarehouseHandlerRepo) GetActiveWarehouses() ([]*domain.Warehouse, error) {
	return nil, nil
}
func (s *stubWarehouseHandlerRepo) CountReferences(id uint64) (int64, error) {
	return s.referenceCount, nil
}

func TestDeleteWarehouseReturnsBadRequestWhenReferenced(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewWarehouseHandler(usecase.NewWarehouseUsecase(&stubWarehouseHandlerRepo{referenceCount: 1}))
	router := gin.New()
	router.DELETE("/api/v1/inventory/warehouses/:id", handler.DeleteWarehouse)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/inventory/warehouses/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
