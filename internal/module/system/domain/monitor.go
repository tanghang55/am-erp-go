package domain

import "time"

type MonitorStatus string

const (
	MonitorStatusOK    MonitorStatus = "OK"
	MonitorStatusWarn  MonitorStatus = "WARN"
	MonitorStatusError MonitorStatus = "ERROR"
)

type MonitorComponent struct {
	Name         string        `json:"name"`
	Status       MonitorStatus `json:"status"`
	Message      string        `json:"message,omitempty"`
	CheckedAt    time.Time     `json:"checked_at"`
	PendingCount int           `json:"pending_count,omitempty"`
}

type MonitorTaskSnapshot struct {
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Message    string     `json:"message,omitempty"`
}

type MonitorOverview struct {
	OverallStatus MonitorStatus      `json:"overall_status"`
	GeneratedAt   time.Time          `json:"generated_at"`
	Database      MonitorComponent   `json:"database"`
	Migrations    MonitorComponent   `json:"migrations"`
	Tasks         []MonitorTaskState `json:"tasks"`
	Alerts        []string           `json:"alerts"`
}

type MonitorTaskState struct {
	Name       string        `json:"name"`
	Status     MonitorStatus `json:"status"`
	LastStatus string        `json:"last_status,omitempty"`
	LastRunAt  *time.Time    `json:"last_run_at,omitempty"`
	Message    string        `json:"message,omitempty"`
}

type MonitorRecentJob struct {
	ID           uint64     `json:"id"`
	TraceID      string     `json:"trace_id"`
	JobType      string     `json:"job_type"`
	JobName      string     `json:"job_name"`
	Status       string     `json:"status"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
	DurationMs   *uint      `json:"duration_ms,omitempty"`
	TotalRows    *uint      `json:"total_rows,omitempty"`
	SuccessRows  *uint      `json:"success_rows,omitempty"`
	FailedRows   *uint      `json:"failed_rows,omitempty"`
	ErrorMessage *string    `json:"error_message,omitempty"`
	GmtCreate    time.Time  `json:"gmt_create"`
}

type MonitorRecentLog struct {
	ID        uint64    `json:"id"`
	TraceID   string    `json:"trace_id"`
	Level     string    `json:"level"`
	Module    string    `json:"module"`
	Message   string    `json:"message"`
	Context   *string   `json:"context,omitempty"`
	Exception *string   `json:"exception,omitempty"`
	GmtCreate time.Time `json:"gmt_create"`
}
