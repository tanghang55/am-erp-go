package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	productDomain "am-erp-go/internal/module/product/domain"
	shipmentDomain "am-erp-go/internal/module/shipment/domain"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"

	"github.com/gin-gonic/gin"
)

type stubShipmentHandlerRepo struct {
	getByID func(id uint64) (*shipmentDomain.Shipment, error)
	updated *shipmentDomain.Shipment
}

func (s *stubShipmentHandlerRepo) Create(shipment *shipmentDomain.Shipment) error {
	shipment.ID = 1
	return nil
}

func (s *stubShipmentHandlerRepo) Update(shipment *shipmentDomain.Shipment) error {
	s.updated = shipment
	return nil
}

func (s *stubShipmentHandlerRepo) GetByID(id uint64) (*shipmentDomain.Shipment, error) {
	if s.getByID != nil {
		return s.getByID(id)
	}
	return nil, nil
}

func (s *stubShipmentHandlerRepo) GetByShipmentNumber(shipmentNumber string) (*shipmentDomain.Shipment, error) {
	return nil, nil
}

func (s *stubShipmentHandlerRepo) List(params *shipmentDomain.ShipmentListParams) ([]*shipmentDomain.Shipment, int64, error) {
	return nil, 0, nil
}

func (s *stubShipmentHandlerRepo) Delete(id uint64) error {
	return nil
}

type stubShipmentHandlerItemRepo struct{}

func (s *stubShipmentHandlerItemRepo) Create(item *shipmentDomain.ShipmentItem) error {
	return nil
}

func (s *stubShipmentHandlerItemRepo) CreateBatch(items []shipmentDomain.ShipmentItem) error {
	return nil
}

func (s *stubShipmentHandlerItemRepo) UpdateBatch(items []shipmentDomain.ShipmentItem) error {
	return nil
}

func (s *stubShipmentHandlerItemRepo) GetByShipmentID(shipmentID uint64) ([]shipmentDomain.ShipmentItem, error) {
	return nil, nil
}

func (s *stubShipmentHandlerItemRepo) DeleteByShipmentID(shipmentID uint64) error {
	return nil
}

type stubShipmentHandlerInventoryService struct{}

func (s *stubShipmentHandlerInventoryService) CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error) {
	return nil, nil
}

func (s *stubShipmentHandlerInventoryService) GetProductBalance(productID, warehouseID uint64) (*inventoryDomain.InventoryBalance, error) {
	return &inventoryDomain.InventoryBalance{PendingShipment: 100, PendingShipmentReserved: 100}, nil
}

type stubShipmentHandlerInventoryServiceLowStock struct{}

func (s *stubShipmentHandlerInventoryServiceLowStock) CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error) {
	return nil, nil
}

func (s *stubShipmentHandlerInventoryServiceLowStock) GetProductBalance(productID, warehouseID uint64) (*inventoryDomain.InventoryBalance, error) {
	return &inventoryDomain.InventoryBalance{PendingShipment: 1, PendingShipmentReserved: 1}, nil
}

type stubShipmentHandlerProductRepo struct{}

func (s *stubShipmentHandlerProductRepo) ListByIDs(ids []uint64) ([]productDomain.Product, error) {
	return []productDomain.Product{
		{ID: 1001, SellerSku: "COMBO-MAIN", Status: productDomain.ProductStatusOffShelf},
	}, nil
}

type stubShipmentHandlerProductRepoActive struct{}

func (s *stubShipmentHandlerProductRepoActive) ListByIDs(ids []uint64) ([]productDomain.Product, error) {
	return []productDomain.Product{
		{ID: 1001, SellerSku: "SKU-1001", Status: productDomain.ProductStatusOnSale},
	}, nil
}

type stubShipmentHandlerWarehouseRepo struct{}

func (s *stubShipmentHandlerWarehouseRepo) GetByID(id uint64) (*inventoryDomain.Warehouse, error) {
	return &inventoryDomain.Warehouse{ID: id}, nil
}

func TestCreateShipmentReturnsBadRequestWhenProductInactive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uc := shipmentUsecase.NewShipmentUsecase(
		&stubShipmentHandlerRepo{},
		&stubShipmentHandlerItemRepo{},
		&stubShipmentHandlerInventoryService{},
		&stubShipmentHandlerProductRepo{},
		&stubShipmentHandlerWarehouseRepo{},
	)
	handler := NewShipmentHandler(uc)

	router := gin.New()
	router.POST("/api/v1/shipments", handler.CreateShipment)

	body, _ := json.Marshal(map[string]any{
		"warehouse_id": 1,
		"items": []map[string]any{
			{
				"product_id":       1001,
				"quantity_planned": 1,
			},
		},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shipments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCreateShipmentReturnsBadRequestWhenPendingShipmentInsufficient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uc := shipmentUsecase.NewShipmentUsecase(
		&stubShipmentHandlerRepo{},
		&stubShipmentHandlerItemRepo{},
		&stubShipmentHandlerInventoryServiceLowStock{},
		&stubShipmentHandlerProductRepoActive{},
		&stubShipmentHandlerWarehouseRepo{},
	)
	handler := NewShipmentHandler(uc)

	router := gin.New()
	router.POST("/api/v1/shipments", handler.CreateShipment)

	body, _ := json.Marshal(map[string]any{
		"warehouse_id": 1,
		"items": []map[string]any{
			{
				"product_id":       1001,
				"quantity_planned": 2,
			},
		},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shipments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestUpdateShipmentReturnsBadRequestWhenStatusInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubShipmentHandlerRepo{
		getByID: func(id uint64) (*shipmentDomain.Shipment, error) {
			return &shipmentDomain.Shipment{
				ID:          id,
				WarehouseID: 1,
				Status:      shipmentDomain.ShipmentStatusShipped,
			}, nil
		},
	}
	uc := shipmentUsecase.NewShipmentUsecase(
		repo,
		&stubShipmentHandlerItemRepo{},
		&stubShipmentHandlerInventoryService{},
		&stubShipmentHandlerProductRepoActive{},
		&stubShipmentHandlerWarehouseRepo{},
	)
	handler := NewShipmentHandler(uc)

	router := gin.New()
	router.PUT("/api/v1/shipments/:id", handler.UpdateShipment)

	body, _ := json.Marshal(map[string]any{
		"logistics_provider_id": 3,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/shipments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestUpdateShipmentBindsDestinationAndLogisticsFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubShipmentHandlerRepo{
		getByID: func(id uint64) (*shipmentDomain.Shipment, error) {
			return &shipmentDomain.Shipment{
				ID:          id,
				WarehouseID: 1,
				Status:      shipmentDomain.ShipmentStatusDraft,
			}, nil
		},
	}
	uc := shipmentUsecase.NewShipmentUsecase(
		repo,
		&stubShipmentHandlerItemRepo{},
		&stubShipmentHandlerInventoryService{},
		&stubShipmentHandlerProductRepoActive{},
		&stubShipmentHandlerWarehouseRepo{},
	)
	handler := NewShipmentHandler(uc)

	router := gin.New()
	router.PUT("/api/v1/shipments/:id", handler.UpdateShipment)

	body, _ := json.Marshal(map[string]any{
		"warehouse_id":             2,
		"destination_name":         "FBA US",
		"logistics_provider_id":    3,
		"shipping_rate_id":         5,
		"transport_mode":           "SEA",
		"expected_delivery_date":   "2026-03-30",
		"remark":                   "updated",
		"items": []map[string]any{
			{
				"product_id":       1001,
				"quantity_planned": 2,
			},
		},
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/shipments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if repo.updated == nil || repo.updated.LogisticsProviderID == nil || *repo.updated.LogisticsProviderID != 3 {
		t.Fatalf("expected logistics provider to update, got %+v", repo.updated)
	}
	if repo.updated.ShippingRateID == nil || *repo.updated.ShippingRateID != 5 {
		t.Fatalf("expected shipping rate to update, got %+v", repo.updated)
	}
}
