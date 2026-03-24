package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"am-erp-go/internal/module/shipment/domain"
	"am-erp-go/internal/module/shipment/usecase"

	"github.com/gin-gonic/gin"
)

type stubPackageSpecRepo struct {
	referenceCount int64
}

func (s *stubPackageSpecRepo) Create(_ *domain.PackageSpec) error { return nil }
func (s *stubPackageSpecRepo) Update(_ *domain.PackageSpec) error { return nil }
func (s *stubPackageSpecRepo) GetByID(_ uint64) (*domain.PackageSpec, error) {
	return &domain.PackageSpec{ID: 1, Name: "标准箱"}, nil
}
func (s *stubPackageSpecRepo) List(_ *domain.PackageSpecListParams) ([]*domain.PackageSpec, int64, error) {
	return []*domain.PackageSpec{}, 0, nil
}
func (s *stubPackageSpecRepo) Delete(_ uint64) error { return nil }
func (s *stubPackageSpecRepo) ListByIDs(_ []uint64) ([]*domain.PackageSpec, error) {
	return []*domain.PackageSpec{}, nil
}
func (s *stubPackageSpecRepo) CountReferences(_ uint64) (int64, error) {
	return s.referenceCount, nil
}

type stubPackageSpecPackagingRepo struct{}

func (s *stubPackageSpecPackagingRepo) ListByPackageSpecID(_ uint64) ([]domain.PackageSpecPackagingItem, error) {
	return []domain.PackageSpecPackagingItem{}, nil
}

func (s *stubPackageSpecPackagingRepo) ReplaceAll(_ uint64, _ []domain.PackageSpecPackagingItem) error {
	return nil
}

func TestDeletePackageSpecReturnsBadRequestWhenReferenced(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubPackageSpecRepo{referenceCount: 1}
	uc := usecase.NewPackageSpecUseCase(repo, &stubPackageSpecPackagingRepo{})
	handler := NewPackageSpecHandler(uc)

	router := gin.New()
	router.DELETE("/api/v1/package-specs/:id", handler.Delete)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/package-specs/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
