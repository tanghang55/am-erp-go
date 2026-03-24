package domain

import "time"

type FieldLabelRepository interface {
	GetAll() ([]*FieldLabel, error)
	List(page, pageSize int, keyword string) ([]*FieldLabel, int64, error)
	GetByID(id uint64) (*FieldLabel, error)
	GetByKey(key string) (*FieldLabel, error)
	Create(label *FieldLabel) error
	Update(label *FieldLabel) error
	Delete(id uint64) error
}

type ConfigCenterRepository interface {
	ListDefinitions(moduleCode string) ([]*ConfigDefinition, error)
	ListValues(scopeType string, scopeRefID uint64, keys []string) ([]*ConfigValue, error)
	UpsertValues(items []*ConfigValue) error
}

type AuditLogRepository interface {
	List(params AuditLogListParams) ([]*AuditLog, int64, error)
	Create(log *AuditLog) error
}

type MonitorRepository interface {
	Ping() error
	ListAppliedMigrationVersions() (map[string]struct{}, error)
	GetLatestReplenishmentRun() (*MonitorTaskSnapshot, error)
	GetLatestPackagingRun() (*MonitorTaskSnapshot, error)
	GetLatestOrderSyncRun() (*MonitorTaskSnapshot, error)
	GetLatestLogRetentionRun() (*MonitorTaskSnapshot, error)
	ListRecentJobRuns(status string, traceID string, limit int) ([]*MonitorRecentJob, error)
	ListRecentSystemLogs(level string, traceID string, limit int) ([]*MonitorRecentLog, error)
	CountExpiringAuthorizations(within time.Duration) (int, error)
}

type JobRunRepository interface {
	Create(run *JobRun) error
	Update(run *JobRun) error
	DeleteOlderThan(cutoff time.Time) (int64, error)
}

type SystemLogRepository interface {
	Create(log *SystemLog) error
	DeleteOlderThan(cutoff time.Time) (int64, error)
}
