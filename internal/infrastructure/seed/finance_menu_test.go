package seed

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMinimalSeedIncludesFinanceOrderProfitMenu(t *testing.T) {
	path := filepath.Join("..", "..", "..", "baseline", "minimal", "minimal_seed.sql")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read minimal seed: %v", err)
	}

	text := string(content)
	wantSnippets := []string{
		"('48','财务总览'",
		"'/finance/profit','TrendCharts',NULL,'610'",
		"('21','现金流水'",
		"'/finance/cash-ledger',NULL,NULL,'620'",
		"('22','成本中心'",
		"'/finance/costing',NULL,NULL,'630'",
		"('57','订单利润'",
		"'/finance/order-profit',NULL,NULL,'640'",
		"('55','汇率管理'",
		"'/finance/exchange-rates','Coin',NULL,'650'",
	}
	for _, snippet := range wantSnippets {
		if !strings.Contains(text, snippet) {
			t.Fatalf("expected minimal seed to contain %q", snippet)
		}
	}
}
