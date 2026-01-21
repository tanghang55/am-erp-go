package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	systemdomain "am-erp-go/internal/module/system/domain"

	"github.com/gin-gonic/gin"
)

type stubFieldLabelUsecase struct {
	labels    map[string]string
	listItems []*systemdomain.FieldLabel
	total     int64
	listArgs  struct {
		page     int
		pageSize int
		keyword  string
	}
	created *systemdomain.FieldLabel
	updated *systemdomain.FieldLabel
	deleted uint64
	err     error
}

func (s *stubFieldLabelUsecase) GetLabels(locale string) (map[string]string, error) {
	return s.labels, s.err
}

func (s *stubFieldLabelUsecase) List(page, pageSize int, keyword string) ([]*systemdomain.FieldLabel, int64, error) {
	s.listArgs.page = page
	s.listArgs.pageSize = pageSize
	s.listArgs.keyword = keyword
	return s.listItems, s.total, s.err
}

func (s *stubFieldLabelUsecase) Create(label *systemdomain.FieldLabel) error {
	s.created = label
	return s.err
}

func (s *stubFieldLabelUsecase) Update(label *systemdomain.FieldLabel) error {
	s.updated = label
	return s.err
}

func (s *stubFieldLabelUsecase) Delete(id uint64) error {
	s.deleted = id
	return s.err
}

func TestGetFieldLabelsReturnsLocale(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubFieldLabelUsecase{
		labels: map[string]string{"product.list.title": "Product List"},
	}
	handler := NewFieldLabelHandler(stub)

	router := gin.New()
	router.GET("/api/v1/system/field-labels", handler.GetLabels)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/system/field-labels?locale=zh-CN", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestListFieldLabelsParsesParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubFieldLabelUsecase{}
	handler := NewFieldLabelHandler(stub)

	router := gin.New()
	router.GET("/api/v1/system/field-labels/manage", handler.List)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/system/field-labels/manage?page=2&page_size=10&keyword=abc", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if stub.listArgs.page != 2 || stub.listArgs.pageSize != 10 || stub.listArgs.keyword != "abc" {
		t.Fatalf("unexpected list args: %+v", stub.listArgs)
	}
}

func TestCreateFieldLabelCallsUsecase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubFieldLabelUsecase{}
	handler := NewFieldLabelHandler(stub)

	router := gin.New()
	router.POST("/api/v1/system/field-labels", handler.Create)

	payload := []byte(`{"label_key":"product.list.title","labels":{"zh-CN":"Product"}}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/system/field-labels", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if stub.created == nil || stub.created.LabelKey != "product.list.title" {
		t.Fatalf("expected created label to be captured")
	}
}

func TestUpdateFieldLabelCallsUsecase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubFieldLabelUsecase{}
	handler := NewFieldLabelHandler(stub)

	router := gin.New()
	router.PUT("/api/v1/system/field-labels/:id", handler.Update)

	payload := []byte(`{"label_key":"product.list.title","labels":{"zh-CN":"Product"}}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/system/field-labels/9", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if stub.updated == nil || stub.updated.ID != 9 {
		t.Fatalf("expected updated label id to be 9")
	}
}

func TestDeleteFieldLabelCallsUsecase(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stub := &stubFieldLabelUsecase{}
	handler := NewFieldLabelHandler(stub)

	router := gin.New()
	router.DELETE("/api/v1/system/field-labels/:id", handler.Delete)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/system/field-labels/7", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if stub.deleted != 7 {
		t.Fatalf("expected deleted id to be 7")
	}
}
