package repository

import (
	"time"

	systemdomain "am-erp-go/internal/module/system/domain"

	"gorm.io/gorm"
)

type monitorRepository struct {
	db *gorm.DB
}

func NewMonitorRepository(db *gorm.DB) systemdomain.MonitorRepository {
	return &monitorRepository{db: db}
}

func (r *monitorRepository) Ping() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (r *monitorRepository) ListAppliedMigrationVersions() (map[string]struct{}, error) {
	rows, err := r.db.Raw("SELECT version FROM schema_migration").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]struct{})
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		result[version] = struct{}{}
	}
	return result, nil
}

func (r *monitorRepository) GetLatestReplenishmentRun() (*systemdomain.MonitorTaskSnapshot, error) {
	type row struct {
		Status     string
		StartedAt  *time.Time
		FinishedAt *time.Time
		Message    *string
	}
	var item row
	if err := r.db.Table("procurement_replenishment_run").
		Select("status, started_at, finished_at, error_message AS message").
		Order("id DESC").
		Limit(1).
		Take(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &systemdomain.MonitorTaskSnapshot{
		Name:       "采购计划任务",
		Status:     item.Status,
		StartedAt:  item.StartedAt,
		FinishedAt: item.FinishedAt,
		Message:    valueOrEmpty(item.Message),
	}, nil
}

func (r *monitorRepository) GetLatestPackagingRun() (*systemdomain.MonitorTaskSnapshot, error) {
	type row struct {
		Status     string
		StartedAt  *time.Time
		FinishedAt *time.Time
		Message    *string
	}
	var item row
	if err := r.db.Table("packaging_procurement_run").
		Select("status, started_at, finished_at, error_message AS message").
		Order("id DESC").
		Limit(1).
		Take(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &systemdomain.MonitorTaskSnapshot{
		Name:       "包材采购计划任务",
		Status:     item.Status,
		StartedAt:  item.StartedAt,
		FinishedAt: item.FinishedAt,
		Message:    valueOrEmpty(item.Message),
	}, nil
}

func (r *monitorRepository) GetLatestOrderSyncRun() (*systemdomain.MonitorTaskSnapshot, error) {
	type row struct {
		Status     string
		StartedAt  *time.Time
		FinishedAt *time.Time
		Message    *string
	}
	var item row
	if err := r.db.Table("third_party_order_sync_run").
		Select("status, started_at, finished_at, message").
		Order("id DESC").
		Limit(1).
		Take(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &systemdomain.MonitorTaskSnapshot{
		Name:       "订单同步任务",
		Status:     item.Status,
		StartedAt:  item.StartedAt,
		FinishedAt: item.FinishedAt,
		Message:    valueOrEmpty(item.Message),
	}, nil
}

func (r *monitorRepository) GetLatestLogRetentionRun() (*systemdomain.MonitorTaskSnapshot, error) {
	type row struct {
		Status     string
		StartedAt  *time.Time
		FinishedAt *time.Time
		Message    *string
	}
	var item row
	if err := r.db.Table("job_run").
		Select("status, started_at, finished_at, error_message AS message").
		Where("job_type = ?", "LOG_RETENTION_CLEANUP").
		Order("id DESC").
		Limit(1).
		Take(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &systemdomain.MonitorTaskSnapshot{
		Name:       "日志清理任务",
		Status:     item.Status,
		StartedAt:  item.StartedAt,
		FinishedAt: item.FinishedAt,
		Message:    valueOrEmpty(item.Message),
	}, nil
}

func (r *monitorRepository) ListRecentJobRuns(status string, traceID string, limit int) ([]*systemdomain.MonitorRecentJob, error) {
	if limit <= 0 {
		limit = 10
	}

	query := r.db.Table("job_run").
		Select("id, trace_id, job_type, job_name, status, started_at, finished_at, duration_ms, total_rows, success_rows, failed_rows, error_message, gmt_create")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if traceID != "" {
		query = query.Where("trace_id = ?", traceID)
	}

	items := make([]*systemdomain.MonitorRecentJob, 0, limit)
	if err := query.Order("id DESC").Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *monitorRepository) ListRecentSystemLogs(level string, traceID string, limit int) ([]*systemdomain.MonitorRecentLog, error) {
	if limit <= 0 {
		limit = 10
	}

	query := r.db.Table("system_log").
		Select("id, trace_id, level, module, message, context, exception, gmt_create")

	if level != "" {
		query = query.Where("level = ?", level)
	}
	if traceID != "" {
		query = query.Where("trace_id = ?", traceID)
	}

	items := make([]*systemdomain.MonitorRecentLog, 0, limit)
	if err := query.Order("id DESC").Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *monitorRepository) CountExpiringAuthorizations(within time.Duration) (int, error) {
	deadline := time.Now().Add(within)
	var count int64
	err := r.db.Table("integration_authorization").
		Where("status = ?", "AUTHORIZED").
		Where("access_token_expire_at IS NOT NULL").
		Where("access_token_expire_at <= ?", deadline).
		Count(&count).Error
	return int(count), err
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
