package repository

import "testing"

func TestBuildProductQuoteQueryQualifiesGmtModified(t *testing.T) {
	expected := "product.id, product.seller_sku, product.asin, product.marketplace, product.title, product.image_url, product.supplier_id, product.gmt_modified AS gmt_modified"
	if got := productQuoteSelectClause(); got != expected {
		t.Fatalf("expected qualified gmt_modified select, got %s", got)
	}
}
