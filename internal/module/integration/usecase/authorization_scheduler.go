package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	systemUsecase "am-erp-go/internal/module/system/usecase"
)

type AuthorizationRefreshScheduler struct {
	enabled  bool
	interval time.Duration
	usecase  *AuthorizationUsecase
	recorder *systemUsecase.JobRecorder
}

func NewAuthorizationRefreshScheduler(enabled bool, intervalMinutes int, usecase *AuthorizationUsecase) *AuthorizationRefreshScheduler {
	if intervalMinutes <= 0 {
		intervalMinutes = 5
	}
	return &AuthorizationRefreshScheduler{
		enabled:  enabled,
		interval: time.Duration(intervalMinutes) * time.Minute,
		usecase:  usecase,
	}
}

func (s *AuthorizationRefreshScheduler) BindJobRecorder(recorder *systemUsecase.JobRecorder) {
	s.recorder = recorder
}

func (s *AuthorizationRefreshScheduler) Start() {
	if s == nil || !s.enabled || s.usecase == nil {
		return
	}
	go s.loop()
}

func (s *AuthorizationRefreshScheduler) loop() {
	s.tick()
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for range ticker.C {
		s.tick()
	}
}

func (s *AuthorizationRefreshScheduler) tick() {
	exec := (*systemUsecase.JobExecution)(nil)
	if s.recorder != nil {
		exec = s.recorder.Start("AUTH_REFRESH", "平台授权自动刷新", "Integration", map[string]any{"trigger": "scheduler"}, nil)
	}
	result, err := s.usecase.RefreshExpiringTokens(context.Background())
	if err != nil {
		log.Printf("[integration-auth] refresh scheduler failed: %v", err)
		if exec != nil {
			s.recorder.FinishFailure(exec, err, nil)
		}
		return
	}
	if result != nil && result.Failed > 0 {
		err := fmt.Errorf("authorization refresh has %d failures", result.Failed)
		log.Printf("[integration-auth] refresh scheduler warning: %v", err)
		if exec != nil {
			s.recorder.FinishFailure(exec, err, map[string]any{
				"total":   result.Total,
				"success": result.Success,
				"failed":  result.Failed,
			})
		}
		return
	}
	if exec != nil {
		output := map[string]any{"refreshed": true}
		if result != nil {
			output["total"] = result.Total
			output["success"] = result.Success
			output["failed"] = result.Failed
		}
		s.recorder.FinishSuccess(exec, output)
	}
}
