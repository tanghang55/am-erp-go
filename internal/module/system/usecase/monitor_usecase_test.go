package usecase

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	systemdomain "am-erp-go/internal/module/system/domain"
)

type stubMonitorRepo struct {
	pingErr                    error
	appliedVersions            map[string]struct{}
	replenishmentRun           *MonitorTaskSnapshot
	packagingRun               *MonitorTaskSnapshot
	orderSyncRun               *MonitorTaskSnapshot
	logRetentionRun            *MonitorTaskSnapshot
	recentJobs                 []*systemdomain.MonitorRecentJob
	recentJobsFilterStatus     string
	recentJobsFilterTraceID    string
	recentJobsFilterLimit      int
	recentLogs                 []*systemdomain.MonitorRecentLog
	recentLogsFilterLevel      string
	recentLogsFilterTraceID    string
	recentLogsFilterLimit      int
	expiringAuthorizationCount int
}

func (s *stubMonitorRepo) Ping() error { return s.pingErr }
func (s *stubMonitorRepo) ListAppliedMigrationVersions() (map[string]struct{}, error) {
	return s.appliedVersions, nil
}
func (s *stubMonitorRepo) GetLatestReplenishmentRun() (*MonitorTaskSnapshot, error) {
	return s.replenishmentRun, nil
}
func (s *stubMonitorRepo) GetLatestPackagingRun() (*MonitorTaskSnapshot, error) {
	return s.packagingRun, nil
}
func (s *stubMonitorRepo) GetLatestOrderSyncRun() (*MonitorTaskSnapshot, error) {
	return s.orderSyncRun, nil
}
func (s *stubMonitorRepo) GetLatestLogRetentionRun() (*MonitorTaskSnapshot, error) {
	return s.logRetentionRun, nil
}
func (s *stubMonitorRepo) ListRecentJobRuns(status string, traceID string, limit int) ([]*systemdomain.MonitorRecentJob, error) {
	s.recentJobsFilterStatus = status
	s.recentJobsFilterTraceID = traceID
	s.recentJobsFilterLimit = limit
	return s.recentJobs, nil
}
func (s *stubMonitorRepo) ListRecentSystemLogs(level string, traceID string, limit int) ([]*systemdomain.MonitorRecentLog, error) {
	s.recentLogsFilterLevel = level
	s.recentLogsFilterTraceID = traceID
	s.recentLogsFilterLimit = limit
	return s.recentLogs, nil
}
func (s *stubMonitorRepo) CountExpiringAuthorizations(within time.Duration) (int, error) {
	return s.expiringAuthorizationCount, nil
}

func TestMonitorOverviewReportsPendingMigrationsAndTaskFailure(t *testing.T) {
	migrationDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(migrationDir, "20260399_pending.sql"), []byte("-- test"), 0o644); err != nil {
		t.Fatalf("write migration file: %v", err)
	}

	now := time.Now()
	repo := &stubMonitorRepo{
		appliedVersions: map[string]struct{}{
			"20260302_sales_order.sql": {},
		},
		replenishmentRun: &MonitorTaskSnapshot{
			Name:       "采购计划",
			Status:     "FAILED",
			StartedAt:  &now,
			FinishedAt: &now,
			Message:    "db timeout",
		},
		packagingRun: &MonitorTaskSnapshot{
			Name:       "包材采购计划",
			Status:     "SUCCESS",
			StartedAt:  &now,
			FinishedAt: &now,
		},
		orderSyncRun: &MonitorTaskSnapshot{
			Name:       "订单同步",
			Status:     "SUCCESS",
			StartedAt:  &now,
			FinishedAt: &now,
		},
		expiringAuthorizationCount: 2,
	}

	uc := NewMonitorUsecase(repo, migrationDir)
	overview, err := uc.GetOverview()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if overview.Database.Status != "OK" {
		t.Fatalf("expected database OK, got %s", overview.Database.Status)
	}
	if overview.Migrations.PendingCount < 1 {
		t.Fatalf("expected pending migrations, got %d", overview.Migrations.PendingCount)
	}
	if overview.OverallStatus != "ERROR" {
		t.Fatalf("expected overall ERROR, got %s", overview.OverallStatus)
	}
	if len(overview.Alerts) < 2 {
		t.Fatalf("expected alerts, got %#v", overview.Alerts)
	}
}

func TestMonitorOverviewWarnsWhenAppliedMigrationFilesAreMissing(t *testing.T) {
	migrationDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(migrationDir, "20260302_sales_order.sql"), []byte("-- test"), 0o644); err != nil {
		t.Fatalf("write migration file: %v", err)
	}

	repo := &stubMonitorRepo{
		appliedVersions: map[string]struct{}{
			"20260302_sales_order.sql": {},
			"20260303_inventory_lot.sql": {},
		},
	}

	uc := NewMonitorUsecase(repo, migrationDir)
	overview, err := uc.GetOverview()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.Migrations.Status != systemdomain.MonitorStatusWarn {
		t.Fatalf("expected migration status WARN, got %s", overview.Migrations.Status)
	}
	if overview.Migrations.PendingCount != 1 {
		t.Fatalf("expected missing migration count 1, got %d", overview.Migrations.PendingCount)
	}
	if overview.Migrations.Message == "" {
		t.Fatalf("expected migration warning message")
	}
	if overview.OverallStatus != systemdomain.MonitorStatusWarn {
		t.Fatalf("expected overall WARN, got %s", overview.OverallStatus)
	}
}

func TestMonitorOverviewReportsDatabaseFailure(t *testing.T) {
	migrationDir := t.TempDir()
	repo := &stubMonitorRepo{
		pingErr:         errors.New("connection refused"),
		appliedVersions: map[string]struct{}{},
	}

	uc := NewMonitorUsecase(repo, migrationDir)
	overview, err := uc.GetOverview()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overview.Database.Status != "ERROR" {
		t.Fatalf("expected database ERROR, got %s", overview.Database.Status)
	}
	if overview.OverallStatus != "ERROR" {
		t.Fatalf("expected overall ERROR, got %s", overview.OverallStatus)
	}
}

func TestMonitorRecentJobsUsesStatusAndLimit(t *testing.T) {
	now := time.Now()
	repo := &stubMonitorRepo{
		recentJobs: []*systemdomain.MonitorRecentJob{
			{
				ID:          3,
				JobType:     "REPLENISH_CALC",
				JobName:     "采购计划自动生成",
				Status:      "SUCCESS",
				StartedAt:   &now,
				FinishedAt:  &now,
				DurationMs:  uintPtr(120),
				GmtCreate:   now,
			},
		},
	}

	uc := NewMonitorUsecase(repo, t.TempDir())
	list, err := uc.ListRecentJobs("SUCCESS", "trace-job-1", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.recentJobsFilterStatus != "SUCCESS" {
		t.Fatalf("expected status filter SUCCESS, got %s", repo.recentJobsFilterStatus)
	}
	if repo.recentJobsFilterLimit != 5 {
		t.Fatalf("expected limit 5, got %d", repo.recentJobsFilterLimit)
	}
	if repo.recentJobsFilterTraceID != "trace-job-1" {
		t.Fatalf("expected trace_id filter trace-job-1, got %s", repo.recentJobsFilterTraceID)
	}
	if len(list) != 1 || list[0].JobType != "REPLENISH_CALC" {
		t.Fatalf("unexpected list: %#v", list)
	}
}

func TestMonitorRecentJobsUsesDefaultLimitWhenInvalid(t *testing.T) {
	repo := &stubMonitorRepo{}
	uc := NewMonitorUsecase(repo, t.TempDir())

	if _, err := uc.ListRecentJobs("", "", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.recentJobsFilterLimit != 10 {
		t.Fatalf("expected default limit 10, got %d", repo.recentJobsFilterLimit)
	}
}

func TestMonitorRecentLogsUsesLevelAndLimit(t *testing.T) {
	now := time.Now()
	repo := &stubMonitorRepo{
		recentLogs: []*systemdomain.MonitorRecentLog{
			{
				ID:        1,
				Level:     "ERROR",
				Module:    "AUTH_REFRESH",
				Message:   "refresh failed",
				GmtCreate: now,
			},
		},
	}

	uc := NewMonitorUsecase(repo, t.TempDir())
	list, err := uc.ListRecentLogs("error", "trace-log-1", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.recentLogsFilterLevel != "ERROR" {
		t.Fatalf("expected level ERROR, got %s", repo.recentLogsFilterLevel)
	}
	if repo.recentLogsFilterLimit != 5 {
		t.Fatalf("expected limit 5, got %d", repo.recentLogsFilterLimit)
	}
	if repo.recentLogsFilterTraceID != "trace-log-1" {
		t.Fatalf("expected trace_id trace-log-1, got %s", repo.recentLogsFilterTraceID)
	}
	if len(list) != 1 || list[0].Module != "AUTH_REFRESH" {
		t.Fatalf("unexpected list: %#v", list)
	}
}

func uintPtr(v uint) *uint {
	return &v
}
