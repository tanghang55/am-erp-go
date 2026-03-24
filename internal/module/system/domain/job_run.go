package domain

import "time"

type JobRun struct {
	ID            uint64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID       string     `json:"trace_id" gorm:"column:trace_id"`
	JobType       string     `json:"job_type" gorm:"column:job_type"`
	JobName       string     `json:"job_name" gorm:"column:job_name"`
	Status        string     `json:"status" gorm:"column:status"`
	StartedAt     *time.Time `json:"started_at" gorm:"column:started_at"`
	FinishedAt    *time.Time `json:"finished_at" gorm:"column:finished_at"`
	DurationMs    *uint      `json:"duration_ms" gorm:"column:duration_ms"`
	TotalRows     *uint      `json:"total_rows" gorm:"column:total_rows"`
	SuccessRows   *uint      `json:"success_rows" gorm:"column:success_rows"`
	FailedRows    *uint      `json:"failed_rows" gorm:"column:failed_rows"`
	InputSummary  *string    `json:"input_summary" gorm:"column:input_summary"`
	OutputSummary *string    `json:"output_summary" gorm:"column:output_summary"`
	ErrorMessage  *string    `json:"error_message" gorm:"column:error_message"`
	CreatedBy     *uint64    `json:"created_by" gorm:"column:created_by"`
	GmtCreate     time.Time  `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified   time.Time  `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (JobRun) TableName() string {
	return "job_run"
}

type SystemLog struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	TraceID     string    `json:"trace_id" gorm:"column:trace_id"`
	Level       string    `json:"level" gorm:"column:level"`
	Module      string    `json:"module" gorm:"column:module"`
	Message     string    `json:"message" gorm:"column:message"`
	Context     *string   `json:"context" gorm:"column:context"`
	Exception   *string   `json:"exception" gorm:"column:exception"`
	File        *string   `json:"file" gorm:"column:file"`
	Line        *uint     `json:"line" gorm:"column:line"`
	GmtCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create;autoCreateTime"`
	GmtModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified;autoUpdateTime"`
}

func (SystemLog) TableName() string {
	return "system_log"
}
