package repository

import (
	"strings"
	"testing"
)

func TestProfitSummaryBaseCurrencyExprDoesNotUseHardcodedFallback(t *testing.T) {
	expr := profitSummaryBaseCurrencyExpr()
	if strings.Contains(expr, "'USD'") || strings.Contains(expr, "'CNY'") {
		t.Fatalf("expected no hardcoded currency fallback, got %s", expr)
	}
}

func TestProductCostSummaryBaseCurrencyExprDoesNotUseHardcodedFallback(t *testing.T) {
	expr := productCostSummaryBaseCurrencyExpr()
	if strings.Contains(expr, "'USD'") || strings.Contains(expr, "'CNY'") {
		t.Fatalf("expected no hardcoded currency fallback, got %s", expr)
	}
}
