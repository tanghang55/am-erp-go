package repository

import (
	"testing"

	"am-erp-go/internal/module/product/domain"
)

func TestProductImageRepositoryImplementsInterface(t *testing.T) {
	var _ domain.ProductImageRepository = (*productImageRepository)(nil)
}
