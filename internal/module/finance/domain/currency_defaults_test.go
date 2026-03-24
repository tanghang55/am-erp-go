package domain

import (
	"reflect"
	"strings"
	"testing"
)

func TestCashLedgerCurrencyFieldsDoNotUseHardcodedDefaults(t *testing.T) {
	typ := reflect.TypeOf(CashLedger{})
	for _, fieldName := range []string{"Currency", "OriginalCurrency", "BaseCurrency"} {
		field, ok := typ.FieldByName(fieldName)
		if !ok {
			t.Fatalf("missing field %s", fieldName)
		}
		tag := field.Tag.Get("gorm")
		if strings.Contains(tag, "default:'USD'") || strings.Contains(tag, "default:'CNY'") {
			t.Fatalf("expected %s gorm tag to avoid hardcoded currency default, got %s", fieldName, tag)
		}
	}
}

func TestCostingSnapshotCurrencyFieldDoesNotUseHardcodedDefault(t *testing.T) {
	field, ok := reflect.TypeOf(CostingSnapshot{}).FieldByName("Currency")
	if !ok {
		t.Fatal("missing field Currency")
	}
	tag := field.Tag.Get("gorm")
	if strings.Contains(tag, "default:'USD'") || strings.Contains(tag, "default:'CNY'") {
		t.Fatalf("expected Currency gorm tag to avoid hardcoded currency default, got %s", tag)
	}
}
