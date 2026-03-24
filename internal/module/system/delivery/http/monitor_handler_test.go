package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	systemdomain "am-erp-go/internal/module/system/domain"

	"github.com/gin-gonic/gin"
)

type stubMonitorUsecase struct {
	overview *systemdomain.MonitorOverview
	jobs     []*systemdomain.MonitorRecentJob
	logs     []*systemdomain.MonitorRecentLog
	err      error
	jobErr   error
	logErr   error
	jobTraceID string
	logTraceID string
	jobStatus  string
	logLevel   string
	jobLimit   int
	logLimit   int
}

func (s *stubMonitorUsecase) GetOverview() (*systemdomain.MonitorOverview, error) {
	return s.overview, s.err
}

func (s *stubMonitorUsecase) ListRecentJobs(status string, traceID string, limit int) ([]*systemdomain.MonitorRecentJob, error) {
	s.jobStatus = status
	s.jobTraceID = traceID
	s.jobLimit = limit
	return s.jobs, s.jobErr
}

func (s *stubMonitorUsecase) ListRecentLogs(level string, traceID string, limit int) ([]*systemdomain.MonitorRecentLog, error) {
	s.logLevel = level
	s.logTraceID = traceID
	s.logLimit = limit
	return s.logs, s.logErr
}

func TestMonitorOverviewReturnsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	handler := NewMonitorHandler(&stubMonitorUsecase{
		overview: &systemdomain.MonitorOverview{
			OverallStatus: systemdomain.MonitorStatusOK,
			GeneratedAt:   now,
		},
	})

	router := gin.New()
	router.GET("/api/v1/system/monitor/overview", handler.GetOverview)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/system/monitor/overview", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestMonitorRecentJobsReturnsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	handler := NewMonitorHandler(&stubMonitorUsecase{
		jobs: []*systemdomain.MonitorRecentJob{
			{
				ID:         1,
				JobType:    "AUTH_REFRESH",
				JobName:    "平台授权自动刷新",
				Status:     "SUCCESS",
				StartedAt:  &now,
				FinishedAt: &now,
			},
		},
	})

	router := gin.New()
	router.GET("/api/v1/system/monitor/jobs", handler.ListRecentJobs)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/system/monitor/jobs?status=SUCCESS&limit=5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Code int `json:"code"`
		Data []struct {
			JobType string `json:"job_type"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}
	if len(resp.Data) != 1 || resp.Data[0].JobType != "AUTH_REFRESH" {
		t.Fatalf("unexpected response data: %s", w.Body.String())
	}
}

func TestMonitorRecentJobsRejectsInvalidLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewMonitorHandler(&stubMonitorUsecase{})

	router := gin.New()
	router.GET("/api/v1/system/monitor/jobs", handler.ListRecentJobs)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/system/monitor/jobs?limit=abc", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestMonitorRecentLogsReturnsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	usecase := &stubMonitorUsecase{
		logs: []*systemdomain.MonitorRecentLog{
			{
				ID:        1,
				TraceID:   "trace-1",
				Level:     "ERROR",
				Module:    "AUTH_REFRESH",
				Message:   "token refresh failed",
				GmtCreate: now,
			},
		},
	}
	handler := NewMonitorHandler(usecase)

	router := gin.New()
	router.GET("/api/v1/system/monitor/logs", handler.ListRecentLogs)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/system/monitor/logs?level=ERROR&trace_id=trace-1&limit=5", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Code int `json:"code"`
		Data []struct {
			Level string `json:"level"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}
	if len(resp.Data) != 1 || resp.Data[0].Level != "ERROR" {
		t.Fatalf("unexpected response data: %s", w.Body.String())
	}
	if usecase.logLevel != "ERROR" || usecase.logTraceID != "trace-1" || usecase.logLimit != 5 {
		t.Fatalf("unexpected log filters: level=%s trace=%s limit=%d", usecase.logLevel, usecase.logTraceID, usecase.logLimit)
	}
}
