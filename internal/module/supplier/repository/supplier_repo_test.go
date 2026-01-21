package repository

import (
	"testing"

	"am-erp-go/internal/module/supplier/domain"
)

func TestSupplierRepositoriesImplementInterfaces(t *testing.T) {
	var _ domain.SupplierTypeRepository = (*supplierTypeRepository)(nil)
	var _ domain.SupplierContactRepository = (*supplierContactRepository)(nil)
	var _ domain.SupplierAccountRepository = (*supplierAccountRepository)(nil)
	var _ domain.SupplierTagRepository = (*supplierTagRepository)(nil)
}
