package repository

import "testing"

func TestQualifiedComboIDColumn(t *testing.T) {
	if got := qualifiedComboIDColumn(); got != "product_combo.combo_id" {
		t.Fatalf("expected qualified combo_id column, got %s", got)
	}
}
