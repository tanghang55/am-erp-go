package usecase

import "testing"

func TestIsExcludedProductStatus(t *testing.T) {
	cases := map[string]bool{
		"DRAFT":        true,
		"OFF_SHELF":    true,
		"ON_SALE":      false,
		"REPLENISHING": false,
		"":             false,
	}

	for input, want := range cases {
		if got := isExcludedProductStatus(input); got != want {
			t.Fatalf("status %q expected %v, got %v", input, want, got)
		}
	}
}
