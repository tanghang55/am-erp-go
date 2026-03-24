package usecase

import (
	"path/filepath"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/migration"
	systemdomain "am-erp-go/internal/module/system/domain"
)

type MonitorUsecase struct {
	repo         systemdomain.MonitorRepository
	migrationDir string
}

type MonitorTaskSnapshot = systemdomain.MonitorTaskSnapshot

func NewMonitorUsecase(repo systemdomain.MonitorRepository, migrationDir string) *MonitorUsecase {
	return &MonitorUsecase{
		repo:         repo,
		migrationDir: filepath.Clean(migrationDir),
	}
}

func (uc *MonitorUsecase) ListRecentJobs(status string, traceID string, limit int) ([]*systemdomain.MonitorRecentJob, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return uc.repo.ListRecentJobRuns(strings.ToUpper(strings.TrimSpace(status)), strings.TrimSpace(traceID), limit)
}

func (uc *MonitorUsecase) ListRecentLogs(level string, traceID string, limit int) ([]*systemdomain.MonitorRecentLog, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	return uc.repo.ListRecentSystemLogs(strings.ToUpper(strings.TrimSpace(level)), strings.TrimSpace(traceID), limit)
}

func (uc *MonitorUsecase) GetOverview() (*systemdomain.MonitorOverview, error) {
	now := time.Now()
	overview := &systemdomain.MonitorOverview{
		OverallStatus: systemdomain.MonitorStatusOK,
		GeneratedAt:   now,
		Database: systemdomain.MonitorComponent{
			Name:      "数据库连接",
			Status:    systemdomain.MonitorStatusOK,
			CheckedAt: now,
			Message:   "连接正常",
		},
		Migrations: systemdomain.MonitorComponent{
			Name:      "数据库迁移",
			Status:    systemdomain.MonitorStatusOK,
			CheckedAt: now,
			Message:   "无待执行迁移",
		},
		Tasks:  make([]systemdomain.MonitorTaskState, 0, 3),
		Alerts: make([]string, 0, 4),
	}

	if err := uc.repo.Ping(); err != nil {
		overview.Database.Status = systemdomain.MonitorStatusError
		overview.Database.Message = err.Error()
		overview.OverallStatus = systemdomain.MonitorStatusError
		overview.Alerts = append(overview.Alerts, "数据库连接失败")
	}

	applied, err := uc.repo.ListAppliedMigrationVersions()
	if err != nil {
		overview.Migrations.Status = systemdomain.MonitorStatusError
		overview.Migrations.Message = err.Error()
		overview.OverallStatus = systemdomain.MonitorStatusError
		overview.Alerts = append(overview.Alerts, "无法读取迁移执行记录")
	} else {
		pending, listErr := migration.ListPendingMigrations(uc.migrationDir, applied)
		if listErr != nil {
			overview.Migrations.Status = systemdomain.MonitorStatusError
			overview.Migrations.Message = listErr.Error()
			overview.OverallStatus = systemdomain.MonitorStatusError
			overview.Alerts = append(overview.Alerts, "无法读取本地迁移文件")
		} else if len(pending) > 0 {
			overview.Migrations.Status = systemdomain.MonitorStatusWarn
			overview.Migrations.PendingCount = len(pending)
			overview.Migrations.Message = "存在待执行迁移"
			overview.Alerts = append(overview.Alerts, "存在待执行迁移")
			if overview.OverallStatus == systemdomain.MonitorStatusOK {
				overview.OverallStatus = systemdomain.MonitorStatusWarn
			}
		} else {
			missing, missingErr := migration.ListMissingAppliedMigrations(uc.migrationDir, applied)
			if missingErr != nil {
				overview.Migrations.Status = systemdomain.MonitorStatusError
				overview.Migrations.Message = missingErr.Error()
				overview.OverallStatus = systemdomain.MonitorStatusError
				overview.Alerts = append(overview.Alerts, "无法校验历史迁移文件")
			} else if len(missing) > 0 {
				overview.Migrations.Status = systemdomain.MonitorStatusWarn
				overview.Migrations.PendingCount = len(missing)
				overview.Migrations.Message = "已执行迁移文件缺失"
				overview.Alerts = append(overview.Alerts, "历史迁移文件缺失")
				if overview.OverallStatus == systemdomain.MonitorStatusOK {
					overview.OverallStatus = systemdomain.MonitorStatusWarn
				}
			}
		}
	}

	taskLoaders := []struct {
		name   string
		loader func() (*systemdomain.MonitorTaskSnapshot, error)
	}{
		{name: "采购计划任务", loader: uc.repo.GetLatestReplenishmentRun},
		{name: "包材采购计划任务", loader: uc.repo.GetLatestPackagingRun},
		{name: "订单同步任务", loader: uc.repo.GetLatestOrderSyncRun},
		{name: "日志清理任务", loader: uc.repo.GetLatestLogRetentionRun},
	}
	for _, item := range taskLoaders {
		task, taskErr := item.loader()
		if taskErr != nil {
			overview.OverallStatus = systemdomain.MonitorStatusError
			overview.Alerts = append(overview.Alerts, "任务监控读取失败")
			continue
		}
		overview.Tasks = append(overview.Tasks, buildTaskState(item.name, task))
	}

	expiring, err := uc.repo.CountExpiringAuthorizations(24 * time.Hour)
	if err == nil && expiring > 0 {
		overview.Alerts = append(overview.Alerts, "存在即将过期的平台授权 token")
		if overview.OverallStatus == systemdomain.MonitorStatusOK {
			overview.OverallStatus = systemdomain.MonitorStatusWarn
		}
	}

	for _, task := range overview.Tasks {
		if task.Status == systemdomain.MonitorStatusError {
			overview.OverallStatus = systemdomain.MonitorStatusError
			overview.Alerts = append(overview.Alerts, task.Name+" 最近一次执行失败")
		} else if task.Status == systemdomain.MonitorStatusWarn && overview.OverallStatus == systemdomain.MonitorStatusOK {
			overview.OverallStatus = systemdomain.MonitorStatusWarn
		}
	}

	return overview, nil
}

func buildTaskState(defaultName string, task *systemdomain.MonitorTaskSnapshot) systemdomain.MonitorTaskState {
	if task == nil {
		return systemdomain.MonitorTaskState{
			Name:    defaultName,
			Status:  systemdomain.MonitorStatusWarn,
			Message: "暂无运行记录",
		}
	}

	state := systemdomain.MonitorTaskState{
		Name:       task.Name,
		LastStatus: task.Status,
		Message:    task.Message,
	}
	if task.FinishedAt != nil {
		state.LastRunAt = task.FinishedAt
	} else {
		state.LastRunAt = task.StartedAt
	}

	switch strings.ToUpper(strings.TrimSpace(task.Status)) {
	case "SUCCESS":
		state.Status = systemdomain.MonitorStatusOK
	case "FAILED":
		state.Status = systemdomain.MonitorStatusError
	case "RUNNING":
		state.Status = systemdomain.MonitorStatusWarn
	default:
		state.Status = systemdomain.MonitorStatusWarn
		if state.Message == "" {
			state.Message = "暂无运行记录"
		}
	}
	return state
}
