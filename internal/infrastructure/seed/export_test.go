package seed

import (
	"strings"
	"testing"
	"time"
)

func TestMinimalSeedTablesReturnsSortedCopy(t *testing.T) {
	tables := MinimalSeedTables()
	expected := []string{"finance_exchange_rate", "menu", "permission", "role", "role_permission"}
	if len(tables) != len(expected) {
		t.Fatalf("expected %d tables, got %d", len(expected), len(tables))
	}
	for i := range expected {
		if tables[i] != expected[i] {
			t.Fatalf("expected table %s at %d, got %s", expected[i], i, tables[i])
		}
	}

	tables[0] = "changed"
	again := MinimalSeedTables()
	if again[0] != "finance_exchange_rate" {
		t.Fatalf("expected defensive copy, got %#v", again)
	}
}

func TestSQLLiteralEscapesSpecialCharacters(t *testing.T) {
	value := sqlLiteral("a'b\\c\n")
	if value != "'a\\'b\\\\c\\n'" {
		t.Fatalf("unexpected literal: %s", value)
	}
}

func TestSQLLiteralFormatsTimeWithoutTimezone(t *testing.T) {
	value := sqlLiteral(time.Date(2026, 3, 8, 15, 16, 17, 0, time.FixedZone("CST", 8*3600)))
	if value != "'2026-03-08 15:16:17'" {
		t.Fatalf("unexpected time literal: %s", value)
	}
}

func TestExportMinimalSeedSQLContainsRequiredDeletes(t *testing.T) {
	sqlText := strings.Join([]string{
		"SET NAMES utf8mb4;",
		"SET FOREIGN_KEY_CHECKS = 0;",
		"",
		"DELETE FROM `user_role`;",
		"DELETE FROM `role_permission`;",
		"DELETE FROM `user`;",
		"DELETE FROM `finance_exchange_rate`;",
		"DELETE FROM `menu`;",
		"DELETE FROM `permission`;",
		"DELETE FROM `role`;",
	}, "\n")

	for _, required := range []string{
		"DELETE FROM `user_role`;",
		"DELETE FROM `role_permission`;",
		"DELETE FROM `user`;",
		"DELETE FROM `finance_exchange_rate`;",
		"DELETE FROM `menu`;",
		"DELETE FROM `permission`;",
		"DELETE FROM `role`;",
	} {
		if !strings.Contains(sqlText, required) {
			t.Fatalf("expected delete statement %s", required)
		}
	}
}
