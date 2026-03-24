package repository

import (
	"strings"
	"testing"
)

func TestPackagingItemColumnPrefixesTable(t *testing.T) {
	if got := packagingItemColumn("status"); got != "packaging_item.status" {
		t.Fatalf("expected qualified packaging_item column, got %s", got)
	}
}

func TestPackagingItemReferenceCountSQLUsesPluralProductPackagingItemsTable(t *testing.T) {
	sql := packagingItemReferenceCountSQL()
	if !strings.Contains(sql, "product_packaging_items") {
		t.Fatalf("expected SQL to query product_packaging_items, got %s", sql)
	}
	if strings.Contains(sql, "product_packaging_item ") {
		t.Fatalf("expected SQL not to query old singular product_packaging_item table, got %s", sql)
	}
	if !strings.Contains(sql, "package_spec_packaging_items") {
		t.Fatalf("expected SQL to query package_spec_packaging_items, got %s", sql)
	}
	if strings.Contains(sql, "package_spec_packaging_item ") {
		t.Fatalf("expected SQL not to query old singular package_spec_packaging_item table, got %s", sql)
	}
}
