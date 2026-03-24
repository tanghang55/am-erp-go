package usecase

import (
	"errors"
	"testing"
	"time"

	systemdomain "am-erp-go/internal/module/system/domain"
)

type stubRetentionJobRepo struct {
	cutoff      time.Time
	deletedRows int64
	err         error
}

func (s *stubRetentionJobRepo) Create(run *systemdomain.JobRun) error { return nil }
func (s *stubRetentionJobRepo) Update(run *systemdomain.JobRun) error { return nil }
func (s *stubRetentionJobRepo) DeleteOlderThan(cutoff time.Time) (int64, error) {
	s.cutoff = cutoff
	return s.deletedRows, s.err
}

type stubRetentionLogRepo struct {
	cutoff      time.Time
	deletedRows int64
	err         error
}

func (s *stubRetentionLogRepo) Create(log *systemdomain.SystemLog) error { return nil }
func (s *stubRetentionLogRepo) DeleteOlderThan(cutoff time.Time) (int64, error) {
	s.cutoff = cutoff
	return s.deletedRows, s.err
}

func TestLogRetentionRunOnceDeletesByConfiguredRetention(t *testing.T) {
	jobRepo := &stubRetentionJobRepo{deletedRows: 5}
	logRepo := &stubRetentionLogRepo{deletedRows: 7}
	scheduler := NewLogRetentionScheduler(true, 60, 30, 15, jobRepo, logRepo)

	before := time.Now()
	if err := scheduler.RunOnce(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if jobRepo.cutoff.IsZero() || logRepo.cutoff.IsZero() {
		t.Fatalf("expected cutoffs to be recorded")
	}

	jobDelta := before.Sub(jobRepo.cutoff)
	if jobDelta < (29*24*time.Hour) || jobDelta > (31*24*time.Hour) {
		t.Fatalf("unexpected job retention cutoff delta: %v", jobDelta)
	}

	logDelta := before.Sub(logRepo.cutoff)
	if logDelta < (14*24*time.Hour) || logDelta > (16*24*time.Hour) {
		t.Fatalf("unexpected system log retention cutoff delta: %v", logDelta)
	}
}

func TestLogRetentionRunOnceReturnsJobDeleteError(t *testing.T) {
	jobRepo := &stubRetentionJobRepo{err: errors.New("job delete failed")}
	logRepo := &stubRetentionLogRepo{}
	scheduler := NewLogRetentionScheduler(true, 60, 30, 30, jobRepo, logRepo)

	if err := scheduler.RunOnce(); err == nil || err.Error() != "job delete failed" {
		t.Fatalf("expected job delete error, got %v", err)
	}
}

func TestLogRetentionRunOnceReturnsSystemLogDeleteError(t *testing.T) {
	jobRepo := &stubRetentionJobRepo{deletedRows: 3}
	logRepo := &stubRetentionLogRepo{err: errors.New("system log delete failed")}
	scheduler := NewLogRetentionScheduler(true, 60, 30, 30, jobRepo, logRepo)

	if err := scheduler.RunOnce(); err == nil || err.Error() != "system log delete failed" {
		t.Fatalf("expected system log delete error, got %v", err)
	}
}
