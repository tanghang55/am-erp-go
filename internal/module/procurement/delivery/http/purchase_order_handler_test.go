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

type stubPurchaseOrderRepo struct{}

func (s *stubPurchaseOrderRepo) List(_ *domain.PurchaseOrderListParams) ([]domain.PurchaseOrder, int64, error) {
	return nil, 0, nil
}

func (s *stubPurchaseOrderRepo) GetByID(_ uint64) (*domain.PurchaseOrder, error) {
	return nil, nil
}

func (s *stubPurchaseOrderRepo) Create(_ *domain.PurchaseOrder) error {
	return nil
}

func (s *stubPurchaseOrderRepo) Update(_ *domain.PurchaseOrder) error {
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

	payload := []byte(`{"supplier_id":1,"marketplace":"US","currency":"USD","items":[{"sku_id":300,"qty_ordered":1,"unit_cost":10}]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/procurement/purchase-orders", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
