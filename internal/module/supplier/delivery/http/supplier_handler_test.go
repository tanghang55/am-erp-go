package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/supplier/domain"
	"am-erp-go/internal/module/supplier/usecase"

	"github.com/gin-gonic/gin"
)

type stubSupplierRepo struct {
	listParams *domain.SupplierListParams
	createItem *domain.Supplier
	err        error
}

func (s *stubSupplierRepo) List(params *domain.SupplierListParams) ([]domain.Supplier, int64, error) {
	s.listParams = params
	return nil, 0, s.err
}
func (s *stubSupplierRepo) GetByID(_ uint64) (*domain.Supplier, error) { return nil, s.err }
func (s *stubSupplierRepo) Create(supplier *domain.Supplier) error {
	s.createItem = supplier
	if supplier.ID == 0 {
		supplier.ID = 1
	}
	return s.err
}
func (s *stubSupplierRepo) Update(_ *domain.Supplier) error { return s.err }
func (s *stubSupplierRepo) Delete(_ uint64) error           { return s.err }

type stubSupplierTypeRepo struct {
	replace   []string
	replaceID uint64
	err       error
}

func (s *stubSupplierTypeRepo) ListBySupplierID(_ uint64) ([]string, error) {
	return nil, s.err
}
func (s *stubSupplierTypeRepo) ListBySupplierIDs(_ []uint64) (map[uint64][]string, error) {
	return map[uint64][]string{}, s.err
}
func (s *stubSupplierTypeRepo) ReplaceBySupplierID(id uint64, types []string) error {
	s.replaceID = id
	s.replace = types
	return s.err
}

type stubSupplierContactRepo struct {
	create           *domain.SupplierContact
	update           *domain.SupplierContact
	deleteID         uint64
	deleteSupplierID uint64
	err              error
}

func (s *stubSupplierContactRepo) ListBySupplierID(_ uint64) ([]domain.SupplierContact, error) {
	return nil, s.err
}
func (s *stubSupplierContactRepo) Create(contact *domain.SupplierContact) error {
	s.create = contact
	return s.err
}
func (s *stubSupplierContactRepo) Update(contact *domain.SupplierContact) error {
	s.update = contact
	return s.err
}
func (s *stubSupplierContactRepo) Delete(id uint64, supplierID uint64) error {
	s.deleteID = id
	s.deleteSupplierID = supplierID
	return s.err
}

type stubSupplierAccountRepo struct {
	create           *domain.SupplierAccount
	update           *domain.SupplierAccount
	deleteID         uint64
	deleteSupplierID uint64
	err              error
}

func (s *stubSupplierAccountRepo) ListBySupplierID(_ uint64) ([]domain.SupplierAccount, error) {
	return nil, s.err
}
func (s *stubSupplierAccountRepo) Create(account *domain.SupplierAccount) error {
	s.create = account
	return s.err
}
func (s *stubSupplierAccountRepo) Update(account *domain.SupplierAccount) error {
	s.update = account
	return s.err
}
func (s *stubSupplierAccountRepo) Delete(id uint64, supplierID uint64) error {
	s.deleteID = id
	s.deleteSupplierID = supplierID
	return s.err
}

type stubSupplierTagRepo struct {
	create           *domain.SupplierTag
	update           *domain.SupplierTag
	deleteID         uint64
	deleteSupplierID uint64
	err              error
}

func (s *stubSupplierTagRepo) ListBySupplierID(_ uint64) ([]domain.SupplierTag, error) {
	return nil, s.err
}
func (s *stubSupplierTagRepo) Create(tag *domain.SupplierTag) error {
	s.create = tag
	return s.err
}
func (s *stubSupplierTagRepo) Update(tag *domain.SupplierTag) error {
	s.update = tag
	return s.err
}
func (s *stubSupplierTagRepo) Delete(id uint64, supplierID uint64) error {
	s.deleteID = id
	s.deleteSupplierID = supplierID
	return s.err
}

func TestListSuppliersParsesTypeFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	supplierRepo := &stubSupplierRepo{}
	typeRepo := &stubSupplierTypeRepo{}
	uc := usecase.NewSupplierUsecase(supplierRepo, typeRepo, nil, nil, nil)
	handler := NewSupplierHandler(uc)

	router := gin.New()
	router.GET("/api/v1/suppliers", handler.ListSuppliers)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/suppliers?type=PRODUCT,LOGISTICS", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if supplierRepo.listParams == nil || len(supplierRepo.listParams.Types) != 2 {
		t.Fatalf("expected type filter parsed")
	}
}

func TestCreateSupplierParsesTypes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	supplierRepo := &stubSupplierRepo{}
	typeRepo := &stubSupplierTypeRepo{}
	uc := usecase.NewSupplierUsecase(supplierRepo, typeRepo, nil, nil, nil)
	handler := NewSupplierHandler(uc)

	router := gin.New()
	router.POST("/api/v1/suppliers", handler.CreateSupplier)

	payload := []byte(`{"supplier_code":"S-001","name":"Demo","types":["PRODUCT"]}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if supplierRepo.createItem == nil || typeRepo.replaceID == 0 || len(typeRepo.replace) != 1 {
		t.Fatalf("expected supplier create + types")
	}
}

func TestSupplierContactEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	contactRepo := &stubSupplierContactRepo{}
	uc := usecase.NewSupplierUsecase(&stubSupplierRepo{}, &stubSupplierTypeRepo{}, contactRepo, nil, nil)
	handler := NewSupplierHandler(uc)

	router := gin.New()
	router.POST("/api/v1/suppliers/:id/contacts", handler.CreateSupplierContact)
	router.PUT("/api/v1/suppliers/:id/contacts", handler.UpdateSupplierContact)
	router.DELETE("/api/v1/suppliers/:id/contacts", handler.DeleteSupplierContact)

	createPayload := []byte(`{"name":"Kate","phone":"123","email":"a@b.com","position":"QA","is_primary":1}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers/5/contacts", bytes.NewReader(createPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || contactRepo.create == nil || contactRepo.create.SupplierID != 5 {
		t.Fatalf("expected contact create")
	}

	updatePayload := []byte(`{"id":3,"name":"Kate","phone":"456","email":"a@b.com","position":"QA","is_primary":0}`)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/api/v1/suppliers/5/contacts", bytes.NewReader(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || contactRepo.update == nil || contactRepo.update.SupplierID != 5 {
		t.Fatalf("expected contact update")
	}

	deletePayload := []byte(`{"id":3}`)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/suppliers/5/contacts", bytes.NewReader(deletePayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || contactRepo.deleteID != 3 || contactRepo.deleteSupplierID != 5 {
		t.Fatalf("expected contact delete")
	}
}

func TestSupplierAccountEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	accountRepo := &stubSupplierAccountRepo{}
	uc := usecase.NewSupplierUsecase(&stubSupplierRepo{}, &stubSupplierTypeRepo{}, nil, accountRepo, nil)
	handler := NewSupplierHandler(uc)

	router := gin.New()
	router.POST("/api/v1/suppliers/:id/accounts", handler.CreateSupplierAccount)
	router.PUT("/api/v1/suppliers/:id/accounts", handler.UpdateSupplierAccount)
	router.DELETE("/api/v1/suppliers/:id/accounts", handler.DeleteSupplierAccount)

	createPayload := []byte(`{"bank_name":"HSBC","bank_account":"123","currency":"USD","tax_no":"T1","payment_terms":"NET30"}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers/6/accounts", bytes.NewReader(createPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || accountRepo.create == nil || accountRepo.create.SupplierID != 6 {
		t.Fatalf("expected account create")
	}

	updatePayload := []byte(`{"id":2,"bank_name":"HSBC","bank_account":"234","currency":"USD","tax_no":"T1","payment_terms":"NET45"}`)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/api/v1/suppliers/6/accounts", bytes.NewReader(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || accountRepo.update == nil || accountRepo.update.SupplierID != 6 {
		t.Fatalf("expected account update")
	}

	deletePayload := []byte(`{"id":2}`)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/suppliers/6/accounts", bytes.NewReader(deletePayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || accountRepo.deleteID != 2 || accountRepo.deleteSupplierID != 6 {
		t.Fatalf("expected account delete")
	}
}

func TestSupplierTagEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tagRepo := &stubSupplierTagRepo{}
	uc := usecase.NewSupplierUsecase(&stubSupplierRepo{}, &stubSupplierTypeRepo{}, nil, nil, tagRepo)
	handler := NewSupplierHandler(uc)

	router := gin.New()
	router.POST("/api/v1/suppliers/:id/tags", handler.CreateSupplierTag)
	router.PUT("/api/v1/suppliers/:id/tags", handler.UpdateSupplierTag)
	router.DELETE("/api/v1/suppliers/:id/tags", handler.DeleteSupplierTag)

	createPayload := []byte(`{"tag":"VIP"}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/suppliers/4/tags", bytes.NewReader(createPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || tagRepo.create == nil || tagRepo.create.SupplierID != 4 {
		t.Fatalf("expected tag create")
	}

	updatePayload := []byte(`{"id":9,"tag":"VIP-2"}`)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/api/v1/suppliers/4/tags", bytes.NewReader(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || tagRepo.update == nil || tagRepo.update.SupplierID != 4 {
		t.Fatalf("expected tag update")
	}

	deletePayload := []byte(`{"id":9}`)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/suppliers/4/tags", bytes.NewReader(deletePayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || tagRepo.deleteID != 9 || tagRepo.deleteSupplierID != 4 {
		t.Fatalf("expected tag delete")
	}
}
