package migration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListPendingMigrationsSortsAndSkipsNonVersionedFiles(t *testing.T) {
	dir := t.TempDir()
	files := []string{
		"am-erp.sql",
		"20260313_rbac_permissions.sql",
		"20260302_sales_order.sql",
		"README.txt",
		"20260309_packaging_procurement.sql",
	}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("-- test"), 0o644); err != nil {
			t.Fatalf("write file %s: %v", name, err)
		}
	}

	got, err := ListPendingMigrations(dir, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		"20260302_sales_order.sql",
		"20260309_packaging_procurement.sql",
		"20260313_rbac_permissions.sql",
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d files, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i].Version != want[i] {
			t.Fatalf("expected %s at %d, got %s", want[i], i, got[i].Version)
		}
	}
}

func TestListPendingMigrationsSkipsAppliedVersions(t *testing.T) {
	dir := t.TempDir()
	files := []string{
		"20260302_sales_order.sql",
		"20260303_inventory_lot.sql",
	}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("-- test"), 0o644); err != nil {
			t.Fatalf("write file %s: %v", name, err)
		}
	}

	got, err := ListPendingMigrations(dir, map[string]struct{}{
		"20260302_sales_order.sql": {},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Version != "20260303_inventory_lot.sql" {
		t.Fatalf("expected only unapplied migration, got %#v", got)
	}
}

func TestRequireBaselineReturnsTrueForExistingSchemaWithoutHistory(t *testing.T) {
	if !RequireBaseline(0, 3, 5) {
		t.Fatalf("expected baseline to be required")
	}
}

func TestRequireBaselineReturnsFalseForFreshDatabase(t *testing.T) {
	if RequireBaseline(0, 3, 0) {
		t.Fatalf("expected fresh database not to require baseline")
	}
}

func TestListMissingAppliedMigrationsFindsVersionsNotPresentOnDisk(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "20260302_sales_order.sql"), []byte("-- test"), 0o644); err != nil {
		t.Fatalf("write migration file: %v", err)
	}

	missing, err := ListMissingAppliedMigrations(dir, map[string]struct{}{
		"20260302_sales_order.sql":            {},
		"20260303_inventory_lot.sql":          {},
		"20260313_rbac_permissions.sql":       {},
		"20260308_system_monitor_menu.sql":    {},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		"20260303_inventory_lot.sql",
		"20260308_system_monitor_menu.sql",
		"20260313_rbac_permissions.sql",
	}
	if len(missing) != len(want) {
		t.Fatalf("expected %d missing files, got %#v", len(want), missing)
	}
	for i := range want {
		if missing[i] != want[i] {
			t.Fatalf("expected %s at %d, got %s", want[i], i, missing[i])
		}
	}
}

func TestLoadVersionManifestSkipsBlankAndCommentLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline_versions.txt")
	content := "# baseline versions\n\n20260302_sales_order.sql\n 20260303_inventory_lot.sql \n\n# keep going\n20260304_procurement_replenishment.sql\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	versions, err := LoadVersionManifest(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		"20260302_sales_order.sql",
		"20260303_inventory_lot.sql",
		"20260304_procurement_replenishment.sql",
	}
	if len(versions) != len(want) {
		t.Fatalf("expected %d versions, got %#v", len(want), versions)
	}
	for i := range want {
		if versions[i] != want[i] {
			t.Fatalf("expected %s at %d, got %s", want[i], i, versions[i])
		}
	}
}

func TestSplitSQLStatementsSupportsMultipleStatements(t *testing.T) {
	sqlText := `
INSERT INTO test_table (name, remark) VALUES ('alpha', 'first');
INSERT INTO test_table (name, remark) VALUES ('beta', 'semi; inside');
`

	statements := splitSQLStatements(sqlText)
	if len(statements) != 2 {
		t.Fatalf("expected 2 statements, got %#v", statements)
	}
	if statements[0] != "INSERT INTO test_table (name, remark) VALUES ('alpha', 'first')" {
		t.Fatalf("unexpected first statement: %s", statements[0])
	}
	if statements[1] != "INSERT INTO test_table (name, remark) VALUES ('beta', 'semi; inside')" {
		t.Fatalf("unexpected second statement: %s", statements[1])
	}
}
