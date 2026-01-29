package bootstrap

import (
	"am-erp-go/internal/infrastructure/auth"
	"am-erp-go/internal/infrastructure/config"
	"am-erp-go/internal/infrastructure/db"
	"am-erp-go/internal/infrastructure/router"
	"am-erp-go/internal/infrastructure/upload"
	identityHttp "am-erp-go/internal/module/identity/delivery/http"
	identityRepo "am-erp-go/internal/module/identity/repository"
	identityUsecase "am-erp-go/internal/module/identity/usecase"
	menuHttp "am-erp-go/internal/module/menu/delivery/http"
	menuRepo "am-erp-go/internal/module/menu/repository"
	menuUsecase "am-erp-go/internal/module/menu/usecase"
	procurementHttp "am-erp-go/internal/module/procurement/delivery/http"
	procurementRepo "am-erp-go/internal/module/procurement/repository"
	procurementUsecase "am-erp-go/internal/module/procurement/usecase"
	productHttp "am-erp-go/internal/module/product/delivery/http"
	productRepo "am-erp-go/internal/module/product/repository"
	productUsecase "am-erp-go/internal/module/product/usecase"
	supplierHttp "am-erp-go/internal/module/supplier/delivery/http"
	supplierRepo "am-erp-go/internal/module/supplier/repository"
	supplierUsecase "am-erp-go/internal/module/supplier/usecase"
	systemHttp "am-erp-go/internal/module/system/delivery/http"
	systemRepo "am-erp-go/internal/module/system/repository"
	systemUsecase "am-erp-go/internal/module/system/usecase"
	inventoryHttp "am-erp-go/internal/module/inventory/delivery/http"
	inventoryRepo "am-erp-go/internal/module/inventory/repository"
	inventoryUsecase "am-erp-go/internal/module/inventory/usecase"
	shipmentHttp "am-erp-go/internal/module/shipment/delivery/http"
	shipmentRepo "am-erp-go/internal/module/shipment/repository"
	shipmentUsecase "am-erp-go/internal/module/shipment/usecase"
	logisticsHttp "am-erp-go/internal/module/logistics/delivery/http"
	logisticsRepo "am-erp-go/internal/module/logistics/repository"
	logisticsUsecase "am-erp-go/internal/module/logistics/usecase"
	packagingHttp "am-erp-go/internal/module/packaging/delivery/http"
	packagingRepo "am-erp-go/internal/module/packaging/repository"
	packagingUsecase "am-erp-go/internal/module/packaging/usecase"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine *gin.Engine
	Config *config.Config
}

func Build() (*App, error) {
	cfg := config.Load()

	database, err := db.NewMySQL(&cfg.Database)
	if err != nil {
		return nil, err
	}

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHour)

	userRepo := identityRepo.NewUserRepository(database)
	authUsecase := identityUsecase.NewAuthUsecase(userRepo, jwtManager)
	authHandler := identityHttp.NewAuthHandler(authUsecase)

	menuRepository := menuRepo.NewMenuRepository(database)
	menuSvc := menuUsecase.NewMenuUsecase(menuRepository, userRepo)
	menuHandler := menuHttp.NewMenuHandler(menuSvc)

	fieldLabelRepository := systemRepo.NewFieldLabelRepository(database)
	fieldLabelSvc := systemUsecase.NewFieldLabelUseCase(fieldLabelRepository)
	fieldLabelHandler := systemHttp.NewFieldLabelHandler(fieldLabelSvc)

	auditLogRepository := systemRepo.NewAuditLogRepository(database)
	auditLogUsecase := systemUsecase.NewAuditLogUsecase(auditLogRepository)
	auditLogHandler := systemHttp.NewAuditLogHandler(auditLogUsecase)

	productRepository := productRepo.NewProductRepository(database)
	productParentRepository := productRepo.NewProductParentRepository(database)
	productPackagingRepository := productRepo.NewProductPackagingRepository(database)
	productSvc := productUsecase.NewProductUsecase(
		productRepository,
		productParentRepository,
		productPackagingRepository,
	)
	imageRepository := productRepo.NewProductImageRepository(database)
	imageUsecase := productUsecase.NewProductImageUsecase(imageRepository, productRepository)
	productHandler := productHttp.NewProductHandler(productSvc, imageUsecase)
	comboRepository := productRepo.NewProductComboRepository(database)
	comboUsecase := productUsecase.NewProductComboUsecase(comboRepository, productRepository, productRepository)
	comboHandler := productHttp.NewComboHandler(comboUsecase)

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

	warehouseRepository := inventoryRepo.NewWarehouseRepository(database)
	warehouseUsecase := inventoryUsecase.NewWarehouseUsecase(warehouseRepository)
	warehouseHandler := inventoryHttp.NewWarehouseHandler(warehouseUsecase)

	balanceRepository := inventoryRepo.NewInventoryBalanceRepository(database)
	inventoryMovementRepository := inventoryRepo.NewInventoryMovementRepository(database)
	inventoryUsecaseObj := inventoryUsecase.NewInventoryUsecase(balanceRepository, inventoryMovementRepository)
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
	purchaseOrderHandler := procurementHttp.NewPurchaseOrderHandler(purchaseOrderUsecase)

	// 发货模块
	shipmentRepository := shipmentRepo.NewShipmentRepo(database)
	shipmentItemRepository := shipmentRepo.NewShipmentItemRepo(database)
	shipmentUsecaseObj := shipmentUsecase.NewShipmentUsecase(
		shipmentRepository,
		shipmentItemRepository,
		inventoryUsecaseObj,   // 使用 InventoryUsecase 进行库存流转
		productRepository,     // 用于加载产品信息
		warehouseRepository,   // 用于加载仓库信息
	)
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
	logisticsServiceUsecase := logisticsUsecase.NewLogisticsServiceUsecase(logisticsServiceRepository)

	logisticsProviderHandler := logisticsHttp.NewLogisticsProviderHandler(logisticsProviderUsecase)
	shippingRateHandler := logisticsHttp.NewShippingRateHandler(shippingRateUsecase)
	logisticsServiceHandler := logisticsHttp.NewLogisticsServiceHandler(logisticsServiceUsecase)

	// 包材模块
	packagingItemRepository := packagingRepo.NewPackagingItemRepository(database)
	packagingLedgerRepository := packagingRepo.NewPackagingLedgerRepository(database)
	packagingUsecase := packagingUsecase.NewPackagingUsecase(packagingItemRepository, packagingLedgerRepository)
	packagingHandler := packagingHttp.NewPackagingHandler(packagingUsecase)

	r := router.NewRouter(
		jwtManager,
		authHandler,
		menuHandler,
		productHandler,
		comboHandler,
		supplierHandler,
		quoteHandler,
		purchaseOrderHandler,
		fieldLabelHandler,
		auditLogHandler,
		uploadHandler,
		warehouseHandler,
		inventoryHandler,
		shipmentHandler,
		packageSpecHandler,
		logisticsProviderHandler,
		shippingRateHandler,
		logisticsServiceHandler,
		packagingHandler,
	)
	engine := r.Setup()

	return &App{
		Engine: engine,
		Config: cfg,
	}, nil
}
