package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsDotEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("DB_DATABASE=erp_test_db\n"), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	origCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origCwd)
	})

	origDb, hadDb := os.LookupEnv("DB_DATABASE")
	_ = os.Unsetenv("DB_DATABASE")
	t.Cleanup(func() {
		if hadDb {
			_ = os.Setenv("DB_DATABASE", origDb)
		} else {
			_ = os.Unsetenv("DB_DATABASE")
		}
	})

	cfg := Load()
	if cfg.Database.Database != "erp_test_db" {
		t.Fatalf("expected DB_DATABASE from .env, got %s", cfg.Database.Database)
	}
}
