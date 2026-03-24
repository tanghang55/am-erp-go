package usecase

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/numbering"
	inventoryDomain "am-erp-go/internal/module/inventory/domain"
	"am-erp-go/internal/module/sales/domain"
	systemdomain "am-erp-go/internal/module/system/domain"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type SalesOrderUsecase struct {
	repo                domain.SalesOrderRepository
	importRepo          domain.SalesImportRepository
	inventoryService    InventoryService
	shipCostDetailWrite SalesShipCostRecorder
	shipProfitWriter    SalesShipProfitRecorder
	auditLogger         AuditLogger
	configProvider      SalesConfigProvider
	allocateTxManager   SalesAllocateTransactionManager
	returnTxManager     SalesReturnTransactionManager
	shipTxManager       SalesShipTransactionManager
}

func NewSalesOrderUsecase(repo domain.SalesOrderRepository, importRepo ...domain.SalesImportRepository) *SalesOrderUsecase {
	var importRepository domain.SalesImportRepository
	if len(importRepo) > 0 {
		importRepository = importRepo[0]
	}
	if importRepository == nil {
		if casted, ok := repo.(domain.SalesImportRepository); ok {
			importRepository = casted
		}
	}

	return &SalesOrderUsecase{
		repo:       repo,
		importRepo: importRepository,
	}
}

type InventoryService interface {
	CreateMovement(ctx context.Context, params *inventoryDomain.CreateMovementParams) (*inventoryDomain.InventoryMovement, error)
}

type AuditLogger interface {
	RecordFromContext(c *gin.Context, payload systemUsecase.AuditLogPayload) error
}

type SalesConfigProvider interface {
	GetDefaultBaseCurrency() string
	GetSalesImportDefaults() systemdomain.ConfigCenterSalesImport
}

type SalesShipCostAllocation struct {
	InventoryLotID uint64
	Qty            uint64
	UnitCost       float64
}

type SalesShipCostRecordParams struct {
	SalesOrderID     uint64
	SalesOrderItemID uint64
	ProductID        uint64
	WarehouseID      uint64
	Marketplace      string
	Currency         string
	OccurredAt       time.Time
	OperatorID       *uint64
	Allocations      []SalesShipCostAllocation
}

type SalesReturnCostRecordParams struct {
	SalesOrderID     uint64
	SalesOrderItemID uint64
	ProductID        uint64
	WarehouseID      uint64
	Marketplace      string
	Currency         string
	QtyReturned      uint64
	OccurredAt       time.Time
	OperatorID       *uint64
}

type SalesShipCostRecorder interface {
	RecordSalesShipCost(params *SalesShipCostRecordParams) error
	ResolveSalesReturnUnitCost(params *SalesReturnCostRecordParams) (*float64, error)
	RecordSalesReturnCost(params *SalesReturnCostRecordParams) (float64, error)
}

type SalesShipProfitRecordParams struct {
	SalesOrderID     uint64
	SalesOrderItemID uint64
	OrderNo          string
	Marketplace      string
	IncomeCurrency   string
	IncomeAmount     float64
	COGSAmount       float64
	OccurredAt       time.Time
	OperatorID       *uint64
}

type SalesReturnProfitRecordParams struct {
	SalesOrderID     uint64
	SalesOrderItemID uint64
	OrderNo          string
	Marketplace      string
	IncomeCurrency   string
	IncomeAmount     float64
	COGSAmount       float64
	OccurredAt       time.Time
	OperatorID       *uint64
}

type SalesShipProfitRecorder interface {
	RecordSalesShipProfit(params *SalesShipProfitRecordParams) error
	RecordSalesReturnProfit(params *SalesReturnProfitRecordParams) error
}

type SalesShipTransactionalDeps struct {
	Repo             domain.SalesOrderRepository
	InventoryService InventoryService
	CostWriter       SalesShipCostRecorder
	ProfitWriter     SalesShipProfitRecorder
}

type SalesAllocateTransactionalDeps struct {
	Repo             domain.SalesOrderRepository
	InventoryService InventoryService
}

type SalesAllocateTransactionManager interface {
	Run(ctx context.Context, fn func(SalesAllocateTransactionalDeps) error) error
}

type SalesReturnTransactionalDeps struct {
	Repo             domain.SalesOrderRepository
	InventoryService InventoryService
	CostWriter       SalesShipCostRecorder
	ProfitWriter     SalesShipProfitRecorder
}

type SalesReturnTransactionManager interface {
	Run(ctx context.Context, fn func(SalesReturnTransactionalDeps) error) error
}

type SalesShipTransactionManager interface {
	Run(ctx context.Context, fn func(SalesShipTransactionalDeps) error) error
}

func (u *SalesOrderUsecase) BindInventoryService(svc InventoryService) {
	u.inventoryService = svc
}

func (u *SalesOrderUsecase) BindShipCostDetailWriter(writer SalesShipCostRecorder) {
	u.shipCostDetailWrite = writer
}

func (u *SalesOrderUsecase) BindShipProfitWriter(writer SalesShipProfitRecorder) {
	u.shipProfitWriter = writer
}

func (u *SalesOrderUsecase) BindAuditLogger(logger AuditLogger) {
	u.auditLogger = logger
}

func (u *SalesOrderUsecase) BindConfigProvider(provider SalesConfigProvider) {
	u.configProvider = provider
}

func (u *SalesOrderUsecase) BindAllocateTransactionManager(manager SalesAllocateTransactionManager) {
	u.allocateTxManager = manager
}

func (u *SalesOrderUsecase) BindReturnTransactionManager(manager SalesReturnTransactionManager) {
	u.returnTxManager = manager
}

func (u *SalesOrderUsecase) BindShipTransactionManager(manager SalesShipTransactionManager) {
	u.shipTxManager = manager
}

func (u *SalesOrderUsecase) List(params *domain.SalesOrderListParams) ([]domain.SalesOrder, int64, error) {
	if params == nil {
		params = &domain.SalesOrderListParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return u.repo.List(params)
}

func (u *SalesOrderUsecase) Get(orderID uint64) (*domain.SalesOrder, error) {
	order, err := u.repo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}
	return order, nil
}

func (u *SalesOrderUsecase) Create(c *gin.Context, order *domain.SalesOrder) (*domain.SalesOrder, error) {
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}
	if len(order.Items) == 0 {
		return nil, domain.ErrInvalidQuantity
	}
	if order.OrderStatus == "" {
		order.OrderStatus = domain.SalesOrderStatusDraft
	}
	if err := u.repo.Create(order); err != nil {
		return nil, err
	}
	u.recordAudit(c, "CREATE", order.ID, nil, order)
	return order, nil
}

func (u *SalesOrderUsecase) Update(c *gin.Context, orderID uint64, updates *domain.SalesOrder) (*domain.SalesOrder, error) {
	if updates == nil {
		return nil, domain.ErrOrderNotFound
	}
	existing, err := u.repo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, domain.ErrOrderNotFound
	}
	if existing.OrderStatus != domain.SalesOrderStatusDraft && existing.OrderStatus != domain.SalesOrderStatusConfirmed {
		return nil, domain.ErrInvalidTransition
	}

	existing.ExternalOrderNo = updates.ExternalOrderNo
	existing.SalesChannel = updates.SalesChannel
	existing.Marketplace = updates.Marketplace
	existing.OrderDate = updates.OrderDate
	existing.Currency = updates.Currency
	existing.OrderAmount = updates.OrderAmount
	existing.Remark = updates.Remark
	existing.UpdatedBy = updates.UpdatedBy
	if updates.Items != nil {
		existing.Items = updates.Items
	}

	if err := u.repo.Update(existing); err != nil {
		return nil, err
	}
	u.recordAudit(c, "UPDATE", orderID, nil, existing)
	return existing, nil
}

func (u *SalesOrderUsecase) ImportCSV(ctx context.Context, fileName string, content []byte, operatorID *uint64) (*domain.ReportImport, error) {
	_ = ctx
	if u.importRepo == nil {
		return nil, domain.ErrImportInvalidFile
	}
	if strings.TrimSpace(fileName) == "" || len(content) == 0 {
		return nil, domain.ErrImportInvalidFile
	}

	fileHash := hashBytes(content)
	existing, err := u.importRepo.GetImportByFileHash(fileHash)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrImportDuplicateFile
	}

	now := time.Now()
	batch := &domain.ReportImport{
		BatchNo:     buildImportBatchNo(),
		ReportType:  "ORDERS",
		FileName:    fileName,
		FileHash:    fileHash,
		Status:      domain.ReportImportStatusProcessing,
		OperatorID:  operatorID,
		StartedAt:   &now,
		TotalRows:   0,
		SuccessRows: 0,
		ErrorRows:   0,
	}
	if err := u.importRepo.CreateImport(batch); err != nil {
		return nil, err
	}

	lines, rowErrors, totalRows, parseErr := parseImportRows(content, u.getImportDefaults())
	batch.TotalRows = totalRows
	if parseErr != nil {
		msg := parseErr.Error()
		finishedAt := time.Now()
		batch.Status = domain.ReportImportStatusFailed
		batch.Message = &msg
		batch.ErrorRows = totalRows
		batch.FinishedAt = &finishedAt
		_ = u.importRepo.UpdateImport(batch)
		return nil, parseErr
	}

	successRows := uint32(0)
	for _, line := range lines {
		productID, resolveErr := u.importRepo.ResolveProductIDBySellerSKU(line.SellerSKU, line.Marketplace)
		if resolveErr != nil {
			rowErrors = append(rowErrors, domain.ReportImportRowError{
				ReportImportID: batch.ID,
				RowNo:          line.RowNo,
				ErrorCode:      strPtr("SKU_LOOKUP_ERROR"),
				ErrorMessage:   resolveErr.Error(),
				RawRow:         strPtr(line.RawRow),
			})
			continue
		}
		if productID == 0 {
			rowErrors = append(rowErrors, domain.ReportImportRowError{
				ReportImportID: batch.ID,
				RowNo:          line.RowNo,
				ErrorCode:      strPtr("SKU_NOT_FOUND"),
				ErrorMessage:   fmt.Sprintf("seller_sku not found: %s", line.SellerSKU),
				RawRow:         strPtr(line.RawRow),
			})
			continue
		}

		if err := u.importRepo.UpsertImportedOrderLine(&line, productID, batch.BatchNo, operatorID); err != nil {
			rowErrors = append(rowErrors, domain.ReportImportRowError{
				ReportImportID: batch.ID,
				RowNo:          line.RowNo,
				ErrorCode:      strPtr("UPSERT_FAILED"),
				ErrorMessage:   err.Error(),
				RawRow:         strPtr(line.RawRow),
			})
			continue
		}

		successRows++
	}

	for i := range rowErrors {
		rowErrors[i].ReportImportID = batch.ID
	}
	if err := u.importRepo.InsertImportRowErrors(rowErrors); err != nil {
		return nil, err
	}

	batch.SuccessRows = successRows
	batch.ErrorRows = uint32(len(rowErrors))
	batch.Status = calcImportStatus(batch.SuccessRows, batch.ErrorRows)
	if batch.ErrorRows > 0 {
		message := fmt.Sprintf("import completed with %d errors", batch.ErrorRows)
		batch.Message = &message
	}
	finishedAt := time.Now()
	batch.FinishedAt = &finishedAt
	if err := u.importRepo.UpdateImport(batch); err != nil {
		return nil, err
	}

	return batch, nil
}

func (u *SalesOrderUsecase) ListImports(page int, pageSize int) ([]domain.ReportImport, int64, error) {
	if u.importRepo == nil {
		return []domain.ReportImport{}, 0, nil
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return u.importRepo.ListImports(page, pageSize)
}

func (u *SalesOrderUsecase) GetImport(importID uint64) (*domain.ReportImport, error) {
	if u.importRepo == nil {
		return nil, domain.ErrImportNotFound
	}
	item, err := u.importRepo.GetImportByID(importID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, domain.ErrImportNotFound
	}
	return item, nil
}

func (u *SalesOrderUsecase) ListImportErrors(importID uint64) ([]domain.ReportImportRowError, error) {
	if u.importRepo == nil {
		return nil, domain.ErrImportNotFound
	}
	return u.importRepo.ListImportErrors(importID)
}

func (u *SalesOrderUsecase) Confirm(c *gin.Context, orderID uint64, operatorID *uint64) error {
	_ = operatorID

	order, err := u.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}
	if order.OrderStatus != domain.SalesOrderStatusDraft {
		return domain.ErrInvalidTransition
	}

	now := time.Now()
	order.OrderStatus = domain.SalesOrderStatusConfirmed
	order.ConfirmAt = &now
	if err := u.repo.Update(order); err != nil {
		return err
	}
	u.recordAudit(c, "CONFIRM", orderID, nil, order)
	return nil
}

func (u *SalesOrderUsecase) Cancel(c *gin.Context, orderID uint64, operatorID *uint64) error {
	_ = operatorID

	order, err := u.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}
	if order.OrderStatus == domain.SalesOrderStatusShipped || order.OrderStatus == domain.SalesOrderStatusDelivered {
		return domain.ErrInvalidTransition
	}

	now := time.Now()
	order.OrderStatus = domain.SalesOrderStatusCancelled
	order.CancelledAt = &now
	if err := u.repo.Update(order); err != nil {
		return err
	}
	u.recordAudit(c, "CANCEL", orderID, nil, order)
	return nil
}

func (u *SalesOrderUsecase) Allocate(c *gin.Context, orderID uint64, params domain.AllocateParams, operatorID *uint64) error {
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if u.allocateTxManager != nil {
		if err := u.allocateTxManager.Run(ctx, func(deps SalesAllocateTransactionalDeps) error {
			return u.allocateWithDeps(c, orderID, params, operatorID, deps)
		}); err != nil {
			return err
		}
		if after, err := u.repo.GetByID(orderID); err == nil {
			u.recordAudit(c, "ALLOCATE", orderID, nil, after)
		}
		return nil
	}
	if err := u.allocateWithDeps(c, orderID, params, operatorID, SalesAllocateTransactionalDeps{
		Repo:             u.repo,
		InventoryService: u.inventoryService,
	}); err != nil {
		return err
	}
	if after, err := u.repo.GetByID(orderID); err == nil {
		u.recordAudit(c, "ALLOCATE", orderID, nil, after)
	}
	return nil
}

func (u *SalesOrderUsecase) allocateWithDeps(
	c *gin.Context,
	orderID uint64,
	params domain.AllocateParams,
	operatorID *uint64,
	deps SalesAllocateTransactionalDeps,
) error {
	execUC := *u
	execUC.repo = deps.Repo
	execUC.inventoryService = deps.InventoryService

	order, err := execUC.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}
	if order.OrderStatus != domain.SalesOrderStatusConfirmed && order.OrderStatus != domain.SalesOrderStatusAllocated {
		return domain.ErrInvalidTransition
	}

	for _, line := range params.Lines {
		item := findItem(order.Items, line.ItemID)
		if item == nil {
			return domain.ErrItemNotFound
		}
		if line.QtyAllocated == 0 || line.QtyAllocated > item.QtyOrdered {
			return domain.ErrInvalidQuantity
		}

		delta := int64(line.QtyAllocated) - int64(item.QtyAllocated)
		if delta > 0 {
			if _, err := execUC.createInventoryMovement(
				c,
				domain.SalesOrderStatusAllocated,
				inventoryDomain.MovementTypeSalesAllocate,
				order.StockPool,
				params.WarehouseID,
				item.ProductID,
				uint64(delta),
				nil,
				orderID,
				order.OrderNo,
				operatorID,
				"sales allocate",
			); err != nil {
				return err
			}
		} else if delta < 0 {
			if _, err := execUC.createInventoryMovement(
				c,
				domain.SalesOrderStatusAllocated,
				inventoryDomain.MovementTypeSalesRelease,
				order.StockPool,
				params.WarehouseID,
				item.ProductID,
				uint64(-delta),
				nil,
				orderID,
				order.OrderNo,
				operatorID,
				"sales release",
			); err != nil {
				return err
			}
		}

		item.QtyAllocated = line.QtyAllocated
	}

	now := time.Now()
	order.OrderStatus = domain.SalesOrderStatusAllocated
	order.AllocatedAt = &now
	if err := execUC.repo.Update(order); err != nil {
		return err
	}
	return nil
}

func (u *SalesOrderUsecase) Ship(c *gin.Context, orderID uint64, params domain.ShipParams, operatorID *uint64) error {
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if u.shipTxManager != nil {
		if err := u.shipTxManager.Run(ctx, func(deps SalesShipTransactionalDeps) error {
			return u.shipWithDeps(c, orderID, params, operatorID, deps)
		}); err != nil {
			return err
		}
		if after, err := u.repo.GetByID(orderID); err == nil {
			u.recordAudit(c, "SHIP", orderID, nil, after)
		}
		return nil
	}
	if err := u.shipWithDeps(c, orderID, params, operatorID, SalesShipTransactionalDeps{
		Repo:             u.repo,
		InventoryService: u.inventoryService,
		CostWriter:       u.shipCostDetailWrite,
		ProfitWriter:     u.shipProfitWriter,
	}); err != nil {
		return err
	}
	if after, err := u.repo.GetByID(orderID); err == nil {
		u.recordAudit(c, "SHIP", orderID, nil, after)
	}
	return nil
}

func (u *SalesOrderUsecase) shipWithDeps(
	c *gin.Context,
	orderID uint64,
	params domain.ShipParams,
	operatorID *uint64,
	deps SalesShipTransactionalDeps,
) error {
	execUC := *u
	execUC.repo = deps.Repo
	execUC.inventoryService = deps.InventoryService
	execUC.shipCostDetailWrite = deps.CostWriter
	execUC.shipProfitWriter = deps.ProfitWriter

	order, err := execUC.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}
	if order.OrderStatus != domain.SalesOrderStatusAllocated && order.OrderStatus != domain.SalesOrderStatusShipped {
		return domain.ErrInvalidTransition
	}

	for _, line := range params.Lines {
		item := findItem(order.Items, line.ItemID)
		if item == nil {
			return domain.ErrItemNotFound
		}
		if line.QtyShipped == 0 || item.QtyShipped+line.QtyShipped > item.QtyAllocated {
			return domain.ErrInvalidQuantity
		}

		movement, err := execUC.createInventoryMovement(
			c,
			domain.SalesOrderStatusShipped,
			inventoryDomain.MovementTypeSalesShip,
			order.StockPool,
			params.WarehouseID,
			item.ProductID,
			line.QtyShipped,
			nil,
			orderID,
			order.OrderNo,
			operatorID,
			"sales ship",
		)
		if err != nil {
			return err
		}
		if err := execUC.recordShipCostDetail(order, item, movement, params.WarehouseID, operatorID); err != nil {
			return err
		}
		if err := execUC.recordShipProfit(order, item, line.QtyShipped, movement, operatorID); err != nil {
			return err
		}

		item.QtyShipped += line.QtyShipped
	}

	now := time.Now()
	order.OrderStatus = domain.SalesOrderStatusShipped
	order.ShippedAt = &now
	if err := execUC.repo.Update(order); err != nil {
		return err
	}
	return nil
}

func (u *SalesOrderUsecase) Deliver(c *gin.Context, orderID uint64, operatorID *uint64) error {
	_ = operatorID

	order, err := u.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}
	if order.OrderStatus != domain.SalesOrderStatusShipped {
		return domain.ErrInvalidTransition
	}

	now := time.Now()
	order.OrderStatus = domain.SalesOrderStatusDelivered
	order.DeliveredAt = &now
	if err := u.repo.Update(order); err != nil {
		return err
	}
	u.recordAudit(c, "DELIVER", orderID, nil, order)
	return nil
}

func (u *SalesOrderUsecase) Return(c *gin.Context, orderID uint64, params domain.ReturnParams, operatorID *uint64) error {
	ctx := context.Background()
	if c != nil && c.Request != nil {
		ctx = c.Request.Context()
	}
	if u.returnTxManager != nil {
		if err := u.returnTxManager.Run(ctx, func(deps SalesReturnTransactionalDeps) error {
			return u.returnWithDeps(c, orderID, params, operatorID, deps)
		}); err != nil {
			return err
		}
		if after, err := u.repo.GetByID(orderID); err == nil {
			u.recordAudit(c, "RETURN", orderID, nil, after)
		}
		return nil
	}
	if err := u.returnWithDeps(c, orderID, params, operatorID, SalesReturnTransactionalDeps{
		Repo:             u.repo,
		InventoryService: u.inventoryService,
		CostWriter:       u.shipCostDetailWrite,
		ProfitWriter:     u.shipProfitWriter,
	}); err != nil {
		return err
	}
	if after, err := u.repo.GetByID(orderID); err == nil {
		u.recordAudit(c, "RETURN", orderID, nil, after)
	}
	return nil
}

func (u *SalesOrderUsecase) returnWithDeps(
	c *gin.Context,
	orderID uint64,
	params domain.ReturnParams,
	operatorID *uint64,
	deps SalesReturnTransactionalDeps,
) error {
	execUC := *u
	execUC.repo = deps.Repo
	execUC.inventoryService = deps.InventoryService
	execUC.shipCostDetailWrite = deps.CostWriter
	execUC.shipProfitWriter = deps.ProfitWriter

	order, err := execUC.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return domain.ErrOrderNotFound
	}
	if order.OrderStatus != domain.SalesOrderStatusDelivered && order.OrderStatus != domain.SalesOrderStatusReturned {
		return domain.ErrInvalidTransition
	}

	for _, line := range params.Lines {
		item := findItem(order.Items, line.ItemID)
		if item == nil {
			return domain.ErrItemNotFound
		}
		if line.QtyReturned == 0 {
			return domain.ErrInvalidQuantity
		}
		if item.QtyReturned+line.QtyReturned > item.QtyShipped {
			return domain.ErrInvalidQuantity
		}

		returnUnitCost, err := execUC.resolveReturnUnitCost(order, item, line.QtyReturned, params.WarehouseID, operatorID)
		if err != nil {
			return err
		}

		movement, err := execUC.createInventoryMovement(
			c,
			domain.SalesOrderStatusReturned,
			inventoryDomain.MovementTypeReturnReceipt,
			order.StockPool,
			params.WarehouseID,
			item.ProductID,
			line.QtyReturned,
			returnUnitCost,
			orderID,
			order.OrderNo,
			operatorID,
			"sales return",
		)
		if err != nil {
			return err
		}
		if err := execUC.recordReturnFinance(order, item, line.QtyReturned, params.WarehouseID, movement, operatorID); err != nil {
			return err
		}

		item.QtyReturned += line.QtyReturned
	}

	order.OrderStatus = domain.SalesOrderStatusReturned
	if err := execUC.repo.Update(order); err != nil {
		return err
	}
	return nil
}

func (u *SalesOrderUsecase) recordReturnFinance(
	order *domain.SalesOrder,
	item *domain.SalesOrderItem,
	qtyReturned uint64,
	warehouseID uint64,
	movement *inventoryDomain.InventoryMovement,
	operatorID *uint64,
) error {
	if order == nil || item == nil || qtyReturned == 0 {
		return nil
	}

	marketplace := ""
	if order.Marketplace != nil {
		marketplace = *order.Marketplace
	}
	currency := strings.ToUpper(strings.TrimSpace(order.Currency))
	if currency == "" {
		currency = u.getFallbackCurrency()
	}
	occurredAt := time.Now()
	if movement != nil && !movement.OperatedAt.IsZero() {
		occurredAt = movement.OperatedAt
	}

	cogsAmount := 0.0
	if u.shipCostDetailWrite != nil {
		var err error
		cogsAmount, err = u.shipCostDetailWrite.RecordSalesReturnCost(&SalesReturnCostRecordParams{
			SalesOrderID:     order.ID,
			SalesOrderItemID: item.ID,
			ProductID:        item.ProductID,
			WarehouseID:      warehouseID,
			Marketplace:      marketplace,
			Currency:         currency,
			QtyReturned:      qtyReturned,
			OccurredAt:       occurredAt,
			OperatorID:       operatorID,
		})
		if err != nil {
			return err
		}
	}

	if u.shipProfitWriter == nil {
		return nil
	}
	return u.shipProfitWriter.RecordSalesReturnProfit(&SalesReturnProfitRecordParams{
		SalesOrderID:     order.ID,
		SalesOrderItemID: item.ID,
		OrderNo:          order.OrderNo,
		Marketplace:      marketplace,
		IncomeCurrency:   currency,
		IncomeAmount:     round6(float64(qtyReturned) * item.UnitPrice),
		COGSAmount:       round6(cogsAmount),
		OccurredAt:       occurredAt,
		OperatorID:       operatorID,
	})
}

func (u *SalesOrderUsecase) resolveReturnUnitCost(
	order *domain.SalesOrder,
	item *domain.SalesOrderItem,
	qtyReturned uint64,
	warehouseID uint64,
	operatorID *uint64,
) (*float64, error) {
	if u.shipCostDetailWrite == nil || order == nil || item == nil || qtyReturned == 0 {
		return nil, nil
	}

	marketplace := ""
	if order.Marketplace != nil {
		marketplace = *order.Marketplace
	}
	currency := strings.ToUpper(strings.TrimSpace(order.Currency))
	if currency == "" {
		currency = u.getFallbackCurrency()
	}

	return u.shipCostDetailWrite.ResolveSalesReturnUnitCost(&SalesReturnCostRecordParams{
		SalesOrderID:     order.ID,
		SalesOrderItemID: item.ID,
		ProductID:        item.ProductID,
		WarehouseID:      warehouseID,
		Marketplace:      marketplace,
		Currency:         currency,
		QtyReturned:      qtyReturned,
		OccurredAt:       time.Now(),
		OperatorID:       operatorID,
	})
}

func (u *SalesOrderUsecase) createInventoryMovement(
	ctx context.Context,
	_ domain.SalesOrderStatus,
	movementType inventoryDomain.MovementType,
	stockPool domain.StockPool,
	warehouseID uint64,
	productID uint64,
	qty uint64,
	unitCost *float64,
	orderID uint64,
	orderNo string,
	operatorID *uint64,
	action string,
) (*inventoryDomain.InventoryMovement, error) {
	if u.inventoryService == nil || qty == 0 || warehouseID == 0 {
		return nil, nil
	}

	referenceType := "SALES_ORDER"
	referenceNumber := orderNo
	referenceID := orderID
	remark := fmt.Sprintf("%s: %s", action, orderNo)

	movement, err := u.inventoryService.CreateMovement(ctx, &inventoryDomain.CreateMovementParams{
		ProductID:       productID,
		WarehouseID:     warehouseID,
		MovementType:    movementType,
		Quantity:        int(qty),
		ReferenceType:   &referenceType,
		ReferenceID:     &referenceID,
		ReferenceNumber: &referenceNumber,
		StockPool:       toInventoryStockPool(stockPool),
		UnitCost:        unitCost,
		Remark:          &remark,
		OperatorID:      operatorID,
	})
	return movement, err
}

func toInventoryStockPool(pool domain.StockPool) *inventoryDomain.StockPool {
	switch pool {
	case domain.StockPoolSellable:
		p := inventoryDomain.StockPoolSellable
		return &p
	case domain.StockPoolAvailable:
		fallthrough
	default:
		p := inventoryDomain.StockPoolAvailable
		return &p
	}
}

func (u *SalesOrderUsecase) recordShipCostDetail(
	order *domain.SalesOrder,
	item *domain.SalesOrderItem,
	movement *inventoryDomain.InventoryMovement,
	warehouseID uint64,
	operatorID *uint64,
) error {
	if u.shipCostDetailWrite == nil || order == nil || item == nil || movement == nil || len(movement.LotAllocations) == 0 {
		return nil
	}

	allocations := make([]SalesShipCostAllocation, 0, len(movement.LotAllocations))
	for _, lot := range movement.LotAllocations {
		if lot.InventoryLotID == 0 || lot.Qty == 0 {
			continue
		}
		allocations = append(allocations, SalesShipCostAllocation{
			InventoryLotID: lot.InventoryLotID,
			Qty:            lot.Qty,
			UnitCost:       lot.UnitCost,
		})
	}
	if len(allocations) == 0 {
		return nil
	}

	marketplace := ""
	if order.Marketplace != nil {
		marketplace = *order.Marketplace
	}
	currency := strings.ToUpper(strings.TrimSpace(order.Currency))
	if currency == "" {
		currency = u.getFallbackCurrency()
	}

	return u.shipCostDetailWrite.RecordSalesShipCost(&SalesShipCostRecordParams{
		SalesOrderID:     order.ID,
		SalesOrderItemID: item.ID,
		ProductID:        item.ProductID,
		WarehouseID:      warehouseID,
		Marketplace:      marketplace,
		Currency:         currency,
		OccurredAt:       movement.OperatedAt,
		OperatorID:       operatorID,
		Allocations:      allocations,
	})
}

func (u *SalesOrderUsecase) recordShipProfit(
	order *domain.SalesOrder,
	item *domain.SalesOrderItem,
	qtyShipped uint64,
	movement *inventoryDomain.InventoryMovement,
	operatorID *uint64,
) error {
	if u.shipProfitWriter == nil || order == nil || item == nil || qtyShipped == 0 || movement == nil {
		return nil
	}
	marketplace := ""
	if order.Marketplace != nil {
		marketplace = *order.Marketplace
	}
	currency := strings.ToUpper(strings.TrimSpace(order.Currency))
	if currency == "" {
		currency = u.getFallbackCurrency()
	}
	incomeAmount := round6(float64(qtyShipped) * item.UnitPrice)
	cogsAmount := 0.0
	for _, lot := range movement.LotAllocations {
		if lot.Qty == 0 {
			continue
		}
		cogsAmount += float64(lot.Qty) * lot.UnitCost
	}

	return u.shipProfitWriter.RecordSalesShipProfit(&SalesShipProfitRecordParams{
		SalesOrderID:     order.ID,
		SalesOrderItemID: item.ID,
		OrderNo:          order.OrderNo,
		Marketplace:      marketplace,
		IncomeCurrency:   currency,
		IncomeAmount:     incomeAmount,
		COGSAmount:       round6(cogsAmount),
		OccurredAt:       movement.OperatedAt,
		OperatorID:       operatorID,
	})
}

func round6(v float64) float64 {
	return float64(int64(v*1_000_000+0.5)) / 1_000_000
}

func findItem(items []domain.SalesOrderItem, itemID uint64) *domain.SalesOrderItem {
	for i := range items {
		if items[i].ID == itemID {
			return &items[i]
		}
	}
	return nil
}

func hashBytes(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func (u *SalesOrderUsecase) recordAudit(c *gin.Context, action string, orderID uint64, before, after any) {
	if u.auditLogger == nil || c == nil || orderID == 0 {
		return
	}
	_ = u.auditLogger.RecordFromContext(c, systemUsecase.AuditLogPayload{
		Module:     "Sales",
		Action:     action,
		EntityType: "SalesOrder",
		EntityID:   fmt.Sprintf("%d", orderID),
		Before:     before,
		After:      after,
	})
}

func buildImportBatchNo() string {
	return numbering.Generate("IMP", time.Now())
}

func calcImportStatus(successRows uint32, errorRows uint32) domain.ReportImportStatus {
	switch {
	case successRows > 0 && errorRows == 0:
		return domain.ReportImportStatusSuccess
	case successRows > 0 && errorRows > 0:
		return domain.ReportImportStatusPartialSuccess
	case successRows == 0 && errorRows > 0:
		return domain.ReportImportStatusFailed
	default:
		return domain.ReportImportStatusSuccess
	}
}

func parseImportRows(content []byte, defaults importDefaults) ([]domain.ImportOrderLine, []domain.ReportImportRowError, uint32, error) {
	reader := csv.NewReader(bytes.NewReader(content))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, 0, err
	}
	if len(records) <= 1 {
		return nil, nil, 0, domain.ErrImportInvalidFile
	}

	headerMap := map[string]int{}
	for i, header := range records[0] {
		headerMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	required := []string{"order_date", "order_no", "line_no", "seller_sku", "qty"}
	for _, field := range required {
		if _, ok := headerMap[field]; !ok {
			return nil, nil, 0, fmt.Errorf("%w: missing required column %s", domain.ErrImportInvalidFile, field)
		}
	}

	totalRows := uint32(0)
	lines := make([]domain.ImportOrderLine, 0, len(records)-1)
	rowErrors := make([]domain.ReportImportRowError, 0)

	for i := 1; i < len(records); i++ {
		record := records[i]
		totalRows++
		rowNo := uint32(i + 1)
		rawRow := strings.Join(record, ",")

		line, parseErr := parseImportLine(record, rowNo, rawRow, headerMap, defaults)
		if parseErr != nil {
			rowErrors = append(rowErrors, domain.ReportImportRowError{
				RowNo:        rowNo,
				ErrorCode:    strPtr("INVALID_ROW"),
				ErrorMessage: parseErr.Error(),
				RawRow:       strPtr(rawRow),
			})
			continue
		}
		lines = append(lines, *line)
	}

	return lines, rowErrors, totalRows, nil
}

func parseImportLine(record []string, rowNo uint32, rawRow string, headerMap map[string]int, defaults importDefaults) (*domain.ImportOrderLine, error) {
	get := func(key string) string {
		idx, ok := headerMap[key]
		if !ok || idx < 0 || idx >= len(record) {
			return ""
		}
		return strings.TrimSpace(record[idx])
	}

	orderDateRaw := get("order_date")
	orderDate, err := parseOrderDate(orderDateRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid order_date: %s", orderDateRaw)
	}

	orderNo := get("order_no")
	if orderNo == "" {
		return nil, fmt.Errorf("order_no is required")
	}

	lineNoRaw := get("line_no")
	lineNoVal, err := strconv.ParseUint(lineNoRaw, 10, 32)
	if err != nil || lineNoVal == 0 {
		return nil, fmt.Errorf("invalid line_no: %s", lineNoRaw)
	}

	sellerSKU := get("seller_sku")
	if sellerSKU == "" {
		return nil, fmt.Errorf("seller_sku is required")
	}

	qtyRaw := get("qty")
	qtyVal, err := strconv.ParseUint(qtyRaw, 10, 64)
	if err != nil || qtyVal == 0 {
		return nil, fmt.Errorf("invalid qty: %s", qtyRaw)
	}

	marketplace := get("marketplace")
	if marketplace == "" {
		marketplace = defaults.DefaultMarketplace
	}
	if marketplace == "" {
		return nil, fmt.Errorf("marketplace is required")
	}

	externalOrderNo := get("external_order_no")
	if externalOrderNo == "" {
		externalOrderNo = orderNo
	}

	currency := get("currency")
	if currency == "" {
		currency = defaults.DefaultCurrency
	}
	if currency == "" {
		currency = defaults.DefaultCurrency
	}

	unitPriceRaw := get("unit_price")
	unitPrice := 0.0
	if unitPriceRaw != "" {
		parsed, parseErr := strconv.ParseFloat(unitPriceRaw, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid unit_price: %s", unitPriceRaw)
		}
		unitPrice = parsed
	}

	var salesChannel *string
	salesChannelRaw := get("sales_channel")
	if salesChannelRaw == "" {
		salesChannelRaw = defaults.DefaultChannel
	}
	sourceType := resolveImportedSourceType(salesChannelRaw)
	if salesChannelRaw != "" {
		salesChannel = &salesChannelRaw
	}

	return &domain.ImportOrderLine{
		RowNo:           rowNo,
		OrderNo:         orderNo,
		SourceType:      sourceType,
		ExternalOrderNo: externalOrderNo,
		LineNo:          uint32(lineNoVal),
		SellerSKU:       sellerSKU,
		Qty:             qtyVal,
		Marketplace:     marketplace,
		OrderDate:       orderDate,
		SalesChannel:    salesChannel,
		Currency:        currency,
		UnitPrice:       unitPrice,
		RawRow:          rawRow,
	}, nil
}

func resolveImportedSourceType(salesChannel string) string {
	channel := strings.ToUpper(strings.TrimSpace(salesChannel))
	switch {
	case strings.Contains(channel, "AMAZON"):
		return "AMAZON_IMPORT"
	default:
		return "MANUAL_IMPORT"
	}
}

func parseOrderDate(raw string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"2006/01/02 15:04:05",
	}
	raw = strings.TrimSpace(raw)
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", raw)
}

func strPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

type importDefaults struct {
	DefaultCurrency    string
	DefaultChannel     string
	DefaultMarketplace string
}

func (u *SalesOrderUsecase) getImportDefaults() importDefaults {
	defaults := importDefaults{
		DefaultCurrency:    "USD",
		DefaultChannel:     "MANUAL",
		DefaultMarketplace: "US",
	}
	if u.configProvider == nil {
		return defaults
	}
	if currency := strings.ToUpper(strings.TrimSpace(u.configProvider.GetDefaultBaseCurrency())); currency != "" {
		defaults.DefaultCurrency = currency
	}
	importConfig := u.configProvider.GetSalesImportDefaults()
	if value := strings.TrimSpace(importConfig.DefaultChannel); value != "" {
		defaults.DefaultChannel = value
	}
	if value := strings.ToUpper(strings.TrimSpace(importConfig.DefaultMarketplace)); value != "" {
		defaults.DefaultMarketplace = value
	}
	return defaults
}

func (u *SalesOrderUsecase) getFallbackCurrency() string {
	if u.configProvider != nil {
		if currency := strings.ToUpper(strings.TrimSpace(u.configProvider.GetDefaultBaseCurrency())); currency != "" {
			return currency
		}
	}
	return "USD"
}
