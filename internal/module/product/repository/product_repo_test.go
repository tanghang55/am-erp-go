package repository

import (
	"strings"
	"testing"
)

func TestProductDisplaySelectIncludesUpdatedByName(t *testing.T) {
	selectSQL := productDisplaySelect()

	if !strings.Contains(selectSQL, "AS updated_by_name") {
		t.Fatalf("expected product display select to include updated_by_name, got: %s", selectSQL)
	}

	if !strings.Contains(selectSQL, "audit_log.username") {
		t.Fatalf("expected product display select to resolve updater from audit_log, got: %s", selectSQL)
	}

	if !strings.Contains(selectSQL, "COLLATE utf8mb4_unicode_ci") {
		t.Fatalf("expected product display select to normalize collation for updated_by_name subquery, got: %s", selectSQL)
	}

	if !strings.Contains(selectSQL, "COALESCE(default_quote.price, product.unit_cost) AS unit_cost") {
		t.Fatalf("expected product display select to derive unit_cost from default supplier quote, got: %s", selectSQL)
	}
}

func TestProductListOrderByUsesCreateTimeDesc(t *testing.T) {
	orderBy := productListOrderBy()

	if orderBy != "product.gmt_create DESC, product.id DESC" {
		t.Fatalf("expected product list order by create time desc, got: %s", orderBy)
	}
}

func TestProductListByIDsSelectIncludesDisplayFields(t *testing.T) {
	selectSQL := productListByIDsSelect()

	expectedFragments := []string{
		"COALESCE(default_quote.price, product.unit_cost) AS unit_cost",
		"supplier.name AS supplier_name",
		"supplier.supplier_code AS supplier_code",
		"brand_cfg.item_name AS brand_name",
		"category_cfg.category_name AS category_name",
		"dimension_unit_cfg.item_name AS dimension_unit_name",
		"weight_unit_cfg.item_name AS weight_unit_name",
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(selectSQL, fragment) {
			t.Fatalf("expected product list-by-ids select to include %q, got: %s", fragment, selectSQL)
		}
	}
}

func TestMapProductPackagingRowsBuildsPackagingItemDetail(t *testing.T) {
	rows := []productPackagingRow{
		{
			ID:              1,
			ProductID:       18,
			PackagingItemID: 5,
			QuantityPerUnit: 2,
			PackagingItemDetailData: PackagingItemDetailData{
				ID:             5,
				ItemCode:       "BOX-001",
				ItemName:       "纸箱",
				Specification:  ptrString("20x15x10cm"),
				Unit:           "PCS",
				UnitCost:       2.5,
				Currency:       "CNY",
				QuantityOnHand: 100,
			},
		},
	}

	items := mapProductPackagingRows(rows)
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].PackagingItem == nil {
		t.Fatalf("expected packaging item detail to be populated")
	}
	if items[0].PackagingItem.ItemCode != "BOX-001" {
		t.Fatalf("expected item code BOX-001, got %#v", items[0].PackagingItem.ItemCode)
	}
	if items[0].PackagingItem.ItemName != "纸箱" {
		t.Fatalf("expected item name 纸箱, got %#v", items[0].PackagingItem.ItemName)
	}
}

func ptrString(value string) *string {
	return &value
}
