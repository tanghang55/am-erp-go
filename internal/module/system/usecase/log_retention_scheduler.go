package usecase

import (
	"fmt"
	"log"
	"time"

	systemdomain "am-erp-go/internal/module/system/domain"
)

type LogRetentionScheduler struct {
	enabled         bool
	interval        time.Duration
	jobRetention    time.Duration
	systemRetention time.Duration
	jobRepo         systemdomain.JobRunRepository
	logRepo         systemdomain.SystemLogRepository
	recorder        *JobRecorder
}

func NewLogRetentionScheduler(
	enabled bool,
	intervalMinutes int,
	jobRetentionDays int,
	systemLogRetentionDays int,
	jobRepo systemdomain.JobRunRepository,
	logRepo systemdomain.SystemLogRepository,
) *LogRetentionScheduler {
	if intervalMinutes <= 0 {
		intervalMinutes = 1440
	}
	if jobRetentionDays <= 0 {
		jobRetentionDays = 30
	}
	if systemLogRetentionDays <= 0 {
		systemLogRetentionDays = 30
	}

	return &LogRetentionScheduler{
		enabled:         enabled,
		interval:        time.Duration(intervalMinutes) * time.Minute,
		jobRetention:    time.Duration(jobRetentionDays) * 24 * time.Hour,
		systemRetention: time.Duration(systemLogRetentionDays) * 24 * time.Hour,
		jobRepo:         jobRepo,
		logRepo:         logRepo,
	}
}

func (s *LogRetentionScheduler) BindJobRecorder(recorder *JobRecorder) {
	s.recorder = recorder
}

func (s *LogRetentionScheduler) Start() {
	if s == nil || !s.enabled || s.jobRepo == nil || s.logRepo == nil {
		return
	}
	go s.loop()
}

func (s *LogRetentionScheduler) loop() {
	s.tick()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for range ticker.C {
		s.tick()
	}
}

func (s *LogRetentionScheduler) tick() {
	exec := (*JobExecution)(nil)
	if s.recorder != nil {
		exec = s.recorder.Start("LOG_RETENTION_CLEANUP", "系统日志自动清理", "System", map[string]any{
			"trigger":                    "scheduler",
			"job_run_retention_days":    int(s.jobRetention.Hours() / 24),
			"system_log_retention_days": int(s.systemRetention.Hours() / 24),
		}, nil)
	}

	now := time.Now()
	jobCutoff := now.Add(-s.jobRetention)
	logCutoff := now.Add(-s.systemRetention)

	deletedJobRuns, deletedSystemLogs, err := s.cleanup(jobCutoff, logCutoff)
	if err != nil {
		stage := "delete_job_run"
		output := map[string]any{
			"job_run_cutoff":    jobCutoff.Format(time.DateTime),
			"system_log_cutoff": logCutoff.Format(time.DateTime),
		}
	if deletedJobRuns > 0 {
		output["deleted_job_runs"] = deletedJobRuns
		stage = "delete_system_log"
	}
	log.Printf("[log-retention] %s failed: %v", stage, err)
	if exec != nil {
		output["stage"] = stage
		s.recorder.FinishFailure(exec, err, output)
		}
		return
	}

	if exec != nil {
		total := uint(deletedJobRuns + deletedSystemLogs)
		exec.TotalRows = &total
		exec.SuccessRows = &total
		s.recorder.FinishSuccess(exec, map[string]any{
			"deleted_job_runs":    deletedJobRuns,
			"deleted_system_logs": deletedSystemLogs,
			"job_run_cutoff":      jobCutoff.Format(time.DateTime),
			"system_log_cutoff":   logCutoff.Format(time.DateTime),
		})
	}
}

func (s *LogRetentionScheduler) RunOnce() error {
	if s == nil || s.jobRepo == nil || s.logRepo == nil {
		return fmt.Errorf("log retention scheduler not configured")
	}
	_, _, err := s.cleanup(time.Now().Add(-s.jobRetention), time.Now().Add(-s.systemRetention))
	return err
}

func (s *LogRetentionScheduler) cleanup(jobCutoff time.Time, logCutoff time.Time) (int64, int64, error) {
	deletedJobRuns, err := s.jobRepo.DeleteOlderThan(jobCutoff)
	if err != nil {
		return 0, 0, err
	}
	deletedSystemLogs, err := s.logRepo.DeleteOlderThan(logCutoff)
	if err != nil {
		return deletedJobRuns, 0, err
	}
	return deletedJobRuns, deletedSystemLogs, nil
}
