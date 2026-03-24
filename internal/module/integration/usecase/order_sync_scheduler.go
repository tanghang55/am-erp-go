package usecase

import (
	"context"
	"log"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"
)

type OrderSyncScheduler struct {
	enabled      bool
	interval     time.Duration
	service      *OrderSyncService
	lastFinished time.Time
	recorder     *systemUsecase.JobRecorder
}

func NewOrderSyncScheduler(enabled bool, intervalMinutes int, service *OrderSyncService) *OrderSyncScheduler {
	if intervalMinutes <= 0 {
		intervalMinutes = 30
	}
	return &OrderSyncScheduler{
		enabled:  enabled,
		interval: time.Duration(intervalMinutes) * time.Minute,
		service:  service,
	}
}

func (s *OrderSyncScheduler) BindJobRecorder(recorder *systemUsecase.JobRecorder) {
	s.recorder = recorder
}

func (s *OrderSyncScheduler) Start() {
	if s == nil || !s.enabled || s.service == nil {
		return
	}
	go s.loop()
}

func (s *OrderSyncScheduler) loop() {
	s.tick()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.tick()
	}
}

func (s *OrderSyncScheduler) tick() {
	now := time.Now()
	if !s.lastFinished.IsZero() && time.Since(s.lastFinished) < s.interval {
		return
	}
	exec := (*systemUsecase.JobExecution)(nil)
	if s.recorder != nil {
		exec = s.recorder.Start("ORDER_SYNC", "第三方订单自动同步", "Integration", map[string]any{"trigger": "scheduler"}, nil)
	}
	if _, err := s.service.SyncOrders(context.Background(), integrationDomain.SyncTriggerScheduled, nil); err != nil {
		log.Printf("[third-party-order-sync] scheduled run failed: %v", err)
		if exec != nil {
			s.recorder.FinishFailure(exec, err, nil)
		}
		return
	}
	if exec != nil {
		s.recorder.FinishSuccess(exec, map[string]any{"synced": true})
	}
	s.lastFinished = now
}
