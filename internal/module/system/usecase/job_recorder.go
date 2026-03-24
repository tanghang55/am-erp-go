package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"

	systemdomain "am-erp-go/internal/module/system/domain"
)

type JobRecorder struct {
	jobRepo systemdomain.JobRunRepository
	logRepo systemdomain.SystemLogRepository
}

type JobExecution struct {
	TraceID    string
	JobType    string
	JobName    string
	Module     string
	StartedAt  time.Time
	CreatedBy  *uint64
	TotalRows  *uint
	SuccessRows *uint
	FailedRows *uint
	Input      any
	Output     any
}

func NewJobRecorder(jobRepo systemdomain.JobRunRepository, logRepo systemdomain.SystemLogRepository) *JobRecorder {
	return &JobRecorder{jobRepo: jobRepo, logRepo: logRepo}
}

func (r *JobRecorder) Start(jobType string, jobName string, module string, input any, createdBy *uint64) *JobExecution {
	return &JobExecution{
		TraceID:   newJobTraceID(),
		JobType:   jobType,
		JobName:   jobName,
		Module:    module,
		StartedAt: time.Now(),
		CreatedBy: createdBy,
		Input:     input,
	}
}

func (r *JobRecorder) FinishSuccess(exec *JobExecution, output any) {
	if r == nil || exec == nil || r.jobRepo == nil {
		return
	}
	finishedAt := time.Now()
	duration := uint(finishedAt.Sub(exec.StartedAt).Milliseconds())
	run := &systemdomain.JobRun{
		TraceID:       exec.TraceID,
		JobType:       exec.JobType,
		JobName:       exec.JobName,
		Status:        "SUCCESS",
		StartedAt:     &exec.StartedAt,
		FinishedAt:    &finishedAt,
		DurationMs:    &duration,
		TotalRows:     exec.TotalRows,
		SuccessRows:   exec.SuccessRows,
		FailedRows:    exec.FailedRows,
		InputSummary:  marshalJSONString(exec.Input),
		OutputSummary: marshalJSONString(output),
		CreatedBy:     exec.CreatedBy,
	}
	_ = r.jobRepo.Create(run)
	r.log(exec.TraceID, "INFO", exec.Module, exec.JobName+" 执行成功", output, nil)
}

func (r *JobRecorder) FinishFailure(exec *JobExecution, err error, output any) {
	if r == nil || exec == nil || r.jobRepo == nil {
		return
	}
	finishedAt := time.Now()
	duration := uint(finishedAt.Sub(exec.StartedAt).Milliseconds())
	var errorMessage *string
	if err != nil {
		msg := err.Error()
		errorMessage = &msg
	}
	run := &systemdomain.JobRun{
		TraceID:       exec.TraceID,
		JobType:       exec.JobType,
		JobName:       exec.JobName,
		Status:        "FAILED",
		StartedAt:     &exec.StartedAt,
		FinishedAt:    &finishedAt,
		DurationMs:    &duration,
		TotalRows:     exec.TotalRows,
		SuccessRows:   exec.SuccessRows,
		FailedRows:    exec.FailedRows,
		InputSummary:  marshalJSONString(exec.Input),
		OutputSummary: marshalJSONString(output),
		ErrorMessage:  errorMessage,
		CreatedBy:     exec.CreatedBy,
	}
	_ = r.jobRepo.Create(run)
	r.log(exec.TraceID, "ERROR", exec.Module, exec.JobName+" 执行失败", output, errorMessage)
}

func (r *JobRecorder) log(traceID string, level string, module string, message string, context any, exception *string) {
	if r == nil || r.logRepo == nil {
		return
	}
	log := &systemdomain.SystemLog{
		TraceID:   traceID,
		Level:     level,
		Module:    module,
		Message:   message,
		Context:   marshalJSONString(context),
		Exception: exception,
	}
	_ = r.logRepo.Create(log)
}

func marshalJSONString(value any) *string {
	if value == nil {
		return nil
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	text := string(raw)
	return &text
}

func newJobTraceID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().Format("20060102150405.000")
	}
	return hex.EncodeToString(b[:])
}
