package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/product/domain"
	"am-erp-go/internal/module/product/usecase"
	supplierdomain "am-erp-go/internal/module/supplier/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type stubProductAuditLogger struct {
	payloads []systemUsecase.AuditLogPayload
}

func (s *stubProductAuditLogger) RecordFromContext(_ *gin.Context, payload systemUsecase.AuditLogPayload) error {
	s.payloads = append(s.payloads, payload)
	return nil
}

type stubProductRepo struct {
	lastListParams *domain.ProductListParams
	referenceCount map[uint64]int64
	product        domain.Product
}

func (s *stubProductRepo) List(params *domain.ProductListParams) ([]domain.Product, int64, error) {
	s.lastListParams = params
	return nil, 0, nil
}

func (s *stubProductRepo) GetByID(id uint64) (*domain.Product, error) {
	if s.product.ID == 0 {
		return nil, errors.New("not implemented")
	}
	copy := s.product
	return &copy, nil
}

func (s *stubProductRepo) ListByIDs(ids []uint64) ([]domain.Product, error) {
	return nil, nil
}

func (s *stubProductRepo) ListByParentID(parentID uint64) ([]domain.Product, error) {
	return nil, nil
}

func (s *stubProductRepo) CountReferencesByIDs(ids []uint64) (map[uint64]int64, error) {
	result := make(map[uint64]int64, len(ids))
	for _, id := range ids {
		result[id] = 0
		if s.referenceCount != nil {
			result[id] = s.referenceCount[id]
		}
	}
	return result, nil
}

func (s *stubProductRepo) CountByConfigReference(configType domain.ProductConfigType, configID uint64) (int64, error) {
	return 0, nil
}

func (s *stubProductRepo) CountByCategoryID(categoryID uint64) (int64, error) {
	return 0, nil
}

func (s *stubProductRepo) Create(product *domain.Product) error {
	if product.ID == 0 {
		product.ID = 1
	}
	s.product = *product
	return nil
}

func (s *stubProductRepo) Update(product *domain.Product) error {
	s.product = *product
	return nil
}

func (s *stubProductRepo) Delete(id uint64) error {
	return nil
}

func (s *stubProductRepo) UpdateParentIDBatch(productIDs []uint64, parentID *uint64) error {
	return nil
}

func (s *stubProductRepo) UpdateImageUrl(id uint64, imageUrl string) error {
	return nil
}

func (s *stubProductRepo) UpdateUnitCost(productID uint64, unitCost float64) error {
	s.product.UnitCost = &unitCost
	return nil
}

func (s *stubProductRepo) GetDefaultSupplierID(productID uint64) (uint64, error) {
	return 0, nil
}

func (s *stubProductRepo) UpdateDefaultSupplierID(productID, supplierID uint64) error {
	return nil
}

func (s *stubProductRepo) UpdateComboInfo(comboID uint64, mainProductID uint64, productIDs []uint64) error {
	return nil
}

func (s *stubProductRepo) ClearComboInfo(comboID uint64) error {
	return nil
}

type stubProductParentRepo struct {
	parent domain.ProductParent
}

type stubHandlerQuoteRepo struct {
	quotes map[string]supplierdomain.ProductSupplierQuote
}

func (s *stubHandlerQuoteRepo) GetByProductSupplier(productID, supplierID uint64) (*supplierdomain.ProductSupplierQuote, error) {
	if s.quotes != nil {
		if quote, ok := s.quotes[fmt.Sprintf("%d:%d", productID, supplierID)]; ok {
			copy := quote
			return &copy, nil
		}
	}
	return nil, supplierdomain.ErrQuoteNotFound
}

func (s *stubHandlerQuoteRepo) Create(quote *supplierdomain.ProductSupplierQuote) error {
	return nil
}

func (s *stubHandlerQuoteRepo) Update(quote *supplierdomain.ProductSupplierQuote) error {
	if s.quotes == nil {
		s.quotes = map[string]supplierdomain.ProductSupplierQuote{}
	}
	s.quotes[fmt.Sprintf("%d:%d", quote.ProductID, quote.SupplierID)] = *quote
	return nil
}

type stubHandlerDefaultsProvider struct{}

func (s *stubHandlerDefaultsProvider) GetDefaultBaseCurrency() string {
	return "USD"
}

func newTestProductUsecase(
	productRepo domain.ProductRepository,
	parentRepo domain.ProductParentRepository,
	configRepo domain.ProductConfigRepository,
	categoryRepo domain.ProductCategoryRepository,
) *usecase.ProductUsecase {
	uc := usecase.NewProductUsecase(productRepo, parentRepo, configRepo, categoryRepo, nil)
	quoteRepo := &stubHandlerQuoteRepo{}
	switch repo := productRepo.(type) {
	case *stubProductRepo:
		if repo.product.ID != 0 && repo.product.SupplierID != nil && *repo.product.SupplierID != 0 {
			quoteRepo.quotes = map[string]supplierdomain.ProductSupplierQuote{
				fmt.Sprintf("%d:%d", repo.product.ID, *repo.product.SupplierID): {
					ID:         1,
					ProductID:  repo.product.ID,
					SupplierID: *repo.product.SupplierID,
					Status:     supplierdomain.ProductSupplierQuoteStatusActive,
				},
			}
		}
	}
	uc.BindQuoteRepository(quoteRepo)
	uc.BindDefaultsProvider(&stubHandlerDefaultsProvider{})
	return uc
}

type stubProductConfigRepo struct {
	item    domain.ProductConfigItem
	getErr  error
	deleted uint64
}

func (s *stubProductConfigRepo) List(params *domain.ProductConfigListParams) ([]domain.ProductConfigItem, int64, error) {
	return nil, 0, nil
}

func (s *stubProductConfigRepo) GetByID(id uint64) (*domain.ProductConfigItem, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	copy := s.item
	return &copy, nil
}

func (s *stubProductConfigRepo) Create(item *domain.ProductConfigItem) error {
	return nil
}

func (s *stubProductConfigRepo) Update(item *domain.ProductConfigItem) error {
	s.item = *item
	return nil
}

func (s *stubProductConfigRepo) Delete(id uint64) error {
	s.deleted = id
	return nil
}

type stubProductCategoryRepo struct {
	item          domain.ProductCategory
	getErr        error
	childrenCount int64
	deleted       uint64
}

func (s *stubProductCategoryRepo) ListAll() ([]domain.ProductCategory, error) {
	return nil, nil
}

func (s *stubProductCategoryRepo) GetByID(id uint64) (*domain.ProductCategory, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	copy := s.item
	return &copy, nil
}

func (s *stubProductCategoryRepo) Create(item *domain.ProductCategory) error {
	return nil
}

func (s *stubProductCategoryRepo) Update(item *domain.ProductCategory) error {
	s.item = *item
	return nil
}

func (s *stubProductCategoryRepo) Delete(id uint64) error {
	s.deleted = id
	return nil
}

func (s *stubProductCategoryRepo) CountChildren(id uint64) (int64, error) {
	return s.childrenCount, nil
}

func (s *stubProductParentRepo) List(params *domain.ProductParentListParams) ([]domain.ProductParent, int64, error) {
	return nil, 0, nil
}

func (s *stubProductParentRepo) GetByID(id uint64) (*domain.ProductParent, error) {
	copy := s.parent
	return &copy, nil
}

func (s *stubProductParentRepo) Create(parent *domain.ProductParent) error {
	return nil
}

func (s *stubProductParentRepo) Update(parent *domain.ProductParent) error {
	s.parent = *parent
	return nil
}

func (s *stubProductParentRepo) Delete(id uint64) error {
	return nil
}

func TestCreateProductReturnsBadRequestWhenSupplierMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(newTestProductUsecase(&stubProductRepo{}, nil, nil, nil), nil)

	router := gin.New()
	router.POST("/api/v1/products", handler.CreateProduct)

	body, _ := json.Marshal(map[string]any{
		"seller_sku":  "SKU-1",
		"asin":        "ASIN-1",
		"title":       "Product 1",
		"marketplace": "US",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCreateProductReturnsBadRequestWhenAsinTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(newTestProductUsecase(&stubProductRepo{}, nil, nil, nil), nil)

	router := gin.New()
	router.POST("/api/v1/products", handler.CreateProduct)

	body, _ := json.Marshal(map[string]any{
		"seller_sku":  "SKU-1",
		"asin":        "ASIN-TOO-LONG-123456789",
		"title":       "Product 1",
		"marketplace": "US",
		"supplier_id": 9,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCreateProductWritesAudit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(newTestProductUsecase(&stubProductRepo{}, nil, nil, nil), nil)
	auditLogger := &stubProductAuditLogger{}
	handler.BindAuditLogger(auditLogger)

	router := gin.New()
	router.POST("/api/v1/products", handler.CreateProduct)

	body, _ := json.Marshal(map[string]any{
		"seller_sku":   "SKU-1",
		"asin":         "ASIN-1",
		"title":        "Product 1",
		"marketplace":  "US",
		"supplier_id":  9,
		"status":       "ON_SALE",
		"dimension_id": nil,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if len(auditLogger.payloads) != 2 {
		t.Fatalf("expected 2 audit payloads, got %d", len(auditLogger.payloads))
	}
	if auditLogger.payloads[0].Action != "CREATE_PRODUCT" {
		t.Fatalf("unexpected audit action: %+v", auditLogger.payloads[0])
	}
	if auditLogger.payloads[1].Action != "AUTO_CREATE_DEFAULT_SUPPLIER_QUOTE" {
		t.Fatalf("unexpected auto quote action: %+v", auditLogger.payloads[1])
	}
}

func TestUpdateProductSkipsAuditWhenNothingChanged(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productRepo := &stubProductRepo{
		product: domain.Product{
			ID:          1,
			SellerSku:   "SKU-1",
			Asin:        "ASIN-1",
			Title:       "Product 1",
			Marketplace: "US",
			SupplierID:  func() *uint64 { v := uint64(9); return &v }(),
			Status:      "ON_SALE",
		},
	}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, nil), nil)
	auditLogger := &stubProductAuditLogger{}
	handler.BindAuditLogger(auditLogger)

	router := gin.New()
	router.PUT("/api/v1/products/:id", handler.UpdateProduct)

	body, _ := json.Marshal(map[string]any{
		"seller_sku":  "SKU-1",
		"asin":        "ASIN-1",
		"title":       "Product 1",
		"marketplace": "US",
		"supplier_id": 9,
		"status":      "ON_SALE",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if len(auditLogger.payloads) != 0 {
		t.Fatalf("expected 0 audit payloads, got %d", len(auditLogger.payloads))
	}
}

func TestUpdateProductPreservesComboRelationsWhenFieldsOmitted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	supplierID := uint64(9)
	parentID := uint64(33)
	comboID := uint64(77)
	productRepo := &stubProductRepo{
		product: domain.Product{
			ID:          1,
			SellerSku:   "SKU-1",
			Asin:        "ASIN-1",
			Title:       "Product 1",
			Marketplace: "US",
			SupplierID:  &supplierID,
			ParentID:    &parentID,
			ComboID:     &comboID,
			IsComboMain: 1,
			Status:      "ON_SALE",
		},
	}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, nil), nil)
	auditLogger := &stubProductAuditLogger{}
	handler.BindAuditLogger(auditLogger)

	router := gin.New()
	router.PUT("/api/v1/products/:id", handler.UpdateProduct)

	body, _ := json.Marshal(map[string]any{
		"seller_sku":  "SKU-1",
		"asin":        "ASIN-1",
		"title":       "Product 1",
		"marketplace": "US",
		"supplier_id": 9,
		"status":      "ON_SALE",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if productRepo.product.ParentID == nil || *productRepo.product.ParentID != parentID {
		t.Fatalf("expected parent_id to be preserved, got %#v", productRepo.product.ParentID)
	}
	if productRepo.product.ComboID == nil || *productRepo.product.ComboID != comboID {
		t.Fatalf("expected combo_id to be preserved, got %#v", productRepo.product.ComboID)
	}
	if productRepo.product.IsComboMain != 1 {
		t.Fatalf("expected is_combo_main to be preserved, got %d", productRepo.product.IsComboMain)
	}
	if len(auditLogger.payloads) != 0 {
		t.Fatalf("expected 0 audit payloads, got %d", len(auditLogger.payloads))
	}
}

func TestBuildAuditDiffIgnoresProductDerivedMetaFields(t *testing.T) {
	before := map[string]any{
		"title":               "产品A",
		"reference_count":     12,
		"deletable":           false,
		"delete_block_reason": "已被引用",
		"inventory_available": 9,
		"updated_by_name":     "admin",
	}
	after := map[string]any{
		"title":               "产品A",
		"reference_count":     45,
		"deletable":           true,
		"delete_block_reason": "",
		"inventory_available": 20,
		"updated_by_name":     "tester",
	}

	beforeDiff, afterDiff, changed := buildAuditDiff(before, after)
	if changed {
		t.Fatalf("expected derived meta fields to be ignored, got before=%v after=%v", beforeDiff, afterDiff)
	}
}

func TestUpdateProductWritesAutoQuoteAuditWhenSupplierChanged(t *testing.T) {
	gin.SetMode(gin.TestMode)
	oldSupplierID := uint64(9)
	productRepo := &stubProductRepo{
		product: domain.Product{
			ID:          1,
			SellerSku:   "SKU-1",
			Asin:        "ASIN-1",
			Title:       "Product 1",
			Marketplace: "US",
			SupplierID:  &oldSupplierID,
			Status:      "ON_SALE",
		},
	}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, nil), nil)
	auditLogger := &stubProductAuditLogger{}
	handler.BindAuditLogger(auditLogger)

	router := gin.New()
	router.PUT("/api/v1/products/:id", handler.UpdateProduct)

	body, _ := json.Marshal(map[string]any{
		"seller_sku":  "SKU-1",
		"asin":        "ASIN-1",
		"title":       "Product 1",
		"marketplace": "US",
		"supplier_id": 11,
		"status":      "ON_SALE",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if len(auditLogger.payloads) != 2 {
		t.Fatalf("expected 2 audit payloads, got %d", len(auditLogger.payloads))
	}
	if auditLogger.payloads[0].Action != "UPDATE_PRODUCT" {
		t.Fatalf("unexpected first audit action: %+v", auditLogger.payloads[0])
	}
	if auditLogger.payloads[1].Action != "AUTO_CREATE_DEFAULT_SUPPLIER_QUOTE" {
		t.Fatalf("unexpected second audit action: %+v", auditLogger.payloads[1])
	}
}

func TestUpdateProductWritesOnlyChangedFieldsToAudit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	supplierID := uint64(9)
	productRepo := &stubProductRepo{
		product: domain.Product{
			ID:          1,
			SellerSku:   "SKU-1",
			Asin:        "ASIN-1",
			Title:       "Product 1",
			Marketplace: "US",
			SupplierID:  &supplierID,
			Status:      "ON_SALE",
			Remark:      "",
		},
	}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, nil), nil)
	auditLogger := &stubProductAuditLogger{}
	handler.BindAuditLogger(auditLogger)

	router := gin.New()
	router.PUT("/api/v1/products/:id", handler.UpdateProduct)

	body, _ := json.Marshal(map[string]any{
		"seller_sku":  "SKU-1",
		"asin":        "ASIN-1",
		"title":       "Product 1",
		"marketplace": "US",
		"supplier_id": 9,
		"status":      "ON_SALE",
		"remark":      "changed-remark",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if len(auditLogger.payloads) != 1 {
		t.Fatalf("expected 1 audit payload, got %d", len(auditLogger.payloads))
	}

	before, ok := auditLogger.payloads[0].Before.(map[string]any)
	if !ok {
		t.Fatalf("expected before to be map[string]any, got %#v", auditLogger.payloads[0].Before)
	}
	after, ok := auditLogger.payloads[0].After.(map[string]any)
	if !ok {
		t.Fatalf("expected after to be map[string]any, got %#v", auditLogger.payloads[0].After)
	}
	if len(before) != 1 || len(after) != 1 {
		t.Fatalf("expected only changed field in audit payload, before=%v after=%v", before, after)
	}
	if before["remark"] != "" || after["remark"] != "changed-remark" {
		t.Fatalf("unexpected diff payload, before=%v after=%v", before, after)
	}
}

func TestUpdateProductReturnsBadRequestWhenSupplierMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productRepo := &stubProductRepo{
		product: domain.Product{
			ID:          1,
			SellerSku:   "SKU-1",
			Asin:        "ASIN-1",
			Title:       "Product 1",
			Marketplace: "US",
			Status:      "ON_SALE",
		},
	}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, nil), nil)

	router := gin.New()
	router.PUT("/api/v1/products/:id", handler.UpdateProduct)

	body, _ := json.Marshal(map[string]any{
		"seller_sku":  "SKU-1",
		"asin":        "ASIN-1",
		"title":       "Product 1",
		"marketplace": "US",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestUpdateProductReturnsBadRequestWhenAsinTooLong(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productRepo := &stubProductRepo{
		product: domain.Product{
			ID:          1,
			SellerSku:   "SKU-1",
			Asin:        "ASIN-1",
			Title:       "Product 1",
			Marketplace: "US",
			Status:      "DRAFT",
		},
	}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, nil), nil)

	router := gin.New()
	router.PUT("/api/v1/products/:id", handler.UpdateProduct)

	body, _ := json.Marshal(map[string]any{
		"seller_sku":  "SKU-1",
		"asin":        "ASIN-TOO-LONG-123456789",
		"title":       "Product 1",
		"marketplace": "US",
		"supplier_id": 9,
		"status":      "DRAFT",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDeleteProductReturnsBadRequestWhenReferenced(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(newTestProductUsecase(&stubProductRepo{
		referenceCount: map[uint64]int64{1: 2},
	}, nil, nil, nil), nil)

	router := gin.New()
	router.DELETE("/api/v1/products/:id", handler.DeleteProduct)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestUpdateProductParentReturnsBadRequestWhenIdentityChanged(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productRepo := &stubProductRepo{}
	parentRepo := &stubProductParentRepo{
		parent: domain.ProductParent{
			ID:          1,
			ParentAsin:  "PARENT-1",
			Marketplace: "US",
		},
	}
	handler := NewProductHandler(newTestProductUsecase(productRepo, parentRepo, nil, nil), nil)

	router := gin.New()
	router.PUT("/api/v1/product-parents/:id", handler.UpdateProductParent)

	body, _ := json.Marshal(map[string]any{
		"parent_asin": "PARENT-2",
		"title":       "Parent 2",
		"marketplace": "CA",
		"status":      "ON_SALE",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/product-parents/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestListProductsParsesStatusesQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productRepo := &stubProductRepo{}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, nil), nil)

	router := gin.New()
	router.GET("/api/v1/products", handler.ListProducts)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?statuses=ON_SALE,REPLENISHING", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if productRepo.lastListParams == nil {
		t.Fatal("expected list params to be captured")
	}
	if len(productRepo.lastListParams.Statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %+v", productRepo.lastListParams.Statuses)
	}
	if productRepo.lastListParams.Statuses[0] != "ON_SALE" || productRepo.lastListParams.Statuses[1] != "REPLENISHING" {
		t.Fatalf("unexpected statuses: %+v", productRepo.lastListParams.Statuses)
	}
}

func TestListProductsParsesBrandAndCategoryQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productRepo := &stubProductRepo{}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, nil), nil)

	router := gin.New()
	router.GET("/api/v1/products", handler.ListProducts)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?brand_id=12&category_id=34", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}

	if productRepo.lastListParams == nil {
		t.Fatal("expected list params to be captured")
	}
	if productRepo.lastListParams.BrandID == nil || *productRepo.lastListParams.BrandID != 12 {
		t.Fatalf("unexpected brand_id: %+v", productRepo.lastListParams.BrandID)
	}
	if productRepo.lastListParams.CategoryID == nil || *productRepo.lastListParams.CategoryID != 34 {
		t.Fatalf("unexpected category_id: %+v", productRepo.lastListParams.CategoryID)
	}
}

func TestListProductsParsesPackingFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productRepo := &stubProductRepo{}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, nil), nil)

	router := gin.New()
	router.GET("/api/v1/products", handler.ListProducts)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?only_with_packaging=true&packing_required=1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if productRepo.lastListParams == nil {
		t.Fatal("expected list params to be captured")
	}
	if !productRepo.lastListParams.OnlyWithPackaging {
		t.Fatalf("expected only_with_packaging=true, got %+v", productRepo.lastListParams)
	}
	if productRepo.lastListParams.PackingRequired == nil || *productRepo.lastListParams.PackingRequired != 1 {
		t.Fatalf("unexpected packing_required: %+v", productRepo.lastListParams.PackingRequired)
	}
}

func TestDeleteProductConfigReturnsBadRequestWhenReferenced(t *testing.T) {
	gin.SetMode(gin.TestMode)
	brandID := uint64(9)
	productRepo := &stubProductRepo{}
	configRepo := &stubProductConfigRepo{
		item: domain.ProductConfigItem{ID: brandID, ConfigType: domain.ProductConfigTypeBrand, ItemCode: "BRAND-1", ItemName: "品牌1"},
	}
	handler := NewProductHandler(newTestProductUsecase(&stubProductRepoReferenced{
		stubProductRepo: productRepo,
		configRefCount:  1,
	}, nil, configRepo, nil), nil)

	router := gin.New()
	router.DELETE("/api/v1/product-configs/:id", handler.DeleteProductConfig)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/product-configs/9", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDeleteProductConfigWritesAudit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	configRepo := &stubProductConfigRepo{
		item: domain.ProductConfigItem{
			ID:         9,
			ConfigType: domain.ProductConfigTypeBrand,
			ItemCode:   "BRAND_A",
			ItemName:   "Brand A",
		},
	}
	handler := NewProductHandler(newTestProductUsecase(&stubProductRepo{}, nil, configRepo, nil), nil)
	auditLogger := &stubProductAuditLogger{}
	handler.BindAuditLogger(auditLogger)

	router := gin.New()
	router.DELETE("/api/v1/product-configs/:id", handler.DeleteProductConfig)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/product-configs/9", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if len(auditLogger.payloads) != 1 {
		t.Fatalf("expected 1 audit payload, got %d", len(auditLogger.payloads))
	}
	if auditLogger.payloads[0].Action != "DELETE_PRODUCT_CONFIG" {
		t.Fatalf("unexpected audit action: %+v", auditLogger.payloads[0])
	}
}

func TestUpdateProductConfigSkipsAuditWhenNothingChanged(t *testing.T) {
	gin.SetMode(gin.TestMode)
	configRepo := &stubProductConfigRepo{
		item: domain.ProductConfigItem{
			ID:         9,
			ConfigType: domain.ProductConfigTypeBrand,
			ItemCode:   "BRAND_A",
			ItemName:   "Brand A",
			Status:     "ACTIVE",
			Sort:       1,
		},
	}
	handler := NewProductHandler(newTestProductUsecase(&stubProductRepo{}, nil, configRepo, nil), nil)
	auditLogger := &stubProductAuditLogger{}
	handler.BindAuditLogger(auditLogger)

	router := gin.New()
	router.PUT("/api/v1/product-configs/:id", handler.UpdateProductConfig)

	body, _ := json.Marshal(map[string]any{
		"config_type": "BRAND",
		"item_code":   "BRAND_A",
		"item_name":   "Brand A",
		"status":      "ACTIVE",
		"sort":        1,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/product-configs/9", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", w.Code, w.Body.String())
	}
	if len(auditLogger.payloads) != 0 {
		t.Fatalf("expected 0 audit payloads, got %d", len(auditLogger.payloads))
	}
}

func TestDeleteProductCategoryReturnsBadRequestWhenReferenced(t *testing.T) {
	gin.SetMode(gin.TestMode)
	productRepo := &stubProductRepoReferenced{configRefCount: 0, categoryRefCount: 1}
	categoryRepo := &stubProductCategoryRepo{
		item: domain.ProductCategory{ID: 9, CategoryCode: "CAT-2", CategoryName: "二级品类", Level: 2, Status: "ACTIVE"},
	}
	handler := NewProductHandler(newTestProductUsecase(productRepo, nil, nil, categoryRepo), nil)

	router := gin.New()
	router.DELETE("/api/v1/product-categories/:id", handler.DeleteProductCategory)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/product-categories/9", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", w.Code, w.Body.String())
	}
}

type stubProductRepoReferenced struct {
	*stubProductRepo
	configRefCount   int64
	categoryRefCount int64
}

func (s *stubProductRepoReferenced) CountByConfigReference(configType domain.ProductConfigType, configID uint64) (int64, error) {
	return s.configRefCount, nil
}

func (s *stubProductRepoReferenced) CountByCategoryID(categoryID uint64) (int64, error) {
	return s.categoryRefCount, nil
}
