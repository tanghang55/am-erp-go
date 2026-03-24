package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	integrationDomain "am-erp-go/internal/module/integration/domain"
)

type RefundProductResolver interface {
	ResolveProductIDBySellerSKU(sellerSKU string, marketplace string) (uint64, error)
}

type RefundSyncConfig struct {
	ProviderCode       string
	Channel            string
	MarketplaceIDs     []string
	LookbackMinutes    int
	InitialLookbackDay int
}

type RefundSyncResult struct {
	TriggerType     integrationDomain.SyncTrigger `json:"trigger_type"`
	RequestCursorAt time.Time                     `json:"request_cursor_at"`
	FetchedRefunds  uint32                        `json:"fetched_refunds"`
	ImportedRefunds uint32                        `json:"imported_refunds"`
	ErrorRefunds    uint32                        `json:"error_refunds"`
	FinishedAt      time.Time                     `json:"finished_at"`
}

type RefundSyncService struct {
	repo           integrationDomain.RefundSyncRepository
	provider       integrationDomain.RefundsProvider
	cfg            RefundSyncConfig
	nowFn          func() time.Time
	mappingResolve SKUMappingResolver
	productResolve RefundProductResolver
}

func NewRefundSyncService(
	repo integrationDomain.RefundSyncRepository,
	provider integrationDomain.RefundsProvider,
	cfg RefundSyncConfig,
	nowFn func() time.Time,
	mappingResolver SKUMappingResolver,
	productResolver RefundProductResolver,
) *RefundSyncService {
	if cfg.ProviderCode == "" && provider != nil {
		cfg.ProviderCode = provider.Code()
	}
	if cfg.ProviderCode == "" {
		cfg.ProviderCode = "UNKNOWN"
	}
	if cfg.Channel == "" {
		cfg.Channel = "DEFAULT"
	}
	if cfg.LookbackMinutes <= 0 {
		cfg.LookbackMinutes = 10
	}
	if cfg.InitialLookbackDay <= 0 {
		cfg.InitialLookbackDay = 7
	}
	if nowFn == nil {
		nowFn = time.Now
	}
	return &RefundSyncService{
		repo:           repo,
		provider:       provider,
		cfg:            cfg,
		nowFn:          nowFn,
		mappingResolve: mappingResolver,
		productResolve: productResolver,
	}
}

func (u *RefundSyncService) SyncRefunds(ctx context.Context, trigger integrationDomain.SyncTrigger, operatorID *uint64) (*RefundSyncResult, error) {
	_ = operatorID
	if u.repo == nil || u.provider == nil {
		return nil, fmt.Errorf("refund sync dependencies not configured")
	}
	if trigger == "" {
		trigger = integrationDomain.SyncTriggerManual
	}
	now := u.nowFn().UTC()

	state, err := u.repo.GetState(u.cfg.ProviderCode, u.cfg.Channel)
	if err != nil {
		return nil, err
	}
	cursor := now.AddDate(0, 0, -u.cfg.InitialLookbackDay)
	if state != nil && state.LastPostedAfter != nil && !state.LastPostedAfter.IsZero() {
		cursor = state.LastPostedAfter.Add(-time.Duration(u.cfg.LookbackMinutes) * time.Minute)
	}

	run := &integrationDomain.RefundSyncRun{
		Provider:           u.cfg.ProviderCode,
		Channel:            u.cfg.Channel,
		TriggerType:        trigger,
		Status:             integrationDomain.OrderSyncRunStatusRunning,
		RequestPostedAfter: timePtr(cursor),
		StartedAt:          now,
	}
	if err := u.repo.CreateRun(run); err != nil {
		return nil, err
	}

	refunds, err := u.provider.ListRefunds(ctx, integrationDomain.ExternalListRefundsRequest{
		MarketplaceIDs: u.cfg.MarketplaceIDs,
		PostedAfter:    cursor,
	})
	if err != nil {
		u.finishRunFailed(run, err)
		return nil, err
	}

	sort.Slice(refunds, func(i, j int) bool {
		return refunds[i].PostedAt.Before(refunds[j].PostedAt)
	})
	maxPostedAt := cursor
	events := make([]integrationDomain.ThirdPartyRefundEvent, 0, len(refunds))
	imported := uint32(0)
	errorCount := uint32(0)
	for _, refund := range refunds {
		if refund.PostedAt.After(maxPostedAt) {
			maxPostedAt = refund.PostedAt
		}
		event, eventErr := u.buildRefundEvent(refund)
		if eventErr != nil {
			errorCount++
			continue
		}
		events = append(events, event)
		if event.Status == integrationDomain.RefundEventStatusUnmapped {
			errorCount++
			continue
		}
		imported++
	}
	if err := u.repo.UpsertEvents(events); err != nil {
		u.finishRunFailed(run, err)
		return nil, err
	}

	finishedAt := u.nowFn().UTC()
	run.FetchedRefunds = uint32(len(refunds))
	run.ImportedRefunds = imported
	run.ErrorRefunds = errorCount
	run.Status = refundRunStatus(imported, errorCount)
	run.FinishedAt = timePtr(finishedAt)
	if errorCount > 0 {
		msg := fmt.Sprintf("%s refund sync completed with %d errors", strings.ToUpper(u.cfg.ProviderCode), errorCount)
		run.Message = &msg
	}
	if err := u.repo.UpdateRun(run); err != nil {
		return nil, err
	}

	if state == nil {
		state = &integrationDomain.RefundSyncState{
			Provider: u.cfg.ProviderCode,
			Channel:  u.cfg.Channel,
		}
	}
	if len(refunds) == 0 && maxPostedAt.Before(now) {
		maxPostedAt = now
	}
	state.LastPostedAfter = timePtr(maxPostedAt)
	state.LastSyncStarted = timePtr(now)
	state.LastSyncFinished = timePtr(finishedAt)
	if err := u.repo.SaveState(state); err != nil {
		return nil, err
	}

	return &RefundSyncResult{
		TriggerType:     trigger,
		RequestCursorAt: cursor,
		FetchedRefunds:  uint32(len(refunds)),
		ImportedRefunds: imported,
		ErrorRefunds:    errorCount,
		FinishedAt:      finishedAt,
	}, nil
}

func (u *RefundSyncService) GetState() (*integrationDomain.RefundSyncState, error) {
	return u.repo.GetState(u.cfg.ProviderCode, u.cfg.Channel)
}

func (u *RefundSyncService) ListRuns(page int, pageSize int) ([]integrationDomain.RefundSyncRun, int64, error) {
	return u.repo.ListRuns(u.cfg.ProviderCode, u.cfg.Channel, &integrationDomain.ListRunsParams{
		Page:     page,
		PageSize: pageSize,
	})
}

func (u *RefundSyncService) buildRefundEvent(refund integrationDomain.ExternalRefund) (integrationDomain.ThirdPartyRefundEvent, error) {
	refundID := strings.TrimSpace(refund.RefundID)
	if refundID == "" {
		return integrationDomain.ThirdPartyRefundEvent{}, fmt.Errorf("refund_id is required")
	}
	orderID := strings.TrimSpace(refund.OrderID)
	if orderID == "" {
		return integrationDomain.ThirdPartyRefundEvent{}, fmt.Errorf("order_id is required")
	}
	sellerSKU := strings.TrimSpace(refund.SellerSKU)
	if sellerSKU == "" {
		return integrationDomain.ThirdPartyRefundEvent{}, fmt.Errorf("seller_sku is required")
	}
	qty := refund.Quantity
	if qty == 0 {
		return integrationDomain.ThirdPartyRefundEvent{}, fmt.Errorf("refund quantity must be greater than 0")
	}
	postedAt := refund.PostedAt
	if postedAt.IsZero() {
		postedAt = u.nowFn().UTC()
	}
	marketplace := normalizeMarketplace(refund.MarketplaceID)
	currency := strings.ToUpper(strings.TrimSpace(refund.Currency))
	if currency == "" {
		currency = "USD"
	}

	event := integrationDomain.ThirdPartyRefundEvent{
		Provider:     u.cfg.ProviderCode,
		Channel:      u.cfg.Channel,
		RefundID:     refundID,
		OrderID:      orderID,
		OrderItemID:  strPtr(strings.TrimSpace(refund.OrderItemID)),
		SellerSKU:    sellerSKU,
		Marketplace:  marketplace,
		QtyRefunded:  qty,
		RefundAmount: refund.Amount,
		Currency:     currency,
		PostedAt:     postedAt.UTC(),
		Status:       integrationDomain.RefundEventStatusMapped,
	}
	raw := map[string]any{
		"refund_id":       refundID,
		"order_id":        orderID,
		"order_item_id":   strings.TrimSpace(refund.OrderItemID),
		"seller_sku":      sellerSKU,
		"marketplace":     marketplace,
		"qty_refunded":    qty,
		"refund_amount":   refund.Amount,
		"currency":        currency,
		"posted_at":       postedAt.UTC().Format(time.RFC3339),
		"marketplace_raw": refund.MarketplaceID,
	}
	if encoded, err := json.Marshal(raw); err == nil {
		value := string(encoded)
		event.RawPayload = &value
	}

	productID, err := u.resolveProductID(sellerSKU, marketplace)
	if err != nil {
		msg := err.Error()
		event.Status = integrationDomain.RefundEventStatusUnmapped
		event.ErrorMessage = &msg
		return event, nil
	}
	if productID == 0 {
		msg := fmt.Sprintf("seller_sku not found: %s", sellerSKU)
		event.Status = integrationDomain.RefundEventStatusUnmapped
		event.ErrorMessage = &msg
		return event, nil
	}
	event.ProductID = &productID
	return event, nil
}

func (u *RefundSyncService) resolveProductID(sellerSKU string, marketplace string) (uint64, error) {
	if u.productResolve == nil {
		return 0, nil
	}
	if u.mappingResolve != nil {
		productID, err := u.mappingResolve.ResolveMappedProductID(u.cfg.ProviderCode, sellerSKU, marketplace)
		if err != nil {
			return 0, err
		}
		if productID > 0 {
			return productID, nil
		}
	}
	return u.productResolve.ResolveProductIDBySellerSKU(sellerSKU, marketplace)
}

func (u *RefundSyncService) finishRunFailed(run *integrationDomain.RefundSyncRun, err error) {
	if run == nil {
		return
	}
	msg := err.Error()
	finishedAt := u.nowFn().UTC()
	run.Status = integrationDomain.OrderSyncRunStatusFailed
	run.Message = &msg
	run.FinishedAt = timePtr(finishedAt)
	_ = u.repo.UpdateRun(run)
}

func refundRunStatus(imported uint32, failed uint32) integrationDomain.OrderSyncRunStatus {
	switch {
	case imported > 0 && failed == 0:
		return integrationDomain.OrderSyncRunStatusSuccess
	case imported == 0 && failed > 0:
		return integrationDomain.OrderSyncRunStatusFailed
	case imported > 0 && failed > 0:
		return integrationDomain.OrderSyncRunStatusPartial
	default:
		return integrationDomain.OrderSyncRunStatusSuccess
	}
}
