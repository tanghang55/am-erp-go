package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/product/domain"
	"am-erp-go/internal/module/product/usecase"

	"github.com/gin-gonic/gin"
)

type stubComboRepo struct {
	lastListParams *domain.ComboListParams
	itemsByComboID map[uint64][]domain.ProductComboItem
}

func (s *stubComboRepo) ListComboIDs(params *domain.ComboListParams) ([]uint64, int64, error) {
	s.lastListParams = params
	return []uint64{}, 0, nil
}

func (s *stubComboRepo) GetItemsByComboID(comboID uint64) ([]domain.ProductComboItem, error) {
	if s.itemsByComboID == nil {
		return []domain.ProductComboItem{}, nil
	}
	return s.itemsByComboID[comboID], nil
}

func (s *stubComboRepo) CreateCombo(mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) (uint64, error) {
	comboID := uint64(1)
	if s.itemsByComboID == nil {
		s.itemsByComboID = map[uint64][]domain.ProductComboItem{}
	}
	items := make([]domain.ProductComboItem, 0, len(productIDs)+1)
	items = append(items, domain.ProductComboItem{ComboID: comboID, MainProductID: mainProductID, ProductID: mainProductID})
	for _, productID := range productIDs {
		items = append(items, domain.ProductComboItem{ComboID: comboID, MainProductID: mainProductID, ProductID: productID, QtyRatio: qtyRatios[productID]})
	}
	s.itemsByComboID[comboID] = items
	return comboID, nil
}

func (s *stubComboRepo) ReplaceComboItems(comboID uint64, mainProductID uint64, productIDs []uint64, qtyRatios map[uint64]uint64) error {
	return nil
}

func (s *stubComboRepo) DeleteCombo(comboID uint64) error {
	delete(s.itemsByComboID, comboID)
	return nil
}

func TestListCombosParsesStatusesQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	comboRepo := &stubComboRepo{itemsByComboID: map[uint64][]domain.ProductComboItem{
		9: {
			{ComboID: 9, MainProductID: 1, ProductID: 1},
			{ComboID: 9, MainProductID: 1, ProductID: 2, QtyRatio: 3},
		},
	}}
	handler := NewComboHandler(usecase.NewProductComboUsecase(comboRepo, nil, nil, nil))

	router := gin.New()
	router.GET("/api/v1/product-combos", handler.ListCombos)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/product-combos?statuses=ON_SALE,REPLENISHING", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if comboRepo.lastListParams == nil {
		t.Fatal("expected list params to be captured")
	}
	if len(comboRepo.lastListParams.Statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %+v", comboRepo.lastListParams.Statuses)
	}
	if comboRepo.lastListParams.Statuses[0] != "ON_SALE" || comboRepo.lastListParams.Statuses[1] != "REPLENISHING" {
		t.Fatalf("unexpected statuses: %+v", comboRepo.lastListParams.Statuses)
	}
}

func TestDeleteComboWritesAudit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	comboRepo := &stubComboRepo{itemsByComboID: map[uint64][]domain.ProductComboItem{
		9: {
			{ComboID: 9, MainProductID: 1, ProductID: 1},
			{ComboID: 9, MainProductID: 1, ProductID: 2, QtyRatio: 3},
		},
	}}
	handler := NewComboHandler(usecase.NewProductComboUsecase(comboRepo, nil, nil, nil))
	auditLogger := &stubProductAuditLogger{}
	handler.BindAuditLogger(auditLogger)

	router := gin.New()
	router.DELETE("/api/v1/product-combos/:id", handler.DeleteCombo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/product-combos/9", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if len(auditLogger.payloads) != 1 {
		t.Fatalf("expected 1 audit payload, got %d", len(auditLogger.payloads))
	}
	if auditLogger.payloads[0].Action != "DELETE_PRODUCT_COMBO" {
		t.Fatalf("unexpected audit action: %+v", auditLogger.payloads[0])
	}
}

func TestCreateComboWritesAudit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	comboRepo := &stubComboRepo{}
	handler := NewComboHandler(usecase.NewProductComboUsecase(comboRepo, nil, nil, nil))
	auditLogger := &stubProductAuditLogger{}
	handler.BindAuditLogger(auditLogger)

	router := gin.New()
	router.POST("/api/v1/product-combos", handler.CreateCombo)

	payload := []byte(`{"main_product_id":1,"children":[{"product_id":2,"qty_ratio":3}]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/product-combos", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if len(auditLogger.payloads) != 1 {
		t.Fatalf("expected 1 audit payload, got %d", len(auditLogger.payloads))
	}
	if auditLogger.payloads[0].Action != "CREATE_PRODUCT_COMBO" {
		raw, _ := json.Marshal(auditLogger.payloads)
		t.Fatalf("unexpected audit action: %s", raw)
	}
}

func TestCreateComboReturnsBadRequestWhenMainProductMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	comboRepo := &stubComboRepo{}
	handler := NewComboHandler(usecase.NewProductComboUsecase(comboRepo, nil, nil, nil))

	router := gin.New()
	router.POST("/api/v1/product-combos", handler.CreateCombo)

	payload := []byte(`{"children":[{"product_id":2,"qty_ratio":3}]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/product-combos", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}
