package repository

import (
	"strings"
	"testing"
)

func TestCostingSnapshotSelectSQLIncludesProductDisplayFields(t *testing.T) {
	selectSQL := costingSnapshotSelectSQL()

	expectedFragments := []string{
		"costing_snapshot.*",
		"COALESCE(product.seller_sku, '') AS seller_sku",
		"COALESCE(product.title, '') AS product_title",
		"COALESCE(product.image_url, '') AS product_image_url",
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(selectSQL, fragment) {
			t.Fatalf("expected costing snapshot select to include %q, got: %s", fragment, selectSQL)
		}
	}
}

func TestCurrentCostTimeWindowWhereSQLQualifiesSnapshotColumns(t *testing.T) {
	whereSQL := currentCostTimeWindowWhereSQL()

	expectedFragments := []string{
		"costing_snapshot.effective_from",
		"costing_snapshot.effective_to",
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(whereSQL, fragment) {
			t.Fatalf("expected current cost where to include %q, got: %s", fragment, whereSQL)
		}
	}
}

func TestCurrentCostOrderByQualifiesSnapshotColumns(t *testing.T) {
	orderBy := currentCostOrderBy()
	if orderBy != "costing_snapshot.cost_type ASC, costing_snapshot.effective_from DESC, costing_snapshot.id DESC" {
		t.Fatalf("expected qualified current cost order, got: %s", orderBy)
	}
}
