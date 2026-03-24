package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/procurement/domain"
	"am-erp-go/internal/module/procurement/usecase"
	productdomain "am-erp-go/internal/module/product/domain"

	"github.com/gin-gonic/gin"
)

type stubPurchaseOrderRepo struct {
	createdOrders  []*domain.PurchaseOrder
	lastListParams *domain.PurchaseOrderListParams
}

func (s *stubPurchaseOrderRepo) List(params *domain.PurchaseOrderListParams) ([]domain.PurchaseOrder, int64, error) {
	s.lastListParams = params
	return nil, 0, nil
}

func (s *stubPurchaseOrderRepo) GetByID(id uint64) (*domain.PurchaseOrder, error) {
	for _, order := range s.createdOrders {
		if order != nil && order.ID == id {
			cp := *order
			if len(order.Items) > 0 {
				cp.Items = make([]domain.PurchaseOrderItem, len(order.Items))
				copy(cp.Items, order.Items)
			}
			return &cp, nil
		}
	}
	return nil, nil
}

func (s *stubPurchaseOrderRepo) Create(order *domain.PurchaseOrder) error {
	if order != nil && order.ID == 0 {
		order.ID = uint64(len(s.createdOrders) + 1)
	}
	if order != nil {
		cp := *order
		if len(order.Items) > 0 {
			cp.Items = make([]domain.PurchaseOrderItem, len(order.Items))
			copy(cp.Items, order.Items)
		}
		s.createdOrders = append(s.createdOrders, &cp)
	}
	return nil
}

func (s *stubPurchaseOrderRepo) Update(_ *domain.PurchaseOrder) error {
	return nil
}

func (s *stubPurchaseOrderRepo) UpdateProgress(order *domain.PurchaseOrder) error {
	if order == nil {
		return nil
	}
	for index, existing := range s.createdOrders {
		if existing != nil && existing.ID == order.ID {
			cp := *existing
			cp.Status = order.Status
			cp.WarehouseID = order.WarehouseID
			cp.OrderedAt = order.OrderedAt
			cp.OrderedBy = order.OrderedBy
			cp.ShippedAt = order.ShippedAt
			cp.ShippedBy = order.ShippedBy
			cp.ReceivedAt = order.ReceivedAt
			cp.ReceivedBy = order.ReceivedBy
			cp.InspectedAt = order.InspectedAt
			cp.InspectedBy = order.InspectedBy
			cp.ClosedAt = order.ClosedAt
			cp.CompletedBy = order.CompletedBy
			cp.IsForceCompleted = order.IsForceCompleted
			cp.ForceCompletedAt = order.ForceCompletedAt
			cp.ForceCompletedBy = order.ForceCompletedBy
			cp.ForceCompleteReason = order.ForceCompleteReason
			if order.Items != nil {
				cp.Items = make([]domain.PurchaseOrderItem, len(order.Items))
				copy(cp.Items, order.Items)
			}
			s.createdOrders[index] = &cp
			return nil
		}
	}
	return nil
}

func (s *stubPurchaseOrderRepo) Delete(_ uint64) error {
	return nil
}

type stubProductLookup struct {
	items map[uint64]productdomain.Product
}

func (s *stubProductLookup) ListByIDs(ids []uint64) ([]productdomain.Product, error) {
	result := make([]productdomain.Product, 0, len(ids))
	for _, id := range ids {
		if item, ok := s.items[id]; ok {
			result = append(result, item)
		}
	}
	return result, nil
}

type stubComboProvider struct {
	items map[uint64][]productdomain.ProductComboItem
}

func (s *stubComboProvider) GetItemsByComboID(comboID uint64) ([]productdomain.ProductComboItem, error) {
	return s.items[comboID], nil
}

func floatPtr(v float64) *float64 {
	return &v
}

func TestCreatePurchaseOrderRejectsMissingItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubPurchaseOrderRepo{}
	uc := usecase.NewPurchaseOrderUsecase(repo, nil, nil, nil, nil)
	handler := NewPurchaseOrderHandler(uc)

	router := gin.New()
	router.POST("/api/procurement/purchase-orders", handler.CreatePurchaseOrder)

	payload := []byte(`{"supplier_id":1,"marketplace":"US","currency":"USD","items":[]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/procurement/purchase-orders", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreatePurchaseOrderRejectsComboWithoutComponents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	comboID := uint64(20)
	mainID := uint64(300)

	productLookup := &stubProductLookup{
		items: map[uint64]productdomain.Product{
			mainID: {ID: mainID, ComboID: &comboID, IsComboMain: 1},
		},
	}

	comboProvider := &stubComboProvider{
		items: map[uint64][]productdomain.ProductComboItem{
			comboID: {
				{ComboID: comboID, MainProductID: mainID, ProductID: mainID, QtyRatio: 1},
			},
		},
	}

	repo := &stubPurchaseOrderRepo{}
	uc := usecase.NewPurchaseOrderUsecase(repo, productLookup, comboProvider, nil, nil)
	handler := NewPurchaseOrderHandler(uc)

	router := gin.New()
	router.POST("/api/procurement/purchase-orders", handler.CreatePurchaseOrder)

	payload := []byte(`{"supplier_id":1,"marketplace":"US","currency":"USD","items":[{"product_id":300,"qty_ordered":1,"unit_cost":10}]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/procurement/purchase-orders", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreatePurchaseOrderRejectsZeroUnitCost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubPurchaseOrderRepo{}
	uc := usecase.NewPurchaseOrderUsecase(repo, nil, nil, nil, nil)
	handler := NewPurchaseOrderHandler(uc)

	router := gin.New()
	router.POST("/api/procurement/purchase-orders", handler.CreatePurchaseOrder)

	payload := []byte(`{"supplier_id":1,"marketplace":"US","currency":"USD","items":[{"product_id":300,"qty_ordered":1,"unit_cost":0}]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/procurement/purchase-orders", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreatePurchaseOrderRejectsMissingSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubPurchaseOrderRepo{}
	uc := usecase.NewPurchaseOrderUsecase(repo, nil, nil, nil, nil)
	handler := NewPurchaseOrderHandler(uc)

	router := gin.New()
	router.POST("/api/procurement/purchase-orders", handler.CreatePurchaseOrder)

	payload := []byte(`{"marketplace":"US","currency":"USD","items":[{"product_id":300,"qty_ordered":1,"unit_cost":10}]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/procurement/purchase-orders", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreatePurchaseOrderBatchSplitsComboChildrenBySupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)

	comboID := uint64(40)
	mainID := uint64(500)
	childAID := uint64(501)
	childBID := uint64(502)
	supplierAID := uint64(8)
	supplierBID := uint64(9)

	productLookup := &stubProductLookup{
		items: map[uint64]productdomain.Product{
			mainID:   {ID: mainID, ComboID: &comboID, IsComboMain: 1},
			childAID: {ID: childAID, SupplierID: &supplierAID, UnitCost: floatPtr(10)},
			childBID: {ID: childBID, SupplierID: &supplierBID, UnitCost: floatPtr(20)},
		},
	}
	comboProvider := &stubComboProvider{
		items: map[uint64][]productdomain.ProductComboItem{
			comboID: {
				{ComboID: comboID, MainProductID: mainID, ProductID: mainID, QtyRatio: 1},
				{ComboID: comboID, MainProductID: mainID, ProductID: childAID, QtyRatio: 1},
				{ComboID: comboID, MainProductID: mainID, ProductID: childBID, QtyRatio: 2},
			},
		},
	}

	repo := &stubPurchaseOrderRepo{}
	uc := usecase.NewPurchaseOrderUsecase(repo, productLookup, comboProvider, nil, nil)
	handler := NewPurchaseOrderHandler(uc)

	router := gin.New()
	router.POST("/api/procurement/purchase-orders/batch", handler.CreatePurchaseOrderBatch)

	payload := []byte(`{"orders":[{"marketplace":"US","currency":"USD","items":[{"product_id":500,"qty_ordered":1,"unit_cost":0}]}]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/procurement/purchase-orders/batch", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	if len(repo.createdOrders) != 2 {
		t.Fatalf("expected 2 created split orders, got %d", len(repo.createdOrders))
	}
	if repo.createdOrders[0].BatchNo == "" || repo.createdOrders[0].BatchNo != repo.createdOrders[1].BatchNo {
		t.Fatalf("expected shared batch no, got %+v", repo.createdOrders)
	}
}

func TestListPurchaseOrdersDefaultsToPageSize10AndParsesFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubPurchaseOrderRepo{}
	uc := usecase.NewPurchaseOrderUsecase(repo, nil, nil, nil, nil)
	handler := NewPurchaseOrderHandler(uc)

	router := gin.New()
	router.GET("/api/procurement/purchase-orders", handler.ListPurchaseOrders)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/procurement/purchase-orders?keyword=PO2026&supplier_id=8&marketplace=US&status=ORDERED",
		nil,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	if repo.lastListParams == nil {
		t.Fatalf("expected list params captured")
	}
	if repo.lastListParams.PageSize != 10 {
		t.Fatalf("expected default page size 10, got %d", repo.lastListParams.PageSize)
	}
	if repo.lastListParams.Keyword != "PO2026" {
		t.Fatalf("expected keyword parsed, got %q", repo.lastListParams.Keyword)
	}
	if repo.lastListParams.Marketplace != "US" {
		t.Fatalf("expected marketplace US, got %q", repo.lastListParams.Marketplace)
	}
	if repo.lastListParams.Status != domain.PurchaseOrderStatusOrdered {
		t.Fatalf("expected status ORDERED, got %s", repo.lastListParams.Status)
	}
	if repo.lastListParams.SupplierID == nil || *repo.lastListParams.SupplierID != 8 {
		t.Fatalf("expected supplier_id 8, got %+v", repo.lastListParams.SupplierID)
	}
}

func TestClosePurchaseOrderReturnsBadRequestWhenPendingInspectionRemains(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubPurchaseOrderRepo{
		createdOrders: []*domain.PurchaseOrder{
			{
				ID:       12,
				PoNumber: "PO-TEST-12",
				Status:   domain.PurchaseOrderStatusReceived,
				Items: []domain.PurchaseOrderItem{
					{
						ID:                1201,
						ProductID:         19,
						QtyOrdered:        2,
						QtyReceived:       2,
						QtyInspectionPass: 1,
					},
				},
			},
		},
	}
	uc := usecase.NewPurchaseOrderUsecase(repo, nil, nil, nil, nil)
	handler := NewPurchaseOrderHandler(uc)

	router := gin.New()
	router.POST("/api/procurement/purchase-orders/:id/close", handler.ClosePurchaseOrder)

	req := httptest.NewRequest(http.MethodPost, "/api/procurement/purchase-orders/12/close", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
}
