package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	systemdomain "am-erp-go/internal/module/system/domain"

	"github.com/gin-gonic/gin"
)

type stubAuditLogUsecase struct {
	params systemdomain.AuditLogListParams
	list   []*systemdomain.AuditLog
	total  int64
	err    error
}

func (s *stubAuditLogUsecase) List(params systemdomain.AuditLogListParams) ([]*systemdomain.AuditLog, int64, error) {
	s.params = params
	return s.list, s.total, s.err
}

func TestListAuditLogsParsesQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	usecase := &stubAuditLogUsecase{}
	handler := NewAuditLogHandler(usecase)

	router := gin.New()
	router.GET("/api/v1/system/logs", handler.List)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/system/logs?page=2&page_size=10&module=Product&action=UPDATE&username=demo&entity_type=SKU&entity_id=1&keyword=abc&date_from=2026-01-01&date_to=2026-01-02", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if usecase.params.Page != 2 || usecase.params.PageSize != 10 || usecase.params.Module != "Product" || usecase.params.Action != "UPDATE" {
		t.Fatalf("unexpected params: %+v", usecase.params)
	}
	if usecase.params.DateFrom == "" || usecase.params.DateTo == "" {
		t.Fatalf("expected date range parsed")
	}
}
