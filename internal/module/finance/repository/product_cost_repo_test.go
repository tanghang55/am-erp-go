package repository

import (
	"strings"
	"testing"
)

func TestProductCostUnionSQLIncludesShipmentAllocated(t *testing.T) {
	sql := productCostUnionSQL()
	if !strings.Contains(sql, "SHIPMENT_ALLOCATED") {
		t.Fatalf("expected shipment allocated event in union sql")
	}
	if !strings.Contains(sql, "WHEN ce.event_type IN ('SHIPMENT_ALLOCATED', 'PACKING_MATERIAL') THEN 'NEUTRAL'") {
		t.Fatalf("expected neutral direction mapping for allocation events")
	}
	if !strings.Contains(sql, "WHEN ce.event_type = 'SHIPMENT_ALLOCATED' THEN 'SHIPMENT'") {
		t.Fatalf("expected shipment reference mapping in union sql")
	}
	if !strings.Contains(sql, "ce.event_type IN ('SHIPMENT_ALLOCATED', 'PACKING_MATERIAL') THEN 0") {
		t.Fatalf("expected quantity neutralization for shipment allocation")
	}
}

func TestProductCostUnionSQLIncludesPackingMaterial(t *testing.T) {
	sql := productCostUnionSQL()
	if !strings.Contains(sql, "PACKING_MATERIAL") {
		t.Fatalf("expected packing material event in union sql")
	}
	if !strings.Contains(sql, "COALESCE(im.reference_type, 'PRODUCT_PACKING')") {
		t.Fatalf("expected product packing reference mapping in union sql")
	}
	if !strings.Contains(sql, "ce.event_type IN ('SHIPMENT_ALLOCATED', 'PACKING_MATERIAL') THEN 0") {
		t.Fatalf("expected quantity neutralization for packing material allocation")
	}
}

func TestProductCostSummarySQLIncludesNeutralAmounts(t *testing.T) {
	sql := productCostSummarySQL("1=1", productCostUnionSQL())

	expectedFragments := []string{
		"COALESCE(SUM(CASE WHEN direction = 'NEUTRAL' THEN base_amount ELSE 0 END), 0) AS neutral_amount",
		"COALESCE(SUM(CASE WHEN direction = 'INBOUND' THEN base_amount ELSE 0 END), 0) - COALESCE(SUM(CASE WHEN direction = 'OUTBOUND' THEN base_amount ELSE 0 END), 0) - COALESCE(SUM(CASE WHEN direction = 'NEUTRAL' THEN base_amount ELSE 0 END), 0) AS net_amount",
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(sql, fragment) {
			t.Fatalf("expected product cost summary sql to include %q, got: %s", fragment, sql)
		}
	}
}
