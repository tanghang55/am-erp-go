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

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Database.Database != "erp_test_db" {
		t.Fatalf("expected DB_DATABASE from .env, got %s", cfg.Database.Database)
	}
}

func TestLoadReturnsErrorWhenIntegrationConfigInvalid(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "integrations.json"), []byte("{invalid"), 0o644); err != nil {
		t.Fatalf("write integrations: %v", err)
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

	origConfigFile, hadConfigFile := os.LookupEnv("INTEGRATION_CONFIG_FILE")
	if err := os.Setenv("INTEGRATION_CONFIG_FILE", "integrations.json"); err != nil {
		t.Fatalf("set env: %v", err)
	}
	t.Cleanup(func() {
		if hadConfigFile {
			_ = os.Setenv("INTEGRATION_CONFIG_FILE", origConfigFile)
		} else {
			_ = os.Unsetenv("INTEGRATION_CONFIG_FILE")
		}
	})

	if _, err := Load(); err == nil {
		t.Fatalf("expected invalid integrations config error")
	}
}

func TestLoadUsesLogRetentionDefaults(t *testing.T) {
	dir := t.TempDir()

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

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.Operations.LogRetention.Enabled {
		t.Fatalf("expected log retention enabled by default")
	}
	if cfg.Operations.LogRetention.CleanupIntervalMinutes != 1440 {
		t.Fatalf("expected cleanup interval 1440, got %d", cfg.Operations.LogRetention.CleanupIntervalMinutes)
	}
	if cfg.Operations.LogRetention.JobRunRetentionDays != 30 {
		t.Fatalf("expected job retention 30, got %d", cfg.Operations.LogRetention.JobRunRetentionDays)
	}
	if cfg.Operations.LogRetention.SystemLogRetentionDays != 30 {
		t.Fatalf("expected system log retention 30, got %d", cfg.Operations.LogRetention.SystemLogRetentionDays)
	}
}
