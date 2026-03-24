package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/numbering"
	integrationDomain "am-erp-go/internal/module/integration/domain"
	salesDomain "am-erp-go/internal/module/sales/domain"
)

type OrderImportGateway interface {
	CreateImport(batch *salesDomain.ReportImport) error
	UpdateImport(batch *salesDomain.ReportImport) error
	InsertImportRowErrors(errors []salesDomain.ReportImportRowError) error
	ResolveProductIDBySellerSKU(sellerSKU string, marketplace string) (uint64, error)
	UpsertImportedOrderLine(line *salesDomain.ImportOrderLine, productID uint64, batchNo string, operatorID *uint64) error
}

type OrderSyncConfig struct {
	ProviderCode       string
	Channel            string
	SourceType         string
	SalesChannel       string
	DefaultCurrency    string
	MarketplaceIDs     []string
	LookbackMinutes    int
	InitialLookbackDay int
}

type OrderSyncResult struct {
	TriggerType     integrationDomain.SyncTrigger `json:"trigger_type"`
	RequestCursorAt time.Time                     `json:"request_cursor_at"`
	FetchedOrders   uint32                        `json:"fetched_orders"`
	FetchedLines    uint32                        `json:"fetched_lines"`
	ImportedLines   uint32                        `json:"imported_lines"`
	ErrorLines      uint32                        `json:"error_lines"`
	BatchNo         string                        `json:"batch_no"`
	FinishedAt      time.Time                     `json:"finished_at"`
}

type OrderSyncCurrencyProvider interface {
	GetDefaultBaseCurrency() string
}

type SKUMappingResolver interface {
	ResolveMappedProductID(providerCode string, sellerSKU string, marketplace string) (uint64, error)
}

type OrderSyncService struct {
	repo             integrationDomain.OrderSyncRepository
	importer         OrderImportGateway
	provider         integrationDomain.OrdersProvider
	cfg              OrderSyncConfig
	currencyProvider OrderSyncCurrencyProvider
	mappingResolver  SKUMappingResolver
	nowFn            func() time.Time
}

func NewOrderSyncService(
	repo integrationDomain.OrderSyncRepository,
	importer OrderImportGateway,
	provider integrationDomain.OrdersProvider,
	cfg OrderSyncConfig,
	nowFn func() time.Time,
	currencyProvider OrderSyncCurrencyProvider,
) *OrderSyncService {
	if cfg.ProviderCode == "" && provider != nil {
		cfg.ProviderCode = provider.Code()
	}
	if cfg.ProviderCode == "" {
		cfg.ProviderCode = "UNKNOWN"
	}
	if cfg.Channel == "" {
		cfg.Channel = "DEFAULT"
	}
	if cfg.SourceType == "" {
		cfg.SourceType = "THIRD_PARTY_API"
	}
	if cfg.SalesChannel == "" {
		cfg.SalesChannel = strings.ToUpper(cfg.ProviderCode)
	}
	if cfg.DefaultCurrency == "" {
		if currencyProvider != nil {
			cfg.DefaultCurrency = strings.ToUpper(strings.TrimSpace(currencyProvider.GetDefaultBaseCurrency()))
		}
		if cfg.DefaultCurrency == "" {
			cfg.DefaultCurrency = "USD"
		}
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
	return &OrderSyncService{
		repo:             repo,
		importer:         importer,
		provider:         provider,
		cfg:              cfg,
		currencyProvider: currencyProvider,
		nowFn:            nowFn,
	}
}

func (u *OrderSyncService) BindSKUMappingResolver(resolver SKUMappingResolver) {
	u.mappingResolver = resolver
}

func (u *OrderSyncService) SyncOrders(ctx context.Context, trigger integrationDomain.SyncTrigger, operatorID *uint64) (*OrderSyncResult, error) {
	if u.repo == nil || u.importer == nil || u.provider == nil {
		return nil, fmt.Errorf("order sync dependencies not configured")
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
	if state != nil && state.LastUpdatedAfter != nil && !state.LastUpdatedAfter.IsZero() {
		cursor = state.LastUpdatedAfter.Add(-time.Duration(u.cfg.LookbackMinutes) * time.Minute)
	}

	run := &integrationDomain.OrderSyncRun{
		Provider:                u.cfg.ProviderCode,
		Channel:                 u.cfg.Channel,
		TriggerType:             trigger,
		Status:                  integrationDomain.OrderSyncRunStatusRunning,
		RequestLastUpdatedAfter: timePtr(cursor),
		StartedAt:               now,
	}
	if err := u.repo.CreateRun(run); err != nil {
		return nil, err
	}

	batch := &salesDomain.ReportImport{
		BatchNo:    buildBatchNo(now, u.cfg.ProviderCode),
		ReportType: fmt.Sprintf("%s_ORDERS", strings.ToUpper(u.cfg.ProviderCode)),
		FileName:   "THIRD-PARTY-API",
		FileHash:   buildBatchHash(u.cfg.ProviderCode, u.cfg.Channel, cursor, now),
		Status:     salesDomain.ReportImportStatusProcessing,
		OperatorID: operatorID,
		StartedAt:  timePtr(now),
	}
	if err := u.importer.CreateImport(batch); err != nil {
		u.finishFailedRun(run, err)
		return nil, err
	}
	run.BatchNo = &batch.BatchNo

	orders, err := u.provider.ListOrders(ctx, integrationDomain.ExternalListOrdersRequest{
		MarketplaceIDs:   u.cfg.MarketplaceIDs,
		LastUpdatedAfter: cursor,
	})
	if err != nil {
		u.finishImportFailed(batch, err)
		u.finishFailedRun(run, err)
		return nil, err
	}

	fetchedOrders := uint32(len(orders))
	fetchedLines := uint32(0)
	importedLines := uint32(0)
	rowErrors := make([]salesDomain.ReportImportRowError, 0)
	maxUpdatedAt := cursor

	for _, order := range orders {
		if order.LastUpdatedAt.After(maxUpdatedAt) {
			maxUpdatedAt = order.LastUpdatedAt
		}

		items := append([]integrationDomain.ExternalOrderItem(nil), order.Items...)
		sort.Slice(items, func(i, j int) bool {
			return items[i].OrderItemID < items[j].OrderItemID
		})

		for idx, item := range items {
			fetchedLines++
			line, lineErr := u.buildImportLine(order, item, uint32(idx+1), fetchedLines)
			if lineErr != nil {
				rowErrors = append(rowErrors, salesDomain.ReportImportRowError{
					RowNo:        fetchedLines,
					ErrorCode:    strPtr("LINE_BUILD_ERROR"),
					ErrorMessage: lineErr.Error(),
					RawRow:       strPtr(rawItem(order.OrderID, item.OrderItemID, item.SellerSKU)),
				})
				continue
			}

			productID, resolveErr := u.resolveProductID(line.SellerSKU, line.Marketplace)
			if resolveErr != nil {
				rowErrors = append(rowErrors, salesDomain.ReportImportRowError{
					RowNo:        fetchedLines,
					ErrorCode:    strPtr("SKU_LOOKUP_ERROR"),
					ErrorMessage: resolveErr.Error(),
					RawRow:       strPtr(line.RawRow),
				})
				continue
			}
			if productID == 0 {
				rowErrors = append(rowErrors, salesDomain.ReportImportRowError{
					RowNo:        fetchedLines,
					ErrorCode:    strPtr("SKU_NOT_FOUND"),
					ErrorMessage: fmt.Sprintf("seller_sku not found: %s", line.SellerSKU),
					RawRow:       strPtr(line.RawRow),
				})
				continue
			}

			if err := u.importer.UpsertImportedOrderLine(line, productID, batch.BatchNo, operatorID); err != nil {
				rowErrors = append(rowErrors, salesDomain.ReportImportRowError{
					RowNo:        fetchedLines,
					ErrorCode:    strPtr("UPSERT_FAILED"),
					ErrorMessage: err.Error(),
					RawRow:       strPtr(line.RawRow),
				})
				continue
			}
			importedLines++
		}
	}

	for i := range rowErrors {
		rowErrors[i].ReportImportID = batch.ID
	}
	if err := u.importer.InsertImportRowErrors(rowErrors); err != nil {
		u.finishImportFailed(batch, err)
		u.finishFailedRun(run, err)
		return nil, err
	}

	finishedAt := u.nowFn().UTC()
	batch.TotalRows = fetchedLines
	batch.SuccessRows = importedLines
	batch.ErrorRows = uint32(len(rowErrors))
	batch.Status = calcImportStatus(batch.SuccessRows, batch.ErrorRows)
	if batch.ErrorRows > 0 {
		msg := fmt.Sprintf("%s sync completed with %d error lines", strings.ToUpper(u.cfg.ProviderCode), batch.ErrorRows)
		batch.Message = &msg
	}
	batch.FinishedAt = timePtr(finishedAt)
	if err := u.importer.UpdateImport(batch); err != nil {
		u.finishFailedRun(run, err)
		return nil, err
	}

	run.FetchedOrders = fetchedOrders
	run.FetchedItems = fetchedLines
	run.ImportedItems = importedLines
	run.ErrorItems = uint32(len(rowErrors))
	run.Status = mapRunStatus(batch.Status)
	run.FinishedAt = timePtr(finishedAt)
	run.Message = batch.Message
	if err := u.repo.UpdateRun(run); err != nil {
		return nil, err
	}

	if fetchedOrders == 0 && maxUpdatedAt.Before(now) {
		maxUpdatedAt = now
	}
	if state == nil {
		state = &integrationDomain.OrderSyncState{
			Provider: u.cfg.ProviderCode,
			Channel:  u.cfg.Channel,
		}
	}
	state.LastUpdatedAfter = timePtr(maxUpdatedAt)
	state.LastSyncStarted = timePtr(now)
	state.LastSyncFinished = timePtr(finishedAt)
	if err := u.repo.SaveState(state); err != nil {
		return nil, err
	}

	return &OrderSyncResult{
		TriggerType:     trigger,
		RequestCursorAt: cursor,
		FetchedOrders:   fetchedOrders,
		FetchedLines:    fetchedLines,
		ImportedLines:   importedLines,
		ErrorLines:      uint32(len(rowErrors)),
		BatchNo:         batch.BatchNo,
		FinishedAt:      finishedAt,
	}, nil
}

func (u *OrderSyncService) GetState() (*integrationDomain.OrderSyncState, error) {
	return u.repo.GetState(u.cfg.ProviderCode, u.cfg.Channel)
}

func (u *OrderSyncService) ListRuns(page int, pageSize int) ([]integrationDomain.OrderSyncRun, int64, error) {
	return u.repo.ListRuns(u.cfg.ProviderCode, u.cfg.Channel, &integrationDomain.ListRunsParams{
		Page:     page,
		PageSize: pageSize,
	})
}

func (u *OrderSyncService) buildImportLine(
	order integrationDomain.ExternalOrder,
	item integrationDomain.ExternalOrderItem,
	lineNo uint32,
	rowNo uint32,
) (*salesDomain.ImportOrderLine, error) {
	sellerSKU := strings.TrimSpace(item.SellerSKU)
	if sellerSKU == "" {
		return nil, fmt.Errorf("seller_sku is required")
	}
	if item.Quantity == 0 {
		return nil, fmt.Errorf("quantity must be greater than 0")
	}
	if strings.TrimSpace(order.OrderID) == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	orderDate := order.PurchaseAt
	if orderDate.IsZero() {
		orderDate = order.LastUpdatedAt
	}
	if orderDate.IsZero() {
		orderDate = u.nowFn().UTC()
	}

	currency := strings.TrimSpace(order.Currency)
	if currency == "" {
		currency = u.cfg.DefaultCurrency
	}

	salesChannel := u.cfg.SalesChannel
	marketplace := normalizeMarketplace(order.MarketplaceID)
	unitPrice := 0.0
	if item.Amount > 0 && item.Quantity > 0 {
		unitPrice = item.Amount / float64(item.Quantity)
	}

	return &salesDomain.ImportOrderLine{
		RowNo:           rowNo,
		OrderNo:         order.OrderID,
		SourceType:      u.cfg.SourceType,
		ExternalOrderNo: order.OrderID,
		LineNo:          lineNo,
		SellerSKU:       sellerSKU,
		Qty:             item.Quantity,
		Marketplace:     marketplace,
		OrderDate:       orderDate,
		SalesChannel:    &salesChannel,
		Currency:        currency,
		UnitPrice:       unitPrice,
		RawRow:          rawItem(order.OrderID, item.OrderItemID, sellerSKU),
	}, nil
}

func (u *OrderSyncService) finishImportFailed(batch *salesDomain.ReportImport, err error) {
	if batch == nil {
		return
	}
	finishedAt := u.nowFn().UTC()
	msg := err.Error()
	batch.Status = salesDomain.ReportImportStatusFailed
	batch.Message = &msg
	batch.FinishedAt = timePtr(finishedAt)
	_ = u.importer.UpdateImport(batch)
}

func (u *OrderSyncService) finishFailedRun(run *integrationDomain.OrderSyncRun, err error) {
	if run == nil {
		return
	}
	finishedAt := u.nowFn().UTC()
	msg := err.Error()
	run.Status = integrationDomain.OrderSyncRunStatusFailed
	run.Message = &msg
	run.FinishedAt = timePtr(finishedAt)
	_ = u.repo.UpdateRun(run)
}

func (u *OrderSyncService) resolveProductID(sellerSKU string, marketplace string) (uint64, error) {
	if u.mappingResolver != nil {
		productID, err := u.mappingResolver.ResolveMappedProductID(u.cfg.ProviderCode, sellerSKU, marketplace)
		if err != nil {
			return 0, err
		}
		if productID > 0 {
			return productID, nil
		}
	}
	return u.importer.ResolveProductIDBySellerSKU(sellerSKU, marketplace)
}

func buildBatchNo(now time.Time, providerCode string) string {
	code := strings.ToUpper(strings.TrimSpace(providerCode))
	if code == "" {
		code = "TP"
	}
	return numbering.Generate(code, now.UTC())
}

func buildBatchHash(providerCode string, channel string, cursor time.Time, now time.Time) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%d", providerCode, channel, cursor.UTC().Format(time.RFC3339Nano), now.UnixNano())))
	return hex.EncodeToString(sum[:])
}

func calcImportStatus(success uint32, failed uint32) salesDomain.ReportImportStatus {
	switch {
	case success > 0 && failed == 0:
		return salesDomain.ReportImportStatusSuccess
	case success == 0 && failed > 0:
		return salesDomain.ReportImportStatusFailed
	case success > 0 && failed > 0:
		return salesDomain.ReportImportStatusPartialSuccess
	default:
		return salesDomain.ReportImportStatusSuccess
	}
}

func mapRunStatus(status salesDomain.ReportImportStatus) integrationDomain.OrderSyncRunStatus {
	switch status {
	case salesDomain.ReportImportStatusSuccess:
		return integrationDomain.OrderSyncRunStatusSuccess
	case salesDomain.ReportImportStatusPartialSuccess:
		return integrationDomain.OrderSyncRunStatusPartial
	default:
		return integrationDomain.OrderSyncRunStatusFailed
	}
}

func normalizeMarketplace(marketplaceID string) string {
	raw := strings.TrimSpace(strings.ToUpper(marketplaceID))
	switch raw {
	case "ATVPDKIKX0DER":
		return "US"
	case "A2EUQ1WTGCTBG2":
		return "CA"
	case "A39IBJ37TRP1C6":
		return "AU"
	case "A1F83G8C2ARO7P":
		return "UK"
	case "A1PA6795UKMFR9":
		return "DE"
	case "A1VC38T7YXB528":
		return "JP"
	default:
		return raw
	}
}

func rawItem(orderID string, orderItemID string, sellerSKU string) string {
	return fmt.Sprintf("order_id=%s,order_item_id=%s,seller_sku=%s", orderID, orderItemID, sellerSKU)
}

func strPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

func timePtr(v time.Time) *time.Time {
	if v.IsZero() {
		return nil
	}
	t := v.UTC()
	return &t
}
