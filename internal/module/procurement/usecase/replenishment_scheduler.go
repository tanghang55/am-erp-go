package usecase

import (
	"log"
	"time"

	"am-erp-go/internal/module/procurement/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"
)

type ReplenishmentScheduler struct {
	usecase      *ReplenishmentUsecase
	lastFinished time.Time
	jobRecorder  *systemUsecase.JobRecorder
}

func NewReplenishmentScheduler(usecase *ReplenishmentUsecase) *ReplenishmentScheduler {
	return &ReplenishmentScheduler{usecase: usecase}
}

func (s *ReplenishmentScheduler) BindJobRecorder(recorder *systemUsecase.JobRecorder) {
	s.jobRecorder = recorder
}

func (s *ReplenishmentScheduler) Start() {
	if s == nil || s.usecase == nil {
		return
	}
	go s.loop()
}

func (s *ReplenishmentScheduler) loop() {
	s.runTick()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.runTick()
	}
}

func (s *ReplenishmentScheduler) runTick() {
	exec := (*systemUsecase.JobExecution)(nil)
	if s.jobRecorder != nil {
		exec = s.jobRecorder.Start("REPLENISH_CALC", "采购计划自动生成", "Procurement", map[string]any{"trigger": "scheduler"}, nil)
	}
	if err := s.usecase.EnsureDailyCleanup(); err != nil {
		log.Printf("[replenishment-scheduler] cleanup failed: %v", err)
	}

	cfg, err := s.usecase.GetConfig()
	if err != nil {
		log.Printf("[replenishment-scheduler] load config failed: %v", err)
		if exec != nil {
			s.jobRecorder.FinishFailure(exec, err, map[string]any{"stage": "load_config"})
		}
		return
	}
	if cfg == nil || cfg.IsEnabled == 0 {
		if exec != nil {
			s.jobRecorder.FinishSuccess(exec, map[string]any{"skipped": true, "reason": "disabled"})
		}
		return
	}

	interval := time.Duration(clampU32(cfg.IntervalMinutes, 1, 10080, 1440)) * time.Minute
	now := time.Now()
	if !s.lastFinished.IsZero() && sameDay(s.lastFinished, now) && time.Since(s.lastFinished) < interval {
		if exec != nil {
			s.jobRecorder.FinishSuccess(exec, map[string]any{"skipped": true, "reason": "interval_not_reached"})
		}
		return
	}

	plans, generatedCount, err := s.usecase.GenerateDailyPlans(nil, &domain.ReplenishmentGenerateParams{
		TriggerType: domain.ReplenishmentTriggerScheduled,
	})
	if err != nil {
		log.Printf("[replenishment-scheduler] generate failed: %v", err)
		if exec != nil {
			total := uint(0)
			exec.TotalRows = &total
			s.jobRecorder.FinishFailure(exec, err, map[string]any{"stage": "generate"})
		}
		return
	}
	if generatedCount > 0 {
		log.Printf("[replenishment-scheduler] generated %d new plans, current %d plans", generatedCount, len(plans))
	}
	if exec != nil {
		total := uint(len(plans))
		exec.TotalRows = &total
		exec.SuccessRows = &total
		s.jobRecorder.FinishSuccess(exec, map[string]any{"generated": generatedCount > 0, "generated_count": generatedCount, "plan_count": len(plans)})
	}
	s.lastFinished = now
}
