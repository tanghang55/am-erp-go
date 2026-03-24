package bootstrap

import (
	"context"
	"log"
	"strings"
	"time"

	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/config"
	"am-erp-go/internal/infrastructure/db"
	"am-erp-go/internal/infrastructure/router"
	infraTx "am-erp-go/internal/infrastructure/transaction"
	"am-erp-go/internal/infrastructure/upload"
	financeHttp "am-erp-go/internal/module/finance/delivery/http"
	financeDomain "am-erp-go/internal/module/finance/domain"
	financeRepo "am-erp-go/internal/module/finance/repository"
	financeUsecase "am-erp-go/internal/module/finance/usecase"
	identityHttp "am-erp-go/internal/module/identity/delivery/http"
	identityRepo "am-erp-go/internal/module/identity/repository"
	identityUsecase "am-erp-go/internal/module/identity/usecase"
	integrationHttp "am-erp-go/internal/module/integration/delivery/http"
	integrationDomain "am-erp-go/internal/module/integration/domain"
	integrationProviderAmazon "am-erp-go/internal/module/integration/providers/amazon"
	integrationRepo "am-erp-go/internal/module/integration/repository"
	integrationUsecase "am-erp-go/internal/module/integration/usecase"
	inventoryHttp "am-erp-go/internal/module/inventory/delivery/http"
	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	logisticsHttp "am-erp-go/internal/module/logistics/delivery/http"
	logisticsRepo "am-erp-go/internal/module/logistics/repository"
	logisticsUsecase "am-erp-go/internal/module/logistics/usecase"
	menuHttp "am-erp-go/internal/module/menu/delivery/http"
	menuRepo "am-erp-go/internal/module/menu/repository"
	menuUsecase "am-erp-go/internal/module/menu/usecase"
	packagingHttp "am-erp-go/internal/module/packaging/delivery/http"
	packagingRepo "am-erp-go/internal/module/packaging/repository"
	packagingUsecase "am-erp-go/internal/module/packaging/usecase"
	procurementHttp "am-erp-go/internal/module/procurement/delivery/http"
	procurementRepo "am-erp-go/internal/module/procurement/repository"
	procurementUsecase "am-erp-go/internal/module/procurement/usecase"
	productHttp "am-erp-go/internal/module/product/delivery/http"
	productRepo "am-erp-go/internal/module/product/repository"
	productUsecase "am-erp-go/internal/module/product/usecase"
	salesHttp "am-erp-go/internal/module/sales/delivery/http"
	salesRepo "am-erp-go/internal/module/sales/repository"
	salesUsecase "am-erp-go/internal/module/sales/usecase"
	shipmentHttp "am-erp-go/internal/module/shipment/delivery/http"
	shipmentRepo "am-erp-go/internal/module/shipment/repository"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"
	supplierHttp "am-erp-go/internal/module/supplier/delivery/http"
	supplierRepo "am-erp-go/internal/module/supplier/repository"
	supplierUsecase "am-erp-go/internal/module/supplier/usecase"
	systemHttp "am-erp-go/internal/module/system/delivery/http"
	systemRepo "am-erp-go/internal/module/system/repository"
	systemUsecase "am-erp-go/internal/module/system/usecase"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine *gin.Engine
	Config *config.Config
}

func Build() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	database, err := db.NewMySQL(&cfg.Database)
	if err != nil {
		return nil, err
	}

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHour)

	userRepo := identityRepo.NewUserRepository(database)
	authUsecase := identityUsecase.NewAuthUsecase(userRepo, jwtManager)
	authHandler := identityHttp.NewAuthHandler(authUsecase)
	userUsecase := identityUsecase.NewUserUsecase(userRepo)
	userHandler := identityHttp.NewUserHandler(userUsecase)

	menuRepository := menuRepo.NewMenuRepository(database)
	menuSvc := menuUsecase.NewMenuUsecase(menuRepository, userRepo)
	menuHandler := menuHttp.NewMenuHandler(menuSvc)

	fieldLabelRepository := systemRepo.NewFieldLabelRepository(database)
	fieldLabelSvc := systemUsecase.NewFieldLabelUseCase(fieldLabelRepository)
	fieldLabelHandler := systemHttp.NewFieldLabelHandler(fieldLabelSvc)

	auditLogRepository := systemRepo.NewAuditLogRepository(database)
	auditLogUsecase := systemUsecase.NewAuditLogUsecase(auditLogRepository)
	auditLogHandler := systemHttp.NewAuditLogHandler(auditLogUsecase)
	configCenterRepository := systemRepo.NewConfigCenterRepository(database)
	configCenterUsecase := systemUsecase.NewConfigCenterUsecase(configCenterRepository, auditLogUsecase)
	configCenterHandler := systemHttp.NewConfigCenterHandler(configCenterUsecase)
	monitorRepository := systemRepo.NewMonitorRepository(database)
	monitorUsecase := systemUsecase.NewMonitorUsecase(monitorRepository, "migrations")
	monitorHandler := systemHttp.NewMonitorHandler(monitorUsecase)
	jobRunRepository := systemRepo.NewJobRunRepository(database)
	systemLogRepository := systemRepo.NewSystemLogRepository(database)
	jobRecorder := systemUsecase.NewJobRecorder(jobRunRepository, systemLogRepository)
	logRetentionScheduler := systemUsecase.NewLogRetentionScheduler(
		cfg.Operations.LogRetention.Enabled,
		cfg.Operations.LogRetention.CleanupIntervalMinutes,
		cfg.Operations.LogRetention.JobRunRetentionDays,
		cfg.Operations.LogRetention.SystemLogRetentionDays,
		jobRunRepository,
		systemLogRepository,
	)
	logRetentionScheduler.BindJobRecorder(jobRecorder)
	logRetentionScheduler.Start()

	productRepository := productRepo.NewProductRepository(database)
	productParentRepository := productRepo.NewProductParentRepository(database)
	productConfigRepository := productRepo.NewProductConfigRepository(database)
	productCategoryRepository := productRepo.NewProductCategoryRepository(database)
	productPackagingRepository := productRepo.NewProductPackagingRepository(database)
	productSvc := productUsecase.NewProductUsecase(
		productRepository,
		productParentRepository,
		productConfigRepository,
		productCategoryRepository,
		productPackagingRepository,
	)
	imageRepository := productRepo.NewProductImageRepository(database)
	imageUsecase := productUsecase.NewProductImageUsecase(imageRepository, productRepository)
	productHandler := productHttp.NewProductHandler(productSvc, imageUsecase)
	productHandler.BindAuditLogger(auditLogUsecase)
	comboRepository := productRepo.NewProductComboRepository(database)
	comboUsageRepository := productRepo.NewComboUsageRepository(database)
	comboUsecase := productUsecase.NewProductComboUsecase(comboRepository, productRepository, productRepository, comboUsageRepository)
	comboHandler := productHttp.NewComboHandler(comboUsecase)
	comboHandler.BindAuditLogger(auditLogUsecase)

	uploadHandler := upload.NewUploadHandler()

	supplierRepository := supplierRepo.NewSupplierRepository(database)
	supplierTypeRepository := supplierRepo.NewSupplierTypeRepository(database)
	supplierContactRepository := supplierRepo.NewSupplierContactRepository(database)
	supplierAccountRepository := supplierRepo.NewSupplierAccountRepository(database)
	supplierTagRepository := supplierRepo.NewSupplierTagRepository(database)
	supplierSvc := supplierUsecase.NewSupplierUsecase(
		supplierRepository,
		supplierTypeRepository,
		supplierContactRepository,
		supplierAccountRepository,
		supplierTagRepository,
	)
	supplierHandler := supplierHttp.NewSupplierHandler(supplierSvc)
	quoteRepository := supplierRepo.NewQuoteRepository(database)
	quoteUsecase := supplierUsecase.NewQuoteUsecase(quoteRepository, productRepository, auditLogUsecase)
	quoteHandler := supplierHttp.NewQuoteHandler(quoteUsecase)
	productSvc.BindQuoteRepository(quoteRepository)
	productSvc.BindImageRepository(imageRepository)
	productSvc.BindDefaultsProvider(configCenterUsecase)
	productSvc.BindUpsertTransactionManager(infraTx.NewProductUpsertTxManager(database))

	warehouseRepository := inventoryRepo.NewWarehouseRepository(database)
	warehouseUsecase := inventoryUsecase.NewWarehouseUsecase(warehouseRepository)
	warehouseHandler := inventoryHttp.NewWarehouseHandler(warehouseUsecase)

	balanceRepository := inventoryRepo.NewInventoryBalanceRepository(database)
	inventoryMovementRepository := inventoryRepo.NewInventoryMovementRepository(database)
	inventoryLotRepository := inventoryRepo.NewInventoryLotRepository(database)
	inventoryUsecaseObj := inventoryUsecase.NewInventoryUsecase(
		balanceRepository,
		inventoryMovementRepository,
		inventoryLotRepository,
	)
	inventoryUsecaseObj.BindPackingRequirementResolver(
		inventoryUsecase.NewProductPackingRequirementResolver(productPackagingRepository),
	)
	inventoryUsecaseObj.BindAssemblyTransactionManager(
		infraTx.NewAssemblyTxManager(database),
	)
	inventoryUsecaseObj.BindAuditLogger(auditLogUsecase)
	inventoryHandler := inventoryHttp.NewInventoryHandler(inventoryUsecaseObj)

	// 采购模块使用 InventoryUsecase 来正确创建库存流水并更新余额
	purchaseOrderRepository := procurementRepo.NewPurchaseOrderRepository(database)
	purchaseOrderUsecase := procurementUsecase.NewPurchaseOrderUsecase(
		purchaseOrderRepository,
		productRepository,
		comboRepository,
		inventoryUsecaseObj, // 使用 InventoryUsecase 而不是直接的 movementRepo
		auditLogUsecase,
	)
	purchaseOrderUsecase.BindDefaultsProvider(configCenterUsecase)
	purchaseOrderHandler := procurementHttp.NewPurchaseOrderHandler(purchaseOrderUsecase)
	replenishmentRepository := procurementRepo.NewReplenishmentRepository(database)
	replenishmentUsecase := procurementUsecase.NewReplenishmentUsecase(replenishmentRepository, purchaseOrderUsecase)
	replenishmentUsecase.BindAuditLogger(auditLogUsecase)
	replenishmentUsecase.BindDefaultsProvider(configCenterUsecase)
	purchaseOrderUsecase.BindPlanCleaner(replenishmentUsecase)
	replenishmentHandler := procurementHttp.NewReplenishmentHandler(replenishmentUsecase)
	replenishmentScheduler := procurementUsecase.NewReplenishmentScheduler(replenishmentUsecase)
	replenishmentScheduler.BindJobRecorder(jobRecorder)
	replenishmentScheduler.Start()

	// 发货模块
	shipmentRepository := shipmentRepo.NewShipmentRepo(database)
	shipmentItemRepository := shipmentRepo.NewShipmentItemRepo(database)
	shipmentUsecaseObj := shipmentUsecase.NewShipmentUsecase(
		shipmentRepository,
		shipmentItemRepository,
		inventoryUsecaseObj, // 使用 InventoryUsecase 进行库存流转
		productRepository,   // 用于加载产品信息
		warehouseRepository, // 用于加载仓库信息
	)
	shipmentUsecaseObj.BindDefaultsProvider(configCenterUsecase)
	shipmentUsecaseObj.BindAuditLogger(auditLogUsecase)
	shipmentHandler := shipmentHttp.NewShipmentHandler(shipmentUsecaseObj)

	// 装箱规格
	packageSpecRepository := shipmentRepo.NewPackageSpecRepository(database)
	packageSpecPackagingRepository := shipmentRepo.NewPackageSpecPackagingRepository(database)
	packageSpecUsecaseObj := shipmentUsecase.NewPackageSpecUseCase(packageSpecRepository, packageSpecPackagingRepository)
	packageSpecHandler := shipmentHttp.NewPackageSpecHandler(packageSpecUsecaseObj)

	// 物流模块
	logisticsProviderRepository := logisticsRepo.NewLogisticsProviderRepository(database)
	shippingRateRepository := logisticsRepo.NewShippingRateRepository(database)
	logisticsServiceRepository := logisticsRepo.NewLogisticsServiceRepository(database)

	logisticsProviderUsecase := logisticsUsecase.NewLogisticsProviderUsecase(logisticsProviderRepository)
	shippingRateUsecase := logisticsUsecase.NewShippingRateUsecase(
		shippingRateRepository,
		logisticsProviderRepository,
	)
	shipmentUsecaseObj.BindLogisticsProviderRepository(logisticsProviderRepository)
	shipmentUsecaseObj.BindShippingRateRepository(shippingRateRepository)
	shippingRateUsecase.BindDefaultsProvider(configCenterUsecase)
	logisticsServiceUsecase := logisticsUsecase.NewLogisticsServiceUsecase(logisticsServiceRepository)

	logisticsProviderHandler := logisticsHttp.NewLogisticsProviderHandler(logisticsProviderUsecase)
	shippingRateHandler := logisticsHttp.NewShippingRateHandler(shippingRateUsecase)
	logisticsServiceHandler := logisticsHttp.NewLogisticsServiceHandler(logisticsServiceUsecase)

	// 包材模块
	packagingItemRepository := packagingRepo.NewPackagingItemRepository(database)
	packagingLedgerRepository := packagingRepo.NewPackagingLedgerRepository(database)
	packagingProcurementRepository := packagingRepo.NewPackagingProcurementRepository(database)
	packagingUC := packagingUsecase.NewPackagingUsecase(packagingItemRepository, packagingLedgerRepository)
	packagingProcurementUC := packagingUsecase.NewPackagingProcurementUsecase(
		packagingProcurementRepository,
		packagingItemRepository,
		packagingLedgerRepository,
	)
	packagingProcurementUC.BindAuditLogger(auditLogUsecase)
	packagingProcurementUC.BindDefaultsProvider(configCenterUsecase)
	packagingProcurementUC.BindConvertTransactionManager(
		infraTx.NewPackagingPlanConvertTxManager(database),
	)
	packagingProcurementUC.BindReceiveTransactionManager(
		infraTx.NewPackagingPurchaseReceiveTxManager(database),
	)
	packagingHandler := packagingHttp.NewPackagingHandler(packagingUC)
	packagingProcurementHandler := packagingHttp.NewPackagingProcurementHandler(packagingProcurementUC)

	// 销售订单模块
	salesOrderRepository := salesRepo.NewSalesOrderRepository(database)
	salesOrderUsecase := salesUsecase.NewSalesOrderUsecase(salesOrderRepository, salesOrderRepository)
	salesOrderUsecase.BindInventoryService(inventoryUsecaseObj)
	salesOrderUsecase.BindConfigProvider(configCenterUsecase)
	salesOrderHandler := salesHttp.NewSalesOrderHandler(salesOrderUsecase)

	// 财务模块
	cashLedgerRepository := financeRepo.NewCashLedgerRepository(database)
	costingSnapshotRepository := financeRepo.NewCostingSnapshotRepository(database)
	costEventRepository := financeRepo.NewCostEventRepository(database)
	orderCostDetailRepository := financeRepo.NewOrderCostDetailRepository(database)
	profitLedgerRepository := financeRepo.NewProfitLedgerRepository(database)
	profitQueryRepository := financeRepo.NewProfitQueryRepository(database)
	productCostRepository := financeRepo.NewProductCostRepository(database)
	dailyProfitSnapshotRepository := financeRepo.NewDailyProfitSnapshotRepository(database)
	exchangeRateRepository := financeRepo.NewExchangeRateRepository(database)
	cashLedgerUsecase := financeUsecase.NewCashLedgerUsecase(cashLedgerRepository)
	costingSnapshotUsecase := financeUsecase.NewCostingSnapshotUsecase(costingSnapshotRepository)
	dailyProfitUsecase := financeUsecase.NewDailyProfitUsecase(
		profitLedgerRepository,
		dailyProfitSnapshotRepository,
		orderCostDetailRepository,
	)
	profitQueryUsecase := financeUsecase.NewProfitQueryUsecase(profitQueryRepository)
	productCostUsecase := financeUsecase.NewProductCostUsecase(productCostRepository)
	exchangeRateUsecase := financeUsecase.NewExchangeRateUsecase(exchangeRateRepository)
	exchangeRateUsecase.BindAuditLogger(auditLogUsecase)
	costEventWriter := financeUsecase.NewCostEventWriter(costEventRepository)
	orderCostWriter := financeUsecase.NewOrderCostWriter(orderCostDetailRepository)
	profitLedgerWriter := financeUsecase.NewProfitLedgerWriter(profitLedgerRepository)
	shipmentLandedSnapshotWriter := financeUsecase.NewShipmentLandedSnapshotWriter(costingSnapshotRepository, costEventRepository)
	financeUsecase.SetDefaultBaseCurrencyResolver(configCenterUsecase.GetDefaultBaseCurrency)
	financeUsecase.SetExchangeRateScaleResolver(configCenterUsecase.GetExchangeRateScale)
	financeUsecase.SetFXRateResolver(exchangeRateUsecase.Resolve)
	shipmentUsecaseObj.BindFXResolver(func(baseCurrency, originalCurrency string, occurredAt time.Time) (*shipmentUsecase.ShipmentFXSnapshot, error) {
		snapshot, err := exchangeRateUsecase.Resolve(baseCurrency, originalCurrency, occurredAt)
		if err != nil {
			return nil, err
		}
		return &shipmentUsecase.ShipmentFXSnapshot{
			Rate:        snapshot.Rate,
			Source:      snapshot.Source,
			Version:     snapshot.Version,
			EffectiveAt: snapshot.EffectiveAt,
		}, nil
	})
	shipmentPlatformReceiveCostResolver := shipmentUsecase.NewShipmentPlatformReceiveUnitCostResolver(
		shipmentRepository,
		shipmentItemRepository,
		configCenterUsecase,
		func(baseCurrency, originalCurrency string, occurredAt time.Time) (*shipmentUsecase.ShipmentFXSnapshot, error) {
			snapshot, err := exchangeRateUsecase.Resolve(baseCurrency, originalCurrency, occurredAt)
			if err != nil {
				return nil, err
			}
			return &shipmentUsecase.ShipmentFXSnapshot{
				Rate:        snapshot.Rate,
				Source:      snapshot.Source,
				Version:     snapshot.Version,
				EffectiveAt: snapshot.EffectiveAt,
			}, nil
		},
		costEventRepository,
	)
	inventoryUsecaseObj.BindPlatformReceiveUnitCostResolver(shipmentPlatformReceiveCostResolver)
	inventoryUsecaseObj.BindShipmentLotUnitCostResolver(shipmentPlatformReceiveCostResolver.ResolveByShipmentReference)
	inventoryUsecaseObj.BindPlatformReceiveRecorder(
		shipmentUsecase.NewShipmentPlatformReceiveRecorder(
			shipmentRepository,
			shipmentItemRepository,
			auditLogUsecase,
		),
	)
	inventoryUsecaseObj.BindPlatformReceiveTransactionManager(
		infraTx.NewPlatformReceiveTxManager(database, auditLogUsecase),
	)
	seedLotCostResolver := func(ctx context.Context, productID, warehouseID uint64) (*float64, error) {
		for _, costType := range []financeDomain.CostType{
			financeDomain.CostTypeLanded,
			financeDomain.CostTypePurchase,
			financeDomain.CostTypeAverage,
		} {
			snapshot, err := costingSnapshotUsecase.GetCurrent(productID, costType)
			if err != nil {
				return nil, err
			}
			if snapshot != nil && snapshot.UnitCost > 0 {
				value := snapshot.UnitCost
				return &value, nil
			}
		}
		return nil, nil
	}
	inventoryUsecaseObj.BindSeedLotUnitCostResolver(seedLotCostResolver)
	shipmentStateTxManager := infraTx.NewShipmentStateTxManager(database, seedLotCostResolver)
	shipmentUsecaseObj.BindConfirmTransactionManager(shipmentStateTxManager)
	shipmentUsecaseObj.BindCancelTransactionManager(shipmentStateTxManager)
	shipmentUsecaseObj.BindMarkShippedTransactionManager(
		infraTx.NewShipmentMarkShippedTxManager(database, seedLotCostResolver),
	)
	salesOrderUsecase.BindAllocateTransactionManager(
		infraTx.NewSalesAllocateTxManager(database, seedLotCostResolver),
	)
	salesOrderUsecase.BindReturnTransactionManager(
		infraTx.NewSalesReturnTxManager(database, seedLotCostResolver),
	)
	salesOrderUsecase.BindShipTransactionManager(
		infraTx.NewSalesShipTxManager(database, seedLotCostResolver),
	)
	purchaseOrderUsecase.BindCostEventRecorder(costEventWriter)
	purchaseOrderUsecase.BindSubmitTransactionManager(
		infraTx.NewPurchaseOrderSubmitTxManager(database),
	)
	purchaseOrderUsecase.BindShipTransactionManager(
		infraTx.NewPurchaseOrderShipTxManager(database),
	)
	purchaseOrderUsecase.BindReceiveTransactionManager(
		infraTx.NewPurchaseOrderReceiveTxManager(database),
	)
	purchaseOrderUsecase.BindInspectTransactionManager(
		infraTx.NewPurchaseOrderInspectTxManager(database),
	)
	shipmentUsecaseObj.BindCostAllocationRecorder(costEventWriter)
	shipmentUsecaseObj.BindLandedSnapshotRecorder(shipmentLandedSnapshotWriter)
	inventoryUsecaseObj.BindPackingCostRecorder(costEventWriter)
	salesOrderUsecase.BindShipCostDetailWriter(orderCostWriter)
	salesOrderUsecase.BindShipProfitWriter(profitLedgerWriter)
	salesOrderUsecase.BindAuditLogger(auditLogUsecase)
	cashLedgerUsecase.BindProfitLedgerRecorder(profitLedgerWriter)
	financeHandler := financeHttp.NewFinanceHandler(
		cashLedgerUsecase,
		costingSnapshotUsecase,
		dailyProfitUsecase,
		profitQueryUsecase,
		productCostUsecase,
		exchangeRateUsecase,
	)
	financeHandler.BindAuditLogger(auditLogUsecase)

	orderSyncRepo := integrationRepo.NewOrderSyncRepository(database)
	refundSyncRepo := integrationRepo.NewRefundSyncRepository(database)
	authorizationRepo := integrationRepo.NewAuthorizationRepository(database)
	skuMappingRepo := integrationRepo.NewSKUMappingRepository(database)
	skuMappingUsecase := integrationUsecase.NewSKUMappingUsecase(skuMappingRepo, nil)
	orderSyncRegistry := integrationUsecase.NewOrderSyncRegistry()
	refundSyncRegistry := integrationUsecase.NewRefundSyncRegistry()
	authProviders := make([]integrationDomain.AuthorizationProvider, 0)
	for _, providerCfg := range cfg.Integrations.Providers {
		if !providerCfg.Enabled {
			continue
		}
		switch strings.ToLower(providerCfg.Type) {
		case "amazon":
			if providerCfg.Amazon == nil {
				log.Printf("[integration] skip provider %s: missing amazon config block", providerCfg.Code)
				continue
			}

			amazonSPClient := integrationProviderAmazon.NewSPAPIClient(integrationProviderAmazon.SPAPIConfig{
				Endpoint:          providerCfg.Amazon.Endpoint,
				AppID:             providerCfg.Amazon.AppID,
				AuthorizeBaseURL:  providerCfg.Amazon.AuthorizeBaseURL,
				RedirectURI:       providerCfg.Amazon.RedirectURI,
				MarketplaceIDs:    providerCfg.Amazon.MarketplaceIDs,
				LWAClientID:       providerCfg.Amazon.LWAClientID,
				LWAClientSecret:   providerCfg.Amazon.LWAClientSecret,
				LWARefreshToken:   providerCfg.Amazon.LWARefreshToken,
				RequestTimeoutSec: providerCfg.RequestTimeoutSecond,
			})
			amazonProvider := integrationProviderAmazon.NewOrdersProvider(providerCfg.Code, amazonSPClient)
			authProviders = append(authProviders, amazonProvider)
			orderSyncService := integrationUsecase.NewOrderSyncService(
				orderSyncRepo,
				salesOrderRepository,
				amazonProvider,
				integrationUsecase.OrderSyncConfig{
					ProviderCode:       providerCfg.Code,
					Channel:            providerCfg.Channel,
					SourceType:         providerCfg.SourceType,
					SalesChannel:       providerCfg.SalesChannel,
					DefaultCurrency:    providerCfg.DefaultCurrency,
					MarketplaceIDs:     providerCfg.Amazon.MarketplaceIDs,
					LookbackMinutes:    providerCfg.LookbackMinutes,
					InitialLookbackDay: providerCfg.InitialLookbackDays,
				},
				nil,
				configCenterUsecase,
			)
			orderSyncService.BindSKUMappingResolver(skuMappingUsecase)

			orderSyncRegistry.Register(providerCfg.Code, orderSyncService)
			orderSyncRegistry.Register(strings.ToLower(providerCfg.Code), orderSyncService)

			refundSyncService := integrationUsecase.NewRefundSyncService(
				refundSyncRepo,
				amazonProvider,
				integrationUsecase.RefundSyncConfig{
					ProviderCode:       providerCfg.Code,
					Channel:            providerCfg.Channel,
					MarketplaceIDs:     providerCfg.Amazon.MarketplaceIDs,
					LookbackMinutes:    providerCfg.LookbackMinutes,
					InitialLookbackDay: providerCfg.InitialLookbackDays,
				},
				nil,
				skuMappingUsecase,
				salesOrderRepository,
			)
			refundSyncRegistry.Register(providerCfg.Code, refundSyncService)
			refundSyncRegistry.Register(strings.ToLower(providerCfg.Code), refundSyncService)

			if providerCfg.AutoSyncEnabled {
				orderSyncScheduler := integrationUsecase.NewOrderSyncScheduler(true, providerCfg.SyncIntervalMinutes, orderSyncService)
				orderSyncScheduler.BindJobRecorder(jobRecorder)
				orderSyncScheduler.Start()
			}
		default:
			log.Printf("[integration] unsupported provider type %s (code=%s)", providerCfg.Type, providerCfg.Code)
		}
	}

	authorizationUsecase := integrationUsecase.NewAuthorizationUsecase(
		authorizationRepo,
		authProviders,
		nil,
	)
	authorizationScheduler := integrationUsecase.NewAuthorizationRefreshScheduler(true, 5, authorizationUsecase)
	authorizationScheduler.BindJobRecorder(jobRecorder)
	authorizationScheduler.Start()

	orderSyncHandler := integrationHttp.NewOrderSyncHandler(orderSyncRegistry)
	refundSyncHandler := integrationHttp.NewRefundSyncHandler(refundSyncRegistry)
	authorizationHandler := integrationHttp.NewAuthorizationHandler(authorizationUsecase)
	authorizationHandler.BindAuditLogger(auditLogUsecase)
	skuMappingHandler := integrationHttp.NewSKUMappingHandler(skuMappingUsecase)
	skuMappingHandler.BindAuditLogger(auditLogUsecase)

	r := router.NewRouter(
		jwtManager,
		userRepo,
		authHandler,
		userHandler,
		menuHandler,
		productHandler,
		comboHandler,
		supplierHandler,
		quoteHandler,
		purchaseOrderHandler,
		replenishmentHandler,
		configCenterHandler,
		fieldLabelHandler,
		auditLogHandler,
		monitorHandler,
		uploadHandler,
		warehouseHandler,
		inventoryHandler,
		shipmentHandler,
		packageSpecHandler,
		logisticsProviderHandler,
		shippingRateHandler,
		logisticsServiceHandler,
		packagingHandler,
		packagingProcurementHandler,
		salesOrderHandler,
		financeHandler,
		orderSyncHandler,
		refundSyncHandler,
		authorizationHandler,
		skuMappingHandler,
	)
	engine := r.Setup()

	return &App{
		Engine: engine,
		Config: cfg,
	}, nil
}
